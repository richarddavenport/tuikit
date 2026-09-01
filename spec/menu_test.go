package spec

import (
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/comp"
)

func targeted() Command {
	return Command{Name: "democtl", Commands: []Command{
		{Name: "logs", Short: "View logs", Target: "services.row", Key: "L", Run: func(Call) int { return OK }},
		{Name: "deploy", Short: "Deploy", Target: "services.row", Key: "D", Run: func(Call) int { return OK }},
		{Name: "restart", Short: "Restart", Target: "nodes.row", Key: "R", Run: func(Call) int { return OK }},
		{Name: "status", Short: "no menu entry", Run: func(Call) int { return OK }},
		{Name: "secret", Short: "hidden", Target: "services.row", Key: "S", Hidden: true, Run: func(Call) int { return OK }},
	}}
}

// The menu for a region IS the commands that target it, so one added to the CLI
// appears and one removed disappears, without anybody remembering either.
func TestAMenuIsTheCommandsThatTargetTheRegion(t *testing.T) {
	got := MenuFor(targeted(), "services.row")

	if len(got) != 2 {
		t.Fatalf("the menu has %d entries: %+v", len(got), got)
	}
	if got[0].Key != "L" || got[0].Label != "View logs" {
		t.Errorf("first entry is %+v", got[0])
	}
	if joined := comp.Hints(got...); !strings.Contains(joined, "D Deploy") {
		t.Errorf("got %q", joined)
	}
}

func TestACommandWithNoTargetIsInNoMenu(t *testing.T) {
	for _, region := range []comp.Name{"services.row", "nodes.row"} {
		for _, hint := range MenuFor(targeted(), region) {
			if hint.Label == "no menu entry" {
				t.Errorf("an untargeted command is in %s's menu", region)
			}
		}
	}
}

func TestAHiddenCommandIsInNoMenu(t *testing.T) {
	for _, hint := range MenuFor(targeted(), "services.row") {
		if hint.Key == "S" {
			t.Error("a hidden command is in the menu")
		}
	}
}

func TestEachRegionGetsItsOwn(t *testing.T) {
	if got := MenuFor(targeted(), "nodes.row"); len(got) != 1 || got[0].Key != "R" {
		t.Errorf("nodes.row's menu is %+v", got)
	}
	if got := MenuFor(targeted(), "nothing.here"); len(got) != 0 {
		t.Errorf("a region nothing targets has a menu: %+v", got)
	}
}

// An agent cannot click, and a right-click may never arrive: a multiplexer that
// captures the mouse gets the event first and the application never learns it
// happened. A command whose only path is a menu is broken for anyone inside
// herdr or tmux, and the tool cannot detect that to warn them.
func TestACommandInAMenuNeedsAKeyboardPath(t *testing.T) {
	if got := Unreachable(targeted()); len(got) != 0 {
		t.Errorf("a reachable tree reports %v", got)
	}

	broken := targeted()
	broken.Commands[1].Key = ""
	got := Unreachable(broken)
	if len(got) != 1 || got[0] != "democtl deploy" {
		t.Errorf("got %v, want [democtl deploy]", got)
	}
}

// A hidden command still needs one: hidden means unadvertised, not unreachable.
func TestEvenAHiddenCommandNeedsAKeyboardPath(t *testing.T) {
	broken := targeted()
	broken.Commands[4].Key = ""
	if got := Unreachable(broken); len(got) != 1 {
		t.Errorf("a hidden command with no key reports %v", got)
	}
}
