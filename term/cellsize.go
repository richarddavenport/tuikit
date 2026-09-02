package term

import (
	"os"
	"strconv"
	"strings"
	"time"
)

// CellSize is the pixel size of one character cell.
//
// A picture is requested in cells and drawn in pixels, and nothing but the
// terminal knows the ratio between them: it depends on the font, its size, and
// the line spacing the reader chose. Guessing it wrong does not fail, it
// stretches — which is worse, because a stretched picture looks like a bad
// design rather than a missing measurement.
func CellSize() (w, h int) {
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return DefaultCellW, DefaultCellH
	}
	defer f.Close() //nolint:errcheck // a query we are done with
	return QueryCellSize(f, DefaultTimeout)
}

// Default cell size, for a terminal that will not answer. A common 14pt
// monospace face, so the error is a few pixels rather than a factor.
const (
	DefaultCellW = 10
	DefaultCellH = 20
)

// QueryCellSize asks with CSI 16 t, which is answered CSI 6 ; height ; width t.
//
// Note the reply order: HEIGHT first. The obvious reading is width-then-height
// and it is wrong, which is the kind of mistake that produces a picture that is
// merely a bit off rather than obviously broken.
func QueryCellSize(f *os.File, timeout time.Duration) (w, h int) {
	if !raw(f) {
		return DefaultCellW, DefaultCellH
	}
	restore := makeRaw(f)
	if restore == nil {
		return DefaultCellW, DefaultCellH
	}
	defer restore()

	if _, err := f.WriteString("\x1b[16t\x1b[c"); err != nil {
		return DefaultCellW, DefaultCellH
	}
	return ParseCellSize(readUntilDA1(f, timeout))
}

// ParseCellSize reads the reply. Exported for the same reason Parse is: the
// parsing is where this goes wrong and it should be checkable without a tty.
func ParseCellSize(reply string) (w, h int) {
	body, ok := cut(reply, "\x1b[6;", "t")
	if !ok {
		return DefaultCellW, DefaultCellH
	}
	parts := strings.Split(body, ";")
	if len(parts) != 2 {
		return DefaultCellW, DefaultCellH
	}
	height, err1 := strconv.Atoi(strings.TrimSpace(parts[0]))
	width, err2 := strconv.Atoi(strings.TrimSpace(parts[1]))
	// A zero or a negative would divide or stretch to nothing later; treat an
	// implausible answer as no answer.
	if err1 != nil || err2 != nil || width <= 0 || height <= 0 {
		return DefaultCellW, DefaultCellH
	}
	return width, height
}
