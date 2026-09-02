package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func demoPalette() *Palette {
	muted := lipgloss.NewStyle()
	return &Palette{
		Title: "run", Name: "palette", Item: "palette.row",
		NoteStyle: &muted,
		Groups: []PaletteGroup{
			{Name: "service · api_api", Note: "2 tasks on 2 nodes", Items: []PaletteItem{
				{Label: "restart", Key: "R", Hint: "rolling, start-first"},
				{Label: "follow logs", Hint: "both tasks, interleaved"},
			}},
			{Name: "node · vm-qat-0", Note: "ready · leader", Items: []PaletteItem{
				{Label: "open a shell here", Key: "s", Hint: "ssh, over the route in use"},
				{Label: "drain", Refused: true, Hint: "refused · leader of a 5-node swarm"},
			}},
			{Name: "global", Note: "works on any screen", Items: []PaletteItem{
				{Label: "switch environment", Key: "E", Hint: "latest · qat · prd · local"},
				{Label: "prune unused images", Warn: true, Hint: "asks first · 8M reclaimable"},
			}},
		},
		Hints: []Hint{{Key: "enter", Label: "run"}, {Key: "esc", Label: "close"}},
	}
}

func render(p *Palette, w, h int) string {
	c := NewCanvas(w, h)
	p.Draw(c, Rect{X: 0, Y: 0, W: w, H: h})
	return c.String()
}

func TestPaletteShowsEveryGroupWithNothingTyped(t *testing.T) {
	got := render(demoPalette(), 92, 20)
	for _, want := range []string{
		"run",               // the title, in the border
		"service · api_api", // the headings
		"node · vm-qat-0", "global",
		"2 tasks on 2 nodes", // and their facts
		"restart", "R", "rolling, start-first",
		"type to search 6 commands",
		"enter run",
	} {
		if !strings.Contains(got, want) {
			t.Errorf("the palette is missing %q:\n%s", want, got)
		}
	}
}

// A refused item is LISTED, with the reason where its key would be. One that
// vanished is indistinguishable from one that never existed.
func TestARefusedCommandIsListedWithoutItsKey(t *testing.T) {
	got := render(demoPalette(), 92, 20)
	if !strings.Contains(got, "drain") {
		t.Fatal("a refused command was hidden")
	}
	for _, line := range strings.Split(got, "\n") {
		if strings.Contains(line, "drain") && !strings.Contains(line, "refused · leader") {
			t.Errorf("the refusal does not say why: %q", line)
		}
	}
}

// The cursor cannot land on a heading.
func TestTheCursorSkipsGroupHeadings(t *testing.T) {
	p := demoPalette()
	render(p, 92, 20)

	first, ok := p.Selected()
	if !ok {
		t.Fatal("nothing is selected")
	}
	if first.Label != "restart" {
		t.Errorf("the palette opened on %q, want the first command", first.Label)
	}

	// Down twice reaches the next group's first command, over its heading.
	p.Move(1)
	render(p, 92, 20)
	p.Move(1)
	render(p, 92, 20)
	got, _ := p.Selected()
	if got.Label != "open a shell here" {
		t.Errorf("two moves reached %q, want the next group's first command", got.Label)
	}
}

// Typing flattens the list and ranks it, and the matched letters are marked so
// the ORDER is checkable rather than trusted.
func TestTypingRanksAndFlattens(t *testing.T) {
	p := demoPalette()
	p.Query = "res"
	got := render(p, 92, 20)

	if strings.Contains(got, "2 tasks on 2 nodes") {
		t.Error("the group headings survived a query")
	}
	if !strings.Contains(got, "restart") {
		t.Errorf("the obvious match is missing:\n%s", got)
	}
	if n := p.Matches(); n == 0 || n == p.Total() {
		t.Errorf("a query matched %d of %d — it filtered nothing or everything", n, p.Total())
	}
	// The group each match came from is a column now.
	if !strings.Contains(got, "service · api_api") {
		t.Errorf("a flat row does not say which group it came from:\n%s", got)
	}
}

// Ranked, never FILTERED by group: the groups are an order of preference, and
// a reader who typed the name of a global command should find it.
func TestAQueryReachesEveryGroup(t *testing.T) {
	p := demoPalette()
	p.Query = "prune"
	render(p, 92, 20)

	got, ok := p.Selected()
	if !ok {
		t.Fatal("nothing matched")
	}
	if got.Label != "prune unused images" {
		t.Errorf("selected %q, want the global command", got.Label)
	}
}

func TestAQueryMatchingNothing(t *testing.T) {
	p := demoPalette()
	p.Query = "zzzzz"
	got := render(p, 92, 20)
	if !strings.Contains(got, "nothing matches") {
		t.Errorf("a query matching nothing says nothing:\n%s", got)
	}
	if _, ok := p.Selected(); ok {
		t.Error("something was selected out of no matches")
	}
}

// Selected reports false on a heading rather than a half-built item.
func TestSelectedIsFalseOnAHeading(t *testing.T) {
	p := demoPalette()
	render(p, 92, 20)
	p.Select(0) // the first heading
	render(p, 92, 20)
	if got, ok := p.Selected(); ok && got.Label == "" {
		t.Error("Selected returned an empty item as if it were real")
	}
}

func TestPaletteInNoRoom(t *testing.T) {
	p := demoPalette()
	c := NewCanvas(10, 3)
	p.Draw(c, Rect{X: 0, Y: 0, W: 10, H: 3})
	// Not panicking is the assertion; a palette in a tiny terminal is an
	// ordinary state.
}

// The rules are rules. Chrome.Divider is the SPACE between panes, and filling a
// row with it draws a blank that looks like a layout bug.
func TestThePaletteRulesAreDrawn(t *testing.T) {
	got := render(demoPalette(), 92, 20)
	rules := 0
	for _, line := range strings.Split(got, "\n") {
		if strings.Count(line, "─") > 20 && !strings.Contains(line, "┌") && !strings.Contains(line, "└") {
			rules++
		}
	}
	if rules != 2 {
		t.Errorf("found %d rules inside the box, want 2 (under the query, above the footer):\n%s", rules, got)
	}
}
