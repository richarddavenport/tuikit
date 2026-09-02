package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/theme"
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
	// Drawn once so the list knows what a frame can show, then scrolled and
	// drawn again — the offset clamps against what was actually drawn.
	draw(l, 20, 6, 20)
	l.Scroll(4)
	c := draw(l, 20, 6, 20)

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

// The count moves as you scroll. A number that says how MUCH is shown is the
// same wherever you are in the list, so it answers nothing about where that is.
func TestTheCountMovesWithTheViewport(t *testing.T) {
	l := &List{Name: services}
	before := draw(l, 20, 6, 20).String()

	l.Scroll(5)
	after := draw(l, 20, 6, 20).String()

	if !strings.Contains(before, "5/20") {
		t.Errorf("at the top the count is not 5/20:\n%s", before)
	}
	if !strings.Contains(after, "10/20") {
		t.Errorf("after scrolling five the count is not 10/20:\n%s", after)
	}
}

// A resize is not a wheel. The reveal flag keeps the wheel from snapping back
// to the cursor, but when the PANE changes size the view moved underneath the
// reader rather than because they asked — and a cursor left off screen means
// whatever is drawn beside the list describes something invisible, and the next
// key acts on it.
func TestAResizeBringsTheCursorBack(t *testing.T) {
	l := &List{Name: services, Focused: true}
	draw(l, 20, 32, 39) // thirty body rows
	l.Move(35)
	draw(l, 20, 32, 39)

	if l.Cursor() != 35 {
		t.Fatalf("the cursor is on %d", l.Cursor())
	}

	// The pane shrinks: eighteen rows, one of them the status line.
	const body = 18 - 1
	draw(l, 20, 18, 39)

	if got := l.Cursor() - l.Offset(); got < 0 || got >= body {
		t.Errorf("after the resize the cursor is %d rows into a %d-row pane", got, body)
	}
}

// The wheel still does not snap back, at the same size.
func TestTheWheelStillDoesNotSnapBackAfterAResizeRule(t *testing.T) {
	l := &List{Name: services, Focused: true}
	draw(l, 20, 18, 39)
	l.Move(30)
	draw(l, 20, 18, 39)

	before := l.Offset()
	l.Scroll(-5)
	draw(l, 20, 18, 39)

	if l.Offset() != before-5 {
		t.Errorf("the wheel snapped back: offset %d, want %d", l.Offset(), before-5)
	}
}

// A row can be more than one colour: a name with a dim count after it, a
// timestamp then a message. democtl and azctl both had to draw their own lists
// for want of this.
func TestARowCanBeSeveralStyles(t *testing.T) {
	forceColour()
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	l := &List{Name: services}
	c := NewCanvas(30, 4)
	l.Draw(c, c.Bounds(), []Row{{Spans: []Segment{
		{Text: "rg-forge", Style: &accent},
		{Text: " 5", Style: &dim},
	}}})

	got := c.String()
	if !strings.Contains(got, "38;5;205") || !strings.Contains(got, "38;5;241") {
		t.Errorf("the spans did not keep their own colours: %q", got)
	}
	if !strings.Contains(harnessStrip(got), "rg-forge 5") {
		t.Errorf("got %q", got)
	}
}

// The selection paints over them. The cursor is the reader's own mark, and a
// row that kept its colours under it would be hard to find in exactly the list
// where finding it matters.
func TestTheSelectedRowIsOneColourWhateverItsSpansSay(t *testing.T) {
	forceColour()
	dim := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	selected := lipgloss.NewStyle().Background(lipgloss.Color("57"))

	l := &List{Name: services, Selected: &selected, Focused: true}
	c := NewCanvas(30, 4)
	l.Draw(c, c.Bounds(), []Row{
		{Spans: []Segment{{Text: "rg-forge"}, {Text: " 5", Style: &dim}}},
		{Text: " another"},
	})

	first := strings.Split(c.String(), "\n")[0]
	if strings.Contains(first, "38;5;241") {
		t.Errorf("the selected row kept a span's colour: %q", first)
	}
	if !strings.Contains(first, "48;5;57") {
		t.Errorf("the selected row is not painted: %q", first)
	}
}

