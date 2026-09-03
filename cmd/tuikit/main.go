// Command tuikit is the framework's own tooling.
//
//	tuikit new           scaffold a tool that builds and passes its own checks
//	tuikit designsystem  write the foundations bundle as HTML
//	tuikit frames        turn a captured run of frames into a page
//	tuikit watch         recapture on save and reload the browser
//	tuikit news          what tuikit decided since this tool last looked
//	tuikit gallery       open every component, running (-list to print it)
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
	"github.com/richarddavenport/tuikit/news"
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
	case "news":
		newsCmd(os.Args[2:])
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
	out := fs.String("out", "", "file to write; defaults to frames.html, or frames.md with -md")
	title := fs.String("title", "", "the page's title; defaults to the directory's name")
	lede := fs.String("lede", "", "a sentence under the title")
	md := fs.Bool("md", false, "write Markdown with an SVG per frame, for a repository")

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
		fmt.Fprintln(os.Stderr, "usage: tuikit frames <capture-dir> [-out page.html] [-md]")
		os.Exit(2)
	}

	// The default follows the format rather than being one name for both, so
	// `-md` on its own does the obvious thing instead of writing Markdown into
	// a file called frames.html.
	if *out == "" {
		*out = map[bool]string{false: "frames.html", true: "frames.md"}[*md]
	}

	// No groups from the command line: grouping carries meaning a flag cannot
	// express — what the reader is doing — so a tool that wants it calls
	// docgen.Frames itself. Ungrouped, every frame still appears.
	page := docgen.Frames{Title: *title, Lede: *lede}
	write := page.Write
	if *md {
		write = page.Markdown
	}
	if err := write(dir, *out); err != nil {
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
// newsCmd prints what tuikit has decided since a tool last looked.
//
// Run from inside the tool, with no arguments: it reads the tool's own go.mod
// to find the tuikit it builds against, and the tool's AGENTS.md for the
// decision it recorded. Both defaults exist so the command is one word — a
// check nobody has to remember the arguments for is a check that gets run.
func newsCmd(args []string) {
	fs := flag.NewFlagSet("news", flag.ExitOnError)
	since := fs.Int("since", -1, "decision number to report from; default is the one recorded in AGENTS.md")
	dir := fs.String("tuikit", "", "the tuikit checkout; default is the replace directive in ./go.mod")
	notes := fs.String("marker", "AGENTS.md", "the file holding \"reconciled with tuikit through decision N\"")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	root := *dir
	if root == "" {
		found, ok := news.Checkout("go.mod")
		if !ok {
			fmt.Fprintln(os.Stderr, "tuikit news: no tuikit replace directive in ./go.mod — run this from a tool, or pass -tuikit <path>")
			os.Exit(2)
		}
		root = found
	}

	all, err := news.Read(filepath.Join(root, "design", "decisions.md"))
	if err != nil {
		fmt.Fprintln(os.Stderr, "tuikit news:", err)
		os.Exit(1)
	}

	from := *since
	if from < 0 {
		n, ok := news.Marker(*notes)
		if !ok {
			fmt.Fprintf(os.Stderr, "tuikit news: %s does not say which decision this tool reconciled with.\n\n"+
				"Add a line to it:\n\n    Reconciled with tuikit through decision %d.\n\n"+
				"Then this command reports only what lands after it. To read everything, -since 0.\n",
				*notes, news.Latest(all))
			os.Exit(2)
		}
		from = n
	}

	fresh := news.Since(all, from)
	if len(fresh) == 0 {
		fmt.Printf("Up to date with tuikit through decision %d.\n", news.Latest(all))
		return
	}

	fmt.Printf("%d decisions since %d:\n", len(fresh), from)
	for _, d := range fresh {
		fmt.Printf("\n  %d. %s\n", d.Number, d.Title)
		for _, line := range wrap(d.Lede, 74) {
			fmt.Println("     " + line)
		}
	}
	// Said every time, because the cost of keying this to decisions is that a
	// change nobody wrote a decision for does not appear above.
	fmt.Printf("\n  Not every addition gets a decision. `tuikit gallery -list` is the\n"+
		"  complete inventory, and a test holds it closed against the package.\n\n"+
		"  When you have read these, record it:\n\n      Reconciled with tuikit through decision %d.\n",
		news.Latest(all))
}

// wrap breaks a lede to width, so a terminal shows a paragraph rather than one
// line it has to scroll.
func wrap(s string, width int) []string {
	var out []string
	line := ""
	for _, word := range strings.Fields(s) {
		if line == "" {
			line = word
			continue
		}
		if len(line)+1+len(word) > width {
			out = append(out, line)
			line = word
			continue
		}
		line += " " + word
	}
	if line != "" {
		out = append(out, line)
	}
	return out
}

func galleryCmd(args []string) {
	fs := flag.NewFlagSet("gallery", flag.ExitOnError)
	list := fs.Bool("list", false, "print the inventory as text instead of running it")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	// The gallery is the one COMPLETE list of what tuikit offers —
	// TestEveryComponentIsInTheGallery holds it closed against the package, so
	// a component cannot exist without an entry. That makes it the right answer
	// to "what does tuikit already have", and running a TUI is the wrong way to
	// ask: an agent building a tool cannot see a terminal, and a hand-written
	// list in a tool's AGENTS.md goes stale the week after it is written.
	if *list {
		for _, e := range gallery.New(theme.Default).Entries() {
			fmt.Printf("%-12s %s\n", e.Name, e.Summary)
			if e.From != "" {
				fmt.Printf("%-12s from %s\n", "", e.From)
			}
		}
		return
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

  tuikit frames <capture-dir> [-out page.html] [-title t] [-lede l] [-md]
        turn a captured run of frames into a page you can look at,
        or with -md into Markdown with an SVG per frame, for a repository

  tuikit watch <dir> -capture "<command>" -frames <dir>
        recapture on save, rebuild the page, reload the browser

  tuikit news [-since N] [-tuikit path]
        what tuikit has decided since this tool last looked.
        Run it from inside a tool, with no arguments.

  tuikit gallery [-list]
        open every component, running, with its states and keys.
        -list prints the inventory as text — what tuikit already has,
        which is the list to read before hand-rolling anything

  tuikit pixels
        what this terminal can draw, and the same bar with and without it

  tuikit version

`)
}
