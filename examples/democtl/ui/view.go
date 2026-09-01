package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/app"
	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
)

var tabNames = []string{"Overview", "Config", "Events"}

// View draws the whole window.
//
// Everything below draws CELLS into a canvas rather than returning strings to
// be joined. That is what makes a click land somewhere: every cell records who
// drew it, so the frame is its own map. It also deletes a category of bug —
// there is no way to draw past the edge, because the coordinate does not exist.
//
// The canvas is one row shorter than the terminal, which is where the old
// bodyHeight arithmetic already put it: two rows of header, the body, one of
// footer, and the bottom line left alone.
func (m *Model) View() string {
	c := comp.NewCanvas(m.width, 3+m.bodyHeight())

	m.header(c)
	// The default case used to be a hand-written complaint. app.Screens owns
	// it now, so every tool's missing screen says the same thing in the same
	// place — and guard.Screens can ask the question of the map.
	m.screens.Draw(app.Screen(m.screen), c,
		comp.Rect{X: 0, Y: 2, W: m.width, H: m.bodyHeight()}, &m.sty.danger)
	m.footer(c)
	if m.menu != nil {
		m.drawMenu(c)
	}
	if m.confirm != nil {
		m.modal(c)
	}

	// Kept so a mouse event can ask what it landed on. The frame IS the region
	// list, so there is nothing else to keep in step with it.
	m.canvas = c
	return c.String()
}

func (m *Model) header(c *comp.Canvas) {
	state, style := "all services running", &m.sty.success
	if !m.fleet.Healthy() {
		state, style = "needs attention", &m.sty.pending
	}
	comp.Bar{
		Left: []comp.Segment{
			{Text: "democtl", Style: &m.sty.title},
			{Text: "  a tuikit example", Style: &m.sty.muted},
			{Text: "  "},
			{Text: state, Style: style},
		},
		Right: []comp.Segment{{Text: m.now.Format("15:04:05"), Style: &m.sty.muted}},
	}.Draw(c, comp.Rect{X: 0, Y: 0, W: m.width, H: 1}, comp.Region(regHeader))

	c.Fill(comp.Rect{X: 0, Y: 1, W: m.width, H: 1}, "─", &m.sty.border, comp.Region(regRule))
}

// footer names the keys that act on WHAT IS FOCUSED, right now.
//
// swarmctl learned the rule the expensive way and wrote it down: its footer
// used to list every action on every panel, which grew a letter per feature and
// read as a menu of things mostly not applicable. And while something is
// capturing keys, the keys it is not taking do not belong here — listing them
// is a lie.
func (m *Model) footer(c *comp.Canvas) {
	var hints []comp.Hint
	switch {
	case m.confirm != nil:
		hints = []comp.Hint{{Key: "y", Label: "confirm"}, {Key: "n", Label: "cancel"}}
	case m.menu != nil:
		hints = []comp.Hint{{Key: "↑↓", Label: "choose"}, {Key: "enter", Label: "do it"}, {Key: "esc", Label: "close"}}
	case m.typing:
		hints = []comp.Hint{{Label: "type to filter"}, {Key: "enter", Label: "keep"}, {Key: "esc", Label: "clear"}}
	case m.screen == screenLogs:
		hints = []comp.Hint{{Key: "↑↓", Label: "scroll"}, {Key: "e", Label: "stderr only"}, {Key: "esc", Label: "back"}, {Key: "q", Label: "quit"}}
	case m.screen == screenRun:
		hints = []comp.Hint{{Key: "r", Label: "run again"}, {Key: "esc", Label: "back"}, {Key: "q", Label: "quit"}}
	default:
		hints = []comp.Hint{
			{Key: "↑↓", Label: "move"}, {Key: "tab", Label: "pane"}, {Key: "‹›", Label: "tabs"},
			{Key: "/", Label: "filter"}, {Key: "m", Label: "menu"},
			{Key: "L", Label: "logs"}, {Key: "D", Label: "deploy"}, {Key: "q", Label: "quit"},
		}
	}
	comp.KeyHints(c, comp.Rect{X: 0, Y: 2 + m.bodyHeight(), W: m.width, H: 1},
		comp.Region(regFooter), &m.sty.muted, hints...)
}

