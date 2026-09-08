# The wider TUI landscape — Go and Rust, read

> **Go and Rust only.** The pool is fourteen GitHub queries across those two
> languages, so Textual, Spectre.Console, FTXUI, Terminal.Gui, Ink and brick are
> all outside it — and htop, tmux and btop, which are among the most-used TUIs
> there are. `other-tuis.md` is the list that covers them, marked by how much of
> it has been measured. Check that before making a claim about "the field".

Read on 2026-09-01. This exists because an earlier survey answered "what are we
offering that these are not?" against **the wrong field**.

That one drew from a single GitHub full-text query — `--language=go bubbletea`
— so the five "competitors" it found were five one-person Bubble Tea side
projects, the largest with 18 stars. Its conclusion, that nobody had built the
layer above a cell buffer, came from a pool that structurally could not contain
`rivo/tview`, `mum4k/termdash`, or anything written in Rust.

That earlier survey is not published here. Its per-project readings are fair and
carefully evidenced, and they are also a close reading of five hobby repositories
by name, which is not a thing to publish about projects that size. What it
concluded is corrected below; what it taught is in `design/decisions.md`.

This survey fixes the pool. Fourteen queries across both languages
(`tui`, `terminal ui`, `terminal user interface`, `ratatui`, `tview`, `tcell`,
`gocui`), deduped to **903 repos**; substrate detection by reading each top
repo's `go.mod`/`Cargo.toml`, recorded in `substrates.tsv`.

The raw pull is not kept here — it is a scrape of public metadata that is stale
within a month, and the finding rather than the fixture is what is worth
reading. `substrates.tsv` is the evidence for the tables below.

Everything below came from reading source in a clone, not from READMEs.

## What the earlier survey missed

**The two most popular Go TUIs in the world are not Bubble Tea apps.**

