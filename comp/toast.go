package comp

import "github.com/charmbracelet/lipgloss"

// Toast is something the interface has to say: a failure, a warning, a result.
//
// # Where this came from
//
// The deploy tool's errorView. Its shape is what failed, against which
// target, the cleaned root cause wrapped to the terminal — and, in its own
// words, "a hint for the failures we know how to fix".
//
// That last part is the component. An error that names what went wrong and
// stops is a dead end: the reader is told they have a problem and left to
// guess. A Hint is the difference between "connection refused" and "connection
// refused — is the tunnel up? try `the database tool connect`", and it is the
// field most likely to be left empty, so it is a field rather than something
// you are expected to append to Body.
//
// Bounded the way a modal is, and for the same reason: the root cause of a
// failure is arbitrary text from somewhere else, and an unbounded box drew
// itself off the side of the screen the first time anyone tried it.
type Toast struct {
	Title string
	Body  string
	// Hint is what to do about it. Empty is allowed and usually means nobody
	// has worked out the answer yet, which is worth being able to see.
	Hint string

	// Anchor is the corner of the rect it sits in. Zero is the top right,
	// which is out of the way of a cursor and of anything being read.
	Anchor Anchor

	// Max, Min and Margin bound it, as Confirm does.
	Max, Min, Margin int

	// Styles. Accent draws the title, and is where a caller distinguishes a
	// failure from a result — by role, not by a Kind enum this package would
	// have to keep in step with every palette.
	Accent, Border, BodyStyle, HintStyle *lipgloss.Style
}

// Anchor is which corner something sits in.
type Anchor int

// The four corners. TopRight is the zero value.
const (
	TopRight Anchor = iota
	TopLeft
	BottomRight
	BottomLeft
)

// Draw puts the toast in a corner of r and returns the box it took.
func (t Toast) Draw(c *Canvas, r Rect, id ID) Rect {
	// Nothing leaves r, whatever the arithmetic below concludes. Min is a
	// preference and the rect is a fact: a toast asked for 24 columns inside a
	// pane 10 wide used to draw 24 of them, over whatever was beside it.
	c = c.Clip(r)

	maxW, minW, margin := or(t.Max, 48), or(t.Min, 24), or(t.Margin, 2)
	w := min(clamp(r.W-margin*2, minW, maxW), r.W)

	body := Wrap(t.Body, w-4)
	lines := len(body)
	if t.Hint != "" {
		lines += 1 + len(Wrap(t.Hint, w-4))
	}
	h := lines + 4 // title, blank, and two of border

	box := Rect{W: w, H: min(h, r.H)}
	switch t.Anchor {
	case TopLeft:
		box.X, box.Y = r.X+margin, r.Y+margin/2
	case BottomRight:
		box.X, box.Y = r.Right()-w-margin+1, r.Bottom()-h-margin/2+1
	case BottomLeft:
		box.X, box.Y = r.X+margin, r.Bottom()-h-margin/2+1
	default:
		box.X, box.Y = r.Right()-w-margin+1, r.Y+margin/2
	}
	box.X, box.Y = max(r.X, box.X), max(r.Y, box.Y)

	inner := Pane{Border: t.Border, Focus: t.Border}.Draw(c, box, id)
	if inner.Empty() {
		return box
	}

	c.Text(inner.X, inner.Y, Truncate(t.Title, inner.W), t.Accent, id)
	y := inner.Y + 2
	for _, line := range body {
		c.Text(inner.X, y, line, t.BodyStyle, id)
		y++
	}
	if t.Hint != "" {
		y++
		for _, line := range Wrap(t.Hint, inner.W) {
			c.Text(inner.X, y, line, t.HintStyle, id)
			y++
		}
	}
	return box
}
