package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

const services Name = "services"

func rows(n int) []Row {
	out := make([]Row, n)
	for i := range out {
		out[i] = Row{Text: " service-" + itoa(i)}
	}
	return out
}

func draw(l *List, w, h, n int) *Canvas {
	c := NewCanvas(w, h)
	l.Draw(c, c.Bounds(), rows(n))
	return c
}

// Two fields, not one. Scrolling is looking around; it does not choose.
func TestTheWheelMovesTheViewportAndNotTheCursor(t *testing.T) {
	l := &List{Name: services, Focused: true}
	draw(l, 20, 6, 20)

	l.Scroll(3)
	draw(l, 20, 6, 20)

	if l.Offset() != 3 {
		t.Errorf("the offset is %d, not 3", l.Offset())
	}
	if l.Cursor() != 0 {
		t.Errorf("scrolling moved the cursor to %d", l.Cursor())
	}
}

// The cursor pulls the viewport far enough to see it, and no further.
func TestTheCursorPullsTheViewportWhenFocused(t *testing.T) {
	l := &List{Name: services, Focused: true}
	draw(l, 20, 6, 20) // five body rows

	l.Move(7)
	draw(l, 20, 6, 20)

	if l.Cursor() != 7 {
		t.Fatalf("the cursor is on %d", l.Cursor())
	}
	if want := 3; l.Offset() != want { // 7 - 5 + 1
		t.Errorf("the offset is %d, want %d — just enough to show the cursor", l.Offset(), want)
	}
}

// swarmctl follows the cursor only when the pane is focused and pgctl always
// does. swarmctl is right: an unfocused pane whose viewport jumps because its
// cursor is elsewhere moves while you are reading it.
func TestAnUnfocusedListDoesNotChaseItsCursor(t *testing.T) {
	l := &List{Name: services, Focused: false}
	l.Select(15)
	draw(l, 20, 6, 20)

	if l.Offset() != 0 {
		t.Errorf("an unfocused list scrolled to %d to chase its cursor", l.Offset())
	}
}

// Scrolling past the end has to be impossible, not discouraged — and the limit
// is what the last frame drew, not a constant.
func TestScrollingClampsToWhatWasDrawn(t *testing.T) {
	l := &List{Name: services}
	draw(l, 20, 6, 8) // five body rows, eight items

	l.Scroll(100)
	if l.Offset() != l.Max() {
		t.Errorf("scrolled to %d, past the %d there is", l.Offset(), l.Max())
	}
	if got := draw(l, 20, 6, 8); strings.TrimSpace(got.String()) == "" {
		t.Error("the list scrolled itself blank")
	}
	l.Scroll(-100)
	if l.Offset() != 0 {
		t.Errorf("scrolling back up stopped at %d", l.Offset())
	}
}

// A short list cannot scroll at all.
func TestAListThatFitsDoesNotScroll(t *testing.T) {
	l := &List{Name: services}
	draw(l, 20, 10, 3)

	l.Scroll(5)
	if l.Offset() != 0 {
		t.Errorf("a list that fits scrolled to %d", l.Offset())
	}
}

// The count is drawn ALWAYS. On a tall terminal where everything fits, a wheel
// that correctly does nothing is otherwise indistinguishable from a broken one.
func TestTheCountIsDrawnEvenWhenNothingIsHidden(t *testing.T) {
	l := &List{Name: services}
	got := draw(l, 20, 10, 3).String()

	if !strings.Contains(got, "3/3") {
		t.Errorf("a list that fits does not say so:\n%s", got)
	}
}

func TestTheCountSaysHowMuchIsShown(t *testing.T) {
	l := &List{Name: services}
	got := draw(l, 20, 6, 20).String() // five body rows of twenty

	if !strings.Contains(got, "5/20") {
		t.Errorf("the count does not say five of twenty are shown:\n%s", got)
	}
}

// A selection scrolled out of view says which way it went, rather than being
// dragged back — which is the same conflation as scrolling with the cursor.
func TestASelectionScrolledOutOfViewSaysWhichWay(t *testing.T) {
	l := &List{Name: services, Focused: true}
	draw(l, 30, 6, 20)

	l.Scroll(6) // the cursor is on 0, now above the window
	got := draw(l, 30, 6, 20)

	if !strings.Contains(got.String(), "↑ selected above") {
		t.Errorf("an off-screen selection left no trace:\n%s", got.String())
	}
	if l.Cursor() != 0 {
		t.Errorf("the cursor was dragged to %d instead", l.Cursor())
	}

	l.Scroll(-100)
	l.Select(19)
	l.Focused = false // so the viewport does not follow
	if got := draw(l, 30, 6, 20).String(); !strings.Contains(got, "↓ selected below") {
		t.Errorf("a selection below the window left no trace:\n%s", got)
	}
}

