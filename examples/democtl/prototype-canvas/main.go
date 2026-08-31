// PROTOTYPE — throwaway. Run it:
//
//	go run ./examples/democtl/prototype-canvas
//
// Then use the mouse: click a row, click a tab, drag the divider, scroll the
// wheel over either pane, right-click a row. The bottom line prints the raw
// mouse event and what the canvas says is under the pointer, so the hit-testing
// is visible rather than inferred.
package main

import (
	"fmt"
	"os"
	"strings"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
	"github.com/richarddavenport/tuikit/theme"
)

// Owner IDs. Strings in the prototype; the real thing would make these typed
// and generated, so a menu could not name a region that does not exist.
const (
	ownerSplit = "split.main"
	ownerList  = "pane.list"
	ownerDetl  = "pane.detail"
	ownerMenu  = "menu"
)

func rowOwner(i int) string  { return fmt.Sprintf("list.row[%d]", i) }
func tabOwner(i int) string  { return fmt.Sprintf("detail.tab[%d]", i) }
func menuOwner(i int) string { return fmt.Sprintf("menu.item[%d]", i) }

type model struct {
	fleet fleet.Fleet
	sty   map[string]*lipgloss.Style

	w, h int

	ratio  float64 // left pane's share of the width
	sel    int
	tab    int
	focus  string
	scroll int // detail pane's log scroll, to prove wheel targeting

	dragging bool
	menu     []string
	menuAt   [2]int

	// last is surfaced on screen so the prototype shows its own state.
	last  string
	under string

	// canvas is kept from the last View so Update can hit-test against exactly
	// what the user is looking at.
	canvas *Canvas
}

func newModel() *model {
	p := theme.Default
	sty := map[string]*lipgloss.Style{}
	add := func(name string, s lipgloss.Style) { sty[name] = &s }
	add("title", lipgloss.NewStyle().Foreground(p.Accent).Bold(true))
	add("muted", lipgloss.NewStyle().Foreground(p.Muted))
	add("border", lipgloss.NewStyle().Foreground(p.Border))
	add("focus", lipgloss.NewStyle().Foreground(p.Accent))
	add("sel", lipgloss.NewStyle().Foreground(p.SelectionFG).Background(p.SelectionBG).Bold(true))
	add("ok", lipgloss.NewStyle().Foreground(p.Success))
	add("pending", lipgloss.NewStyle().Foreground(p.Pending))
	add("danger", lipgloss.NewStyle().Foreground(p.Danger))
	return &model{
		fleet: fleet.New(1),
		sty:   sty,
		w:     132, h: 38,
		ratio: 0.34,
		focus: ownerList,
		last:  "move the mouse",
	}
}

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.w, m.h = msg.Width, msg.Height

	case tea.KeyMsg:
		switch msg.String() {
		case "q", "ctrl+c":
			return m, tea.Quit
		case "esc":
			m.menu = nil
		}

	case tea.MouseMsg:
		m.mouse(msg)
	}
	return m, nil
}

