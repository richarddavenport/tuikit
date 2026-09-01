package comp

import "github.com/charmbracelet/lipgloss"

// Pane is a bordered box with a title, and the rect inside it.
//
// # Where this came from
//
// swarmctl has TWO of these — drawBox (dashboard.go:1013) and drawBoxRaw
// (pane.go:454) — and they disagree with each other. One truncates a title too
// long for its frame and styles it; the other drops the title entirely and
// takes it pre-styled. Neither is wrong; they were written months apart for the
// same job. A tool duplicating a function inside itself is the clearest
// argument for extraction there is, because there is no second tool's
// requirements to blame it on.
//
// democtl draws its own, and puts the title on a ROW INSIDE the box where
// swarmctl puts it in the top edge. That difference is real rather than
// accidental: swarmctl's panes are small and numerous, so a row costs it
// something, while democtl's title row carries the focus highlight across the
// full width of the pane. So it is a field rather than a decision taken here.
//
// A title too long to fit is truncated with an ellipsis, following drawBox.
// drawBoxRaw's answer — drop the title — loses the name of the pane you are
// looking at exactly when the window is too small to work out what it is.
type Pane struct {
	Title string
	// Where the title goes. The zero value puts it in the top edge, which is
	// the more common arrangement and the cheaper one.
	TitleAt TitlePlacement

	Focused bool

	// Border is the frame; Focus replaces it when the pane has focus.
	Border, Focus *lipgloss.Style
	// TitleStyle draws the title. FocusTitle replaces it when focused, and
	// falls back to TitleStyle when nil.
	TitleStyle, FocusTitle *lipgloss.Style
}

// TitlePlacement is where a pane's title sits.
type TitlePlacement int

const (
	// TitleInEdge writes the title into the top border: ┌─ name ───┐.
	TitleInEdge TitlePlacement = iota
	// TitleOnRow gives the title the first row inside the box, painted the
	// full width so a focus highlight reads as the pane's rather than the
	// text's.
	TitleOnRow
)

// Draw renders the pane into r and returns the rect inside it — after the
// title row, if the title has one.
//
// The inside is blanked, so a pane drawn over something else covers it. That is
// what makes a modal or a menu a draw rather than a composite.
func (p Pane) Draw(c *Canvas, r Rect, id ID) Rect {
	c = c.Clip(r)
	edge := p.Border
	if p.Focused && p.Focus != nil {
		edge = p.Focus
	}
	title := p.TitleStyle
	if p.Focused && p.FocusTitle != nil {
		title = p.FocusTitle
	}

	if r.W < 2 || r.H < 2 {
		return Rect{}
	}
	c.Box(r, edge, id)
	inner := r.Inset(c.Chrome().Inset)
	c.Fill(inner, " ", nil, ID{})

	if p.Title == "" {
		return inner
	}
	switch p.TitleAt {
	case TitleOnRow:
		row := Rect{X: inner.X, Y: inner.Y, W: inner.W, H: 1}
		c.Fill(row, " ", title, id)
		c.Text(row.X, row.Y, " "+truncate(p.Title, row.W-1, c.Chrome().Ellipsis), title, id)
		return Rect{X: inner.X, Y: inner.Y + 1, W: inner.W, H: inner.H - 1}
	default:
		// Two columns in, so the title never sits against a corner.
		c.Text(r.X+2, r.Y, truncate(" "+p.Title+" ", max(0, inner.W-2), c.Chrome().Ellipsis), title, id)
		return inner
	}
}
