package comp

import (
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
)

// ANSI turns a subprocess's output into [Segment]s.
//
// # Why this is in comp and not in harness
//
// Every operator tool shells out to something, and what comes back has color
// in it. gcpeasy shows `gcloud`, `kubectl` and a Rails console in a pane, and
// wrote its own SGR parser to do it (issue 63). Without this a tool either
// strips the color — losing information the other program went to the trouble
// of sending — or writes the parser again.
//
// [harness.Strip] parses the same sequences and is the wrong side of the fence:
// it reads back a frame tuikit itself wrote, where the styles came from
// [theme.Palette] in the first place. This is the open world.
//
// # Why the colors are passed through
//
// Decision 28 puts a tool's own color on ANSI 0-15 so the reader's theme wins.
// That rule is about DESIGN, and a subprocess's output is not the tool's design
// — it is data, the same way the words are. A tool does not rewrite `kubectl`'s
// nouns and should not rewrite its colors.
//
// So 0-15 become [lipgloss.ANSIColor], which the reader's terminal theme still
// renders. 256-color and truecolor are passed through as they arrived, because
// dropping them would lose the distinction the program was drawing.
//
// # What it does not do
//
// Cursor addressing, scroll regions, the alternate screen. That is a terminal
// emulator; decision 27 says tuikit does not host one, and gcpeasy's
// interactive pane is what it costs when you try.
//
// Carriage return, backspace and tab ARE handled, because they are how a
// program draws a progress line and dropping them leaves the frames unreadable.
func ANSI(s string) []Segment {
	lines := ANSILines(s)
	if len(lines) == 0 {
		return nil
	}
	return lines[0]
}

// ANSILines is [ANSI] over output that spans several lines.
//
// The style carries ACROSS a newline, because that is what a terminal does: a
// program that sets green and prints three lines gets three green lines, and a
// parser that reset at each newline would show one.
//
// A line is assembled in a buffer of cells before it becomes segments, so a `\r`
// can go back and overwrite what came before it — which is how every progress
// bar in the world is drawn.
func ANSILines(s string) [][]Segment {
	raws := strings.Split(s, "\n")
	out := make([][]Segment, 0, len(raws))
	style := (*lipgloss.Style)(nil)

	for _, raw := range raws {
		raw = strings.TrimSuffix(raw, "\r")
		line, next := ansiLine(raw, style)
		style = next
		out = append(out, line)
	}
	return out
}

// cell is one column of a line being assembled: what is in it and how it looks.
type cell struct {
	text  string
	style *lipgloss.Style
}

// ansiLine renders one line into cells, then runs them together into segments.
//
// Cells rather than a string builder, because `\r` and `\b` MOVE, and a builder
// can only append. A progress line that redraws itself ten times must come out
// as its final state, not as ten concatenated states.
func ansiLine(s string, style *lipgloss.Style) ([]Segment, *lipgloss.Style) {
	var cells []cell
	col := 0

	put := func(text string) {
		for col >= len(cells) {
			cells = append(cells, cell{text: " "})
		}
		cells[col] = cell{text: text, style: style}
		col++
	}

	for i := 0; i < len(s); {
		switch c := s[i]; {
		case c == 0x1b:
			n, next, ok := sgr(s[i:], style)
			if !ok {
				// Some other escape — a cursor move, a title. Skipped whole
				// rather than printed: half an escape sequence on screen is
				// worse than none of it.
				i += n
				continue
			}
			style = next
			i += n
		case c == '\r':
			col = 0
			i++
		case c == '\b':
			col = max(0, col-1)
			i++
		case c == '\t':
			// To the next eight-column stop, the way a terminal does it.
			for stop := (col/DefaultTab + 1) * DefaultTab; col < stop; {
				put(" ")
			}
			i++
		case c < 0x20:
			// Any other control character. Dropped: it means something to a
			// terminal and nothing to a cell grid.
			i++
		default:
			w := runeLen(s[i:])
			put(s[i : i+w])
			i += w
		}
	}
	return runTogether(cells), style
}

// runTogether merges neighboring cells that share a style, so a line of one
// color is one segment rather than eighty.
func runTogether(cells []cell) []Segment {
	var out []Segment
	for _, c := range cells {
		if n := len(out); n > 0 && out[n-1].Style == c.style {
			out[n-1].Text += c.text
			continue
		}
		out = append(out, Segment{Text: c.text, Style: c.style})
	}
	return out
}