// mouse is the whole interaction layer. Every branch starts by asking the
// canvas what is under the pointer — there is no other source of truth.
func (m *model) mouse(e tea.MouseMsg) {
	if m.canvas == nil {
		return
	}
	owner := m.canvas.OwnerAt(e.X, e.Y)
	m.under = owner
	m.last = fmt.Sprintf("%s at %d,%d", e.String(), e.X, e.Y)

	// A drag in progress owns the mouse until it is released, whatever it is
	// now over. Without this the splitter is dropped the moment the pointer
	// outruns it, which it always does.
	if m.dragging {
		if e.Action == tea.MouseActionRelease {
			m.dragging = false
			return
		}
		m.setRatio(float64(e.X) / float64(m.w))
		return
	}

	// A menu is modal, and the modal check comes first — the same rule the
	// keyboard has.
	if m.menu != nil {
		if e.Action == tea.MouseActionPress {
			if strings.HasPrefix(owner, "menu.item[") {
				var i int
				fmt.Sscanf(owner, "menu.item[%d]", &i)
				m.last = "chose: " + m.menu[i]
			}
			m.menu = nil
		}
		return
	}

	if tea.MouseEvent(e).IsWheel() {
		// Scroll what the pointer is OVER, not what has focus. Getting this
		// backwards is the commonest wheel bug in a TUI.
		delta := 1
		if e.Button == tea.MouseButtonWheelUp {
			delta = -1
		}
		switch {
		case strings.HasPrefix(owner, "list.row[") || owner == ownerList:
			m.sel = clampInt(m.sel+delta, 0, len(m.fleet.Services)-1)
		case owner == ownerDetl || strings.HasPrefix(owner, "detail."):
			m.scroll = clampInt(m.scroll+delta, 0, 20)
		}
		return
	}

	if e.Action != tea.MouseActionPress {
		return
	}

	switch e.Button {
	case tea.MouseButtonLeft:
		switch {
		case owner == ownerSplit:
			m.dragging = true
		case strings.HasPrefix(owner, "list.row["):
			var i int
			fmt.Sscanf(owner, "list.row[%d]", &i)
			m.sel, m.focus = i, ownerList
		case strings.HasPrefix(owner, "detail.tab["):
			var i int
			fmt.Sscanf(owner, "detail.tab[%d]", &i)
			m.tab, m.focus = i, ownerDetl
		case owner == ownerList:
			m.focus = ownerList
		case owner == ownerDetl:
			m.focus = ownerDetl
		}
	case tea.MouseButtonRight:
		if strings.HasPrefix(owner, "list.row[") {
			var i int
			fmt.Sscanf(owner, "list.row[%d]", &i)
			m.sel = i
			// In the real thing these come from spec.Command, so the menu
			// cannot drift from the CLI or the keymap.
			m.menu = []string{"View logs", "Deploy", "Restart", "Remove"}
			m.menuAt = [2]int{e.X, e.Y}
		}
	}
}

func (m *model) setRatio(f float64) {
	m.ratio = clampFloat(f, 0.18, 0.7)
}

func (m *model) View() string {
	c := NewCanvas(m.w, m.h)
	listW := int(float64(m.w) * m.ratio)

	m.drawHeader(c)
	m.drawList(c, Rect{0, 2, listW, m.h - 4})
	m.drawSplitter(c, listW, Rect{listW, 2, 1, m.h - 4})
	m.drawDetail(c, Rect{listW + 1, 2, m.w - listW - 1, m.h - 4})
	m.drawStatus(c)
	if m.menu != nil {
		m.drawMenu(c)
	}

	m.canvas = c
	return c.String()
}

func (m *model) drawHeader(c *Canvas) {
	c.Text(0, 0, "democtl", m.sty["title"], "header")
	c.Text(8, 0, "canvas + mouse prototype", m.sty["muted"], "header")
	c.Fill(Rect{0, 1, m.w, 1}, '─', m.sty["border"], "header")
}

func (m *model) drawList(c *Canvas, r Rect) {
	style := m.sty["border"]
	if m.focus == ownerList {
		style = m.sty["focus"]
	}
	c.Box(r, fmt.Sprintf("Services (%d)", len(m.fleet.Services)), style, m.sty["title"], ownerList)

	in := r.Inset(1)
	for i, s := range m.fleet.Services {
		if i >= in.H {
			break
		}
		row := fmt.Sprintf(" %-14s %d/%d %s", s.Name, s.Ready, s.Want, mark(s.State))
		st := m.stateStyle(s.State)
		if i == m.sel {
			st = m.sty["sel"]
		}
		// Fill first so the whole row belongs to the row, not just its text:
		// clicking the blank space after a short name has to select it.
		c.Fill(Rect{in.X, in.Y + i, in.W, 1}, ' ', st, rowOwner(i))
		c.Text(in.X, in.Y+i, truncate(row, in.W), st, rowOwner(i))
	}
}

