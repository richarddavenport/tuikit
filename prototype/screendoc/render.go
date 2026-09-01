package screendoc

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"

	"github.com/richarddavenport/tuikit/theme"
)

// Renderer turns a document into a frame.
type Renderer struct {
	styles map[string]lipgloss.Style
	glyphs map[rune]bool
	data   Data
}

// New builds a renderer, resolving the document's style names through a
// palette. This is the only place a role becomes a colour.
func New(doc Doc, p theme.Palette, data Data) (*Renderer, error) {
	r := &Renderer{styles: map[string]lipgloss.Style{}, glyphs: map[rune]bool{}, data: data}
	for name, s := range doc.Styles {
		fg, ok := role(p, s.Role)
		if !ok {
			return nil, fmt.Errorf("style %q names role %q, which is not in the palette", name, s.Role)
		}
		st := lipgloss.NewStyle().Foreground(fg).Bold(s.Bold)
		if s.BG != "" {
			bg, ok := role(p, s.BG)
			if !ok {
				return nil, fmt.Errorf("style %q names background role %q, which is not in the palette", name, s.BG)
			}
			st = st.Background(bg)
		}
		r.styles[name] = st
	}
	for _, g := range doc.Glyphs {
		for _, c := range g {
			r.glyphs[c] = true
		}
	}
	return r, nil
}

// role resolves a palette role by name. Named rather than indexed so a
// document says Accent and survives someone deciding the interface is blue.
func role(p theme.Palette, name string) (lipgloss.Color, bool) {
	switch name {
	case "Accent":
		return p.Accent, true
	case "Muted":
		return p.Muted, true
	case "Border":
		return p.Border, true
	case "Success":
		return p.Success, true
	case "Pending":
		return p.Pending, true
	case "Danger":
		return p.Danger, true
	case "Stderr":
		return p.Stderr, true
	case "SelectionFG":
		return p.SelectionFG, true
	case "SelectionBG":
		return p.SelectionBG, true
	}
	// Extra is the palette's escape hatch for a tenth meaning, held as named
	// roles rather than as a colour.
	for _, e := range p.Extra {
		if e.Name == name {
			return e.Color, true
		}
	}
	return "", false
}

func (r *Renderer) style(name string) lipgloss.Style {
	if s, ok := r.styles[name]; ok {
		return s
	}
	return lipgloss.NewStyle()
}

func (r *Renderer) paint(name, s string) string {
	if name == "" {
		return s
	}
	return r.style(name).Render(s)
}

// Render draws the document at a size.
func (r *Renderer) Render(doc Doc, w, h int) string {
	return strings.Join(r.node(doc.Root, w, h), "\n")
}

// matches evaluates a condition against the data.
func (r *Renderer) matches(w *When) bool {
	if w == nil {
		return true
	}
	v := r.data.Value(w.Key)
	if w.Not != nil {
		return v != *w.Not
	}
	return v == w.Is
}

// live drops the children whose condition does not hold, before any space is
// distributed — a hidden node must not reserve a column.
func (r *Renderer) live(nodes []Node) []Node {
	out := make([]Node, 0, len(nodes))
	for _, n := range nodes {
		if r.matches(n.When) {
			out = append(out, n)
		}
	}
	return out
}

// node draws one element into exactly the space it was given.
func (r *Renderer) node(n Node, w, h int) []string {
	switch n.Type {
	case "column":
		return r.column(n, w, h)
	case "row":
		return r.row(n, w, h)
	case "statusbar":
		return []string{r.statusbar(n, w)}
	case "rule":
		return []string{r.paint(n.Style, strings.Repeat("─", w))}
	case "text":
		return []string{r.paint(n.Style, clip(r.expand(n.Text, w), w))}
	case "box":
		return r.box(n, w, h)
	case "list":
		return r.list(n, w)
	case "tabs":
		return []string{r.tabs(n, w)}
	case "fields":
		return r.fields(n, w)
	case "group":
		// A stack of children under one condition. Conditions do not compose —
		// there is no "and" — and this is why they do not need to: "the
		// Overview tab, when something is selected" is a group inside a group,
		// which is also the shape a builder would produce by dragging one box
		// into another.
		var out []string
		for _, c := range r.live(n.Children) {
			out = append(out, r.node(c, w, 0)...)
		}
		return out
	case "spacer":
		return []string{""}
	}
	return []string{r.paint("danger", fmt.Sprintf("node type %q has no renderer", n.Type))}
}

// column stacks children, giving each what its constraint asks for.
func (r *Renderer) column(n Node, w, h int) []string {
	children := r.live(n.Children)
	sizes := distribute(children, h-n.Reserve, n.Gap*(len(children)-1))
	var out []string
	for i, c := range children {
		if i > 0 && n.Gap > 0 {
			out = append(out, make([]string, n.Gap)...)
		}
		out = append(out, r.node(c, w, sizes[i])...)
	}
	return out
}

