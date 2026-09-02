package comp

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// List is a scrollable, selectable list of rows.
//
// # Where this came from
//
// pgctl's `window` (view.go:232), swarmctl's pane scrolling (pane.go:409) and
// azctl's cursor/top pair (dashboard.go) are the same function written three
// times. All three keep the offset in a separate field from the cursor, and all
// three pull the offset along to keep the cursor in view.
//
// The one place two of them differ meaningfully is worth a parameter rather
// than a choice: swarmctl follows the cursor ONLY WHEN THE PANE IS FOCUSED, and
// pgctl always does. swarmctl is right. An unfocused pane whose viewport jumps
// because its cursor is somewhere else is a pane that moves while you are
// reading it, and the cursor there is a memory of where you were rather than a
// thing you are pointing at. So Focused decides it.
//
// What none of the three had, and the canvas-mouse prototype found by being
// used:
//
//   - The wheel moves the OFFSET and never the cursor. Wiring it to the cursor
//     is the obvious shortcut when a list has no viewport, and it makes
//     scrolling appear to select items at random.
//   - The offset clamps against what was actually drawn last frame rather than
//     a constant. Clamping to a constant is how a five-line pane scrolls into
//     empty space and blanks itself.
//   - The count is drawn ALWAYS, not only when the list overflows. On a tall
//     terminal where everything fits, a wheel that correctly does nothing is
//     indistinguishable from a wheel that is broken; `7/7` answers that in a
//     glance and an absent indicator answers nothing.
//   - A selection scrolled out of view says which way it went. Not by dragging
//     the cursor back into the viewport, which is the same conflation again: a
//     selection that is merely invisible leaves the detail pane describing
//     something nothing on screen points at, and the next key press acts on
//     something the reader cannot see.
type List struct {
	// Name is the region rows are drawn under: row i is owned by Name[i],
	// carrying its index in the LIST, not the row it landed on. After a scroll
	// those differ, and an ID that means "row 3 of the screen" acts on
	// whatever moved into row 3 rather than failing visibly.
	Name Name

	// Focused decides whether a cursor MOVE pulls the viewport along.
	Focused bool

	// Empty is drawn when there are no rows at all. A filter that matches
	// nothing is an ordinary state, not an error.
	Empty string

	// Marker is drawn against the cursor row and Blank against the rest, for a
	// list whose selection is a CHARACTER rather than only a colour. They
	// should be the same width, or the rows jump as you move.
	//
	// comp.Form has had these since it was written and a List did not, which
	// azctl's migration found the hard way: its resource rows are marked with
	// › and the port silently dropped them. Two components with a cursor
	// should agree about how a cursor is shown.
	Marker, Blank string

	// Styles. Selected is the cursor row when focused, Unfocused when not.
	Selected, Unfocused, Status, EmptyStyle *lipgloss.Style

	// The cursor and the offset are not settable, only movable. Two fields
	// because looking around and choosing are different operations — and
	// unexported because the invariant between them is the whole component. A
	// caller that can assign either can reintroduce the bug this exists to
	// prevent, and it will, because assigning looks harmless.
	cursor, offset int
	// pending is moves not yet resolved against the rows. See Move.
	pending int

	// reveal says the cursor MOVED and should be brought into view on the next
	// draw. A flag rather than doing it in Draw, because "keep the cursor
	// visible" applied on every frame means the wheel snaps back the instant
	// it scrolls past the selection — the same conflation as scrolling with the
	// cursor, wearing the other hat.
	reveal bool

	// shown and count are what the last frame actually drew, and what the
	// offset is clamped against.
	shown, count int
	// drawn says a frame has been drawn, so the first one is not mistaken for
	// a resize.
	drawn bool
}

// Cursor is the selected row's index in the list.
func (l *List) Cursor() int { return l.cursor }

// Offset is the first row shown.
func (l *List) Offset() int { return l.offset }

