package comp

import (
	"reflect"
	"strings"
	"testing"
)

// The zero value is unsorted, and unsorted is distinguishable from "sorted by
// column 0". That distinction is the reason for the unexported field.
func TestTheZeroSortIsUnsorted(t *testing.T) {
	var s Sort
	if s.On() {
		t.Error("a fresh Sort claims to be sorting")
	}
	if _, _, ok := s.By(); ok {
		t.Error("By reports a column before one was chosen")
	}
	s.Toggle(0)
	if col, desc, ok := s.By(); !ok || col != 0 || desc {
		t.Errorf("after sorting by column 0: %d %v %v", col, desc, ok)
	}
}

// Click once for ascending, again for descending. A third click keeps the
// direction rather than clearing, because a table that intermittently forgets
// its order is worse than one holding an order you did not want.
func TestToggleReversesTheSameColumn(t *testing.T) {
	var s Sort
	s.Toggle(2)
	if _, desc, _ := s.By(); desc {
		t.Error("the first toggle was descending")
	}
	s.Toggle(2)
	if _, desc, _ := s.By(); !desc {
		t.Error("the second toggle did not reverse")
	}
	s.Toggle(2)
	if _, _, ok := s.By(); !ok {
		t.Error("the third toggle cleared the sort")
	}
}

// A different column starts ascending again, rather than inheriting the last
// column's direction.
func TestADifferentColumnStartsAscending(t *testing.T) {
	var s Sort
	s.Toggle(1)
	s.Toggle(1) // now descending
	s.Toggle(3)
	if col, desc, _ := s.By(); col != 3 || desc {
		t.Errorf("moving to column 3 gave %d desc=%v", col, desc)
	}
}

func TestClearReturnsToTheToolsOrder(t *testing.T) {
	var s Sort
	s.Toggle(1)
	s.Clear()
	if s.On() {
		t.Error("clear left it sorting")
	}
	idx := s.Apply(3, func(a, b, _ int) bool { return a > b })
	if !reflect.DeepEqual(idx, []int{0, 1, 2}) {
		t.Errorf("an unsorted Apply reordered to %v", idx)
	}
}

// Apply hands back indices so the caller's rows are never copied, and it
// reverses for Desc rather than making the caller write the comparison twice.
func TestApplyOrdersAndReverses(t *testing.T) {
	names := []string{"charlie", "alpha", "bravo"}
	less := func(a, b, _ int) bool { return names[a] < names[b] }

	var s Sort
	s.Toggle(0)
	if got := s.Apply(3, less); !reflect.DeepEqual(got, []int{1, 2, 0}) {
		t.Errorf("ascending gave %v", got)
	}
	s.Toggle(0)
	if got := s.Apply(3, less); !reflect.DeepEqual(got, []int{0, 2, 1}) {
		t.Errorf("descending gave %v", got)
	}
}

// Stable, so rows that compare equal keep the order the tool built them in.
// Without that, re-sorting by a column with ties shuffles the rows under the
// reader every time the data refreshes.
func TestApplyIsStable(t *testing.T) {
	state := []string{"ready", "ready", "ready", "failed"}
	less := func(a, b, _ int) bool { return state[a] < state[b] }

	var s Sort
	s.Toggle(0)
	if got := s.Apply(4, less); !reflect.DeepEqual(got, []int{3, 0, 1, 2}) {
		t.Errorf("ties were reordered: %v", got)
	}
}

// A nil comparison is a caller that has not decided how to order this column
// yet. It leaves the rows alone rather than panicking.
func TestApplyWithNoComparisonIsQuiet(t *testing.T) {
	var s Sort
	s.Toggle(0)
	if got := s.Apply(3, nil); !reflect.DeepEqual(got, []int{0, 1, 2}) {
		t.Errorf("a nil comparison gave %v", got)
	}
}

// The arrow marks the sorted column and nothing else.
func TestHeaderMarksOnlyTheSortedColumn(t *testing.T) {
	c := NewCanvas(40, 4)
	labels := []string{"name", "state", "age"}

	var s Sort
	if got := s.Header(c, labels); !reflect.DeepEqual(got, labels) {
		t.Errorf("an unsorted table marked a header: %v", got)
	}
	s.Toggle(1)
	got := s.Header(c, labels)
	if got[1] != "state"+c.Chrome().SortAsc {
		t.Errorf("the sorted header is %q", got[1])
	}
	if got[0] != "name" || got[2] != "age" {
		t.Errorf("another header was marked: %v", got)
	}
	if labels[1] != "state" {
		t.Error("Header modified the caller's slice")
	}
}

// The arrow comes from the chrome, so a tool that changed it changed it
// everywhere rather than in the one place it remembered.
func TestTheArrowIsTheChromes(t *testing.T) {
	ch := NewCanvas(10, 2).Chrome()
	ch.SortAsc, ch.SortDesc = "^", "v"
	c := NewCanvas(10, 2).WithChrome(ch)

	var s Sort
	s.Toggle(0)
	if got := s.Marker(c, 0); got != "^" {
		t.Errorf("marker is %q", got)
	}
	s.Toggle(0)
	if got := s.Marker(c, 0); got != "v" {
		t.Errorf("reversed marker is %q", got)
	}
	if got := s.Marker(c, 1); got != "" {
		t.Errorf("an unsorted column has marker %q", got)
	}
}

// A column index that is not in the header does not panic, which is what a
// stale sort looks like after the columns changed.
func TestHeaderOutOfRangeIsQuiet(t *testing.T) {
	c := NewCanvas(40, 4)
	var s Sort
	s.Toggle(9)
	got := s.Header(c, []string{"a", "b"})
	if strings.Join(got, "") != "ab" {
		t.Errorf("a stale sort column changed the header: %v", got)
	}
}
