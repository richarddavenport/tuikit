package term

import (
	"fmt"
	"image/color"
	"os"
	"strings"
	"time"
)

// Colors asks the terminal what its palette actually is, by index.
//
// This is what keeps the pixel layer honest. Decision 28 put the palette on
// ANSI 0–15 so that a tuikit tool follows the reader's terminal theme rather
// than overriding it — and a picture with a hardcoded `#5f00ff → #ff5faf`
// gradient would walk straight back out of that, looking identical under all 22
// Omarchy themes while the characters beside it changed.
//
// So the gradient is not chosen here. It is read from the same sixteen the
// characters use, with OSC 4, and a reader who changes their theme changes the
// picture too.
//
// Indices with no answer are absent from the map; the caller decides what to do
// about that rather than being handed a plausible lie.
func Colors(indices ...int) map[int]color.RGBA {
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return nil
	}
	defer f.Close()
	return QueryColors(f, DefaultTimeout, indices...)
}

// QueryColors sends one OSC 4 per index, then DA1 to end the read.
func QueryColors(f *os.File, timeout time.Duration, indices ...int) map[int]color.RGBA {
	if len(indices) == 0 || !raw(f) {
		return nil
	}
	restore := makeRaw(f)
	if restore == nil {
		return nil
	}
	defer restore()

	var q strings.Builder
	for _, i := range indices {
		fmt.Fprintf(&q, "\x1b]4;%d;?\x1b\\", i)
	}
	q.WriteString("\x1b[c") // the terminator, as everywhere else here
	if _, err := f.WriteString(q.String()); err != nil {
		return nil
	}
	return ParseColors(readUntilDA1(f, timeout))
}

// ParseColors reads a run of OSC 4 replies: ESC ] 4 ; N ; rgb:R/G/B ST.
//
// Replies may arrive in any order and a terminal may answer only some of them,
// so each is matched by the index it carries rather than by position.
func ParseColors(reply string) map[int]color.RGBA {
	out := map[int]color.RGBA{}
	for _, part := range strings.Split(reply, "\x1b]4;")[1:] {
		var idx int
		if _, err := fmt.Sscanf(part, "%d;", &idx); err != nil {
			continue
		}
		if c, ok := ParseRGB(part); ok {
			out[idx] = c
		}
	}
	if len(out) == 0 {
		return nil
	}
	return out
}
