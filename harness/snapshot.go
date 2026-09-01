package harness

import (
	"fmt"
	"os"
	"strings"
)

// Snapshot drives a model through a script and writes the frames — from the
// BINARY, with no test involved.
//
// # Why there are two capture mechanisms
//
// Decision 10, and neither replaces the other. A test helper reaches any state
// at all — an error, a half-loaded pane, a modal over a forty-row table,
// something no key sequence gets to yet — because the model's fields are
// unexported and `go test` is the only thing that can reach inside the package.
// This reaches only what a keystroke reaches, and needs no test.
//
// The difference that matters is discoverability. `MYTOOL_FRAMES=/tmp/x go test
// ./internal/tui -run CaptureFrames` requires knowing an environment variable,
// a package and a test name — three facts that live in the source. An agent
// that wants to see what it just built should be able to find out how from
// --help, and this is how.
//
// # The script
//
//	press E j j
//	shot deep-in-the-list
//	click inventory.row[3]
//	shot a-resource-selected
//
// `shot` is the verb that makes this a capture rather than a keystroke replay:
// the script says WHERE the interesting frames are, so one invocation produces
// the set worth looking at. A script with no shot in it captures one frame at
// the end, because a caller who asked for a snapshot and gets nothing has been
// told nothing.
//
// Lines are separated by newlines or semicolons, so it survives being a shell
// argument.
func Snapshot(m Pointer, dir, script string, opts ...Option) (frames []Frame, err error) {
	if err := os.MkdirAll(dir, 0o755); err != nil {
		return nil, err
	}

	rec := &collector{}
	// The Session reports write failures through T, which is a testing
	// interface. Fatalf panics with a sentinel so a binary gets an error
	// instead of a dead process.
	defer func() {
		if p := recover(); p != nil {
			if _, ok := p.(stop); !ok {
				panic(p)
			}
			frames, err = nil, rec.err
		}
	}()

	s := Capture(rec, dir, opts...)
	shots := 0
	for i, line := range lines(script) {
		verb, rest, _ := strings.Cut(line, " ")
		if verb == "shot" {
			name := strings.TrimSpace(rest)
			if name == "" {
				return nil, fmt.Errorf("line %d: shot needs a name", i+1)
			}
			s.Shot(name, m)
			shots++
			continue
		}
		if err := step(m, line); err != nil {
			return nil, fmt.Errorf("line %d: %s: %w", i+1, line, err)
		}
	}
	if shots == 0 {
		s.Shot("screen", m)
	}
	return s.Done(), rec.err
}

// lines splits a script on newlines and semicolons, dropping blanks and
// comments — so the same text works in a file and as a shell argument.
func lines(script string) []string {
	var out []string
	for _, part := range strings.FieldsFunc(script, func(r rune) bool { return r == '\n' || r == ';' }) {
		if part = strings.TrimSpace(part); part != "" && !strings.HasPrefix(part, "#") {
			out = append(out, part)
		}
	}
	return out
}

// ScriptFile reads a script from a path, or returns the text unchanged when it
// is not one — so --script takes either without a second flag to say which.
func ScriptFile(s string) (string, error) {
	if strings.ContainsAny(s, "\n;") || strings.Contains(s, " ") {
		return s, nil
	}
	body, err := os.ReadFile(s)
	if err != nil {
		return "", err
	}
	return string(body), nil
}

// collector is a T for a binary: it records rather than failing a test.
type collector struct{ err error }

type stop struct{}

func (c *collector) Helper() {}
func (c *collector) Errorf(format string, args ...any) {
	if c.err == nil {
		c.err = fmt.Errorf(format, args...)
	}
}
func (c *collector) Fatalf(format string, args ...any) {
	c.Errorf(format, args...)
	panic(stop{})
}
