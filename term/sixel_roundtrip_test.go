package term_test

import (
	"image"
	"image/color"
	"strconv"
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/term"
)

// decodeSixel is a minimal Sixel reader, written for this test alone.
//
// The encoder cannot be checked by looking at it in CI — the only terminal that
// could show it is on someone's desk. So the test decodes the bytes back to
// pixels and compares them to what went in, which catches the failures that
// actually happen: a band written in the wrong order, run-length counts off by
// one, colour registers scaled as bytes instead of percentages.
func decodeSixel(t *testing.T, s string) *image.RGBA {
	t.Helper()
	body, ok := strings.CutPrefix(s, "\x1bP")
	if !ok {
		t.Fatal("no DCS introducer")
	}
	body, ok = strings.CutSuffix(body, "\x1b\\")
	if !ok {
		t.Fatal("unterminated DCS")
	}
	i := strings.Index(body, "q")
	if i < 0 {
		t.Fatal("no sixel introducer")
	}
	body = body[i+1:]

	var img *image.RGBA
	palette := map[int]color.RGBA{}
	cur, x, top := 0, 0, 0

	for len(body) > 0 {
		switch c := body[0]; c {
		case '"': // raster attributes: "ratio;ratio;W;H
			end := strings.IndexAny(body[1:], "#-$")
			if end < 0 {
				end = len(body) - 1
			}
			f := strings.Split(body[1:1+end], ";")
			w, _ := strconv.Atoi(f[2])
			h, _ := strconv.Atoi(f[3])
			img = image.NewRGBA(image.Rect(0, 0, w, h))
			body = body[1+end:]

		case '#': // either a definition (#n;2;r;g;b) or a selection (#n)
			j := 1
			for j < len(body) && body[j] >= '0' && body[j] <= '9' {
				j++
			}
			n, _ := strconv.Atoi(body[1:j])
			cur = n
			if j < len(body) && body[j] == ';' {
				f := strings.SplitN(body[j+1:], ";", 4)
				r, _ := strconv.Atoi(f[1])
				g, _ := strconv.Atoi(f[2])
				k := 0
				for k < len(f[3]) && f[3][k] >= '0' && f[3][k] <= '9' {
					k++
				}
				b, _ := strconv.Atoi(f[3][:k])
				// Components are percentages on the wire.
				palette[n] = color.RGBA{pct8(r), pct8(g), pct8(b), 255}
				j += 1 + len(f[0]) + 1 + len(f[1]) + 1 + len(f[2]) + 1 + k
			}
			x = 0
			body = body[j:]

		case '$': // carriage return: same band, next colour
			x = 0
			body = body[1:]

		case '-': // next band
			top += 6
			x = 0
			body = body[1:]

		case '!': // run length
			j := 1
			for j < len(body) && body[j] >= '0' && body[j] <= '9' {
				j++
			}
			n, _ := strconv.Atoi(body[1:j])
			for k := 0; k < n; k++ {
				putSixel(img, palette[cur], x, top, body[j])
				x++
			}
			body = body[j+1:]

		default:
			putSixel(img, palette[cur], x, top, c)
			x++
			body = body[1:]
		}
	}
	return img
}

func putSixel(img *image.RGBA, c color.RGBA, x, top int, ch byte) {
	if img == nil || ch < 0x3F {
		return
	}
	bits := ch - 0x3F
	for r := 0; r < 6; r++ {
		if bits&(1<<r) != 0 && x < img.Bounds().Dx() && top+r < img.Bounds().Dy() {
			img.SetRGBA(x, top+r, c)
		}
	}
}

func pct8(v int) uint8 { return uint8(v * 255 / 100) }

// TestSixelRoundTrip encodes an image, decodes it back, and compares.
//
// Colours are compared after quantising to the 6-cube, because that is the
// documented lossy step — anything else being different is a bug.
func TestSixelRoundTrip(t *testing.T) {
	const w, h = 37, 19 // deliberately not multiples of 6
	src := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			src.SetRGBA(x, y, color.RGBA{
				R: uint8(x * 255 / w), G: uint8(y * 255 / h), B: 128, A: 255,
			})
		}
	}

	got := decodeSixel(t, term.EncodeSixel(src))
	if b := got.Bounds(); b.Dx() != w || b.Dy() != h {
		t.Fatalf("decoded %v, want %dx%d", b, w, h)
	}

	bad := 0
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			want, have := snap(src.RGBAAt(x, y)), got.RGBAAt(x, y)
			if !near(want, have) {
				if bad < 5 {
					t.Errorf("pixel %d,%d = %v, want ~%v", x, y, have, want)
				}
				bad++
			}
		}
	}
	if bad > 0 {
		t.Fatalf("%d of %d pixels wrong", bad, w*h)
	}
}

// snap is the encoder's documented quantisation: the 6×6×6 cube.
func snap(c color.RGBA) color.RGBA {
	q := func(v uint8) uint8 { return uint8(int(v) * 5 / 255 * 255 / 5) }
	return color.RGBA{q(c.R), q(c.G), q(c.B), 255}
}

// near allows one step of rounding in each direction: the wire format carries
// percentages, so 255 makes the trip as 100% and comes back as 255 only after
// two integer divisions.
func near(a, b color.RGBA) bool {
	d := func(x, y uint8) bool { return int(x)-int(y) < 4 && int(y)-int(x) < 4 }
	return d(a.R, b.R) && d(a.G, b.G) && d(a.B, b.B)
}

// TestSixelRoundTripFlat is the case the RLE path takes.
func TestSixelRoundTripFlat(t *testing.T) {
	src := fill(64, 12, color.RGBA{0, 204, 102, 255})
	got := decodeSixel(t, term.EncodeSixel(src))
	for y := 0; y < 12; y++ {
		for x := 0; x < 64; x++ {
			if !near(snap(src.RGBAAt(x, y)), got.RGBAAt(x, y)) {
				t.Fatalf("pixel %d,%d = %v, want ~%v", x, y, got.RGBAAt(x, y), snap(src.RGBAAt(x, y)))
			}
		}
	}
}
