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
