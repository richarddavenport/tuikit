package guard

import (
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/spec"
)

func tool() spec.Command {
	return spec.Command{Name: "democtl", Commands: []spec.Command{
		{Name: "logs", Short: "follow a service", Key: "L", Target: "services.row"},
		{Name: "deploy", Short: "deploy a service", Key: "D", Target: "services.row"},
		{Name: "version", Short: "print the version"},
	}}
}

func TestKeysPassesWhenTheHelpNamesThemAll(t *testing.T) {
	rec := &recorder{}
	Keys(rec, tool(), []comp.KeySection{
		{Name: "A service", Keys: []comp.Hint{
			{Key: "L", Label: "logs"},
			{Key: "D", Label: "deploy"},
		}},
		// Screen keys the spec has never heard of are fine, and must be: they
		// are not commands and never will be.
		{Name: "Everywhere", Keys: []comp.Hint{
			{Key: "j/k", Label: "move"},
			{Key: "q", Label: "quit"},
		}},
	})
	if len(rec.msgs) != 0 {
		t.Errorf("a complete help screen failed: %v", rec.msgs)
	}
}

// The guard has to FIRE, or it is indistinguishable from a guard that never
// runs.
func TestKeysFiresOnAMissingBinding(t *testing.T) {
	rec := &recorder{}
	Keys(rec, tool(), []comp.KeySection{
		{Name: "A service", Keys: []comp.Hint{{Key: "L", Label: "logs"}}},
	})
	if len(rec.msgs) != 1 {
		t.Fatalf("got %d complaints, want 1: %v", len(rec.msgs), rec.msgs)
	}
	if !strings.Contains(rec.msgs[0], "deploy") || !strings.Contains(rec.msgs[0], `"D"`) {
		t.Errorf("the complaint does not name the command and its key: %s", rec.msgs[0])
	}
}

// A hint naming several keys for one action documents both. Otherwise a tool
// that wrote "j/k" would be told its second key is undocumented.
func TestAPairOfKeysInOneHintCountsAsBoth(t *testing.T) {
	root := spec.Command{Name: "t", Commands: []spec.Command{
		{Name: "up", Key: "k", Target: "row"},
		{Name: "down", Key: "j", Target: "row"},
	}}
	rec := &recorder{}
	Keys(rec, root, []comp.KeySection{
		{Keys: []comp.Hint{{Key: "j/k", Label: "move"}}},
	})
	if len(rec.msgs) != 0 {
		t.Errorf("a hint naming both keys was not accepted: %v", rec.msgs)
	}
}

// A command with no Key is not a keyboard path, so it is not this guard's
// business — Reachable is the one that asks whether it should have one.
func TestKeysIgnoresCommandsWithNoKey(t *testing.T) {
	rec := &recorder{}
	Keys(rec, tool(), []comp.KeySection{
		{Keys: []comp.Hint{{Key: "L"}, {Key: "D"}}},
	})
	if len(rec.msgs) != 0 {
		t.Errorf("a command with no Key was demanded: %v", rec.msgs)
	}
}
