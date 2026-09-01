package guard

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/theme"
)

// A guard that never fires is indistinguishable from a clean package, so every
// case here is checked in both directions: the offending source must fail, and
// the innocent source that looks like it must not.

func TestTokensRejectsAColourLiteral(t *testing.T) {
	dir := pkg(t, `package ui

import "github.com/charmbracelet/lipgloss"

var border = lipgloss.Color("240")
`)
	got := run(t, func(rec T) { Tokens(rec, dir, allRolesUsed()) })
	want(t, got, "builds a colour from a literal")
}

func TestTokensAcceptsAColourFromThePalette(t *testing.T) {
	dir := pkg(t, `package ui

import "github.com/charmbracelet/lipgloss"

var pal = tuikitPalette()
var style = lipgloss.NewStyle().Foreground(pal.Border)
`+allRoleUses())
	silent(t, run(t, func(rec T) { Tokens(rec, dir, theme.Default) }))
}

// The palette may be reached through any name — a tool can import theme under
// an alias, or hold a Palette value of its own. What is constant is the dot.
func TestTokensSeesARoleReachedThroughAnyName(t *testing.T) {
	dir := pkg(t, "package ui\n\n"+strings.ReplaceAll(allRoleUses(), "pal.", "myTheme."))
	silent(t, run(t, func(rec T) { Tokens(rec, dir, theme.Default) }))
}

func TestTokensRejectsARoleNothingDrawsWith(t *testing.T) {
	dir := pkg(t, "package ui\n\n"+allRoleUses())

	p := theme.Default
	p.Extra = []theme.Role{{Name: "Info", Color: lipgloss.Color("39"), Why: "a note"}}

	got := run(t, func(rec T) { Tokens(rec, dir, p) })
	want(t, got, "Info is an Extra role you named and never drew with")
}

func TestGlyphsRejectsACharacterOutsideTheSet(t *testing.T) {
	dir := pkg(t, "package ui\n\nvar cursor = \"█\"\n")
	got := run(t, func(rec T) { Glyphs(rec, dir, theme.DefaultGlyphs) })
	want(t, got, "U+2588")
}

func TestGlyphsAcceptsTheAllowList(t *testing.T) {
	dir := pkg(t, "package ui\n\nvar box = \"┌─┐\"\nvar ok = \"✓ done · now\"\n")
	silent(t, run(t, func(rec T) { Glyphs(rec, dir, theme.DefaultGlyphs) }))
}

// This is the bug the parsed-not-matched decision was made for: a comment
// explaining a glyph was reported as printing it, and "a comment can say
// anything" is exactly what the guard claims.
func TestAGlyphInACommentIsNotAGlyphOnScreen(t *testing.T) {
	dir := pkg(t, "package ui\n\n// A comment quoting \"api_api ×3\" says nothing to a terminal.\nvar s = \"plain\"\n")
	silent(t, run(t, func(rec T) { Glyphs(rec, dir, theme.DefaultGlyphs) }))
}

// A tool that adds a glyph deliberately is the supported path; the guard has to
// follow the set it was given rather than the default.
func TestGlyphsFollowsAnExtendedSet(t *testing.T) {
	dir := pkg(t, "package ui\n\nvar n = \"×3\"\n")

	want(t, run(t, func(rec T) { Glyphs(rec, dir, theme.DefaultGlyphs) }), "U+00D7")
	silent(t, run(t, func(rec T) {
		Glyphs(rec, dir, theme.DefaultGlyphs.With('×', "multiplication sign — replica counts"))
	}))
}

// A test file may name a colour it is checking for, or print a glyph in a
// failure message. Neither reaches a terminal.
func TestGuardsIgnoreTestFiles(t *testing.T) {
	dir := pkg(t, "package ui\n\n"+allRoleUses())
	write(t, dir, "ui_test.go", "package ui\n\nimport \"github.com/charmbracelet/lipgloss\"\n\nvar c = lipgloss.Color(\"1\")\nvar g = \"█\"\n")

	silent(t, run(t, func(rec T) { Tokens(rec, dir, theme.Default) }))
	silent(t, run(t, func(rec T) { Glyphs(rec, dir, theme.DefaultGlyphs) }))
}

