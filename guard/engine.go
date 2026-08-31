package guard

import (
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
// exist when this was written: democtl's fleet imports fmt, math/rand and time,
// and swarmctl's internal/engine imports nothing from charmbracelet across some
// forty files. What an engine imports is stdlib plus its own domain SDK — pgx,
// the Azure SDK, the Docker SDK.
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
