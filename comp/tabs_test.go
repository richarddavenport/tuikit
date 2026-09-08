package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

const tabsName Name = "detail.tab"

func tabsOf(names ...string) []Tab {
	out := make([]Tab, len(names))
	for i, n := range names {
		out[i] = Tab{Name: n}
	}
	return out
}

// The chevrons go around the STRIP, because they say the strip cycles and which
// keys do it — a fact about the keymap that no styling can carry. Around the
// active tab they would only repeat what its color already says.
func TestTheChevronsWrapTheStripNotTheCurrentTab(t *testing.T) {
	c := NewCanvas(40, 1)
	Tabs{Tabs: tabsOf("Overview", "Config", "Events"), Active: 0}.
		Draw(c, c.Bounds(), tabsName)

	if got, want := c.String(), "‹ Overview · Config · Events ›"; got != want {
		t.Errorf("got  %q\nwant %q", got, want)
	}
}

// A tab can say how much is behind it.
func TestATabCanCarryACount(t *testing.T) {
	c := NewCanvas(40, 1)
	Tabs{Tabs: []Tab{{Name: "Edits", Count: 3}, {Name: "Log"}}}.
		Draw(c, c.Bounds(), tabsName)

	if got := c.String(); !strings.Contains(got, "Edits(3)") {
		t.Errorf("the count is missing: %q", got)
	}
	if got := c.String(); strings.Contains(got, "Log(") {
		t.Errorf("a tab with nothing behind it drew a count: %q", got)
	}
}

// Each tab owns its own region, so a click lands on the tab and not the strip.
func TestEachTabIsClickableSeparately(t *testing.T) {
	c := NewCanvas(40, 1)
	Tabs{Tabs: tabsOf("Overview", "Config", "Events"), Active: 0}.
		Draw(c, c.Bounds(), tabsName)

	for i, want := range []string{"Overview", "Config", "Events"} {
		r, ok := c.Region(Region(tabsName).At(i))
		if !ok {
			t.Errorf("tab %d (%s) was not drawn", i, want)
			continue
		}
		if got := c.OwnerAt(r.X+1, r.Y); got.Index != i {
			t.Errorf("the middle of tab %d is owned by %v", i, got)
		}
	}
}

// A strip you cannot operate should not look like one you can.
func TestTheCurrentTabLooksDifferentWhenFocused(t *testing.T) {
	forceColor()
	quiet := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))
	on := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	lit := lipgloss.NewStyle().Foreground(lipgloss.Color("229"))

	strip := Tabs{
		Tabs: tabsOf("Overview", "Config"), Active: 0,
		Style: &quiet, Selected: &on, FocusSelected: &lit,
	}
	away := NewCanvas(40, 1)
	strip.Draw(away, away.Bounds(), tabsName)

	strip.Focused = true
	here := NewCanvas(40, 1)
	strip.Draw(here, here.Bounds(), tabsName)

	if !strings.Contains(away.String(), "38;5;205") {
		t.Error("an unfocused strip does not mark its current tab")
	}
	if !strings.Contains(here.String(), "38;5;229") {
		t.Error("a focused strip marks its current tab the same way")
	}
}

// A strip wider than its pane is clipped by the canvas rather than wrapping or
// pushing the border out.
func TestAStripTooWideIsClipped(t *testing.T) {
	c := NewCanvas(12, 1)
	Tabs{Tabs: tabsOf("Overview", "Config", "Events")}.Draw(c, c.Bounds(), tabsName)

	if got := Width(c.String()); got > 12 {
		t.Errorf("the strip drew %d columns into 12: %q", got, c.String())
	}
}
