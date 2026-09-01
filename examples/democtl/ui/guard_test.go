package ui

import (
	"testing"

	"github.com/richarddavenport/tuikit/guard"
	"github.com/richarddavenport/tuikit/theme"
)

// This is the whole of what a tuikit tool has to write to keep its interface in
// its vocabulary, and the scaffolder will generate it.
//
// It is also tuikit's first real consumer: until democtl existed the guards had
// only ever been run against swarmctl, the package they were extracted from.
func TestTheInterfaceStaysInItsVocabulary(t *testing.T) {
	guard.Tokens(t, ".", Palette)
	guard.Glyphs(t, ".", Glyphs)
	// Every action reachable by mouse has a keyboard path. An agent cannot
	// click, and a multiplexer may eat the right-click before democtl sees it.
	guard.Reachable(t, Commands(1))
	// And the furniture the components draw on democtl's behalf — the box
	// corners, the chevrons, the scroll markers — which guard.Glyphs cannot
	// see, because they are literals in comp rather than here.
	guard.Chrome(t, theme.DefaultChrome, Glyphs)
}
