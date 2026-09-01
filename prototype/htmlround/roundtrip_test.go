package htmlround

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/richarddavenport/tuikit/harness"
)

// forceColour makes lipgloss emit sequences into a pipe, the same way
// harness.Capture does. Without it a test binary renders everything plain and
// the round trip passes by having nothing to lose.
func forceColour() {
	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)
}

// roundTrip is the property under test: a frame, rendered to HTML and read
// back, is the same frame.
func roundTrip(t *testing.T, frame string) {
	t.Helper()

	want := FromANSI(frame)
	got, err := FromHTML(harness.HTML(frame))
	if err != nil {
		t.Fatalf("parsing the page: %v", err)
	}
	if d := Diff(want, got); d != "" {
		t.Errorf("the frame did not survive the trip to HTML and back:\n%s", d)
	}
}

// The forms a tuikit frame is made of, one at a time, so a failure names the
// construct rather than a screen.
func TestTheFormsAFrameIsMadeOf(t *testing.T) {
	forceColour()

	role := func(index string) lipgloss.Style {
		return lipgloss.NewStyle().Foreground(lipgloss.Color(index))
	}

	for _, tc := range []struct {
		name  string
		frame string
	}{
		{"plain text", "api_gateway  3/3"},
		{"a palette foreground", role("205").Render("selected")},
		{"bold", lipgloss.NewStyle().Bold(true).Render("HEADER")},
		{"a selection background", lipgloss.NewStyle().
			Foreground(lipgloss.Color("229")).
			Background(lipgloss.Color("57")).
			Render(" api_gateway ")},
		{"two colours on one line", role("34").Render("ok") + "  " + role("196").Render("failed")},
		{"box drawing", "┌──────────┐\n│ services │\n└──────────┘"},
		{"the glyph set", role("241").Render("↑↓ move · ⏎ select · q quit")},
		{"an empty line between two full ones", role("205").Render("a") + "\n\n" + role("205").Render("b")},
		{"trailing spaces", lipgloss.NewStyle().Background(lipgloss.Color("57")).Render("row   ")},
		{"markup in the content", "if a < b && c > d"},
	} {
		t.Run(tc.name, func(t *testing.T) { roundTrip(t, tc.frame) })
	}
}

// A 24-bit colour is not something theme produces — the palette is nine 256
// indices — but it is something a tool author writes the moment they reach for
// a hex value, and lipgloss emits it verbatim under a TrueColor profile.
//
// Split out from the table above because it is a KNOWN loss rather than a
// property: harness.css handles 38;5;N and falls through on 38;2;R;G;B, so the
// colour reaches the page as nothing at all. Written as a test so the day it
// starts passing is recorded.
func TestTwentyFourBitColourIsLost(t *testing.T) {
	forceColour()
	frame := lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5faf")).Render("hand-picked")

	got, err := FromHTML(harness.HTML(frame))
	if err != nil {
		t.Fatalf("parsing the page: %v", err)
	}
	if d := Diff(FromANSI(frame), got); d == "" {
		t.Log("24-bit colour now survives — fold this back into the table above")
	} else {
		t.Logf("known loss, unfixed:\n%s", d)
	}
}

// A style left open at the end of a line continues onto the next one in a
// terminal, and harness.writeLine closes every span at a line boundary without
// reopening it. Whether that is a loss depends on whether anything emits it.
func TestAStyleThatStraddlesALineBreak(t *testing.T) {
	forceColour()
	frame := "\x1b[48;5;57mrow one\nrow two\x1b[0m"

	got, err := FromHTML(harness.HTML(frame))
	if err != nil {
		t.Fatalf("parsing the page: %v", err)
	}
	if d := Diff(FromANSI(frame), got); d != "" {
		t.Logf("a background does not cross the line break:\n%s", d)
	}
}

// The page must not carry escape sequences, whatever else it does.
func TestNoEscapeReachesThePage(t *testing.T) {
	forceColour()
	frame := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("selected")

	if page := harness.HTML(frame); strings.Contains(page, "\x1b") {
		t.Errorf("an escape sequence reached the page:\n%s", page)
	}
}
