package comp

import (
	"regexp"
	"strings"
	"testing"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/richarddavenport/tuikit/theme"
)

// uncolored strips SGR sequences. Local rather than harness.Strip, because
// harness imports comp and a test cannot import it back.
var uncolored = regexp.MustCompile(`\x1b\[[0-9;]*m`)

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

// the deploy tool follows the cursor only when the pane is focused and the
// database tool always does. The deploy tool is right: an unfocused pane whose
// viewport jumps because its cursor is elsewhere moves while you are reading
// it.
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
// that correctly does nothing is otherwise indistinguishable from a broken
// one.
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
// moves, and an ID meaning "row 3 of the screen" then acts on whatever
// scrolled into row 3.
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

// An empty list is an ordinary state, and its row is clickable across the
// pane.
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
	forceColor()
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
// same wherever you are in the list, so it answers nothing about where that
// is.
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
// whatever is drawn beside the list describes something invisible, and the
// next key acts on it.
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

// A row can be more than one color: a name with a dim count after it, a
// timestamp then a message. democtl and the cloud tool both had to draw their
// own lists for want of this.
func TestARowCanBeSeveralStyles(t *testing.T) {
	forceColor()
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
		t.Errorf("the spans did not keep their own colors: %q", got)
	}
	if !strings.Contains(harnessStrip(got), "rg-forge 5") {
		t.Errorf("got %q", got)
	}
}

