package comp

import (
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/theme"
)

const splitName Name = "split"

func TestASplitDividesByItsRatio(t *testing.T) {
	s := &Split{Name: splitName, Ratio: [2]int{1, 3}}
	first, second := s.Layout(NewCanvas(132, 10), Rect{X: 0, Y: 0, W: 132, H: 10})

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
	narrow, _ := s.Layout(NewCanvas(132, 10), Rect{W: 132, H: 10})
	wide, _ := s.Layout(gapped(4), Rect{W: 132, H: 10})

	if narrow.W != wide.W {
		t.Errorf("the gap changed the ratio: %d then %d", narrow.W, wide.W)
	}
}

// A split that can be dragged to nothing is a pane you cannot get back.
func TestNeitherPaneCanBeDraggedToNothing(t *testing.T) {
	s := &Split{Name: splitName, Min: 24}
	r := Rect{W: 132, H: 10}

	s.MoveTo(0, r)
	if first, _ := s.Layout(NewCanvas(r.W, max(r.H, 1)), r); first.W < 24 {
		t.Errorf("dragged to the left edge the first pane is %d", first.W)
	}
	s.MoveTo(200, r)
	if _, second := s.Layout(NewCanvas(r.W, max(r.H, 1)), r); second.W < 24 {
		t.Errorf("dragged past the right edge the second pane is %d", second.W)
	}
}

// A minimum wider than half the space cannot be honoured on both sides. The
// first pane gets what is left rather than the second going negative.
func TestAMinimumTooBigForTheSpaceDoesNotGoNegative(t *testing.T) {
	s := &Split{Name: splitName, Min: 40}
	first, second := s.Layout(NewCanvas(50, 4), Rect{W: 50, H: 4})

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
	if first, _ := s.Layout(NewCanvas(r.W, max(r.H, 1)), r); first.W != 30 {
		t.Errorf("dragging to column 40 of a rect starting at 10 gave %d", first.W)
	}
}

// gapped is a canvas whose chrome has a different gap, for the cases that are
// about the gap rather than about the split.
func gapped(gap int) *Canvas {
	ch := theme.DefaultChrome
	ch.Gap = gap
	return NewCanvas(132, 10).WithChrome(ch)
}

// Move is the keyboard's way in. MoveTo takes an absolute position, which is
// what a mouse event carries and what a key press does not have.
func TestMoveNudgesTheDivider(t *testing.T) {
	c := NewCanvas(132, 10)
	r := Rect{W: 132, H: 10}
	s := &Split{Name: splitName, Ratio: [2]int{1, 3}}

	before, _ := s.Layout(c, r)
	s.Move(c, +10, r)
	after, _ := s.Layout(c, r)

	if after.W != before.W+10 {
		t.Errorf("the first pane went from %d to %d, want %d", before.W, after.W, before.W+10)
	}
}

// A vertical split slides in ROWS. Moving it by columns is the bug that makes a
// stacked divider ignore the keyboard, and it costs nothing to get right.
func TestMoveOnAVerticalSplitMovesRows(t *testing.T) {
	c := NewCanvas(40, 20)
	r := Rect{W: 40, H: 20}
	s := &Split{Name: splitName, Vertical: true, Ratio: [2]int{1, 2}}

	before, _ := s.Layout(c, r)
	s.Move(c, +3, r)
	after, _ := s.Layout(c, r)

	if after.H != before.H+3 {
		t.Errorf("the first pane went from %d to %d rows, want %d", before.H, after.H, before.H+3)
	}
}

// Move clamps like a drag does, so holding a key cannot push a pane to nothing.
func TestMoveClampsToTheMinimum(t *testing.T) {
	c := NewCanvas(132, 10)
	r := Rect{W: 132, H: 10}
	s := &Split{Name: splitName, Ratio: [2]int{1, 2}, Min: 24}

	for range 20 {
		s.Move(c, -10, r)
	}
	if first, _ := s.Layout(c, r); first.W < 24 {
		t.Errorf("the first pane was squeezed to %d, below its minimum of 24", first.W)
	}
}