// The owner carries the index in the LIST. They differ the moment the viewport
// moves, and an ID meaning "row 3 of the screen" then acts on whatever scrolled
// into row 3.
func TestRowsAreOwnedByTheirIndexNotTheirRow(t *testing.T) {
	l := &List{Name: services}
	l.Scroll(0)
	c := draw(l, 20, 6, 20)
	l.Scroll(4)
	c = draw(l, 20, 6, 20)

	top := c.OwnerAt(2, 0)
	if top.Index != 4 {
		t.Errorf("after scrolling to 4, the top row is owned by %v", top)
	}
	if _, ok := c.Region(Region(services).At(4)); !ok {
		t.Error("the scrolled-to row cannot be named")
	}
	if _, ok := c.Region(Region(services).At(0)); ok {
		t.Error("a row that scrolled off is still nameable")
	}
}

// An empty list is an ordinary state, and its row is clickable across the pane.
func TestAnEmptyListSaysSo(t *testing.T) {
	l := &List{Name: services, Empty: "  nothing matches"}
	c := NewCanvas(30, 6)
	l.Draw(c, c.Bounds(), nil)

	if !strings.Contains(c.String(), "nothing matches") {
		t.Errorf("an empty list says nothing:\n%s", c.String())
	}
	if !strings.Contains(c.String(), "0/0") {
		t.Errorf("an empty list has no count:\n%s", c.String())
	}
}

// The whole row is clickable, including the blank after a short name.
func TestTheWholeRowIsClickable(t *testing.T) {
	l := &List{Name: services}
	c := draw(l, 30, 6, 5)

	if got := c.OwnerAt(28, 0); got.Index != 0 {
		t.Errorf("the far end of row 0 is owned by %v", got)
	}
}

// A list squeezed to nothing is a normal state on a short terminal, not a
// crash.
func TestAListInNoRoomDrawsNothing(t *testing.T) {
	l := &List{Name: services}
	c := NewCanvas(20, 1)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 1}, rows(20))

	if l.Offset() < 0 {
		t.Errorf("the offset went negative: %d", l.Offset())
	}
}

func TestTheCursorRowIsStyledByFocus(t *testing.T) {
	forceColour()
	selected := lipgloss.NewStyle().Background(lipgloss.Color("57"))
	unfocused := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	l := &List{Name: services, Selected: &selected, Unfocused: &unfocused, Focused: true}
	if got := draw(l, 20, 6, 5).String(); !strings.Contains(got, "\x1b[48;5;57m") {
		t.Errorf("a focused list does not show its selection:\n%q", got)
	}

	l.Focused = false
	if got := draw(l, 20, 6, 5).String(); !strings.Contains(got, "\x1b[38;5;205m") {
		t.Errorf("an unfocused list does not show its selection:\n%q", got)
	}
}

// The offset clamps against what the CURRENT frame can show, not what the last
// one could. A pane that grows — a terminal resized, a sibling pane closed —
// leaves an offset that is now past the end, and drawing from it shows blank
// rows below the last item while the items above are unreachable.
//
// This is the same shape as the bug the prototype found by scrolling: a detail
// pane clamped to a constant 20 rather than to its content, so wheeling over a
// five-line Overview scrolled it into empty space and blanked it.
func TestTheOffsetClampsAgainWhenThePaneGrows(t *testing.T) {
	l := &List{Name: services}
	draw(l, 20, 6, 8) // five body rows: the offset can reach 3
	l.Scroll(3)

	if l.Offset() != 3 {
		t.Fatalf("the offset is %d before the pane grows", l.Offset())
	}

	// Now the pane is tall enough for every row, so there is nothing to scroll.
	got := draw(l, 20, 12, 8)

	if l.Offset() != 0 {
		t.Errorf("the offset is still %d in a pane that fits everything", l.Offset())
	}
	if !strings.Contains(got.String(), "service-0") {
		t.Errorf("the first rows are unreachable:\n%s", got.String())
	}
	if !strings.Contains(got.String(), "8/8") {
		t.Errorf("the count does not say everything is shown:\n%s", got.String())
	}
}
