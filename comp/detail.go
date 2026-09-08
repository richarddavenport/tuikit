package comp

import "github.com/charmbracelet/lipgloss"

// Detail is what a pane says about the one thing you have selected.
//
// # Where this came from
//
// The last hand-drawn thing in two tools. Both had a helper that placed one
// label and one value and returned the next row, and both made the caller
// carry `y` from call to call:
//
//	y = m.field(c, inner, y, "state", svc.State.String(), style)   // democtl
//	y := m.fields(c, Rect{X: r.X, Y: r.Y + 3, ...}, fields)        // the cloud tool
//
// That is the shape of the bug they both guarded against and neither tested:
// `if y > inner.Y+inner.H-1 { return y + 1 }`, written once, in one of the two
// helpers, and not in the loops beside it. A detail pane one row too tall
// writes over the border it sits in, and every one of those `y + 1`s is a
// place to get that wrong.
//
// So the arithmetic is here. A caller says what the facts ARE; where they land
// is not a decision anybody was making on purpose.
//
// # Labels line up per block, not per pane
//
// A block's label column takes what its widest label needs. democtl padded to
// ten and the cloud tool to sixteen, both arbitrary, and the cloud tool's
// overflowed: a label longer than the pad pushed its value out of line with
// every other row. A natural width cannot do that, and it keeps a block of
// two-word labels from being dragged wide by a heading somewhere else in the
// pane.
type Detail struct {
	// Title and Subtitle name the thing. Both optional; a Detail that is all
	// blocks is a perfectly ordinary one.
	Title, Subtitle string

	Blocks []Block

	// Styles. Value falls back to nothing, which is the tool's text color.
	TitleStyle, SubtitleStyle, HeadingStyle, LabelStyle, ValueStyle *lipgloss.Style
}

// Block is one run of a detail pane: an optional heading, then either facts or
// prose.
//
// Either, not both, because a heading with facts under it is a table and a
// heading with a paragraph under it is a note, and a block that tried to be
// both would have to decide which came first on the caller's behalf.
type Block struct {
	// Heading names what is under it, and what is under it is INDENTED one
	// further because that is what a heading means. Without that the rows
	// under "tags" are level with the heading and read as three more facts
	// about the thing rather than as the tags.
	Heading string
	Facts   []Fact
	// Text is prose, wrapped to the width. For a note, an error, an
	// explanation — the things a fact cannot hold because they are sentences.
	Text string
	// Indent moves the whole block right, in Chrome.Indent columns per level,
	// for facts that belong to something named above them.
	Indent int
}

// Fact is one label and its value.
type Fact struct {
	Label, Value string
	// Style paints the VALUE. A state that is red in the list and plain here
	// is the same fact told twice, differently.
	Style *lipgloss.Style

	// Secret hides the value behind [Mask].
	//
	// The COMPONENT decides what reaches the canvas, not the caller: a tool
	// that masked by choosing which string to pass would put the real one in
	// the frame every time it got the condition backwards, and nothing could
	// tell. Pass the value and say whether it is showing.
	//
	// Which secrets are showing is the tool's, and app.Toggles is the shape
	// for it — keyed by a stable id, because a filter rebuilds the rows and an
	// index survives none of it. Masked is the state you should be in by
	// default: revealing is a keystroke, hiding should not be something you
	// have to remember.
	Secret bool
}

// shown is what actually reaches the canvas.
func (f Fact) shown() string {
	if f.Secret {
		return Mask()
	}
	return f.Value
}

// Draw renders into r and returns the row after the last one it used, so a
// caller that wants to put something underneath can.
//
// It stops at the bottom of r rather than drawing past it. That is the whole
// reason this exists.
func (d Detail) Draw(c *Canvas, r Rect, id ID) int {
	c = c.Clip(r)
	y := r.Y

	line := func(text string, style *lipgloss.Style, indent int) {
		if y > r.Bottom() {
			y++
			return
		}
		x := indent * c.Chrome().Indent
		c.Text(r.X+x, y, Truncate(text, max(0, r.W-x)), style, id)
		y++
	}

	if d.Title != "" {
		line(d.Title, d.TitleStyle, 0)
	}
	if d.Subtitle != "" {
		line(d.Subtitle, d.SubtitleStyle, 0)
	}

	for _, b := range d.Blocks {
		// A blank row before every block, so the blocks read as blocks. Not
		// before the first one when there is no title, because a pane that
		// opens with an empty line looks like it failed to draw.
		if y > r.Y {
			y++
		}
		if b.Heading != "" {
			line(b.Heading, d.HeadingStyle, b.Indent)
		}
		under := b.Indent
		if b.Heading != "" {
			under++
		}
		switch {
		case len(b.Facts) > 0:
			d.facts(c, r, &y, b, under, id)
		case b.Text != "":
			for _, wrapped := range Wrap(b.Text, max(1, r.W-under*c.Chrome().Indent)) {
				line(wrapped, d.ValueStyle, under)
			}
		}
	}
	return y
}

// facts lays out one block's label/value pairs, aligned to each other.
func (d Detail) facts(c *Canvas, r Rect, y *int, b Block, level int, id ID) {
	indent := level * c.Chrome().Indent
	width := max(0, r.W-indent)

	// The label column is natural and the value column takes the rest, so a
	// long value truncates and a long label does not shove it out of line.
	cells := make([][]string, len(b.Facts))
	for i, f := range b.Facts {
		cells[i] = []string{f.Label, f.shown()}
	}
	lines := Table{Gap: 1, Columns: []Column{{}, {Fill: true}}}.Rows(width, cells)

	// Drawn in two pieces rather than one, so a fact can color its value
	// without coloring its label — which is what makes a failed state read as
	// failed rather than as another gray row.
	labelWidth := 0
	for _, f := range b.Facts {
		labelWidth = max(labelWidth, Width(f.Label))
	}
	for i, f := range b.Facts {
		if *y > r.Bottom() {
			*y++
			continue
		}
		if f.Style == nil {
			c.Text(r.X+indent, *y, lines[i], d.LabelStyle, id)
			*y++
			continue
		}
		x := c.Text(r.X+indent, *y, fit(f.Label, labelWidth, false)+" ", d.LabelStyle, id)
		c.Text(r.X+indent+x, *y, Truncate(f.shown(), max(0, width-x)), f.Style, id)
		*y++
	}
}