| Stars | | Substrate |
| ---: | --- | --- |
| 81,858 | [jesseduffield/lazygit](https://github.com/jesseduffield/lazygit) | tcell (direct) |
| 54,522 | [wagoodman/dive](https://github.com/wagoodman/dive) | gocui + tcell |
| 44,740 | charmbracelet/bubbletea | — |
| 41,856 | [sxyazi/yazi](https://github.com/sxyazi/yazi) | **ratatui** |

lazygit alone has more stars than Bubble Tea and every framework in the previous
comparison combined, and the previous survey never saw it.

Of the top 80 TUI repos, excluding the libraries themselves:

| | Go apps (26) | | Rust apps (41) |
| --- | ---: | --- | ---: |
| bubbletea | 13 | **ratatui** | **31** |
| tcell (incl. via tview) | 11 | `tui` (ratatui's dead predecessor) | 2 |
| — of which tview | 5 | cursive | 2 |
| termui / termdash / gocui | 3 | tui-realm | 1 |

**Rust is not a sideshow; it is the more consolidated half of the field.**
`ratatui` has 22,468 stars and **47.6M crates.io downloads** (17.2M in the last
90 days). It is a single answer where Go has four. Any claim of the form
"nobody has done X in a TUI framework" has to survive ratatui first.

## The substrates that actually matter

| | Stars | Lang | Last commit | Test files | Draw model |
| --- | ---: | --- | --- | ---: | --- |
| [ratatui](https://github.com/ratatui/ratatui) | 22,468 | Rust | 2026-08-31 | ~519 inline buffer asserts | **cells** |
| [tview](https://github.com/rivo/tview) | 14,074 | Go | 2026-08-11 | **0** | **cells** (via tcell) |
| [termui](https://github.com/gizak/termui) | 13,583 | Go | 2025-07-10 | **0** | cells |
| [gocui](https://github.com/jroimartin/gocui) | 10,596 | Go | 2025-05-01 | **0** | cells |
| [cursive](https://github.com/gyscos/cursive) | 4,843 | Rust | 2026-07-14 | inline | cells |
| [termdash](https://github.com/mum4k/termdash) | 3,035 | Go | 2026-06-09 | **99** | **cells** |
| [iocraft](https://github.com/ccbrown/iocraft) | 1,530 | Rust | 2026-08-19 | inline | cells + flexbox |
| [tui-realm](https://github.com/veeso/tui-realm) | 994 | Rust | 2026-07-29 | 34 **committed** `.snap` | cells (ratatui) |
| **tuikit** | — | Go | — | 7.4k test LOC | **cells** |

**The first correction to the earlier doc: a cell buffer is not distinctive.**
It is the majority position in the field and has been since termbox. Bubble
Tea's `View() string` is the *outlier*, and the previous survey mistook the
outlier's neighborhood for the world.

## The finding that survives

Not "cells." **Cells that record who drew them.**

Every framework here hit-tests by *recomputing* where things were drawn. Draw
runs once and throws the mapping away; the click handler runs the layout a
second time, by hand, and the two are kept in agreement by the author's
attention. It is the same defect in five codebases.

**cursive** — `Button` centers its label in `draw`:

```rust
let offset = HAlign::Center.get_offset(self.label.width(), printer.size.x);
```

and centers it again, separately, in `on_event`:

```rust
let self_offset = HAlign::Center.get_offset(width, self.last_size.x);
...
} if position.fits_in_rect(offset + (self_offset, 0), self.req_size()) => {
```

The same expression twice, reading two different size sources
(`printer.size.x` vs `self.last_size.x`).

**tview** — `List.indexAtPoint` reverses the draw loop's arithmetic:

```go
index := y - rectY
if l.showSecondaryText {
	index /= 2
}
index += l.itemOffset
```

Every layout decision `Draw` made — the two-line item, the scroll offset — is
re-encoded here. `Table.CellAt` goes further and replays column widths out of
`t.visibleColumnWidths` / `t.visibleColumnIndices`, fields that exist *only*
so the hit-test can rebuild what `Draw` already knew. And the guard at the top
of `indexAtPoint` is wrong —

```go
if rectX < 0 || rectX >= rectX+width || y < rectY || y >= rectY+height {
```

`rectX >= rectX+width` is false for any positive width, and `x` is never
tested. It is harmless only because `MouseHandler` checked `InRect` first. Nine
years, 14k stars, and the bounds check in the hit-test is dead code nobody
noticed, because it is the *second* copy of something that was already right.

**ratatui — no hit-testing at all.** `MouseEvent` appears in three examples and
in **zero** of its 20-odd widgets. The whole public surface for the problem is
`Rect::contains(Position)`. `ListState` and `TableState` hold `offset` and
`selected`, and mapping a click to a row is left entirely to you — after the
constraint solver has moved the rects. Forty-seven million downloads, and the
answer to "which row did the user click?" is still "do the arithmetic yourself."

**tview and termdash** are the two that route mouse events properly through a
tree, and both stop at the widget boundary: the framework tells a widget *that*
it was clicked, and the widget works out *where* on its own. termdash's
`MouseScopeWidget`/`Container`/`Global` is the most considered version of this
in either language, and the last mile is still the widget's problem.

**iocraft comes closest to the right idea.** Its event hook stashes the box the
component actually occupied, at draw time:

```rust
fn post_component_draw(&mut self, drawer: &mut ComponentDrawer) {
    self.component_location = (drawer.canvas_position(), drawer.size());
}
```

then filters mouse events against it and hands the component *local*
coordinates. Position observed from the render rather than recomputed — the
correct instinct, one granularity too coarse. It is still one rect per
component, so a `Table` learns nothing about rows; and because every hook polls
the same event stream independently, **every** component whose box contains the
point is notified. There is no topmost-wins. An overlay does not occlude what is
under it.

That is the field's ceiling: bubbleapp scanned ownership back out of a rendered
string, iocraft snapshots a rect at draw time, tview keeps a shadow copy of its
column widths. Three people reaching for the same missing thing.

`comp.Cell{Text, Style, Owner}` is that thing, at cell granularity, and it
remains ours after the pool was widened correctly.

## Where the widened field puts us behind

Being wrong about the pool cuts both ways. Three candidates; **one survives
reading our own source**, and it is not one of the three.

### 1. ratatui's constraint layout — mostly already taken, and the rest declined on purpose

`comp/layout.go` has `Length`, `Percent`, `Fill(weight)`, `.Min(n)`, `.Max(n)`,
and a pin-and-resolve loop that re-runs after a Min or Max moves a band. The
file states the rejection outright:

> Measured against what four tools actually lay out, which is a column of bands
> and a pane split in two, a single linear pass is enough. [...] Flex alignment
> and negative spacing exist in ratatui because ratatui is a general library; we
> are not one.

What is genuinely absent: `Ratio`, `Flex` alignment for the under-filled case,
negative `Spacing` for shared borders, and nesting (two axes means `Rows` then
`Cols` by hand). All four are general-library features. The two-of-four rule
covers this; **accept, revisit when a second tool needs one.**

### 2. termdash's mouse scopes — a question owner-cells do not ask

`MouseScopeNone|Widget|Container|Global` exists because termdash routes events
down a container tree and must decide who gets one. `app.Mouse.Route` does not
route: `c.OwnerAt(x, y)` names the region, `Handler.Blocked` takes the mouse for
a modal, and the app dispatches. There is no tree, so there is no scope to
declare.

And on the axis that matters we are ahead of both of them: **a drag owns the
mouse until release** (`Mouse.on`, checked before `Blocked`). termdash has no
drag at all; tview exposes `MouseLeftDown`/`Up` and leaves capture to you.
**Not behind. Withdraw the claim.**

### 3. tui-realm's committed snapshots — we already do this

105 files under `testdata/`. tui-realm's 34 `.snap` make it the only *other*
project doing it. The correction is to the uniqueness claim, not to us.

### 4. The real gap: nothing happens when the pointer moves

Not on the original list, because it comes from reading the field rather than
comparing architectures.

We handle `MouseActionMotion` **only while dragging**. There is no hover, and
`design/mouse.md` never considers one — it is unexamined, not rejected. The
field mostly has it: tview has `MouseMove` and four `*DoubleClick` actions;
bubbleapp's theme has a `Hover` state per component and scans zones back out of
the frame to find it; cursive tracks a hovered view.

Hover is the affordance that tells a user a thing is clickable before they
click it, and in a TUI that is doing more work than in a GUI, because there is
no cursor change to fall back on. We also have **no double-click**, which is the
conventional "open" gesture in every file-manager-shaped tool.

Both are cheap here for the same reason everything else is: `OwnerAt` already
answers "what is under the pointer" for free, and motion events are already
being delivered (`WithMouseCellMotion` is on for drags). The cost is a redraw
per motion event, which is the thing to measure before committing.

## Revised: what is actually ours

Against 903 repos rather than five:

1. **The CLI/TUI pair from one declaration.** Still nothing. Not one substrate in
   either language has an opinion about flags, subcommands, exit codes, or
   `describe --json`. `spec.Command` remains an axis the field is not on.
2. **Cells that record their owner.** Holds, and is now *better* evidenced —
   ratatui declines the problem, tview and cursive solve it twice by hand,
   iocraft solves it at the wrong granularity.
3. **Conventions as test failures.** `guard.Tokens`, `guard.Glyphs`,
   `guard.Engine`, `guard.Reachable`. Nothing in 903 repos has a glyph
   allow-list or a reachability check. tview has no tests at all; ratatui has
   `clippy.toml` and `typos.toml`, which police the source, not the frame.
4. **Semantic color roles.** Weaker than claimed. cursive's `PaletteStyle` (13
   named roles, themable) and tview's `Theme` (11) are real palettes. Ours is
   distinguished by being *load-bearing* — `guard.Tokens` fails a build that
   writes a raw color — not by existing.
5. **Extracted from four production tools.** Every framework here is a framework
   first. lazygit and dive are the opposite — applications that never extracted.
   Nobody is standing where we are.

## What to take, revised

The earlier list was right and the reasons are now stronger.

1. **Constraint layout** (ratatui, not bento). Read `ratatui-core/src/layout/`
   directly.
2. **A screen stack with typed values** (soda). Unchanged — and note that the
   widened field offers nothing better. cursive uses a `StackView` of boxed
   views; tview a `Pages` map of names. Neither returns a value.
3. **Hover and double-click** (tview, bubbleapp). The one place the field has an
   affordance we do not, and the one gap here that is a feature rather than an
   architecture. `OwnerAt` already answers it; measure the redraw cost first.

Not taken: **ratatui's stateless widgets.** `render(area, buf)` with all state
in the caller's `ListState` is what makes its mouse story impossible — the
widget cannot remember where it drew a row because the widget does not persist.
That is the fork in the road where we go the other way, deliberately.

## Method note

`gh search repos` caps at 100 results per query and ranks by stars, so the tail
below ~2,500 stars is sampled, not enumerated. Substrate detection reads the
manifest at the repo root only, so workspaces that declare deps in a member
crate show `unknown` (6 Rust, 2 Go). Neither affects the shape of the finding.
