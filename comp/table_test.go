package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// The marker answers "can I reach it" before the name answers "which is it",
// so the narrow leading column is the ordinary case rather than a special one.
func TestANarrowMarkerColumnComesFirst(t *testing.T) {
	tbl := Table{Gap: 1, Columns: []Column{{Width: 1}, {Fill: true}, {Width: 5, Right: true}}}
	got := tbl.Rows(24, [][]string{
		{"✓", "api_gateway", "3/3"},
		{"✗", "api_migrate", "0/1"},
	})

	for _, line := range got {
		if Width(line) > 24 {
			t.Errorf("%q is %d columns", line, Width(line))
		}
	}
	if !strings.HasPrefix(got[0], "✓ api_gateway") {
		t.Errorf("got %q", got[0])
	}
	// The counts line up, which is the whole point of a column.
	if at(got[0], "3/3") != at(got[1], "0/1") {
		t.Errorf("the last column does not line up:\n%q\n%q", got[0], got[1])
	}
}

// A column with no width takes what its widest cell needs, so a table of short
// names does not leave a trench down the middle.
func TestANaturalColumnFitsItsWidestCell(t *testing.T) {
	tbl := Table{Gap: 1, Columns: []Column{{}, {Fill: true}}}
	got := tbl.Rows(30, [][]string{{"a", "one"}, {"longer", "two"}})

	if !strings.HasPrefix(got[0], "a      one") {
		t.Errorf("got %q", got[0])
	}
	if !strings.HasPrefix(got[1], "longer one") && !strings.HasPrefix(got[1], "longer two") {
		t.Errorf("got %q", got[1])
	}
}

// The filler takes what is left, so the table is exactly as wide as it was
// told.
func TestTheFillerTakesWhatIsLeft(t *testing.T) {
	tbl := Table{Gap: 1, Columns: []Column{{Width: 3}, {Fill: true}, {Width: 4}}}
	rows := tbl.Rows(40, [][]string{{"abc", "middle", "tail"}})

	if got := Width(fit(rows[0], 40, false)); got != 40 {
		t.Errorf("the row does not fill its width: %d", got)
	}
	if !strings.Contains(rows[0], "middle") || !strings.HasSuffix(rows[0], "tail") {
		t.Errorf("got %q", rows[0])
	}
}

// A cell too long for its column is truncated with an ellipsis, so the columns
// after it stay where they are. A cell that pushed its neighbours along would
// make every row a different shape.
func TestALongCellIsTruncatedRatherThanPushing(t *testing.T) {
	tbl := Table{Gap: 1, Columns: []Column{{Width: 8}, {Width: 6}}}
	got := tbl.Rows(20, [][]string{
		{"short", "ok"},
		{"a very long name indeed", "ok"},
	})

	if at(got[0], "ok") != at(got[1], "ok") {
		t.Errorf("a long cell moved the next column:\n%q\n%q", got[0], got[1])
	}
	if !strings.Contains(got[1], "…") {
		t.Errorf("the cut cell does not say it was cut: %q", got[1])
	}
}

// Wide runes count as two columns, or a CJK name silently knocks a table out
// of alignment for every row below it.
func TestColumnsAreMeasuredInColumns(t *testing.T) {
	tbl := Table{Gap: 1, Columns: []Column{{Width: 8}, {Fill: true}}}
	got := tbl.Rows(20, [][]string{{"世界", "after"}, {"ab", "after"}})

	// Compared in COLUMNS. Comparing byte offsets is the same mistake the
	// table is being tested for: 世界 is six bytes and four columns.
	if at(got[0], "after") != at(got[1], "after") {
		t.Errorf("a wide rune knocked the column out:\n%q\n%q", got[0], got[1])
	}
}