// Row is one line of a list.
//
// Text and Style are the common case: a whole row in one colour. Spans is for a
// row that is more than one — a name with a dim count after it, a timestamp
// then a message — and wins when it is set.
//
// Spans arrived late, from two tools independently. democtl's log lines are a
// muted timestamp then plain text, and azctl's tree rows are a label then a
// dim count; both had to fall back to drawing themselves rather than using a
// List. Two of four needing it is the line at which it stops being a special
// case.
type Row struct {
	Text  string
	Style *lipgloss.Style
	Spans []Segment

	// Skip makes a row one the cursor passes over: a heading, a blank line
	// between groups, a rule.
	//
	// The zero value is selectable, so a list that has never heard of this
	// behaves exactly as it did.
	//
	// Without it, a list containing a heading breaks three ways at once: ↑↓
	// appears to do nothing, whatever sits beside the list has nothing to show,
	// and enter acts on a row that is not a thing. swarmctl hit this and
	// hand-rolled `navigable()`; the command palette needs the same for its
	// tier headers.
	Skip bool

	// Depth indents the row, in levels — the flattened-hierarchy shape that
	// two of the four tools already have in five separate places, every one of
	// them writing `"  " + line` by hand.
	//
	// A NUMBER, not a node. The flattening stays the tool's: azctl's
	// row{bucket, res} and swarmctl's diffRow{service, action, change} look
	// alike and are not, because each carries a domain payload, and a shared
	// []Node would make both of them box their data or keep it twice. What
	// they have in common is that a child sits two columns right of its
	// header, and that is this field.
	Depth int

	// Lead is the row's own prefix, drawn in the cursor marker's column and
	// INSTEAD of it.
	//
	// For a group header that says something there already — ▸ collapsed, ▾
	// expanded. A row carrying its own state glyph does not also want the
	// cursor's mark on top of it: the highlight is what says where the cursor
	// is, and two glyphs fighting for one column is how a tree ends up with
	// its headers a character out of line with its children.
	Lead string
}

// Draw renders the rows into r.
//
// The bottom line of r is the status line — the count, and which way an
// off-screen selection went. It is reserved whether or not the list overflows,
// because a viewport that only looks like one when it is scrolling is a
// viewport you cannot tell from a short list.
func (l *List) Draw(c *Canvas, r Rect, rows []Row) {
	l.DrawFunc(c, r, len(rows), func(i int) Row { return rows[i] })
}

// DrawFunc is Draw for a list whose rows are made on demand.
//
// From the file managers and log viewers this library keeps being compared to,
// where a directory of 200,000 entries or a log of a million lines is ordinary.
// Draw already only PAINTS what fits — the expensive part was never the
// drawing, it was building a []Row for everything so that twenty of them could
// be shown.
//
// row is called only for the rows actually on screen, so the cost of a frame
// is the size of the pane rather than the size of the data. It is called with
// indices in [0, n), and it must be cheap: it runs on every frame, for every
// visible row.
//
// The obvious alternative — a Rows interface with Len and At — was not taken
// because a func is what every call site already has, and an interface would
// make the common case (a slice you already built) into a wrapper type.
func (l *List) DrawFunc(c *Canvas, r Rect, n int, row func(i int) Row) {
	c = c.Clip(r)
	l.count = n
	body := r
	if r.H > 1 {
		body.H = r.H - 1
	}
	// A pane that changed size under the cursor reveals it again. The reveal
	// flag exists so the WHEEL does not snap back, and a resize is not a
	// wheel: the view moved underneath the reader rather than because they
	// asked. Leaving the cursor stranded there means the detail beside it
	// describes something off screen, and the next key acts on it.
	if l.drawn && body.H != l.shown {
		l.reveal = true
	}
	l.shown, l.drawn = body.H, true

	if n == 0 {
		if l.Empty != "" {
			l.fill(c, Rect{X: body.X, Y: body.Y, W: body.W, H: 1}, l.EmptyStyle, Region(l.Name))
			c.Text(body.X, body.Y, l.Empty, l.EmptyStyle, Region(l.Name))
		}
		l.status(c, r)
		return
	}

	l.resolve(n, row)
	if l.Focused && l.reveal {
		// Far enough to see the cursor, and no further.
		if l.cursor < l.offset {
			l.offset = l.cursor
		}
		if l.cursor >= l.offset+body.H {
			l.offset = l.cursor - body.H + 1
		}
	}
	l.reveal = false
	l.offset = clamp(l.offset, 0, l.Max())

	for i := l.offset; i < n && i-l.offset < body.H; i++ {
		y := body.Y + i - l.offset
		id := Region(l.Name).At(i)

		// Asked for once per visible row per frame, and never for a row that
		// is off screen — which is the whole point.
		this := row(i)
		style := this.Style
		if i == l.cursor {
			style = l.Unfocused
			if l.Focused {
				style = l.Selected
			}
		}
		// The row is filled before the text is drawn, so clicking the blank
		// after a short name still selects it: the cell decides ownership, not
		// the glyph.
		l.fill(c, Rect{X: body.X, Y: y, W: body.W, H: 1}, style, id)

		// A selected row is one colour whatever its spans say. The selection is
		// the reader's own mark on the list, and a row that kept its own
		// colours under it would make the cursor hard to find in exactly the
		// list where finding it matters.
		lead := l.lead(c, this, i)
		if len(this.Spans) == 0 || (i == l.cursor && style != nil) {
			text := this.Text
			if text == "" {
				text = spansText(this.Spans)
			}
			c.Text(body.X, y, lead+text, style, id)
			continue
		}
		x := body.X + c.Text(body.X, y, lead, style, id)
		for _, span := range this.Spans {
			x += c.Text(x, y, span.Text, span.Style, id)
		}
	}
	l.status(c, r)
}

