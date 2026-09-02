package comp

import "github.com/charmbracelet/lipgloss"

// Highlight turns a string and a set of rune indices into Spans, so a row can
// mark the letters a search actually matched.
//
// This is the other half of [fuzzy.Match] returning positions. A ranked list
// whose order you have to take on trust is a worse list than an unranked one:
// the reader cannot tell whether the top row is first because it is the best
// match or because it was first already. Underlining what matched turns the
// order into something they can check at a glance.
//
// Runs of the same style are merged, so a five-letter match is one span rather
// than five — which matters because the canvas groups adjacent cells sharing a
// style pointer into a single escape sequence.
func Highlight(text string, at []int, base, hit *lipgloss.Style) []Segment {
	runes := []rune(text)
	if len(runes) == 0 {
		return nil
	}
	marked := make([]bool, len(runes))
	for _, i := range at {
		if i >= 0 && i < len(marked) {
			marked[i] = true
		}
	}

	var spans []Segment
	start := 0
	for i := 1; i <= len(runes); i++ {
		if i < len(runes) && marked[i] == marked[start] {
			continue
		}
		style := base
		if marked[start] {
			style = hit
		}
		spans = append(spans, Segment{Text: string(runes[start:i]), Style: style})
		start = i
	}
	return spans
}
