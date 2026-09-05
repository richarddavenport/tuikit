# Rebuilding lazygit

[jesseduffield/lazygit](https://github.com/jesseduffield/lazygit) — 81,999
stars, Go, on gocui. The most used TUI in the surveyed field, and more starred
than Bubble Tea itself.

Read from the repository on 2026-09-04: `pkg/gui/` and its subpackages, with
file sizes quoted so the claims below can be checked.

## The shape

Four panels down the left — status, files, branches, commits, stash — and a
main panel on the right showing whatever the cursor is on. A command palette,
a menu, confirmations, a diff view, and a lot of git.

## What tuikit already supplies

| lazygit has | comp gives |
| --- | --- |
| the left column of panels | `Layout.Rows` + `Pane`, one `List` each |
| the main/secondary split | `Split`, nested for the two-panel diff |
| its side panel lists | `List`, with `Row.Skip` for headings |
| `pkg/gui/popup` | `Confirm`, `Toast` |
| its menu | `Menu`, from the `spec` declaration |
| keybinding help | `Keys` + `guard.Keys` |
| its status line | `Bar` |
| search / filter | `Input` + `fuzzy` + `Highlight` |
| `pkg/gui/style` | `theme` — and lazygit's is a hand-rolled equivalent |

That is most of the chrome, and it is the part tuikit is already right about.

## Holes

### 1. Tree — `pkg/gui/filetree/`, eleven files

`build_tree.go` (4.3 kB), `node.go` (8.4 kB), `collapsed_paths.go`,
`file_tree_view_model.go`, and a `README.md` explaining it. A file tree with
collapsing, a flat/tree toggle, filtering, and two node types (working-tree
files and commit files) over one traversal.

**This is the biggest single gap in the field**, not just here: six of the
fourteen surveyed tools need one. `comp.Tree` was refused once, on the grounds
that azctl's and swarmctl's versions carried different payloads and a shared
node type would make both worse. That reasoning was sound against a sample of
two. It does not survive lazygit, yazi, superfile, dive, fx, termshark and
ranger each building the same thing.

The refusal also points at the shape: what differs is the **payload**, so the
component must not own it. `List` already solved this exact problem —
`Row.Depth` is a number rather than a node, and the flattening stays the
tool's. A Tree is likely `Row.Depth` plus collapse state plus "which rows does
collapsing hide", not a node type.

### 2. Line-range selection in a diff — `pkg/gui/patch_exploring/`

`state.go` is 13 kB and `focus.go` handles keeping the selection visible. This
is staging by hunk and by line: a range that grows with shift-arrow, snaps to
hunk boundaries, and survives the diff being re-rendered under it.

`comp.List` has a cursor and no selection. **A range selection over rows is a
real gap** and it is not only lazygit's — any tool that stages, bulk-selects or
copies a span needs it, and it is the interaction pgctl asked for in issue 49's
multi-select from the other direction.

Worth noting the size: 13 kB of state for what sounds like "shift-click a
range". Most of it is the snapping and the survival across re-render.

### 3. A text view that is not a log — the main panel

`LogPane` tails a stream and follows the end. A diff is not that: it scrolls
from the top, wraps or does not, highlights syntax, and has a selection over
it. Five of the fourteen tools need this shape.

## Theirs, and rightly

- **`pkg/gui/presentation/graph/`** (10 kB) — the commit graph, the `│ ├ ─ ╯`
  lines beside the log. A framework supplying this would be a framework with an
  opinion about git.
- **`pkg/gui/presentation/`** — sixteen files turning branches, commits,
  stashes and submodules into rows. This is exactly the layer tuikit says
  belongs to the tool, and lazygit agrees by putting it in its own package.
- **`pkg/gui/mergeconflicts/`** — finding and rendering conflict markers.
- **`pkg/gui/context/`, `controllers/`** — its own routing, which `app.Keys`
  and `app.Stack` cover differently rather than better.

## The verdict

**Could tuikit rebuild lazygit today? No — three components short.** Tree, a
range selection, and a text view with syntax.

**Is that a criticism of tuikit's scope?** Only partly. lazygit is a git client,
not an operator tool, and a diff view with hunk staging is git-shaped work. But
the three holes are all general: a tree is not a git idea, a range selection is
not a git idea, and neither is a scrollable syntax-aware view.

**What it does not need** is the thing worth noticing. lazygit has no charts, no
forms, no wizard, no tabs. Its 82k stars come from four lists, a split, and a
very good diff. The components that matter are few and deep, which is an
argument for `comp` staying small and each entry being thorough — not for a
catalogue.
