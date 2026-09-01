// Command tuikit is the framework's own tooling.
//
// Only `designsystem` exists so far. `new`, `watch` and `gallery` are the
// planned surface and are named here rather than left undiscoverable — see
// design/decisions.md.
package main

import (
	"context"
	"flag"
	"fmt"
	"os"
	"os/signal"
	"path/filepath"
	"strings"
	"syscall"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/docgen"
	"github.com/richarddavenport/tuikit/gallery"
	"github.com/richarddavenport/tuikit/theme"
	"github.com/richarddavenport/tuikit/watch"
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
	case "watch":
		watchCmd(os.Args[2:])
	case "gallery":
		galleryCmd(os.Args[2:])
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

// watchCmd runs the inner loop.
func watchCmd(args []string) {
	fs := flag.NewFlagSet("watch", flag.ExitOnError)
	capture := fs.String("capture", "", "the command that captures frames (required)")
	framesDir := fs.String("frames", "", "where the capture writes (required)")
	out := fs.String("out", "", "the page file; defaults to <frames>/index.html")
	addr := fs.String("addr", "127.0.0.1:7654", "address to serve on")
	title := fs.String("title", "", "the page's title")
	every := fs.Duration("every", 0, "how often to check for changes; default 300ms")

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
	if dir == "" || *capture == "" || *framesDir == "" {
		fmt.Fprintln(os.Stderr, `usage: tuikit watch <dir> -capture "<command>" -frames <dir> [-out page.html]

  tuikit watch ./internal/tui \
    -capture "go test ./internal/tui -run CaptureFrames" \
    -frames /tmp/frames`)
		os.Exit(2)
	}
	if *out == "" {
		*out = filepath.Join(*framesDir, "index.html")
	}

	cfg := watch.Config{
		Dir: dir,
		// Split on spaces rather than shelling out: a capture command is a
		// program and its arguments, and going through a shell would make the
		// tool's own quoting somebody else's problem.
		Capture:  strings.Fields(*capture),
		Frames:   *framesDir,
		Out:      *out,
		Interval: *every,
		Log:      os.Stdout,
		Page: func(framesDir string) (string, error) {
			return docgen.Frames{Title: *title}.Page(framesDir)
		},
	}

	// stop is called before exiting rather than deferred: os.Exit skips defers,
	// and a watch that leaves the signal handler installed on its way out is a
	// watch that ignores the second Ctrl-C.
	ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)

	err := cfg.Run(ctx, *addr, func(url string) {
		fmt.Printf("watching %s\n", dir)
		fmt.Printf("  capture  %s\n", *capture)
		fmt.Printf("  serving  %s\n", url)
	})
	stop()
	if err != nil {
		fmt.Fprintln(os.Stderr, "tuikit watch:", err)
		os.Exit(1)
	}
}

// galleryCmd opens every component, running.
//
// A real TUI rather than a page, because what a component is like to USE — what
// it feels like to arrow through, what it does at 80 columns, what it looks
// like empty — is not a thing a screenshot answers. tuikit designsystem renders
// the vocabulary so you can review it; this lets you use it.
func galleryCmd(args []string) {
	fs := flag.NewFlagSet("gallery", flag.ExitOnError)
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	// The default palette, because this binary shows tuikit's own components. A
	// tool checking its OWN palette against them calls gallery.New with it —
	// which is the point of taking one at all.
	p := tea.NewProgram(gallery.New(theme.Default),
		tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tuikit gallery:", err)
		os.Exit(1)
	}
}

func usage(w *os.File) {
	_, _ = fmt.Fprint(w, `tuikit — a TUI framework for developers and agents

  tuikit designsystem [-out dir] [-tool name]
        write the foundations bundle as HTML

  tuikit frames <capture-dir> [-out page.html] [-title t] [-lede l]
        turn a captured run of frames into a page you can look at

  tuikit watch <dir> -capture "<command>" -frames <dir>
        recapture on save, rebuild the page, reload the browser

  tuikit gallery
        open every component, running, with its states and keys

  tuikit version

Planned: new (scaffold a tool). See design/.
`)
}
