package gallery

import (
	"strings"
	"time"

	"github.com/richarddavenport/tuikit/comp"
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
		m.toastEntry(s),
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

func (m *Model) listEntry(s *styles) Entry {
	demo := func(rows []comp.Row, focused bool, prepare func(*comp.List)) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, paneFocused bool) {
			l := &comp.List{
				Name: "demo.list", Empty: "  nothing matches",
				Selected: &s.selected, Unfocused: &s.focused,
				Status: &s.muted, EmptyStyle: &s.muted,
				Focused: focused && paneFocused,
			}
			if prepare != nil {
				l.Draw(c, r, rows) // once, so it knows what it can scroll
				prepare(l)
			}
			l.Draw(c, r, rows)
		}
	}
	return Entry{
		Name:    "List",
		Summary: "A scrollable, selectable list. The viewport and the selection are separate.",
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
		Glyphs: []string{"↑", "↓"},
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
			{Name: "no room", Note: "one row; it draws what it can rather than crashing",
				Draw: func(c *comp.Canvas, r comp.Rect, focused bool) {
					demo(rowsFor(40), true, nil)(c, comp.Rect{X: r.X, Y: r.Y, W: r.W, H: 1}, focused)
				}},
		},
	}
}

func (m *Model) paneEntry(s *styles) Entry {
	draw := func(p comp.Pane, title string) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, focused bool) {
			p.Border, p.Focus, p.TitleStyle = &s.border, &s.focused, &s.title
			p.Focused = p.Focused && focused
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
		},
	}
}

func (m *Model) tabsEntry(s *styles) Entry {
	draw := func(tabs []comp.Tab, active int, focused bool) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, paneFocused bool) {
			comp.Tabs{
				Tabs: tabs, Active: active, Focused: focused && paneFocused,
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
