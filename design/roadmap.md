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
10. pgctl, dugo, swarmctl migrate later or never ([#14](https://github.com/richarddavenport/tuikit/issues/14)).

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

## Starting a tool today

`tuikit new` is step 9, so the supported path until then is copying
`examples/democtl` — [#1](https://github.com/richarddavenport/tuikit/issues/1)
documents it. That is deliberate: scaffolding before azctl's migration has shown
what a tool actually needs produces a template full of guesses.
