# `spec` — one declaration, four surfaces

Describe a command once. Get the command line, the TUI screen it opens, the
context menu it appears in, and a machine-readable manifest — from the same
declaration, so they cannot drift apart.

`comp` and `app` do not import `spec`. A TUI-only tool pays nothing for it.

## The declaration

```go
spec.Command{
	Name: "logs", Short: "follow a service's output",
	Args:  []spec.Arg{{Name: "service", Complete: services}},
	Flags: []spec.Flag{{Name: "since", Kind: spec.Duration}},
	Screen: "logs",              // the TUI screen it opens
	Target: "dashboard.service", // the region whose menu it joins
	Key:    "l",                 // its keyboard path
	Run:    func(c spec.Call) int { ... },
}
```

`Commands` nests subcommands. A nil `Run` makes a grouping command that prints
its children. `Hidden` keeps something working but unadvertised. `PassThrough`
collects flags the command did not declare — from the cloud tool, whose `play`
takes `--<param>` for whatever a runtime-chosen YAML playbook declares.

## What it produces

| | |
| --- | --- |
| `Run(root, argv, out, errw) int` | the CLI: parsing, help, exit codes |
| `Usage(cmd, path, w)` | the help text |
| `Complete(root, argv)` | completions, args and flags both |
| `CompletionScript(shell, tool, w)` | the shell script that calls it |
| `MenuFor(root, region)` | the context menu for a region, as `[]comp.Hint` |
| `WithKeys(root)` | every command that has a key, for a help screen |
| `Describe(root, version, palette, glyphs)` | the whole surface as a `Manifest` |
| `SchemaOf(cmd)` | JSON Schema, for tool-calling APIs |
| `Unreachable(root)` | commands with a `Target` and no `Key` |
| `Parse(cmd, argv)` | a `Call` — typed args and flags |

**Why the manifest matters.** An agent working out what a tool can do otherwise
greps a README that is a version behind, parses `--help`, or reads the source —
three answers that disagree in small ways with no way to tell which is current.
`Describe` is generated from the same declarations the CLI runs and the TUI
draws, so it cannot be stale without them being stale too. It carries the
palette and glyph set as well, because an agent adding a screen needs to know
which colour roles and characters exist *before* it writes anything, and those
are exactly what a guard will fail it for afterwards.

## Exit codes, and why cobra is not here

```
0  it worked, or a dry run found nothing to do
1  it failed
2  a dry run found DRIFT — not an error, a finding
```

`OK` · `Fail` · `Drift`. The two is the point: `the deploy tool diff && deploy`
must not deploy when there is drift.

## Argument and flag kinds

`Bool` · `String` · `Int` · `Duration` — and `Completer`, a `func(prefix)
[]string` so completion asks the tool rather than guessing. `SnapshotFlags` is
the shared set for capture-producing commands.

## Reserved keys

`spec.Reserved` names four keys that always mean one thing: `ctrl+c` and `q`
quit, `esc` goes back, `?` opens help. `Leaves(name)` and `Helps(name)` answer
which is which, and `guard.Reserved` fails a build that binds one to something
else. This is the only thing tuikit imposes on a tool's keyboard.

## What it cannot do

- **No positional-flag interleaving beyond what's declared**, no flag groups, no
  mutually-exclusive sets.
- **Not a router.** `Screen` is a name, not a path with parameters.
- **No i18n.** Help text is the string you wrote.
- **No key sequences** — a `Key` is one keystroke.
