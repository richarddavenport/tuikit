package theme

// Chrome is the third closed vocabulary: what an interface DRAWS WITH.
//
// Palette says which colours a tool may use and GlyphSet which characters it
// may print. Neither says what a box looks like — so every component picked its
// own, and sixteen characters ended up as literals inside comp where no guard
// could see them. A tool narrowing its glyph set for a font without chevrons
// would still have had them printed on its behalf.
//
// Kept as a closed set for the same reason the other two are. This is not a
// configuration bag: it is the small number of decisions that make an interface
// look like itself, and a tool takes [DefaultChrome] and overrides what it
// wants.
//
// # The rule that keeps it from separating
//
// Every character a Chrome draws must be in the GlyphSet — guard.Chrome holds
// that closed. Choosing RoundedBox means adding ╭╮╰╯ to your set deliberately,
// because they are not in DefaultGlyphs and a font without them draws four
// boxes where the corners should be. The coupling is the point: the look and
// the allow-list cannot drift apart, because one is checked against the other.
type Chrome struct {
	// Box is the frame a pane is drawn with.
	Box BoxSet

	// Divider is what fills the gap between panes, and Grabber marks it where
	// it can be dragged. A blank divider is a gap; a drawn one is a seam.
	Divider, VDivider string

	// Ellipsis marks text that was cut.
	Ellipsis string
	// Separator joins key hints: the · in "tab pane · q quit".
	Separator string
	// ScrollUp and ScrollDown say which way an off-screen selection went.
	ScrollUp, ScrollDown string
	// ChevronLeft and ChevronRight wrap a tab strip, saying it cycles and
	// which keys do it — a fact about the keymap that no styling can carry.
	ChevronLeft, ChevronRight string

	// Gap is the columns between panes, and rows between stacked ones.
	Gap int
	// Inset is the distance from a pane's border to its content. One, in every
	// interface anyone has drawn; a field because a tool with no borders wants
	// zero.
	Inset int
}

// BoxSet is the six characters a frame is drawn with.
//
// Six, not eleven: there are no tee or cross pieces, so a title sits inside the
// top edge rather than breaking it. That is a constraint the default glyph set
// imposes and the components were designed around — a set with tees would need
// components that know what to do with them.
type BoxSet struct {
	TopLeft, Top, TopRight          string
	Left, Right                     string
	BottomLeft, Bottom, BottomRight string
}

// The box sets. Light is the default and the only one in DefaultGlyphs; the
// others need their characters added deliberately, which is the decision worth
// making rather than a preference worth defaulting.
var (
	// LightBox is the single-line frame swarmctl and democtl draw.
	LightBox = BoxSet{"┌", "─", "┐", "│", "│", "└", "─", "┘"}
	// RoundedBox is softer and needs ╭╮╰╯ in the glyph set.
	RoundedBox = BoxSet{"╭", "─", "╮", "│", "│", "╰", "─", "╯"}
	// HeavyBox needs ┏┓┗┛━┃ — worth checking against the terminal fonts your
	// readers actually have, because heavy box drawing is patchier than light.
	HeavyBox = BoxSet{"┏", "━", "┓", "┃", "┃", "┗", "━", "┛"}
	// ASCIIBox draws in a font that has nothing. Ugly on purpose: it is for
	// somewhere that would otherwise show replacement boxes, and it should look
	// like a fallback rather than a style.
	ASCIIBox = BoxSet{"+", "-", "+", "|", "|", "+", "-", "+"}
)

// DefaultChrome is what a tuikit tool starts from, and every character in it is
// in [DefaultGlyphs].
var DefaultChrome = Chrome{
	Box: LightBox,
	// A blank divider, so the gap between two panes reads as space rather than
	// as a third thing. A tool that wants a visible seam sets "│".
	Divider:  " ",
	VDivider: " ",

	Ellipsis:   "…",
	Separator:  " · ",
	ScrollUp:   "↑",
	ScrollDown: "↓",

	ChevronLeft:  "‹",
	ChevronRight: "›",

	Gap:   1,
	Inset: 1,
}

// Glyphs is every character this Chrome can draw, for a guard to check against
// the allow-list.
//
// Every field, walked rather than listed, so a field added to Chrome cannot be
// forgotten here — which is the failure this whole type exists to prevent, one
// level up.
func (c Chrome) Glyphs() []rune {
	seen := map[rune]bool{}
	var out []rune
	for _, s := range []string{
		c.Box.TopLeft, c.Box.Top, c.Box.TopRight,
		c.Box.Left, c.Box.Right,
		c.Box.BottomLeft, c.Box.Bottom, c.Box.BottomRight,
		c.Divider, c.VDivider,
		c.Ellipsis, c.Separator, c.ScrollUp, c.ScrollDown,
		c.ChevronLeft, c.ChevronRight,
	} {
		for _, r := range s {
			if !seen[r] {
				seen[r] = true
				out = append(out, r)
			}
		}
	}
	return out
}

// With returns a copy with a different box set, for the common change.
func (c Chrome) With(box BoxSet) Chrome { c.Box = box; return c }
