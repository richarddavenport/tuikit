package comp

import (
	"strings"

	"github.com/charmbracelet/x/ansi"

	"github.com/richarddavenport/tuikit/theme"
)

// Text operations components need, measured in COLUMNS.
//
// These are here rather than in each tool because every one of the four wrote
// its own, and the versions that counted bytes or runes were the source of a
// header off the side of the screen and a tab strip cut mid-escape. There is no
// escape sequence to miscount any more — the canvas takes plain text — but a
// CJK name is still two columns per rune, so the counting has to be right.
//
// Truncating and wrapping are not layout: the canvas clips, which is a
// different thing. Clipping is what happens when a component draws past its
// rect, silently and correctly. Truncation is a component DECIDING that a name
// is too long and saying so with an ellipsis.

// Truncate shortens text to w columns, ending with an ellipsis when it had to
// cut. The ellipsis is a column of its own, so the result is never wider than
// asked for.
func Truncate(s string, w int) string { return truncate(s, w, theme.DefaultChrome.Ellipsis) }

// truncate is Truncate in a given chrome. Components use this, so a tool that
// changed its ellipsis changed it everywhere rather than in the one place it
// remembered.
func truncate(s string, w int, ellipsis string) string {
	if w <= 0 {
		return ""
	}
	if Width(s) <= w {
		return s
	}
	if w <= Width(ellipsis) {
		return ellipsis
	}
	return ansi.Truncate(s, w, ellipsis)
}

// Wrap breaks text into lines of at most w columns, at spaces where it can and
// mid-word when a single word is longer than the line.
//
// A newline in the input is KEPT, and a blank line stays blank. It used to go
// through strings.Fields, which collapses every kind of whitespace equally — so
// a two-paragraph confirmation body came out as one run-on paragraph, and the
// author who put the break there had no way to tell the difference between
// "wrapped" and "ignored".
func Wrap(s string, w int) []string {
	if w <= 0 {
		return nil
	}
	var out []string
	for i, para := range strings.Split(s, "\n") {
		if i > 0 && strings.TrimSpace(para) == "" {
			out = append(out, "")
			continue
		}
		out = append(out, wrapLine(para, w)...)
	}
	return out
}

// wrapLine wraps one line that has no breaks of its own.
func wrapLine(s string, w int) []string {
	var out []string
	line := ""
	for _, word := range strings.Fields(s) {
		switch {
		case line == "":
			line = word
		case Width(line)+1+Width(word) <= w:
			line += " " + word
		default:
			out = append(out, line)
			line = word
		}
		for Width(line) > w {
			head := ansi.Truncate(line, w, "")
			out = append(out, head)
			line = strings.TrimPrefix(line, head)
		}
	}
	if line != "" {
		out = append(out, line)
	}
	return out
}

// WrapFirst is the first wrapped line, for somewhere with room for one.
//
// It stops at the line break and says nothing about the rest, which is worth
// questioning: a note that reads "no suitable node: memory reservation exceeds"
// looks like a complete sentence rather than a truncated one, and everywhere
// else in tuikit an invisible remainder is treated as a bug. Left as it is for
// now because changing it changes what democtl draws, and that belongs in its
// own commit against its own goldens rather than inside a port.
func WrapFirst(s string, w int) string {
	if lines := Wrap(s, w); len(lines) > 0 {
		return lines[0]
	}
	return ""
}
