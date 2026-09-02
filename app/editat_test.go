package app_test

import (
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/app"
)

func key(t tea.KeyType, runes ...rune) tea.KeyMsg {
	return tea.KeyMsg{Type: t, Runes: runes}
}

func TestEditAtInsertsAtTheCaret(t *testing.T) {
	got, cur, ok := app.EditAt("helo", 3, key(tea.KeyRunes, 'l'))
	if !ok || got != "hello" || cur != 4 {
		t.Errorf("got %q caret %d ok %v, want \"hello\" 4 true", got, cur, ok)
	}
}

// TestEditAtSpaceOnce is the bug Edit was written to stop: Bubble Tea sets
// Runes to a space for KeySpace as well as setting the type, so a handler that
// does both puts two in.
func TestEditAtSpaceOnce(t *testing.T) {
	got, cur, _ := app.EditAt("a", 1, key(tea.KeySpace, ' '))
	if got != "a " || cur != 2 {
		t.Errorf("got %q caret %d, want \"a \" 2", got, cur)
	}
}

func TestEditAtBackspaceAndDelete(t *testing.T) {
	if got, cur, _ := app.EditAt("hello", 3, key(tea.KeyBackspace)); got != "helo" || cur != 2 {
		t.Errorf("backspace gave %q caret %d, want \"helo\" 2", got, cur)
	}
	if got, cur, _ := app.EditAt("hello", 3, key(tea.KeyDelete)); got != "helo" || cur != 3 {
		t.Errorf("delete gave %q caret %d, want \"helo\" 3", got, cur)
	}
}

// TestEditAtAtTheEdges: both still TAKE the key. A backspace on an empty prompt
// that falls through would be acted on by the screen underneath.
func TestEditAtAtTheEdges(t *testing.T) {
	if got, cur, ok := app.EditAt("", 0, key(tea.KeyBackspace)); got != "" || cur != 0 || !ok {
		t.Errorf("got %q %d %v, want \"\" 0 true", got, cur, ok)
	}
	if got, cur, ok := app.EditAt("ab", 2, key(tea.KeyDelete)); got != "ab" || cur != 2 || !ok {
		t.Errorf("got %q %d %v, want \"ab\" 2 true", got, cur, ok)
	}
}

func TestEditAtMoves(t *testing.T) {
	for _, tc := range []struct {
		key  tea.KeyType
		from int
		want int
	}{
		{tea.KeyLeft, 3, 2}, {tea.KeyLeft, 0, 0},
		{tea.KeyRight, 3, 4}, {tea.KeyRight, 5, 5},
		{tea.KeyHome, 3, 0}, {tea.KeyEnd, 1, 5},
		{tea.KeyCtrlA, 3, 0}, {tea.KeyCtrlE, 1, 5},
	} {
		got, cur, ok := app.EditAt("hello", tc.from, key(tc.key))
		if got != "hello" || cur != tc.want || !ok {
			t.Errorf("%v from %d: %q caret %d, want caret %d", tc.key, tc.from, got, cur, tc.want)
		}
	}
}

// TestEditAtRunesNotBytes. A caret counted in bytes lands inside a multi-byte
// character, and the next backspace cuts it in half.
func TestEditAtRunesNotBytes(t *testing.T) {
	const s = "héllo" // é is two bytes
	got, cur, _ := app.EditAt(s, 2, key(tea.KeyBackspace))
	if got != "hllo" || cur != 1 {
		t.Errorf("got %q caret %d, want \"hllo\" 1", got, cur)
	}
	if got, _, _ := app.EditAt("日本", 1, key(tea.KeyDelete)); got != "日" {
		t.Errorf("got %q, want 日", got)
	}
}

func TestEditAtCtrlU(t *testing.T) {
	// Readline's ctrl-u: everything BEFORE the caret, not the whole line.
	got, cur, _ := app.EditAt("hello world", 6, key(tea.KeyCtrlU))
	if got != "world" || cur != 0 {
		t.Errorf("got %q caret %d, want \"world\" 0", got, cur)
	}
}

func TestEditAtCtrlW(t *testing.T) {
	for _, tc := range []struct {
		in    string
		at    int
		want  string
		caret int
	}{
		{"hello world", 11, "hello ", 6},
		{"hello world  ", 13, "hello ", 6}, // eats the spaces, then the word
		{"one", 3, "", 0},
		{"", 0, "", 0},
	} {
		got, cur, _ := app.EditAt(tc.in, tc.at, key(tea.KeyCtrlW))
		if got != tc.want || cur != tc.caret {
			t.Errorf("ctrl-w on %q at %d = %q caret %d, want %q %d",
				tc.in, tc.at, got, cur, tc.want, tc.caret)
		}
	}
}

// TestEditAtDoesNotTakeEnterOrEscape: they mean different things in different
// tools, so they stay the caller's.
func TestEditAtDoesNotTakeEnterOrEscape(t *testing.T) {
	for _, k := range []tea.KeyType{tea.KeyEnter, tea.KeyEsc, tea.KeyTab, tea.KeyUp} {
		if _, _, ok := app.EditAt("x", 1, key(k)); ok {
			t.Errorf("EditAt took %v", k)
		}
	}
}

// TestEditAtClampsAStrayCaret. A caller that resized or replaced the text can
// hand in a caret past the end; that must not panic.
func TestEditAtClampsAStrayCaret(t *testing.T) {
	if got, cur, _ := app.EditAt("ab", 99, key(tea.KeyRunes, 'c')); got != "abc" || cur != 3 {
		t.Errorf("got %q caret %d, want \"abc\" 3", got, cur)
	}
	if got, cur, _ := app.EditAt("ab", -5, key(tea.KeyRunes, 'c')); got != "cab" || cur != 1 {
		t.Errorf("got %q caret %d, want \"cab\" 1", got, cur)
	}
}
