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

## 3. ANSI 256 indices, never truecolour hex

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
