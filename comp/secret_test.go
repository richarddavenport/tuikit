package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/theme"
)

const theSecret = "ghp_R3alT0kenValue"

// TestAMaskedFactNeverReachesTheCanvas is the property the whole thing exists
// for, and the one a caller cannot be trusted with: a tool that masked by
// choosing which string to pass would put the real one in the frame every time
// it got the condition backwards, and nothing could tell.
func TestAMaskedFactNeverReachesTheCanvas(t *testing.T) {
	c := NewCanvas(60, 4)
	Detail{Blocks: []Block{{Facts: []Fact{
		{Label: "token", Value: theSecret, Secret: true},
		{Label: "owner", Value: "platform"},
	}}}}.Draw(c, Rect{X: 0, Y: 0, W: 60, H: 4}, Region("detail"))

	frame := c.String()
	if strings.Contains(frame, theSecret) {
		t.Fatalf("the secret is in the frame:\n%s", frame)
	}
	if strings.Contains(frame, "R3al") {
		t.Errorf("part of the secret is in the frame:\n%s", frame)
	}
	if !strings.Contains(frame, Mask()) {
		t.Errorf("nothing was drawn in its place:\n%s", frame)
	}
	// The fact beside it is untouched — masking is per fact, not per pane.
	if !strings.Contains(frame, "platform") {
		t.Errorf("an ordinary fact was masked too:\n%s", frame)
	}
}

// Revealed, it is the value and nothing else.
func TestARevealedFactIsTheValue(t *testing.T) {
	c := NewCanvas(60, 3)
	Detail{Blocks: []Block{{Facts: []Fact{
		{Label: "token", Value: theSecret},
	}}}}.Draw(c, Rect{X: 0, Y: 0, W: 60, H: 3}, Region("detail"))

	if !strings.Contains(c.String(), theSecret) {
		t.Error("a fact that is not secret was hidden anyway")
	}
}

// TestTheMaskDoesNotLeakTheLength. Masking to the value's own length tells a
// reader the password is six characters, which is most of what a guesser wants.
func TestTheMaskDoesNotLeakTheLength(t *testing.T) {
	short := Fact{Label: "a", Value: "x", Secret: true}
	long := Fact{Label: "a", Value: strings.Repeat("y", 200), Secret: true}
	if short.shown() != long.shown() {
		t.Errorf("a one-character secret masks to %q and a long one to %q",
			short.shown(), long.shown())
	}
	if Width(Mask()) != maskWidth {
		t.Errorf("the mask is %d columns, want %d", Width(Mask()), maskWidth)
	}
}

// The bullet is in the default glyph set, so a masked value is not a row of
// replacement boxes on a font without it.
func TestTheMaskIsInTheGlyphSet(t *testing.T) {
	for _, r := range Mask() {
		if _, ok := theme.DefaultGlyphs[r]; !ok {
			t.Errorf("the mask uses %q (%U), which is not in DefaultGlyphs", r, r)
		}
	}
}

// A secret form field is masked while it is being typed into, and draws no
// caret: a caret moving over eight identical bullets says nothing, and one
// that stops early says how long the secret is.
func TestASecretFieldIsMaskedWhileTyping(t *testing.T) {
	block := lipgloss.NewStyle().Reverse(true)
	c := NewCanvas(50, 2)
	Form{
		Fields:  []Field{{Label: "token", Kind: FieldText, Text: theSecret, Secret: true}},
		Focused: true, Marker: "> ", Blank: "  ", Caret: 4, CursorBG: &block,
	}.Draw(c, Rect{X: 0, Y: 0, W: 50, H: 2}, "form")

	frame := c.String()
	if strings.Contains(frame, theSecret) {
		t.Fatalf("the secret is in the frame while being typed:\n%s", frame)
	}
	if !strings.Contains(frame, Mask()) {
		t.Errorf("nothing was drawn in its place:\n%s", frame)
	}
	for x := 0; x < 50; x++ {
		if cell, _ := c.CellAt(x, 0); cell.Style == &block {
			t.Error("a caret was drawn on a masked field")
		}
	}
}

// An empty secret field still shows its placeholder — an unanswered field
// should look unanswered, not like a value that is hidden.
func TestAnEmptySecretFieldShowsItsPlaceholder(t *testing.T) {
	c := NewCanvas(50, 2)
	Form{
		Fields:  []Field{{Label: "token", Kind: FieldText, Placeholder: "required", Secret: true}},
		Focused: true, Marker: "> ", Blank: "  ",
	}.Draw(c, Rect{X: 0, Y: 0, W: 50, H: 2}, "form")

	if got := c.String(); !strings.Contains(got, "required") {
		t.Errorf("an empty secret field hid its placeholder:\n%s", got)
	}
}
