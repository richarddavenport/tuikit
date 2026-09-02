package ui

import (
	"testing"

	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
	"github.com/richarddavenport/tuikit/harness"
)

// Every screen democtl can show, captured as a frame you can look at.
//
// Gated, so a plain `go test ./...` skips it — rendering frames is something
// you ask for:
//
//	DEMOCTL_FRAMES=/tmp/frames go test ./examples/democtl/ui -run CaptureFrames
//
// democtl is a fixture by construction: the fleet comes from a seed and the
// clock is fixed at fleet.Epoch, so these frames say the same thing tomorrow.
// A tool with a backend would capture the same screens in live mode as well,
// which is what finds the things a fixture hides.
func TestCaptureFrames(t *testing.T) {
	dir := harness.Enabled("DEMOCTL_FRAMES")
	if dir == "" {
		t.Skip("set DEMOCTL_FRAMES to capture frames")
	}

	s := harness.Capture(t, dir,
		harness.Size(132, 38),
		harness.At(fleet.Epoch),
		harness.WithMode(harness.Fixture))

	for _, st := range states() {
		s.Shot(st.name, run(st.build()))
	}
	t.Log(s)
	s.Done()
}

// The same screens as goldens, so a layout change is an ordinary test failure
// rather than something somebody has to notice.
//
// Not gated: this is the half that runs in CI. Regenerate deliberately with
// `go test ./examples/democtl/ui -update-goldens` when a change is intended.
func TestFramesMatchTheirGoldens(t *testing.T) {
	for _, st := range states() {
		t.Run(st.name, func(t *testing.T) {
			m := st.build()
			m.SetSize(132, 38)
			m.Now(fleet.Epoch)
			harness.Golden(t, "testdata", st.name, view(m))
		})
	}
}

// A frame has to survive a narrow terminal too, and 80 columns is where a TUI
// stops being able to hide.
func TestFramesAtEightyColumns(t *testing.T) {
	for _, st := range states() {
		t.Run(st.name, func(t *testing.T) {
			m := st.build()
			m.SetSize(80, 24)
			m.Now(fleet.Epoch)
			harness.Golden(t, "testdata/80", st.name, view(m))
		})
	}
}

// Colour must not change the shape.
//
// The half of the golden suite that the goldens cannot be: they run uncoloured,
// where a helper that measures a styled string by counting runes is correct.
// This draws every screen both ways and compares what a reader sees — which is
// how democtl's tab strip was caught printing "‹ Config …" at 80 columns in a
// real terminal while its golden showed the whole strip.
func TestColourDoesNotChangeTheShape(t *testing.T) {
	for _, size := range []struct {
		name string
		w, h int
	}{
		{"132x38", 132, 38},
		{"80x24", 80, 24},
	} {
		t.Run(size.name, func(t *testing.T) {
			for _, st := range states() {
				harness.ShapeSurvivesColour(t, st.name, func() string {
					m := st.build()
					m.SetSize(size.w, size.h)
					m.Now(fleet.Epoch)
					return view(m)
				})
			}
		})
	}
}
