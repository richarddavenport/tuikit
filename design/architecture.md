# Architecture

How tuikit is put together, and why the seams fall where they do. For the dated
record of choices — including what was rejected — see [decisions.md](decisions.md).

## Why this exists

Four tools now share a shape: `swarmctl` (65k lines), `pgctl` (11k), `dugo`
(3.8k), `azctl` (3.2k). Each has `cmd/` + `internal/engine` + `internal/tui` +
`internal/cli`, bubbletea and lipgloss, a `design/` directory, a Makefile with
`## target:` self-help and ldflags version stamping. azctl's own notes say it is
"modelled on swarmctl's engine/tui/cli split"; dugo's say its project
infrastructure was "ported from swarmctl".

That travel happens by copy and by re-explanation. A fix to a pane lands in one
repo. The theme was extracted once, in swarmctl, and the other three did not get
it. This is the argument for a module.

## The two ideas

### 1. One declaration, three surfaces

A command is declared once and produces a CLI command, a TUI screen, and an
entry in a machine-readable manifest.

```go
spec.Command{
    Name:  "status",
    Short: "current state of the stack",
    Args:  spec.Args{{Name: "env", Required: true, Complete: envs}},
    Flags: spec.Flags{{Name: "watch", Kind: spec.Bool, Short: "follow changes"}},
    Run:   func(ctx spec.Ctx, in Input) (Output, error) { ... },
}
```

- **CLI** — parsing, `--help`, `--json`, shell completions, and the exit-code
  contract, all generated. No hand-maintained usage heredoc to drift.
- **TUI** — a screen with the standard pane chrome, keymap, and async
  conventions already wired.
- **Manifest** — `tool describe --json` emits every command, flag, screen, key
  binding, colour role and glyph. An agent reads the tool's whole surface in one
  call instead of grepping for it.

The third surface is the point. It is what "agents understand it out of the
gate" actually means in practice.

### 2. Nothing in the interface is unnamed

Lifted wholesale from swarmctl's `internal/tui/theme`, which already solved
this and wrote down why:

- **Nine colour roles** — `Accent`, `Muted`, `Border`, `Danger`, `Pending`,
  `Success`, `Stderr`, `SelectionBG`, `SelectionFG`. ANSI 256 indices, not
  truecolour hex, "because that is what a terminal understands and what every
  terminal has agreed on; a truecolour hex would look right on this machine and
  wrong over ssh from another." Named by role, never by hue.
- **A glyph allow-list** — every non-ASCII character the interface may print,
  each with a note on why it is safe.
- **`Hex()`** computed from the 256-palette formula, not a table, so anything
  rendering outside a terminal starts from the same numbers.
- **Two guard tests** — `guard.Tokens(t, pkg)` fails on any `lipgloss.Color`
  literal in the UI package and on any role declared but never drawn with;
  `guard.Glyphs(t, pkg)` AST-parses the package and fails on a printed
  character outside the set. Parsed rather than pattern-matched, so a comment
  explaining a glyph is not read as printing one.

These become exported test helpers. Every tuikit tool gets both in one line.

## Mechanism vs policy — where the line falls

The framework owes two different kinds of thing, and confusing them is how a
library ends up with opinions about somebody else's CI.

**Mechanism lives in the library.** `harness.Capture`, `harness.Golden`,
`guard.Width`, `guard.Tokens`, `guard.Glyphs` are functions. They understand
bubbletea models, terminal columns, ANSI, and Go ASTs. They understand nothing
about Postgres or Docker Swarm, and they do not *run* anywhere — the tool calls
them. If it has to work for every tool, it is mechanism.

**Policy lives in the scaffolder.** `tuikit new` writes the tool a `_test.go`
that already calls the guards and a workflow that already runs `go test ./...`.
That is an opinion about how a tool ought to be set up, and it is delivered as
generated code the tool owns and can edit on day one. If it is a default rather
than a law, it is policy.

**Everything backend-shaped belongs to the tool.** Its fixtures, whether it
captures live frames at all, and what it needs standing up to do so. pgctl needs
a database; swarmctl would need a swarm; a tool with no backend needs neither.
No answer tuikit could give would be right for all three. Its only obligation is
not to make it hard — which the skip-unless-env-var pattern already discharges.

This is also the answer to a question an earlier draft asked and should not
have: *how does live-mode capture get a backend in CI?* It does not,
and it is not tuikit's question. Live capture is something a person or an agent
runs deliberately to get frames to look at. What goes in a tool's CI is fixtures,
goldens, and the guards — and even that is the scaffolder's default, not a rule.

## tuikit's own CI

The library needs a tool to test itself against, so `examples/democtl` is part of
the module: a small, backendless tool with a handful of screens, a table, a
modal, a log pane and a step list. It exists to be captured, golden-tested and
guard-checked by tuikit's own suite — proving the harness works without a
database anywhere near it.

It has a second job. It is the reference implementation an agent reads when it
wants to know what using tuikit looks like, which is worth more than a page of
prose about the same thing.

## Packages

All mechanism. Nothing here decides when it runs.

| Package | What it holds |
|---|---|
| `theme` | Roles, glyphs, `Hex`, light/dark handling. Extracted from swarmctl. |
| `comp` | Components written 2–4 times already: pane, tabs, table, list with filter, logs pane (search/filter/timestamps/scroll), step list, form, confirm, toast, spinner, key hint bar. |
| `app` | The bubbletea shell: screen router, focus model, the async conventions as *types* rather than prose — generation counters, single-flight polling, bounded poll failure, the `capturesKeys()` mode split. |
| `spec` | Command declarations → CLI + TUI screen + manifest. stdlib `flag`, no cobra. |
| `guard` | `Tokens`, `Glyphs`, and a `Screens` check that every screen constant has a `View()` case. |
| `harness` | Headless drive, deterministic frames, golden tests, ANSI capture, ANSI→HTML. Generalised from pgctl `57a13ad`. |
| `docgen` | Manifest → mkdocs-material site, and a capture directory → a frames page. Docs cannot drift from the code. |
| `cmd/tuikit` | `new` (scaffold), `watch` (recapture on save, serve, reload), `gallery` (the running component browser). |
| `safety` | `secret: true` scrubbing; the `assert` primitive that runs under `--dry-run`; protected-environment refusal. Both learned the hard way in azctl. |

## Self-documenting

`docgen` reads the same manifest and writes the mkdocs site: CLI reference, key
reference, component sheet, and the generated design-system bundle swarmctl's
`cmd/designsystem` already produces as HTML. swarmctl's own reasoning for
generating rather than maintaining it applies to all of it: *"a palette written
by hand is a palette that drifts, and a drifted design system describes a tool
that does not exist."*

## Scaffolder

`tuikit new <name>` lays down: `cmd/<name>`, `internal/engine`,
`internal/tui`, `internal/cli`, the Makefile with `## target:` help and ldflags
stamping, three workflows (CI: gofmt/golangci-lint/test/cross-build · release on
semver tag → darwin+linux × arm64/amd64 → gh release · docs → Pages),
`.golangci.yml`, `install.sh`, `design/decisions.md`, and the guard tests
already wired. Plus an `AGENTS.md` that tells an agent about `describe`,
`--snapshot`, and the goldens.

The library carries what can be upgraded; the scaffolder carries what cannot —
and it carries the defaults, since a default is an opinion and opinions belong in
code the tool owns.
