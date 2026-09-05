# Rebuilding gcpeasy

[scttymn/gcpeasy](https://github.com/scttymn/gcpeasy) — Go, on Bubble Tea 1.3.10,
Bubbles, **lipgloss v2** and cobra. A GCP operator tool: pick a project, a GKE
cluster and a pod, then run things against them.

Read from a shallow clone on 2026-09-05.

## Why this one is not like the others

Every other rebuild here is a famous tool on a substrate we do not use. This is
a small one on **exactly ours** — Bubble Tea, lipgloss, cobra, a Go domain layer
shelling out to a cloud CLI. It is what azctl, pgctl and swarmctl looked like
before tuikit, written by somebody who has never heard of it.

That makes it the most direct evidence in the survey, and the least flattering
to read.

## The shape, in one number

**`cmd/tui.go` is 78,596 bytes. 2,972 lines. 154 functions. One file.**

The whole program is 5,554 lines of Go. The TUI is 54% of it, and the domain
layer — everything that actually talks to GCP — is 343 lines across two files.

For scale: gh-dash's `ui.go` was 54 kB and the k9s rebuild called it the largest
single UI file in the Go half of this survey. This is half again bigger, in a
tool with a fraction of the features.

## What tuikit already supplies

Named against the actual functions, so the claim can be checked.

| gcpeasy wrote | comp gives |
| --- | --- |
| `renderPanel` — border, title, focus colour, inner clip | `Pane` |
| `renderRow` — cursor glyph, marker, primary, secondary, padding | `List` + `Row.Lead`, `Row.Right` |
| four `render*Panel` funcs, one per pane | one `List` each |
| `renderLeft`, `panelInnerWidth`, `maxInt`, `minInt` | `Layout.Rows`, `Rect` |
| `renderFooter`, `padFooter`, `footerHint`, `footerHints` | `Bar` + `comp.Hint` |
| `renderCommandModal`, `commandItems` | `Palette` |
| `renderHelpModal`, `helpRow` | `Keys` |
| `renderAuthDialog`, `renderAuthButtons`, `renderRefreshModal` | `Confirm` |
| `overlayCentered` | drawing into a rect from `comp.Center` |
| `clip` | `Truncate` |
| `renderBootScreen`, `renderGradientLogo`, `logoGradientColor`, `mixRGB`, `tuiRGB.hex` | `paint.Ramp` + `comp.Picture` |
| `visibleProjects`, `visibleClusters`, `visiblePods` | `Row.Skip`, or a filter over `DrawFunc` |
| `taskSpec`, `runningTask`, `taskStartedMsg`, `taskOutputMsg`, `taskDoneMsg` | `app.Async`, `app.Gen`, `app.Poll` |
| `handleKey` / `handleAuthDialogKey` / `handleCommandModalKey` / `handleHelpModalKey` | `app.Keys` — capture, screen, global |
| `tuiStateCache`, `tuiPreferences` | the tool's, but `app.Toggles` is the hidden-item half |
| its cobra tree | `spec.Command`, plus the menu and manifest for free |

## Four tests that exist because of strings

This is the part worth sitting with. `cmd/tui_test.go` is 30.9 kB, and among
what it asserts:

- `TestLeftAndRightPanelsMatchHeight`
- `TestPanelsRenderAtRequestedWidth`
- `TestSelectedRowsFillPanelInterior`
- `TestLongPodRowsDoNotWrap`

Every one is a test that two rendered **strings** have compatible shapes. On a
cell grid none of them can be written, because none of the failures can happen:
a pane is a rect, a row is filled to its rect's width, and drawing past the
right edge is clipped rather than wrapped.

They are good tests. They are also four tests, in a 31 kB file, defending
properties that a different substrate gives away.

## What the guards would say

Not hypothetical — these are the specific findings.

**`guard.Engine` fires, and the cost is already visible.** `internal/` imports
`bufio`, `fmt` and `os`, and calls `fmt.Print` **21 times**. `SelectCluster`
(`kubernetes.go:47`) and `SelectPod` (`pod.go:113`) *prompt the user on stdin*
from inside the domain layer.

The consequence is not stylistic. **The TUI cannot call them.** A function that
prints to stdout would draw over the frame and then block on a read that will
never come, so `SetupClusterAndSelectPod` is reachable only from `cmd/pod.go`,
`cmd/rails.go` and `cmd/cluster.go` — the CLI paths — and the TUI reimplements
selection itself. One decision in the engine, taken for the CLI's convenience,
cost a second implementation of the tool's central interaction.

That is `guard.Engine`'s whole argument, demonstrated by someone who was not
trying to make it.

**`guard.Tokens` fires, and the instinct behind it was right.** There are no hex
literals anywhere. Instead there are 15 named style vars — `tuiTitleStyle`,
`tuiMutedStyle`, `tuiSelectedRowStyle`, `tuiActiveBorderColor` — built from ANSI
256 numerics (`"81"`, `"229"`, `"62"`, `"149"`). Somebody centralised their
palette on their own, which is the same move `theme.Palette` is.

What they do not get is a set that is *closed*: nothing stops the sixteenth
style, nothing reports a role nothing draws with, and nothing distinguishes the
0–15 band the reader themes from the 16–255 band they do not (decision 28).

**No mouse at all.** Zero occurrences of `MouseMsg`. Not a criticism — it is
work, and `app.Mouse` plus owner IDs is the thing that makes it not work.

## Holes

### 1. ANSI from a subprocess, turned into spans

`appendOutput` (`tui.go:2140`) takes the output of `gcloud`, `kubectl` or
`rails console` and has to cope with what a real program emits: `\x1b` CSI
sequences via `parseTerminalEscape` and `applyCSI`, `\r` moving the column back,
`\b` and `\x7f`, tabs to a stop of four, and a 2,000-line cap.

`comp.LogPane` and `comp.Viewer` both take plain text or `Segment`s. **Nothing
in tuikit turns a subprocess's coloured output into `Segment`s** — and every
operator tool shells out to something.

The logic exists, on the wrong side of the fence: `harness.Strip` and
`harness.Rows` parse exactly these sequences, but they live in the *test*
package and are about reading a frame tuikit itself wrote.

A bounded piece of work: SGR to `[]Segment`, `\r` and `\b` applied, everything
else dropped. Not an emulator. See the issue.

### 2. Panel focus, for the fifth time

`m.focus tuiPanel` plus `m.cursors[panel]` — focus, and a cursor per pane. With
lazygit, termshark, dive and swarmctl that is five, and this one is on our own
substrate. Issue 59.

## Outside the line

**The interactive session pane is a terminal, and we said we do not do those.**

`tuiInteractiveCommand` runs a PTY (`creack/pty`) so `rails console` and
`kubectl exec` work with real line editing and scrollback. `applyCSI`,
`writeOutputRune`, `newOutputLine`, `outputRow`/`outputCol` are a small terminal
emulator, and the file's own comment says why: *"This is what makes copy/paste,
scrollback, colors, and line editing behave"*.

Decision 27 names hosting another terminal as out of scope, and this is the
first rebuild where a real chunk of the target lands on the far side of that
line. The honest verdict is not "tuikit could build gcpeasy" — it is that
tuikit could build **all of gcpeasy except the interactive pane**, and that pane
is a deliberate exclusion rather than a gap.

Worth noticing that it is genuinely separable: four panels, a footer, modals and
a boot screen on one side; a PTY in a box on the other.

## Theirs — the domain, not the shape

- `internal/kubernetes.go`, `internal/pod.go` — GKE, kubectl, namespaces.
- `cmd/auth.go`, `cmd/env.go`, `cmd/rails.go` — gcloud auth, environments, Rails.

## The finding that is about us, not them

**lipgloss v2 ships a cell buffer.** gcpeasy's `overlayCentered` uses it:

```go
canvas := lipgloss.NewCanvas(width, height)
composite := lipgloss.NewCompositor(
    lipgloss.NewLayer(base).X(0).Y(0).Z(0),
    lipgloss.NewLayer(modal).X(x).Y(y).Z(z),
)
return canvas.Compose(composite).Render()
```

`lipgloss.Canvas` has `CellAt`, `SetCell`, `Compose`, `Render`, `Resize`. Under
it, `ultraviolet.Cell` is `{Content, Style, Link, Width}`.

So **a cell grid is no longer what makes tuikit different**, and any claim in
that shape should be retired. What `uv.Cell` does not carry is an **owner**.
There is no `OwnerAt`, so a click still has to be resolved against remembered
rectangles rather than against the frame that was actually drawn — which is the
property `comp.Canvas` exists for and the one the README should lead on.

Checked because `other-tuis.md` exists to make claims of the form "nobody else
does X" checkable before they are made. This one needed narrowing.

## The verdict

**Could tuikit rebuild gcpeasy today? Everything except the PTY pane, and the
78 kB would not survive it.**

The rest of the survey asks whether `comp` is complete. This one asks the other
question — what a tool on our own substrate pays for not having it — and answers
it in a single number: 2,972 lines in one file, 54% of the program, defended by
four tests for properties a cell grid gives free, with its central interaction
implemented twice because the engine prints.
