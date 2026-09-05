# Other TUIs

A running list of what else exists, kept so a claim of the form "nobody does X"
has somewhere to be checked before it is made. Two have already failed that
check — see the note at the bottom.

This is a **reference list, not a review**. Entries say what a thing is, not
whether it is good. `tui-landscape.md` beside this file is the measured survey
of the Go and Rust halves; everything else here is unverified.

## How much to trust each row

**Measured** — Go and Rust, pulled 2026-09-01 into `substrates.tsv`: 903 repos
across fourteen GitHub queries, with the substrate read from each top repo's
`go.mod` or `Cargo.toml`. Star counts are real and dated.

**Unverified** — every other language. Written from memory, with no pull behind
it and deliberately no numbers. **Do not cite a figure from this section**;
check it first. That is not caution for its own sake — issue 52 was somebody
reading a colour off tuikit's own generated page and repeating it as fact, and
an unverified list is the same trap with a longer fuse.

Dates are when a row was last checked, not when the project last moved.

## Go — measured, 2026-09-01

| Stars | | Substrate |
| ---: | --- | --- |
| 81,858 | lazygit | tcell, directly |
| 54,522 | dive | gocui + tcell |
| 44,740 | charmbracelet/bubbletea | — |
| 23,030 | superfile | bubbletea |
| 20,616 | fx | bubbletea |
| 17,070 | wtf | bubbletea + tcell + tview |
| 13,964 | rivo/tview | tcell |
| 12,448 | gh-dash | bubbletea |
| 10,596 | jroimartin/gocui | — |
| 9,984 | termshark | tcell |

Of the top 80 TUI repos, 26 are Go apps: 13 bubbletea, 11 tcell (5 of those via
tview), 3 termui/termdash/gocui.

**The two most popular Go TUIs are not Bubble Tea apps.** lazygit uses tcell
directly and has more stars than Bubble Tea itself.

## Rust — measured, 2026-09-01

| Stars | | Substrate |
| ---: | --- | --- |
| 41,856 | yazi | ratatui |
| 40,888 | CodeWhale | ratatui |
| 34,335 | herdr | ratatui |
| 26,329 | grok-build | ratatui |
| 22,468 | ratatui/ratatui | — |
| 22,453 | gitui | ratatui |
| 19,343 | spotify-tui | tui-rs, ratatui's predecessor |
| 13,964 | bottom | ratatui |
| 10,521 | oha | ratatui |
| 7,763 | trippy | ratatui |

41 of the top-80 apps are Rust, and **31 of those are ratatui**. 47.6M
crates.io downloads, 17.2M in the last 90 days.

**Rust is the consolidated half.** One answer where Go has four.

## Python — unverified

| | |
| --- | --- |
| **Textual** | The most active TUI framework outside Rust. CSS-like stylesheets, a real layout engine, a devtools console, async-first, and it renders the same app to a browser. |
| **Rich** | Textual's rendering sibling; tables, markup, progress. Very widely used on its own. |
| **prompt_toolkit** | What IPython is built on. Line editing, completion, full-screen apps. |
| **urwid** | The older widget toolkit. |

Apps: **posting** (API client), **harlequin** (SQL IDE), **toolong**, **ranger**
(file manager, predates Textual).

Textual is the one to check any framework-level claim against. It has done more
than anything else in this list on styling, layout and developer tooling.

## C and C++ — unverified

| | |
| --- | --- |
| **ncurses** | The substrate almost everything ultimately descends from. |
| **notcurses** | The ambitious modern C one — blitting, multimedia, direct mode. |
| **FTXUI** | Modern C++, functional-flavoured components. |
| **imtui** | Immediate-mode, Dear ImGui's model in a terminal. |

Apps: **htop**, **tmux**, **Midnight Commander**, **ncmpcpp**, **cmus** in C;
**btop** in C++. These are among the most-used TUIs in the world and they
predate every framework above.

## Everything else — unverified

| | |
| --- | --- |
| **Spectre.Console** (C#) | Rich rendering **and** `Spectre.Console.Cli`, which takes a command declaration and produces parsing, help and exit codes. |
| **Terminal.Gui** (C#) | Mature widget toolkit, originally Miguel de Icaza's gui.cs. |
| **Ink** (JS/TS) | React for command-line apps; a lot of npm tooling uses it. |
| **blessed** (JS) | The older one. |
| **brick** (Haskell) | Well regarded; declarative, `vty` underneath. |
| **TTY toolkit** (Ruby) | A suite of small gems covering prompts, tables, boxes, commands. |
| **Lanterna** (Java) | Long-lived; a curses-like layer plus widgets. |
| **libvaxis** (Zig) | Newer, modern protocol support. |
| **Ratatouille** (Elixir) | Built on the Elm architecture. |
| **illwill** (Nim) | Small and deliberately so. |

## Claims that failed this check

Kept because a list of what was wrong is more useful than a list of what is
right.

**"The one axis where somebody else is plainly better than us"** — filed about
constraint layout, in decision 27, before `comp.Layout` existed. Issue 23 built
it and found a linear pass reaches the same answer for the shapes these tools
lay out. The prediction outlived the work that answered it.

**"The CLI is the axis nobody else is on"** — decision 27 again. Spectre.Console
produces a CLI and rich rendering from one declaration, and TTY covers similar
ground. Two surfaces from one declaration is not novel; the narrower claim,
four surfaces plus an enforced engine boundary, is what the evidence supports.

Both survived a while because nobody had looked outside the pool the original
survey drew from — one GitHub query, in one language.

**"A cell grid where every cell records what drew it"**, stated as though the
grid were the distinctive half. lipgloss v2 ships `lipgloss.Canvas` with
`CellAt`, `SetCell` and `Compose`, over an `ultraviolet.Cell` of
`{Content, Style, Link, Width}`. The cell grid is now table stakes. What no
other cell buffer carries is the **owner**, so that is the claim, and it is
narrower than the one being made. Found by the gcpeasy rebuild, 2026-09-05.

**"gh-dash's `ui.go` is 54 kB of exactly the state management `app` was
extracted to remove."** Measured afterwards: 1,917 lines, of which `View` and
its helpers are 275 and a single `Update` function is 741. The claim was right
about the *kind* of code and wrong to attribute the whole file to it. The
sharper version is that one function is 39% of the file.

**"Its 82k stars come from four lists, a split, and a very good diff"**, and the
same move about dive's 54k. A star count has many causes and a UI inventory is
not a measurement of any of them. Both rewritten to describe the interface
without claiming to explain the popularity.

**"gcpeasy's 2,972-line TUI is what a tool pays for not using tuikit."** Written
in the first draft of `rebuilds/gcpeasy.md`. Measured afterwards: 624 of those
lines are drawing, and 1,715 are the program itself. `comp` would replace the
624. The larger number was true and the implication was not.
