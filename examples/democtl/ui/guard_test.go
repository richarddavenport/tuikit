package ui

import (
	"testing"

	"github.com/richarddavenport/tuikit/guard"
	"github.com/richarddavenport/tuikit/theme"
)

// This is the whole of what a tuikit tool has to write to keep its interface
// in its vocabulary, and the scaffolder will generate it.
//
// It is also tuikit's first real consumer: until democtl existed the guards
// had only ever been run against the deploy tool, the package they were
// extracted from.
func TestTheInterfaceStaysInItsVocabulary(t *testing.T) {
	guard.Tokens(t, ".", Palette)
	guard.Glyphs(t, ".", Glyphs)
	// Every action reachable by mouse has a keyboard path. An agent cannot
	// click, and a multiplexer may eat the right-click before democtl sees it.
	guard.Reachable(t, Commands(1))
	// And nothing has taken a key that means something else in every other
	// tuikit tool — the one rule tuikit imposes rather than offers.
	guard.Reserved(t, Commands(1))
	// And the furniture the components draw on democtl's behalf — the box
	// corners, the chevrons, the scroll markers — which guard.Glyphs cannot
	// see, because they are literals in comp rather than here.
	guard.Chrome(t, theme.DefaultChrome, Glyphs)
	// And the other direction: a chrome character typed out HERE rather than
	// read from the chrome. guard.Glyphs cannot see it, because the character
	// is allowed — being in the box set is the point of it. What is wrong is
	// the source, and a hardcoded one stays put when the box set changes.
	guard.Furniture(t, ".", theme.DefaultChrome)
}
