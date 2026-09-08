package comp

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/theme"
)

// Bar is one line with content at each end: a header, or a key hint bar.
//
// # Where this came from
//
// The deploy tool's footerLine, the cloud tool's footer and
// democtl's header and footer are the same line four times over. The shared
// shape is content left, optional content right, and a decision about what to
// do when they will not both fit.
//
// That decision is the one meaningful difference. democtl drops the right side
// only when the two would COLLIDE; the deploy tool drops it unless the left
// still gets 24 columns, because its right side is a version number and its
// left side is what you can press. The deploy tool is protecting the primary
// content from being squeezed to nothing by decoration, which is a real
// requirement — so MinLeft carries it, and zero is democtl's rule.
//
// The left side wins in every case. A bar too narrow for both is a bar whose
// right-hand side was never load-bearing.
type Bar struct {
	Left, Right []Segment

	// MinLeft is the columns the left side must still get for the right side
	// to be drawn at all. Zero drops the right side only on a collision.
	MinLeft int
}

// Segment is a run of text with a style. A bar's ends are lists of them
// because a header is rarely one colour: a tool's name, its description and
// its health are three things, and each says something different.
type Segment struct {
	Text  string
	Style *lipgloss.Style
}

// Draw renders the bar along the top row of r.
func (b Bar) Draw(c *Canvas, r Rect, id ID) {
	c = c.Clip(r)
	right := width(b.Right)
	room := r.W - right

	if right > 0 && room > max(b.MinLeft, width(b.Left)) {
		drawSegments(c, r.Right()-right+1, r.Y, b.Right, id)
	}
	drawSegments(c, r.X, r.Y, b.Left, id)
}

func width(segs []Segment) int {
	var n int
	for _, s := range segs {
		n += Width(s.Text)
	}
	return n
}

func drawSegments(c *Canvas, x, y int, segs []Segment, id ID) {
	for _, s := range segs {
		x += c.Text(x, y, s.Text, s.Style, id)
	}
}

// Hint is one key and what it does.
//
// A type rather than a formatted string, so the key hint bar and a context
// menu are built from the same thing. mouse.md's rule is that the two paths to
// an action must be ONE LIST rather than a list and a keymap maintained beside
// it, and this is the list.
type Hint struct {
	Key, Label string
}

// Hints renders hints the way all four tools already write them: `j/k move ·
// enter connect · q quit`. The separator is in the glyph set once, here,
// instead of in every footer string in every tool.
//
// What goes IN the bar is the tool's business, and the deploy tool learned the
// rule the expensive way: the footer used to list every action on every panel,
// which grew a letter per feature and read as a menu of things mostly not
// applicable. Hints name what acts on what is focused, right now. While a
// prompt is capturing keys, the scroll keys are not among them, and listing
// them is a lie.
func Hints(hints ...Hint) string {
	parts := make([]string, 0, len(hints))
	for _, h := range hints {
		switch {
		case h.Key == "":
			parts = append(parts, h.Label)
		case h.Label == "":
			parts = append(parts, h.Key)
		default:
			parts = append(parts, h.Key+" "+h.Label)
		}
	}
	return strings.Join(parts, theme.DefaultChrome.Separator)
}

// KeyHints is the whole bar, for the common case of hints and nothing else.
func KeyHints(c *Canvas, r Rect, id ID, style *lipgloss.Style, hints ...Hint) {
	Bar{Left: []Segment{{Text: "  " + Hints(hints...), Style: style}}}.Draw(c, r, id)
}
