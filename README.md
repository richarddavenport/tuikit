# tuikit

A Go framework for terminal applications that **show you state and let you act
on it** — a git client, a file manager, a cluster browser, a system monitor, an
API client, a deploy tool.

Not for hosting a text buffer or another terminal. If you are building an editor
or a multiplexer, the things you need most — modal editing, undo history, PTY
management — are not here and are not planned.

Inside that line it aims to be complete. Draw into a cell grid where every cell
records what drew it, so clicks resolve without a hit-test map. Hold the whole
interface inside a vocabulary a test can enforce. And, if the tool has a command
line, declare each command once and get the CLI, the TUI screen, the context
menu and a machine-readable manifest from the same declaration — `comp` and
`app` do not depend on `spec`, so a TUI-only tool pays nothing for it.

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

**23 components** — `List` `Pane` `Tabs` `Bar` `Confirm` `StepList` `LogPane`
`Spinner` `Table` `Detail` `Menu` `Toast` `Form` `Split` `Layout` `Meter`
`Input` `Waiting` `Breadcrumb` `Scrollbar` `Keys` `Palette` `Rule`. Run
`tuikit gallery` to use every one of them in every state it has, or
`tuikit gallery -list` to read the inventory as text.

**One declaration, four surfaces.** A `spec.Command` carries its flags, args,
help, key binding and target region. From that, `spec` generates the CLI, the
completion script, `describe --json`, the context menu for whatever region it
acts on, and — via `spec.SchemaOf` — a JSON Schema for tool-calling APIs.

**Clicks that resolve to identity, not coordinates.** `comp.Canvas` is a grid
of cells; each records the `ID` of what drew it. `OwnerAt(x, y)` is the entire
hit test. An ID names *the item*, not the row it landed on, so a list that
scrolls does not act on the wrong thing.

**Panes compose.** `comp.Layout` divides a rect into any number of bands down
one axis; `comp.Split` gives two of them a divider you can drag, with a `Min` so
neither can be dragged to nothing. Splits nest, so three panes with two
independent dividers is a composition rather than a component.

**Guards that fail the build.** A colour outside the palette, a character
outside the glyph set, chrome the glyph set cannot print, an engine that
imported a terminal library, a key binding with no help entry, a reserved key
bound to the wrong meaning, a component drawing outside its rect — each is a
test failure rather than a review comment.

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

| | |
| --- | --- |
| `comp` | the cell canvas and the components |
| `app` | the shell: key routing, mouse, screen stack, async generations |
| `spec` | one command declaration, and the surfaces it produces |
| `theme` | nine colour roles, a closed glyph set, the box characters |
| `guard` | the tests that hold all of the above closed |
| `harness` | drive a model, capture frames, compare goldens |
| `term` | what this terminal can do — graphics, cell size, palette |
| `paint` | gradients, rounded panels and bars, as images |
| `fuzzy` | ranked matching that reports *where* it matched |
| `docgen` | frames and the design system, as HTML or Markdown |
| `scaffold` | `tuikit new` |

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

`design/decisions.md` is the numbered record of every choice and the reasoning
behind it, including the things that were refused. `design/architecture.md` is
how the packages fit together. `AGENTS.md` is what an agent working in a tuikit
tool needs to know.

## Licence

MIT. See [LICENSE](LICENSE).
