package ui

import (
	"github.com/richarddavenport/tuikit/app"
)

// view draws the model the way the runner does.
//
// The model has no View since issue 21 — the runner owns the canvas — so a test
// that wants the frame as a string asks for it here. The runner is built at the
// model's own size so that a test which resized the model still gets the frame
// it was describing.
func view(m *Model) string { return run(m).View() }

// run wraps a model so the harness can drive it: Update, View and Canvas all
// belong to the runner now.
func run(m *Model) *app.Runner { return app.New(m, app.WithSize(m.width, m.height)) }
