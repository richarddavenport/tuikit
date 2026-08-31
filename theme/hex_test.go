package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func lipglossColor(s string) lipgloss.Color { return lipgloss.Color(s) }

// The formula is the point: a table of hand-copied values would be a table of
// chances to copy one wrong. These are the anchors that prove it is the right
// formula — one from each of the three regions, plus every role the interface
// actually uses.
func TestHexFollowsTheAnsiPalette(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"0", "#000000"},   // the base sixteen, which are a table
		{"15", "#ffffff"},  //
		{"16", "#000000"},  // the cube's first corner
		{"231", "#ffffff"}, // and its last
		{"232", "#080808"}, // the grey ramp's ends
		{"255", "#eeeeee"}, //
		// Every role, so a change to one is a change to a value someone can see.
		{"205", "#ff5faf"}, // Accent
		{"241", "#626262"}, // Muted
		{"240", "#585858"}, // Border
		{"34", "#00af00"},  // Success
		{"214", "#ffaf00"}, // Pending
		{"196", "#ff0000"}, // Danger
		{"203", "#ff5f5f"}, // Stderr
		{"229", "#ffffaf"}, // SelectionFG
		{"57", "#5f00ff"},  // SelectionBG
	} {
		if got := Hex(lipglossColor(tc.in)); got != tc.want {
			t.Errorf("Hex(%s) = %s, want %s", tc.in, got, tc.want)
		}
	}
}

// Anything that is not a palette index is black rather than a panic: this runs
// in a generator, and a bad value should produce a visibly wrong swatch, not a
// crashed build.
func TestHexIsTotal(t *testing.T) {
	for _, in := range []string{"", "nope", "-1", "256", "#ff00ff"} {
		if got := Hex(lipglossColor(in)); got != "#000000" {
			t.Errorf("Hex(%q) = %s, want black", in, got)
		}
	}
}

// Every role must convert, or the design system draws a black square where a
// colour should be.
func TestEveryRoleHasAColour(t *testing.T) {
	for _, r := range Default.Roles() {
		if hex := Hex(r.Color); hex == "#000000" {
			t.Errorf("%s (%s) has no colour", r.Name, r.Color)
		}
		if r.Why == "" {
			t.Errorf("%s has no reason for existing — that is what makes it a role", r.Name)
		}
	}
}
