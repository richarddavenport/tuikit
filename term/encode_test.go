package term_test

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/term"
)

func fill(w, h int, c color.RGBA) *image.RGBA {
	img := image.NewRGBA(image.Rect(0, 0, w, h))
	for y := 0; y < h; y++ {
		for x := 0; x < w; x++ {
			img.SetRGBA(x, y, c)
		}
	}
	return img
}

func TestSixelShape(t *testing.T) {
	out := term.EncodeSixel(fill(12, 12, color.RGBA{255, 0, 0, 255}))

	if !strings.HasPrefix(out, "\x1bP") {
		t.Error("no device control string introducer")
	}
	if !strings.HasSuffix(out, "\x1b\\") {
		t.Error("unterminated device control string")
	}
	if !strings.Contains(out, `"1;1;12;12`) {
		t.Error("no raster attributes; a terminal has to be told the size")
	}
	// Red is 5,0,0 in the 6-cube → index 180, and 255 is 100%.
	if !strings.Contains(out, "#180;2;100;0;0") {
		t.Errorf("red was not defined as a colour register:\n%q", out)
	}
	// 12 rows is exactly two bands of six, so there is one band separator and
	// it is not trailing.
	if n := strings.Count(out, "-"); n != 1 {
		t.Errorf("%d band separators for 12 rows, want 1", n)
	}
}

// TestSixelRunLength: a flat panel is the common case and the whole reason a
// Sixel redraw is affordable at all.
func TestSixelRunLength(t *testing.T) {
	out := term.EncodeSixel(fill(200, 6, color.RGBA{0, 0, 255, 255}))
	if !strings.Contains(out, "!200") {
		t.Errorf("200 identical columns were not run-length encoded:\n%q", out)
	}
	if len(out) > 200 {
		t.Errorf("encoded length %d for a flat 200×6 image; RLE is not working", len(out))
	}
}

// TestSixelSkipsColoursNotInABand. Emitting every register for every band is
// correct and enormous; on a two-colour panel it doubles the payload.
func TestSixelSkipsColoursNotInABand(t *testing.T) {
	img := fill(6, 12, color.RGBA{255, 0, 0, 255})
	for y := 6; y < 12; y++ { // the second band only
		for x := 0; x < 6; x++ {
			img.SetRGBA(x, y, color.RGBA{0, 255, 0, 255})
		}
	}
	out := term.EncodeSixel(img)
	first := strings.Index(out, "#")
	if first < 0 {
		t.Fatal("no colour registers in the output")
	}
	bands := strings.Split(out[first:], "-")
	if len(bands) != 2 {
		t.Fatalf("want 2 bands, got %d", len(bands))
	}
	if strings.Contains(bands[1], "$") {
		t.Error("the second band emitted a colour that is not in it")
	}
}

func TestSixelEmptyImage(t *testing.T) {
	if out := term.EncodeSixel(image.NewRGBA(image.Rect(0, 0, 0, 0))); out != "" {
		t.Errorf("empty image encoded as %q, want the empty string", out)
	}
}

func TestKittyFlags(t *testing.T) {
	out := term.EncodeKitty(fill(4, 4, color.RGBA{1, 2, 3, 255}), 7)

	for _, want := range []string{
		"\x1b_G",  // application programmable command
		"a=T",     // transmit and display in one step
		"f=32",    // RGBA — the alpha channel is the point
		"o=z",     // zlib, because RGBA is four bytes a pixel
		"s=4,v=4", // the size
		"i=7",     // the id, so it can be deleted later
		"z=-1",    // BEHIND the text: the whole reason to prefer this protocol
		"C=1",     // do not move the cursor
		"q=2",     // no reply — one arriving mid-frame would read as a keystroke
		"m=0",     // the last chunk
		"\x1b\\",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("missing %q in:\n%q", want, out)
		}
	}
}

