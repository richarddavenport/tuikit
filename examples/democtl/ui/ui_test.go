package ui

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
)

// states is every distinct thing democtl can be showing, named. The width test
// walks it, and so will the capture harness — one list, so a screen added
// without a frame is a screen added without either.
func states() []struct {
	name  string
	build func() *Model
} {
	return []struct {
		name  string
		build func() *Model
	}{
		{"dashboard", func() *Model { return New(1) }},
		{"dashboard-detail-focused", func() *Model {
			m := New(1)
			press(m, "tab")
			return m
		}},
		{"dashboard-config-tab", func() *Model {
			m := New(1)
			press(m, "tab", "right")
			return m
		}},
		{"dashboard-events-tab", func() *Model {
			m := New(1)
			press(m, "tab", "right", "right")
			return m
		}},
		{"failed-service-selected", func() *Model {
			m := New(1)
			press(m, "j", "j")
			return m
		}},
		{"filtering", func() *Model {
			m := New(1)
			press(m, "/", "w", "e")
			return m
		}},
		{"filter-matches-nothing", func() *Model {
			m := New(1)
			press(m, "/", "z", "z", "z")
			return m
		}},
		{"logs", func() *Model {
			m := New(1)
			press(m, "L")
			return m
		}},
		{"logs-stderr-only", func() *Model {
			m := New(1)
			press(m, "L", "e")
			return m
		}},
		{"confirm", func() *Model {
			m := New(1)
			press(m, "D")
			return m
		}},
		{"confirm-danger", func() *Model {
			m := New(1)
			press(m, "j", "j", "D")
			return m
		}},
		{"run-in-flight", func() *Model {
			m := New(1)
			press(m, "D", "y")
			settle(m, 2)
			return m
		}},
		{"run-succeeded", func() *Model {
			m := New(1)
			press(m, "D", "y")
			settle(m, 20)
			return m
		}},
		{"run-failed", func() *Model {
			m := New(1)
			press(m, "j", "j", "D", "y")
			settle(m, 20)
			return m
		}},
	}
}

// Nothing may be wider than the terminal.
//
// The commonest bug in a TUI and the one assertions miss, because content that
// is fine in a fixture is not fine on a real machine: a description listing
// every affected row, a path from someone's home directory. Two real tools
// shipped it. This test is the shape guard.Width will take.
func TestNothingIsWiderThanTheTerminal(t *testing.T) {
	for _, width := range []int{80, 100, 132} {
		for _, s := range states() {
			t.Run(s.name, func(t *testing.T) {
				m := s.build()
				m.SetSize(width, 38)
				for i, line := range strings.Split(m.View(), "\n") {
					if w := lipgloss.Width(line); w > width {
						t.Errorf("at %d columns, line %d is %d wide:\n%q", width, i+1, w, line)
					}
				}
			})
		}
	}
}

// A screen without a View case renders as an empty terminal, and nothing on
// screen says a case is missing. The default branch complains in colour; this
// checks no real screen reaches it.
func TestEveryScreenHasAViewCase(t *testing.T) {
	for s := screenDashboard; s <= screenRun; s++ {
		m := New(1)
		m.screen = s
		if got := m.View(); strings.Contains(got, "has no View case") {
			t.Errorf("screen %d has no View case", s)
		}
	}
}

// The mode split. While the filter is being typed, j is the letter j — without
// this the list scrolls behind whatever is capturing keys, which is the kind of
// bug that survives review because nobody tries it.
func TestTypingDoesNotAlsoDriveTheList(t *testing.T) {
	m := New(1)
	press(m, "/", "j", "j")

	if m.cur != 0 {
		t.Errorf("the cursor moved to %d while the filter was being typed", m.cur)
	}
	if m.filter != "jj" {
		t.Errorf("filter = %q, want %q", m.filter, "jj")
	}
}

// The same for a modal: a dialog you can scroll the list behind is a dialog
// that is not modal.
func TestAModalCapturesKeysToo(t *testing.T) {
	m := New(1)
	press(m, "D", "j")

	if m.cur != 0 {
		t.Error("the list moved behind the modal")
	}
	if m.confirm == nil {
		t.Error("j dismissed the modal")
	}
}

// The generation counter. A step that finishes after the run was abandoned must
// not draw into its successor — easy to write, almost impossible to see once
// written, which is why it is a test rather than a comment.
func TestAStepFromAnAbandonedRunIsDropped(t *testing.T) {
	m := New(1)
	press(m, "D", "y")
	stale := m.gen

	press(m, "esc") // walks away, taking a new generation

	m.done = make([]stepState, len(m.plan.Steps))
	m.Update(stepDoneMsg{gen: stale, index: 0, state: stepOK})

	if m.done[0] != stepWaiting {
		t.Error("a result from an abandoned run was applied to the next one")
	}
}

