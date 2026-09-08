package comp

// Layout divides a rect along one axis.
//
// # Where this came from
//
// The last arithmetic every tool writes by hand. Each of them has some version
// of a header, an optional strip, a body, and a footer, computed twice — once
// as a rect and once as a height — with the two kept in step by whoever
// remembers:
//
//	func (m *Model) body() Rect     { return Rect{X: 0, Y: 2, W: m.width, H: m.listHeight() + 2} }
//	func (m *Model) listHeight() int { return max(3, m.height-5) }
//
// the cloud tool got that wrong this week: adding one row for a pivot strip
// means changing a literal in one function and a different literal in another,
// and nothing relates them. There is no test that can catch it, because both
// numbers are equally plausible.
//
// # Why not ratatui's solver
//
// bento ports Cassowary — a real constraint solver, with ratatui's strength
// constants — and it is the right answer for a general-purpose library that
// must lay out anything. Measured against what four tools actually lay out,
// which is a column of bands and a pane split in two, a single linear pass is
// enough. Fill, Min and Max cover every case in the deploy tool, the database
// tool, the cloud tool and democtl.
//
// Generalise past that when two of the four need it, which is the rule every
// component here was extracted under. Flex alignment and negative spacing
// exist in ratatui because ratatui is a general library; we are not one.
type Layout struct {
	// Constraints are the bands, in order.
	Constraints []Constraint
	// Gap is the space left between one band and the next. Zero for none.
	Gap int
}

// Constraint is how much of the axis one band takes.
//
// A closed set of four, built by the functions below rather than by a struct
// literal, so a band cannot be two kinds of thing at once — which is what a
// struct with a Length and a Weight and a Min on it would allow.
type Constraint struct {
	kind   constraintKind
	value  int
	lo, hi int
}

type constraintKind int

const (
	kindLength constraintKind = iota
	kindPercent
	kindFill
)

// Length is a fixed size: a header of one row, a footer of one.
func Length(n int) Constraint { return Constraint{kind: kindLength, value: n} }

// Percent is a share of the whole, before any fill is distributed.
func Percent(n int) Constraint { return Constraint{kind: kindPercent, value: n} }

// Fill takes what is left, split between the fills by weight. Weight zero is
// treated as one, so Fill(0) is an ordinary equal share rather than a band
// that silently disappears.
func Fill(weight int) Constraint { return Constraint{kind: kindFill, value: max(1, weight)} }

// Min floors a constraint. A body clamped to three rows still shows something
// on a terminal too short for the layout, which beats a negative height.
func (c Constraint) Min(n int) Constraint { c.lo = n; return c }

// Max ceils a constraint, for a pane that should not grow past what it has to
// show.
func (c Constraint) Max(n int) Constraint { c.hi = n; return c }

// Rows divides r top to bottom.
func (l Layout) Rows(r Rect) []Rect {
	sizes := l.solve(r.H)
	out := make([]Rect, len(sizes))
	y := r.Y
	for i, h := range sizes {
		out[i] = Rect{X: r.X, Y: y, W: r.W, H: h}
		y += h + l.Gap
	}
	return out
}

// Cols divides r left to right.
func (l Layout) Cols(r Rect) []Rect {
	sizes := l.solve(r.W)
	out := make([]Rect, len(sizes))
	x := r.X
	for i, w := range sizes {
		out[i] = Rect{X: x, Y: r.Y, W: w, H: r.H}
		x += w + l.Gap
	}
	return out
}

// solve resolves the bands against the space available.
//
// Fixed sizes first, then the remainder split between the fills by weight —
// and then again, because a fill whose Min or Max moved it is no longer taking
// its share, and the space it gave up or took has to come from somewhere. Each
// pass pins at least one band, so it terminates in at most one round per fill.
//
// That re-solve is the whole reason ratatui reaches for a constraint solver.
// At this scale — a handful of bands on one axis — pinning and repeating gets
// the same answer, and a reader can follow it.
func (l Layout) solve(total int) []int {
	n := len(l.Constraints)
	sizes := make([]int, n)
	if n == 0 {
		return sizes
	}
	space := max(0, total-l.Gap*(n-1))

	fixed := 0
	pinned := make([]bool, n)
	for i, c := range l.Constraints {
		switch c.kind {
		case kindLength:
			sizes[i] = c.clamp(c.value)
		case kindPercent:
			sizes[i] = c.clamp(space * c.value / 100)
		case kindFill:
			continue
		}
		pinned[i] = true
		fixed += sizes[i]
	}

	for {
		weight, taken := 0, fixed
		for i, c := range l.Constraints {
			switch {
			case c.kind != kindFill:
			case pinned[i]:
				taken += sizes[i]
			default:
				weight += c.value
			}
		}
		if weight == 0 {
			break
		}
		// Negative when the fixed bands already overflow, in which case every
		// fill is zero and the overflow is the canvas's to clip — a band that
		// is not there beats a rect with a negative height.
		rest := max(0, space-taken)

		clamped := false
		for i, c := range l.Constraints {
			if c.kind != kindFill || pinned[i] {
				continue
			}
			want := rest * c.value / weight
			if got := c.clamp(want); got != want {
				// A Min or a Max moved it, so it is no longer taking its
				// share: pin it and give the rest back to the others.
				sizes[i], pinned[i], clamped = got, true, true
				continue
			}
			sizes[i] = want
		}
		if clamped {
			continue
		}

		// The rounding remainder, one unit at a time to the heaviest fill that
		// can still take it — so a band never loses a unit while another gains
		// two, which is what plain integer division does and what a reader
		// notices as a layout that is one off at some sizes and not others.
		given := 0
		for i, c := range l.Constraints {
			if c.kind == kindFill && !pinned[i] {
				given += sizes[i]
			}
		}
		for spare := rest - given; spare > 0; spare-- {
			i := l.hungriest(sizes, pinned)
			if i < 0 {
				break
			}
			sizes[i]++
		}
		break
	}
	return sizes
}

// hungriest is the fill band furthest below what its weight asked for, and
// that its Max still allows to grow.
func (l Layout) hungriest(sizes []int, pinned []bool) int {
	best, bestWeight := -1, 0
	for i, c := range l.Constraints {
		if c.kind != kindFill || pinned[i] || (c.hi > 0 && sizes[i] >= c.hi) {
			continue
		}
		if c.value > bestWeight {
			best, bestWeight = i, c.value
		}
	}
	return best
}

// clamp applies Min and Max, in that order: a Min above a Max is a caller
// asking for something impossible, and the floor is the one that keeps a pane
// usable.
func (c Constraint) clamp(n int) int {
	if c.hi > 0 {
		n = min(n, c.hi)
	}
	return max(n, c.lo)
}