// sgr reads one escape sequence and returns its length, the style it produces,
// and whether it was an SGR at all.
func sgr(s string, style *lipgloss.Style) (n int, next *lipgloss.Style, ok bool) {
	if len(s) < 2 || s[1] != '[' {
		// Not a CSI. Skip the escape and whatever single byte follows it.
		return min(2, len(s)), style, false
	}
	end := 2
	for end < len(s) && (s[end] == ';' || s[end] == ':' || (s[end] >= '0' && s[end] <= '9')) {
		end++
	}
	if end >= len(s) {
		return len(s), style, false
	}
	final := s[end]
	end++
	if final != 'm' {
		return end, style, false
	}
	return end, apply(style, params(s[2:end-1])), true
}

func params(s string) []int {
	if s == "" {
		// A bare CSI m is CSI 0 m — reset.
		return []int{0}
	}
	fields := strings.FieldsFunc(s, func(r rune) bool { return r == ';' || r == ':' })
	out := make([]int, 0, len(fields))
	for _, f := range fields {
		n, err := strconv.Atoi(f)
		if err != nil {
			n = 0
		}
		out = append(out, n)
	}
	return out
}

// apply folds SGR parameters into a style.
//
// Returns nil for a full reset rather than an empty style, so an unstyled run
// compares equal to the zero value and runTogether merges it with its
// neighbors.
func apply(style *lipgloss.Style, ps []int) *lipgloss.Style {
	cur := lipgloss.NewStyle()
	if style != nil {
		cur = *style
	}
	reset := true

	for i := 0; i < len(ps); i++ {
		switch p := ps[i]; {
		case p == 0:
			cur = lipgloss.NewStyle()
		case p == 1:
			cur, reset = cur.Bold(true), false
		case p == 2:
			cur, reset = cur.Faint(true), false
		case p == 3:
			cur, reset = cur.Italic(true), false
		case p == 4:
			cur, reset = cur.Underline(true), false
		case p == 7:
			cur, reset = cur.Reverse(true), false
		case p == 9:
			cur, reset = cur.Strikethrough(true), false
		case p == 22:
			cur, reset = cur.Bold(false).Faint(false), false
		case p == 23:
			cur, reset = cur.Italic(false), false
		case p == 24:
			cur, reset = cur.Underline(false), false
		case p == 27:
			cur, reset = cur.Reverse(false), false
		case p >= 30 && p <= 37:
			cur, reset = cur.Foreground(lipgloss.ANSIColor(p-30)), false
		case p == 39:
			cur, reset = cur.UnsetForeground(), false
		case p >= 40 && p <= 47:
			cur, reset = cur.Background(lipgloss.ANSIColor(p-40)), false
		case p == 49:
			cur, reset = cur.UnsetBackground(), false
		case p >= 90 && p <= 97:
			cur, reset = cur.Foreground(lipgloss.ANSIColor(p-90+8)), false
		case p >= 100 && p <= 107:
			cur, reset = cur.Background(lipgloss.ANSIColor(p-100+8)), false
		case p == 38 || p == 48:
			col, used := extended(ps[i:])
			if col != nil {
				if p == 38 {
					cur = cur.Foreground(col)
				} else {
					cur = cur.Background(col)
				}
				reset = false
			}
			i += used - 1
		}
	}
	if reset && len(ps) > 0 && ps[len(ps)-1] == 0 {
		return nil
	}
	if reset && style == nil {
		return nil
	}
	return &cur
}

// extended reads a 5;n (256-color) or 2;r;g;b (truecolor) argument and returns
// how many parameters it consumed.
func extended(ps []int) (lipgloss.TerminalColor, int) {
	if len(ps) < 2 {
		return nil, len(ps)
	}
	switch ps[1] {
	case 5:
		if len(ps) < 3 {
			return nil, len(ps)
		}
		return lipgloss.Color(strconv.Itoa(ps[2])), 3
	case 2:
		if len(ps) < 5 {
			return nil, len(ps)
		}
		return lipgloss.Color("#" + hex2(ps[2]) + hex2(ps[3]) + hex2(ps[4])), 5
	}
	return nil, 2
}

func hex2(n int) string {
	const digits = "0123456789abcdef"
	n = clamp(n, 0, 255)
	return string([]byte{digits[n>>4], digits[n&0xf]})
}

// runeLen is the bytes in the UTF-8 sequence starting at s[0].
func runeLen(s string) int {
	switch b := s[0]; {
	case b < 0x80:
		return 1
	case b < 0xe0:
		return min(2, len(s))
	case b < 0xf0:
		return min(3, len(s))
	default:
		return min(4, len(s))
	}
}
