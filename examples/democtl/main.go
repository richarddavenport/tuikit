// Command democtl is tuikit's example tool and its own test fixture.
//
// It manages a fictional service fleet generated from a seed, so it needs no
// backend and shows the same thing every run. That is what lets tuikit's capture
// harness render it as a fixture, and what lets a reader compare a frame in the
// documentation against what they see on their own terminal.
//
//	go run ./examples/democtl
//
// It is also the reference implementation: the engine/ui split, the palette and
// glyph set held closed by guards, the mode split that stops a list scrolling
// behind a modal, and the generation counter that stops an abandoned run
// drawing into its successor.
package main

import (
	"flag"
	"fmt"
	"os"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/examples/democtl/ui"
)

func main() {
	seed := flag.Int64("seed", 1, "which fleet to generate; the same seed is the same fleet")
	flag.Parse()

	if _, err := tea.NewProgram(ui.New(*seed), tea.WithAltScreen()).Run(); err != nil {
		fmt.Fprintln(os.Stderr, "democtl:", err)
		os.Exit(1)
	}
}
