package ui

import (
	"testing"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
	"github.com/richarddavenport/tuikit/harness"
)

// drawn builds a model and renders once, because a region cannot be clicked
// before it has been drawn — which is the point, not a limitation.
func drawn(w, h int) *Model {
	m := New(1)
	m.SetSize(w, h)
	m.Now(fleet.Epoch)
	m.View()
	return m
}

// Clicking a row selects THAT service, addressed by name.
func TestClickingARowSelectsIt(t *testing.T) {
	m := drawn(132, 38)
	m.focus = paneDetail

	harness.Click(t, m, "services.row[3]")

	if m.cur != 3 {
		t.Errorf("the cursor is on %d, not 3", m.cur)
	}
	if m.focus != paneList {
		t.Errorf("clicking the list did not focus it")
	}
}

// Clicking the blank space after a short name still selects the row. The row
// fills its width before the text is drawn, so the CELL decides ownership
// rather than the glyph — in a string world this is a manual rect calculation
// that is wrong the first time.
func TestClickingPastTheNameStillSelectsTheRow(t *testing.T) {
	m := drawn(132, 38)
	c := m.Canvas()

	r, ok := c.Region(comp.Region(regServicesRow).At(2))
	if !ok {
		t.Fatal("no third row was drawn")
	}
	if got := c.OwnerAt(r.X+r.W-1, r.Y); got.Index != 2 {
		t.Errorf("the last column of row 2 is owned by %v", got)
	}
}

// Clicking a tab selects it and focuses the pane it belongs to.
func TestClickingATabSelectsIt(t *testing.T) {
	m := drawn(132, 38)
	harness.Click(t, m, "detail.tab[2]")

	if m.tab != 2 {
		t.Errorf("the tab is %d, not 2", m.tab)
	}
	if m.focus != paneDetail {
		t.Errorf("clicking a tab did not focus the detail pane")
	}
}

// The wheel scrolls the pane under the POINTER, and it moves a viewport rather
// than the selection. Wiring it to the cursor instead makes scrolling appear to
// pick services at random.
func TestTheWheelScrollsTheViewportNotTheSelection(t *testing.T) {
	m := drawn(132, 12) // short enough that the list overflows
	before := m.cur

	harness.Wheel(t, m, "services", +2)

	if m.listOffset == 0 {
		t.Errorf("the wheel did not scroll the list")
	}
	if m.cur != before {
		t.Errorf("the wheel moved the selection from %d to %d", before, m.cur)
	}
}

// Scrolling past the end has to be impossible, not discouraged. Clamping to a
// constant instead is how a five-line pane scrolls into empty space and blanks
// itself.
func TestScrollingStopsAtTheEnd(t *testing.T) {
	m := drawn(132, 12)
	harness.Wheel(t, m, "services", +50)

	if m.listOffset > m.listMax {
		t.Errorf("scrolled to %d, past the %d there is", m.listOffset, m.listMax)
	}
	if got := m.View(); got == "" {
		t.Error("the pane scrolled itself blank")
	}

	harness.Wheel(t, m, "services", -50)
	if m.listOffset != 0 {
		t.Errorf("scrolling back up stopped at %d", m.listOffset)
	}
}

// After scrolling, a click has to select the service that is THERE. This is the
// one the owner ID design exists for: an ID that meant "row 3 of the screen"
// passes every test until something scrolls, and then acts on the wrong item
// rather than failing visibly.
func TestAfterScrollingAClickSelectsWhatIsUnderIt(t *testing.T) {
	m := drawn(132, 12)
	harness.Wheel(t, m, "services", +1)

	c := m.Canvas()
	inner, ok := c.Region(comp.Region(regServices))
	if !ok {
		t.Fatal("no list was drawn")
	}
	// The first row on screen is now the second service, so clicking it must
	// select index 1 and not index 0.
	top := c.OwnerAt(inner.X+2, inner.Y+2)
	if top.Index != m.listOffset {
		t.Fatalf("the top row is owned by %v, want index %d", top, m.listOffset)
	}
	harness.Click(t, m, top.String())
	if m.cur != m.listOffset {
		t.Errorf("clicking the top row selected %d, want %d", m.cur, m.listOffset)
	}
}

// Dragging the divider moves it, and a drag survives the pointer outrunning
// what it grabbed — which happens on every real drag.
func TestDraggingTheDividerMovesIt(t *testing.T) {
	m := drawn(132, 38)
	before := m.listWidth()

	harness.Drag(t, m, "split", +10)

	if got := m.listWidth(); got != before+10 {
		t.Errorf("the divider moved to %d, want %d", got, before+10)
	}
	if m.dragging {
		t.Error("the drag never ended")
	}
}

