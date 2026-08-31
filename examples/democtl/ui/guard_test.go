package ui

import (
	"testing"

	"github.com/richarddavenport/tuikit/guard"
)

// This is the whole of what a tuikit tool has to write to keep its interface in
// its vocabulary, and the scaffolder will generate it.
//
// It is also tuikit's first real consumer: until democtl existed the guards had
// only ever been run against swarmctl, the package they were extracted from.
func TestTheInterfaceStaysInItsVocabulary(t *testing.T) {
	guard.Tokens(t, ".", Palette)
	guard.Glyphs(t, ".", Glyphs)
}
