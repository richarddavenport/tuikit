package app

import (
	"strings"

	tea "github.com/charmbracelet/bubbletea"
)

// EditAt applies one keystroke to a string AND a caret, and reports whether it
// took the key.
//
// [Edit] is this without a caret: every change happens at the end. That is the
// right shape for a filter you type and clear, and the wrong one for anything
// you correct — a typo six characters back costs six backspaces, and ← does
// nothing at all.
//
// The caret is a RUNE index: 0 is before the first character, len is after the
// last. Bytes would put the caret halfway through a multi-byte character, which
// is not a position, and would then split it on the next backspace.
//
// Enter and escape are not handled here, deliberately and for the same reason
// [Edit] leaves them: they mean accept, run, keep, cancel, or clear depending on
// the tool, and a component guessing is a component getting it wrong somewhere.
func EditAt(s string, cursor int, msg tea.KeyMsg) (string, int, bool) {
	r := []rune(s)
	cursor = clampInt(cursor, 0, len(r))

	switch msg.Type {
	case tea.KeyRunes, tea.KeySpace:
		// Runes already holds the space for KeySpace; appending one as well is
		// the bug Edit's doc comment exists to record.
		in := []rune(string(msg.Runes))
		out := make([]rune, 0, len(r)+len(in))
		out = append(out, r[:cursor]...)
		out = append(out, in...)
		out = append(out, r[cursor:]...)
		return string(out), cursor + len(in), true

	case tea.KeyBackspace:
		if cursor == 0 {
			// Still taken: a backspace at the start of a prompt is not a key
			// for the screen underneath to act on.
			return s, cursor, true
		}
		return string(append(append([]rune{}, r[:cursor-1]...), r[cursor:]...)), cursor - 1, true

	case tea.KeyDelete:
		if cursor >= len(r) {
			return s, cursor, true
		}
		return string(append(append([]rune{}, r[:cursor]...), r[cursor+1:]...)), cursor, true

	case tea.KeyLeft:
		return s, max(0, cursor-1), true
	case tea.KeyRight:
		return s, min(len(r), cursor+1), true
	case tea.KeyHome, tea.KeyCtrlA:
		return s, 0, true
	case tea.KeyEnd, tea.KeyCtrlE:
		return s, len(r), true

	case tea.KeyCtrlU:
		// Everything before the caret, which is what readline does — not the
		// whole line. A tool wanting "clear it all" says so itself.
		return string(r[cursor:]), 0, true

	case tea.KeyCtrlW:
		start := wordStart(r, cursor)
		return string(append(append([]rune{}, r[:start]...), r[cursor:]...)), start, true
	}
	return s, cursor, false
}

// wordStart is where ctrl-w should delete back to: past any spaces immediately
// behind the caret, then past the word before them.
func wordStart(r []rune, cursor int) int {
	i := cursor
	for i > 0 && isSpace(r[i-1]) {
		i--
	}
	for i > 0 && !isSpace(r[i-1]) {
		i--
	}
	return i
}

func isSpace(r rune) bool { return strings.ContainsRune(" \t", r) }

func clampInt(v, lo, hi int) int { return max(lo, min(v, hi)) }
