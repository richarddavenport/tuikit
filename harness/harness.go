// Package harness renders a tool's screens to files, so a person or an agent
// can look at them.
//
// This is the framework's headline affordance rather than a testing utility.
// Capture exists because someone building a screen wants to see it; goldens,
// documentation and catching overflow are all downstream of that, and treating
// any of them as the point produces a harness nobody runs.
//
// It generalises a working prototype — the database tool's
// screenshot_probe_test.go, whose seventeen frames found three bugs that
// forty-odd assertions had not. Two of the three were something drawn wider
// than its container, which no assertion checking content rather than shape
// can see.
//
// # Two modes, and both are needed
//
// Fixture mode builds the model by hand: deterministic, backendless, and what
// goldens are made of. Live mode fills it from a real backend: documentation,
// and the thing that finds what a fixture hides.
//
// The second is not a luxury. The database tool's header-overflow bug was
// masked in its first capture because the fixture shortened a config path to
// "the database tool.yaml"; only a real temp-dir path was long enough to
// overflow. A fixture encodes the author's assumptions, which is exactly what
// a capture is meant to catch.
package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"sort"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// T is the part of *testing.T the harness uses.
//
// Narrower than testing.TB for the same reason guard.T is: so the harness can
// be tested against a recorder. *testing.T satisfies it.
type T interface {
	Helper()
	Errorf(format string, args ...any)
	Fatalf(format string, args ...any)
}

// Model is anything that can draw itself. Narrower than tea.Model on purpose —
// capturing a frame needs a View and nothing else.
type Model interface{ View() string }

// Sizer is a model that can be told how big the terminal is. Implement it and
// the harness sizes every frame; otherwise it captures whatever size the model
// already believes in.
type Sizer interface{ SetSize(w, h int) }

// Clock is a model whose sense of time can be fixed.
//
// Worth implementing. A frame that reads "47s ago" has to read that tomorrow,
// or every golden fails the day after it is written.
type Clock interface{ Now(time.Time) }

// Mode records where a frame's data came from.
//
// A page of frames that quietly mixes the two misreports the tool, so the
// distinction is carried through to the manifest rather than left to a
// caption.
type Mode string

const (
	// Fixture is a model built by hand: deterministic, no backend.
	Fixture Mode = "fixture"
	// Live is a model filled from a real backend moments before capture.
	Live Mode = "live"
	// Composed is a frame whose data could not be obtained and was staged — the
	// database tool marked its snapshot-manifest frames this way, because no
	// production snapshot existed on that machine to photograph.
	Composed Mode = "composed"
)

// Frame is one captured screen.
type Frame struct {
	Name  string `json:"name"`
	File  string `json:"file"`
	Mode  Mode   `json:"mode"`
	Width int    `json:"width"`
	Lines int    `json:"lines"`
}

// Session writes a run of frames into a directory.
type Session struct {
	t      T
	dir    string
	w, h   int
	at     time.Time
	mode   Mode
	frames []Frame
}

// Option configures a Session.
type Option func(*Session)

// Size fixes the terminal. A frame is reproducible only at a known width.
func Size(w, h int) Option { return func(s *Session) { s.w, s.h = w, h } }

// At freezes the clock.
func At(t time.Time) Option { return func(s *Session) { s.at = t } }

// WithMode sets what every frame in this session is, unless a shot says
// otherwise.
func WithMode(m Mode) Option { return func(s *Session) { s.mode = m } }

// Enabled returns the directory named by an environment variable, or "".
//
// Capture is gated so a plain `go test ./...` skips it: rendering frames is
// something you ask for. The caller does the skipping explicitly rather than
// the harness no-opping, because a capture that silently does nothing is the
// same failure as a guard that scans nothing.
func Enabled(env string) string { return os.Getenv(env) }

// Capture starts a session writing into dir.
//
// It forces lipgloss's colour profile, which is the trick that makes any of
// this work: a test has no TTY, so lipgloss strips every colour and you
// capture a grey rectangle. This does not contradict the palette being ANSI
// 256 — the values stay 256 indices; forcing the profile only stops them being
// discarded on a pipe.
func Capture(t T, dir string, opts ...Option) *Session {
	t.Helper()

	if dir == "" {
		t.Fatalf("harness: no directory — capture into somewhere, or skip explicitly")
	}
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("harness: %v", err)
	}

	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)

	s := &Session{t: t, dir: dir, w: 132, h: 38, at: time.Unix(0, 0).UTC(), mode: Fixture}
	for _, opt := range opts {
		opt(s)
	}
	return s
}

// Shot renders the model and writes it as <name>.ansi.
func (s *Session) Shot(name string, m Model) { s.t.Helper(); s.shot(name, m, s.mode) }

// ShotAs renders one frame with a mode of its own — for a screen that had to
// be staged in an otherwise live run.
func (s *Session) ShotAs(name string, m Model, mode Mode) { s.t.Helper(); s.shot(name, m, mode) }

func (s *Session) shot(name string, m Model, mode Mode) {
	s.t.Helper()

	if sz, ok := m.(Sizer); ok {
		sz.SetSize(s.w, s.h)
	}
	if c, ok := m.(Clock); ok {
		c.Now(s.at)
	}

	view := m.View()
	file := name + ".ansi"
	if err := os.WriteFile(filepath.Join(s.dir, file), []byte(view), 0o600); err != nil {
		s.t.Fatalf("harness: %v", err)
	}
	s.frames = append(s.frames, Frame{
		Name:  name,
		File:  file,
		Mode:  mode,
		Width: Width(view),
		Lines: len(Lines(view)),
	})
}

// Frames returns what has been captured, in order.
func (s *Session) Frames() []Frame { return s.frames }

// Done writes the manifest and returns the frames.
//
// The manifest is what docgen reads to build a page and what an agent reads to
// know which frames exist without listing a directory and guessing at names.
func (s *Session) Done() []Frame {
	s.t.Helper()

	sorted := append([]Frame(nil), s.frames...)
	sort.Slice(sorted, func(i, j int) bool { return sorted[i].Name < sorted[j].Name })

	body, err := json.MarshalIndent(struct {
		Width  int     `json:"width"`
		Height int     `json:"height"`
		Frames []Frame `json:"frames"`
	}{s.w, s.h, sorted}, "", "  ")
	if err != nil {
		s.t.Fatalf("harness: %v", err)
	}
	if err := os.WriteFile(filepath.Join(s.dir, "frames.json"), append(body, '\n'), 0o600); err != nil {
		s.t.Fatalf("harness: %v", err)
	}
	return s.frames
}

// Dir is where the session is writing.
func (s *Session) Dir() string { return s.dir }

func (s *Session) String() string {
	return fmt.Sprintf("harness: %d frames at %dx%d in %s", len(s.frames), s.w, s.h, s.dir)
}
