package comp

import (
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/theme"
)

func ruled(w, h int, rl Rule, ch theme.Chrome) string {
	c := NewCanvas(w, h).WithChrome(ch)
	rl.Draw(c, Rect{X: 0, Y: 0, W: w, H: h}, Region("rule"))
	return c.String()
}

// The character comes from the chrome, so a tool that changes its box set
// changes its rules with it. This is the whole reason the component exists:
// five of the seven hand-rolled sites wrote "─" and got it wrong on a box set
// that does not have one.
func TestTheRuleTakesItsCharacterFromTheChrome(t *testing.T) {
	for _, tc := range []struct {
		name string
		box  theme.BoxSet
		want string
	}{
		{"light", theme.LightBox, "─"},
		{"heavy", theme.HeavyBox, "━"},
		{"ascii", theme.ASCIIBox, "-"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			ch := theme.DefaultChrome
			ch.Box = tc.box
			got := strings.TrimRight(ruled(6, 1, Rule{}, ch), "\n")
			if want := strings.Repeat(tc.want, 6); got != want {
				t.Errorf("drew %q, want %q", got, want)
			}
		})
	}
}

// A caller may override it, for the rare rule that means something different
// from the others on the screen.
func TestARuleCanBeToldWhatToDrawWith(t *testing.T) {
	got := strings.TrimRight(ruled(4, 1, Rule{Rune: "="}, theme.DefaultChrome), "\n")
	if got != "====" {
		t.Errorf("drew %q", got)
	}
}

// One row, the way Bar does. A band taller than a line is a layout with room to
// spare, not a request for a block of ─ that reads as a rendering fault.
func TestARuleIsOneRowHoweverTallItsBandIs(t *testing.T) {
	out := ruled(4, 3, Rule{}, theme.DefaultChrome)

	lines := strings.Split(out, "\n")
	if strings.TrimSpace(lines[0]) == "" {
		t.Errorf("the first row is blank; the rule was not drawn: %q", out)
	}
	// Whatever the canvas does with trailing blanks, no row below the first may
	// carry the character.
	for i, line := range lines[1:] {
		if strings.Contains(line, "─") {
			t.Errorf("row %d is %q; a rule is a line, not a fill", i+1, line)
		}
	}
}

// An empty rect draws nothing rather than panicking, like every other
// component handed a band that a narrow terminal left no room for.
func TestARuleInNoSpaceDrawsNothing(t *testing.T) {
	c := NewCanvas(4, 1)
	Rule{}.Draw(c, Rect{}, Region("rule"))
	if got := strings.TrimSpace(c.String()); got != "" {
		t.Errorf("drew %q into no space", got)
	}
}

// The rule owns its cells, so a click on one resolves to whatever drew it
// rather than to nothing.
func TestTheRuleOwnsItsCells(t *testing.T) {
	c := NewCanvas(6, 1)
	Rule{}.Draw(c, Rect{X: 1, Y: 0, W: 4, H: 1}, Region("header"))

	if got := c.OwnerAt(2, 0); got.Name != "header" {
		t.Errorf("cell 2 is owned by %v", got)
	}
	if got := c.OwnerAt(0, 0); got.Name == "header" {
		t.Error("the rule claimed a cell outside its rect")
	}
}
