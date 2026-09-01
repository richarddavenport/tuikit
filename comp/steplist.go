package comp

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// StepState is where one step of a run has got to.
type StepState int

// The five states a step can be in. Waiting is the zero value, because a step
// list is usually built before anything has started.
const (
	StepWaiting StepState = iota
	StepRunning
	StepSkipped
	StepDone
	StepFailed
)

// StepList is a run in progress: what will happen, what has, and what it cost.
//
// # Where this came from
//
// swarmctl's renderdeploy.go (809 lines), pgctl's actionrun.go (370) and
// azctl's runner.go (262) each have one, and democtl's is a fourth. All four
// agree on the shape: a glyph per state, the step's name, and something
// trailing that says what it cost or why it did not run.
//
// The meaningful difference is WHERE that trailing thing goes. azctl puts it
// inline after the label; democtl right-aligns it against the pane's edge. It
// turns out they are not the same thing badly agreed on — they are two things.
// A duration is a number you scan down a column, so it goes right. A skip
// reason or a dry-run command is prose about that step, so it sits beside it.
// Both, and no flag.
//
// democtl's comment on the right-aligned column is worth keeping too: it was
// measured against the TERMINAL rather than the box, which put the durations
// two columns past the edge where they clipped to "400m…". Visible in a
// captured frame and in nothing else.
//
// The glyphs are the tool's, not this component's. A closed glyph set is
// per-tool and a component that hard-codes ✓ has taken a decision belonging to
// whoever owns the set — and would smuggle a character past guard.Glyphs.
type StepList struct {
	Steps []Step

	// Look is the glyph and style for each state, indexed by StepState.
	Look [5]StepLook

	// Status is the line under the list — running…, done, stopped — with its
	// own style, and hints that follow it.
	Status      string
	StatusStyle *lipgloss.Style
	Hints       []Hint

	// Muted draws the trailing detail and the durations.
	Muted *lipgloss.Style
}

// Step is one thing that will happen, or has.
type Step struct {
	Label string
	State StepState

	// Detail is prose about this step — a skip reason, a dry-run command — set
	// beside the label because it is about the step rather than about the run.
	Detail string
	// Took is preformatted, and right-aligned into a column. A string rather
	// than a duration because "400ms" and "0.4s" are a tool's decision about
	// its own readers, and there is no version of it that suits all four.
	Took string
	// Note is a line of its own beneath the step, for the explanation a
	// failure owes.
	Note string
}

// StepLook is how one state is drawn.
type StepLook struct {
	Glyph string
	Style *lipgloss.Style
	// LabelStyle draws the step's name. Nil leaves it unstyled, which is what
	// a step that has already happened wants — it is the badge that carries
	// the state, not the name.
	LabelStyle *lipgloss.Style
}

// Rows lays the steps out as List rows, for a run longer than its pane.
//
// A StepList draws every step and does not scroll; a List scrolls and knows
// nothing about steps. azctl's playbooks are longer than its pane, so it needs
// both, and it was writing this adapter by hand — twenty lines of mechanism
// living in a tool because the component could not do the job.
//
// It is the answer comp.Table's doc comment already gives for rows: a Table
// lays out, a List selects and scrolls, and rows that do both are two
// components composed. So this lays out, and `l.Draw(c, r, steps.Rows(r.W))`
// scrolls.
//
// A note becomes a row of its own and so carries its own index, where Draw can
// give it the step's. Nothing sets a note today; a tool that does, and wants a
// click on it to select the step it belongs to, has a real requirement — and it
// is not this function's to guess at.
func (s StepList) Rows(w int) []Row {
	out := make([]Row, 0, len(s.Steps))
	for i := range s.Steps {
		out = append(out, s.rowsFor(i, w)...)
	}
	return out
}

// rowsFor is one step's lines: the step, and its note if it has one.
func (s StepList) rowsFor(i, w int) []Row {
	step := s.Steps[i]
	look := s.Look[step.State]

	spans := []Segment{
		{Text: "  "},
		{Text: look.Glyph, Style: look.Style},
		{Text: " "},
		{Text: step.Label, Style: look.LabelStyle},
	}
	if step.Detail != "" {
		spans = append(spans, Segment{Text: "  " + step.Detail, Style: s.Muted})
	}
	if step.Took != "" {
		// Right-aligned against the RECT it is given, not the terminal's
		// width. democtl measured against the terminal, which put the durations
		// two columns past the edge where they clipped to "400m…" — visible in
		// a captured frame and in nothing else.
		took := step.Took + " "
		if pad := w - width(spans) - Width(took); pad > 0 {
			spans = append(spans, Segment{Text: strings.Repeat(" ", pad)})
		}
		spans = append(spans, Segment{Text: took, Style: s.Muted})
	}

	out := []Row{{Spans: spans}}
	if step.Note != "" {
		out = append(out, Row{Spans: []Segment{
			{Text: "      " + step.Note, Style: look.Style},
		}})
	}
	return out
}

// Draw renders the list into r, for a run that fits. One that does not is a
// List drawing Rows.
func (s StepList) Draw(c *Canvas, r Rect, name Name) {
	c = c.Clip(r)
	y := r.Y
	for i := range s.Steps {
		// A note shares its step's ID here, because Draw can afford to know
		// which rows belong together and a List cannot.
		id := Region(name).At(i)
		for _, row := range s.rowsFor(i, r.W) {
			if y > r.Bottom() {
				return
			}
			x := r.X
			for _, span := range row.Spans {
				x += c.Text(x, y, span.Text, span.Style, id)
			}
			y++
		}
	}

	if s.Status == "" || y+1 > r.Bottom() {
		return
	}
	y++
	x := r.X + c.Text(r.X, y, "  ", nil, Region(name))
	x += c.Text(x, y, s.Status, s.StatusStyle, Region(name))
	if len(s.Hints) > 0 {
		c.Text(x, y, "  "+Hints(s.Hints...), s.Muted, Region(name))
	}
}

// Look builds a StepLook, for the common case of a glyph and one style.
func Look(glyph string, style *lipgloss.Style) StepLook {
	return StepLook{Glyph: glyph, Style: style}
}
