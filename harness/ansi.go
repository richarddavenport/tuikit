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

	for i, row := range Rows(frame) {
		if i > 0 {
			b.WriteByte('\n')
		}
		for _, span := range row {
			css := span.Style.css()
			if css == "" {
				b.WriteString(html.EscapeString(span.Text))
				continue
			}
			fmt.Fprintf(&b, `<span style="%s">%s</span>`, css, html.EscapeString(span.Text))
		}
	}
	b.WriteString(`</div>`)
	return b.String()
}

// Style is how a terminal was drawing when it wrote a run of text: a foreground
// and background as hex (empty meaning the terminal's own default), and bold.
//
// Exported because a frame is written once and read by more than one renderer —
// HTML for a page, SVG for a document that needs an image. Both need the same
// answer to "what colour is this character", and two parsers would eventually
// disagree about a frame neither of them drew.
type Style struct {
	FG, BG string
	Bold   bool
}

// Span is a stretch of one line drawn with one style.
//
// Named Span rather than Run because harness.Run already drives a key through a
// model, and one package with two meanings of "run" is a package nobody can
// skim.
type Span struct {
	Text  string
	Style Style
}

// Rows reads a frame into styled runs, one slice per line.
//
// This is the ANSI reader every renderer shares. The state a terminal carries
// is threaded through the lines rather than reset at each one, because it
// carries in a terminal: a background left switched on at the end of a row
// paints the start of the next. What does NOT carry is the run itself — a run
// never straddles a newline, since a colour left open across one paints the
// page's margin rather than the frame's.
//
// Cursor movement and every other non-SGR sequence is dropped: a static frame
// has already been laid out, so a sequence that moves a cursor has nothing left
// to say.
func Rows(frame string) [][]Span {
	lines := Lines(frame)
	rows := make([][]Span, 0, len(lines))

	var st sgr
	for _, line := range lines {
		var row []Span
		var text strings.Builder

		flush := func(with sgr) {
			if text.Len() == 0 {
				return
			}
			row = append(row, Span{Text: text.String(), Style: Style{FG: with.fg, BG: with.bg, Bold: with.bold}})
			text.Reset()
		}

		runes := []rune(line)
		for i := 0; i < len(runes); i++ {
			if runes[i] != 0x1b {
				text.WriteRune(runes[i])
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
				continue
			}
			flush(st)
			st = st.apply(params)
		}
		flush(st)
		rows = append(rows, row)
	}
	return rows
}

// sgr is the drawing state a terminal carries between characters.
//
// State rather than a style-per-sequence, which is the bug this replaced: a
// terminal ACCUMULATES, so `\x1b[1m` then `\x1b[38;5;205m` is bold pink, and a
// reader that treats each sequence as a complete style renders the second as
// pink with the bold silently dropped.
type sgr struct {
	fg, bg string // hex, or empty for the terminal's default
	bold   bool
}

// css is the style as a declaration list, or empty when the terminal was
// drawing with its own defaults and the markup should say nothing at all.
func (s Style) css() string {
	var out []string
	if s.FG != "" {
		out = append(out, "color:"+s.FG)
	}
	if s.BG != "" {
		out = append(out, "background:"+s.BG)
	}
	if s.Bold {
		out = append(out, "font-weight:600")
	}
	return strings.Join(out, ";")
}

// apply folds one SGR sequence into the state.
//
// Only what a terminal interface emits is handled: reset, bold, the default
// colours, and both extended colour forms. Anything else is ignored rather than
// guessed at.
func (s sgr) apply(params string) sgr {
	if params == "" {
		return sgr{} // a bare \x1b[m is a reset
	}
	fields := strings.Split(params, ";")
	for i := 0; i < len(fields); i++ {
		n, err := strconv.Atoi(fields[i])
		if err != nil {
			continue
		}
		switch {
		case n == 0:
			s = sgr{}
		case n == 1:
			s.bold = true
		case n == 22:
			s.bold = false
		case n == 39:
			s.fg = ""
		case n == 49:
			s.bg = ""
		case n == 38 || n == 48:
			hex, used := extended(fields[i+1:])
			if hex != "" {
				if n == 38 {
					s.fg = hex
				} else {
					s.bg = hex
				}
			}
			i += used
		case n >= 30 && n <= 37:
			s.fg = indexHex(n - 30)
		case n >= 90 && n <= 97:
			s.fg = indexHex(n - 90 + 8)
		case n >= 40 && n <= 47:
			s.bg = indexHex(n - 40)
		case n >= 100 && n <= 107:
			s.bg = indexHex(n - 100 + 8)
		}
	}
	return s
}

// extended reads the tail of a 38 or 48 sequence: `5;N` for a palette index,
// `2;R;G;B` for 24-bit. It returns the colour and how many fields it consumed.
//
// The 24-bit form is why this is a parser rather than a lookup. Reading the
// fields of `38;2;255;95;175` one at a time and asking what each MEANS finds
// 95 in the range of the bright-colour codes, and renders a hand-picked pink as
// bright magenta — a plausible wrong colour, which nobody questions.
func extended(rest []string) (string, int) {
	if len(rest) == 0 {
		return "", 0
	}
	switch rest[0] {
	case "5":
		if len(rest) < 2 {
			return "", len(rest)
		}
		n, err := strconv.Atoi(rest[1])
		if err != nil {
			return "", 2
		}
		return indexHex(n), 2
	case "2":
		if len(rest) < 4 {
			return "", len(rest)
		}
		var c [3]int
		for k := range c {
			v, err := strconv.Atoi(rest[k+1])
			if err != nil {
				return "", 4
			}
			c[k] = v
		}
		return fmt.Sprintf("#%02x%02x%02x", c[0], c[1], c[2]), 4
	}
	return "", 1
}

// indexHex is a palette index as the browser needs it.
func indexHex(n int) string { return theme.Hex(lipgloss.Color(strconv.Itoa(n))) }

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
