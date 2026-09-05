package comp

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Viewer is a document you scroll, look through, and act on ranges of.
//
// # Where this came from
//
// Seven of the fourteen tools surveyed in the rebuilds repository each built
// this, and none of them could get it from their toolkit: lazygit's main panel,
// gitui's diff.rs and syntax_text.rs, k9s's live_view.go, termshark's
// scrollabletext and fileviewer, dive's layer detail, gh-dash's PR body — and
// fx, which is nothing else at all.
//
// # Why not LogPane
//
// [LogPane] tails. It opens at the newest line, follows the end, and its count
// means how far from the bottom you have scrolled. A document opens at the TOP,
// does not follow, and its count means where you are. Adding a Follow field to
// one of them makes a component that is wrong in both directions.
//
// # Why not List
//
// Closer, and the difference is one rule. [List] paints the cursor row in a
// single style, deliberately: "a selected row is one colour whatever its spans
// say", because the selection is the reader's own mark and a row keeping its
// colours under it makes the cursor hard to find.
//
// That rule is right for a list of services and wrong for a document. The line
// under the cursor in a diff is the line you are READING, and dropping its
// syntax colours to show where the cursor is trades the content for the
// pointer. So the cursor style here goes UNDER the spans instead — see [under].
//
// The other difference is sideways. A list is as wide as its pane and a
// document is not: a diff line, a JSON string and a hex dump all run off the
// edge, and truncating them means the content cannot be read at all. [List] has
// no answer to that and does not need one.
//
// # What it is not
//
// Not a text editor. Decision 27 draws that line and lazygit itself respects
// it: it stages hunks and does not edit them. Showing a buffer and acting on
// ranges of it is in; being the place you type the buffer is out.
type Viewer struct {
	// Name is the region lines are drawn under: line i is owned by Name[i],
	// carrying its index in the DOCUMENT rather than the row it landed on.
	Name Name

	// Focused decides whether a cursor move pulls the view along, and whether
	// the cursor is drawn as the reader's or as a remembered position.
	Focused bool

	// Empty is drawn when there are no lines. A filter matching nothing, or a
	// file that is genuinely empty, are ordinary states.
	Empty string

	// Numbers draws the line number gutter.
	//
	// Width comes from the TOTAL number of lines, not from the ones on screen,
	// so the gutter does not change width as you scroll past line 999 and the
	// text does not shift under the reader.
	Numbers bool

	// NoCursor is a document with nothing selected in it — k9s's YAML view and
	// termshark's scrollabletext are both this. Scrolling still works; there is
	// simply nothing for a range to be about.
	NoCursor bool

	// NoStatus gives the status row back to the document. See [List.NoStatus]
	// for when that is the right trade.
	NoStatus bool

	// Tab is the tab stop, in columns. Zero means [DefaultTab].
	//
	// A viewer that draws a tab as one cell shows Go source indented by a
	// single column, and a Makefile as nonsense. Expansion is to the next stop
	// rather than a fixed number of spaces, which is the difference between
	// code that lines up and code that nearly does.
	Tab int

	// Styles. Selected is the cursor line, and it is applied UNDERNEATH the
	// line's own spans: set a background on it and the syntax colours survive.
	// Ranged is the rest of a range selection, falling back to Selected.
	Selected, Ranged, Number, Status, EmptyStyle *lipgloss.Style

	// The cursor, the vertical offset and the horizontal one. Not settable,
	// only movable, for the same reason as [List]: the invariant between them
	// is the component.
	cursor, offset, col int

	sel rangeSel

	// reveal says the cursor moved and should be brought into view. A flag
	// rather than doing it every frame, so the wheel does not snap back the
	// instant it scrolls past the cursor.
	reveal bool

	// shown, count and widest are what the last frame drew, and what the
	// offsets are clamped against.
	shown, count, widest int
	drawn                bool
}

// DefaultTab is the tab stop when [Viewer.Tab] is zero. Eight is the terminal's
// own default and what a diff from git already assumes.
const DefaultTab = 8

// Line is one line of a document.
//
// Text and Style are the common case, Spans the highlighted one, and Spans wins
// when it is set — the same arrangement as [Row], so a tool that has built rows
// for a [List] already knows this type.
type Line struct {
	Text  string
	Style *lipgloss.Style
	Spans []Segment
}

// Cursor is the selected line's index in the document.
func (v *Viewer) Cursor() int { return v.cursor }

// Offset is the first line shown.
func (v *Viewer) Offset() int { return v.offset }

// Col is the first column shown, when the view has been scrolled sideways.
func (v *Viewer) Col() int { return v.col }

