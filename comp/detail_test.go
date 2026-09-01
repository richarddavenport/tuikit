package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

const inspector Name = "inspector"

func drawDetail(d Detail, w, h int) []string {
	c := NewCanvas(w, h)
	d.Draw(c, c.Bounds(), Region(inspector))
	return strings.Split(c.String(), "\n")
}

// The whole reason this exists: both tools carried `y` from call to call and
// guarded the bottom in one of the two helpers that needed it. A pane one row
// too tall writes over the border it sits in.
func TestItStopsAtTheBottomOfItsRect(t *testing.T) {
	c := NewCanvas(30, 6)
	// Two rows of room inside a box that owns the rest.
	c.Fill(c.Bounds(), "#", nil, Region("background"))

	d := Detail{Title: "vm-forge-0", Blocks: []Block{{Facts: []Fact{
		{Label: "group", Value: "rg-forge"},
		{Label: "location", Value: "uksouth"},
		{Label: "kind", Value: "StorageV2"},
		{Label: "type", Value: "microsoft.compute/virtualmachines"},
	}}}}
	d.Draw(c, Rect{X: 0, Y: 0, W: 30, H: 3}, Region(inspector))

	// Row 3 is outside the rect and must still be the background.
	if got := c.OwnerAt(0, 3); got.Name != "background" {
		t.Errorf("row 3 is owned by %v — the detail drew past its rect", got)
	}
}

// A caller that wants to put something underneath needs to know where the
// detail finished, INCLUDING the rows it could not draw — otherwise two panes
// stacked in a short window overlap silently.
func TestItReportsTheRowAfterTheLastEvenWhenItRanOut(t *testing.T) {
	c := NewCanvas(30, 20)
	d := Detail{Blocks: []Block{{Facts: []Fact{
		{Label: "a", Value: "1"}, {Label: "b", Value: "2"}, {Label: "c", Value: "3"},
	}}}}

	if got := d.Draw(c, Rect{X: 0, Y: 0, W: 30, H: 1}, Region(inspector)); got != 3 {
		t.Errorf("finished at %d, want 3 — the rows it could not draw still count", got)
	}
}

// A block's labels line up with each other and not with the rest of the pane.
// democtl padded to ten and azctl to sixteen, both arbitrary, and azctl's
// overflowed: a longer label pushed its value out of line with every other row.
func TestLabelsLineUpPerBlock(t *testing.T) {
	lines := drawDetail(Detail{Blocks: []Block{
		{Facts: []Fact{{Label: "id", Value: "one"}, {Label: "at", Value: "two"}}},
		{Heading: "tags", Facts: []Fact{{Label: "environment", Value: "prod"}}},
	}}, 40, 8)

	if lines[0] != "id one" || lines[1] != "at two" {
		t.Errorf("the first block did not use its own width: %q, %q", lines[0], lines[1])
	}
	// Indented under its heading, and using its own width rather than the two
	// columns the block above needed.
	if lines[4] != "  environment prod" {
		t.Errorf("the second block was dragged to the first's width: %q", lines[4])
	}
}

// A fact can colour its VALUE without colouring its label. A state that is red
// in the list and grey here is the same fact told twice, differently.
func TestAFactCanColourItsValueAlone(t *testing.T) {
	label := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	value := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	c := NewCanvas(30, 2)
	Detail{
		LabelStyle: &label,
		Blocks:     []Block{{Facts: []Fact{{Label: "state", Value: "failed", Style: &value}}}},
	}.Draw(c, c.Bounds(), Region(inspector))

	at := func(x int) *lipgloss.Style {
		cell, _ := c.CellAt(x, 0)
		return cell.Style
	}
	if at(0) != &label {
		t.Error("the label did not keep the label style")
	}
	if at(6) != &value {
		t.Error("the value did not take its own style")
	}
}

// Blocks read as blocks, but a pane that opens with an empty line looks like it
// failed to draw.
func TestThereIsNoBlankRowBeforeTheFirstThing(t *testing.T) {
	lines := drawDetail(Detail{Blocks: []Block{
		{Facts: []Fact{{Label: "group", Value: "rg-forge"}}},
		{Facts: []Fact{{Label: "kind", Value: "vm"}}},
	}}, 30, 6)

	if lines[0] != "group rg-forge" {
		t.Errorf("the pane opens with %q", lines[0])
	}
	if lines[1] != "" || lines[2] != "kind vm" {
		t.Errorf("the blocks did not separate: %q then %q", lines[1], lines[2])
	}
}

// Prose wraps; a fact truncates. A note is a sentence and losing the end of it
// loses the point, where a resource ID that runs off the edge is still
// recognisable from its start.
func TestProseWrapsAndAFactTruncates(t *testing.T) {
	lines := drawDetail(Detail{Blocks: []Block{
		{Text: "the tasks are replaced one at a time"},
		{Facts: []Fact{{Label: "type", Value: "microsoft.compute/virtualmachines"}}},
	}}, 20, 8)

	if lines[0] == "" || lines[1] == "" {
		t.Errorf("the prose did not wrap: %q, %q", lines[0], lines[1])
	}
	for _, l := range lines {
		if Width(l) > 20 {
			t.Errorf("%q is wider than the rect", l)
		}
	}
}

// A heading indents what is under it, by Chrome.Indent — the same rule a nested
// list row follows, so a pane and a tree agree about what nesting looks like.
// Without it the rows under "tags" are level with the heading and read as three
// more facts about the thing rather than as the tags.
func TestAnIndentedBlockUsesChromesIndent(t *testing.T) {
	lines := drawDetail(Detail{Blocks: []Block{
		{Heading: "tags", Facts: []Fact{{Label: "env", Value: "prod"}}},
	}}, 30, 6)

	found := false
	for _, l := range lines {
		if l == "  env prod" {
			found = true
		}
	}
	if !found {
		t.Errorf("no indented fact row in:\n%s", strings.Join(lines, "\n"))
	}
}
