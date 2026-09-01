package ui

import (
	"fmt"

	"github.com/charmbracelet/lipgloss"

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
	body := comp.Rect{X: 0, Y: 2, W: m.width, H: m.bodyHeight()}
	switch m.screen {
	case screenDashboard:
		m.dashboard(c, body)
	case screenLogs:
		m.logs(c, body)
	case screenRun:
		m.run(c, body)
	default:
		// A visible complaint rather than an empty string: a screen that
		// renders nothing looks like a hang, and the reader has no way to know
		// a case is missing.
		c.Text(body.X, body.Y, fmt.Sprintf("screen %d has no View case", m.screen), &m.sty.danger, comp.ID{})
	}
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
	id := comp.Region(regHeader)
	x := 0
	x += c.Text(x, 0, "democtl", &m.sty.title, id)
	x += c.Text(x, 0, "  a tuikit example", &m.sty.muted, id)
	x += c.Text(x, 0, "  ", nil, id)

	state, style := "all services running", &m.sty.success
	if !m.fleet.Healthy() {
		state, style = "needs attention", &m.sty.pending
	}
	x += c.Text(x, 0, state, style, id)

	// The clock goes right, and only if it does not collide. Measured in
	// columns by the canvas rather than by len — this is the bug that put a
	// header off the side of the screen in a real tool.
	clock := m.now.Format("15:04:05")
	if at := m.width - comp.Width(clock); at > x {
		c.Text(at, 0, clock, &m.sty.muted, id)
	}
	c.Fill(comp.Rect{X: 0, Y: 1, W: m.width, H: 1}, "─", &m.sty.border, comp.Region(regRule))
}

