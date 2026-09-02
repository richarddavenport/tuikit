package harness

import (
	"fmt"
	"strconv"
	"strings"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/comp"
)

// Pointer is a Driver that can say what it drew.
//
// That second half is what makes naming possible. A model that only returns a
// string can be sent a click at column 31 and nobody — not a person reading the
// script, not an agent writing one — can say what that was meant to hit.
type Pointer interface {
	Driver
	Canvas() *comp.Canvas
}

// Script drives a model with keys and mouse actions, addressed by region NAME.
//
//	harness.Script(t, m, `
//	    press E j j
//	    click services.row[2]
//	    wheel logs -3
//	    drag split +10
//	    rclick services.row[2]
//	    doubleclick services.row[2]
//	    hover services.row[2]
//	`)
//
// A coordinate is a guess that happens to work today. A name is a claim about
// the interface, and one that has to be drawn to be clickable — so a script
// that has drifted from the tool says which region it can no longer find,
// rather than clicking empty space and reporting success. That is the whole
// reason the canvas records ownership at draw time.
//
// An agent writes the second kind and cannot write the first.
func Script(t T, m Pointer, script string) {
	t.Helper()
	for i, line := range strings.Split(script, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		if err := step(m, line); err != nil {
			t.Fatalf("script line %d: %s: %v", i+1, line, err)
		}
	}
}

// Click presses the left button in the middle of a region.
func Click(t T, m Pointer, region string) { t.Helper(); do(t, m, "click "+region) }

// RClick presses the right button, which is how a context menu opens — when
// the terminal lets the event through at all. A multiplexer that keeps
// right-click for itself means the application never sees this, which is why
// every menu also opens from the keyboard.
func RClick(t T, m Pointer, region string) { t.Helper(); do(t, m, "rclick "+region) }

// DoubleClick presses twice on a region.
//
// Whether that counts as a double is app.Mouse's decision and its clock's; this
// only sends what a terminal sends. A helper that asserted the answer would be
// testing itself.
func DoubleClick(t T, m Pointer, region string) { t.Helper(); do(t, m, "doubleclick "+region) }

// Hover moves the pointer over a region without pressing anything.
func Hover(t T, m Pointer, region string) { t.Helper(); do(t, m, "hover "+region) }

// Wheel turns the wheel over a region, in notches. Negative is up.
func Wheel(t T, m Pointer, region string, notches int) {
	t.Helper()
	do(t, m, fmt.Sprintf("wheel %s %+d", region, notches))
}

// Drag presses in a region, moves it by n, and releases.
//
// The AXIS comes from the region's shape rather than from a second argument: a
// divider is one column wide or one row tall, so a tall thin one moves in
// columns and a wide flat one moves in rows. A vertical split's divider takes
// `drag run.split +4` and moves down four rows, which is the only reading of
// that line anyone has — and the alternative, an X-only Drag, silently does
// nothing at all to a stacked divider.
func Drag(t T, m Pointer, region string, n int) {
	t.Helper()
	do(t, m, fmt.Sprintf("drag %s %+d", region, n))
}

func do(t T, m Pointer, line string) {
	t.Helper()
	if err := step(m, line); err != nil {
		t.Fatalf("%s: %v", line, err)
	}
}

