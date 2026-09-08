package guard

import (
	"go/ast"
	"strconv"
	"strings"
)

// Denied is one import prefix an engine may not reach for, and why.
//
// The reasons are not decoration. A guard that reports "forbidden import" and
// stops has told you what happened; one that says why tells you whether the
// rule is wrong, which is occasionally the right conclusion.
type Denied struct {
	Prefix string
	Why    string
}

// TerminalPackages is what an engine may not import.
//
// Exported and documented rather than buried in the function, because a tool
// being held to a rule should be able to read it.
//
// The list is by CATEGORY, not by library. Decision 22 sets the bar at "no
// terminal concepts", not "no Bubble Tea" — an engine that measures display
// width has learned about columns, whichever package it used to do it.
var TerminalPackages = []Denied{
	{"github.com/charmbracelet", "a terminal UI library — the engine does not draw"},
	{"github.com/gdamore/tcell", "a terminal UI library — the engine does not draw"},
	{"github.com/rivo/tview", "a terminal UI library — the engine does not draw"},
	{"github.com/nsf/termbox-go", "a terminal UI library — the engine does not draw"},
	{"github.com/muesli/termenv", "terminal capabilities — the engine has no terminal"},
	{"github.com/muesli/reflow", "wrapping and width — the engine has no columns"},
	{"github.com/muesli/ansi", "escape sequences — the engine has no terminal"},
	{"github.com/mattn/go-isatty", "asks whether output is a terminal, which the engine must not care about"},
	{"github.com/mattn/go-runewidth", "display width — the engine has no columns"},
	{"github.com/rivo/uniseg", "grapheme clustering, used for display width — the engine has no columns"},
	{"golang.org/x/term", "terminal control — the engine has no terminal"},
	{"github.com/richarddavenport/tuikit", "tuikit itself — the engine gets nothing from the framework"},
}

// Engine reports an import that means the engine has learned about terminals.
//
// The split every tool in the family keeps is that the engine knows the domain
// and the UI knows the terminal. That is worth as much as it is enforced, and
// it is mechanically checkable, so it is checked.
//
// The bar is higher than "no Bubble Tea imports". Measured on the tools that
// exist when this was written: democtl's fleet imports fmt, math/rand and
// time, and the deploy tool's internal/engine imports nothing from
// charmbracelet across some forty files. What an engine imports is stdlib plus
// its own domain SDK — pgx, the Azure SDK, the Docker SDK.
//
// Pass denied to hold a package to a different list; the default is
// TerminalPackages.
func Engine(t T, dir string, denied ...Denied) {
	t.Helper()

	if len(denied) == 0 {
		denied = TerminalPackages
	}
	for _, f := range sources(t, dir) {
		for _, imp := range f.ast.Imports {
			path, err := strconv.Unquote(imp.Path.Value)
			if err != nil {
				continue
			}
			for _, d := range denied {
				if !underPrefix(path, d.Prefix) {
					continue
				}
				t.Errorf("%s:%d imports %s — %s. "+
					"The engine knows the domain; the terminal is the UI's business.",
					f.name, f.fset.Position(imp.Pos()).Line, path, d.Why)
			}
		}
	}
}

// underPrefix matches a path against an import prefix on a path boundary, so
// "github.com/muesli/ansi" does not also condemn a hypothetical
// "github.com/muesli/ansiblе-something".
func underPrefix(path, prefix string) bool {
	return path == prefix || strings.HasPrefix(path, prefix+"/")
}

// Derived fails when a test builds one of the engine's own answers by hand.
//
// # The failure it catches
//
// Decision 1 says the fixture the screens are rendered from is a value the
// engine returns. Nothing enforced it, and in the board it quietly stopped
// being one: `fixture()` set a field by hand while the engine composed the
// same field as a sentence, and the renderer composed it again. The golden
// showed the row disagreeing with itself on every run and looked right,
// because it was checking the renderer against a world invented three hundred
// lines away in the same file.
//
// Four bugs shipped past 72 goldens, a color check and a narrow-terminal run.
// Three of them for this reason.
//
// The compounding part is the worst of it: a hand-made fixture makes a wrong
// screen look correct AND hands you an easy way to keep it that way. Editing
// the fixture to match the renderer "fixes" the test and hides the bug.
//
// # What to pass
//
// The types that are the engine's ANSWERS rather than its inputs — what a
// gather-and-reconcile produced, not the arguments it was called with. In the
// board those are `Truth` and `Live`: neither is ever written by a person, so
// neither should ever be typed by one in a test.
//
//	guard.Derived(t, ".", "engine", "Truth", "Live")
//
// A type that a test SHOULD construct — a request, an option, a filter — must
// not be listed. This is not "tests may not build structs"; it is "these
// particular structs are conclusions, and a conclusion typed by hand is a
// conclusion nobody checked".
//
// # The seam that makes obeying it possible
//
// An engine that only exposes "run the commands and judge the result" as one
// step leaves a test no way to get a real value without real commands. Split
// gathering from reconciling and the fixture becomes the reconciler applied to
// canned input, which is a value the engine returned and a world that exists.
func Derived(t T, dir, pkg string, types ...string) {
	t.Helper()
	if len(types) == 0 {
		t.Fatalf("guard: Derived with no types — a guard that checks nothing passes silently")
	}

	derived := map[string]bool{}
	for _, name := range types {
		derived[name] = true
	}

	for _, f := range testSources(t, dir) {
		ast.Inspect(f.ast, func(n ast.Node) bool {
			lit, ok := n.(*ast.CompositeLit)
			if !ok {
				return true
			}
			sel, ok := lit.Type.(*ast.SelectorExpr)
			if !ok || !derived[sel.Sel.Name] {
				return true
			}
			if id, ok := sel.X.(*ast.Ident); !ok || id.Name != pkg {
				return true
			}
			// An empty literal is a zero value, not an invented world: it is
			// what a caller writes to say "nothing yet", and refusing it would
			// push tests into a worse shape to satisfy a guard.
			if len(lit.Elts) == 0 {
				return true
			}
			t.Errorf("%s:%d builds %s.%s by hand — it is something the engine "+
				"RETURNS, and a fixture typed here is a world that does not exist. "+
				"Call the engine's own reconcile on canned input instead",
				f.name, f.fset.Position(lit.Pos()).Line, pkg, sel.Sel.Name)
			return true
		})
	}
}
