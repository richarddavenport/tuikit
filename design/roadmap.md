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
4. **`tuikit watch`** — the inner loop. Cheap once `harness` and `docgen` exist,
   and it is what makes everything after this pleasant to build.
5. **`comp` + `app`** — extract components, generalising only where two of the
   four already differ meaningfully. **`tuikit gallery`** grows alongside: a
   component that is not in the gallery is not finished.
6. **`spec`** — declarations, CLI generation, manifest, completions.
7. **Migrate azctl** (3.2k lines, youngest, least to lose). This is the proof;
   a framework that has never met a real tool is a guess.
8. **`tuikit new`**, seeded from what azctl's migration actually needed.
9. pgctl, dugo, swarmctl migrate later or never.

## Open

- Whether `engine` gets any framework support at all, or stays entirely the
  tool's own. Current lean: entirely the tool's own — the split works *because*
  the engine has no UI imports.
- Whether `guard.Width` can be mechanism alone, or needs each tool to enumerate
  its screens. Leaning on a `Screens()` method the tool implements — the guard
  then walks it, and `guard.Screens` already needs that list to check every
  screen constant has a `View()` case.