// The selection paints over them. The cursor is the reader's own mark, and a
// row that kept its colors under it would be hard to find in exactly the list
// where finding it matters.
func TestTheSelectedRowIsOneColorWhateverItsSpansSay(t *testing.T) {
	forceColor()
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
		t.Errorf("the selected row kept a span's color: %q", first)
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

// A list whose selection is a CHARACTER, not only a color.
//
// comp.Form has had a cursor marker since it was written and a List did not,
// which the cloud tool's migration found the hard way: its resource rows are
// marked with › and the port silently dropped them. Two components with a
// cursor should agree about how a cursor is shown.
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

// A list without one is unchanged: most selections are a color.
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

// The indent is Chrome's, so a tool that wants a tighter tree changes it in
// one place rather than reindenting every row it builds.
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

// A row carrying its own state glyph gets the cursor's mark AS WELL, in its
// own column to the left of the indent.
//
// This used to be one or the other, on the grounds that two glyphs fighting
// for one column is how a tree's headers end up out of line with its children.
// That worry is real and it is answered by giving the marker a column of its
// own rather than by dropping it: the header's text and the child's still line
// up, and the cursor is now visible on a list whose rows have leads.
//
// Which matters because such a list otherwise carried its selection entirely
// in Selected's background — invisible in a pipe, in a golden, and to a reader
// who cannot see color.
func TestAMarkAndALeadBothGetADrawn(t *testing.T) {
	c := NewCanvas(30, 4)
	l := &List{Name: services, Marker: "> ", Blank: "  "}
	l.Draw(c, c.Bounds(), []Row{
		{Text: "rg-forge", Lead: "▾ "},
		{Text: "vm-forge-0", Depth: 1},
	})

	lines := strings.Split(c.String(), "\n")
	// The cursor is on row 0: its marker, then no indent, then its own glyph.
	if got := lines[0]; got != "> ▾ rg-forge" {
		t.Errorf("the cursor row is %q", got)
	}
	// Row 1 is not the cursor and has no lead: the blank, then the indent.
	if got := lines[1]; got != "    vm-forge-0" {
		t.Errorf("child is %q, want the blank then the indent", got)
	}
	// And the thing the old rule was protecting still holds. Measured in
	// COLUMNS: strings.Index counts bytes, and ▾ is three of them, which is
	// the miscount Width exists to prevent.
	head, child := columnOf(t, lines[0], "rg-forge"), columnOf(t, lines[1], "vm-forge-0")
	if head != child {
		t.Errorf("the header's text starts at column %d and its child's at %d:\n%q\n%q",
			head, child, lines[0], lines[1])
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
		// The cursor's row is asked about twice: once to settle the cursor
		// against the rows that exist, and once to draw it. Every other row is
		// built exactly once.
		want := 1
		if i == l.Cursor() {
			want = 2
		}
		if asked[i] > want {
			t.Errorf("row %d was built %d times in one frame, want %d", i, asked[i], want)
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

// grouped is a list with headings the cursor may not hold: rows 0, 4 and 7.
func grouped() []Row {
	return []Row{
		{Text: "GROUP ONE", Skip: true},
		{Text: "alpha"}, {Text: "beta"}, {Text: "gamma"},
		{Text: "GROUP TWO", Skip: true},
		{Text: "delta"}, {Text: "epsilon"},
		{Text: "GROUP THREE", Skip: true},
		{Text: "zeta"},
	}
}

func drawGrouped(l *List, rows []Row) {
	c := NewCanvas(30, 12)
	l.Draw(c, Rect{X: 0, Y: 0, W: 30, H: 12}, rows)
}

// The first frame lands on the first row the cursor may hold, not on row 0 —
// which in a grouped list is a heading.
func TestTheCursorStartsOnASelectableRow(t *testing.T) {
	l := &List{Name: services, Focused: true}
	drawGrouped(l, grouped())
	if l.Cursor() != 1 {
		t.Errorf("the cursor opened on row %d, want 1 — row 0 is a heading", l.Cursor())
	}
}

// Moving passes over headings, so j/k never appear to do nothing.
func TestMovingPassesOverHeadings(t *testing.T) {
	l := &List{Name: services, Focused: true}
	rows := grouped()
	drawGrouped(l, rows)

	for _, want := range []int{2, 3, 5, 6, 8} {
		l.Move(1)
		drawGrouped(l, rows)
		if l.Cursor() != want {
			t.Fatalf("moving down reached row %d, want %d", l.Cursor(), want)
		}
	}
	for _, want := range []int{6, 5, 3, 2, 1} {
		l.Move(-1)
		drawGrouped(l, rows)
		if l.Cursor() != want {
			t.Fatalf("moving up reached row %d, want %d", l.Cursor(), want)
		}
	}
}

// At the end it stays put, rather than sticking on a trailing heading — and it
// does not queue the moves it could not make.
func TestMovingOffTheEndStaysPut(t *testing.T) {
	l := &List{Name: services, Focused: true}
	rows := grouped()
	drawGrouped(l, rows)

	for range 10 {
		l.Move(1)
	}
	drawGrouped(l, rows)
	if l.Cursor() != 8 {
		t.Fatalf("the cursor is on %d, want the last selectable row 8", l.Cursor())
	}
	// One press back should move one row, not undo ten queued ones.
	l.Move(-1)
	drawGrouped(l, rows)
	if l.Cursor() != 6 {
		t.Errorf("after ten downs and one up the cursor is on %d, want 6", l.Cursor())
	}
}

// Clicking a heading does nothing. A cursor that lands somewhere you did not
// click is worse than a click that is ignored.
func TestSelectingAHeadingMovesOffIt(t *testing.T) {
	l := &List{Name: services, Focused: true}
	rows := grouped()
	drawGrouped(l, rows)

	l.Select(4) // GROUP TWO
	drawGrouped(l, rows)
	if l.Cursor() == 4 {
		t.Error("the cursor sat on a heading")
	}
	if l.Cursor() != 5 {
		t.Errorf("the cursor went to %d, want 5 — the row under the heading", l.Cursor())
	}
}

// A list of nothing but headings has no cursor to find, and must not spin
// looking for one.
func TestAListOfOnlyHeadings(t *testing.T) {
	l := &List{Name: services, Focused: true}
	rows := []Row{{Text: "ONE", Skip: true}, {Text: "TWO", Skip: true}}

	done := make(chan struct{})
	go func() {
		drawGrouped(l, rows)
		l.Move(1)
		drawGrouped(l, rows)
		l.Move(-5)
		drawGrouped(l, rows)
		close(done)
	}()
	select {
	case <-done:
	case <-time.After(2 * time.Second):
		t.Fatal("a list of only headings hung looking for a cursor")
	}
}

// The wheel moves the viewport and never the cursor, so a skipped row is just
// a row to it.
func TestTheWheelIgnoresSkippedRows(t *testing.T) {
	l := &List{Name: services, Focused: true}
	rows := grouped()
	c := NewCanvas(30, 5)
	l.Draw(c, Rect{X: 0, Y: 0, W: 30, H: 5}, rows)

	before := l.Cursor()
	l.Scroll(2)
	l.Draw(c, Rect{X: 0, Y: 0, W: 30, H: 5}, rows)
	if l.Cursor() != before {
		t.Errorf("the wheel moved the cursor from %d to %d", before, l.Cursor())
	}
	if l.Offset() != 2 {
		t.Errorf("the viewport is at %d, want 2", l.Offset())
	}
}

// The status line counts every row: it describes the viewport, not the
// selection.
func TestTheCountIncludesHeadings(t *testing.T) {
	muted := lipgloss.NewStyle()
	l := &List{Name: services, Focused: true, Status: &muted}
	c := NewCanvas(30, 5)
	l.Draw(c, Rect{X: 0, Y: 0, W: 30, H: 5}, grouped())
	// The status is last-shown/total. Nine rows, headings included.
	if got := c.String(); !strings.Contains(got, "/9") {
		t.Errorf("the status does not count all nine rows:\n%s", got)
	}
}

// A list with no skipped rows behaves exactly as it did.
func TestAnUnmarkedListIsUnchanged(t *testing.T) {
	rows := make([]Row, 6)
	for i := range rows {
		rows[i] = Row{Text: "row " + itoa(i)}
	}
	l := &List{Name: services, Focused: true}
	drawGrouped(l, rows)
	l.Move(3)
	drawGrouped(l, rows)
	if l.Cursor() != 3 {
		t.Errorf("the cursor is on %d, want 3", l.Cursor())
	}
}

// A glyph that IS the state keeps its color on the selected row.
//
// From issue 44, and from a person: "when highlighting I can't see the color
// of the dot." A selected row is otherwise one color whatever its spans say,
// which is right for a label and wrong for a status glyph — the one row a
// reader is looking at became the one row whose status they could not read.
func TestALeadKeepsItsColorWhenTheRowIsSelected(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)

	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	sel := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("7"))

	l := &List{Name: "conns", Selected: &sel, Focused: true}
	c := NewCanvas(20, 3)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 3}, []Row{
		{Lead: "●", LeadStyle: &green, Text: " local"},
	})

	out := c.String()
	if !strings.Contains(out, green.Render("●")) {
		t.Errorf("the state glyph lost its color under the selection:\n%q", out)
	}
	// The label still takes the selection — this must not become "every row
	// keeps its own colors", which is the change that makes a cursor hard to
	// find in exactly the list where finding it matters.
	if !strings.Contains(out, "local") {
		t.Fatal("the label was not drawn")
	}
	if !strings.Contains(out, opening(sel)+" local") {
		t.Errorf("the label did not take the selection style:\n%q", out)
	}
}

