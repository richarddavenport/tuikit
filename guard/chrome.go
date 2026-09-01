package guard

import "github.com/richarddavenport/tuikit/theme"

// Chrome fails when the interface's own furniture draws a character the tool
// has not allowed.
//
// The hole this closes: guard.Glyphs reads the TOOL's source, and the box
// corners, chevrons, separators and scroll markers are drawn by comp on the
// tool's behalf. They never appear in a string literal anyone scans. A tool
// that narrowed its glyph set for a font without chevrons would still have had
// them printed for it, and nothing would have said so.
//
// So the rule is checked where the characters actually live. Choosing
// theme.RoundedBox means adding ╭╮╰╯ to the set — a deliberate decision rather
// than four replacement boxes discovered on somebody else's terminal.
func Chrome(t T, c theme.Chrome, g theme.GlyphSet) {
	t.Helper()

	for _, r := range c.Glyphs() {
		if !g.Printable(r) {
			t.Errorf("the chrome draws %q (%U) but it is not in the glyph set — "+
				"add it deliberately, or choose a box set the set already allows", r, r)
		}
	}
}
