// Command tuikit is the framework's own tooling.
//
//	tuikit new           scaffold a tool that builds and passes its own checks
//	tuikit designsystem  write the foundations bundle as HTML
//	tuikit frames        turn a captured run of frames into a page
//	tuikit watch         recapture on save and reload the browser
//	tuikit gallery       open every component, running
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

	"github.com/richarddavenport/tuikit/app"
	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/docgen"
	"github.com/richarddavenport/tuikit/gallery"
	"github.com/richarddavenport/tuikit/scaffold"
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
	case "new":
		newTool(os.Args[2:])
	case "designsystem":
		designsystem(os.Args[2:])
	case "frames":
		frames(os.Args[2:])
	case "watch":
		watchCmd(os.Args[2:])
	case "pixels":
		pixels(os.Args[2:])
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

// newTool scaffolds a tool.
//
// It generates and then BOOTSTRAPS: resolves modules and writes the goldens for
// the screens the tool was born with. A generated tool that fails `go test
// ./...` on its first run has taught its owner, in the first thirty seconds,
// that the tests are noise.
func newTool(args []string) {
	fs := flag.NewFlagSet("new", flag.ExitOnError)
	module := fs.String("module", "", "the go.mod path; defaults to the tool's name")
	short := fs.String("short", "", "one sentence saying what the tool does")
	dir := fs.String("dir", ".", "where to create the tool's directory")
	tuikitPath := fs.String("tuikit", "", "what the `replace` points at; defaults to ../tuikit")
	skip := fs.Bool("no-bootstrap", false, "write the files and stop, without running go")

	// The name is taken before parsing, so it may come first or last. stdlib
	// flag stops at the first non-flag argument, which would make `tuikit new
	// widgetctl -short "..."` silently parse no flags at all — the same trap
	// `frames` works around below, and the reason spec does the separation for
	// every generated CLI.
	var name string
	if len(args) > 0 && !strings.HasPrefix(args[0], "-") {
		name, args = args[0], args[1:]
	}
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}
	if name == "" && fs.NArg() == 1 {
		name = fs.Arg(0)
	}
	if name == "" || fs.NArg() > 1 {
		fmt.Fprintln(os.Stderr, "usage: tuikit new <name> [-module path] [-short text] [-dir where]")
		os.Exit(2)
	}

	tool := scaffold.Tool{
		Name:   name,
		Module: *module,
		Short:  *short,
		Tuikit: *tuikitPath,
	}
	written, err := tool.Write(*dir)
	for _, path := range written {
		fmt.Println("  " + path)
	}
	if err != nil {
		fmt.Fprintln(os.Stderr, "tuikit new:", err)
		os.Exit(1)
	}

	root := filepath.Join(*dir, tool.Name)
	if !*skip {
		if err := scaffold.Bootstrap(root); err != nil {
			fmt.Fprintln(os.Stderr, "tuikit new:", err)
			os.Exit(1)
		}
	}
	fmt.Printf("\n%d files written to %s\n\n  cd %s && make check\n", len(written), root, root)
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
	// Detection is asked here rather than inside the gallery, because it reads
	// /dev/tty: a model that queried in its constructor would ask the
	// developer's real terminal during `go test`. Nothing is asked of a
	// terminal that turns out to draw only characters.
	//
	// The runner owns the canvas — see app.Model. The gallery has no View.
	g := app.New(gallery.New(theme.Default), app.WithPixels(comp.Detect()))

	p := tea.NewProgram(g, tea.WithAltScreen(), tea.WithMouseCellMotion())
	if _, err := p.Run(); err != nil {
		fmt.Fprintln(os.Stderr, "tuikit gallery:", err)
		os.Exit(1)
	}
}

func usage(w *os.File) {
	_, _ = fmt.Fprint(w, `tuikit — a TUI framework for developers and agents

  tuikit new <name> [-module path] [-short text] [-dir where]
        scaffold a tool that builds, runs and passes its own checks

  tuikit designsystem [-out dir] [-tool name]
        write the foundations bundle as HTML

  tuikit frames <capture-dir> [-out page.html] [-title t] [-lede l]
        turn a captured run of frames into a page you can look at

  tuikit watch <dir> -capture "<command>" -frames <dir>
        recapture on save, rebuild the page, reload the browser

  tuikit gallery
        open every component, running, with its states and keys

  tuikit pixels
        what this terminal can draw, and the same bar with and without it

  tuikit version

`)
}
