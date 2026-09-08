package comp

import (
	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/fuzzy"
)

// Palette is one key to everything a tool can do right now.
//
// From the deploy tool's navigation design. It replaces having to remember
// twenty single-key bindings while still teaching them, because every row
// shows its key — a directory of the keyboard rather than a replacement for
// it.
//
// # What it knows and does not know
//
// It takes GROUPS of items and knows nothing about what a group is. The deploy
// tool's are scope tiers — service, node, environment, global — and that word
// does not appear here, because "what scopes exist" is a fact about the deploy
// tool and not about palettes. A tool with two groups called "commands" and
// "go to" gets the same component.
//
// It does the ranking itself, with [fuzzy], because a palette whose caller
// ranked its own items is a palette that could disagree with the letters it
// underlines.
//
// # A refused item is shown, not hidden
//
// `drain` on a swarm leader belongs in the list, in Danger, with the reason
// where its key would be. Hiding it is the easy thing and the wrong one: a
// command that vanished is indistinguishable from one that never existed, and
// the reader goes looking for it in the documentation instead of reading the
// sentence that explains it.
type Palette struct {
	// Title is painted into the top border — "run", "go to".
	Title string
	// Query is what has been typed, and Caret where in it.
	Query string
	Caret int
	// Groups are in the order they should be offered. With nothing typed they
	// are headings; with a query they become a column.
	Groups []PaletteGroup

	// Note is the line under the query: what the groups mean, or what pinning
	// does. Right is right-aligned on the same line.
	Note, Right string

	// Name owns the box, Item each row, so a click resolves to a command by
	// name rather than by counting.
	Name, Item Name

	Border, TitleStyle, PromptStyle, QueryStyle            *lipgloss.Style
	GroupStyle, LabelStyle, KeyStyle, HintStyle, NoteStyle *lipgloss.Style
	MatchStyle, WarnStyle, DangerStyle, SelectedFG, Cursor *lipgloss.Style

	// Hints is the footer.
	Hints []Hint

	list  List
	shown []PaletteItem
}

// PaletteGroup is a heading and the things under it.
type PaletteGroup struct {
	Name string
	// Note is a fact about the group, right-aligned on its heading: "2 tasks on
	// 2 nodes", "works on any screen, connected or not".
	Note  string
	Items []PaletteItem
}

// PaletteItem is one thing the reader could do.
type PaletteItem struct {
	Label string
	// Key is the single-key binding, if there is one. This is what makes the
	// palette teach rather than replace the keyboard.
	Key string
	// Hint is the consequence: "rolling, start-first", "measures 5 nodes over
	// ssh · ~60s".
	Hint string

	// Warn marks a hint as something to read before pressing enter — "asks
	// first · 8M reclaimable".
	Warn bool
	// Refused means it cannot run, and Hint says why. It is still listed.
	Refused bool

	// Group is filled in by the palette when a query flattens the list, so a
	// row can say which group it came from.
	Group string
	// match is where the query matched the label.
	match []int
}

// Runnable reports whether pressing enter on this would do anything.
func (i PaletteItem) Runnable() bool { return !i.Refused }

const (
	paletteCursorW = 2
	paletteKeyW    = 5
	paletteLabelW  = 28
	paletteFlatW   = 34
	paletteGroupW  = 20
)

// Move steps the selection. Group headings are skipped, because [List] is told
// they are not selectable.
func (p *Palette) Move(by int) { p.list.Move(by) }

// Select puts the cursor on a row of the LAST frame.
func (p *Palette) Select(i int) { p.list.Select(i) }

// Selected is the item under the cursor, or false when the cursor is on a
// heading or the list is empty.
func (p *Palette) Selected() (PaletteItem, bool) {
	i := p.list.Cursor()
	if i < 0 || i >= len(p.shown) {
		return PaletteItem{}, false
	}
	item := p.shown[i]
	if item.Label == "" {
		return PaletteItem{}, false // a heading or a blank
	}
	return item, true
}

// Matches is how many items the query left, and Total how many there are.
func (p *Palette) Matches() int { return len(p.items()) }

// Total is every item, whatever is typed.
func (p *Palette) Total() int {
	n := 0
	for _, g := range p.Groups {
		n += len(g.Items)
	}
	return n
}

// items is the query applied: everything when nothing is typed, ranked matches
// otherwise.
//
// Ranked, never filtered BY GROUP: the groups are an order of preference, not
// a filter, and a reader who typed the name of a global command should find it
// whatever they had selected.
func (p *Palette) items() []PaletteItem {
	var all []PaletteItem
	for _, g := range p.Groups {
		for _, item := range g.Items {
			item.Group = g.Name
			all = append(all, item)
		}
	}
	if p.Query == "" {
		return all
	}
	labels := make([]string, len(all))
	for i, item := range all {
		labels[i] = item.Label
	}
	ranked := fuzzy.Rank(p.Query, labels)
	out := make([]PaletteItem, 0, len(ranked))
	for _, r := range ranked {
		item := all[r.Index]
		item.match = r.At
		out = append(out, item)
	}
	return out
}

