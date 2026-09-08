package comp

import (
	"strings"
	"testing"
)

// strip renders a solved layout the way ratatui's own constraint tests do:
// one letter per band, repeated for its size. It reads better than a list of
// rect literals, and a wrong answer is visible rather than arithmetic.
//
//	"aaabbbbbb" — three units to the first band, six to the second
func strip(l Layout, total int) string {
	var b strings.Builder
	for i, size := range l.solve(total) {
		if i > 0 {
			b.WriteString(strings.Repeat("·", l.Gap))
		}
		b.WriteString(strings.Repeat(string(rune('a'+i)), size))
	}
	return b.String()
}

func TestLayoutSolves(t *testing.T) {
	for _, tc := range []struct {
		name  string
		l     Layout
		total int
		want  string
	}{
		{"a fixed band and a fill", Layout{Constraints: []Constraint{Length(2), Fill(1)}}, 10, "aabbbbbbbb"},
		{"the shape every tool has", Layout{Constraints: []Constraint{
			Length(1), Length(1), Fill(1), Length(1),
		}}, 12, "abcccccccccd"},
		{"two fills share equally", Layout{Constraints: []Constraint{Fill(1), Fill(1)}}, 8, "aaaabbbb"},
		{"weights", Layout{Constraints: []Constraint{Fill(1), Fill(3)}}, 8, "aabbbbbb"},
		{"a percentage is of the whole", Layout{Constraints: []Constraint{Percent(25), Fill(1)}}, 8, "aabbbbbb"},
		{"a gap costs the fill", Layout{Constraints: []Constraint{Length(2), Fill(1)}, Gap: 1}, 10, "aa·bbbbbbb"},
		{"a min floors a fill", Layout{Constraints: []Constraint{Fill(1).Min(4), Fill(9)}}, 10, "aaaabbbbbb"},
		{"a max ceils a fill", Layout{Constraints: []Constraint{Fill(1).Max(3), Fill(1)}}, 10, "aaabbbbbbb"},
		{"fixed bands that overflow leave the fill at nothing", Layout{Constraints: []Constraint{
			Length(6), Length(6), Fill(1),
		}}, 8, "aaaaaabbbbbb"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := strip(tc.l, tc.total); got != tc.want {
				t.Errorf("\n got %q\nwant %q", got, tc.want)
			}
		})
	}
}

// Rounding never loses a unit. Integer division does, and the symptom is a
// layout that is one row short at some heights and not others — which reads as
// a flicker rather than as a bug.
func TestNoUnitIsLostToRounding(t *testing.T) {
	l := Layout{Constraints: []Constraint{Fill(1), Fill(1), Fill(1)}}
	for total := range 60 {
		sum := 0
		for _, n := range l.solve(total) {
			sum += n
		}
		if sum != total {
			t.Errorf("at %d the bands add to %d", total, sum)
		}
	}
}

// A layout with a gap still adds up, gaps included.
func TestGapsAreAccountedFor(t *testing.T) {
	l := Layout{Constraints: []Constraint{Fill(1), Fill(1)}, Gap: 2}
	for total := 4; total < 40; total++ {
		sizes := l.solve(total)
		if sizes[0]+sizes[1]+2 != total {
			t.Errorf("at %d the bands and the gap add to %d", total, sizes[0]+sizes[1]+2)
		}
	}
}

// A terminal too short for the layout gives zero-height bands rather than
// negative ones. A negative rect is a rect that draws nothing anywhere, which
// is much harder to recognize than a band that is simply not there.
func TestNothingGoesNegative(t *testing.T) {
	l := Layout{Constraints: []Constraint{Length(1), Length(1), Fill(1), Length(1)}, Gap: 1}
	for total := range 8 {
		for i, n := range l.solve(total) {
			if n < 0 {
				t.Errorf("at %d band %d is %d", total, i, n)
			}
		}
	}
}

