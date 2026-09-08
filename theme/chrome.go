package theme

// Chrome is the third closed vocabulary: what an interface DRAWS WITH.
//
// Palette says which colors a tool may use and GlyphSet which characters it
// may print. Neither says what a box looks like — so every component picked
// its own, and sixteen characters ended up as literals inside comp where no
// guard could see them. A tool narrowing its glyph set for a font without
// chevrons would still have had them printed on its behalf.
//
// Kept as a closed set for the same reason the other two are. This is not a
// configuration bag: it is the small number of decisions that make an
// interface look like itself, and a tool takes [DefaultChrome] and overrides
// what it wants.
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

	// SortAsc and SortDesc mark the column a table is ordered by. They are
	// drawn on the end of a header label rather than in a column of their own,
	// so they should be one column wide or the header shifts when the sort
	// changes.
	SortAsc, SortDesc string

	// Tree is the connectors a tree draws its shape with.
	Tree TreeSet

	// ScrollTrack and ScrollThumb are a scrollbar's two characters.
	//
	// The same vocabulary as comp.Meter rotated a quarter turn: a dotted track
	// and a solid bar. Block elements (U+2580–U+259F) are the obvious choice
	// and are excluded, because a font without them draws a scrollbar as a
	// column of replacement boxes — which is the failure this whole set exists
	// to prevent.
	ScrollTrack, ScrollThumb string
	// ChevronLeft and ChevronRight wrap a tab strip, saying it cycles and
	// which keys do it — a fact about the keymap that no styling can carry.
	ChevronLeft, ChevronRight string

	// Collapsed and Expanded lead a group header in a list of nested rows.
	//
	// The same category as the chevrons: a fact no styling can carry. A header
	// that is merely bold says it is a header; ▸ says there is something under
	// it you have not seen, which is the difference between a group worth
	// opening and one that is empty.
	Collapsed, Expanded string

	// Gap is the columns between panes, and rows between stacked ones.
	Gap int
	// Inset is the distance from a pane's border to its content. One, in every
	// interface anyone has drawn; a field because a tool with no borders wants
	// zero.
	Inset int
	// Indent is the columns a nested row moves right per level of depth. Two,
	// in both tools that have grouped rows, in five separate places — and a
	// field rather than a constant because it is a spacing decision like the
	// other two, and a tool that wants a tighter tree should not have to
	// reindent every row itself.
	Indent int
}

// BoxSet is the six characters a frame is drawn with.
//
// Six, not eleven: there are no tee or cross pieces, so a title sits inside
// the top edge rather than breaking it. That is a constraint the default glyph
// set imposes and the components were designed around — a set with tees would
// need components that know what to do with them.
type BoxSet struct {
	TopLeft, Top, TopRight          string
	Left, Right                     string
	BottomLeft, Bottom, BottomRight string
}

// The box sets. Light is the default and the only one in DefaultGlyphs; the
// others need their characters added deliberately, which is the decision worth
// making rather than a preference worth defaulting.
var (
	// LightBox is the single-line frame the deploy tool and democtl draw.
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

// DefaultChrome is what a tuikit tool starts from, and every character in it
// is in [DefaultGlyphs].
var DefaultChrome = Chrome{
	Box: LightBox,
	// A blank divider, so the gap between two panes reads as space rather than
	// as a third thing. A tool that wants a visible seam sets "│".
	Divider:  " ",
	VDivider: " ",

	Ellipsis:    "…",
	Separator:   " · ",
	ScrollTrack: "·",
	ScrollThumb: "│",
	ScrollUp:    "↑",
	ScrollDown:  "↓",
	SortAsc:     "↑",
	SortDesc:    "↓",

	ChevronLeft:  "‹",
	ChevronRight: "›",

	// A trailing space, because these sit in the marker's column and the row's
	// text follows immediately. Part of the glyph rather than something every
	// caller remembers to add.
	Collapsed: "▸ ",
	Expanded:  "▾ ",

	Tree: UnicodeTree,

	Gap:    1,
	Inset:  1,
	Indent: 2,
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
		c.SortAsc, c.SortDesc,
		c.ScrollTrack, c.ScrollThumb,
		c.ChevronLeft, c.ChevronRight,
		c.Collapsed, c.Expanded,
		c.Tree.Vertical, c.Tree.Branch, c.Tree.Last, c.Tree.Gap,
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

// TreeSet is what [comp.Branches] draws a hierarchy's shape with.
//
// Four strings, and all four must be the same width or the rows below a branch
// stop lining up with the rows beside it. Two columns is the usual choice and
// is what DefaultChrome uses.
//
// # Why this is themeable rather than fixed
//
// htop carries two of these — `CRT_treeStrUtf8` and `CRT_treeStrAscii` — and
// picks between them from the locale, because a terminal running under
// `LANG=C` draws box-drawing characters as replacement boxes and a tree
// becomes a column of them. That is the same failure the whole GlyphSet exists
// to prevent, found in a tool that has been shipping since 2004.
//
// So [ASCIITree] is here beside the default, and a tool that has to run on
// somebody else's server can choose it.
type TreeSet struct {
	// Vertical passes through the column of an ancestor that has more children
	// coming: the │ in "│  └─".
	Vertical string
	// Branch is a child with siblings after it.
	Branch string
	// Last is the final child, which closes the line rather than continuing it.
	Last string
	// Gap is the blank under an ancestor with nothing left below, and it is a
	// character rather than an absence so that all four widths match.
	Gap string
}

// UnicodeTree is the default: box-drawing connectors, two columns each.
var UnicodeTree = TreeSet{Vertical: "│ ", Branch: "├─", Last: "└─", Gap: "  "}

// ASCIITree is htop's fallback set, for a terminal that cannot draw the
// others.
var ASCIITree = TreeSet{Vertical: "| ", Branch: "|-", Last: "`-", Gap: "  "}

// WithTree returns a copy drawing its tree with the given connectors.
//
// Separate from [Chrome.With] because a box set and a tree set are independent
// choices: a tool can want rounded corners and ASCII branches, or the reverse.
// Choosing the ASCII fallback for a terminal that has nothing means both.
func (c Chrome) WithTree(t TreeSet) Chrome { c.Tree = t; return c }
