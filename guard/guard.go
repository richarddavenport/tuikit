// Package guard holds the interface closed.
//
// A tuikit tool's UI package calls these from a test. They are the mechanism —
// they know about Go source and terminal vocabulary, and nothing about what
// the tool does. Deciding to run them is the tool's business; the scaffolder
// writes the test that does, so a tool has them from its first commit.
//
// All three enforce rules that are obvious, that everyone agrees with, and
// that got broken anyway — one file at a time, by people who each had a
// reason.
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

// Tokens reports a colour that did not come from the palette, and a role in
// the palette that nothing draws with.
//
// The first was broken in six places in the deploy tool before it was a rule:
// the panel border was built from raw numbers twice, and the log pane kept its
// own pair. Nobody decided that; it happened one file at a time. A raw index
// also says what a colour IS instead of what it is FOR, which is how one value
// came to mean both "title" and "focused border" without anyone choosing that
// they should move together.
//
// The second is the other direction. A role nothing draws with is a decision
// nobody made.
func Tokens(t T, dir string, p theme.Palette, except ...Exemption) {
	t.Helper()

	sources := sources(t, dir)
	used := map[string]int{}
	for _, f := range sources {
		_, exempt := reasonFor(except, f.name)
		ast.Inspect(f.ast, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok || !isLipglossColor(call.Fun) {
				return true
			}
			if exempt {
				used[filepath.Base(f.name)]++
				return true
			}
			t.Errorf("%s:%d builds a colour from a literal — "+
				"name the ROLE in your palette and use that instead",
				f.name, f.fset.Position(call.Pos()).Line)
			return true
		})
	}
	checkExemptions(t, except, sources, used)

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
	// check came from the deploy tool, where the palette was the tool's own and
	// an unused role really was dead — the cloud tool's migration is where that
	// stopped being true.
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
// A terminal font without a glyph draws a replacement box, which reads as a
// bug rather than as decoration.
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

// Exemption names a file [Tokens] may not check, and why.
//
// # When the guard is wrong
//
// lazygit's `presentation/icons/file_icons.go` is 794 lines holding 743 hex
// colour literals. They are file-type BRAND colours — the Go gopher's blue,
// the Rust orange — and the whole point of them is that they are the same
// everywhere. Decision 28 puts colour on ANSI 0-15 so the reader's theme wins,
// and that reasoning does not reach a brand mark.
//
// Before this there were three options and all were bad: drop the guard and
// lose it everywhere, move the colours somewhere unscanned, or give up the
// icons.
//
// # Why a reason is required
//
// Because the value of the guard is that a colour outside the palette is a
// DECISION, and an exemption that records nothing has thrown that away — it is
// a suppression flag wearing a better name. This is the same bargain
// [theme.GlyphSet.With] strikes for characters, and for the same reason.
//
// # Why it cannot rot
//
// Two checks, both failures:
//
//   - An exemption naming a file that is not there. The file was renamed or
//     deleted and the exemption outlived it, which means the next file to take
//     that name is silently unguarded.
//   - An exemption on a file with no literals in it. Nothing is being excused,
//     so it should be deleted — and until it is, it is a hole nobody is using
//     and nobody will notice opening.
//
// A stale exemption is worse than no exemption, because it reads as though
// somebody checked.
type Exemption struct {
	File   string
	Reason string
}

// Except exempts one file, by base name, for a stated reason.
//
// The reason is not optional: an empty one panics, the way
// [theme.GlyphSet.With] does on a malformed pair, because both are a typo at
// the call site rather than a condition to handle at runtime.
func Except(file, reason string) Exemption {
	if file == "" || reason == "" {
		panic("guard: Except wants a file and a reason — an exemption with no reason is a suppression flag")
	}
	return Exemption{File: filepath.Base(file), Reason: reason}
}

func reasonFor(except []Exemption, name string) (string, bool) {
	for _, e := range except {
		if e.File == filepath.Base(name) {
			return e.Reason, true
		}
	}
	return "", false
}

// checkExemptions fails on an exemption that has stopped meaning anything.
func checkExemptions(t T, except []Exemption, sources []source, used map[string]int) {
	t.Helper()

	present := map[string]bool{}
	for _, f := range sources {
		present[filepath.Base(f.name)] = true
	}
	for _, e := range except {
		switch {
		case !present[e.File]:
			t.Errorf("guard: %s is exempted (%q) but is not in the scanned directory — "+
				"delete the exemption, or the next file to take that name is unguarded",
				e.File, e.Reason)
		case used[e.File] == 0:
			t.Errorf("guard: %s is exempted (%q) but builds no colours from literals — "+
				"delete the exemption rather than leaving a hole nobody is using",
				e.File, e.Reason)
		}
	}
}

// testSources is the _test.go files, which [sources] deliberately skips.
//
// The only guard that reads tests is [Derived], and it has to: the world a
// fixture invents is invented in a test file, which is exactly where nothing
// else looks.
func testSources(t T, dir string) []source {
	t.Helper()

	names, err := filepath.Glob(filepath.Join(dir, "*_test.go"))
	if err != nil {
		t.Fatalf("guard: %v", err)
	}
	if len(names) == 0 {
		t.Fatalf("guard: no test files in %s — a guard that scans nothing passes silently", dir)
	}

	out := make([]source, 0, len(names))
	for _, name := range names {
		src, err := os.ReadFile(name)
		if err != nil {
			t.Fatalf("guard: %v", err)
		}
		fset := token.NewFileSet()
		file, err := parser.ParseFile(fset, name, src, 0)
		if err != nil {
			t.Fatalf("guard: %v", err)
		}
		out = append(out, source{name: name, src: src, fset: fset, ast: file})
	}
	return out
}
