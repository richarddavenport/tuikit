package spec

import (
	"flag"
	"fmt"
	"io"
	"strings"
)

// Run parses argv against a command tree and runs what it names, returning the
// exit code.
//
// A code returned rather than os.Exit called, which is the whole reason this is
// not cobra: main decides when the process ends, and a test can run a command
// and read what it printed. Nothing in this package exits.
func Run(root Command, argv []string, out, errw io.Writer) int {
	cmd, rest, path := resolve(root, argv)

	// Help is answered here rather than by each command, so no command can
	// forget to support it and none can answer it differently.
	if wantsHelp(rest) {
		Usage(cmd, path, out)
		return OK
	}

	if cmd.Run == nil {
		// A grouping command with no subcommand named. Not an error if they
		// asked for the group itself, but they gave us nothing to do.
		Usage(cmd, path, errw)
		return Fail
	}

	call, err := Parse(cmd, rest)
	if err != nil {
		printf(errw, "%s: %v\n\n", strings.Join(path, " "), err)
		Usage(cmd, path, errw)
		return Fail
	}
	call.Out, call.Err = out, errw
	return cmd.Run(call)
}

// resolve walks the tree as far as the arguments name subcommands, and returns
// what is left.
func resolve(cmd Command, argv []string) (Command, []string, []string) {
	path := []string{cmd.Name}
	for len(argv) > 0 && !strings.HasPrefix(argv[0], "-") {
		next, ok := child(cmd, argv[0])
		if !ok {
			break
		}
		cmd, argv, path = next, argv[1:], append(path, next.Name)
	}
	return cmd, argv, path
}

func child(cmd Command, name string) (Command, bool) {
	for _, sub := range cmd.Commands {
		if sub.Name == name {
			return sub, true
		}
	}
	return Command{}, false
}

func wantsHelp(argv []string) bool {
	for _, a := range argv {
		if a == "-h" || a == "--help" || a == "help" {
			return true
		}
	}
	return false
}

// Parse turns arguments into a Call.
//
// stdlib flag does the flags; the positional arguments are checked here,
// because flag has no idea what they mean and a missing required argument
// should say which one rather than "expected 2 args, got 1".
func Parse(cmd Command, argv []string) (Call, error) {
	fs := flag.NewFlagSet(cmd.Name, flag.ContinueOnError)
	fs.SetOutput(io.Discard) // the error is reported by Run, with usage

	values := map[string]*string{}
	for _, f := range flagsOf(cmd) {
		v := new(string)
		*v = f.Default
		if f.Kind == Bool && f.Default == "" {
			*v = "false"
		}
		values[f.Name] = v
		fs.Var(&flagValue{v: v, boolean: f.Kind == Bool}, f.Name, f.Help)
		if f.Short != "" {
			fs.Var(&flagValue{v: v, boolean: f.Kind == Bool}, f.Short, f.Help)
		}
	}
	// Every command takes --json, so an agent never has to find out which ones
	// support it.
	jsonFlag := new(string)
	*jsonFlag = "false"
	fs.Var(&flagValue{v: jsonFlag, boolean: true}, "json", "machine-readable output")

	// Flags and positionals are separated before stdlib flag sees them.
	//
	// flag.Parse STOPS at the first non-flag argument, so `deploy api_gateway
	// --dry-run` would parse no flags at all and silently take --dry-run as a
	// positional. This repo has already been bitten: cmd/tuikit's frames
	// command works around it by hand, with the note that "the directory goes
	// after the flags" is a rule nobody remembers and nothing enforces. A
	// generated CLI must not make every tool learn that.
	flagArgs, positional, extra, err := split(cmd, argv)
	if err != nil {
		return Call{}, err
	}
	if err := fs.Parse(flagArgs); err != nil {
		return Call{}, err
	}

	call := Call{
		Args:  map[string]string{},
		Flags: map[string]string{},
		Extra: pairs(extra),
		JSON:  *jsonFlag == "true",
	}
	for name, v := range values {
		call.Flags[name] = *v
	}

	for i, arg := range cmd.Args {
		switch {
		case arg.Variadic:
			if i < len(positional) {
				call.Args[arg.Name] = positional[i]
				call.Rest = positional[i:]
			}
		case i < len(positional):
			call.Args[arg.Name] = positional[i]
		case arg.Required:
			return Call{}, fmt.Errorf("missing argument %s", arg.Name)
		}
	}
	if last := len(cmd.Args); len(positional) > last && (last == 0 || !cmd.Args[last-1].Variadic) {
		return Call{}, fmt.Errorf("unexpected argument %q", positional[last])
	}
	return call, nil
}

