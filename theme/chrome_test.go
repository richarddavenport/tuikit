package theme

import (
	"strings"
	"testing"
)

// Every character the default chrome draws is in the default glyph set, or the
// interface prints something a guard would reject if a tool had written it.
func TestTheDefaultChromeIsInTheDefaultGlyphSet(t *testing.T) {
	for _, r := range DefaultChrome.Glyphs() {
		if !DefaultGlyphs.Printable(r) {
			t.Errorf("the default chrome draws %q (%U), which is not in the default glyph set", r, r)
		}
	}
}

// The other box sets are NOT in the default set, and that is the point:
// choosing one is adding its characters deliberately, rather than discovering
// four replacement boxes on somebody else's terminal.
func TestTheOtherBoxSetsNeedTheirCharactersAdded(t *testing.T) {
	for name, box := range map[string]BoxSet{"rounded": RoundedBox, "heavy": HeavyBox} {
		var missing int
		for _, r := range DefaultChrome.With(box).Glyphs() {
			if !DefaultGlyphs.Printable(r) {
				missing++
			}
		}
		if missing == 0 {
			t.Errorf("%s box needs nothing added, so the coupling proves nothing", name)
		}
	}
}

// ASCII needs nothing added, because it is the fallback for a font that has
// nothing.
func TestTheASCIIBoxNeedsNothing(t *testing.T) {
	for _, r := range DefaultChrome.With(ASCIIBox).Glyphs() {
		if r < 128 {
			continue
		}
		if !DefaultGlyphs.Printable(r) {
			t.Errorf("the ASCII fallback needs %q added", r)
		}
	}
}

// Glyphs walks the fields rather than listing them, so a field added to Chrome
// cannot be forgotten — which is the failure this type exists to prevent, one
// level up.
func TestGlyphsCoversEveryField(t *testing.T) {
	c := Chrome{
		Box:     BoxSet{"1", "2", "3", "4", "5", "6", "7", "8"},
		Divider: "9", VDivider: "a",
		Ellipsis: "b", Separator: "c", ScrollUp: "d", ScrollDown: "e",
		ChevronLeft: "f", ChevronRight: "g",
	}
	got := string(c.Glyphs())
	for _, want := range strings.Split("123456789abcdefg", "") {
		if !strings.Contains(got, want) {
			t.Errorf("%q is not reported by Glyphs: got %q", want, got)
		}
	}
}

func TestGlyphsDoesNotRepeatItself(t *testing.T) {
	seen := map[rune]bool{}
	for _, r := range DefaultChrome.Glyphs() {
		if seen[r] {
			t.Errorf("%q is reported twice", r)
		}
		seen[r] = true
	}
}
