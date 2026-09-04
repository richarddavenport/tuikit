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

// Spans lays out STYLED cells, keeping each cell's own style.
//
// The same layout as [Table.Rows] and the same widths; what changes is that a
// cell arrives as a [Segment] and leaves as one, so a column whose colour
// carries meaning still has it.
//
// Rows returns joined strings, which is right when a whole line is one colour —
// and it silently deletes information when it is not. pgctl ported a manifest
// onto Table and had to make three columns plain: a rule's data mode was amber
// for "none" and "filtered", a carried count was amber when a table came across
// filtered rather than whole, and an unknown size was a muted dash. All three
// became words with a comment explaining the loss (issue 51). Words are a
// reasonable fallback and they are not the same thing: colour is read without
// being looked at, which is the entire job of a status column.
//
// Padding is emitted as an unstyled segment rather than folded into the cell,
// so a cell with a background does not paint the gap after it.
func (t Table) Spans(w int, rows [][]Segment) [][]Segment {
	widths := t.widths(w, textOf(rows))

	out := make([][]Segment, 0, len(rows))
	for _, cells := range rows {
		var line []Segment
		for i, width := range widths {
			if i > 0 && t.Gap > 0 {
				line = append(line, Segment{Text: strings.Repeat(" ", t.Gap)})
			}
			var cell Segment
			if i < len(cells) {
				cell = cells[i]
			}
			line = append(line, pad(cell, width, t.Columns[i].Right)...)
		}
		out = append(out, trimTrailing(line))
	}
	return out
}

// pad fits one styled cell to a width, truncating what is too long and
// surrounding what is too short.
func pad(cell Segment, width int, right bool) []Segment {
	text := Truncate(cell.Text, width)
	gap := width - Width(text)
	if gap <= 0 {
		return []Segment{{Text: text, Style: cell.Style}}
	}
	blank := Segment{Text: strings.Repeat(" ", gap)}
	if right {
		return []Segment{blank, {Text: text, Style: cell.Style}}
	}
	return []Segment{{Text: text, Style: cell.Style}, blank}
}

// textOf drops the styles, so the widths are measured by the one function that
// already knows how — a second width calculation is a second answer.
func textOf(rows [][]Segment) [][]string {
	out := make([][]string, len(rows))
	for i, cells := range rows {
		out[i] = make([]string, len(cells))
		for j, cell := range cells {
			out[i][j] = cell.Text
		}
	}
	return out
}

// trimTrailing drops the blank segments off the end, matching Rows's
// TrimRight: a row that paints its trailing padding claims cells it is not
// using, and a click on one then lands on the row rather than past it.
func trimTrailing(line []Segment) []Segment {
	for len(line) > 0 {
		last := line[len(line)-1]
		if strings.TrimRight(last.Text, " ") != "" {
			break
		}
		line = line[:len(line)-1]
	}
	return line
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
