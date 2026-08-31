package harness

import (
	"fmt"
	"html"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/theme"
)

// Lines splits a frame.
func Lines(frame string) []string { return strings.Split(frame, "\n") }

// Width is the widest visible line of a frame, in columns.
//
// Measured with lipgloss.Width rather than len: a right arrow is three bytes
// and one column, and the escape sequences are not columns at all. Anything
// that counts bytes here reports a frame as wider than it is and then
// "corrects" it into a broken one.
func Width(frame string) int {
	var w int
	for _, line := range Lines(frame) {
		if n := lipgloss.Width(line); n > w {
			w = n
		}
	}
	return w
}

// Strip removes every escape sequence, leaving what a reader sees.
//
// This is what a golden holds. Colour is captured for looking at; a diff in a
// pull request wants the shape, and a golden full of escape sequences is a
// golden nobody reads.
//
// The trap, which the prototype hit and which is why this is one function with
// a test rather than a regex at each call site: an escape sequence is
// terminated by a letter, and the SGR terminator is `m`. A stripper that ends a
// sequence at any of [A-Za-z] therefore eats the `m` and leaves `[0;38;5;205`
// behind — stripping exactly the colour it was meant to preserve, and leaving
// the digits on screen.
func Strip(frame string) string {
	var b strings.Builder
	b.Grow(len(frame))

	runes := []rune(frame)
	for i := 0; i < len(runes); i++ {
		if runes[i] != 0x1b {
			b.WriteRune(runes[i])
			continue
		}
		// CSI: ESC [ parameters... final-byte, where the final byte is 0x40-0x7E.
		j := i + 1
		if j < len(runes) && runes[j] == '[' {
			j++
			for j < len(runes) && runes[j] >= 0x20 && runes[j] <= 0x3f {
				j++
			}
			if j < len(runes) {
				j++ // the final byte, `m` included
			}
			i = j - 1
			continue
		}
		// Anything else introduced by ESC: drop the introducer and one byte.
		if j < len(runes) {
			i = j
		}
	}
	return b.String()
}

// HTML renders a frame as a self-contained block, colour intact.
//
// Three rules the prototype settled, encoded here rather than left to whoever
// publishes the page:
//
//   - The frame keeps a dark ground in both light and dark themes. The ANSI was
//     captured for a dark terminal; recolouring it reports colours the tool does
//     not have.
//   - The system monospace stack, with no webfont. Box-drawing and braille must
//     share one set of advance widths, and a webfont plus a fallback for the
//     glyphs it lacks guarantees a frame that does not line up.
//   - Nothing is measured in pixels. One character is one cell.
func HTML(frame string) string {
	var b strings.Builder
	b.WriteString(`<div class="tuikit-frame">`)

	var open bool
	for i, line := range Lines(frame) {
		if i > 0 {
			b.WriteByte('\n')
		}
		open = writeLine(&b, line, open)
	}
	if open {
		b.WriteString("</span>")
	}
	b.WriteString(`</div>`)
	return b.String()
}

// writeLine converts one line, returning whether a span is still open. Spans do
// not straddle lines — a colour left switched on across a newline paints the
// page's background, which is how a frame ends up with a coloured margin.
func writeLine(b *strings.Builder, line string, open bool) bool {
	if open {
		b.WriteString("</span>")
		open = false
	}
	runes := []rune(line)
	for i := 0; i < len(runes); i++ {
		if runes[i] != 0x1b {
			b.WriteString(html.EscapeString(string(runes[i])))
			continue
		}
		j := i + 1
		if j >= len(runes) || runes[j] != '[' {
			if j < len(runes) {
				i = j
			}
			continue
		}
		j++
		start := j
		for j < len(runes) && runes[j] >= 0x20 && runes[j] <= 0x3f {
			j++
		}
		params := string(runes[start:j])
		final := ' '
		if j < len(runes) {
			final = runes[j]
			j++
		}
		i = j - 1

		if final != 'm' {
			continue // cursor movement and friends have no meaning in a static frame
		}
		if open {
			b.WriteString("</span>")
			open = false
		}
		if style := css(params); style != "" {
			fmt.Fprintf(b, `<span style="%s">`, style)
			open = true
		}
	}
	return open
}

// css turns SGR parameters into a style. Only what a tuikit interface emits is
// handled: reset, bold, and the 256-colour foreground and background forms.
// Anything else is ignored rather than guessed at.
func css(params string) string {
	fields := strings.Split(params, ";")
	var out []string
	for i := 0; i < len(fields); i++ {
		switch fields[i] {
		case "", "0":
			return "" // reset closes the span
		case "1":
			out = append(out, "font-weight:600")
		case "38", "48":
			// 38;5;N and 48;5;N — foreground and background from the palette.
			if i+2 < len(fields) && fields[i+1] == "5" {
				prop := "color"
				if fields[i] == "48" {
					prop = "background"
				}
				out = append(out, prop+":"+theme.Hex(lipgloss.Color(fields[i+2])))
				i += 2
			}
		default:
			if n, err := strconv.Atoi(fields[i]); err == nil {
				if hex := basicColor(n); hex != "" {
					out = append(out, hex)
				}
			}
		}
	}
	return strings.Join(out, ";")
}

// basicColor handles the sixteen direct SGR colours, which lipgloss emits for a
// palette index below 16.
func basicColor(n int) string {
	switch {
	case n >= 30 && n <= 37:
		return "color:" + theme.Hex(lipgloss.Color(strconv.Itoa(n-30)))
	case n >= 90 && n <= 97:
		return "color:" + theme.Hex(lipgloss.Color(strconv.Itoa(n-90+8)))
	case n >= 40 && n <= 47:
		return "background:" + theme.Hex(lipgloss.Color(strconv.Itoa(n-40)))
	case n >= 100 && n <= 107:
		return "background:" + theme.Hex(lipgloss.Color(strconv.Itoa(n-100+8)))
	}
	return ""
}

// FrameCSS is the style a page needs for HTML's output. Dark ground in both
// themes, system mono, no ligatures, one character one cell.
const FrameCSS = `
.tuikit-frame {
  font-family: ui-monospace, "SF Mono", SFMono-Regular, Menlo, Consolas, "DejaVu Sans Mono", monospace;
  font-size: 13px; line-height: 1.25;
  font-variant-ligatures: none; font-feature-settings: "liga" 0, "calt" 0;
  white-space: pre; tab-size: 4;
  background: #0c0c0e; color: #d0d0d0;
  padding: 14px 16px; border-radius: 6px;
  overflow-x: auto; border: 1px solid #26262c;
}
`

// widthOf and sprintf keep golden.go free of imports it would otherwise need
// for two calls.
func widthOf(line string) int { return lipgloss.Width(line) }

func sprintf(format string, args ...any) string { return fmt.Sprintf(format, args...) }
