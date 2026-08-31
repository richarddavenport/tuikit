# democtl

tuikit's example tool, its reference implementation, the fixture its capture
harness renders — and, until `tuikit new` exists, **the supported way to start a
tool**. See "Starting a tool" in the root README.

```sh
go run ./examples/democtl          # a fictional service fleet
go run ./examples/democtl -seed 7  # a different one, just as reproducible
```

## Why it has no backend

Everything democtl shows is generated from a seed. That is not a shortcut around
writing a real integration: the harness captures democtl's screens as tuikit's
own test fixture, and a fixture that changes between runs is not one. Given the
same seed, `fleet` returns the same services, the same log lines and the same
step timings, forever. `fleet.Epoch` fixes the clock for the same reason — a
frame reading "47s ago" has to read that tomorrow too.

## What it is a reference for

| Convention | Where |
|---|---|
| Engine and UI split, engine with no UI imports | `fleet/` and `ui/` |
| Every colour from a `theme.Palette`, every glyph from a `theme.GlyphSet` | `ui/style.go` |
| The guards, as a tool actually writes them | `ui/guard_test.go` — six lines |
| The mode split: while something captures keys, `j` is the letter j | `ui/model.go`, `capturesKeys` |
| A generation counter, so an abandoned run cannot draw into its successor | `ui/update.go`, `stepDoneMsg.gen` |
| Every screen constant needs a `View()` case | `ui/view.go`, and the test for it |
| Truncating plain text and clipping coloured text are different operations | `ui/text.go` |
| A modal bounded to the terminal | `ui/view.go`, `overlay` |

## What the glyph set costs, on purpose

democtl uses tuikit's default allow-list unchanged, which contains six
box-drawing characters: `┌ ─ ┐ │ └ ┘`. There are no tee or cross pieces, so the
panes have no dividers and a titled box carries its title *inside* the top edge
rather than breaking it. That is the design system working. A shape the set does
not have needs a different design, not a different font.

## Two bugs its own frames found

Both were fixed before this landed, and both have regression tests:

1. **The modal punched a hole through the panes.** The overlay replaced whole
   background lines instead of compositing over them, so the frame vanished for
   the dialog's height. Eight of eighteen framed rows destroyed, and no assertion
   noticed.
2. **Step durations were clipped to `400m…`.** They were right-aligned against
   the terminal's width rather than the box's inside, putting them two columns
   past the edge.

Neither is subtle once seen. Neither was visible in a test that checked content
rather than shape. This is the argument for capture, made by the tool that exists
to demonstrate it.
