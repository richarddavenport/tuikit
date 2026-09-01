package guard

import (
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/theme"
)

func TestChromeIsQuietWhenTheDefaultsAgree(t *testing.T) {
	rec := &recorder{}
	Chrome(rec, theme.DefaultChrome, theme.DefaultGlyphs)

	if len(rec.msgs) != 0 {
		t.Errorf("the defaults disagree with each other: %v", rec.msgs)
	}
}

// A guard that never fires is indistinguishable from a clean package.
func TestChromeFiresOnABoxSetTheToolHasNotAllowed(t *testing.T) {
	rec := &recorder{}
	Chrome(rec, theme.DefaultChrome.With(theme.RoundedBox), theme.DefaultGlyphs)

	if len(rec.msgs) == 0 {
		t.Fatal("a rounded box passed against a glyph set without ╭╮╰╯")
	}
	if !strings.Contains(rec.msgs[0], "╭") {
		t.Errorf("the failure does not name the character: %q", rec.msgs[0])
	}
}

// And it passes once they have been added deliberately, which is the whole
// shape of the decision.
func TestChromeIsQuietOnceTheCharactersAreAdded(t *testing.T) {
	glyphs := theme.DefaultGlyphs.With(
		'╭', "rounded box drawing", '╮', "rounded box drawing",
		'╰', "rounded box drawing", '╯', "rounded box drawing",
	)
	rec := &recorder{}
	Chrome(rec, theme.DefaultChrome.With(theme.RoundedBox), glyphs)

	if len(rec.msgs) != 0 {
		t.Errorf("still complaining after the characters were added: %v", rec.msgs)
	}
}
