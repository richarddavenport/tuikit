# tuikit

**Terminal tools that look like someone designed them.** A Go framework for the
interfaces that show you state and let you act on it — a git client, a cluster
browser, a deploy tool. Built on
[Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss).

```
democtl                                                           estate: prod

┌─ services ───────────────────┐ ┌─ api-gateway ─────────────────────────────┐
│ ● api-gateway    running     │ │ image      ghcr.io/acme/api:1.9.3         │
│ ● auth           running     │ │ replicas   4/4                            │
│ ● billing      → queued      │ │ updated    2m ago                         │
│ ● search         running     │ │                                           │
│ ✗ mailer         failed      │ │ ▾ health                                  │
│ ● cdn            running     │ │   ✓ readiness        passing              │
│ ● queue          running     │ │   ✓ liveness         passing              │
│                              │ │   • migrations       running…             │
│ 7 services · 1 failing       │ │                                           │
└──────────────────────────────┘ └───────────────────────────────────────────┘

j/k move · enter open · d deploy · / filter · q quit
```

That is one `Draw` function. Every box, every glyph, every color role in it is
a component you did not write.

```go
func (m *Model) Draw(c *comp.Canvas, r comp.Rect) {
	band := comp.Layout{Constraints: []comp.Constraint{
		comp.Length(1), // title
		comp.Fill(1),   // body
		comp.Length(1), // hints
	}}.Rows(r)

	comp.Bar{
		Left:  []comp.Segment{{Text: "democtl"}},
		Right: []comp.Segment{{Text: "estate: " + m.estate}},
	}.Draw(c, band[0], comp.Region("title"))

	left, right := m.split.Draw(c, band[1])

	inside := comp.Pane{Title: "services", Focused: true}.Draw(c, left, comp.Region("services"))
	m.list.Draw(c, inside, m.rows())

	inside = comp.Pane{Title: m.selected().Name}.Draw(c, right, comp.Region("detail"))
	comp.Detail{Blocks: m.blocks()}.Draw(c, inside, comp.Region("detail.body"))

	comp.KeyHints(c, band[2], comp.Region("hints"), nil, m.hints()...)
}
```

No state of its own, no reactivity, no virtual DOM: pass the state, get the
frame. A component is handed a rectangle, it paints into it, it returns. There
is no tree, so there is nothing to keep in step with the drawing.

## It is themed by whatever themed your terminal

No config file. No loader. Nothing to ask the user for. A tuikit interface sends
the first sixteen ANSI color **indices** — the only colors a terminal lets its
user redefine — so the tool your reader opens is already in their palette, on
the first frame. Nine roles, named by role and never by hue: `Accent` survives
someone deciding the interface should be blue.

## 25 components, and that's the whole set

```
frame    Pane · Rule · Layout · Split · Bar
data     List · Table · Detail · LogPane · Viewer · Sparkline
nav      Tabs · Breadcrumb · Menu · Scrollbar · Keys · Palette
status   StepList · Spinner · Waiting · Meter · Toast · Confirm
entry    Input · Form
```

Held closed by a test: a component `comp` has and the gallery does not is a test
failure, and so is the reverse. `tuikit gallery` runs every one of them in every
state it has.

A component exists when two tools hand-rolled the same thing and their versions
did not meaningfully differ. `Pane` is the clearest case: one tool had **two**
box-drawing functions that disagreed with each other, written months apart for
the same job. A tool duplicating a function inside itself is the strongest
argument for extraction there is, because there is no second tool's
requirements to blame it on.

## What it cannot do

- **Not a text editor, not a multiplexer.** No modal editing, no undo history,
  no PTYs. If your interface needs one, you want theirs.
- **No shadows, no blur, no easing.** A terminal has no z-axis: a menu is on top
  because it was drawn last, and that is the entire implementation — no
  compositing step.
- **29 glyphs, and the list is closed**, enforced by a test. No icon font, no
  emoji. A terminal font without a glyph draws a replacement box, which reads as
  a bug rather than as decoration. Adding one is a line naming the character and
  the reason.

## What else is in it

**One declaration, four surfaces.** A `spec.Command` carries its flags, args,
help, key binding and target region. From that, `spec` generates the CLI, the
completion script, `describe --json`, the context menu for whatever region it
acts on, and — via `spec.SchemaOf` — a JSON Schema for tool-calling APIs.

**Three viewports, because they are not the same thing.** `List` has a cursor
and rows. `LogPane` tails a stream and re-attaches when you scroll back to the
bottom. `Viewer` is a document: it opens at the top, scrolls sideways, and puts
the cursor style *underneath* the line's own spans — so the line you are reading
in a diff keeps its syntax colors. Seven of the nine tools in the rebuild
studies had written that last one themselves.

**Clicks that resolve to identity, not coordinates.** `comp.Canvas` is a grid
of cells; each records the `ID` of what drew it. `OwnerAt(x, y)` is the entire
hit test. An ID names *the item*, not the row it landed on, so a list that
scrolls does not act on the wrong thing.

**Panes compose.** `comp.Layout` divides a rect into any number of bands down
one axis; `comp.Split` gives two of them a divider you can drag, with a `Min` so
neither can be dragged to nothing. Splits nest, so three panes with two
independent dividers is a composition rather than a component.

**Guards that fail the build.** Each of these is a test failure rather than a
review comment:

- a color outside the palette
- a character outside the glyph set
- chrome the glyph set cannot print
- an engine that imported a terminal library
- a key binding with no help entry
- a reserved key bound to the wrong meaning
- a component drawing outside its rect

