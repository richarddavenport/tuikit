// Package paint draws the pictures a terminal graphics protocol carries.
//
// It is deliberately not a drawing library. The shape list is short — a
// rounded rectangle, a linear ramp, a soft edge, a row of bars — because those
// are the things a terminal draws badly and pixels draw well. Everything a
// terminal already draws well stays as characters, which is what keeps this
// package small enough to read and keeps the fallback honest.
//
// # Why there is no text here
//
// The obvious next thing is a heading in a real typeface, and it is deliberately
// absent. Sixel is opaque: it replaces the cells it covers, so any text inside
// the picture's rectangle must be IN the picture, which means an embedded font,
// which means a licence decision and a real number on the binary — and it means
// the same words render differently depending on which terminal the reader has.
//
// Keeping text as cells avoids all of it. The picture is a band; the words live
// above and below it in the reader's own font at their own size, legible at any
// terminal width, selectable, and identical on all four terminals. See issue 30.
package paint

import (
	"image"
	"image/color"
	"math"
)

// Ramp is a linear gradient between two colours.
//
// Two stops rather than an arbitrary list, because the handoff's design uses
// two and a list would be a general answer to a question nobody has asked yet.
type Ramp struct{ From, To color.RGBA }

// At samples the ramp. t outside 0..1 clamps rather than wrapping — a caller
// that is off by a pixel at the edge should get the edge colour, not the
// opposite one.
func (r Ramp) At(t float64) color.RGBA {
	t = math.Max(0, math.Min(1, t))
	lerp := func(a, b uint8) uint8 { return uint8(float64(a) + (float64(b)-float64(a))*t) }
	return color.RGBA{lerp(r.From.R, r.To.R), lerp(r.From.G, r.To.G), lerp(r.From.B, r.To.B), lerp(r.From.A, r.To.A)}
}

// Panel is a rounded rectangle filled with a horizontal ramp.
//
// The alpha channel is real: corners and edges fade to transparent rather than
// to a background colour. A terminal that composites (kitty) gets that
// directly; one that does not (Sixel) bakes it against a known background with
// [Flatten], which is why the background is the caller's business and not this
// package's.
type Panel struct {
	W, H int
	Ramp Ramp
	// Radius is the corner radius in pixels. Zero is a plain rectangle.
	Radius int
	// Bars, if set, draws a chart inside the panel: one bar per value, each
	// 0..1 of the usable height. The panel answers a question rather than
	// merely decorating, which is the only reason a picture earns its cells.
	Bars []float64
	// BarColor is the bars' colour. Zero value uses the ramp's end, which
	// reads as "the same family, one step brighter".
	BarColor color.RGBA
}

// Image draws the panel.
func (p Panel) Image() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, p.W, p.H))
	if p.W <= 0 || p.H <= 0 {
		return img
	}
	for y := 0; y < p.H; y++ {
		for x := 0; x < p.W; x++ {
			c := p.Ramp.At(float64(x) / float64(p.W-1))
			// coverage is antialiasing and rounding in one number: how much of
			// this pixel is inside the rounded rectangle at all.
			c.A = uint8(float64(c.A) * p.coverage(float64(x)+0.5, float64(y)+0.5))
			img.SetRGBA(x, y, c)
		}
	}
	p.bars(img)
	return img
}

// coverage is 1 inside the shape, 0 outside, and a fraction across one pixel of
// the corner arcs — which is the whole antialiasing story for a shape made only
// of straight edges and quarter circles.
func (p Panel) coverage(x, y float64) float64 {
	r := float64(p.Radius)
	if r <= 0 {
		return 1
	}
	// cx, cy is the nearest corner arc's centre, if this pixel is in a corner.
	cx, cy := x, y
	switch {
	case x < r:
		cx = r
	case x > float64(p.W)-r:
		cx = float64(p.W) - r
	}
	switch {
	case y < r:
		cy = r
	case y > float64(p.H)-r:
		cy = float64(p.H) - r
	}
	if cx == x && cy == y {
		return 1 // not in a corner: fully inside
	}
	d := math.Hypot(x-cx, y-cy)
	return math.Max(0, math.Min(1, r-d+0.5))
}

