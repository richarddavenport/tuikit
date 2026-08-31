// Package theme is what a tuikit interface is allowed to look like: the colour
// roles it may use, and the characters it may print.
//
// It exists as its own package for two reasons, both learned in swarmctl before
// this was a library. The first is that these were not written down anywhere —
// nine raw ANSI numbers in one file plus six more hardcoded across three others,
// with the SAME number meaning "focused border" in one place and "title" in
// another purely by coincidence. A role with a name can be changed once; a
// number cannot be changed at all.
//
// The second is that a design system has to be able to read them. Anything that
// draws the interface somewhere other than a terminal — a mockup, a palette, a
// component sheet, a captured frame turned into HTML — has to start from the
// same table the terminal starts from, or it is drawing a tool that does not
// exist.
package theme

import "github.com/charmbracelet/lipgloss"

// Palette is the closed set of colour roles an interface draws with.
//
// Naming is by ROLE, never by hue. "Accent" survives someone deciding the
// interface should be blue; "pink" does not. A raw index also says what a
// colour IS instead of what it is FOR, which is how one value came to mean both
// "title" and "focused border" without anyone choosing that they move together.
//
// A tool takes [Default] and overrides the fields it wants. It does not build a
// Palette from scratch, because the point of the set being closed is that
// nine decisions is the whole vocabulary.
type Palette struct {
	// Accent is the interface's own colour: titles, the selected row, the
	// focused panel's border. It marks WHERE YOU ARE, which is why the same
	// value carries all three.
	Accent lipgloss.Color
	// Muted is text that is present but not the point: hints, footers,
	// explanations under a value.
	Muted lipgloss.Color
	// Border is an unfocused panel's edge — a shade below Muted, so the box is
	// visible without competing with the text inside it.
	Border lipgloss.Color
	// Success is a finished action and a live log stream: the state you wanted.
	Success lipgloss.Color
	// Pending is a queued change, not yet applied — amber reads as "waiting on
	// you", distinct from both the running value and an error.
	Pending lipgloss.Color
	// Danger is a refusal, an error, or a guarded environment. It is never
	// decoration: something coloured Danger is something that stopped or will.
	Danger lipgloss.Color
	// Stderr marks a log line from stderr. Distinct from Danger on purpose —
	// most programs write ordinary progress to stderr, and colouring that as a
	// failure would make every run look broken.
	Stderr lipgloss.Color
	// SelectionFG and SelectionBG are the one place the interface paints a
	// background: the selected row of a table, where a foreground colour alone
	// cannot be seen against the surrounding rows.
	SelectionFG lipgloss.Color
	SelectionBG lipgloss.Color

	// Extra is for a role the nine do not cover.
	//
	// It exists because the alternative is worse: a tool that needs a tenth
	// meaning and cannot name it reaches for a literal instead, and the guard
	// that keeps colours honest is the first casualty. Adding one should still
	// feel like a decision — if it is a shade of an existing role rather than a
	// different meaning, it is not one.
	Extra []Role
}

// Default is the palette swarmctl arrived at, and the one a tuikit tool gets
// unless it says otherwise.
//
// Values are ANSI 256 palette indices, because that is what a terminal
// understands and what every terminal has agreed on; a truecolour hex would
// look right on this machine and wrong over ssh from another.
var Default = Palette{
	Accent:      "205",
	Muted:       "241",
	Border:      "240",
	Success:     "34",
	Pending:     "214",
	Danger:      "196",
	Stderr:      "203",
	SelectionFG: "229",
	SelectionBG: "57",
}

// Role is one entry in the palette, for anything rendering the palette itself.
type Role struct {
	Name  string
	Color lipgloss.Color
	Why   string
}

