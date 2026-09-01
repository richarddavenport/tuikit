package comp

import "github.com/charmbracelet/lipgloss"

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

	// Styles. Selected is the cursor row when focused, Unfocused when not.
	Selected, Unfocused, Status, EmptyStyle *lipgloss.Style

	// The cursor and the offset are not settable, only movable. Two fields
	// because looking around and choosing are different operations — and
	// unexported because the invariant between them is the whole component. A
	// caller that can assign either can reintroduce the bug this exists to
	// prevent, and it will, because assigning looks harmless.
	cursor, offset int

	// reveal says the cursor MOVED and should be brought into view on the next
	// draw. A flag rather than doing it in Draw, because "keep the cursor
	// visible" applied on every frame means the wheel snaps back the instant
	// it scrolls past the selection — the same conflation as scrolling with the
	// cursor, wearing the other hat.
	reveal bool

	// shown and count are what the last frame actually drew, and what the
	// offset is clamped against.
	shown, count int
}

// Cursor is the selected row's index in the list.
func (l *List) Cursor() int { return l.cursor }

// Offset is the first row shown.
func (l *List) Offset() int { return l.offset }

// Row is one line of a list.
type Row struct {
	Text  string
	Style *lipgloss.Style
}

// Draw renders the rows into r.
//
// The bottom line of r is the status line — the count, and which way an
// off-screen selection went. It is reserved whether or not the list overflows,
// because a viewport that only looks like one when it is scrolling is a
// viewport you cannot tell from a short list.
func (l *List) Draw(c *Canvas, r Rect, rows []Row) {
	l.count = len(rows)
	body := r
	if r.H > 1 {
		body.H = r.H - 1
	}
	l.shown = body.H

	if len(rows) == 0 {
		if l.Empty != "" {
			l.fill(c, Rect{X: body.X, Y: body.Y, W: body.W, H: 1}, l.EmptyStyle, Region(l.Name))
			c.Text(body.X, body.Y, l.Empty, l.EmptyStyle, Region(l.Name))
		}
		l.status(c, r)
		return
	}

	l.cursor = clamp(l.cursor, 0, len(rows)-1)
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

	for i := l.offset; i < len(rows) && i-l.offset < body.H; i++ {
		y := body.Y + i - l.offset
		id := Region(l.Name).At(i)

		style := rows[i].Style
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
		c.Text(body.X, y, rows[i].Text, style, id)
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
		marker = "↑ selected above"
	case l.cursor >= l.offset+l.shown:
		marker = "↓ selected below"
	}
	if marker != "" {
		c.Text(r.X+1, y, marker, l.Status, id)
	}

	count := itoa(min(l.shown, l.count)) + "/" + itoa(l.count)
	c.Text(r.Right()-Width(count), y, count, l.Status, id)
}

// Max is the furthest the list can be scrolled, given what the last frame drew.
func (l *List) Max() int { return max(0, l.count-l.shown) }

// Scroll moves the viewport. It never moves the cursor, and it never asks for
// the cursor to be revealed.
func (l *List) Scroll(by int) { l.offset = clamp(l.offset+by, 0, l.Max()) }

// Move moves the cursor, and asks for it to be brought into view.
func (l *List) Move(by int) { l.Select(l.cursor + by) }

// Select puts the cursor on a row by its index in the LIST.
//
// Clamped at draw time rather than here, so a caller may point at row 40 of a
// list this component has not been shown yet — which is what happens whenever
// state is restored before the first frame.
func (l *List) Select(i int) { l.cursor, l.reveal = max(0, i), true }

// Reset puts the list back to the top, for when the rows underneath it have
// changed out from under the cursor — a filter, usually.
func (l *List) Reset() { l.cursor, l.offset, l.reveal = 0, 0, false }

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
