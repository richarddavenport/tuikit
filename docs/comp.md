# `comp` — the canvas and the components

Everything you look at. A grid of cells, 25 components that draw into it, and
the text measurement they all share.

## The canvas

`comp.Canvas` is a grid of cells. Each cell holds one grapheme cluster, a style
pointer, and **the ID of what drew it**.

| | |
| --- | --- |
| `NewCanvas(w, h)` | a grid. `app.Runner` makes this for you; a test makes its own |
| `Set(x, y, cluster, style, owner) int` | one cell, returns the columns it took |
| `Text(x, y, s, style, owner) int` | a run of them, returns its width |
| `Fill(r, cluster, style, owner)` | a rect |
| `Box(r, style, owner)` | a border from the chrome's box set |
| `Clip(r) *Canvas` | a canvas that cannot draw outside `r` |
| `OwnerAt(x, y) ID` | **the entire hit test** |
| `CellAt(x, y) (Cell, bool)` | what is there, style included |
| `Region(id) (Rect, bool)` | where something ended up, this frame |
| `String() string` | **the frame** — text and escape sequences, ready to print, save or render |
| `Chrome()`, `WithChrome(ch)` | the box characters, gap width and markers in force |

**What an ID means.** `Region("services").At(7)` is *the seventh service*, not
row seven of the screen. After a scroll those differ, and an ID meaning
screen-position does not error — it performs the right action on the wrong row.
`ParseID` reads one back from a capture script.

**Clipping is not truncation.** A component drawing past its rect is clipped,
silently and correctly. Truncation is a component *deciding* a name is too long
and saying so with an ellipsis. Different jobs; `Clip` does the first, `Truncate`
the second.

## Measuring text

Columns, not bytes and not runes — a CJK glyph is two columns per rune and
`fmt.Sprintf("%-20s", s)` counts bytes, so the column that was supposed to line
up does not.

`Width` · `Truncate` · `Pad` · `Wrap` · `WrapFirst` · `Mask` (for secrets) ·
`Highlight` (a string plus matched indices becomes styled `Segment`s)

## Layout

| | |
| --- | --- |
| `Layout` + `Rows(r)` / `Cols(r)` | bands down an axis from `[]Constraint` |
| `Length(n)`, `Percent(n)`, `Fill(weight)` | the constraints |
| `.Min(n)`, `.Max(n)` | clamps on a constraint; a clamp is paid for by the others |
| `Split` | two panes and a **draggable** divider — `Draw`, `Layout`, `Move`, `MoveTo`, `Min` |
| `Rect`, `Center(r, w, h)` | rects, and centring that clamps position without shrinking size |

Splits nest. Three panes with two independent dividers is a composition, not a
component — but an inner split's rect comes out of the outer one, so after the
outer moves the inner must be handed the **new** rect.

## The components

Twenty-five. Run `tuikit gallery` to use every one in every state it has, or
`tuikit gallery -list` for the inventory as text.

### Lists and documents — three viewports, because they are not one thing

| | |
| --- | --- |
| `List` | a cursor and rows. Viewport and selection are separate, so scrolling away does not move the cursor |
| `LogPane` | tails a stream. Following is a place, not a mode: scroll up to detach, scroll back to re-attach |
| `Viewer` | a document. Opens at the **top**, scrolls sideways, and puts the cursor style *under* the line's spans so a diff keeps its syntax colors |

`List` can do:

- **grouped rows.** `Row.Skip` marks a heading the cursor passes over.
- **a lead glyph that survives the selection.** `Row.LeadStyle` keeps a status
  color readable on the row you are pointing at.
