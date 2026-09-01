package comp

import (
	"reflect"
	"testing"
)

func TestTruncateCountsColumnsNotRunes(t *testing.T) {
	for _, tc := range []struct {
		name, in string
		w        int
		want     string
	}{
		{"short enough is untouched", "api_gateway", 20, "api_gateway"},
		{"exactly fitting is untouched", "api", 3, "api"},
		{"too long ends in an ellipsis", "api_gateway", 6, "api_g…"},
		{"one column is just the ellipsis", "api", 1, "…"},
		{"nothing fits in nothing", "api", 0, ""},
		{"a wide rune counts two", "世界世界", 5, "世界…"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Truncate(tc.in, tc.w)
			if got != tc.want {
				t.Errorf("Truncate(%q, %d) = %q, want %q", tc.in, tc.w, got, tc.want)
			}
			if Width(got) > tc.w {
				t.Errorf("Truncate(%q, %d) = %q, which is %d columns", tc.in, tc.w, got, Width(got))
			}
		})
	}
}

func TestWrapBreaksAtSpacesAndThenAnywhere(t *testing.T) {
	for _, tc := range []struct {
		name, in string
		w        int
		want     []string
	}{
		{"at spaces", "no suitable node available", 12, []string{"no suitable", "node", "available"}},
		{"a word longer than the line", "unreasonablylongword", 8, []string{"unreason", "ablylong", "word"}},
		{"nothing wraps to nothing", "", 10, nil},
	} {
		t.Run(tc.name, func(t *testing.T) {
			got := Wrap(tc.in, tc.w)
			if !reflect.DeepEqual(got, tc.want) {
				t.Errorf("Wrap(%q, %d) = %q, want %q", tc.in, tc.w, got, tc.want)
			}
			for _, line := range got {
				if Width(line) > tc.w {
					t.Errorf("%q is %d columns, over the %d asked for", line, Width(line), tc.w)
				}
			}
		})
	}
}

func TestWrapFirstIsTheFirstLineAndNoWider(t *testing.T) {
	if got, want := WrapFirst("all of it fits", 20), "all of it fits"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
	if got := WrapFirst("", 20); got != "" {
		t.Errorf("nothing wrapped to %q", got)
	}
	got := WrapFirst("no suitable node: memory reservation exceeds every node", 20)
	if Width(got) > 20 {
		t.Errorf("%q is %d columns", got, Width(got))
	}
}

// A newline in the input is where the author said the break goes, and a wrapper
// that collapses it has thrown away the one piece of formatting they were able
// to express. A two-paragraph confirmation body used to come out as one run-on
// paragraph.
func TestWrapKeepsTheBreaksItWasGiven(t *testing.T) {
	got := Wrap("docker stack rm old\n\nIt is guarded: a check decided this needs doing.", 30)
	want := []string{
		"docker stack rm old",
		"",
		"It is guarded: a check decided",
		"this needs doing.",
	}
	if len(got) != len(want) {
		t.Fatalf("got %d lines, want %d:\n%q", len(got), len(want), got)
	}
	for i := range want {
		if got[i] != want[i] {
			t.Errorf("line %d is %q, want %q", i, got[i], want[i])
		}
	}
}

// A line long enough to wrap still wraps; keeping breaks is not the same as
// only breaking where it was told to.
func TestWrapStillWrapsWithinAParagraph(t *testing.T) {
	got := Wrap("one two three four five six seven", 12)
	for _, line := range got {
		if Width(line) > 12 {
			t.Errorf("%q is wider than 12", line)
		}
	}
	if len(got) < 3 {
		t.Errorf("got %d lines, expected the text to wrap: %q", len(got), got)
	}
}