func (l *List) fill(c *Canvas, r Rect, s *lipgloss.Style, id ID) {
	if s == nil {
		// An unstyled fill would still claim the cells, which is what makes
		// the whole row clickable.
		c.Fill(r, " ", nil, id)
		return
	}
	c.Fill(r, " ", s, id)
}

// status draws the count, and which way an off-screen selection went.
func (l *List) status(c *Canvas, r Rect) {
	if r.H < 2 {
		return
	}
	y := r.Bottom()
	id := Region(l.Name)

	var marker string
	switch {
	case l.count == 0:
	case l.cursor < l.offset:
		marker = c.Chrome().ScrollUp + " selected above"
	case l.cursor >= l.offset+l.shown:
		marker = c.Chrome().ScrollDown + " selected below"
	}
	if marker != "" {
		c.Text(r.X+1, y, marker, l.Status, id)
	}

	// The last row on screen, over the total — a POSITION, not a proportion.
	// "5 of 20 shown" is the same number wherever you scroll to, so it says
	// nothing about where you are; "8/20" moves under your hand, which is what
	// makes a viewport feel like one.
	count := itoa(min(l.offset+l.shown, l.count)) + "/" + itoa(l.count)
	c.Text(r.Right()-Width(count), y, count, l.Status, id)
}

// Count is how many rows the last frame was given, and Shown how many of them
// fit. Both, because Max is their difference and a caller cannot get back to
// the pair from it — which is what a [Scrollbar] needs to size its thumb and
// place it.
func (l *List) Count() int { return l.count }

// Shown is how many rows the last frame had room for.
func (l *List) Shown() int { return l.shown }

// Max is the furthest the list can be scrolled, given what the last frame drew.
func (l *List) Max() int { return max(0, l.count-l.shown) }

// Scroll moves the viewport. It never moves the cursor, and it never asks for
// the cursor to be revealed.
func (l *List) Scroll(by int) { l.offset = clamp(l.offset+by, 0, l.Max()) }

// Move moves the cursor, and asks for it to be brought into view.
//
// The move is RECORDED and resolved at draw time, because the list does not
// hold its rows — they arrive at Draw — so nothing here can know which of them
// the cursor is allowed to land on. That is the same arrangement Select and the
// viewport already use: "clamped at draw time rather than here".
//
// It means Cursor() between a Move and a Draw is the old value. Every caller in
// this repo and in the tools moves in Update and reads in Draw, which is the
// order a Bubble Tea program runs in anyway.
func (l *List) Move(by int) { l.pending += by; l.reveal = true }

