// Package spec is one declaration and the four surfaces it produces.
//
// A command is declared once and becomes: the CLI command, the TUI screen, the
// entry an agent reads in `tool describe --json`, and the right-click menu for
// the regions it targets. None of the four can drift from the others, because
// there is only one of them.
//
// That is the idea the whole framework is organised around. Four tools have
// four copies of the same command described four times — in a cobra tree, in a
// keymap, in a footer string, and in a README — and the copies disagree. Not
// dramatically: a flag renamed in one place, a key that does something slightly
// different, a README a version behind. Enough that an agent reading any one of
// them is reading a lie.
//
// # No cobra
//
// Decision 12. Once spec is the source of truth, cobra is a second description
// of the same tree — and adapting into it costs the exit-code contract, because
// cobra wants to own os.Exit. Nothing here does: Run returns a code and main
// decides. The CLI is built on stdlib flag.
package spec

import "github.com/richarddavenport/tuikit/comp"

// Kind is what a flag holds.
type Kind int

// The kinds. Bool is the zero value because most flags are switches, and a flag
// declared without a kind should be the harmless one.
const (
	Bool Kind = iota
	String
	Int
	Duration
)

func (k Kind) String() string {
	switch k {
	case String:
		return "string"
	case Int:
		return "int"
	case Duration:
		return "duration"
	}
	return "bool"
}

// Completer offers the values an argument or flag can take.
//
// A function rather than a list, because the interesting completions are live:
// the environments in a config, the services in a stack. Shell completions and
// the TUI's own pickers come from the same one, so they cannot offer different
// answers.
type Completer func(prefix string) []string

// Arg is a positional argument.
type Arg struct {
	Name     string
	Help     string
	Required bool
	// Variadic takes the rest. Only the last argument may be.
	Variadic bool
	Complete Completer
}

// Flag is an option.
type Flag struct {
	Name string
	// Short is a single letter, without the dash. Empty for none.
	Short    string
	Kind     Kind
	Default  string
	Help     string
	Complete Completer
}

// Command is one thing a tool can do, declared once.
type Command struct {
	Name  string
	Short string
	// Long is the paragraph under the usage line. Optional.
	Long string

	Args  []Arg
	Flags []Flag

	// Run does the work and returns the exit code. A code rather than an
	// error, because the contract is richer than ok-or-not: see Codes.
	//
	// Nil makes this a grouping command, which prints its children's usage.
	Run func(Call) int

	// Commands are subcommands.
	Commands []Command

	// --- the other three surfaces ---

	// Screen names the TUI screen this command opens, if it opens one.
	Screen string

	// Target is the region type this command acts on. A command with a Target
	// appears in that region's context menu, and nowhere else has to be told.
	Target comp.Name

	// Key is the keyboard path to this command inside the TUI.
	//
	// Required for anything with a Target, and guard.Reachable holds that
	// closed. An agent cannot click, and a right-click may never arrive at all
	// — a multiplexer that captures it gets the event first and the
	// application never learns it happened.
	Key string

	// Hidden keeps a command out of help and completions without removing it.
	// For things that still work but should not be advertised.
	Hidden bool
}

// Call is one invocation, parsed.
type Call struct {
	// Args are the positional arguments by name. A variadic argument's extra
	// values are in Rest.
	Args map[string]string
	Rest []string
	// Flags are every flag by name, including defaults not given.
	Flags map[string]string
	// JSON says the caller asked for machine-readable output. Every command
	// gets it, so an agent never has to find out which ones support it.
	JSON bool
	// Out and Err are where to write. Fields rather than os.Stdout, so a test
	// reads what a command printed instead of capturing a global.
	Out, Err Writer
}

// Writer is the little of io.Writer this package needs, kept narrow so a test
// can pass a buffer without importing anything.
type Writer interface{ Write(p []byte) (int, error) }

// Flag reads a flag's value.
func (c Call) Flag(name string) string { return c.Flags[name] }

// Bool reads a boolean flag.
func (c Call) Bool(name string) bool { return c.Flags[name] == "true" }

// Arg reads a positional argument.
func (c Call) Arg(name string) string { return c.Args[name] }
