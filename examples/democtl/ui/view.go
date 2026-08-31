package ui

import (
	"fmt"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
)

var tabNames = []string{"Overview", "Config", "Events"}

// View draws the whole window.
//
// Every screen constant needs a case here. The default is a visible complaint
// rather than an empty string: a screen that renders nothing looks like a hang,
// and the reader has no way to know a case is missing.
func (m *Model) View() string {
	var body string
	switch m.screen {
	case screenDashboard:
		body = m.dashboard()
	case screenLogs:
		body = m.logs()
	case screenRun:
		body = m.run()
	default:
		body = m.sty.danger.Render(fmt.Sprintf("screen %d has no View case", m.screen))
	}

	out := m.header() + "\n" + body + "\n" + m.footer()
	if m.confirm != nil {
		out = m.overlay(out)
	}
	return out
}

func (m *Model) header() string {
	state := m.sty.success.Render("all services running")
	if !m.fleet.Healthy() {
		state = m.sty.pending.Render("needs attention")
	}
	left := m.sty.title.Render("democtl") + m.sty.muted.Render("  a tuikit example") + "  " + state

	right := m.sty.muted.Render(m.now.Format("15:04:05"))
	return spread(m.width, left, right) + "\n" + m.sty.border.Render(rule(m.width))
}

// spread puts left and right on one line, width columns apart, and gives the
// space to the left half when there is not enough for both.
//
// Measured with lipgloss.Width rather than len: the styles above are escape
// sequences, and len counts them as visible. This is the bug that put a header
// off the side of the screen in a real tool, found by looking at a captured
// frame rather than by any assertion.
func spread(w int, left, right string) string {
	gap := w - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		return clip(left, w)
	}
	return left + strings.Repeat(" ", gap) + right
}

func (m *Model) dashboard() string {
	listWidth := m.width / 3
	if listWidth < 24 {
		listWidth = 24
	}
	detailWidth := m.width - listWidth - 1

	list := m.servicePane(listWidth, m.bodyHeight())
	detail := m.detailPane(detailWidth, m.bodyHeight())
	return lipgloss.JoinHorizontal(lipgloss.Top, list, " ", detail)
}

func (m *Model) servicePane(w, h int) string {
	title := fmt.Sprintf("Services (%d)", len(m.visible()))
	if m.filter != "" || m.typing {
		title = "Filter: " + m.filter
		if m.typing {
			title += "_"
		}
	}

	var rows []string
	for i, s := range m.visible() {
		row := fmt.Sprintf(" %-14s %d/%d %s", trim(s.Name, 14), s.Ready, s.Want, mark(s.State))
		row = pad(row, w-2)
		switch {
		case i == m.cur && m.focus == paneList:
			row = m.sty.selected.Render(row)
		case i == m.cur:
			row = m.sty.focused.Render(row)
		default:
			row = m.stateStyle(s.State).Render(row)
		}
		rows = append(rows, row)
	}
	if len(rows) == 0 {
		rows = append(rows, m.sty.muted.Render(pad("  nothing matches", w-2)))
	}
	return m.box(title, rows, w, h, m.focus == paneList)
}

func (m *Model) detailPane(w, h int) string {
	svc, ok := m.selected()
	if !ok {
		return m.box("Detail", []string{m.sty.muted.Render("  no service selected")}, w, h, m.focus == paneDetail)
	}

	var rows []string
	rows = append(rows, m.tabStrip(w-2), "")
	switch m.tab {
	case 0:
		rows = append(rows,
			m.field("state", m.stateStyle(svc.State).Render(svc.State.String())),
			m.field("replicas", fmt.Sprintf("%d/%d", svc.Ready, svc.Want)),
			m.field("stack", svc.Stack),
			m.field("node", svc.Node),
			m.field("updated", ago(m.now.Sub(svc.Updated))))
		if svc.Note != "" {
			rows = append(rows, "", m.sty.muted.Render("  "+wrapFirst(svc.Note, w-4)))
		}
	case 1:
		rows = append(rows,
			m.field("image", trim(svc.Image, w-14)),
			m.field("desired", fmt.Sprintf("%d", svc.Want)),
			m.field("restart", "on-failure"),
			m.field("order", "start-first"))
	case 2:
		for _, l := range m.fleet.Logs(svc.Name, 6) {
			rows = append(rows, m.sty.muted.Render("  "+l.At.Format("15:04:05")+" ")+trim(l.Text, w-14))
		}
	}
	return m.box(svc.Name, rows, w, h, m.focus == paneDetail)
}

// tabStrip draws the tabs. The angle quotes are in the glyph set; the tee
// pieces a nicer strip would want are not, so it does not have any.
func (m *Model) tabStrip(w int) string {
	var parts []string
	for i, name := range tabNames {
		label := " " + name + " "
		if i == m.tab {
			parts = append(parts, m.sty.title.Render("‹"+label+"›"))
		} else {
			parts = append(parts, m.sty.muted.Render(" "+label+" "))
		}
	}
	return pad(" "+strings.Join(parts, ""), w)
}

func (m *Model) field(name, value string) string {
	return m.sty.muted.Render(fmt.Sprintf("  %-10s ", name)) + value
}

