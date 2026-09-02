package guard

import (
	"strings"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/spec"
)

// Keys fails when a command has a keyboard path the help screen does not name.
//
// A help screen written out by hand drifts the first time somebody adds a
// binding and forgets, and nothing catches it — because a help screen that is
// slightly wrong still renders perfectly, and the key it omits is exactly the
// one a reader came looking for.
//
// That is the same argument as [Reachable]: the fix for two lists that must
// agree is not diligence, it is one list. This is the cheaper version of one
// list — two lists and a test that says when they disagree.
//
// # Only one direction
//
// Every command with a Key must be listed. The reverse is NOT checked, and
// deliberately: a tool's screen-level keys — scroll, filter, close the pane —
// are not commands and never will be, so a help screen that names more than the
// spec does is a help screen doing its job. Failing on those would push a tool
// into declaring `j/k move` as a command, which would then appear in the CLI
// and in `describe --json`, and the guard would have made the surface worse to
// keep itself happy.
func Keys(t T, root spec.Command, sections []comp.KeySection) {
	t.Helper()

	listed := map[string]bool{}
	for _, section := range sections {
		for _, hint := range section.Keys {
			// A hint may name several keys for one action — "j/k", "↑↓". Each
			// counts, or a tool that documented a pair would be told its
			// second key is missing.
			for _, key := range strings.FieldsFunc(hint.Key, func(r rune) bool {
				return r == '/' || r == ' ' || r == ','
			}) {
				listed[key] = true
			}
			listed[hint.Key] = true
		}
	}

	for _, cmd := range spec.WithKeys(root) {
		if !listed[cmd.Key] {
			t.Errorf("%q is bound to %q and the help screen does not name it — "+
				"a help screen that is slightly wrong still renders perfectly, "+
				"and the key it omits is the one a reader came looking for.",
				cmd.Name, cmd.Key)
		}
	}
}
