package spec

import (
	"fmt"
	"io"
	"strings"
)

// Complete offers what could come next.
//
// One function, used by every shell and by the TUI's own pickers, so bash and
// the interface cannot offer different answers. The shell scripts below do not
// know anything about the commands — they call back into the tool, which is
// what keeps them from going stale when a flag is added.
func Complete(root Command, argv []string) []string {
	// The word being completed is the last one, and it is a prefix rather than
	// a whole word. An empty final argument means "what could I type here",
	// which is what a shell sends when the cursor is after a space.
	prefix := ""
	if len(argv) > 0 {
		prefix = argv[len(argv)-1]
		argv = argv[:len(argv)-1]
	}

	cmd, rest, _ := resolve(root, argv)

	// After a flag that takes a value, the value is what comes next — not
	// another flag, and not a subcommand.
	if flag, ok := awaitingValue(cmd, rest); ok {
		return offer(flag.Complete, prefix)
	}

	if strings.HasPrefix(prefix, "-") {
		return flagNames(cmd, prefix)
	}

	var out []string
	for _, sub := range cmd.Commands {
		if !sub.Hidden && strings.HasPrefix(sub.Name, prefix) {
			out = append(out, sub.Name)
		}
	}

	// A positional argument's own completions, for the position we are at.
	if arg, ok := nextArg(cmd, rest); ok {
		out = append(out, offer(arg.Complete, prefix)...)
	}
	return out
}

// awaitingValue reports that the last argument was a flag still wanting one.
func awaitingValue(cmd Command, argv []string) (Flag, bool) {
	if len(argv) == 0 {
		return Flag{}, false
	}
	last := argv[len(argv)-1]
	if !strings.HasPrefix(last, "-") || strings.Contains(last, "=") {
		return Flag{}, false
	}
	name := strings.TrimLeft(last, "-")
	for _, f := range flagsOf(cmd) {
		if (f.Name == name || f.Short == name) && f.Kind != Bool {
			return f, true
		}
	}
	return Flag{}, false
}

// nextArg is the positional argument a new word would fill.
func nextArg(cmd Command, argv []string) (Arg, bool) {
	var given int
	for i := 0; i < len(argv); i++ {
		if !strings.HasPrefix(argv[i], "-") {
			given++
			continue
		}
		if f, ok := awaitingValue(cmd, argv[:i+1]); ok && f.Name != "" {
			i++ // the value belongs to the flag, not to a position
		}
	}
	if given < len(cmd.Args) {
		return cmd.Args[given], true
	}
	if n := len(cmd.Args); n > 0 && cmd.Args[n-1].Variadic {
		return cmd.Args[n-1], true
	}
	return Arg{}, false
}

func offer(c Completer, prefix string) []string {
	if c == nil {
		return nil
	}
	var out []string
	for _, v := range c(prefix) {
		if strings.HasPrefix(v, prefix) {
			out = append(out, v)
		}
	}
	return out
}

func flagNames(cmd Command, prefix string) []string {
	var out []string
	for _, f := range flagsOf(cmd) {
		if name := "--" + f.Name; strings.HasPrefix(name, prefix) {
			out = append(out, name)
		}
	}
	// The constant is the CANDIDATE and prefix is what was typed, so the
	// arguments read backwards and are the right way round.
	if strings.HasPrefix("--json", prefix) { //nolint:gocritic // see above
		out = append(out, "--json")
	}
	return out
}

// CompletionScript writes the shell's side of it.
//
// Deliberately tiny, and deliberately ignorant: it calls the tool back rather
// than listing anything itself. A generated script that enumerated the flags
// would be a fourth copy of the declaration, and the one nobody regenerates
// after adding one.
func CompletionScript(shell, tool string, w io.Writer) error {
	switch shell {
	case "bash":
		_, err := fmt.Fprintf(w, `# %[1]s completion for bash. Load with: source <(%[1]s completion bash)
_%[1]s() {
  local IFS=$'\n'
  COMPREPLY=($(%[1]s __complete "${COMP_WORDS[@]:1}"))
}
complete -o default -F _%[1]s %[1]s
`, tool)
		return err
	case "zsh":
		_, err := fmt.Fprintf(w, `# %[1]s completion for zsh. Load with: source <(%[1]s completion zsh)
_%[1]s() {
  local -a reply
  reply=(${(f)"$(%[1]s __complete ${words[2,-1]})"})
  _describe '%[1]s' reply
}
compdef _%[1]s %[1]s
`, tool)
		return err
	case "fish":
		_, err := fmt.Fprintf(w, `# %[1]s completion for fish. Load with: %[1]s completion fish | source
complete -c %[1]s -f -a '(%[1]s __complete (commandline -opc)[2..-1] (commandline -ct))'
`, tool)
		return err
	}
	return fmt.Errorf("no completion for %q — bash, zsh or fish", shell)
}
