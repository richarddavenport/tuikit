package guard_test

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/guard"
	"github.com/richarddavenport/tuikit/theme"
)

func pkg(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "ui.go"), []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

// The bug it exists for: a chrome character typed into a draw.
func TestFurnitureCatchesAHardcodedRule(t *testing.T) {
	rec := &fake{}
	guard.Furniture(rec, pkg(t, `package ui

func draw(c *Canvas, r Rect) {
	c.Fill(r, "─", nil, id)
}
`), theme.DefaultChrome)

	if len(rec.msgs) != 1 {
		t.Fatalf("reported %d, want 1: %v", len(rec.msgs), rec.msgs)
	}
	if !strings.Contains(rec.msgs[0], "ui.go:4") {
		t.Errorf("the message does not point at the line: %s", rec.msgs[0])
	}
}

// The two false-positive classes that made this look intractable when it was
// filed. Both NAME a chrome character in order to allow it; neither draws one.
func TestFurnitureIgnoresDeclarations(t *testing.T) {
	rec := &fake{}
	guard.Furniture(rec, pkg(t, `package ui

var entry = Entry{
	Glyphs: []string{"┌", "─", "┐", "│", "└", "┘"},
}

var set = theme.DefaultGlyphs.With(
	'╭', "rounded box drawing", '╮', "rounded box drawing",
)

var box = theme.BoxSet{"┌", "─", "┐", "│", "│", "└", "─", "┘"}
`), theme.DefaultChrome)

	if len(rec.msgs) != 0 {
		t.Errorf("a declaration was reported as a draw: %v", rec.msgs)
	}
}

// Reading it from the chrome is the fix, so it must not also be an error.
func TestFurnitureAllowsTheChromesOwnCharacter(t *testing.T) {
	rec := &fake{}
	guard.Furniture(rec, pkg(t, `package ui

func draw(c *Canvas, r Rect) {
	c.Fill(r, c.Chrome().Box.Top, nil, id)
	c.Text(r.X, r.Y, "plain ascii text", nil, id)
}
`), theme.DefaultChrome)

	if len(rec.msgs) != 0 {
		t.Errorf("reported: %v", rec.msgs)
	}
}

// A non-chrome glyph in a draw is the tool's own business — a spinner frame, a
// status dot. This guard is about furniture, and Glyphs already covers whether
// a character is allowed at all.
func TestFurnitureLeavesNonChromeGlyphsAlone(t *testing.T) {
	rec := &fake{}
	guard.Furniture(rec, pkg(t, `package ui

func draw(c *Canvas, r Rect) {
	c.Text(r.X, r.Y, "●", nil, id)
}
`), theme.DefaultChrome)

	if len(rec.msgs) != 0 {
		t.Errorf("reported a status glyph as furniture: %v", rec.msgs)
	}
}

type fake struct{ msgs []string }

func (f *fake) Helper()                        {}
func (f *fake) Errorf(format string, a ...any) { f.msgs = append(f.msgs, sprint(format, a...)) }
func (f *fake) Fatalf(format string, a ...any) { f.msgs = append(f.msgs, sprint(format, a...)) }

func sprint(format string, a ...any) string { return fmt.Sprintf(format, a...) }