- **a right-aligned tail.** `Row.Right`.
- **rows in more than one color.** `Row.Spans`.
- **a range selection.** `Extend`, `Range`, `ClearRange`. The anchor is
  [`comp.Range`](#a-range-over-anything), which anything with a cursor can hold.
- **a status column that survives an indent.** `Row.Status` is drawn before the
  indent and `Row.Lead` after it, so a git status letter forms a column down the
  screen while a fold marker travels with its row. `List.StatusWidth` reserves
  the column and `LeadWidth()` is how wide the whole prefix is, for a header
  that has to start in the same place.
- **branches instead of an indent.** `Row.Indent` replaces the spaces
  `Row.Depth` would have drawn — see [trees](#trees).
- **a cursor that follows the row.** `Row.Key`, opt-in. Set it and a filter
  applied or cleared leaves the cursor on the same *thing* rather than the same
  *line*. Leave it empty and the cursor is an index, which is right and free for
  a list whose rows never move.
- **rows produced on demand.** `DrawFunc` asks only for what is visible, so a
  million rows cost a frame the size of the pane.
- **giving the counter row back.** `NoStatus`.

`Viewer` can do:

- **line numbers.** The gutter's width comes from the whole document, so it
  does not change as you scroll.
- **tab expansion** to the next stop, measured across the whole line.
- **no cursor at all.** `NoCursor`, for a document with nothing selected.
- **jumping to a line.** `Goto`, for a search hit or a `:` line number.
- **sideways scrolling.** `ScrollX` cuts inside a span rather than on its
  boundary, so the offset does not jump by a syntax token.
- **a highlight nobody dragged.** `SetRange(lo, hi)` works with `NoCursor` set,
  for a span a *tool* chose — the bytes belonging to the field selected in
  another pane. `Extend` still refuses, because dragging needs somewhere to drag
  from.
- **the same lazy `DrawFunc`** as `List`.

### Trees

`Tree` holds collapse state over a flattened hierarchy, keyed by the tool's own
identity for a row rather than by index. `Visible(nodes)` returns the indices to
draw.

`Branches(nodes, chrome)` returns the connector prefix for each row — `│`, `├─`,
`└─`. Which one a row gets cannot be derived from its depth: it depends on
whether ancestors still have siblings coming, which is a fact about the rows in
between. It returns **strings**, so they can go in `Row.Indent` for a list or
inside one cell of a `Table` — which is where a process tree belongs, because a
prefix in front of the PID makes every numeric column ragged.

The connectors live in `theme.Chrome.Tree`, with `theme.ASCIITree` beside the
default for a terminal that draws box characters as replacement boxes.

### A range over anything

`comp.Range` holds an **anchor**; the caller supplies the cursor. That is what
lets one type serve a list, a document, a table row number and a tree node index
without learning what any of them are.

The zero value means no selection, so a fresh component cannot report a range
from row 0. The span is derived from the cursor passed in and never stored,
which is what keeps it correct under a deferred `Move`.

It composes with `Marks` rather than competing: `Marks.Span` takes the keys a
`Range` covers.

### Subprocess output

`ANSI(s)` and `ANSILines(s)` turn SGR escape sequences into `[]Segment`, for a
pane showing what `kubectl` or `gcloud` just said.

Colors pass through rather than being remapped. A tool's own color comes from
the palette so the reader's theme wins; a subprocess's output is **data**, the
same way its words are, and a tool that does not rewrite the nouns should not
rewrite the colors either.

`\r`, `\b` and tabs are handled, because they are how a progress line is drawn
— a line is assembled in cells rather than appended to, so ten redraws come out
as the final state. Cursor addressing and scroll regions are skipped whole:
that is a terminal emulator, and hosting one is out of scope.

### Structure

`Pane` (a bordered box with a title, in the edge or on a row) · `Tabs` ·
`Bar` (content at each end, with `MinLeft` deciding when the right end is
dropped) · `Rule` · `Breadcrumb` · `Scrollbar` · `Layout` · `Split`

### Data

`Table` (columns with an alignment and a width rule; `Rows` for strings,
`Spans` for styled cells) · `Detail` (label/value blocks that line up *per
block*, so a long key does not drag the facts above it wide) · `Meter`
(characters everywhere, pixels where they exist) · `Tree` (**collapse state
only** — it draws nothing and owns no node type; `Visible`, `Toggle`,
`IsCollapsed`, `Expand`, `Collapse`, `HasChildren`)

### Input and action

`Input` (one line, with a caret you can move) · `Form` (every field at once —
`FieldText`, `FieldChoice`, `FieldToggle`, any of them `Secret`, and `Must` for
a phrase that has to be typed before the action is allowed) · `Confirm` · `Menu` (at
a point, or on a *region* so it finds out where that is at draw time) ·
`Palette` (one key to everything the tool can do now) · `Keys` (every binding by
screen, windowed with `↑ N above` / `↓ N more`)

### Progress and state

`StepList` (what will happen, what has, and what it cost) · `Spinner` (its frame
comes from the clock, so every spinner turns together) · `Waiting` (a region
whose contents have not arrived — drawn *inside* the interface, not instead of
it) · `Toast` (what the interface has to say, and where known what to do about
it)

### Pixels

`Pixels` — `Detect()`, and a picture drawn under or over the characters where
the terminal supports it. See [`term`](term.md) and [`paint`](paint.md).

## What it cannot do

- **No soft wrap in `Viewer`.** It breaks one line = one row, which the cursor,
  the gutter and every index depend on.
- **No retained tree.** Components are function calls: `Draw(canvas, rect)` and
  return. Parent and child is only who passes which rect to whom. There is
  nothing to diff and nothing to hold a `Parent` pointer.
- **No constraint solver.** `Layout` is a single linear pass down one axis.
  Two-dimensional constraints between siblings are ratatui's and bento's answer,
  not this one — nesting reaches the same arrangement for the shapes these tools
  lay out.
- **No text editing.** Decision 27's line: showing a buffer and acting on ranges
  of it is in; being the place you type the buffer is out.
- **No sort state on `Table`**, no chart over time, no selection *set*. Each is
  an open issue with the evidence attached.

## Read more

`go doc github.com/richarddavenport/tuikit/comp.<Name>` — every component's doc
comment names the tools it came from and what their versions disagreed about.
