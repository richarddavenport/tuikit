# Rebuilds

Could tuikit rebuild the TUIs people actually use? Each file here works one out
on paper and comes back with a list of what is missing.

**Why this and not "what components should we add".** A component built because
somebody thought it would be useful has no shape to be right or wrong about.
Decision 31 refuses those, and the extraction rule asks for two independent
implementations before generalising — because one is an anecdote.

The rule was never "wait forever". It was "do not generalise from one example",
and these tools are examples. Six of the fourteen surveyed need a collapsible
tree; five need a syntax-aware text view. That is a better sample than four
private tools, and it is available now.

**What a rebuild is not.** Not a plan to build a clone, and not a judgement of
the tool. It is a coverage test: read the real source, list what it draws, and
mark each one *have*, *hole*, or *theirs*.

**"Theirs" is the important column.** Most of what a good TUI does is its own
domain, and a framework that tried to supply it would be wrong. A commit graph
belongs to a git client. The test is not "could tuikit supply everything" — it
is "is what tuikit supplies the right set, and is anything missing that several
tools each had to build alone".

## Done

- [lazygit.md](lazygit.md) — 82k stars, the most used TUI in the field

## Method

Read the actual repository, not the README, and say which. Cite paths and file
sizes so a claim can be checked. Where something was reasoned about rather than
read, mark it — `other-tuis.md` explains why that matters.
