# Decisions

Newest last. A decision is here when it closed off an alternative someone would
otherwise reach for.

## 1. The framework is a library and a scaffolder, split on mechanism vs policy

Four tools — swarmctl, pgctl, azctl, dugo — already share a shape, and it travels
by copy. A fix to a pane lands in one repo.

**Mechanism goes in the library.** `guard.Tokens`, `guard.Glyphs`, `theme`,
`docgen` are functions. They know Go source, terminal columns and ANSI, and they
know nothing about Postgres or Docker Swarm. They do not run anywhere; the tool
calls them.

**Policy goes in the scaffolder.** `tuikit new` will write the tool a test that
calls the guards and a workflow that runs it. That is an opinion about setup,
delivered as generated code the tool owns and can edit on day one.

**Backend-shaped things belong to the tool** — its fixtures, whether it captures
live frames at all, what has to be standing up for that. No answer tuikit could
give would be right for pgctl (needs a database), swarmctl (needs a swarm) and a
backendless tool at once.

The corollary: tuikit's CI tests tuikit. A consuming tool's CI is that tool's
business.

## 2. Roles are a closed struct with an escape hatch, not a map

`theme.Palette` has nine named fields, so a tool takes `theme.Default` and
overrides what it wants. A map of names would make every lookup a question of
whether the key exists.

`Extra` exists because the alternative is worse. A tool that needs a tenth
meaning and cannot name it reaches for `lipgloss.Color("39")` instead — and
`guard.Tokens`, the thing that keeps colours honest, is the first casualty. It
should still feel like a decision: a shade of an existing role is not one.

`Roles()` returns the nine in a fixed reading order with `Extra` appended, and
the `Why` text describes the role rather than the hue, so it survives a tool
recolouring the palette. That is the whole reason roles are named.

## 3. ANSI indices, never truecolour hex

> **Narrowed by decision 28 (2026-09-01): the first SIXTEEN indices, not 256.**
> The reasoning below is unchanged and still the reason not to use hex; what it
> got wrong is that it treated all 256 as equivalent. They are not — only the
> first sixteen follow the reader's theme.

Inherited from swarmctl. A truecolour hex looks right on the machine it was
picked on and wrong over ssh from another; the 256 palette is what every terminal
has agreed on. `theme.Hex` converts for anything drawing outside a terminal, and
is computed from the palette formula rather than tabled — a table of nine
hand-copied values is nine chances to copy one wrong, and wrong here means a
design system quietly describing a different tool.

Capture will force lipgloss's *TrueColor profile* so those values are not
stripped when there is no TTY. That is not a contradiction: the values stay 256.

## 4. A GlyphSet is copied on extension, never mutated

`GlyphSet` is a map, and a map is a reference. A tool adding a glyph to
`DefaultGlyphs` in place would add it to every other tool in the process — and to
the guard that is supposed to catch it. `With` returns a copy, and panics on a
malformed call because that is a typo at the call site rather than a condition to
handle.

## 5. The guards take an interface, not `*testing.T`

`guard.T` is `Helper`/`Errorf`/`Fatalf`. `*testing.T` satisfies it, and so does a
recorder — which is what lets tuikit test that the guards actually fire. A guard
that never fires is indistinguishable from a clean package, and "the rule is
enforced" is exactly the claim worth checking.

For the same reason a guard pointed at a directory with no Go files calls
`Fatalf` rather than passing. Scanning nothing silently is how everyone comes to
believe a rule is on when it is not.

## 6. The guards parse, they do not match

A regex over source cannot tell a string literal from a comment that quotes one,
so a comment *explaining* a glyph got reported as printing it. `go/parser` hands
back only the literals. The test for this is kept as its own case, because the
bug is easy to reintroduce by "simplifying" the guard.

Role usage is matched on a selector — `\.Accent\b` — rather than a literal
`theme.` prefix. A tool may alias the import or hold a palette value of its own;
what is constant is that a role is reached through a dot.

## 7. The design system is generated, never maintained beside the code

A design tool draws pixels; a TUI draws a character grid. A mockup made without
the column budget, the 256 colours and the glyph allow-list is a picture of a
tool that cannot be built. So `docgen` renders from `theme`, and a test fails if
any terminal block in the bundle draws a character outside the set.

Only the foundations exist so far — colours, glyphs, the grid. Component and
screen cards arrive with the components. A card depicting a component that does
not exist is exactly the drift this package prevents.

`docgen` draws every role in a line of interface, generated from the palette
rather than written out, so an `Extra` role appears without editing the page and
a role cannot satisfy the guard with a swatch alone.

## 8. The name is tuikit, not ctlkit

`ctlkit` echoed the family it came from — swarmctl, pgctl, azctl — but the `-ctl`
suffix is a claim about what a tool *does*: it controls something. The framework
controls nothing. What it is about is the terminal interface, which is what the
four have in common and the only thing this module knows how to help with.

Renamed before anything imported it, which is the only cheap time to do it.

## 9. Capture is a building tool first

Its purpose is seeing the screen you are writing. Goldens, documentation and
catching overflow are all downstream of that, and treating any of them as the
point produces a harness nobody runs. It is why `tuikit watch` matters more than
any assertion the harness could offer: a command you have to remember to type
with two environment variables set is not a loop.

## 10. Two capture mechanisms, deliberately

A test helper reaches any state at all — an error, a half-loaded pane, a modal
over a 40-row table, something no key sequence gets to yet — because the model's
fields are unexported and `go test` is the only thing that can reach inside the
package. A `--snapshot` flag on the binary needs no test and is discoverable by
an agent from `describe --json`, but only reaches what a keystroke reaches.

Neither replaces the other, and pretending one does means either losing the
states that matter most or making an agent read the test suite to take a picture.

## 11. Extract from swarmctl, prove on azctl

swarmctl is 65k lines and mid-flight, so it is the source rather than the first
migration. azctl is 3.2k lines and the youngest, so it has the least to lose and
is the honest test of whether the framework fits a real tool. A framework that
has never met one is a guess.

## 12. No cobra

Once `spec` is the source of truth for a command, cobra is a second description
of the same tree and the work becomes adapting one into the other. That costs two
things worth keeping: swarmctl's exit-code contract — `0`, `1`, `2` meaning "a
dry run found drift", and a script's own status passed through — and the
dual-mode entry where `main` decides whether an argument is a subcommand or the
start of the TUI's flags. Cobra wants to own `os.Exit` and returns 1 for
everything.

What cobra would have given free — completions, generated help — generates from
`spec` directly. dugo drops it when it migrates.

## 13. ANSI files and HTML, not SVG

The ANSI file is both the archive and the golden: stripped of colour it diffs in
a pull request, and kept whole it is exactly what the terminal emitted. HTML with
real text keeps a published frame selectable and searchable. SVG would be a third
representation of the same frame, earning nothing and drifting from the other two.

## 14. The local check runs the same linter CI does, pinned, via `go run`

The first push failed CI on a lint config the runner could not load, because
`make lint` skipped golangci-lint when it was not installed and said so in a line
nobody reads. A local check that quietly omits what CI enforces is a check that
lies.

`make lint` now runs `go run github.com/golangci/golangci-lint/v2/cmd/...@v2.12.0`,
so the version is the same whether or not anything is installed, and there is no
path where the step is skipped.

The pin itself is inherited from swarmctl, which hit the same wall and left the
comment: golangci-lint-action's `version: latest` resolves to the v1 line, built
with go1.24, which refuses to run against a `go 1.25.0` directive rather than
degrading. The action must be `@v8` and the version an explicit v2.

This is the copy-by-hand problem the scaffolder exists to end: swarmctl knew the
answer, and tuikit rediscovered it by breaking.

## 15. democtl is seeded and clock-frozen, not merely backendless

The example tool needs no database or swarm, which is convenient. What matters
more is that it is *deterministic*: `fleet.New(seed)` returns the same services,
`Logs` the same lines, `Deploy` the same timings, and `fleet.Epoch` fixes the
clock so a frame reading "47s ago" still reads that tomorrow.

