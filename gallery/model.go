package gallery

import (
	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/theme"
)

// Regions the gallery draws, named so its own frames can be driven by a script
// the same way a tool's are.
const (
	regHeader   comp.Name = "header"
	regRule     comp.Name = "header.rule"
	regFooter   comp.Name = "footer"
	regIndex    comp.Name = "index"
	regIndexRow comp.Name = "index.row"
	regEntry    comp.Name = "entry"
	regStateTab comp.Name = "entry.state"
	regPreview  comp.Name = "preview"
	regFacts    comp.Name = "facts"
)

// Model is the gallery.
type Model struct {
	sty     styles
	entries []Entry

	index comp.List
	state int
	// focus says whether the index or the preview has the keys. The preview
	// takes them so a component can be operated rather than only looked at.
	preview bool

	width, height int
	canvas        *comp.Canvas
}

// New builds the gallery over a palette.
//
// A palette rather than the default hard-coded, because the gallery is where
// someone checks that their own vocabulary works: a component drawn in a
// palette that does not name the roles it needs is exactly the failure this is
// meant to surface, and it should surface HERE rather than in their tool.
func New(p theme.Palette) *Model {
	m := &Model{sty: newStyles(p), width: 132, height: 38}
	m.entries = m.Entries()
	m.index = comp.List{
		Name:      regIndexRow,
		Empty:     "  no components",
		Selected:  &m.sty.selected,
		Unfocused: &m.sty.focused,
		Status:    &m.sty.muted,
	}
	return m
}

// SetSize is what the capture harness calls instead of waiting for a terminal.
func (m *Model) SetSize(w, h int) { m.width, m.height = w, h }

// Canvas is the last frame, so a script can address a region by name.
func (m *Model) Canvas() *comp.Canvas { return m.canvas }

// Entry is the component being shown, and State the variant of it.
func (m *Model) entry() Entry {
	return m.entries[clamp(m.index.Cursor(), 0, len(m.entries)-1)]
}

func (m *Model) Init() tea.Cmd { return nil }

func (m *Model) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.WindowSizeMsg:
		m.width, m.height = msg.Width, msg.Height
	case tea.MouseMsg:
		m.mouse(msg)
	case tea.KeyMsg:
		return m, m.key(msg)
	}
	return m, nil
}

func (m *Model) key(msg tea.KeyMsg) tea.Cmd {
	switch msg.String() {
	case "q", "ctrl+c":
		return tea.Quit
	case "j", "down":
		m.index.Move(1)
		m.state = 0
	case "k", "up":
		m.index.Move(-1)
		m.state = 0
	case "tab":
		m.preview = !m.preview
	case "right", "l", "]":
		m.state = min(m.state+1, len(m.entry().States)-1)
	case "left", "h", "[":
		m.state = max(m.state-1, 0)
	}
	return nil
}

func (m *Model) mouse(msg tea.MouseMsg) {
	if m.canvas == nil {
		return
	}
	id := m.canvas.OwnerAt(msg.X, msg.Y)
	switch {
	case msg.Button == tea.MouseButtonWheelUp:
		m.index.Scroll(-3)
	case msg.Button == tea.MouseButtonWheelDown:
		m.index.Scroll(3)
	case msg.Action == tea.MouseActionPress && msg.Button == tea.MouseButtonLeft:
		switch id.Name {
		case regIndexRow:
			m.index.Select(id.Index)
			m.state, m.preview = 0, false
		case regStateTab:
			m.state, m.preview = id.Index, true
		case regPreview, regEntry:
			m.preview = true
		}
	}
}

// View draws the gallery.
func (m *Model) View() string {
	c := comp.NewCanvas(m.width, 3+m.bodyHeight())
	m.header(c)

	body := comp.Rect{X: 0, Y: 2, W: m.width, H: m.bodyHeight()}
	indexWidth := max(m.width/4, 18)
	m.drawIndex(c, comp.Rect{X: 0, Y: body.Y, W: indexWidth, H: body.H})
	m.drawEntry(c, comp.Rect{X: indexWidth + 1, Y: body.Y, W: body.W - indexWidth - 1, H: body.H})

	m.footer(c)
	m.canvas = c
	return c.String()
}

func (m *Model) bodyHeight() int { return max(6, m.height-4) }

func (m *Model) header(c *comp.Canvas) {
	comp.Bar{
		Left: []comp.Segment{
			{Text: "tuikit gallery", Style: &m.sty.title},
			{Text: "  every component, running", Style: &m.sty.muted},
		},
		Right: []comp.Segment{
			{Text: itoa(len(m.entries)) + " components", Style: &m.sty.muted},
		},
	}.Draw(c, comp.Rect{X: 0, Y: 0, W: m.width, H: 1}, comp.Region(regHeader))
	c.Fill(comp.Rect{X: 0, Y: 1, W: m.width, H: 1}, "─", &m.sty.border, comp.Region(regRule))
}

