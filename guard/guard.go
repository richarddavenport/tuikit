// Package guard holds the interface closed.
//
// A tuikit tool's UI package calls these from a test. They are the mechanism —
// they know about Go source and terminal vocabulary, and nothing about what the
// tool does. Deciding to run them is the tool's business; the scaffolder writes
// the test that does, so a tool has them from its first commit.
//
// All three enforce rules that are obvious, that everyone agrees with, and that
// got broken anyway — one file at a time, by people who each had a reason.
package guard

import (
	"go/ast"
	"go/parser"
	"go/token"
	"os"
	"path/filepath"
	"regexp"
	"sort"
	"strings"

	"github.com/richarddavenport/tuikit/theme"
)

// T is the part of *testing.T these guards use.
//
// Narrower than testing.TB so the guards can be tested against a recorder —
// a guard that never fires is indistinguishable from a clean package, and
// "the rule is enforced" is exactly the claim worth checking. *testing.T
// satisfies it.
type T interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

// Tokens reports a colour that did not come from the palette, and a role in the
// palette that nothing draws with.
//
// The first was broken in six places in swarmctl before it was a rule: the
// panel border was built from raw numbers twice, and the log pane kept its own
// pair. Nobody decided that; it happened one file at a time. A raw index also
// says what a colour IS instead of what it is FOR, which is how one value came
// to mean both "title" and "focused border" without anyone choosing that they
// should move together.
//
// The second is the other direction. A role nothing draws with is a decision
// nobody made.
func Tokens(t T, dir string, p theme.Palette) {
	t.Helper()

	sources := sources(t, dir)
	for _, f := range sources {
		ast.Inspect(f.ast, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || !isLipglossColor(call.Fun) {
				return true
			}
			t.Errorf("%s:%d builds a colour from a literal — "+
				"name the ROLE in your palette and use that instead",
				f.name, f.fset.Position(call.Pos()).Line)
			return true
		})
	}

	var all strings.Builder
	for _, f := range sources {
		all.Write(f.src)
	}
	joined := all.String()

	// Only the EXTRA roles have to be used.
	//
	// The nine are the framework's vocabulary, not the tool's promise to draw
	// with all of them: a tool that never paints a selection background cannot
	// "drop" SelectionBG, because it is a field on a struct it inherited. This
	// check came from swarmctl, where the palette was the tool's own and an
	// unused role really was dead — azctl's migration is where that stopped
	// being true.
	//
	// An unused Extra is still dead, and still worth failing on. A tool that
	// invented a tenth meaning and then did not use it has left a name for the
	// next person to wonder about.
	for _, r := range p.Extra {
		// Matched on a selector rather than a fixed "theme." prefix: a tool may
		// import theme under any name, or hold its palette in a value of its
		// own. What is constant is that the role is reached through a dot.
		if !regexp.MustCompile(`\.` + regexp.QuoteMeta(r.Name) + `\b`).MatchString(joined) {
			t.Errorf("%s is an Extra role you named and never drew with — use it or drop it", r.Name)
		}
	}
}

// Glyphs reports a character printed by the package that the allow-list does
// not cover.
//
// A terminal font without a glyph draws a replacement box, which reads as a bug
// rather than as decoration.
func Glyphs(t T, dir string, g theme.GlyphSet) {
	t.Helper()

	offenders := map[rune][]string{}
	for _, f := range sources(t, dir) {
		ast.Inspect(f.ast, func(n ast.Node) bool {
			lit, ok := n.(*ast.BasicLit)
			if !ok || lit.Kind != token.STRING {
				return true
			}
			for _, r := range lit.Value {
				if r < 128 {
					continue
				}
				if !g.Printable(r) {
					offenders[r] = append(offenders[r], f.name)
				}
			}
			return true
		})
	}

	runes := make([]rune, 0, len(offenders))
	for r := range offenders {
		runes = append(runes, r)
	}
	sort.Slice(runes, func(i, j int) bool { return runes[i] < runes[j] })
	for _, r := range runes {
		t.Errorf("glyph %q (U+%04X) is printed by %s but not in the glyph set — "+
			"a font without it shows a box; add it deliberately or use ASCII",
			r, r, strings.Join(unique(offenders[r]), ", "))
	}
}

// source is one non-test file of the package, parsed and raw.
type source struct {
	name string
	src  []byte
	fset *token.FileSet
	ast  *ast.File
}

// sources reads the package's own files. Test files are excluded: a test may
// legitimately name a colour it is checking for, or print a glyph in a failure
// message, and neither reaches a terminal.
//
// Parsed, not matched. A regex over the source cannot tell a string literal
// from a comment that quotes one, so a comment explaining a glyph was reported
// as printing it — and "a comment can say anything" is exactly what these
// guards claim. go/parser hands back only the literals.
func sources(t T, dir string) []source {
	t.Helper()

	names, err := filepath.Glob(filepath.Join(dir, "*.go"))
	if err != nil {
		t.Fatalf("guard: %v", err)
	}
	if len(names) == 0 {
		t.Fatalf("guard: no Go files in %s — a guard that scans nothing passes silently", dir)
	}

	var out []source
	for _, name := range names {
		if strings.HasSuffix(name, "_test.go") {
			continue
		}
		src, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("guard: %v", err)
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, name, src, 0) // 0: comments dropped
		if err != nil {
			t.Fatalf("guard: %v", err)
		}
		out = append(out, source{name: name, src: src, fset: fset, ast: file})
	}
	if len(out) == 0 {
		t.Fatalf("guard: %s has only test files — nothing to guard", dir)
	}
	return out
}

// isLipglossColor reports a `lipgloss.Color(…)` conversion.
func isLipglossColor(fun ast.Expr) bool {
	sel, ok := fun.(*ast.SelectorExpr)
	if !ok || sel.Sel.Name != "Color" {
		return false
	}
	pkg, ok := sel.X.(*ast.Ident)
	return ok && pkg.Name == "lipgloss"
}

func unique(in []string) []string {
	seen := map[string]bool{}
	out := make([]string, 0, len(in))
	for _, s := range in {
		if !seen[s] {
			seen[s] = true
			out = append(out, s)
		}
	}
	return out
}
