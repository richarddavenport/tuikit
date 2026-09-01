package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

const run Name = "run.step"

func look() [5]StepLook {
	return [5]StepLook{
		StepWaiting: {Glyph: "•"},
		StepRunning: {Glyph: "→"},
		StepSkipped: {Glyph: "●"},
		StepDone:    {Glyph: "✓"},
		StepFailed:  {Glyph: "✗"},
	}
}

// A duration is a number you scan down a column, so it goes right. Prose about
// a step sits beside it. They are two things, not one thing badly agreed on.
func TestADurationGoesRightAndProseStaysBeside(t *testing.T) {
	c := NewCanvas(40, 4)
	StepList{Look: look(), Steps: []Step{
		{Label: "pull image", State: StepDone, Took: "1.2s"},
		{Label: "health check", State: StepSkipped, Detail: "already true"},
	}}.Draw(c, c.Bounds(), run)

	lines := strings.Split(c.String(), "\n")
	// One column short of the edge: the trailing blank is real in the cells
	// and trimmed on the way out, which is what every other right-aligned
	// thing in a frame does.
	if !strings.HasSuffix(lines[0], "1.2s") || Width(lines[0]) != 39 {
		t.Errorf("the duration is not against the right edge: %q (%d columns)", lines[0], Width(lines[0]))
	}
	if !strings.Contains(lines[1], "health check  already true") {
		t.Errorf("the reason is not beside its step: %q", lines[1])
	}
}

// Right-aligned against the BOX, not the terminal. Measuring against the
// terminal put the durations two columns past the edge, where they clipped to
// "400m…" — visible in a captured frame and in nothing else.
func TestTheDurationColumnIsAgainstTheRectNotTheCanvas(t *testing.T) {
	c := NewCanvas(60, 3)
	inner := Rect{X: 2, Y: 0, W: 30, H: 3}
	StepList{Look: look(), Steps: []Step{{Label: "pull image", State: StepDone, Took: "1.2s"}}}.
		Draw(c, inner, run)

	line := strings.Split(c.String(), "\n")[0]
	if Width(line) > inner.X+inner.W {
		t.Errorf("the duration ran past the box's edge at %d: %q", inner.X+inner.W, line)
	}
	if !strings.Contains(line, "1.2s") {
		t.Errorf("the duration is missing: %q", line)
	}
}

// A failure owes an explanation, on a line of its own.
func TestAFailureGetsALineForItsReason(t *testing.T) {
	c := NewCanvas(60, 4)
	StepList{Look: look(), Steps: []Step{
		{Label: "wait for tasks", State: StepFailed, Note: "task 3 exited 137 before the health check passed"},
	}}.Draw(c, c.Bounds(), run)

	lines := strings.Split(c.String(), "\n")
	if !strings.Contains(lines[0], "✗ wait for tasks") {
		t.Errorf("the step is wrong: %q", lines[0])
	}
	if !strings.Contains(lines[1], "exited 137") {
		t.Errorf("the reason is not on its own line: %q", lines[1])
	}
}

// The glyphs belong to the tool. A component that hard-codes ✓ has taken a
// decision belonging to whoever owns the glyph set, and would smuggle a
// character past guard.Glyphs.
func TestTheGlyphsComeFromTheCaller(t *testing.T) {
	c := NewCanvas(40, 3)
	ascii := [5]StepLook{
		StepDone:   {Glyph: "[ok]"},
		StepFailed: {Glyph: "[!!]"},
	}
	StepList{Look: ascii, Steps: []Step{
		{Label: "one", State: StepDone},
		{Label: "two", State: StepFailed},
	}}.Draw(c, c.Bounds(), run)

	got := c.String()
	if !strings.Contains(got, "[ok] one") || !strings.Contains(got, "[!!] two") {
		t.Errorf("the caller's glyphs were not used:\n%s", got)
	}
	if strings.ContainsAny(got, "✓✗") {
		t.Errorf("the component drew a glyph of its own:\n%s", got)
	}
}

// Each step is clickable as itself.
func TestEachStepOwnsItsRow(t *testing.T) {
	c := NewCanvas(40, 4)
	StepList{Look: look(), Steps: []Step{
		{Label: "one", State: StepDone},
		{Label: "two", State: StepRunning},
	}}.Draw(c, c.Bounds(), run)

	for i := range 2 {
		if got := c.OwnerAt(4, i); got.Index != i {
			t.Errorf("row %d is owned by %v", i, got)
		}
	}
}

func TestTheStatusLineFollowsTheSteps(t *testing.T) {
	forceColour()
	danger := lipgloss.NewStyle().Foreground(lipgloss.Color("196"))

	c := NewCanvas(60, 6)
	StepList{
		Look:        look(),
		Steps:       []Step{{Label: "one", State: StepFailed}},
		Status:      "stopped",
		StatusStyle: &danger,
		Hints:       []Hint{{Key: "r", Label: "run again"}, {Key: "esc", Label: "go back"}},
	}.Draw(c, c.Bounds(), run)

	got := c.String()
	if !strings.Contains(got, "stopped") {
		t.Errorf("no status:\n%s", got)
	}
	if !strings.Contains(got, "r run again · esc go back") {
		t.Errorf("the status hints are missing or not joined:\n%s", got)
	}
}

// A list longer than its pane stops at the edge rather than drawing past it.
func TestAListTallerThanItsRectStops(t *testing.T) {
	steps := make([]Step, 20)
	for i := range steps {
		steps[i] = Step{Label: "step", State: StepDone}
	}
	c := NewCanvas(40, 5)
	StepList{Look: look(), Steps: steps, Status: "done"}.Draw(c, Rect{X: 0, Y: 0, W: 40, H: 3}, run)

	lines := strings.Split(c.String(), "\n")
	for i := 3; i < len(lines); i++ {
		if strings.TrimSpace(lines[i]) != "" {
			t.Errorf("line %d is outside the rect: %q", i+1, lines[i])
		}
	}
}
