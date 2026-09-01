// Package htmlround is a throwaway prototype: it asks whether harness.HTML is
// lossless by parsing its output back into cells and comparing.
//
// Not for merge. It exists to answer one question before anything is designed
// on top of it — if a frame survives the trip to HTML and back, then a frame
// and a page are the same artefact in two encodings, and an HTML mock can be
// read as a specification of a screen. If it does not survive, the places it
// loses are exactly the places a design tool would silently lie.
//
// The comparison is over CELLS, not bytes. Byte equality of the re-emitted
// ANSI would be the wrong question: lipgloss is free to write \x1b[1;38;5;205m
// or \x1b[1m\x1b[38;5;205m for the same thing, and a test that fails on that
// is a test about lipgloss's formatting rather than about lost information.
// What must survive is what a reader sees: which rune sits in which cell, and
// in what colour.
package htmlround

import (
	"fmt"
	"strconv"
	"strings"

	"github.com/charmbracelet/lipgloss"
	"golang.org/x/net/html"

	"github.com/richarddavenport/tuikit/theme"
)

// Cell is one character position as a reader sees it. Colours are held as hex
// rather than as palette indices because that is the form both sides can
// speak: the ANSI side resolves an index through theme.Hex, and the HTML side
// already has hex. An empty string is the terminal's default.
type Cell struct {
	Rune rune
	FG   string
	BG   string
	Bold bool
}

// Grid is a frame as cells, one slice per line. Ragged on purpose — a frame's
// lines are not padded to equal width, and padding them here would invent
// content that the comparison would then confirm.
type Grid [][]Cell

// style is the SGR state a scanner carries between runes.
type style struct {
	fg, bg string
	bold   bool
}

func (s style) cell(r rune) Cell { return Cell{Rune: r, FG: s.fg, BG: s.bg, Bold: s.bold} }