// Without a LeadStyle nothing changes, so a list that has never heard of this
// draws exactly what it drew before.
func TestALeadWithNoStyleOfItsOwnTakesTheRowsStyle(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	sel := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("7"))

	l := &List{Name: "conns", Selected: &sel, Focused: true}
	c := NewCanvas(20, 3)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 3}, []Row{{Lead: "▸", Text: " group"}})

	if out := c.String(); !strings.Contains(out, opening(sel)+"▸ group") {
		t.Errorf("an unstyled lead did not take the selection:\n%q", out)
	}
}

// opening is the escape sequence a style starts with, so a test can assert
// "this text is drawn in that style" without depending on the row's padding.
func opening(s lipgloss.Style) string {
	before, _, _ := strings.Cut(s.Render("\x00"), "\x00")
	return before
}

// The lead is still in the marker's column on an unselected row, so a styled
// lead does not shift the text of one row out of line with the others.
func TestAStyledLeadDoesNotMoveTheText(t *testing.T) {
	green := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))

	styled := NewCanvas(20, 4)
	plain := NewCanvas(20, 4)
	rows := []Row{{Lead: "●", Text: " one"}, {Lead: "○", Text: " two"}}

	(&List{Name: "a"}).Draw(plain, Rect{X: 0, Y: 0, W: 20, H: 4}, rows)
	rows[0].LeadStyle = &green
	(&List{Name: "a"}).Draw(styled, Rect{X: 0, Y: 0, W: 20, H: 4}, rows)

	if uncolored.ReplaceAllString(styled.String(), "") != uncolored.ReplaceAllString(plain.String(), "") {
		t.Errorf("styling a lead moved the row:\n%q\n%q", styled.String(), plain.String())
	}
}

