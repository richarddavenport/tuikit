# tuikit

A TUI framework for developers and agents. Go, Bubble Tea, Lip Gloss.

It exists because four tools — [swarmctl], [pgctl], [azctl], [dugo] — arrived at
the same shape independently, and that shape currently travels by copy. tuikit is
that shape as a module: the vocabulary an interface is allowed to use, the guards
that hold it closed, and (soon) the components, the capture harness, and a
scaffolder.

[swarmctl]: https://github.com/richarddavenport/swarmctl
[pgctl]: https://github.com/richarddavenport/pgctl
[azctl]: https://github.com/richarddavenport/azctl
[dugo]: https://github.com/richarddavenport/dugo

## What is here now

**`theme`** — the vocabulary. Nine colour roles named by role and never by hue,
as ANSI 256 indices; a closed glyph allow-list; and `Hex` for anything drawing
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

`tuikit new` does not exist yet — it is [#13], and it comes last on purpose, so
that it generates what a real migration turned out to need rather than what
seemed likely beforehand. Until then the supported path is copying the example,
which is what the scaffolder will generate anyway:

```sh
cp -r examples/democtl ../mytool
cd ../mytool
go mod init github.com/you/mytool     # then fix the import paths
```

`examples/democtl` is a complete working tool, not a snippet: a fictional
service fleet with a dashboard, a log pane, a confirm modal and a step run, in
about 1,700 lines. Delete the screens you do not want.

**What to keep, and why:**

| Keep | Because |
|---|---|
| `fleet/` and `ui/` as separate packages | The engine knows the domain and has **no UI imports**; the UI never calls the domain directly. Every tool in the family keeps this split, and it is what makes both halves testable. |
| `ui/guard_test.go` | Six lines. It holds your palette and glyph set closed from your first commit rather than from whenever someone thinks of it. |
| `ui/style.go` | Styles built from a `theme.Palette` rather than declared at package level, so changing the palette changes the interface without editing the file. |
| `capturesKeys` in `ui/model.go` | The mode split: while a modal is open or a filter is being typed, `j` is the letter j. Without it your list scrolls behind the dialog, and nobody finds that in review. |
| The generation counter on `stepDoneMsg` | An async result from something the user walked away from must not draw into its successor. Easy to write, almost impossible to see once written. |
| The `View()` switch with a loud `default` | A screen constant without a case renders as an empty terminal and says nothing about why. |
| `fleet.Epoch` and the seeded data | Only if you want capture. A frame reading "47s ago" has to read that tomorrow, or every golden fails the day after it is written. |

`examples/democtl/README.md` says what each part is a reference for, and lists
the two bugs its own frames caught.

[#13]: https://github.com/richarddavenport/tuikit/issues/13

## What is coming

`comp` (the
components each of the four wrote separately), `app` (the Bubble Tea shell and
its async conventions), `spec` (one command declaration → CLI, TUI screen, and a
`describe --json` manifest an agent reads), `tuikit watch`, `tuikit gallery`,
`tuikit new`. See `design/` for the reasoning and
[the issues](https://github.com/richarddavenport/tuikit/issues) for the state.

## Why an agent gets on with it

The interface has a written-down vocabulary, the vocabulary is machine-readable,
and the guards turn "I used a colour that does not exist" into a test failure
rather than a review comment. An agent can also *see* what it built:
`harness` renders any screen to a file, and the goldens tell it whether the
layout moved.
