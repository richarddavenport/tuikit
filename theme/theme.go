// Package theme is what a tuikit interface is allowed to look like: the colour
// roles it may use, and the characters it may print.
//
// It exists as its own package for two reasons, both learned in the deploy
// tool before this was a library. The first is that these were not written
// down anywhere — nine raw ANSI numbers in one file plus six more hardcoded
// across three others, with the SAME number meaning "focused border" in one
// place and "title" in another purely by coincidence. A role with a name can
// be changed once; a number cannot be changed at all.
//
// The second is that a design system has to be able to read them. Anything
// that draws the interface somewhere other than a terminal — a mockup, a
// palette, a component sheet, a captured frame turned into HTML — has to start
// from the same table the terminal starts from, or it is drawing a tool that
// does not exist.
package theme

import "github.com/charmbracelet/lipgloss"

// Palette is the closed set of colour roles an interface draws with.
//
// Naming is by ROLE, never by hue. "Accent" survives someone deciding the
// interface should be blue; "pink" does not. A raw index also says what a
// colour IS instead of what it is FOR, which is how one value came to mean
// both "title" and "focused border" without anyone choosing that they move
// together.
//
// A tool takes [Default] and overrides the fields it wants. It does not build
// a Palette from scratch, because the point of the set being closed is that
// nine decisions is the whole vocabulary.
//
// # Why the colours are an interface
//
// A role holds a lipgloss.TerminalColor rather than a lipgloss.Color, so a
// tool may supply an AdaptiveColor and have its interface read on a light
// terminal as well as a dark one. That is the cloud tool's requirement, found
// by migrating it: its palette is light/dark pairs, and a Palette that could
// not hold them would have forced it to choose between tuikit's vocabulary and
// working in daylight.
//
// Depth and adaptation are different questions. Decision 13 chose ANSI 256
// over truecolour because ssh decides the profile; that is about how many
// colours there are. Which of them to use on a pale background is a separate
// decision, and one a design system has no business taking for a tool.
type Palette struct {
	// Accent is the interface's own colour: titles, the selected row, the
	// focused panel's border. It marks WHERE YOU ARE, which is why the same
	// value carries all three.
	Accent lipgloss.TerminalColor
	// Muted is text that is present but not the point: hints, footers,
	// explanations under a value.
	Muted lipgloss.TerminalColor
	// Border is an unfocused panel's edge — a shade below Muted, so the box is
	// visible without competing with the text inside it.
	Border lipgloss.TerminalColor
	// Success is a finished action and a live log stream: the state you wanted.
	Success lipgloss.TerminalColor
	// Pending is a queued change, not yet applied — amber reads as "waiting on
	// you", distinct from both the running value and an error.
	Pending lipgloss.TerminalColor
	// Danger is a refusal, an error, or a guarded environment. It is never
	// decoration: something coloured Danger is something that stopped or will.
	Danger lipgloss.TerminalColor
	// Stderr marks a log line from stderr. Distinct from Danger on purpose —
	// most programs write ordinary progress to stderr, and colouring that as a
	// failure would make every run look broken.
	Stderr lipgloss.TerminalColor
	// SelectionFG and SelectionBG are the one place the interface paints a
	// background: the selected row of a table, where a foreground colour alone
	// cannot be seen against the surrounding rows.
	SelectionFG lipgloss.TerminalColor
	SelectionBG lipgloss.TerminalColor

	// Extra is for a role the nine do not cover.
	//
	// It exists because the alternative is worse: a tool that needs a tenth
	// meaning and cannot name it reaches for a literal instead, and the guard
	// that keeps colours honest is the first casualty. Adding one should still
	// feel like a decision — if it is a shade of an existing role rather than a
	// different meaning, it is not one.
	Extra []Role
}

