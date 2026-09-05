package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

const viewerName Name = "viewer"

func doc(lines ...string) []Line {
	out := make([]Line, len(lines))
	for i, s := range lines {
		out[i] = Line{Text: s}
	}
	return out
}

func numbered(n int) []Line {
	out := make([]Line, n)
	for i := range out {
		out[i] = Line{Text: "line " + itoa(i+1)}
	}
	return out
}

// A document opens at the top. This is the whole difference from a LogPane,
// which opens at the tail — and democtl's log view was filed as a bug for
// opening at the oldest line, so the two really are opposite defaults.
func TestADocumentOpensAtTheTop(t *testing.T) {
	v := &Viewer{Name: viewerName}
	c := NewCanvas(20, 5)
	v.Draw(c, c.Bounds(), numbered(100))

	if v.Offset() != 0 {
		t.Errorf("opened at line %d", v.Offset())
	}
	if !strings.Contains(c.String(), "line 1") {
		t.Errorf("the first line is not on screen:\n%s", c.String())
	}
}

// The zero Viewer draws, which is what makes it usable before a tool has
// decided on any of its styling.
func TestTheZeroViewerIsUsable(t *testing.T) {
	v := &Viewer{}
	c := NewCanvas(20, 4)
	v.Draw(c, c.Bounds(), doc("alpha", "beta"))

	if !strings.Contains(c.String(), "alpha") {
		t.Errorf("a zero viewer drew nothing:\n%s", c.String())
	}
	if _, _, ok := v.Range(); ok {
		t.Error("a fresh viewer reports a range")
	}
}

// The reason this is not a List. A list paints the cursor row in one colour on
// purpose; a document must not, because the line under the cursor is the line
// being read and its syntax colours are the content.
func TestTheCursorLineKeepsItsSpans(t *testing.T) {
	keyword := lipgloss.NewStyle().Foreground(lipgloss.Color("5"))
	cursor := lipgloss.NewStyle().Background(lipgloss.Color("8"))

	v := &Viewer{Name: viewerName, Focused: true, Selected: &cursor}
	c := NewCanvas(20, 3)
	v.Draw(c, c.Bounds(), []Line{{Spans: []Segment{
		{Text: "func", Style: &keyword}, {Text: " main"},
	}}})

	cell, ok := c.CellAt(0, 0)
	if !ok || cell.Style == nil {
		t.Fatal("the keyword has no style")
	}
	fg, _, _, _ := cell.Style.GetForeground().RGBA()
	if fg == 0 {
		t.Error("the cursor line lost its foreground — this is the List rule, and it is wrong here")
	}
	if _, _, _, a := cell.Style.GetBackground().RGBA(); a == 0 {
		t.Error("the cursor line did not take the cursor's background")
	}
}

// The cursor style goes under the spans, so a span that sets its OWN background
// keeps it. A diff's added-line green must survive the cursor sitting on it.
func TestASpansOwnBackgroundBeatsTheCursors(t *testing.T) {
	added := lipgloss.NewStyle().Background(lipgloss.Color("2"))
	cursor := lipgloss.NewStyle().Background(lipgloss.Color("8"))

	merged := under(&cursor, &added)
	if got := merged.GetBackground(); got != added.GetBackground() {
		t.Errorf("the span's background became %v", got)
	}
}

// A document is wider than its pane and a list is not. Truncating a diff line
// at the edge means the content cannot be read at all.
func TestTheViewScrollsSideways(t *testing.T) {
	v := &Viewer{Name: viewerName, NoStatus: true}
	c := NewCanvas(10, 2)
	long := "abcdefghijklmnopqrstuvwxyz"
	v.Draw(c, c.Bounds(), doc(long))

	v.ScrollX(6)
	c = NewCanvas(10, 2)
	v.Draw(c, c.Bounds(), doc(long))

	if !strings.Contains(c.String(), "ghij") {
		t.Errorf("scrolled to column %d and got:\n%s", v.Col(), c.String())
	}
}

// Sideways scrolling stops at the longest line there actually is, so a document
// of short lines cannot be pushed into empty space.
func TestSidewaysScrollingStopsAtTheLongestLine(t *testing.T) {
	v := &Viewer{Name: viewerName, NoStatus: true}
	c := NewCanvas(40, 3)
	v.Draw(c, c.Bounds(), doc("short", "tiny"))

	v.ScrollX(500)
	if v.Col() > 4 {
		t.Errorf("scrolled to column %d past a five-column document", v.Col())
	}
}

// A tab is expanded to the next STOP rather than to a fixed number of spaces.
// Code that nearly lines up is worse than code that does not.
func TestTabsGoToTheNextStop(t *testing.T) {
	got := expandTabs([]Segment{{Text: "ab\tc"}}, 4)
	if got[0].Text != "ab  c" {
		t.Errorf("ab<tab>c at a stop of 4 became %q, want %q", got[0].Text, "ab  c")
	}
}