// FromANSI reads a captured frame.
//
// Deliberately handles MORE than harness.css does — 24-bit colour, underline's
// neighbours, the default-colour codes 39 and 49. The point of the round trip
// is to find what the HTML side drops, and a reader that shares the writer's
// blind spots cannot see them: both sides would agree that a colour the page
// never rendered was never there.
func FromANSI(frame string) Grid {
	var g Grid
	// The style carries ACROSS lines, because it does in a terminal: a
	// background left switched on at the end of one row paints the start of the
	// next. Resetting it per line here would give the reader the same blind
	// spot as the writer, and a round trip in which both sides forget the same
	// thing reports success.
	var st style
	for _, line := range strings.Split(frame, "\n") {
		var cells []Cell

		runes := []rune(line)
		for i := 0; i < len(runes); i++ {
			if runes[i] != 0x1b {
				cells = append(cells, st.cell(runes[i]))
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
			var final rune
			if j < len(runes) {
				final = runes[j]
				j++
			}
			i = j - 1

			if final == 'm' {
				st = apply(st, params)
			}
		}
		g = append(g, cells)
	}
	return g
}

// apply folds one SGR sequence into the running style.
func apply(st style, params string) style {
	fields := strings.Split(params, ";")
	if params == "" {
		return style{} // a bare \x1b[m is a reset
	}
	for i := 0; i < len(fields); i++ {
		n, err := strconv.Atoi(fields[i])
		if err != nil {
			continue
		}
		switch {
		case n == 0:
			st = style{}
		case n == 1:
			st.bold = true
		case n == 22:
			st.bold = false
		case n == 39:
			st.fg = ""
		case n == 49:
			st.bg = ""
		case n == 38 || n == 48:
			hex, used := extended(fields[i+1:])
			if hex != "" {
				if n == 38 {
					st.fg = hex
				} else {
					st.bg = hex
				}
			}
			i += used
		case n >= 30 && n <= 37:
			st.fg = indexHex(n - 30)
		case n >= 90 && n <= 97:
			st.fg = indexHex(n - 90 + 8)
		case n >= 40 && n <= 47:
			st.bg = indexHex(n - 40)
		case n >= 100 && n <= 107:
			st.bg = indexHex(n - 100 + 8)
		}
	}
	return st
}

// extended reads the tail of a 38/48 sequence: `5;N` for a palette index,
// `2;R;G;B` for 24-bit. Returns the hex and how many fields it consumed.
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
		for k := 0; k < 3; k++ {
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

func indexHex(n int) string { return theme.Hex(lipgloss.Color(strconv.Itoa(n))) }

// FromHTML reads what harness.HTML wrote.
//
// A real parser rather than a scanner tuned to our own output, because the
// point of the exercise is a format a person or a design tool can author. A
// reader that only understands markup it generated itself proves nothing about
// a hand-written mock.
func FromHTML(page string) (Grid, error) {
	doc, err := html.Parse(strings.NewReader(page))
	if err != nil {
		return nil, err
	}
	frame := find(doc, "tuikit-frame")
	if frame == nil {
		return nil, fmt.Errorf("no element with class tuikit-frame")
	}

	g := Grid{nil}
	walk(frame, style{}, &g)
	return g, nil
}

// find returns the first element carrying the given class.
func find(n *html.Node, class string) *html.Node {
	if n.Type == html.ElementNode {
		for _, a := range n.Attr {
			if a.Key == "class" {
				for _, f := range strings.Fields(a.Val) {
					if f == class {
						return n
					}
				}
			}
		}
	}
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		if found := find(c, class); found != nil {
			return found
		}
	}
	return nil
}

// walk appends every rune under n, carrying the enclosing style down. A
// newline in a text node starts a line, which is what white-space: pre means.
func walk(n *html.Node, st style, g *Grid) {
	for c := n.FirstChild; c != nil; c = c.NextSibling {
		switch c.Type {
		case html.TextNode:
			for _, r := range c.Data {
				if r == '\n' {
					*g = append(*g, nil)
					continue
				}
				line := len(*g) - 1
				(*g)[line] = append((*g)[line], st.cell(r))
			}
		case html.ElementNode:
			walk(c, merge(st, attr(c, "style")), g)
		}
	}
}

func attr(n *html.Node, key string) string {
	for _, a := range n.Attr {
		if a.Key == key {
			return a.Val
		}
	}
	return ""
}

// merge folds a CSS declaration block into the inherited style. Only the three
// properties harness.HTML emits are understood; anything else is ignored
// rather than guessed at, which is the same rule the writer follows.
func merge(st style, decls string) style {
	for _, d := range strings.Split(decls, ";") {
		prop, val, ok := strings.Cut(d, ":")
		if !ok {
			continue
		}
		prop, val = strings.TrimSpace(prop), strings.TrimSpace(val)
		switch prop {
		case "color":
			st.fg = val
		case "background", "background-color":
			st.bg = val
		case "font-weight":
			st.bold = val == "600" || val == "bold" || val == "700"
		}
	}
	return st
}

// Diff reports the first cell that differs, in the shape golden.go uses: what
// it was, what it is now, and where. A frame is a picture, and the first
// difference is the useful one.
func Diff(want, got Grid) string {
	if len(want) != len(got) {
		return fmt.Sprintf("  height: was %d lines, now %d", len(want), len(got))
	}
	for y := range want {
		if len(want[y]) != len(got[y]) {
			return fmt.Sprintf("  line %d: was %d cells, now %d\n    was %s\n    now %s",
				y+1, len(want[y]), len(got[y]), text(want[y]), text(got[y]))
		}
		for x := range want[y] {
			if want[y][x] == got[y][x] {
				continue
			}
			return fmt.Sprintf("  line %d, column %d\n    was %s\n    now %s\n    in line: %s",
				y+1, x+1, describe(want[y][x]), describe(got[y][x]), text(want[y]))
		}
	}
	return ""
}

func describe(c Cell) string {
	out := fmt.Sprintf("%q", string(c.Rune))
	if c.FG != "" {
		out += " fg " + c.FG
	} else {
		out += " fg default"
	}
	if c.BG != "" {
		out += " bg " + c.BG
	}
	if c.Bold {
		out += " bold"
	}
	return out
}

func text(cells []Cell) string {
	var b strings.Builder
	for _, c := range cells {
		b.WriteRune(c.Rune)
	}
	return strconv.Quote(b.String())
}
