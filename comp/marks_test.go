package comp

import (
	"reflect"
	"testing"
)

// The zero value works, which is what lets a tool declare one as a field and
// not think about it again.
func TestTheZeroMarksIsUsable(t *testing.T) {
	var m Marks
	if m.Len() != 0 || m.Has("a") || m.Keys() != nil && len(m.Keys()) != 0 {
		t.Error("a fresh Marks claims to hold something")
	}
	m.Toggle("a")
	if !m.Has("a") {
		t.Error("toggling a fresh Marks did nothing")
	}
}

func TestToggleGoesBothWays(t *testing.T) {
	var m Marks
	m.Toggle("a")
	m.Toggle("a")
	if m.Has("a") || m.Len() != 0 {
		t.Error("toggling twice left it marked")
	}
}

// An empty key is not a mark. A row with no identity — a heading, a blank —
// must not become selectable by accident.
func TestAnEmptyKeyIsNotAMark(t *testing.T) {
	var m Marks
	m.Toggle("")
	m.Add("")
	if m.Len() != 0 {
		t.Errorf("an empty key was marked, giving %v", m.Keys())
	}
}

// Sorted, so an action on several things happens in the same order twice.
func TestKeysAreSorted(t *testing.T) {
	var m Marks
	m.Add("charlie", "alpha", "bravo")
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"alpha", "bravo", "charlie"}) {
		t.Errorf("keys are %v", got)
	}
}

// The method the component exists for: one code path for "act on the marks" and
// "act on the one under the cursor".
func TestActingFallsBackToTheCursor(t *testing.T) {
	var m Marks
	if got := m.Acting("web-1"); !reflect.DeepEqual(got, []string{"web-1"}) {
		t.Errorf("with nothing marked, acting on %v", got)
	}
	m.Add("worker-2", "worker-3")
	if got := m.Acting("web-1"); !reflect.DeepEqual(got, []string{"worker-2", "worker-3"}) {
		t.Errorf("with marks, acting on %v — the cursor should not be in it", got)
	}
}

// Nothing marked and no cursor is an ordinary state: an empty list, not an
// error.
func TestActingOnNothingIsNil(t *testing.T) {
	var m Marks
	if got := m.Acting(""); got != nil {
		t.Errorf("acting on %v with nothing at all", got)
	}
}

// A mark survives the rows moving, which is the whole reason for a key.
func TestMarksSurviveAFilterMovingTheRows(t *testing.T) {
	var m Marks
	before := []string{"alpha", "bravo", "charlie", "delta"}
	m.Add(before[2])

	// A filter removes the first row. Every index has moved by one.
	after := []string{"bravo", "charlie", "delta"}
	if !m.Has(after[1]) {
		t.Error("the mark moved to a different row when the rows shifted")
	}
	if m.Has(after[0]) || m.Has(after[2]) {
		t.Error("something else became marked")
	}
}

// yazi's gesture: a range is dragged, then committed.
func TestSpanCommitsARange(t *testing.T) {
	var m Marks
	rows := []string{"a", "b", "c", "d", "e"}
	m.Span(rows[1:4])
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"b", "c", "d"}) {
		t.Errorf("the span marked %v", got)
	}
}

// k9s's gesture: no anchor at all. Scan back for the nearest mark and fill.
func TestFillReachesBackToTheNearestMark(t *testing.T) {
	var m Marks
	rows := []string{"a", "b", "c", "d", "e"}
	m.Add("b")
	if !m.Fill(rows, 4) {
		t.Fatal("fill found nothing to reach back to")
	}
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"b", "c", "d", "e"}) {
		t.Errorf("fill marked %v", got)
	}
}

// And forwards when there is nothing behind, which is the half that makes it
// work on the first row.
func TestFillReachesForwardWhenNothingIsBehind(t *testing.T) {
	var m Marks
	rows := []string{"a", "b", "c", "d"}
	m.Add("c")
	if !m.Fill(rows, 0) {
		t.Fatal("fill found nothing ahead")
	}
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"a", "b", "c"}) {
		t.Errorf("fill marked %v", got)
	}
}

// Nothing marked at all means there is nothing to fill from, and saying so lets
// a tool fall back to marking the cursor instead of filling silently.
func TestFillWithNoMarksReportsSo(t *testing.T) {
	var m Marks
	if m.Fill([]string{"a", "b"}, 1) {
		t.Error("fill claimed to have filled from nothing")
	}
	if m.Len() != 0 {
		t.Errorf("fill marked %v anyway", m.Keys())
	}
}

// A cursor past the end of the rows is what a stale index looks like, and it
// must not panic.
func TestFillOutOfRangeIsQuiet(t *testing.T) {
	var m Marks
	m.Add("a")
	if m.Fill([]string{"a", "b"}, 9) || m.Fill([]string{"a", "b"}, -1) || m.Fill(nil, 0) {
		t.Error("fill acted on an index it does not have")
	}
}

func TestClearAndRemove(t *testing.T) {
	var m Marks
	m.Add("a", "b", "c")
	m.Remove("b")
	if got := m.Keys(); !reflect.DeepEqual(got, []string{"a", "c"}) {
		t.Errorf("after removing b: %v", got)
	}
	m.Clear()
	if m.Len() != 0 {
		t.Errorf("clear left %v", m.Keys())
	}
}
