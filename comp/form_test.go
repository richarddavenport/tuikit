package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

const formName Name = "form"

func sample() Form {
	return Form{
		Marker: "> ",
		Fields: []Field{
			{Label: "target", Kind: FieldChoice, Choices: []string{"staging", "production"}, Choice: 0},
			{Label: "tag", Kind: FieldText, Text: "v2.4.1"},
			{Label: "note", Kind: FieldText, Placeholder: "(none)"},
			{Label: "dry run", Kind: FieldToggle, On: true},
		},
	}
}

// Every field at once, so the operator can see what they have chosen rather
// than remembering it. A wizard asking one question per screen asks someone to
// hold the answers in their head exactly when they should be able to look.
func TestEveryFieldIsVisibleAtOnce(t *testing.T) {
	c := NewCanvas(60, 6)
	sample().Draw(c, c.Bounds(), formName)

	got := c.String()
	for _, want := range []string{"target", "staging", "tag", "v2.4.1", "note", "(none)", "dry run", "yes"} {
		if !strings.Contains(got, want) {
			t.Errorf("%q is missing:\n%s", want, got)
		}
	}
}

// The values line up, or a form reads as a list of unrelated sentences.
func TestTheValuesLineUp(t *testing.T) {
	c := NewCanvas(60, 6)
	sample().Draw(c, c.Bounds(), formName)

	lines := strings.Split(c.String(), "\n")
	first := at(lines[1], "v2.4.1")
	if second := at(lines[3], "yes"); first != second {
		t.Errorf("values start at %d and %d:\n%s", first, second, c.String())
	}
}

// An unanswered field looks unanswered rather than broken.
func TestAnEmptyFieldShowsItsPlaceholder(t *testing.T) {
	c := NewCanvas(60, 6)
	sample().Draw(c, c.Bounds(), formName)

	if !strings.Contains(c.String(), "(none)") {
		t.Errorf("an empty field shows nothing:\n%s", c.String())
	}
}

// The marker comes from the caller. pgctl uses ▸, which is not in tuikit's
// default glyph set, so a component hard-coding it would smuggle a character
// past guard.Glyphs and draw a replacement box on a font without it.
func TestTheMarkerComesFromTheCaller(t *testing.T) {
	f := sample()
	f.Focused, f.Cursor = true, 1
	c := NewCanvas(60, 6)
	f.Draw(c, c.Bounds(), formName)

	lines := strings.Split(c.String(), "\n")
	if !strings.HasPrefix(lines[1], "> tag") {
		t.Errorf("the cursor row is %q", lines[1])
	}
	if strings.HasPrefix(lines[0], ">") {
		t.Errorf("a row that is not the cursor is marked: %q", lines[0])
	}
	// The blank is the marker's width, or the labels jump as you move.
	if at(lines[0], "target") != at(lines[1], "tag") {
		t.Errorf("the labels move with the cursor:\n%q\n%q", lines[0], lines[1])
	}
}

// swarmctl requires a destructive action to be confirmed by typing the
// subject's name. The difference between a keystroke and a decision.
func TestAPhraseHasToMatchBeforeTheFormIsComplete(t *testing.T) {
	f := Form{Fields: []Field{{Label: "type the name", Must: "api_gateway"}}}

	if f.Complete() {
		t.Error("an unmatched phrase counts as complete")
	}
	f.Fields[0].Text = "api_gatewa"
	if f.Complete() {
		t.Error("half the phrase counts as complete")
	}
	f.Fields[0].Text = "api_gateway"
	if !f.Complete() {
		t.Error("the matching phrase does not complete the form")
	}
}

// A field with no phrase is not a field the form is waiting on. Whether an
// empty text field is acceptable is the tool's business — plenty are optional.
func TestOnlyPhrasesBlockCompletion(t *testing.T) {
	if !(Form{Fields: []Field{{Label: "note", Kind: FieldText}}}).Complete() {
		t.Error("an empty ordinary field blocks the form")
	}
}

// A phrase half typed is a sentence in progress, not an error to complain
// about — so it is marked rather than rejected.
func TestAHalfTypedPhraseIsMarkedNotRejected(t *testing.T) {
	f := Form{Fields: []Field{{Label: "name", Must: "api_gateway", Text: "api_gate"}}}
	c := NewCanvas(40, 2)
	f.Draw(c, c.Bounds(), formName)

	if !strings.Contains(c.String(), "api_gate") {
		t.Errorf("what was typed is not shown:\n%s", c.String())
	}
}