// A divider cannot be dragged to nothing, or it is a pane you cannot get back.
func TestTheDividerKeepsBothPanesUsable(t *testing.T) {
	m := drawn(132, 38)
	harness.Drag(t, m, "split", -100)

	if m.listWidth() < minPane {
		t.Errorf("the list pane is %d columns, under the %d minimum", m.listWidth(), minPane)
	}
}

// Right-click opens the menu for what is under it.
func TestRightClickOpensTheMenuForThatRow(t *testing.T) {
	m := drawn(132, 38)
	harness.RClick(t, m, "services.row[4]")

	if m.menu == nil {
		t.Fatal("no menu opened")
	}
	if m.cur != 4 {
		t.Errorf("the menu opened on row 4 but the cursor is on %d", m.cur)
	}
	if len(m.menu.items) == 0 {
		t.Error("the menu has no actions")
	}
}

// The menu opens from the keyboard too, at the cursor.
//
// Not a convenience. A multiplexer that captures right-click gets the event
// first and this application never learns it happened — there is no protocol
// for asking — so an action whose only path is a context menu is broken for
// everyone inside herdr or tmux, and the tool cannot detect it to say so.
func TestTheMenuOpensFromTheKeyboard(t *testing.T) {
	m := drawn(132, 38)
	harness.Press(m, "j", "j", "m")

	if m.menu == nil {
		t.Fatal("m did not open the menu")
	}
	if m.menu.on.Index != 2 {
		t.Errorf("the menu opened on %v, not the cursor", m.menu.on)
	}
}

// Every action in the menu carries the key that does the same thing, so the two
// paths are one list and cannot drift.
func TestEveryMenuActionNamesItsKey(t *testing.T) {
	m := drawn(132, 38)
	harness.Press(m, "m")

	for _, item := range m.menu.items {
		if item.key == "" {
			t.Errorf("%q has no keyboard path", item.label)
		}
		if item.do == nil {
			t.Errorf("%q does nothing", item.label)
		}
	}
}

// Choosing from the menu does the same thing the key does.
func TestTheMenuAndTheKeyDoTheSameThing(t *testing.T) {
	byKey := drawn(132, 38)
	harness.Press(byKey, "D")

	byMenu := drawn(132, 38)
	harness.Press(byMenu, "m")
	harness.Click(t, byMenu, "menu.item[1]") // Deploy

	if byMenu.confirm == nil {
		t.Fatal("the menu did not open the confirm")
	}
	if byKey.confirm.title != byMenu.confirm.title {
		t.Errorf("the key asks %q and the menu asks %q", byKey.confirm.title, byMenu.confirm.title)
	}
}

// A modal takes the mouse the way it takes the keyboard: clicking the frame
// behind a question is not an answer to it.
func TestAModalTakesTheMouse(t *testing.T) {
	m := drawn(132, 38)
	harness.Press(m, "D")
	before := m.cur

	harness.Click(t, m, "services.row[5]")

	if m.cur != before {
		t.Errorf("a click behind the modal moved the cursor to %d", m.cur)
	}
	if m.confirm == nil {
		t.Error("a click behind the modal dismissed it")
	}
}

// A script names regions, and a name that was not drawn is an error rather than
// a click into empty space that reports success.
func TestAScriptCannotNameARegionThatWasNotDrawn(t *testing.T) {
	m := drawn(132, 38)
	rec := &recorder{}

	harness.Click(rec, m, "services.row[999]")

	if !rec.failed {
		t.Error("clicking a region nobody drew was accepted")
	}
}

// The whole thing as a script, which is how a capture will read.
func TestAScriptDrivesTheInterfaceByName(t *testing.T) {
	m := drawn(132, 38)

	harness.Script(t, m, `
		# select a service, look at its Config tab, then open its menu
		click services.row[2]
		click detail.tab[1]
		rclick services.row[2]
	`)

	if m.cur != 2 || m.tab != 1 || m.menu == nil {
		t.Errorf("after the script: cur=%d tab=%d menu=%v", m.cur, m.tab, m.menu != nil)
	}
}

// recorder stands in for *testing.T so a failure can be asserted rather than
// causing one.
type recorder struct{ failed bool }

func (r *recorder) Helper()                           {}
func (r *recorder) Errorf(string, ...any)             { r.failed = true }
func (r *recorder) Fatalf(format string, args ...any) { r.failed = true }
