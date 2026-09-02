package term

import (
	"os"
	"strings"
)

// Multiplexer names the terminal multiplexer between this program and the real
// terminal, or "" when there is none.
//
// It matters because a multiplexer is a terminal emulator too. It parses what
// programs write, keeps its own screen, and repaints panes — so a graphics
// escape sequence has to be understood and re-emitted by it, not merely passed
// along. What actually happens depends on the multiplexer and the protocol:
//
//   - tmux 3.4 and later parses Sixel natively and renders it into its own
//     buffer, so Sixel usually survives.
//   - The kitty protocol is carried in APC sequences, which tmux does not
//     forward unless `allow-passthrough` is on AND the program wraps them.
//     tuikit does not wrap them, so kitty images are dropped inside tmux even
//     though the QUERY — being small and answered directly — gets through and
//     reports support.
//
// That last combination is the confusing one: detection says kitty, and no
// image ever appears. Naming the multiplexer is what turns that into an
// explanation instead of a bug hunt.
func Multiplexer() string {
	switch {
	case os.Getenv("ZELLIJ") != "":
		return "zellij"
	case os.Getenv("TMUX") != "":
		return "tmux"
	case os.Getenv("STY") != "":
		return "screen"
	case os.Getenv("HERD_SESSION") != "", os.Getenv("HERDR_SESSION") != "":
		return "herdr"
	}
	// TERM is the last resort: a multiplexer usually sets its own, and a
	// program attached from elsewhere may not have inherited the variables.
	if t := os.Getenv("TERM"); strings.HasPrefix(t, "screen") || strings.HasPrefix(t, "tmux") {
		return "tmux or screen"
	}
	return ""
}
