// Package ui is democtl's interface.
//
// It is written the way a tuikit tool is written today: every colour comes from
// a theme.Palette, every non-ASCII character from a theme.GlyphSet, and
// guard_test.go holds both closed. There is no component library yet, so the
// panes here are hand-drawn — and that is deliberate. tuikit's comp package will
// be extracted from working code rather than designed in the abstract, and this
// is some of the code it will be extracted from.
package ui

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/theme"
)

// Palette and Glyphs are democtl's vocabulary.
//
// The palette is tuikit's default unchanged, because democtl has no reason to
// disagree with it. The glyph set adds nothing, which is the more interesting
// half: every box below is drawn with the six characters the default allows,
// and the tee and cross pieces a nicer border would want are simply not
// available. A design that needs a shape the set does not have needs a
// different design, not a different font.
var (
	Palette = theme.Default
	Glyphs  = theme.DefaultGlyphs
)

// styles is every style democtl draws with, built once from the palette.
//
// Built rather than declared at package level so a different palette produces a
// different interface without touching this file — which is what makes the
// palette a knob rather than a suggestion.
type styles struct {
	title    lipgloss.Style
	muted    lipgloss.Style
	border   lipgloss.Style
	focused  lipgloss.Style
	selected lipgloss.Style
	success  lipgloss.Style
	pending  lipgloss.Style
	danger   lipgloss.Style
	stderr   lipgloss.Style
}

func newStyles(p theme.Palette) styles {
	return styles{
		title:    lipgloss.NewStyle().Foreground(p.Accent).Bold(true),
		muted:    lipgloss.NewStyle().Foreground(p.Muted),
		border:   lipgloss.NewStyle().Foreground(p.Border),
		focused:  lipgloss.NewStyle().Foreground(p.Accent),
		selected: lipgloss.NewStyle().Foreground(p.SelectionFG).Background(p.SelectionBG).Bold(true),
		success:  lipgloss.NewStyle().Foreground(p.Success),
		pending:  lipgloss.NewStyle().Foreground(p.Pending),
		danger:   lipgloss.NewStyle().Foreground(p.Danger),
		stderr:   lipgloss.NewStyle().Foreground(p.Stderr),
	}
}
