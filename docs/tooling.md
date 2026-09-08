# The commands

```
tuikit new <name>          scaffold a tool that passes its own checks
tuikit gallery [-list]     every component, running — or as text
tuikit pixels              what this terminal can draw
tuikit frames <dir> [-md]  a captured run as a page, or as Markdown + SVG
tuikit watch <dir>         recapture on save, rebuild, reload
tuikit designsystem        the colors and glyphs as HTML
tuikit news                what changed since a tool last looked
```

## `gallery` — the one to run first

Every component, in every state it has, including the states a screenshot of the
happy path never shows: empty, overflowing, and no room. A component that is not
in the gallery is not finished, and `gallery_test.go` holds that closed by
reading `comp`'s own source — adding a component without an entry fails a test.

`-list` prints the inventory as text, each entry with the tools it was extracted
from.

## `news` — how a tool learns tuikit changed

Backed by the `news` package. `design/decisions.md` is a numbered record; a
generated tool records the number it was born reconciled with, and `tuikit news`
prints the decisions since.

`Read(path)` · `Since(ds, n)` · `Latest(ds)` · `Marker(path)` ·
`Checkout(gomod)`

This is the feedback loop in the other direction from `guard`: guards tell a
tool when it has drifted from the rules, and news tells it when the rules moved.

## `watch` — the inner loop

Backed by the `watch` package (`NewPoller(dir)`, `Config`). Change a line,
recapture, rebuild the page, reload. The point is that you look at the frame
rather than imagining it — which is how the ghost in `comp.Viewer`'s sideways
scroll was found, after its tests were green.

## `pixels` and `designsystem`

`pixels` reports what the terminal you ran it in can draw: graphics protocol,
cell size, whether either was measured or guessed. Run it when an image looks
wrong.

`designsystem` renders `theme` as HTML. Generated, never hand-written — see
[docgen](docgen.md).
