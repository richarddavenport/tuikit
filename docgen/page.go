package docgen

import (
	"fmt"
	"html"
	"sort"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/theme"
)

// A page is one card in the design system. Everything it draws goes through
// here, so there is exactly one place that decides what a terminal looks like
// when it is not a terminal.

// renderer is a resolved DesignSystem: no page has to ask whether it was given
// a palette.
type renderer struct {
	tool   string
	pal    theme.Palette
	glyphs theme.GlyphSet
}

// terminalCSS is the grid itself.
//
// The three rules that make this a terminal rather than a picture of one:
// a monospace font with no ligatures, so one character is one cell; a line
// height that stacks box-drawing characters into unbroken lines; and no
// padding inside the grid, because a terminal has none. Everything else is
// chrome around it.
//
// The mono stack is the system's own and deliberately carries no webfont.
// Box-drawing and braille glyphs must share one set of advance widths, and
// mixing a webfont with a fallback for the glyphs it lacks guarantees a frame
// that does not line up.
const terminalCSS = `
:root { color-scheme: dark; }
* { box-sizing: border-box; }
body {
  margin: 0; padding: 32px;
  background: #16161a; color: #d0d0d0;
  font: 14px/1.55 ui-sans-serif, -apple-system, "Segoe UI", system-ui, sans-serif;
}
h1 { font-size: 20px; margin: 0 0 4px; color: #f0f0f0; letter-spacing: -0.01em; }
h2 { font-size: 13px; margin: 28px 0 10px; color: #a0a0a0;
     text-transform: uppercase; letter-spacing: 0.08em; font-weight: 600; }
p  { margin: 0 0 16px; max-width: 62ch; color: #9a9a9a; }
p.lede { color: #b8b8b8; }
code { font-family: var(--mono); font-size: 0.92em; color: #c8c8c8; }

/* The grid. One character is one cell; nothing here is measured in pixels. */
:root { --mono: ui-monospace, "SF Mono", SFMono-Regular, Menlo, Consolas, "DejaVu Sans Mono", monospace; }
.term {
  font-family: var(--mono);
  font-size: 13px; line-height: 1.25;
  font-variant-ligatures: none; font-feature-settings: "liga" 0, "calt" 0;
  white-space: pre; tab-size: 4;
  background: #0c0c0e; color: #d0d0d0;
  padding: 14px 16px; border-radius: 6px;
  overflow-x: auto; border: 1px solid #26262c;
}
.term.light { background: #fbfbfa; color: #303030; border-color: #e2e2e0; }

/* Roles. The only colors in this document — see foundations/colors. */
`

// paletteCSS turns the roles into the classes every page draws with, so a page
// literally cannot use a color the interface does not have.
func (r renderer) paletteCSS() string {
	var b strings.Builder
	for _, role := range r.pal.Roles() {
		fmt.Fprintf(&b, ".%s { color: %s; }\n", cssClass(role.Name), theme.Hex(role.Color))
	}
	// The one painted background, which needs both halves at once.
	fmt.Fprintf(&b, ".selected { color: %s; background: %s; font-weight: 600; }\n",
		theme.Hex(r.pal.SelectionFG), theme.Hex(r.pal.SelectionBG))
	b.WriteString(".bold { font-weight: 600; }\n")
	return b.String()
}

func cssClass(role string) string { return strings.ToLower(role) }

// hexOf is theme.Hex, named short because the color pages call it constantly.
func hexOf(c lipgloss.TerminalColor) string { return theme.Hex(c) }

// valueOf is the color as the tool declared it, which is what a design system
// page should say alongside what it resolves to.
func valueOf(c lipgloss.TerminalColor) string { return theme.Value(c) }

// classFor is the class a role is actually drawn with.
//
// The two selection roles are the one pair that is never used apart — a
// foreground on the painted row and the paint behind it are one decision — so
// they share a class rather than pretending to be independent.
func classFor(role string) string {
	if role == "SelectionFG" || role == "SelectionBG" {
		return "selected"
	}
	return cssClass(role)
}

// page renders one card. The @dsCard marker on the first line is what puts it
// in the Design System pane, and must stay there.
func (r renderer) page(group, title, lede, body string) string {
	return fmt.Sprintf(`<!-- @dsCard group=%q -->
<meta charset="utf-8">
<title>%s</title>
<style>%s%s</style>
<h1>%s</h1>
<p class="lede">%s</p>
%s
`, group, html.EscapeString(title), terminalCSS, r.paletteCSS(),
		html.EscapeString(title), html.EscapeString(lede), body)
}

// term wraps terminal content. Lines are given as already-marked-up strings —
// see span, which is the only way to color one.
func term(lines ...string) string {
	return `<div class="term">` + strings.Join(lines, "\n") + `</div>`
}

// span colors a run of text with a role.
//
// It PANICS on a role that does not exist. A design system whose swatch says
// one thing and whose components quietly say another is worse than no design
// system, and this is a generator — a panic is a build failure, which is
// exactly the right loudness.
func (r renderer) span(role, text string) string {
	if !r.knownRole(role) {
		panic("docgen: no such role: " + role + " — add it to the palette or use one that exists")
	}
	return fmt.Sprintf(`<span class="%s">%s</span>`, cssClass(role), html.EscapeString(text))
}

// plain is text at the terminal's default color — which a tuikit interface
// never sets, so it is whatever the reader's terminal is.
func plain(text string) string { return html.EscapeString(text) }

func (r renderer) knownRole(role string) bool {
	if role == "selected" || role == "bold" {
		return true
	}
	for _, have := range r.pal.Roles() {
		if cssClass(have.Name) == cssClass(role) {
			return true
		}
	}
	return false
}

// glyphRow is one entry of the allow-list, in a stable order — a map is not one.
type glyphRow struct {
	R   rune
	Why string
}

func (r renderer) glyphRows() []glyphRow {
	out := make([]glyphRow, 0, len(r.glyphs))
	for g, why := range r.glyphs {
		out = append(out, glyphRow{g, why})
	}
	sort.Slice(out, func(i, j int) bool { return out[i].R < out[j].R })
	return out
}
