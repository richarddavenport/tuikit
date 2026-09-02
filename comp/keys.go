package comp

import "github.com/charmbracelet/lipgloss"

// Keys is the `?` overlay: every binding, by screen.
//
// Every tool writes its own footer hints and none of them has a help screen,
// because a footer names what acts on what is in FRONT of you and this names
// what exists. The two are different questions and both are worth answering.
//
// # Grouped by screen, because that is how a reader looks
//
// Not alphabetically and not by what the key does. Somebody opening this is
// asking "what can I do here", and "here" is a screen. A flat list of forty
// bindings is a reference; a list under four headings is an answer.
//
// # Where the sections should come from
//
// Ideally from [spec.Command], which already carries a Key — "the keyboard path
// to this command inside the TUI". A help screen built from the spec cannot
// drift from what the keys actually do. One written out by hand drifts the
// first time somebody adds a binding and forgets, and nothing catches it,
// because a help screen that is slightly wrong still renders perfectly.
//
// This component does not read the spec itself, because a tool's screen-level
// keys — scroll, filter, close — are not commands and never will be. Building
// the sections is the tool's; keeping them honest is [guard.Keys]'.
type Keys struct {
	Sections []KeySection

	// Overlay draws it as a centred box over the interface, like Confirm.
	// Otherwise it fills the rect it is given. An overlay is usually right:
	// "what was I looking at" is half of what a reader opens help to answer.
	Overlay bool
	Title   string

	// KeyWidth aligns the keys into a column. Zero measures them.
	KeyWidth int

	TitleStyle, SectionStyle, KeyStyle, LabelStyle, Border *lipgloss.Style
}

// KeySection is one screen's bindings.
type KeySection struct {
	Name string
	Keys []Hint
}

// Draw paints the help. Returns the row after the last one used.
func (k Keys) Draw(c *Canvas, r Rect, id ID) int {
	if k.Overlay {
		r = k.box(c, r, id)
	}
	if r.Empty() {
		return r.Y
	}
	c = c.Clip(r)
	width := k.KeyWidth
	if width == 0 {
		width = k.measure()
	}

	y := r.Y
	if k.Title != "" {
		c.Text(r.X, y, Truncate(k.Title, r.W), k.TitleStyle, id)
		y += 2
	}

	// out replaces the last visible row with a marker and stops.
	//
	// A help screen that quietly omits half the keys is worse than one that
	// says it was cut, because a reader believes it — and the keys it dropped
	// are exactly the ones they came looking for.
	out := func() bool {
		if y <= r.Bottom() {
			return false
		}
		c.Text(r.X, r.Bottom(), Truncate(c.Chrome().Ellipsis+" more", r.W), k.SectionStyle, id)
		return true
	}

	for i, section := range k.Sections {
		if i > 0 {
			y++
		}
		if section.Name != "" {
			if out() {
				return r.Bottom() + 1
			}
			c.Text(r.X, y, Truncate(section.Name, r.W), k.SectionStyle, id)
			y++
		}
		for _, hint := range section.Keys {
			if out() {
				return r.Bottom() + 1
			}
			x := r.X + c.Text(r.X, y, Pad(hint.Key, width)+"  ", k.KeyStyle, id)
			c.Text(x, y, Truncate(hint.Label, max(0, r.Right()-x+1)), k.LabelStyle, id)
			y++
		}
	}
	return y
}

// box is the overlay's rect: as tall as the content wants, inside the canvas.
func (k Keys) box(c *Canvas, r Rect, id ID) Rect {
	h := k.rows() + 2 // and two of border
	if k.Title != "" {
		h += 2
	}
	w := clamp(r.W-confirmMargin, confirmMin, confirmMax)
	box := Rect{
		X: r.X + max(0, (r.W-w)/2),
		Y: r.Y + max(0, (r.H-h)/2),
		W: w,
		H: min(h, r.H),
	}
	return Pane{Focused: true, Focus: k.Border, Border: k.Border}.Draw(c, box, id)
}

// rows is how many lines the sections need, blank separators included.
func (k Keys) rows() int {
	n := 0
	for i, section := range k.Sections {
		if i > 0 {
			n++
		}
		if section.Name != "" {
			n++
		}
		n += len(section.Keys)
	}
	return n
}

// measure is the widest key, so the labels line up in one column.
func (k Keys) measure() int {
	w := 0
	for _, section := range k.Sections {
		for _, hint := range section.Keys {
			w = max(w, Width(hint.Key))
		}
	}
	return w
}
