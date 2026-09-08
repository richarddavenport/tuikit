package theme

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
)

// The palette is a closed set with a fixed reading order. A tool that recolors
// it still gets the same nine roles in the same order, because the order is
// what a palette page reads and a reshuffle would silently redraw it.
func TestRolesAreTheNineInOrder(t *testing.T) {
	want := []string{"Accent", "Muted", "Border", "Success", "Pending",
		"Danger", "Stderr", "SelectionFG", "SelectionBG"}
	got := Default.Roles()
	if len(got) != len(want) {
		t.Fatalf("got %d roles, want %d", len(got), len(want))
	}
	for i, name := range want {
		if got[i].Name != name {
			t.Errorf("role %d is %s, want %s", i, got[i].Name, name)
		}
	}
}

// Overriding a color keeps the role's reason, because the reason describes
// what the color is FOR. A tool that decides its accent is blue has not
// changed what an accent means.
func TestOverridingAColorKeepsTheReason(t *testing.T) {
	blue := Default
	blue.Accent = lipgloss.Color("4")

	before, after := Default.Roles()[0], blue.Roles()[0]
	if after.Color != lipgloss.Color("4") {
		t.Errorf("Accent = %s, want 4", after.Color)
	}
	if after.Why != before.Why {
		t.Errorf("the reason changed with the hue:\n got %q\nwant %q", after.Why, before.Why)
	}
	if Default.Accent != lipgloss.Color("13") {
		t.Errorf("overriding a copy mutated Default: Accent = %s", Default.Accent)
	}
}

// Extra roles come last and are visible to everything that walks the palette —
// the guard included, which is the point. A tenth meaning that the guard cannot
// see is a literal waiting to happen.
func TestExtraRolesAreWalkedLikeAnyOther(t *testing.T) {
	p := Default
	p.Extra = []Role{{"Info", lipgloss.Color("39"), "a note the tool wants to make"}}

	roles := p.Roles()
	last := roles[len(roles)-1]
	if len(roles) != 10 || last.Name != "Info" {
		t.Fatalf("Extra did not land last: got %d roles ending %s", len(roles), last.Name)
	}
	if len(Default.Roles()) != 9 {
		t.Error("Extra on a copy leaked into Default")
	}
}

// A GlyphSet is a map, and a map is a reference. A tool adding a glyph in place
// would be adding it to every other tool in the process — and to the guard that
// is meant to catch it.
func TestWithCopiesRatherThanMutates(t *testing.T) {
	extended := DefaultGlyphs.With('×', "multiplication sign — replica counts")

	if !extended.Printable('×') {
		t.Error("With did not add the glyph")
	}
	if DefaultGlyphs.Printable('×') {
		t.Error("With mutated DefaultGlyphs — every tool in the process now allows ×")
	}
	if len(extended) != len(DefaultGlyphs)+1 {
		t.Errorf("got %d glyphs, want %d", len(extended), len(DefaultGlyphs)+1)
	}
}

// An odd count or a wrong type is a typo at the call site, not a condition to
// handle at runtime.
func TestWithPanicsOnAMalformedCall(t *testing.T) {
	for _, tc := range []struct {
		name  string
		pairs []any
	}{
		{"odd count", []any{'×'}},
		{"string key", []any{"×", "a reason"}},
		{"non-string reason", []any{'×', 3}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			defer func() {
				if recover() == nil {
					t.Error("did not panic")
				}
			}()
			DefaultGlyphs.With(tc.pairs...)
		})
	}
}

// ASCII always, the spinner's Braille always, the allow-list by membership, and
// nothing else. The spinner is the exception worth a test because it is the one
// range no string literal in the tool will ever contain.
func TestPrintable(t *testing.T) {
	for _, tc := range []struct {
		r    rune
		want bool
		why  string
	}{
		{'a', true, "ASCII"},
		{' ', true, "ASCII"},
		{'─', true, "in the allow-list"},
		{'✓', true, "in the allow-list"},
		{0x2800, true, "the spinner's first Braille cell"},
		{0x28FF, true, "the spinner's last Braille cell"},
		{'█', true, "a block element — CP437, and Sparkline draws with it"},
		{'▁', true, "the lowest eighth, the other end of the same ramp"},
		{'▓', false, "a shade, which nothing draws with"},
		{'▲', false, "a geometric shape beyond the plain bullet"},
		{'×', false, "not in the set until a tool adds it"},
		{'🚀', false, "emoji"},
	} {
		if got := DefaultGlyphs.Printable(tc.r); got != tc.want {
			t.Errorf("Printable(%q) = %v, want %v — %s", tc.r, got, tc.want, tc.why)
		}
	}
}

// A glyph without a reason is a glyph nobody decided on.
func TestEveryDefaultGlyphSaysWhyItIsSafe(t *testing.T) {
	for r, why := range DefaultGlyphs {
		if why == "" {
			t.Errorf("%q (U+%04X) is allowed but unexplained", r, r)
		}
	}
}