// Layout and Draw cannot disagree about the gap, because both read it from the
// canvas. A Layout that could be handed a different one is a Layout whose
// answer is not the layout.
func TestLayoutAndDrawAgreeAboutTheGap(t *testing.T) {
	c := gapped(4)
	r := Rect{W: 132, H: 10}
	s := &Split{Name: splitName, Ratio: [2]int{1, 3}}

	wantFirst, wantSecond := s.Layout(c, r)
	gotFirst, gotSecond := s.Draw(c, r)

	if gotFirst != wantFirst || gotSecond != wantSecond {
		t.Errorf("Draw gave %v/%v, Layout gave %v/%v", gotFirst, gotSecond, wantFirst, wantSecond)
	}
}

// A key can arrive before the first draw, so there is no frame to read the
// current position from. Doing nothing is the same answer app.Mouse.Route gives
// a nil canvas, and it beats panicking on the first keystroke.
func TestMoveBeforeTheFirstFrameDoesNothing(t *testing.T) {
	s := &Split{Name: splitName, Ratio: [2]int{1, 3}}
	s.Move(nil, +10, Rect{W: 132, H: 10})
	if s.At != 0 {
		t.Errorf("At = %d after moving with no frame, want 0", s.At)
	}
}

// Splits nest, so three panes with two draggable dividers is a composition
// rather than a component.
//
// Each Split carries its own position, so the dividers move independently. The
// thing to know when nesting: an inner split's rect comes out of the outer
// one, so after the outer moves the inner must be handed the NEW rect. A rect
// captured before the outer moved describes where the inner used to be.
func TestSplitsNest(t *testing.T) {
	full := Rect{X: 0, Y: 0, W: 60, H: 10}
	outer := &Split{Name: "outer", Ratio: [2]int{1, 2}, Min: 8}
	inner := &Split{Name: "inner", Vertical: true, Min: 3}

	c := NewCanvas(60, 10)
	left, right := outer.Draw(c, full)
	top, bottom := inner.Draw(c, right)

	for name, r := range map[string]Rect{"left": left, "top": top, "bottom": bottom} {
		if r.Empty() {
			t.Errorf("%s pane is empty", name)
		}
	}
	if _, ok := c.Region(Region("outer")); !ok {
		t.Error("the outer divider is not clickable")
	}
	if _, ok := c.Region(Region("inner")); !ok {
		t.Error("the inner divider is not clickable")
	}
	if top.Y+top.H > bottom.Y {
		t.Errorf("the inner panes overlap: top %v, bottom %v", top, bottom)
	}

	// Move each, re-deriving the inner rect from the outer as a caller must:
	// the inner split lives inside the outer's second pane, so a rect captured
	// before the outer moved describes where the inner used to be.
	outer.Move(c, 6, full)
	c2 := NewCanvas(60, 10)
	left2, right2 := outer.Draw(c2, full)
	inner.Move(c2, 1, right2)

	c3 := NewCanvas(60, 10)
	_, right3 := outer.Draw(c3, full)
	top3, bottom3 := inner.Draw(c3, right3)

	if left2.W != left.W+6 {
		t.Errorf("the outer divider moved to %d, want %d", left2.W, left.W+6)
	}
	if top3.H != top.H+1 {
		t.Errorf("the inner divider moved to %d, want %d", top3.H, top.H+1)
	}
	// And the inner Min still holds inside the nested rect, which is the thing
	// that would quietly stop being true if a rect were threaded through wrong:
	// asked to move further than Min allows, it stops rather than obeying.
	inner.Move(c3, 20, right3)
	c4 := NewCanvas(60, 10)
	_, right4 := outer.Draw(c4, full)
	_, bottom4 := inner.Draw(c4, right4)
	if bottom4.H < 3 {
		t.Errorf("the bottom pane was dragged to %d rows, under its Min of 3", bottom4.H)
	}
	_ = bottom3
}
