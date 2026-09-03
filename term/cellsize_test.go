package term

import "testing"

// The kernel's answer is refused when it is not a plausible cell, so a terminal
// that fills the pixel fields with something else does not take the picture
// with it — and the escape, which is right on most terminals, still gets asked.
func TestAnImplausibleKernelCellIsRefused(t *testing.T) {
	for _, tc := range []struct {
		name               string
		px, py, cols, rows int
		want               bool
	}{
		{"a real HiDPI window", 3776, 2006, 236, 59, true},
		{"a real non-HiDPI window", 1888, 1003, 236, 59, true},
		{"no pixels reported", 0, 0, 236, 59, false},
		{"no cells reported", 3776, 2006, 0, 0, false},
		{"a cell one pixel wide", 236, 2006, 236, 59, false},
		{"a cell the size of the screen", 3776, 2006, 2, 2, false},
	} {
		t.Run(tc.name, func(t *testing.T) {
			w, h, ok := cellFrom(tc.px, tc.py, tc.cols, tc.rows)
			if ok != tc.want {
				t.Fatalf("accepted=%v, want %v (got %dx%d)", ok, tc.want, w, h)
			}
		})
	}
}

// Rounded rather than truncated. 3832 pixels over 239 columns is 16.03, and a
// cell one pixel narrow tiles a whole row of pictures short.
func TestTheKernelCellIsRoundedNotTruncated(t *testing.T) {
	if w, _, _ := cellFrom(3832, 2028, 239, 59); w != 16 {
		t.Errorf("3832/239 came out %d, want 16", w)
	}
	// 17.6 rounds up rather than down to 17.
	if w, _, _ := cellFrom(3520, 2028, 200, 59); w != 18 {
		t.Errorf("3520/200 came out %d, want 18", w)
	}
}
