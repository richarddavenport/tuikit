package theme

import (
	"fmt"
	"strconv"

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
func Hex(c lipgloss.Color) string {
	i, err := strconv.Atoi(string(c))
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
