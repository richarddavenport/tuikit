package comp

import (
	"github.com/charmbracelet/lipgloss"
)

// Input is one line being typed into.
//
// From three places that each drew their own. democtl's filter is a pane title
// with an underscore stuck on the end; azctl's is the same trick; and the
// swarmctl palette needs a real one — a caret you can move, text that scrolls
// when it outgrows the box, and a cursor you can see.
//
// The underscore version is worth naming, because it is what everyone writes
// first and it is wrong in a specific way: it cannot show WHERE you are. Every
// edit happens at the end, so ← and ← do nothing, and a typo six characters
// back means deleting six characters. That is not a missing feature so much as
// a different, smaller idea of what a text field is.
//
// # What it does not do
//
// It does not own the keystrokes. [app.EditAt] applies one key to a string and
// a caret; this draws the result. Splitting them is what lets a tool decide
// that enter means "accept", or "run", or "add another" — decisions no
// component can make — while nobody re-implements what backspace does.
type Input struct {
	Text string
	// Cursor is the caret, as a RUNE index into Text: 0 is before the first
	// character, len(runes) is after the last. Runes rather than bytes because
	// a caret halfway through a multi-byte character is not a position.
	Cursor int

	// Prompt sits before the text and is not part of it — "> ", "/", "filter ".
	Prompt string
	// Placeholder is shown instead of the text when there is none, and is not
	// itself editable. A prompt showing nothing at all cannot say what it wants.
	Placeholder string

	// Focused decides whether the block cursor is drawn. An unfocused input
	// showing a cursor is claiming keystrokes it will not receive.
	Focused bool

	PromptStyle, TextStyle, PlaceholderStyle *lipgloss.Style
	// CursorFG and CursorBG paint the one cell the caret is on. The default
	// palette makes these the background and the foreground, so the caret is
	// reverse video and reads on a light terminal without being told about it.
	CursorFG, CursorBG *lipgloss.Style
}

// Draw paints the input into one row of r and returns the row below.
func (in Input) Draw(c *Canvas, r Rect, id ID) int {
	if r.Empty() {
		return r.Y
	}
	x := r.X
	if in.Prompt != "" {
		x += c.Text(x, r.Y, in.Prompt, in.PromptStyle, id)
	}
	right := r.X + r.W - 1
	if x > right {
		return r.Y + 1
	}

	// An empty field shows its placeholder, focused or not.
	//
	// It used to hide it under a live caret, on the reasoning that a prompt
	// showing two things at once is showing two things at once. The palette
	// design settled it the other way and is right: the caret says where you
	// are typing and the hint says what to type, and an empty focused prompt
	// with neither is the least useful state a field can be in. Focused, the
	// hint follows the caret rather than replacing it.
	if in.Text == "" && in.Placeholder != "" {
		at := x
		if in.Focused {
			c.Set(x, r.Y, " ", in.CursorBG, id)
			at = x + 2
		}
		if at <= right {
			c.Text(at, r.Y, Truncate(in.Placeholder, right-at+1), in.PlaceholderStyle, id)
		}
		return r.Y + 1
	}

	runes := []rune(in.Text)
	cursor := clamp(in.Cursor, 0, len(runes))
	// One column is kept for the caret so that a cursor at the end of a full
	// line is still visible. Without it the caret falls off the edge exactly
	// when the text is longest, which is when you most need to see it.
	width := right - x + 1
	from := in.scroll(runes, cursor, width)
	// When text continues past the edge the last column belongs to the
	// ellipsis, so the window shifts once more rather than letting the caret
	// land there and be painted over. A caret you cannot see because the
	// "there is more" marker ate it is the worst of both.
	if in.trailing(runes, from, width) {
		for from < cursor && in.columns(runes, from, cursor) >= width-1 {
			from++
		}
	}

	for i := from; i < len(runes) && x <= right; i++ {
		style := in.TextStyle
		if in.Focused && i == cursor {
			style = in.CursorBG
		}
		x += c.Set(x, r.Y, string(runes[i]), style, id)
	}
	// A caret past the last character has no character to invert, so it gets a
	// painted blank of its own.
	if in.Focused && cursor >= len(runes) && x <= right {
		c.Set(x, r.Y, " ", in.CursorBG, id)
	}
	// Text still runs past the right edge: say so rather than cutting silently.
	if in.trailing(runes, from, width) {
		c.Set(right, r.Y, c.Chrome().Ellipsis, in.PlaceholderStyle, id)
	}
	return r.Y + 1
}

// scroll is the first rune to draw, chosen so the caret is always on screen.
//
// The window follows the caret rather than the text: a field scrolled to the
// end and then arrowed back to the start has to come with you, or the thing you
// are editing is somewhere you cannot see.
func (in Input) scroll(runes []rune, cursor, width int) int {
	if width <= 0 {
		return 0
	}
	// Walk back from the caret until the visible run fills the width.
	from, used := cursor, 0
	for from > 0 {
		w := Width(string(runes[from-1]))
		if used+w >= width {
			break
		}
		used += w
		from--
	}
	return from
}

// columns is how many columns runes[from:to] occupies.
func (in Input) columns(runes []rune, from, to int) int {
	used := 0
	for i := from; i < to && i < len(runes); i++ {
		used += Width(string(runes[i]))
	}
	return used
}

// trailing reports whether anything is left over past the right edge.
func (in Input) trailing(runes []rune, from, width int) bool {
	used := 0
	for i := from; i < len(runes); i++ {
		used += Width(string(runes[i]))
		if used > width {
			return true
		}
	}
	return false
}
