package comp

import "github.com/charmbracelet/lipgloss"

// Split divides a rect into two panes with a divider you can drag.
//
// # Where this came from
//
// democtl grew a draggable divider during the mouse work and kept it in the
// tool: the gap width, the minimum pane, and the arithmetic that turns a
// pointer position into a pane width were all literals in its dashboard. The
// next tool writes them again, slightly differently. This is that, extracted
// before there is a second copy to reconcile.
//
// # The gap is the divider
//
// The column between two panes is not empty space, it is the thing you grab.
// Giving it a region is what makes a drag possible at all, and it costs the
// frame nothing: an owned blank is still trimmed from the output, so a tool
// that never drags renders exactly as it did before.
//
// Its width and its character come from the canvas's chrome, so a tool that
// wants two columns between panes, or a visible seam rather than a gap,
// changes them in one place for every split it has.
type Split struct {
	// Name is the divider's region, and what a drag names.
	Name Name

	// At is the first pane's size, once somebody has dragged it. Zero means
	// Ratio — so a tool that is never dragged has no state to capture, and a
	// golden does not depend on whether anyone touched it.
	At int

	// Ratio is the default division: {1, 3} is a third. Zero is half.
	Ratio [2]int

	// Min is the narrowest either pane may become. A split that can be dragged
	// to nothing is a pane you cannot get back.
	Min int

	// Vertical stacks the panes and puts the divider between them.
	Vertical bool

	Style *lipgloss.Style
}

// Layout divides r without drawing, for a tool that needs the rects first.
//
// It takes the CANVAS rather than a gap, so it cannot be handed a different
// one from the one Draw will use. That asymmetry was the actual bug: two ways
// of getting the same rects that could disagree, and both democtl and the
// cloud tool were building a throwaway 0x0 canvas to read the gap back out of
// the chrome.
func (s *Split) Layout(c *Canvas, r Rect) (first, second Rect) {
	gap := c.Chrome().Gap
	total, at := r.W, s.at(r.W, gap)
	if s.Vertical {
		total, at = r.H, s.at(r.H, gap)
	}
	rest := max(0, total-at-gap)

	if s.Vertical {
		return Rect{X: r.X, Y: r.Y, W: r.W, H: at},
			Rect{X: r.X, Y: r.Y + at + gap, W: r.W, H: rest}
	}
	return Rect{X: r.X, Y: r.Y, W: at, H: r.H},
		Rect{X: r.X + at + gap, Y: r.Y, W: rest, H: r.H}
}

// at resolves the first pane's size, clamped so both stay usable.
func (s *Split) at(total, gap int) int {
	at := s.At
	if at == 0 {
		num, den := s.Ratio[0], s.Ratio[1]
		if den == 0 {
			num, den = 1, 2
		}
		at = total * num / den
	}
	// A minimum wider than half the space cannot be honored on both sides, so
	// the first pane gets what is left rather than the second going negative.
	room := max(0, total-gap)
	return clamp(at, min(s.Min, room), max(0, room-s.Min))
}

// Draw renders the divider and returns the two panes.
//
// The divider owns its column, so a press on it can begin a drag. It draws the
// chrome's Divider, which is a blank by default: the space between two panes
// reads as space rather than as a third thing.
func (s *Split) Draw(c *Canvas, r Rect) (first, second Rect) {
	ch := c.Chrome()
	first, second = s.Layout(c, r)
	if ch.Gap <= 0 {
		return first, second
	}

	divider, gap := ch.Divider, Rect{X: r.X + first.W, Y: r.Y, W: ch.Gap, H: r.H}
	if s.Vertical {
		divider, gap = ch.VDivider, Rect{X: r.X, Y: r.Y + first.H, W: r.W, H: ch.Gap}
	}
	if divider == "" {
		divider = " "
	}
	c.Fill(gap, divider, s.Style, Region(s.Name))
	return first, second
}

// Move nudges the divider by n, in the direction it slides.
//
// The keyboard's way in. MoveTo takes an absolute position, which is right for
// a mouse — that is what the event carries — but a key press has no position,
// so every tool with a draggable split was resolving the current layout and
// adding to it by hand:
//
//	first, _ := m.split.Layout(m.body(), m.chromeGap())
//	m.split.MoveTo(m.body().X+first.W+by, m.body())
//
// Three lines and two concepts to move a divider four columns.
func (s *Split) Move(c *Canvas, n int, r Rect) {
	// No frame yet, so there is no current position to move from. A key can
	// arrive before the first draw — the capture harness presses one and then
	// redraws — and app.Mouse.Route takes the same view of a nil canvas.
	if c == nil {
		return
	}
	first, _ := s.Layout(c, r)
	if s.Vertical {
		s.MoveTo(r.Y+first.H+n, r)
		return
	}
	s.MoveTo(r.X+first.W+n, r)
}

// MoveTo puts the divider under a pointer, for a drag.
//
// The position is absolute — the column or row the pointer is on — because
// that is what a mouse event carries, and converting it here means no tool
// does the arithmetic twice.
func (s *Split) MoveTo(pos int, r Rect) {
	if s.Vertical {
		s.At = max(1, pos-r.Y)
		return
	}
	s.At = max(1, pos-r.X)
}
