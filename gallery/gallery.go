// Package gallery is every component in comp, running.
//
// Deliberately not a static sheet. A sheet tells you what exists; a running
// gallery tells you what it feels like to arrow through, what it does at 80
// columns, and what it looks like when the list is empty. Both a person picking
// a component and an agent deciding how to build a screen start here.
//
// It is built OUT OF comp, which is the cheapest honesty available: a gallery
// that drew its own panes and lists could show a component that no longer
// works. Its own list is a comp.List, its own panes are comp.Panes, its state
// picker is a comp.Tabs and its footer is a comp.Bar — so the gallery failing
// to open is itself a test result.
package gallery

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/theme"
)

// Entry is one component.
//
// The metadata is not decoration. A component's keys, its mouse behaviours, the
// colour ROLES it draws with and the glyphs it needs are the things you have to
// know before choosing it, and they are exactly what a static sheet leaves out.
// The roles and glyphs are also what a tool has to be able to supply: a
// component needing a role a tool's palette does not name is a component that
// tool cannot use.
type Entry struct {
	Name string
	// Summary is one line: what it is for.
	Summary string
	// From records which tools it was extracted from, so the gallery answers
	// "why does it work like that" as well as "what does it do".
	From string

	Keys   []comp.Hint
	Mouse  []comp.Hint
	Roles  []string
	Glyphs []string

	// States are the ways it can look. Every entry has more than one, because
	// the happy path is the state that never needed a gallery.
	States []State
}

// State is one way a component can look, and how to draw it.
//
// Draw takes the rect it may use and whether the gallery's preview has focus,
// so a component whose selection looks different when focused can show both
// without a second entry.
type State struct {
	Name string
	// Note says what this state is FOR — the condition that produces it —
	// because "empty" is obvious and "clamped" is not.
	Note string
	Draw func(c *comp.Canvas, r comp.Rect, focused bool)
}

// styles is the gallery's own vocabulary, built from the default palette.
// tuikit's own tooling draws with tuikit's own roles.
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
