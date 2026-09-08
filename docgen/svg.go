package docgen

import (
	"fmt"
	"html"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"github.com/richarddavenport/tuikit/harness"
)

// SVG renders a captured frame as an image.
//
// # Why SVG and not PNG
//
// A PNG needs a rasterizer, which needs a typeface, which is the open question
// in decision 30 — and answering it in order to write documentation would be
// deciding it for the wrong reason. SVG needs no typeface of its own: it names
// the same system monospace stack the HTML page does and lets the reader's
// machine draw it. It also keeps the text as text, so a frame in a document
// stays greppable and a screen reader can read it.
//
// # Why the grid is not left to the font
//
// A frame is a character grid, and a renderer that trusts a font's advance
// width to reproduce it is one box-drawing character away from a frame that
// does not line up — the failure the HTML page avoids by measuring nothing in
// pixels. So every span states the width it must occupy (`textLength`), and the
// viewer is required to fit it. Box-drawing, braille and CJK land on the cell
// boundary because they are told to, not because the font agreed.
//
// `lengthAdjust="spacing"` adjusts the gaps and leaves the glyphs alone.
// `spacingAndGlyphs` would distort the box-drawing characters into not meeting.
//
// # Why no <style> element
//
// Presentation attributes, not CSS. An SVG referenced from a Markdown document
// is rendered in a sanitiser's idea of SVG, and a stripped <style> block leaves
// a frame that is all one color with no error to explain it. Attributes
// survive; a stylesheet is a bet.
func SVG(frame string) string { return svgFrame(frame, defaultMetrics) }

// metrics is the cell geometry, in pixels at font size 1.
//
// Not configurable. A frame is captured at a fixed size for the same reason a
// golden is, and a page whose frames are each a different scale is a page where
// two screenshots cannot be compared.
type metrics struct {
	fontSize float64
	cellW    float64
	lineH    float64
	padX     float64
	padY     float64
	radius   float64
}

// The 0.6 advance ratio is what every monospace face in the stack uses; it is
// the value the grid is BUILT on rather than measured from, since textLength
// then holds each span to it whatever the reader has installed.
var defaultMetrics = metrics{fontSize: 13, cellW: 7.8, lineH: 16.25, padX: 14, padY: 12, radius: 7}

// The frame's own ground and default ink, matching the HTML page's rule: the
// ANSI was captured for a dark terminal, and recoloring it would report
// colors the tool does not have.
const (
	svgGround = "#0c0c0e"
	svgInk    = "#d0d0d0"
	svgFont   = "ui-monospace, 'SF Mono', SFMono-Regular, Menlo, Consolas, 'DejaVu Sans Mono', monospace"
)

func svgFrame(frame string, m metrics) string {
	rows := harness.Rows(frame)
	cols := harness.Width(frame)

	w := float64(cols)*m.cellW + 2*m.padX
	h := float64(len(rows))*m.lineH + 2*m.padY

	var b strings.Builder
	fmt.Fprintf(&b, `<svg xmlns="http://www.w3.org/2000/svg" width="%s" height="%s" viewBox="0 0 %s %s" font-family="%s" font-size="%s" role="img">`,
		num(w), num(h), num(w), num(h), html.EscapeString(svgFont), num(m.fontSize))

	// A label rather than nothing, because an image with no text alternative is
	// an image a reader using one is simply not shown.
	fmt.Fprintf(&b, "\n<title>%s</title>", html.EscapeString(fmt.Sprintf("terminal frame, %d columns by %d rows", cols, len(rows))))
	fmt.Fprintf(&b, "\n<rect width=\"%s\" height=\"%s\" rx=\"%s\" fill=\"%s\"/>", num(w), num(h), num(m.radius), svgGround)

	// Backgrounds first, all of them, then the text. Interleaving would let a
	// later row's background paint over the previous row's descenders.
	for y, row := range rows {
		col := 0
		for _, span := range row {
			n := widthOf(span.Text)
			if span.Style.BG != "" && n > 0 {
				fmt.Fprintf(&b, "\n<rect x=\"%s\" y=\"%s\" width=\"%s\" height=\"%s\" fill=\"%s\"/>",
					num(m.padX+float64(col)*m.cellW), num(m.padY+float64(y)*m.lineH),
					num(float64(n)*m.cellW), num(m.lineH), span.Style.BG)
			}
			col += n
		}
	}

	for y, row := range rows {
		col := 0
		for _, span := range row {
			n := widthOf(span.Text)
			if n == 0 {
				continue
			}
			if text := strings.TrimRight(span.Text, " "); text != "" {
				fill := span.Style.FG
				if fill == "" {
					fill = svgInk
				}
				// The baseline sits at 0.78 of the line box: the ratio that
				// puts a box-drawing character's horizontal on the same row as
				// the one beside it in the next span.
				fmt.Fprintf(&b, "\n<text x=\"%s\" y=\"%s\" textLength=\"%s\" lengthAdjust=\"spacing\" fill=\"%s\"%s xml:space=\"preserve\">%s</text>",
					num(m.padX+float64(col)*m.cellW), num(m.padY+float64(y)*m.lineH+m.lineH*0.78),
					num(float64(widthOf(text))*m.cellW), fill, bold(span.Style.Bold), html.EscapeString(text))
			}
			col += n
		}
	}
	b.WriteString("\n</svg>\n")
	return b.String()
}

// widthOf is display columns, not bytes and not runes — a box-drawing character
// is one column and three bytes, and a frame measured in bytes is a frame whose
// grid drifts the moment it draws a border.
func widthOf(s string) int { return lipgloss.Width(s) }

func bold(on bool) string {
	if on {
		return ` font-weight="600"`
	}
	return ""
}

// num formats a coordinate without a trailing ".00", so the file stays readable
// and a diff between two captures shows what changed rather than every number.
func num(f float64) string {
	s := fmt.Sprintf("%.2f", f)
	s = strings.TrimRight(s, "0")
	return strings.TrimSuffix(s, ".")
}
