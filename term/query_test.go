package term_test

import (
	"image/color"
	"testing"

	"github.com/richarddavenport/tuikit/term"
)

func TestParseRGB(t *testing.T) {
	for _, tc := range []struct {
		name  string
		reply string
		want  color.RGBA
	}{
		// Components are a fraction of their own maximum, not a byte. These
		// three are all pure red, and truncating to the first two digits would
		// get the last one badly wrong.
		{"four digits", "\x1b]11;rgb:ffff/0000/0000\x1b\\", color.RGBA{255, 0, 0, 255}},
		{"two digits", "\x1b]11;rgb:ff/00/00\x07", color.RGBA{255, 0, 0, 255}},
		{"one digit", "\x1b]11;rgb:f/0/0\x07", color.RGBA{255, 0, 0, 255}},
		{"a real background", "\x1b]11;rgb:1a1a/1b1b/2626\x1b\\", color.RGBA{26, 27, 38, 255}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, ok := term.ParseRGB(tc.reply)
			if !ok {
				t.Fatalf("ParseRGB(%q) found nothing", tc.reply)
			}
			if got != tc.want {
				t.Errorf("got %v, want %v", got, tc.want)
			}
		})
	}
}

// TestParseRGBSaysWhenItFound_Nothing is the distinction ParseBackground alone
// cannot make: a terminal that answered with the default colour and a terminal
// that did not answer are the same value.
func TestParseRGBSaysWhenItFoundNothing(t *testing.T) {
	for _, bad := range []string{"", "\x1b[?62;4c", "rgb:", "rgb:zz/00/00", "rgb:00/00", "rgb:12345/0/0"} {
		if _, ok := term.ParseRGB(bad); ok {
			t.Errorf("ParseRGB(%q) claimed to find a colour", bad)
		}
	}
	// And the default IS a legitimate answer, distinguishable from silence.
	reply := "\x1b]11;rgb:1a1a/1a1a/1a1a\x1b\\"
	if c, ok := term.ParseRGB(reply); !ok || c != term.DefaultBackground {
		t.Errorf("a terminal answering exactly the default was not recognised: %v %v", c, ok)
	}
}

func TestParseBackgroundFallsBack(t *testing.T) {
	if got := term.ParseBackground("nothing here"); got != term.DefaultBackground {
		t.Errorf("got %v, want the default", got)
	}
}

func TestParseColors(t *testing.T) {
	// Two answers, out of order, plus an index the terminal declined — which
	// must be absent rather than present and wrong.
	// Four-digit components as a terminal actually sends them: each byte
	// repeated, so ff is ffff rather than ff00. (ff00 is 0.996 of full, not
	// full — which is correct arithmetic and not what any terminal means.)
	reply := "\x1b]4;13;rgb:ffff/5f5f/afaf\x1b\\" +
		"\x1b]4;5;rgb:5f5f/0000/ffff\x1b\\" +
		"\x1b[?62;4c"

	got := term.ParseColors(reply)
	if len(got) != 2 {
		t.Fatalf("got %d colours, want 2: %v", len(got), got)
	}
	if want := (color.RGBA{95, 0, 255, 255}); got[5] != want {
		t.Errorf("index 5 = %v, want %v", got[5], want)
	}
	if want := (color.RGBA{255, 95, 175, 255}); got[13] != want {
		t.Errorf("index 13 = %v, want %v", got[13], want)
	}
	if _, ok := got[9]; ok {
		t.Error("an index the terminal never mentioned came back with a value")
	}
}

func TestParseColorsFindsNothing(t *testing.T) {
	if got := term.ParseColors("\x1b[?62;4c"); got != nil {
		t.Errorf("got %v, want nil", got)
	}
}
