# Prototypes

Throwaway spikes, on `prototype/*` branches, kept so the questions they answered
are not mysteries in six months. Not for merge, and not maintained — if one goes
red, read what it concluded and delete it.

## `htmlround` — is `harness.HTML` lossless?

Renders a frame to HTML, parses it back into cells, and compares. Cells rather
than bytes: lipgloss may write `\x1b[1;38;5;205m` or `\x1b[1m\x1b[38;5;205m` for
the same thing, and a test that fails on that is a test about lipgloss.

**Answer: nearly, with three named holes, all in the writer.**

1. `css` handles `38;5;N` and falls through on `38;2;R;G;B`, then reads the
   channels as separate SGR codes — so `#ff5faf` reaches the page as `#ff00ff`.
   A plausible wrong colour, which is worse than no colour.
2. `writeLine` opens a span from each SGR sequence's parameters alone, where a
   terminal ACCUMULATES state. Lossless only for self-contained runs, which is
   what lipgloss happens to emit.
3. Spans close at a line boundary and never reopen, so a style that straddles a
   newline is dropped.

Each is pinned by a test that logs when it starts passing.

## `screendoc` — can a screen be declared rather than drawn?

The question a visual builder rests on. `testdata/dashboard.json` declares
democtl's dashboard; the test renders it and compares CELL FOR CELL against the
hand-written `View`, so a role that resolves to the wrong colour fails.

**Answer: yes, and the interesting part is what the format had to grow.**

Six states and three widths reproduce exactly. Getting there forced: ratios
taken against the whole width including gaps (a third of 132 is 44, of 131 is
43, and the pane is a column narrower for the rest of the screen's life);
conditionals, in three unrelated places; conditional titles; width-relative
formats (`trim(svc.Image, w-14)` is a real line of democtl); and groups, because
conditions do not compose.

Roughly half the format is branching. **That is the finding.** A screen is not a
picture, it is a picture with conditions all through it, and expressing "when
typing, the title becomes `Filter:`" by direct manipulation is the actual design
problem of a builder — not the palette, and not the dragging.

Known gap: a list row carries one style, so the Events tab (`muted(timestamp) +
plain(text)`) is beyond it. Closing that means rows become a sequence of styled
spans everywhere, which is a real change to the format rather than a patch to
the document.

Cost: ~510 lines of Go, and a 339-line document to replace ~200 lines of view
code. The document is not smaller than the code. If a builder is ever worth
building, the case has to be that its output is constrained and machine-editable
— it cannot be brevity.

## What they found in code that is not a prototype

- **democtl, at 80 columns and narrower:** `tabStrip` ends in `pad`, whose
  `trim` counts RUNES. With colour on, the escape bytes inflate the count, the
  strip is cut mid-sequence, the terminating `\x1b[0m` is eaten and bold leaks
  to the end of the row. `padVisible` is the fix. Found independently by both
  spikes.
- **The goldens cannot see that class of bug at all.** They run under a plain
  `go test`, where lipgloss emits nothing and `trim` sees real runes — so the
  80-column golden shows a perfectly plausible `‹ Config ›  Events`. Capture
  forces `TrueColor`, so the `.ansi` frames and the real terminal have a bug the
  goldens say is not there.
