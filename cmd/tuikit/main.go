// Command tuikit is the framework's own tooling.
//
// Only `designsystem` exists so far. `new`, `watch` and `gallery` are the
// planned surface and are named here rather than left undiscoverable — see
// design/decisions.md.
package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/richarddavenport/tuikit/docgen"
)

var version = "dev"

func main() {
	if len(os.Args) < 2 {
		usage(os.Stderr)
		os.Exit(2)
	}
	switch os.Args[1] {
	case "designsystem":
		designsystem(os.Args[2:])
	case "version":
		fmt.Println("tuikit", version)
	case "help", "-h", "--help":
		usage(os.Stdout)
	default:
		fmt.Fprintf(os.Stderr, "tuikit: unknown command %q\n\n", os.Args[1])
		usage(os.Stderr)
		os.Exit(2)
	}
}

func designsystem(args []string) {
	fs := flag.NewFlagSet("designsystem", flag.ExitOnError)
	out := fs.String("out", "design-system", "directory to write the bundle into")
	tool := fs.String("tool", "", "the tool's name, used in the prose")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	// The default palette and glyph set, because this binary renders tuikit's
	// own vocabulary. A tool with its own palette calls docgen from its own
	// generator — the bundle has to come from the code that ships.
	written, err := docgen.DesignSystem{Tool: *tool}.Write(*out)
	for _, path := range written {
		fmt.Println("  " + path)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "tuikit designsystem:", err)
		os.Exit(1)
	}
	fmt.Printf("%d cards written to %s\n", len(written), *out)
}

func usage(w *os.File) {
	fmt.Fprint(w, `tuikit — a TUI framework for developers and agents

  tuikit designsystem [-out dir] [-tool name]
        write the foundations bundle as HTML

  tuikit version

Planned: new (scaffold a tool), watch (recapture on save), gallery (browse the
components). See design/.
`)
}
