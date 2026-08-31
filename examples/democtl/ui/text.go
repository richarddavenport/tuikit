package ui

import (
	"fmt"
	"strings"
	"time"

	"github.com/charmbracelet/lipgloss"
	"github.com/charmbracelet/x/ansi"
)

// Two ways to shorten a line, and they are not interchangeable.
//
// trim counts runes and is right for plain text. clip is for a line that
// already carries colour: trimming counts the escape sequences as width, so a
// coloured row gets cut short of the pane and, worse, cut mid-escape — leaving
// the colour switched on for everything after it. Anything that has been
// through a style is clipped.

// trim shortens plain text, marking the cut.
func trim(s string, w int) string {
	if w <= 0 {
		return ""
	}
	r := []rune(s)
	if len(r) <= w {
		return s
	}
	if w == 1 {
		return "…"
	}
	return string(r[:w-1]) + "…"
}

// clip shortens a line that may contain escape sequences, measuring what is
// visible and leaving the sequences intact.
func clip(s string, w int) string {
	if w <= 0 {
		return ""
	}
	if lipgloss.Width(s) <= w {
		return s
	}
	return ansi.Truncate(s, w, "…")
}

// pad extends plain text to exactly w columns.
func pad(s string, w int) string {
	s = trim(s, w)
	if gap := w - len([]rune(s)); gap > 0 {
		return s + strings.Repeat(" ", gap)
	}
	return s
}

// padVisible extends a possibly-coloured line to exactly w columns, measuring
// what shows rather than what is stored.
func padVisible(s string, w int) string {
	if lipgloss.Width(s) > w {
		return clip(s, w)
	}
	return s + strings.Repeat(" ", w-lipgloss.Width(s))
}

// rule is a horizontal line of the one box-drawing character available for it.
func rule(w int) string {
	if w < 0 {
		return ""
	}
	return strings.Repeat("─", w)
}

// wrap breaks text at word boundaries into lines of at most w columns. A word
// longer than the line is cut rather than allowed to overhang, because the
// caller asked for w and meant it.
func wrap(s string, w int) []string {
	if w <= 0 {
		return nil
	}
	var out []string
	line := ""
	for _, word := range strings.Fields(s) {
		switch {
		case line == "":
			line = word
		case len([]rune(line))+1+len([]rune(word)) <= w:
			line += " " + word
		default:
			out = append(out, line)
			line = word
		}
		for len([]rune(line)) > w {
			r := []rune(line)
			out = append(out, string(r[:w]))
			line = string(r[w:])
		}
	}
	if line != "" {
		out = append(out, line)
	}
	return out
}

// wrapFirst is the first wrapped line, for somewhere with room for one.
func wrapFirst(s string, w int) string {
	if lines := wrap(s, w); len(lines) > 0 {
		if len(lines) > 1 {
			return trim(lines[0], w)
		}
		return lines[0]
	}
	return ""
}

// ago is a duration as a person would say it. Coarse on purpose: "4m ago" is
// what the reader wants, and "4m13.204s ago" is the same fact made unreadable.
func ago(d time.Duration) string {
	switch {
	case d < time.Minute:
		return fmt.Sprintf("%ds ago", int(d.Seconds()))
	case d < time.Hour:
		return fmt.Sprintf("%dm ago", int(d.Minutes()))
	default:
		return fmt.Sprintf("%dh ago", int(d.Hours()))
	}
}

// took is how long a step ran, at the precision a step list can use.
func took(d time.Duration) string {
	if d < time.Second {
		return fmt.Sprintf("%dms", d.Milliseconds())
	}
	return fmt.Sprintf("%.1fs", d.Seconds())
}

// composite puts over on top of under, starting at column left, and keeps what
// under had on either side.
//
// Replacing the whole line instead — which is what this did first — wipes the
// frame the modal is floating over, so a dialog appears to punch a hole through
// the panes rather than sit on them. Two captured frames made that obvious and
// no assertion would have.
func composite(under, over string, left, width int) string {
	over = clip(over, max(0, width-left))
	prefix := ansi.Truncate(under, left, "")
	if w := lipgloss.Width(prefix); w < left {
		prefix += strings.Repeat(" ", left-w)
	}
	rest := ansi.TruncateLeft(under, left+lipgloss.Width(over), "")
	return prefix + over + rest
}
