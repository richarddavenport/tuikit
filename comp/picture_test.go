package comp_test

import (
	"image"
	"image/color"
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/term"
)

const reg = comp.Name("panel")

// solid is a draw function that records the size it was asked for.
func solid(got *[2]int) func(w, h int) *image.RGBA {
	return func(w, h int) *image.RGBA {
		if got != nil {
			*got = [2]int{w, h}
		}
		img := image.NewRGBA(image.Rect(0, 0, w, h))
		for i := range img.Pix {
			img.Pix[i] = 255
		}
		return img
	}
}

func canvasWith(t *testing.T, mode term.Graphics) *comp.Canvas {
	t.Helper()
	c := comp.NewCanvas(20, 6).WithGraphics(comp.Pixels{Mode: mode, CellW: 10, CellH: 20})
	c.Fill(comp.Rect{X: 2, Y: 1, W: 8, H: 3}, "x", nil, comp.Region(reg))
	return c
}

// TestNoGraphicsIsByteIdentical is the property every golden rests on.
//
// It is not "close enough" or "visually the same": the frame must be the exact
// bytes it was before a pixel layer existed, because the alternative is that
// every golden in every tool moves the day someone adds a picture.
func TestNoGraphicsIsByteIdentical(t *testing.T) {
	plain := comp.NewCanvas(20, 6)
	plain.Fill(comp.Rect{X: 2, Y: 1, W: 8, H: 3}, "x", nil, comp.Region(reg))
	want := plain.String()

	c := comp.NewCanvas(20, 6) // no WithGraphics at all
	c.Fill(comp.Rect{X: 2, Y: 1, W: 8, H: 3}, "x", nil, comp.Region(reg))
	if took := c.Picture(comp.Region(reg), solid(nil)); took {
		t.Error("Picture reported taking a picture on a canvas with no graphics")
	}
	if got := c.String(); got != want {
		t.Errorf("frame changed:\n got %q\nwant %q", got, want)
	}
	if c.Graphics() != term.None {
		t.Errorf("Graphics() = %v, want None", c.Graphics())
	}
}

// TestNoneModeIsAlsoByteIdentical: explicitly detecting "this terminal cannot"
// must be as inert as never having asked.
func TestNoneModeIsAlsoByteIdentical(t *testing.T) {
	want := canvasWith(t, term.None).String()
	c := canvasWith(t, term.None)
	c.Picture(comp.Region(reg), solid(nil))
	if got := c.String(); got != want {
		t.Errorf("None mode changed the frame:\n got %q\nwant %q", got, want)
	}
	if strings.Contains(c.String(), "\x1b7") {
		t.Error("None mode emitted a cursor save")
	}
}

// TestPixelSizeComesFromTheRegion: the draw function is handed pixels, because
// it is the one thing it cannot work out for itself.
func TestPixelSizeComesFromTheRegion(t *testing.T) {
	var got [2]int
	c := canvasWith(t, term.Kitty)
	if !c.Picture(comp.Region(reg), solid(&got)) {
		t.Fatal("Picture declined on a kitty canvas")
	}
	// The region is 8×3 cells and the cell is 10×20 pixels.
	if want := [2]int{80, 60}; got != want {
		t.Errorf("draw asked for %v pixels, want %v", got, want)
	}
}

// TestSixelKeepsTheCharactersUnderneath.
//
// An opaque image covers them, so blanking looks harmless — and it is, right up
// until the image does not arrive. Forced onto a terminal that cannot draw
// Sixel, or stripped by a multiplexer, blanking leaves an EMPTY bar where the
// character bar would have been: worse than having no pixel layer at all.
func TestSixelKeepsTheCharactersUnderneath(t *testing.T) {
	c := canvasWith(t, term.Sixel)
	if !c.Picture(comp.Region(reg), solid(nil)) {
		t.Fatal("Picture declined on a sixel canvas")
	}
	cell, _ := c.CellAt(3, 2)
	if cell.Text != "x" {
		t.Errorf("cell under a sixel picture is %q; it should still hold the character "+
			"so that a picture which never arrives degrades to the cells", cell.Text)
	}
}

// TestKittyKeepsTheText is the whole reason to prefer the protocol: at z=-1 the
// image is under the text layer, so the words stay real text.
func TestKittyKeepsTheText(t *testing.T) {
	c := canvasWith(t, term.Kitty)
	c.Picture(comp.Region(reg), solid(nil))
	if cell, _ := c.CellAt(3, 2); cell.Text != "x" {
		t.Errorf("cell inside a kitty picture = %q, want the text to survive", cell.Text)
	}
}

