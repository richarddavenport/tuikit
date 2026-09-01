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

// Anything that is neither a palette index nor a hex value is black rather than
// a panic: this runs in a generator, and a bad value should produce a visibly
// wrong swatch, not a crashed build.
func TestHexIsTotal(t *testing.T) {
	for _, in := range []string{"", "nope", "-1", "256"} {
		if got := Hex(lipglossColor(in)); got != "#000000" {
			t.Errorf("Hex(%q) = %s, want black", in, got)
		}
	}
}

// A tool that named its roles in hex has already answered the question, so Hex
// hands it back rather than pretending not to understand. azctl's palette is
// hex pairs, and a design system that drew nine black squares for it would be
// describing a tool that does not exist.
func TestHexPassesThroughAHexValue(t *testing.T) {
	for in, want := range map[string]string{
		"#ff00ff": "#ff00ff",
		"#FF00FF": "#ff00ff",
		"#d2a8ff": "#d2a8ff",
		"#fff":    "#ffffff",
	} {
		if got := Hex(lipglossColor(in)); got != want {
			t.Errorf("Hex(%q) = %s, want %s", in, got, want)
		}
	}
}

// An adaptive colour answers with its dark value, because everything that
// renders a palette outside a terminal draws on a dark ground — that is what
// the frame was captured for, and recolouring it would report a tool that does
// not exist.
func TestAnAdaptiveColourAnswersDark(t *testing.T) {
	c := lipgloss.AdaptiveColor{Light: "#1f2328", Dark: "#e6edf3"}
	if got := Hex(c); got != "#e6edf3" {
		t.Errorf("Hex(adaptive) = %s, want the dark value", got)
	}
	if got := Value(c); got != "#1f2328 / #e6edf3" {
		t.Errorf("Value(adaptive) = %q, want both", got)
	}
}

// Value is what the tool DECLARED, which is what a design system page should
// say alongside what it resolves to: "205" tells a reader the palette follows
// their terminal's own scheme, and "#ff5faf" tells them it does not.
func TestValueIsWhatWasWritten(t *testing.T) {
	if got := Value(lipglossColor("205")); got != "205" {
		t.Errorf("Value = %q", got)
	}
	if got := Value(lipglossColor("#d2a8ff")); got != "#d2a8ff" {
		t.Errorf("Value = %q", got)
	}
	if got := Value(nil); got != "" {
		t.Errorf("Value(nil) = %q", got)
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
