package paint_test

import (
	"image/color"
	"testing"

	"github.com/richarddavenport/tuikit/paint"
)

var ramp = paint.Ramp{
	From: color.RGBA{0, 0, 0, 255},
	To:   color.RGBA{100, 200, 255, 255},
}

func TestRampEndsAndMiddle(t *testing.T) {
	if got := ramp.At(0); got != ramp.From {
		t.Errorf("At(0) = %v, want %v", got, ramp.From)
	}
	if got := ramp.At(1); got != ramp.To {
		t.Errorf("At(1) = %v, want %v", got, ramp.To)
	}
	if got := ramp.At(0.5); got.G != 100 {
		t.Errorf("At(0.5).G = %d, want 100", got.G)
	}
}

// TestRampClamps: a caller a pixel off the end should get the end colour, not
// the opposite one.
func TestRampClamps(t *testing.T) {
	if got := ramp.At(-3); got != ramp.From {
		t.Errorf("At(-3) = %v, want the From end", got)
	}
	if got := ramp.At(4); got != ramp.To {
		t.Errorf("At(4) = %v, want the To end", got)
	}
}

func TestPanelSizeAndOpacity(t *testing.T) {
	img := paint.Panel{W: 40, H: 20, Ramp: ramp}.Image()
	if b := img.Bounds(); b.Dx() != 40 || b.Dy() != 20 {
		t.Fatalf("bounds %v, want 40x20", b)
	}
	// No radius: every pixel is fully opaque, corners included.
	for _, p := range [][2]int{{0, 0}, {39, 0}, {0, 19}, {39, 19}, {20, 10}} {
		if a := img.RGBAAt(p[0], p[1]).A; a != 255 {
			t.Errorf("alpha at %v = %d, want 255 with no radius", p, a)
		}
	}
}

// TestRoundedCornersAreTransparent is what makes the kitty path worth having:
// the corner is genuinely absent, so it blends with whatever is behind it.
func TestRoundedCornersAreTransparent(t *testing.T) {
	img := paint.Panel{W: 40, H: 20, Ramp: ramp, Radius: 6}.Image()
	if a := img.RGBAAt(0, 0).A; a != 0 {
		t.Errorf("corner alpha = %d, want 0", a)
	}
	if a := img.RGBAAt(20, 10).A; a != 255 {
		t.Errorf("centre alpha = %d, want 255", a)
	}
	// The arc is antialiased rather than stepped, so somewhere along it there
	// is a partial pixel.
	partial := false
	for y := 0; y < 20 && !partial; y++ {
		for x := 0; x < 40; x++ {
			if a := img.RGBAAt(x, y).A; a > 0 && a < 255 {
				partial = true
				break
			}
		}
	}
	if !partial {
		t.Error("the corner arc has no partial pixels; it is aliased")
	}
}

func TestBarsAreDrawn(t *testing.T) {
	img := paint.Panel{
		W: 60, H: 30, Ramp: ramp, Radius: 4,
		Bars:     []float64{1, 0, 1},
		BarColor: color.RGBA{255, 0, 0, 255},
	}.Image()

	red := 0
	for y := 0; y < 30; y++ {
		for x := 0; x < 60; x++ {
			if img.RGBAAt(x, y) == (color.RGBA{255, 0, 0, 255}) {
				red++
			}
		}
	}
	if red == 0 {
		t.Fatal("no bars drawn")
	}
	// A zero-height bar draws nothing, so two full bars and one empty should
	// be well under half the panel.
	if red > 60*30/2 {
		t.Errorf("%d bar pixels; the empty bar was drawn", red)
	}
}

// TestBarsClamp: data outside 0..1 must not draw outside the panel.
func TestBarsClamp(t *testing.T) {
	paint.Panel{W: 20, H: 10, Ramp: ramp, Bars: []float64{-5, 12}}.Image()
	// Not panicking IS the assertion; an unclamped bar indexes out of range.
}

// TestMoreBarsThanPixels: a chart that silently vanishes when the panel gets
// narrow is worse than one that gets chunky.
func TestMoreBarsThanPixels(t *testing.T) {
	bars := make([]float64, 50)
	for i := range bars {
		bars[i] = 1
	}
	img := paint.Panel{W: 12, H: 12, Ramp: ramp, Bars: bars, BarColor: color.RGBA{255, 0, 0, 255}}.Image()
	for y := 0; y < 12; y++ {
		for x := 0; x < 12; x++ {
			if img.RGBAAt(x, y) == (color.RGBA{255, 0, 0, 255}) {
				return
			}
		}
	}
	t.Error("50 bars in 12 pixels drew nothing at all")
}

// TestFlatten is Sixel's second concession made concrete.
func TestFlatten(t *testing.T) {
	src := paint.Panel{W: 20, H: 20, Ramp: ramp, Radius: 8}.Image()
	bg := color.RGBA{30, 30, 40, 255}
	out := paint.Flatten(src, bg)

	if got := out.RGBAAt(0, 0); got != bg {
		t.Errorf("transparent corner flattened to %v, want the background %v", got, bg)
	}
	for y := 0; y < 20; y++ {
		for x := 0; x < 20; x++ {
			if a := out.RGBAAt(x, y).A; a != 255 {
				t.Fatalf("alpha %d survived flattening at %d,%d; Sixel has no alpha", a, x, y)
			}
		}
	}
}
