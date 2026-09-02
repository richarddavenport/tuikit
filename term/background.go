package term

import (
	"image/color"
	"os"
	"strconv"
	"strings"
	"time"
)

// DefaultBackground is assumed when the terminal will not say: a dark grey
// rather than pure black, because most dark themes are, and a shadow baked
// against the wrong black shows up as a lighter rectangle.
var DefaultBackground = color.RGBA{R: 0x1a, G: 0x1a, B: 0x1a, A: 0xff}

// Background is the terminal's background colour, as RGB.
//
// Only Sixel needs it, and it needs it badly. Sixel has no alpha channel, so a
// soft edge or a rounded corner has to be composited against the background
// BEFORE it is sent — and getting it wrong does not soften anything, it draws a
// visible rectangle of not-quite-background around the picture.
//
// This is also the one piece of theme information tuikit reads back. Decision
// 28 put the palette on ANSI 0–15 so the reader's terminal theme wins; a Sixel
// picture cannot honour that theme without knowing what it actually is.
func Background() color.RGBA {
	f, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		return DefaultBackground
	}
	defer f.Close() //nolint:errcheck // a query we are done with
	return QueryBackground(f, DefaultTimeout)
}

// QueryBackground asks with OSC 11.
func QueryBackground(f *os.File, timeout time.Duration) color.RGBA {
	if !raw(f) {
		return DefaultBackground
	}
	restore := makeRaw(f)
	if restore == nil {
		return DefaultBackground
	}
	defer restore()

	// OSC 11 ; ? then DA1, so the read ends even on a terminal that ignores it.
	if _, err := f.WriteString("\x1b]11;?\x1b\\\x1b[c"); err != nil {
		return DefaultBackground
	}
	return ParseBackground(readUntilDA1(f, timeout))
}

// ParseBackground reads an OSC 11 reply: rgb:RRRR/GGGG/BBBB.
//
// The components are variable width — 1 to 4 hex digits each — and are a
// FRACTION of their own maximum rather than a byte. rgb:ffff/0000/0000 and
// rgb:ff/00/00 are the same red, which is why each component is scaled by its
// own length instead of truncated to two digits.
func ParseBackground(reply string) color.RGBA {
	if c, ok := ParseRGB(reply); ok {
		return c
	}
	return DefaultBackground
}

// ParseRGB reads an X11 rgb: colour out of a reply, and says whether it found
// one.
//
// The bool is not decoration. ParseBackground has to answer with SOME colour,
// so it substitutes a default on failure — and a caller that needs to tell
// "the terminal said #1a1a1a" from "the terminal said nothing" cannot, because
// those are the same value. That is precisely the collision decision 28 found
// in Hex, where #000000 meant both index 0 and failure, and it is worth not
// building twice.
func ParseRGB(reply string) (color.RGBA, bool) {
	i := strings.Index(reply, "rgb:")
	if i < 0 {
		return color.RGBA{}, false
	}
	body := reply[i+len("rgb:"):]
	// The reply is terminated by BEL or ST; either ends the colour.
	for _, end := range []string{"\x07", "\x1b\\", "\x1b"} {
		if j := strings.Index(body, end); j >= 0 {
			body = body[:j]
		}
	}
	parts := strings.Split(strings.TrimSpace(body), "/")
	if len(parts) != 3 {
		return color.RGBA{}, false
	}
	var out [3]uint8
	for i, p := range parts {
		v, err := strconv.ParseUint(p, 16, 32)
		if err != nil || len(p) == 0 || len(p) > 4 {
			return color.RGBA{}, false
		}
		// Rounded, not truncated. A terminal may answer 5f00 or 5f5f for the
		// same colour — both are 0x5f as a fraction of their own maximum — and
		// truncation turns the first into 94 while the second stays 95. One
		// off is invisible in a gradient and very visible in a test.
		full := uint64(1)<<(4*len(p)) - 1
		out[i] = uint8((v*255 + full/2) / full)
	}
	return color.RGBA{R: out[0], G: out[1], B: out[2], A: 0xff}, true
}
