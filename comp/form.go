package comp

import (
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// Form is fields the reader fills in before something happens.
//
// # Where this came from
//
// The database tool's viewForm and the deploy tool's applyedits.go. Both
// render EVERY FIELD AT ONCE, and the database tool says why: "so the operator
// can see what they have chosen rather than remembering it". A wizard that
// asks one question per screen is asking someone to hold the answers in their
// head while deciding whether to go ahead, which is exactly when they should
// be able to look.
//
// # The phrase
//
// The deploy tool requires destructive actions to be confirmed by TYPING the
// subject's name — env/service where the environment is guarded (action.go
// confirmPhrase, confirmed). This is where it lives, because it is a text
// field with one extra rule: the value has to match. comp.Confirm's doc
// comment sends you here.
//
// It is not a nag. It is the difference between a keystroke and a decision,
// and it belongs on the removal of a thing whose name you should be able to
// type if you are sure you mean that one.
//
// # The marker
//
// The glyph on the focused row comes from the caller, like a StepList's
// badges. The database tool uses ▸, which is not in tuikit's default glyph set
// — so a component that hard-coded it would smuggle a character past
// guard.Glyphs and produce a replacement box on a font without it.
type Form struct {
	Fields []Field
	// Cursor is the field being edited.
	Cursor int
	// Focused says the form has the keys. An unfocused form still shows every
	// value, because seeing what you have chosen is the point.
	Focused bool

	// Marker is drawn against the focused field; Blank is drawn against the
	// rest, and should be the same width or the labels jump as you move.
	Marker, Blank string

	// LabelWidth aligns the values into a column. Zero measures the labels.
	LabelWidth int

	// Caret is where you are typing in the focused text field, as a rune
	// index, and CursorFG/CursorBG paint that one cell.
	//
	// A form without a caret is a form you can only correct by deleting back
	// to the mistake — which is fine for a phrase you type once and wrong for
	// a value you are editing. Set them and the focused text field is drawn by
	// [Input], so it scrolls when the value outgrows the column and shows
	// where you are; leave them and it is drawn flat, as it always was.
	Caret              int
	CursorFG, CursorBG *lipgloss.Style

	Label, FocusLabel, Value, Muted, Danger *lipgloss.Style
}

// FieldKind is what sort of answer a field takes.
type FieldKind int

// The kinds, from the database tool's form: free text, one of a list, and a
// flag.
const (
	FieldText FieldKind = iota
	FieldChoice
	FieldToggle
)

// Field is one question.
type Field struct {
	Label string
	Kind  FieldKind

	// Text is the value of a FieldText.
	Text string
	// Placeholder is drawn when Text is empty, so an unanswered field looks
	// unanswered rather than broken.
	Placeholder string

	// Choices and Choice are a FieldChoice.
	Choices []string
	Choice  int

	// On is a FieldToggle.
	On bool

	// Must is a phrase Text has to match before the form is Complete — the deploy
	// tool's type-the-name-to-confirm. Empty means no such requirement.
	Must string

	// Disabled greys a field out. It still shows its value: a choice you
	// cannot change is one you may still need to read.
	Disabled bool

	// Secret hides a text field's value behind [Mask], including while it is
	// being typed into. See Fact.Secret for why the component decides rather
	// than the caller.
	//
	// The caret is not drawn on a masked field. A caret moving over eight
	// identical bullets says nothing, and one that stops early says how long
	// the secret is.
	Secret bool
}

// Complete reports whether every field that has to be answered has been.
//
// Only phrases are checked. Whether an empty text field is acceptable is the
// tool's business — plenty are optional — but a phrase exists precisely to be
// compared, so a form that did not check it would be decoration.
func (f Form) Complete() bool {
	for _, field := range f.Fields {
		if field.Must != "" && field.Text != field.Must {
			return false
		}
	}
	return true
}

// Draw renders every field into r, one per line.
func (f Form) Draw(c *Canvas, r Rect, name Name) {
	c = c.Clip(r)
	width := f.LabelWidth
	if width == 0 {
		for _, field := range f.Fields {
			width = max(width, Width(field.Label))
		}
	}
	marker, blank := f.Marker, f.Blank
	if blank == "" {
		blank = strings.Repeat(" ", Width(marker))
	}

	for i, field := range f.Fields {
		y := r.Y + i
		if y > r.Bottom() {
			return
		}
		id := Region(name).At(i)

		mark, label := blank, f.Label
		if i == f.Cursor && f.Focused {
			mark, label = marker, f.FocusLabel
		}
		if field.Disabled {
			label = f.Muted
		}

		x := r.X + c.Text(r.X, y, mark, f.FocusLabel, id)
		x += c.Text(x, y, fit(field.Label, width, false)+"  ", label, id)
		// Only the focused, enabled text field is being typed into, and only
		// when the caller supplied a cursor to paint it with.
		editing := i == f.Cursor && f.Focused && !field.Disabled && !field.Secret &&
			field.Kind == FieldText && f.CursorBG != nil
		f.value(c, x, y, r.Right(), field, id, editing)
	}
}

func (f Form) value(c *Canvas, x, y, right int, field Field, id ID, editing bool) {
	switch field.Kind {
	case FieldToggle:
		text, style := "no", f.Muted
		if field.On {
			text, style = "yes", f.Value
		}
		c.Text(x, y, text, style, id)

	case FieldChoice:
		for i, choice := range field.Choices {
			if x > right {
				return
			}
			style := f.Muted
			if i == field.Choice {
				style = f.Value
			}
			x += c.Text(x, y, choice+"  ", style, id)
		}

	default:
		if field.Secret && field.Text != "" {
			c.Text(x, y, Mask(), f.Value, id)
			return
		}
		// The field being typed into is an Input, so the caret, the scrolling
		// and the ellipsis are one implementation rather than two.
		if editing {
			Input{
				Text: field.Text, Cursor: f.Caret, Placeholder: field.Placeholder,
				Focused:          true,
				TextStyle:        f.Value,
				PlaceholderStyle: f.Muted,
				CursorFG:         f.CursorFG, CursorBG: f.CursorBG,
			}.Draw(c, Rect{X: x, Y: y, W: right - x + 1, H: 1}, id)
			return
		}
		if field.Text == "" {
			c.Text(x, y, field.Placeholder, f.Muted, id)
			return
		}
		style := f.Value
		// A phrase that does not match yet is not an error — it is a sentence
		// half typed. It is marked, not complained about.
		if field.Must != "" && field.Text != field.Must {
			style = f.Danger
		}
		c.Text(x, y, Truncate(field.Text, right-x+1), style, id)
	}
}