// harnessStrip is Strip without the import cycle — comp cannot import harness.
func harnessStrip(s string) string {
	var b strings.Builder
	for i := 0; i < len(s); i++ {
		if s[i] != 0x1b {
			b.WriteByte(s[i])
			continue
		}
		for i < len(s) && s[i] != 'm' {
			i++
		}
	}
	return b.String()
}

// A list whose selection is a CHARACTER, not only a colour.
//
// comp.Form has had a cursor marker since it was written and a List did not,
// which azctl's migration found the hard way: its resource rows are marked with
// › and the port silently dropped them. Two components with a cursor should
// agree about how a cursor is shown.
func TestAListCanMarkItsCursor(t *testing.T) {
	l := &List{Name: services, Marker: "› ", Blank: "  ", Focused: true}
	l.Move(1)
	c := draw(l, 24, 6, 4)

	lines := strings.Split(c.String(), "\n")
	if !strings.HasPrefix(lines[1], "› ") {
		t.Errorf("the cursor row is %q", lines[1])
	}
	if strings.Contains(lines[0], "›") {
		t.Errorf("a row that is not the cursor is marked: %q", lines[0])
	}
	// The blank is the marker's width, or the rows jump as you move.
	if at(lines[0], "service-0") != at(lines[1], "service-1") {
		t.Errorf("the rows move with the cursor:\n%q\n%q", lines[0], lines[1])
	}
}

// A list without one is unchanged: most selections are a colour.
func TestAListWithoutAMarkerDrawsNoIndent(t *testing.T) {
	l := &List{Name: services}
	c := draw(l, 24, 6, 3)
	if got := strings.Split(c.String(), "\n")[0]; !strings.HasPrefix(got, " service-0") {
		t.Errorf("got %q", got)
	}
}

// A child sits Chrome.Indent columns right of its header. The shape two of the
// four tools already had, in five places, each writing `"  " + line` by hand.
func TestDepthIndentsARow(t *testing.T) {
	c := NewCanvas(30, 4)
	l := &List{Name: services}
	l.Draw(c, c.Bounds(), []Row{
		{Text: "rg-forge", Lead: "▾ "},
		{Text: "vm-forge-0", Depth: 1},
		{Text: "nic-forge-0", Depth: 1},
	})

	lines := strings.Split(c.String(), "\n")
	if got := lines[0]; got != "▾ rg-forge" {
		t.Errorf("header is %q", got)
	}
	if got := lines[1]; got != "  vm-forge-0" {
		t.Errorf("child is %q, want two columns of indent", got)
	}
}

// The indent is Chrome's, so a tool that wants a tighter tree changes it in one
// place rather than reindenting every row it builds.
func TestTheIndentComesFromChrome(t *testing.T) {
	ch := theme.DefaultChrome
	ch.Indent = 4
	c := NewCanvas(30, 3).WithChrome(ch)

	l := &List{Name: services}
	l.Draw(c, c.Bounds(), []Row{{Text: "child", Depth: 1}})

	if got := strings.Split(c.String(), "\n")[0]; got != "    child" {
		t.Errorf("got %q, want four columns of indent", got)
	}
}

// A row carrying its own state glyph does not also get the cursor's mark. Two
// glyphs fighting for one column is how a tree's headers end up a character out
// of line with its children.
func TestALeadReplacesTheCursorMark(t *testing.T) {
	c := NewCanvas(30, 4)
	l := &List{Name: services, Marker: "> ", Blank: "  "}
	l.Draw(c, c.Bounds(), []Row{
		{Text: "rg-forge", Lead: "▾ "},
		{Text: "vm-forge-0", Depth: 1},
	})

	lines := strings.Split(c.String(), "\n")
	// The cursor is on row 0, which has a Lead: it keeps its own glyph.
	if got := lines[0]; got != "▾ rg-forge" {
		t.Errorf("the cursor mark displaced the header's own: %q", got)
	}
	// Row 1 has no Lead, so it gets the blank that keeps rows in line.
	if got := lines[1]; got != "    vm-forge-0" {
		t.Errorf("child is %q, want indent then the marker's blank", got)
	}
}

