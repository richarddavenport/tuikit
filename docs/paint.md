# `paint` — pixels

A small raster library for the things a terminal UI actually draws in pixels.
Not a general 2-D graphics package: the shape list is short, and a library that
draws everything costs more than it saves.

| | |
| --- | --- |
| `Panel` | a rounded rectangle, filled, with optional border and alpha |
| `Bar` | a progress bar — what `comp.Meter` draws when pixels exist |
| `Ramp` | a gradient between colors |
| `Flatten(src, bg)` | composite onto an opaque background |

Everything returns an `*image.RGBA`, which `term.EncodeSixel` or
`term.EncodeKitty` turns into an escape sequence and `comp.Pixels` places.

**`Flatten` is not optional for Sixel.** Sixel has no alpha, so a rounded corner
composited against nothing comes out as a hard black notch. The kitty path
composites for real and skips it. `term.Background()` is where the color to
flatten against comes from — the terminal's own, read via OSC 11, rather than an
assumed black.

**Gradients follow the reader.** `Ramp` takes colors, and the caller gets them
from `term.Colors` — the reader's actual ANSI palette — so a generated picture
matches the theme they chose instead of overriding it.

## What it cannot do

No text rendering, no paths, no transforms, no image loading or decoding, no
anti-aliased strokes beyond what the shapes do themselves. If you need a chart,
there is nothing here for it yet — that is an open issue.
