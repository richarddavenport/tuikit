// Package comp is tuikit's component substrate: a grid of cells that components
// draw into, instead of strings they return.
//
// A component that returns a string cannot be clicked. Hit-testing needs
// geometry, and `lipgloss.JoinHorizontal` throws it away — by the time View()
// has a string, nothing knows that row 4, columns 2 to 31 was services[2]. So
// every cell carries the ID of whatever drew it, and OwnerAt is the whole hit
// test. There is no region list to keep in step with the drawing, because the
// drawing IS the region list: a component that moved but forgot to update its
// rect is not a bug that can exist here.
//
// Two properties fall out of that, and both were surprises worth stating:
//
//   - Overflow stops being a class of bug. Set clips to the canvas, so drawing
//     past the edge is not an error to catch — it is a coordinate that does not
//     exist. Two of the three bugs pgctl's first capture found, and both of
//     democtl's, were something drawn wider than its container.
//   - The trim-versus-clip problem stops existing. Style lives on the cell, so
//     the canvas never parses ANSI and there is no escape sequence to miscount.
//     Lip Gloss keeps styling; its layout helpers do not come along.
//
// See design/mouse.md for how this was settled, and decision 18.
package comp

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
	"github.com/rivo/uniseg"

	"github.com/richarddavenport/tuikit/theme"
)

// Rect is a region of the grid.
type Rect struct{ X, Y, W, H int }

// Contains reports whether a point is inside the rect — the question a click
// asks of a pane before it asks anything else.
func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

// Inset shrinks a rect by n on every side: the inside-the-border case, which is
// most of what a component wants.
func (r Rect) Inset(n int) Rect {
	return Rect{r.X + n, r.Y + n, max(0, r.W-2*n), max(0, r.H-2*n)}
}

// Empty reports whether the rect has no cells. Worth asking before drawing into
// it: a pane squeezed to nothing by a narrow terminal is a normal state, not a
// failure.
func (r Rect) Empty() bool { return r.W <= 0 || r.H <= 0 }

// Right and Bottom are the last column and row inside the rect.
func (r Rect) Right() int  { return r.X + r.W - 1 }
func (r Rect) Bottom() int { return r.Y + r.H - 1 }

// Name is a region's name, declared once by the tool as a constant.
//
// A named type rather than a string so that the set of regions a tool has is
// something you can read in one place, and so a capture script or a right-click
// menu cannot quietly invent one. Canvas.Region answers whether a name was
// actually drawn, which is the other half of that guarantee.
type Name string

// NoIndex marks an ID that names a region as a whole rather than one item in it.
const NoIndex = -1

// ID is the identity of whatever drew a cell.
//
// An IDENTITY, never a screen position. The index is the item's absolute index
// in its collection, not the row it happens to occupy — "the ID is where it is
// on screen" works perfectly until something scrolls, and then it does not fail
// visibly, it performs the WRONG ACTION on the item that moved into that row.
// This is the first place the canvas design could have gone quietly wrong.
type ID struct {
	Name  Name
	Index int
}

// Region names a whole region.
func Region(n Name) ID { return ID{Name: n, Index: NoIndex} }

// At names one item within a region, by its absolute index.
func (id ID) At(i int) ID { return ID{Name: id.Name, Index: i} }

// Zero reports whether the ID names nothing — what OwnerAt returns for a cell
// nobody claimed, and for a coordinate off the canvas.
func (id ID) Zero() bool { return id.Name == "" }

// String is the name a capture script or a menu uses: `services.row[2]`.
func (id ID) String() string {
	if id.Zero() {
		return ""
	}
	if id.Index == NoIndex {
		return string(id.Name)
	}
	return string(id.Name) + "[" + strconv.Itoa(id.Index) + "]"
}

