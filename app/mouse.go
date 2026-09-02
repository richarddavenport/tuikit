package app

import (
	"time"

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
	// hovering is what the pointer was last over, so Hover fires on entering
	// and leaving rather than on every event.
	hovering comp.ID
	// pressed and pressedAt are the last left press, for double-click.
	pressed   comp.ID
	pressedAt time.Time

	// Now is the clock double-click measures against. Nil takes time.Now.
	//
	// A field rather than a call, so a capture script can double-click without
	// real time passing — the same reason comp.Spinner takes a moment rather
	// than counting frames.
	Now func() time.Time
	// DoubleClickWithin is how close two presses must be. Zero takes 400ms,
	// which is the usual desktop default and is not worth inventing.
	DoubleClickWithin time.Duration
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

	// Hover is the region under the pointer, or the zero ID when it has left
	// everything.
	//
	// Called only when the answer CHANGES, so a tool redraws on entering and
	// leaving rather than on every pixel of movement. That is not a nicety: a
	// callback per motion event is a redraw per motion event, and the cost of
	// that is the thing to measure before leaning on it.
	//
	// Hover is the affordance that says a thing is clickable BEFORE you click
	// it, and in a terminal it does more work than in a window, because there
	// is no cursor shape to fall back on.
	Hover func(id comp.ID) tea.Cmd

	// DoubleClick is a second press on the same region, inside the interval.
	//
	// Press still fires for both, because the first click of a double-click is
	// a real click — a list that only selected on single clicks would flicker
	// its selection off on the second.
	//
	// Decision 19 applies with force here: double-click is the conventional
	// "open" gesture, and an open that exists only for a mouse is an open an
	// agent cannot perform. guard.Reachable will say so.
	DoubleClick func(id comp.ID, msg tea.MouseMsg) tea.Cmd
}

// defaultDoubleClick is the usual desktop interval.
const defaultDoubleClick = 400 * time.Millisecond

// now is the clock, defaulting to the real one.
func (m *Mouse) now() time.Time {
	if m.Now != nil {
		return m.Now()
	}
	return time.Now()
}

// double reports whether this press completes a double-click on id, and
// records it either way.
func (m *Mouse) double(id comp.ID) bool {
	within := m.DoubleClickWithin
	if within == 0 {
		within = defaultDoubleClick
	}
	at := m.now()
	// The same REGION, not the same pixel: a row is one thing however wide it
	// is, and asking a reader to hit the same cell twice is asking for a skill
	// rather than a gesture.
	is := id == m.pressed && !m.pressedAt.IsZero() && at.Sub(m.pressedAt) <= within
	if is {
		// Consumed, so three clicks are a double and a single rather than two
		// doubles — which would fire "open" twice for one gesture.
		m.pressed, m.pressedAt = comp.ID{}, time.Time{}
		return true
	}
	m.pressed, m.pressedAt = id, at
	return false
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

	// Motion with no button is the pointer moving over things. Reported only
	// when what is under it changes.
	if msg.Action == tea.MouseActionMotion && msg.Button == tea.MouseButtonNone {
		if id == m.hovering {
			return nil
		}
		m.hovering = id
		if h.Hover != nil {
			return h.Hover(id)
		}
		return nil
	}

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
		// The first click of a double-click is a real click, so Press fires for
		// both and DoubleClick is the extra meaning on the second.
		second := m.double(id)
		var cmds []tea.Cmd
		if h.Press != nil {
			cmds = append(cmds, h.Press(id, msg))
		}
		if second && h.DoubleClick != nil {
			cmds = append(cmds, h.DoubleClick(id, msg))
		}
		return tea.Batch(cmds...)
	case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonRight:
		if h.RightPress != nil {
			return h.RightPress(id, msg)
		}
	}
	return nil
}

// Dragging reports whether a drag is in progress, for a view that draws it.
func (m *Mouse) Dragging() bool { return m.dragging }
