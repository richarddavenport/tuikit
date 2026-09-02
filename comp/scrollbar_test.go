package comp

import (
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/theme"
)

// bar draws a scrollbar and returns its rows top to bottom, as T for track and
// H for thumb.
func bar(s Scrollbar, h int) string {
	c := NewCanvas(1, h)
	s.Draw(c, Rect{X: 0, Y: 0, W: 1, H: h}, Region("bar"))
	var out strings.Builder
	for y := 0; y < h; y++ {
		cell, _ := c.CellAt(0, y)
		switch cell.Text {
		case theme.DefaultChrome.ScrollThumb:
			out.WriteByte('H')
		case theme.DefaultChrome.ScrollTrack:
			out.WriteByte('T')
		default:
			out.WriteByte('.')
		}
	}
	return out.String()
}

func TestTheThumbIsAProportion(t *testing.T) {
	// Ten rows of forty, so a quarter of a twelve-row bar: three.
	got := bar(Scrollbar{Total: 40, Shown: 10, Offset: 0}, 12)
	if want := "HHHTTTTTTTTT"; got != want {
		t.Errorf("got %s, want %s", got, want)
	}
}

func TestTheThumbIsAPosition(t *testing.T) {
	s := Scrollbar{Total: 40, Shown: 10}

	s.Offset = 0
	if got := bar(s, 12); !strings.HasPrefix(got, "H") {
		t.Errorf("at the top the thumb is at %s", got)
	}
	// The last row of the list must put the thumb on the last row of the bar.
	// "Nearly at the end" and "at the end" are different answers.
	s.Offset = 30
	if got := bar(s, 12); !strings.HasSuffix(got, "H") {
		t.Errorf("at the bottom the thumb is at %s", got)
	}
	s.Offset = 15
	got := bar(s, 12)
	if strings.HasPrefix(got, "H") || strings.HasSuffix(got, "H") {
		t.Errorf("halfway the thumb is at %s, want it in the middle", got)
	}
}

// A list of two hundred thousand in a ten-row pane gives a thumb of 0.0005
// rows. A scrollbar that vanishes exactly when the list is longest is worse
// than none.
func TestTheThumbIsNeverZeroRows(t *testing.T) {
	got := bar(Scrollbar{Total: 200000, Shown: 10, Offset: 0}, 10)
	if !strings.Contains(got, "H") {
		t.Errorf("the thumb disappeared: %s", got)
	}
	if strings.Count(got, "H") != 1 {
		t.Errorf("the thumb is %d rows, want 1", strings.Count(got, "H"))
	}
	// And it still reaches the bottom.
	if got := bar(Scrollbar{Total: 200000, Shown: 10, Offset: 199990}, 10); !strings.HasSuffix(got, "H") {
		t.Errorf("at the end of a huge list the thumb is at %s", got)
	}
}

// Nothing to scroll draws nothing. A full-height thumb says only that there is
// a scrollbar.
func TestNothingToScrollDrawsNothing(t *testing.T) {
	for _, s := range []Scrollbar{
		{Total: 5, Shown: 10, Offset: 0},
		{Total: 10, Shown: 10, Offset: 0},
		{Total: 0, Shown: 10},
		{Total: 10, Shown: 0},
	} {
		if got := bar(s, 6); strings.ContainsAny(got, "HT") {
			t.Errorf("%+v drew %s", s, got)
		}
	}
}

// An offset past the end is a caller's arithmetic, not a reason to draw
// outside the bar.
func TestAnOffsetOutOfRangeClamps(t *testing.T) {
	for _, off := range []int{-50, 999} {
		got := bar(Scrollbar{Total: 40, Shown: 10, Offset: off}, 12)
		if strings.Count(got, "H")+strings.Count(got, "T") != 12 {
			t.Errorf("offset %d drew %s", off, got)
		}
	}
}

func TestScrollbarInNoRoom(t *testing.T) {
	c := NewCanvas(4, 4)
	Scrollbar{Total: 40, Shown: 10}.Draw(c, Rect{}, Region("bar"))
	if strings.TrimSpace(c.String()) != "" {
		t.Error("an empty rect drew something")
	}
}

// No block elements: a font without them draws a scrollbar as a column of
// replacement boxes, which is what the glyph set exists to prevent.
func TestTheScrollbarUsesNoBlockElements(t *testing.T) {
	for _, s := range []string{theme.DefaultChrome.ScrollTrack, theme.DefaultChrome.ScrollThumb} {
		for _, r := range s {
			if r >= 0x2580 && r <= 0x259F {
				t.Errorf("the scrollbar uses block element %U", r)
			}
			if _, ok := theme.DefaultGlyphs[r]; !ok {
				t.Errorf("the scrollbar uses %q (%U), which is not in DefaultGlyphs", r, r)
			}
		}
	}
}
