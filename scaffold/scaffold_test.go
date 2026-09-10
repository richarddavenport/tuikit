package scaffold

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/news"
)

// The acceptance criterion for the scaffolder, run rather than asserted about:
// a generated tool builds, runs, and passes its own checks with no edits.
//
// It is slow, because it is a compile and a test run of another module. It is
// also the only check that means anything — a scaffolder is a claim about code
// nobody has looked at yet, and every other test here would pass on a template
// that does not compile.
func TestAGeneratedToolPassesItsOwnChecks(t *testing.T) {
	if testing.Short() {
		t.Skip("compiles and tests another module")
	}
	root := generate(t, Tool{Name: "widgetctl", Short: "watches some widgets"})

	// Bootstrap is part of generating, so this is what `tuikit new` does and
	// not a step the test invented to make itself pass.
	if err := Bootstrap(root); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}

	for _, args := range [][]string{
		{"build", "./..."},
		{"vet", "./..."},
		{"test", "./..."},
	} {
		if out, err := run(root, "go", args...); err != nil {
			t.Errorf("go %s failed:\n%s", strings.Join(args, " "), out)
		}
	}

	// It RUNS. A tool whose tests pass and whose binary cannot answer
	// `describe` has an interface an agent cannot read.
	out, err := run(root, "go", "run", "./cmd/widgetctl", "describe", "--json")
	if err != nil {
		t.Fatalf("describe: %v\n%s", err, out)
	}
	for _, want := range []string{`"widgetctl list"`, `"items.row"`, `"Accent"`} {
		if !strings.Contains(out, want) {
			t.Errorf("describe --json does not mention %s", want)
		}
	}
}

// Every generated file is gofmt-clean, checked without a toolchain so it fails
// fast and says which file.
func TestEveryGeneratedFileIsFormatted(t *testing.T) {
	root := generate(t, Tool{Name: "widgetctl"})
	out, err := run(root, "gofmt", "-l", ".")
	if err != nil {
		t.Fatalf("gofmt: %v\n%s", err, out)
	}
	if strings.TrimSpace(out) != "" {
		t.Errorf("not gofmt-clean:\n%s", out)
	}
}

// The name is substituted into paths as well as contents. A tool called
// widgetctl whose binary lives in cmd/NAME is a tool that does not build, and
// the failure would be a puzzling one.
func TestTheNameReachesThePaths(t *testing.T) {
	root := generate(t, Tool{Name: "widgetctl"})
	if _, err := os.Stat(filepath.Join(root, "cmd", "widgetctl", "main.go")); err != nil {
		t.Errorf("cmd/widgetctl/main.go was not written: %v", err)
	}
	if _, err := os.Stat(filepath.Join(root, "cmd", "NAME")); err == nil {
		t.Error("cmd/NAME survived — the placeholder is still in a path")
	}
}

// A directory with anything in it is left alone. "It emptied my repo" is not a
// thing to find out by doing it.
func TestWritingIntoSomethingElseRefuses(t *testing.T) {
	dir := t.TempDir()
	if err := os.MkdirAll(filepath.Join(dir, "widgetctl"), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(dir, "widgetctl", "README.md"), []byte("mine"), 0o600); err != nil {
		t.Fatal(err)
	}
	if _, err := (Tool{Name: "widgetctl"}).Write(dir); err == nil {
		t.Fatal("wrote into a directory that already had something in it")
	}
	body, err := os.ReadFile(filepath.Join(dir, "widgetctl", "README.md"))
	if err != nil || string(body) != "mine" {
		t.Errorf("the existing file did not survive: %q, %v", body, err)
	}
}

// A name that is not a legal directory, binary and package path is refused
// here, where the message can say so, rather than as a compile error in
// generated code somebody did not write.
func TestABadNameIsRefused(t *testing.T) {
	for _, name := range []string{"", "Widget", "widget ctl", "widget/ctl", "widget.ctl"} {
		if _, err := (Tool{Name: name}).Write(t.TempDir()); err == nil {
			t.Errorf("%q was accepted", name)
		}
	}
}

// generate writes a tool into a temporary directory and returns its root.
//
// Against THIS checkout, not against the released version the default would
// pull: the acceptance test's whole job is to fail when a change here breaks
// generated code, and a generated tool compiled against the last tag cannot
// see the change being tested. That makes it the one caller that wants
// -tuikit, and the reason the flag still exists.
func generate(t *testing.T, tool Tool) string {
	t.Helper()
	tool.Tuikit = tuikitRoot(t)

	dir := t.TempDir()
	if _, err := tool.Write(dir); err != nil {
		t.Fatalf("write: %v", err)
	}
	return filepath.Join(dir, tool.Name)
}