// The failure worth guarding against most: a guard pointed at the wrong
// directory scans nothing and passes, and everyone believes the rule is on.
func TestAGuardWithNothingToScanFailsRatherThanPasses(t *testing.T) {
	empty := t.TempDir()
	want(t, run(t, func(rec T) { Tokens(rec, empty, theme.Default) }), "no Go files")

	only := t.TempDir()
	write(t, only, "ui_test.go", "package ui\n")
	want(t, run(t, func(rec T) { Glyphs(rec, only, theme.DefaultGlyphs) }), "only test files")
}

// --- fixtures and the recorder ------------------------------------------

// allRoleUses is a line drawing with every role, so a fixture testing one guard
// does not trip the other's "role never used" arm.
func allRoleUses() string {
	var b strings.Builder
	b.WriteString("var _ = []any{\n")
	for _, r := range theme.Default.Roles() {
		fmt.Fprintf(&b, "\tpal.%s,\n", r.Name)
	}
	b.WriteString("}\n")
	return b.String()
}

func allRolesUsed() theme.Palette { return theme.Default }

func pkg(t *testing.T, src string) string {
	t.Helper()
	dir := t.TempDir()
	write(t, dir, "ui.go", src)
	return dir
}

func write(t *testing.T, dir, name, src string) {
	t.Helper()
	if err := os.WriteFile(filepath.Join(dir, name), []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
}

// recorder stands in for *testing.T so a guard's own failures can be asserted
// on. Fatalf panics because that is what Fatalf means — stop here — and run
// recovers it.
type recorder struct{ msgs []string }

type fatal struct{}

func (r *recorder) Helper() {}
func (r *recorder) Errorf(format string, args ...any) {
	r.msgs = append(r.msgs, fmt.Sprintf(format, args...))
}
func (r *recorder) Fatalf(format string, args ...any) {
	r.Errorf(format, args...)
	panic(fatal{})
}

func run(t *testing.T, f func(T)) []string {
	t.Helper()
	rec := &recorder{}
	func() {
		defer func() {
			if p := recover(); p != nil {
				if _, ok := p.(fatal); !ok {
					panic(p)
				}
			}
		}()
		f(rec)
	}()
	return rec.msgs
}

func want(t *testing.T, msgs []string, substr string) {
	t.Helper()
	for _, m := range msgs {
		if strings.Contains(m, substr) {
			return
		}
	}
	t.Errorf("no failure mentioned %q; got %v", substr, msgs)
}

func silent(t *testing.T, msgs []string) {
	t.Helper()
	if len(msgs) != 0 {
		t.Errorf("guard fired on acceptable source: %v", msgs)
	}
}

// The nine are the framework's vocabulary, not the tool's promise to draw with
// all of them.
//
// A tool that never paints a selection background cannot "drop" SelectionBG —
// it is a field on a struct it inherited. This check came from swarmctl, where
// the palette was the tool's own and an unused role really was dead; azctl's
// migration is where that stopped being true, and it failed on four roles at
// once for the crime of being a browser rather than a table.
func TestTokensDoesNotDemandAToolUseEveryInheritedRole(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "ui.go", `package ui

import "github.com/richarddavenport/tuikit/theme"

var p = theme.Default
var _ = p.Accent
`)
	if got := run(t, func(rec T) { Tokens(rec, dir, theme.Default) }); len(got) != 0 {
		t.Errorf("a tool using one of the nine was told off for the other eight: %v", got)
	}
}

// The spinner range is the documented escape hatch, and the guard has to honour
// the same answer GlyphSet.Printable gives. Two functions answering "may this
// be printed" differently is worse than either answer.
func TestGlyphsHonoursTheSpinnerRange(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "ui.go", "package ui\n\nvar frame = \"⠿\"\n")

	if got := run(t, func(rec T) { Glyphs(rec, dir, theme.DefaultGlyphs) }); len(got) != 0 {
		t.Errorf("a braille spinner frame was rejected: %v", got)
	}
}
