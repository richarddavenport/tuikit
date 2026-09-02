package gallery

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/fuzzy"
	"github.com/richarddavenport/tuikit/theme"
)

// Entries is every component in comp, with the states worth looking at.
//
// The rule this list exists to serve is that a component not in the gallery is
// not finished — and gallery_test.go holds it closed by reading comp's own
// source, so adding a component without an entry fails a test rather than
// quietly shipping something nobody has ever seen run.
//
// Every entry has an EMPTY state and an OVERFLOWING one wherever those mean
// anything, because they are the two a component is most likely to get wrong
// and the two a screenshot of the happy path will never show.
func (m *Model) Entries() []Entry {
	s := &m.sty
	return []Entry{
		m.listEntry(s),
		m.paneEntry(s),
		m.tabsEntry(s),
		m.barEntry(s),
		m.confirmEntry(s),
		m.stepListEntry(s),
		m.logPaneEntry(s),
		m.spinnerEntry(s),
		m.tableEntry(s),
		m.detailEntry(s),
		m.menuEntry(s),
		m.toastEntry(s),
		m.formEntry(s),
		m.splitEntry(s),
		m.layoutEntry(s),
		m.meterEntry(s),
		m.inputEntry(s),
		m.waitingEntry(s),
		m.breadcrumbEntry(s),
		m.scrollbarEntry(s),
	}
}

func (m *Model) layoutEntry(s *styles) Entry {
	draw := func(l comp.Layout, labels ...string) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			style := []*lipgloss.Style{&s.title, &s.muted, &s.focused, &s.pending}
			for i, band := range l.Rows(r) {
				if band.H <= 0 {
					continue
				}
				c.Fill(band, " ", style[i%len(style)], comp.Region("demo.layout").At(i))
				label := labels[i%len(labels)]
				c.Text(band.X+1, band.Y, label+"  "+itoa(band.H),
					style[i%len(style)], comp.Region("demo.layout").At(i))
			}
		}
	}
	window := comp.Layout{Constraints: []comp.Constraint{
		comp.Length(1), comp.Length(1), comp.Fill(1).Min(3), comp.Length(1),
	}}
	return Entry{
		Name:    "Layout",
		Summary: "Bands down an axis, so nobody computes a height twice.",
		From:    "the body()/bodyHeight() pair every one of the four tools writes",
		Roles:   []string{"Accent", "Muted", "Pending"},
		States: []State{
			{Name: "the shape every tool has", Note: "a header, a rule, the body, the hints — one declaration instead of a rect and a height",
				Draw: draw(window, "header", "rule", "body", "hints")},
			{Name: "weights", Note: "Fill(1) and Fill(3): the leftover splits by weight, and the rounding never loses a row",
				Draw: draw(comp.Layout{Constraints: []comp.Constraint{comp.Fill(1), comp.Fill(3)}}, "Fill(1)", "Fill(3)")},
			{Name: "a clamp is paid for by the others", Note: "Fill(1).Min(6) takes its floor and the rest is re-solved — solving once and clamping after would overflow",
				Draw: draw(comp.Layout{Constraints: []comp.Constraint{
					comp.Fill(1).Min(6), comp.Fill(9),
				}}, "Fill(1).Min(6)", "Fill(9)")},
			{Name: "a gap between bands", Note: "the gap comes out of the fill, not out of the rect",
				Draw: draw(comp.Layout{Constraints: []comp.Constraint{
					comp.Length(2), comp.Fill(1), comp.Length(2),
				}, Gap: 2}, "Length(2)", "Fill(1)", "Length(2)")},
			{Name: "no room", Note: "a terminal too short gives bands of zero rather than negative ones",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					draw(window, "header", "rule", "body", "hints")(c, comp.Rect{X: r.X, Y: r.Y, W: r.W, H: 3}, focused)
				}},
		},
	}
}

func (m *Model) splitEntry(s *styles) Entry {
	draw := func(sp comp.Split, chrome func(theme.Chrome) theme.Chrome) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			if chrome != nil {
				c = c.WithChrome(chrome(c.Chrome()))
			}
			sp.Style = &s.border
			first, second := sp.Draw(c, r)
			for _, pane := range []struct {
				r     comp.Rect
				title string
			}{{first, "first"}, {second, "second"}} {
				inner := comp.Pane{
					Title: pane.title, TitleAt: comp.TitleOnRow,
					Border: &s.border, Focus: &s.focused, TitleStyle: &s.title,
				}.Draw(c, pane.r, comp.Region("demo.split"))
				if !inner.Empty() {
					c.Text(inner.X+1, inner.Y, comp.Truncate(itoa(pane.r.W)+" columns", inner.W-2),
						&s.muted, comp.Region("demo.split"))
				}
			}
		}
	}
	third := comp.Split{Name: "demo.divider", Ratio: [2]int{1, 3}, Min: 10}
	return Entry{
		Name:    "Split",
		Summary: "Two panes and the divider between them, which you can grab and drag.",
		From:    "democtl's dashboard, before it was a component",
		Mouse:   []comp.Hint{{Key: "drag", Label: "the gap between the panes IS the handle"}},
		Roles:   []string{"Border"},
		States: []State{
			{Name: "a third", Note: "the ratio is taken against the whole width, gaps included",
				Draw: draw(third, nil)},
			{Name: "dragged", Note: "At overrides the ratio; zero means nobody has touched it",
				Draw: draw(comp.Split{Name: "demo.divider", At: 60, Min: 10}, nil)},
			{Name: "at its minimum", Note: "a split draggable to nothing is a pane you cannot get back",
				Draw: draw(comp.Split{Name: "demo.divider", At: 1, Min: 20}, nil)},
			{Name: "stacked", Note: "the same divider, between rows",
				Draw: draw(comp.Split{Name: "demo.divider", Vertical: true, Ratio: [2]int{1, 2}, Min: 3}, nil)},
			{Name: "a wider gap", Note: "the gap is the chrome's, so it is one number for every split",
				Draw: draw(third, func(ch theme.Chrome) theme.Chrome { ch.Gap = 5; return ch })},
			{Name: "a visible seam", Note: "and so is what fills it — a blank reads as space, a line as a join",
				Draw: draw(third, func(ch theme.Chrome) theme.Chrome { ch.Divider = "│"; return ch })},
		},
	}
}

