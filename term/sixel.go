package term

import (
	"fmt"
	"image"
	"image/color"
	"strings"
)

// sixelColors is the palette every image is quantised to: the 6×6×6 cube.
//
// Sixel is colour-register based rather than truecolour, so a gradient is
// stepped no matter what — the only question is how coarsely. 216 registers is
// inside what every Sixel terminal supports (the usual floor is 256) and needs
// no per-image analysis, which matters because this runs on a redraw.
//
// A median-cut palette per image would look better and would also mean the
// panel's colours shifting slightly whenever its data changed. Stepping that
// stays put reads as a design; stepping that moves reads as a bug.
const sixelLevels = 6

// EncodeSixel encodes an image as a Sixel escape sequence.
//
// The image must be opaque: Sixel has no alpha, and a transparent pixel here
// would be written as whatever colour it happens to carry rather than showing
// what is behind it. Use [paint.Flatten] against the terminal's background
// first — the caller knows that colour and this package does not.
func EncodeSixel(img *image.RGBA) string {
	b := img.Bounds()
	w, h := b.Dx(), b.Dy()
	if w <= 0 || h <= 0 {
		return ""
	}

	var out strings.Builder
	// P1=0 (aspect from raster attrs), P2=1 (0 bits are transparent, so the
	// image is exactly its own shape), P3=0.
	out.WriteString("\x1bP0;1;0q")
	fmt.Fprintf(&out, "\"1;1;%d;%d", w, h)

	// Quantise once. Doing it per band would read every pixel six times.
	idx := make([]int, w*h)
	used := map[int]bool{}
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			i := quantise(img.RGBAAt(b.Min.X+x, b.Min.Y+y))
			idx[y*w+x] = i
			used[i] = true
		}
	}
	for i := range used {
		r, g, bl := unquantise(i)
		// Sixel colour components are percentages, not bytes.
		fmt.Fprintf(&out, "#%d;2;%d;%d;%d", i, pct(r), pct(g), pct(bl))
	}

	// A sixel is six vertical pixels in one character, so the image is walked
	// in bands of six rows and each band is written once per colour in it.
	for top := 0; top < h; top += 6 {
		first := true
		for i := range used {
			row := bandRow(idx, w, h, top, i)
			if row == "" {
				continue // this colour is not in this band
			}
			if !first {
				out.WriteByte('$') // back to the start of the same band
			}
			first = false
			fmt.Fprintf(&out, "#%d%s", i, row)
		}
		if top+6 < h {
			out.WriteByte('-') // next band
		}
	}
	out.WriteString("\x1b\\")
	return out.String()
}

// bandRow is one colour's contribution to one six-row band, run-length encoded.
//
// Returns "" when the colour does not appear, so the caller can skip it
// entirely — on a panel with a handful of colours per band that is the
// difference between an image and a wall of empty runs.
func bandRow(idx []int, w, h, top, want int) string {
	var b strings.Builder
	runChar, runLen, any := byte(0), 0, false
	flush := func() {
		if runLen == 0 {
			return
		}
		// !n repeats the next character n times. Below four it costs more than
		// it saves, which is why the threshold is not one.
		if runLen > 3 {
			fmt.Fprintf(&b, "!%d%c", runLen, runChar)
		} else {
			b.WriteString(strings.Repeat(string(runChar), runLen))
		}
		runLen = 0
	}
	for x := 0; x < w; x++ {
		var bits byte
		for r := 0; r < 6; r++ {
			if y := top + r; y < h && idx[y*w+x] == want {
				bits |= 1 << r
			}
		}
		if bits != 0 {
			any = true
		}
		ch := 0x3F + bits
		if ch == runChar {
			runLen++
			continue
		}
		flush()
		runChar, runLen = ch, 1
	}
	flush()
	if !any {
		return ""
	}
	return b.String()
}

// quantise maps a colour to the 6×6×6 cube.
func quantise(c color.RGBA) int {
	q := func(v uint8) int { return int(v) * (sixelLevels - 1) / 255 }
	return q(c.R)*sixelLevels*sixelLevels + q(c.G)*sixelLevels + q(c.B)
}

// unquantise is quantise's inverse, back to the centre of the cube cell.
func unquantise(i int) (r, g, b uint8) {
	v := func(n int) uint8 { return uint8(n * 255 / (sixelLevels - 1)) }
	return v(i / (sixelLevels * sixelLevels) % sixelLevels), v(i / sixelLevels % sixelLevels), v(i % sixelLevels)
}

func pct(v uint8) int { return int(v) * 100 / 255 }
