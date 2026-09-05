# `scaffold` — starting a tool

```sh
tuikit new mytool -dir ..
cd ../mytool && make check   # it builds, runs, and passes its own guards
```

`scaffold.Bootstrap(root)` writes a tool that is **already correct**: a
`spec.Command` tree, an `app.Model` with a screen, a `theme` wired up, a
`Makefile`, and a guard test that fails if any of it drifts.

The point is the last part. A generated tool that passed no checks would teach
the wrong habits from the first commit; this one starts with `guard.Tokens`,
`guard.Glyphs`, `guard.Keys` and the rest already green, so the first thing you
break, you notice.

`scaffold.Tool` carries the choices: `Name`, `Module` (defaults to
`github.com/<user>/<name>`), `Short` — one sentence that goes in the README, the
CLI help *and* the manifest an agent reads, said once — and `Decision`, the
tuikit decision number the tool records as reconciled.

**`Decision` defaults to the highest in the checkout it builds against.** A tool
generated today is by definition up to date, and one born stale would report
news it has no history with. See [tooling](tooling.md).

## What it cannot do

It scaffolds once. There is no `tuikit upgrade` that rewrites a tool in place —
`tuikit news` tells you what changed and you apply it, because a generator that
edits code it did not write is a generator that will eventually overwrite
something you meant.