// bars draws the chart, inset so it never touches the rounded corners.
func (p Panel) bars(img *image.RGBA) {
	if len(p.Bars) == 0 {
		return
	}
	col := p.BarColor
	if col == (color.RGBA{}) {
		col = p.Ramp.To
		col.A = 255
	}
	pad := p.Radius + 2
	area := image.Rect(pad, pad, p.W-pad, p.H-pad)
	if area.Dx() <= 0 || area.Dy() <= 0 {
		return
	}
	// A gap of one pixel between bars — but only once there is a pixel to
	// spare. At one pixel per bar the gap would consume the whole bar and the
	// chart would silently draw nothing, which is how this read before a test
	// asked for fifty bars in twelve pixels.
	w := max(1, area.Dx()/len(p.Bars))
	bw := max(1, w-1)
	for i, v := range p.Bars {
		v = math.Max(0, math.Min(1, v))
		h := int(float64(area.Dy()) * v)
		x0 := area.Min.X + i*w
		for x := x0; x < min(x0+bw, area.Max.X); x++ {
			for y := area.Max.Y - h; y < area.Max.Y; y++ {
				img.SetRGBA(x, y, col)
			}
		}
	}
}

// Flatten composites an image with alpha onto an opaque background.
//
// This is Sixel's second concession made concrete: a protocol that cannot blend
// with the cells underneath needs the blending done in advance, against a
// background the caller has to KNOW rather than sample. Guess it wrong and the
// soft edge becomes a visible halo — which is why the alternative to knowing is
// not guessing but square corners.
func Flatten(src *image.RGBA, bg color.RGBA) *image.RGBA {
	b := src.Bounds()
	out := image.NewRGBA(b)
	for y := b.Min.Y; y < b.Max.Y; y++ {
		for x := b.Min.X; x < b.Max.X; x++ {
			c := src.RGBAAt(x, y)
			a := float64(c.A) / 255
			mix := func(fg, back uint8) uint8 {
				return uint8(float64(fg)*a + float64(back)*(1-a))
			}
			out.SetRGBA(x, y, color.RGBA{mix(c.R, bg.R), mix(c.G, bg.G), mix(c.B, bg.B), 255})
		}
	}
	return out
}

// Bar is a progress bar: a rounded track with a rounded ramp-filled portion.
//
// The design system's cell version of this is `[` `─`×filled `·`×remaining `]`,
// and that is what a reader sees on a terminal without graphics. This is the
// same information with the steps taken out — which is the honest test of
// whether a pixel layer is worth having at all. If the picture does not say
// something the characters cannot, it is decoration.
type Bar struct {
	W, H  int
	Value float64
	Ramp  Ramp
	// Track is the unfilled remainder. Alpha is respected, so a track can be
	// a faint tint rather than a colour.
	Track color.RGBA
	// Radius defaults to half the height — a bar with square ends reads as a
	// container rather than as a level.
	Radius int
}

// Image draws the bar.
func (b Bar) Image() *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, b.W, b.H))
	if b.W <= 0 || b.H <= 0 {
		return img
	}
	r := b.Radius
	if r == 0 {
		r = b.H / 2
	}
	shape := Panel{W: b.W, H: b.H, Radius: r}
	filled := int(math.Round(math.Max(0, math.Min(1, b.Value)) * float64(b.W)))

	for y := 0; y < b.H; y++ {
		for x := 0; x < b.W; x++ {
			c := b.Track
			if x < filled {
				// The ramp is sampled across the WHOLE bar, not across the
				// filled part, so the colour at a given level does not change
				// as the level moves. A gradient that slides is a gradient
				// that reads as motion nobody asked for.
				c = b.Ramp.At(float64(x) / float64(b.W-1))
			}
			c.A = uint8(float64(c.A) * shape.coverage(float64(x)+0.5, float64(y)+0.5))
			img.SetRGBA(x, y, c)
		}
	}
	return img
}