func (m *Model) footer(c *comp.Canvas) {
	comp.KeyHints(c, comp.Rect{X: 0, Y: 2 + m.bodyHeight(), W: m.width, H: 1},
		comp.Region(regFooter), &m.sty.muted,
		comp.Hint{Key: "↑↓", Label: "component"},
		comp.Hint{Key: "‹›", Label: "state"},
		comp.Hint{Key: "tab", Label: "pane"},
		comp.Hint{Key: "q", Label: "quit"})
}

func (m *Model) drawIndex(c *comp.Canvas, r comp.Rect) {
	inner := comp.Pane{
		Title: "Components", TitleAt: comp.TitleOnRow, Focused: !m.preview,
		Border: &m.sty.border, Focus: &m.sty.focused, TitleStyle: &m.sty.title,
	}.Draw(c, r, comp.Region(regIndex))

	rows := make([]comp.Row, len(m.entries))
	for i, e := range m.entries {
		rows[i] = comp.Row{Text: " " + e.Name}
	}
	m.index.Focused = !m.preview
	m.index.Draw(c, inner, rows)
}

// drawEntry is the right-hand half: what the component is, the state picker,
// the component itself, and the facts you need before choosing it.
func (m *Model) drawEntry(c *comp.Canvas, r comp.Rect) {
	e := m.entry()
	inner := comp.Pane{
		Title: e.Name, TitleAt: comp.TitleOnRow, Focused: m.preview,
		Border: &m.sty.border, Focus: &m.sty.focused, TitleStyle: &m.sty.title,
	}.Draw(c, r, comp.Region(regEntry))
	if inner.Empty() {
		return
	}

	y := inner.Y
	y += c.Text(inner.X+1, y, comp.Truncate(e.Summary, inner.W-2), &m.sty.muted, comp.Region(regEntry))
	y = inner.Y + 1

	if len(e.States) > 0 {
		m.state = clamp(m.state, 0, len(e.States)-1)
		tabs := make([]comp.Tab, len(e.States))
		for i, s := range e.States {
			tabs[i] = comp.Tab{Name: s.Name}
		}
		comp.Tabs{
			Tabs: tabs, Active: m.state, Focused: m.preview,
			Style: &m.sty.muted, Selected: &m.sty.focused,
			FocusSelected: &m.sty.title, Chrome: &m.sty.muted,
		}.Draw(c, comp.Rect{X: inner.X + 1, Y: y, W: inner.W - 1, H: 1}, regStateTab)
		y++

		if note := e.States[m.state].Note; note != "" {
			c.Text(inner.X+1, y, comp.Truncate(note, inner.W-2), &m.sty.muted, comp.Region(regEntry))
		}
		y += 2
	}

	facts := m.facts(e)
	previewH := max(3, inner.H-(y-inner.Y)-len(facts)-1)
	preview := comp.Rect{X: inner.X + 1, Y: y, W: inner.W - 2, H: previewH}

	if len(e.States) > 0 && e.States[m.state].Draw != nil {
		// Into a canvas of its own, so a component that overruns its rect
		// cannot scribble on the gallery's own chrome — and the gallery finds
		// out about the overrun rather than absorbing it.
		sub := comp.NewCanvas(preview.W, preview.H)
		e.States[m.state].Draw(sub, sub.Bounds(), m.preview)
		blit(c, preview, sub)
	}

	y = preview.Y + preview.H
	for _, line := range facts {
		if y > inner.Bottom() {
			return
		}
		c.Text(inner.X+1, y, comp.Truncate(line.text, inner.W-2), line.style, comp.Region(regFacts))
		y++
	}
}

type fact struct {
	text  string
	style *lipgloss.Style
}

func (m *Model) facts(e Entry) []fact {
	var out []fact
	add := func(label, body string) {
		if body != "" {
			out = append(out, fact{text: label + "  " + body, style: &m.sty.muted})
		}
	}
	add("keys   ", comp.Hints(e.Keys...))
	add("mouse  ", comp.Hints(e.Mouse...))
	add("roles  ", join(e.Roles, " · "))
	add("glyphs ", join(e.Glyphs, " "))
	add("from   ", e.From)
	return out
}

func join(parts []string, sep string) string {
	var out string
	for i, p := range parts {
		if i > 0 {
			out += sep
		}
		out += p
	}
	return out
}

func clamp(v, lo, hi int) int { return max(lo, min(v, hi)) }

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var d [20]byte
	i := len(d)
	for n > 0 {
		i--
		d[i] = byte('0' + n%10)
		n /= 10
	}
	return string(d[i:])
}