// Max is the furthest the view can scroll down.
func (v *Viewer) Max() int { return max(0, v.count-v.shown) }

// Move moves the cursor, and asks for it to be brought into view.
//
// With NoCursor set there is no cursor to move, so this scrolls instead —
// which is what the arrow keys should do in a document that has no selection,
// and saves every such tool a branch.
func (v *Viewer) Move(by int) {
	if v.NoCursor {
		v.Scroll(by)
		return
	}
	v.cursor = max(0, v.cursor+by)
	v.reveal = true
	v.sel.clear()
}

// Goto puts the cursor on a line and brings it into view.
//
// The landing place for a search hit and for `:` line numbers, which fx and
// gitui both have. Clamped at draw time against the count, so a caller may hand
// over a match index without knowing how long the document turned out to be.
func (v *Viewer) Goto(i int) {
	v.cursor = max(0, i)
	v.reveal = true
	v.sel.clear()
}

// Scroll moves the view without moving the cursor.
//
// The wheel, and page up and down. Unlike [LogPane.Scroll] there is no tail to
// re-attach to: a document has a bottom, and arriving at it means you are at
// the bottom rather than in a mode.
func (v *Viewer) Scroll(by int) {
	v.offset = clamp(v.offset+by, 0, v.Max())
}

// ScrollX moves the view sideways.
//
// Clamped against the widest line the last frame actually drew, so a document
// of short lines cannot be scrolled into empty space — and so the clamp costs
// nothing to compute, since the widths were measured to draw them anyway.
func (v *Viewer) ScrollX(by int) {
	v.col = clamp(v.col+by, 0, max(0, v.widest-1))
}

// Extend moves the cursor and keeps the other end of the selection where it is.
//
// The shift-arrow half of a range, and the reason a Viewer has a cursor at all:
// staging a hunk, copying a span of a log, and yanking a subtree are the same
// gesture. What the range commits INTO is the tool's, and is the open question
// in issue 55.
func (v *Viewer) Extend(by int) {
	if v.NoCursor {
		return
	}
	v.sel.start(v.cursor)
	v.cursor = max(0, v.cursor+by)
	v.reveal = true
}

// Range is the selected span, low to high, and whether there is one.
func (v *Viewer) Range() (lo, hi int, ok bool) {
	if v.NoCursor {
		return 0, 0, false
	}
	return v.sel.span(v.cursor)
}

// ClearRange drops the selection, leaving the cursor.
func (v *Viewer) ClearRange() { v.sel.clear() }

// Draw renders a document held as a slice.
func (v *Viewer) Draw(c *Canvas, r Rect, lines []Line) {
	v.DrawFunc(c, r, len(lines), func(i int) Line { return lines[i] })
}

// DrawFunc renders a document whose lines are produced on demand.
//
// line is called only for the lines actually on screen, so a frame costs the
// height of the pane rather than the length of the document — which is the
// difference between opening a 40 MB log and waiting for it. Same arrangement
// as [List.DrawFunc], and for the same reason.
func (v *Viewer) DrawFunc(c *Canvas, r Rect, n int, line func(i int) Line) {
	c = c.Clip(r)
	v.count = n

	body := r
	if r.H > 1 && !v.NoStatus {
		body.H = r.H - 1
	}
	if v.drawn && body.H != v.shown {
		v.reveal = true
	}
	v.shown, v.drawn = body.H, true

	if n == 0 {
		if v.Empty != "" {
			c.Text(body.X, body.Y, v.Empty, v.EmptyStyle, Region(v.Name))
		}
		v.status(c, r)
		return
	}

	v.cursor = clamp(v.cursor, 0, n-1)
	if !v.NoCursor && v.reveal {
		if v.cursor < v.offset {
			v.offset = v.cursor
		}
		if v.cursor >= v.offset+body.H {
			v.offset = v.cursor - body.H + 1
		}
	}
	v.reveal = false
	v.offset = clamp(v.offset, 0, v.Max())

	gutter := 0
	if v.Numbers {
		gutter = len(itoa(n)) + 1
	}
	text := body
	text.X, text.W = body.X+gutter, max(0, body.W-gutter)

	v.widest = 0
	lo, hi, ranged := v.Range()

	for i := v.offset; i < n && i-v.offset < body.H; i++ {
		y := body.Y + i - v.offset
		id := Region(v.Name).At(i)

		this := line(i)
		spans := this.Spans
		if len(spans) == 0 {
			spans = []Segment{{Text: this.Text, Style: this.Style}}
		}
		spans = expandTabs(spans, v.tab())
		if w := spansWidth(spans); w > v.widest {
			v.widest = w
		}

		// The line style goes under the spans rather than over them, so a
		// cursor with a background keeps the syntax colours it sits on.
		var beneath *lipgloss.Style
		switch {
		case v.NoCursor:
		case i == v.cursor && v.Focused:
			beneath = v.Selected
		case ranged && i >= lo && i <= hi:
			beneath = or2(v.Ranged, v.Selected)
		}

		if v.Numbers {
			num := Pad(itoa(i+1), gutter-1) + " "
			c.Text(body.X, y, num, or2(v.Number, beneath), id)
		}
		// Filled before the text, and ALWAYS — an unstyled fill still claims
		// the cells, which is what makes the blank after a short line belong to
		// that line and what stops a shorter line leaving the tail of a longer
		// one behind it. Scrolling sideways is where that shows: every line
		// gets shorter at once, and the ghosts are the previous frame's.
		c.Fill(Rect{X: text.X, Y: y, W: text.W, H: 1}, " ", beneath, id)

		x := text.X
		for _, span := range fromCol(spans, v.col) {
			if x >= text.Right()+1 {
				break
			}
			x += c.Text(x, y, Truncate(span.Text, text.Right()+1-x), under(beneath, span.Style), id)
		}
	}
	v.col = clamp(v.col, 0, max(0, v.widest-1))
	v.status(c, r)
}

