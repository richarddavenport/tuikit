# `guard` — the tests that hold it closed

Nine functions. Each reads your own source, or your own declarations, and fails
a build when the interface would lie about itself.

They are ordinary Go tests. You call them from your own package:

```go
func TestGuards(t *testing.T) {
	guard.Tokens(t, "internal/tui", theme.Default)
	guard.Glyphs(t, "internal/tui", theme.DefaultGlyphs)
	guard.Chrome(t, theme.DefaultChrome, theme.DefaultGlyphs)
	guard.Furniture(t, "internal/tui", theme.DefaultChrome)
	guard.Engine(t, "internal/estate")
	guard.Reachable(t, root)
	guard.Reserved(t, root)
	guard.Keys(t, root, sections)
	guard.Screens(t, "internal/tui", "screen", drawn)
}
```

## What each one catches

| | |
| --- | --- |
| `Tokens` | a colour literal instead of a palette role — and a role nothing draws with |
| `Glyphs` | a character outside the glyph set. Comments don't count; only what reaches the screen |
| `Chrome` | chrome configured with characters the glyph set cannot print |
| `Furniture` | a `─` or `▸` typed by hand into `Set`/`Text`/`Fill` instead of taken from the chrome |
| `Engine` | your domain layer importing a terminal library. The engine has never heard of a terminal |
| `Reachable` | a command with a `Target` and no `Key` — mouse-only, and an agent cannot click |
| `Reserved` | `ctrl+c`, `q`, `esc` or `?` bound to something other than what they mean |
| `Keys` | a binding with no help entry, or a help entry with no binding |
| `Screens` | a screen constant nothing draws |

Two more live in `harness` because they need to run the program rather than read
it: `harness.Hints` presses every advertised key and fails on one that does
nothing, and `harness.ShapeSurvivesColour` checks a frame still reads with every
escape sequence stripped.

## The properties they buy

**Colours are the reader's.** ANSI 0–15 come from their terminal theme. A hex
literal overrides a choice they made on purpose, and `Tokens` is why that cannot
happen by accident.

**The help is true.** `Keys` and `harness.Hints` between them mean an advertised
key exists *and* does something. A footer that lies is worse than no footer.

**The engine is portable.** `Engine` scans imports for a deny-list
(`TerminalPackages` by default, or your own `Denied` list, each with a reason).
Your domain layer stays testable, scriptable and reusable because it cannot
learn what a terminal is.

**Everything is reachable by keyboard.** Not an accessibility gesture: a
multiplexer that keeps right-click gets the event first, and an agent driving
the tool has no pointer at all.

## What they cannot do

- **They read the package you point them at**, not your dependencies. A library
  you import can print whatever it likes.
- **They are AST scans**, so a character assembled at runtime from a variable is
  invisible to `Glyphs` and `Furniture`. They catch the mistake people actually
  make — typing the character — not a determined evasion.
- **`Screens` needs to be told** which values are drawn, via the `drawn` func,
  because it cannot execute your draw switch.
- **A guard pointed at nothing fails**, deliberately. A guard that silently
  scanned zero files would pass forever.
