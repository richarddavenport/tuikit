package ui

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/comp"

	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
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

	screen screen
	focus  pane
	cur    int // index into visible()
	tab    int

	// filter narrows the list. typing is the mode split: while it is true the
	// list's own keys are text, and only Esc and Enter mean anything else.
	filter string
	typing bool

	// logs
	logOffset int
	logStderr bool // show stderr lines only

	// run
	plan    fleet.Plan
	done    []stepState
	running int // index of the step in flight, or -1
	// gen invalidates async results from a run that has been abandoned. A step
	// that finishes after you pressed Esc must not draw into the next run — the
	// bug is easy to write and almost impossible to see once written.
	gen int

	confirm *confirmState

	width, height int
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

// New builds the model. seed picks the fleet; the same seed is the same fleet.
func New(seed int64) *Model {
	return &Model{
		fleet:   fleet.New(seed),
		sty:     newStyles(Palette),
		running: -1,
		width:   132,
		height:  38,
		now:     fleet.Epoch,
	}
}

// SetPalette swaps the vocabulary. A tool would not normally expose this; democtl
// does because tuikit's own tests render it under a changed palette to prove the
// interface follows.
func (m *Model) SetPalette(p theme.Palette) { m.sty = newStyles(p) }

// SetSize is what the capture harness calls instead of waiting for a terminal.
func (m *Model) SetSize(w, h int) { m.width, m.height = w, h }

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
	if m.cur >= len(list) {
		m.cur = len(list) - 1
	}
	return list[m.cur], true
}

// capturesKeys reports that something on screen is eating keystrokes.
//
// The split that every screen's key handling hangs off: while a modal is open or
// the filter is being typed, j is the letter j. Without this the list scrolls
// behind a dialog, which is the kind of bug that survives review because nobody
// tries it.
func (m *Model) capturesKeys() bool { return m.typing || m.confirm != nil }

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
