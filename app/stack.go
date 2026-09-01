package app

// Stack is where you are and how you got there.
//
// # Why a map was not enough
//
// Screens started as `map[Screen]View`, which answers "what does this screen
// draw" and nothing else. So every tool wrote its own back:
//
//	case "esc":
//		if m.screen != screenDashboard {
//			m.screen = screenDashboard
//		}
//
// Back as a hardcoded constant, correct only because every screen happened to
// be reached from the dashboard. Open the logs from a run and it goes to the
// wrong place, silently — the same failure shape as an owner ID that means "row
// 3 of the screen" rather than "row 3 of the list". And azctl did not even have
// that: its runner was a second tea.NewProgram, so a playbook run lost the
// browser entirely.
//
// # What this is not
//
// Not a router. bubblyui ports Vue Router into a Go TUI — path patterns,
// specificity scoring, nested outlets, guards, ~5,000 lines — and its routed
// components are permanently-mounted singletons that never receive a lifecycle
// call, its history's state slot is tested and never populated, and its route
// table cannot be listed from outside. Fifteen thousand lines routing nothing
// the framework knows about.
//
// A screen is not a URL. It is a value the tool already has, so parameters are
// typed and free and there is nothing to serialise.
//
// # The label
//
// An entry carries how you would ASK for this screen from the command line —
// "logs api_gateway". Not to rebuild it from: to say where you are. That makes
// the position readable by an agent, printable in a manifest, and reachable by
// a capture script that would otherwise press three keys to arrive somewhere.
//
// A string rather than a spec.Call, so this package still knows nothing about
// declarations. A tool that wants the structure keeps the Call itself.
type Stack struct {
	entries []Entry
}

// Entry is one screen on the stack.
type Entry struct {
	Screen Screen
	// Label is the command line that would open this screen. Optional.
	Label string
	// temp marks a screen that should not be gone back TO — a modal, a picker,
	// a confirmation. The flag is read when something is pushed OVER it, which
	// is soda's arrangement and the one that works: a transient screen never
	// becomes a back-stop, and it does not have to know in advance whether
	// anything will be pushed over it.
	temp bool
}

// Push opens a screen, remembering where you were.
func (s *Stack) Push(screen Screen, label string) { s.push(Entry{Screen: screen, Label: label}) }

// PushTemp opens a screen that back should skip over.
func (s *Stack) PushTemp(screen Screen, label string) {
	s.push(Entry{Screen: screen, Label: label, temp: true})
}

func (s *Stack) push(e Entry) {
	// Pushing over a temporary screen replaces it rather than stacking on it.
	if n := len(s.entries); n > 0 && s.entries[n-1].temp {
		s.entries[n-1] = e
		return
	}
	s.entries = append(s.entries, e)
}

// Replace swaps the current screen without remembering it, for a screen that
// becomes another rather than opening one.
func (s *Stack) Replace(screen Screen, label string) {
	if n := len(s.entries); n > 0 {
		s.entries[n-1] = Entry{Screen: screen, Label: label}
		return
	}
	s.Push(screen, label)
}

// Back goes to the previous screen and reports whether it moved.
//
// The bool is what a key handler needs: esc on the first screen is not "go
// back", it is a key nothing took, and a tool may want to quit on it.
//
// The root is never popped. A stack with nothing under it is a tool showing no
// screen at all, which no tool means, and every caller would otherwise write
// the same Depth() > 1 guard.
func (s *Stack) Back() bool {
	if len(s.entries) <= 1 {
		return false
	}
	s.entries = s.entries[:len(s.entries)-1]
	return true
}

// BackToRoot unwinds to the first screen.
func (s *Stack) BackToRoot() {
	if len(s.entries) > 1 {
		s.entries = s.entries[:1]
	}
}

// Current is the screen showing. The zero value of Screen on an empty stack,
// which is a tool's first constant and therefore its first screen.
func (s *Stack) Current() Screen {
	if len(s.entries) == 0 {
		return Screen(0)
	}
	return s.entries[len(s.entries)-1].Screen
}

// At reports the current entry, and whether there is one.
func (s *Stack) At() (Entry, bool) {
	if len(s.entries) == 0 {
		return Entry{}, false
	}
	return s.entries[len(s.entries)-1], true
}

// Depth is how many screens deep, so a footer can say whether esc goes back or
// quits without asking twice.
func (s *Stack) Depth() int { return len(s.entries) }

// Path is the way in, oldest first — what an agent reads to find out where it
// is, and what a manifest prints.
func (s *Stack) Path() []Entry {
	out := make([]Entry, len(s.entries))
	copy(out, s.entries)
	return out
}
