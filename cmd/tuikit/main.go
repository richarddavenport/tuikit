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
	"strings"

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
	case "frames":
		frames(os.Args[2:])
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

// frames turns a capture directory into a page.
func frames(args []string) {
	fs := flag.NewFlagSet("frames", flag.ExitOnError)
	out := fs.String("out", "frames.html", "file to write")
	title := fs.String("title", "", "the page's title; defaults to the directory's name")
	lede := fs.String("lede", "", "a sentence under the title")

	// The directory is taken before parsing, so it may come first or last.
	// stdlib flag stops at the first non-flag argument, which would make
	// `tuikit frames ./frames -out page.html` silently parse no flags at all —
	// and "the directory goes after the flags" is a rule nobody remembers and
	// nothing enforces.
	var dir string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		dir, args = args[0], args[1:]
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if dir == "" && fs.NArg() == 1 {
		dir = fs.Arg(0)
	}
	if dir == "" || fs.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "usage: tuikit frames <capture-dir> [-out page.html]")
		os.Exit(2)
	}

	// No groups from the command line: grouping carries meaning a flag cannot
	// express — what the reader is doing — so a tool that wants it calls
	// docgen.Frames itself. Ungrouped, every frame still appears.
	if err := (docgen.Frames{Title: *title, Lede: *lede}).Write(dir, *out); err != nil {
		fmt.Fprintln(os.Stderr, "tuikit frames:", err)
		os.Exit(1)
	}
	fmt.Println("wrote " + *out)
}

func usage(w *os.File) {
	_, _ = fmt.Fprint(w, `tuikit — a TUI framework for developers and agents

  tuikit designsystem [-out dir] [-tool name]
        write the foundations bundle as HTML

  tuikit frames <capture-dir> [-out page.html] [-title t] [-lede l]
        turn a captured run of frames into a page you can look at

  tuikit version

Planned: new (scaffold a tool), watch (recapture on save, serve, reload),
gallery (browse the components). See design/.
`)
}
