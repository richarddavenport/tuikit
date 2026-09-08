package comp

import "github.com/charmbracelet/lipgloss"

// Sparkline is a number over time.
//
// # Where this came from
//
// Two tools, and the second one says something the first could not.
//
// k9s wrote `internal/tchart/` — a sparkline, a gauge and a dot matrix — and
// `view/pulse.go` is a dashboard of them.
//
// bottom **vendored 55 kB of ratatui's own Chart widget** into
// `src/canvas/components/time_series/vendored.rs`, with a note saying to keep
// it in sync, and did the same to `grid.rs`. It forked its framework's chart
// rather than use it, and took on the maintenance knowingly.
//
// So the interesting question is not whether to have one. It is what bottom
// forked FOR, and the answer is the rule below.
//
// # Right-aligned against now
//
// A metric over time has a fixed right edge and a ragged left one. The newest
// sample is always in the last column; a series shorter than the pane leaves
// the LEFT blank, and one longer drops its oldest.
//
// A general chart centers its data or stretches it to fit, and both are wrong
// here: the bar under your cursor moves when a sample arrives, and two
// sparklines of different lengths stop being comparable. That is what bottom's
// vendored copy specializes for and it is the whole reason this type exists
// rather than a `Chart`.
//
// # What it is not
//
// Not a chart. No axes, no legend, no labels, no second series. Those are a
// bigger question and only one tool in the survey has them, so [Meter] stays
// the answer for one number and this is the answer for one number over time.
type Sparkline struct {
	// Max is the top of the scale. Zero takes the largest value on screen,
	// which is what a sparkline usually wants — the shape of the last minute,
	// not its absolute size.
	//
	// Set it when two sparklines are meant to be compared, because auto-scaled
	// series with different peaks look identical and are not.
	Max float64

	// Style draws the bars. Nil takes the canvas's default.
	Style *lipgloss.Style
}

// bars is the ramp, one eighth per step. Index 0 is a value greater than zero
// but too small to see, so a live series never looks like a dead one.
var bars = []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"}

// Draw renders values into r, newest at the right edge.
//
// A rect taller than one row draws a column chart: full blocks stack, and the
// topmost one is the partial. One row is the classic sparkline and is the same
// arithmetic with a height of one.
func (s Sparkline) Draw(c *Canvas, r Rect, values []float64, id ID) {
	if r.Empty() || len(values) == 0 {
		return
	}
	c = c.Clip(r)

	// The newest sample is the last one, and it belongs in the last column.
	// Everything else follows from that.
	shown := values
	if len(shown) > r.W {
		shown = shown[len(shown)-r.W:]
	}
	left := r.X + r.W - len(shown)

	top := s.Max
	if top <= 0 {
		for _, v := range shown {
			if v > top {
				top = v
			}
		}
	}
	if top <= 0 {
		// Every sample is zero. Drawing nothing is right: a flat line at the
		// bottom of an auto-scaled chart would be indistinguishable from a
		// series at its peak.
		return
	}

	// Eighths of a row, across every row of the rect.
	steps := r.H * len(bars)
	for i, v := range shown {
		eighths := int(v / top * float64(steps))
		eighths = clamp(eighths, 0, steps)
		if v > 0 && eighths == 0 {
			// Present but too small to see. One eighth beats nothing, because
			// a gap in a sparkline means no data.
			eighths = 1
		}
		s.column(c, r, left+i, eighths, id)
	}
}

// column draws one bar, from the bottom up.
func (s Sparkline) column(c *Canvas, r Rect, x, eighths int, id ID) {
	for row := 0; row < r.H; row++ {
		y := r.Bottom() - row
		remaining := eighths - row*len(bars)
		switch {
		case remaining <= 0:
			return
		case remaining >= len(bars):
			c.Set(x, y, bars[len(bars)-1], s.Style, id)
		default:
			c.Set(x, y, bars[remaining-1], s.Style, id)
		}
	}
}
