// Command democtl is tuikit's example tool and its own test fixture.
//
// It manages a fictional service fleet generated from a seed, so it needs no
// backend and shows the same thing every run. That is what lets tuikit's capture
// harness render it as a fixture, and what lets a reader compare a frame in the
// documentation against what they see on their own terminal.
//
//	go run ./examples/democtl
//
// It is also the reference implementation: the engine/ui split, the palette and
// glyph set held closed by guards, the mode split that stops a list scrolling
// behind a modal, and the generation counter that stops an abandoned run
// drawing into its successor.
package main

import (
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/richarddavenport/tuikit/examples/democtl/ui"
	"github.com/richarddavenport/tuikit/spec"
)

var version = "dev"

// main decides when the process ends.
//
// Nothing in spec calls os.Exit, which is why it is not cobra: the exit-code
// contract is richer than ok-or-not — 2 means a dry run found drift — and a
// framework that owns the exit is a framework that flattens it.
//
// The dual-mode entry lives here too: no arguments opens the TUI, and anything
// else is a command. main is the only place that can decide which, because it
// is the only place that knows a bare `democtl` means "show me".
func main() {
	// -seed is taken by hand rather than with flag.Parse, and the reason is
	// what a reader of this file is most likely to copy.
	//
	// flag.CommandLine owns -h and --help. A flag.Parse here answers them
	// itself, prints the stdlib's usage and exits, so democtl's own help never
	// runs and `democtl --help` says nothing about democtl. It also prints
	// every flag any linked package registered at init: harness registers
	// -update-goldens, ui imports harness for --snapshot, and a shipped tool's
	// --help ends up advertising a testing flag.
	//
	// The seed has to be read before the tree exists, because the tree is built
	// from it. That is the whole reason it is not a declared flag.
	seed, argv := takeSeed(os.Args[1:])
	if len(argv) == 0 {
		argv = []string{"tui"}
	}

	root := ui.Commands(seed)
	switch argv[0] {
	case "describe":
		// The whole surface in one call, so an agent never has to grep for it.
		out, err := spec.Describe(root, version, ui.Palette, ui.Glyphs).JSON()
		if err != nil {
			fmt.Fprintln(os.Stderr, "democtl:", err)
			os.Exit(spec.Fail)
		}
		fmt.Println(string(out))
		os.Exit(spec.OK)

	case "completion":
		if len(argv) < 2 {
			fmt.Fprintln(os.Stderr, "usage: democtl completion <bash|zsh|fish>")
			os.Exit(spec.Fail)
		}
		if err := spec.CompletionScript(argv[1], "democtl", os.Stdout); err != nil {
			fmt.Fprintln(os.Stderr, "democtl:", err)
			os.Exit(spec.Fail)
		}
		os.Exit(spec.OK)

	case "__complete":
		// What the shell calls back into. The scripts know nothing about the
		// commands, so they cannot go stale when one is added.
		for _, word := range spec.Complete(root, argv[1:]) {
			fmt.Println(word)
		}
		os.Exit(spec.OK)

	}
	// `tui` is a declared command now, so it goes through spec like the rest —
	// which is what gives it --snapshot and --script, in help and in the
	// manifest, without democtl declaring either.
	os.Exit(spec.Run(root, argv, os.Stdout, os.Stderr))
}

// takeSeed pulls -seed out of the arguments and returns the rest.
//
// Both spellings and both separators, because a reader who types --seed=7 has
// not made a mistake worth an error message. A value that is not a number is
// ignored rather than rejected: this is a fixture knob, and the run it would
// abort is a run somebody wanted to look at.
func takeSeed(argv []string) (int64, []string) {
	seed := int64(1)
	rest := make([]string, 0, len(argv))
	for i := 0; i < len(argv); i++ {
		name, value, split := strings.Cut(argv[i], "=")
		if name != "-seed" && name != "--seed" {
			rest = append(rest, argv[i])
			continue
		}
		if !split && i+1 < len(argv) {
			value, i = argv[i+1], i+1
		}
		if n, err := strconv.ParseInt(value, 10, 64); err == nil {
			seed = n
		}
	}
	return seed, rest
}
