package comp

import "github.com/charmbracelet/lipgloss"

// Rule is a horizontal line across a band — the thing under a header, and
// above a footer.
//
// It exists because four codebases drew it by hand and disagreed about what it
// was made of. The cloud tool drew it twice, the deploy tool once, and
// tuikit's own gallery, democtl and [Palette] between them three more times;
// five of those seven wrote the literal "─" and two asked the chrome. A tool
// on [theme.ASCIIBox] — which exists so an interface can be drawn in a font
// that has nothing — got "-" from its palette and "─" from its header in the
// same frame.
//
// So the character comes from Chrome and not from the caller. That is the
// whole point: [guard.Glyphs] cannot catch a hardcoded "─", because "─" is a
// legal glyph. It was the wrong SOURCE, not an illegal character, and the only
// fix for a wrong source is to have one source.
//
// Not a field on [Bar]. Bar means "one line with content at each end", and a
// Bar with no content that fills its own width is a second meaning wearing the
// first one's name. Nor does Bar underline itself, tempting as that is when
// six of the seven sites are a Bar with a rule beneath it: the seventh is a
// rule ABOVE a footer, and which side of a band its line falls on belongs to
// the layout rather than to whichever bar happens to be next to it.
type Rule struct {
	Style *lipgloss.Style
	// Rune replaces the chrome's own character, for the rare rule that means
	// something different from the others on the screen. Empty takes
	// Chrome.Box.Top, so a tool that changes its box set changes its rules.
	Rune string
}

// Draw paints the line across r.
//
// One row, at r.Y, the way [Bar] does — a band taller than a line is a layout
// that has room to spare, not a request for a thicker rule.
func (rl Rule) Draw(c *Canvas, r Rect, id ID) {
	if r.Empty() {
		return
	}
	glyph := rl.Rune
	if glyph == "" {
		glyph = c.Chrome().Box.Top
	}
	c.Fill(Rect{X: r.X, Y: r.Y, W: r.W, H: 1}, glyph, rl.Style, id)
}