// NoStatus gives the row back. From issue 46: the database tool stacks five
// lists in one column, and five status rows are a quarter of an 80x24 body
// spent on counters that read 3/3 beside panel titles already saying (3).
func TestNoStatusGivesTheRowBackToTheRows(t *testing.T) {
	const h = 5

	with := &List{Name: "a"}
	c := NewCanvas(20, h)
	with.Draw(c, Rect{X: 0, Y: 0, W: 20, H: h}, rows(10))

	without := &List{Name: "a", NoStatus: true}
	c2 := NewCanvas(20, h)
	without.Draw(c2, Rect{X: 0, Y: 0, W: 20, H: h}, rows(10))

	if without.Shown() != with.Shown()+1 {
		t.Errorf("NoStatus showed %d rows, want one more than %d", without.Shown(), with.Shown())
	}
	if strings.Contains(c2.String(), "/10") {
		t.Errorf("the counter was drawn anyway:\n%s", c2.String())
	}
	if !strings.Contains(c.String(), "/10") {
		t.Errorf("the default stopped drawing its counter:\n%s", c.String())
	}
}

// Overhead is the component answering the question a tool was encoding as a
// constant, so a tool laying out several lists cannot go silently wrong when
// the answer changes.
func TestOverheadSaysWhatTheListSpendsOnItself(t *testing.T) {
	if got := (&List{}).Overhead(); got != 1 {
		t.Errorf("a list with a status row reports %d", got)
	}
	if got := (&List{NoStatus: true}).Overhead(); got != 0 {
		t.Errorf("a list without one reports %d", got)
	}

	// And it agrees with what Draw actually does, which is the half that keeps
	// it from becoming another untested assertion.
	for _, l := range []*List{{Name: "a"}, {Name: "a", NoStatus: true}} {
		const h = 6
		l.Draw(NewCanvas(20, h), Rect{X: 0, Y: 0, W: 20, H: h}, rows(20))
		if want := h - l.Overhead(); l.Shown() != want {
			t.Errorf("Overhead says %d but Draw showed %d of %d rows", l.Overhead(), l.Shown(), h)
		}
	}
}

// Clamping the cursor to where it already is must not eat a pending move.
//
// Issue 45: every one of these tools arrived at "clamp every cursor whenever
// the data changes" independently, from when a cursor was a plain int. With
// the deferred Move that clamp reads as Select(Cursor()), which zeroed the
// move the arrow key had just recorded. The key did nothing and nothing
// errored.
func TestClampingToWhereTheCursorAlreadyIsKeepsAPendingMove(t *testing.T) {
	l := &List{Name: "a", Focused: true}
	c := NewCanvas(20, 6)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 6}, rows(10))

	l.Move(1)
	l.Select(clampTo(l.Cursor(), 9)) // the tools' clampCursors, in one line
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 6}, rows(10))

	if l.Cursor() != 1 {
		t.Errorf("the cursor is on %d; the move was discarded by the clamp", l.Cursor())
	}
}

