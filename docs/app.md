# `app` — the shell

What holds a tool together between keystrokes: routing, the frame, where you
are, and work that is still in flight.

## The model

```go
type Model interface {
	Init() tea.Cmd
	Update(tea.Msg) (Model, tea.Cmd)
	Draw(c *comp.Canvas, r comp.Rect)
}
```

Bubble Tea's third step is `View() string`. That string is a seam. A tool can
return a canvas frame, or a hand-joined pile of lipgloss, or a canvas frame with
something concatenated on the end. Nothing can tell those apart, and the goldens
will record whatever comes out.

`Draw` closes the seam. It hands you a canvas and a rect, and gives you nowhere
else to put anything.

`Update` returns `app.Model` rather than `tea.Model` for the same reason: a
signature accepting `tea.Model` accepts anything with a `View`.

## `Runner` — owns the frame

`app.New(model, opts...)` adapts a `Model` to Bubble Tea and is the **only**
place a tuikit program calls `comp.NewCanvas`. Size, chrome and the pixel layer
have one right answer per program, and a tool that made them itself would make
them four times.

| | |
| --- | --- |
| `WithSize(w, h)` | a fixed canvas, for tests and captures |
| `WithChrome(ch)` | box characters, gap width, markers |
| `WithPixels(p)` | the graphics layer |
| `WithFullHeight()` | use the whole terminal rather than reserving the bottom row |
| `Canvas()`, `Model()`, `Size()`, `SetSize()`, `Now(t)` | for a harness driving it |

## `Keys` — routing, in the only order that works

```go
app.Keys{Capture: modal, Screen: screen, Global: quit}.Route(msg)
```

Whatever is **capturing** gets the key first, then the screen, then the globals.
A global handled before a modal is a modal you cannot type into — press `q` in a
filter box and the program exits. A capture takes the key whether or not it does
anything with it: an unrecognised key inside a filter box is a character, not a
chance for the screen underneath to act.

The order is not something you write. It is the shape of the struct, and a nil
`Capture` is a tool with nothing capturing rather than a tool that forgot.

## `Mouse` — hits, not coordinates

`Mouse.Route(msg, canvas, handler)` resolves what was hit through
`Canvas.OwnerAt` and dispatches. It carries one piece of state, and that state
is the reason it is a type: **a drag in progress owns the mouse until release**,
whatever it is now over. The pointer outruns the divider it grabbed on every
real drag.

`Handler` is the set of things that can happen: `Press`, `RightPress`, `Wheel`,
`Hover`, `Release`, `Drag` — plus `Drags(id) bool`, which says whether a region
is draggable at all, and `Blocked`, so clicking the frame behind a question is
not an answer to it.
Anything reachable by `RightPress` needs a keyboard path too, because a
multiplexer that keeps right-click for itself gets the event first.

Constants: `WheelRows = 3`. Double-click window defaults to 400ms and its clock
is a field, so a capture script can double-click without real time passing.

## `Stack` and `Screens` — where you are, and how you got there

`Screens` is `map[Screen]View`, which answers only "what does this screen draw".
`Stack` answers the other question. Every tool otherwise writes:

```go
case "esc":
	if m.screen != screenDashboard { m.screen = screenDashboard }
```

— back as a hardcoded constant, correct only because every screen happened to be
reached from the dashboard.

Each entry carries **how you would ask for this screen from a command line** —
`logs api_gateway`. Not to rebuild it from: to say where you are. That makes the
position readable by an agent, printable in a manifest, and reachable by a
capture script. `comp.Breadcrumb` draws it.

It is **not a router**. A screen is a value the tool already has, so parameters
are typed and free and there is nothing to serialise.

## Async

| | |
| --- | --- |
| `Gen` | invalidates results from work you abandoned — a tick from a closed session must not drive its replacement |
| `Poll` | a repeating background read: single-flight, and it **gives up** after `Limit` consecutive failures (default five) |
| `Toggles` | a set of things that are on, keyed by a stable id — open tree buckets, unmasked secrets |

`Gen` is the one worth understanding. The failure it prevents looks like
nothing: a step that finishes after you pressed Esc quietly writes its result
into the screen you left.

## Editing

`Edit(s, msg)` and `EditAt(s, cursor, msg)` apply one keystroke to a string,
returning whether it was consumed. The caret-aware version backs `comp.Input`.

## What it cannot do

- **No router.** No path patterns, no nested outlets, no history serialisation.
- **No focus manager.** Nothing holds which of several panes on one screen is
  focused; `Focused` is a field you set on `List`, `Pane` and `Tabs`. Open issue,
  four tools deep.
- **No key sequences.** `Keys` routes one keystroke; there is no pending prefix,
  so `g` `g` cannot be expressed.
- **It does not own Bubble Tea.** You still call `tea.NewProgram` and choose
  your own options.
