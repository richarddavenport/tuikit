package fuzzy_test

import (
	"testing"

	"github.com/richarddavenport/tuikit/fuzzy"
)

func TestMatchFindsASubsequence(t *testing.T) {
	for _, tc := range []struct {
		query, text string
		want        []int
	}{
		{"dsk", "disk usage", []int{0, 2, 3}},
		{"env", "switch environment", []int{7, 8, 9}},
		{"", "anything", nil},
		{"disk usage", "disk usage", []int{0, 1, 2, 3, 4, 5, 6, 7, 8, 9}},
	} {
		got, ok := fuzzy.Match(tc.query, tc.text)
		if !ok {
			t.Errorf("%q did not match %q", tc.query, tc.text)
			continue
		}
		if len(got.At) != len(tc.want) {
			t.Errorf("%q in %q matched at %v, want %v", tc.query, tc.text, got.At, tc.want)
			continue
		}
		for i := range got.At {
			if got.At[i] != tc.want[i] {
				t.Errorf("%q in %q matched at %v, want %v", tc.query, tc.text, got.At, tc.want)
				break
			}
		}
	}
}

func TestMatchRejects(t *testing.T) {
	for _, tc := range [][2]string{
		{"xyz", "disk usage"},
		{"ksd", "disk"},   // right letters, wrong order
		{"diskk", "disk"}, // one letter too many
	} {
		if _, ok := fuzzy.Match(tc[0], tc[1]); ok {
			t.Errorf("%q should not match %q", tc[0], tc[1])
		}
	}
}

func TestMatchIsCaseInsensitive(t *testing.T) {
	if _, ok := fuzzy.Match("DSK", "disk usage"); !ok {
		t.Error("upper-case query did not match lower-case text")
	}
	// ...but an exact-case match scores higher, so `Env` prefers the one that
	// is actually capitalized.
	exact, _ := fuzzy.Match("Env", "Environment")
	loose, _ := fuzzy.Match("Env", "environment")
	if exact.Score <= loose.Score {
		t.Errorf("exact case scored %d, loose %d — exact should win", exact.Score, loose.Score)
	}
}

// TestARunIsNotBrokenUp is the regression for the backward pass this replaced.
//
// Pulling each matched letter as far right as it would go moved the v of "env"
// onto the v in "service" — turning a perfect three-letter run at the start
// into a match with a thirteen-rune gap, and underlining letters nobody typed.
func TestARunIsNotBrokenUp(t *testing.T) {
	got, ok := fuzzy.Match("env", "env of this service")
	if !ok {
		t.Fatal("no match")
	}
	for i, want := range []int{0, 1, 2} {
		if got.At[i] != want {
			t.Fatalf("matched at %v, want [0 1 2] — the run at the start", got.At)
		}
	}
}

// TestTheWholeAlignmentIsWeighed, which is why this is a table and not a rule
// applied per letter. Two adjacent letters at the start beat the same two
// letters spread across the string, even though the spread pair lands on a
// word boundary.
func TestTheWholeAlignmentIsWeighed(t *testing.T) {
	got, ok := fuzzy.Match("ge", "generate env")
	if !ok {
		t.Fatal("no match")
	}
	if got.At[0] != 0 || got.At[1] != 1 {
		t.Errorf("matched at %v, want [0 1] — someone typing \"ge\" means generate", got.At)
	}
}

// TestShapeBeatsPosition: consecutive letters at a word start are what someone
// typing an abbreviation means.
func TestShapeBeatsPosition(t *testing.T) {
	good, _ := fuzzy.Match("du", "disk usage")      // two word starts
	poor, _ := fuzzy.Match("du", "shutdown lambda") // buried
	if good.Score <= poor.Score {
		t.Errorf("word-start match scored %d, buried match %d", good.Score, poor.Score)
	}

	run, _ := fuzzy.Match("res", "restart")         // consecutive
	split, _ := fuzzy.Match("res", "refresh scope") // scattered
	if run.Score <= split.Score {
		t.Errorf("consecutive scored %d, scattered %d", run.Score, split.Score)
	}
}

func TestRankOrdersAndFilters(t *testing.T) {
	texts := []string{"switch environment", "edit environments", "disk usage", "env of this service"}
	got := fuzzy.Rank("env", texts)
	if len(got) != 3 {
		t.Fatalf("got %d matches, want 3 (disk usage has no v after an e-n)", len(got))
	}
	// "env of this service" starts with the word itself.
	if texts[got[0].Index] != "env of this service" {
		t.Errorf("best match was %q", texts[got[0].Index])
	}
	for _, r := range got {
		if texts[r.Index] == "disk usage" {
			t.Error("disk usage should not have matched")
		}
	}
}

// TestRankIsStable. With nothing typed, the palette must stay in the order the
// caller built it — its tiers are already sorted nearest-first.
func TestRankIsStable(t *testing.T) {
	texts := []string{"one", "two", "three", "four"}
	got := fuzzy.Rank("", texts)
	if len(got) != len(texts) {
		t.Fatalf("empty query returned %d of %d", len(got), len(texts))
	}
	for i, r := range got {
		if r.Index != i {
			t.Errorf("position %d holds input %d — the order moved", i, r.Index)
		}
	}
}

// TestRunesNotBytes: At is used to style characters, so an index into the
// middle of a multi-byte character would underline half of one.
func TestRunesNotBytes(t *testing.T) {
	got, ok := fuzzy.Match("né", "un néant")
	if !ok {
		t.Fatal("no match")
	}
	for _, i := range got.At {
		if i >= len([]rune("un néant")) {
			t.Errorf("index %d is past the end in runes", i)
		}
	}
	if got.At[0] != 3 || got.At[1] != 4 {
		t.Errorf("matched at %v, want [3 4] — rune indices, not byte offsets", got.At)
	}
}