// A real Select still discards, because a click means that row and a queued
// arrow key landing on top of a click would be worse than the bug.
func TestSelectingADifferentRowStillDiscardsAPendingMove(t *testing.T) {
	l := &List{Name: "a", Focused: true}
	c := NewCanvas(20, 6)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 6}, rows(10))

	l.Move(3)
	l.Select(7)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 6}, rows(10))

	if l.Cursor() != 7 {
		t.Errorf("the cursor is on %d, not the row that was clicked", l.Cursor())
	}
}

// And a clamp that really does move the cursor — the case the pattern was
// written for — still moves it.
func TestAClampThatActuallyMovesTheCursorStillDoes(t *testing.T) {
	l := &List{Name: "a", Focused: true}
	c := NewCanvas(20, 6)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 6}, rows(10))
	l.Select(8)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 6}, rows(10))

	// The list shrinks under the cursor, and the tool clamps.
	l.Select(clampTo(l.Cursor(), 2))
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 6}, rows(3))

	if l.Cursor() != 2 {
		t.Errorf("the cursor is on %d after the list shrank to 3 rows", l.Cursor())
	}
}

func clampTo(i, hi int) int {
	if i > hi {
		return hi
	}
	return max(0, i)
}

// A row can pin content to its right edge, which is what Bar has always done
// and what three call sites in the board were doing by hand (issue 53).
func TestARowCanPinContentToItsRightEdge(t *testing.T) {
	l := &List{Name: "lanes"}
	c := NewCanvas(30, 3)
	l.Draw(c, Rect{X: 0, Y: 0, W: 30, H: 3}, []Row{
		{Text: " ACTIVE", Right: []Segment{{Text: "3/5"}}},
	})

	line := strings.Split(c.String(), "\n")[0]
	if !strings.HasSuffix(strings.TrimRight(line, " "), "3/5") {
		t.Errorf("the right-hand content is not at the edge: %q", line)
	}
	if !strings.Contains(line, "ACTIVE") {
		t.Errorf("the left-hand content was lost: %q", line)
	}
}

// The width is measured in columns, which is the bug the reporter hit twice:
// len() on a string holding a wide glyph puts the content past the edge, where
// the canvas silently does not draw it.
func TestTheRightEdgeIsMeasuredInColumns(t *testing.T) {
	l := &List{Name: "lanes"}
	c := NewCanvas(20, 2)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 2}, []Row{
		{Text: "lane", Right: []Segment{{Text: "› 12"}}},
	})

	line := strings.Split(c.String(), "\n")[0]
	if !strings.Contains(line, "› 12") {
		t.Errorf("a row ending in a wide glyph lost it: %q", line)
	}
}

// Lead comes out of the same width. The reporter's second bug: forgetting it
// pushed every count one column past the edge.
func TestTheRightEdgeAccountsForTheLead(t *testing.T) {
	l := &List{Name: "lanes"}
	c := NewCanvas(16, 2)
	l.Draw(c, Rect{X: 0, Y: 0, W: 16, H: 2}, []Row{
		{Lead: "▾", Text: " lane", Right: []Segment{{Text: "9"}}},
	})

	line := strings.Split(c.String(), "\n")[0]
	if !strings.Contains(line, "9") {
		t.Errorf("the right-hand content was pushed off the row: %q", line)
	}
	if Width(strings.TrimRight(line, " ")) > 16 {
		t.Errorf("the row is wider than its rect: %q", line)
	}
}

// Dropped rather than overlapped when there is no room, because a count
// written over the end of a name is two pieces of information and neither is
// readable.
func TestTheRightEdgeIsDroppedWhenThereIsNoRoom(t *testing.T) {
	l := &List{Name: "lanes"}
	c := NewCanvas(4, 2)
	l.Draw(c, Rect{X: 0, Y: 0, W: 4, H: 2}, []Row{
		{Text: "lane", Right: []Segment{{Text: "a very long count"}}},
	})

	if line := strings.Split(c.String(), "\n")[0]; strings.Contains(line, "very") {
		t.Errorf("it overlapped instead of being dropped: %q", line)
	}
}

