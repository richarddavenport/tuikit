package app

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/comp"
)

// Screen is one full-window view.
type Screen int

// View draws one screen into the space it is given.
type View func(c *comp.Canvas, r comp.Rect)

// Screens is the router.
//
// # Why a map and not a switch
//
// A switch with a missing case renders an empty terminal and says nothing about
// why. This draws a visible complaint naming the screen, which is what democtl
// already does by hand — and guard.Screens (#18) turns the same question into a
// test, so a screen constant added without a view fails a build rather than
// showing somebody a blank window and letting them wonder whether it hung.
//
// The complaint is deliberately ugly. A missing screen is a bug, and a bug that
// looks like a design decision gets shipped.
type Screens map[Screen]View

// Draw renders one screen, or says what is missing.
func (s Screens) Draw(screen Screen, c *comp.Canvas, r comp.Rect, style *lipgloss.Style) {
	if view, ok := s[screen]; ok && view != nil {
		view(c, r)
		return
	}
	c.Text(r.X, r.Y, "screen "+itoa(int(screen))+" has no View — see app.Screens",
		style, comp.Region("app.missing"))
}

// Has reports whether a screen has a view, for a tool's own exhaustiveness test
// until guard.Screens exists.
func (s Screens) Has(screen Screen) bool {
	view, ok := s[screen]
	return ok && view != nil
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	neg := n < 0
	if neg {
		n = -n
	}
	var d [20]byte
	i := len(d)
	for n > 0 {
		i--
		d[i] = byte('0' + n%10)
		n /= 10
	}
	if neg {
		i--
		d[i] = '-'
	}
	return string(d[i:])
}
