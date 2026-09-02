package app

import "testing"

const (
	dashboard Screen = iota
	logs
	run
	confirm
)

// The zero value is a tool sitting on its first screen, which is the constant a
// Screen enumeration starts at.
func TestTheZeroValueIsTheFirstScreen(t *testing.T) {
	var s Stack
	if s.Current() != dashboard {
		t.Errorf("a fresh stack is on %d", s.Current())
	}
	if s.Back() {
		t.Error("back moved on an empty stack — esc on the first screen is a key nothing took")
	}
}

// The root is never popped: a stack with nothing under it is a tool showing no
// screen, which no tool means. Every caller would otherwise write the same
// guard, and one of them would forget.
func TestBackOnTheRootDoesNotMove(t *testing.T) {
	var s Stack
	s.Push(dashboard, "browse")

	if s.Back() {
		t.Error("back popped the root")
	}
	if s.Current() != dashboard || s.Depth() != 1 {
		t.Errorf("the root is %d at depth %d", s.Current(), s.Depth())
	}
}

// Back returns to where you came FROM, which is the whole point. A hardcoded
// constant is right only until a screen can be reached two ways.
func TestBackReturnsToWhereYouCameFrom(t *testing.T) {
	for _, from := range []Screen{dashboard, run} {
		var s Stack
		s.Push(dashboard, "browse")
		if from != dashboard {
			s.Push(from, "")
		}
		s.Push(logs, "logs api_gateway")

		if !s.Back() {
			t.Fatal("back did not move")
		}
		if s.Current() != from {
			t.Errorf("opened logs from %d and came back to %d", from, s.Current())
		}
	}
}

// A temporary screen never becomes a back-stop. The flag is read when something
// is pushed OVER it, so a modal does not have to know in advance whether
// anything will be.
func TestATemporaryScreenIsNotSomethingToGoBackTo(t *testing.T) {
	var s Stack
	s.Push(dashboard, "")
	s.PushTemp(confirm, "")
	s.Push(run, "deploy api_gateway")

	if !s.Back() {
		t.Fatal("back did not move")
	}
	if s.Current() != dashboard {
		t.Errorf("back landed on %d — a confirmation is not a place to return to", s.Current())
	}
}

// Dismissing a modal is an ordinary back, because until something is pushed
// over it, it is where you are.
func TestBackDismissesATemporaryScreen(t *testing.T) {
	var s Stack
	s.Push(dashboard, "")
	s.PushTemp(confirm, "")

	if !s.Back() {
		t.Fatal("back did not dismiss the modal")
	}
	if s.Current() != dashboard {
		t.Errorf("dismissing the modal landed on %d", s.Current())
	}
}

func TestBackToRootUnwinds(t *testing.T) {
	var s Stack
	s.Push(dashboard, "")
	s.Push(logs, "")
	s.Push(run, "")

	s.BackToRoot()
	if s.Current() != dashboard || s.Depth() != 1 {
		t.Errorf("root is %d at depth %d", s.Current(), s.Depth())
	}
}

// Replace is a screen BECOMING another, not opening one, so there is nothing
// new to come back to.
func TestReplaceDoesNotDeepenTheStack(t *testing.T) {
	var s Stack
	s.Push(dashboard, "")
	s.Push(logs, "logs api")
	s.Replace(logs, "logs web")

	if s.Depth() != 2 {
		t.Errorf("depth is %d, want 2", s.Depth())
	}
	if e, _ := s.At(); e.Label != "logs web" {
		t.Errorf("the label is %q", e.Label)
	}
	s.Back()
	if s.Current() != dashboard {
		t.Errorf("back from a replaced screen landed on %d", s.Current())
	}
}

// The path is what an agent reads to find out where it is — as command lines,
// because that is a thing it can act on.
func TestThePathReadsAsCommandLines(t *testing.T) {
	var s Stack
	s.Push(dashboard, "browse")
	s.Push(logs, "logs api_gateway")

	path := s.Path()
	if len(path) != 2 || path[0].Label != "browse" || path[1].Label != "logs api_gateway" {
		t.Errorf("path is %+v", path)
	}
	// A copy, so a caller cannot reach in and rewrite history.
	path[0].Label = "something else"
	if again := s.Path(); again[0].Label != "browse" {
		t.Error("Path handed out the stack's own slice")
	}
}

// Depth is what a footer asks before it promises esc goes back.
func TestDepthSaysWhetherThereIsAWayBack(t *testing.T) {
	var s Stack
	if s.Depth() != 0 {
		t.Errorf("a fresh stack has depth %d", s.Depth())
	}
	s.Push(dashboard, "")
	s.Push(logs, "")
	if s.Depth() != 2 {
		t.Errorf("depth is %d, want 2", s.Depth())
	}
}

// BackTo is what clicking a breadcrumb means. The stack had Back and
// BackToRoot and nothing between.
func TestBackToADepth(t *testing.T) {
	var s Stack
	s.Push(1, "democtl")
	s.Push(2, "services")
	s.Push(3, "api_gateway")
	s.Push(4, "Logs")

	if !s.BackTo(1) {
		t.Fatal("BackTo(1) reported nothing to do")
	}
	if s.Depth() != 2 {
		t.Errorf("the stack is %d deep, want 2", s.Depth())
	}
	if got := s.Current(); got != 2 {
		t.Errorf("we are on screen %v, want 2", got)
	}
}

// Out of range does nothing rather than clamping: guessing which screen a
// reader meant is how a click ends up somewhere they did not point.
func TestBackToOutOfRangeDoesNothing(t *testing.T) {
	var s Stack
	s.Push(1, "root")
	s.Push(2, "second")

	for _, depth := range []int{-1, 2, 99} {
		if s.BackTo(depth) {
			t.Errorf("BackTo(%d) claimed to move", depth)
		}
		if s.Depth() != 2 {
			t.Fatalf("BackTo(%d) changed the stack to %d deep", depth, s.Depth())
		}
	}
	// And going where you already are is not a move either.
	if s.BackTo(1) {
		t.Error("BackTo on the current depth claimed to move")
	}
}