// row places children side by side and joins them line by line.
func (r *Renderer) row(n Node, w, h int) []string {
	children := r.live(n.Children)
	sizes := distribute(children, w, n.Gap*(len(children)-1))
	cols := make([][]string, len(children))
	for i, c := range children {
		cols[i] = r.node(c, sizes[i], h)
	}
	out := make([]string, h)
	for y := 0; y < h; y++ {
		var b strings.Builder
		for i := range cols {
			if i > 0 {
				b.WriteString(strings.Repeat(" ", n.Gap))
			}
			line := ""
			if y < len(cols[i]) {
				line = cols[i][y]
			}
			b.WriteString(padVisible(line, sizes[i]))
		}
		out[y] = strings.TrimRight(b.String(), " ")
	}
	return out
}

// distribute hands out space: fixed first, then ratios, then whatever is left
// to the fillers. Integer arithmetic throughout — a cell is not divisible, and
// a layout that rounds differently on two runs has goldens that flicker.
//
// A ratio is taken against the WHOLE width, gaps included; only the filler
// pays for them. Which way round that goes is not a detail — a third of 132 is
// 44 and a third of 131 is 43, and the pane is a column narrower for the rest
// of the screen's life. The format has to state it, and this states it.
func distribute(children []Node, total, gaps int) []int {
	sizes := make([]int, len(children))
	left := total - gaps
	for i, c := range children {
		switch {
		case c.Constraint == nil:
			sizes[i] = 1
			left -= 1
		case c.Constraint.Fixed > 0:
			sizes[i] = c.Constraint.Fixed
			left -= sizes[i]
		case len(c.Constraint.Ratio) == 2:
			v := total * c.Constraint.Ratio[0] / c.Constraint.Ratio[1]
			if v < c.Constraint.Min {
				v = c.Constraint.Min
			}
			sizes[i] = v
			left -= v
		}
	}
	for i, c := range children {
		if c.Constraint != nil && c.Constraint.Fill {
			sizes[i] = left
			if c.Constraint.Min > 0 && sizes[i] < c.Constraint.Min {
				sizes[i] = c.Constraint.Min
			}
		}
	}
	return sizes
}

// statusbar puts spans left and right, w columns apart.
func (r *Renderer) statusbar(n Node, w int) string {
	left, right := r.spans(n.Left), r.spans(n.Right)
	gap := w - lipgloss.Width(left) - lipgloss.Width(right)
	if gap < 1 {
		return clip(left, w)
	}
	return left + strings.Repeat(" ", gap) + right
}

func (r *Renderer) spans(spans []Span) string {
	var b strings.Builder
	for _, s := range spans {
		if s.Switch != "" {
			s = s.Cases[strconv.FormatBool(r.data.Bool(s.Switch))]
		}
		b.WriteString(r.paint(s.Style, r.expand(s.Text, 0)))
	}
	return b.String()
}

// box draws a titled frame with the six characters the glyph set allows.
func (r *Renderer) box(n Node, w, h int) []string {
	edge := "border"
	if n.Focused != "" && r.data.Bool(n.Focused) {
		edge = "focused"
	}
	inner := w - 2

	var rows []string
	for _, c := range r.live(n.Body) {
		rows = append(rows, r.node(c, inner, 0)...)
	}

	out := []string{r.paint(edge, "┌"+strings.Repeat("─", inner)+"┐")}
	body := h - 2
	if title := r.title(n, inner); title != "" {
		out = append(out, r.paint(edge, "│")+r.paint("title", pad(" "+trim(title, inner-1), inner))+r.paint(edge, "│"))
		body--
	}
	for i := 0; i < body; i++ {
		line := ""
		if i < len(rows) {
			line = rows[i]
		}
		out = append(out, r.paint(edge, "│")+padVisible(line, inner)+r.paint(edge, "│"))
	}
	return append(out, r.paint(edge, "└"+strings.Repeat("─", inner)+"┘"))
}

// title picks the first candidate whose condition holds.
func (r *Renderer) title(n Node, w int) string {
	for _, t := range n.Titles {
		if r.matches(t.When) {
			return r.expand(t.Text, w)
		}
	}
	return r.expand(n.Title, w)
}

// list draws bound rows, one per line, styled by a field rather than by a
// colour written into the document.
func (r *Renderer) list(n Node, w int) []string {
	rows := r.data.Rows(n.Rows)
	if len(rows) == 0 {
		return []string{r.paint(n.EmptyStyle, pad(n.Empty, w))}
	}
	cur, focused := -1, false
	if n.Cursor != "" {
		cur = r.data.Index(n.Cursor)
	}
	if n.Selected != "" && n.Unfocused != "" {
		focused = r.data.Bool(n.Rows + ".focused")
	}

	out := make([]string, 0, len(rows))
	for i, row := range rows {
		if n.GlyphBy != "" {
			row["mark"] = n.GlyphMap[row[n.GlyphBy]]
		}
		line := pad(r.fill(n.Template, row, w), w)

		style := ""
		if n.StyleBy != "" {
			style = n.StyleMap[row[n.StyleBy]]
		}
		if i == cur {
			style = n.Unfocused
			if focused {
				style = n.Selected
			}
		}
		out = append(out, r.paint(style, line))
	}
	return out
}

