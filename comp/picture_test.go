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
	c := comp.NewCanvas(20, 6).WithGraphics(mode, 10, 20)
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

// TestSixelBlanksTheRegion. Sixel is opaque, so a character left inside the
// rectangle is a character that flashes back between redraws.
func TestSixelBlanksTheRegion(t *testing.T) {
	c := canvasWith(t, term.Sixel)
	if !c.Picture(comp.Region(reg), solid(nil)) {
		t.Fatal("Picture declined on a sixel canvas")
	}
	cell, _ := c.CellAt(3, 2)
	if cell.Text != " " {
		t.Errorf("cell inside a sixel picture = %q, want a blank", cell.Text)
	}
	if owner := c.OwnerAt(3, 2); owner.Name != reg {
		t.Errorf("blanking lost the owner: %v", owner)
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
	if !strings.Contains(frame, "\x1b[2;3H") { // region at 2,1 → 1-based 3,2
		t.Error("no absolute positioning before the picture")
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

// TestDrawReturningNil: a rasteriser that gives up must not produce an empty
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
