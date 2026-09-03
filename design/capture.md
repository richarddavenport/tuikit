# Capture

How someone building a tuikit tool sees what they are building, and how an agent
sees it too. This is the framework's headline affordance; goldens and
documentation are downstream of it.

## Seeing what you are building

This is the framework's headline affordance, not a testing footnote. Capture
exists because someone building a screen wants to look at it — everything else it
enables (documentation, goldens, catching overflow) is downstream of that.

### The constraint that shapes it

pgctl's capture lives in a `_test.go` for a real reason: the model's fields are
unexported, so `go test` is the only thing that can reach into the package and
put the UI in a chosen state. Any "just run a command" design has to answer that,
and the honest answer is that two mechanisms are needed because there are two
questions.

| Question | Mechanism | Reaches |
|---|---|---|
| What does *this state* look like? | `harness.Shots` in a `_test.go` | Any state at all — an error, a half-loaded pane, a modal over a 40-row table, something no key sequence gets to yet |
| What does the *tool* look like? | `tool tui --snapshot out/ --script keys.txt` | Only states reachable by keystroke, but it is the real binary and needs no test |

The first is what you use while building. The second is what an agent finds in
`describe --json` and can run without knowing anything about the test suite.
The scaffolder generates the first, so a tuikit tool has capture from its first
commit rather than at the point someone thinks to add it.

### `tuikit watch` — the actual loop

Capture is only a building tool if it is fast enough to leave running:

```
tuikit watch ./internal/tui
```

Re-runs the capture on save, regenerates the page, serves it on localhost, and
reloads the browser. Write a line of view code, look at the frame. That is the
loop; a command you have to remember to type with two environment variables set
is not.

It also makes the agent loop symmetrical — the same directory of frames is what I
read to see what I just changed.

### `tuikit gallery` — the design system, running

A binary in the module that opens a real TUI of every component: pane, tabs,
table, filtered list, log pane, step list, form, confirm, toast, spinner, key
hint bar — each with its variants and states, its keybindings, and the colour
roles and glyphs it draws with.

Not a static component sheet. A sheet tells you what exists; a running gallery
tells you what it feels like to arrow through, what it does at 80 columns, and
what it looks like when the list is empty. Both a person picking a component and
an agent deciding how to build a screen start here.

This is the piece that makes tuikit a design system rather than a library of
widgets, and it is why `examples/democtl` is worth building early — the gallery
and the demo tool are close relatives.

## The agent layer