func (m *Model) dashboard(c *comp.Canvas, r comp.Rect) {
	listWidth := m.listWidth()
	m.servicePane(c, comp.Rect{X: r.X, Y: r.Y, W: listWidth, H: r.H})

	// The divider owns its column so a drag has something to grab. It draws a
	// blank, and an owned blank is still trimmed from the output, so naming it
	// costs the frame nothing.
	c.Fill(comp.Rect{X: r.X + listWidth, Y: r.Y, W: 1, H: r.H}, " ", nil, comp.Region(regSplit))

	m.detailPane(c, comp.Rect{X: r.X + listWidth + 1, Y: r.Y, W: r.W - listWidth - 1, H: r.H})
}

// listWidth is the divider's position: a third of the window until someone
// drags it, and then wherever they left it.
func (m *Model) listWidth() int {
	if m.split > 0 {
		return clamp(m.split, minPane, max(minPane, m.width-minPane-1))
	}
	return max(m.width/3, minPane)
}

func (m *Model) servicePane(c *comp.Canvas, r comp.Rect) {
	title := fmt.Sprintf("Services (%d)", len(m.visible()))
	if m.filter != "" || m.typing {
		title = "Filter: " + m.filter
		if m.typing {
			title += "_"
		}
	}
	inner := m.box(c, r, title, m.focus == paneList, comp.Region(regServices))

	list := m.visible()
	rows := make([]comp.Row, len(list))
	for i, s := range list {
		rows[i] = comp.Row{
			Text: fmt.Sprintf(" %-14s %d/%d %s",
				comp.Truncate(s.Name, 14), s.Ready, s.Want, mark(s.State)),
			Style: m.stateStyle(s.State),
		}
	}
	m.list.Focused = m.focus == paneList
	m.list.Draw(c, inner, rows)
}

func (m *Model) detailPane(c *comp.Canvas, r comp.Rect) {
	id := comp.Region(regDetail)
	svc, ok := m.selected()
	if !ok {
		inner := m.box(c, r, "Detail", m.focus == paneDetail, id)
		c.Text(inner.X, inner.Y, "  no service selected", &m.sty.muted, id)
		return
	}
	inner := m.box(c, r, svc.Name, m.focus == paneDetail, id)

	m.tabStrip(c, comp.Rect{X: inner.X, Y: inner.Y, W: inner.W, H: 1})
	y := inner.Y + 2

	switch m.tab {
	case 0:
		y = m.field(c, inner, y, "state", svc.State.String(), m.stateStyle(svc.State))
		y = m.field(c, inner, y, "replicas", fmt.Sprintf("%d/%d", svc.Ready, svc.Want), nil)
		y = m.field(c, inner, y, "stack", svc.Stack, nil)
		y = m.field(c, inner, y, "node", svc.Node, nil)
		y = m.field(c, inner, y, "updated", ago(m.now.Sub(svc.Updated)), nil)
		if svc.Note != "" {
			c.Text(inner.X, y+1, "  "+comp.WrapFirst(svc.Note, inner.W-2), &m.sty.muted, id)
		}
	case 1:
		y = m.field(c, inner, y, "image", comp.Truncate(svc.Image, inner.W-12), nil)
		y = m.field(c, inner, y, "desired", fmt.Sprintf("%d", svc.Want), nil)
		y = m.field(c, inner, y, "restart", "on-failure", nil)
		m.field(c, inner, y, "order", "start-first", nil)
	case 2:
		for i, l := range m.fleet.Logs(svc.Name, 6) {
			if y+i > inner.Y+inner.H-1 {
				break
			}
			x := c.Text(inner.X, y+i, "  "+l.At.Format("15:04:05")+" ", &m.sty.muted, id)
			c.Text(inner.X+x, y+i, comp.Truncate(l.Text, inner.W-12), nil, id)
		}
	}
}

// field draws one name/value row and returns the next row.
func (m *Model) field(c *comp.Canvas, inner comp.Rect, y int, name, value string, style *lipgloss.Style) int {
	if y > inner.Y+inner.H-1 {
		return y + 1
	}
	id := comp.Region(regDetail)
	x := c.Text(inner.X, y, fmt.Sprintf("  %-10s ", name), &m.sty.muted, id)
	c.Text(inner.X+x, y, value, style, id)
	return y + 1
}

