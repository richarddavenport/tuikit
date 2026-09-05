package comp

import (
	"reflect"
	"testing"
)

// A small filesystem, the shape lazygit's filetree and yazi both flatten to.
func nodes() []Node {
	return []Node{
		{0, "src"},          // 0
		{1, "src/app"},      // 1
		{2, ""},             // 2  src/app/main.go
		{2, ""},             // 3  src/app/run.go
		{1, "src/lib"},      // 4
		{2, ""},             // 5  src/lib/util.go
		{0, "docs"},         // 6
		{1, ""},             // 7  docs/readme.md
		{0, ""},             // 8  go.mod
	}
}

func TestEverythingIsVisibleUntilSomethingIsCollapsed(t *testing.T) {
	tr := &Tree{}
	if got := tr.Visible(nodes()); len(got) != 9 {
		t.Errorf("%d of 9 rows visible with nothing collapsed", len(got))
	}
}

// Collapsing a directory hides its whole subtree, and stops at the first row
// shallow enough not to belong to it.
func TestCollapsingHidesTheSubtreeAndNoMore(t *testing.T) {
	tr := &Tree{}
	tr.Toggle("src")

	want := []int{0, 6, 7, 8} // src, then docs and its child, then go.mod
	if got := tr.Visible(nodes()); !reflect.DeepEqual(got, want) {
		t.Errorf("visible %v, want %v", got, want)
	}
}

// A collapse inside a collapse is already hidden, and re-opening the outer one
// must not open the inner.
func TestANestedCollapseSurvivesTheOuterOne(t *testing.T) {
	tr := &Tree{}
	tr.Toggle("src/app")
	tr.Toggle("src")

	if got := tr.Visible(nodes()); !reflect.DeepEqual(got, []int{0, 6, 7, 8}) {
		t.Fatalf("visible %v", got)
	}

	tr.Toggle("src") // open the outer one again
	want := []int{0, 1, 4, 5, 6, 7, 8}
	if got := tr.Visible(nodes()); !reflect.DeepEqual(got, want) {
		t.Errorf("after reopening, visible %v, want %v — the inner collapse was lost", got, want)
	}
}

// Collapsing a leaf hides nothing. A tool that lets one be toggled should not
// get a subtree that swallows its siblings.
func TestCollapsingALeafHidesNothing(t *testing.T) {
	ns := []Node{{0, "a"}, {1, "a/leaf"}, {1, "a/other"}}
	tr := &Tree{}
	tr.Toggle("a/leaf")

	if got := tr.Visible(ns); !reflect.DeepEqual(got, []int{0, 1, 2}) {
		t.Errorf("collapsing a leaf hid %v", got)
	}
}

// The key is the tool's identity, so a row appearing above a collapsed one does
// not open a different one. Indices move; keys do not.
func TestCollapseSurvivesRowsMovingAboveIt(t *testing.T) {
	tr := &Tree{}
	tr.Toggle("docs")

	before := tr.Visible(nodes())
	withExtra := append([]Node{{0, "new"}}, nodes()...)
	after := tr.Visible(withExtra)

	// docs was index 6 and is now 7; its child must still be hidden.
	if len(after) != len(before)+1 {
		t.Errorf("inserting a row above changed what was collapsed: %d then %d", len(before), len(after))
	}
	for _, i := range after {
		if withExtra[i].Depth == 1 && i == 8 {
			t.Error("docs' child became visible when a row moved above it")
		}
	}
}

// HasChildren is read off the depths, so it cannot disagree with them.
func TestHasChildrenReadsTheDepths(t *testing.T) {
	ns := nodes()
	for i, want := range map[int]bool{0: true, 1: true, 2: false, 5: false, 6: true, 8: false} {
		if got := HasChildren(ns, i); got != want {
			t.Errorf("HasChildren(%d) = %v, want %v", i, got, want)
		}
	}
	if HasChildren(ns, len(ns)-1) {
		t.Error("the last row reported children")
	}
}

// Collapse shuts everything with children and nothing without, so no marker is
// drawn beside a row that has nothing to open.
func TestCollapseAllShutsOnlyWhatHasChildren(t *testing.T) {
	ns := nodes()
	tr := &Tree{}
	tr.Collapse(ns)

	if got := tr.Visible(ns); !reflect.DeepEqual(got, []int{0, 6, 8}) {
		t.Errorf("collapse-all left %v visible, want the three roots", got)
	}
	for i, n := range ns {
		if n.Key == "" {
			continue
		}
		if tr.IsCollapsed(n.Key) != HasChildren(ns, i) {
			t.Errorf("%q collapsed=%v but HasChildren=%v", n.Key, tr.IsCollapsed(n.Key), HasChildren(ns, i))
		}
	}
	tr.Expand()
	if len(tr.Visible(ns)) != len(ns) {
		t.Error("expand-all did not reopen everything")
	}
}

// The zero value works, which is what makes it usable as a field rather than
// something a tool has to remember to construct.
func TestTheZeroTreeIsUsable(t *testing.T) {
	var tr Tree
	if got := tr.Visible(nodes()); len(got) != 9 {
		t.Errorf("the zero tree hid rows: %v", got)
	}
	tr.Toggle("src")
	if len(tr.Visible(nodes())) != 4 {
		t.Error("the zero tree could not be toggled")
	}
}