**Screens you can look at without a terminal.** `harness` renders any screen to
a file. `tuikit frames` turns a captured run into an HTML page or into Markdown
with an SVG per frame, for a README or a pull request. Goldens are
color-stripped, so a layout change is an ordinary test failure.

**Themed by the reader's terminal.** The nine color roles resolve to ANSI
indices 0–15 — the only colors a terminal lets its user redefine — so a tuikit
tool looks like the rest of that person's terminal rather than like tuikit.

**Pixels where the terminal has them.** On kitty-protocol or Sixel terminals,
components can offer a picture — a gradient meter, a rounded panel — drawn
*behind* the characters. The characters are always the drawing; the picture is a
decoration pass, so a terminal that cannot show one loses nothing but the
gradient.

## What it cannot do, at length

**No text editor, no multiplexer.** Modal editing, undo history, syntax-aware
buffers over large files, PTY hosting and process management are out of scope and
not planned. `Input` is one line; there is no text area yet.

**No reactivity, no virtual DOM.** Redraw a snapshot each frame. `app.Gen`
discards the result of a read the reader walked away from, which is the whole
async problem a tool like this has.

**Coverage is honest, not complete.** The [rebuilds repository](https://github.com/richarddavenport/tuikit-rebuilds) works out
whether tuikit could rebuild the TUIs people actually use, one at a time, from
their source. lazygit needs three components that do not exist yet. Those files
are the current list of holes, kept public because a framework claiming to cover
a field should show its own gaps.

**Layout is bands, not a constraint solver.** `comp.Layout` divides a rect down
one axis into `Length`, `Percent` and `Fill` bands, each with an optional
`.Min()` and `.Max()`, resolved in a single pass. That covers a header, a body
and a footer, and a pane split in two.

It will not resolve competing rules — "B is twice A", "prefer 20 but accept 12",
"these three are equal unless the fourth needs room". [ratatui](https://ratatui.rs)
and [bento](https://github.com/metafates/bento) use a Cassowary solver for
exactly that, which is the right answer for a library that has to lay out
anything at all. If your interface needs one, you want theirs.

**Text is always terminal characters.** No embedded typeface, no text larger
than one cell. Pictures do gradients and curves; type is cells.

**Terminal queries are unix-only.** Graphics detection, cell size, background
and palette reads need a tty ioctl. On Windows they return nothing and the
interface falls back to characters — which works, and is untested.

**Images do not survive a multiplexer.** tmux and similar do not forward the
kitty protocol's APC sequences without configuration, while the capability query
gets through — so it reports support and draws nothing. Sixel survives.
`term.Multiplexer()` names the one you are in.

**Bubble Tea is not abstracted.** `app.Model` is `Init`/`Update`/`Draw`, and
`Update` takes a `tea.Msg`. This is a shape on top of Bubble Tea, not a
replacement runtime.

**Components are extracted, not designed.** A component exists here when two
tools hand-rolled the same thing and their versions did not meaningfully differ.
So there is no `Tree`, no general `Dialog`, no form validation and no layout DSL
— not because they are hard, but because one tool needing something is not
evidence about its shape. `design/decisions.md` records each refusal and what
would reverse it.

**No v1, no tag.** The API moves. Consumers currently resolve it with a
`replace` directive.

## The parts

Each links to a page on what it can do, what its API is, and what it will not
do. [`docs/`](docs/) is the index.

| | |
| --- | --- |
| [`comp`](docs/comp.md) | the cell canvas and the components |
| [`app`](docs/app.md) | the shell: key routing, mouse, screen stack, async generations |
| [`spec`](docs/spec.md) | one command declaration, and the surfaces it produces |
| [`theme`](docs/theme.md) | nine color roles, a closed glyph set, the box characters |
| [`guard`](docs/guard.md) | the tests that hold all of the above closed |
| [`harness`](docs/harness.md) | drive a model, capture frames, compare goldens |
| [`term`](docs/term.md) | what this terminal can do — graphics, cell size, palette |
| [`paint`](docs/paint.md) | gradients, rounded panels and bars, as images |
| [`fuzzy`](docs/fuzzy.md) | ranked matching that reports *where* it matched |
| [`docgen`](docs/docgen.md) | frames and the design system, as HTML or Markdown |
| [`scaffold`](docs/scaffold.md) | `tuikit new` |

## Commands

```
tuikit new <name>          scaffold a tool that passes its own checks
tuikit gallery [-list]     every component, running — or as text
tuikit pixels              what this terminal can draw
tuikit frames <dir> [-md]  a captured run as a page, or as Markdown + SVG
tuikit watch <dir>         recapture on save, rebuild, reload
tuikit designsystem        the colors and glyphs as HTML
tuikit news                what changed since a tool last looked
```

## Start

```sh
go get github.com/richarddavenport/tuikit
```

`tuikit gallery` shows every component in every state. `tuikit new` scaffolds a
tool. Go 1.25. MIT.

## Reading further

[`docs/`](docs/) is the user guide: one page per package, what it can do and
what it cannot. [`docs/tooling.md`](docs/tooling.md) covers the commands.

`design/decisions.md` is the numbered record of every choice and the reasoning
behind it, including the things that were refused. `design/architecture.md` is
how the packages fit together. `AGENTS.md` is what an agent working in a tuikit
tool needs to know.

## License

MIT. See [LICENSE](LICENSE).
