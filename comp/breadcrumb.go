package comp

import "github.com/charmbracelet/lipgloss"

// Breadcrumb is how you got here, and how to go back.
//
// [app.Stack] has kept a Path since it was written and nothing ever drew it —
// a public, tested accessor that no caller called, which is the clearest sign
// available that a component was missing rather than undecided (decision 31).
//
// It takes strings rather than app.Entry because app imports comp, so comp
// cannot import app. That turns out to be the right shape anyway: what a crumb
// is CALLED is the tool's business — "democtl", "api_gateway", "Logs" — and a
// component that took the stack's own type would be claiming otherwise.
type Breadcrumb struct {
	// Crumbs is root first, current last.
	Crumbs []string

	// Name owns the whole strip; Item owns crumb i, so a click resolves to a
	// DEPTH by name rather than by counting columns. Leave Item empty and the
	// crumbs are not separately clickable.
	Name, Item Name

	Style, Current, Separator *lipgloss.Style
}

// Draw paints the trail into one row of r and returns the row below.
//
// Too narrow, it elides from the LEFT, keeping the root and the current screen:
// where you are and how far in are the two things a breadcrumb is for, and the
// middle is what you can afford to lose.
func (b Breadcrumb) Draw(c *Canvas, r Rect) int {
	if r.Empty() || len(b.Crumbs) == 0 {
		return r.Y
	}
	c = c.Clip(r)
	// Chrome.Separator already carries its own spaces.
	sep := c.Chrome().Separator
	from := b.fit(r.W, Width(sep))

	x := r.X
	if from > 0 {
		x += c.Text(x, r.Y, c.Chrome().Ellipsis+sep, b.Style, Region(b.Name))
	}
	for i := from; i < len(b.Crumbs); i++ {
		if i > from {
			x += c.Text(x, r.Y, sep, b.Separator, Region(b.Name))
		}
		style := b.Style
		if i == len(b.Crumbs)-1 {
			style = b.Current
		}
		x += c.Text(x, r.Y, b.Crumbs[i], style, b.crumb(i))
	}
	return r.Y + 1
}

// crumb is who owns crumb i: itself if the caller named the items, otherwise
// the strip, so the whole trail is still one region a click can land on.
func (b Breadcrumb) crumb(i int) ID {
	if b.Item == "" {
		return Region(b.Name)
	}
	return Region(b.Item).At(i)
}

// fit is the first crumb to draw so the trail ends inside the width.
//
// The LAST crumb always survives, even alone and even truncated: a breadcrumb
// that dropped the screen you are on to keep the ones you are not would be
// answering the wrong question.
func (b Breadcrumb) fit(width, sep int) int {
	if len(b.Crumbs) == 0 {
		return 0
	}
	used := Width(b.Crumbs[len(b.Crumbs)-1])
	for i := len(b.Crumbs) - 2; i >= 0; i-- {
		next := used + sep + Width(b.Crumbs[i])
		// Room for the ellipsis that says something was dropped, unless this
		// is the root and nothing was.
		if i > 0 {
			next += sep + 1
		}
		if next > width {
			return i + 1
		}
		used += sep + Width(b.Crumbs[i])
	}
	return 0
}
