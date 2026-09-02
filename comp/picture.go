package comp

import (
	"fmt"
	"image"
	"image/color"
	"strings"

	"github.com/richarddavenport/tuikit/paint"
	"github.com/richarddavenport/tuikit/term"
)

// Pixels is everything a canvas needs to know to draw a picture.
//
// One struct rather than four arguments, because three of the four are
// measurements only the terminal can supply and they arrive together — a call
// site that has the mode but not the cell size is a call site that is about to
// stretch something.
type Pixels struct {
	Mode         term.Graphics
	CellW, CellH int
	// Background is what a Sixel picture is flattened against, since Sixel has
	// no alpha. Ignored by the kitty path, which composites for real.
	Background color.RGBA
	// Ramp is the gradient a picture draws with, read from the terminal's OWN
	// palette rather than chosen here.
	//
	// This is the whole reason a component asks for "a picture" instead of
	// passing colours. Decision 28 put the characters on ANSI 0–15 so the
	// reader's theme wins; a picture carrying a literal gradient would look
	// identical under all 22 Omarchy themes while the text beside it changed,
	// which is a worse result than having no picture at all.
	Ramp paint.Ramp
}

// Detect asks the terminal all three questions at once.
//
// Called by main, never by a model. Detection reads /dev/tty, and a constructor
// that did it would query the developer's own terminal during `go test` — which
// on a graphics-capable one would write escape sequences into the goldens.
func Detect() Pixels {
	mode := term.Detect()
	if mode == term.None {
		return Pixels{} // do not spend two more round trips to learn nothing
	}
	w, h := term.CellSize()
	p := Pixels{Mode: mode, CellW: w, CellH: h, Background: term.Background(), Ramp: fallbackRamp}
	// 5 and 13 are magenta and bright magenta — Accent's own family, per
	// decision 28. Read from the terminal, so changing theme changes the
	// picture as well as the text.
	if got := term.Colors(rampFrom, rampTo); got != nil {
		if c, ok := got[rampFrom]; ok {
			p.Ramp.From = c
		}
		if c, ok := got[rampTo]; ok {
			p.Ramp.To = c
		}
	}
	return p
}

// The two palette entries a picture ramps between, and what to use when the
// terminal will not say: xterm's own defaults for those indices.
const (
	rampFrom = 5  // magenta
	rampTo   = 13 // bright magenta — the index Accent maps to
)

var fallbackRamp = paint.Ramp{
	From: color.RGBA{R: 0xcd, B: 0xcd, A: 0xff},
	To:   color.RGBA{R: 0xff, B: 0xff, A: 0xff},
}

// graphics is the pixel state, shared by every view of a canvas.
//
// Behind a pointer because Clip and WithChrome copy the Canvas struct. A
// picture requested through a clipped view — which is every component, since
// components are handed clips — has to land on the frame that is actually
// printed, and a slice field on a copy would append into a value nobody reads.
type graphics struct {
	Pixels
	pictures []picture
	nextID   int
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
func (c *Canvas) WithGraphics(p Pixels) *Canvas {
	p.CellW, p.CellH = max(1, p.CellW), max(1, p.CellH)
	if p.Background == (color.RGBA{}) {
		p.Background = term.DefaultBackground
	}
	if p.Ramp == (paint.Ramp{}) {
		p.Ramp = fallbackRamp
	}
	view := *c
	view.gfx = &graphics{Pixels: p}
	return &view
}

// Ramp is the gradient a picture on this canvas draws with.
func (c *Canvas) Ramp() paint.Ramp {
	if c.gfx == nil {
		return fallbackRamp
	}
	return c.gfx.Ramp
}

// Graphics is what this canvas can draw beyond characters.
func (c *Canvas) Graphics() term.Graphics {
	if c.gfx == nil {
		return term.None
	}
	return c.gfx.Mode
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
	if c.gfx == nil || c.gfx.Mode == term.None || draw == nil {
		return false
	}
	r, ok := c.Region(id)
	if !ok || r.Empty() {
		return false
	}
	img := draw(r.W*c.gfx.CellW, r.H*c.gfx.CellH)
	if img == nil || img.Bounds().Empty() {
		return false
	}

	// Sixel replaces the cells it covers and cannot blend with them, so any
	// character left inside the rectangle would be painted over — except in the
	// frames between a redraw and the image being re-sent, where it would flash
	// back. Blanking is not tidiness; it is the only way the region has one
	// appearance rather than two.
	if c.gfx.Mode == term.Sixel {
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
// Positioning is RELATIVE to where the frame ended, not absolute.
//
// Absolute placement (CSI row;col H) was the first attempt and it is wrong
// wherever the frame does not start at the terminal's home cell — which is any
// program printing into a scrollback rather than running alt-screen. It put
// every picture in the top-left corner of the window while the characters it
// belonged to sat further down, and because the Sixel path blanks its region
// first, what you saw was an empty bar and a stray image somewhere else.
//
// These bytes are appended after the last row, so the cursor is on row h-1 of
// the frame. Moving up from there costs nothing in an alt-screen program and is
// the only thing that works outside one. Columns use CSI n G, which is relative
// to the line rather than the screen, so the horizontal half needs no
// arithmetic at all.
func (c *Canvas) pixels() string {
	if c.gfx == nil || len(c.gfx.pictures) == 0 {
		return ""
	}
	var b strings.Builder
	b.WriteString("\x1b7") // save cursor
	for _, p := range c.gfx.pictures {
		if up := c.h - 1 - p.r.Y; up > 0 {
			fmt.Fprintf(&b, "\x1b[%dA", up)
		}
		fmt.Fprintf(&b, "\x1b[%dG", p.r.X+1) // column, 1-based, within the line
		switch c.gfx.Mode {
		case term.Sixel:
			// Sixel has no alpha, so the soft edges are composited here
			// against the colour the terminal reported. This is the handoff's
			// second concession — the shadow is baked, and it can only be
			// baked against a background that is KNOWN rather than guessed.
			b.WriteString(term.EncodeSixel(paint.Flatten(p.img, c.gfx.Background)))
		case term.Kitty:
			b.WriteString(term.EncodeKitty(p.img, p.id))
		}
		// Back to where the frame ended, so the next picture's "up" is
		// measured from the same place this one's was.
		b.WriteString("\x1b8\x1b7")
	}
	b.WriteString("\x1b8") // restore cursor
	return b.String()
}
