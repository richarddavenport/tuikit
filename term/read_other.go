//go:build !unix

package term

import (
	"os"
	"time"
)

// readReply has no implementation off unix, because neither does /dev/tty.
// Every caller treats the empty string as "this terminal draws characters",
// which is the correct answer for a platform this package cannot ask.
func readReply(*os.File, time.Duration) string { return "" }