// The indent belongs to the row, so a row of spans is indented too — otherwise
// the one kind of row a tree needs most would be the one that ignores it.
func TestSpansAreIndentedLikeAnythingElse(t *testing.T) {
	c := NewCanvas(30, 3)
	l := &List{Name: services}
	l.Draw(c, c.Bounds(), []Row{
		{Text: "header"},
		{Depth: 1, Spans: []Segment{{Text: "vm-forge-0"}, {Text: " 3"}}},
	})

	if got := strings.Split(c.String(), "\n")[1]; got != "  vm-forge-0 3" {
		t.Errorf("got %q", got)
	}
}

// The whole row is clickable whatever its depth: the indent is drawn INTO the
// row's own cells, so clicking the blank left of a child still selects it.
func TestTheIndentIsStillTheRow(t *testing.T) {
	c := NewCanvas(30, 3)
	l := &List{Name: services}
	l.Draw(c, c.Bounds(), []Row{{Text: "header"}, {Text: "child", Depth: 1}})

	if got := c.OwnerAt(0, 1); got.Name != services || got.Index != 1 {
		t.Errorf("the indent before a child is owned by %v", got)
	}
}

// TestDrawFuncOnlyAsksForVisibleRows is the property that makes a list of two
// hundred thousand things affordable.
//
// Draw already only PAINTED what fits; the expensive part was building a []Row
// for everything so twenty of them could be shown.
func TestDrawFuncOnlyAsksForVisibleRows(t *testing.T) {
	const huge = 200000
	asked := map[int]int{}
	l := List{Name: "row"}
	c := NewCanvas(30, 10)

	l.DrawFunc(c, Rect{X: 0, Y: 0, W: 30, H: 10}, huge, func(i int) Row {
		asked[i]++
		return Row{Text: "row " + itoa(i)}
	})

	// Nine rows of body plus a status line.
	if len(asked) > 12 {
		t.Errorf("asked for %d rows to fill a 10-row pane", len(asked))
	}
	if len(asked) == 0 {
		t.Fatal("asked for no rows at all")
	}
	for i := range asked {
		if i >= 12 {
			t.Errorf("asked for row %d, which is nowhere near the screen", i)
		}
		if asked[i] != 1 {
			t.Errorf("row %d was built %d times in one frame", i, asked[i])
		}
	}
}

// TestDrawFuncScrolls: the window moves over the data without the data moving.
func TestDrawFuncScrolls(t *testing.T) {
	l := List{Name: "row", Focused: true}
	c := NewCanvas(30, 6)
	rect := Rect{X: 0, Y: 0, W: 30, H: 6}
	row := func(i int) Row { return Row{Text: "row " + itoa(i)} }

	l.DrawFunc(c, rect, 100000, row)
	l.Select(50000)

	c2 := NewCanvas(30, 6)
	l.DrawFunc(c2, rect, 100000, row)
	// No styles are set, so the frame carries no escape sequences and a plain
	// contains is enough.
	if got := c2.String(); !strings.Contains(got, "row 50000") {
		t.Errorf("after selecting row 50000 the frame is:\n%s", got)
	}
}

// TestDrawAndDrawFuncAgree, since Draw is now written in terms of DrawFunc and
// nothing else would notice if it drifted.
func TestDrawAndDrawFuncAgree(t *testing.T) {
	rows := []Row{{Text: "one"}, {Text: "two"}, {Text: "three"}}
	rect := Rect{X: 0, Y: 0, W: 20, H: 5}

	a := List{Name: "row"}
	ca := NewCanvas(20, 5)
	a.Draw(ca, rect, rows)

	b := List{Name: "row"}
	cb := NewCanvas(20, 5)
	b.DrawFunc(cb, rect, len(rows), func(i int) Row { return rows[i] })

	if ca.String() != cb.String() {
		t.Errorf("Draw and DrawFunc differ:\n%q\n%q", ca.String(), cb.String())
	}
}
