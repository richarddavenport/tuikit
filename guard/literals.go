package guard

import (
	"go/ast"
	"go/token"
	"strconv"
	"strings"

	"github.com/richarddavenport/tuikit/theme"
)

// Furniture fails when a tool draws a chrome character it typed out itself.
//
// [Glyphs] cannot catch this. It asks whether a printed character is ALLOWED,
// and `─` is allowed — being in the box set is the point of it. The defect is
// the SOURCE, not the character: a tool that writes `c.Fill(band, "─", …)` has
// hardcoded the light box set, and on [theme.ASCIIBox] — which exists so an
// interface can be drawn in a font that has nothing — it draws `─` beside the
// `-` that everything reading the chrome draws. One frame, two box sets.
//
// That shipped: five sites across three codebases, including tuikit's own
// gallery and democtl, until [comp.Rule] and decision 35.
//
// # Only inside a draw
//
// The check is a chrome character in a string passed to a canvas draw — Set,
// Text, Fill. Not every literal anywhere, which is what made this look
// intractable when it was filed (issue 48): a chrome character in a
// `Glyphs: []string{"┌", "─"}` field, or in a glyph-set declaration, is NAMING
// the character in order to allow it, and flagging those is a guard every tool
// learns to suppress.
//
// Narrowing to the call excludes both, because neither is one. A tool that
// genuinely needs a specific character in a draw passes it deliberately through
// its own named constant, which is also not a literal at the call.
func Furniture(t T, dir string, ch theme.Chrome) {
	t.Helper()

	furniture := map[rune]bool{}
	for _, r := range ch.Glyphs() {
		if r >= 128 {
			furniture[r] = true
		}
	}

	for _, f := range sources(t, dir) {
		ast.Inspect(f.ast, func(n ast.Node) bool {
			call, ok := n.(*ast.CallExpr)
			if !ok {
				return true
			}
			sel, ok := call.Fun.(*ast.SelectorExpr)
			if !ok || !drawMethod(sel.Sel.Name) {
				return true
			}
			for _, arg := range call.Args {
				lit, ok := arg.(*ast.BasicLit)
				if !ok || lit.Kind != token.STRING {
					continue
				}
				text, err := strconv.Unquote(lit.Value)
				if err != nil {
					continue
				}
				for _, r := range text {
					if !furniture[r] {
						continue
					}
					t.Errorf("%s:%d draws %q as a literal — it is chrome, so read it from "+
						"c.Chrome() (or use the component that owns it). Hardcoded, it stays %q "+
						"on a box set that does not have one, beside everything else that changed.",
						f.name, f.fset.Position(lit.Pos()).Line, string(r), string(r))
					break
				}
			}
			return true
		})
	}
}

// drawMethod names the canvas calls that put a character on screen. A method
// list rather than a type check, because the guard parses one package without
// resolving imports — the same trade [Glyphs] makes.
func drawMethod(name string) bool {
	switch name {
	case "Set", "Text", "Fill":
		return true
	}
	return strings.HasPrefix(name, "Fill")
}
