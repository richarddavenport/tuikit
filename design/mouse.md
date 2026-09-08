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

Branch `prototype/canvas-mouse` — throwaway, not for merge.
It renders democtl's dashboard through a cell grid and wires up all four
behaviors.

```sh
git checkout prototype/canvas-mouse
go run ./examples/democtl/prototype-canvas
```

A `Cell` is a rune, a style pointer, and the ID of whatever drew it. `Set` is
bounds-checked into a no-op. `OwnerAt(x, y)` is the hit test.

### What it cost

| | Lines |
|---|---|
| The whole substrate | **153** |
| democtl's dashboard redrawn on it, with all four mouse behaviors | 369 |
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

**All four behaviors?** All four, and two details that only showed up by
building it:

- Clicking the blank space after a short name has to select the row. It does,
  because the row fills its width before drawing text — the cell decides
  ownership, not the glyph. In a string world that is a manual rect calculation.
- A drag has to survive the pointer outrunning the divider it grabbed, which it
  always does. So a drag in progress owns the mouse until release, whatever it is
  now over. That is four lines at the top of the mouse handler, and it is the
  same shape as the modal check.

### What scrolling it found

Two gaps, both found by Richard using the wheel rather than by any test, and both
requirements for `comp` rather than prototype details.

**The wheel scrolls a viewport; it does not move the cursor.** The list had no
viewport at all — it drew from index 0 always — so the wheel was wired to the
selection as a stand-in, and scrolling appeared to pick services at random. The
targeting was right the whole time: the wheel does act on the pane under the
pointer. What was wrong is what scrolling *meant*. Looking around and choosing
are different operations and a list needs both.

**Scrolling past the end has to be impossible, not discouraged.** The detail
pane clamped to a hardcoded 20 rather than to its content, so wheeling over a
five-line Overview scrolled it into empty space and blanked it.

Three rules for `comp.List` fall out:

- The viewport offset is a separate field from the selection.
- **Owner IDs carry the absolute index, not the screen row.** An ID is an
  identity, so a click after scrolling selects the service that is *there*. This
  is the first place the canvas design could have gone quietly wrong, and it is
  worth stating because "the ID is where it is on screen" is the easy mistake.
- Both offsets clamp against what was actually drawn last frame, not a constant.

A viewport also has to look like one — a count, three rows a notch — or a list
that scrolls looks like a list that lost rows. And the count has to be drawn
**always**, not only when the list overflows: on a tall terminal where everything
fits, a wheel that correctly does nothing is indistinguishable from a wheel that
is broken. `112/112` answers that in a glance; an absent indicator answers
nothing. Found by Richard scrolling a 58-row pane that held the whole list.

And a **selection scrolled out of view has to leave a trace**. Not by dragging
the cursor into the viewport, which is the bug above wearing a different hat, but
by saying which way it went. A selection that is merely invisible is its own
problem: the detail pane goes on describing an item nothing on screen points at,
and the next key press acts on something the reader cannot see. This was the
third thing scrolling found, and the only one that is a design decision rather
than a missing implementation.

**An unexpected result: overflow stops being a class of bug.** `Set` clips to
the canvas, so drawing outside the terminal is not an error to catch — it is a
coordinate that does not exist. A menu deliberately placed at column 88 of 96
renders clipped at the edge with every line still exactly 96 columns.
`guard.Width` becomes a belt-and-braces check on `comp` itself rather than
something every tool has to run. Two of the three bugs the database tool's
capture found, and both of democtl's, were overflow.

## Decisions this settles

**Style lives on the cell, so the canvas never parses ANSI.** Components draw
plain runes plus a style. That deletes the trim-versus-clip problem — the one
the design system's grid page is about — because there is no escape sequence to
miscount. Lip Gloss keeps styling; its layout helpers (`JoinHorizontal`, `Place`)
go.

**Every mouse action needs a keyboard path, and a guard will enforce it.**
`guard.Reachable`. Four reasons now, and the last two decide it: ssh and
keyboard-only users; muscle memory; **an agent cannot click**; and **a
right-click may never arrive at all**.

That fourth one was found by running the prototype inside herdr, which showed its
own context menu instead. A multiplexer that captures the mouse gets the event
first, and the application never learns it happened — there is no protocol for
"did my right-click arrive?". So a tool whose only path to an action is a context
menu is simply broken for anyone inside herdr or tmux, and it cannot detect or
report that.

herdr's own answer is `ui.right_click_passthrough_modifier` in
`~/.config/herdr/config.toml`, empty by default; setting it to `ctrl` makes
Ctrl+right-click reach the pane's app. That fixes the developer's machine. It
does not fix the tool for anyone who has not set it, which is the point.

**The context menu opens from the keyboard too, at the cursor.** Not merely "each
action in the menu also has a binding somewhere" — the *menu itself* is
keyboard-openable, the way a menu key works. Two reasons. It is the only way the
mouse and keyboard paths cannot drift, because they are the same list rather than
a list and a keymap maintained beside it. And it is what makes the whole surface
reachable when right-click is being eaten upstream, which is a condition the tool
cannot detect.

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
- **A multiplexer may keep right-click for itself**, as herdr does by default.
  Nothing the application can do about it, which is why the menu has a keyboard
  opener rather than a note in the README.

One incidental confirmation: herdr's `ui.mouse_scroll_lines` defaults to `3`,
which is the three-rows-a-notch the prototype picked on the grounds that it is
what every other terminal program does. Worth having that checked rather than
assumed.

## What changes in the plan

A new step 4.5, before `comp`: land `comp.Canvas` for real — wide runes, a
`Rect`-based layout helper, typed owner IDs — and port democtl onto it. `comp`
is then extracted on top of a substrate that already knows about the mouse,
rather than being rewritten to learn.


## Hover and double-click

Added after the landscape survey found them missing — `design/research/tui-landscape.md`
calls hover "the real gap", because tview, bubbleapp and cursive all have it and
we handled `MouseActionMotion` only while dragging.

**Hover fires on change, not on movement.** `app.Handler.Hover` is called when
what is under the pointer becomes a different region, and once more with the
zero ID when it leaves everything. A callback per motion event would be a redraw
per motion event, which is the cost that has to be measured rather than assumed.
Reporting the leave matters as much as the enter: a tool that highlighted a row
needs to be told to stop.

Hover is the affordance that says a thing is clickable before you click it, and
it does more work in a terminal than in a window, because there is no cursor
shape to fall back on.

**Double-click is per REGION, not per cell.** A row is one thing however wide it
is, and asking a reader to hit the same cell twice is asking for a skill rather
than a gesture. Two presses on the same owner within 400ms.

Three things it deliberately does:

- `Press` fires for both clicks. The first click of a double-click is a real
  click, and a list that only selected on singles would flicker its selection
  off on the second.
- A double is CONSUMED, so three clicks are a double and a single rather than
  two doubles — which would fire "open" twice for one gesture.
- The clock is a field (`Mouse.Now`), so a capture script can double-click
  without real time passing. Same reason `comp.Spinner` takes a moment rather
  than counting frames.

**Decision 19 applies with force here.** Double-click is the conventional
"open", and an open that exists only for a mouse is an open an agent cannot
perform — `guard.Reachable` will say so. Whatever a double-click does, enter
must do too.

The harness has `hover <region>` and `doubleclick <region>`; the latter sends
two presses and lets `app.Mouse` decide, because a helper that asserted the
answer would be testing itself.
