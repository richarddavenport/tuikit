package app

import (
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/comp"
)

func key(s string) tea.KeyMsg { return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(s)} }

func took(name string, log *[]string) Handled {
	return func(tea.KeyMsg) (tea.Cmd, bool) { *log = append(*log, name); return nil, true }
}

func passed(name string, log *[]string) Handled {
	return func(tea.KeyMsg) (tea.Cmd, bool) { *log = append(*log, name); return nil, false }
}

// A global handled before a modal is a modal you cannot type into: press q in a
// filter box and the program exits.
func TestWhateverIsCapturingGetsTheKeyFirst(t *testing.T) {
	var log []string
	Keys{
		Capture: took("capture", &log),
		Screen:  took("screen", &log),
		Global:  took("global", &log),
	}.Route(key("q"))

	if len(log) != 1 || log[0] != "capture" {
		t.Errorf("the key went to %v", log)
	}
}

// A capture takes the key whether or not it does anything with it. An
// unrecognised key inside a filter box is a character, not a chance for the
// screen underneath to act on it.
func TestACaptureSwallowsKeysItDoesNotUse(t *testing.T) {
	var log []string
	Keys{
		Capture: passed("capture", &log),
		Screen:  took("screen", &log),
		Global:  took("global", &log),
	}.Route(key("j"))

	if len(log) != 1 {
		t.Errorf("a key fell past the capture to %v", log)
	}
}

// With nothing capturing, the screen gets it, then the globals.
func TestTheOrderWithNothingCapturing(t *testing.T) {
	var log []string
	Keys{Screen: passed("screen", &log), Global: took("global", &log)}.Route(key("q"))

	if want := []string{"screen", "global"}; strings.Join(log, ",") != strings.Join(want, ",") {
		t.Errorf("the key went to %v, want %v", log, want)
	}
}

func TestAScreenThatTakesTheKeyStopsIt(t *testing.T) {
	var log []string
	Keys{Screen: took("screen", &log), Global: took("global", &log)}.Route(key("j"))

	if len(log) != 1 || log[0] != "screen" {
		t.Errorf("the key went to %v", log)
	}
}

// --- mouse --------------------------------------------------------------

func frame() *comp.Canvas {
	c := comp.NewCanvas(20, 4)
	c.Fill(comp.Rect{X: 0, Y: 0, W: 20, H: 1}, " ", nil, comp.Region("rows").At(0))
	c.Fill(comp.Rect{X: 0, Y: 1, W: 1, H: 3}, " ", nil, comp.Region("split"))
	return c
}

func press(x, y int, b tea.MouseButton) tea.MouseMsg {
	return tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: b}
}

// The pointer outruns the divider it grabbed on every real drag. A handler that
// re-resolves the owner on each motion drops it the moment the cursor leaves.
func TestADragOwnsTheMouseUntilRelease(t *testing.T) {
	var m Mouse
	c := frame()
	var to []string

	h := Handler{
		Drags: func(id comp.ID) bool { return id.Name == "split" },
		Drag:  func(id comp.ID, _ tea.MouseMsg) tea.Cmd { to = append(to, id.String()); return nil },
		Press: func(id comp.ID, _ tea.MouseMsg) tea.Cmd { to = append(to, "press:"+id.String()); return nil },
	}
	m.Route(press(0, 2, tea.MouseButtonLeft), c, h)
	// Far away, over a different region entirely.
	m.Route(tea.MouseMsg{X: 18, Y: 0, Action: tea.MouseActionMotion, Button: tea.MouseButtonLeft}, c, h)

	for _, got := range to {
		if got != "split" {
			t.Errorf("the drag was delivered to %q", got)
		}
	}
	if !m.Dragging() {
		t.Error("the drag ended when the pointer left the divider")
	}
	m.Route(tea.MouseMsg{X: 18, Y: 0, Action: tea.MouseActionRelease, Button: tea.MouseButtonLeft}, c, h)
	if m.Dragging() {
		t.Error("the drag survived its release")
	}
}

