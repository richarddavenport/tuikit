# `docgen` — the interface, as documentation

Turns what the tool actually drew into something you can send to someone. Both
halves are generated, never hand-maintained.

The input is a directory of `.ansi` files — frames written by
[`harness`](harness.md), each one the exact string the program would have
printed. No screenshot is taken and no terminal is involved: `docgen` parses
those escape sequences back into styled spans and re-emits them as SVG or HTML.

## Frames — a captured run as a page

```go
docgen.Frames{
	Title: "Deploying",
	Lede:  "What a deploy looks like from the dashboard.",
	Groups: []docgen.Group{{Name: "Choosing", Frames: []string{"dashboard", "picked"}}},
}.Markdown(dir, out)
```

Reads a capture directory's `.ansi` files and its manifest. Produces either an
HTML page or Markdown with one SVG per frame.

Markdown is usually what you want for a README or a pull request. It renders on
GitHub, and there is no screenshot to keep in sync.

**A frame not named by any group is appended rather than dropped.** That is the
failure mode of a hand-maintained list: a new screen silently vanishes from the
docs. Grouping is given rather than inferred because it carries something a
filename cannot — what the reader is doing.

`SVG(frame)` converts one frame on its own. It uses `textLength` with
`lengthAdjust="spacing"` so a proportional fallback font still lands on the
grid, and carries no `<style>` element, so it survives GitHub's sanitiser.

Driven by `tuikit frames <dir> [-md]`, and `tuikit watch <dir>` re-runs it on
save.

## DesignSystem — the palette and glyphs as pages

`DesignSystem` renders `theme` into `Page` cards — colors, glyphs, box sets.

There is deliberately **no hand-written palette page** anywhere in this repo,
and there should never be one: a page describing colors is a second source of
truth that goes stale the first time a role changes. `tuikit designsystem`
generates it from `theme` itself.

## What it cannot do

- **No API documentation.** That is `go doc`, and the doc comments are where the
  reasoning lives.
- **No animation.** Frames are stills; a capture is a sequence of them.
- **The SVG is text, not a screenshot.** Graphics protocols (Sixel, kitty) are
  escape sequences a browser cannot render, so a frame with a picture in it
  shows the characters around the picture.