// A brand-new list has no range. The zero value has to mean "nothing
// selected", or every list reports a selection of row 0 before anyone touches
// it — and a caller acting on Range() would act on it.
func TestTheZeroListHasNoRange(t *testing.T) {
	if lo, hi, ok := (&List{Name: "a"}).Range(); ok {
		t.Errorf("a brand-new list reports a range %d..%d", lo, hi)
	}
}

// Shift-down once selects two rows, which is what every other list does.
func TestExtendAnchorsAtTheCursor(t *testing.T) {
	l := &List{Name: "a", Focused: true}
	c := NewCanvas(20, 8)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 8}, rows(10))

	l.Extend(1)
	l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 8}, rows(10))

	lo, hi, ok := l.Range()
	if !ok || lo != 0 || hi != 1 {
		t.Errorf("one extend gave %d..%d ok=%v, want 0..1", lo, hi, ok)
	}
}

// The range is ordered whichever way it was dragged, so a caller never asks.
func TestARangeIsOrderedEitherWay(t *testing.T) {
	l := &List{Name: "a", Focused: true}
	c := NewCanvas(20, 8)
	draw := func() { l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 8}, rows(10)) }
	draw()

	l.Select(5)
	draw()
	l.Extend(-2)
	draw()

	if lo, hi, ok := l.Range(); !ok || lo != 3 || hi != 5 {
		t.Errorf("extending upward gave %d..%d ok=%v, want 3..5", lo, hi, ok)
	}
}

// Extend, turn round, and extend past the start: the case a pair of bounds
// gets wrong, because it does not remember which end was the anchor.
func TestARangeCanBeDraggedBackThroughItsAnchor(t *testing.T) {
	l := &List{Name: "a", Focused: true}
	c := NewCanvas(20, 12)
	draw := func() { l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 12}, rows(20)) }
	draw()

	l.Select(5)
	draw()
	l.Extend(3) // 5..8
	draw()
	l.Extend(-6) // through the anchor to 2, so the range is 2..5
	draw()

	if lo, hi, ok := l.Range(); !ok || lo != 2 || hi != 5 {
		t.Errorf("got %d..%d ok=%v, want 2..5 — the anchor moved", lo, hi, ok)
	}
}

// An ordinary arrow key drops the selection. A range that survived one would
// be a range a reader cannot get rid of.
func TestMovingWithoutExtendingDropsTheRange(t *testing.T) {
	l := &List{Name: "a", Focused: true}
	c := NewCanvas(20, 8)
	draw := func() { l.Draw(c, Rect{X: 0, Y: 0, W: 20, H: 8}, rows(10)) }
	draw()

	l.Extend(2)
	draw()
	if _, _, ok := l.Range(); !ok {
		t.Fatal("no range to drop")
	}

	l.Move(1)
	draw()
	if lo, hi, ok := l.Range(); ok {
		t.Errorf("an arrow key left a range %d..%d", lo, hi)
	}
}

// A row with its own lead still gets the cursor's marker. It used to get one
// or the other, so a list with a status glyph had no cursor at all once the
// color was stripped.
func TestALeadDoesNotSwallowTheMarker(t *testing.T) {
	c := NewCanvas(30, 4)
	l := &List{Name: services, Marker: "> ", Blank: "  ", Focused: true, NoStatus: true}
	rows := []Row{{Lead: "M ", Text: "changed"}, {Lead: "A ", Text: "added"}}
	l.Draw(c, c.Bounds(), rows)

	lines := strings.Split(c.String(), "\n")
	if !strings.HasPrefix(lines[0], "> M changed") {
		t.Errorf("the cursor row is %q", lines[0])
	}
	if !strings.HasPrefix(lines[1], "  A added") {
		t.Errorf("the other row is %q", lines[1])
	}
}

