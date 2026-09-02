package comp

import (
	"fmt"
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Waiting is a region whose contents have not arrived yet.
//
// [Spinner] is the inline case — one turning glyph beside a row that is doing
// something. This is the other one: a pane, or a whole screen, that has nothing
// to show because it is still reading.
//
// # Draw it INSIDE the interface, not instead of it
//
// The mistake this exists to stop is replacing the frame with the message.
// azctl's first estate read drew "reading the estate…" and nothing else — no
// header, no panes, no key hints — so the interface did not look busy, it
// looked absent, and a slow subscription looked like a crash.
//
// A pane that is empty and framed says "the thing you asked for is coming".
// A blank terminal with a sentence on it says "something went wrong". They are
// one line of code apart and they are not the same interface.
//
// So: draw your header, your panes and your footer as usual, and put this in
// the rect the contents would have filled.
type Waiting struct {
	// Label is what is being waited FOR, in the present participle: "reading
	// the estate", "connecting to the swarm". Not "please wait", which says
	// nothing a spinner has not already said.
	Label string
	// Detail is an optional second line — where it is reading from, or what it
	// will do next. A wait is more tolerable when it is specific.
	Detail string

	Spinner Spinner
	// Since is when the wait started. Zero means do not count.
	Since time.Time
	// After is how long a wait has to last before the elapsed time is shown.
	// Zero takes two seconds.
	//
	// A threshold rather than always, because a counter that appears for 300ms
	// and vanishes is a flicker, and a wait that is over before you read it
	// never needed a number. Past the threshold it is the difference between
	// "working" and "hung".
	After time.Duration

	Style, DetailStyle *lipgloss.Style
}

const waitingAfter = 2 * time.Second

// Draw centres the wait in r.
//
// Centred because the region is otherwise empty: a message in the top-left of a
// large blank pane reads as content that failed to fill it, and the same words
// in the middle read as a placeholder. Nothing is drawn if there is no room,
// which is an ordinary state for a pane squeezed by a narrow terminal.
func (w Waiting) Draw(c *Canvas, r Rect, at time.Time, id ID) {
	if r.Empty() {
		return
	}
	c = c.Clip(r)

	line := w.Spinner.Frame(at) + " " + w.Label
	if e := w.elapsed(at); e != "" {
		line += "  " + e
	}
	y := r.Y + r.H/2
	if w.Detail != "" && r.H > 1 {
		y-- // keep the pair centred, not the first line
	}
	c.Text(centre(r, line), y, Truncate(line, r.W), w.Style, id)

	if w.Detail != "" && y+1 <= r.Bottom() {
		detail := Truncate(w.Detail, r.W)
		c.Text(centre(r, detail), y+1, detail, w.DetailStyle, id)
	}
}

// elapsed is how long this has been going, once that is worth saying.
func (w Waiting) elapsed(at time.Time) string {
	if w.Since.IsZero() {
		return ""
	}
	after := w.After
	if after == 0 {
		after = waitingAfter
	}
	d := at.Sub(w.Since)
	if d < after {
		return ""
	}
	if d < time.Minute {
		return fmt.Sprintf("%ds", int(d.Seconds()))
	}
	return fmt.Sprintf("%dm %02ds", int(d.Minutes()), int(d.Seconds())%60)
}

// centre is the x a string starts at to sit in the middle of r.
func centre(r Rect, s string) int {
	return r.X + max(0, (r.W-Width(s))/2)
}