// tabStrip draws the tabs.
//
// The chevrons now wrap the strip rather than the current tab, which is
// swarmctl's arrangement and the better one: the current tab is already
// coloured, so chevrons around it repeat what the colour says, while chevrons
// around the strip say that ‹ and › cycle it — which nothing else on screen
// does.
func (m *Model) tabStrip(c *comp.Canvas, r comp.Rect) {
	tabs := make([]comp.Tab, len(tabNames))
	for i, name := range tabNames {
		tabs[i] = comp.Tab{Name: name}
	}
	comp.Tabs{
		Tabs:          tabs,
		Active:        m.tab,
		Focused:       m.focus == paneDetail,
		Style:         &m.sty.muted,
		Selected:      &m.sty.focused,
		FocusSelected: &m.sty.title,
		Chrome:        &m.sty.muted,
	}.Draw(c, comp.Rect{X: r.X + 1, Y: r.Y, W: r.W - 1, H: 1}, regDetailTab)
}

func (m *Model) logs(c *comp.Canvas, r comp.Rect) {
	svc, ok := m.selected()
	if !ok {
		return
	}
	lines := m.fleet.Logs(svc.Name, 40)
	if m.logStderr {
		var only []fleet.LogLine
		for _, l := range lines {
			if l.Stderr {
				only = append(only, l)
			}
		}
		lines = only
	}
	rows := make([]comp.LogLine, len(lines))
	for i, l := range lines {
		rows[i] = comp.LogLine{At: l.At.Format("15:04:05"), Text: l.Text, Stderr: l.Stderr}
	}

	title := "Logs " + svc.Name
	if m.logStderr {
		title += " (stderr only)"
	}
	inner := m.box(c, r, title, true, comp.Region(regLogs))

	m.log.Empty = "  no lines on this stream"
	m.log.Time = &m.sty.muted
	m.log.Stderr = &m.sty.stderr
	m.log.Status = &m.sty.muted
	m.log.EmptyStyle = &m.sty.muted
	m.log.Draw(c, inner, rows, regLogsRow)
}

func (m *Model) run(c *comp.Canvas, r comp.Rect) {
	inner := m.box(c, r, m.plan.Name, true, comp.Region(regRun))

	steps := make([]comp.Step, len(m.plan.Steps))
	for i, step := range m.plan.Steps {
		s := comp.Step{Label: step.Name, State: stepStates[m.done[i]]}
		switch m.done[i] {
		case stepSkipped:
			s.Detail, s.Took = "already true", took(step.Took)
		case stepOK:
			s.Took = took(step.Took)
		case stepFailed:
			s.Note = "task 3 exited 137 before the health check passed"
		}
		steps[i] = s
	}

	list := comp.StepList{
		Steps: steps,
		Look: [5]comp.StepLook{
			comp.StepWaiting: {Glyph: "•", Style: &m.sty.muted, LabelStyle: &m.sty.muted},
			comp.StepRunning: {Glyph: "→", Style: &m.sty.pending, LabelStyle: &m.sty.pending},
			comp.StepSkipped: {Glyph: "●", Style: &m.sty.muted},
			comp.StepDone:    {Glyph: "✓", Style: &m.sty.success},
			comp.StepFailed:  {Glyph: "✗", Style: &m.sty.danger},
		},
		Muted: &m.sty.muted,
	}
	switch {
	case m.running >= 0:
		list.Status, list.StatusStyle = "running…", &m.sty.pending
	case len(m.done) > 0 && m.done[len(m.done)-1] == stepOK:
		list.Status, list.StatusStyle = "done", &m.sty.success
	default:
		list.Status, list.StatusStyle = "stopped", &m.sty.danger
		list.Hints = []comp.Hint{{Key: "r", Label: "to run again"}, {Key: "esc", Label: "to go back"}}
	}
	list.Draw(c, inner, regRunStep)
}

// stepStates maps democtl's step states onto the component's. A table rather
// than matching integers, so renumbering either side is a compile error instead
// of a silently wrong badge.
var stepStates = map[stepState]comp.StepState{
	stepWaiting: comp.StepWaiting,
	stepRunning: comp.StepRunning,
	stepSkipped: comp.StepSkipped,
	stepOK:      comp.StepDone,
	stepFailed:  comp.StepFailed,
}