// split separates flags from positional arguments, wherever they appear.
//
// A flag that takes a value swallows the token after it, which is the only
// thing that needs knowing about the declaration — and the reason this cannot
// be a generic argument shuffle.
func split(cmd Command, argv []string) (flags, positional, extra []string, err error) {
	takesValue := map[string]bool{}
	known := map[string]bool{"json": true}
	for _, f := range flagsOf(cmd) {
		known[f.Name] = true
		if f.Kind != Bool {
			takesValue[f.Name] = true
		}
		if f.Short != "" {
			known[f.Short] = true
			takesValue[f.Short] = f.Kind != Bool
		}
	}

	for i := 0; i < len(argv); i++ {
		arg := argv[i]
		switch {
		case arg == "--":
			// Everything after is positional, however it looks. A service
			// really can be called -w.
			return flags, append(positional, argv[i+1:]...), extra, nil
		case !strings.HasPrefix(arg, "-") || arg == "-":
			positional = append(positional, arg)
			continue
		}

		name, _, hasValue := strings.Cut(strings.TrimLeft(arg, "-"), "=")
		if !known[name] {
			if !cmd.PassThrough {
				return nil, nil, nil, fmt.Errorf("unknown flag -%s", name)
			}
			// Collected rather than rejected, and it takes a value: a
			// pass-through command cannot know which of its flags are
			// switches, so it treats them all as values and lets whatever
			// consumes them decide.
			extra = append(extra, arg)
			if !hasValue && i+1 < len(argv) && !strings.HasPrefix(argv[i+1], "-") {
				i++
				extra = append(extra, argv[i])
			}
			continue
		}
		flags = append(flags, arg)
		if takesValue[name] && !hasValue {
			if i+1 >= len(argv) {
				return nil, nil, nil, fmt.Errorf("flag -%s needs a value", name)
			}
			i++
			flags = append(flags, argv[i])
		}
	}
	return flags, positional, extra, nil
}

// flagsOf is a command's declared flags plus the ones every command of its kind
// gets. A command that opens a screen can be captured, and saying so once here
// beats every tool remembering to declare two flags identically.
func flagsOf(cmd Command) []Flag {
	if cmd.Screen == "" {
		return cmd.Flags
	}
	return append(append([]Flag(nil), cmd.Flags...), SnapshotFlags...)
}

// flagValue adapts a string to flag.Value, so every flag is read the same way
// whatever its kind. The kinds are for help, completion and describe; a command
// that wants an int parses one.
type flagValue struct {
	v       *string
	boolean bool
}

func (f *flagValue) String() string {
	if f.v == nil {
		return ""
	}
	return *f.v
}

func (f *flagValue) Set(s string) error { *f.v = s; return nil }

// IsBoolFlag makes -watch work without a value, the way a switch should.
func (f *flagValue) IsBoolFlag() bool { return f.boolean }

// Usage writes help for a command, generated from its declaration.
//
// Generated, so it cannot describe a flag the command does not have — which is
// what a hand-written usage string does about six months in.
func Usage(cmd Command, path []string, w io.Writer) {
	line := strings.Join(path, " ")
	for _, a := range cmd.Args {
		switch {
		case a.Variadic:
			line += " [" + a.Name + "...]"
		case a.Required:
			line += " <" + a.Name + ">"
		default:
			line += " [" + a.Name + "]"
		}
	}
	if len(flagsOf(cmd)) > 0 {
		line += " [flags]"
	}
	printf(w, "usage: %s\n", line)
	if cmd.Short != "" {
		printf(w, "\n%s\n", cmd.Short)
	}
	if cmd.Long != "" {
		printf(w, "\n%s\n", cmd.Long)
	}

	if visible := visible(cmd.Commands); len(visible) > 0 {
		printf(w, "\ncommands:\n")
		width := 0
		for _, sub := range visible {
			width = max(width, len(sub.Name))
		}
		for _, sub := range visible {
			printf(w, "  %-*s  %s\n", width, sub.Name, sub.Short)
		}
	}

	if len(cmd.Args) > 0 {
		printf(w, "\narguments:\n")
		for _, a := range cmd.Args {
			printf(w, "  %-12s  %s\n", a.Name, a.Help)
		}
	}

	printf(w, "\nflags:\n")
	for _, f := range flagsOf(cmd) {
		name := "--" + f.Name
		if f.Short != "" {
			name = "-" + f.Short + ", " + name
		}
		if f.Kind != Bool {
			name += " " + f.Kind.String()
		}
		help := f.Help
		if f.Default != "" {
			help += " (default " + f.Default + ")"
		}
		printf(w, "  %-24s  %s\n", name, help)
	}
	printf(w, "  %-24s  %s\n", "--json", "machine-readable output")
}

// printf writes and drops the error.
//
// Nothing useful can be done with a failed write to stderr on the way out of a
// command, and threading one through every line of a usage message would say
// nothing a reader needs. Named here so that is a decision taken once rather
// than twelve unchecked calls.
func printf(w io.Writer, format string, args ...any) {
	_, _ = fmt.Fprintf(w, format, args...)
}

func visible(cmds []Command) []Command {
	var out []Command
	for _, c := range cmds {
		if !c.Hidden {
			out = append(out, c)
		}
	}
	return out
}

// pairs turns collected pass-through arguments into names and values.
func pairs(args []string) map[string]string {
	out := map[string]string{}
	for i := 0; i < len(args); i++ {
		name, value, hasValue := strings.Cut(strings.TrimLeft(args[i], "-"), "=")
		if !hasValue && i+1 < len(args) && !strings.HasPrefix(args[i+1], "-") {
			i++
			value = args[i]
		}
		out[name] = value
	}
	return out
}
