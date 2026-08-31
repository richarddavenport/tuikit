package harness

import (
	tea "github.com/charmbracelet/bubbletea"
)

// Driver is a model the harness can send messages to — tea.Model, essentially,
// but stated here so the harness does not require one.
type Driver interface {
	Model
	Update(tea.Msg) (tea.Model, tea.Cmd)
}

// Press sends keystrokes, so a frame can be reached by the sequence a person
// would type rather than by reaching into the model's fields.
//
// Named keys ("tab", "esc", "enter", the arrows, "backspace") become their key
// types; anything else is sent as runes, so "j" is j and "hello" is five
// characters' worth of one message.
func Press(m Driver, keys ...string) {
	for _, k := range keys {
		m.Update(key(k))
	}
}

// Run sends a key and drains the command it returns, one level deep.
//
// One level, deliberately. A command that returns a command that returns a
// command is a chain the harness cannot know the end of, and a harness that
// loops until quiescent hangs on the first tea.Tick — which is how a step list
// is usually built. A tool that needs a chain driven should drive it itself and
// capture the states it cares about.
func Run(m Driver, k string) {
	_, cmd := m.Update(key(k))
	if cmd == nil {
		return
	}
	if msg := cmd(); msg != nil {
		m.Update(msg)
	}
}

// Resize tells a model the terminal changed, for capturing the same screen at
// several widths.
func Resize(m Driver, w, h int) { m.Update(tea.WindowSizeMsg{Width: w, Height: h}) }

func key(name string) tea.KeyMsg {
	if t, ok := named[name]; ok {
		return tea.KeyMsg{Type: t}
	}
	return tea.KeyMsg{Type: tea.KeyRunes, Runes: []rune(name)}
}

var named = map[string]tea.KeyType{
	"tab":       tea.KeyTab,
	"shift+tab": tea.KeyShiftTab,
	"esc":       tea.KeyEsc,
	"enter":     tea.KeyEnter,
	"space":     tea.KeySpace,
	"backspace": tea.KeyBackspace,
	"delete":    tea.KeyDelete,
	"up":        tea.KeyUp,
	"down":      tea.KeyDown,
	"left":      tea.KeyLeft,
	"right":     tea.KeyRight,
	"home":      tea.KeyHome,
	"end":       tea.KeyEnd,
	"pgup":      tea.KeyPgUp,
	"pgdown":    tea.KeyPgDown,
	"ctrl+c":    tea.KeyCtrlC,
}