func (m *Model) formEntry(s *styles) Entry {
	draw := func(f comp.Form) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, focused bool) {
			f.Marker = "> "
			f.Label, f.FocusLabel, f.Value = &s.muted, &s.title, &s.focused
			f.Muted, f.Danger = &s.muted, &s.danger
			f.Draw(c, r, "demo.form")
		}
	}
	fields := []comp.Field{
		{Label: "target", Kind: comp.FieldChoice, Choices: []string{"staging", "production"}},
		{Label: "tag", Kind: comp.FieldText, Text: "v2.4.1"},
		{Label: "note", Kind: comp.FieldText, Placeholder: "(none)"},
		{Label: "dry run", Kind: comp.FieldToggle, On: true},
	}
	return Entry{
		Name:    "Form",
		Summary: "Every field at once, so you can see what you have chosen rather than remember it.",
		From:    "pgctl viewForm, swarmctl applyedits.go and confirmPhrase",
		Keys:    []comp.Hint{{Key: "↑↓", Label: "field"}, {Key: "‹›", Label: "choice"}, {Key: "space", Label: "toggle"}},
		Roles:   []string{"Accent", "Muted", "Danger"},
		States: []State{
			{Name: "focused", Note: "the marker is the caller's glyph, not one this package chose",
				Draw: draw(comp.Form{Fields: fields, Cursor: 1, Focused: true})},
			{Name: "unfocused", Note: "still shows every value — seeing what you chose is the point",
				Draw: draw(comp.Form{Fields: fields, Cursor: 1})},
			{Name: "type the name", Note: "swarmctl's rule: the difference between a keystroke and a decision",
				Draw: draw(comp.Form{Focused: true, Fields: []comp.Field{
					{Label: "service", Kind: comp.FieldText, Text: "api_gateway", Disabled: true},
					{Label: "type the name", Kind: comp.FieldText, Must: "api_gateway", Text: "api_gate"},
				}})},
			{Name: "phrase matched", Note: "and now the action is allowed to happen",
				Draw: draw(comp.Form{Focused: true, Fields: []comp.Field{
					{Label: "service", Kind: comp.FieldText, Text: "api_gateway", Disabled: true},
					{Label: "type the name", Kind: comp.FieldText, Must: "api_gateway", Text: "api_gateway"},
				}})},
			{Name: "a secret being typed", Note: "masked while you type it, and NO caret — one moving over eight identical bullets says nothing, and one that stops early says how long the secret is",
				Draw: draw(comp.Form{Focused: true, Cursor: 1, Fields: []comp.Field{
					{Label: "registry", Kind: comp.FieldText, Text: "ghcr.io"},
					{Label: "token", Kind: comp.FieldText, Text: "ghp_R3alT0kenValue", Secret: true},
				}})},
			{Name: "a secret not yet given", Note: "an unanswered field looks unanswered rather than like a value that is hidden",
				Draw: draw(comp.Form{Focused: true, Cursor: 1, Fields: []comp.Field{
					{Label: "registry", Kind: comp.FieldText, Text: "ghcr.io"},
					{Label: "token", Kind: comp.FieldText, Placeholder: "required", Secret: true},
				}})},
			{Name: "no room", Note: "stops at the edge like everything else",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					draw(comp.Form{Fields: fields, Cursor: 0, Focused: true})(
						c, comp.Rect{X: r.X, Y: r.Y, W: min(r.W, 22), H: 2}, focused)
				}},
		},
	}
}

func (m *Model) toastEntry(s *styles) Entry {
	draw := func(t comp.Toast) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			t.Border, t.BodyStyle, t.HintStyle = &s.border, &s.muted, &s.muted
			if t.Accent == nil {
				t.Accent = &s.danger
			}
			sub := comp.NewCanvas(r.W, r.H)
			t.Draw(sub, sub.Bounds(), comp.Region("demo.toast"))
			blit(c, r, sub)
		}
	}
	return Entry{
		Name:    "Toast",
		Summary: "Something the interface has to say — and, where we know it, what to do about it.",
		From:    "swarmctl errorView",
		Roles:   []string{"Danger", "Success", "Border", "Muted"},
		Glyphs:  []string{"┌", "─", "┐", "│", "└", "┘"},
		States: []State{
			{Name: "a failure with a hint", Note: "the field most likely to be left empty, so it is a field",
				Draw: draw(comp.Toast{
					Title: "cannot reach staging",
					Body:  "dial tcp 10.0.0.4:5432: connection refused",
					Hint:  "is the tunnel up? try `pgctl connect staging`",
				})},
			{Name: "a failure without one", Note: "allowed, and usually means nobody has worked out the answer yet",
				Draw: draw(comp.Toast{
					Title: "deploy failed",
					Body:  "task 3 exited 137 before the health check passed",
				})},
			{Name: "a result", Note: "the role is the caller's, so the same box says good news",
				Draw: draw(comp.Toast{
					Accent: &s.success,
					Title:  "deployed api_gateway",
					Body:   "3 of 3 tasks converged in 8.4s.",
				})},
			{Name: "bottom left", Note: "the corner is the caller's — what it must not cover is theirs to know",
				Draw: draw(comp.Toast{
					Anchor: comp.BottomLeft,
					Title:  "cannot reach staging",
					Body:   "dial tcp 10.0.0.4:5432: connection refused",
				})},
			{Name: "a long root cause", Note: "bounded, because the cause is arbitrary text from somewhere else",
				Draw: draw(comp.Toast{
					Title: "apply failed",
					Body:  strings.Repeat("a root cause that came from somewhere else and does not know how wide your terminal is. ", 2),
					Hint:  "run with -v for the full trace",
				})},
		},
	}
}

func (m *Model) menuEntry(s *styles) Entry {
	items := []comp.Hint{
		{Key: "L", Label: "View logs"},
		{Key: "D", Label: "Deploy"},
		{Key: "R", Label: "Restart"},
	}
	menu := func(cursor int) comp.Menu {
		return comp.Menu{
			Name: "demo.menu", Item: "demo.menu.item",
			Items: items, Cursor: cursor,
			Border: &s.focused, Style: &s.muted, Selected: &s.selected,
		}
	}
	return Entry{
		Name:    "Menu",
		Summary: "A short list of actions, at a point or on the thing they act on.",
		From:    "democtl's context menu; azctl needs the same one",
		Keys: []comp.Hint{
			{Key: "↑↓", Label: "choose"}, {Key: "enter", Label: "do it"}, {Key: "esc", Label: "close"},
		},
		Mouse: []comp.Hint{
			{Key: "rclick", Label: "open it on what is under the pointer"},
			{Key: "click", Label: "choose — and a click on the border is not a click on an action"},
		},
		Roles: []string{"Accent", "Muted", "SelectionFG", "SelectionBG"},
		States: []State{
			{Name: "at a point", Note: "where a right-click landed",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					menu(0).DrawAt(c, r.X+2, r.Y+1)
				}},
			{Name: "moved down", Note: "the key sits beside the action, because it is the same list",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					menu(1).DrawAt(c, r.X+2, r.Y+1)
				}},
			{Name: "nudged back on screen", Note: "opened past the edge; half a menu is a list of actions you cannot read",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					menu(0).DrawAt(c, r.Right()-4, r.Bottom()-1)
				}},
			{Name: "on a region", Note: "the keyboard path — at the thing the cursor is on, wherever that is in THIS frame",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					// Something to anchor to, drawn first, so the menu can find
					// it the way it would find a scrolled row.
					row := comp.Rect{X: r.X + 6, Y: r.Y + 4, W: 24, H: 1}
					c.Fill(row, " ", &s.selected, comp.Region("demo.row"))
					c.Text(row.X+1, row.Y, "api_gateway", &s.selected, comp.Region("demo.row"))
					menu(0).DrawOn(c, comp.Region("demo.row"))
				}},
		},
	}
}

