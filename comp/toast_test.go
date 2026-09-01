package comp

import (
	"strings"
	"testing"
)

const toastName Name = "toast"

// A hint is the difference between telling someone they have a problem and
// telling them what to do about it. swarmctl's error view carries "a hint for
// the failures we know how to fix", and it is the field most likely to be left
// off — so it is a field rather than something you append to the body.
func TestAToastCanSayWhatToDoAboutIt(t *testing.T) {
	c := NewCanvas(60, 14)
	Toast{
		Title: "cannot reach staging",
		Body:  "dial tcp 10.0.0.4:5432: connection refused",
		Hint:  "is the tunnel up? try `pgctl connect staging`",
	}.Draw(c, c.Bounds(), Region(toastName))

	got := c.String()
	for _, want := range []string{"cannot reach staging", "connection refused", "pgctl connect"} {
		if !strings.Contains(got, want) {
			t.Errorf("%q is missing:\n%s", want, got)
		}
	}
}

// A toast with no hint is allowed, and usually means nobody has worked out the
// answer yet — which is worth being able to see.
func TestAToastWithoutAHintIsShorter(t *testing.T) {
	with := NewCanvas(60, 14)
	tall := Toast{Title: "failed", Body: "something broke", Hint: "try again"}.
		Draw(with, with.Bounds(), Region(toastName))

	without := NewCanvas(60, 14)
	short := Toast{Title: "failed", Body: "something broke"}.
		Draw(without, without.Bounds(), Region(toastName))

	if short.H >= tall.H {
		t.Errorf("a toast with no hint is %d rows and one with is %d", short.H, tall.H)
	}
}

// Bounded like a modal, and for the same reason: a root cause is arbitrary text
// from somewhere else.
func TestAToastIsBoundedToItsContainer(t *testing.T) {
	for _, w := range []int{200, 80, 40, 20} {
		c := NewCanvas(w, 16)
		Toast{
			Title: "cannot reach staging",
			Body:  strings.Repeat("a root cause that came from somewhere else and is very long. ", 4),
		}.Draw(c, c.Bounds(), Region(toastName))

		for i, line := range strings.Split(c.String(), "\n") {
			if Width(line) > w {
				t.Errorf("at %d columns, line %d is %d wide", w, i+1, Width(line))
			}
		}
	}
}

// The corner is the caller's, because what a toast must not cover depends on
// what is underneath it.
func TestAToastSitsInTheCornerItIsGiven(t *testing.T) {
	for _, tc := range []struct {
		anchor          Anchor
		leftish, topish bool
	}{
		{TopRight, false, true},
		{TopLeft, true, true},
		{BottomRight, false, false},
		{BottomLeft, true, false},
	} {
		c := NewCanvas(80, 24)
		got := Toast{Title: "note", Body: "short", Anchor: tc.anchor}.Draw(c, c.Bounds(), Region(toastName))

		// Which side it is ON, by comparing the space either side of it. A
		// 48-column box on an 80-column canvas starts at x=30 and is still
		// right-anchored, so an absolute threshold tests the wrong thing.
		before, after := got.X, 79-got.Right()
		above, below := got.Y, 23-got.Bottom()
		if left := before <= after; left != tc.leftish {
			t.Errorf("anchor %d left %d columns before it and %d after", tc.anchor, before, after)
		}
		if top := above <= below; top != tc.topish {
			t.Errorf("anchor %d left %d rows above it and %d below", tc.anchor, above, below)
		}
	}
}

// Drawn over whatever is behind it, like any other box.
func TestAToastCoversWhatIsBehindIt(t *testing.T) {
	c := NewCanvas(60, 12)
	c.Fill(c.Bounds(), "x", nil, Region("behind"))
	r := Toast{Title: "failed", Body: "something broke"}.Draw(c, c.Bounds(), Region(toastName))

	// Inside the box, cell by cell — the x's on either side of it are meant
	// to still be there.
	for x := r.X + 1; x < r.Right(); x++ {
		if cell, ok := c.CellAt(x, r.Y+1); ok && cell.Text == "x" {
			t.Fatalf("the frame behind shows through at column %d", x)
		}
	}
}
