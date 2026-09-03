package app

import (
	"time"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/theme"
)

// Model is what a tuikit tool implements instead of tea.Model.
//
// The difference is the last method. TEA's third step is `View() string`, and
// that string is a seam: a tool can hand back a canvas frame, or a hand-joined
// pile of lipgloss, or a canvas frame with something concatenated onto it, and
// nothing can tell the difference. The goldens will happily record whatever
// comes out.
//
// azctl was got to the point where `grep -c 'c.Text(|c.Fill(|c.Set('
// internal/tui/*.go` is zero — the tool computes no coordinates at all. Nothing
// enforced that. It was a habit, and habits are not properties.
//
// Draw makes it one. A model is handed a canvas and a rect and has nowhere else
// to put anything: there is no return value to smuggle a string through. The
// runner owns the canvas, its size, its chrome, and its pixel layer, so those
// stop being four things every tool sets up the same way and gets subtly
// different.
//
// Update returns a Model rather than a tea.Model for the same reason. A
// signature that accepts tea.Model accepts anything with a View, which is the
// escape hatch this exists to close.
type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) (Model, tea.Cmd)
	Draw(c *comp.Canvas, r comp.Rect)
}

// Runner adapts a Model to Bubble Tea, and owns the frame.
//
// It is the only place a tuikit program calls comp.NewCanvas. That is the
// point: the canvas's size, its chrome and its pixel layer are decisions with
// one right answer per program, and a tool that made them itself would make
// them four times.
type Runner struct {
	model  Model
	w, h   int
	canvas *comp.Canvas
	chrome theme.Chrome
	pixels comp.Pixels
	// fullHeight gives the reserved bottom row back to the canvas.
	fullHeight bool
}

// Option configures a Runner.
type Option func(*Runner)

// WithFullHeight draws into every row of the terminal, including the last.
//
// The default keeps one row back. Nobody had written down why (issue 37), and
// the answer turned out to be that nobody decided it: democtl inherited the
// arithmetic from swarmctl and it was then promoted to a rule on the grounds
// that two tools did it "independently".
//
// Measured since, in tmux 3.5a at 24x10 under tea.WithAltScreen: a frame of
// exactly the terminal height, with a bordered pane so the bottom-right cell is
// genuinely written, renders with its top line intact and does not scroll. The
// pending-wrap hazard is real in general and does not fire here, because
// nothing is written after the last cell. azctl also ran full height for months
// before it migrated, with no report of a lost line.
//
// So this is safe as far as anyone has looked, and the default still keeps the
// row: one emulator family has been measured, the failure mode is a top line
// eaten on some OTHER terminal, and that is a bad trade against one row. Turn it
// on for a tool where the row matters and say which terminals you checked.
func WithFullHeight() Option { return func(r *Runner) { r.fullHeight = true } }

// WithChrome sets the characters components draw with.
func WithChrome(ch theme.Chrome) Option { return func(r *Runner) { r.chrome = ch } }

// WithPixels turns on the pixel layer, with what [comp.Detect] answered.
//
// Left unset there is none, which is what every test gets — detection reads
// /dev/tty and a test has no controlling terminal, so goldens cannot pick up a
// developer's terminal by accident.
func WithPixels(p comp.Pixels) Option { return func(r *Runner) { r.pixels = p } }

// WithSize sets the terminal size before the first WindowSizeMsg arrives.
//
// Mostly for capture, where there is no terminal to ask. A running program is
// told within a frame of starting.
func WithSize(w, h int) Option { return func(r *Runner) { r.w, r.h = w, h } }

// New wraps a model.
func New(m Model, opts ...Option) *Runner {
	r := &Runner{model: m, w: 132, h: 38, chrome: theme.DefaultChrome}
	for _, opt := range opts {
		opt(r)
	}
	return r
}

// Init starts the model.
func (r *Runner) Init() tea.Cmd { return r.model.Init() }

// Update passes the message on and keeps the size.
//
// The size is watched here as well as being passed through, because the runner
// is what makes the canvas and cannot ask the model how big it thinks it is.
func (r *Runner) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if size, ok := msg.(tea.WindowSizeMsg); ok {
		r.w, r.h = size.Width, size.Height
	}
	m, cmd := r.model.Update(msg)
	if m != nil {
		r.model = m
	}
	return r, cmd
}

// View builds the frame.
//
// The canvas is one row shorter than the terminal unless [WithFullHeight] says
// otherwise. Decision 38 has the measurement and the argument; the short version
// is that the reserved row is a hedge, not a requirement, and it is kept by
// default because the cost of being wrong is a lost top line and the cost of
// being right is one row.
func (r *Runner) View() string {
	h := r.h - 1
	if r.fullHeight {
		h = r.h
	}
	c := comp.NewCanvas(r.w, max(1, h)).WithChrome(r.chrome).WithGraphics(r.pixels)
	r.model.Draw(c, c.Bounds())
	// Kept so a mouse event can ask what it landed on. The frame IS the region
	// list, so there is nothing else to keep in step with it.
	r.canvas = c
	return c.String()
}

// Canvas is the last frame, so a mouse event or a capture script can address a
// region by name.
func (r *Runner) Canvas() *comp.Canvas { return r.canvas }

// Size is the terminal as the runner understands it.
func (r *Runner) Size() (w, h int) { return r.w, r.h }

// SetSize is what the capture harness calls instead of waiting for a terminal.
//
// It goes through Update rather than setting the fields, so a model that keeps
// its own idea of the size — most do, for layout — hears about it by the same
// route it would in a running program. Two paths to one fact is how they come
// to disagree.
func (r *Runner) SetSize(w, h int) {
	r.Update(tea.WindowSizeMsg{Width: w, Height: h})
}

// Now forwards a fixed clock to a model that has one, so the harness can pin
// time without knowing whether this particular model cares.
func (r *Runner) Now(t time.Time) {
	if c, ok := r.model.(interface{ Now(time.Time) }); ok {
		c.Now(t)
	}
}

// Model is the model underneath, for a test that wants to look at it.
func (r *Runner) Model() Model { return r.model }
