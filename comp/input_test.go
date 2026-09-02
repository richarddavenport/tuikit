package comp_test

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/harness"
)

const inputID = comp.Name("query")

var caret = lipgloss.NewStyle().Reverse(true)

func drawInput(t *testing.T, in comp.Input, w int) (*comp.Canvas, string) {
	t.Helper()
	in.CursorBG = &caret
	c := comp.NewCanvas(w, 1)
	in.Draw(c, comp.Rect{X: 0, Y: 0, W: w, H: 1}, comp.Region(inputID))
	return c, strings.TrimRight(harness.Strip(c.String()), " ")
}

func TestInputPromptAndText(t *testing.T) {
	_, got := drawInput(t, comp.Input{Prompt: "> ", Text: "env"}, 20)
	if want := "> env"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

// TestInputPlaceholderOnlyWhenIdle. A placeholder under a live caret is a
// prompt showing two things at once.
func TestInputPlaceholderOnlyWhenIdle(t *testing.T) {
	_, idle := drawInput(t, comp.Input{Prompt: "> ", Placeholder: "type to search"}, 30)
	if want := "> type to search"; idle != want {
		t.Errorf("unfocused got %q, want %q", idle, want)
	}
	_, live := drawInput(t, comp.Input{Prompt: "> ", Placeholder: "type to search", Focused: true}, 30)
	if strings.Contains(live, "type to search") {
		t.Errorf("focused input still shows the placeholder: %q", live)
	}
}

// TestInputCaretIsAPaintedCell, which is what an underscore on the end cannot
// be: it has to sit ON a character, in the middle of the text.
func TestInputCaretIsAPaintedCell(t *testing.T) {
	c, _ := drawInput(t, comp.Input{Text: "hello", Cursor: 1, Focused: true}, 20)
	cell, ok := c.CellAt(1, 0)
	if !ok || cell.Text != "e" {
		t.Fatalf("cell under the caret is %q", cell.Text)
	}
	if cell.Style != &caret {
		t.Error("the character under the caret is not painted")
	}
	if other, _ := c.CellAt(0, 0); other.Style == &caret {
		t.Error("a character that is not under the caret was painted")
	}
}

// TestInputCaretPastTheEnd has no character to invert, so it needs a blank of
// its own — otherwise the caret disappears exactly when you are typing.
func TestInputCaretPastTheEnd(t *testing.T) {
	c, _ := drawInput(t, comp.Input{Text: "ab", Cursor: 2, Focused: true}, 20)
	cell, ok := c.CellAt(2, 0)
	if !ok || cell.Text != " " || cell.Style != &caret {
		t.Errorf("cell after the text is %q painted=%v, want a painted blank", cell.Text, cell.Style == &caret)
	}
}

func TestInputUnfocusedDrawsNoCaret(t *testing.T) {
	c, _ := drawInput(t, comp.Input{Text: "ab", Cursor: 1}, 20)
	for x := 0; x < 4; x++ {
		if cell, _ := c.CellAt(x, 0); cell.Style == &caret {
			t.Errorf("an unfocused input painted a caret at %d", x)
		}
	}
}

// TestInputScrollsToKeepTheCaretVisible is the property an underscore-on-the-end
// field cannot have: the window follows the caret, not the text.
func TestInputScrollsToKeepTheCaretVisible(t *testing.T) {
	long := "abcdefghijklmnopqrstuvwxyz"

	_, atEnd := drawInput(t, comp.Input{Text: long, Cursor: 26, Focused: true}, 10)
	if !strings.Contains(atEnd, "z") {
		t.Errorf("caret at the end, but the end is not shown: %q", atEnd)
	}
	_, atStart := drawInput(t, comp.Input{Text: long, Cursor: 0, Focused: true}, 10)
	if !strings.HasPrefix(atStart, "a") {
		t.Errorf("caret at the start, but the view did not come back: %q", atStart)
	}
	_, middle := drawInput(t, comp.Input{Text: long, Cursor: 13, Focused: true}, 10)
	if !strings.Contains(middle, "n") {
		t.Errorf("caret in the middle is not on screen: %q", middle)
	}
}

// TestInputSaysWhenTextRunsOff. Cutting silently makes a field look like it
// holds less than it does.
func TestInputSaysWhenTextRunsOff(t *testing.T) {
	_, got := drawInput(t, comp.Input{Text: "abcdefghijklmnop", Cursor: 0, Focused: true}, 8)
	if !strings.HasSuffix(got, "…") {
		t.Errorf("got %q, want it to end in an ellipsis", got)
	}
}

func TestInputTooNarrow(t *testing.T) {
	for _, w := range []int{0, 1} {
		_, got := drawInput(t, comp.Input{Prompt: "> ", Text: "x"}, w)
		if strings.Contains(got, "x") {
			t.Errorf("width %d drew the text: %q", w, got)
		}
	}
}

func TestInputClampsAStrayCaret(t *testing.T) {
	if _, got := drawInput(t, comp.Input{Text: "ab", Cursor: 99, Focused: true}, 10); got != "ab" {
		t.Errorf("got %q, want \"ab\" and no panic", got)
	}
	if _, got := drawInput(t, comp.Input{Text: "ab", Cursor: -3, Focused: true}, 10); got != "ab" {
		t.Errorf("got %q", got)
	}
}
