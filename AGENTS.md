# tuikit

A TUI framework for developers and agents. Go, Bubble Tea, Lip Gloss.

See `design/index.md` for how it is put together and why.

## Agent skills

### Issue tracker

GitHub Issues on `richarddavenport/tuikit`, via the `gh` CLI. See `docs/agents/issue-tracker.md`.

### Triage labels

The five canonical roles, unchanged: `needs-triage`, `needs-info`, `ready-for-agent`,
`ready-for-human`, `wontfix`. See `docs/agents/triage-labels.md`.

### Domain docs

Single-context. Vocabulary in `CONTEXT.md`; decisions in `design/decisions.md`, not `docs/adr/`.
See `docs/agents/domain.md`.

## When comp cannot do something

You will hit this. The rule is **not** to work around it quietly.

If you are about to write `c.Text`, `c.Fill`, `c.Set`, or arithmetic on a Rect,
stop: the canvas exists so a tool never does that, and a tool that had to is a
tool the framework failed. Work around it if you must — you have a thing to
ship — but **report it in the same commit**, or it will not be reported.

```sh
gh issue create -R richarddavenport/tuikit --template from-a-tool.md
```

The template asks for four things, and the second is the one that matters:

1. what the tool was doing
2. **the code you wrote by hand** — the actual code, not a description
3. which component you reached for and what stopped you
4. whether you computed a coordinate

Why the code: a component gets built when two tools have hand-rolled the SAME
thing, and "the same" is a judgement nobody can make from a description.
`app.Toggles` exists because two tools' versions turned out to be identical.
`comp.Tree` does *not* exist because two tools' versions turned out to differ in
a way that mattered — one carried a bucket and a resource, the other a service
and a change, and a shared node type would have made both of them worse.

A report with the code is evidence. A report without it is a feature request,
and tuikit does not take those.

**A guess at the API is welcome and a wrong guess is useful.** Say if you have
none.

## Where a tool's files go

Config: `$XDG_CONFIG_HOME/<tool>/`, else `~/.config/<tool>/`. State:
`$XDG_STATE_HOME/<tool>/`, else `~/.local/state/<tool>/`. On **every** platform,
macOS included.

**Not `os.UserConfigDir()`.** It ignores XDG on darwin and answers
`~/Library/Application Support`, which is right for an application with a bundle
identifier and wrong for a command-line tool — and unsyncable, so a config
written there has to be written again on the next machine. Decision 33 has the
evidence, including the two tools whose identical search code lands somewhere a
third tool's never will.

Config is written by a person and belongs in a dotfiles repository. State is
written by the tool and must not follow anyone to another machine. They are
different directories.

`tuikit new` writes `internal/engine/paths.go` and its test into a generated
tool. That file is **copied, not imported** — config is engine work, and
`guard.Engine` denies the engine every tuikit import. It is fifteen lines the
tool owns; change them.
