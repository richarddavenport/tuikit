package comp_test

import (
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/fuzzy"
)

var (
	plain = lipgloss.NewStyle()
	bold  = lipgloss.NewStyle().Bold(true)
)

func TestHighlightSplitsOnTheMatch(t *testing.T) {
	spans := comp.Highlight("disk usage", []int{0, 1, 2, 3}, &plain, &bold)
	if len(spans) != 2 {
		t.Fatalf("got %d spans, want 2: %v", len(spans), spans)
	}
	if spans[0].Text != "disk" || spans[0].Style != &bold {
		t.Errorf("first span is %q styled=%v", spans[0].Text, spans[0].Style == &bold)
	}
	if spans[1].Text != " usage" || spans[1].Style != &plain {
		t.Errorf("second span is %q", spans[1].Text)
	}
}

// TestHighlightMergesRuns. A five-letter match must be one span, not five: the
// canvas groups adjacent cells sharing a style POINTER into one escape
// sequence, so unmerged spans are five times the bytes for the same picture.
func TestHighlightMergesRuns(t *testing.T) {
	spans := comp.Highlight("abcdef", []int{1, 2, 3, 4}, &plain, &bold)
	if len(spans) != 3 {
		t.Fatalf("got %d spans, want 3", len(spans))
	}
	if spans[1].Text != "bcde" {
		t.Errorf("the run came out as %q", spans[1].Text)
	}
}

func TestHighlightScattered(t *testing.T) {
	spans := comp.Highlight("abcdef", []int{0, 2, 4}, &plain, &bold)
	want := []string{"a", "b", "c", "d", "e", "f"}
	if len(spans) != len(want) {
		t.Fatalf("got %d spans, want %d", len(spans), len(want))
	}
	for i, s := range spans {
		if s.Text != want[i] {
			t.Errorf("span %d is %q, want %q", i, s.Text, want[i])
		}
	}
}

func TestHighlightEdges(t *testing.T) {
	if got := comp.Highlight("", []int{0}, &plain, &bold); got != nil {
		t.Errorf("empty text gave %v, want nil", got)
	}
	// No matches: one span, all base.
	if got := comp.Highlight("abc", nil, &plain, &bold); len(got) != 1 || got[0].Style != &plain {
		t.Errorf("unmatched text gave %v", got)
	}
	// Out-of-range indices are ignored rather than panicking: they come from a
	// match against text that has since been re-rendered.
	if got := comp.Highlight("abc", []int{99, -1}, &plain, &bold); len(got) != 1 {
		t.Errorf("stray indices gave %v", got)
	}
}

// TestHighlightUsesRuneIndices end to end with the matcher, since the two only
// agree if both count in runes.
func TestHighlightUsesRuneIndices(t *testing.T) {
	const text = "un néant"
	m, ok := fuzzy.Match("né", text)
	if !ok {
		t.Fatal("no match")
	}
	spans := comp.Highlight(text, m.At, &plain, &bold)
	var marked string
	for _, s := range spans {
		if s.Style == &bold {
			marked += s.Text
		}
	}
	if marked != "né" {
		t.Errorf("highlighted %q, want %q", marked, "né")
	}
}