// A disabled field still shows its value: a choice you cannot change is one you
// may still need to read.
func TestADisabledFieldStillShowsItsValue(t *testing.T) {
	f := Form{Fields: []Field{{Label: "target", Kind: FieldText, Text: "production", Disabled: true}}}
	c := NewCanvas(40, 2)
	f.Draw(c, c.Bounds(), formName)

	if !strings.Contains(c.String(), "production") {
		t.Errorf("a disabled field hides its value:\n%s", c.String())
	}
}

func TestAFormStaysInsideItsRect(t *testing.T) {
	c := NewCanvas(60, 8)
	sample().Draw(c, Rect{X: 0, Y: 0, W: 20, H: 2}, formName)

	lines := strings.Split(c.String(), "\n")
	for i, line := range lines {
		if Width(line) > 20 {
			t.Errorf("line %d is %d columns", i+1, Width(line))
		}
		if i >= 2 && strings.TrimSpace(line) != "" {
			t.Errorf("drew row %d into a two-row rect: %q", i+1, line)
		}
	}
}

// TestAFormCanBeTypedInto: the focused text field is drawn by Input, so it has
// a caret and it scrolls. A form without one can only be corrected by deleting
// back to the mistake.
func TestAFormCanBeTypedInto(t *testing.T) {
	block := lipgloss.NewStyle().Reverse(true)
	f := Form{
		Fields: []Field{
			{Label: "env", Kind: FieldText, Text: "production"},
			{Label: "note", Kind: FieldText, Text: "second"},
		},
		Cursor: 0, Focused: true, Marker: "> ", Blank: "  ",
		Caret: 4, CursorBG: &block,
	}
	c := NewCanvas(40, 3)
	f.Draw(c, Rect{X: 0, Y: 0, W: 40, H: 3}, "form")

	// The caret sits ON the fifth character of "production", not after it.
	found := false
	for x := 0; x < 40; x++ {
		if cell, _ := c.CellAt(x, 0); cell.Style == &block {
			if cell.Text != "u" {
				t.Errorf("the caret is on %q, want the rune at index 4", cell.Text)
			}
			found = true
		}
	}
	if !found {
		t.Error("no caret was painted on the focused field")
	}
	// The unfocused row has none.
	for x := 0; x < 40; x++ {
		if cell, _ := c.CellAt(x, 1); cell.Style == &block {
			t.Error("a caret was painted on a field that is not being edited")
		}
	}
}

// TestAFormWithoutACursorIsUnchanged, so every existing form draws as it did.
func TestAFormWithoutACursorIsUnchanged(t *testing.T) {
	fields := []Field{{Label: "env", Kind: FieldText, Text: "production"}}
	plain := Form{Fields: fields, Focused: true, Marker: "> ", Blank: "  "}
	withCaret := plain
	withCaret.Caret = 4 // but no CursorBG

	a, b := NewCanvas(40, 2), NewCanvas(40, 2)
	plain.Draw(a, Rect{X: 0, Y: 0, W: 40, H: 2}, "form")
	withCaret.Draw(b, Rect{X: 0, Y: 0, W: 40, H: 2}, "form")
	if a.String() != b.String() {
		t.Error("a form with no cursor style drew differently")
	}
}

// A choice or a toggle is never "typed into", whatever the caret says.
func TestOnlyTextFieldsGetACaret(t *testing.T) {
	block := lipgloss.NewStyle().Reverse(true)
	f := Form{
		Fields:  []Field{{Label: "how", Kind: FieldChoice, Choices: []string{"a", "b"}}},
		Focused: true, Marker: "> ", Blank: "  ", Caret: 0, CursorBG: &block,
	}
	c := NewCanvas(40, 2)
	f.Draw(c, Rect{X: 0, Y: 0, W: 40, H: 2}, "form")
	for x := 0; x < 40; x++ {
		if cell, _ := c.CellAt(x, 0); cell.Style == &block {
			t.Error("a choice field was drawn with a caret")
		}
	}
}
