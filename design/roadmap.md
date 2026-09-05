# Roadmap

Where tuikit is going, and what is still undecided. Step 1 is done; the rest is
not commitment.

## Plan

1. ~~**Extract** `theme` + `guard` + `docgen`'s design-system bundle out of
   swarmctl into the module.~~ **Done** — commit `0d0e560`. Verified by running
   the extracted guards against swarmctl's own `internal/tui`: both clean, and a
   role nothing draws with still reported, so they are not passing vacuously.
2. ~~**`examples/democtl`**~~ **Done** — commit `609ee55`. A fictional fleet with
   a dashboard, logs, a confirm modal and a step run; seeded and clock-frozen so
   the harness can use it as a fixture. First real consumer of the guards. Its
   own frames found two bugs before the harness that will automate looking at
   them exists.
3. **`harness`** ([#3](https://github.com/richarddavenport/tuikit/issues/3)) — generalise pgctl's `screenshot_probe_test.go`: fixture and
   live modes, ANSI capture, ANSI→HTML, `guard.Width`, goldens. Prove it against
   democtl, then against swarmctl's existing screens.
4. **`comp.Canvas`** ([#5](https://github.com/richarddavenport/tuikit/issues/5), needs [#4](https://github.com/richarddavenport/tuikit/issues/4) wide runes) — the cell-grid substrate, settled by the prototype on
   `prototype/canvas-mouse`. Wide runes, `Rect` layout, typed owner IDs, and
   democtl ported onto it. Before `comp`, because a component that returns a
   string cannot be clicked. See [mouse.md](mouse.md).
5. **`tuikit watch`** ([#7](https://github.com/richarddavenport/tuikit/issues/7)) — the inner loop. Cheap once `harness` and `docgen` exist,
   and it is what makes everything after this pleasant to build.
6. **`comp` + `app`** ([#8](https://github.com/richarddavenport/tuikit/issues/8), [#9](https://github.com/richarddavenport/tuikit/issues/9), gallery [#10](https://github.com/richarddavenport/tuikit/issues/10)) — extract components, generalising only where two of the
   four already differ meaningfully. **`tuikit gallery`** grows alongside: a
   component that is not in the gallery is not finished.
7. **`spec`** ([#11](https://github.com/richarddavenport/tuikit/issues/11)) — declarations, CLI generation, manifest, completions.
8. **Migrate azctl** ([#12](https://github.com/richarddavenport/tuikit/issues/12)) (3.2k lines, youngest, least to lose). This is the proof;
   a framework that has never met a real tool is a guess.
9. **`tuikit new`** ([#13](https://github.com/richarddavenport/tuikit/issues/13)), seeded from what azctl's migration actually needed.
10. pgctl, dugo, swarmctl migrate later or never. Tracked privately, because it
    is work on those tools rather than on this one — dugo and pgctl are done,
    swarmctl is two steps in.

## Open

Each is an issue, so it gets closed by a decision rather than forgotten.

- Whether `guard.Width` is needed once the canvas clips structurally —
  [#6](https://github.com/richarddavenport/tuikit/issues/6). Leaning on keeping
  it, aimed at `comp` rather than at every tool: the canvas guarantees nothing
  is drawn outside the *canvas*, not that a component stayed inside the *rect*
  it was given.
- Wide runes in the canvas — folded into
  [#5](https://github.com/richarddavenport/tuikit/issues/5), because a canvas
  that is wrong for CJK and then corrected means writing `Set` twice.

## From reading the field

Steps 1–9 are done. The five most-starred Bubble Tea app frameworks were then
read in full ([research/framework-comparison.md](research/framework-comparison.md),
decision 27), and three things came back worth taking:

11. **Delete `View() string`** ([#21](https://github.com/richarddavenport/tuikit/issues/21))
    — bento's `Model` embeds the thing that renders. Ours returns a string,
    which is a seam where a tool can hand back a hand-joined frame and no guard
    would notice. Cheapest of the three, and it hardens what `comp` already is.
12. **A screen stack with history** ([#22](https://github.com/richarddavenport/tuikit/issues/22))
    — `app.Screens` calls itself the router and is a flat map. democtl's `esc`
    is a hardcoded constant; azctl's runner is a second program you cannot
    return from. Take soda's mechanism, add a `spec.Call` label per entry so an
    agent can read where it is, and stop at about a hundred lines.
13. **Constraint layout** ([#23](https://github.com/richarddavenport/tuikit/issues/23))
    — the one axis where the field is plainly ahead of us. Not a Cassowary
    solver; `Length/Min/Max/Percentage/Fill` in one pass covers every layout the
    four tools have.

## Starting a tool today

```sh
make install
tuikit new mytool -short "what it does" -module github.com/you/mytool
cd mytool && make check
```

It came last on purpose, seeded from what azctl's migration actually needed
rather than from a guess — scaffolding written before a real tool had been
migrated would have been a template full of them.
