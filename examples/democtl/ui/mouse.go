package ui

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
)

// wheelRows is how far one notch scrolls. Three, because that is what every
// other terminal program does — confirmed rather than assumed: herdr's own
// ui.mouse_scroll_lines defaults to 3.
const wheelRows = 3

// mouse routes one mouse event.
//
// It asks the LAST FRAME what is at the pointer, which is the whole benefit of
// drawing into a canvas: there is no region list to keep in step, because the
// frame is the region list. A pane that moved cannot be clicked in its old
// place, since its old place is not what was drawn.
func (m *Model) mouse(msg tea.MouseMsg) tea.Cmd {
	if m.canvas == nil {
		return nil
	}

	// A drag owns the mouse until release, whatever it is now over. The
	// pointer outruns the divider it grabbed on every real drag, and without
	// this the divider is dropped the moment the cursor leaves it.
	if m.dragging {
		if msg.Action == tea.MouseActionRelease {
			m.dragging = false
			return nil
		}
		m.setSplit(msg.X)
		return nil
	}

	id := m.canvas.OwnerAt(msg.X, msg.Y)

	// A modal takes the mouse the way it takes the keyboard. Clicking the
	// frame behind a question is not an answer to it.
	if m.confirm != nil {
		return nil
	}
	if m.menu != nil {
		return m.menuMouse(msg, id)
	}

	switch {
	case msg.Button == tea.MouseButtonWheelUp:
		m.scroll(id, -wheelRows)
	case msg.Button == tea.MouseButtonWheelDown:
		m.scroll(id, wheelRows)

	case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft:
		return m.press(id, msg)

	case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonRight:
		m.openMenuAt(id, msg.X, msg.Y)
	}
	return nil
}

func (m *Model) press(id comp.ID, msg tea.MouseMsg) tea.Cmd {
	switch id.Name {
	case regServicesRow:
		m.focus = paneList
		m.list.Select(id.Index)
	case regServices:
		m.focus = paneList
	case regDetailTab:
		m.focus, m.tab = paneDetail, id.Index
	case regDetail:
		m.focus = paneDetail
	case regSplit:
		m.dragging = true
		m.setSplit(msg.X)
	}
	return nil
}

// scroll moves the viewport of whatever the pointer is over — the pane under
// the pointer, not the focused one, because looking somewhere is not the same
// as working there.
//
// It moves a VIEWPORT and never the selection. Wiring the wheel to the cursor
// instead, which is the obvious shortcut when a list has no viewport, makes
// scrolling appear to pick items at random.
func (m *Model) scroll(id comp.ID, by int) {
	switch id.Name {
	case regLogs, regLogsRow:
		m.log.Scroll(by)
	case regServices, regServicesRow:
		m.list.Scroll(by)
	}
}

// setSplit moves the divider, keeping both panes usable. A split that can be
// dragged to nothing is a pane you cannot get back.
func (m *Model) setSplit(x int) {
	m.split = clamp(x, minPane, m.width-minPane-1)
}

// minPane is the narrowest either pane may be dragged to.
const minPane = 24

func clamp(v, lo, hi int) int { return max(lo, min(v, hi)) }

// --- the context menu ---------------------------------------------------

// menuState is an open context menu.
//
// A menu opened by right-click knows its point. One opened from the keyboard
// knows only which region it is ON, and finds out where that is when the frame
// is drawn — the row it belongs to may have scrolled, or the pane may have been
// resized, since the key was pressed. Asking the frame being drawn is the only
// answer that cannot be stale.
type menuState struct {
	on       comp.ID
	x, y     int
	onRegion bool
	items    []menuItem
	cur      int
}

// menuItem is one action: the hint that names it, and what it does.
//
// comp.Hint rather than a label and a key of its own, because the footer is
// built from the same type. The two paths to an action are one list, which is
// the only arrangement in which they cannot drift.
type menuItem struct {
	comp.Hint
	do func(*Model) tea.Cmd
}

// openMenuAt opens the menu for a region at a point — where a right-click
// landed.
func (m *Model) openMenuAt(id comp.ID, x, y int) {
	items := m.actionsFor(id)
	if len(items) == 0 {
		return
	}
	if id.Name == regServicesRow {
		m.focus = paneList
		m.list.Select(id.Index)
	}
	m.menu = &menuState{on: id, x: x, y: y, items: items}
}

// openMenuOn opens the menu ON a region, wherever that region currently is.
// This is the keyboard path, and "at the cursor" means at the thing the cursor
// is on rather than at the corner of the screen.
//
// Not a convenience. A multiplexer that captures right-click gets the event
// first and this application never learns it happened — there is no protocol
// for asking — so an action whose only path is a context menu is broken for
// everyone inside herdr or tmux, and the tool cannot detect that to warn
// anyone.
func (m *Model) openMenuOn(id comp.ID) {
	items := m.actionsFor(id)
	if len(items) == 0 {
		return
	}
	m.menu = &menuState{on: id, onRegion: true, items: items}
}

// actionsFor is what a region can do. One list, so the menu, the keymap and —
// once spec exists — the CLI command and the manifest all come from it.
func (m *Model) actionsFor(id comp.ID) []menuItem {
	switch id.Name {
	case regServicesRow, regServices:
		if _, ok := m.selected(); !ok {
			return nil
		}
		return []menuItem{
			{Hint: comp.Hint{Key: "L", Label: "View logs"}, do: func(m *Model) tea.Cmd {
				m.screen, m.log.Follow = screenLogs, true
				return nil
			}},
			{Hint: comp.Hint{Key: "D", Label: "Deploy"}, do: func(m *Model) tea.Cmd {
				m.confirmDeploy()
				return nil
			}},
		}
	}
	return nil
}

func (m *Model) menuMouse(msg tea.MouseMsg, id comp.ID) tea.Cmd {
	if msg.Action != tea.MouseActionPress {
		return nil
	}
	if id.Name == regMenuItem {
		m.menu.cur = id.Index
		return m.menuChoose()
	}
	m.menu = nil // a click anywhere else dismisses it
	return nil
}

func (m *Model) menuKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "j", "down":
		m.menu.cur = min(m.menu.cur+1, len(m.menu.items)-1)
	case "k", "up":
		m.menu.cur = max(m.menu.cur-1, 0)
	case "enter":
		return m.menuChoose()
	case "esc", "q":
		m.menu = nil
	}
	return nil
}

func (m *Model) menuChoose() tea.Cmd {
	item := m.menu.items[m.menu.cur]
	m.menu = nil
	return item.do(m)
}

// confirmDeploy is the modal both D and the menu open, so the two paths cannot
// ask different questions.
func (m *Model) confirmDeploy() {
	svc, ok := m.selected()
	if !ok {
		return
	}
	m.confirm = &confirmState{
		title:  "Deploy " + svc.Name + "?",
		body:   "Pushes a new service spec and waits for the tasks to converge. The current tasks are replaced one at a time.",
		danger: svc.State == fleet.Failed,
		do: func(m *Model) tea.Cmd {
			// The failed service is the one whose deploy breaks, so the step
			// list's failure state is reachable by hand rather than only from
			// a test.
			m.startRun(fleet.Deploy(svc.Name, svc.State == fleet.Failed))
			return m.startStep(0)
		},
	}
}
