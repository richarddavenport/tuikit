# tuikit

A TUI framework for developers and agents. Go, Bubble Tea, Lip Gloss.

It exists because four tools — [swarmctl], [pgctl], [azctl], [dugo] — arrived at
the same shape independently, and that shape currently travels by copy. tuikit is
that shape as a module: the vocabulary an interface is allowed to use, the guards
that hold it closed, the components, the capture harness, and a scaffolder that
writes a tool which builds and passes its own checks on the first run.

[swarmctl]: https://github.com/richarddavenport/swarmctl
[pgctl]: https://github.com/richarddavenport/pgctl
[azctl]: https://github.com/richarddavenport/azctl
[dugo]: https://github.com/richarddavenport/dugo

## What is here now

**`theme`** — the vocabulary. Nine colour roles named by role and never by hue,
as the terminal's own sixteen ANSI indices, so a tool is themed by whatever
themed the terminal; a closed glyph allow-list; and `Hex` for anything drawing
the palette outside a terminal.

```go
pal := theme.Default
pal.Accent = "33" // this tool is blue

glyphs := theme.DefaultGlyphs.With('×', "multiplication sign — replica counts")
```

**`guard`** — tests that hold the rules closed. Call them from your own packages:

```go
// in your UI package
func TestTheInterfaceStaysInItsVocabulary(t *testing.T) {
    guard.Tokens(t, ".", palette)  // no colour literals; no unused roles
    guard.Glyphs(t, ".", glyphs)   // nothing printed that a font may not have
}

// in your engine package
func TestTheEngineHasNeverHeardOfATerminal(t *testing.T) {
    guard.Engine(t, ".")           // no colour, no width, no keys, no framework
}
```

They parse the package rather than pattern-matching it, so a comment explaining a
glyph is not read as printing one. They fail rather than pass when pointed at a
directory with nothing to scan — scanning nothing silently is how everyone comes
to believe a rule is on when it is not.

`guard.Engine` holds the split that every tool in the family keeps: the engine
knows the domain, the UI knows the terminal. Its deny-list is by *category* —
drawing libraries, terminal capabilities, width measurement — because an engine
that measures display width has learned about columns whichever package it used.

**`harness`** — render your screens to files, so you can look at them:

```go
func TestCaptureFrames(t *testing.T) {
    dir := harness.Enabled("MYTOOL_FRAMES")
    if dir == "" {
        t.Skip("set MYTOOL_FRAMES to capture frames")
    }
    s := harness.Capture(t, dir, harness.Size(132, 38), harness.At(epoch))
    s.Shot("dashboard", newModel())
    s.Done()
}

func TestGoldens(t *testing.T) {
    harness.Golden(t, "testdata", "dashboard", newModel().View())
}
```

Capture writes one `.ansi` per frame, in colour, plus a manifest carrying each
frame's provenance — fixture, live, or composed. Goldens hold the *shape*,
colour stripped, so a layout change is an ordinary test failure and the report
names the line and its column count. `harness.HTML` turns a frame into a block
for a page.

**Two capture mechanisms, deliberately.** The test above reaches any state at
all — an error, a half-loaded pane, a modal over a forty-row table — because it
can touch unexported fields. The BINARY reaches what a keystroke reaches, needs
no test, and an agent finds it from `--help`:

```
mytool browse --snapshot /tmp/frames \
  --script 'press j j; shot moved; click items.row[2]; shot picked'
```

A command declared with a `Screen` gets those two flags from `spec`, so they are
in the help and in `describe --json` without the tool declaring either — and
`shot` is what makes this a capture rather than a keystroke replay: the script
says where the interesting frames are.

Capture is a **building tool first**: its purpose is seeing the screen you are
writing. Goldens and documentation are downstream of that — which is why the
loop matters more than any assertion the harness could offer:

```
tuikit watch ./internal/tui \
  -capture "go test ./internal/tui -run CaptureFrames" \
  -frames /tmp/frames
```