// A table lays out; a list selects and scrolls. Rows that do both are the two
// composed, which is how democtl's service list already works.
func TestRowsCanBeHandedToAList(t *testing.T) {
	tbl := Table{Gap: 1, Columns: []Column{{Width: 1}, {Fill: true}}}
	lines := tbl.Rows(20, [][]string{{"✓", "one"}, {"✗", "two"}})

	rows := make([]Row, len(lines))
	for i, line := range lines {
		rows[i] = Row{Text: line}
	}
	c := NewCanvas(20, 4)
	l := &List{Name: "t"}
	l.Draw(c, c.Bounds(), rows)

	if !strings.Contains(c.String(), "✓ one") {
		t.Errorf("the table's rows did not reach the list:\n%s", c.String())
	}
}

func TestDrawStaysInsideItsRect(t *testing.T) {
	c := NewCanvas(30, 4)
	Table{Gap: 1, Columns: []Column{{Fill: true}}}.
		Draw(c, Rect{X: 0, Y: 0, W: 10, H: 1}, [][]string{{"far too long for ten"}, {"second"}}, nil, Region("t"))

	lines := strings.Split(c.String(), "\n")
	if Width(lines[0]) > 10 {
		t.Errorf("drew %d columns into 10: %q", Width(lines[0]), lines[0])
	}
	if strings.TrimSpace(lines[1]) != "" {
		t.Errorf("drew a second row into a one-row rect: %q", lines[1])
	}
}

// at is where a cell starts, in columns rather than bytes.
func at(line, cell string) int {
	i := strings.Index(line, cell)
	if i < 0 {
		return -1
	}
	return Width(line[:i])
}

// A cell keeps its own style. Issue 51: Rows joins to strings, so the database
// tool had to make three columns plain and write a comment explaining the
// loss.
func TestSpansKeepEachCellsStyle(t *testing.T) {
	amber := lipgloss.NewStyle().Foreground(lipgloss.Color("214"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	tbl := Table{Columns: []Column{{Width: 10}, {Width: 8}}, Gap: 1}
	out := tbl.Spans(24, [][]Segment{
		{{Text: "orders"}, {Text: "filtered", Style: &amber}},
		{{Text: "customers"}, {Text: "all", Style: &muted}},
	})

	if len(out) != 2 {
		t.Fatalf("got %d rows", len(out))
	}
	var found bool
	for _, span := range out[0] {
		if span.Text == "filtered" && span.Style == &amber {
			found = true
		}
	}
	if !found {
		t.Errorf("the amber cell lost its style: %+v", out[0])
	}
}

// The layout is the same as Rows's, or a tool that switches loses its columns.
func TestSpansLaysOutTheSameAsRows(t *testing.T) {
	tbl := Table{Columns: []Column{{Width: 6}, {Fill: true}, {Width: 5, Right: true}}, Gap: 2}
	cells := [][]string{{"api", "the api service", "2/2"}, {"web", "web", "10/10"}}

	styled := make([][]Segment, len(cells))
	for i, row := range cells {
		styled[i] = make([]Segment, len(row))
		for j, cell := range row {
			styled[i][j] = Segment{Text: cell}
		}
	}

	plain := tbl.Rows(40, cells)
	for i, line := range tbl.Spans(40, styled) {
		var b string
		for _, span := range line {
			b += span.Text
		}
		if b != plain[i] {
			t.Errorf("row %d:\n  Spans %q\n  Rows  %q", i, b, plain[i])
		}
	}
}

// Padding is its own segment, so a cell with a background does not paint the
// gap after it.
func TestPaddingIsNotPartOfTheCell(t *testing.T) {
	sel := lipgloss.NewStyle().Background(lipgloss.Color("7"))
	tbl := Table{Columns: []Column{{Width: 10}}}

	for _, line := range tbl.Spans(10, [][]Segment{{{Text: "hi", Style: &sel}}}) {
		for _, span := range line {
			if span.Style == &sel && span.Text != "hi" {
				t.Errorf("the styled segment carries padding: %q", span.Text)
			}
		}
	}
}