func (m *Model) detailEntry(s *styles) Entry {
	draw := func(d comp.Detail) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			d.TitleStyle, d.SubtitleStyle = &s.title, &s.muted
			d.HeadingStyle, d.LabelStyle = &s.muted, &s.muted
			d.Draw(c, r, comp.Region("demo.detail"))
		}
	}
	full := comp.Detail{
		Title:    "vm-forge-0",
		Subtitle: "Virtual machines",
		Blocks: []comp.Block{
			{Facts: []comp.Fact{
				{Label: "group", Value: "rg-forge"},
				{Label: "location", Value: "uksouth"},
				{Label: "state", Value: "failed", Style: &s.danger},
			}},
			{Heading: "tags", Facts: []comp.Fact{
				{Label: "env", Value: "prod"},
				{Label: "owner", Value: "platform"},
			}},
		},
	}
	return Entry{
		Name:    "Detail",
		Summary: "What a pane says about the one thing you have selected.",
		From:    "democtl field(), azctl fields() — both carried y from call to call",
		Roles:   []string{"Accent", "Muted", "Danger"},
		States: []State{
			{Name: "a thing", Note: "labels line up per BLOCK, so a long tag key does not drag the facts above it wide",
				Draw: draw(full)},
			{Name: "a value with its own colour", Note: "a state that is red in the list and grey here is one fact told twice",
				Draw: draw(comp.Detail{Title: "api_migrate", Blocks: []comp.Block{{Facts: []comp.Fact{
					{Label: "state", Value: "failed", Style: &s.danger},
					{Label: "replicas", Value: "0/1"},
				}}}})},
			{Name: "a secret", Note: "masked is the state you should be in by default — revealing is a keystroke, hiding should not be something you remember",
				Draw: draw(comp.Detail{Title: "ghcr-bot", Blocks: []comp.Block{{Facts: []comp.Fact{
					{Label: "user", Value: "mbp-ci"},
					{Label: "token", Value: "ghp_R3alT0kenValue", Secret: true},
					{Label: "scopes", Value: "read:packages"},
				}}}})},
			{Name: "revealed", Note: "the same fact, shown. Which secrets are showing is the TOOL's — app.Toggles, keyed by a stable id",
				Draw: draw(comp.Detail{Title: "ghcr-bot", Blocks: []comp.Block{{Facts: []comp.Fact{
					{Label: "user", Value: "mbp-ci"},
					{Label: "token", Value: "ghp_R3alT0kenValue"},
					{Label: "scopes", Value: "read:packages"},
				}}}})},
			{Name: "prose", Note: "a note wraps; a fact truncates — losing the end of a sentence loses the point",
				Draw: draw(comp.Detail{
					Title:  "api_migrate",
					Blocks: []comp.Block{{Text: "Pushes a new service spec and waits for the tasks to converge. The current tasks are replaced one at a time."}},
				})},
			{Name: "no room", Note: "it stops at the bottom of its rect rather than writing over the border it sits in",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					draw(full)(c, comp.Rect{X: r.X, Y: r.Y, W: r.W, H: 4}, focused)
				}},
		},
	}
}

func (m *Model) tableEntry(s *styles) Entry {
	svc := [][]string{
		{"✓", "api_gateway", "api", "3/3"},
		{"●", "api_worker", "api", "1/2"},
		{"✗", "api_migrate", "api", "0/1"},
		{"✓", "web_frontend", "web", "4/4"},
	}
	long := [][]string{
		{"✓", "a_service_with_a_very_long_name", "platform", "12/12"},
		{"✓", "short", "web", "1/1"},
	}
	wide := [][]string{
		{"✓", "世界のサービス", "api", "3/3"},
		{"✓", "ascii_service", "api", "3/3"},
	}
	draw := func(tbl comp.Table, rows [][]string) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			tbl.Draw(c, r, rows, &s.muted, comp.Region("demo.table"))
		}
	}
	standard := comp.Table{Gap: 1, Columns: []comp.Column{
		{Width: 1}, {Fill: true}, {Width: 8}, {Width: 5, Right: true},
	}}
	return Entry{
		Name:    "Table",
		Summary: "Rows of aligned columns. It lays out; a List selects and scrolls.",
		From:    "pgctl rows.go, swarmctl pane and detail tables",
		Roles:   []string{"Muted"},
		Glyphs:  []string{"✓", "✗", "●", "…"},
		States: []State{
			{Name: "marker first", Note: "'can I reach it' before 'which is it' — an unreachable one changes everything below",
				Draw: draw(standard, svc)},
			{Name: "a cell too long", Note: "truncated, so the columns after it stay where they are",
				Draw: draw(standard, long)},
			{Name: "wide runes", Note: "measured in columns, or one CJK name knocks every row below it out",
				Draw: draw(standard, wide)},
			{Name: "narrow", Note: "the filler gives its space up first",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					standard.Draw(c, comp.Rect{X: r.X, Y: r.Y, W: min(r.W, 26), H: r.H}, svc,
						&s.muted, comp.Region("demo.table"))
				}},
		},
	}
}

