package comp

import "testing"

func ring() *Focus {
	return &Focus{Ring: []Name{"projects", "clusters", "pods"}}
}

// The zero value is usable, and lands on the head of the ring rather than
// nowhere — so the first key press has somewhere to go.
func TestTheZeroFocusLandsOnTheFirstPane(t *testing.T) {
	f := ring()
	if got := f.Current(); got != "projects" {
		t.Errorf("current is %q before anything was set", got)
	}
	if !f.Is("projects") || f.Is("pods") {
		t.Error("Is disagrees with Current")
	}
}

// A Focus with no ring at all does not panic, which is the state a screen is in
// before it has decided what its panes are.
func TestAnEmptyRingIsQuiet(t *testing.T) {
	var f Focus
	f.Next()
	f.Prev()
	if f.Current() != "" || f.Is("anything") {
		t.Error("an empty focus claims to have one")
	}
	if f.Set("pods") {
		t.Error("Set took a name with no ring to put it in")
	}
}

func TestNextAndPrevWrap(t *testing.T) {
	f := ring()
	f.Prev()
	if got := f.Current(); got != "pods" {
		t.Errorf("prev from the head gave %q, want the tail", got)
	}
	f.Next()
	if got := f.Current(); got != "projects" {
		t.Errorf("next from the tail gave %q, want the head", got)
	}
}

// Set refuses a name outside the ring, which is what makes it safe to call with
// whatever a click landed on.
func TestSetRefusesSomethingNotInTheRing(t *testing.T) {
	f := ring()
	if f.Set("footer") {
		t.Error("Set took a region that is not a pane")
	}
	if got := f.Current(); got != "projects" {
		t.Errorf("a refused Set moved the focus to %q", got)
	}
	if !f.Set("pods") || f.Current() != "pods" {
		t.Error("Set refused a name that is in the ring")
	}
}

// A click lands on a ROW, and the ring holds the pane. On resolves one to the
// other so no tool has to write that mapping.
func TestAClickOnARowFocusesItsPane(t *testing.T) {
	f := ring()
	if !f.On(Region("pods").At(3)) {
		t.Error("a click on pods.row did not move the focus")
	}
	if got := f.Current(); got != "pods" {
		t.Errorf("focus is %q after clicking a pod row", got)
	}
}

// Matching is at a dot boundary, so a region that merely starts with a pane's
// name does not steal its focus.
func TestOnMatchesOnADotBoundary(t *testing.T) {
	// Two panes, focused on the second, so a match is visible as a MOVE. With
	// one pane the focus is already there and On correctly reports nothing
	// moved, which says nothing about whether it matched.
	f := &Focus{Ring: []Name{"pod", "log"}}
	f.Set("log")

	if f.On(Region("podsummary")) {
		t.Error("podsummary was treated as part of pod")
	}
	if got := f.Current(); got != "log" {
		t.Errorf("podsummary moved the focus to %q", got)
	}
	if !f.On(Region("pod.row")) {
		t.Error("pod.row was not matched to pod")
	}
	if got := f.Current(); got != "pod" {
		t.Errorf("focus is %q after clicking pod.row", got)
	}
}

// On reports whether the focus MOVED, so a caller can tell a click that changed
// panes from a click inside the pane it was already on.
func TestOnSaysWhetherItMoved(t *testing.T) {
	f := ring()
	f.Set("pods")
	if f.On(Region("pods").At(1)) {
		t.Error("clicking the focused pane reported a move")
	}
	if !f.On(Region("clusters").At(0)) {
		t.Error("clicking another pane reported no move")
	}
}

// A click on something outside the ring leaves the focus alone.
func TestOnIgnoresARegionOutsideTheRing(t *testing.T) {
	f := ring()
	f.Set("pods")
	if f.On(Region("footer")) {
		t.Error("the footer took the focus")
	}
	if got := f.Current(); got != "pods" {
		t.Errorf("focus is %q after clicking the footer", got)
	}
}

// The ring can change between frames — a pane that is not on screen should not
// be tabbed into — and the focus survives as long as its name is still there.
func TestTheFocusSurvivesTheRingChanging(t *testing.T) {
	f := ring()
	f.Set("pods")
	f.Ring = []Name{"projects", "pods"}
	if got := f.Current(); got != "pods" {
		t.Errorf("focus is %q after a pane was removed from the ring", got)
	}
	f.Next()
	if got := f.Current(); got != "projects" {
		t.Errorf("next from the new tail gave %q", got)
	}
}

// A focused pane that LEAVES the ring hands the keyboard back rather than
// leaving it pointing at something no longer drawn.
func TestAFocusDroppedFromTheRingFallsBack(t *testing.T) {
	f := ring()
	f.Set("pods")
	f.Ring = []Name{"projects", "clusters"}
	if got := f.Current(); got != "projects" {
		t.Errorf("focus is %q after its pane left the ring", got)
	}
}

// Indices were refused for this reason: a ring that gains a pane at the front
// would move an index-based focus to a different pane without erroring.
func TestAPaneAddedInFrontDoesNotStealTheFocus(t *testing.T) {
	f := ring()
	f.Set("clusters")
	f.Ring = append([]Name{"accounts"}, f.Ring...)
	if got := f.Current(); got != "clusters" {
		t.Errorf("focus moved to %q when a pane was inserted before it", got)
	}
}
