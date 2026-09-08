package comp

import (
	"time"

	"github.com/charmbracelet/lipgloss"
)

// Spinner is a turning indicator for work in flight.
//
// # Where this came from
//
// The database tool's spinner and the cloud tool's ⠿. The rule worth keeping
// is the database tool's, and its comment states it: the frame comes from THE
// CLOCK rather than from a counter, "so every spinner on screen turns together
// and at a steady rate however often the view is rendered".
//
// That is not a nicety. A counter advanced in View turns faster when more is
// happening, which is exactly backwards, and two spinners on one screen drift
// apart until the interface looks broken. A clock-driven spinner also makes a
// captured frame deterministic, because the harness freezes the clock — a
// counter would make every golden depend on how many times View had been
// called.
//
// The cloud tool's single ⠿ is a spinner that has stopped, which reads as hung
// rather than as working. Not carried.
type Spinner struct {
	// Frames turn in one direction at one dot per frame, so a dropped frame
	// slows it rather than reversing it. Empty takes the braille set.
	Frames []string
	// Every is how long a frame lasts. Zero takes 100ms, which is fast enough
	// to read as motion and slow enough not to strobe.
	Every time.Duration
	Style *lipgloss.Style
}

// The default frames, from the database tool. Braille, so the dot travels
// round a ring rather than flickering between unrelated shapes.
var spinnerFrames = []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"}

const spinnerEvery = 100 * time.Millisecond

// Frame is the glyph for a moment.
func (s Spinner) Frame(at time.Time) string {
	frames := s.Frames
	if len(frames) == 0 {
		frames = spinnerFrames
	}
	every := s.Every
	if every == 0 {
		every = spinnerEvery
	}
	return frames[int(at.UnixNano()/int64(every))%len(frames)]
}

// Draw puts the frame at the top left of r and returns the columns used.
func (s Spinner) Draw(c *Canvas, r Rect, at time.Time, id ID) int {
	c = c.Clip(r)
	return c.Text(r.X, r.Y, s.Frame(at), s.Style, id)
}
