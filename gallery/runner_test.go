package gallery

import "github.com/richarddavenport/tuikit/app"

// view draws the gallery the way the runner does.
//
// The model has no View since issue 21 — the runner owns the canvas — so a test
// that wants the frame as a string asks for it here, at the size the model was
// told about.
func view(m *Model) string { return run(m).View() }

func run(m *Model) *app.Runner { return app.New(m, app.WithSize(m.width, m.height)) }
