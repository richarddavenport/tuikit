package comp

import "testing"

// Issue 72: the range works over anything with a cursor, not only a List.
//
// A tree, here — node indices into a flattened hierarchy, which is what
// gowid's copymodetree needs and what ours could not do.
func TestARangeWorksOverAnythingWithACursor(t *testing.T) {
	var sel Range
	if _, _, ok := sel.Span(0); ok {
		t.Error("a fresh Range reported a selection from row 0 — the zero value must mean none")
	}
	if sel.Anchored() {
		t.Error("a fresh Range says it is anchored")
	}

	sel.Start(4)
	sel.Start(9) // every key of a drag calls it; only the first counts
	if !sel.Anchored() {
		t.Fatal("Start did not anchor")
	}
	if lo, hi, _ := sel.Span(7); lo != 4 || hi != 7 {
		t.Errorf("span is %d..%d, want 4..7 — the second Start moved the anchor", lo, hi)
	}

	// Turn round and go past the start. This is the case a pair of bounds gets
	// wrong, and the reason the anchor is kept rather than both ends.
	if lo, hi, _ := sel.Span(1); lo != 1 || hi != 4 {
		t.Errorf("dragging back past the anchor gave %d..%d, want 1..4", lo, hi)
	}

	sel.Clear()
	if _, _, ok := sel.Span(7); ok {
		t.Error("Clear left the selection behind")
	}
}

// It composes with Marks: the span picks the keys, the set remembers them.
func TestARangeFeedsTheMarkSet(t *testing.T) {
	keys := []string{"a", "b", "c", "d", "e"}
	var sel Range
	sel.Start(1)
	lo, hi, ok := sel.Span(3)
	if !ok {
		t.Fatal("no span")
	}

	var m Marks
	m.Span(keys[lo : hi+1])
	if got := m.Len(); got != 3 {
		t.Errorf("marked %d keys, want 3", got)
	}
	if !m.Has("b") || !m.Has("d") {
		t.Errorf("the span marked %v, want b..d", m.Keys())
	}
}