// Default is the palette the deploy tool arrived at, and the one a tuikit tool
// gets unless it says otherwise.
//
// Values are the FIRST SIXTEEN ANSI indices, and the sixteen are the whole
// point: they are the only colours a terminal lets its user redefine.
//
// Everything from 16 up is a fixed formula — a 6x6x6 cube and a grey ramp —
// identical in every terminal and untouched by every theme. A palette built
// from those indices looks the same under gruvbox, tokyo-night and solarized,
// which is another way of saying it ignores what the reader chose. This
// palette used to be 205/241/240; it was themeable in the sense that a Go
// programmer could edit it.
//
// The sixteen are already semantic, which is what makes this a mapping rather
// than a guess. Terminal themes agree that 0 is the background and 7 the
// foreground; 8 is the dimmed grey comments are drawn in; 15 is the brightest
// text. Omarchy's templates say so literally — `palette = 0={{ background }}`,
// `palette = 8={{ muted }}` — and every other theme system does the same thing
// under other names.
//
// So a tuikit tool is themed by whatever themed the terminal, with no config
// format, no loader, and nothing to reload: an index is resolved by the
// terminal at paint time, so changing the theme retints the next frame.
//
// # The selection is reverse video, deliberately
//
// SelectionFG is the background and SelectionBG the foreground, which inverts
// correctly on a light theme BY CONSTRUCTION rather than by detecting one. The
// deploy tool's design notes reached the same place independently: it is "the
// one treatment that reads identically in both profiles".
//
// # What this costs
//
// One grey. Muted and Border are the same index, where they used to be 241 and
// 240 — one step apart on the ramp and near-indistinguishable anyway.
//
// And the accent is the terminal's magenta rather than the theme's own accent
// colour, because ANSI has no accent slot. A tool that wants the real one
// reads it from wherever its desktop keeps it and overrides the role; that is
// what Extra and a plain assignment are for.
var Default = Palette{
	Accent:      lipgloss.Color("13"), // bright magenta
	Muted:       lipgloss.Color("8"),  // the dimmed grey — literally named muted
	Border:      lipgloss.Color("8"),  // the same grey; see above
	Success:     lipgloss.Color("2"),  // green
	Pending:     lipgloss.Color("3"),  // yellow
	Danger:      lipgloss.Color("1"),  // red
	Stderr:      lipgloss.Color("9"),  // bright red
	SelectionFG: lipgloss.Color("0"),  // the background
	SelectionBG: lipgloss.Color("7"),  // the foreground
}

// Role is one entry in the palette, for anything rendering the palette itself.
type Role struct {
	Name  string
	Color lipgloss.TerminalColor
	Why   string
}

// Roles is the palette in the order it is worth reading: what the interface
// is, then what it says, then what it warns about. Extra roles come last, in
// the order the tool declared them.
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
// Adding one means adding it here and deciding, deliberately, that it is
// common enough.
//
// Box-drawing, arrows, bullets and typographic punctuation are in every font
// shipped with a terminal. Block elements (U+2580–U+259F), geometric shapes
// beyond the plain bullet, emoji and Nerd Font private-use icons are not.
type GlyphSet map[rune]string

// DefaultGlyphs is the deploy tool's allow-list, and the one a tuikit tool
// starts from.
var DefaultGlyphs = GlyphSet{
	'·': "middle dot — key separator in footers",
	'—': "em dash — clause separator in messages",
	'…': "ellipsis — truncation marker",
	'→': "rightwards arrow — pending edit",
	'↑': "up arrow — key hints",
	'↓': "down arrow — key hints",
	'‹': "single left angle quote — tab strip",
	'›': "single right angle quote — tab strip",
	'▸': "black right-pointing small triangle — a collapsed group",
	'▾': "black down-pointing small triangle — an expanded group",
	'●': "black circle — badges and the dirty marker",
	'•': "bullet — masked secret",
	'✓': "check mark — success",
	'✗': "ballot X — failure",

	// Block elements, for Sparkline. Eight steps of an eighth each.
	//
	// SpinnerRange's comment used to say a spinner reaches for Braille rather
	// than "the block elements that would be a box on someone's terminal". That
	// has it backwards, and the dates are checkable: ▀ ▄ █ are in CP437, the
	// original IBM PC set, and the eighths are Unicode 1.0.1 Block Elements.
	// Braille Patterns arrived in Unicode 3.0, seven years later. Block
	// elements are the older and better-supported of the two.
	'▁': "lower one eighth block — sparkline", '▂': "lower one quarter block — sparkline",
	'▃': "lower three eighths block — sparkline", '▄': "lower half block — sparkline",
	'▅': "lower five eighths block — sparkline", '▆': "lower three quarters block — sparkline",
	'▇': "lower seven eighths block — sparkline", '█': "full block — sparkline",
	'┌': "box drawing", '─': "box drawing", '┐': "box drawing",
	'│': "box drawing", '└': "box drawing", '┘': "box drawing",
	// The tee, which a box never needs and a tree always does. A rectangle has
	// four corners and no junctions, so this was missing until comp.Branches
	// was written and the chrome guard said so.
	'├': "tree branch with more siblings below",
}

// With returns a copy of the set with additions.
//
// A copy, because a GlyphSet is a map and a map is a reference: a tool that
// added a glyph to DefaultGlyphs in place would be adding it to every other
// tool in the process, and to the guard that is supposed to catch it. Pairs
// are rune, reason, rune, reason; an odd count or a non-rune key panics,
// because both are typos rather than conditions to handle.
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
// rather than leaving it as a hole: Braille patterns are in every terminal
// font — which is exactly why a spinner reaches for them instead of the block
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
