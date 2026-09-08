package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

const footer Name = "footer"

func seg(text string) []Segment { return []Segment{{Text: text}} }

func TestABarPutsContentAtEachEnd(t *testing.T) {
	c := NewCanvas(30, 1)
	Bar{Left: seg("democtl"), Right: seg("09:14:03")}.Draw(c, c.Bounds(), Region(footer))

	if got, want := c.String(), "democtl"+strings.Repeat(" ", 15)+"09:14:03"; got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

// The left side wins. A bar too narrow for both is a bar whose right-hand side
// was never load-bearing.
func TestTheRightSideGoesWhenTheyWouldCollide(t *testing.T) {
	c := NewCanvas(14, 1)
	Bar{Left: seg("a long left side"), Right: seg("09:14:03")}.Draw(c, c.Bounds(), Region(footer))

	got := c.String()
	if strings.Contains(got, "09:14") {
		t.Errorf("the right side survived a collision: %q", got)
	}
	if Width(got) > 14 {
		t.Errorf("the bar drew %d columns: %q", Width(got), got)
	}
}

// the deploy tool drops the version unless the hints still get room, because a
// version squeezing the keys you can press is decoration winning over content.
func TestMinLeftProtectsTheLeftSideFromDecoration(t *testing.T) {
	c := NewCanvas(30, 1)
	Bar{Left: seg("q quit"), Right: seg("v2.4.1"), MinLeft: 24}.Draw(c, c.Bounds(), Region(footer))

	if got := c.String(); strings.Contains(got, "v2.4.1") {
		t.Errorf("the right side drew with only 24 columns for the left: %q", got)
	}

	// With room for both, it comes back.
	wide := NewCanvas(40, 1)
	Bar{Left: seg("q quit"), Right: seg("v2.4.1"), MinLeft: 24}.Draw(wide, wide.Bounds(), Region(footer))
	if got := wide.String(); !strings.Contains(got, "v2.4.1") {
		t.Errorf("the right side is missing with room for both: %q", got)
	}
}

// Without MinLeft, only a collision drops it — democtl's rule.
func TestWithoutMinLeftOnlyACollisionDropsTheRightSide(t *testing.T) {
	c := NewCanvas(15, 1)
	Bar{Left: seg("q quit"), Right: seg("v2.4.1")}.Draw(c, c.Bounds(), Region(footer))

	if got := c.String(); !strings.Contains(got, "v2.4.1") {
		t.Errorf("the right side went without a collision: %q", got)
	}
}

func TestEachEndCanBeSeveralStyles(t *testing.T) {
	forceColour()
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	c := NewCanvas(40, 1)
	Bar{Left: []Segment{
		{Text: "democtl", Style: &accent},
		{Text: "  a tuikit example", Style: &muted},
	}}.Draw(c, c.Bounds(), Region(footer))

	got := c.String()
	if !strings.Contains(got, "38;5;205") || !strings.Contains(got, "38;5;241") {
		t.Errorf("the segments do not keep their own styles: %q", got)
	}
}

// One list, so a footer and a context menu cannot describe the same action
// differently.
func TestHintsReadTheWayAllFourToolsWriteThem(t *testing.T) {
	got := Hints(Hint{"j/k", "move"}, Hint{"enter", "connect"}, Hint{"q", "quit"})
	if want := "j/k move · enter connect · q quit"; got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

func TestAHintCanBeAllLabelOrAllKey(t *testing.T) {
	if got, want := Hints(Hint{Label: "connecting"}, Hint{Key: "esc", Label: "cancel"}), "connecting · esc cancel"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestKeyHintsClipsRatherThanOverflowing(t *testing.T) {
	c := NewCanvas(12, 1)
	KeyHints(c, c.Bounds(), Region(footer), nil,
		Hint{"j/k", "move"}, Hint{"enter", "connect"}, Hint{"q", "quit"})

	if got := Width(c.String()); got > 12 {
		t.Errorf("the hints drew %d columns into 12: %q", got, c.String())
	}
}
