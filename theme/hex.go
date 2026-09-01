package theme

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Hex is what an ANSI index looks like, for anything drawing the palette
// somewhere that has no terminal — a design system, a component sheet, a
// screenshot in the docs.
//
// It answers "" for a colour it cannot read, and NOT "#000000": index 0 is a
// legitimate role now that the palette is built from the first sixteen, and a
// sentinel that collides with a real answer is a check that stops checking. The
// test that every role converts used to look for a black square and would now
// pass over a broken one.
//
// The first sixteen are the terminal's OWN, so the hex here is what a default
// xterm would draw and not what the reader will see. That is the honest answer
// for a page rendered outside a terminal, and it is why [Value] exists: a
// design system should say the palette is `13`, which follows the reader's
// theme, rather than implying it is a fixed pink.
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
		return ""
	default:
		return ""
	}
	// A tool that named its roles in hex has already answered the question.
	if strings.HasPrefix(s, "#") {
		return normalise(s)
	}

	i, err := strconv.Atoi(s)
	if err != nil || i < 0 || i > 255 {
		return ""
	}
	switch {
	case i < 16:
		// The only part that is a table: the first sixteen are the terminal's
		// own, and predate any formula. These are XTERM's defaults, which is
		// what the palette means when nothing has been reconfigured — and now
		// that the palette IS these sixteen, they are what every swatch on a
		// design system page comes from, so being the right table matters more
		// than it did. The set here used to be the VGA/"system colors" one
		// (#800000 red, #c0c0c0 white), which the comment already claimed was
		// xterm's and was not.
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

// base16 is xterm's own defaults for the sixteen, which are what a terminal
// draws before any theme touches it. Named the way xterm names them, because
// the values on their own look like typos.
var base16 = [16]string{
	"#000000", // 0  black
	"#cd0000", // 1  red3
	"#00cd00", // 2  green3
	"#cdcd00", // 3  yellow3
	"#0000ee", // 4  blue2
	"#cd00cd", // 5  magenta3
	"#00cdcd", // 6  cyan3
	"#e5e5e5", // 7  gray90
	"#7f7f7f", // 8  gray50
	"#ff0000", // 9  red
	"#00ff00", // 10 green
	"#ffff00", // 11 yellow
	"#5c5cff", // 12 rgb:5c/5c/ff
	"#ff00ff", // 13 magenta
	"#00ffff", // 14 cyan
	"#ffffff", // 15 white
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

// Value is a role's colour as the tool DECLARED it — "13", "#d2a8ff", or a
// light/dark pair.
//
// A design system page should say what was written, not only what it resolves
// to, because the two say different things. "13" is one of the sixteen the
// reader's terminal defines, so it FOLLOWS THEIR THEME; "205" is fixed in every
// terminal that has ever existed; "#ff5faf" is fixed and will not survive being
// looked at over ssh from a machine with a different profile.
//
// This comment used to say that any index follows the terminal's own scheme.
// That is wrong, and it is the kind of wrong that decides a design: everything
// from 16 up is a formula — a 6x6x6 cube and a grey ramp — that no theme
// touches. The line is not index versus hex. It is fifteen.
//
// Hex answers what it looks like; this answers what it is.
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