The harness renders democtl's screens as tuikit's own fixture. A fixture that
changes between runs is not one, and a golden that fails the day after it is
written teaches everyone to ignore goldens.

The consequence worth stating: nothing in `fleet` may call `time.Now`, iterate a
map into output, or use `rand` without a source. Its test asserts all three by
comparing two independently generated fleets.

## 16. The example is a tool, not a widget showcase

democtl manages a fictional fleet with a dashboard, logs, a modal and a step run
— the four shapes swarmctl, pgctl, azctl and dugo all have. It is deliberately
not a gallery of components in a grid.

A showcase demonstrates that a widget renders. A tool demonstrates the things
that actually go wrong: a modal that has to bound itself to the terminal, an
abandoned run whose steps must not draw into its successor, a filter that must
stop `j` from scrolling the list behind it. Those are what a reader — human or
agent — needs the reference for, and none of them appear in a widget grid.

`tuikit gallery` is still worth building, and it will be seeded from democtl's
components. But the components come from a working tool first, the same way
`theme` came from swarmctl rather than from a palette designed in the abstract.

## 17. Writing the example before the components is the right order

democtl's panes are hand-drawn against `theme`, because `comp` does not exist
yet. That is not work to be redone — it is the material `comp` gets extracted
from, and extracting from two or three real usages is what stops a component
library from being a set of guesses with configuration options nobody needs.

The evidence arrived immediately: democtl's first frames contained two bugs, and
both are the kind a component would have to solve properly rather than a widget
would paper over — an overlay that composites rather than replacing lines, and a
right-aligned column measured against the box rather than the terminal.

## 18. Components draw cells, not strings

Settled by the prototype on `prototype/canvas-mouse` (`1cdd3f8`), and written up
in [mouse.md](mouse.md).

A component that returns a string cannot be clicked: hit-testing needs geometry,
and `lipgloss.JoinHorizontal` throws it away. So `comp` sits on a cell grid where
each cell carries a rune, a style and the ID of whatever drew it. `OwnerAt(x, y)`
is the whole hit test, and there is no region list that can drift from the
drawing because the drawing is the region list.

The substrate cost 153 lines. It deletes `clip`, `pad`, `padVisible`, `rule` and
`composite`, and turns a 28-line compositing overlay into an 8-line draw. Lip
Gloss keeps styling; its layout helpers go.

Two consequences beyond the mouse. Style lives on the cell, so nothing parses
ANSI — the trim-versus-clip distinction the grid page documents stops existing.
And `Set` clips to the canvas, so drawing past the edge is not an error to catch
but a coordinate that does not exist: overflow, which caused two of pgctl's three
capture bugs and both of democtl's, stops being a class of bug.

## 19. Every mouse action needs a keyboard path, and a guard enforces it

`guard.Reachable`. Keyboard-only users and ssh are the usual reasons and the weak
ones. Two stronger ones.

**An agent cannot click.** In a framework whose thesis is that an agent gets on
with a tool immediately, an action reachable only by right-click is an action
agents cannot take.

**A right-click may never arrive.** Running the prototype inside herdr showed
herdr's context menu instead of the tool's. A multiplexer that captures the mouse
gets the event first and the application never learns it happened — there is no
protocol for asking. So this is not a preference about input styles: a
context-menu-only action is broken for a whole class of users, and the tool
cannot detect or report it.

The corollary: **the context menu opens from the keyboard too, at the cursor.**
Not "every action in the menu also has a binding" — the menu itself, the way a
menu key works. It is the only arrangement where the two paths cannot drift,
because they are the same list rather than a list plus a keymap maintained
beside it.

## 20. An owner ID is an identity, not a screen position

Found by scrolling the prototype. A list's owner IDs have to carry the absolute
index of the item, not the row it happens to occupy — otherwise a click after
scrolling selects whatever used to be there.

It matters beyond lists because it is the first place the canvas design could go
quietly wrong: "the ID is where it is on screen" is the easy mistake, it works
perfectly until something scrolls, and it produces a wrong action rather than a
visible fault.

The related rule, from the same session: a viewport offset is a separate field
from the selection. Scrolling is looking around; it does not choose. Conflating
them made the wheel appear to pick items at random.

And its corollary, which is the part that is a judgement rather than a fix: a
selection scrolled out of view must say which way it went. The temptation is to
drag the cursor back into the viewport, which is the same conflation again. The
answer is an edge marker — the selection stays where it is, and the interface
admits it is off screen.

## 21. Regions have names, and the harness drives them by name

The prototype's tests already click by owner ID rather than coordinate. That is
the shape a capture script takes — `click list.row[2]; drag split.main +10` —
and an agent naming a region beats an agent guessing a coordinate. It only works
because ownership is recorded at draw time, which is another thing the canvas
gets for free and a string cannot.

## 22. The engine gets nothing from tuikit

The split every tool in the family keeps: the engine knows the domain, the UI
never calls the domain directly, and the CLI is a peer of the TUI over the same
engine. The open question was whether tuikit should offer the engine anything.

It should not, and the bar is higher than "no Bubble Tea imports". Measured on
the tools that exist: democtl's `fleet` imports `fmt`, `math/rand` and `time` —
stdlib only. swarmctl's `internal/engine` imports nothing from charmbracelet at
all. What an engine imports is stdlib plus its own domain SDK: pgx, the Azure
SDK, the Docker SDK.

So the rule is **no terminal concepts**, not merely no terminal library. No
colour, no width, no keys, no framework. An engine that has never heard of a
column is one that can be tested, reused from a CLI, driven from a cron job, and
read by someone who does not know what tuikit is.

**Where the cost lands.** `spec.Command.Run` cannot be the engine's own function
signature, because matching it would make the engine import `spec`. A thin
adapter in the tool — in `cmd/` or beside the CLI — translates between a
declaration and the engine's plain functions. That is a few lines of glue per
command, paid deliberately, and it is the price of the split being load-bearing
rather than decorative.

The alternative considered and rejected: a "minimal UI-free surface" the engine
may import — a context, progress reporting, error kinds. It removes the glue,
but "no UI imports" quietly becomes "no imports we currently consider UI", and
the boundary stops being checkable.

Checkable is the point, and `guard.Engine` does it — a deny-list by *category*
rather than by library, because an engine that measures display width has
learned about columns whichever package it used. swarmctl's, pgctl's and azctl's
engines all pass it; swarmctl's TUI reports 37 imports, so it is not passing
vacuously.

## 23. `tuikit watch` polls, rather than depending on fsnotify

The obvious way to watch a directory is `fsnotify`. This module has four direct
dependencies and they are all Charm, and for a dev tool watching one package the
difference between an inotify callback and a digest every 300ms is imperceptible
— the rebuild itself takes 200-600ms. A dependency that buys nothing a person
can feel is a dependency that only costs.

Digesting the tree rather than comparing a single mtime also makes a save that
rewrites a file with identical content a non-event, which is what an editor with
format-on-save does constantly.

Two rules the loop follows that are not about watching:

**A failed build puts the error on the page.** Leaving the last good frames up
when the code no longer compiles is a page that describes a tool which does not
exist — the same failure as a design system maintained beside the code rather
than generated from it. The error page carries the compiler's own output,
file and line intact.

**The rebuild is debounced.** An editor that writes a file in two syscalls is one
save, and rebuilding twice makes the page flicker through a state nobody asked
for.

## 24. A palette role holds a `TerminalColor`, not a `Color`

Found by migrating azctl, which is what that migration is for.

azctl's palette is `lipgloss.AdaptiveColor` pairs — a light value and a dark
one, so the tool reads on a pale terminal as well as a dark one. `theme.Palette`
held `lipgloss.Color`, which cannot express that at all, so azctl's choice was
between tuikit's vocabulary and working in daylight. That is not a choice a
design system should be imposing.

**Depth and adaptation are different questions.** Decision 13 chose ANSI 256
over truecolour because ssh decides the profile — that is about how many colours
there are. Which of them to use against a pale background is a separate
decision, and the tool's.

So the nine roles hold `lipgloss.TerminalColor`. `theme.Default` still names
ANSI indices and nothing about it changed; a tool that wants adaptive pairs can
now have them.

