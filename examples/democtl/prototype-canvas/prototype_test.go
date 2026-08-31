// PROTOTYPE — throwaway. Drives the thing with synthetic mouse events, which is
// also the point: if a test can click, so can the capture harness and so can an
// agent.
package main

import (
	"fmt"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

// find returns the first coordinate owned by id, which is how a script would say
// "click list.row[3]" without knowing where that is.
func find(m *model, id string) (int, int, bool) {
	for y := range m.canvas.H {
		for x := range m.canvas.W {
			if m.canvas.OwnerAt(x, y) == id {
				return x, y, true
			}
		}
	}
	return 0, 0, false
}

func click(m *model, id string, button tea.MouseButton) bool {
	x, y, ok := find(m, id)
	if !ok {
		return false
	}
	m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: button})
	m.View()
	return true
}

func boot() *model {
	m := newModel()
	m.w, m.h = 100, 30
	m.View()
	return m
}

func TestClickSelectsTheRowItLandedOn(t *testing.T) {
	m := boot()
	if !click(m, "list.row[3]", tea.MouseButtonLeft) {
		t.Fatal("row 3 is not on the canvas")
	}
	if m.sel != 3 {
		t.Errorf("sel = %d, want 3", m.sel)
	}
	if m.focus != ownerList {
		t.Errorf("focus = %s, want the list", m.focus)
	}
}

// Clicking the blank space after a short name has to select the row. That only
// works because the row fills its width before drawing text — the cell decides
// ownership, not the glyph.
func TestTheWholeRowIsClickableNotJustItsText(t *testing.T) {
	m := boot()
	x, y, _ := find(m, "list.row[2]")
	far := x
	for m.canvas.OwnerAt(far+1, y) == "list.row[2]" {
		far++
	}
	m.Update(tea.MouseMsg{X: far, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft})
	if m.sel != 2 {
		t.Errorf("clicking the padding at column %d selected %d, want 2", far, m.sel)
	}
}

func TestClickingATabSwitchesIt(t *testing.T) {
	m := boot()
	if !click(m, "detail.tab[2]", tea.MouseButtonLeft) {
		t.Fatal("tab 2 is not on the canvas")
	}
	if m.tab != 2 || m.focus != ownerDetl {
		t.Errorf("tab = %d focus = %s", m.tab, m.focus)
	}
}

func TestDraggingTheSplitterResizes(t *testing.T) {
	m := boot()
	before := m.ratio

	x, y, ok := find(m, ownerSplit)
	if !ok {
		t.Fatal("the splitter owns no cells")
	}
	m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft})
	if !m.dragging {
		t.Fatal("pressing the splitter did not start a drag")
	}
	// Motion well past the divider: a drag has to survive the pointer outrunning
	// the thing it grabbed, which it always does.
	m.Update(tea.MouseMsg{X: x + 18, Y: y, Action: tea.MouseActionMotion, Button: tea.MouseButtonLeft})
	m.Update(tea.MouseMsg{X: x + 18, Y: y, Action: tea.MouseActionRelease, Button: tea.MouseButtonLeft})
	m.View()

	if m.dragging {
		t.Error("release did not end the drag")
	}
	if m.ratio <= before {
		t.Errorf("ratio %.2f did not grow from %.2f", m.ratio, before)
	}
	if _, _, ok := find(m, ownerSplit); !ok {
		t.Error("the splitter vanished after the drag")
	}
}

// Scroll what the pointer is over, not what has focus. The list has focus here;
// the wheel is over the detail pane.
func TestTheWheelScrollsWhatItIsOver(t *testing.T) {
	m := boot()
	click(m, "detail.tab[2]", tea.MouseButtonLeft) // events tab, which has enough lines to scroll
	m.focus = ownerList
	sel := m.sel

	x, y, ok := find(m, ownerDetl)
	if !ok {
		t.Fatal("the detail pane owns no cells")
	}
	m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})

	if m.scroll == 0 {
		t.Error("the wheel over the detail pane did not scroll it")
	}
	if m.sel != sel {
		t.Error("the wheel scrolled the focused pane instead of the one under the pointer")
	}
}

func TestRightClickOpensAMenuThatIsModal(t *testing.T) {
	m := boot()
	if !click(m, "list.row[1]", tea.MouseButtonRight) {
		t.Fatal("row 1 is not on the canvas")
	}
	if m.menu == nil {
		t.Fatal("right click opened no menu")
	}

	// The menu is drawn last, so it owns its cells outright — no compositing.
	x, y, ok := find(m, "menu.item[1]")
	if !ok {
		t.Fatal("the menu's items own no cells")
	}

	// A click behind the menu must not reach what is behind it.
	before := m.sel
	m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft})
	m.View()
	if m.menu != nil {
		t.Error("choosing an item left the menu open")
	}
	if !strings.HasPrefix(m.last, "chose: ") {
		t.Errorf("the click did not select an item: %q", m.last)
	}
	if m.sel != before {
		t.Error("the click fell through the menu to the row underneath")
	}
}

