package comp

import "github.com/charmbracelet/lipgloss"

// Scrollbar is where you are in a list, as a position rather than a count.
//
// [List] already exposed everything this needs — Offset, Max and the height it
// last drew into — and drew none of it. Its status line says how much is hidden
// ("8/20"); a scrollbar says WHERE in it you are, which is the question a long
// list makes you ask and a number cannot answer at a glance.
type Scrollbar struct {
	// Total is how many rows exist, Shown how many fit, Offset the first one
	// visible. All three, because a thumb is a proportion and a position and
	// neither can be worked out from the other.
	Total, Shown, Offset int

	Track, Thumb *lipgloss.Style
}

// Draw paints the bar down r. Nothing is drawn when nothing can scroll: a
// full-height thumb says only that there is a scrollbar.
func (s Scrollbar) Draw(c *Canvas, r Rect, id ID) {
	if r.Empty() || s.Total <= 0 || s.Shown <= 0 || s.Shown >= s.Total {
		return
	}
	c = c.Clip(r)
	ch := c.Chrome()

	// At least one row. A list of two hundred thousand in a ten-row pane gives
	// a thumb of 0.0005 rows, and a scrollbar that vanishes exactly when the
	// list is longest is worse than none.
	thumb := max(1, r.H*s.Shown/s.Total)
	if thumb > r.H {
		thumb = r.H
	}

	// The travel is the rows the thumb can occupy, and the offset's range is
	// Total-Shown. Rounding to nearest rather than down, so the last row of a
	// list puts the thumb on the last row of the bar rather than one short —
	// "nearly at the end" and "at the end" are different answers.
	travel, span := r.H-thumb, s.Total-s.Shown
	top := 0
	if travel > 0 && span > 0 {
		top = (clamp(s.Offset, 0, span)*travel + span/2) / span
	}

	for i := 0; i < r.H; i++ {
		glyph, style := ch.ScrollTrack, s.Track
		if i >= top && i < top+thumb {
			glyph, style = ch.ScrollThumb, s.Thumb
		}
		c.Set(r.X, r.Y+i, glyph, style, id)
	}
}