Two things follow:

**`theme.Hex` answers with the DARK value of an adaptive colour**, and passes a
hex value straight through. Everything that renders a palette outside a terminal
— a design system page, a captured frame turned into HTML — draws on a dark
ground, because that is what the frame was captured for. A palette page for a
light interface is a real thing to want, and is not this.

**`theme.Value` is new**: the colour as the tool *declared* it, for the design
system page to show beside what it resolves to. "205" tells a reader the palette
is ANSI indices and will follow their terminal's own scheme; "#d2a8ff" tells
them it will not. Hex answers what it looks like; Value answers what it is.

## 25. The nine roles are the framework's vocabulary, not the tool's promise

Two guard bugs, both found in the first ten minutes of azctl's migration.

**`guard.Tokens` demanded that every role be drawn with.** azctl failed on four
at once — `Stderr`, `SelectionFG`, `SelectionBG` and its own `Text` — for the
crime of being a resource browser rather than a table. "Either use it or drop
it" is advice a tool cannot take: the nine are fields on a struct it inherited
from `theme.Default`, and there is no dropping them.

The check came from swarmctl, where the palette was the tool's own file and an
unused role really was dead code. Inheriting a palette is what made that stop
being true, and nothing noticed until a second tool inherited one.

So the check now applies to `Extra` only. An unused Extra is still dead and
still fails: a tool that invented a tenth meaning and then did not use it has
left a name for the next person to wonder about.

**`guard.Glyphs` did not honour `GlyphSet.Printable`.** It looked the rune up in
the map directly, so `SpinnerRange` — documented in `theme` as the one thing the
allow-list does not cover — did nothing, and azctl's `⠿` was rejected. Two
functions answering "may this be printed" differently is worse than either
answer on its own; the guard now asks `Printable`, which is the exported one.

Worth noting what the second bug was hiding behind: `SpinnerRange`'s comment
says the spinner's characters "never appear in a string literal and the guard
never sees them". That was true when the spinner came from bubbles. It stopped
being true when `comp.Spinner` started drawing braille itself, and a tool
writing its own frame was never covered at all.

## 26. A depth and a toggle set, not a `comp.Tree`

azctl has a tree, and the question was whether `comp` needs one. The five tools
were surveyed rather than guessed at:

| | Hierarchy in the UI? |
|---|---|
| azctl | **Yes** — real expand/collapse, `expanded map[string]bool`, ▸/▾, two levels |
| swarmctl | The flattened shape in **four** places — diff rows, apply rows, pane rows, disk detail — and always fully expanded |
| dugo | Recursive *data*, presented as a **drill-down** with a nav stack |
| pgctl | No. Five flat parallel panels; hierarchy is panel-to-panel scoping |
| democtl | No. `stack` is a detail field, never a grouping |

A tree widget has one consumer. That is not the bar.

**dugo is why the answer is not "one consumer for now".** It has the deepest
hierarchy of the five and answers it with `navStack []NavigationState` — push on
descend, pop on ascend, the cursor remembered per level. Shipping `comp.Tree`
would commit the framework to the two-level flattened model, and the tool with
the most hierarchy in it would not use the component. A drill-down is not a
worse tree; it is the right answer when the depth is unbounded and the fan-out
is large, which is what a filesystem is.

### What did have two consumers

**A row that knows its depth.** azctl left the clearest possible statement of a
missing field, as a comment explaining a workaround:

> The marker is built into the row rather than set on the list, because azctl
> marks a RESOURCE row and not a bucket header — comp.List's Marker is one
> character for the whole list, which is right for a flat list and not for a
> tree.

So it built the indent into the row text and recomputed `i == cursor` itself, to
place its own cursor glyph — doing the list's job, inside the list's input.
swarmctl does the same by hand in four places, as a literal `"  " + line`.
`Row.Depth` and `Row.Lead` are that comment, as two fields.

**A keyed toggle set.** azctl's `expanded map[string]bool` and swarmctl's
`revealState{all bool; rows map[string]bool}` are the same type under two names,
for two unrelated purposes — expansion, and unmasking secrets. A shape arrived
at twice independently, for different reasons, is the strongest form of this
project's signal. `app.Toggles` took swarmctl's semantics, which were the better
of the two: turning the global override off also clears the per-key set, so the
key twice is a reliable way back to nothing.

### What must not be generalised

The flattening. `row{bucket, res}`, `diffRow{service, action, change}` and
`applyRow{service, stack, edits, edit}` look alike and are not: each carries a
domain payload, and a shared `[]Node` would make every tool box its data or keep
it twice. It is 12–45 lines of domain code per site, and it is the part that is
genuinely different each time. A list of rows with a depth on them is the widest
interface that is still honest.

### Two rules the survey turned up, worth keeping

azctl's `rebuild()` does two things nobody else does and everybody eventually
wants:

- **The selection is preserved by identity, not by row index.** After a
  re-pivot, a filter or a refresh, it re-selects the row whose resource ID
  matches the one that was selected. This is the same rule as owner IDs carrying
  the absolute index (decision 20): a position is a fact about the screen and an
  ID is a fact about the thing, and only one of them survives a rebuild.
- **A filter implies expansion.** You searched for a thing; you want to see the
  thing, not a list of folders it might be in.

Both need the tool's own identity concept, so they are rules here rather than
code in `comp`.

## 27. tuikit is a shape for one kind of tool, not a TUI framework

Read the five most-starred Bubble Tea app frameworks in full — see
[research/framework-comparison.md](research/framework-comparison.md) — to answer
the only question that matters about a project like this: what is it offering
that the alternatives are not?

The field splits in two, and tuikit is in neither half.

**bento is a rendering substrate.** A faithful `ratatui-core` port — grapheme
cells, `Buffer.Diff`, a Cassowary solver — with no opinions about applications
at all. It does not even depend on Bubble Tea.

**The other four are app frameworks over strings.** bubbleapp is React,
bubblyui is Vue, the two sodas are a screen stack. Every one composes with
`lipgloss.JoinVertical`.

tuikit is an opinionated shape for **one kind of tool**: an operator CLI *and*
TUI over a domain engine, with `guard.Engine` failing the build if the engine
learns what a terminal is. That is a narrower claim than "a TUI framework", and
the narrowness is the product.

### What follows from that

**We do not compete on component count.** bubblyui has 31 and they are real;
ours are 14 and thicker, because each was extracted from two of four working
tools rather than designed. A component that exists in one tool is not a
component, it is that tool's code.

**We do not need reactivity.** Refs, computed values and watchers answer "what
changed" for a graph of UI state. Our shape is a snapshot redrawn — `app.Gen`
drops the result of a read the reader walked away from, and that is the whole
async problem an operator tool has.

**The CLI is not a side feature.** It is the axis nobody else is on. Not one of
the five has anything to say about flags, subcommands or exit codes, because
none of them is built for a tool that must also work in a pipe.

**The cell buffer is a means, not the thesis.** bento has cells and stores no
owner, with mouse reporting never switched on. bubbleapp wanted ownership badly
enough to embed invisible ANSI markers and scan them back out of the rendered
frame, with an overlap heuristic its own source calls guesswork. Each has half
of what a cell that records its owner gives you for nothing. If the substrate
were the point, bento would already have won; it is the layer above that was
missing.

**Tests of appearance are the differentiator nobody replicated.** bubblyui has
4,816 tests, a `SnapshotManager` with diffing and normalizers, and a written
guide — and zero committed snapshot files. Not one component's appearance is
pinned. Building the machinery is easy; committing the frames is the part that
requires believing they matter.

### What this decision does not license