// TestKittyChunks: kitty takes at most 4096 bytes of payload per escape
// sequence, and a payload silently truncated is an image that never appears.
func TestKittyChunks(t *testing.T) {
	// Genuine noise. An arithmetic pattern is not enough: zlib crushed
	// byte(i*7%251) over 160kB down under one chunk, and the test passed for
	// the wrong reason until it was asked to prove chunking happened at all.
	img := image.NewRGBA(image.Rect(0, 0, 200, 200))
	seed := uint32(1)
	for i := range img.Pix {
		seed = seed*1664525 + 1013904223
		img.Pix[i] = byte(seed >> 24)
	}
	out := term.EncodeKitty(img, 1)

	if n := strings.Count(out, "m=1"); n < 1 {
		t.Fatalf("a large image produced %d continuation chunks, want at least 1", n)
	}
	if n := strings.Count(out, "m=0"); n != 1 {
		t.Errorf("%d final chunks, want exactly 1", n)
	}
	// Only the first chunk carries the metadata.
	if n := strings.Count(out, "a=T"); n != 1 {
		t.Errorf("metadata repeated on %d chunks, want 1", n)
	}
	// Every chunk's control data must be well formed. Written as ",m=1" the
	// continuations came out as `ESC_G,m=1;` — a leading comma — and the
	// terminal kept the first chunk and refused the rest. One-chunk images were
	// unaffected, so this only broke pictures big enough to be worth drawing.
	for i, part := range strings.Split(out, "\x1b_G")[1:] {
		if strings.HasPrefix(part, ",") {
			t.Errorf("chunk %d starts with a comma: %q", i, part[:min(40, len(part))])
		}
	}

	for _, part := range strings.Split(out, "\x1b_G")[1:] {
		semi, end := strings.Index(part, ";"), strings.Index(part, "\x1b\\")
		if semi < 0 || end < semi {
			t.Fatalf("malformed chunk: %q", part[:min(80, len(part))])
		}
		payload := part[semi+1 : end]
		if len(payload) > 4096 {
			t.Errorf("chunk of %d bytes exceeds kitty's 4096 limit", len(payload))
		}
	}
}

func TestKittyEmptyImage(t *testing.T) {
	if out := term.EncodeKitty(image.NewRGBA(image.Rect(0, 0, 0, 0)), 1); out != "" {
		t.Errorf("empty image encoded as %q", out)
	}
}

func TestDeleteKitty(t *testing.T) {
	if got, want := term.DeleteKitty(3), "\x1b_Ga=d,d=i,i=3,q=2\x1b\\"; got != want {
		t.Errorf("DeleteKitty(3) = %q, want %q", got, want)
	}
}

// TestKittyHandlesStride: a sub-image has a stride wider than its bounds, and
// handing kitty img.Pix directly would shear it.
func TestKittyHandlesStride(t *testing.T) {
	full := fill(100, 100, color.RGBA{9, 9, 9, 255})
	sub := full.SubImage(image.Rect(10, 10, 20, 20)).(*image.RGBA)
	out := term.EncodeKitty(sub, 1)
	if !strings.Contains(out, "s=10,v=10") {
		t.Errorf("sub-image encoded with the wrong size:\n%q", out[:min(120, len(out))])
	}
}

func TestParseCellSize(t *testing.T) {
	// CSI 16 t is answered CSI 6 ; HEIGHT ; WIDTH t — height first, which is
	// the obvious thing to get backwards.
	if w, h, ok := term.ParseCellSize("\x1b[6;38;19t"); w != 19 || h != 38 || !ok {
		t.Errorf("got %dx%d measured=%v, want 19x38 measured", w, h, ok)
	}

	// iTerm2 declines CSI 16 t. The fallback divides the text area in pixels
	// by the text area in cells — without it every picture there came out at
	// the assumed 10x20 while the cells were larger, so the bars were about
	// two thirds the width they should have been.
	reply := "\x1b[4;1200;1680t" + "\x1b[8;60;120t" + "\x1b[?62;4c"
	if w, h, ok := term.ParseCellSize(reply); w != 14 || h != 20 || !ok {
		t.Errorf("fallback gave %dx%d measured=%v, want 14x20 measured", w, h, ok)
	}

	// The direct answer wins when both arrive.
	both := "\x1b[6;38;19t\x1b[4;1200;1680t\x1b[8;60;120t"
	if w, h, _ := term.ParseCellSize(both); w != 19 || h != 38 {
		t.Errorf("got %dx%d, want the direct answer 19x38", w, h)
	}

	for _, bad := range []string{
		"", "garbage", "\x1b[6;0;0t", "\x1b[6;10t", "\x1b[6;a;bt",
		"\x1b[8;60;120t",               // cells but no pixels
		"\x1b[4;1200;1680t",            // pixels but no cells
		"\x1b[4;1200;1680t\x1b[8;0;0t", // would divide by zero
	} {
		w, h, ok := term.ParseCellSize(bad)
		if w != term.DefaultCellW || h != term.DefaultCellH || ok {
			t.Errorf("ParseCellSize(%q) = %dx%d measured=%v, want the default and not measured",
				bad, w, h, ok)
		}
	}
}