// tabs draws a strip, the active one wrapped in the glyphs the set allows.
func (r *Renderer) tabs(n Node, w int) string {
	active := r.data.Index(n.Active)
	var b strings.Builder
	b.WriteString(" ")
	for i, name := range n.Names {
		label := " " + name + " "
		if i == active {
			b.WriteString(r.paint(n.On, n.OpenWith+label+n.CloseWith))
		} else {
			b.WriteString(r.paint(n.Off, " "+label+" "))
		}
	}
	return padVisible(b.String(), w)
}

// fields draws name/value rows.
func (r *Renderer) fields(n Node, w int) []string {
	out := make([]string, 0, len(n.Fields))
	for _, f := range n.Fields {
		style := ""
		if f.StyleBy != "" {
			style = f.StyleMap[r.data.Value(f.StyleBy)]
		}
		label := fmt.Sprintf("  %-*s ", n.LabelWidth, f.Name)
		out = append(out, r.paint(n.LabelStyle, label)+r.paint(style, r.expand(f.Value, w)))
	}
	return out
}

// expand replaces {key} with a value from the data.
func (r *Renderer) expand(s string, w int) string {
	return substitute(s, func(key, spec string) string { return format(r.data.Value(key), spec, w) })
}

// fill is expand against one row rather than against the data.
func (r *Renderer) fill(s string, row map[string]string, w int) string {
	return substitute(s, func(key, spec string) string { return format(row[key], spec, w) })
}

// substitute walks {key} and {key:spec} placeholders.
func substitute(s string, get func(key, spec string) string) string {
	var b strings.Builder
	for {
		open := strings.IndexByte(s, '{')
		if open < 0 {
			b.WriteString(s)
			return b.String()
		}
		close := strings.IndexByte(s[open:], '}')
		if close < 0 {
			b.WriteString(s)
			return b.String()
		}
		close += open
		b.WriteString(s[:open])
		key, spec, _ := strings.Cut(s[open+1:close], ":")
		b.WriteString(get(key, spec))
		s = s[close+1:]
	}
}

// format applies a placeholder's spec.
//
//	-14      left-aligned in 14 columns, trimmed to fit
//	fit-12   trimmed to the container's width less 12
//	wrap-2   the first wrapped line, at the container's width less 2
//
// The width-relative forms exist because the hand-written screen has them:
// trim(svc.Image, w-14) is a real line of democtl. A format with only absolute
// widths cannot say it, and a document that cannot say it is a document that
// only works at one terminal size.
func format(v, spec string, w int) string {
	switch {
	case spec == "":
		return v
	case strings.HasPrefix(spec, "fit-"):
		n, err := strconv.Atoi(spec[4:])
		if err != nil {
			return v
		}
		return trim(v, w-n)
	case strings.HasPrefix(spec, "wrap-"):
		n, err := strconv.Atoi(spec[5:])
		if err != nil {
			return v
		}
		return wrapFirst(v, w-n)
	}
	n, err := strconv.Atoi(strings.TrimPrefix(spec, "-"))
	if err != nil {
		return v
	}
	return fmt.Sprintf("%-*s", n, trim(v, n))
}

// wrap breaks text at word boundaries; wrapFirst is its first line, for
// somewhere with room for one.
func wrap(s string, w int) []string {
	if w <= 0 {
		return nil
	}
	var out []string
	line := ""
	for _, word := range strings.Fields(s) {
		switch {
		case line == "":
			line = word
		case len([]rune(line))+1+len([]rune(word)) <= w:
			line += " " + word
		default:
			out = append(out, line)
			line = word
		}
		for len([]rune(line)) > w {
			r := []rune(line)
			out = append(out, string(r[:w]))
			line = string(r[w:])
		}
	}
	if line != "" {
		out = append(out, line)
	}
	return out
}

func wrapFirst(s string, w int) string {
	lines := wrap(s, w)
	if len(lines) == 0 {
		return ""
	}
	if len(lines) > 1 {
		return trim(lines[0], w)
	}
	return lines[0]
}

// The text helpers, which a real comp package would own. Copied rather than
// imported because democtl's are unexported, and because a prototype that
// reaches into the thing it is meant to replace proves less.

func trim(s string, w int) string {
	if w <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	if w == 1 {
		return "…"
	}
	return string(r[:w-1]) + "…"
}

func clip(s string, w int) string {
	if w <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= w {
		return s
	}
	return ansi.Truncate(s, w, "…")
}

func pad(s string, w int) string {
	s = trim(s, w)
	if gap := w - len([]rune(s)); gap > 0 {
		return s + strings.Repeat(" ", gap)
	}
	return s
}

func padVisible(s string, w int) string {
	if lipgloss.Width(s) > w {
		return clip(s, w)
	}
	return s + strings.Repeat(" ", w-lipgloss.Width(s))
}