// The claim that made the canvas worth testing: drawing past the edge is not a
// bug to catch, it is a coordinate that does not exist.
func TestNothingCanBeDrawnOutsideTheCanvas(t *testing.T) {
	for _, w := range []int{40, 80, 132} {
		m := newModel()
		m.w, m.h = w, 24
		m.menu = []string{"View logs", "Deploy", "Restart", "Remove"}
		m.menuAt = [2]int{w - 3, 20} // a menu deliberately hung off the right edge
		view := m.View()

		lines := strings.Split(view, "\n")
		if len(lines) != 24 {
			t.Errorf("at %d columns: %d lines, want 24", w, len(lines))
		}
		for i, line := range lines {
			// No ANSI here to confuse the count: styles are per cell and the
			// prototype's styles are colour only, so a rune count is the width.
			if got := len([]rune(stripANSI(line))); got != w {
				t.Errorf("at %d columns: line %d is %d wide", w, i+1, got)
			}
		}
	}
}

func stripANSI(s string) string {
	var b strings.Builder
	var in bool
	for _, r := range s {
		switch {
		case r == 0x1b:
			in = true
		case in && (r == 'm' || r == 'K'):
			in = false
		case !in:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// --- what Richard found by scrolling ------------------------------------

// The wheel scrolls a viewport. It does not move the cursor.
//
// The first version had no viewport at all — the list drew from index 0 always
// — so the wheel was wired to the selection as a stand-in, and scrolling
// appeared to pick services at random. Conflating "look around" with "choose"
// is a real bug and not only a missing feature.
func TestTheWheelScrollsTheListWithoutMovingTheCursor(t *testing.T) {
	m := boot()
	sel := m.sel

	x, y, _ := find(m, "list.row[1]")
	m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})
	m.View()

	if m.sel != sel {
		t.Errorf("the wheel moved the selection from %d to %d", sel, m.sel)
	}
	if m.listTop == 0 {
		t.Error("the wheel did not scroll the list")
	}
}

// A click after scrolling has to select the service that is THERE, which only
// works because the owner ID carries the absolute index rather than the screen
// row.
func TestClickingAfterScrollingSelectsWhatIsOnScreen(t *testing.T) {
	m := boot()
	x, y, _ := find(m, "list.row[1]")
	m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})
	m.View()

	top := m.listTop
	if top == 0 {
		t.Fatal("nothing scrolled, so there is nothing to check")
	}

	// The second visible row is now services[top+1].
	want := top + 1
	if !click(m, rowOwner(want), tea.MouseButtonLeft) {
		t.Fatalf("row %d is not on screen after scrolling to %d", want, top)
	}
	if m.sel != want {
		t.Errorf("sel = %d, want %d — the owner ID is a screen position, not an identity", m.sel, want)
	}
}

// Scrolling past the end has to be impossible, not merely discouraged. With the
// clamp missing, wheeling over a five-line Overview blanked the pane.
func TestNeitherPaneCanScrollPastItsContent(t *testing.T) {
	m := boot()

	x, y, _ := find(m, ownerDetl)
	for range 12 {
		m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})
		m.View()
	}
	if !strings.Contains(m.View(), "replicas") {
		t.Error("the detail pane scrolled past its content and went blank")
	}

	lx, ly, _ := find(m, "list.row["+itoa(m.listTop)+"]")
	for range 40 {
		m.Update(tea.MouseMsg{X: lx, Y: ly, Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})
		m.View()
	}
	if m.listTop > len(m.services)-m.listRows {
		t.Errorf("listTop %d is past the last full page (%d)", m.listTop, len(m.services)-m.listRows)
	}
	if _, _, ok := find(m, rowOwner(len(m.services)-1)); !ok {
		t.Error("the last service scrolled off the bottom into empty space")
	}
}

func itoa(i int) string { return fmt.Sprintf("%d", i) }

// A selection scrolled out of view has to leave a trace. Not by dragging the
// cursor into the viewport — that is the bug this replaced — but by saying it
// is up there. Otherwise the detail pane describes a service nothing on screen
// points at, and the next key press acts on something invisible.
func TestASelectionScrolledOutOfViewSaysWhereItWent(t *testing.T) {
	m := boot()
	if m.sel != 0 {
		t.Fatalf("expected to start at the top, got %d", m.sel)
	}

	x, y, _ := find(m, "list.row[1]")
	m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonWheelDown})
	view := m.View()

	if m.listTop == 0 {
		t.Fatal("nothing scrolled")
	}
	if m.sel != 0 {
		t.Errorf("the cursor was dragged along to %d", m.sel)
	}
	if !strings.Contains(view, "selected above") {
		t.Errorf("the selection scrolled off with no indication:\n%s", view)
	}

	// And scrolling back removes it, rather than leaving a permanent banner.
	for range 4 {
		m.Update(tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonWheelUp})
		m.View()
	}
	if strings.Contains(m.View(), "selected above") {
		t.Error("the indicator stayed after the selection came back into view")
	}
}
