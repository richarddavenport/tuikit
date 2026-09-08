package comp

import (
	"image"
	"image/color"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/paint"
)

// Meter is how far along something is.
//
// From the deploy tool's run dialog and its activity strip, which both wanted
// the same bar and had neither — the dialog showed five rows all reading
// "measuring…", which says a thing is happening and nothing about how much of
// it is left.
//
// The characters are `[`, `─` for the filled part and `·` for the rest. Block
// elements (U+2580–U+259F) are deliberately not used: they are the obvious
// choice and they are excluded from the glyph set, because a font without them
// draws a progress bar as a row of empty boxes.
//
// # The pixel layer
//
// Set Ramp and the meter offers a smooth gradient bar to a terminal that can
// draw one — and this is the honest test of whether pixels are worth having,
// because the characters already say the number. What the picture adds is
// resolution: a cell bar has as many steps as it has columns, so at 24 columns
// it cannot tell 51% from 53%, and it snaps. The picture does not.
//
// A terminal that cannot show it shows the characters, which were drawn first
// and unconditionally. Nothing here is written twice.
type Meter struct {
	// Value is 0..1, clamped. Out of range is a caller's arithmetic being
	// wrong, and a bar that overflows its brackets makes that harder to see
	// rather than easier.
	Value float64
	// Label follows the bar: "2 of 5", "24s elapsed". Optional.
	Label string
	// Track names the bar's interior so the pixel layer can find it. A picture
	// is placed on an OWNER, not on coordinates — the meter does not know where
	// it ended up and does not need to.
	Track Name

	Filled, Empty, LabelStyle *lipgloss.Style

	// Pixels asks for a smooth bar where the terminal can draw one.
	//
	// A bool rather than a colour, deliberately. The gradient comes from the
	// canvas, which read it from the terminal's own palette — so a component
	// says WHETHER it wants a picture and never WHAT COLOUR, which is the same
	// arrangement the nine roles give the characters. A component that could
	// pass its own gradient is a component that can escape the theme.
	Pixels bool
}

// Draw paints the meter into one row of r and returns the row below.
func (m Meter) Draw(c *Canvas, r Rect, id ID) int {
	if r.Empty() {
		return r.Y
	}
	label := m.Label
	labelW := 0
	if label != "" {
		labelW = Width(label) + 1
	}
	// Two brackets and at least one interior cell, or there is no bar to draw.
	barW := r.W - labelW
	if barW < 3 {
		return r.Y + 1
	}
	inner := barW - 2
	filled := int(float64(inner)*clamp01(m.Value) + 0.5)

	track := Region(m.Track)
	if m.Track == "" {
		track = id
	}
	x := r.X
	x += c.Text(x, r.Y, "[", m.Empty, id)
	for i := 0; i < inner; i++ {
		glyph, style := "·", m.Empty
		if i < filled {
			glyph, style = "─", m.Filled
		}
		x += c.Set(x, r.Y, glyph, style, track)
	}
	x += c.Text(x, r.Y, "]", m.Empty, id)
	if label != "" {
		c.Text(x+1, r.Y, label, m.LabelStyle, id)
	}

	// Offered after the characters are already down. The canvas declines on a
	// terminal that cannot draw it, and the rasteriser is never even called.
	if m.Pixels {
		ramp := c.Ramp()
		value := clamp01(m.Value)
		c.Picture(track, func(w, h int) *image.RGBA {
			return paint.Bar{W: w, H: h, Value: value, Ramp: ramp, Track: trackTint(ramp)}.Image()
		})
	}
	return r.Y + 1
}

// trackTint is the unfilled remainder: the ramp's own start, mostly
// transparent. Derived rather than configured, because a track colour that
// does not belong to the ramp is a tenth colour role nobody named.
func trackTint(r paint.Ramp) color.RGBA {
	c := r.From
	c.A = 60
	return c
}

func clamp01(v float64) float64 {
	if v < 0 {
		return 0
	}
	if v > 1 {
		return 1
	}
	return v
}