func run(dir, name string, args ...string) (string, error) {
	cmd := exec.Command(name, args...)
	cmd.Dir = dir
	out, err := cmd.CombinedOutput()
	return string(out), err
}

// Somebody who has installed the command and has no checkout gets a tool. That
// is the whole point of a public module, and for a while it was the one thing
// the scaffolder could not do: the default `replace => ../tuikit` meant every
// first run anywhere refused, naming a flag whose only valid argument was a
// clone the person had not made.
func TestAToolIsGeneratedWithNoCheckoutAnywhere(t *testing.T) {
	dir := t.TempDir()
	if _, err := (Tool{Name: "widgetctl"}).Write(dir); err != nil {
		t.Fatalf("a plain run refused: %v", err)
	}

	gomod, err := os.ReadFile(filepath.Join(dir, "widgetctl", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(string(gomod), "replace") {
		t.Errorf("a tool generated without -tuikit still has a replace:\n%s", gomod)
	}
	if want := "require " + tuikitModule + " " + tuikitVersion; !strings.Contains(string(gomod), want) {
		t.Errorf("go.mod does not %q:\n%s", want, gomod)
	}
}

// -tuikit is still checked before anything is written, so a wrong path says so
// while it is still on screen rather than failing four steps later inside
// `go mod tidy`, against a module the reader did not write.
func TestABadTuikitPathIsRefusedBeforeWriting(t *testing.T) {
	dir := t.TempDir()
	_, err := (Tool{Name: "widgetctl", Tuikit: filepath.Join(dir, "nowhere")}).Write(dir)
	if err == nil {
		t.Fatal("accepted a -tuikit that is not a checkout")
	}
	if !strings.Contains(err.Error(), "-tuikit") {
		t.Errorf("the error does not name the flag it is about: %v", err)
	}
	if entries, _ := os.ReadDir(filepath.Join(dir, "widgetctl")); len(entries) > 0 {
		t.Errorf("%d files were written anyway", len(entries))
	}
}

// A stale tuikitVersion generates tools pinned to a tuikit that is not the
// current release.
func TestTheCompiledInVersionIsTheNewestTag(t *testing.T) {
	out, err := run(tuikitRoot(t), "git", "describe", "--tags", "--abbrev=0")
	if err != nil {
		t.Skipf("no git tag to check against: %v", err)
	}
	if want := strings.TrimSpace(out); tuikitVersion != want {
		t.Errorf("tuikitVersion is %s but the newest tag is %s — bump it", tuikitVersion, want)
	}
}

// A generated tool is born reconciled. One generated with a marker of zero
// would, on its first `tuikit news`, be told about every decision tuikit has
// ever taken — including the ones that produced the code it was just given.
func TestAGeneratedToolRecordsTheCurrentDecision(t *testing.T) {
	// Against a checkout, where "current" means this working tree and the
	// marker is read from its decisions.md. The default path stamps the
	// constant instead, and is checked against the decisions its required
	// version ships — the two answers differ between tags, which is why they
	// are two tests and not one.
	dir := t.TempDir()
	if _, err := (Tool{Name: "probectl", Tuikit: tuikitRoot(t)}).Write(dir); err != nil {
		t.Fatal(err)
	}

	body, err := os.ReadFile(filepath.Join(dir, "probectl", "AGENTS.md"))
	if err != nil {
		t.Fatal(err)
	}
	got, ok := news.Marker(filepath.Join(dir, "probectl", "AGENTS.md"))
	if !ok {
		t.Fatalf("the generated AGENTS.md has no marker:\n%s", body)
	}

	ds, err := news.Read(filepath.Join(tuikitRoot(t), "design", "decisions.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := news.Latest(ds); got != want {
		t.Errorf("generated at decision %d, but tuikit is at %d", got, want)
	}
}

// tuikitRoot is this checkout, as an absolute path — the scaffolder resolves a
// relative one against the tool it is writing, not against the test.
func tuikitRoot(t *testing.T) string {
	t.Helper()
	abs, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	return abs
}

// The replace path is relative, which is what go.mod.tmpl's own comment
// promises. Issue 47: `-tuikit /abs/path` was written verbatim under that
// comment, so the generated repository built on exactly one machine.
func TestTheReplacePathIsRelative(t *testing.T) {
	base := t.TempDir()
	// tuikit as a sibling of where the tool will be written, which is the
	// arrangement the whole scheme assumes.
	side := filepath.Join(base, "tuikit")
	if err := os.MkdirAll(side, 0o755); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(filepath.Join(tuikitRoot(t), "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(filepath.Join(side, "go.mod"), body, 0o600); err != nil {
		t.Fatal(err)
	}

	if _, err := (Tool{Name: "probectl", Tuikit: side}).Write(base); err != nil {
		t.Fatal(err)
	}

	gomod, err := os.ReadFile(filepath.Join(base, "probectl", "go.mod"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(gomod), "=> ../tuikit") {
		t.Errorf("an absolute -tuikit was not made relative:\n%s", gomod)
	}
	if strings.Contains(string(gomod), base) {
		t.Errorf("the generating machine's path is in the generated go.mod:\n%s", gomod)
	}
}

// A local-path replace means tuikit is not fetchable as a module, which means
// CI needs a second checkout to put it where the replace points. Two things
// that must agree, and nothing checked they did: the generated workflow had one
// checkout, so the one generated file that cannot work was the one nothing ran.
//
// Checkable without a runner, which is the point — the invariant is between two
// generated files, not between a file and GitHub.
func TestTheWorkflowChecksOutWhateverTheReplaceNeeds(t *testing.T) {
	// Both shapes, because the invariant is not "two checkouts" — it is that
	// the workflow matches the go.mod beside it. Only checking the -tuikit one
	// would let the ordinary tool ship a workflow that clones a repository it
	// has no reason to want and no token to reach.
	for _, tool := range []Tool{
		{Name: "probectl"},
		{Name: "probectl", Tuikit: tuikitRoot(t)},
	} {
		dir := t.TempDir()
		if _, err := tool.Write(dir); err != nil {
			t.Fatal(err)
		}

		gomod, err := os.ReadFile(filepath.Join(dir, "probectl", "go.mod"))
		if err != nil {
			t.Fatal(err)
		}
		ci, err := os.ReadFile(filepath.Join(dir, "probectl", ".github", "workflows", "ci.yml"))
		if err != nil {
			t.Fatal(err)
		}

		local := strings.Contains(string(gomod), "=> ..") || strings.Contains(string(gomod), "=> /")
		checkouts := strings.Count(string(ci), "actions/checkout@")

		switch {
		case local && checkouts < 2:
			t.Errorf("go.mod replaces tuikit with a local path but the workflow checks out %d repositories — "+
				"CI cannot build, and `go build` fails before gofmt, vet or the tests run", checkouts)
		case !local && checkouts > 1:
			t.Errorf("the replace is gone but the workflow still checks out %d repositories", checkouts)
		}

		if local != strings.Contains(string(ci), "TUIKIT_TOKEN") {
			t.Errorf("replace=%v but TUIKIT_TOKEN present=%v; the secret exists only to reach the second checkout",
				local, !local)
		}
	}
}

// The acceptance test above compiles a tool built against this checkout, which
// is the only way it can catch a change made here. Nothing in it exercises the
// go.mod everyone else gets, so a require pointing at a version that does not
// resolve would pass every check and fail on the first stranger.
func TestAToolBuiltAgainstTheReleasedVersionResolves(t *testing.T) {
	if testing.Short() {
		t.Skip("fetches tuikit from the module proxy")
	}
	dir := t.TempDir()
	if _, err := (Tool{Name: "widgetctl"}).Write(dir); err != nil {
		t.Fatal(err)
	}
	root := filepath.Join(dir, "widgetctl")
	if err := Bootstrap(root); err != nil {
		t.Fatalf("bootstrap: %v", err)
	}
	if out, err := run(root, "go", "build", "./..."); err != nil {
		t.Fatalf("a tool depending on tuikit %s does not build:\n%s", tuikitVersion, out)
	}

	// And it is born reconciled with THAT tuikit, not with this working tree.
	// The two differ for every commit between a tag and the next one, and a
	// marker taken from the tree would name decisions the tool's dependency
	// does not contain — a tool told it has read something that does not exist
	// where it is looking.
	dir, err := run(root, "go", "list", "-m", "-f", "{{.Dir}}", tuikitModule)
	if err != nil {
		t.Fatalf("locating the module: %v\n%s", err, dir)
	}
	ds, err := news.Read(filepath.Join(strings.TrimSpace(dir), "design", "decisions.md"))
	if err != nil {
		t.Fatal(err)
	}
	if want := news.Latest(ds); tuikitDecision != want {
		t.Errorf("tuikitDecision is %d but tuikit %s ships %d — bump it with the version",
			tuikitDecision, tuikitVersion, want)
	}
}