func (m *Model) logs() string {
	svc, ok := m.selected()
	if !ok {
		return ""
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

	h := m.bodyHeight()
	if m.logOffset > max(0, len(lines)-h+2) {
		m.logOffset = max(0, len(lines)-h+2)
	}
	lines = lines[min(m.logOffset, len(lines)):]

	var rows []string
	for _, l := range lines {
		text := l.Text
		if l.Stderr {
			text = m.sty.stderr.Render(text)
		}
		rows = append(rows, m.sty.muted.Render("  "+l.At.Format("15:04:05")+" ")+clip(text, m.width-14))
	}
	if len(rows) == 0 {
		rows = append(rows, m.sty.muted.Render("  no lines on this stream"))
	}

	title := "Logs " + svc.Name
	if m.logStderr {
		title += " (stderr only)"
	}
	return m.box(title, rows, m.width, h, true)
}

func (m *Model) run() string {
	var rows []string
	for i, step := range m.plan.Steps {
		var badge, label string
		switch m.done[i] {
		case stepOK:
			badge, label = m.sty.success.Render("✓"), step.Name
		case stepSkipped:
			badge, label = m.sty.muted.Render("●"), step.Name+m.sty.muted.Render("  already true")
		case stepFailed:
			badge, label = m.sty.danger.Render("✗"), step.Name
		case stepRunning:
			badge, label = m.sty.pending.Render("→"), m.sty.pending.Render(step.Name)
		default:
			badge, label = m.sty.muted.Render("•"), m.sty.muted.Render(step.Name)
		}
		row := fmt.Sprintf("  %s %s", badge, label)
		if m.done[i] == stepOK || m.done[i] == stepSkipped {
			// The width is the box's inside, not the terminal's. Measuring
			// against m.width here put the durations two columns past the
			// right edge, where they were clipped to "400m…" — visible in a
			// captured frame and in nothing else.
			row = spread(m.width-2, row, m.sty.muted.Render(took(step.Took)+" "))
		}
		rows = append(rows, row)
		if m.done[i] == stepFailed {
			rows = append(rows, m.sty.danger.Render("      task 3 exited 137 before the health check passed"))
		}
	}

	rows = append(rows, "")
	switch {
	case m.running >= 0:
		rows = append(rows, m.sty.pending.Render("  running…"))
	case len(m.done) > 0 && m.done[len(m.done)-1] == stepOK:
		rows = append(rows, m.sty.success.Render("  done"))
	default:
		rows = append(rows, m.sty.danger.Render("  stopped")+m.sty.muted.Render("  r to run again · esc to go back"))
	}
	return m.box(m.plan.Name, rows, m.width, m.bodyHeight(), true)
}

func (m *Model) footer() string {
	var keys string
	switch {
	case m.confirm != nil:
		keys = "y confirm · n cancel"
	case m.typing:
		keys = "type to filter · enter keep · esc clear"
	case m.screen == screenLogs:
		keys = "↑↓ scroll · e stderr only · esc back · q quit"
	case m.screen == screenRun:
		keys = "r run again · esc back · q quit"
	default:
		keys = "↑↓ move · tab pane · ‹› tabs · / filter · L logs · D deploy · q quit"
	}
	return m.sty.muted.Render(clip("  "+keys, m.width))
}

// overlay puts the modal over the frame, bounded to the terminal.
//
// The bound is the point. A description can be as long as the thing it
// describes, and a box with no maximum draws off the side of the screen — which
// is what happened in a real tool, and what a captured frame showed in a second.
func (m *Model) overlay(behind string) string {
	w := min(m.width-8, 64)
	title := m.sty.title
	if m.confirm.danger {
		title = m.sty.danger
	}

	body := wrap(m.confirm.body, w-4)
	rows := make([]string, 0, len(body)+4)
	rows = append(rows, title.Render(clip(m.confirm.title, w-4)), "")
	for _, line := range body {
		rows = append(rows, m.sty.muted.Render(line))
	}
	rows = append(rows, "", m.sty.muted.Render("y confirm · n cancel"))

	box := m.box("", rows, w, len(rows)+2, true)

	lines := strings.Split(behind, "\n")
	top := max(0, (len(lines)-lipgloss.Height(box))/2)
	left := max(0, (m.width-w)/2)

	for i, boxLine := range strings.Split(box, "\n") {
		if top+i < len(lines) {
			lines[top+i] = composite(lines[top+i], boxLine, left, m.width)
		}
	}
	return strings.Join(lines, "\n")
}

// box draws a titled frame with the six box-drawing characters the glyph set
// allows. There are no tee or cross pieces, so the title sits inside the top
// edge rather than breaking it.
func (m *Model) box(title string, rows []string, w, h int, focused bool) string {
	edge := m.sty.border
	if focused {
		edge = m.sty.focused
	}

	inner := w - 2
	var b strings.Builder
	b.WriteString(edge.Render("┌"+rule(inner)+"┐") + "\n")
	if title != "" {
		b.WriteString(edge.Render("│") + m.sty.title.Render(pad(" "+trim(title, inner-1), inner)) + edge.Render("│") + "\n")
	}
	body := h - 2
	if title != "" {
		body--
	}
	for i := range body {
		line := ""
		if i < len(rows) {
			line = rows[i]
		}
		b.WriteString(edge.Render("│") + padVisible(line, inner) + edge.Render("│") + "\n")
	}
	b.WriteString(edge.Render("└" + rule(inner) + "┘"))
	return b.String()
}

// bodyHeight is what is left after the header's two lines and the footer's one.
func (m *Model) bodyHeight() int { return max(3, m.height-4) }

func (m *Model) stateStyle(s fleet.State) lipgloss.Style {
	switch s {
	case fleet.Running:
		return m.sty.success
	case fleet.Pending:
		return m.sty.pending
	case fleet.Degraded:
		return m.sty.pending
	case fleet.Failed:
		return m.sty.danger
	}
	return m.sty.muted
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
