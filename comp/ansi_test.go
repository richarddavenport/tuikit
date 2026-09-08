package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

func texts(segs []Segment) string {
	var b strings.Builder
	for _, s := range segs {
		b.WriteString(s.Text)
	}
	return b.String()
}

// The case from the issue.
func TestANSITurnsSGRIntoSegments(t *testing.T) {
	got := ANSI("\x1b[32mReady\x1b[0m")
	if len(got) != 1 {
		t.Fatalf("got %d segments, want 1: %#v", len(got), got)
	}
	if got[0].Text != "Ready" {
		t.Errorf("text %q, want Ready", got[0].Text)
	}
	if got[0].Style == nil {
		t.Fatal("the green was dropped")
	}
	if fg := got[0].Style.GetForeground(); fg != lipgloss.ANSIColor(2) {
		t.Errorf("foreground is %v, want ANSI 2 — the reader's own green", fg)
	}
}

// Unstyled text carries no style at all, so it inherits whatever the pane sets
// rather than being pinned to a default.
func TestPlainTextHasNoStyle(t *testing.T) {
	got := ANSI("just words")
	if len(got) != 1 || got[0].Style != nil {
		t.Errorf("plain text produced %#v, want one unstyled segment", got)
	}
}

// A reset returns to unstyled, and the two runs do not merge.
func TestAResetEndsTheRun(t *testing.T) {
	got := ANSI("\x1b[31mbad\x1b[0m fine")
	if len(got) != 2 {
		t.Fatalf("got %d segments, want 2: %#v", len(got), got)
	}
	if got[0].Text != "bad" || got[1].Text != " fine" {
		t.Errorf("split as %q / %q", got[0].Text, got[1].Text)
	}
	if got[1].Style != nil {
		t.Error("the text after the reset kept a style")
	}
}

// Style carries across a newline, the way a terminal does it. A parser that
// reset per line would color one line of three.
func TestStyleCarriesAcrossLines(t *testing.T) {
	lines := ANSILines("\x1b[34mone\ntwo\nthree\x1b[0m")
	if len(lines) != 3 {
		t.Fatalf("got %d lines, want 3", len(lines))
	}
	for i, line := range lines {
		if len(line) == 0 || line[0].Style == nil {
			t.Errorf("line %d lost the color: %#v", i, line)
		}
	}
}

// Carriage return goes back and overwrites, which is how every progress bar is
// drawn. Ten redraws must come out as the final state, not ten states.
func TestCarriageReturnOverwrites(t *testing.T) {
	got := ANSI("loading 10%\rloading 90%")
	if s := texts(got); s != "loading 90%" {
		t.Errorf("got %q, want the final state", s)
	}
}

// A shorter redraw leaves the tail of the longer one, exactly as a terminal
// does — the tool sees what a person would have seen.
func TestCarriageReturnLeavesWhatItDidNotCover(t *testing.T) {
	got := ANSI("aaaaaaa\rbb")
	if s := texts(got); s != "bbaaaaa" {
		t.Errorf("got %q, want bbaaaaa", s)
	}
}

func TestBackspaceGoesBackOneColumn(t *testing.T) {
	if s := texts(ANSI("abc\bX")); s != "abX" {
		t.Errorf("got %q, want abX", s)
	}
}

func TestTabGoesToTheNextStop(t *testing.T) {
	if s := texts(ANSI("ab\tc")); s != "ab      c" {
		t.Errorf("got %q, want ab then padding to column 8", s)
	}
}

// A cursor move is skipped whole. Half an escape sequence on screen is worse
// than none of it, and this is the line decision 27 draws.
func TestNonSGRSequencesAreSkippedNotPrinted(t *testing.T) {
	got := texts(ANSI("a\x1b[2Jb\x1b[10;20Hc"))
	if got != "abc" {
		t.Errorf("got %q, want abc with the cursor moves dropped", got)
	}
}

// 256-color and truecolor survive, because a subprocess's output is data and
// dropping the distinction loses what the program was drawing.
func TestExtendedColorsSurvive(t *testing.T) {
	if got := ANSI("\x1b[38;5;208mamber\x1b[0m"); got[0].Style == nil {
		t.Error("256-color was dropped")
	} else if fg := got[0].Style.GetForeground(); fg != lipgloss.Color("208") {
		t.Errorf("256-color became %v, want 208", fg)
	}
	if got := ANSI("\x1b[38;2;255;128;0mamber\x1b[0m"); got[0].Style == nil {
		t.Error("truecolor was dropped")
	} else if fg := got[0].Style.GetForeground(); fg != lipgloss.Color("#ff8000") {
		t.Errorf("truecolor became %v, want #ff8000", fg)
	}
}

// Neighboring cells with one style come out as one segment, so a line of one
// color is not eighty of them.
func TestARunOfOneColorIsOneSegment(t *testing.T) {
	if got := ANSI("\x1b[31mabcdefghij\x1b[0m"); len(got) != 1 {
		t.Errorf("a ten-character run became %d segments", len(got))
	}
}

// Bold and color compose, and 22 turns bold off without losing the color.
func TestAttributesCompose(t *testing.T) {
	got := ANSI("\x1b[1;31mboth\x1b[22mjust color")
	if len(got) != 2 {
		t.Fatalf("got %d segments, want 2", len(got))
	}
	if !got[0].Style.GetBold() {
		t.Error("bold was lost")
	}
	if got[1].Style == nil || got[1].Style.GetBold() {
		t.Error("22 should end bold and keep the color")
	}
	if fg := got[1].Style.GetForeground(); fg != lipgloss.ANSIColor(1) {
		t.Errorf("the color did not survive 22: %v", fg)
	}
}

// Nothing in, nothing out — and it must not panic on a truncated escape.
func TestTruncatedInputDoesNotPanic(t *testing.T) {
	for _, s := range []string{"", "\x1b", "\x1b[", "\x1b[3", "\x1b[38;5", "\x1b[38;2;1"} {
		ANSI(s)
	}
}