// The column a tab starts at spans the SPANS, so a highlighted keyword followed
// by a tab indents from where it really is rather than from zero.
func TestATabKnowsWhichColumnItStartsAt(t *testing.T) {
	got := expandTabs([]Segment{{Text: "ab"}, {Text: "\tc"}}, 4)
	if got[1].Text != "  c" {
		t.Errorf("the second span became %q, want %q — it expanded from column 0", got[1].Text, "  c")
	}
}

// The gutter's width comes from the whole document, not from what is on screen.
// Otherwise it grows by a column as you scroll past line 999 and every line of
// text shifts under the reader.
func TestTheGutterDoesNotChangeWidthAsYouScroll(t *testing.T) {
	v := &Viewer{Name: viewerName, Numbers: true, NoStatus: true}
	lines := numbered(1000)

	c := NewCanvas(20, 3)
	v.Draw(c, c.Bounds(), lines)
	top := strings.Index(c.String(), "line 1")

	v.Scroll(500)
	c = NewCanvas(20, 3)
	v.Draw(c, c.Bounds(), lines)
	deep := strings.Index(c.String(), "line 5")

	if top != deep {
		t.Errorf("the text starts at column %d at the top and %d at line 500", top, deep)
	}
}

// A viewer with no cursor is an ordinary state — k9s's YAML view and
// termshark's scrollabletext are both this. The arrow keys should still work,
// which is the branch this saves every such tool.
func TestWithNoCursorMoveScrolls(t *testing.T) {
	v := &Viewer{Name: viewerName, NoCursor: true}
	c := NewCanvas(20, 4)
	v.Draw(c, c.Bounds(), numbered(100))

	v.Move(10)
	if v.Offset() != 10 {
		t.Errorf("moved in a cursorless viewer and the offset is %d", v.Offset())
	}
	if _, _, ok := v.Range(); ok {
		t.Error("a cursorless viewer reports a range")
	}
}

// Extending anchors at the cursor on the first press, so shift-down once
// selects two lines rather than none.
func TestViewerExtendAnchorsAtTheCursor(t *testing.T) {
	v := &Viewer{Name: viewerName}
	c := NewCanvas(20, 10)
	v.Draw(c, c.Bounds(), numbered(20))

	v.Goto(4)
	v.Draw(c, c.Bounds(), numbered(20))
	v.Extend(1)

	lo, hi, ok := v.Range()
	if !ok || lo != 4 || hi != 5 {
		t.Errorf("range is %d..%d (%v), want 4..5", lo, hi, ok)
	}
}

// Extending backwards past the anchor and out the other side keeps the anchor
// where it was — the case a pair of bounds gets wrong.
func TestARangeSurvivesTurningRound(t *testing.T) {
	v := &Viewer{Name: viewerName}
	c := NewCanvas(20, 10)
	v.Draw(c, c.Bounds(), numbered(20))

	v.Goto(5)
	v.Draw(c, c.Bounds(), numbered(20))
	v.Extend(2)  // 5..7
	v.Extend(-5) // cursor at 2, anchor still 5

	lo, hi, _ := v.Range()
	if lo != 2 || hi != 5 {
		t.Errorf("range is %d..%d, want 2..5", lo, hi)
	}
}

// An ordinary move drops the range, because a selection that survived an arrow
// key is one the reader cannot get rid of.
func TestAPlainMoveDropsTheRange(t *testing.T) {
	v := &Viewer{Name: viewerName}
	c := NewCanvas(20, 10)
	v.Draw(c, c.Bounds(), numbered(20))

	v.Extend(3)
	v.Move(1)
	if _, _, ok := v.Range(); ok {
		t.Error("the range survived a plain move")
	}
}

// Goto is where a search hit lands. The index may be past the end of a document
// the caller has not measured, and clamping is at draw time for the same reason
// as List's: the lines arrive at Draw.
func TestGotoPastTheEndLandsOnTheLastLine(t *testing.T) {
	v := &Viewer{Name: viewerName}
	c := NewCanvas(20, 5)
	v.Goto(9999)
	v.Draw(c, c.Bounds(), numbered(30))

	if v.Cursor() != 29 {
		t.Errorf("the cursor is on line %d of a 30-line document", v.Cursor())
	}
	if !strings.Contains(c.String(), "line 30") {
		t.Errorf("the last line was not revealed:\n%s", c.String())
	}
}

// Scrolling with the wheel does not drag the cursor, and does not snap back to
// it on the next frame.
func TestScrollingDoesNotSnapBackToTheCursor(t *testing.T) {
	v := &Viewer{Name: viewerName, Focused: true}
	c := NewCanvas(20, 5)
	v.Draw(c, c.Bounds(), numbered(100))

	v.Scroll(40)
	v.Draw(c, c.Bounds(), numbered(100))
	if v.Offset() != 40 {
		t.Errorf("the view snapped back to %d", v.Offset())
	}
	if v.Cursor() != 0 {
		t.Errorf("scrolling moved the cursor to %d", v.Cursor())
	}
}

