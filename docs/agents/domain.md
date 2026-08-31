# Domain Docs

How the engineering skills should consume this repo's documentation when
exploring the codebase.

## Before exploring, read these

- **`CONTEXT.md`** at the repo root — the vocabulary.
- **`design/index.md`** — the pointer page for the design notes, then the topic
  doc for whatever you are about to touch: `architecture.md` (how the packages
  fit, and where the mechanism/policy line falls), `capture.md` (seeing what you
  are building), `mouse.md` (the cell grid and first-class mouse support),
  `roadmap.md` (what is built, what is next, what is undecided).
- **`design/decisions.md`** — the numbered decision log. **This repo's ADRs.**

There is no `docs/adr/` here, and there should not be. `design/decisions.md` is
the single decision record, and a second one beside it is exactly the drift this
project exists to prevent — see decision 7, and the commit that folded an
overlapping "Decisions taken" list back into the log.

If a file you expect is missing, **proceed silently**. Don't flag its absence or
suggest creating it upfront; `/domain-modeling` adds to these lazily, when a term
or a decision actually gets resolved.

## File structure

Single-context repo:

```
/
├── AGENTS.md
├── CONTEXT.md              ← the vocabulary
├── design/
│   ├── index.md            ← pointer page
│   ├── decisions.md        ← the numbered decision log (the ADRs)
│   ├── architecture.md
│   ├── capture.md
│   ├── mouse.md
│   └── roadmap.md
├── docs/agents/            ← this file, and the tracker/label config
├── theme/  guard/  docgen/  cmd/  examples/
```

## Writing a decision

Append to `design/decisions.md` as the next numbered entry, newest last. A
decision belongs there when it **closed off an alternative someone would
otherwise reach for** — record what was rejected and why, not only what was
chosen. Keep the existing voice: prose that explains the reasoning, not a
template with headings.

Do not renumber existing entries. Do not start a second list anywhere else,
including in a design doc's body.

## Use the glossary's vocabulary

When your output names a domain concept — an issue title, a refactor proposal, a
hypothesis, a test name — use the term as defined in `CONTEXT.md`. Don't drift to
synonyms the glossary avoids.

If the concept you need isn't in the glossary, that's a signal: either you're
inventing language the project doesn't use (reconsider), or there's a real gap
(note it for `/domain-modeling`).

## Flag decision conflicts

If your output contradicts an existing entry, surface it explicitly rather than
silently overriding:

> _Contradicts decision 18 (components draw cells, not strings) — but worth
> reopening because…_
