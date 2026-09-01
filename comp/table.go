package comp

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Table is rows of aligned columns.
//
// # Where this came from
//
// pgctl's rows.go (192 lines) and swarmctl's pane and detail tables. Both build
// a row as a status marker, a name, and then whatever else fits — and pgctl's
// comment says why that order: "the marker answers 'can I reach it' before the
// name answers 'which is it', because an unreachable environment changes what
// every panel below is showing."
//
// So a Table is columns with an ALIGNMENT and a WIDTH RULE, and the leading
// narrow column is the ordinary case rather than a special one.
//
// # What a Table is not
//
// It is not a List with columns. A List owns a selection and a viewport; a
// Table owns alignment. Rows that both scroll and align are a List whose rows a
// Table laid out, which is two components composed rather than one component
// with more fields — and it is how democtl's service list already works, having
// formatted its rows before handing them over.
type Table struct {
	Columns []Column
	// Gap is the columns between one column and the next.
	Gap int
}

// Column is one column's rule.
type Column struct {
	// Width is fixed when above zero. Zero means the column takes what the
	// widest cell in it needs, which is what a status marker or a replica
	// count wants.
	Width int
	// Fill gives this column whatever is left over. At most one should have
	// it; the last one wins, because a table with two greedy columns has no
	// answer and silently picking one hides the mistake.
	Fill bool
	// Right aligns the cell to the column's right edge, for numbers.
	Right bool
}

// Rows lays out cells into lines of exactly w columns.
//
// Strings out rather than drawing, so the result can be handed to a List and
// scrolled, given a style per row, or measured. A Table that drew itself would
// be a Table that could not be selected from.
func (t Table) Rows(w int, rows [][]string) []string {
	widths := t.widths(w, rows)

	out := make([]string, 0, len(rows))
	for _, cells := range rows {
		var b strings.Builder
		for i, width := range widths {
			if i > 0 {
				b.WriteString(strings.Repeat(" ", t.Gap))
			}
			cell := ""
			if i < len(cells) {
				cell = cells[i]
			}
			b.WriteString(fit(cell, width, t.Columns[i].Right))
		}
		out = append(out, strings.TrimRight(b.String(), " "))
	}
	return out
}

// widths resolves each column: fixed as asked, natural to its widest cell, and
// the filler to whatever is left.
func (t Table) widths(w int, rows [][]string) []int {
	widths := make([]int, len(t.Columns))
	spare := w - t.Gap*max(0, len(t.Columns)-1)

	for i, col := range t.Columns {
		switch {
		case col.Fill:
			continue
		case col.Width > 0:
			widths[i] = col.Width
		default:
			for _, cells := range rows {
				if i < len(cells) {
					widths[i] = max(widths[i], Width(cells[i]))
				}
			}
		}
		spare -= widths[i]
	}
	for i, col := range t.Columns {
		if col.Fill {
			widths[i] = max(0, spare)
		}
	}
	return widths
}

// fit pads or truncates one cell to its column.
func fit(s string, w int, right bool) string {
	s = Truncate(s, w)
	pad := w - Width(s)
	if pad <= 0 {
		return s
	}
	if right {
		return strings.Repeat(" ", pad) + s
	}
	return s + strings.Repeat(" ", pad)
}

// Draw writes the rows into r, one per line. The common case is handing them to
// a List instead.
func (t Table) Draw(c *Canvas, r Rect, rows [][]string, style *lipgloss.Style, id ID) {
	c = c.Clip(r)
	for i, line := range t.Rows(r.W, rows) {
		if i >= r.H {
			return
		}
		c.Text(r.X, r.Y+i, line, style, id)
	}
}
