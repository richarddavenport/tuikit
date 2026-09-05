# `term` — what this terminal can do

Capability detection, done once, cached, and overridable. Nothing here draws;
`comp.Pixels` and `paint` do that.

## Graphics

`term.Detect()` returns `None`, `Sixel` or `Kitty`.

`Probe(f, timeout)` and `ProbeTTY(timeout)` ask the terminal directly, sending
`QueryBytes` — a kitty graphics query followed by a device-attributes request —
and reading the reply. The second request is the trick: a terminal that ignores
the first still answers the second, so the read always terminates instead of
waiting out the timeout on every start.

`Parse(reply)` and `Parseable(s)` are the pure halves, so the decision is
testable without a terminal.

`EncodeSixel(img)` · `EncodeKitty(img, id, cols, rows)` · `DeleteKitty(id)`

`EnvOverride` (`TUIKIT_GRAPHICS`) forces the answer. Set it to `none` in CI.

## Cell size

`WindowPixels(f)` and `WindowCells(f)` use the `TIOCGWINSZ` ioctl — the kernel
already knows, and it costs no round trip. `QueryCellSize(f, timeout)` falls
back to `CSI 14t` when the kernel does not fill in the pixel fields, and
`CellSize()` puts the two together and caches.

`DefaultCellW`/`DefaultCellH` are the guess when neither works.
`EnvCellSize` (`TUIKIT_CELL_SIZE`) overrides.

Why it matters: an image is sized in pixels and placed in cells, and getting the
ratio wrong stretches every picture.

## Colours

`Background()` and `QueryBackground` read the terminal's actual background via
OSC 11, so a Sixel image — which has no alpha — can be flattened against the
right colour instead of against a guess.

`Colors(indices...)` reads ANSI palette entries via OSC 4. That is how a
generated gradient follows the reader's own theme rather than imposing one.

`Multiplexer()` names tmux or screen when it can, because both intercept
sequences and change what a probe means.

## What it cannot do

- **Two graphics protocols**, not six. yazi ships kitty, kitty-old, iTerm,
  Sixel, chafa and ueberzug; there is **no half-block fallback** here, so a
  terminal with neither gets characters. A known floor, not a boundary.
- **No terminfo.** It asks the terminal, not a database.
- **Probing needs a real TTY.** Under a pipe or in a test it returns the
  default; that is what the environment overrides are for.