// A modal that opens mid-drag must not leave the divider stuck to the pointer,
// so the drag is checked before Blocked.
func TestADragInProgressOutranksAModal(t *testing.T) {
	var m Mouse
	c := frame()
	moved := false
	h := Handler{
		Drags: func(id comp.ID) bool { return id.Name == "split" },
		Drag:  func(comp.ID, tea.MouseMsg) tea.Cmd { moved = true; return nil },
	}
	m.Route(press(0, 2, tea.MouseButtonLeft), c, h)

	h.Blocked = func() bool { return true }
	moved = false
	m.Route(tea.MouseMsg{X: 4, Y: 2, Action: tea.MouseActionMotion, Button: tea.MouseButtonLeft}, c, h)
	if !moved {
		t.Error("a modal opening mid-drag stranded the divider on the pointer")
	}
}

// Clicking the frame behind a question is not an answer to it.
func TestAModalTakesTheMouse(t *testing.T) {
	var m Mouse
	clicked := false
	m.Route(press(4, 0, tea.MouseButtonLeft), frame(), Handler{
		Blocked: func() bool { return true },
		Press:   func(comp.ID, tea.MouseMsg) tea.Cmd { clicked = true; return nil },
	})
	if clicked {
		t.Error("a click reached the frame behind a modal")
	}
}

// The wheel acts on the region under the POINTER, not the focused one.
func TestTheWheelActsOnWhatIsUnderThePointer(t *testing.T) {
	var m Mouse
	var on comp.ID
	var by int
	h := Handler{Wheel: func(id comp.ID, n int) tea.Cmd { on, by = id, n; return nil }}

	m.Route(press(4, 0, tea.MouseButtonWheelDown), frame(), h)
	if on.Name != "rows" || by != WheelRows {
		t.Errorf("wheel down went to %v by %d", on, by)
	}
	m.Route(press(4, 0, tea.MouseButtonWheelUp), frame(), h)
	if by != -WheelRows {
		t.Errorf("wheel up moved by %d", by)
	}
}

func TestNothingHappensWithoutAFrame(t *testing.T) {
	var m Mouse
	if cmd := m.Route(press(0, 0, tea.MouseButtonLeft), nil, Handler{
		Press: func(comp.ID, tea.MouseMsg) tea.Cmd { t.Error("pressed with no frame drawn"); return nil },
	}); cmd != nil {
		t.Error("a command came from nowhere")
	}
}

// --- screens ------------------------------------------------------------

// A missing screen renders an empty terminal and says nothing about why. This
// says which one, and says it ugly: a bug that looks like a design decision
// gets shipped.
func TestAScreenWithNoViewComplainsVisibly(t *testing.T) {
	c := comp.NewCanvas(60, 3)
	Screens{0: func(*comp.Canvas, comp.Rect) {}}.Draw(7, c, c.Bounds(), nil)

	if got := c.String(); !strings.Contains(got, "screen 7 has no View") {
		t.Errorf("a missing screen drew %q", got)
	}
}

func TestAScreenWithAViewDrawsIt(t *testing.T) {
	c := comp.NewCanvas(20, 2)
	Screens{3: func(c *comp.Canvas, r comp.Rect) {
		c.Text(r.X, r.Y, "the screen", nil, comp.Region("s"))
	}}.Draw(3, c, c.Bounds(), nil)

	if got := c.String(); !strings.Contains(got, "the screen") {
		t.Errorf("got %q", got)
	}
}

// A nil view is a screen somebody meant to write, which is the same failure as
// one that is missing.
func TestANilViewIsAMissingView(t *testing.T) {
	s := Screens{1: nil}
	if s.Has(1) {
		t.Error("a nil view counts as present")
	}
	c := comp.NewCanvas(60, 2)
	s.Draw(1, c, c.Bounds(), nil)
	if !strings.Contains(c.String(), "no View") {
		t.Error("a nil view drew nothing and said nothing")
	}
}