// modal draws the confirm question over whatever is behind it.
func (m *Model) modal(c *comp.Canvas) {
	comp.Confirm{
		Title:       m.confirm.title,
		Body:        m.confirm.body,
		Danger:      m.confirm.danger,
		Hints:       []comp.Hint{{Key: "y", Label: "confirm"}, {Key: "n", Label: "cancel"}},
		Border:      &m.sty.focused,
		TitleStyle:  &m.sty.title,
		DangerStyle: &m.sty.danger,
		BodyStyle:   &m.sty.muted,
		HintStyle:   &m.sty.muted,
	}.Draw(c, comp.Region(regConfirm))
}

// drawMenu puts the context menu on top.
//
// Eight lines, and there is no compositing step. On a string frame this was 28
// lines that re-measured every line underneath and spliced the menu into it —
// and the version that replaced whole lines instead punched a hole through the
// panes. Drawn last is on top; that is the entire implementation.
func (m *Model) drawMenu(c *comp.Canvas) {
	w := 4
	for _, item := range m.menu.items {
		w = max(w, comp.Width(item.Label)+comp.Width(item.Key)+6)
	}
	x, y := m.menu.x, m.menu.y
	if m.menu.onRegion {
		// Where the thing it belongs to is IN THIS FRAME. The panes are
		// already drawn, so the canvas can be asked — and a menu anchored to a
		// row that has since scrolled follows it rather than pointing at where
		// it used to be.
		if at, ok := c.Region(m.menu.on); ok {
			x, y = at.X+2, at.Y
		}
	}
	r := comp.Rect{X: x, Y: y, W: w, H: len(m.menu.items) + 2}

	// Nudged back on screen rather than clipped. The canvas would happily draw
	// half a menu off the edge — that is what it is for — but half a menu is a
	// list of actions you cannot read, which is a different thing from a pane
	// that is cut off.
	if bounds := c.Bounds(); true {
		r.X = clamp(r.X, 0, max(0, bounds.W-r.W))
		r.Y = clamp(r.Y, 0, max(0, bounds.H-r.H))
	}
	inner := m.box(c, r, "", true, comp.Region(regMenu))

	for i, item := range m.menu.items {
		id := comp.Region(regMenuItem).At(i)
		style := &m.sty.muted
		if i == m.menu.cur {
			style = &m.sty.selected
		}
		row := comp.Rect{X: inner.X, Y: inner.Y + i, W: inner.W, H: 1}
		c.Fill(row, " ", style, id)
		c.Text(row.X+1, row.Y, item.Label, style, id)
		// The key sits beside the action rather than in a manual somewhere,
		// because it is the same list.
		c.Text(row.X+row.W-comp.Width(item.Key)-1, row.Y, item.Key, style, id)
	}
}

// box draws one of democtl's panes.
//
// A thin wrapper over comp.Pane rather than a call at each site, because the
// styles are the tool's and the placement is a house rule: democtl puts its
// titles on a row of their own so the focus highlight belongs to the pane.
func (m *Model) box(c *comp.Canvas, r comp.Rect, title string, focused bool, id comp.ID) comp.Rect {
	return comp.Pane{
		Title:      title,
		TitleAt:    comp.TitleOnRow,
		Focused:    focused,
		Border:     &m.sty.border,
		Focus:      &m.sty.focused,
		TitleStyle: &m.sty.title,
	}.Draw(c, r, id)
}

// bodyHeight is what is left after the header's two lines and the footer's one.
func (m *Model) bodyHeight() int { return max(3, m.height-4) }

// stateStyle is a pointer into the model's styles rather than a copy, so every
// row in a state shares one address and serialising groups them into a single
// run instead of one escape sequence per row.
func (m *Model) stateStyle(s fleet.State) *lipgloss.Style {
	switch s {
	case fleet.Running:
		return &m.sty.success
	case fleet.Pending, fleet.Degraded:
		return &m.sty.pending
	case fleet.Failed:
		return &m.sty.danger
	}
	return &m.sty.muted
}

func mark(s fleet.State) string {
	switch s {
	case fleet.Running:
		return "✓"
	case fleet.Pending:
		return "•"
	case fleet.Degraded:
		return "●"
	case fleet.Failed:
		return "✗"
	}
	return " "
}