**This is no longer speculative.** pgctl commit `57a13ad` ("Capture the UI's
frames, and fix the three bugs that found", 2026-08-31) is a working prototype —
`internal/tui/screenshot_probe_test.go`, 308 lines, 17 frames. tuikit's `harness`
is a generalisation of that, not an invention.

Note what it is *not*: `swarmctl/internal/tui/capture.go` is **log** capture.
Screen capture existed nowhere until pgctl.

### What the prototype established

**Forcing the colour profile is the whole trick.** A test has no TTY, so lipgloss
strips every colour and you capture a grey rectangle:

```go
lipgloss.SetColorProfile(termenv.TrueColor)
lipgloss.SetHasDarkBackground(true)
```

This does not conflict with "indices, not truecolour" — the theme's *values*
stay palette indices, which is what ships. Forcing the TrueColor *profile* only
stops lipgloss from discarding them on a pipe.

A captured frame therefore holds the reader's own theme for the sixteen roles
(decision 28), and `harness.HTML` renders them as xterm's defaults. A published
frame shows what *a* terminal draws, not what yours does — which is the honest
answer for a page rendered where there is no terminal to ask.

**Fix the frame, freeze the clock.** `m.width, m.height = 132, 38` and
`m.now = <fixed time>` before every shot. Nothing else is reproducible.

**Two capture modes, and both are needed.**

| Mode | Data | For |
|---|---|---|
| Fixture | Constructed model | Goldens. Deterministic, no backend — so a tool can run these in CI if it wants to. |
| Live | Real backend, probed moments before | Documentation, and finding what fixtures hide. Run deliberately, by a person or an agent. Not a CI thing. |

The second is not a luxury. The header-overflow bug **was masked in the first
capture** because the fixture shortened the config path to `pgctl.yaml`; only a
real temp-dir path was long enough to overflow. Fixtures encode the author's
assumptions, which is exactly what a capture is supposed to catch.

**Frames go out as ANSI files**, one per screen, gated on an env var so an
ordinary `go test ./...` skips them:

```
PGCTL_SCREENSHOT_DIR=/tmp/shots go test ./internal/tui/ -run CaptureScreens
```

tuikit standardises this and also exposes it from the binary, so an agent needs
no test invocation: `tool tui --snapshot out/ --script keys.txt`.

### What capture is actually for

Seventeen frames, looked at once, found three bugs that **40-odd assertions had
not**:

1. The plan overlay drew wider than the terminal — the description can contain a
   list of every table a widened selection adds, and the box had no bound.
2. The header did the same with a long config path.
3. `"12 user triggers disabled for the load (1 tables)"`.

None are subtle. All are invisible to assertions that check content rather than
shape. This is the case for capture as a *development* tool and not merely
documentation.

It also found a bug in the capture harness itself — the "largest database" loop
compared against index 0 instead of the running maximum, so it picked vpic
(1.3 GB) over product-development (31 GB). Capture code gets tests like any other.

### What tuikit adds on top

- **`guard.Width`** — the invariant pgctl now tests by hand
  (`lipgloss.Width(line) > m.width`) generalised: every screen and every modal,
  rendered at 80 / 100 / 132 columns, failing on any line wider than the
  terminal. Overflow is a whole *class* of bug and two of the three found were
  instances of it.
- **ANSI → HTML, written once.** The prototype hit a real trap: an
  escape-stripping regex over `[A-Za-z]` eats the `m` that terminates every SGR
  sequence, so it strips the colour it is meant to preserve. That belongs in a
  library with a test, not re-derived per repo.
- **Text goldens.** The ANSI files are already diffable; stripped of colour they
  are golden files, and `harness.Golden(t, ...)` makes a layout regression an
  ordinary test failure. 23k lines of tests across the four repos and not one
  asserts on a rendered frame.
- **A frame log** — one entry per keystroke, so a whole flow is reviewable rather
  than a screen at a time.
- **Labelled provenance.** Every frame declares whether it is live or composed.
  pgctl marked the snapshot-manifest and plan/run frames composed because no
  production snapshot exists on that machine. A page of frames that quietly mixes
  the two is a page that misreports the tool.

### Publishing frames

`docgen` turns a capture directory into a page. Three rules the prototype got
right and that tuikit should encode rather than leave to judgement:

1. **Frames keep a dark ground in both light and dark themes.** The ANSI was
   captured for a dark terminal; recolouring it reports colours the tool does not
   have.
2. **System monospace stack, no webfont.** Box-drawing and braille glyphs must
   share one set of advance widths. A webfont plus a fallback for the glyphs it
   lacks guarantees misalignment.
3. **The chrome is near-silent.** The frames are the loud thing. Group plates by
   what the reader is doing — browsing / inspecting / acting — and give each the
   keystroke that reaches it.

### Two formats, because a page and a repository want different things

`tuikit frames <dir>` writes one self-contained HTML page: everything inline,
nothing to serve, open it and look. That is the right shape for the inner loop
and for sending someone a link.

`tuikit frames <dir> -md` writes Markdown with an SVG per frame, into
`<out>.md` plus a `<out>/` directory beside it. That is the shape for things
that live IN a repository — a README, an mkdocs site, a pull request — where a
self-contained page cannot go and a fenced block of raw ANSI renders as noise.

Both read the same capture and share the same section ordering, including the
rule that a frame no group named still appears. Two copies of that loop is how
one format grows the bug the other fixed.

**The images are SVG, and that is a decision rather than a default.** A PNG
needs a rasteriser, which needs a typeface, which is the open question in
decision 30 — writing documentation is the wrong reason to answer it. SVG names
the same system monospace stack the page does and lets the reader's machine
draw, so it needs no typeface of its own, keeps the text greppable, and scales.

Two properties it has to hold that a naive SVG writer does not:

- **Every span states its width** (`textLength`, with `lengthAdjust="spacing"`).
  A renderer that trusts the font's advance width to reproduce a character grid
  is one box-drawing character away from a frame that does not meet.
  `spacingAndGlyphs` would hit the width by distorting the glyphs, which bends
  the box-drawing characters instead of moving them.
- **No `<style>` element.** An SVG referenced from Markdown is rendered through
  a sanitiser, and a stripped stylesheet leaves a frame that is all one colour
  with no error to explain it. Presentation attributes survive; a stylesheet is
  a bet.
