# tuikit

### Build the tool, not the terminal.

A Go framework for terminal applications that **show you state and let you act
on it** — a git client, a file manager, a cluster browser, a system monitor, an
API client, a deploy tool.

Three things here that a widget library does not give you:

**Nothing is invented.** Every one of the 24 components was pulled out of tools
that had already written it. Usually two tools, and usually differently. Each
component's doc comment names those tools and says what their versions
disagreed about.

That could still be a story we tell ourselves, so it is checked against other
people's code. `design/research/rebuilds/` reads the source of ten widely used
TUIs, including lazygit, yazi, k9s, bottom and gitui. For each feature it asks
one question: do we have it, are we missing it, or does it belong to the tool?
Every gap is either filled or written down with its evidence. Three claims we
had made turned out to be false, and those are recorded too.

**The interface cannot lie about itself.** Nine guards read your own source and
fail the build. They catch a key advertised in the help that no handler takes.
A screen nothing can reach. A colour outside the palette. A box-drawing
character typed by hand instead of taken from the theme. These are ordinary Go
tests, and you call them from your own package.

**One declaration, four surfaces.** Describe a command once and get the CLI, the
TUI screen, the context menu and a machine-readable manifest from it. `comp` and
`app` do not depend on `spec`, so a TUI-only tool pays nothing for it.

Not for hosting a text buffer or another terminal. If you are building an editor
or a multiplexer, the things you need most — modal editing, undo history, PTY
management — are not here and are not planned. Inside that line it aims to be
complete: components draw into a cell grid where every cell records **what** drew
it, so a click resolves to the seventh service rather than to row seven, and
still does after the list scrolls.

Built on [Bubble Tea](https://github.com/charmbracelet/bubbletea) and
[Lip Gloss](https://github.com/charmbracelet/lipgloss). Go 1.25.

```sh
git clone https://github.com/richarddavenport/tuikit
cd tuikit && go run ./cmd/tuikit new mytool -dir ..
cd ../mytool && make check   # it builds, runs, and passes its own guards
```

The generated tool resolves tuikit through a `replace` pointing at a sibling
checkout, which is why the clone comes first. That goes when there is a tag.

## What you get

**24 components** — `List` `Viewer` `Pane` `Tabs` `Bar` `Confirm` `StepList`
`LogPane` `Spinner` `Table` `Detail` `Menu` `Toast` `Form` `Split` `Layout`
`Meter` `Input` `Waiting` `Breadcrumb` `Scrollbar` `Keys` `Palette` `Rule`. Run
`tuikit gallery` to use every one of them in every state it has, or
`tuikit gallery -list` to read the inventory as text.

**One declaration, four surfaces.** A `spec.Command` carries its flags, args,
help, key binding and target region. From that, `spec` generates the CLI, the
completion script, `describe --json`, the context menu for whatever region it
acts on, and — via `spec.SchemaOf` — a JSON Schema for tool-calling APIs.

**Three viewports, because they are not the same thing.** `List` has a cursor
and rows. `LogPane` tails a stream and re-attaches when you scroll back to the
bottom. `Viewer` is a document: it opens at the top, scrolls sideways, and puts
the cursor style *underneath* the line's own spans — so the line you are reading
in a diff keeps its syntax colours. Seven of the nine tools in the rebuild
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

- a colour outside the palette
- a character outside the glyph set
- chrome the glyph set cannot print
- an engine that imported a terminal library
- a key binding with no help entry
- a reserved key bound to the wrong meaning
- a component drawing outside its rect

**Screens you can look at without a terminal.** `harness` renders any screen to
a file. `tuikit frames` turns a captured run into an HTML page or into Markdown
with an SVG per frame, for a README or a pull request. Goldens are
colour-stripped, so a layout change is an ordinary test failure.

**Themed by the reader's terminal.** The nine colour roles resolve to ANSI
indices 0–15 — the only colours a terminal lets its user redefine — so a tuikit
tool looks like the rest of that person's terminal rather than like tuikit.

**Pixels where the terminal has them.** On kitty-protocol or Sixel terminals,
components can offer a picture — a gradient meter, a rounded panel — drawn
*behind* the characters. The characters are always the drawing; the picture is a
decoration pass, so a terminal that cannot show one loses nothing but the
gradient.

## What it cannot do

**No text editor, no multiplexer.** Modal editing, undo history, syntax-aware
buffers over large files, PTY hosting and process management are out of scope and
not planned. `Input` is one line; there is no text area yet.

**No reactivity, no virtual DOM.** Redraw a snapshot each frame. `app.Gen`
discards the result of a read the reader walked away from, which is the whole
async problem a tool like this has.

**Coverage is honest, not complete.** `design/research/rebuilds/` works out
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
| [`theme`](docs/theme.md) | nine colour roles, a closed glyph set, the box characters |
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
tuikit designsystem        the colours and glyphs as HTML
tuikit news                what changed since a tool last looked
```

## Reading further

[`docs/`](docs/) is the user guide: one page per package, what it can do and
what it cannot. [`docs/tooling.md`](docs/tooling.md) covers the commands.

`design/decisions.md` is the numbered record of every choice and the reasoning
behind it, including the things that were refused. `design/architecture.md` is
how the packages fit together. `AGENTS.md` is what an agent working in a tuikit
tool needs to know.

## Licence

MIT. See [LICENSE](LICENSE).