Ignoring the field. Three things came back worth taking, and are filed: deleting
`View() string` (#21), a screen stack with history (#22), and constraint layout
(#23) — the one axis where somebody else is plainly better than us.

## 28. The palette is the terminal's own sixteen

Decision 3 chose ANSI indices over truecolour hex and was right about why. It
was wrong about the scope, and the wrongness was invisible for as long as nobody
asked the next question.

**Only indices 0–15 follow the reader's theme.** Everything from 16 up is a
formula — a 6×6×6 cube and a 24-step grey ramp — identical in every terminal,
which no theme touches and no theme can. `theme/hex.go` has said so in a comment
since it was written: *"the first sixteen are the terminal's own, and predate any
formula."* Nobody joined it to the palette, which was 205/241/240 — three numbers
in the fixed region.

So a tuikit tool looked the same under gruvbox, tokyo-night and solarized. That
is not a neutral choice. It means the interface overrides what the reader
already decided, everywhere, forever.

`theme.Value`'s doc comment had the error written down: *"205 tells a reader the
palette is ANSI indices and will follow their terminal's own scheme."* It does
not. **The line is not index versus hex. It is fifteen.**

### The mapping

The sixteen are already semantic, which is what makes this a mapping rather than
a guess. Terminal themes agree that 0 is the background, 7 the foreground, 8 the
dimmed grey that comments are drawn in, and 15 the brightest text. Omarchy's
templates say it literally — `palette = 0={{ background }}`, `palette = 8={{
muted }}` — and every other theme system does the same under other names.

| Role | Was | Now | |
|---|---|---|---|
| Accent | 205 | **13** | bright magenta |
| Muted | 241 | **8** | the dimmed grey |
| Border | 240 | **8** | the same grey |
| Success | 34 | **2** | green |
| Pending | 214 | **3** | yellow |
| Danger | 196 | **1** | red |
| Stderr | 203 | **9** | bright red |
| SelectionFG | 229 | **0** | the background |
| SelectionBG | 57 | **7** | the foreground |

### The selection is reverse video, and that is the point

SelectionFG is the background and SelectionBG the foreground, so it inverts
correctly on a light theme **by construction** rather than by detecting one. No
`AdaptiveColor`, no `COLORFGBG` sniffing, no light profile to maintain.

swarmctl's design notes reached the same place independently, exploring a light
terminal profile: reverse video is *"the one treatment that reads identically in
both profiles"*. Two routes to one answer is the signal this project runs on.

### What it costs

**One grey.** Muted and Border are the same index. They were 241 and 240 — one
step apart on the ramp, `#626262` and `#585858`, which nobody could tell apart.

**The accent is the terminal's magenta, not the theme's own accent colour.**
ANSI has no accent slot. A tool that wants the real one reads it from wherever
its desktop keeps it and assigns the role; that is what a plain assignment is
for, and `guard.Tokens` still holds every *other* colour closed.

### What it does not cost

Nothing else. The guards are untouched, because they check that a colour came
from the palette and not what the palette contains. **Not one golden moved**,
because goldens are colour-stripped. And there is no reload to write: an index
is resolved by the terminal at paint time, so changing the theme retints the
next frame — where a tool that baked hex into a config file needs a signal
handler and a re-read.

### Two things fixed on the way

`Hex` answered `#000000` both for index 0 and for a colour it could not read. It
answers `""` for the failure now: index 0 is a real role, and a sentinel that
collides with a legitimate answer is a check that has stopped checking — the
test that every role converts would have passed over a broken one.

And `base16` claimed to hold xterm's defaults while holding the VGA set
(`#800000` red, `#c0c0c0` white). It matters more than it did, because those
sixteen are now every swatch on the design system page.

## 29. Pixels are a decoration pass, and the cells are always the drawing

swarmctl's design notes asked for a pixel layer — a panel drawn as a real image
composited over the terminal, with a soft shadow and a smooth gradient. The
obvious reading is that this needs a second renderer and a second view, and that
a tool would have to be written twice.

It does not, and the reason is decision 20. **An owner ID already says which
cells a component drew.** So the pixel layer needs no layout of its own; it
needs one call:

```go
// Picture asks that, if this terminal can, img be painted over the cells
// owned by id. A terminal that cannot simply shows the cells.
func (c *Canvas) Picture(id ID, img image.Image)
```

A component draws its panel as characters — always, unconditionally — and then
says "and if you can do better, here is a picture of that." Nothing is written
twice, because the cell drawing is not a fallback that exists for the pixel
path's benefit. It is the drawing. The picture is a decoration over it.

**The goldens never move**, because capability detection is false under test.
Every golden is permanently a picture of the cells, which makes the fallback the
most heavily regression-tested path in the system rather than the least.

### Omarchy 4 made this mandatory rather than tidy

There are two terminal graphics protocols and no terminal speaks both.

**Sixel** is DEC's, from the 1980s: the image is encoded as text in stripes six
pixels tall and pasted at the cursor. It is opaque — it covers the cells under
it and nothing shows through — colour-register based rather than truecolour, and
it has no z-index. Disturb the screen and it is gone.

**The kitty graphics protocol** carries 32-bit RGBA and a z-index, stores the
image under an id so it can be re-placed cheaply, and can sit *behind* the text
at z = −1. Its Unicode placeholders let an image be anchored to ordinary
printable cells, which means it survives a line-diffing renderer untouched —
bubbletea v1.3.10's included.

(The kitty **keyboard** protocol is an unrelated thing with a colliding name.
foot implements that one, which is why "foot supports kitty" gets said.)

Omarchy 4 "Quattro" supports four terminals, and they land in three places:

| | Speaks | Gets |
|---|---|---|
| **foot** — the Quattro default | Sixel | pixels, opaque |
| Ghostty | kitty protocol | pixels, layered |
| kitty | kitty protocol | pixels, layered |
| Alacritty | nothing | cells |

foot is the default *and* is the one that will never gain the other protocol:
its maintainer closed the question in foot#481 — *"I'm fairly sure I don't want
another image protocol in foot ... our somewhat fuzzy goal of being
'lightweight'."* And Alacritty's own entry in the Omarchy manual reads *"does
not, however, support native tabs, splits, or image rendering."*

So the cell drawing is not the degraded case for a minority. It is the only
thing an Alacritty user ever sees, and it is what every foot user sees for the
frames between a redraw and the next transmission.

### Design to 5a, and let kitty be a quiet upgrade

The best-looking mock in the handoff is 5b: rounded on all four corners, a
translucent scrim dimming the app while the rows stay real text on top, a smooth
gradient. All three of those are the z-index and the alpha channel. **They are
not available on the default terminal.**

The default gets 5a's four concessions — opaque panel, shadow baked into the
image against a known background, origin snapped to a cell boundary, stepped
gradient. Designing to 5b first would make foot look like a broken version of
the real interface instead of a plainer one. Design to 5a; Ghostty and kitty
users get softer edges and nobody is shown a hole.

### Three emitters, one design

The expensive part is shared. The panel is rasterised once into an
`image.RGBA`; only the final encode differs — Sixel bytes, kitty bytes, or
nothing. Ordering follows the audience, not the polish: cells, then Sixel for
the default, then kitty.

### The canvas owns the bytes, not the runner

First reading of this said the runner had to own the final bytes, and that
issue 21 (delete `View() string`) was therefore a prerequisite. That was wrong,
and worth writing down because the mistake is the instructive part.

Every tuikit model already ends `return c.String()`. **The canvas is already the
last thing to touch the frame**, so the decoration pass belongs there — and for
the kitty protocol that is not merely adequate but correct, because placement
bytes are part of the frame and want to travel with it.

Issue 21 remains worth doing for the reason it was filed — a `string` seam lets
a tool hand back a hand-joined frame and no guard will catch it. But it is a
guard question, not a pixel one, and the ladder does not wait for it.

### What it costs

**A raster path and an embedded typeface**, which is a new asset class and a
real number on the binary. This is the actual work; the escape sequences are
the small part.

**A capability query with a timeout** — the one place tuikit reads from the
terminal rather than writing to it.

### The gradient is read from the terminal, not chosen here

Issue 30 asked this to be settled explicitly, because it is where decision 28
and a pixel layer pull against each other. Decision 28 put the characters on
ANSI 0–15 so the reader's terminal theme wins. A picture carrying a literal
`#5f00ff → #ff5faf` would walk straight back out of that: identical under all 22
Omarchy themes while the text beside it changed, which is a worse result than
having no picture at all.

**So a component asks for a picture and never for a colour.** `Meter` has a
`Pixels bool`, not a ramp; the gradient lives on the canvas, and the canvas read
it from the terminal with OSC 4 — indices 5 and 13, Accent's own family. Change
theme and the picture changes with the text. This is the same arrangement the
nine roles give the characters, and it is what keeps `guard.Tokens` meaningful:
a component that could pass its own gradient is a component that can escape the
theme.

Two more things are read back the same way: the background, with OSC 11, because
Sixel has no alpha and a soft edge must be composited against a colour that is
KNOWN rather than guessed; and the cell size, with `CSI 16 t`, because a picture
is asked for in cells and drawn in pixels and nothing else knows the ratio.

### And there is no embedded typeface

The obvious next thing is a heading rendered in a real font, and it is
deliberately absent. Sixel is opaque, so any text inside a picture's rectangle
must be IN the picture — which means an embedded face, a licence decision, a
real number on the binary, and the same words rendering differently depending on
which terminal the reader has.

Keeping text as cells avoids all of it. The picture is a band; the words live
above and below it in the reader's own font at their own size, legible at any
width, selectable, and identical on all four terminals. The pixels do what
pixels are good at — gradients, curves, soft edges, resolution — and nothing
else.

### Nothing is blanked, and that took a wrong turn first

The Sixel path originally blanked the cells under a picture. An opaque image
covers them anyway, so leaving them could only make them flash back between
redraws — which is correct about the good case and badly wrong about the bad
one.

Forcing Sixel on Ghostty, which does not speak it, showed what that costs: an
EMPTY `[     ]` where the character bar would have been. The same happens if a
multiplexer strips the sequence, or a terminal advertises Sixel and then
declines the payload. Worse than having no pixel layer at all.

So the characters are drawn and the image covers them. It costs a frame of text
in a rare in-between case, and it degrades to the cell drawing in every failure —
which is the thesis of this decision. Blanking quietly made the cells a fallback
that only existed if the decoration worked.

### What it does not cost

The guards, the goldens, the harness and the mouse hit-testing are all
untouched, because a decoration pass changes no cell and no owner ID. A frame
with a picture over it is the same frame.

### What shipped

`term` (detection, cell size, background, palette, both encoders), `paint`
(ramps, rounded panels, bars, flattening), `Canvas.Picture`, `comp.Meter`, a
gallery entry, a meter in democtl's run screen, and `tuikit pixels` — a
self-test that prints what your terminal answered and draws the same bar with
and without a picture, because no test in CI can tell you whether foot draws it.

## 30. The model draws into a canvas; there is no `View() string`

TEA's third step is `View() string`, and that string is a seam.

A tuikit model is supposed to compute no coordinates: components take rects,
the canvas owns the grid, and the tool declares what goes where. azctl was got
to the point where `grep -c 'c.Text(|c.Fill(|c.Set(' internal/tui/*.go` is
**zero** — a real result, and one nothing enforced. A model returning a string
can hand back a hand-joined pile of lipgloss, or a canvas frame with something
concatenated onto it, and the goldens will record whatever comes out. It was a
habit, and a habit is not a property.

```go
type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) (Model, tea.Cmd)
	Draw(c *comp.Canvas, r comp.Rect)
}
```

Handed a canvas and a rect, a model has nowhere else to put anything. Update
returns a `Model` rather than a `tea.Model` for the same reason: a signature
accepting `tea.Model` accepts anything with a View, which is the escape hatch
this exists to close.

### What the runner took over

Four things every tool set up the same way and could get subtly different:

- **The canvas**, and the one place `comp.NewCanvas` is called in a program.
- **Its size.** One row shorter than the terminal — which both tools here
  already did, independently, with different arithmetic (democtl through its
  band layout, the gallery through a `bodyHeight` helper). Two of two is a
  convention, not a preference.
- **The chrome**, so one pane cannot end up with different box characters
  because somebody threaded it through five call sites and missed one.
- **The pixel layer**, and with it the rule that detection is asked by main and
  never by a model — see decision 29.

### The one-row margin is inherited, not decided — see issue 37

The runner sizes the canvas at `terminal height - 1`, justified above as a
convention because both tools here did it independently. That justification is
weaker than it reads: democtl's own comment says the row came from "where the
old bodyHeight arithmetic already put it", so it was inherited from swarmctl
rather than chosen, and two tools agreeing is not two decisions when one was
copied from the other.

The plausible reasons are real — pending wrap on the bottom-right cell, and
where the cursor parks — and neither has been measured. The counter-evidence is
also real: azctl ran at full height until it migrated here, with no scrolling
and no lost top line.

So this is recorded as an open question rather than a settled rule. It costs
every tool a row on faith, and issue 37 says how to find out.

### It cost the harness nothing

`*app.Runner` has `View`, `Update` returning a `tea.Model`, and `Canvas`, so it
satisfies `harness.Model`, `Driver` and `Pointer` unchanged. Tests wrap the
model and drive the runner; assertions still read the model's fields. Nothing in
`harness/` was touched.

**Not one golden moved**, in either tool, which is the check that matters: the
runner reproduces frames that four hundred lines of layout used to produce by
hand.

### `SetSize` goes through Update

The runner could set its own fields and be done. It sends a `tea.WindowSizeMsg`
instead, so a model that keeps its own idea of the size — most do, for layout —
hears about it by the same route it would in a running program. Two paths to one
fact is how they come to disagree, and the disagreement would only show up in
captured frames, which is the worst place to find it.

## 31. "No use case" and "we don't know the shape" are different reasons

The extraction rule — generalise only where two of the four tools differ
meaningfully — has been doing more work than it should, because it was being
used to explain three different situations that want three different answers.

**A component we know how to build and nobody has asked for.** No design risk.
Deferring is fine, but it is not a decision, and it quietly becomes "we never
got round to it". `comp.Input` sat here for weeks: a one-line text field is not
a mystery, it was waiting for a second consumer to make it feel earned. The rule
did no work there at all.

**A component two tools need that has not been built.** A backlog item. The rule
says BUILD, so this case should never persist — if it does, the rule is being
misquoted rather than applied.

**A component whose shape we do not know.** Here the rule earns its keep, and
decision 26 is the proof: azctl's `row{bucket, res}` and swarmctl's
`diffRow{service, action, change}` look alike and are not, because each carries
a domain payload. A shared `[]Node` would have forced both to box their data or
keep it twice. `comp.Tree` would have been WRONG, not merely early.

**So: no use case is a reason to wait. Not knowing the shape is a reason to
refuse.** Only the second is a principle.

### The tell

Three of the five things filed after this conversation had their data structures
already built, public, tested and unused:

- `app.Stack.Path()` — a method for a breadcrumb nothing draws. Nothing in the
  repo calls it.
- `List.Offset()` and `List.Max()` — everything a scrollbar needs, exposed, with
  no scrollbar.
- `spec.Command.Key` — "the keyboard path to this command inside the TUI", and
  no help screen reads it.

A public, tested, uncalled accessor is the clearest evidence available that
something is in the first case rather than the third. The shape was worked out
when the accessor was written; only the drawing is missing.

### What stays refused

A general `Dialog`, form validation, a layout DSL beyond `comp.Layout`, and
`Tree` again. Not because nobody wants them — because two tools wanting one is
not the same as two tools wanting the SAME one, and that difference is only
visible from a migration.


## 32. The comments are load-bearing, and that cuts both ways

Three bugs in one evening had the same shape: **prose describing behaviour the
code never had.**

- azctl's mouse handler: *"a bucket header expands with enter or a click on
  it"*, above a handler that only moved the cursor. Clicking a header did
  nothing visible, so the mouse looked broken.
- azctl's runner footer: `q abort (the running step finishes)`. `q` cancelled
  the run and quit the whole program, with no question asked.
- `comp.Spinner`'s own doc: *"azctl's single ⠿ is a spinner that has stopped,
  which reads as hung rather than as working. Not carried."* It was carried, by
  omission — the runner never ticked, so the spinner froze on whichever glyph
  the last event left it on.

This project writes down WHY more than most, and that is worth keeping: every
decision here exists because a comment recorded a reason somebody would
otherwise have had to rediscover. But a comment is a claim nothing checks, and
three of them were claims about behaviour rather than about reasons.

**A comment explaining why is documentation. A comment describing what the code
does is an untested assertion.** The first kind ages into insight; the second
ages into a lie, and a confident one, because it was written by someone who
meant it.

Two of the three were caught by a person running the tool, not by the 105 files
under `testdata/`. The goldens hold the SHAPE of a frame, and none of these
changed a frame — a click that does nothing draws the same cells as a click that
was never wired.

Issue 39 asks whether the footer case is mechanically checkable, since a footer
is a list of `comp.Hint{Key, Label}` and a handler is a switch on
`msg.String()`. The other two may only be a discipline: when a comment says what
happens, there should be a test named after the sentence.

## 33. Config goes in `~/.config/<tool>/`, and it is copied rather than imported

Three tools, three answers, and on the machine this was written on you can see
the disagreement on disk.

| tool | how it resolves | where that lands on macOS |
| --- | --- | --- |
| pgctl | `os.UserConfigDir()` | `~/Library/Application Support/pgctl/config.yaml` |
| swarmctl | `os.UserConfigDir()` | `~/Library/Application Support/swarmctl/config.yaml` |
| azctl | `os.UserHomeDir()` + a hardcoded `.config` | `~/.config/azctl/playbooks` |

`os.UserConfigDir()` honours `$XDG_CONFIG_HOME` on Linux and ignores it on
darwin, where it returns `~/Library/Application Support`. So pgctl and swarmctl
have byte-for-byte identical search code that can never land where azctl's does.

The tell that the stdlib answer is the wrong one is what was actually found
there: swarmctl's own `state.yaml` sitting among `com.apple.ContextStoreAgent`
and `com.apple.avfoundation`, while the configs a person had written by hand —
`orb.yaml`, `orb-monorepo.yaml` — were in `~/.config/swarmctl/`, where the
search order does not look. The user put them where the convention says they go.
That convention is not a preference: the same `~/.config` held `gh`, `git`,
`nvim`, `fish`, `tmux`, `btop`, `sops` and `gcloud`.

**So: `$XDG_CONFIG_HOME` if set, else `~/.config/<tool>/`, on every platform
including macOS.** `os.UserConfigDir()` is right for an application with a
bundle identifier and wrong for a developer's command-line tool. The deciding
property is that `~/.config` is *syncable* — it is the directory people symlink
into a dotfiles repository, and a config you cannot carry to the next machine is
a config you will write twice.

**State is not config.** Config is written by a person and belongs in that
dotfiles repository; state is written by the tool and must not follow you to
another machine — swarmctl's `state.yaml` maps an environment to the SSH key
installed *on this laptop*. State goes in `$XDG_STATE_HOME`, else
`~/.local/state/<tool>/`. One tool has state today, which is exactly when the
pattern is cheap to set.

### Why this is not a `tuikit/conf` package

Two tools have hand-rolled the same fifteen lines, which is the trigger in
decision 31 for extracting a component. It is still refused, and the reason is
decision 22: `guard.TerminalPackages` denies
`github.com/richarddavenport/tuikit` outright — *the engine gets nothing from
the framework*. Config loading is engine work in both tools that do it
(`pgctl/internal/config`, `swarmctl/internal/engine`), so a shared package would
force one of two things, and both are worse than the duplication:

- weaken the deny-list to "tuikit, except the parts we decided are not really
  the framework", which is the "minimal UI-free surface" already rejected in
  decision 22 — the boundary stops being checkable the moment it has an
  exception; or
- move config loading into the UI, where it is not domain work and the CLI
  cannot reach it.

**The extraction rule answers "is this shared?", not "should it be a package."**
Here the answer is shared and copied. The mechanism is the scaffolder: `tuikit
new` writes `internal/engine/paths.go` into the tool, so a new tool starts
correct without a runtime dependency, and the code it starts with is fifteen
lines it owns and can change.

That is a real cost, honestly stated: fixing a bug in those fifteen lines means
fixing it in every tool. It is accepted because the alternative is a layering
rule that no longer means anything, and because the thing being copied is a
policy that should almost never change — if it does, it is because an operating
system moved, and that is not a patch anybody applies silently.

Existing tools were deliberately **not** changed. Moving pgctl's and swarmctl's
search order orphans a file that exists right now, and each tool's own repository
is where that migration gets weighed.

## 34. Frames publish as SVG, and the typeface question stays open

`docgen` produced one thing: a self-contained HTML page. It is the right shape
for the inner loop and for sending someone a link, and the wrong shape for
everything that lives in a repository — a README, an mkdocs site, a pull request
— where a page cannot go and a fenced block of raw ANSI renders as noise. So
`tuikit frames -md` writes Markdown with an image per frame.

"An image" is where the decision is. A PNG needs a rasteriser, a rasteriser
needs a typeface, and **decision 30 — embed a typeface or decide never to — is
still open.** Rasterising for documentation would answer it by accident, in the
one context where the answer is least considered: the fallback everywhere else
is that text stays cells, and a document is not a good reason to reverse that.

SVG needs no typeface of its own. It names the same system monospace stack the
HTML page already names and lets the reader's machine draw, which is the
existing rule rather than a new one. It also keeps the text as text, so a frame
in a document stays greppable and readable by a screen reader — a PNG of a
terminal is opaque to both.

Two properties a naive SVG writer would not have, both found by looking at real
frames rather than by reasoning:

- **Every span states the width it must occupy** — `textLength`, with
  `lengthAdjust="spacing"`. A character grid reproduced by trusting the font's
  advance width is one box-drawing character away from a frame that does not
  meet, which is the same failure the page avoids by measuring nothing in
  pixels. `spacingAndGlyphs` reaches the width by distorting the glyphs, so the
  borders bend rather than move; `spacing` adjusts the gaps and leaves them
  alone.
- **No `<style>` element.** An SVG referenced from a Markdown document is
  rendered through a sanitiser. A stripped stylesheet leaves a frame that is all
  one colour with no error anybody sees. Presentation attributes survive; a
  stylesheet is a bet, and the losing case is silent.

**What this forced elsewhere, and the better outcome.** Two renderers need the
same answer to "what colour is this character", and the SGR reader lived inside
the HTML writer's line loop. Rather than copy it, it came out as
`harness.Rows(frame) [][]Span` — the ANSI reader every format shares. Two
parsers would eventually disagree about a frame neither of them drew, and the
existing HTML tests (accumulated style, style crossing a line break, a colour
channel that looks like a code) now cover both formats because both go through
the one reader. The same went for the section ordering, which is a rule — a
frame no group named must still appear — and not a layout detail.

Verified by rasterising democtl's frames and looking at them: the dashboard, a
modal drawn over two panes, and a `comp.Meter` — the longest run of box-drawing
in the suite, and the case that would show a broken grid first.

## 35. Chrome characters come from the chrome, and `comp.Rule` is the missing one

Filed from a tool (issue 43): azctl drew a horizontal rule under two headers
with `c.Fill(bands[1], "─", …)`, which broke the invariant azctl's own AGENTS.md
documents about itself — `grep -c 'c.Text(\|c.Fill(\|c.Set(' internal/tui/*.go`
is supposed to be 0.

The report claimed two tools and three instances, and was careful to say it was
discounting four of swarmctl's five `─` hits as `comp.Pane`'s job in a tool that
predates `comp.Pane`. That care is what made it checkable, and checking it found
more: searching for the SHAPE rather than the character turns up seven
instances across four codebases, and **three of them were inside tuikit** —
`gallery/model.go`, `examples/democtl/ui/view.go`, and `comp/palette.go` twice.

Every one is the same three lines: a `Bar` drawn into one band, then a rule
filled into the next.

**The bug this was hiding.** Five of the seven wrote the literal `"─"`, and two
asked `c.Chrome().Box.Top`. `theme.ASCIIBox` is `{"+", "-", "+", …}` and exists
so an interface can be drawn in a font that has nothing. A tool on it got `-`
from its palette's rules and `─` from its header rule **in the same frame** —
precisely the failure ASCIIBox exists to prevent, shipping in tuikit's own
gallery and in democtl. `guard.Glyphs` never had a chance: `─` is a legal glyph.
It was the wrong SOURCE, not an illegal character, and the only fix for a wrong
source is to have one source.

So the rule is: **a character that belongs to the chrome is read from the
chrome, never written as a literal.** `comp.Rule` takes `Chrome.Box.Top` unless
told otherwise, so a tool that changes its box set changes its rules with it.

### The two shapes that were rejected

`comp.Bar{Fill: '─'}` — a field on the component that was already closest. Bar
means "one line with content at each end", and a Bar with no content that fills
its own width is a second meaning wearing the first one's name.

`Bar` underlining itself, which the report preferred and which I wanted to be
right because it is one component fewer. Six of the seven sites are a Bar with a
rule beneath it. The seventh, `comp/palette.go`, is a rule **above** a footer.
Which side of a band the line falls on belongs to the layout, not to whichever
bar happens to be adjacent — and a component that covers six of seven leaves the
seventh hand-rolled, which is how a component ends up half-adopted.

### What it cost, and what it did not

Not one frame changed. The goldens moved only where the gallery gained a row and
a count, which is the evidence that the seven sites were drawing the same thing:
if any had differed, a golden would have said so.

The guard that would have caught this is NOT built, and issue 48 records why
rather than leaving it implied. A literal `─` in a tool's source is not always a
draw: both `gallery.Entry.Glyphs` and a tool's own glyph-set declaration list
chrome characters in order to ALLOW them. Telling a declaration from a draw is
the same distinction `guard.Glyphs` already makes by parsing rather than
matching, so it is probably tractable — but it is a design question with two
known false-positive classes, and decision 31 says not knowing the shape is a
reason to refuse rather than to guess.

## 36. A default that was right for the first tool is evidence, not a law

Three reports against `comp.List` landed together (issues 44, 45, 46), from
pgctl migrating onto the canvas. They look unrelated and are the same thing: a
choice made when one tool used the component, meeting the second tool.

**The selection swallowed a row's state glyph (44).** A selected row is drawn in
one colour whatever its spans say, because "a row that kept its own colours
under it would make the cursor hard to find in exactly the list where finding it
matters." That is right for a LABEL. pgctl's connection list marks reachability
with `●` `○` `✗` in the first column, and on the cursor row all three came out
bold black on white — the one row a reader is looking at was the one row whose
status they could not read. **A person using it reported this**, not a test:
*"when highlighting I can't see the color of the dot."*

It is also an accessibility defect and not only a legibility one. A black `●` on
light grey does not read as "a green one, highlighted"; it reads as a DIFFERENT
state — off, disabled. pgctl was saved by using four distinct shapes as well as
four colours. A tool encoding state in colour alone would have lost it outright
and nothing in the API would have warned it.

`Row.LeadStyle` keeps the lead column's own colour through the selection. Only
the lead. Letting every styled span survive is more elegant and makes the
cursor's prominence depend on how colourful a row happens to be — strong on a
plain list, nearly invisible on a busy one, which is the opposite of what a
selection is for.

**The status row is charged per list (46).** It is reserved whether or not the
list overflows, because "a viewport that only looks like one when it is
scrolling is a viewport you cannot tell from a short list." True for azctl's one
big tree, where `18/18` earns its row. pgctl stacks FIVE lists in a column: at
80×24 that is five of about twenty-one body rows, a quarter of the column, on
counters reading `3/3`, `3/3`, `1/1`, `1/1` and blank — beside panel titles that
already say the same number. `NoStatus` turns it off, and `Overhead()` reports
what the list spends on itself so a tool stops encoding `const chrome = 3`.

**`Select` cancelled a pending `Move`, silently (45).** Every one of these tools
independently arrived at "clamp every cursor when the data changes", from when a
cursor was a plain int that could point past a list that had shrunk. Against the
deferred `Move` that clamp reads as `Select(Cursor())`, which zeroed the move the
arrow key recorded in the same `Update`. The move was applied and immediately
discarded: **the key did nothing and nothing errored.**

The fix is narrower than the report's options. Selecting the row the cursor is
already on is not a selection — it is a caller saying the cursor is fine where
it is — so it keeps the pending move. Out of range it still selects for real,
which is the case the clamp was written for. Deliberately NOT the wider "apply
pending on top of any Select": a click means that row, and a queued arrow key
landing on top of a click would be worse than the bug.

### The rule

**A component's default is evidence from one tool until a second tool has used
it.** All three defaults were argued for in a doc comment and all three
arguments were sound — they were just sound about a shape only one tool had.
None of the three was found by a test, and one was found by a person. The
migration is doing what decision 14 said it would; this is what "prove it
against a second tool" looks like when it works.

What this does not license is a flag per preference. Each of these has a default
that stays, a reason the default is right where it was right, and a second shape
that the first reasoning demonstrably does not cover.

## 37. The overrun guard is kept, and it lives in the gallery

Issue 6, open since the first week, blocked on the canvas existing.

`guard.Width` was planned because overflow is the commonest bug in a TUI and the
one assertions miss — two of pgctl's three capture bugs and both of democtl's
were something drawn wider than its container. Then the canvas made it look
unnecessary: `Set` clips, so drawing past the edge is a coordinate that does not
exist rather than an error to catch.

**The canvas guarantees the wrong thing.** It guarantees nothing lands outside
the CANVAS. It guarantees nothing about a component staying inside the RECT it
was handed, and that is the failure that matters: a component which overruns
paints over its neighbour rather than failing. The frame is still well-formed,
every golden still passes, and the pane beside it is simply wrong.

`Canvas.Clip` closes it structurally — a clipped canvas cannot draw outside its
rect — but only for components that call it. Measured: fourteen of twenty-two
did, and the eight that did not were not all bugs, because several take no rect
at all. `Confirm.Draw(c, id)` centres itself and returns where it landed;
`Menu.DrawAt` nudges itself back on screen; `Split.Draw` divides a rect and
returns two. Those position against the canvas by contract.

So it is kept. Not as `guard.Width`, and not in `guard`.

### Where it lives, and why not in guard

**The gallery.** It is already the complete list of components — held closed
against `comp` by `TestEveryComponentIsInTheGallery`, which fails if an exported
type with a `Draw*` method has no entry — and it already draws every one of them
into a rect, in every state it has. A guard in `guard` would need its own list of
components and its own way to construct each one, and a second list is a list
that drifts.

That also answers the issue's other open question. It asked whether the guard
needs each tool to enumerate its screens, and whether `guard.Screens`'s list
could serve both. It does not arise: the check is against `comp`, not against a
tool, because a tool that draws its own rectangle wrong is a tool bug and a
component that overruns is everybody's.

`State.Overlay` is the single exemption, and it has to be declared per state
rather than per component — `comp.Keys` is bounded normally and unbounded with
`Overlay: true`, which is the same distinction the field name already carries.

### What it found

`comp.Toast`, immediately. `Min` defaults to 24 columns and the width was
`clamp(r.W-margin*2, minW, maxW)`, which raises the width UP to the minimum: a
toast in a pane 10 columns wide drew 24 of them, 114 cells over whatever was
beside it. The origin was clamped to the rect and the width never was. It now
clips and bounds both dimensions — `Min` is a preference, the rect is a fact.

Three more were latent overruns in the gallery's own demo drawing, invisible at
132 and 80 columns and real at 24. Not one golden moved when they were fixed,
which is the point: **nothing that existed could have caught any of these**, and
a frame that is wrong only at a width nobody captured is exactly the bug this
project keeps finding by looking.

## 38. The reserved bottom row is a hedge, and now it says so

Issue 37: `app.Runner` sized the canvas at `terminal height - 1` and nothing in
the repo said why. Every mention was a description — *"the bottom line is left
for the terminal"*, *"which is where the old bodyHeight arithmetic already put
it"*. The second is the tell: democtl inherited it from swarmctl. It was then
promoted to a framework rule in decision 30 on the grounds that two tools did it
independently, **which is not two decisions if one was copied**.

**Measured.** tmux 3.5a, 24×10, `tea.WithAltScreen`, a frame of exactly the
terminal height with a bordered pane so the bottom-right cell is genuinely
written:

```
┌──────────────────────┐   ← row 1, intact
│ ROW01                │
…
└──────────────────────┘   ← row 10, the last cell written
```

No scroll, no lost top line. The reserved version simply leaves row 10 blank.

The pending-wrap hazard is real in general and does not fire here, because
nothing is written *after* the last cell — the flag is set and the frame ends.
The issue's own hypothesis looks right: this is an INLINE-renderer workaround
carried into alt-screen code, where it does not apply. azctl is corroborating
evidence, having run at full height for months before it migrated with no report
of a lost line.

### So why is it still the default

Because one emulator family has been measured and the failure mode is a top line
eaten on some *other* terminal. That asymmetry decides it: being wrong costs a
reader the top of their interface and a bug nobody can reproduce; being right
costs one row. Flipping the default would also move every golden in every tool,
which is a large diff to buy a row on the strength of a single measurement.

`app.WithFullHeight()` gives the row back to any tool that wants it, and the
option's doc carries the measurement so the next person inherits evidence rather
than a convention. The default flips if foot and one Windows terminal agree —
the same missing datapoint as issue 31, and worth collecting in the same sitting.

**What was actually wrong here was not the row.** It was that a hedge had been
written down as a rule, and a tool paying a cost could not find out why. That is
fixed whichever way the measurement eventually goes.

## 39. No leave-confirm in the framework: the second tool wanted the opposite

Issue 41 recorded azctl's "leaving would lose something" pattern and parked it
for a second consumer, per decision 31. The second consumer arrived, and it
settles the question the other way.

**azctl**, mid-playbook, `q` or `esc` opens a confirm:

> The run is still going. Leaving stops it where it is — **the steps that have
> already run are not undone.**

**pgctl**, mid-operation, `q` cancels immediately with no question at all
(`internal/tui/app.go:330`):

```go
case "q":
    if m.active != nil && m.active.running {
        // A running operation is cancelled rather than abandoned, so the
        // engine's failure hooks get to bring an environment back up.
        m.active.cancel()
```

Both are right. The difference is not taste and not maturity: **azctl's work
cannot be undone and pgctl's can.** A half-run playbook leaves an environment
neither finished nor untouched, so the reader has to be told before it happens.
A cancelled pgctl operation runs its failure hooks and brings the database back
up, so a confirmation would be a dialog standing between a reader and the safest
available action — and one that makes cancelling *slower* in exactly the moment
someone is trying to stop something.

An `app.Keys{Leaving: …}` field, or a stack that refuses to be popped, would
have imposed azctl's answer on pgctl. The framework cannot tell these apart,
because the question is whether the domain's work is recoverable, and that is
the engine's knowledge — decision 22 says the engine tells the UI nothing about
terminals, and this is the same boundary from the other side.

So: **refused**, and this is the `comp.Tree` outcome rather than the
`app.Toggles` one. Two tools wrote something similar-looking, and the parts that
differ are the parts that matter.

What is worth keeping is the wording rule, which both tools already follow and
which is about writing rather than about mechanism: **"are you sure" is a
question about nothing.** A confirm names what is lost, or it is theatre. That
belongs in the design notes and not in a component.

The second half of azctl's pattern also survives as a general point and is now
in `design/keys.md`: `ctrl+c` never asks — a confirmation on the universal
escape hatch is a program arguing with it — and the key that OPENED a question
must not also answer yes.

## 40. Two guards for the lies an interface tells about itself

Issues 39 and 48, both filed as "probably not tractable, do not guess". Both are
tractable once narrowed to the part that is actually checkable, and the
narrowing is the decision.

### `guard.Furniture` — a chrome character typed out by hand

Issue 48. `guard.Glyphs` asks whether a printed character is ALLOWED, and `─` is
allowed: being in the box set is the point of it. The defect is the SOURCE. A
tool that writes `c.Fill(band, "─", …)` has hardcoded the light box set, and on
`theme.ASCIIBox` it draws `─` beside the `-` that everything reading the chrome
draws — one frame, two box sets.

What made it look intractable was that a chrome character in a source file is
not always a draw: `gallery.Entry.Glyphs` and a tool's own glyph-set declaration
both NAME these characters in order to permit them, and a guard that flags those
is one every tool learns to suppress.

**Narrowing to a literal passed to a canvas draw — `Set`, `Text`, `Fill` —
excludes both, because neither is a call.** Run against six real packages it
reported two findings, both known and both real (azctl's, now `comp.Rule`), and
zero false positives. It also found a **sixth** hardcoded site nobody had
counted: `scaffold/templates/internal/tui/view.go.tmpl`, the file every new tool
starts from.

### `harness.Hints` — an advertised key that does nothing

Issue 39, and the answer to its "guard, harness check, or discipline" is: a
harness check, for less than it hoped.

The AST route dies on the same rock `guard.Keys` documents — a tool's
screen-level keys are not commands and never will be, so a footer legitimately
names more than the spec does. The runtime route works: press each advertised
key, report the ones that change nothing.

**But it does not catch the bug that prompted it.** azctl's runner promised
`q abort (the running step finishes)` while `q` cancelled the run and quit the
program. `q` did plenty — it just did not do what the footer said. That is a
claim about English, and no mechanical check reaches it.

So the guard catches the lesser sibling, the DEAD advertised key, and its doc
says so rather than implying more. The issue's own worry about false positives
is handled by the caller passing the hints rather than the frame being scraped:
a key that is legitimately inert in this state is left out, deliberately, in a
test that says why.

### The pattern in both

**A check that catches part of a class of bug is worth having if it says which
part.** The failure mode to avoid is not a narrow guard; it is a guard that
looks like it covers the class and does not, because then nobody looks for the
rest. Decision 32 remains the discipline for everything unreachable this way —
a comment or a label describing behaviour is an untested assertion, and the
answer to one is a test named after the sentence.

## 41. The kernel knows what a pixel is

Issue 31. On a HiDPI Mac, Ghostty reported a cell as 16×34 and drew at 16×34;
iTerm2 reported 8×17 and drew at 16×34, so every Sixel picture came out at
exactly half size. Both answered `CSI 14 t` honestly and disagreed about the
unit, and nothing in-band says which. The issue refused to add a scale-factor
guess until three questions were answered.

**The first question answered the other two.** `TIOCGWINSZ` is a different
channel: `CSI 14 t` is answered by the terminal application, which picks points
or device pixels; `ws_xpixel`/`ws_ypixel` are what that same terminal wrote into
the tty. Measured on one display, both terminals:

| | `CSI 14 t` | `TIOCGWINSZ` ÷ cells | draws at |
| --- | --- | --- | --- |
| Ghostty | 16×34 | **16×34** | 16×34 |
| iTerm2 | 8×17 | **16×34** | 16×34 |

The kernel channel is right in both cases, including the one where the escape
lies. So this was never a scale factor to guess — **it is a division**, and the
answer was available all along on a channel nobody had asked.

`CellSize` now asks the kernel first and falls back to the escape. It refuses an
implausible answer (a cell under 2px or over 200) rather than trusting the
fields blindly, so a terminal that fills them with something else does not take
the picture with it and the escape still gets asked. The division rounds rather
than truncates: 3832 pixels over 239 columns is 16.03, and a cell one pixel
narrow tiles a whole row of pictures short.

### What this closes and what it does not

The other two questions — is it iTerm2-specific, and does it matter on Linux —
**no longer need answering to fix the bug**, because the fix does not depend on
knowing which terminals lie. It depends on a channel that was right on both. A
terminal where the kernel is also wrong would still need `TUIKIT_CELL_SIZE`, and
that is now what the variable is for: the override of last resort, rather than
the answer to HiDPI.

Worth recording that the visible symptom was Sixel-only and always would have
been. The kitty placement states its footprint in CELLS (`c=`/`r=`), so a wrong
pixel measurement there makes a blurrier picture rather than a misplaced one —
which is why both terminals looked right in the screenshots that settled this,
both having negotiated kitty.