func (m *Model) spinnerEntry(s *styles) Entry {
	// A fixed moment, because a gallery whose frames differ between runs has
	// goldens nobody can review. The harness freezes the clock for the same
	// reason, and a clock-driven spinner is what makes that possible.
	at := time.Date(2026, 8, 31, 9, 14, 3, 0, time.UTC)
	draw := func(sp comp.Spinner, offset time.Duration, label string) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			sp.Style = &s.pending
			x := sp.Draw(c, r, at.Add(offset), comp.Region("demo.spinner"))
			c.Text(r.X+x, r.Y, " "+label, &s.muted, comp.Region("demo.spinner"))
		}
	}
	return Entry{
		Name:    "Spinner",
		Summary: "Work in flight. Its frame comes from the clock, so every spinner turns together.",
		From:    "pgctl spinner(), azctl's single ⠿",
		Roles:   []string{"Pending", "Muted"},
		Glyphs:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏"},
		States: []State{
			{Name: "turning", Note: "the frame is a function of the time, not of how often View ran",
				Draw: draw(comp.Spinner{}, 0, "reading the target's foreign keys…")},
			{Name: "a moment later", Note: "the same spinner, 300ms on — every spinner on screen agrees",
				Draw: draw(comp.Spinner{}, 300*time.Millisecond, "reading the target's foreign keys…")},
			{Name: "the caller's frames", Note: "braille is the default, not a requirement",
				Draw: draw(comp.Spinner{Frames: []string{"-", "\\", "|", "/"}}, 0, "for a font without braille")},
		},
	}
}

func rowsFor(n int) []comp.Row {
	names := []string{
		"api_gateway", "api_worker", "api_migrate", "web_frontend", "web_assets",
		"data_indexer", "data_archiver", "edge_router", "edge_cache", "auth_session",
		"auth_tokens", "media_upload", "media_transcode", "ops_metrics", "ops_logs",
	}
	out := make([]comp.Row, 0, n)
	for i := 0; i < n; i++ {
		out = append(out, comp.Row{Text: " " + names[i%len(names)]})
	}
	return out
}

// groupedRows is the flattened-hierarchy shape: headers with children indented
// under them, which is what two of the four tools build in five places.
func groupedRows(c *comp.Canvas, s *styles, open bool) []comp.Row {
	groups := []struct {
		name     string
		children []string
	}{
		{"rg-forge", []string{"vm-forge-0", "nic-forge-0", "sambpforge"}},
		{"rg-platform", []string{"kv-platform", "app-web", "app-api"}},
		{"rg-data", []string{"sql-reporting", "sadatalake"}},
	}

	ch := c.Chrome()
	var out []comp.Row
	for _, g := range groups {
		lead := ch.Collapsed
		if open {
			lead = ch.Expanded
		}
		out = append(out, comp.Row{Lead: lead, Spans: []comp.Segment{
			{Text: g.name, Style: &s.title},
			{Text: " " + itoa(len(g.children)), Style: &s.muted},
		}})
		if !open {
			continue
		}
		for _, child := range g.children {
			out = append(out, comp.Row{Depth: 1, Text: child})
		}
	}
	return out
}

func (m *Model) listEntry(s *styles) Entry {
	demo := func(rows []comp.Row, focused bool, prepare func(*comp.List)) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, paneFocused bool) {
			l := &comp.List{
				Name: "demo.list", Empty: "  nothing matches",
				Selected: &s.selected, Unfocused: &s.focused,
				Status: &s.muted, EmptyStyle: &s.muted,
				Focused: focused,
			}
			if prepare != nil {
				l.Draw(c, r, rows) // once, so it knows what it can scroll
				prepare(l)
			}
			l.Draw(c, r, rows)
		}
	}
	// The same list, built on demand rather than up front.
	huge := func(n int, at int) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			l := &comp.List{
				Name: "demo.list", Selected: &s.selected, Unfocused: &s.focused,
				Status: &s.muted, Focused: true,
			}
			row := func(i int) comp.Row {
				return comp.Row{Text: fmt.Sprintf("  entry %06d", i)}
			}
			l.DrawFunc(c, r, n, row) // once, so it knows what it can scroll
			l.Select(at)
			l.DrawFunc(c, r, n, row)
		}
	}
	// A query, ranked, with the letters that matched marked.
	filtered := func(query string) func(*comp.Canvas, comp.Rect, bool) {
		commands := []string{
			"switch environment", "edit environments", "env of this service",
			"disk usage", "restart", "follow logs", "capture logs",
			"prune unused images", "open a shell here", "reconnect",
		}
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			hits := fuzzy.Rank(query, commands)
			l := &comp.List{
				Name: "demo.list", Empty: "  nothing matches",
				Selected: &s.selected, Unfocused: &s.focused,
				Status: &s.muted, EmptyStyle: &s.muted, Focused: true,
			}
			l.DrawFunc(c, r, len(hits), func(i int) comp.Row {
				h := hits[i]
				return comp.Row{Spans: comp.Highlight(" "+commands[h.Index], shift(h.At), nil, &s.title)}
			})
		}
	}
	return Entry{
		Name:    "List",
		Summary: "A scrollable, selectable list, flat or grouped. The viewport and the selection are separate.",
		From:    "pgctl window(), swarmctl pane scrolling, azctl cursor/top",
		Keys: []comp.Hint{
			{Key: "↑↓", Label: "move the selection"},
			{Key: "j/k", Label: "the same"},
		},
		Mouse: []comp.Hint{
			{Key: "click", Label: "select the row, anywhere across it"},
			{Key: "wheel", Label: "scroll the viewport, never the selection"},
		},
		Roles:  []string{"SelectionFG", "SelectionBG", "Accent", "Muted"},
		Glyphs: []string{"↑", "↓", "▸", "▾"},
		States: []State{
			{Name: "focused", Note: "the selection is yours to move",
				Draw: demo(rowsFor(6), true, nil)},
			{Name: "unfocused", Note: "the cursor is a memory of where you were",
				Draw: demo(rowsFor(6), false, nil)},
			{Name: "empty", Note: "a filter that matches nothing is an ordinary state",
				Draw: demo(nil, true, nil)},
			{Name: "overflowing", Note: "more rows than the pane; the count says where you are",
				Draw: demo(rowsFor(40), true, nil)},
			{Name: "scrolled away", Note: "the selection is off screen, and says which way",
				Draw: demo(rowsFor(40), true, func(l *comp.List) { l.Scroll(8) })},
			{Name: "grouped", Note: "Depth indents a child; Lead is the header's own glyph, in the marker's column instead of it",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					l := &comp.List{
						Name: "demo.list", Marker: "› ", Blank: "  ",
						Selected: &s.selected, Unfocused: &s.focused,
						Status: &s.muted, Focused: true,
					}
					l.Move(1) // onto a child, so the marker is visible under a header
					l.Draw(c, r, groupedRows(c, s, true))
				}},
			{Name: "grouped, collapsed", Note: "the headers say there is something under them you have not seen",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					l := &comp.List{
						Name: "demo.list", Marker: "› ", Blank: "  ",
						Selected: &s.selected, Unfocused: &s.focused,
						Status: &s.muted, Focused: true,
					}
					l.Draw(c, r, groupedRows(c, s, false))
				}},
			{Name: "no room", Note: "one row; it draws what it can rather than crashing",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					demo(rowsFor(40), true, nil)(c, comp.Rect{X: r.X, Y: r.Y, W: r.W, H: 1}, focused)
				}},
			{Name: "200,000 rows", Note: "built on demand — a frame costs the size of the pane, not the size of the data",
				Draw: huge(200000, 0)},
			{Name: "200,000 rows, deep in", Note: "the window moved; nothing else did",
				Draw: huge(200000, 13742)},
			{Name: "grouped with headings", Note: "↑↓ passes over the headings, so j/k never appear to do nothing — and a click on one is ignored rather than selecting its neighbour",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					l := &comp.List{
						Name: "demo.list", Selected: &s.selected, Unfocused: &s.focused,
						Status: &s.muted, Focused: true,
					}
					rows := []comp.Row{
						{Text: " SERVICES", Skip: true, Style: &s.muted},
						{Text: "   api_gateway"}, {Text: "   api_migrate"},
						{Text: "", Skip: true},
						{Text: " NODES", Skip: true, Style: &s.muted},
						{Text: "   vm-qat-0"}, {Text: "   vm-qat-1"},
					}
					l.Draw(c, r, rows)
					l.Move(2) // over api_migrate, the blank and the heading
					l.Draw(c, r, rows)
				}},
			{Name: "ranked by a query", Note: "the letters that matched are marked, so the order is something you can check rather than trust",
				Draw: filtered("env")},
			{Name: "a query matching nothing", Note: "an ordinary state, not an error",
				Draw: filtered("zzz")},
		},
	}
}

