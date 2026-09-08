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
	forceColor()
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

// A run longer than its pane is a List drawing StepList's rows. One layout,
// two ways of showing it — which is the answer comp.Table already gives.
func TestRowsAreTheSameLayoutDrawGives(t *testing.T) {
	steps := StepList{Look: look(), Steps: []Step{
		{Label: "create the machine", State: StepDone, Took: "4.2s"},
		{Label: "install docker", State: StepRunning},
		{Label: "already there", State: StepSkipped, Detail: "already satisfied"},
	}}

	drawn := NewCanvas(40, 4)
	steps.Draw(drawn, drawn.Bounds(), run)

	scrolled := NewCanvas(40, 5)
	l := &List{Name: run}
	l.Draw(scrolled, scrolled.Bounds(), steps.Rows(40))

	want := strings.Split(drawn.String(), "\n")[:3]
	got := strings.Split(scrolled.String(), "\n")[:3]
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("row %d:\n drawn %q\n  list %q", i, want[i], got[i])
		}
	}
}

// The duration is right-aligned against the WIDTH IT IS GIVEN. democtl measured
// against the terminal instead of the box, which put durations two columns past
// the edge where they clipped to "400m…".
func TestRowsRightAlignTheDurationAgainstTheirWidth(t *testing.T) {
	steps := StepList{Look: look(), Steps: []Step{{Label: "pull", State: StepDone, Took: "1.2s"}}}

	for _, w := range []int{30, 60} {
		c := NewCanvas(w, 1)
		l := &List{Name: run}
		l.Draw(c, c.Bounds(), steps.Rows(w))

		line := strings.Split(c.String(), "\n")[0]
		if !strings.HasSuffix(strings.TrimRight(line, " "), "1.2s") {
			t.Errorf("at %d columns the duration is not at the right edge: %q", w, line)
		}
		if Width(line) > w {
			t.Errorf("at %d columns the row is %d wide", w, Width(line))
		}
	}
}

// A scrolling step list still owns its rows by step index, so a click acts on
// the step it landed on after the viewport has moved.
func TestAScrolledStepListStillOwnsItsRows(t *testing.T) {
	var steps []Step
	for i := range 20 {
		steps = append(steps, Step{Label: "step " + itoa(i), State: StepWaiting})
	}
	list := StepList{Look: look(), Steps: steps}

	c := NewCanvas(40, 6)
	l := &List{Name: run}
	l.Draw(c, c.Bounds(), list.Rows(40))
	l.Scroll(4)
	l.Draw(c, c.Bounds(), list.Rows(40))

	if got := c.OwnerAt(2, 0); got.Index != 4 {
		t.Errorf("after scrolling to 4 the top row is owned by %v", got)
	}
}
