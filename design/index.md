# Design notes

Rationale, rejected options and plans — for people changing tuikit, not for
people using it.

- [decisions.md](decisions.md) — dated record of choices, including what was
  rejected and why
- [architecture.md](architecture.md) — how the packages fit together, and where
  the line between library and scaffolder falls
- [capture.md](capture.md) — how someone building a tool sees what they are
  building, and how an agent sees it too
- [mouse.md](mouse.md) — first-class mouse support, and the cell-grid substrate
  the prototype settled on
- [roadmap.md](roadmap.md) — what is built, what is next, what is undecided
- [research/tui-landscape.md](research/tui-landscape.md) — the measured survey:
  903 Go and Rust repos, what they are built on, and where tuikit is not special
- [research/other-tuis.md](research/other-tuis.md) — everything outside Go and
  Rust, marked by how much of it has been checked. Read it before writing
  "nobody does X"; two claims have already failed there
- [research/rebuilds/](research/rebuilds/) — could tuikit rebuild the TUIs
  people actually use? One file per tool, read from the source, listing what is
  missing. This is where a component earns its way in

The design system bundle is generated rather than written: `make designsystem`
renders `theme` as HTML. There is no hand-maintained palette page here, and there
should never be one — see decision 7.
- [keys.md](keys.md) — the keyboard meanings that turned out to be worth
  writing down, and the one that is still open
