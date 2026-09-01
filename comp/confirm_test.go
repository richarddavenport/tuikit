package comp

import (
	"strings"
	"testing"
)

const confirmName Name = "confirm"

func confirmOn(w, h int, cf Confirm) (*Canvas, Rect) {
	c := NewCanvas(w, h)
	c.Fill(c.Bounds(), "x", nil, Region("behind"))
	return c, cf.Draw(c, Region(confirmName))
}

// Bounded to its container by construction. An unbounded box drew itself off
// the side of the screen in a real tool, and the screenshots caught it before
// anyone else did.
func TestAModalIsBoundedToItsContainer(t *testing.T) {
	for _, w := range []int{200, 132, 80, 40, 24, 10} {
		c, r := confirmOn(w, 24, Confirm{
			Title: "Deploy api_gateway?",
			Body:  strings.Repeat("a very long explanation that would happily run off the side. ", 6),
		})
		if r.X+r.W > w && w >= confirmMin {
			t.Errorf("at %d columns the modal spans %d..%d", w, r.X, r.X+r.W)
		}
		for i, line := range strings.Split(c.String(), "\n") {
			if Width(line) > w {
				t.Errorf("at %d columns line %d is %d wide", w, i+1, Width(line))
			}
		}
	}
}

// pgctl floors the width at 32 and democtl does not, so democtl's modal shrinks
// on a narrow terminal until the question no longer reads. A box wider than the
// terminal is clipped by the canvas, which is a better failure than one too
// narrow to read.
func TestAModalHasAFloorAsWellAsACap(t *testing.T) {
	_, narrow := confirmOn(30, 20, Confirm{Title: "Remove?", Body: "gone for good"})
	if narrow.W < confirmMin {
		t.Errorf("on a 30-column terminal the modal is %d wide, under the %d floor", narrow.W, confirmMin)
	}

	_, wide := confirmOn(300, 20, Confirm{Title: "Remove?", Body: "gone for good"})
	if wide.W > confirmMax {
		t.Errorf("on a 300-column terminal the modal is %d wide, over the %d cap", wide.W, confirmMax)
	}
}

// Drawn last is on top, and the box covers what it sits on rather than
// compositing with it.
func TestAModalCoversWhatIsBehindIt(t *testing.T) {
	c, r := confirmOn(60, 20, Confirm{Title: "Deploy?", Body: "one at a time"})

	row := strings.Split(c.String(), "\n")[r.Y+1]
	if strings.Contains(row[len(row)/3:2*len(row)/3], "x") {
		t.Errorf("the frame behind shows through the modal: %q", row)
	}
}

// The body wraps to the box rather than being cut off at the border.
func TestTheBodyWrapsInsideTheBox(t *testing.T) {
	c, _ := confirmOn(80, 24, Confirm{
		Title: "Deploy api_gateway?",
		Body: "Pushes a new service spec and waits for the tasks to converge. " +
			"The current tasks are replaced one at a time.",
	})
	got := c.String()
	if !strings.Contains(got, "Pushes a new service spec") {
		t.Errorf("the body is missing:\n%s", got)
	}
	if !strings.Contains(got, "one at a time.") {
		t.Errorf("the end of the body was cut off:\n%s", got)
	}
}

// A modal grows for its content, so the keys are never pushed off the bottom.
func TestTheKeysAreAlwaysBelowTheBody(t *testing.T) {
	short, rs := confirmOn(80, 24, Confirm{Title: "Go?", Body: "short", Hints: []Hint{{"y", "confirm"}}})
	long, rl := confirmOn(80, 24, Confirm{
		Title: "Go?",
		Body:  strings.Repeat("a longer explanation that needs several lines. ", 4),
		Hints: []Hint{{"y", "confirm"}},
	})
	if rl.H <= rs.H {
		t.Errorf("a longer body did not make a taller box: %d then %d", rs.H, rl.H)
	}
	for _, c := range []*Canvas{short, long} {
		if !strings.Contains(c.String(), "y confirm") {
			t.Errorf("the keys are missing:\n%s", c.String())
		}
	}
}
