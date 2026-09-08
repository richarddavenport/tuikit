package comp

// Focus is which region has the keyboard.
//
// # Where this came from
//
// Five tools in the rebuild survey wrote this, and every one of them is a
// multi-pane screen: lazygit's `pkg/gui/context/`, termshark's five widgets
// (framefocus, trackfocus, renderfocused, keepselected, enableselected),
// dive's `ui/v1/app/controller.go`, gcpeasy's `focus tuiPanel`, and the deploy
// tool's `focus int`.
//
// It was refused once on a count of one, because only the deploy tool of the
// four private tools had it. That was the wrong pool to count in: three of
// those four have a single list per screen, which is the shape that does not
// need this.
//
// # What was already here, and what was not
//
// [List], [Pane], [Tabs] and [Viewer] all have a Focused field, so every
// component can be TOLD. Nothing held the answer. This holds it.
//
// # Why a name rather than an index
//
// A ring of indices breaks the moment a pane is added, removed or hidden — the
// focus silently moves to a different pane, which is the same failure shape as
// an owner ID that means "row 3 of the screen". A [Name] survives the ring
// changing under it, and it is the name the region is already drawn under, so
// a click resolves against the same value.
//
// # What it does not do
//
// It does not own the components and it does not draw. termshark's version is
// a wrapper widget that frames whatever has focus; in immediate mode that
// would be a field, not a wrapper, and Pane.Focused already is one. The state
// is the part worth extracting.
type Focus struct {
	// Ring is the order [Focus.Next] moves through. Usually declared once
	// beside the region names.
	//
	// A tool may change it between frames — a pane that is not on screen this
	// time should not be tabbed into — and the current focus survives as long
	// as its name is still somewhere in the new ring.
	Ring []Name

	// at is the focused name. The zero Focus has none, and Is reports false for
	// everything until something is focused, which is what makes a screen with
	// nothing selected drawable rather than a special case.
	at Name
}

// Current is the focused region, or "" if nothing is.
func (f *Focus) Current() Name { return f.resolve() }

// Is reports whether a region has the keyboard.
//
// The call site is `l.Focused = f.Is(regPods)`, once per component per frame,
// which is the whole integration.
func (f *Focus) Is(n Name) bool { return n != "" && f.resolve() == n }

// resolve is the current name, falling back to the head of the ring.
//
// Falling back rather than staying empty, so a tool that never called Set
// still has somewhere for the first key press to go — and so a focused pane
// that has left the ring does not leave the keyboard pointing at nothing.
func (f *Focus) resolve() Name {
	if f.at != "" && f.inRing(f.at) {
		return f.at
	}
	if len(f.Ring) == 0 {
		return f.at
	}
	return f.Ring[0]
}

func (f *Focus) inRing(n Name) bool {
	for _, r := range f.Ring {
		if r == n {
			return true
		}
	}
	return false
}

// Set puts the keyboard on a region, and reports whether it took.
//
// Refuses a name that is not in the ring, which is what makes it safe to call
// with whatever a click landed on: a click on the footer does not move the
// focus, and the caller does not have to know that.
func (f *Focus) Set(n Name) bool {
	if !f.inRing(n) {
		return false
	}
	f.at = n
	return true
}

// On puts the keyboard on whatever region an ID belongs to.
//
// A click lands on `pods.row[3]`, and the ring holds `pods`. Rather than make
// every tool map one to the other, this matches a ring entry that is a prefix
// of the ID's name at a dot boundary — so `pods.row` finds `pods`, and `podsx`
// does not.
//
// Returns whether the focus moved, so a caller can tell a click that changed
// panes from one that did not.
func (f *Focus) On(id ID) bool {
	for _, r := range f.Ring {
		if id.Name == r || (len(id.Name) > len(r) && id.Name[:len(r)] == r && id.Name[len(r)] == '.') {
			if f.resolve() == r {
				return false
			}
			f.at = r
			return true
		}
	}
	return false
}

// Next and Prev move round the ring.
//
// They wrap, because a ring that stops at the end is a ring you have to press
// shift-tab to get out of, and every tool in the survey wraps.
func (f *Focus) Next() { f.move(1) }

// Prev moves the other way.
func (f *Focus) Prev() { f.move(-1) }

func (f *Focus) move(by int) {
	if len(f.Ring) == 0 {
		return
	}
	cur := f.resolve()
	for i, r := range f.Ring {
		if r != cur {
			continue
		}
		f.at = f.Ring[((i+by)%len(f.Ring)+len(f.Ring))%len(f.Ring)]
		return
	}
	f.at = f.Ring[0]
}
