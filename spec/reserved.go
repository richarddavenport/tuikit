package spec

// Reserved is the handful of keys whose meaning is fixed across every tuikit
// tool, and the one thing tuikit imposes rather than offers.
//
// # Why impose anything
//
// Everything else here is a part you take or leave. This is a rule, because the
// value is entirely in it being the same everywhere: a reader who has used one
// of these tools presses `q` expecting to leave, and a tool where `q` means
// "queue" has set a trap using the other tools' credibility. One tool can be
// idiosyncratic; a family cannot.
//
// Kept to four, and each one is a key a reader arrives already believing. The
// bar for adding a fifth is that a reader would be SURPRISED to find it meaning
// something else — not that it would be tidy.
//
// # What is reserved is the meaning, not the behaviour
//
// A tool may put a question in front of `q` (decision 39 — azctl asks before
// abandoning a half-run playbook, pgctl cancels immediately, and both are
// right). It may not make `q` mean something that is not leaving.
//
// `ctrl+c` is the exception with no latitude at all: it stops things now. A
// confirmation on the universal escape hatch is a program arguing with the one
// key a reader is entitled to expect works.
var Reserved = map[string]string{
	"ctrl+c": "stop what is happening and quit, immediately and without asking",
	"q":      "quit, or leave the screen you are on",
	"esc":    "go back, or dismiss what is open",
	"?":      "show the keys",
}

// Leaves reports whether a command's name reads as leaving, which is what `q`
// and `esc` are allowed to be bound to.
//
// Matched on the name rather than a flag on Command, because a flag is a second
// thing to keep in step with the name and this is not a close call in practice:
// a command bound to `q` is called quit, exit, back, close or cancel, or it is
// the bug this catches.
func Leaves(name string) bool {
	switch name {
	case "quit", "exit", "back", "close", "cancel", "dismiss", "leave", "stop":
		return true
	}
	return false
}

// Helps reports whether a command's name reads as showing the keys.
func Helps(name string) bool {
	switch name {
	case "help", "keys", "shortcuts", "bindings":
		return true
	}
	return false
}