// The skip state is the idempotency case, and a step list that draws it the
// same as a success lies about the work it did.
func TestASkippedStepIsNotDrawnAsSuccess(t *testing.T) {
	m := New(1)
	m.SetSize(100, 38)
	press(m, "D", "y")
	settle(m, 20)

	if !strings.Contains(m.View(), "already true") {
		t.Error("the skipped step is indistinguishable from one that ran")
	}
}

// An empty result is an ordinary state, not an error, and has to say so rather
// than draw an empty box.
func TestAFilterMatchingNothingSaysSo(t *testing.T) {
	m := New(1)
	m.SetSize(100, 38)
	press(m, "/", "z", "z", "z")

	if !strings.Contains(m.View(), "nothing matches") {
		t.Error("a filter matching nothing drew an empty pane with no explanation")
	}
}

// --- driving the model --------------------------------------------------

// press sends keystrokes and runs whatever commands they return to completion,
// so a test reads as the sequence a person would type.
func press(m *Model, keys ...string) {
	for _, k := range keys {
		var msg tea.KeyMsg
		switch k {
		case "tab", "esc", "enter", "up", "down", "left", "right", "backspace":
			msg = tea.KeyMsg{Type: keyType(k)}
		default:
			msg = tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(k)}
		}
		_, cmd := m.Update(msg)
		drain(m, cmd, 1)
	}
}

func keyType(name string) tea.KeyType {
	switch name {
	case "tab":
		return tea.KeyTab
	case "esc":
		return tea.KeyEsc
	case "enter":
		return tea.KeyEnter
	case "up":
		return tea.KeyUp
	case "down":
		return tea.KeyDown
	case "left":
		return tea.KeyLeft
	case "right":
		return tea.KeyRight
	case "backspace":
		return tea.KeyBackspace
	}
	return tea.KeyRunes
}

// settle runs the pending step commands n times, so a run can be looked at
// partway through or at its end. The steps are tea.Tick commands, which block
// for their duration — settle calls them directly rather than waiting.
func settle(m *Model, n int) {
	for range n {
		if m.running < 0 {
			return
		}
		i := m.running
		state := stepOK
		switch {
		case m.plan.Steps[i].Skip:
			state = stepSkipped
		case m.plan.Steps[i].Fails:
			state = stepFailed
		}
		_, cmd := m.Update(stepDoneMsg{gen: m.gen, index: i, state: state})
		drain(m, cmd, 0)
	}
}

// drain runs a command's message back into the model. depth stops a command
// that returns a command that returns a command from running forever; the step
// chain is the only thing here that nests, and one level is enough.
func drain(m *Model, cmd tea.Cmd, depth int) {
	if cmd == nil || depth < 0 {
		return
	}
	// A tea.Tick would block for the step's duration. The run tests drive the
	// chain through settle instead, so nothing here has to wait on a clock.
	if m.screen == screenRun {
		return
	}
	if msg := cmd(); msg != nil {
		if _, next := m.Update(msg); next != nil {
			drain(m, next, depth-1)
		}
	}
}

var _ = fleet.Running

// --- what looking at frames found ---------------------------------------

// A modal floats over the panes; it does not punch a hole through them.
//
// The first version replaced the whole background line, so the frame the dialog
// sat on vanished for its height. No assertion noticed. A captured frame made it
// obvious immediately, which is the argument for capture in one test.
func TestAModalDoesNotWipeTheFrameBehindIt(t *testing.T) {
	closed, open := New(1), New(1)
	closed.SetSize(100, 24)
	open.SetSize(100, 24)
	press(open, "D")

	// Counted rather than checked line by line: the assertion is that the
	// modal destroyed no framed row, and counting says that without having to
	// name which rows the modal happens to land on.
	before, after := framedRows(closed.View()), framedRows(open.View())
	if before != after {
		t.Errorf("the modal wiped %d framed rows out of %d — it is punching a hole "+
			"through the panes rather than sitting on them", before-after, before)
	}
	if !strings.Contains(open.View(), "Deploy api_gateway?") {
		t.Fatal("the modal is not on screen at all")
	}
}

// framedRows counts the lines that still have a pane edge on both sides.
func framedRows(view string) int {
	var n int
	for _, line := range strings.Split(view, "\n") {
		if strings.HasPrefix(line, "│") && strings.HasSuffix(line, "│") {
			n++
		}
	}
	return n
}

// A duration measured against the terminal rather than the box is a duration
// two columns past the edge, clipped to "400m…".
func TestAStepsDurationFitsInsideTheBox(t *testing.T) {
	m := New(1)
	m.SetSize(100, 24)
	press(m, "D", "y")
	settle(m, 20)

	view := m.View()
	if !strings.Contains(view, "400ms") {
		t.Error("the first step's duration is truncated")
	}
	if strings.Contains(view, "m…") {
		t.Errorf("a duration is being clipped:\n%s", view)
	}
}
