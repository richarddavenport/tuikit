//go:build unix

package term

import (
	"os"

	"golang.org/x/sys/unix"
)

// WindowPixels asks the KERNEL for the window's size in pixels, via TIOCGWINSZ.
//
// A second opinion, and the point is that it is a different channel. `CSI 14 t`
// is answered by the terminal application, which chooses whether to speak in
// points or in device pixels and says which nowhere. ws_xpixel/ws_ypixel are
// whatever the terminal wrote into the tty when it last set the window size.
//
// If a terminal reports points in-band and device pixels here, the ratio
// between them IS the scale factor, and the guess issue 31 refuses to make
// becomes a measurement. If both agree, this channel is no help and the
// environment variable stays the answer.
//
// Zero, zero, false when the kernel has no answer — which is common: plenty of
// terminals never set the pixel fields at all.
func WindowPixels(f *os.File) (w, h int, ok bool) {
	sz, err := unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ)
	if err != nil || sz.Xpixel == 0 || sz.Ypixel == 0 {
		return 0, 0, false
	}
	return int(sz.Xpixel), int(sz.Ypixel), true
}

// WindowCells is the window in rows and columns, from the same call, so a
// caller can divide one by the other and get a cell size the terminal never
// stated.
func WindowCells(f *os.File) (cols, rows int, ok bool) {
	sz, err := unix.IoctlGetWinsize(int(f.Fd()), unix.TIOCGWINSZ)
	if err != nil || sz.Col == 0 || sz.Row == 0 {
		return 0, 0, false
	}
	return int(sz.Col), int(sz.Row), true
}
