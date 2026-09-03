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

## Showing what a tool looks like

A capture writes `.ansi` frames and a manifest. Two things read them:

```sh
tuikit frames /tmp/frames -out page.html          # one self-contained page
tuikit frames /tmp/frames -md -out docs/screens.md  # Markdown + an SVG per frame
```

The page is for looking and for sending someone a link. The Markdown is for
what lives in a repository — a README, an mkdocs site, a pull request — and
writes its images into a directory named after it (`docs/screens.md` →
`docs/screens/`).

Both are generated from the capture and hold no state of their own: run it
again and what changed in the tool is what changes in the diff. That is what
makes the Markdown worth committing. Do not hand-edit it.

Images are SVG, never PNG — decision 34. A PNG needs a typeface and decision 30
is deliberately still open.

## Telling the tools what changed

`comp` cannot do something → a tool files an issue here. This is the other
direction, and it needs its own mechanism because there is no upgrade event to
attach one to: a tool resolves tuikit through `replace => ../tuikit`, so a pull
here changes its behaviour with no version to bump and nothing to read.

`design/decisions.md` is the marker. It is numbered and append-only, and a tool
records the number it has reconciled with:

```
Reconciled with tuikit through decision 34.
```

Run from inside the tool, with no arguments — it reads the tool's own `go.mod`
for the tuikit it builds against, and its `AGENTS.md` for that line:

```sh
tuikit news
```

**So when a change here affects the tools, write it a decision.** A change with
no decision is one `tuikit news` cannot report, and the cost of that is paid by
whoever is surprised by it later. `tuikit gallery -list` is the other half — the
complete inventory, held closed against `comp` by a test — for the additions
that are new components rather than new rules.

Do not hand-write "what tuikit has that this tool has not taken up" into a
tool's AGENTS.md. One of those existed, said "One thing", and named two of the
nine that had landed since (decision 32: a comment describing what the code does
is an untested assertion, and a list of features is the same claim).