// status says where in the document the view is, and how far sideways.
func (v *Viewer) status(c *Canvas, r Rect) {
	if r.H < 2 || v.NoStatus {
		return
	}
	// The last line on screen over the total — a POSITION rather than a
	// proportion, for the reason given at [List.status]: "5 of 20 shown" is the
	// same number wherever you scroll to.
	label := itoa(min(v.offset+v.shown, v.count)) + "/" + itoa(v.count)
	if v.col > 0 {
		label = "+" + itoa(v.col) + " " + label
	}
	c.Text(r.Right()-Width(label), r.Bottom(), label, v.Status, Region(v.Name))
}

func (v *Viewer) tab() int {
	if v.Tab > 0 {
		return v.Tab
	}
	return DefaultTab
}

// under puts base beneath s, so s keeps every value it set and takes the rest
// from base.
//
// This is the whole difference between a document and a list. lipgloss's
// Inherit copies only what the receiver has NOT set, which is exactly the
// question being asked: the cursor supplies a background, the span keeps its
// foreground, and the line under the cursor is still readable code.
func under(base, s *lipgloss.Style) *lipgloss.Style {
	if base == nil {
		return s
	}
	if s == nil {
		return base
	}
	merged := s.Inherit(*base)
	return &merged
}

// spansWidth is the columns a line occupies.
func spansWidth(spans []Segment) int {
	var w int
	for _, s := range spans {
		w += Width(s.Text)
	}
	return w
}

// fromCol drops col columns from the left of a line.
//
// Cutting BETWEEN spans where it can and inside one where it must, because a
// horizontal offset that could only land on a span boundary would jump by the
// length of a syntax token.
func fromCol(spans []Segment, col int) []Segment {
	if col <= 0 {
		return spans
	}
	out := make([]Segment, 0, len(spans))
	var x int
	for _, s := range spans {
		w := Width(s.Text)
		switch {
		case x+w <= col: // entirely to the left
		case x >= col:
			out = append(out, s)
		default:
			out = append(out, Segment{Text: ansi.TruncateLeft(s.Text, col-x, ""), Style: s.Style})
		}
		x += w
	}
	return out
}

// expandTabs replaces tabs with spaces to the next stop.
//
// Across the whole line rather than per span, because a tab's width depends on
// the column it starts at and a span does not know where it begins. A viewer
// that expanded each span from zero would indent the second token of every
// line wrongly.
func expandTabs(spans []Segment, tab int) []Segment {
	if tab < 1 {
		tab = DefaultTab
	}
	var any bool
	for _, s := range spans {
		if strings.ContainsRune(s.Text, '\t') {
			any = true
			break
		}
	}
	if !any {
		return spans
	}

	out := make([]Segment, 0, len(spans))
	var col int
	for _, s := range spans {
		var b strings.Builder
		for _, r := range s.Text {
			if r == '\t' {
				n := tab - col%tab
				b.WriteString(strings.Repeat(" ", n))
				col += n
				continue
			}
			b.WriteRune(r)
			col += Width(string(r))
		}
		out = append(out, Segment{Text: b.String(), Style: s.Style})
	}
	return out
}