func (m *Model) paneEntry(s *styles) Entry {
	draw := func(p comp.Pane, title string) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, focused bool) {
			p.Border, p.Focus, p.TitleStyle = &s.border, &s.focused, &s.title
			inner := p.Draw(c, comp.Rect{X: r.X, Y: r.Y, W: r.W, H: min(r.H, 6)}, comp.Region("demo.pane"))
			if !inner.Empty() {
				c.Text(inner.X+1, inner.Y, title, &s.muted, comp.Region("demo.pane"))
			}
		}
	}
	return Entry{
		Name:    "Pane",
		Summary: "A bordered box with a title, and the rect inside it.",
		From:    "swarmctl drawBox and drawBoxRaw (which disagreed), democtl box",
		Mouse:   []comp.Hint{{Key: "click", Label: "the border belongs to the pane"}},
		Roles:   []string{"Border", "Accent"},
		Glyphs:  []string{"┌", "─", "┐", "│", "└", "┘"},
		States: []State{
			{Name: "title in the edge", Note: "swarmctl's arrangement; costs no row",
				Draw: draw(comp.Pane{Title: "Services"}, "the room inside")},
			{Name: "title on a row", Note: "democtl's; the focus highlight is the pane's",
				Draw: draw(comp.Pane{Title: "Services", TitleAt: comp.TitleOnRow}, "the room inside")},
			{Name: "focused", Note: "the border and the title change together",
				Draw: draw(comp.Pane{Title: "Services", TitleAt: comp.TitleOnRow, Focused: true}, "the room inside")},
			{Name: "title too long", Note: "truncated, never dropped — you still know what you are looking at",
				Draw: draw(comp.Pane{Title: "a service with a name far longer than this box", TitleAt: comp.TitleOnRow}, "")},
			{Name: "no room", Note: "under two columns it draws nothing at all",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					draw(comp.Pane{Title: "Services"}, "")(c, comp.Rect{X: r.X, Y: r.Y, W: 1, H: 1}, focused)
				}},
			{Name: "rounded", Note: "a chrome the tool chose — ╭╮╰╯ have to be added to its glyph set",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					draw(comp.Pane{Title: "Services", TitleAt: comp.TitleOnRow}, "one field, every box")(
						c.WithChrome(c.Chrome().With(theme.RoundedBox)), r, focused)
				}},
			{Name: "ASCII", Note: "the fallback for a font with nothing — ugly on purpose",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					draw(comp.Pane{Title: "Services", TitleAt: comp.TitleOnRow}, "it should look like a fallback")(
						c.WithChrome(c.Chrome().With(theme.ASCIIBox)), r, focused)
				}},
		},
	}
}

func (m *Model) tabsEntry(s *styles) Entry {
	draw := func(tabs []comp.Tab, active int, focused bool) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, paneFocused bool) {
			comp.Tabs{
				Tabs: tabs, Active: active, Focused: focused,
				Style: &s.muted, Selected: &s.focused, FocusSelected: &s.title, Chrome: &s.muted,
			}.Draw(c, r, "demo.tabs")
		}
	}
	three := []comp.Tab{{Name: "Overview"}, {Name: "Config"}, {Name: "Events"}}
	return Entry{
		Name:    "Tabs",
		Summary: "A strip of names, one of them current. The chevrons say it cycles.",
		From:    "swarmctl tabStrip, pgctl tab bar, democtl",
		Keys:    []comp.Hint{{Key: "‹›", Label: "cycle"}},
		Mouse:   []comp.Hint{{Key: "click", Label: "each tab is its own region"}},
		Roles:   []string{"Accent", "Muted"},
		Glyphs:  []string{"‹", "›", "·"},
		States: []State{
			{Name: "focused", Note: "a strip you can operate looks like one",
				Draw: draw(three, 0, true)},
			{Name: "unfocused", Note: "current, but not yours to change right now",
				Draw: draw(three, 1, false)},
			{Name: "with counts", Note: "a tab can say how much is behind it",
				Draw: draw([]comp.Tab{{Name: "Edits", Count: 3}, {Name: "Log"}, {Name: "Plan", Count: 12}}, 0, true)},
			{Name: "too wide", Note: "clipped by the canvas rather than wrapping",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					draw([]comp.Tab{{Name: "Overview"}, {Name: "Configuration"}, {Name: "Events"}, {Name: "History"}}, 0, true)(
						c, comp.Rect{X: r.X, Y: r.Y, W: min(r.W, 24), H: 1}, focused)
				}},
		},
	}
}

