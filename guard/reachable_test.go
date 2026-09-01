package guard

import (
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/spec"
)

// A guard that never fires is indistinguishable from a clean package, so it is
// run against a tree that breaks the rule and one that does not.
func TestReachableFiresOnAMouseOnlyCommand(t *testing.T) {
	rec := &recorder{}
	Reachable(rec, spec.Command{Name: "t", Commands: []spec.Command{
		{Name: "logs", Target: "services.row", Key: "L"},
		{Name: "deploy", Target: "services.row"},
	}})

	if len(rec.msgs) != 1 {
		t.Fatalf("reported %d problems: %v", len(rec.msgs), rec.msgs)
	}
	if !strings.Contains(rec.msgs[0], "t deploy") {
		t.Errorf("the failure does not name the command: %q", rec.msgs[0])
	}
	if !strings.Contains(rec.msgs[0], "agent cannot click") {
		t.Errorf("the failure does not say why: %q", rec.msgs[0])
	}
}

func TestReachableIsQuietWhenEveryActionHasAKey(t *testing.T) {
	rec := &recorder{}
	Reachable(rec, spec.Command{Name: "t", Commands: []spec.Command{
		{Name: "logs", Target: "services.row", Key: "L"},
		{Name: "status"}, // no target, so no menu, so no key needed
	}})

	if len(rec.msgs) != 0 {
		t.Errorf("a reachable tree reported %v", rec.msgs)
	}
}
