package spec

import "github.com/richarddavenport/tuikit/comp"

// MenuFor is the context menu for a region: the commands that target it.
//
// # Why this is the fourth surface and not a fifth list
//
// A hand-maintained menu drifts from the CLI within a release. Built from the
// same declarations, the menu for a region IS the commands whose Target matches
// it — so a command added to the CLI appears in the menu, and one removed
// disappears, without anybody remembering either.
//
// The hint carries the KEY as well as the label, because comp.Hint is what the
// footer is built from too. The menu and the keymap are one list.
func MenuFor(root Command, region comp.Name) []comp.Hint {
	var out []comp.Hint
	walk(root, nil, func(cmd Command, path []string) {
		if cmd.Hidden || cmd.Target != region {
			return
		}
		label := cmd.Short
		if label == "" {
			label = cmd.Name
		}
		out = append(out, comp.Hint{Key: cmd.Key, Label: label})
	})
	return out
}

// Unreachable is every command that can be reached by mouse and not by
// keyboard.
//
// # Why this is a hard rule
//
// Four reasons, and the last two decide it. ssh and keyboard-only users. Muscle
// memory. An AGENT CANNOT CLICK. And a right-click may never arrive at all: a
// multiplexer that captures the mouse gets the event first and the application
// never learns it happened — there is no protocol for asking "did my
// right-click reach you?". So a command whose only path is a context menu is
// simply broken for anyone inside herdr or tmux, and the tool cannot detect
// that to warn them.
//
// This is the check behind guard.Reachable. It is here rather than in guard
// because it is a question about a declaration rather than about source code:
// there is nothing to parse.
func Unreachable(root Command) []string {
	var out []string
	walk(root, nil, func(cmd Command, path []string) {
		if cmd.Target != "" && cmd.Key == "" {
			out = append(out, join(path))
		}
	})
	return out
}

// Keyed is a command that has a keyboard path, and the path.
type Keyed struct {
	Name string
	Key  string
}

// WithKeys is every command in the tree that has a Key.
//
// For guard.Keys, which asks whether the help screen names them all. Separate
// from Unreachable because that asks the opposite question — a command with a
// Target and no Key — and a function answering both would be answering neither
// clearly.
func WithKeys(root Command) []Keyed {
	var out []Keyed
	walk(root, nil, func(cmd Command, path []string) {
		if cmd.Key != "" {
			out = append(out, Keyed{Name: join(path), Key: cmd.Key})
		}
	})
	return out
}
