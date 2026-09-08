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
// Ideally from [spec.Command], which already carries a Key — "the keyboard
// path to this command inside the TUI". A help screen built from the spec
// cannot drift from what the keys actually do. One written out by hand drifts
// the first time somebody adds a binding and forgets, and nothing catches it,
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

	// Offset is the first line shown, for a list taller than its space.
	//
	// Without it this said "… more" and offered no way to see the rest, which is
	// the least helpful state a help screen has: it tells a reader there is
	// something they cannot reach. The database tool declined to adopt the
	// component for exactly this — 35 lines across four sections, on a 24-row
	// terminal, which is an ordinary terminal rather than a hard case (issue 50).
	//
	// [Keys.Rows] is the total, so a caller can clamp this and decide whether
	// to offer scrolling at all.
	Offset int

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

	lines := k.lines()
	off := clamp(k.Offset, 0, max(0, len(lines)-1))
	room := r.Bottom() - y + 1
	if room <= 0 {
		return y
	}

	// A marker costs a row, at each end that has something beyond it. Counted
	// before anything is drawn, or the last line is drawn and then covered.
	above, below := off > 0, false
	shown := room
	if above {
		shown--
	}
	if off+shown < len(lines) {
		below = true
		shown--
	}
	shown = max(0, shown)

	if above {
		c.Text(r.X, y, Truncate(c.Chrome().ScrollUp+" "+itoa(off)+" above", r.W), k.SectionStyle, id)
		y++
	}

	for _, line := range lines[off:min(len(lines), off+shown)] {
		switch {
		case line.blank:
		case line.section != "":
			c.Text(r.X, y, Truncate(line.section, r.W), k.SectionStyle, id)
		default:
			x := r.X + c.Text(r.X, y, Pad(line.hint.Key, width)+"  ", k.KeyStyle, id)
			c.Text(x, y, Truncate(line.hint.Label, max(0, r.Right()-x+1)), k.LabelStyle, id)
		}
		y++
	}

	// A help screen that quietly omits half the keys is worse than one that
	// says it was cut, because a reader believes it — and the keys it dropped
	// are exactly the ones they came looking for. Now it also says how many,
	// which is the difference between "there is more" and "there are nine
	// more", and Offset is how a reader gets to them.
	if below {
		rest := len(lines) - (off + shown)
		c.Fill(Rect{X: r.X, Y: y, W: r.W, H: 1}, " ", nil, id)
		c.Text(r.X, y, Truncate(c.Chrome().ScrollDown+" "+itoa(rest)+" more", r.W), k.SectionStyle, id)
		y++
	}
	return y
}

// keyLine is one drawn row: a section heading, a binding, or the blank between
// sections. Flattened before drawing so the list can be windowed — which is
// what the caller's own workaround did, and the reason it existed.
type keyLine struct {
	section string
	hint    Hint
	blank   bool
}

// lines flattens the sections into the rows they occupy.
func (k Keys) lines() []keyLine {
	var out []keyLine
	for i, section := range k.Sections {
		if i > 0 {
			out = append(out, keyLine{blank: true})
		}
		if section.Name != "" {
			out = append(out, keyLine{section: section.Name})
		}
		for _, hint := range section.Keys {
			out = append(out, keyLine{hint: hint})
		}
	}
	return out
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

// Rows is how many lines the sections need, blank separators included — the
// number a caller clamps [Keys.Offset] against.
func (k Keys) Rows() int { return k.rows() }

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
