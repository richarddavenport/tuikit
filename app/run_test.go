package app_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/app"
	"github.com/richarddavenport/tuikit/comp"
)

// The reserved bottom row is a choice, and both sides of it are drawn.
//
// Issue 37: nobody had written down why the canvas is one row short, and the
// answer was that nobody decided it. Measured in tmux 3.5a under alt-screen, a
// full-height frame with a bordered pane renders with its top line intact — so
// the row is a hedge rather than a requirement, and a tool can have it back.
func TestFullHeightUsesTheLastRow(t *testing.T) {
	for _, tc := range []struct {
		name string
		opts []app.Option
		want int
	}{
		{"the default keeps a row back", nil, 9},
		{"WithFullHeight uses it", []app.Option{app.WithFullHeight()}, 10},
	} {
		t.Run(tc.name, func(t *testing.T) {
			r := app.New(sizeModel{}, tc.opts...)
			r.Update(tea.WindowSizeMsg{Width: 20, Height: 10})
			r.View()

			if got := r.Canvas().Bounds().H; got != tc.want {
				t.Errorf("the canvas is %d rows in a 10-row terminal, want %d", got, tc.want)
			}
		})
	}
}

// sizeModel draws nothing; the canvas's size is the whole assertion.
type sizeModel struct{}

func (sizeModel) Init() tea.Cmd                         { return nil }
func (m sizeModel) Update(tea.Msg) (app.Model, tea.Cmd) { return m, nil }
func (sizeModel) Draw(*comp.Canvas, comp.Rect)          {}
