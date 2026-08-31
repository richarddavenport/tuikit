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

The design system bundle is generated rather than written: `make designsystem`
renders `theme` as HTML. There is no hand-maintained palette page here, and there
should never be one — see decision 7.
