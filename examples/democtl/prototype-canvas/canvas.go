// PROTOTYPE — throwaway. Answers one question: should tuikit's comp package be
// canvas-based or string-based? Not on main; see design/mouse.md for the verdict.
//
// The idea under test: components draw cells into a grid instead of returning
// strings. Each cell carries a rune, a style, and the ID of whatever drew it.
// If that works, hit-testing is a lookup rather than a bookkeeping exercise.
package main

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

type Rect struct{ X, Y, W, H int }

func (r Rect) Contains(x, y int) bool {
	return x >= r.X && x < r.X+r.W && y >= r.Y && y < r.Y+r.H
}

// Inset shrinks a rect by n on every side — the "inside the border" case, which
// is most of what a component wants.
func (r Rect) Inset(n int) Rect {
	return Rect{r.X + n, r.Y + n, max(0, r.W-2*n), max(0, r.H-2*n)}
}

// Cell is one character position.
//
// Owner is the whole point. A string knows what it looks like; a cell knows what
// it IS, and that is what a mouse click needs.
type Cell struct {
	R     rune
	S     *lipgloss.Style // pointer, so runs can be grouped by identity when serialising
	Owner string
}

type Canvas struct {
	W, H  int
	cells []Cell
}

func NewCanvas(w, h int) *Canvas {
	c := &Canvas{W: w, H: h, cells: make([]Cell, w*h)}
	for i := range c.cells {
		c.cells[i].R = ' '
	}
	return c
}

func (c *Canvas) at(x, y int) *Cell {
	if x < 0 || y < 0 || x >= c.W || y >= c.H {
		return nil
	}
	return &c.cells[y*c.W+x]
}

// Set writes one cell. Out-of-bounds is a no-op rather than a panic: clipping to
// the canvas is the behaviour every caller wants, and doing it here means no
// component ever needs a bounds check.
//
// This is what makes "nothing wider than the terminal" structural. There is no
// way to draw past the edge — the coordinate simply does not exist.
func (c *Canvas) Set(x, y int, r rune, s *lipgloss.Style, owner string) {
	if cell := c.at(x, y); cell != nil {
		*cell = Cell{R: r, S: s, Owner: owner}
	}
}

// Text draws a plain string. Plain is the point: the canvas never parses ANSI,
// because style lives on the cell rather than in the text. That deletes the
// whole trim-versus-clip problem — there is no escape sequence to miscount.
func (c *Canvas) Text(x, y int, text string, s *lipgloss.Style, owner string) {
	for i, r := range []rune(text) {
		c.Set(x+i, y, r, s, owner)
	}
}

func (c *Canvas) Fill(r Rect, ru rune, s *lipgloss.Style, owner string) {
	for y := r.Y; y < r.Y+r.H; y++ {
		for x := r.X; x < r.X+r.W; x++ {
			c.Set(x, y, ru, s, owner)
		}
	}
}

// Box draws a frame with the six characters the glyph set allows.
func (c *Canvas) Box(r Rect, title string, s *lipgloss.Style, ts *lipgloss.Style, owner string) {
	if r.W < 2 || r.H < 2 {
		return
	}
	right, bottom := r.X+r.W-1, r.Y+r.H-1
	for x := r.X + 1; x < right; x++ {
		c.Set(x, r.Y, '─', s, owner)
		c.Set(x, bottom, '─', s, owner)
	}
	for y := r.Y + 1; y < bottom; y++ {
		c.Set(r.X, y, '│', s, owner)
		c.Set(right, y, '│', s, owner)
	}
	c.Set(r.X, r.Y, '┌', s, owner)
	c.Set(right, r.Y, '┐', s, owner)
	c.Set(r.X, bottom, '└', s, owner)
	c.Set(right, bottom, '┘', s, owner)
	if title != "" {
		c.Text(r.X+2, r.Y, " "+title+" ", ts, owner)
	}
}

// OwnerAt is the entire hit-testing implementation.
//
// No region list to keep in step with the drawing, because the drawing IS the
// region list. A component that moved but forgot to update its rect is not a
// bug that can exist here.
func (c *Canvas) OwnerAt(x, y int) string {
	if cell := c.at(x, y); cell != nil {
		return cell.Owner
	}
	return ""
}

// String serialises the grid, grouping runs that share a style so the output is
// not one escape sequence per character.
func (c *Canvas) String() string {
	var b strings.Builder
	for y := range c.H {
		var run []rune
		var cur *lipgloss.Style
		flush := func() {
			if len(run) == 0 {
				return
			}
			if cur == nil {
				b.WriteString(string(run))
			} else {
				b.WriteString(cur.Render(string(run)))
			}
			run = run[:0]
		}
		for x := range c.W {
			cell := c.cells[y*c.W+x]
			if cell.S != cur {
				flush()
				cur = cell.S
			}
			run = append(run, cell.R)
		}
		flush()
		if y < c.H-1 {
			b.WriteByte('\n')
		}
	}
	return b.String()
}
