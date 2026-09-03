package guard_test

import (
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/guard"
	"github.com/richarddavenport/tuikit/spec"
)

func root(cmds ...spec.Command) spec.Command {
	return spec.Command{Name: "democtl", Commands: cmds}
}

// The trap this exists to stop: q meaning something that is not leaving.
func TestReservedCatchesAKeyMeaningSomethingElse(t *testing.T) {
	rec := &fake{}
	guard.Reserved(rec, root(
		spec.Command{Name: "queue", Key: "q", Target: comp.Name("row")},
	))

	if len(rec.msgs) != 1 {
		t.Fatalf("reported %d, want 1: %v", len(rec.msgs), rec.msgs)
	}
	if got := rec.msgs[0]; !strings.Contains(got, "queue") || !strings.Contains(got, "reserved") {
		t.Errorf("the message does not say what is wrong: %s", got)
	}
}

// And the ordinary case passes, or the guard is indistinguishable from a clean
// package.
func TestReservedAllowsTheMeaningItReserves(t *testing.T) {
	rec := &fake{}
	guard.Reserved(rec, root(
		spec.Command{Name: "quit", Key: "q"},
		spec.Command{Name: "back", Key: "esc"},
		spec.Command{Name: "help", Key: "?"},
		spec.Command{Name: "deploy", Key: "D", Target: comp.Name("row")},
	))

	if len(rec.msgs) != 0 {
		t.Errorf("reported a legitimate binding: %v", rec.msgs)
	}
}

// ? is help, not "anything at all".
func TestReservedHoldsQuestionMarkToHelp(t *testing.T) {
	rec := &fake{}
	guard.Reserved(rec, root(spec.Command{Name: "quit", Key: "?"}))

	if len(rec.msgs) != 1 {
		t.Errorf("a leaving command took the help key without complaint: %v", rec.msgs)
	}
}

// Nothing declares ctrl+c. It is handled before a tool sees it, so a command
// claiming it is describing a binding it does not have — decision 32's shape.
func TestNothingMayDeclareCtrlC(t *testing.T) {
	rec := &fake{}
	guard.Reserved(rec, root(spec.Command{Name: "quit", Key: "ctrl+c"}))

	if len(rec.msgs) != 1 {
		t.Errorf("ctrl+c was declarable: %v", rec.msgs)
	}
}