func (m *Model) footer(c *comp.Canvas) {
	var keys string
	switch {
	case m.confirm != nil:
		keys = "y confirm · n cancel"
	case m.menu != nil:
		keys = "↑↓ choose · enter do it · esc close"
	case m.typing:
		keys = "type to filter · enter keep · esc clear"
	case m.screen == screenLogs:
		keys = "↑↓ scroll · e stderr only · esc back · q quit"
	case m.screen == screenRun:
		keys = "r run again · esc back · q quit"
	default:
		keys = "↑↓ move · tab pane · ‹› tabs · / filter · m menu · L logs · D deploy · q quit"
	}
	c.Text(0, 2+m.bodyHeight(), comp.Truncate("  "+keys, m.width), &m.sty.muted, comp.Region(regFooter))
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
	if len(list) == 0 {
		// Painted the full width, like any other row, so the empty state is a
		// row of the list rather than a sentence floating in the pane.
		empty := comp.Rect{X: inner.X, Y: inner.Y, W: inner.W, H: 1}
		c.Fill(empty, " ", &m.sty.muted, comp.Region(regServices))
		c.Text(empty.X, empty.Y, "  nothing matches", &m.sty.muted, comp.Region(regServices))
		return
	}
	// Both offsets clamp against what was actually drawn last frame rather
	// than against a constant, which is how a five-line pane ends up scrolled
	// into empty space.
	m.listMax = max(0, len(list)-inner.H)
	m.listOffset = clamp(m.listOffset, 0, m.listMax)

	for i := m.listOffset; i < len(list); i++ {
		s := list[i]
		y := inner.Y + i - m.listOffset
		if y > inner.Y+inner.H-1 {
			break
		}
		// The owner carries the index into the LIST, not the row it landed on.
		// They differ the moment the viewport moves, and an ID that means "row
		// 3 of the screen" then acts on whatever scrolled into row 3.
		id := comp.Region(regServicesRow).At(i)
		style := m.stateStyle(s.State)
		switch {
		case i == m.cur && m.focus == paneList:
			style = &m.sty.selected
		case i == m.cur:
			style = &m.sty.focused
		}
		row := comp.Rect{X: inner.X, Y: y, W: inner.W, H: 1}
		c.Fill(row, " ", style, id)
		c.Text(row.X, row.Y, fmt.Sprintf(" %-14s %d/%d %s",
			comp.Truncate(s.Name, 14), s.Ready, s.Want, mark(s.State)), style, id)
	}
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

// tabStrip draws the tabs. The angle quotes are in the glyph set; the tee
// pieces a nicer strip would want are not, so it does not have any.
//
// No padding arithmetic any more. The strip draws until it runs out of rect,
// and the canvas stops it there — which is what the old rune-counting pad got
// wrong at 80 columns, cutting the strip mid-escape and leaking bold down the
// row.
func (m *Model) tabStrip(c *comp.Canvas, r comp.Rect) {
	x := r.X + 1
	for i, name := range tabNames {
		id := comp.Region(regDetailTab).At(i)
		label := " " + name + " "
		if i == m.tab {
			x += c.Text(x, r.Y, "‹"+label+"›", &m.sty.title, id)
		} else {
			x += c.Text(x, r.Y, " "+label+" ", &m.sty.muted, id)
		}
	}
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
	if m.logOffset > max(0, len(lines)-r.H+2) {
		m.logOffset = max(0, len(lines)-r.H+2)
	}
	lines = lines[min(m.logOffset, len(lines)):]

	title := "Logs " + svc.Name
	if m.logStderr {
		title += " (stderr only)"
	}
	inner := m.box(c, r, title, true, comp.Region(regLogs))

	if len(lines) == 0 {
		c.Text(inner.X, inner.Y, "  no lines on this stream", &m.sty.muted, comp.Region(regLogs))
		return
	}
	for i, l := range lines {
		if i >= inner.H {
			break
		}
		id := comp.Region(regLogsRow).At(m.logOffset + i)
		x := c.Text(inner.X, inner.Y+i, "  "+l.At.Format("15:04:05")+" ", &m.sty.muted, id)

		var style *lipgloss.Style
		if l.Stderr {
			style = &m.sty.stderr
		}
		c.Text(inner.X+x, inner.Y+i, comp.Truncate(l.Text, m.width-14), style, id)
	}
}

func (m *Model) run(c *comp.Canvas, r comp.Rect) {
	inner := m.box(c, r, m.plan.Name, true, comp.Region(regRun))
	y := inner.Y

	for i, step := range m.plan.Steps {
		if y > inner.Y+inner.H-1 {
			return
		}
		id := comp.Region(regRunStep).At(i)

		var badge, label string
		var badgeStyle, labelStyle *lipgloss.Style
		switch m.done[i] {
		case stepOK:
			badge, badgeStyle, label = "✓", &m.sty.success, step.Name
		case stepSkipped:
			badge, badgeStyle, label = "●", &m.sty.muted, step.Name
		case stepFailed:
			badge, badgeStyle, label = "✗", &m.sty.danger, step.Name
		case stepRunning:
			badge, badgeStyle, label, labelStyle = "→", &m.sty.pending, step.Name, &m.sty.pending
		default:
			badge, badgeStyle, label, labelStyle = "•", &m.sty.muted, step.Name, &m.sty.muted
		}

		x := inner.X + c.Text(inner.X, y, "  ", nil, id)
		x += c.Text(x, y, badge, badgeStyle, id)
		x += c.Text(x, y, " ", nil, id)
		x += c.Text(x, y, label, labelStyle, id)
		if m.done[i] == stepSkipped {
			c.Text(x, y, "  already true", &m.sty.muted, id)
		}

		// The duration is right-aligned against the box INSIDE, not the
		// terminal. Measuring against the terminal put it two columns past the
		// right edge, where it clipped to "400m…" — visible in a captured
		// frame and in nothing else.
		if m.done[i] == stepOK || m.done[i] == stepSkipped {
			right := took(step.Took) + " "
			c.Text(inner.X+inner.W-comp.Width(right), y, right, &m.sty.muted, id)
		}
		y++

		if m.done[i] == stepFailed && y <= inner.Y+inner.H-1 {
			c.Text(inner.X, y, "      task 3 exited 137 before the health check passed", &m.sty.danger, id)
			y++
		}
	}

	y++
	if y > inner.Y+inner.H-1 {
		return
	}
	id := comp.Region(regRun)
	switch {
	case m.running >= 0:
		c.Text(inner.X, y, "  running…", &m.sty.pending, id)
	case len(m.done) > 0 && m.done[len(m.done)-1] == stepOK:
		c.Text(inner.X, y, "  done", &m.sty.success, id)
	default:
		x := c.Text(inner.X, y, "  stopped", &m.sty.danger, id)
		c.Text(inner.X+x, y, "  r to run again · esc to go back", &m.sty.muted, id)
	}
}

// modal draws the confirm box over whatever is behind it.
//
// Drawn LAST, so it is on top. There is no compositing step, no re-measuring of
// the lines beneath and nothing to get wrong — which is the whole of the bug
// that once wiped 8 of 18 framed rows.
func (m *Model) modal(c *comp.Canvas) {
	w := min(m.width-8, 64)
	title := &m.sty.title
	if m.confirm.danger {
		title = &m.sty.danger
	}

	body := comp.Wrap(m.confirm.body, w-4)
	h := len(body) + 6 // title, blank, body, blank, keys, and two of border

	r := comp.Rect{
		X: max(0, (m.width-w)/2),
		Y: max(0, (3+m.bodyHeight()-h)/2),
		W: w,
		H: h,
	}
	id := comp.Region(regConfirm)
	inner := m.box(c, r, "", true, id)

	c.Text(inner.X, inner.Y, comp.Truncate(m.confirm.title, inner.W), title, id)
	for i, line := range body {
		c.Text(inner.X, inner.Y+2+i, line, &m.sty.muted, id)
	}
	c.Text(inner.X, inner.Y+len(body)+3, "y confirm · n cancel", &m.sty.muted, id)
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
		w = max(w, comp.Width(item.label)+comp.Width(item.key)+6)
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
		c.Text(row.X+1, row.Y, item.label, style, id)
		// The key sits beside the action rather than in a manual somewhere,
		// because it is the same list.
		c.Text(row.X+row.W-comp.Width(item.key)-1, row.Y, item.key, style, id)
	}
}

// box draws a titled frame and returns the rect inside it.
//
// The inside is blanked, so a box drawn over something else covers it. That is
// what makes the modal a draw rather than a composite.
func (m *Model) box(c *comp.Canvas, r comp.Rect, title string, focused bool, id comp.ID) comp.Rect {
	edge := &m.sty.border
	if focused {
		edge = &m.sty.focused
	}
	c.Box(r, edge, id)

	inner := r.Inset(1)
	c.Fill(inner, " ", nil, comp.ID{})
	if title == "" {
		return inner
	}
	// The title row is painted in full, trailing spaces included, because the
	// style is the row's rather than the text's.
	c.Fill(comp.Rect{X: inner.X, Y: inner.Y, W: inner.W, H: 1}, " ", &m.sty.title, id)
	c.Text(inner.X, inner.Y, " "+comp.Truncate(title, inner.W-1), &m.sty.title, id)

	return comp.Rect{X: inner.X, Y: inner.Y + 1, W: inner.W, H: inner.H - 1}
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
