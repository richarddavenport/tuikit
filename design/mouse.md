# Mouse

First-class mouse support: click to focus and select, wheel scrolling, drag to
resize panes, right-click menus. Prompted by [herdr], which does all four and
puts it well — "tmux-style prefix keys *and* click, drag, split; pick per
moment, not per tool."

[herdr]: https://github.com/herdrdev/herdr

## Why this is an architecture question, not a feature

**A component that returns a string cannot be clicked.** Hit-testing needs
geometry, and `lipgloss.JoinHorizontal(...)` throws geometry away. By the time
`View()` returns a string, nothing knows that row 4, columns 2–31 was
`services[2]`.

Rust TUI libraries are cell-buffer based — you get a rect, you draw cells into a
buffer, and the rect is still there to test a click against. Bubble Tea and Lip
Gloss are string-composition based. That difference is most of why mouse support
feels native in one and bolted-on in the other, and it is why this had to be
settled before `comp` rather than after.

Bubble Tea's input side is complete: `tea.MouseMsg` carries `X`, `Y`, `Action`
(press/release/motion), `Button` (left/right/middle/wheel) and modifiers, and
`WithMouseCellMotion()` delivers drags. Nothing is missing at the terminal.
What was missing is *where the click lands*.

## The prototype

Branch `prototype/canvas-mouse`, commit `1cdd3f8` — throwaway, not for merge.
It renders democtl's dashboard through a cell grid and wires up all four
behaviours.

```sh
git checkout prototype/canvas-mouse
go run ./examples/democtl/prototype-canvas
```

A `Cell` is a rune, a style pointer, and the ID of whatever drew it. `Set` is
bounds-checked into a no-op. `OwnerAt(x, y)` is the hit test.

### What it cost

| | Lines |
|---|---|
| The whole substrate (`canvas.go`) | **153** |
| democtl's dashboard redrawn on it, with all four mouse behaviours | 369 |
| Tests driving it by synthetic mouse events | 201 |

Against that, from the string-based version it deletes `clip`, `pad`,
`padVisible`, `rule`, and `composite` outright, and the 28-line `overlay` becomes
an 8-line `drawMenu` with no compositing step at all.

### What it answered

**Does it kill the hand-rolled overlay?** Yes, completely. A menu is drawn last,
so it is on top. There is no compositing, no re-measuring of the lines beneath,
and nothing to get wrong. The bug that wiped 8 of 18 framed rows in democtl is
not a bug that can be written here.

**Does hit-testing fall out for free?** Yes — `OwnerAt` is four lines, and there
is no region list that can drift from the drawing, because the drawing *is* the
region list. A component that moved but forgot to update its rect is not a
failure mode that exists.

**All four behaviours?** All four, and two details that only showed up by
building it:

- Clicking the blank space after a short name has to select the row. It does,
  because the row fills its width before drawing text — the cell decides
  ownership, not the glyph. In a string world that is a manual rect calculation.
- A drag has to survive the pointer outrunning the divider it grabbed, which it
  always does. So a drag in progress owns the mouse until release, whatever it is
  now over. That is four lines at the top of the mouse handler, and it is the
  same shape as the modal check.

**An unexpected result: overflow stops being a class of bug.** `Set` clips to the
canvas, so drawing outside the terminal is not an error to catch — it is a
coordinate that does not exist. A menu deliberately placed at column 88 of 96
renders clipped at the edge with every line still exactly 96 columns. `guard.Width`
becomes a belt-and-braces check on `comp` itself rather than something every tool
has to run. Two of the three bugs pgctl's capture found, and both of democtl's,
were overflow.

## Decisions this settles

**Style lives on the cell, so the canvas never parses ANSI.** Components draw
plain runes plus a style. That deletes the trim-versus-clip problem — the one
the design system's grid page is about — because there is no escape sequence to
miscount. Lip Gloss keeps styling; its layout helpers (`JoinHorizontal`, `Place`)
go.

**Every mouse action needs a keyboard path, and a guard will enforce it.**
`guard.Reachable`. Three reasons, and the third decides it: ssh and
keyboard-only users, muscle memory, and **an agent cannot click**. In a framework
whose thesis is that agents get on with it immediately, a right-click-only action
is one agents cannot reach.

**Regions get names, and the harness drives them by name.** The prototype's tests
already find targets by owner ID rather than coordinate, which is how a capture
script should read:

```
democtl tui --snapshot out/ --script 'click list.row[2]; drag split.main +10; rclick list.row[2]'
```

An agent clicking by name is strictly better than an agent guessing coordinates,
and it only works because ownership is recorded at draw time.

**Right-click menus are a fourth surface of `spec.Command`.** The menu for a
region is the commands whose target matches that region's type, built from the
same declaration that produces the CLI command, the keybinding and the
`describe --json` entry. A menu that cannot drift from the CLI is worth more
than a hand-maintained one.

## Costs and caveats, stated plainly

- **Lip Gloss's layout helpers are out.** Styling stays. This is a real
  departure from the idiom most Bubble Tea code is written in, and anyone
  reading tuikit's components will notice.
- **Wide runes are not handled yet.** The prototype assumes one rune, one cell.
  CJK and emoji need a width-aware `Set` that claims two cells and a continuation
  marker. Known, not hard, not done.
- **Mouse capture breaks native text selection.** Shift usually bypasses it, and
  a toggle key is conventional. This is a documentation problem, not one to
  solve.
- **tmux needs `mouse on`.** Also documentation.

## What changes in the plan

A new step 4.5, before `comp`: land `comp.Canvas` for real — wide runes, a
`Rect`-based layout helper, typed owner IDs — and port democtl onto it. `comp`
is then extracted on top of a substrate that already knows about the mouse,
rather than being rewritten to learn.
