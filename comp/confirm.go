package comp

import "github.com/charmbracelet/lipgloss"

// Confirm is a question in a box over the frame.
//
// # Where this came from
//
// swarmctl's action.go (777 lines) and pgctl's actionview.go (278) both have
// one, and democtl's overlay is a third. What they share is the shape — a
// title, an explanation, and the keys — and a hard-won rule about width.
//
// The rule is pgctl's, and its comment records why: "the plan's description
// contains lines as long as a list of every table a widened selection adds, and
// an unbounded box drew itself off the side of the screen — which the
// screenshots caught before anyone else did." So a modal is BOUNDED TO ITS
// CONTAINER by construction. There is no way to ask for one that is too wide.
//
// The meaningful difference is the floor. pgctl clamps between 32 and 104
// columns; democtl only caps at 64 and has no minimum, so on a narrow terminal
// its modal shrinks until the question no longer reads. pgctl is right, and a
// floor costs nothing: a box wider than the terminal is clipped by the canvas
// anyway, which is a better failure than a box too narrow to read.
//
// # Typing the name
//
// swarmctl requires DESTRUCTIVE actions to be confirmed by typing the subject's
// name — the service, or env/service when the environment is guarded (action.go
// confirmPhrase, confirmed). That lives in comp.Form, as a Field with a Must:
// it is a text field with one extra rule, and putting it here would have been a
// second implementation of typed input. A confirm that needs it draws a Form in
// its body and asks Form.Complete before acting.
type Confirm struct {
	Title, Body string

	// Danger marks an action that destroys something, which changes how the
	// title is drawn.
	Danger bool

	Hints []Hint

	// Max, Min and Margin bound the box. Zero takes the defaults, which are
	// pgctl's numbers with democtl's margin.
	Max, Min, Margin int

	// Styles.
	Border, TitleStyle, DangerStyle, BodyStyle, HintStyle *lipgloss.Style
}

// Default bounds, from pgctl, whose modal met the longest content of the four.
const (
	confirmMax    = 64
	confirmMin    = 32
	confirmMargin = 8
)

// Draw centres the box on the canvas and returns the rect it took.
//
// Drawn last is on top. There is no compositing step, no re-measuring of the
// lines beneath and nothing to get wrong — the bug that once wiped 8 of 18
// framed rows in democtl is not one that can be written here.
func (cf Confirm) Draw(c *Canvas, id ID) Rect {
	maxW, minW, margin := or(cf.Max, confirmMax), or(cf.Min, confirmMin), or(cf.Margin, confirmMargin)
	bounds := c.Bounds()

	w := clamp(bounds.W-margin, minW, maxW)
	body := Wrap(cf.Body, w-4)

	// title, blank, body, blank, keys, and two of border.
	h := len(body) + 6
	r := Center(bounds, w, h)

	title := cf.TitleStyle
	if cf.Danger && cf.DangerStyle != nil {
		title = cf.DangerStyle
	}
	inner := Pane{Focused: true, Focus: cf.Border, Border: cf.Border}.Draw(c, r, id)
	if inner.Empty() {
		return r
	}

	c.Text(inner.X, inner.Y, Truncate(cf.Title, inner.W), title, id)
	for i, line := range body {
		c.Text(inner.X, inner.Y+2+i, line, cf.BodyStyle, id)
	}
	c.Text(inner.X, inner.Y+len(body)+3, Hints(cf.Hints...), cf.HintStyle, id)
	return r
}

func or(v, fallback int) int {
	if v == 0 {
		return fallback
	}
	return v
}