// The marker sits OUTSIDE the indent, so the cursor is in the same column on
// every row. Indented, it moves with the nesting and a tree cannot be scanned.
func TestTheMarkerIsNotIndentedWithTheRow(t *testing.T) {
	c := NewCanvas(30, 4)
	l := &List{Name: services, Marker: "> ", Blank: "  ", Focused: true, NoStatus: true}
	l.Move(1)
	l.Draw(c, c.Bounds(), []Row{{Text: "top"}, {Depth: 2, Text: "deep"}})

	lines := strings.Split(c.String(), "\n")
	if !strings.HasPrefix(lines[1], "> ") {
		t.Errorf("a nested cursor row starts %q, so the marker moved with the indent", lines[1][:4])
	}
}

// columnOf is which COLUMN a substring starts at. strings.Index counts bytes,
// and a box-drawing glyph is three of them — the miscount comp.Width exists to
// prevent, and one this test made on its first attempt.
func columnOf(t *testing.T, line, want string) int {
	t.Helper()
	i := strings.Index(line, want)
	if i < 0 {
		t.Fatalf("%q is not in %q", want, line)
	}
	return Width(line[:i])
}

// Issue 66: a status column and a tree indent at the same time.
//
// lazygit's file rows lead with git status letters, which have to form a
// column down the screen. dive's layer tree indents. Before this, one Lead
// field had to serve both, and Depth pushed the status letters out of line.
func TestAStatusColumnSurvivesTheIndent(t *testing.T) {
	l := &List{Name: "files", StatusWidth: 2, NoStatus: true}
	rows := []Row{
		{Status: "M", Text: "main.go"},
		{Status: "A", Depth: 1, Text: "nested.go"},
		{Status: "D", Depth: 2, Text: "deep.go"},
	}
	c := NewCanvas(40, 4)
	l.Draw(c, c.Bounds(), rows)
	lines := strings.Split(uncolored.ReplaceAllString(c.String(), ""), "\n")

	// Every status letter in the same column, whatever the depth.
	col := -1
	for i, want := range []string{"M", "A", "D"} {
		at := strings.Index(lines[i], want)
		if at < 0 {
			t.Fatalf("row %d has no %q: %q", i, want, lines[i])
		}
		if col == -1 {
			col = at
		} else if at != col {
			t.Errorf("row %d has its status at column %d, want %d — the indent moved it", i, at, col)
		}
	}

	// And the text still indents.
	if strings.Index(lines[1], "nested.go") <= strings.Index(lines[0], "main.go") {
		t.Error("the depth-1 row did not indent")
	}
}

// Issue 76: the columns before a row's text, in columns.
func TestLeadWidthIsColumnsAndStatusRowsIsRows(t *testing.T) {
	l := &List{Marker: "> ", StatusWidth: 2}
	if got, want := l.LeadWidth(), 4; got != want {
		t.Errorf("LeadWidth = %d, want %d (a 2-column marker plus a 2-column status)", got, want)
	}
	if got, want := l.StatusRows(), 1; got != want {
		t.Errorf("StatusRows = %d, want %d", got, want)
	}
	// The old name still answers, so the private tools keep building.
	if l.Overhead() != l.StatusRows() {
		t.Error("Overhead stopped agreeing with StatusRows")
	}
	// A list with neither is zero, not one.
	if got := (&List{}).LeadWidth(); got != 0 {
		t.Errorf("a list with no marker and no status has LeadWidth %d, want 0", got)
	}
}

// Issue 78 through the list: branches go in Row.Indent and replace the spaces.
func TestRowIndentReplacesTheDepthSpaces(t *testing.T) {
	l := &List{Name: "tree", NoStatus: true}
	nodes := []Node{{Depth: 0, Key: "a"}, {Depth: 1, Key: "b"}, {Depth: 1, Key: "c"}}
	prefixes := Branches(nodes, theme.DefaultChrome)

	rows := make([]Row, len(nodes))
	for i, n := range nodes {
		rows[i] = Row{Indent: prefixes[i], Text: n.Key}
	}
	c := NewCanvas(40, 4)
	l.Draw(c, c.Bounds(), rows)
	frame := uncolored.ReplaceAllString(c.String(), "")

	if !strings.Contains(frame, "├─b") {
		t.Errorf("no branch before the middle child:\n%s", frame)
	}
	if !strings.Contains(frame, "└─c") {
		t.Errorf("no closing branch before the last child:\n%s", frame)
	}
}

