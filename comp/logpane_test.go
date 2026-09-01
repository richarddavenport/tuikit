package comp

import (
	"strings"
	"testing"
)

const logs Name = "logs.row"

func logLines(n int) []LogLine {
	out := make([]LogLine, n)
	for i := range out {
		out[i] = LogLine{At: "09:14:0" + itoa(i%10), Text: "line " + itoa(i)}
	}
	return out
}

func drawLogs(p *LogPane, w, h, n int) *Canvas {
	c := NewCanvas(w, h)
	p.Draw(c, c.Bounds(), logLines(n), logs)
	return c
}

// A log you have just opened is one you want the end of. democtl's opened on
// the oldest lines in the buffer, which for a log is the wrong end.
func TestFollowingShowsTheNewestLines(t *testing.T) {
	p := &LogPane{Follow: true}
	got := drawLogs(p, 40, 6, 20).String()

	if !strings.Contains(got, "line 19") {
		t.Errorf("following does not show the newest line:\n%s", got)
	}
	if strings.Contains(got, "line 0 ") {
		t.Errorf("following is showing the oldest line:\n%s", got)
	}
}

// Scrolling up detaches; scrolling back to the bottom re-attaches. That is what
// makes following a place rather than a mode.
func TestScrollingDetachesAndArrivingBackReattaches(t *testing.T) {
	p := &LogPane{Follow: true}
	drawLogs(p, 40, 6, 20)

	p.Scroll(-3)
	if p.Follow {
		t.Error("scrolling up did not detach from the tail")
	}
	drawLogs(p, 40, 6, 20)
	if p.Offset() == p.bottom() {
		t.Error("the view did not move")
	}

	p.Scroll(3)
	if !p.Follow {
		t.Error("scrolling back to the bottom did not re-attach")
	}
}

// Following pins to the tail as lines arrive.
func TestFollowingStaysAtTheTailAsLinesArrive(t *testing.T) {
	p := &LogPane{Follow: true}
	drawLogs(p, 40, 6, 20)
	got := drawLogs(p, 40, 6, 40).String()

	if !strings.Contains(got, "line 39") {
		t.Errorf("following did not keep up with the stream:\n%s", got)
	}
}

// The status is drawn always. A pane that only says where it is while it is
// moving is one you cannot tell from a pane that is stuck.
func TestTheStatusSaysWhetherItIsAtTheTail(t *testing.T) {
	p := &LogPane{Follow: true}
	if got := drawLogs(p, 40, 6, 20).String(); !strings.Contains(got, "following") {
		t.Errorf("a following pane does not say so:\n%s", got)
	}

	p.Scroll(-4)
	if got := drawLogs(p, 40, 6, 20).String(); !strings.Contains(got, "4 below") {
		t.Errorf("a detached pane does not say how far from the tail it is:\n%s", got)
	}
}

func TestScrollingStopsAtBothEnds(t *testing.T) {
	p := &LogPane{Follow: true}
	drawLogs(p, 40, 6, 20)

	p.Scroll(-100)
	if p.Offset() != 0 {
		t.Errorf("scrolling up stopped at %d", p.Offset())
	}
	drawLogs(p, 40, 6, 20)
	p.Scroll(100)
	if !p.Follow {
		t.Error("scrolling past the bottom did not re-attach")
	}
}

// A stream with nothing in it is an ordinary state, not an error.
func TestAnEmptyStreamSaysSo(t *testing.T) {
	p := &LogPane{Follow: true, Empty: "  no lines on this stream"}
	c := NewCanvas(40, 6)
	p.Draw(c, c.Bounds(), nil, logs)

	if !strings.Contains(c.String(), "no lines on this stream") {
		t.Errorf("an empty stream says nothing:\n%s", c.String())
	}
}

// A long line is truncated to the pane rather than pushing its border out.
func TestALongLineIsTruncatedToThePane(t *testing.T) {
	c := NewCanvas(30, 4)
	p := &LogPane{Follow: true}
	p.Draw(c, c.Bounds(), []LogLine{{At: "09:14:03", Text: strings.Repeat("wide ", 30)}}, logs)

	for _, line := range strings.Split(c.String(), "\n") {
		if Width(line) > 30 {
			t.Errorf("a line drew %d columns: %q", Width(line), line)
		}
	}
}

// Each line is clickable as itself, by its index in the stream.
func TestEachLineOwnsItsRow(t *testing.T) {
	p := &LogPane{Follow: true}
	c := drawLogs(p, 40, 6, 20)

	top := c.OwnerAt(2, 0)
	if top.Index != p.Offset() {
		t.Errorf("the top row is owned by %v, want index %d", top, p.Offset())
	}
}
