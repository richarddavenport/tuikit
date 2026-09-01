package ui

import (
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
	"github.com/richarddavenport/tuikit/harness"
	"github.com/richarddavenport/tuikit/prototype/htmlround"
)

// Every democtl screen, rendered to HTML and read back.
//
// Part of the throwaway prototype in prototype/htmlround, and it lives here
// rather than there because states() is what makes it worth running: fourteen
// screens nobody designed to be easy to parse, including a modal over a table,
// a log pane with two colours per line, and a filter that matches nothing.
//
// The synthetic table in the prototype tests the constructs one at a time.
// This tests whether they interact.
func TestPrototypeFramesSurviveTheTripThroughHTML(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)

	for _, size := range []struct {
		name string
		w, h int
	}{
		{"132x38", 132, 38},
		{"80x24", 80, 24},
	} {
		t.Run(size.name, func(t *testing.T) {
			for _, st := range states() {
				t.Run(st.name, func(t *testing.T) {
					m := st.build()
					m.SetSize(size.w, size.h)
					m.Now(fleet.Epoch)
					frame := m.View()

					got, err := htmlround.FromHTML(harness.HTML(frame))
					if err != nil {
						t.Fatalf("parsing the page: %v", err)
					}
					if d := htmlround.Diff(htmlround.FromANSI(frame), got); d != "" {
						t.Errorf("%s did not survive the trip:\n%s", st.name, d)
					}
				})
			}
		})
	}
}
