package scaffold

import (
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"testing"
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
func generate(t *testing.T, tool Tool) string {
	t.Helper()
	// The generated go.mod must point at THIS checkout, not at wherever a
	// sibling ../tuikit happens to be — otherwise the test compiles somebody
	// else's copy of the framework.
	here, err := filepath.Abs("..")
	if err != nil {
		t.Fatal(err)
	}
	tool.Tuikit = here

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