func (m *Model) barEntry(s *styles) Entry {
	return Entry{
		Name:    "Bar",
		Summary: "One line with content at each end. Its Hints are a key and what it does.",
		From:    "swarmctl footerLine, azctl footer, democtl header and footer",
		Roles:   []string{"Accent", "Muted", "Success"},
		Glyphs:  []string{"·"},
		States: []State{
			{Name: "both ends", Note: "the common case",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					comp.Bar{
						Left: []comp.Segment{
							{Text: "democtl", Style: &s.title},
							{Text: "  a tuikit example", Style: &s.muted},
							{Text: "  all services running", Style: &s.success},
						},
						Right: []comp.Segment{{Text: "09:14:03", Style: &s.muted}},
					}.Draw(c, r, comp.Region("demo.bar"))
				}},
			{Name: "hints", Note: "the separator lives in one place, not in every footer string",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					comp.KeyHints(c, r, comp.Region("demo.bar"), &s.muted,
						comp.Hint{Key: "↑↓", Label: "move"}, comp.Hint{Key: "tab", Label: "pane"},
						comp.Hint{Key: "/", Label: "filter"}, comp.Hint{Key: "q", Label: "quit"})
				}},
			{Name: "squeezed", Note: "the left side wins; decoration goes first",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					comp.Bar{
						Left:    []comp.Segment{{Text: "q quit", Style: &s.muted}},
						Right:   []comp.Segment{{Text: "v2.4.1", Style: &s.muted}},
						MinLeft: 24,
					}.Draw(c, comp.Rect{X: r.X, Y: r.Y, W: min(r.W, 28), H: 1}, comp.Region("demo.bar"))
				}},
		},
	}
}

func (m *Model) confirmEntry(s *styles) Entry {
	draw := func(cf comp.Confirm) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			cf.Border, cf.TitleStyle = &s.focused, &s.title
			cf.DangerStyle, cf.BodyStyle, cf.HintStyle = &s.danger, &s.muted, &s.muted
			cf.Hints = []comp.Hint{{Key: "y", Label: "confirm"}, {Key: "n", Label: "cancel"}}
			// Drawn into a canvas of its own so it centres on the preview
			// rather than on the whole window.
			sub := comp.NewCanvas(r.W, r.H)
			cf.Draw(sub, comp.Region("demo.confirm"))
			blit(c, r, sub)
		}
	}
	return Entry{
		Name:    "Confirm",
		Summary: "A question in a box, bounded to its container by construction.",
		From:    "swarmctl action.go, pgctl actionview.go, democtl overlay",
		Keys:    []comp.Hint{{Key: "y", Label: "confirm"}, {Key: "n", Label: "cancel"}},
		Roles:   []string{"Accent", "Danger", "Muted"},
		Glyphs:  []string{"┌", "─", "┐", "│", "└", "┘", "·"},
		States: []State{
			{Name: "ordinary", Note: "grows for its body; the keys never fall off the bottom",
				Draw: draw(comp.Confirm{Title: "Deploy api_gateway?",
					Body: "Pushes a new service spec and waits for the tasks to converge."})},
			{Name: "danger", Note: "an action that destroys something says so in its title",
				Draw: draw(comp.Confirm{Title: "Remove web_assets?", Danger: true,
					Body: "The service and its tasks are removed. This cannot be undone."})},
			{Name: "a long body", Note: "wraps inside the box rather than through the border",
				Draw: draw(comp.Confirm{Title: "Apply the plan?",
					Body: strings.Repeat("The plan's description can be as long as a list of every table a widened selection adds. ", 3)})},
		},
	}
}

func (m *Model) stepListEntry(s *styles) Entry {
	look := [5]comp.StepLook{
		comp.StepWaiting: {Glyph: "•", Style: &s.muted, LabelStyle: &s.muted},
		comp.StepRunning: {Glyph: "→", Style: &s.pending, LabelStyle: &s.pending},
		comp.StepSkipped: {Glyph: "●", Style: &s.muted},
		comp.StepDone:    {Glyph: "✓", Style: &s.success},
		comp.StepFailed:  {Glyph: "✗", Style: &s.danger},
	}
	draw := func(list comp.StepList) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			list.Look, list.Muted = look, &s.muted
			list.Draw(c, r, "demo.steps")
		}
	}
	steps := []comp.Step{
		{Label: "resolve image digest", State: comp.StepDone, Took: "412ms"},
		{Label: "check node capacity", State: comp.StepSkipped, Detail: "already true", Took: "1ms"},
		{Label: "push service spec", State: comp.StepDone, Took: "1.1s"},
	}
	return Entry{
		Name:    "StepList",
		Summary: "A run in progress: what will happen, what has, and what it cost.",
		From:    "swarmctl renderdeploy.go, pgctl actionrun.go, azctl runner.go",
		Keys:    []comp.Hint{{Key: "r", Label: "run again"}, {Key: "esc", Label: "back"}},
		Roles:   []string{"Success", "Danger", "Pending", "Muted"},
		Glyphs:  []string{"✓", "✗", "●", "•", "→"},
		States: []State{
			{Name: "running", Note: "one step in flight, the rest waiting",
				Draw: draw(comp.StepList{Steps: append(append([]comp.Step{}, steps...),
					comp.Step{Label: "wait for tasks to converge", State: comp.StepRunning},
					comp.Step{Label: "verify health endpoint"}),
					Status: "running…", StatusStyle: &s.pending})},
			{Name: "succeeded", Note: "the durations line up in a column you can scan",
				Draw: draw(comp.StepList{Steps: append(append([]comp.Step{}, steps...),
					comp.Step{Label: "wait for tasks to converge", State: comp.StepDone, Took: "8.4s"}),
					Status: "done", StatusStyle: &s.success})},
			{Name: "failed", Note: "a failure owes an explanation, on a line of its own",
				Draw: draw(comp.StepList{Steps: append(append([]comp.Step{}, steps...),
					comp.Step{Label: "wait for tasks to converge", State: comp.StepFailed,
						Note: "task 3 exited 137 before the health check passed"},
					comp.Step{Label: "verify health endpoint"}),
					Status: "stopped", StatusStyle: &s.danger,
					Hints: []comp.Hint{{Key: "r", Label: "to run again"}}})},
			{Name: "no room", Note: "stops at the edge rather than drawing past it",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					draw(comp.StepList{Steps: steps, Status: "running…", StatusStyle: &s.pending})(
						c, comp.Rect{X: r.X, Y: r.Y, W: r.W, H: 2}, focused)
				}},
		},
	}
}