// Draw paints the palette into r and returns the box it took.
func (p *Palette) Draw(c *Canvas, r Rect) Rect {
	inner := Pane{Focused: true, Focus: p.Border, Border: p.Border,
		Title: p.Title, TitleStyle: p.TitleStyle}.Draw(c, r, Region(p.Name))
	if inner.Empty() {
		return r
	}
	id := Region(p.Name)

	// The query, then a note, then a rule.
	x := inner.X + c.Text(inner.X, inner.Y, "> ", p.PromptStyle, id)
	Input{
		Text: p.Query, Cursor: p.Caret, Focused: true,
		Placeholder:      "type to search " + itoa(p.Total()) + " commands",
		TextStyle:        p.QueryStyle,
		PlaceholderStyle: p.NoteStyle,
		CursorFG:         p.SelectedFG, CursorBG: p.Cursor,
	}.Draw(c, Rect{X: x, Y: inner.Y, W: inner.Right() - x + 1, H: 1}, id)

	y := inner.Y + 1
	if p.Note != "" || p.Right != "" {
		Bar{
			Left:  []Segment{{Text: p.Note, Style: p.NoteStyle}},
			Right: []Segment{{Text: p.Right, Style: p.NoteStyle}},
		}.Draw(c, Rect{X: inner.X, Y: y, W: inner.W, H: 1}, id)
		y++
	}
	Rule{Style: p.Border}.Draw(c, Rect{X: inner.X, Y: y, W: inner.W, H: 1}, id)
	y++

	// The rows, then a rule and the footer.
	body := Rect{X: inner.X, Y: y, W: inner.W, H: max(0, inner.Bottom()-y-1)}
	p.rows(c, body)

	Rule{Style: p.Border}.Draw(c, Rect{X: inner.X, Y: inner.Bottom() - 1, W: inner.W, H: 1}, id)
	KeyHints(c, Rect{X: inner.X, Y: inner.Bottom(), W: inner.W, H: 1}, id, p.NoteStyle, p.Hints...)
	return r
}

// rows builds the list and draws it.
func (p *Palette) rows(c *Canvas, r Rect) {
	items := p.items()
	flat := p.Query != ""

	var rows []Row
	p.shown = nil
	if flat {
		for _, item := range items {
			rows = append(rows, p.row(item, true))
			p.shown = append(p.shown, item)
		}
	} else {
		// Grouped, in the order the caller gave. A heading is a row the cursor
		// passes over — comp.List's Skip.
		for _, g := range p.Groups {
			if len(g.Items) == 0 {
				continue
			}
			rows = append(rows, Row{Skip: true, Spans: []Segment{
				{Text: Pad(" "+g.Name, max(0, r.W-Width(g.Note)-1)), Style: p.GroupStyle},
				{Text: g.Note, Style: p.NoteStyle},
			}})
			p.shown = append(p.shown, PaletteItem{})
			for _, item := range g.Items {
				item.Group = g.Name
				rows = append(rows, p.row(item, false))
				p.shown = append(p.shown, item)
			}
		}
	}

	p.list.Name = p.Item
	p.list.Focused = true
	p.list.Selected = p.SelectedFG
	p.list.Unfocused = p.SelectedFG
	p.list.Empty = "  nothing matches"
	p.list.EmptyStyle = p.NoteStyle
	p.list.Draw(c, r, rows)
}

// row is one command, in the columns the design asks for.
func (p *Palette) row(item PaletteItem, flat bool) Row {
	label, hint := p.LabelStyle, p.HintStyle
	key := item.Key
	switch {
	case item.Refused:
		// The key column goes blank: a refusal has no keystroke, and leaving
		// one there offers something that will not happen.
		label, hint, key = p.NoteStyle, p.DangerStyle, ""
	case item.Warn:
		hint = p.WarnStyle
	}

	width := paletteLabelW
	if flat {
		width = paletteFlatW
	}
	spans := []Segment{{Text: Pad("", paletteCursorW)}}
	// The matched letters are marked, so the ORDER is something a reader can
	// check rather than take on trust.
	spans = append(spans, Highlight(Pad(item.Label, width), item.match, label, p.MatchStyle)...)
	spans = append(spans, Segment{Text: Pad(key, paletteKeyW), Style: p.KeyStyle})
	if flat {
		spans = append(spans, Segment{Text: Pad(item.Group, paletteGroupW), Style: p.NoteStyle})
	}
	spans = append(spans, Segment{Text: item.Hint, Style: hint})
	return Row{Spans: spans}
}
