package news

import (
	"os"
	"path/filepath"
	"testing"
)

const decisions = `# Decisions

## 1. The first one

Its opening paragraph, which is
wrapped across lines.

A second paragraph nobody needs in a summary.

## 2. The second one

Only one line here.

### A subheading that is not a decision

More prose.

## 10. Ten sorts after two, not before it

Numbered, not lexical.
`

func write(t *testing.T, name, body string) string {
	t.Helper()
	path := filepath.Join(t.TempDir(), name)
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return path
}

// The lede is the first paragraph and stops there. A decision runs to several
// hundred words, and a listing that prints all of them is one nobody reads.
func TestReadTakesTheFirstParagraphOnly(t *testing.T) {
	ds, err := Read(write(t, "decisions.md", decisions))
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) != 3 {
		t.Fatalf("read %d decisions, want 3: %+v", len(ds), ds)
	}
	if ds[0].Title != "The first one" {
		t.Errorf("title is %q", ds[0].Title)
	}
	if want := "Its opening paragraph, which is wrapped across lines."; ds[0].Lede != want {
		t.Errorf("lede is %q, want %q", ds[0].Lede, want)
	}
}

// A `###` inside a decision is part of it, not a new one. decisions.md uses
// them for the sections inside a long entry.
func TestASubheadingIsNotADecision(t *testing.T) {
	ds, _ := Read(write(t, "decisions.md", decisions))
	for _, d := range ds {
		if d.Title == "A subheading that is not a decision" {
			t.Error("a ### was read as a decision")
		}
	}
}

// Numerically, so 10 is after 2. Sorting these as strings puts decision 10
// between 1 and 2 and silently hides eight of them from every listing.
func TestSinceComparesNumbersNotStrings(t *testing.T) {
	ds, _ := Read(write(t, "decisions.md", decisions))

	fresh := Since(ds, 2)
	if len(fresh) != 1 || fresh[0].Number != 10 {
		t.Errorf("since 2 gave %+v", fresh)
	}
	if n := Latest(ds); n != 10 {
		t.Errorf("latest is %d, want 10", n)
	}
	if got := Since(ds, 10); len(got) != 0 {
		t.Errorf("a tool that has read everything was told about %+v", got)
	}
}

// A file with no numbered decisions is an error, not an empty list. Silently
// reporting "up to date" from a file that was never parsed is the one answer
// that must not be reachable by accident.
func TestAFileWithNoDecisionsIsAnError(t *testing.T) {
	if _, err := Read(write(t, "decisions.md", "# Notes\n\nNothing numbered.\n")); err == nil {
		t.Error("a file with no decisions parsed cleanly")
	}
	if _, err := Read(filepath.Join(t.TempDir(), "absent.md")); err == nil {
		t.Error("a missing file parsed cleanly")
	}
}

// The marker is written by hand in prose, so it is read loosely.
func TestTheMarkerIsReadOutOfProse(t *testing.T) {
	for _, body := range []string{
		"Reconciled with tuikit through decision 32.",
		"reconciled with tuikit through decision 32",
		"# mytool\n\nSome prose.\n\nReconciled with tuikit through decision 32. More prose.\n",
		"Reconciled  with   tuikit\nthrough decision 32.",
	} {
		got, ok := Marker(write(t, "AGENTS.md", body))
		if !ok || got != 32 {
			t.Errorf("%q read as (%d, %v)", body, got, ok)
		}
	}

	if _, ok := Marker(write(t, "AGENTS.md", "# mytool\n\nNo marker here.\n")); ok {
		t.Error("a file with no marker reported one")
	}
}

// The checkout comes from the go.mod, because that is the path the compiler
// used — a second way to say where tuikit is, is a second way to be told about
// a tuikit the tool does not build against.
func TestTheCheckoutComesFromTheReplaceDirective(t *testing.T) {
	dir := t.TempDir()
	gomod := filepath.Join(dir, "go.mod")
	body := "module example.com/mytool\n\ngo 1.25\n\nreplace github.com/richarddavenport/tuikit => ../tuikit\n"
	if err := os.WriteFile(gomod, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}

	got, ok := Checkout(gomod)
	if !ok {
		t.Fatal("no checkout found")
	}
	// Relative to the go.mod, which is what Go does — not to whatever directory
	// the command happened to be run from.
	want, _ := filepath.Abs(filepath.Join(dir, "..", "tuikit"))
	if got != want {
		t.Errorf("resolved to %s, want %s", got, want)
	}
}

func TestAGoModWithNoReplaceIsNotAToolWeCanHelp(t *testing.T) {
	gomod := write(t, "go.mod", "module example.com/thing\n\ngo 1.25\n")
	if _, ok := Checkout(gomod); ok {
		t.Error("a go.mod with no tuikit replace reported a checkout")
	}
}

// The real file, so the parser is held against the thing it exists to read
// rather than only against a fixture that agrees with it.
func TestTheRealDecisionsFileParses(t *testing.T) {
	ds, err := Read(filepath.Join("..", "design", "decisions.md"))
	if err != nil {
		t.Fatal(err)
	}
	if len(ds) < 30 {
		t.Fatalf("only %d decisions parsed out of the real file", len(ds))
	}
	// Numbered from 1 with no gaps: a gap means a heading the parser missed.
	for i, d := range ds {
		if d.Number != i+1 {
			t.Fatalf("decision %d in order is numbered %d (%q) — a heading was missed or misread", i+1, d.Number, d.Title)
		}
		if d.Title == "" || d.Lede == "" {
			t.Errorf("decision %d has no %s", d.Number, map[bool]string{true: "title", false: "lede"}[d.Title == ""])
		}
	}
}
