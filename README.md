# tuikit

A TUI framework for developers and agents. Go, Bubble Tea, Lip Gloss.

It exists because four tools — [swarmctl], [pgctl], [azctl], [dugo] — arrived at
the same shape independently, and that shape currently travels by copy. tuikit is
that shape as a module: the vocabulary an interface is allowed to use, the guards
that hold it closed, and (soon) the components, the capture harness, and a
scaffolder.

[swarmctl]: https://github.com/richarddavenport/swarmctl
[pgctl]: https://github.com/richarddavenport/pgctl
[azctl]: https://github.com/richarddavenport/azctl
[dugo]: https://github.com/richarddavenport/dugo

## What is here now

**`theme`** — the vocabulary. Nine colour roles named by role and never by hue,
as ANSI 256 indices; a closed glyph allow-list; and `Hex` for anything drawing
the palette outside a terminal.

```go
pal := theme.Default
pal.Accent = "33" // this tool is blue

glyphs := theme.DefaultGlyphs.With('×', "multiplication sign — replica counts")
```

**`guard`** — two tests that hold it closed. Call them from your UI package:

```go
func TestTheInterfaceStaysInItsVocabulary(t *testing.T) {
    guard.Tokens(t, ".", palette)  // no colour literals; no unused roles
    guard.Glyphs(t, ".", glyphs)   // nothing printed that a font may not have
}
```

Both parse the package rather than pattern-matching it, so a comment explaining a
glyph is not read as printing one. Both fail rather than pass when pointed at a
directory with nothing to scan.

**`docgen`** — the vocabulary as HTML, generated from the code so it cannot
drift:

```
go run ./cmd/tuikit designsystem -out design-system -tool mytool
```

## What is coming

`harness` (deterministic frame capture, goldens, ANSI→HTML), `comp` (the
components each of the four wrote separately), `app` (the Bubble Tea shell and
its async conventions), `spec` (one command declaration → CLI, TUI screen, and a
`describe --json` manifest an agent reads), `tuikit watch`, `tuikit gallery`,
`tuikit new`. See `design/`.

## Why an agent gets on with it

The interface has a written-down vocabulary, the vocabulary is machine-readable,
and the guards turn "I used a colour that does not exist" into a test failure
rather than a review comment. What is still missing is the part that lets an
agent *see* what it built — that is `harness`, and it is next.
