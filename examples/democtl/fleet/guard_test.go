package fleet

import (
	"testing"

	"github.com/richarddavenport/tuikit/guard"
)

// The engine knows the domain and nothing about terminals — no colour, no
// width, no keys, no framework, and no tuikit (decision 22).
//
// The rule is worth as much as it is enforced, and it is mechanically
// checkable, so it is checked. This file is what the scaffolder generates
// alongside the UI package's guards.
func TestTheEngineHasNeverHeardOfATerminal(t *testing.T) {
	guard.Engine(t, ".")
}
