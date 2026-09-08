package comp

// Sort is which column a table is ordered by, and which way.
//
// # Where this came from
//
// Two tools. bottom's `canvas/components/data_table/sortable.rs` is 15 kB with
// a `SortOrder` enum, a sort column and header rendering that shows which.
// k9s's `ui/table.go` does the same.
//
// # What is here and what is not
//
// The COMPARISON is not here and never will be. A framework guessing how to
// order a Kubernetes age against a byte count would be wrong, and the tool owns
// its data and its types.
//
// What is here is the state, which is the same three fields and the same header
// arrow in every tool that has a sortable table: which column, which way, and
// what a click or a key does to that.
//
// # Why it is not a field on Table
//
// [Table] lays out and does not scroll; a scrolling aligned table is a [List]
// whose rows a Table laid out. Sort state belongs with the cursor rather than
// with the alignment, and putting it on Table would tie it to a type that is
// re-created every frame.
type Sort struct {
	// Col is the column being sorted by. Negative means unsorted, which is the
	// zero value — a table shows its rows in the order the tool built them
	// until somebody asks otherwise.
	Col int
	// Desc reverses it.
	Desc bool

	// unsorted distinguishes "no sort" from "sorted by column 0", because the
	// zero value has to mean the first of those.
	sorted bool
}

// On reports whether the table is sorted by a column at all.
func (s *Sort) On() bool { return s.sorted }

// By is the column and direction, and whether there is one.
func (s *Sort) By() (col int, desc, ok bool) { return s.Col, s.Desc, s.sorted }

// Toggle sorts by a column, reversing it if it was already the one.
//
// The behavior every spreadsheet has and every tool re-implements: click once
// for ascending, again for descending. A third click does NOT clear it, because
// a table that intermittently forgets its order is worse than one that keeps a
// direction you did not want.
func (s *Sort) Toggle(col int) {
	if col < 0 {
		return
	}
	if s.sorted && s.Col == col {
		s.Desc = !s.Desc
		return
	}
	s.Col, s.Desc, s.sorted = col, false, true
}

// Clear returns to the order the tool built the rows in.
func (s *Sort) Clear() { s.Col, s.Desc, s.sorted = 0, false, false }

// Marker is the arrow for a header cell, or "" for a column that is not the
// sorted one.
//
// From the chrome, so a tool that changed its arrows changed them everywhere.
// The blank for an unsorted column is the CALLER's problem rather than this
// one's: a header row is laid out by [Table], which pads columns to a width, so
// an empty string here is padded like anything else.
func (s *Sort) Marker(c *Canvas, col int) string {
	if !s.sorted || s.Col != col {
		return ""
	}
	ch := c.Chrome()
	if s.Desc {
		return ch.SortDesc
	}
	return ch.SortAsc
}

// Header decorates a header row with the arrow.
//
// Takes the labels and gives them back with the sorted one marked, so a caller
// hands the result straight to [Table.Rows]. The arrow goes on the END of the
// label rather than in a column of its own, because a column of arrows is a
// column of blanks for every table that is not sorted by it.
func (s *Sort) Header(c *Canvas, labels []string) []string {
	out := make([]string, len(labels))
	copy(out, labels)
	if !s.sorted || s.Col < 0 || s.Col >= len(out) {
		return out
	}
	out[s.Col] += s.Marker(c, s.Col)
	return out
}

// Apply orders indices by a comparison the caller supplies.
//
// less answers "does row a come before row b, ascending, by column col". The
// caller keeps its own data and its own types; this reverses for Desc, leaves
// the order alone when nothing is sorted, and is stable so rows that compare
// equal stay in the order the tool built them.
//
// Indices rather than the rows themselves, for the same reason [Tree.Visible]
// returns indices: a table of two hundred thousand should not be copied to be
// sorted.
func (s *Sort) Apply(n int, less func(a, b, col int) bool) []int {
	out := make([]int, n)
	for i := range out {
		out[i] = i
	}
	if !s.sorted || less == nil {
		return out
	}
	stableSort(out, func(a, b int) bool {
		if s.Desc {
			return less(b, a, s.Col)
		}
		return less(a, b, s.Col)
	})
	return out
}

// stableSort is an insertion sort, which is stable and is the right shape here.
//
// Not sort.SliceStable, because this runs on the caller's comparison over the
// caller's data and the allocation-free version is short. For the row counts a
// terminal shows it is also faster.
func stableSort(idx []int, less func(a, b int) bool) {
	for i := 1; i < len(idx); i++ {
		for j := i; j > 0 && less(idx[j], idx[j-1]); j-- {
			idx[j], idx[j-1] = idx[j-1], idx[j]
		}
	}
}
