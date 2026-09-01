package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/app"
	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
)

// stepDoneMsg reports a finished step. It carries the generation it belongs to,
// which is the only thing standing between an abandoned run and a step drawing
// into its successor.
type stepDoneMsg struct {
	gen   int
	index int
	state stepState
}

// Update handles one message.
func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
		return m, nil

	case stepDoneMsg:
		if m.gen.Stale(msg.gen) {
			// A result from a run the user walked away from. Dropping it here
			// is the whole reason gen exists.
			return m, nil
		}
		m.done[msg.index] = msg.state
		if msg.state == stepFailed {
			m.running = -1
			return m, nil
		}
		return m, m.startStep(msg.index + 1)

	case tea.MouseMsg:
		return m, m.onMouse(msg)

	case tea.KeyMsg:
		return m, m.key(msg)
	}
	return m, nil
}

// key routes one keystroke.
//
// The order used to be a comment and two early returns anyone could delete. It
// is app.Keys' struct now: whatever is capturing goes in Capture, and a nil
// Capture is a tool with nothing capturing rather than a tool that forgot.
func (m *Model) key(msg tea.KeyMsg) tea.Cmd {
	return app.Keys{
		Capture: m.capture(),
		Screen:  m.screenKey,
		Global:  func(msg tea.KeyMsg) (tea.Cmd, bool) { return m.globalKey(msg), true },
	}.Route(msg)
}

// capture is whatever is eating keystrokes, or nil.
func (m *Model) capture() app.Handled {
	switch {
	case m.confirm != nil:
		return func(msg tea.KeyMsg) (tea.Cmd, bool) { return m.confirmKey(msg), true }
	case m.menu != nil:
		return func(msg tea.KeyMsg) (tea.Cmd, bool) { return m.menuKey(msg), true }
	case m.typing:
		return func(msg tea.KeyMsg) (tea.Cmd, bool) { return m.filterKey(msg), true }
	}
	return nil
}

func (m *Model) screenKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	switch m.screen {
	case screenDashboard:
		return nil, m.dashboardKey(msg)
	case screenLogs:
		return nil, m.logsKey(msg)
	case screenRun:
		return m.runKey(msg)
	}
	return nil, false
}

func (m *Model) globalKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q", "ctrl+c":
		return tea.Quit
	case "esc":
		if m.screen != screenDashboard {
			m.screen = screenDashboard
			// Abandoning a run invalidates everything still in flight for it.
			m.gen.Next()
			m.running = -1
		}
	}
	return nil
}

// dashboardKey returns whether it handled the key. It never returns a command:
// the one action here opens a modal, and the command comes from confirming it.
func (m *Model) dashboardKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "j", "down":
		m.list.Move(1)
	case "k", "up":
		m.list.Move(-1)
	case "tab":
		if m.focus == paneList {
			m.focus = paneDetail
		} else {
			m.focus = paneList
		}
	case "left", "h":
		if m.focus == paneDetail && m.tab > 0 {
			m.tab--
		}
	case "right", "l":
		if m.focus == paneDetail && m.tab < len(tabNames)-1 {
			m.tab++
		}
	case "/":
		m.typing = true
	case "L":
		if _, ok := m.selected(); ok {
			m.screen, m.log.Follow = screenLogs, true
		}
	case "D":
		m.confirmDeploy()
	case "m":
		// The menu opens from the keyboard, at the cursor. Not a convenience:
		// it is the only way the mouse and keyboard paths cannot drift, since
		// they are one list rather than a list and a keymap maintained beside
		// it.
		m.openMenuOn(comp.Region(regServicesRow).At(m.list.Cursor()))
	default:
		return false
	}
	return true
}

func (m *Model) filterKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "enter", "esc":
		m.typing = false
		if msg.String() == "esc" {
			m.filter = ""
		}
		m.list.Reset()
	case "backspace":
		if m.filter != "" {
			m.filter = m.filter[:len(m.filter)-1]
			m.list.Reset()
		}
	default:
		if len(msg.Runes) == 1 {
			m.filter += string(msg.Runes)
			// The rows changed out from under the cursor, so it goes back to
			// the top rather than pointing at whatever is now in its place.
			m.list.Reset()
		}
	}
	return nil
}

func (m *Model) confirmKey(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "y", "enter":
		do := m.confirm.do
		m.confirm = nil
		if do != nil {
			return do(m)
		}
	case "n", "esc":
		m.confirm = nil
	}
	return nil
}

func (m *Model) logsKey(msg tea.KeyMsg) bool {
	switch msg.String() {
	case "j", "down":
		m.log.Scroll(1)
	case "k", "up":
		m.log.Scroll(-1)
	case "e":
		m.logStderr = !m.logStderr
		// A different stream is a different buffer, so it opens at ITS tail
		// rather than at wherever the last one had been scrolled to.
		m.log.Follow = true
	default:
		return false
	}
	return true
}

func (m *Model) runKey(msg tea.KeyMsg) (tea.Cmd, bool) {
	if msg.String() == "r" && m.running < 0 {
		m.startRun(m.plan)
		return m.startStep(0), true
	}
	return nil, false
}

// startRun resets the step list and takes a new generation, so anything still
// in flight from the last one lands in the void.
func (m *Model) startRun(p fleet.Plan) {
	m.screen, m.plan = screenRun, p
	m.done = make([]stepState, len(p.Steps))
	m.gen.Next()
}

// startStep runs one step and reports the next. Sequential on purpose: a step
// list whose steps overlap cannot say which one you are waiting on.
func (m *Model) startStep(i int) tea.Cmd {
	if i >= len(m.plan.Steps) {
		m.running = -1
		return nil
	}
	step := m.plan.Steps[i]
	m.running = i
	m.done[i] = stepRunning

	state := stepOK
	switch {
	case step.Skip:
		state = stepSkipped
	case step.Fails:
		state = stepFailed
	}

	gen := m.gen.Current()
	// Tick rather than a goroutine: the duration is the step's own, so a
	// captured run reports the timings the plan declares.
	return tea.Tick(step.Took, func(time.Time) tea.Msg {
		return stepDoneMsg{gen: gen, index: i, state: state}
	})
}