// TestPicturesComeAfterTheText, because Sixel is opaque and a space printed
// over it erases it.
func TestPicturesComeAfterTheText(t *testing.T) {
	c := canvasWith(t, term.Sixel)
	c.Picture(comp.Region(reg), solid(nil))
	frame := c.String()

	save := strings.Index(frame, "\x1b7")
	if save < 0 {
		t.Fatal("no cursor save in the frame")
	}
	if text := strings.Index(frame, "x"); text > save {
		t.Error("text is emitted after the picture; Sixel would erase it")
	}
	if !strings.HasSuffix(frame, "\x1b8") {
		t.Error("the frame does not restore the cursor")
	}
	// The canvas is 6 rows and the region starts at y=1, so the picture is 4
	// rows above where the frame ended, at column 3 (1-based).
	if !strings.Contains(frame, "\x1b[4A") {
		t.Error("no upward move to the picture's row")
	}
	if !strings.Contains(frame, "\x1b[3G") {
		t.Error("no column move to the picture's start")
	}
	// Absolute positioning would put the picture at the top of the WINDOW
	// rather than on the cells it belongs to, wherever the frame happens to
	// have been printed.
	if strings.Contains(frame, "H\x1bP") {
		t.Error("absolute cursor positioning is back")
	}
}

// TestAViewRecordsOntoTheFrame. Clip copies the Canvas struct, so a component —
// which is always handed a clip — must still land its picture on the frame that
// is actually printed.
func TestAViewRecordsOntoTheFrame(t *testing.T) {
	c := canvasWith(t, term.Kitty)
	view := c.Clip(comp.Rect{X: 0, Y: 0, W: 20, H: 6})
	if !view.Picture(comp.Region(reg), solid(nil)) {
		t.Fatal("Picture declined on a clipped view")
	}
	if !strings.Contains(c.String(), "\x1b_G") {
		t.Error("a picture taken through a clip did not reach the parent's frame")
	}
}

func TestPictureOnAnUnknownRegion(t *testing.T) {
	c := canvasWith(t, term.Kitty)
	if c.Picture(comp.Region("nothing-drew-this"), solid(nil)) {
		t.Error("Picture claimed to draw on a region that does not exist")
	}
}

func TestNilDrawIsDeclined(t *testing.T) {
	c := canvasWith(t, term.Kitty)
	if c.Picture(comp.Region(reg), nil) {
		t.Error("Picture accepted a nil draw function")
	}
}

// TestDrawReturningNil: a rasterizer that gives up must not produce an empty
// escape sequence that still moves the cursor.
func TestDrawReturningNil(t *testing.T) {
	c := canvasWith(t, term.Kitty)
	if c.Picture(comp.Region(reg), func(int, int) *image.RGBA { return nil }) {
		t.Error("Picture accepted a nil image")
	}
	if strings.Contains(c.String(), "\x1b7") {
		t.Error("a declined picture still emitted cursor movement")
	}
}

var _ = color.RGBA{}

// TestKittyFootprintComesFromTheRegion is what makes the kitty path immune to
// the cell-size disagreement in issue 31.
//
// Without c and r the terminal divides the image's pixel size by its own idea
// of a cell — the one number nobody agrees on. Ghostty reports it in device
// pixels and iTerm2 in points, so the same picture came out at two different
// scales, one of them overlapping the text above it. The region knows how many
// cells it is, exactly, and that is what gets sent.
func TestKittyFootprintComesFromTheRegion(t *testing.T) {
	c := canvasWith(t, term.Kitty) // region is 8×3 cells
	if !c.Picture(comp.Region(reg), solid(nil)) {
		t.Fatal("Picture declined")
	}
	if got := c.String(); !strings.Contains(got, "c=8,r=3") {
		t.Error("the kitty placement does not state its footprint in cells")
	}
}

// TestAWrongCellSizeStillPlacesTheRightFootprint: the pixels change, the cells
// do not. That is the property — a bad measurement makes a blurrier picture,
// not a misplaced one.
func TestAWrongCellSizeStillPlacesTheRightFootprint(t *testing.T) {
	for _, cell := range [][2]int{{8, 18}, {16, 34}, {1, 1}} {
		c := comp.NewCanvas(20, 6).WithGraphics(comp.Pixels{
			Mode: term.Kitty, CellW: cell[0], CellH: cell[1],
		})
		c.Fill(comp.Rect{X: 2, Y: 1, W: 8, H: 3}, "x", nil, comp.Region(reg))
		c.Picture(comp.Region(reg), solid(nil))
		if got := c.String(); !strings.Contains(got, "c=8,r=3") {
			t.Errorf("cell %v changed the footprint", cell)
		}
	}
}