func (m *Model) logPaneEntry(s *styles) Entry {
	lines := func(n int) []comp.LogLine {
		out := make([]comp.LogLine, 0, n)
		for i := 0; i < n; i++ {
			l := comp.LogLine{
				At:   "09:1" + string(rune('0'+i%10)) + ":0" + string(rune('0'+(i*7)%10)),
				Text: "api_gateway: handled request in " + string(rune('1'+i%9)) + "21ms",
			}
			if i%7 == 3 {
				l.Stderr, l.Text = true, "api_gateway: upstream timed out after 3s"
			}
			out = append(out, l)
		}
		return out
	}
	demo := func(n int, prepare func(*comp.LogPane)) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			p := &comp.LogPane{
				Follow: true, Empty: "  no lines on this stream",
				Time: &s.muted, Stderr: &s.stderr, Status: &s.muted, EmptyStyle: &s.muted,
			}
			rows := lines(n)
			if prepare != nil {
				p.Draw(c, r, rows, "demo.logs")
				prepare(p)
			}
			p.Draw(c, r, rows, "demo.logs")
		}
	}
	return Entry{
		Name:    "LogPane",
		Summary: "A stream of lines, tailing. Following is a place, not a mode.",
		From:    "swarmctl logspane.go",
		Keys:    []comp.Hint{{Key: "↑↓", Label: "scroll, detaching from the tail"}},
		Mouse:   []comp.Hint{{Key: "wheel", Label: "the same"}},
		Roles:   []string{"Muted", "Stderr"},
		States: []State{
			{Name: "following", Note: "pinned to the newest line as they arrive",
				Draw: demo(40, nil)},
			{Name: "detached", Note: "scrolled up; it says how far below the tail you are",
				Draw: demo(40, func(p *comp.LogPane) { p.Scroll(-6) })},
			{Name: "empty", Note: "a stderr filter matching nothing",
				Draw: demo(0, nil)},
			{Name: "one line", Note: "no scrolling to do, and it still says so",
				Draw: demo(1, nil)},
		},
	}
}

// blit copies one canvas into a rect of another, which is how a component that
// centres itself is previewed inside a pane rather than on the whole window.
func blit(dst *comp.Canvas, r comp.Rect, src *comp.Canvas) {
	for y := 0; y < src.Bounds().H && y < r.H; y++ {
		for x := 0; x < src.Bounds().W && x < r.W; x++ {
			if cell, ok := src.CellAt(x, y); ok && !cell.Continuation() {
				dst.Set(r.X+x, r.Y+y, cell.Text, cell.Style, cell.Owner)
			}
		}
	}
}

// meterEntry is the only component in this gallery with a pixel layer, and the
// states are arranged to show what that buys: the same value, drawn twice.
func (m *Model) meterEntry(s *styles) Entry {
	draw := func(meters ...comp.Meter) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			y := r.Y
			for i, meter := range meters {
				meter.Filled, meter.Empty, meter.LabelStyle = &s.focused, &s.border, &s.muted
				meter.Track = comp.Name("demo.track" + itoa(i))
				y = meter.Draw(c, comp.Rect{X: r.X, Y: y, W: r.W, H: 1}, comp.Region("demo.meter"))
				y++ // a blank row between them
			}
		}
	}
	return Entry{
		Name:    "Meter",
		Summary: "How far along something is — characters everywhere, pixels where they exist.",
		From:    "swarmctl's run dialog and activity strip, which both wanted it and had neither",
		Roles:   []string{"Accent", "Border", "Muted"},
		Glyphs:  []string{"─", "·"},
		States: []State{
			{Name: "part way", Note: "the bar is ─ and ·; block elements are excluded, so a font without them still draws a bar",
				Draw: draw(comp.Meter{Value: 0.4, Label: "2 of 5"})},
			{Name: "the range", Note: "empty, part, full — a bar that cannot reach either end is a bar nobody trusts",
				Draw: draw(
					comp.Meter{Value: 0, Label: "queued"},
					comp.Meter{Value: 0.62, Label: "measuring images"},
					comp.Meter{Value: 1, Label: "done"},
				)},
			{Name: "no label", Note: "the label is optional; the bar still says the number, coarsely",
				Draw: draw(comp.Meter{Value: 0.75})},
			{Name: "with a pixel layer", Note: "identical here — a test has no terminal, so this is the fallback and always will be",
				Draw: draw(comp.Meter{Value: 0.62, Label: "2 of 5", Pixels: true})},
			{Name: "resolution is the point", Note: "51% and 53% are the same cell bar; run this in foot or Ghostty and they are not",
				Draw: draw(
					comp.Meter{Value: 0.51, Label: "51%", Pixels: true},
					comp.Meter{Value: 0.53, Label: "53%", Pixels: true},
				)},
		},
	}
}

// inputEntry is the states that separate a real text field from an underscore
// stuck on the end of a title.
func (m *Model) inputEntry(s *styles) Entry {
	draw := func(in comp.Input) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, focused bool) {
			in.PromptStyle, in.TextStyle = &s.focused, nil
			in.PlaceholderStyle, in.CursorBG = &s.muted, &s.selected
			in.Focused = in.Focused && focused
			in.Draw(c, comp.Rect{X: r.X, Y: r.Y, W: min(r.W, 44), H: 1}, comp.Region("demo.input"))
		}
	}
	const long = "environments/production/services/api_gateway"

	return Entry{
		Name:    "Input",
		Summary: "One line being typed into, with a caret you can move.",
		From:    "democtl and azctl both fake it with a trailing underscore; the swarmctl palette cannot",
		Keys: []comp.Hint{
			{Key: "←→", Label: "move the caret"},
			{Key: "^a ^e", Label: "start, end"},
			{Key: "^w", Label: "delete the word before it"},
			{Key: "^u", Label: "delete back to the start"},
		},
		Roles:  []string{"Accent", "Muted", "SelectionFG", "SelectionBG"},
		Glyphs: []string{"…"},
		States: []State{
			{Name: "empty", Note: "the placeholder says what the field wants; it is not editable text",
				Draw: draw(comp.Input{Prompt: "> ", Placeholder: "type to search 41 commands"})},
			{Name: "typing", Note: "the caret is a painted cell, so it can sit in the middle of the text",
				Draw: draw(comp.Input{Prompt: "> ", Text: "env", Cursor: 3, Focused: true})},
			{Name: "caret in the middle", Note: "what an underscore on the end cannot do — fixing a typo six back costs six backspaces without it",
				Draw: draw(comp.Input{Prompt: "> ", Text: "production", Cursor: 5, Focused: true})},
			{Name: "longer than the box", Note: "the window follows the caret, and the ellipsis keeps its own column",
				Draw: draw(comp.Input{Prompt: "> ", Text: long, Cursor: 0, Focused: true})},
			{Name: "scrolled to the end", Note: "same text, caret at the end: the view came with it",
				Draw: draw(comp.Input{Prompt: "> ", Text: long, Cursor: len([]rune(long)), Focused: true})},
			{Name: "unfocused", Note: "no caret — an input showing one is claiming keystrokes it will not get",
				Draw: draw(comp.Input{Prompt: "> ", Text: "env", Cursor: 3})},
		},
	}
}

// shift moves match indices right by one, because the rows above are drawn with
// a leading space and the match was made against the text without it.
func shift(at []int) []int {
	out := make([]int, len(at))
	for i, v := range at {
		out[i] = v + 1
	}
	return out
}

