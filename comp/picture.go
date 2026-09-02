package comp

import (
	"fmt"
	"image"
	"strings"

	"github.com/richarddavenport/tuikit/term"
)

// graphics is the pixel state, shared by every view of a canvas.
//
// Behind a pointer because Clip and WithChrome copy the Canvas struct. A
// picture requested through a clipped view — which is every component, since
// components are handed clips — has to land on the frame that is actually
// printed, and a slice field on a copy would append into a value nobody reads.
type graphics struct {
	mode         term.Graphics
	cellW, cellH int
	pictures     []picture
	nextID       int
}

type picture struct {
	r   Rect
	img *image.RGBA
	id  int
}

// WithGraphics returns a view that may carry pictures.
//
// A tool calls this once, at the top, with what [term.Detect] answered. Leave
// it unset and every Picture call is a no-op — which is the behaviour under
// test, and the reason goldens cannot move.
func (c *Canvas) WithGraphics(mode term.Graphics, cellW, cellH int) *Canvas {
	view := *c
	view.gfx = &graphics{mode: mode, cellW: max(1, cellW), cellH: max(1, cellH)}
	return &view
}

// Graphics is what this canvas can draw beyond characters.
func (c *Canvas) Graphics() term.Graphics {
	if c.gfx == nil {
		return term.None
	}
	return c.gfx.mode
}

// Picture asks that, if this terminal can, a picture be painted over the cells
// owned by id.
//
// The component has ALREADY drawn those cells as characters by the time it
// calls this. That ordering is the whole design: the cell drawing is not a
// fallback maintained for the pixel path's benefit, it is the drawing, and the
// picture is a decoration over it. A terminal that cannot show pictures shows
// what was there anyway, and no code was written twice to make that true.
//
// draw is a function rather than an image so that a terminal which cannot show
// it never pays to rasterise it. It is handed the region's size in PIXELS,
// which it cannot compute itself: only the canvas knows the cell size the
// terminal reported.
//
// Returns whether a picture was taken, so a component can adjust — a panel
// might leave its chart band empty rather than drawing bars underneath a
// picture of bars.
func (c *Canvas) Picture(id ID, draw func(w, h int) *image.RGBA) bool {
	if c.gfx == nil || c.gfx.mode == term.None || draw == nil {
		return false
	}
	r, ok := c.Region(id)
	if !ok || r.Empty() {
		return false
	}
	img := draw(r.W*c.gfx.cellW, r.H*c.gfx.cellH)
	if img == nil || img.Bounds().Empty() {
		return false
	}

	// Sixel replaces the cells it covers and cannot blend with them, so any
	// character left inside the rectangle would be painted over — except in the
	// frames between a redraw and the image being re-sent, where it would flash
	// back. Blanking is not tidiness; it is the only way the region has one
	// appearance rather than two.
	if c.gfx.mode == term.Sixel {
		c.Fill(r, " ", nil, id)
	}

	c.gfx.nextID++
	c.gfx.pictures = append(c.gfx.pictures, picture{r: r, img: img, id: c.gfx.nextID})
	return true
}

// pixels serialises the pictures, to be appended AFTER the frame's text.
//
// After, not before, because Sixel is opaque: a space printed over it erases
// it. The kitty protocol does not care — at z=-1 the image is under the text
// layer either way — so one ordering serves both, and the one that serves both
// is the one Sixel demands.
//
// Positioning is absolute, which assumes the frame's top-left is the terminal's
// home cell. That is true of an alt-screen program and false of one printing
// into a scrollback, so pictures are for alt-screen interfaces. The cursor is
// saved and restored around the whole block so that whatever the renderer
// believed about the cursor is still true afterwards.
func (c *Canvas) pixels() string {
	if c.gfx == nil || len(c.gfx.pictures) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\x1b7") // save cursor
	for _, p := range c.gfx.pictures {
		fmt.Fprintf(&b, "\x1b[%d;%dH", p.r.Y+1, p.r.X+1) // 1-based
		switch c.gfx.mode {
		case term.Sixel:
			b.WriteString(term.EncodeSixel(p.img))
		case term.Kitty:
			b.WriteString(term.EncodeKitty(p.img, p.id))
		}
	}
	b.WriteString("\x1b8") // restore cursor
	return b.String()
}
