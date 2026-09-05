# `harness` — testing and capture

Drive a model without a terminal, look at what it drew, and hold it to a golden.

## Driving

`Press(m, "j", "j", "enter")` · `Run(m, key)` · `Resize(m, w, h)`

Mouse, addressed by **region name** rather than by coordinates — the whole point
of owner IDs:

`Click` · `RClick` · `DoubleClick` · `Hover` · `Drag(t, m, region, n)` ·
`Wheel(t, m, region, notches)`

`Script(t, m, script)` runs a whole sequence as text, and `ScriptFile` loads one
from disk. A script says `click dashboard.service[3]`, not `click 41,12`, so it
survives the layout changing.

## Looking at a frame

| | |
| --- | --- |
| `Strip(frame)` | the text with every escape sequence removed |
| `Lines(frame)` | those, split |
| `Width(frame)` | the widest line, in columns |
| `Rows(frame) [][]Span` | text **and** style, per run — for asserting a colour |
| `HTML(frame)`, `FrameCSS` | the frame as a page |

## Goldens

```go
harness.Golden(t, "testdata", "dashboard", frame)
```

Colour-stripped, so a palette change does not rewrite every fixture and a
diff shows what moved rather than what was recoloured. `-update-goldens`
rewrites them.

`ShapeSurvivesColour(t, name, render)` is the other half: it asserts the frame
still reads with colour removed — that no distinction is carried by hue alone.

`Hints(t, build, hints...)` presses every advertised key against a fresh model
and fails on one that does nothing. Together with `guard.Keys` — which checks a
binding has a help entry — that closes the loop on a lying footer.

## Capturing a run

```go
s := harness.Capture(t, "testdata/frames", harness.Size(132, 38), harness.At(fixedTime))
```

`Snapshot(m, dir, script, opts...)` runs a script and writes one `Frame` per
step, with its note. `tuikit frames <dir>` turns those into an HTML page or into
Markdown with an SVG per frame, for a README or a pull request. `tuikit watch
<dir>` recaptures on save.

`Mode` picks where the data comes from — `Fixture` for a canned estate, `Live`
against the real thing. `Enabled(env)` reads the environment variable that gates
the live one, so a capture suite is opt-in rather than flaky by default.

`At(t)` fixes the clock. Spinners and double-clicks both take their time from
it, so a capture is byte-identical between runs.

## What it cannot do

- **No real terminal.** It never opens a PTY. What it tests is what your model
  drew, not what a terminal did with it.
- **No pixel comparison.** Graphics are escape sequences in the frame; the
  goldens hold the characters around them, not the image.
- **No timing assertions.** The clock is a fixture, which is what makes captures
  reproducible and what stops you testing that something took 200ms.