// waitingEntry shows the states that separate "busy" from "broken".
func (m *Model) waitingEntry(s *styles) Entry {
	// A fixed moment, so the goldens hold. The spinner is clock-driven, which
	// is what makes freezing the clock enough.
	at := time.Date(2026, 9, 2, 9, 30, 0, 0, time.UTC)
	draw := func(w comp.Waiting, on time.Duration, framed bool) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			w.Style, w.DetailStyle = &s.muted, &s.border
			w.Spinner.Style = &s.pending
			area := r
			if framed {
				// What the component's doc asks for: the interface is drawn,
				// and the wait sits in the hole its contents will fill.
				area = comp.Pane{
					Title: "Resources", Border: &s.border, TitleStyle: &s.title,
				}.Draw(c, r, comp.Region("demo.waitpane"))
			}
			w.Draw(c, area, at.Add(on), comp.Region("demo.waiting"))
		}
	}
	const label = "reading the estate"

	return Entry{
		Name:    "Waiting",
		Summary: "A region whose contents have not arrived. Draw it inside the interface, not instead of it.",
		From:    "azctl's first estate read; swarmctl's connecting screen",
		Roles:   []string{"Muted", "Pending", "Border"},
		Glyphs:  []string{"⠋", "⠙", "⠹", "⠸", "⠼", "⠴", "⠦", "⠧", "⠇", "⠏", "…"},
		States: []State{
			{Name: "in a pane", Note: "the shape everyone should use — framed, so it reads as pending rather than absent",
				Draw: draw(comp.Waiting{Label: label}, 0, true)},
			{Name: "bare", Note: "the same wait with no interface around it. This is what looked like a crash",
				Draw: draw(comp.Waiting{Label: label}, 0, false)},
			{Name: "with a detail", Note: "a wait is more tolerable when it is specific about what it is doing",
				Draw: draw(comp.Waiting{Label: label, Detail: "one Resource Graph query, then grouping is free"}, 0, true)},
			{Name: "under two seconds", Note: "no counter — one that appears and vanishes is a flicker",
				Draw: draw(comp.Waiting{Label: label, Since: at}, 900*time.Millisecond, true)},
			{Name: "taking a while", Note: "past the threshold the count is the difference between working and hung",
				Draw: draw(comp.Waiting{Label: label, Since: at}, 9*time.Second, true)},
			{Name: "a long wait", Note: "minutes, so a number does not run away",
				Draw: draw(comp.Waiting{Label: label, Since: at}, 135*time.Second, true)},
		},
	}
}

// breadcrumbEntry: how you got here, and how to go back.
func (m *Model) breadcrumbEntry(s *styles) Entry {
	draw := func(crumbs []string, width int) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			w := r.W
			if width > 0 {
				w = min(w, width)
			}
			comp.Breadcrumb{
				Crumbs: crumbs, Name: "demo.trail", Item: "demo.crumb",
				Style: &s.muted, Current: &s.title, Separator: &s.border,
			}.Draw(c, comp.Rect{X: r.X, Y: r.Y, W: w, H: 1})
		}
	}
	deep := []string{"democtl", "services", "api_gateway", "Logs"}

	return Entry{
		Name:    "Breadcrumb",
		Summary: "How you got here, and how to go back. app.Stack has kept a Path since it was written.",
		From:    "app.Stack.Path(), which nothing drew until this existed",
		Keys:    []comp.Hint{{Key: "esc", Label: "back one — the crumbs are the mouse's path"}},
		Mouse:   []comp.Hint{{Key: "click", Label: "back to that depth (app.Stack.BackTo)"}},
		Roles:   []string{"Accent", "Muted", "Border"},
		Glyphs:  []string{"·", "…"},
		States: []State{
			{Name: "the root", Note: "one crumb is still a trail — it says you are as far out as you go",
				Draw: draw([]string{"democtl"}, 0)},
			{Name: "two deep", Note: "the last is the screen you are on, and is the only one coloured",
				Draw: draw([]string{"democtl", "Logs"}, 0)},
			{Name: "four deep", Note: "root first, current last",
				Draw: draw(deep, 0)},
			{Name: "too narrow", Note: "elided from the LEFT: where you are and how far in are what a breadcrumb is for; the middle is what you can lose",
				Draw: draw(deep, 26)},
			{Name: "very narrow", Note: "the screen you are on survives alone — dropping it to keep the ones you are not would answer the wrong question",
				Draw: draw(deep, 8)},
		},
	}
}

// scrollbarEntry: where you are, rather than how much is hidden.
func (m *Model) scrollbarEntry(s *styles) Entry {
	// twelve rows in every state, so the bar is the only thing that changes.
	const shown = 12
	draw := func(total, offset int) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			// A list beside it, so the bar is read against something.
			l := &comp.List{Name: "demo.sbrow", Selected: &s.selected, Unfocused: &s.focused}
			body := comp.Rect{X: r.X, Y: r.Y, W: r.W - 2, H: min(r.H, shown)}
			l.Scroll(offset)
			l.DrawFunc(c, body, total, func(i int) comp.Row {
				return comp.Row{Text: " row " + itoa(i)}
			})
			comp.Scrollbar{
				Total: l.Count(), Shown: l.Shown(), Offset: l.Offset(),
				Track: &s.border, Thumb: &s.focused,
			}.Draw(c, comp.Rect{X: r.X + r.W - 1, Y: body.Y, W: 1, H: body.H}, comp.Region("demo.sb"))
		}
	}
	return Entry{
		Name:    "Scrollbar",
		Summary: "Where you are in a list. Its status line says how much is hidden; this says where.",
		From:    "comp.List, which exposed Offset and Max and drew neither",
		Mouse:   []comp.Hint{{Key: "wheel", Label: "the list's, not the bar's — this reports, it does not act"}},
		Roles:   []string{"Border", "Accent"},
		Glyphs:  []string{"·", "│"},
		States: []State{
			{Name: "at the top", Note: "the thumb is a proportion and a position; neither can be worked out from the other",
				Draw: draw(40, 0)},
			{Name: "halfway", Note: "it moves under your hand, which is what makes a viewport feel like one",
				Draw: draw(40, 14)},
			{Name: "at the end", Note: "the last row of the list puts the thumb on the last row of the bar — nearly there and there are different answers",
				Draw: draw(40, 28)},
			{Name: "200,000 rows", Note: "a thumb of 0.0005 rows clamps to one: a scrollbar that vanishes when the list is longest is worse than none",
				Draw: draw(200000, 0)},
			{Name: "nothing to scroll", Note: "draws nothing at all — a full-height thumb says only that there is a scrollbar",
				Draw: draw(8, 0)},
		},
	}
}
