package comp

import "github.com/charmbracelet/lipgloss"

// Menu is a short list of actions, at a point or on a thing.
//
// # Where this came from
//
// democtl's context menu, which was forty lines of drawing with three real
// rules buried in it — and azctl needs the same menu, which is the second
// consumer that makes this an extraction rather than a guess.
//
// The rules, all of them learned rather than designed:
//
//   - A menu opened from the KEYBOARD knows only which region it is on, and
//     finds out where that is when the frame is drawn. The row it belongs to
//     may have scrolled, or the pane may have been resized, since the key was
//     pressed. Asking the frame being drawn is the only answer that cannot be
//     stale — which is DrawOn.
//   - It is nudged back on screen rather than clipped. The canvas would happily
//     draw half a menu off the edge, which is what a canvas is for, but half a
//     menu is a list of actions you cannot read — a different thing from a pane
//     that is cut off.
//   - The key sits beside the action, because it is the SAME LIST. A menu and a
//     keymap maintained separately is how the two paths to an action come to
//     disagree, and an action reachable only by right-click is unreachable
//     inside a multiplexer that keeps right-click for itself.
type Menu struct {
	// Name is the frame; Item is the rows, indexed. Two names rather than one,
	// because a click on the border is not a click on an action and a menu
	// that treats them alike closes when you grab its edge.
	Name, Item Name

	// Items are the actions. A Hint, so the menu and the key hint bar are
	// built from one type.
	Items  []Hint
	Cursor int

	// Min is the narrowest the box may be, so a menu of one short word still
	// looks like a menu.
	Min int

	Border, Style, Selected *lipgloss.Style
}

// DrawAt places the menu with its top-left corner at x, y — where a right-click
// landed.
func (m Menu) DrawAt(c *Canvas, x, y int) Rect {
	if len(m.Items) == 0 {
		return Rect{}
	}
	w := max(m.Min, 4)
	for _, item := range m.Items {
		// The label, the key, a column between them, and the two borders and
		// two insets around the pair.
		w = max(w, Width(item.Label)+Width(item.Key)+6)
	}
	r := Rect{X: x, Y: y, W: w, H: len(m.Items) + 2}

	bounds := c.Bounds()
	r.X = min(max(r.X, 0), max(0, bounds.W-r.W))
	r.Y = min(max(r.Y, 0), max(0, bounds.H-r.H))

	inner := Pane{Border: m.Border, Focused: true, Focus: m.Border}.Draw(c, r, Region(m.Name))
	for i, item := range m.Items {
		id := Region(m.Item).At(i)
		style := m.Style
		if i == m.Cursor {
			style = m.Selected
		}
		row := Rect{X: inner.X, Y: inner.Y + i, W: inner.W, H: 1}
		c.Fill(row, " ", style, id)
		c.Text(row.X+1, row.Y, item.Label, style, id)
		c.Text(row.X+row.W-Width(item.Key)-1, row.Y, item.Key, style, id)
	}
	return r
}

// DrawOn places the menu on a region, wherever that region is IN THIS FRAME.
//
// The keyboard path. "At the cursor" means at the thing the cursor is on, and
// the only place that is knowable is the frame being drawn — the row may have
// scrolled since the key was pressed. Falls back to the top-left when the
// region was not drawn at all, which is a menu somewhere odd rather than a menu
// that vanished.
func (m Menu) DrawOn(c *Canvas, on ID) Rect {
	x, y := 0, 0
	if at, ok := c.Region(on); ok {
		x, y = at.X+2, at.Y
	}
	return m.DrawAt(c, x, y)
}
