//go:build unix

package term

import (
	"os"
	"strings"
	"time"

	"golang.org/x/sys/unix"
)

// readReply reads until the DA1 terminator arrives or the deadline passes.
//
// # Why this is not os.File.SetReadDeadline
//
// It was, and that was the bug behind "graphics none" on a terminal that had
// just answered. On macOS a file opened on /dev/tty is not registered with Go's
// runtime poller, so SetReadDeadline returns ErrNoDeadline — and the code
// treated that as "no reply" and returned immediately, having already SENT the
// query. The terminal's answer then arrived with nobody reading it, sat in the
// input buffer, and was echoed by the shell after the program exited:
//
//	^[_Gi=1;invalid payload^[\
//	^[[?64;1;2;4;6;17;18;21;22c
//
// Both answers were there. Nothing read them. Worse, the leak is the loud
// symptom of a silent one — every capability came back false, so the whole
// pixel layer switched itself off on a terminal that supports two protocols.
//
// So the deadline is enforced here instead: the descriptor is put in
// non-blocking mode and read directly, with EAGAIN meaning "nothing yet" rather
// than "nothing coming". Opening /dev/tty creates its own open file
// description, so O_NONBLOCK applies to ours alone and the shell's terminal is
// untouched.
func readReply(f *os.File, timeout time.Duration) string {
	fd := int(f.Fd())
	if err := unix.SetNonblock(fd, true); err == nil {
		defer unix.SetNonblock(fd, false) //nolint:errcheck // ours alone
	}

	deadline := time.Now().Add(timeout)
	var b strings.Builder
	buf := make([]byte, 256)
	for b.Len() < maxReply && time.Now().Before(deadline) {
		n, err := unix.Read(fd, buf)
		if n > 0 {
			b.Write(buf[:n])
			if da1Complete(b.String()) {
				break // answered; do not spend the rest of the timeout
			}
			continue
		}
		if err == unix.EAGAIN || err == unix.EWOULDBLOCK || n == 0 {
			// Nothing yet. A short sleep rather than a spin, because this runs
			// once at start-up and the terminal is a person's, not a server's.
			time.Sleep(2 * time.Millisecond)
			continue
		}
		break
	}
	return b.String()
}