// Cell is one character position.
//
// Text holds one grapheme cluster rather than a rune, because a rune is not a
// character: an emoji built from a zero-width joiner is five runes and one
// thing a reader sees, and a grid that splits it draws two halves of nothing.
//
// An empty Text marks the CONTINUATION of a wide cluster in the cell to its
// left. It is not a space, and serialising skips it — emitting a space there
// would push everything after it one column right, which is exactly the class
// of bug the canvas exists to make impossible.
type Cell struct {
	Text  string
	Style *lipgloss.Style
	Owner ID
}

// Continuation reports whether this cell is the right-hand half of a wide one.
func (c Cell) Continuation() bool { return c.Text == "" }

// Canvas is the grid.
type Canvas struct {
	w, h  int
	cells []Cell
	// clip is what may be drawn into. The whole canvas by default; narrower
	// for the view a component is handed.
	clip Rect
	// chrome is what the components draw WITH: the box characters, the
	// separators, the markers. Carried by the canvas rather than passed to
	// each component, because a component handed its own chrome is a component
	// that can be handed the wrong one — or none, and draw a frame out of
	// empty strings.
	chrome theme.Chrome
}

// NewCanvas makes a canvas of blanks.
func NewCanvas(w, h int) *Canvas {
	c := &Canvas{w: max(0, w), h: max(0, h)}
	c.clip = Rect{0, 0, c.w, c.h}
	c.chrome = theme.DefaultChrome
	c.cells = make([]Cell, c.w*c.h)
	for i := range c.cells {
		c.cells[i].Text = " "
	}
	return c
}

// Bounds is the whole canvas as a rect, for a component handed the window.
func (c *Canvas) Bounds() Rect { return Rect{0, 0, c.w, c.h} }

func (c *Canvas) in(x, y int) bool {
	return x >= 0 && y >= 0 && x < c.w && y < c.h && c.clip.Contains(x, y)
}

// Chrome is what this canvas draws with.
func (c *Canvas) Chrome() theme.Chrome { return c.chrome }

// WithChrome returns a view drawing in a different vocabulary.
//
// A view, like Clip, so a tool sets it once at the top and every component
// underneath inherits it. There is no way for one pane to end up with rounded
// corners because somebody threaded the chrome through five call sites and
// missed one.
func (c *Canvas) WithChrome(ch theme.Chrome) *Canvas {
	view := *c
	view.chrome = ch
	return &view
}

// Clip returns a view of the canvas that cannot draw outside r.
//
// The canvas already makes drawing off the TERMINAL impossible. This makes
// drawing outside the rect a component was GIVEN impossible too, which is the
// other half — and the half design/roadmap.md flagged as still open: "the
// canvas guarantees nothing is drawn outside the canvas, not that a component
// stayed inside the rect it was given."
//
// A view rather than a copy: the cells are shared, so a component draws into
// the real frame and simply cannot reach past its own box. Coordinates stay
// absolute, so a component's rect arithmetic is unchanged and nothing has to be
// translated back for a hit test.
func (c *Canvas) Clip(r Rect) *Canvas {
	view := *c
	view.clip = intersect(c.clip, r)
	return &view
}

func intersect(a, b Rect) Rect {
	x, y := max(a.X, b.X), max(a.Y, b.Y)
	right, bottom := min(a.Right(), b.Right()), min(a.Bottom(), b.Bottom())
	return Rect{X: x, Y: y, W: max(0, right-x+1), H: max(0, bottom-y+1)}
}

func (c *Canvas) at(x, y int) *Cell {
	if !c.in(x, y) {
		return nil
	}
	return &c.cells[y*c.w+x]
}

// Width measures a cluster the way Lip Gloss does.
//
// The same function lipgloss.Width calls, deliberately. If the canvas measured
// with one table and the harness checked with another, a frame would be laid
// out to one width and reported as another, and the disagreement would show up
// as a mysterious column of drift in a capture rather than as an error.
func Width(s string) int { return ansi.StringWidth(s) }

