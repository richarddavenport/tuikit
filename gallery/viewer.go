package gallery

import (
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/comp"
)

// viewerEntry is the component seven of the fourteen surveyed tools each built
// alone, so its states are the ones those tools argued about: a cursor over
// syntax, a range, a line too wide for the pane, and no cursor at all.
func (m *Model) viewerEntry(s *styles) Entry {
	// A diff, because it is the case that made the cursor rule matter: the line
	// under the cursor is the line being read.
	added := lipgloss.NewStyle().Foreground(lipgloss.Color("2"))
	removed := lipgloss.NewStyle().Foreground(lipgloss.Color("1"))
	diff := []comp.Line{
		{Text: "@@ -12,7 +12,9 @@ func (l *List) Extend(by int) {", Style: &s.muted},
		{Text: " func (l *List) Extend(by int) {"},
		{Text: "-\tif !l.anchored {", Style: &removed},
		{Text: "-\t\tl.anchor, l.anchored = l.cursor, true", Style: &removed},
		{Text: "-\t}", Style: &removed},
		{Text: "+\tl.sel.start(l.cursor)", Style: &added},
		{Text: " \tl.pending += by"},
		{Text: " \tl.reveal = true"},
		{Text: " }"},
	}

	// One line, coloured in pieces, so the cursor can be shown sitting on it
	// without eating the colours.
	code := []comp.Line{
		{Spans: []comp.Segment{
			{Text: "func ", Style: &s.title}, {Text: "under"},
			{Text: "(base, s ", Style: &s.muted}, {Text: "*lipgloss.Style"},
			{Text: ") *lipgloss.Style {", Style: &s.muted},
		}},
		{Spans: []comp.Segment{
			{Text: "\tif ", Style: &s.title}, {Text: "base == "},
			{Text: "nil", Style: &s.pending}, {Text: " {"},
		}},
		{Spans: []comp.Segment{
			{Text: "\t\treturn ", Style: &s.title}, {Text: "s"},
		}},
		{Text: "\t}"},
		{Spans: []comp.Segment{
			{Text: "\tmerged := s.", Style: nil}, {Text: "Inherit", Style: &s.focused},
			{Text: "(*base)"},
		}},
	}

	view := func(prepare func(*comp.Viewer), lines []comp.Line, numbers bool) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			v := &comp.Viewer{
				Name: "demo.viewer", Focused: true, Numbers: numbers, Tab: 4,
				Selected: &s.selected, Ranged: &s.focused,
				Status: &s.muted, Number: &s.muted, EmptyStyle: &s.muted,
			}
			if prepare != nil {
				v.Draw(c, r, lines) // once, so it knows what it can scroll
				prepare(v)
			}
			v.Draw(c, r, lines)
		}
	}

	long := []comp.Line{
		{Text: `{"name":"api-gateway","image":"registry.example.com/estate/api-gateway:2026.9.1","replicas":3}`},
		{Text: `{"name":"worker","image":"registry.example.com/estate/worker:2026.9.1","replicas":12}`},
		{Text: `{"name":"scheduler","image":"registry.example.com/estate/scheduler:2026.8.7","replicas":1}`},
	}

	return Entry{
		Name: "Viewer",
		Summary: "A document you scroll, look through, and act on ranges of. " +
			"Opens at the top, does not follow, and keeps its colours under the cursor.",
		From: "seven of the fourteen tools surveyed each built this and none could get it " +
			"from their toolkit — lazygit, gitui, k9s, fx, termshark, dive, gh-dash",
		Keys: []comp.Hint{
			{Key: "↑/↓", Label: "move the cursor"},
			{Key: "shift+↑/↓", Label: "extend a range"},
			{Key: "←/→", Label: "scroll sideways"},
			{Key: "g", Label: "go to a line"},
		},
		Mouse:  []comp.Hint{{Key: "wheel", Label: "scroll without moving the cursor"}},
		Roles:  []string{"SelectionBG", "Accent", "Muted"},
		Glyphs: []string{},
		States: []State{
			{Name: "a diff", Note: "the cursor is a background and the line keeps its own colour — a List would repaint the whole row and lose it",
				Draw: view(func(v *comp.Viewer) { v.Goto(5) }, diff, true)},
			{Name: "a range", Note: "shift-arrow, for staging a hunk or copying a span; ordered, so the caller never asks which way it was dragged",
				Draw: view(func(v *comp.Viewer) { v.Goto(2); v.Extend(2) }, diff, true)},
			{Name: "syntax under the cursor", Note: "spans keep their styles and take the cursor's background — the one rule that makes this not a List",
				Draw: view(func(v *comp.Viewer) { v.Goto(1) }, code, false)},
			{Name: "wider than the pane", Note: "a document is not as wide as its pane; truncating it means the content cannot be read at all",
				Draw: view(nil, long, false)},
			{Name: "scrolled sideways", Note: "the cut lands inside a span rather than on its boundary, so the offset does not jump by a token",
				Draw: view(func(v *comp.Viewer) { v.ScrollX(46) }, long, false)},
			{Name: "no cursor", Note: "k9s's YAML view and termshark's scrollabletext are both this — the arrows still scroll",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					v := &comp.Viewer{Name: "demo.viewer", NoCursor: true, Numbers: true, Status: &s.muted, Number: &s.muted}
					v.Draw(c, r, yamlLines(s))
				}},
			{Name: "a long document", Note: "lines are asked for one at a time and only when visible, so the frame costs the pane rather than the file",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					v := &comp.Viewer{
						Name: "demo.viewer", Focused: true, Numbers: true,
						Selected: &s.selected, Status: &s.muted, Number: &s.muted,
					}
					line := func(i int) comp.Line { return comp.Line{Text: "  entry " + itoa(i+1)} }
					v.DrawFunc(c, r, 1_000_000, line)
					v.Goto(654_321)
					v.DrawFunc(c, r, 1_000_000, line)
				}},
			{Name: "nothing to show", Note: "an empty document is an ordinary state — a filter matching nothing, or a file that is genuinely empty",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					v := &comp.Viewer{Name: "demo.viewer", Empty: "  no differences", EmptyStyle: &s.muted, Status: &s.muted}
					v.Draw(c, r, nil)
				}},
			{Name: "no room", Note: "the state every component gets wrong first",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					v := &comp.Viewer{Name: "demo.viewer", Numbers: true, Status: &s.muted}
					v.Draw(c, comp.Rect{X: r.X, Y: r.Y, W: r.W, H: 1}, diff)
				}},
		},
	}
}

func yamlLines(s *styles) []comp.Line {
	src := `apiVersion: apps/v1
kind: Deployment
metadata:
  name: api-gateway
  namespace: estate
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: api-gateway
          image: registry.example.com/estate/api-gateway:2026.9.1`
	var out []comp.Line
	for _, raw := range strings.Split(src, "\n") {
		key, rest, found := strings.Cut(raw, ":")
		if !found {
			out = append(out, comp.Line{Text: raw})
			continue
		}
		out = append(out, comp.Line{Spans: []comp.Segment{
			{Text: key + ":", Style: &s.title}, {Text: rest},
		}})
	}
	return out
}