// drawSplitter claims one column. Making the divider a real owner is what turns
// dragging from a coordinate guess into a hit test.
func (m *model) drawSplitter(c *Canvas, _ int, r Rect) {
	style := m.sty["border"]
	if m.dragging {
		style = m.sty["focus"]
	}
	c.Fill(r, '│', style, ownerSplit)
}

func (m *model) drawDetail(c *Canvas, r Rect) {
	style := m.sty["border"]
	if m.focus == ownerDetl {
		style = m.sty["focus"]
	}
	svc := m.fleet.Services[m.sel]
	c.Box(r, svc.Name, style, m.sty["title"], ownerDetl)

	in := r.Inset(1)
	x := in.X + 1
	for i, name := range []string{"Overview", "Config", "Events"} {
		label := " " + name + " "
		st := m.sty["muted"]
		if i == m.tab {
			st = m.sty["title"]
			label = "‹" + name + "›"
		}
		c.Text(x, in.Y, label, st, tabOwner(i))
		x += len([]rune(label)) + 1
	}

	lines := []string{
		fmt.Sprintf("state      %s", svc.State),
		fmt.Sprintf("replicas   %d/%d", svc.Ready, svc.Want),
		fmt.Sprintf("stack      %s", svc.Stack),
		fmt.Sprintf("node       %s", svc.Node),
		fmt.Sprintf("image      %s", svc.Image),
	}
	if m.tab == 2 {
		lines = nil
		for _, l := range m.fleet.Logs(svc.Name, 30) {
			lines = append(lines, l.At.Format("15:04:05")+"  "+l.Text)
		}
	}
	for i := m.scroll; i < len(lines) && i-m.scroll < in.H-2; i++ {
		c.Text(in.X+1, in.Y+2+i-m.scroll, truncate(lines[i], in.W-2), m.sty["muted"], ownerDetl)
	}
}

// drawMenu is the reason a canvas beats splicing strings: a menu is drawn last,
// so it is on top. There is no compositing step, no re-measuring of the lines
// underneath, and nothing to get wrong.
func (m *model) drawMenu(c *Canvas) {
	w := 16
	r := Rect{m.menuAt[0], m.menuAt[1], w, len(m.menu) + 2}
	c.Fill(r, ' ', m.sty["muted"], ownerMenu)
	c.Box(r, "", m.sty["focus"], m.sty["title"], ownerMenu)
	for i, item := range m.menu {
		c.Fill(Rect{r.X + 1, r.Y + 1 + i, w - 2, 1}, ' ', m.sty["muted"], menuOwner(i))
		c.Text(r.X+1, r.Y+1+i, " "+item, m.sty["muted"], menuOwner(i))
	}
}

func (m *model) drawStatus(c *Canvas) {
	y := m.h - 1
	c.Text(0, y-1, truncate("  "+m.last, m.w), m.sty["muted"], "status")
	state := fmt.Sprintf("  under: %-18s ratio: %.2f  sel: %d  tab: %d  scroll: %d  focus: %s",
		orNone(m.under), m.ratio, m.sel, m.tab, m.scroll, m.focus)
	c.Text(0, y, truncate(state, m.w), m.sty["title"], "status")
}

func (m *model) stateStyle(s fleet.State) *lipgloss.Style {
	switch s {
	case fleet.Running:
		return m.sty["ok"]
	case fleet.Failed:
		return m.sty["danger"]
	default:
		return m.sty["pending"]
	}
}

func mark(s fleet.State) string {
	switch s {
	case fleet.Running:
		return "✓"
	case fleet.Failed:
		return "✗"
	case fleet.Pending:
		return "•"
	default:
		return "●"
	}
}

func truncate(s string, w int) string {
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	if w < 1 {
		return ""
	}
	return string(r[:w-1]) + "…"
}

func orNone(s string) string {
	if s == "" {
		return "(nothing)"
	}
	return s
}

func clampInt(v, lo, hi int) int {
	return max(lo, min(v, hi))
}

func clampFloat(v, lo, hi float64) float64 {
	return max(lo, min(v, hi))
}

func main() {
	p := tea.NewProgram(newModel(), tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "prototype:", err)
		os.Exit(1)
	}
}
