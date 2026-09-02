// Package term is the one place tuikit reads FROM the terminal rather than
// writing to it.
//
// Everything else in this library is a one-way conversation: a canvas is
// serialised and printed, and what the terminal does with it is the terminal's
// business. Graphics detection cannot work that way. There is no environment
// variable that answers "can you draw a picture", because the answer is a
// property of the emulator rather than of the shell that launched it, so the
// only honest way to know is to ask and wait.
//
// Waiting is the whole difficulty. A terminal that does not implement a query
// does not decline it — it says nothing at all, forever. So every read here has
// a deadline, and the deadline expiring is a legitimate answer rather than an
// error: it means [None], and [None] is a mode the interface fully supports.
package term

import (
	"os"
	"strings"
	"time"

	xterm "github.com/charmbracelet/x/term"
)

// Graphics is what a terminal can draw beyond characters.
//
// The three values are not a quality ladder, they are three different
// programs. Sixel replaces the cells it covers; kitty composites with them.
// Code that treats Kitty as "Sixel but nicer" will produce a panel that only
// looks right on one of them.
type Graphics int

const (
	// None is a terminal that draws characters and nothing else — Alacritty,
	// anything unrecognised, anything too slow to answer, and every test.
	//
	// It is first so that the zero value is the safe one. A Graphics nobody
	// set means "draw cells", which is the mode that always works.
	None Graphics = iota
	// Sixel is DEC's protocol, and foot's — the default terminal in Omarchy 4.
	Sixel
	// Kitty is the kitty graphics protocol: kitty, Ghostty, WezTerm, Konsole.
	Kitty
)

func (g Graphics) String() string {
	switch g {
	case Sixel:
		return "sixel"
	case Kitty:
		return "kitty"
	}
	return "none"
}

// EnvOverride is the variable that skips detection: none, sixel or kitty.
//
// It exists for three different people. Someone whose terminal answers a query
// it does not honour needs a way out; someone reproducing a bug needs to pin
// the mode; and someone capturing a frame for documentation needs the cells
// even on a terminal that could do better.
const EnvOverride = "TUIKIT_GRAPHICS"

// DefaultTimeout is how long to wait for a terminal to answer.
//
// Long enough for a transatlantic ssh round trip, short enough that a terminal
// which will never answer does not become a hang at start-up. It is spent once,
// before the first frame, and never again.
const DefaultTimeout = 250 * time.Millisecond

// Detect asks the controlling terminal what it can draw.
//
// The env override wins, then a terminal that is not a terminal — a pipe, a
// test, CI — is None without a query being sent. Everything else is asked.
func Detect() Graphics {
	if g, ok := override(); ok {
		return g
	}
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return None // no controlling terminal: a pipe, or a test
	}
	defer f.Close()
	return Query(f, DefaultTimeout)
}

// override reads EnvOverride. The bool distinguishes "set to none" from unset,
// which matters: the first is a decision and the second is a question.
func override() (Graphics, bool) {
	switch strings.ToLower(strings.TrimSpace(os.Getenv(EnvOverride))) {
	case "none", "off", "cells":
		return None, true
	case "sixel":
		return Sixel, true
	case "kitty":
		return Kitty, true
	}
	return None, false
}

// Query sends both questions at once and reads one reply.
//
// Sending them together is the trick that makes this terminate. The kitty query
// is answered only by terminals that implement it and ignored in silence by
// every other, which on its own would mean waiting out the full timeout on most
// terminals in the world. Primary Device Attributes is answered by essentially
// everything, so appending it gives the read a reliable end: whatever arrived
// before the DA1 reply is the kitty answer, and the DA1 reply itself carries
// the Sixel answer as attribute 4.
//
// The terminal must be raw for this. A cooked terminal buffers by line and the
// reply has no newline in it, so the read would block until the user pressed
// return — which is not a bug anyone finds quickly.
func Query(f *os.File, timeout time.Duration) Graphics {
	if !raw(f) {
		return None
	}
	restore := makeRaw(f)
	if restore == nil {
		return None
	}
	defer restore()

	// \e_Gi=1,a=q\e\\ asks kitty about image id 1; \e[c is DA1.
	if _, err := f.WriteString("\x1b_Gi=1,a=q\x1b\\\x1b[c"); err != nil {
		return None
	}
	return Parse(readUntilDA1(f, timeout))
}

// readUntilDA1 reads until the DA1 reply ends or the deadline passes.
//
// The deadline is on the file, not on a goroutine, so a terminal that answers
// nothing costs exactly one timeout and leaves nothing running behind it.
func readUntilDA1(f *os.File, timeout time.Duration) string {
	if err := f.SetReadDeadline(time.Now().Add(timeout)); err != nil {
		return "" // a platform that cannot deadline a tty gets cells
	}
	defer f.SetReadDeadline(time.Time{}) //nolint:errcheck // best effort

	var b strings.Builder
	buf := make([]byte, 64)
	for b.Len() < 4096 {
		n, err := f.Read(buf)
		b.Write(buf[:n])
		// 'c' terminates a DA1 reply. Stopping on it rather than on the
		// timeout is what makes start-up fast on a terminal that does answer.
		if err != nil || strings.Contains(b.String(), "c") && strings.Contains(b.String(), "\x1b[?") {
			break
		}
	}
	return b.String()
}

// Parse reads a reply. Exported so a test can exercise the decision without a
// terminal — the parsing is where the bugs are, and it should not need a tty
// to be checked.
//
// Kitty wins a tie. Nothing is known to answer both, but if something does, the
// protocol that composites with text is the better one to use.
func Parse(reply string) Graphics {
	if strings.Contains(reply, "\x1b_G") && strings.Contains(reply, "OK") {
		return Kitty
	}
	if da1, ok := cut(reply, "\x1b[?", "c"); ok {
		for _, attr := range strings.Split(da1, ";") {
			if attr == "4" {
				return Sixel
			}
		}
	}
	return None
}

// cut returns what lies between the first open and the next close.
func cut(s, open, close string) (string, bool) {
	i := strings.Index(s, open)
	if i < 0 {
		return "", false
	}
	rest := s[i+len(open):]
	j := strings.Index(rest, close)
	if j < 0 {
		return "", false
	}
	return rest[:j], true
}

// raw reports whether f is a terminal at all. A pipe cannot answer a question,
// and asking one would spend the whole timeout finding that out.
func raw(f *os.File) bool { return xterm.IsTerminal(f.Fd()) }

// makeRaw puts the terminal in raw mode and returns how to undo it, or nil.
//
// Raw is not optional. A cooked terminal buffers by line and a reply carries no
// newline, so the read would block until the user pressed return — a hang that
// looks like a slow start-up and is actually a wait for input nobody knows to
// give.
func makeRaw(f *os.File) func() {
	state, err := xterm.MakeRaw(f.Fd())
	if err != nil {
		return nil
	}
	return func() { xterm.Restore(f.Fd(), state) } //nolint:errcheck // nothing useful to do
}
