# `theme` — the vocabulary

Nine colour roles, a closed glyph set, and the box characters. Small on purpose:
this is the list `guard` holds a tool to.

## Palette

`theme.Palette` carries roles, not colours — nine of them: `Accent`, `Muted`,
`Border`, `Success`, `Pending`, `Danger`, `Stderr`, `SelectionFG`,
`SelectionBG`. `theme.Default` is the shipped one, and a tool overrides the
fields it wants rather than building a palette from scratch, because the point
of the set being closed is that nine names cover it.

`Extra []Role` is the escape hatch for a role a tool genuinely needs and the
nine do not cover. It is named and declared, so `guard.Tokens` still sees it.

A tool styles with `p.Danger`, never with `lipgloss.Color("#ff5555")` — and
`guard.Tokens` fails the build on the literal. The reason is not tidiness: ANSI
0–15 are the **reader's** colours, set in their terminal, and a hex literal
overrides a choice they made deliberately.

`Hex(c)` and `Value(c)` render a role for a design-system page.

## Glyphs

`theme.GlyphSet` is a `map[rune]string`: every non-ASCII character a tool is
allowed to print, each with a name. `theme.DefaultGlyphs` is the set;
`SpinnerRange` carves out the braille block so a spinner does not need 256
entries.

`guard.Glyphs` fails on anything outside it. A closed set is what makes "will
this render in the reader's terminal" answerable at all.

## Chrome

`theme.Chrome` is the furniture: the box set, the gap between panes, the
ellipsis, the divider, the scroll track and thumb, the markers, the indent
width.

```go
ch := theme.DefaultChrome
ch.Divider = "│"   // a seam instead of a gap, everywhere, at once
ch.Gap = 4
```

Box sets: `LightBox`, `RoundedBox`, `HeavyBox`, `ASCIIBox`. A tool
that wants a heavier frame changes one field; `comp.Rule` and `comp.Pane` follow
it, because the rule under a header is the same character the pane's top is.

`guard.Chrome` checks the chrome's own characters are in the glyph set —
otherwise you have configured furniture the terminal cannot print.
`guard.Furniture` catches the other direction: a `─` typed by hand into a `Text`
call instead of taken from here.

## What it cannot do

- **No stylesheets.** There is no cascade, no selectors, no CSS-like layer.
  Textual has that and does it well; this is nine names and a map.
- **No per-component theming.** A role means the same thing everywhere. That is
  what makes a design-system page truthful.
- **No runtime palette editing.** `term.Colors` reads the terminal's actual ANSI
  values so a picture can match them, but the roles themselves are compiled in.
