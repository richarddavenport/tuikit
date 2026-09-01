package comp

import (
	"strings"
	"testing"
)

const splitName Name = "split"

func TestASplitDividesByItsRatio(t *testing.T) {
	s := &Split{Name: splitName, Ratio: [2]int{1, 3}}
	first, second := s.Layout(Rect{X: 0, Y: 0, W: 132, H: 10}, 1)

	if first.W != 44 {
		t.Errorf("a third of 132 is %d", first.W)
	}
	if second.X != 45 || second.W != 87 {
		t.Errorf("the second pane is %+v", second)
	}
	if first.W+1+second.W != 132 {
		t.Errorf("the panes and the gap are %d columns, not 132", first.W+1+second.W)
	}
}

// The ratio is taken against the whole width, gaps included — the same rule the
// screendoc prototype found the hard way. A third of 132 is 44 and a third of
// 131 is 43, and the pane is a column narrower for the rest of the screen's
// life.
func TestTheRatioIncludesTheGap(t *testing.T) {
	s := &Split{Name: splitName, Ratio: [2]int{1, 3}}
	narrow, _ := s.Layout(Rect{W: 132, H: 10}, 1)
	wide, _ := s.Layout(Rect{W: 132, H: 10}, 4)

	if narrow.W != wide.W {
		t.Errorf("the gap changed the ratio: %d then %d", narrow.W, wide.W)
	}
}

// A split that can be dragged to nothing is a pane you cannot get back.
func TestNeitherPaneCanBeDraggedToNothing(t *testing.T) {
	s := &Split{Name: splitName, Min: 24}
	r := Rect{W: 132, H: 10}

	s.MoveTo(0, r)
	if first, _ := s.Layout(r, 1); first.W < 24 {
		t.Errorf("dragged to the left edge the first pane is %d", first.W)
	}
	s.MoveTo(200, r)
	if _, second := s.Layout(r, 1); second.W < 24 {
		t.Errorf("dragged past the right edge the second pane is %d", second.W)
	}
}

// A minimum wider than half the space cannot be honoured on both sides. The
// first pane gets what is left rather than the second going negative.
func TestAMinimumTooBigForTheSpaceDoesNotGoNegative(t *testing.T) {
	s := &Split{Name: splitName, Min: 40}
	first, second := s.Layout(Rect{W: 50, H: 4}, 1)

	if first.W < 0 || second.W < 0 {
		t.Errorf("panes are %d and %d", first.W, second.W)
	}
	if first.W+1+second.W != 50 {
		t.Errorf("the panes and gap are %d columns, not 50", first.W+1+second.W)
	}
}

// The gap is the thing you grab, and giving it a region costs the frame
// nothing: an owned blank is still trimmed from the output.
func TestTheDividerOwnsItsColumnAndDrawsNothing(t *testing.T) {
	c := NewCanvas(20, 3)
	s := &Split{Name: splitName, Ratio: [2]int{1, 2}}
	first, _ := s.Draw(c, c.Bounds())

	if got := c.OwnerAt(first.W, 1); got != Region(splitName) {
		t.Errorf("the gap column is owned by %v", got)
	}
	if got := c.String(); strings.TrimSpace(got) != "" {
		t.Errorf("a blank divider drew something: %q", got)
	}
	if _, ok := c.Region(Region(splitName)); !ok {
		t.Error("the divider cannot be named by a script")
	}
}

// A tool that wants a seam rather than a gap changes one field, for every split
// it has.
func TestTheDividerIsTheChromes(t *testing.T) {
	ch := NewCanvas(20, 3).Chrome()
	ch.Divider = "│"
	c := NewCanvas(20, 3).WithChrome(ch)

	s := &Split{Name: splitName, Ratio: [2]int{1, 2}}
	s.Draw(c, c.Bounds())

	if !strings.Contains(c.String(), "│") {
		t.Errorf("the seam is missing:\n%s", c.String())
	}
}

// The gap's width is the chrome's too, so "more air between the panes" is one
// number in one place.
func TestTheGapWidthIsTheChromes(t *testing.T) {
	ch := NewCanvas(40, 3).Chrome()
	ch.Gap = 4
	c := NewCanvas(40, 3).WithChrome(ch)

	s := &Split{Name: splitName, Ratio: [2]int{1, 2}}
	first, second := s.Draw(c, c.Bounds())

	if second.X-first.W != 4 {
		t.Errorf("the gap is %d columns", second.X-first.W)
	}
	if first.W+4+second.W != 40 {
		t.Errorf("the panes and gap are %d columns, not 40", first.W+4+second.W)
	}
}

// Stacked, for a tool that splits top from bottom.
func TestAVerticalSplitDividesRows(t *testing.T) {
	c := NewCanvas(20, 11)
	s := &Split{Name: splitName, Vertical: true, Ratio: [2]int{1, 2}}
	top, bottom := s.Draw(c, c.Bounds())

	if top.H != 5 || bottom.Y != 6 || bottom.H != 5 {
		t.Errorf("top %+v bottom %+v", top, bottom)
	}
	if got := c.OwnerAt(3, top.H); got != Region(splitName) {
		t.Errorf("the divider row is owned by %v", got)
	}
}

// The pointer's position is what a mouse event carries, so converting it here
// means no tool does the arithmetic twice.
func TestMoveToTakesAnAbsolutePosition(t *testing.T) {
	s := &Split{Name: splitName, Min: 4}
	r := Rect{X: 10, Y: 0, W: 60, H: 4}

	s.MoveTo(40, r)
	if first, _ := s.Layout(r, 1); first.W != 30 {
		t.Errorf("dragging to column 40 of a rect starting at 10 gave %d", first.W)
	}
}
