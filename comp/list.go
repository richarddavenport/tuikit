package comp

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// List is a scrollable, selectable list of rows.
//
// # Where this came from
//
// The database tool's `window`, the deploy tool's pane scrolling and the cloud
// tool's cursor/top pair are the same function written three times. All three
// keep the offset in a separate field from the cursor, and all three pull the
// offset along to keep the cursor in view.
//
// The one place two of them differ meaningfully is worth a parameter rather
// than a choice: the deploy tool follows the cursor ONLY WHEN THE PANE IS
// FOCUSED, and the database tool always does. The deploy tool is right. An
// unfocused pane whose viewport jumps because its cursor is somewhere else is
// a pane that moves while you are reading it, and the cursor there is a memory
// of where you were rather than a thing you are pointing at. So Focused
// decides it.
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

	// StatusWidth is the columns reserved for [Row.Status], and zero means no
	// status column at all.
	//
	// Declared on the list rather than measured from the rows, because
	// DrawFunc is handed one row at a time and never sees the widest. A number
	// the caller states is a column the same width on every frame, including
	// the frame where the only row with a two-column glyph has scrolled off.
	StatusWidth int

	// Marker is drawn against the cursor row and Blank against the rest, for a
	// list whose selection is a CHARACTER rather than only a colour. They
	// should be the same width, or the rows jump as you move.
	//
	// comp.Form has had these since it was written and a List did not, which the
	// cloud tool's migration found the hard way: its resource rows are marked
	// with › and the port silently dropped them. Two components with a cursor
	// should agree about how a cursor is shown.
	Marker, Blank string

	// Styles. Selected is the cursor row when focused, Unfocused when not.
	Selected, Unfocused, Status, EmptyStyle *lipgloss.Style

	// Ranged draws the rows in a range selection that are not the cursor.
	// Nil falls back to Unfocused, so a list that sets no range style still
	// shows one rather than nothing.
	Ranged *lipgloss.Style

	// NoStatus gives the status row back to the rows.
	//
	// The default reserves it whether or not the list overflows, because a
	// viewport that only looks like one when it is scrolling is a viewport you
	// cannot tell from a short list. That reasoning holds for one big list and
	// stops holding for several small ones: the database tool stacks FIVE lists
	// in a column, and at 80x24 their status rows are five of about twenty-one
	// body rows — a quarter of the column spent on counters reading 3/3, 3/3,
	// 1/1, 1/1 and blank, next to panel titles that already say the same number.
	//
	// Off by default, so a list that has never heard of this keeps the row and
	// the guarantee that comes with it. Turn it off only where something else
	// on screen says how much there is.
	NoStatus bool

	// The cursor and the offset are not settable, only movable. Two fields
	// because looking around and choosing are different operations — and
	// unexported because the invariant between them is the whole component. A
	// caller that can assign either can reintroduce the bug this exists to
	// prevent, and it will, because assigning looks harmless.
	cursor, offset int

	// key is the identity of the row the cursor was on at the last draw, when
	// the rows carry one. Matched against the incoming rows to find where the
	// cursor belongs now.
	key string
	// sel is the range selection's anchor. Shared with [Viewer], because the
	// invariant is easy to get wrong the same way twice.
	sel Range
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
// Text and Style are the common case: a whole row in one colour. Spans is for
// a row that is more than one — a name with a dim count after it, a timestamp
// then a message — and wins when it is set.
//
// Spans arrived late, from two tools independently. democtl's log lines are a
// muted timestamp then plain text, and the cloud tool's tree rows are a label
// then a dim count; both had to fall back to drawing themselves rather than
// using a List. Two of four needing it is the line at which it stops being a
// special case.
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
	// and enter acts on a row that is not a thing. The deploy tool hit this and
	// hand-rolled `navigable()`; the command palette needs the same for its
	// tier headers.
	Skip bool

	// Depth indents the row, in levels — the flattened-hierarchy shape that
	// two of the four tools already have in five separate places, every one of
	// them writing `"  " + line` by hand.
	//
	// A NUMBER, not a node. The flattening stays the tool's: the cloud tool's
	// row{bucket, res} and the deploy tool's diffRow{service, action, change}
	// look alike and are not, because each carries a domain payload, and a shared
	// []Node would make both of them box their data or keep it twice. What they
	// have in common is that a child sits two columns right of its header, and
	// that is this field.
	Depth int

	// Key is this row's identity, for a list whose rows come and go.
	//
	// Set it and the cursor follows the ROW rather than the line. A filter
	// applied, a filter cleared, a view toggled, a refresh that reorders — the
	// cursor lands back on the thing the reader was looking at, because it
	// remembers what that was rather than where it sat.
	//
	// Empty, and the cursor is an index, which is correct and free for a list
	// whose rows never move. Opt-in for that reason: a tool that would have to
	// invent a key per row per frame to get behaviour it already had should
	// not have to.
	//
	// Everything else in comp is keyed by identity — owner IDs, [Marks],
	// [Tree.Collapsed] — and their docs all say why: indices move. The cursor
	// was the last thing that did not, and dive reported the consequence twice
	// five years apart (issue 80).
	//
	// When the key is gone from the list, the cursor stays at its index and is
	// clamped, which is the old behaviour and the only sane fallback: the row
	// a reader was on has been deleted, and its neighbour is the best guess.
	Key string

	// Status is a fixed column between the cursor marker and the indent.
	//
	// For a glyph that means the same thing on every row whatever its depth: a
	// git status letter, a change mark, a reachability dot. It is padded to
	// [List.StatusWidth], so the glyphs line up down the screen and a header
	// can be positioned with [List.LeadWidth].
	//
	// The distinction from Lead is where the INDENT falls. Status is drawn
	// before it and Lead after it, so a status column stays a column while a
	// tree marker travels with its row. Both at once is what lazygit's file
	// list and dive's layer tree each need and neither could have: a row with
	// one Lead had to choose, and Depth pushed the chosen one out of line
	// (issue 66).
	Status string

	// StatusStyle colours Status and keeps that colour under the selection,
	// the same way LeadStyle does and for the same reason.
	StatusStyle *lipgloss.Style

	// Indent replaces the spaces Depth would have drawn.
	//
	// For a tree that draws BRANCHES rather than whitespace. [Branches] builds
	// these from the same []Node [Tree] already takes; anything else a tool
	// wants in that space works too.
	//
	// Empty falls back to Depth, so nothing that does not set it changes.
	Indent string

	// Lead is the row's own prefix, drawn after the indent and before the text.
	//
	// For a marker that belongs to the row rather than to the column — ▸
	// collapsed, ▾ expanded. It travels with the row as it indents, which is
	// what a tree wants and what a status column does not.
	Lead string

	// Right is drawn hard against the row's right-hand edge, after Spans.
	//
	// [Bar] has done this since it was written and was not reachable from
	// inside a list, so three call sites in the board did the arithmetic by hand
	// — and it produced two bugs that a golden caught and a reader would not
	// have: len(s) where Width(s) was meant, on a string holding a `,` and
	// forgetting that Lead comes out of the same width, which pushed every lane
	// count one column past the edge where it was silently not drawn (issue
	// 53).
	//
	// Both are the same mistake: the row knows its width and the caller does
	// not. So the row does it.
	//
	// Dropped rather than overlapped when there is not room for both, because a
	// count written over the end of a name is two pieces of information and
	// neither is readable.
	Right []Segment

	// LeadStyle draws Lead in the row's own colour, and keeps it there when the
	// row is selected.
	//
	// Nil means the lead takes whatever the rest of the row takes, which is
	// what every list did before this existed.
	//
	// It is here because a selected row is otherwise one colour whatever its
	// spans say, and that is right for a LABEL and wrong for a glyph that IS the
	// state. The database tool's connection list marks reachability with ● ○ ✗ in
	// the first column; on the cursor row all three came out bold black on white,
	// so the one row a reader is looking at was the one row whose status they
	// could not read. A person using it said so.
	//
	// It is also an accessibility rule and not only a legibility one. A black ●
	// on light grey does not read as "green ● that is highlighted", it reads as a
	// DIFFERENT state — off, disabled. The database tool was saved by using
	// distinct shapes as well as colours; a tool encoding state in colour alone
	// would have lost it outright, and nothing in the API would have said so.
	//
	// Only the lead, deliberately. Letting every styled span survive selection
	// is more elegant and makes the cursor's prominence depend on how colourful
	// a row happens to be — the selection would be strong on a plain list and
	// nearly invisible on a busy one, which is the opposite of what it is for.
	LeadStyle *lipgloss.Style
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
// where a directory of 200,000 entries or a log of a million lines is
// ordinary. Draw already only PAINTS what fits — the expensive part was never
// the drawing, it was building a []Row for everything so that twenty of them
// could be shown.
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
	if r.H > 1 && !l.NoStatus {
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
		if lo, hi, ok := l.Range(); ok && i >= lo && i <= hi && i != l.cursor {
			style = or2(l.Ranged, l.Unfocused)
		}
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
		// The lead keeps its own colour through the selection when it has one,
		// because it is the row's state and not the row's label.
		leadStyle := style
		if this.LeadStyle != nil {
			leadStyle = this.LeadStyle
		}

		if len(this.Spans) == 0 || (i == l.cursor && style != nil) {
			text := this.Text
			if text == "" {
				text = spansText(this.Spans)
			}
			x := body.X + l.drawLead(c, body.X, y, this, i, leadStyle, id)
			c.Text(x, y, text, style, id)
			l.right(c, body, y, this, style, id)
			continue
		}
		x := body.X + l.drawLead(c, body.X, y, this, i, leadStyle, id)
		for _, span := range this.Spans {
			x += c.Text(x, y, span.Text, span.Style, id)
		}
		l.right(c, body, y, this, style, id)
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

// right draws Row.Right against the row's right edge.
//
// A selected row's Right takes the selection style like everything else: it is
// part of the row rather than a mark on it, which is the distinction LeadStyle
// exists to make on the other side.
func (l *List) right(c *Canvas, body Rect, y int, row Row, style *lipgloss.Style, id ID) {
	if len(row.Right) == 0 {
		return
	}
	w := 0
	for _, span := range row.Right {
		w += Width(span.Text)
	}
	if w <= 0 || w > body.W {
		return
	}
	x := body.Right() - w + 1
	for _, span := range row.Right {
		st := span.Style
		if style != nil {
			st = style
		}
		x += c.Text(x, y, span.Text, st, id)
	}
}

// StatusRows is how many rows of a band the list spends on itself rather than
// on rows — 1 for the status row, 0 with NoStatus.
//
// ROWS. [List.LeadWidth] is the columns question, and the two were confusable
// enough under the old name that a rebuild used this one to indent a table
// header (issue 76).
//
// It exists so a tool laying out several lists does not encode this
// component's internals as a constant. The database tool had `const chrome =
// 3` (two borders and "the row comp.List keeps for its position counter"),
// which is a number that goes silently wrong the moment the answer changes —
// the class decision 32 is about.
func (l *List) StatusRows() int {
	if l.NoStatus {
		return 0
	}
	return 1
}

// Overhead is [List.StatusRows] under its old name.
//
// Renamed because "overhead" does not say ROWS, and a caller aligning a table
// header read it as columns, used it, and got a header indented by one — which
// compiled, drew, and was wrong (issue 76). [List.LeadWidth] is the columns
// answer.
//
// Deprecated: use [List.StatusRows] for rows or [List.LeadWidth] for columns.
func (l *List) Overhead() int { return l.StatusRows() }

// status draws the count, and which way an off-screen selection went.
func (l *List) status(c *Canvas, r Rect) {
	if r.H < 2 || l.NoStatus {
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

// Max is the furthest the list can be scrolled, given what the last frame
// drew.
func (l *List) Max() int { return max(0, l.count-l.shown) }

// Scroll moves the viewport. It never moves the cursor, and it never asks for
// the cursor to be revealed.
func (l *List) Scroll(by int) { l.offset = clamp(l.offset+by, 0, l.Max()) }

// Move moves the cursor, and asks for it to be brought into view.
//
// The move is RECORDED and resolved at draw time, because the list does not
// hold its rows — they arrive at Draw — so nothing here can know which of them
// the cursor is allowed to land on. That is the same arrangement Select and
// the viewport already use: "clamped at draw time rather than here".
//
// It means Cursor() between a Move and a Draw is the old value. Every caller
// in this repo and in the tools moves in Update and reads in Draw, which is
// the order a Bubble Tea program runs in anyway.
func (l *List) Move(by int) { l.pending += by; l.reveal = true; l.ClearRange() }

// Select puts the cursor on a row by its index in the LIST.
//
// Immediate, unlike Move: a caller that just clicked row 3 means row 3, and
// asks about it in the same breath. If that row turns out to be one the cursor
// may not hold, the draw moves off it — a click on a heading does nothing
// rather than selecting its neighbour, because a cursor that lands somewhere
// you did not click is worse than a click that is ignored.
//
// Selecting the row the cursor is ALREADY on keeps a pending Move, because it
// is not a selection — it is a caller saying the cursor is fine where it is.
// The distinction matters because "clamp every cursor whenever the data
// changes" is a pattern all four tools arrived at independently, back when a
// cursor was a plain int that could point past a list which had shrunk under
// it:
//
//	m.list.Select(clamp(m.list.Cursor(), n-1))
//
// With the cursor in range that clamp is `Select(Cursor())`, and it used to
// zero the pending move recorded by the arrow key in the same Update. The move
// was applied and immediately discarded, so the key did nothing and nothing
// errored (issue 45). Out of range it still selects for real, which is the
// case the clamp was written for.
//
// Not the wider fix of applying pending on top of any Select: a click means
// that row, and a queued arrow key landing on top of a click would be worse
// than the bug.
func (l *List) Select(i int) {
	if at := max(0, i); at != l.cursor {
		// The remembered key goes with it: a click means THAT row, and a key
		// from the row the cursor used to be on would pull it straight back on
		// the next frame.
		l.cursor, l.pending, l.key = at, 0, ""
	}
	l.reveal = true
	l.ClearRange()
}

// Extend moves the cursor and keeps the other end of the selection where it
// is.
//
// The shift-arrow half of a range. The first Extend from no selection anchors
// at the cursor, so a reader who presses shift-down once has selected two rows
// rather than none — which is what every other list in the world does.
//
// Deferred like [List.Move], and for the same reason: the rows a move lands on
// are not known until the frame that draws them.
func (l *List) Extend(by int) {
	l.sel.Start(l.cursor)
	l.pending += by
	l.reveal = true
}

// Range is the selected span, low to high, and whether there is one.
//
// Ordered, so a caller acting on it never has to ask which way the reader
// dragged. lazygit's patch_exploring carries 13 kB of state for this and most
// of it is keeping a range sane across a re-render; the cursor and one anchor
// are enough when both are clamped by the same resolve.
func (l *List) Range() (lo, hi int, ok bool) { return l.sel.Span(l.cursor) }

// ClearRange drops the selection, leaving the cursor.
//
// Called by [List.Move] and [List.Select], because moving without extending is
// how every list says "start again" — a range that survived an ordinary arrow
// key would be a range a reader cannot get rid of.
func (l *List) ClearRange() { l.sel.Clear() }

// or2 picks the first style that is set.
func or2(a, b *lipgloss.Style) *lipgloss.Style {
	if a != nil {
		return a
	}
	return b
}

// Reset puts the list back to the top, for when the rows underneath it have
// changed out from under the cursor — a filter, usually.
func (l *List) Reset() {
	l.cursor, l.offset, l.pending, l.reveal = 0, 0, 0, false
	l.key = ""
}

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
	cur := l.locate(n, row)
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
	settled := row(cur)
	if settled.Skip {
		if i, ok := l.nearest(n, row, cur); ok {
			cur = i
			settled = row(cur)
		}
	}
	l.cursor = cur
	// Remembered from the row already fetched, not a second call: a list of two
	// hundred thousand pays for the rows it draws and nothing else, which is
	// what DrawFunc is for.
	l.key = settled.Key
}

// locate is where the cursor belongs in the rows that just arrived.
//
// By KEY when the rows carry one and the remembered key is still present,
// otherwise by index, clamped. Those are the two halves of issue 80: the index
// keeps the cursor in range and the key keeps it on the same thing, and a list
// needs both because a row can be filtered out entirely.
//
// A linear scan, once per frame, over rows already being asked for. A map
// would cost a build of every key per frame to save a walk of the same length,
// and DrawFunc exists so that a list of two hundred thousand is not walked at
// all — so this stays O(n) and the tools that care about that do not set Key.
func (l *List) locate(n int, row func(int) Row) int {
	at := clamp(l.cursor, 0, n-1)
	if l.key == "" {
		return at
	}
	for i := 0; i < n; i++ {
		if row(i).Key == l.key {
			return i
		}
	}
	// Gone. Stay at the index and let the clamp above stand — the row the
	// reader was on has been removed, and its neighbour is the nearest thing
	// to what they were looking at.
	return at
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

// lead is everything drawn before a row's text: the cursor's mark, then the
// indent, then the row's own prefix.
//
// # The marker comes first, and outside the indent
//
// It used to be one or the other: a row with a Lead got no marker at all. That
// is wrong in the ordinary case rather than an edge one, because a status
// glyph in the lead column is exactly what Row.LeadStyle exists for. Such a
// list had its selection carried entirely by Selected's background, so it
// vanished in a pipe, in a golden, and for a reader who cannot see the colour
// — and harness.ShapeSurvivesColour could not catch it, because the shape was
// fine and the information was what went missing.
//
// Found by building a tool rather than by reading one: three lists in gcpeasy
// each re-implemented this badly, and lazygit's file rows need it too.
//
// Outside the indent because a nested row's marker still belongs in the
// cursor's column. Indenting it puts the cursor somewhere different on every
// row and makes a tree impossible to scan.
// drawLead paints the marker, the status column, the indent and the row's own
// prefix, and returns the columns they took.
//
// Four pieces rather than one string, because Status carries its own colour.
// Concatenating them would mean a git status letter and a fold marker sharing
// one style, which is the thing LeadStyle exists to prevent one level up.
func (l *List) drawLead(c *Canvas, x, y int, row Row, i int, leadStyle *lipgloss.Style, id ID) int {
	w := c.Text(x, y, l.mark(i), leadStyle, id)
	if status := l.statusOf(row); status != "" {
		style := leadStyle
		if row.StatusStyle != nil {
			style = row.StatusStyle
		}
		w += c.Text(x+w, y, status, style, id)
	}
	w += c.Text(x+w, y, l.indentOf(c, row), leadStyle, id)
	return w + c.Text(x+w, y, row.Lead, leadStyle, id)
}

// statusOf pads Status to StatusWidth, so the column is the same width on
// every row including the rows with nothing to put in it.
func (l *List) statusOf(row Row) string {
	if l.StatusWidth <= 0 {
		return ""
	}
	return Pad(Truncate(row.Status, l.StatusWidth), l.StatusWidth)
}

// indentOf is the row's own indent string, or the spaces its Depth asks for.
func (l *List) indentOf(c *Canvas, row Row) string {
	if row.Indent != "" {
		return row.Indent
	}
	return strings.Repeat(" ", max(0, row.Depth)*c.Chrome().Indent)
}

// LeadWidth is the columns drawn before a row's own text, for a header that
// has to start in the same place.
//
// The FIXED part only: the cursor marker and the status column. An indent
// varies per row by definition, so a caller aligning a table header wants this
// number and not that one.
//
// It exists because there was no way to ask. Two rebuilds computed
// `comp.Width(list.Marker) + comp.Width(glyph)` by hand, and one of them first
// reached for [List.StatusRows] — which answers a different question in a
// different unit and compiles perfectly (issue 76).
func (l *List) LeadWidth() int {
	w := 0
	if l.Marker != "" {
		w = Width(l.Marker)
		if l.Blank != "" {
			w = max(w, Width(l.Blank))
		}
	}
	return w + max(0, l.StatusWidth)
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

// spansText is a row's words without its colours, for when the selection
// paints over them.
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
