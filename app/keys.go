package app

import tea "github.com/charmbracelet/bubbletea"

// Handled is what a key handler returns: the command to run, and whether it
// took the key at all. A handler that did not take it lets the next one try.
type Handled func(tea.KeyMsg) (tea.Cmd, bool)

// Keys routes one keystroke, in the only order that works.
//
// # The order is the contract
//
// Whatever is CAPTURING keys gets them first, then the screen, then the
// globals. A global handled before a modal is a modal you cannot type into —
// press q in a filter box and the program exits.
//
// # Why this is a type
//
// Every one of the four tools has this split, and every one writes it as an
// early return somebody has to remember. democtl's capturesKeys is documented
// as its central convention and is still one `if` that could be deleted without
// anything failing to compile. Here the order is not something you write; it is
// the shape of the struct, and a Capture handler that is nil is a tool with
// nothing capturing rather than a tool that forgot.
type Keys struct {
	// Capture is whatever is eating keystrokes right now — a modal, a filter
	// being typed, an open menu. Nil when nothing is. While it is set, j is
	// the letter j.
	Capture Handled
	// Screen is the current screen's own keys.
	Screen Handled
	// Global is quit, escape, and anything else that works everywhere. It runs
	// last, and only if nothing above took the key.
	Global Handled
}

// Route sends the key through the three in order and returns the command from
// whichever took it.
func (k Keys) Route(msg tea.KeyMsg) tea.Cmd {
	// A capture takes the key whether or not it does anything with it. That is
	// the point: an unrecognised key inside a filter box is a character, not a
	// chance for the screen underneath to act on it.
	if k.Capture != nil {
		cmd, _ := k.Capture(msg)
		return cmd
	}
	if k.Screen != nil {
		if cmd, took := k.Screen(msg); took {
			return cmd
		}
	}
	if k.Global != nil {
		cmd, _ := k.Global(msg)
		return cmd
	}
	return nil
}
