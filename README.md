# tuikit

A Go framework for **operator tools** — the kind with a command line *and* a
terminal interface over the same domain logic. `kubectl` and `k9s` as one
binary, for your own infrastructure.

It is not a general TUI framework. It is a narrow shape, and the narrowness is
the point: declare a command once and get a CLI, a TUI screen, a context-menu
entry and a machine-readable manifest from the same declaration; draw into a
cell grid where every cell records what drew it, so clicks resolve without
maintaining a hit-test map; and hold the whole interface inside a vocabulary
that a test can enforce.

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

**It is not a general TUI framework.** No reactivity, no virtual DOM, no
component marketplace. If you are building a text editor or a game, use
[tview](https://github.com/rivo/tview), [tcell](https://github.com/gdamore/tcell)
or [ratatui](https://ratatui.rs).

**Layout is bands and splits, not a constraint solver.** `comp.Layout` divides
a rect into `Length`, `Percent` and `Fill` bands down one axis, each with an
optional `.Min()` and `.Max()`. ratatui and
[bento](https://github.com/metafates/bento) have a Cassowary solver and are
plainly better at this.

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
