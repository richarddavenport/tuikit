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
3. **`harness`** — generalise pgctl's `screenshot_probe_test.go`: fixture and
   live modes, ANSI capture, ANSI→HTML, `guard.Width`, goldens. Prove it against
   democtl, then against swarmctl's existing screens.
4. **`comp.Canvas`** — the cell-grid substrate, settled by the prototype on
   `prototype/canvas-mouse`. Wide runes, `Rect` layout, typed owner IDs, and
   democtl ported onto it. Before `comp`, because a component that returns a
   string cannot be clicked. See [mouse.md](mouse.md).
5. **`tuikit watch`** — the inner loop. Cheap once `harness` and `docgen` exist,
   and it is what makes everything after this pleasant to build.
6. **`comp` + `app`** — extract components, generalising only where two of the
   four already differ meaningfully. **`tuikit gallery`** grows alongside: a
   component that is not in the gallery is not finished.
7. **`spec`** — declarations, CLI generation, manifest, completions.
8. **Migrate azctl** (3.2k lines, youngest, least to lose). This is the proof;
   a framework that has never met a real tool is a guess.
9. **`tuikit new`**, seeded from what azctl's migration actually needed.
10. pgctl, dugo, swarmctl migrate later or never.

## Open

- Whether `engine` gets any framework support at all, or stays entirely the
  tool's own. Current lean: entirely the tool's own — the split works *because*
  the engine has no UI imports.
- Whether `guard.Width` is needed at all once the canvas clips structurally, or
  stays as a check on `comp` itself. Leaning on keeping it, cheaply, aimed at
  `comp` rather than at every tool.
- Wide runes in the canvas: one rune is not one cell for CJK or emoji. A
  width-aware `Set` claiming two cells with a continuation marker is the answer;
  it is not written.
