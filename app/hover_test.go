package app_test

import (
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/app"
	"github.com/richarddavenport/tuikit/comp"
)

// frame is two regions side by side, so the pointer has somewhere to move
// between.
func frame() *comp.Canvas {
	c := comp.NewCanvas(20, 2)
	c.Fill(comp.Rect{X: 0, Y: 0, W: 8, H: 2}, " ", nil, comp.Region("left"))
	c.Fill(comp.Rect{X: 10, Y: 0, W: 8, H: 2}, " ", nil, comp.Region("right"))
	return c
}

func move(x, y int) tea.MouseMsg {
	return tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionMotion, Button: tea.MouseButtonNone}
}

func press(x, y int) tea.MouseMsg {
	return tea.MouseMsg{X: x, Y: y, Action: tea.MouseActionPress, Button: tea.MouseButtonLeft}
}

// Hover fires on entering and leaving, not on every event. A callback per
// motion event is a redraw per motion event.
func TestHoverFiresOnlyOnChange(t *testing.T) {
	var m app.Mouse
	c := frame()
	var got []comp.Name
	h := app.Handler{Hover: func(id comp.ID) tea.Cmd {
		got = append(got, id.Name)
		return nil
	}}

	m.Route(move(1, 0), c, h)  // into left
	m.Route(move(2, 0), c, h)  // still left
	m.Route(move(5, 1), c, h)  // still left
	m.Route(move(11, 0), c, h) // into right
	m.Route(move(12, 1), c, h) // still right
	m.Route(move(9, 0), c, h)  // into the gap, which is nothing

	want := []comp.Name{"left", "right", ""}
	if len(got) != len(want) {
		t.Fatalf("Hover fired %d times (%v), want %d", len(got), got, len(want))
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("hover %d was %q, want %q", i, got[i], want[i])
		}
	}
}

// Leaving everything reports the zero ID rather than nothing at all: a tool
// that highlighted a row needs to be told to stop.
func TestLeavingEverythingIsReported(t *testing.T) {
	var m app.Mouse
	c := frame()
	var last comp.ID
	fired := 0
	h := app.Handler{Hover: func(id comp.ID) tea.Cmd { last, fired = id, fired+1; return nil }}

	m.Route(move(1, 0), c, h)
	m.Route(move(9, 0), c, h)
	if fired != 2 {
		t.Fatalf("Hover fired %d times, want 2", fired)
	}
	if !last.Zero() {
		t.Errorf("leaving reported %v, want the zero ID", last)
	}
}

// A tool with no Hover handler pays nothing, and the state still tracks so
// adding one later needs no other change.
func TestHoverIsOptional(t *testing.T) {
	var m app.Mouse
	if cmd := m.Route(move(1, 0), frame(), app.Handler{}); cmd != nil {
		t.Error("a handler with no Hover returned a command")
	}
}

func TestDoubleClick(t *testing.T) {
	at := time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)
	m := app.Mouse{Now: func() time.Time { return at }}
	c := frame()

	presses, doubles := 0, 0
	h := app.Handler{
		Press:       func(comp.ID, tea.MouseMsg) tea.Cmd { presses++; return nil },
		DoubleClick: func(comp.ID, tea.MouseMsg) tea.Cmd { doubles++; return nil },
	}

	m.Route(press(1, 0), c, h)
	at = at.Add(100 * time.Millisecond)
	m.Route(press(2, 0), c, h)

	// Press fires for BOTH: the first click of a double-click is a real click,
	// and a list that only selected on singles would flicker its selection off.
	if presses != 2 {
		t.Errorf("Press fired %d times, want 2", presses)
	}
	if doubles != 1 {
		t.Errorf("DoubleClick fired %d times, want 1", doubles)
	}
}

// Too slow is two clicks.
func TestTwoSlowClicksAreNotADouble(t *testing.T) {
	at := time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)
	m := app.Mouse{Now: func() time.Time { return at }}
	doubles := 0
	h := app.Handler{DoubleClick: func(comp.ID, tea.MouseMsg) tea.Cmd { doubles++; return nil }}

	m.Route(press(1, 0), frame(), h)
	at = at.Add(2 * time.Second)
	m.Route(press(1, 0), frame(), h)

	if doubles != 0 {
		t.Errorf("two clicks two seconds apart made %d doubles", doubles)
	}
}

// A different region is two clicks, however fast. The same REGION rather than
// the same pixel, though: a row is one thing however wide it is.
func TestADoubleClickIsPerRegion(t *testing.T) {
	at := time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)
	m := app.Mouse{Now: func() time.Time { return at }}
	doubles := 0
	h := app.Handler{DoubleClick: func(comp.ID, tea.MouseMsg) tea.Cmd { doubles++; return nil }}
	c := frame()

	m.Route(press(1, 0), c, h)
	m.Route(press(11, 0), c, h) // the other region, immediately
	if doubles != 0 {
		t.Errorf("clicking two different regions made %d doubles", doubles)
	}

	// The far end of the SAME region does count.
	m.Route(press(7, 1), c, h)
	m.Route(press(11, 0), c, h)
	if doubles != 0 {
		t.Errorf("got %d doubles, want 0 so far", doubles)
	}
	m.Route(press(1, 0), c, h)
	m.Route(press(7, 1), c, h)
	if doubles != 1 {
		t.Errorf("two clicks at either end of one region made %d doubles, want 1", doubles)
	}
}

// Three clicks are a double and a single, not two doubles — which would fire
// "open" twice for one gesture.
func TestThreeClicksAreOneDouble(t *testing.T) {
	at := time.Date(2026, 9, 2, 12, 0, 0, 0, time.UTC)
	m := app.Mouse{Now: func() time.Time { return at }}
	doubles := 0
	h := app.Handler{DoubleClick: func(comp.ID, tea.MouseMsg) tea.Cmd { doubles++; return nil }}
	c := frame()

	for range 3 {
		m.Route(press(1, 0), c, h)
		at = at.Add(50 * time.Millisecond)
	}
	if doubles != 1 {
		t.Errorf("three clicks made %d doubles, want 1", doubles)
	}
}
