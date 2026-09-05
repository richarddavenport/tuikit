package comp

// rangeSel is the anchor half of a range selection.
//
// Shared by [List] and [Viewer] rather than written twice. The invariant is
// small and easy to get wrong in the same way both times: the zero value must
// mean NO selection, so a fresh component cannot report a range from row 0.
//
// Two fields rather than an anchor of -1, because -1 is not a zero value and a
// struct that has to be constructed to be correct is one a caller will forget
// to construct.
//
// Kept beside the cursor rather than as a pair of bounds, because a range grows
// from one end and a pair does not remember which end that was: extend, turn
// round, and extend past the start is the case that gets it wrong.
type rangeSel struct {
	anchor   int
	anchored bool
}

// start anchors at cursor if nothing is anchored yet.
//
// The first extend from no selection anchors where the cursor IS, so a reader
// who presses shift-down once has selected two rows rather than none — which is
// what every other list in the world does.
func (s *rangeSel) start(cursor int) {
	if !s.anchored {
		s.anchor, s.anchored = cursor, true
	}
}

// span is the selection low to high, and whether there is one.
//
// Ordered, so a caller acting on it never has to ask which way the reader
// dragged.
func (s *rangeSel) span(cursor int) (lo, hi int, ok bool) {
	if !s.anchored {
		return 0, 0, false
	}
	if s.anchor <= cursor {
		return s.anchor, cursor, true
	}
	return cursor, s.anchor, true
}

// clear drops the selection, leaving the cursor.
func (s *rangeSel) clear() { s.anchor, s.anchored = 0, false }
