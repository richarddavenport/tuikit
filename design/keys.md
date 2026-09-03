# Keys

What tuikit has settled about keyboard behaviour, and what it has not.

`app.Keys` routes in one order — Capture, then Screen, then Global — and that
order is the mechanism. This file is about the handful of MEANINGS that turned
out to be worth writing down, because two tools got them wrong the same way.

## `ctrl+c` never asks

A confirmation on the universal escape hatch is a program arguing with the one
key a reader is entitled to expect works. Whatever else is in flight, `ctrl+c`
stops it now.

`q` and `esc` are different: they are the tool's own keys, they mean different
things per screen, and a tool is entitled to put a question in front of them.

## The key that opened a question must not answer yes

If `q` opens "stop the run and quit?", then `q` is the obvious key to press
twice, and pressing it twice must not quit. It means no, or it means nothing.

This is not hypothetical: it is the natural way to press a key you think did not
register.

## A confirm names what is lost

"Are you sure?" is a question about nothing. The body says what happens if the
reader says yes, in the terms of their own work:

> The run is still going. Leaving stops it where it is — the steps that have
> already run are not undone.

A reader who cannot picture the consequence is not answering the question, they
are guessing at it.

## Whether leaving should ask at all is the TOOL's decision

Decision 39. azctl asks mid-playbook because its steps cannot be undone; pgctl
cancels immediately without asking because its failure hooks bring the database
back up and a dialog would stand between a reader and the safest action. The
framework cannot tell these apart — whether the work is recoverable is the
engine's knowledge.


## The four that are reserved (decision 42)

These mean one thing each, in every tuikit tool, and `guard.Reserved` fails a
build that binds them to anything else:

| key | means |
| --- | --- |
| `ctrl+c` | stop what is happening and quit, immediately and without asking |
| `q` | quit, or leave the screen you are on |
| `esc` | go back, or dismiss what is open |
| `?` | show the keys |

**What is reserved is the meaning, not the behaviour.** A tool may put a
question in front of `q` — what leaving costs is the tool's business, and
decision 39 is the whole argument for that. It may not make `q` mean something
that is not leaving.

`ctrl+c` has no latitude at all, and nothing may even declare it: it is handled
before a tool sees it, so a command claiming it describes a binding it does not
have.