// Set writes one grapheme cluster and returns the columns it claimed.
//
// Out of bounds is a no-op rather than a panic: clipping is what every caller
// wants, and doing it here means no component needs a bounds check. A wide
// cluster that would hang off the right edge draws nothing at all — there is no
// half of a glyph to draw, and leaving the row a column short would be worse
// than leaving it blank.
func (c *Canvas) Set(x, y int, cluster string, s *lipgloss.Style, owner ID) int {
	w := Width(cluster)
	if w <= 0 || !c.in(x, y) || x+w > c.w {
		return 0
	}

	// Whatever is being overwritten may be half of a wide cluster. Both of its
	// cells have to go, or the leftover half is an orphan: a continuation with
	// nothing to continue, which serialises as a missing column.
	c.clear(x, y)
	c.clear(x+w-1, y)

	*c.at(x, y) = Cell{Text: cluster, Style: s, Owner: owner}
	for i := 1; i < w; i++ {
		// The continuation carries the same owner, so a click on the right half
		// of a wide glyph hits the thing that drew it.
		*c.at(x+i, y) = Cell{Text: "", Style: s, Owner: owner}
	}
	return w
}

// clear blanks the whole cluster occupying a cell, however wide it is and
// whichever half was named.
func (c *Canvas) clear(x, y int) {
	cell := c.at(x, y)
	if cell == nil {
		return
	}
	lead := x
	for lead > 0 && c.cells[y*c.w+lead].Continuation() {
		lead--
	}
	blank := Cell{Text: " "}
	w := max(1, Width(c.cells[y*c.w+lead].Text))
	for i := 0; i < w && lead+i < c.w; i++ {
		c.cells[y*c.w+lead+i] = blank
	}
}

// Text draws a string, one grapheme at a time, and returns the columns drawn.
//
// Plain text is the point. The canvas never parses ANSI, so there is nothing
// here that can miscount an escape sequence as a column — the bug that put a
// header off the side of the screen in a real tool, and the one that ate
// democtl's tab strip at 80 columns.
func (c *Canvas) Text(x, y int, text string, s *lipgloss.Style, owner ID) int {
	drawn, rest, state := 0, text, -1
	for rest != "" {
		var cluster string
		cluster, rest, _, state = uniseg.FirstGraphemeClusterInString(rest, state)
		n := c.Set(x+drawn, y, cluster, s, owner)
		if n == 0 {
			// Either the canvas ran out or the cluster will not fit. Stop
			// rather than skipping it, so a clipped string is a prefix of what
			// was asked for and never a version with a hole in the middle.
			break
		}
		drawn += n
	}
	return drawn
}

// Fill covers a rect with one cluster: the background of a pane, or the blank
// a row is padded with.
func (c *Canvas) Fill(r Rect, cluster string, s *lipgloss.Style, owner ID) {
	w := max(1, Width(cluster))
	for y := r.Y; y <= r.Bottom(); y++ {
		for x := r.X; x <= r.Right(); x += w {
			c.Set(x, y, cluster, s, owner)
		}
	}
}

// Box draws a frame in the canvas's chrome.
//
// Six characters, not eleven: there are no tee or cross pieces, so a component
// that wants a title writes it into the top edge itself rather than breaking
// the line. That is a constraint the glyph set imposes and the components were
// designed around — a set with tees would need components that know what to do
// with them.
func (c *Canvas) Box(r Rect, s *lipgloss.Style, owner ID) {
	if r.W < 2 || r.H < 2 {
		return
	}
	box := c.chrome.Box
	for x := r.X + 1; x < r.Right(); x++ {
		c.Set(x, r.Y, box.Top, s, owner)
		c.Set(x, r.Bottom(), box.Bottom, s, owner)
	}
	for y := r.Y + 1; y < r.Bottom(); y++ {
		c.Set(r.X, y, box.Left, s, owner)
		c.Set(r.Right(), y, box.Right, s, owner)
	}
	c.Set(r.X, r.Y, box.TopLeft, s, owner)
	c.Set(r.Right(), r.Y, box.TopRight, s, owner)
	c.Set(r.X, r.Bottom(), box.BottomLeft, s, owner)
	c.Set(r.Right(), r.Bottom(), box.BottomRight, s, owner)
}

