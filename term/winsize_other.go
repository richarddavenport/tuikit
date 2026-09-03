//go:build !unix

package term

import "os"

// WindowPixels has no answer off unix: TIOCGWINSZ is a tty ioctl.
func WindowPixels(*os.File) (w, h int, ok bool) { return 0, 0, false }

// WindowCells has no answer off unix, for the same reason.
func WindowCells(*os.File) (cols, rows int, ok bool) { return 0, 0, false }
