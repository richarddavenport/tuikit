package comp

import (
	"testing"
	"time"
)

// Every spinner on screen turns together and at a steady rate, however often
// the view is rendered. A counter advanced in View turns faster when more is
// happening — exactly backwards — and two spinners drift apart until the
// interface looks broken.
func TestSpinnersTurnTogetherWhateverTheRenderRate(t *testing.T) {
	at := time.Unix(0, 0).Add(350 * time.Millisecond)
	one, two := Spinner{}, Spinner{}

	// One of them has been "rendered" many more times; it makes no difference,
	// because neither is counting.
	for range 40 {
		two.Frame(at)
	}
	if one.Frame(at) != two.Frame(at) {
		t.Errorf("two spinners disagree at the same moment: %q and %q", one.Frame(at), two.Frame(at))
	}
}

// The clock drives it, so a frozen clock gives a fixed frame and a captured
// golden does not depend on how many times View was called.
func TestAFrozenClockGivesAFixedFrame(t *testing.T) {
	at := time.Unix(0, 0).Add(700 * time.Millisecond)
	first := Spinner{}.Frame(at)
	for range 5 {
		if got := (Spinner{}).Frame(at); got != first {
			t.Errorf("the frame moved without the clock: %q then %q", first, got)
		}
	}
}

func TestItTurnsAsTheClockRuns(t *testing.T) {
	seen := map[string]bool{}
	base := time.Unix(0, 0)
	for i := range 10 {
		seen[Spinner{}.Frame(base.Add(time.Duration(i)*100*time.Millisecond))] = true
	}
	if len(seen) != len(spinnerFrames) {
		t.Errorf("a second of turning showed %d of %d frames", len(seen), len(spinnerFrames))
	}
}

func TestTheFramesAreTheCallersIfItHasSome(t *testing.T) {
	s := Spinner{Frames: []string{"-", "\\", "|", "/"}, Every: time.Second}
	if got := s.Frame(time.Unix(2, 0)); got != "|" {
		t.Errorf("got %q, want %q", got, "|")
	}
}

func TestDrawStaysInItsRect(t *testing.T) {
	c := NewCanvas(10, 1)
	if n := (Spinner{}).Draw(c, Rect{X: 0, Y: 0, W: 1, H: 1}, time.Unix(0, 0), Region("demo")); n != 1 {
		t.Errorf("drew %d columns", n)
	}
	if got := Width(c.String()); got != 1 {
		t.Errorf("the row is %d columns: %q", got, c.String())
	}
}