// OwnerAt is the entire hit test.
func (c *Canvas) OwnerAt(x, y int) ID {
	if cell := c.at(x, y); cell != nil {
		return cell.Owner
	}
	return ID{}
}

// CellAt reads one cell, for a test or a guard that wants to know what is drawn
// where without serialising the frame first.
func (c *Canvas) CellAt(x, y int) (Cell, bool) {
	if cell := c.at(x, y); cell != nil {
		return *cell, true
	}
	return Cell{}, false
}

// Region is the box a region occupies, and whether it was drawn at all.
//
// The false is the useful half: it is what lets a capture script that says
// `click services.row[2]` fail loudly when nothing drew that, rather than
// clicking column zero and reporting success.
func (c *Canvas) Region(id ID) (Rect, bool) {
	minX, minY, maxX, maxY := c.w, c.h, -1, -1
	for y := 0; y < c.h; y++ {
		for x := 0; x < c.w; x++ {
			if c.cells[y*c.w+x].Owner != id {
				continue
			}
			minX, minY = min(minX, x), min(minY, y)
			maxX, maxY = max(maxX, x), max(maxY, y)
		}
	}
	if maxX < 0 {
		return Rect{}, false
	}
	return Rect{minX, minY, maxX - minX + 1, maxY - minY + 1}, true
}

// String serialises the grid, grouping runs that share a style so the output is
// not one escape sequence per character.
//
// Styles are compared by POINTER. A tool holds its styles in a struct and hands
// out the same address every time, so identity is both cheaper and more honest
// than equality: two styles that happen to look alike are still two decisions.
func (c *Canvas) String() string {
	var b strings.Builder
	for y := 0; y < c.h; y++ {
		var run strings.Builder
		var cur *lipgloss.Style
		flush := func() {
			if run.Len() == 0 {
				return
			}
			if cur == nil {
				b.WriteString(run.String())
			} else {
				b.WriteString(cur.Render(run.String()))
			}
			run.Reset()
		}
		for x := 0; x <= c.lastInk(y); x++ {
			cell := c.cells[y*c.w+x]
			if cell.Continuation() {
				continue // its lead already wrote it
			}
			if cell.Style != cur {
				flush()
				cur = cell.Style
			}
			run.WriteString(cell.Text)
		}
		flush()
		if y < c.h-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}

// lastInk is the last column of a row worth emitting.
//
// Trailing blanks are dropped, so a footer does not arrive with sixty spaces
// after it — they are invisible in a terminal, and they make a golden diff
// noise. A blank that carries a STYLE is not trailing whitespace though: it is
// a row painted to its edge, which is how a selected row gets a background all
// the way across, so those are kept.
func (c *Canvas) lastInk(y int) int {
	for x := c.w - 1; x >= 0; x-- {
		cell := c.cells[y*c.w+x]
		if cell.Style != nil || (cell.Text != " " && !cell.Continuation()) {
			return x
		}
	}
	return -1
}

// ParseID reads a name as a capture script writes it: `services.row[2]`, or
// `services` for a region as a whole.
//
// The inverse of ID.String, so a script and a frame speak the same language.
// Whether the region exists is a separate question, and Canvas.Region answers
// it — a name can be well-formed and still name nothing.
func ParseID(s string) (ID, error) {
	if s == "" {
		return ID{}, fmt.Errorf("comp: no region named")
	}
	open := strings.IndexByte(s, '[')
	if open < 0 {
		return Region(Name(s)), nil
	}
	if !strings.HasSuffix(s, "]") {
		return ID{}, fmt.Errorf("comp: %q opens an index and never closes it", s)
	}
	i, err := strconv.Atoi(s[open+1 : len(s)-1])
	if err != nil {
		return ID{}, fmt.Errorf("comp: %q has an index that is not a number", s)
	}
	return Region(Name(s[:open])).At(i), nil
}