Change a line of view code, look at the frame. It recaptures on save, rebuilds
the page, and reloads the browser. When the build fails the error goes **on the
page** — leaving the last good frames up would describe a tool that no longer
exists.

**`tuikit gallery`** — every component, running, with its states and its keys.
The fastest way to find out whether this framework is for you:

```
make gallery
```

A real TUI rather than a page, because what a component is like to USE — what it
feels like to arrow through, what it does at 80 columns, what it looks like
empty — is not a thing a screenshot answers. It is built out of `comp`, so it
cannot show you a component that does not work; it found two bugs on its first
run, one of them in itself.

**`docgen`** — generated from the code, so it cannot drift. The vocabulary as
HTML, and a capture as a page you can look at:

```
tuikit designsystem -out design-system -tool mytool
tuikit frames ./frames -out page.html -title mytool
```

The page carries three rules that are properties of terminal frames rather than
editorial taste, so they live in the code: the frames keep a dark ground in both
light and dark themes (the ANSI was captured for a dark terminal), they use the
system monospace stack with no webfont (box-drawing and braille must share one
set of advance widths), and every frame shows its provenance — fixture, live or
composed.

## Starting a tool

```sh
make install    # puts tuikit on your PATH, once

tuikit new mytool -short "what it does" -module github.com/you/mytool
cd mytool && make check
```

That is the whole of it. The generated tool **builds, runs and passes its own
checks with no edits** — a test in `scaffold/` generates one into a temporary
directory and puts `go build`, `go vet`, `go test` and `describe --json`
through it, so the claim is run rather than asserted.

`tuikit new` came last on purpose. It generates what a real migration turned
out to need rather than what seemed likely beforehand, and the order it writes
things in is the order of how much trouble each one saves:

| You get | Because |
|---|---|
| `internal/engine`, `internal/tui`, `internal/cli` as separate packages | The engine knows the domain and has **no terminal imports**; the CLI is a peer of the interface over the same engine, not a wrapper around it. `guard.Engine` holds it closed from the first commit. |
| A **pointer** model with the components already in it | tuikit's components own state — a cursor, a viewport, a divider — and that state cannot survive being copied on every message. A value model compiles, runs, and silently forgets every scroll. It is the one change azctl's migration had to make in every file. |
| A fixture and a capture test | Screens you can look at without whatever the real backend needs. azctl went years with no screenshot of its own interface because that split was not there. |
| `guard_test.go`, already calling the guards | Five lines. Your palette, your glyph set and your engine's ignorance are closed from the first commit rather than from whenever someone thinks of it. |
| A `spec.Command` tree, and a `main` that returns its exit code | One declaration becomes the CLI, the screens, `describe --json` and the context menus. Nothing calls `os.Exit` but `main`, so the richer exit-code contract survives. |
| A Makefile, three workflows, a linter config | Repo infrastructure cannot be a dependency, which is the whole reason a scaffolder exists. `make check` is exactly what CI runs. |
| `README.md`, `AGENTS.md`, `CONTEXT.md`, `design/decisions.md` | An agent arriving at the repo is told where things live, what will fail its change, and what the words mean. |

`examples/democtl` is still there, and it is the bigger reference: a complete
working tool with a dashboard, a log pane, a context menu, a confirm modal and
a step run, in about 1,700 lines. Read it for a screen the scaffolder does not
generate; do not copy it to start.

## What is coming

The pieces are all here — `theme`, `guard`, `comp`, `app`, `spec`, `harness`,
`docgen`, and `tuikit` itself with `new`, `watch`, `gallery`, `frames` and
`designsystem`. What is left is the work of using them: more guards, and the
remaining three tools moving across. See `design/` for the reasoning and
[the issues](https://github.com/richarddavenport/tuikit/issues) for the state.

## Why an agent gets on with it

The interface has a written-down vocabulary, the vocabulary is machine-readable,
and the guards turn "I used a colour that does not exist" into a test failure
rather than a review comment. An agent can also *see* what it built:
`harness` renders any screen to a file, and the goldens tell it whether the
layout moved.
