package comp

// Range is the anchor half of a range selection, over anything with a cursor.
//
// # What it is for
//
// A reader holding shift and moving, over a list, a document, a table or a
// tree. The gesture is identical in all four and the thing underneath is not,
// which is why this holds an ANCHOR and a caller supplies the cursor rather
// than this owning both.
//
// [List] and [Viewer] embed it. Anything else can too:
//
//	var sel comp.Range
//	sel.Start(cursor)          // shift pressed
//	lo, hi, ok := sel.Span(cursor)
//
// gowid ships copy mode twice — `widgets/copymodetable` and
// `widgets/copymodetree` — because its range lives on a widget and a table is
// not a list. Ours lived on [List] for the same reason and had the same
// consequence: a range over a table or a tree was the tool's to write (issue
// 72). It was already extracted and wired to two callers; exporting it is the
// rest of that work.
//
// # Why an anchor and not a pair of bounds
//
// A range grows from one end and a pair does not remember which end that was.
// Extend, turn round, extend past the start is the case a pair gets wrong.
//
// The invariant is small and easy to get wrong in the same way every time: the
// ZERO VALUE means no selection, so a fresh component cannot report a range
// from row 0. Two fields rather than an anchor of -1, because -1 is not a zero
// value and a struct that has to be constructed to be correct is one a caller
// will forget to construct.
//
// # What it does not do
//
// It does not hold the cursor, move anything, or know what a row is. That is
// what lets a tree pass a node index and a table pass a row number without
// this type learning about either.
//
// For a selection SET rather than a span — marks that survive the cursor moving
// away — see [Marks]. The two compose: [Marks.Span] takes the keys a Range
// covers.
type Range struct {
	anchor   int
	anchored bool
}

// Start anchors at cursor if nothing is anchored yet.
//
// The first extend from no selection anchors where the cursor IS, so a reader
// who presses shift-down once has selected two rows rather than none — which is
// what every other list in the world does.
//
// Idempotent on purpose: every key of a drag calls it, and only the first does
// anything.
func (s *Range) Start(cursor int) {
	if !s.anchored {
		s.anchor, s.anchored = cursor, true
	}
}

// Span is the selection low to high, and whether there is one.
//
// Ordered, so a caller acting on it never has to ask which way the reader
// dragged.
//
// DERIVED from the cursor passed in rather than stored, which is the property
// that makes it correct under a deferred move: a far end accumulated in a key
// handler is always one frame behind what the reader sees, and this cannot be
// (issues 45 and 77).
func (s *Range) Span(cursor int) (lo, hi int, ok bool) {
	if !s.anchored {
		return 0, 0, false
	}
	if s.anchor <= cursor {
		return s.anchor, cursor, true
	}
	return cursor, s.anchor, true
}

// Anchored reports whether a selection is in progress, for a caller drawing a
// hint or deciding what a key means.
func (s *Range) Anchored() bool { return s.anchored }

// Clear drops the selection, leaving the cursor.
func (s *Range) Clear() { s.anchor, s.anchored = 0, false }
