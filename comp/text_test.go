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
