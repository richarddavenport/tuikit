package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/app"
	"github.com/richarddavenport/tuikit/comp"

	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
	"github.com/richarddavenport/tuikit/spec"
	"github.com/richarddavenport/tuikit/theme"
)

// screen is what the window is showing. Every constant here needs a case in
// View — see the exhaustiveness test, which exists because a screen added
// without one renders as an empty terminal and nothing says why.
type screen int

const (
	screenDashboard screen = iota
	screenLogs
	screenRun
)

// pane is which half of the dashboard has focus.
type pane int

const (
	paneList pane = iota
	paneDetail
)

// Model is democtl's whole state.
//
// One struct rather than a model per screen: a screen is a way of looking at the
// fleet, not a separate program, and threading the fleet through three models
// costs more than the branch in View.
type Model struct {
	fleet fleet.Fleet
	sty   styles

	// stack is where you are and how you got there. A field rather than a
	// `screen` — esc used to go to screenDashboard by name, which is right
	// only while every screen is reached from the dashboard.
	stack app.Stack
	focus pane
	tab   int

	// list is the service list's cursor and viewport, which used to be three
	// fields here and the arithmetic to keep them honest.
	list comp.List
	// screens is the router. A screen constant with no entry draws a visible
	// complaint rather than an empty terminal.
	screens app.Screens
	// commands is democtl declared: the CLI, the menu, the keys and the
	// manifest, from one place.
	commands spec.Command

	// filter narrows the list. typing is the mode split: while it is true the
	// list's own keys are text, and only Esc and Enter mean anything else.
	filter string
	typing bool

	// logs. The pane owns the viewport and whether it is following the tail;
	// logStderr is the tool's own filter over what it is handed.
	log       comp.LogPane
	logStderr bool // show stderr lines only

	// run
	plan    fleet.Plan
	done    []stepState
	running int // index of the step in flight, or -1
	// gen invalidates async results from a run that has been abandoned. A step
	// that finishes after you pressed Esc must not draw into the next run — the
	// bug is easy to write and almost impossible to see once written.
	gen app.Gen

	confirm *confirmState
	// menu is the open context menu. Who owns the mouse mid-drag is app.Mouse's
	// business, not a flag here.
	menu  *menuState
	mouse app.Mouse

	// split is the divider between the panes: where it sits, how far it may be
	// dragged, and the gap it lives in.
	split comp.Split

	width, height int
	// The pixel layer, set by main and never here. Detection talks to
	// /dev/tty, and a model that queried in its constructor would ask the
	// developer's real terminal during `go test`.
	pixels comp.Pixels

	// canvas is the last frame drawn, kept so a mouse event can ask what it
	// landed on. The frame is its own region list.
	canvas *comp.Canvas
	// now is frozen by the capture harness so a rendered frame says the same
	// thing tomorrow. Nothing here calls time.Now directly.
	now time.Time
}

type stepState int

const (
	stepWaiting stepState = iota
	stepRunning
	stepSkipped
	stepOK
	stepFailed
)

// confirmState is a modal question. It holds the action rather than a flag, so
// there is no switch elsewhere deciding what "yes" meant.
type confirmState struct {
	title  string
	body   string
	danger bool
	do     func(*Model) tea.Cmd
}

// at is the screen showing, as democtl's own constant. app.Screen is an int so
// a tool can use its own enumeration, and this is the one place the cast lives.
func (m *Model) at() screen { return screen(m.stack.Current()) }

// New builds the model. seed picks the fleet; the same seed is the same fleet.
func New(seed int64) *Model {
	m := &Model{
		fleet:   fleet.New(seed),
		sty:     newStyles(Palette),
		running: -1,
		width:   132,
		height:  38,
		now:     fleet.Epoch,
	}
	// Built once, so "every screen has a view" is a property of the model
	// rather than of whoever last edited a switch.
	m.screens = app.Screens{
		app.Screen(screenDashboard): m.dashboard,
		app.Screen(screenLogs):      m.logs,
		app.Screen(screenRun):       m.run,
	}
	// The root goes on the stack, so back has somewhere to land and nothing
	// has to special-case "the first screen".
	m.stack.Push(app.Screen(screenDashboard), "democtl")
	m.split = comp.Split{Name: regSplit, Ratio: [2]int{1, 3}, Min: minPane}
	m.commands = Commands(seed)
	m.log = comp.LogPane{Follow: true}
	m.list = comp.List{
		Name:       regServicesRow,
		Empty:      "  nothing matches",
		Selected:   &m.sty.selected,
		Unfocused:  &m.sty.focused,
		Status:     &m.sty.muted,
		EmptyStyle: &m.sty.muted,
	}
	return m
}

// SetPalette swaps the vocabulary. A tool would not normally expose this; democtl
// does because tuikit's own tests render it under a changed palette to prove the
// interface follows.
func (m *Model) SetPalette(p theme.Palette) { m.sty = newStyles(p) }

// SetSize is what the capture harness calls instead of waiting for a terminal.
func (m *Model) SetSize(w, h int) { m.width, m.height = w, h }

// Canvas is the last frame drawn, so a script can ask what is where. This is
// what harness.Driver wants, and it is three lines because the frame already
// knows.
func (m *Model) Canvas() *comp.Canvas { return m.canvas }

// SetGraphics turns the pixel layer on, with what the terminal answered.
func (m *Model) SetGraphics(p comp.Pixels) { m.pixels = p }

// Now freezes the clock.
func (m *Model) Now(t time.Time) { m.now = t }

// Init satisfies tea.Model. democtl loads nothing at startup, because there is
// nothing to load — the fleet is already in the model.
func (m *Model) Init() tea.Cmd { return nil }

// visible is the services the filter admits, and the list every index refers to.
func (m *Model) visible() []fleet.Service {
	if m.filter == "" {
		return m.fleet.Services
	}
	var out []fleet.Service
	for _, s := range m.fleet.Services {
		if contains(s.Name, m.filter) || contains(s.Stack, m.filter) {
			out = append(out, s)
		}
	}
	return out
}

// selected is the service under the cursor, and whether there is one. A filter
// that matches nothing is an ordinary state, not an error, so every caller has
// to handle the empty case rather than index into nothing.
func (m *Model) selected() (fleet.Service, bool) {
	list := m.visible()
	if len(list) == 0 {
		return fleet.Service{}, false
	}
	if m.list.Cursor() >= len(list) {
		m.list.Select(len(list) - 1)
	}
	return list[m.list.Cursor()], true
}

func contains(haystack, needle string) bool {
	if needle == "" {
		return true
	}
	for i := 0; i+len(needle) <= len(haystack); i++ {
		if equalFold(haystack[i:i+len(needle)], needle) {
			return true
		}
	}
	return false
}

func equalFold(a, b string) bool {
	if len(a) != len(b) {
		return false
	}
	for i := range len(a) {
		if lower(a[i]) != lower(b[i]) {
			return false
		}
	}
	return true
}

func lower(c byte) byte {
	if c >= 'A' && c <= 'Z' {
		return c + 32
	}
	return c
}
