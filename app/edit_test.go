package app

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"
)

func TestEdit(t *testing.T) {
	for _, tc := range []struct {
		name  string
		start string
		msg   tea.KeyMsg
		want  string
		took  bool
	}{
		{"a letter", "fo", tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune("r")}, "for", true},
		// Bubble Tea sets Runes for KeySpace as well as the type. The cloud tool
		// appended both and put two spaces in for every one pressed.
		{"a space goes in once", "rg", tea.KeyMsg{Type: tea.KeySpace, Runes: []rune{' '}}, "rg ", true},
		{"backspace", "rg-", tea.KeyMsg{Type: tea.KeyBackspace}, "rg", true},
		{"backspace on empty is still taken", "", tea.KeyMsg{Type: tea.KeyBackspace}, "", true},
		{"enter is the caller's", "rg", tea.KeyMsg{Type: tea.KeyEnter}, "rg", false},
		{"escape is the caller's", "rg", tea.KeyMsg{Type: tea.KeyEsc}, "rg", false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got, took := Edit(tc.start, tc.msg)
			if got != tc.want || took != tc.took {
				t.Errorf("Edit(%q) = %q, %v; want %q, %v", tc.start, got, took, tc.want, tc.took)
			}
		})
	}
}

// Backspace removes a CHARACTER, not a byte. Cutting a multi-byte rune in half
// leaves a string that is no longer valid UTF-8, and the frame it is drawn
// into shows a replacement mark for something the reader typed correctly.
func TestBackspaceRemovesARune(t *testing.T) {
	got, _ := Edit("café", tea.KeyMsg{Type: tea.KeyBackspace})
	if got != "caf" {
		t.Errorf("got %q, want %q", got, "caf")
	}
}
