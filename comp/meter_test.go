package comp_test

import (
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/term"
)

const meterID = comp.Name("meter")
const trackID = comp.Name("track")

func drawMeter(t *testing.T, m comp.Meter, w int) (*comp.Canvas, string) {
	t.Helper()
	m.Track = trackID
	c := comp.NewCanvas(w, 1)
	m.Draw(c, comp.Rect{X: 0, Y: 0, W: w, H: 1}, comp.Region(meterID))
	return c, c.String()
}

func TestMeterFillsProportionally(t *testing.T) {
	for _, tc := range []struct {
		value float64
		want  string
	}{
		{0, "[··········]"},
		{0.5, "[─────·····]"},
		{1, "[──────────]"},
		{0.25, "[───·······]"}, // 2.5 rounds up
	} {
		if _, got := drawMeter(t, comp.Meter{Value: tc.value}, 12); got != tc.want {
			t.Errorf("Value %v drew %q, want %q", tc.value, got, tc.want)
		}
	}
}

// TestMeterClamps: a bar that overflows its brackets hides a caller's bad
// arithmetic instead of showing it.
func TestMeterClamps(t *testing.T) {
	if _, got := drawMeter(t, comp.Meter{Value: 5}, 12); got != "[──────────]" {
		t.Errorf("Value 5 drew %q", got)
	}
	if _, got := drawMeter(t, comp.Meter{Value: -2}, 12); got != "[··········]" {
		t.Errorf("Value -2 drew %q", got)
	}
}

func TestMeterLabel(t *testing.T) {
	_, got := drawMeter(t, comp.Meter{Value: 0.5, Label: "2 of 5"}, 20)
	// 20 columns less "2 of 5" and its gap leaves 13 for the bar, so 11 inside
	// the brackets and half of 11 rounds to 6.
	if want := "[──────·····] 2 of 5"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// TestMeterUsesNoBlockElements. They are the obvious choice for a progress bar
// and they are excluded: a font without them draws a row of empty boxes.
func TestMeterUsesNoBlockElements(t *testing.T) {
	_, got := drawMeter(t, comp.Meter{Value: 0.5}, 20)
	for _, r := range got {
		if r >= 0x2580 && r <= 0x259F {
			t.Errorf("meter drew block element %U", r)
		}
	}
}

func TestMeterTooNarrowDrawsNothing(t *testing.T) {
	for _, w := range []int{0, 1, 2} {
		if _, got := drawMeter(t, comp.Meter{Value: 0.5}, w); strings.TrimSpace(got) != "" {
			t.Errorf("width %d drew %q", w, got)
		}
	}
}

// TestMeterTrackIsItsOwnRegion, which is what lets the pixel layer find the bar
// without knowing where the meter ended up.
func TestMeterTrackIsItsOwnRegion(t *testing.T) {
	c, _ := drawMeter(t, comp.Meter{Value: 0.5}, 12)
	r, ok := c.Region(comp.Region(trackID))
	if !ok {
		t.Fatal("the track is not a region")
	}
	if r.X != 1 || r.W != 10 {
		t.Errorf("track is %v, want the ten cells inside the brackets", r)
	}
	// The brackets belong to the meter, not the track.
	if owner := c.OwnerAt(0, 0); owner.Name != meterID {
		t.Errorf("the opening bracket is owned by %v", owner)
	}
}

// TestMeterWithARampIsUnchangedWithoutGraphics is the property that makes the
// pixel layer safe to add to any component: offering a picture costs nothing on
// a terminal that cannot take one.
func TestMeterWithARampIsUnchangedWithoutGraphics(t *testing.T) {
	_, plain := drawMeter(t, comp.Meter{Value: 0.5}, 20)
	_, offered := drawMeter(t, comp.Meter{Value: 0.5, Pixels: true}, 20)
	if plain != offered {
		t.Errorf("offering a picture changed the frame:\n got %q\nwant %q", offered, plain)
	}
}

// TestMeterRasteriserIsNotCalledWithoutGraphics. The draw function is a
// function precisely so a terminal that cannot show a picture never pays to
// make one.
func TestMeterRasteriserIsNotCalledWithoutGraphics(t *testing.T) {
	c := comp.NewCanvas(20, 1)
	called := false
	m := comp.Meter{Value: 0.5, Track: trackID, Pixels: true}
	m.Draw(c, comp.Rect{X: 0, Y: 0, W: 20, H: 1}, comp.Region(meterID))
	// Reaching into the canvas is not possible from here, so this asserts the
	// observable half: no escape sequence, therefore no image was made.
	if strings.Contains(c.String(), "\x1b") {
		called = true
	}
	if called {
		t.Error("a picture was encoded on a canvas with no graphics")
	}
}

func TestMeterOffersAPictureWhenItCan(t *testing.T) {
	c := comp.NewCanvas(20, 1).WithGraphics(comp.Pixels{Mode: term.Kitty, CellW: 10, CellH: 20})
	m := comp.Meter{Value: 0.5, Track: trackID, Pixels: true}
	m.Draw(c, comp.Rect{X: 0, Y: 0, W: 20, H: 1}, comp.Region(meterID))

	frame := c.String()
	if !strings.Contains(frame, "\x1b_G") {
		t.Error("no kitty image in the frame")
	}
	// The characters are still there: they were drawn first, and kitty
	// composites behind them.
	if !strings.Contains(frame, "[─────────·········]") {
		t.Error("the cell bar did not survive the picture")
	}
}