// Issue 80: the cursor follows the row, not the line.
//
// dive reported this twice, five years apart: filter to something, move onto a
// row, clear the filter, and the cursor is on a different file.
func TestTheCursorFollowsItsKeyThroughAFilter(t *testing.T) {
	all := []Row{
		{Key: "a", Text: "alpha"},
		{Key: "b", Text: "bravo"},
		{Key: "c", Text: "charlie"},
		{Key: "d", Text: "delta"},
	}
	filtered := []Row{all[2], all[3]} // charlie, delta

	l := &List{Name: "rows", Focused: true, NoStatus: true}
	c := NewCanvas(20, 4)

	// Filtered, cursor onto delta.
	l.Draw(c, c.Bounds(), filtered)
	l.Move(1)
	l.Draw(c, c.Bounds(), filtered)
	if got := filtered[l.Cursor()].Key; got != "d" {
		t.Fatalf("cursor on %q, want d", got)
	}

	// Filter cleared. Index 1 is now bravo; the key says delta.
	l.Draw(c, c.Bounds(), all)
	if got := all[l.Cursor()].Key; got != "d" {
		t.Errorf("cursor landed on %q after the filter cleared, want d", got)
	}
}

// Without a key the cursor is an index, which is correct and free for a list
// whose rows never move. Opting out has to keep working exactly as before.
func TestWithoutAKeyTheCursorIsStillAnIndex(t *testing.T) {
	l := &List{Name: "rows", Focused: true, NoStatus: true}
	c := NewCanvas(20, 4)

	l.Draw(c, c.Bounds(), rows(4))
	l.Move(2)
	l.Draw(c, c.Bounds(), rows(4))
	if got := l.Cursor(); got != 2 {
		t.Fatalf("cursor at %d, want 2", got)
	}
	// A shorter list clamps, the way it always did.
	l.Draw(c, c.Bounds(), rows(2))
	if got := l.Cursor(); got != 1 {
		t.Errorf("cursor at %d after the list shrank to 2, want 1", got)
	}
}

// The row the reader was on is deleted. Falling back to the index is the only
// sane answer, and it must still be clamped.
func TestADeletedKeyFallsBackToTheIndex(t *testing.T) {
	before := []Row{{Key: "a"}, {Key: "b"}, {Key: "c"}}
	after := []Row{{Key: "a"}, {Key: "c"}}

	l := &List{Name: "rows", Focused: true, NoStatus: true}
	c := NewCanvas(20, 4)
	l.Draw(c, c.Bounds(), before)
	l.Move(1)
	l.Draw(c, c.Bounds(), before)
	if got := before[l.Cursor()].Key; got != "b" {
		t.Fatalf("cursor on %q, want b", got)
	}

	l.Draw(c, c.Bounds(), after)
	if l.Cursor() >= len(after) {
		t.Errorf("cursor at %d with %d rows", l.Cursor(), len(after))
	}
	if got := after[l.Cursor()].Key; got != "c" {
		t.Errorf("cursor on %q after b was deleted, want c — its neighbor", got)
	}
}

// A click means that row, so it drops the remembered key. Otherwise the key
// from the old row pulls the cursor straight back on the next frame.
func TestSelectForgetsTheRememberedKey(t *testing.T) {
	all := []Row{{Key: "a"}, {Key: "b"}, {Key: "c"}}
	l := &List{Name: "rows", Focused: true, NoStatus: true}
	c := NewCanvas(20, 4)

	l.Draw(c, c.Bounds(), all)
	l.Move(2)
	l.Draw(c, c.Bounds(), all)

	l.Select(0)
	l.Draw(c, c.Bounds(), all)
	if got := all[l.Cursor()].Key; got != "a" {
		t.Errorf("after clicking row 0 the cursor is on %q, want a", got)
	}
}
