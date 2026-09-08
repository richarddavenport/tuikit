package comp

import (
	"strings"
	"testing"
	"time"
)

var when = time.Date(2026, 9, 2, 9, 30, 0, 0, time.UTC)

func drawWaiting(w Waiting, cw, ch int, at time.Time) string {
	c := NewCanvas(cw, ch)
	w.Draw(c, Rect{X: 0, Y: 0, W: cw, H: ch}, at, Region("wait"))
	return c.String()
}

func TestWaitingIsCentered(t *testing.T) {
	got := drawWaiting(Waiting{Label: "reading the estate"}, 40, 5, when)
	lines := strings.Split(got, "\n")
	if len(lines) != 5 {
		t.Fatalf("got %d lines, want 5", len(lines))
	}
	// The middle row, and roughly the middle column.
	if strings.TrimSpace(lines[2]) == "" {
		t.Errorf("nothing on the middle row:\n%s", got)
	}
	lead := len(lines[2]) - len(strings.TrimLeft(lines[2], " "))
	if lead < 8 {
		t.Errorf("the message starts at column %d; it is not centered:\n%s", lead, got)
	}
}

// The spinner turns with the clock, so a captured frame is deterministic and
// two waits on one screen agree.
func TestWaitingTurnsWithTheClock(t *testing.T) {
	a := drawWaiting(Waiting{Label: "x"}, 20, 3, when)
	b := drawWaiting(Waiting{Label: "x"}, 20, 3, when.Add(300*time.Millisecond))
	if a == b {
		t.Error("the frame did not change 300ms later")
	}
	if c := drawWaiting(Waiting{Label: "x"}, 20, 3, when); c != a {
		t.Error("the same moment drew differently — a golden could not hold this")
	}
}

// A counter that appears for 300ms and vanishes is a flicker. Past the
// threshold it is the difference between "working" and "hung".
func TestElapsedOnlyAfterTheThreshold(t *testing.T) {
	w := Waiting{Label: "reading", Since: when}

	if got := drawWaiting(w, 40, 3, when.Add(300*time.Millisecond)); strings.Contains(got, "0s") {
		t.Errorf("a 300ms wait showed a counter:\n%s", got)
	}
	if got := drawWaiting(w, 40, 3, when.Add(5*time.Second)); !strings.Contains(got, "5s") {
		t.Errorf("a five-second wait did not say so:\n%s", got)
	}
	if got := drawWaiting(w, 40, 3, when.Add(75*time.Second)); !strings.Contains(got, "1m 15s") {
		t.Errorf("a long wait did not read as minutes:\n%s", got)
	}
}

func TestNoElapsedWithoutASince(t *testing.T) {
	got := drawWaiting(Waiting{Label: "reading"}, 40, 3, when.Add(time.Hour))
	if strings.Contains(got, "s") && strings.Contains(got, "60") {
		t.Errorf("counted without being told when it started:\n%s", got)
	}
}

func TestWaitingDetailIsASecondLine(t *testing.T) {
	got := drawWaiting(Waiting{
		Label: "reading the estate", Detail: "one Resource Graph query",
	}, 44, 6, when)
	if !strings.Contains(got, "one Resource Graph query") {
		t.Errorf("the detail is missing:\n%s", got)
	}
	lines := strings.Split(got, "\n")
	var label, detail int
	for i, l := range lines {
		if strings.Contains(l, "reading the estate") {
			label = i
		}
		if strings.Contains(l, "Resource Graph") {
			detail = i
		}
	}
	if detail != label+1 {
		t.Errorf("the detail is on row %d and the label on %d", detail, label)
	}
}

// A pane squeezed by a narrow terminal is an ordinary state, not a failure.
func TestWaitingInNoRoom(t *testing.T) {
	for _, r := range []Rect{{W: 0, H: 0}, {W: 10, H: 0}, {W: 0, H: 4}} {
		c := NewCanvas(10, 4)
		Waiting{Label: "reading"}.Draw(c, r, when, Region("wait"))
		if strings.TrimSpace(c.String()) != "" {
			t.Errorf("rect %v drew something", r)
		}
	}
}

// One row is enough for the label; the detail is what gets dropped.
func TestWaitingInOneRow(t *testing.T) {
	got := drawWaiting(Waiting{Label: "reading", Detail: "dropped"}, 30, 1, when)
	if !strings.Contains(got, "reading") {
		t.Errorf("the label was lost in one row:\n%q", got)
	}
	if strings.Contains(got, "dropped") {
		t.Errorf("the detail was drawn with no room for it:\n%q", got)
	}
}

func TestWaitingTruncatesRatherThanOverflows(t *testing.T) {
	got := drawWaiting(Waiting{Label: strings.Repeat("long ", 20)}, 20, 3, when)
	for _, line := range strings.Split(got, "\n") {
		if Width(line) > 20 {
			t.Errorf("a line is %d columns wide in a 20-column rect", Width(line))
		}
	}
}
