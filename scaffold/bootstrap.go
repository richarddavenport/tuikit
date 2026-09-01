package scaffold

import (
	"fmt"
	"os/exec"
	"strings"
)

// Bootstrap finishes a generated tool: resolves its modules, then writes the
// goldens for the screens it was born with.
//
// Both steps are part of generating the tool rather than something to do
// afterwards, because a tool that does not pass `go test ./...` on its first
// run has taught its owner, in the first thirty seconds, that the tests are
// noise. A golden cannot ship in a template — it is a render of the code, and
// the code has a name substituted into it — so it is made here instead.
//
// Separate from Write so that writing files and running a toolchain stay
// separate. A caller with no Go available can still generate.
func Bootstrap(root string) error {
	for _, step := range [][]string{
		{"go", "mod", "tidy"},
		// Only the package that has goldens: -update-goldens is registered by
		// the harness, so a package that does not import it rejects the flag.
		{"go", "test", "./internal/tui", "-update-goldens"},
	} {
		cmd := exec.Command(step[0], step[1:]...)
		cmd.Dir = root
		if out, err := cmd.CombinedOutput(); err != nil {
			return fmt.Errorf("%s: %w\n%s", strings.Join(step, " "), err, out)
		}
	}
	return nil
}
