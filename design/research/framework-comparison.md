# The five app frameworks, read

Read on 2026-09-01, from the framework/app-shell bucket of
[bubbletea-ecosystem-buckets.md](bubbletea-ecosystem-buckets.md). The question
was Richard's: **what are we offering that these are not?**

Every claim below came from reading the source, not the README. Where a README
and the code disagree, that is noted, because it happened three times.

| | Stars | Last commit | Tests |
| --- | ---: | --- | ---: |
| [metafates/bento](https://github.com/metafates/bento) | 17 | Mar 2025 | 4 files |
| [alexanderbh/bubbleapp](https://github.com/alexanderbh/bubbleapp) | 18 | Jun 2025 | **0** |
| [newbpydev/bubblyui](https://github.com/newbpydev/bubblyui) | 17 | Jan 2026 | 4,816 |
| [tuiphi/soda](https://github.com/tuiphi/soda) | 7 | Feb 2024 | **0** |
| [metafates/soda](https://github.com/metafates/soda) | 5 | Jan 2024 | **0** (archived) |

All five are one person's project. All five are dormant or nearly.

## The shape of the field

**They split into two camps, and tuikit is in neither.**

**bento is a rendering substrate** — a faithful `ratatui-core` port, with no app
opinions at all. It does not even depend on Bubble Tea; it reimplements the
runtime inline.

**The other four are app frameworks over strings.** bubbleapp is React (hooks,
`func(ctx, props) string`), bubblyui is Vue (`Ref`/`Computed`/`Watch`,
`provide`/`inject`), the two sodas are a screen stack. Every one of them
composes with `lipgloss.JoinVertical`.

tuikit is neither: an opinionated shape for **one kind of tool** — an operator
CLI *and* TUI over a domain engine — sitting on its own cell substrate. Nothing
in the list is trying to be that.

| | bento | bubbleapp | bubblyui | soda ×2 | **tuikit** |
|---|---|---|---|---|---|
| Draw model | **cells** | strings | strings | strings | **cells** |
| Cell records who drew it | no | n/a | n/a | n/a | **yes** |
| Hit-testing | none | bounding rects scanned from ANSI markers | **none** | none | by owner ID |
| Drag | no | no | no | no | yes |
| Layout | **Cassowary, ratatui constraints** | 4-pass grow | flex-ish over strings | 14 lines of alignment | `Split` + rects |
| Committed golden frames | 0 | 0 | **0** | 0 | 96 |
| Semantic colour roles | none | **best of the five** | small `Theme` | 8 chrome slots | 9 roles |
| Glyph allow-list | no | no | no | no | **yes** |
| CLI surface | none | none | none | none | `spec` |
| Scaffolder | none | none | none | none | `tuikit new` |

## bento — the substrate, and where it stops

The thesis is two declarations:

```go
type Widget interface { Render(area Rect, buffer *Buffer) }
type Model  interface { Widget; Init() Cmd; Update(msg Msg) (Model, Cmd) }
```

**`Model` embeds `Widget`. There is no `View() string` at all** — the view step
of TEA is literally a render-into-buffer, so nothing can smuggle raw ANSI past
it, and the top-level app and every nested widget share one contract.

It is a real port: `Cell{Symbol string, Fg, Bg, Modifier, Skip}` with grapheme
clusters, `Buffer.Diff` producing `[]PositionedCell` for minimal repaint, and
`internal/casso` — a Cassowary solver — with ratatui's strength constants
matched exactly (`_lengthSizeEq = Strong * 10.0`, and so on) behind
`Min/Max/Length/Percentage/Ratio/Fill`, `Flex.SpaceBetween`,
`Spacing.Overlap`.

And it takes **none** of the leverage a cell buffer creates:

- **No hit-testing.** A cell stores its symbol and style and no owner. Worse,
  mouse reporting is never switched on: the SGR parser is copied wholesale from
  Bubble Tea, `// a.disableMouse()` is commented out in `tty.go`, and no widget
  references `MouseMsg`. It is dead code.
- **No snapshot testing** — the one thing a cell buffer gives you free. No
  `testdata/`, no goldens. `TerminalBackend` is a 23-method interface with a
  single implementation: the seam exists and nobody used it.
- **No design vocabulary.** `Color = termenv.Color`, a type alias, so anything
  satisfies it.
- **No CLI surface, no scaffolding.**

**bento is evidence the substrate works in Go, and evidence that nobody has
built the layer above it.** It did not try ownership-based hit-testing and
reject it. It stopped.

## bubbleapp — the honest one

Zero tests, abandoned stubs in-tree (`component/modal` returns `""`;
`component/grid/grid.go` is 394 lines commented out under
`// Disclaimer: The size calculations are mostly AI vibed.`). But two ideas are
worth recording.

**Mouse zones as component identity.** Using a fork of BubbleZone, a component
wraps its output in invisible ANSI markers; after render, `zone.Scan()` walks
the finished frame and recovers each marker's bounds. `id###childID` lets one
component own many sub-regions — table rows, tab titles — without a registry.

That is *the string-world approximation of a cell buffer that records
ownership*, and its limits are instructive: `InBounds` is a plain box, so
wrapped or non-rectangular content mis-hits, and overlap is resolved by a
heuristic the source labels honestly:

```go
// Longer ID means it is more specific here and thus hovered.
// Not sure if that is valid always
```

Somebody wanted ownership badly enough to scan it back out of the finished
frame. That is the argument for putting it in the cell.

**The best theme of the five.** A Tailwind ramp behind named semantic roles
(`Primary/Secondary/Tertiary/Info/Danger/Success/Warning`, each with
`Light/Lighter/Dark/Darker/Bg/Fg`), then per-component state tables:
`Button map[Variant]map[ComponentState]lipgloss.Style` over
`Normal|Hover|Focus|Disabled`. Richer than our nine roles on the state axis —
and it still hardcodes `⟨` in a button, because nothing stops it.

## bubblyui — the volume trap

672 commits in ten weeks, 855 Go files, 4,816 tests, 146 markdown files, CI,
codecov, a CHANGELOG and a `CLAUDE.md`. `go build ./...` and `go test ./pkg/...`
pass. And the headline claims do not survive contact:

- **"Type-safe reactivity"** — the README's own counter is
  `count.Set(count.GetTyped().(int) + 1)`, and `ctx.Get("theme").(Theme)`
  assertions are everywhere.
- **Mouse** — the only mouse code in the framework is passing
  `WithMouseCellMotion` through to Bubble Tea. No component handles
  `tea.MouseMsg`. `Button.OnClick` fires from the event bus; *something else*
  has to decide a click happened, and nothing does.
- **Snapshot testing** — there is a real `SnapshotManager` with diffing,
  normalizers and a written guide. `find . -name '*.snap*'` returns **nothing**.
  Not one component's appearance is pinned. The snapshot tests exercise the
  manager against `t.TempDir()`.
- **"30+ built-in components"** — this one holds. 31 exported constructors,
  with tests, that render.

Two ideas are genuinely good. **An MCP server exposing the running app** —
`bubblyui://components`, `bubblyui://state/refs`, `bubblyui://events/log`, with
subscriptions and a `setref` tool to poke live state from outside. And
**`HelpText()` generated from registered `KeyBinding`s**, so help cannot drift
from the bindings.

## The routers

### bubblyui's — a warning

~5,000 lines of source, ~10,000 of tests, and the best-engineered corner of that
repo: `:id`/`:path*`/`:id?` patterns compiled to regex, specificity scoring so
`/user/me` beats `/user/:id`, nested routes producing a `Matched` chain,
`BeforeEach`/`AfterEach`/per-route guards with panic recovery and
circular-redirect detection, a history cursor with forward.

Then:

- **Routed components are permanently-mounted singletons.** Navigation never
  calls `Init`, `Unmount` or any lifecycle hook. `View()` simply stops being
  called on the unmatched ones. A route's component cannot be parameterised per
  navigation, so the example closes over a `**Router` and reads
  `CurrentRoute().Params["id"]` inside its template.
- **`HistoryEntry.State`** — the slot for restoring scroll position — exists, is
  tested, and nothing ever populates it. `PushWithState` has no caller.
- **Zero usage in the component library.** No component is route-aware.
- **`KeyBindings()` returns nil at the router boundary**, so the framework's
  nicest idea dies exactly where navigation begins.
- **You cannot list an app's routes from outside.** `RouteRegistry.GetAll()` has
  zero non-test callers, the MCP server exposes no route resource, and
  `RouterDebugger.RecordNavigation` is never called by anything.
- No test does *Push → assert rendered screen*.

Fifteen thousand lines routing nothing the framework knows about.

### soda's — the mechanism worth taking

A screen **is a value with a constructor**, so parameters are typed and free:

```go
return soda.PushState(New(s.n + 1))
```

No path, no params map, no serialisation. The API is four commands — `PushState`,
`PushTempState`, `Back`/`BackN`, `BackToRoot` — and `SaveToHistory` lives on the
*outgoing* state, so a modal pushed with `PushTempState` never becomes a
back-stop. Lifecycle is explicit: leaving calls `Destroy()` and cancels the
screen's `context.Context`; returning calls `Init(ctx)` on the same value.

**And a screen returns a value by calling forward, not by returning.** The whole
feature costs soda zero lines:

```go
filepicker.WithOnSelect(func(path string) tea.Cmd {
    contents, err := os.ReadFile(path)
    if err != nil { return soda.SendError(err) }
    return soda.PushState(viewport.New(string(contents)))
})
```

The chosen value goes to a closure the *parent* supplied when it built the
picker, and that closure decides what happens next — `soda.Back` to unwind, or
`PushState` to go deeper. The trade-off is real: the parent must build eagerly
and pre-commit to what it does with the answer, so "push any picker, get a value
back" is not expressible generically.

## What is actually ours

1. **The CLI.** Not one of the five has anything to say about flags,
   subcommands or exit codes. `spec.Command` → CLI command + TUI screen +
   `describe --json` + context menu is an axis nobody else is on.
2. **Cells that know who wrote them.** bento has cells without ownership;
   bubbleapp wanted ownership badly enough to scan it out of a rendered string.
   Each has half.
3. **Committed goldens.** 96 frames at two widths, plus a colour-shape check.
   bubblyui built the machinery, wrote the guide, and committed none.
4. **Conventions as test failures.** Nothing else has `guard.Tokens`,
   `guard.Glyphs`, `guard.Engine` or `guard.Reachable` — or a glyph allow-list
   at all.
5. **Extracted from four production tools rather than designed greenfield.**
   Every one of these is a framework first and an application never.

## What to take

Ranked, with the reason each is worth doing rather than admiring.

1. **Delete `View() string`** (bento). `Model` embeds the thing that renders into
   the canvas. Today `View()` returns `c.String()`, which is a seam: a tool could
   hand back a hand-joined string and no guard would catch it. Closing it makes
   "the tool never computes a coordinate" enforceable rather than aspirational.
2. **A screen stack with history** (soda). `app.Screens` calls itself "the
   router" and is a flat map. democtl's `esc` is a hardcoded `screenDashboard`,
   correct only because every screen happens to be reached from the dashboard;
   azctl's runner is a second `tea.NewProgram` you cannot come back from. Take
   soda's mechanism — typed values, `PushTemp`, `Destroy` on pop, results by
   closure — and add a `spec.Call` **label** per entry, so an agent can read
   where it is and a capture script can say `open logs api_gateway`. About a
   hundred lines. bubblyui is the argument for keeping it that small.
3. **Constraint layout** (bento). `Min/Max/Length/Percentage/Fill` over a rect
   would cover every layout in the four tools and end the remaining arithmetic.
   The one place the field is genuinely ahead of us.

Not taken, and why: **URL-space routing** — bubblyui proves it buys nothing in a
Go TUI, since no component is route-aware and the route table is not readable
from outside anyway. **Reactivity** — our shape is a snapshot redrawn, not a
graph of signals. **A component-count race** — bubblyui's 31 are real and
thinner than our 14.

## One convergence worth noting

bento independently arrived at

```go
type TryUpdater interface { TryUpdate(msg Msg) (cmd Cmd, consumed bool) }
```

for input routing without a focus tree. That is `app.Handled` exactly. Two
people solving it the same way, separately, is the signal this project runs on.
