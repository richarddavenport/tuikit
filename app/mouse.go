package app

import (
	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/comp"
)

// WheelRows is how far one notch scrolls.
//
// Three, because that is what every other terminal program does — confirmed
// rather than assumed: herdr's own ui.mouse_scroll_lines defaults to 3.
const WheelRows = 3

// Mouse routes mouse events through the canvas, resolving what was hit rather
// than where it happened.
//
// It carries one piece of state, and that piece is the whole reason it is a
// type: a DRAG IN PROGRESS OWNS THE MOUSE UNTIL RELEASE, whatever it is now
// over. The pointer outruns the divider it grabbed on every real drag, and a
// handler that re-resolves the owner on each motion drops the divider the
// moment the cursor leaves it. Written by hand this is four lines at the top of
// the handler that everyone writes once they have felt the bug.
type Mouse struct {
	dragging bool
	// on is what the drag started on, so every motion is delivered to the
	// thing that was grabbed rather than to whatever is under the pointer now.
	on comp.ID
}

// Handler is what a tool does about a mouse event. Any of these may be nil.
type Handler struct {
	// Blocked reports that something has taken the mouse — a modal, usually.
	// Clicking the frame behind a question is not an answer to it.
	Blocked func() bool

	// Press is a left press on a region.
	Press func(id comp.ID, msg tea.MouseMsg) tea.Cmd
	// RightPress is a right press, which may never arrive: a multiplexer that
	// keeps right-click for itself gets the event first and the application
	// never learns it happened. Whatever this reaches needs a keyboard path.
	RightPress func(id comp.ID, msg tea.MouseMsg) tea.Cmd
	// Wheel is a notch over a region, in rows — negative is up. The region is
	// the one under the POINTER, not the focused one: looking somewhere is not
	// the same as working there.
	Wheel func(id comp.ID, by int) tea.Cmd

	// Drags reports whether a press on this region begins a drag. Nil means
	// nothing drags.
	Drags func(id comp.ID) bool
	// Drag is every motion while dragging, with the region the drag STARTED
	// on. Called for the press too, so a drag that never moves still acts.
	Drag func(id comp.ID, msg tea.MouseMsg) tea.Cmd
	// Release ends it.
	Release func(id comp.ID) tea.Cmd
}

// Route dispatches one event against the last frame drawn.
//
// The canvas is what makes this possible: it records who drew each cell, so the
// question is "what is this" rather than "where is this", and there is no
// region list to keep in step with the drawing.
func (m *Mouse) Route(msg tea.MouseMsg, c *comp.Canvas, h Handler) tea.Cmd {
	if c == nil {
		return nil
	}

	// A drag owns the mouse until release. Checked before anything else,
	// including Blocked: a modal that opens mid-drag must not leave the
	// divider stuck to the pointer.
	if m.dragging {
		if msg.Action == tea.MouseActionRelease {
			m.dragging = false
			if h.Release != nil {
				return h.Release(m.on)
			}
			return nil
		}
		if h.Drag != nil {
			return h.Drag(m.on, msg)
		}
		return nil
	}

	if h.Blocked != nil && h.Blocked() {
		return nil
	}
	id := c.OwnerAt(msg.X, msg.Y)

	switch {
	case msg.Button == tea.MouseButtonWheelUp:
		if h.Wheel != nil {
			return h.Wheel(id, -WheelRows)
		}
	case msg.Button == tea.MouseButtonWheelDown:
		if h.Wheel != nil {
			return h.Wheel(id, WheelRows)
		}
	case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft:
		if h.Drags != nil && h.Drags(id) {
			m.dragging, m.on = true, id
			if h.Drag != nil {
				return h.Drag(id, msg)
			}
			return nil
		}
		if h.Press != nil {
			return h.Press(id, msg)
		}
	case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonRight:
		if h.RightPress != nil {
			return h.RightPress(id, msg)
		}
	}
	return nil
}

// Dragging reports whether a drag is in progress, for a view that draws it.
func (m *Mouse) Dragging() bool { return m.dragging }
