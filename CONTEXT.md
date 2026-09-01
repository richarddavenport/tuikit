# Context

tuikit's vocabulary. Terms here are the ones the code and the design notes
already use — this file records them, it does not invent them.

When your output names one of these concepts — an issue title, a refactor
proposal, a hypothesis, a test name — use the term as defined here rather than a
synonym. If a concept you need is missing, that is a signal: either you are
introducing language the project does not use, or there is a real gap worth
adding.

## The vocabulary

**Role** — a named colour in the palette. Named by *role*, never by hue:
`Accent` survives someone deciding the interface should be blue; `pink` does not.
A raw ANSI index says what a colour IS instead of what it is FOR.

**Palette** — the closed set of nine roles, valued as the terminal's own sixteen
ANSI indices so a tool is themed by whatever themed the terminal (`Accent`, `Muted`, `Border`,
`Success`, `Pending`, `Danger`, `Stderr`, `SelectionFG`, `SelectionBG`), plus
`Extra` for a tenth meaning the nine do not cover. A tool takes `theme.Default`
and overrides what it wants; it does not build one from scratch.

**Glyph set** — the allow-list of non-ASCII characters an interface may print,
each with a reason. Closed, and a guard holds it closed: a font without a glyph
draws a replacement box, which reads as a bug rather than as decoration.

**Guard** — a test function that holds a rule closed against a package.
`Tokens`, `Glyphs`, `Chrome`, `Engine`, `Reachable` and `Screens` exist;
`Width` is still open (#6). A guard is mechanism: it knows Go source and terminal vocabulary, and
nothing about what a tool does.

**Mechanism / policy** — the line the library-versus-scaffolder split falls on.
Mechanism is a library function that must work for every tool. Policy is a
default, delivered as generated code the tool owns and can edit. Anything
backend-shaped belongs to the tool.

**Canvas** — the cell grid components draw into, instead of returning strings.
A component that returns a string cannot be clicked.

**Cell** — one character position: a rune, a style, and an owner.

**Owner ID** — the identity of whatever drew a cell. An *identity*, never a
screen position: after a list scrolls, a click has to select the item that is
there. This is the first place the canvas design can go quietly wrong.

**Region** — contiguous cells sharing an owner. What a click resolves to, and
what a capture script or an agent names instead of a coordinate.

**Viewport** — an offset into content longer than the pane. A separate field
from the selection, because scrolling is looking around and not choosing.

**Frame** — one rendered screen.

**Capture** — rendering frames in order to look at them. **Fixture mode** is
deterministic and backendless, and is what goldens are made of. **Live mode**
fills the model from a real backend, and is what finds the things a fixture
hides. Capture is a building tool first; documentation and goldens are
downstream of that.

**Golden** — a stored frame, diffed in review, so a layout change is an ordinary
test failure rather than something someone has to notice.

**Surface** — one output of a `spec.Command` declaration. There are four: the
CLI command, the TUI screen, the `describe --json` entry an agent reads, and the
right-click menu. One declaration, so none of them can drift from the others.

**Engine / UI split** — the engine knows the domain and has **no terminal
concepts at all**: no colour, no width, no keys, no framework, and no tuikit
import. In practice an engine imports stdlib plus its own domain SDK. The UI
never calls the domain directly; the CLI is a peer of the TUI over the same
engine. Every tool in the family keeps it (decision 22).

**Screen** — one full-window view. Every screen constant needs a `View()` case;
a screen without one renders as an empty terminal and says nothing about why.

## Flag contradictions

If your output contradicts a numbered entry in `design/decisions.md`, say so
rather than silently overriding it:

> _Contradicts decision 18 (components draw cells, not strings) — but worth
> reopening because…_
