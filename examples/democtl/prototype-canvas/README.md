# PROTOTYPE — canvas + mouse

**Throwaway.** Lives on the `prototype/canvas-mouse` branch only. It answers one
question and is not meant to be merged.

```sh
go run ./examples/democtl/prototype-canvas
```

Click a row (including the blank space after a short name), click a tab, drag
the divider, scroll the wheel over either pane, right-click a row. The bottom
two lines print the raw mouse event and what the canvas says is under the
pointer, so the hit-testing is visible rather than inferred.

## The question

Should tuikit's `comp` package be canvas-based (components draw cells into a
grid) or string-based (components return strings, the way lipgloss encourages)?

## What it does

`canvas.go` is the whole substrate in 153 lines: a grid of cells, each holding a
rune, a style pointer, and the ID of whatever drew it. `Set` is bounds-checked
into a no-op, so drawing outside the terminal is not an error to catch — it is a
coordinate that does not exist. `OwnerAt(x, y)` is the entire hit-testing
implementation.

`prototype_test.go` drives it with synthetic `tea.MouseMsg` values, finding
targets by owner ID rather than coordinate. If a test can click, so can the
capture harness and so can an agent.

## What scrolling found

The first version had no viewport at all — the list drew from index 0 always —
so the wheel was wired to the *selection* as a stand-in. Scrolling therefore
appeared to pick services at random, and the detail pane could scroll past its
five lines of content and go blank.

Fixed here, and the fix is a requirement for `comp.List` rather than a prototype
detail:

- A viewport offset is separate from the selection. Scrolling is looking around;
  it does not move the cursor.
- Owner IDs carry the **absolute** index, not the screen row. An ID is an
  identity — after scrolling, a click has to select the service that is there.
- Both offsets clamp to actual content, measured from what was drawn last frame.
- Three rows a notch, and a `20/28` indicator, because a viewport with no sign of
  being one looks like a list that lost rows.

## Verdict

See `design/mouse.md` on `main`.
