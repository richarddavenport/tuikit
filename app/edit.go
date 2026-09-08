package app

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Edit applies one keystroke to a string being typed and reports whether it
// took the key.
//
// # Why this is not four lines in each tool
//
// It was four lines in each tool, and one of them was wrong. The cloud tool's
// filter had:
//
//	case tea.KeyRunes, tea.KeySpace:
//		m.filter += string(msg.Runes)
//		if msg.Type == tea.KeySpace {
//			m.filter += " "
//		}
//
// Bubble Tea sets Runes to []rune{' '} for KeySpace as well as setting the
// type, so a space went in twice. Nobody noticed, because no captured frame
// types a space and nobody filters an Azure estate by a phrase.
//
// That is the whole argument for this being a function: it is small enough
// that everyone writes it, and just fiddly enough that somebody gets it wrong.
//
// What it does NOT handle is enter and escape, which mean different things in
// different tools — keep, accept, cancel, clear the value, close the prompt.
// Those stay with the caller, so nothing has to guess.
func Edit(s string, msg tea.KeyMsg) (string, bool) {
	switch msg.Type {
	case tea.KeyBackspace:
		if s == "" {
			// Still taken. A backspace on an empty prompt is not a key for the
			// screen underneath to act on.
			return s, true
		}
		// Runes, not bytes: a filter containing a non-ASCII character would
		// otherwise be cut in half and stop being valid UTF-8.
		r := []rune(s)
		return string(r[:len(r)-1]), true

	case tea.KeySpace:
		// Runes already holds the space. Appending one AS WELL is the bug this
		// function exists to stop happening again.
		return s + string(msg.Runes), true

	case tea.KeyRunes:
		return s + string(msg.Runes), true
	}
	return s, false
}
