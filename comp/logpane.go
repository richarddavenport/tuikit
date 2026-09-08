package comp

import "github.com/charmbracelet/lipgloss"

// LogPane is a stream of lines, tailing.
//
// # Where this came from
//
// The deploy tool's logspane.go, 681 lines, is the only real one of the four —
// and the only place FOLLOW is modelled. Its comment describes the behaviour
// worth keeping: following pins to the tail; scrolling up detaches; scrolling
// back to the bottom re-attaches. "The less/lazydocker feel."
//
// democtl has a log view with an offset and no follow at all, so it opens on
// the OLDEST lines in the buffer — which for a log is the wrong end. That is
// not a difference between two tools worth a parameter, it is the thing
// the deploy tool already learned.
//
// A log pane is a viewport WITHOUT a selection, which is why it is not a List:
// there is no cursor, so there is nothing for the cursor rules to be about.
// The count means something different too — not where the selection is, but
// how far from the tail you have scrolled.
//
// Not carried over: the deploy tool's line-insertion anchoring, which keeps
// the view still when a line arrives above the viewport, and its 5000-line
// buffer cap. Both matter for a live stream and neither has anything to
// exercise it here — democtl's log lines are a fixture and arrive all at once.
// logspane.go:265 is where the anchoring lives when a tool needs it.
type LogPane struct {
	// Follow pins the view to the newest line. It starts true, because a log
	// you have just opened is one you want the end of.
	Follow bool

	// Empty is drawn when the stream has no lines — a stderr filter matching
	// nothing is an ordinary state.
	Empty string

	// Styles. Time draws the timestamp column, Text the line, and Stderr
	// replaces Text for a line from that stream.
	Time, Text, Stderr, Status, EmptyStyle *lipgloss.Style

	offset       int
	shown, count int
}

// LogLine is one line of a stream.
type LogLine struct {
	// At is preformatted, for the same reason a step's duration is: how much
	// of a timestamp to show is a tool's decision about its own readers.
	At     string
	Text   string
	Stderr bool
}

// Draw renders the pane into r, reserving its bottom row for the status.
func (p *LogPane) Draw(c *Canvas, r Rect, lines []LogLine, name Name) {
	c = c.Clip(r)
	p.count = len(lines)
	body := r
	if r.H > 1 {
		body.H = r.H - 1
	}
	p.shown = body.H

	if len(lines) == 0 {
		if p.Empty != "" {
			c.Text(body.X, body.Y, p.Empty, p.EmptyStyle, Region(name))
		}
		p.status(c, r, name)
		return
	}

	if p.Follow {
		p.offset = p.bottom()
	}
	p.offset = clamp(p.offset, 0, p.bottom())

	for i := p.offset; i < len(lines) && i-p.offset < body.H; i++ {
		y := body.Y + i - p.offset
		id := Region(name).At(i)

		style := p.Text
		if lines[i].Stderr {
			style = p.Stderr
		}
		x := c.Text(body.X, y, "  "+lines[i].At+" ", p.Time, id)
		c.Text(body.X+x, y, Truncate(lines[i].Text, body.W-x), style, id)
	}
	p.status(c, r, name)
}

// status says whether the view is at the tail, and how far from it if not.
//
// Always drawn, for the same reason a list's count is: a pane that only says
// where it is while it is moving is one you cannot tell from a pane that is
// stuck.
func (p *LogPane) status(c *Canvas, r Rect, name Name) {
	if r.H < 2 {
		return
	}
	label := "following"
	if !p.Follow {
		label = itoa(p.count-p.offset-p.shown) + " below"
	}
	c.Text(r.Right()-Width(label), r.Bottom(), label, p.Status, Region(name))
}

func (p *LogPane) bottom() int { return max(0, p.count-p.shown) }

// Scroll moves the view, detaching from the tail — and re-attaching when it
// arrives back at the bottom, which is what makes following feel like a place
// rather than a mode.
func (p *LogPane) Scroll(by int) {
	p.offset = clamp(p.offset+by, 0, p.bottom())
	p.Follow = p.offset >= p.bottom()
}

// Offset is the first line shown.
func (p *LogPane) Offset() int { return p.offset }