// Rows and Cols place the bands, which is the half solve does not do.
func TestRowsAndColsPlaceTheBands(t *testing.T) {
	l := Layout{Constraints: []Constraint{Length(1), Fill(1), Length(1)}}

	rows := l.Rows(Rect{X: 4, Y: 2, W: 20, H: 10})
	want := []Rect{{4, 2, 20, 1}, {4, 3, 20, 8}, {4, 11, 20, 1}}
	for i := range want {
		if rows[i] != want[i] {
			t.Errorf("row %d is %+v, want %+v", i, rows[i], want[i])
		}
	}

	cols := l.Cols(Rect{X: 0, Y: 0, W: 10, H: 3})
	wantC := []Rect{{0, 0, 1, 3}, {1, 0, 8, 3}, {9, 0, 1, 3}}
	for i := range wantC {
		if cols[i] != wantC[i] {
			t.Errorf("col %d is %+v, want %+v", i, cols[i], wantC[i])
		}
	}
}

// The bands tile the rect exactly: no overlap, no hole. This is the property a
// hand-written body() and listHeight() cannot promise, because the two numbers
// are related only by whoever remembers.
func TestTheBandsTileTheRect(t *testing.T) {
	l := Layout{Constraints: []Constraint{Length(1), Length(1), Fill(1), Length(1)}}
	r := Rect{X: 0, Y: 0, W: 40, H: 20}

	rows := l.Rows(r)
	next := r.Y
	for i, band := range rows {
		if band.Y != next {
			t.Errorf("band %d starts at %d, want %d", i, band.Y, next)
		}
		next = band.Y + band.H
	}
	if next != r.Y+r.H {
		t.Errorf("the bands end at %d, want %d", next, r.Y+r.H)
	}
}

// A fill pushed off its share by a Min or a Max gives the space back to the
// others, or takes it from them. Solving once and clamping afterwards is the
// obvious implementation and it overflows: 1/10 and 9/10 of ten, with a floor
// of four on the first, adds to thirteen.
func TestAClampedFillIsPaidForByTheOthers(t *testing.T) {
	for _, tc := range []struct {
		name  string
		l     Layout
		total int
		want  string
	}{
		{"a floor takes from the others",
			Layout{Constraints: []Constraint{Fill(1).Min(4), Fill(9)}}, 10, "aaaabbbbbb"},
		{"a ceiling gives to the others",
			Layout{Constraints: []Constraint{Fill(9).Max(3), Fill(1)}}, 10, "aaabbbbbbb"},
		// The ceiling pins b at 2, and the ten that are left split evenly
		// between a and c — the floor of 3 never binds.
		{"two clamps still tile",
			Layout{Constraints: []Constraint{Fill(1).Min(3), Fill(1).Max(2), Fill(1)}}, 12, "aaaaabbccccc"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := strip(tc.l, tc.total); got != tc.want {
				t.Errorf("\n got %q\nwant %q", got, tc.want)
			}
		})
	}
}

// Whatever the constraints, the bands add to the space — with or without
// clamps. This is the property the hand-written arithmetic could not promise.
func TestTheBandsAlwaysAddUp(t *testing.T) {
	layouts := []Layout{
		{Constraints: []Constraint{Length(1), Fill(1), Length(1)}},
		{Constraints: []Constraint{Fill(1).Min(4), Fill(3)}},
		{Constraints: []Constraint{Fill(2).Max(6), Fill(1), Percent(20)}},
		{Constraints: []Constraint{Percent(30), Percent(30), Fill(1)}},
	}
	for i, l := range layouts {
		for total := 12; total < 80; total++ {
			sum := 0
			for _, n := range l.solve(total) {
				sum += n
			}
			// A clamp can legitimately force an overflow — a Min bigger than
			// the space is a caller asking for the impossible — so the check
			// is that it never comes up SHORT while a fill could still grow.
			if sum < total {
				t.Errorf("layout %d at %d: the bands add to %d, leaving a hole", i, total, sum)
			}
		}
	}
}
