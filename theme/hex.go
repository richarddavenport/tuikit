package theme

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Hex is what an ANSI 256 index looks like, for anything drawing the palette
// somewhere that has no terminal — a design system, a component sheet, a
// screenshot in the docs.
//
// Computed rather than tabled, because the 256-colour palette IS a formula for
// everything above index 15: a 6×6×6 cube, then a 24-step grey ramp. A table
// of nine hand-copied hex values would be nine chances to copy one wrong, and
// wrong here means a design system that quietly describes a different tool.
// An AdaptiveColor answers with its DARK value. Everything that renders a
// palette outside a terminal — a design system page, a captured frame turned
// into HTML — draws on a dark ground, because that is what the frame was
// captured for; recolouring it would report a tool that does not exist. A
// palette page for a light interface is a real thing to want and is not this.
func Hex(c lipgloss.TerminalColor) string {
	var s string
	switch v := c.(type) {
	case lipgloss.Color:
		s = string(v)
	case lipgloss.AdaptiveColor:
		s = v.Dark
	case lipgloss.CompleteColor:
		s = v.TrueColor
	case lipgloss.CompleteAdaptiveColor:
		s = v.Dark.TrueColor
	case nil:
		return "#000000"
	default:
		return "#000000"
	}
	// A tool that named its roles in hex has already answered the question.
	if strings.HasPrefix(s, "#") {
		return normalise(s)
	}

	i, err := strconv.Atoi(s)
	if err != nil || i < 0 || i > 255 {
		return "#000000"
	}
	switch {
	case i < 16:
		// The only part that is a table: the first sixteen are the terminal's
		// own, and predate any formula. These are xterm's defaults, which is
		// what the palette means when nothing has been reconfigured.
		return base16[i]
	case i < 232:
		// The cube: index = 16 + 36r + 6g + b, each channel one of six levels.
		n := i - 16
		return rgb(cubeLevel[n/36], cubeLevel[(n/6)%6], cubeLevel[n%6])
	default:
		// The grey ramp: 24 steps from 8 to 238, ten apart.
		g := 8 + 10*(i-232)
		return rgb(g, g, g)
	}
}

// cubeLevel is the six values a channel takes in the colour cube. Not evenly
// spaced: the gap from nothing to the first step is larger, so dark colours
// stay distinguishable.
var cubeLevel = [6]int{0, 95, 135, 175, 215, 255}

var base16 = [16]string{
	"#000000", "#800000", "#008000", "#808000",
	"#000080", "#800080", "#008080", "#c0c0c0",
	"#808080", "#ff0000", "#00ff00", "#ffff00",
	"#0000ff", "#ff00ff", "#00ffff", "#ffffff",
}

func rgb(r, g, b int) string { return fmt.Sprintf("#%02x%02x%02x", r, g, b) }

// normalise expands the three-digit form, so #fff and #ffffff mean the same
// thing to anything reading the palette.
func normalise(hex string) string {
	if len(hex) != 4 {
		return strings.ToLower(hex)
	}
	var b strings.Builder
	b.WriteByte('#')
	for _, r := range hex[1:] {
		b.WriteRune(r)
		b.WriteRune(r)
	}
	return strings.ToLower(b.String())
}

// Value is a role's colour as the tool DECLARED it — "205", "#d2a8ff", or a
// light/dark pair.
//
// A design system page should say what was written, not only what it resolves
// to: "205" tells a reader the palette is ANSI indices and will follow their
// terminal's own scheme, where "#ff5faf" tells them it will not. Hex answers
// what it looks like; this answers what it is.
func Value(c lipgloss.TerminalColor) string {
	switch v := c.(type) {
	case lipgloss.Color:
		return string(v)
	case lipgloss.AdaptiveColor:
		return v.Light + " / " + v.Dark
	case lipgloss.CompleteColor:
		return v.TrueColor
	case lipgloss.CompleteAdaptiveColor:
		return v.Light.TrueColor + " / " + v.Dark.TrueColor
	case nil:
		return ""
	}
	return ""
}
