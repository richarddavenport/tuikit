// Package docgen writes a tool's interface out as HTML, so it can be designed
// somewhere that is not a terminal.
//
// The problem it solves: a design tool draws pixels, and a TUI draws a
// CHARACTER GRID — one monospace cell per column, 256 ANSI colors, box-drawing
// characters, and a closed set of glyphs. A mockup made without those
// constraints is a picture of a tool that cannot be built. So the bundle is a
// faithful grid: every color comes from the tool's palette, every glyph comes
// from the same allow-list guard.Glyphs enforces, and nothing is positioned in
// pixels.
//
// That direction matters. The design system is GENERATED FROM THE CODE, never
// maintained beside it — a palette written by hand is a palette that drifts,
// and a drifted design system describes a tool that does not exist.
package docgen

import (
	"os"
	"path/filepath"

	"github.com/richarddavenport/tuikit/theme"
)

// Page is one card of the bundle.
type Page struct {
	// Path is relative, with directories: "foundations/colors.html".
	Path string
	// HTML is a self-contained document. Its first line carries the @dsCard
	// marker that puts it in a Design System pane.
	HTML string
}

// DesignSystem renders a tool's vocabulary.
//
// Only the foundations are here so far — the colors, the glyphs, and the grid
// they sit on. Component and screen cards arrive with the components; a card
// depicting a component that does not exist would be the drift this package
// exists to prevent.
type DesignSystem struct {
	// Tool is the name used in prose. Empty means "the interface".
	Tool string
	// Palette and Glyphs default to theme.Default and theme.DefaultGlyphs when
	// zero, so the common case needs no configuration.
	Palette theme.Palette
	Glyphs  theme.GlyphSet
}

// Pages renders the bundle in the order it is worth reading.
func (d DesignSystem) Pages() []Page {
	r := d.renderer()
	return []Page{
		{"foundations/colors.html", r.colorsPage()},
		{"foundations/glyphs.html", r.glyphsPage()},
		{"foundations/grid.html", r.gridPage()},
	}
}

// Write puts the bundle on disk and returns the paths written.
func (d DesignSystem) Write(dir string) ([]string, error) {
	var written []string
	for _, p := range d.Pages() {
		path := filepath.Join(dir, p.Path)
		if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
			return written, err
		}
		if err := os.WriteFile(path, []byte(p.HTML), 0o600); err != nil {
			return written, err
		}
		written = append(written, path)
	}
	return written, nil
}

// renderer resolves the zero values once, so no page has to ask whether it was
// given a palette.
func (d DesignSystem) renderer() renderer {
	r := renderer{tool: d.Tool, pal: d.Palette, glyphs: d.Glyphs}
	if r.tool == "" {
		r.tool = "the interface"
	}
	// Checked on one field rather than against a zero Palette: Palette holds a
	// slice, so it is not comparable, and a tool that set only Extra still has
	// no colors.
	if r.pal.Accent == nil {
		r.pal = theme.Default
	}
	if r.glyphs == nil {
		r.glyphs = theme.DefaultGlyphs
	}
	return r
}