// Select puts the cursor on a row by its index in the LIST.
//
// Immediate, unlike Move: a caller that just clicked row 3 means row 3, and
// asks about it in the same breath. If that row turns out to be one the cursor
// may not hold, the draw moves off it — a click on a heading does nothing
// rather than selecting its neighbour, because a cursor that lands somewhere
// you did not click is worse than a click that is ignored.
func (l *List) Select(i int) { l.cursor, l.pending, l.reveal = max(0, i), 0, true }

// Reset puts the list back to the top, for when the rows underneath it have
// changed out from under the cursor — a filter, usually.
func (l *List) Reset() { l.cursor, l.offset, l.pending, l.reveal = 0, 0, 0, false }

// resolve settles the cursor against the rows that actually exist.
//
// Three jobs, in order, and the order matters: apply the moves that were
// recorded since the last frame, keep the cursor inside the list, and get it
// off a row it may not hold.
//
// row is asked about individual indices rather than handed a slice, so this
// works for DrawFunc's list of two hundred thousand as well: with no skipped
// rows it costs one call per step.
func (l *List) resolve(n int, row func(int) Row) {
	if n == 0 {
		l.cursor, l.pending = 0, 0
		return
	}
	cur := clamp(l.cursor, 0, n-1)
	step, dir := l.pending, 1
	l.pending = 0
	if step < 0 {
		dir = -1
	}

	for ; step != 0; step -= dir {
		next, ok := l.step(n, row, cur, dir)
		if !ok {
			// Off the end. Stay where we are rather than sticking on a
			// trailing heading, and drop the rest of the move: pressing ↓ ten
			// times at the bottom should not queue ten moves back up.
			break
		}
		cur = next
	}

	// A click, a restored cursor, or rows that changed underneath can all leave
	// it on a row it may not hold.
	if row(cur).Skip {
		if i, ok := l.nearest(n, row, cur); ok {
			cur = i
		}
	}
	l.cursor = cur
}

// step is the next selectable row in one direction, or false at the end.
func (l *List) step(n int, row func(int) Row, from, dir int) (int, bool) {
	for i := from + dir; i >= 0 && i < n; i += dir {
		if !row(i).Skip {
			return i, true
		}
	}
	return 0, false
}

// nearest is the closest selectable row to i, forwards first.
//
// Forwards first because a list usually opens with a heading, and the row a
// reader means is the one under it. Returns false when every row is skipped —
// a list of nothing but headings, which is a legitimate empty-ish state and
// must not spin looking for a cursor that cannot exist.
func (l *List) nearest(n int, row func(int) Row, i int) (int, bool) {
	if next, ok := l.step(n, row, i, 1); ok {
		return next, true
	}
	return l.step(n, row, i, -1)
}

// lead is everything drawn before a row's text: its indent, then either its own
// prefix or the cursor's mark.
func (l *List) lead(c *Canvas, row Row, i int) string {
	indent := strings.Repeat(" ", max(0, row.Depth)*c.Chrome().Indent)
	if row.Lead != "" {
		return indent + row.Lead
	}
	return indent + l.mark(i)
}

// mark is the marker for a row, or the blank that keeps the others in line.
func (l *List) mark(i int) string {
	if l.Marker == "" {
		return ""
	}
	if i == l.cursor {
		return l.Marker
	}
	if l.Blank != "" {
		return l.Blank
	}
	return strings.Repeat(" ", Width(l.Marker))
}

// spansText is a row's words without its colours, for when the selection paints
// over them.
func spansText(spans []Segment) string {
	var b strings.Builder
	for _, s := range spans {
		b.WriteString(s.Text)
	}
	return b.String()
}

func clamp(v, lo, hi int) int { return max(lo, min(v, hi)) }

// itoa is small enough to keep the package free of strconv in the draw path.
func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var d [20]byte
	i := len(d)
	for n > 0 {
		i--
		d[i] = byte('0' + n%10)
		n /= 10
	}
	return string(d[i:])
}
