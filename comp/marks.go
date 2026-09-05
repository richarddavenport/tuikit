package comp

import "sort"

// Marks is what a reader has picked out, kept by their own identity.
//
// # Where this came from
//
// Three tools, and they disagree about the shape in a way worth knowing.
//
// yazi splits it in two: `tab/visual.rs` is a live anchored range you drag, and
// `tab/selected.rs` is an IndexMap keyed by URL that the range commits into.
// The range is primary.
//
// k9s does not have an anchor at all. `SelectTable.SpanMark` scans backwards
// from the cursor for the nearest existing mark, forwards if it finds none, and
// fills between. The set is primary and the range is derived from it.
//
// pgctl asked for the same thing from a third direction (issue 49).
//
// # Why it is not a field on List
//
// Because the two gestures want the same set underneath, and because a
// selection outlives the component that drew it. A list is rebuilt every frame
// from rows the tool owns; what a reader ticked has to survive a filter, a
// refresh, and moving to another screen and back. termshark makes the point
// from another angle: it has a copy mode over a table AND over a tree, so the
// thing being selected from is not always a [List].
//
// # Why a key rather than an index
//
// Indices move. A filter that removes one row above a marked one would
// otherwise mark a different row, and it would not error while doing it — the
// same failure shape [Tree] keys around. yazi keys by URL and k9s by row ID for
// exactly this reason.
type Marks struct {
	keys map[string]bool
}

// Toggle marks an unmarked key and unmarks a marked one.
func (m *Marks) Toggle(key string) {
	if key == "" {
		return
	}
	if m.keys[key] {
		delete(m.keys, key)
		return
	}
	m.Add(key)
}

// Add marks a key.
func (m *Marks) Add(keys ...string) {
	if m.keys == nil {
		m.keys = map[string]bool{}
	}
	for _, k := range keys {
		if k != "" {
			m.keys[k] = true
		}
	}
}

// Remove unmarks a key.
func (m *Marks) Remove(keys ...string) {
	for _, k := range keys {
		delete(m.keys, k)
	}
}

// Has reports whether a key is marked, for a row drawing its own tick.
func (m *Marks) Has(key string) bool { return m.keys[key] }

// Len is how many are marked.
func (m *Marks) Len() int { return len(m.keys) }

// Clear unmarks everything.
func (m *Marks) Clear() { m.keys = nil }

// Keys is what is marked, sorted.
//
// Sorted rather than map order, because an action taken on several things
// should happen in the same order twice — and because an interface that
// reorders itself between runs is one nobody can screenshot.
func (m *Marks) Keys() []string {
	out := make([]string, 0, len(m.keys))
	for k := range m.keys {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

// Acting is what an action should act on: the marks, or the cursor's key when
// nothing is marked.
//
// # This is the method that makes marks worth having
//
// k9s's `GetSelectedItems` does exactly this, and it is why `d` can mean
// "delete the marked pods" and "delete this one" through a single code path:
//
//	func (s *SelectTable) GetSelectedItems() []string {
//		if s.marks.Len() == 0 {
//			if item := s.GetSelectedItem(); item != "" {
//				return []string{item}
//			}
//			return nil
//		}
//		return s.marks.UnsortedList()
//	}
//
// Neither yazi's version nor pgctl's request made this visible. A component
// that handed back the set and left every call site to write that `if` would
// have done the easy half.
//
// Returns nil when nothing is marked and there is no cursor, which is an
// ordinary state — an empty list — rather than an error.
func (m *Marks) Acting(cursorKey string) []string {
	if m.Len() > 0 {
		return m.Keys()
	}
	if cursorKey == "" {
		return nil
	}
	return []string{cursorKey}
}

// Span marks every key in a range, which is how a dragged selection commits.
//
// The [List.Range] half of yazi's arrangement: extend with shift-arrow, then
// hand the keys the range covers to this.
func (m *Marks) Span(keys []string) { m.Add(keys...) }

// Fill marks from the nearest already-marked key to the cursor.
//
// k9s's arrangement, where the set is primary and there is no live anchor:
// scan back from the cursor for a mark, then forward if there was none, and
// mark everything between. Both gestures are here because two tools each chose
// one, and a component that picked a side would be wrong for the other.
//
// keys is the visible rows in order, and at is the cursor's index into it.
// Reports whether it found anything to fill from.
func (m *Marks) Fill(keys []string, at int) bool {
	if at < 0 || at >= len(keys) {
		return false
	}
	from := -1
	for i := at - 1; i >= 0; i-- {
		if m.Has(keys[i]) {
			from = i
			break
		}
	}
	if from == -1 {
		for i := at + 1; i < len(keys); i++ {
			if m.Has(keys[i]) {
				from = i
				break
			}
		}
	}
	if from == -1 {
		return false
	}
	lo, hi := from, at
	if lo > hi {
		lo, hi = hi, lo
	}
	m.Add(keys[lo : hi+1]...)
	return true
}