func step(m Pointer, line string) error {
	fields := strings.Fields(line)
	verb, rest := fields[0], fields[1:]
	if len(rest) == 0 {
		return fmt.Errorf("%q names nothing to act on", verb)
	}

	// press names keys rather than a region, so it is answered before anything
	// tries to look one up.
	if verb == "press" {
		Press(m, rest...)
		return nil
	}

	x, y, err := centre(m, rest[0])
	if err != nil {
		return err
	}

	switch verb {
	case "click":
		send(m, x, y, tea.MouseActionPress, tea.MouseButtonLeft)
		send(m, x, y, tea.MouseActionRelease, tea.MouseButtonLeft)
	case "doubleclick":
		// Two presses, which is what a terminal sends. Whether they count as a
		// double is app.Mouse's decision and its clock's — a script that faked
		// the answer would be testing itself.
		send(m, x, y, tea.MouseActionPress, tea.MouseButtonLeft)
		send(m, x, y, tea.MouseActionRelease, tea.MouseButtonLeft)
		send(m, x, y, tea.MouseActionPress, tea.MouseButtonLeft)
		send(m, x, y, tea.MouseActionRelease, tea.MouseButtonLeft)
	case "hover":
		send(m, x, y, tea.MouseActionMotion, tea.MouseButtonNone)
	case "rclick":
		send(m, x, y, tea.MouseActionPress, tea.MouseButtonRight)
		send(m, x, y, tea.MouseActionRelease, tea.MouseButtonRight)
	case "wheel":
		n, err := amount(rest)
		if err != nil {
			return err
		}
		button := tea.MouseButtonWheelDown
		if n < 0 {
			button, n = tea.MouseButtonWheelUp, -n
		}
		for i := 0; i < n; i++ {
			send(m, x, y, tea.MouseActionPress, button)
		}
	case "drag":
		n, err := amount(rest)
		if err != nil {
			return err
		}
		// Wider than it is tall means a divider between STACKED panes, which
		// moves in rows. Anything else moves in columns.
		dx, dy := n, 0
		if wide, err := wider(m, rest[0]); err != nil {
			return err
		} else if wide {
			dx, dy = 0, n
		}

		send(m, x, y, tea.MouseActionPress, tea.MouseButtonLeft)
		// Through the intervening cells, not straight to the end: a drag that
		// only ever arrives is a drag whose motion handling is untested.
		for step := 1; step <= abs(n); step++ {
			send(m, x+step*sign(dx), y+step*sign(dy), tea.MouseActionMotion, tea.MouseButtonLeft)
		}
		send(m, x+dx, y+dy, tea.MouseActionRelease, tea.MouseButtonLeft)
	default:
		return fmt.Errorf("no such action %q — press, click, rclick, wheel, drag or shot", verb)
	}
	return nil
}

// wider reports that a region is wider than it is tall, which is how a divider
// says which way it slides.
func wider(m Pointer, name string) (bool, error) {
	id, err := comp.ParseID(name)
	if err != nil {
		return false, err
	}
	r, ok := m.Canvas().Region(id)
	if !ok {
		return false, fmt.Errorf("%s was not drawn in the last frame", id)
	}
	return r.W > r.H, nil
}

// centre is the middle of a region, which is the safest cell to aim at: an
// edge is where a rounding error lands.
func centre(m Pointer, name string) (int, int, error) {
	id, err := comp.ParseID(name)
	if err != nil {
		return 0, 0, err
	}
	c := m.Canvas()
	if c == nil {
		return 0, 0, fmt.Errorf("nothing has been drawn yet — render a frame first")
	}
	r, ok := c.Region(id)
	if !ok {
		return 0, 0, fmt.Errorf("%s was not drawn in the last frame", id)
	}
	return r.X + r.W/2, r.Y + r.H/2, nil
}

func send(m Pointer, x, y int, action tea.MouseAction, button tea.MouseButton) {
	m.Update(tea.MouseMsg{X: x, Y: y, Action: action, Button: button})
	// The next line of the script addresses the frame this one produced,
	// rather than the one before it.
	redraw(m)
}

func amount(rest []string) (int, error) {
	if len(rest) < 2 {
		return 0, fmt.Errorf("no amount given")
	}
	n, err := strconv.Atoi(strings.TrimPrefix(rest[1], "+"))
	if err != nil {
		return 0, fmt.Errorf("%q is not an amount", rest[1])
	}
	return n, nil
}

func abs(n int) int {
	if n < 0 {
		return -n
	}
	return n
}

func sign(n int) int {
	if n < 0 {
		return -1
	}
	return 1
}
