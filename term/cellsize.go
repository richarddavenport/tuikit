package term

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// Default cell size, for a terminal that will not answer. A common 14pt
// monospace face, so the error is a few pixels rather than a factor.
const (
	DefaultCellW = 10
	DefaultCellH = 20
)

// CellSize is the pixel size of one character cell, and whether the terminal
// actually said so.
//
// A picture is requested in cells and drawn in pixels, and nothing but the
// terminal knows the ratio. Guessing it wrong does not fail, it SCALES — the
// picture comes out a fraction of the cells it was meant to cover, which looks
// like a bug in the drawing rather than a missing measurement. That is why the
// bool is here: an assumed size is worth saying out loud.
func CellSize() (w, h int, measured bool) {
	if w, h, ok := overrideCellSize(); ok {
		return w, h, true
	}
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return DefaultCellW, DefaultCellH, false
	}
	defer f.Close() //nolint:errcheck // a query we are done with

	// The kernel first, because it turned out to be the honest channel.
	// Measured on one HiDPI Mac, both terminals on the same display:
	//
	//	              CSI 14t says   TIOCGWINSZ says   draws at
	//	  Ghostty        16x34            16x34         16x34
	//	  iTerm2          8x17            16x34         16x34
	//
	// The escape is answered by the terminal application, which picks points or
	// device pixels and says which nowhere. ws_xpixel is what that same
	// terminal wrote into the tty, and both wrote device pixels. So the ratio
	// issue 31 refused to guess is not a guess: it is this division.
	if w, h, ok := kernelCellSize(f); ok {
		return w, h, true
	}
	return QueryCellSize(f, DefaultTimeout)
}

// kernelCellSize divides the window's pixels by its cells.
//
// Refused rather than trusted when the answer is not a plausible cell, because
// a terminal that fills the pixel fields with something else would otherwise
// take the picture with it — and the escape, which is right on most terminals,
// would never be asked.
func kernelCellSize(f *os.File) (w, h int, ok bool) {
	px, py, ok := WindowPixels(f)
	if !ok {
		return 0, 0, false
	}
	cols, rows, ok := WindowCells(f)
	if !ok || cols <= 0 || rows <= 0 {
		return 0, 0, false
	}
	return cellFrom(px, py, cols, rows)
}

// cellFrom is the arithmetic, split out so it can be tested without a tty.
func cellFrom(px, py, cols, rows int) (w, h int, ok bool) {
	if px <= 0 || py <= 0 || cols <= 0 || rows <= 0 {
		return 0, 0, false
	}
	// Rounded, not truncated: 3832 pixels over 239 columns is 16.03, and a cell
	// one pixel narrow tiles a whole row of pictures short.
	w, h = (px+cols/2)/cols, (py+rows/2)/rows
	if w < 2 || h < 2 || w > 200 || h > 200 {
		return 0, 0, false
	}
	return w, h, true
}

// EnvCellSize forces the cell size, as WxH in pixels.
//
// It exists because terminals disagree about what a "pixel" is: on a HiDPI Mac
// iTerm2 answers CSI 14t in POINTS and draws at device resolution, so every
// Sixel picture came out at exactly half size. That is now DETECTED rather than
// declared — [CellSize] asks the kernel first, and TIOCGWINSZ reported true
// device pixels on both terminals measured (issue 31).
//
// So this is no longer the answer to HiDPI. It stays as the override of last
// resort, for a terminal that fills neither channel honestly, and because
// something that draws pictures should let a reader state a measurement it got
// wrong.
const EnvCellSize = "TUIKIT_CELL_SIZE"

func overrideCellSize() (w, h int, ok bool) {
	v := strings.TrimSpace(os.Getenv(EnvCellSize))
	if v == "" {
		return 0, 0, false
	}
	parts := strings.SplitN(strings.ToLower(v), "x", 2)
	if len(parts) != 2 {
		return 0, 0, false
	}
	w, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	h, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	if err1 != nil || err2 != nil || w <= 0 || h <= 0 {
		return 0, 0, false
	}
	return w, h, true
}

// QueryCellSize asks three questions at once, because no single one is
// answered everywhere.
//
//   - CSI 16 t is the cell size directly, replied as CSI 6 ; HEIGHT ; WIDTH t.
//     Note the order: height first, which is the obvious thing to get backwards.
//   - CSI 14 t is the text area in PIXELS, replied as CSI 4 ; height ; width t.
//   - CSI 18 t is the text area in CHARACTERS, replied as CSI 8 ; rows ; cols t.
//
// Dividing the second by the third gives the cell size on terminals that
// decline the first — iTerm2 among them, which answered nothing to CSI 16 t and
// so drew every picture at the default 10×20 while its cells were larger. The
// bars came out about two thirds the width they should have been.
func QueryCellSize(f *os.File, timeout time.Duration) (w, h int, measured bool) {
	if !raw(f) {
		return DefaultCellW, DefaultCellH, false
	}
	restore := makeRaw(f)
	if restore == nil {
		return DefaultCellW, DefaultCellH, false
	}
	defer restore()

	if _, err := f.WriteString("\x1b[16t\x1b[14t\x1b[18t\x1b[c"); err != nil {
		return DefaultCellW, DefaultCellH, false
	}
	return ParseCellSize(readReply(f, timeout))
}

// ParseCellSize reads whichever of the three replies arrived.
func ParseCellSize(reply string) (w, h int, measured bool) {
	// The direct answer, if there is one.
	if height, width, ok := twoNumbers(reply, "\x1b[6;"); ok {
		return width, height, true
	}
	// Otherwise divide the pixel size of the text area by its size in cells.
	pxH, pxW, okPx := twoNumbers(reply, "\x1b[4;")
	rows, cols, okCells := twoNumbers(reply, "\x1b[8;")
	if okPx && okCells && rows > 0 && cols > 0 {
		if cw, ch := pxW/cols, pxH/rows; cw > 0 && ch > 0 {
			return cw, ch, true
		}
	}
	return DefaultCellW, DefaultCellH, false
}

// twoNumbers reads `prefix A ; B t` and returns A and B.
func twoNumbers(reply, prefix string) (a, b int, ok bool) {
	body, found := cut(reply, prefix, "t")
	if !found {
		return 0, 0, false
	}
	parts := strings.Split(body, ";")
	if len(parts) != 2 {
		return 0, 0, false
	}
	first, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	second, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	// A zero would divide or stretch to nothing later; an implausible answer is
	// treated as no answer.
	if err1 != nil || err2 != nil || first <= 0 || second <= 0 {
		return 0, 0, false
	}
	return first, second, true
}
