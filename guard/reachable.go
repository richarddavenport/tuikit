package guard

import (
	"github.com/richarddavenport/tuikit/spec"
)

// Reachable fails when a command can be reached by mouse and not by keyboard.
//
// Four reasons, and the last two decide it. ssh and keyboard-only users. Muscle
// memory. **An agent cannot click.** And a right-click may never arrive at all:
// a multiplexer that captures the mouse gets the event first and the
// application never learns it happened — there is no protocol for asking "did
// my right-click reach you?".
//
// That last one is not theoretical. Running the canvas-mouse prototype inside
// herdr showed HERDR's context menu, not the tool's. herdr's own fix is a
// config setting on the machine running it, which fixes one developer's laptop
// and not the tool. So a command whose only path is a context menu is broken
// for everyone inside herdr or tmux, and the tool cannot detect that to say so.
//
// Unlike the other guards this reads a DECLARATION rather than source: the
// question is about a spec.Command tree, so there is nothing to parse and
// nothing that can hide from it.
func Reachable(t T, root spec.Command) {
	t.Helper()

	for _, name := range spec.Unreachable(root) {
		t.Errorf("%s can be reached by mouse and not by keyboard — an agent cannot click, "+
			"and a multiplexer may eat the right-click before your tool sees it. Give it a Key.", name)
	}
}