// A document longer than the pane cannot be scrolled past its end.
func TestScrollingStopsAtTheBottom(t *testing.T) {
	v := &Viewer{Name: viewerName, NoStatus: true}
	c := NewCanvas(20, 5)
	v.Draw(c, c.Bounds(), numbered(12))

	v.Scroll(999)
	v.Draw(c, c.Bounds(), numbered(12))
	if v.Offset() != 7 {
		t.Errorf("offset %d, want 7 — twelve lines in a five-row pane", v.Offset())
	}
}

// Lines are asked for one at a time and only when visible, so a 40 MB log costs
// the height of the pane rather than its own length.
func TestOnlyVisibleLinesAreAskedFor(t *testing.T) {
	v := &Viewer{Name: viewerName}
	c := NewCanvas(20, 6)
	var asked int
	v.DrawFunc(c, c.Bounds(), 1_000_000, func(i int) Line {
		asked++
		return Line{Text: "line " + itoa(i)}
	})
	if asked > 6 {
		t.Errorf("asked for %d lines to draw six rows", asked)
	}
}

// An owner ID carries the index in the DOCUMENT, so a click after scrolling
// acts on the line the reader pointed at rather than the row it landed on.
func TestOwnershipIsTheLineNotTheRow(t *testing.T) {
	v := &Viewer{Name: viewerName, NoStatus: true}
	c := NewCanvas(20, 4)
	v.Draw(c, c.Bounds(), numbered(100))
	v.Scroll(30)
	c = NewCanvas(20, 4)
	v.Draw(c, c.Bounds(), numbered(100))

	if got := c.OwnerAt(2, 0); got != Region(viewerName).At(30) {
		t.Errorf("the top row is owned by %v, want line 30", got)
	}
}

// A click on the blank after a short line still lands on that line, because the
// cell decides ownership rather than the glyph.
func TestTheBlankAfterALineBelongsToIt(t *testing.T) {
	sel := lipgloss.NewStyle().Background(lipgloss.Color("8"))
	v := &Viewer{Name: viewerName, Focused: true, Selected: &sel, NoStatus: true}
	c := NewCanvas(20, 3)
	v.Draw(c, c.Bounds(), doc("ab", "cd"))

	if got := c.OwnerAt(15, 0); got != Region(viewerName).At(0) {
		t.Errorf("the blank past the text is owned by %v", got)
	}
}

// An empty document is an ordinary state, not an error.
func TestAnEmptyDocumentSaysSo(t *testing.T) {
	v := &Viewer{Name: viewerName, Empty: "nothing to show"}
	c := NewCanvas(20, 4)
	v.Draw(c, c.Bounds(), nil)

	if !strings.Contains(c.String(), "nothing to show") {
		t.Errorf("an empty document drew:\n%q", c.String())
	}
}

// A pane too small to hold anything does not panic, which is the state every
// component gets wrong first.
func TestNoRoomDrawsNothing(t *testing.T) {
	v := &Viewer{Name: viewerName, Numbers: true}
	c := NewCanvas(40, 10)
	v.Draw(c, Rect{X: 0, Y: 0, W: 0, H: 0}, numbered(50))
	v.Draw(c, Rect{X: 0, Y: 0, W: 2, H: 1}, numbered(50))
}

// A pane that changed size under the cursor brings it back into view: the view
// moved because the window did, not because the reader asked.
func TestAResizeRevealsTheCursor(t *testing.T) {
	v := &Viewer{Name: viewerName, Focused: true, NoStatus: true}
	tall := Rect{W: 20, H: 20}
	c := NewCanvas(20, 20)
	v.Draw(c, tall, numbered(100))
	v.Goto(15)
	v.Draw(c, tall, numbered(100))

	short := Rect{W: 20, H: 5}
	v.Draw(NewCanvas(20, 5), short, numbered(100))
	if v.Cursor() < v.Offset() || v.Cursor() >= v.Offset()+5 {
		t.Errorf("after the resize the cursor is at %d and the view at %d", v.Cursor(), v.Offset())
	}
}

// fromCol cuts inside a span when it has to, so a horizontal offset does not
// jump by the length of a syntax token.
func TestScrollingSidewaysCutsInsideASpan(t *testing.T) {
	got := fromCol([]Segment{{Text: "package"}, {Text: " main"}}, 3)
	if len(got) != 2 || got[0].Text != "kage" {
		t.Errorf("cutting three columns gave %+v", got)
	}
}

// A shorter line must not leave the tail of a longer one behind it. Every line
// gets shorter at once when the view scrolls sideways, which is where a fill
// that only happened for styled lines shows up as ghosts.
func TestAShorterLineDoesNotLeaveAGhost(t *testing.T) {
	v := &Viewer{Name: viewerName, NoStatus: true}
	r := Rect{W: 40, H: 3}
	c := NewCanvas(40, 3)
	v.Draw(c, r, doc("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))

	v.ScrollX(20)
	v.Draw(c, r, doc("aaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"))

	if got := strings.TrimRight(c.String(), " \n"); Width(got) != 10 {
		t.Errorf("after scrolling 20 columns into a 30-column line, the row is %q", got)
	}
}
