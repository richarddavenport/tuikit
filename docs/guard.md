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
| `Tokens` | a color literal instead of a palette role — and a role nothing draws with |
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
nothing, and `harness.ShapeSurvivesColor` checks a frame still reads with every
escape sequence stripped.

## `guard.Derived` — the fixture has to be the engine's answer

Decision 1 says the fixture the screens are rendered from is a value the engine
returns. Nothing enforced it, and in the board it quietly stopped being one.

```go
guard.Derived(t, ".", "engine", "Truth", "Live")
```

Pass the types that are the engine's **answers** rather than its inputs — what a
gather-and-reconcile produced, not the arguments it was called with. A type a
test *should* construct must not be listed: this is not "tests may not build
structs", it is "these particular structs are conclusions, and a conclusion
typed by hand is a conclusion nobody checked".

**What it catches.** the board's `fixture()` set a field by hand while the
engine composed the same field as a sentence, and the renderer composed it
again. The golden showed a row disagreeing with itself on every run and looked
right, because it was checking the renderer against a world invented three
hundred lines away in the same file. Four bugs shipped past 72 goldens; three
for this reason.

The compounding part is the worst of it: a hand-made fixture makes a wrong
screen look correct **and** hands you an easy way to keep it that way — editing
the fixture to match the renderer.

An empty literal is allowed. A zero value is "nothing yet", not an invented
world.


## When the guard is wrong: `guard.Except`

`guard.Tokens` refuses a color built from a literal, because a color that did
not come from the palette is a decision nobody made. That is right almost
always and there is a real exception.

lazygit's `presentation/icons/file_icons.go` holds **743 hex literals**. They
are file-type brand colors — the Go gopher's blue, the Rust orange — and the
whole point of them is that they are the same everywhere. Decision 28 puts
color on ANSI 0–15 so the reader's theme wins, and that reasoning does not
reach a brand mark.

```go
guard.Tokens(t, ".", Palette,
	guard.Except("file_icons.go", "file-type brand colors, not theme"))
```

**The reason is required**, and an empty one panics. Without it this is a
suppression flag wearing a better name, and the value of the guard is that a
color outside the palette is a *decision*.

**And it cannot rot.** Two things fail:

- an exemption naming a file that is not in the directory — it was renamed or
  deleted, and the next file to take that name would be silently unguarded
- an exemption on a file that builds no colors from literals — nothing is being
  excused, so it is a hole nobody is using and nobody will notice opening

A stale exemption is worse than no exemption, because it reads as though
somebody checked.


## The properties they buy

**Colors are the reader's.** ANSI 0–15 come from their terminal theme. A hex
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