// Roles is the palette in the order it is worth reading: what the interface is,
// then what it says, then what it warns about. Extra roles come last, in the
// order the tool declared them.
//
// The Why text describes the ROLE, not the hue, so it survives a tool
// recolouring the palette — which is the whole reason roles are named.
func (p Palette) Roles() []Role {
	return append([]Role{
		{"Accent", p.Accent, "titles, the selected row, the focused panel's border — where you are"},
		{"Muted", p.Muted, "hints, footers, explanations — present but not the point"},
		{"Border", p.Border, "an unfocused panel's edge, a shade below Muted"},
		{"Success", p.Success, "a finished action, a live log stream"},
		{"Pending", p.Pending, "a queued change, not yet applied — waiting on you"},
		{"Danger", p.Danger, "a refusal, an error, a guarded environment"},
		{"Stderr", p.Stderr, "a log line from stderr — not a failure, just the other stream"},
		{"SelectionFG", p.SelectionFG, "text on the one painted background"},
		{"SelectionBG", p.SelectionBG, "the selected row of a table"},
	}, p.Extra...)
}

// GlyphSet is every non-ASCII character an interface may print, mapped to why
// it is safe to print it.
//
// The set is CLOSED, and guard.Glyphs holds it closed. A terminal font without
// a glyph draws a replacement box, which reads as a bug rather than as
// decoration — that is exactly what happened to a block-character edit cursor.
// Adding one means adding it here and deciding, deliberately, that it is common
// enough.
//
// Box-drawing, arrows, bullets and typographic punctuation are in every font
// shipped with a terminal. Block elements (U+2580–U+259F), geometric shapes
// beyond the plain bullet, emoji and Nerd Font private-use icons are not.
type GlyphSet map[rune]string

// DefaultGlyphs is swarmctl's allow-list, and the one a tuikit tool starts from.
var DefaultGlyphs = GlyphSet{
	'·': "middle dot — key separator in footers",
	'—': "em dash — clause separator in messages",
	'…': "ellipsis — truncation marker",
	'→': "rightwards arrow — pending edit",
	'↑': "up arrow — key hints",
	'↓': "down arrow — key hints",
	'‹': "single left angle quote — tab strip",
	'›': "single right angle quote — tab strip",
	'●': "black circle — badges and the dirty marker",
	'•': "bullet — masked secret",
	'✓': "check mark — success",
	'✗': "ballot X — failure",
	'┌': "box drawing", '─': "box drawing", '┐': "box drawing",
	'│': "box drawing", '└': "box drawing", '┘': "box drawing",
}

// With returns a copy of the set with additions.
//
// A copy, because a GlyphSet is a map and a map is a reference: a tool that
// added a glyph to DefaultGlyphs in place would be adding it to every other
// tool in the process, and to the guard that is supposed to catch it. Pairs are
// rune, reason, rune, reason; an odd count or a non-rune key panics, because
// both are typos rather than conditions to handle.
func (g GlyphSet) With(pairs ...any) GlyphSet {
	if len(pairs)%2 != 0 {
		panic("theme: GlyphSet.With wants rune, reason pairs")
	}
	out := make(GlyphSet, len(g)+len(pairs)/2)
	for r, why := range g {
		out[r] = why
	}
	for i := 0; i < len(pairs); i += 2 {
		r, ok := pairs[i].(rune)
		if !ok {
			panic("theme: GlyphSet.With wants a rune key")
		}
		why, ok := pairs[i+1].(string)
		if !ok {
			panic("theme: GlyphSet.With wants a string reason")
		}
		out[r] = why
	}
	return out
}

// SpinnerRange is the one thing on screen the allow-list does not cover.
//
// The spinner is drawn by bubbles, not by the tool, so its characters never
// appear in a string literal and the guard never sees them. Writing that down
// rather than leaving it as a hole: Braille patterns are in every terminal font
// — which is exactly why a spinner reaches for them instead of the block
// elements that would be a box on someone's terminal — so the dependency
// happens to be making the same decision this list makes.
//
// Anything the tool draws itself still has to be in the GlyphSet.
var SpinnerRange = [2]rune{0x2800, 0x28FF}

// Printable reports a rune the interface may put on screen.
func (g GlyphSet) Printable(r rune) bool {
	if r < 128 {
		return true
	}
	if r >= SpinnerRange[0] && r <= SpinnerRange[1] {
		return true
	}
	_, ok := g[r]
	return ok
}
