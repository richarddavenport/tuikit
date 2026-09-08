package harness

import (
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// ShapeSurvivesColor asserts that turning color on does not change what a
// frame says.
//
// This exists because the goldens cannot see a whole class of bug. They run
// under a plain `go test`, where there is no TTY and lipgloss emits no escape
// sequences at all — so any helper that measures a styled string by counting
// RUNES is correct in every golden and wrong in every terminal. Capture forces
// the TrueColor profile, so the .ansi frames and the real tool have a bug the
// goldens say is not there.
//
// democtl had exactly that: a tab strip padded with a rune-counting helper,
// cut mid-escape at 80 columns, eating the sequence that turned bold off. The
// golden showed `‹ Config ›  Events` and looked perfect.
//
// Color is decoration. If it changes the SHAPE, something measured bytes that
// it should have measured in columns, and this says so.
//
// render must build and draw the frame from scratch each time it is called:
// lipgloss resolves color at Render, so the same model is drawn twice, once
// under each profile.
func ShapeSurvivesColor(t T, name string, render func() string) {
	t.Helper()

	// The profile is process-wide, so it is put back. A helper that leaves
	// color switched on hands the next test in the package a frame full of
	// escape sequences, and the failure lands nowhere near the cause.
	profile, dark := lipgloss.ColorProfile(), lipgloss.HasDarkBackground()
	defer func() {
		lipgloss.SetColorProfile(profile)
		lipgloss.SetHasDarkBackground(dark)
	}()

	lipgloss.SetColorProfile(termenv.Ascii)
	plain := render()

	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)
	colored := Strip(render())

	if plain != colored {
		t.Errorf("%s says something different once color is on:\n%s", name, diff(plain, colored))
	}
}
