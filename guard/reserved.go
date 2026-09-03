package guard

import (
	"strings"

	"github.com/richarddavenport/tuikit/spec"
)

// Reserved fails when a command takes a key that means something else
// everywhere.
//
// The one rule tuikit imposes rather than offers. [spec.Reserved] says which
// keys and why; the short version is that a reader who has used one of these
// tools presses `q` expecting to leave, and a tool where `q` means "queue" has
// set a trap using the other tools' credibility.
//
// # What this catches, and what it cannot
//
// It reads DECLARED commands — the spec tree, which is the same list that feeds
// the CLI, `describe --json` and the context menus. A key bound there is caught
// before it ships.
//
// It does not see a key handled directly in a model's Update switch, and no
// static check can: `case "q":` is the same source whether the branch quits or
// queues, and telling them apart means understanding what the branch does.
// Decision 40's rule applies — a check that catches part of a class is worth
// having if it says which part, and the part it misses is covered by the fact
// that a tool's screen keys are read by whoever writes the help screen.
func Reserved(t T, root spec.Command) {
	t.Helper()

	for _, cmd := range spec.WithKeys(root) {
		meaning, isReserved := spec.Reserved[cmd.Key]
		if !isReserved {
			continue
		}
		// WithKeys reports the full path ("democtl quit"), and what a command
		// MEANS is its own name rather than its parents'.
		leaf := cmd.Name
		if i := strings.LastIndex(leaf, " "); i >= 0 {
			leaf = leaf[i+1:]
		}

		ok := spec.Leaves(leaf)
		if cmd.Key == "?" {
			ok = spec.Helps(leaf)
		}
		if cmd.Key == "ctrl+c" {
			// Nothing may take it, not even a quit: it is handled by the
			// runtime before a tool sees it, and a tool that declares it is
			// describing a binding it does not have.
			ok = false
		}
		if !ok {
			t.Errorf("%q is bound to %q, which is reserved: %s. "+
				"A reader who has used another of these tools will press it expecting that. "+
				"Bind %q to another key, or rename it if it really does leave.",
				cmd.Name, cmd.Key, meaning, cmd.Name)
		}
	}
}
