package term_test

import (
	"os"
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/term"
)

// TestParse works the decision without a terminal, which is the point of Parse
// being exported: the replies below are what real emulators actually send, and
// none of them needs a tty to be checked.
func TestParse(t *testing.T) {
	for _, tc := range []struct {
		name  string
		reply string
		want  term.Graphics
	}{
		{"kitty answers its own query", "\x1b_Gi=31;OK\x1b\\\x1b[?62;c", term.Kitty},
		{"ghostty, same shape", "\x1b_Gi=31;OK\x1b\\\x1b[?62;22c", term.Kitty},

		// An error reply has the same shape as a success, with a code where the
		// OK goes. Matching "OK" loosely would read a refusal as support.
		{"kitty refuses", "\x1b_Gi=31;ENOENT:image not found\x1b\\\x1b[?62;c", term.None},
		{"kitty refuses, EINVAL", "\x1b_Gi=31;EINVAL:bad payload\x1b\\\x1b[?6c", term.None},
		{"foot: DA1 with attribute 4", "\x1b[?62;4;22c", term.Sixel},
		{"xterm with sixel built in", "\x1b[?63;1;2;4;6;9;15;22c", term.Sixel},
		{"alacritty: DA1, no 4", "\x1b[?6c", term.None},
		{"nothing at all — the timeout case", "", term.None},
		{"garbage", "hello", term.None},

		// 14 contains a 4 and 41 starts with one. A contains check would call
		// both of these Sixel; splitting on ';' is what stops it.
		{"attribute 14 is not attribute 4", "\x1b[?62;14c", term.None},
		{"attribute 41 is not attribute 4", "\x1b[?62;41c", term.None},
	} {
		t.Run(tc.name, func(t *testing.T) {
			if got := term.Parse(tc.reply); got != tc.want {
				t.Errorf("Parse(%q) = %v, want %v", tc.reply, got, tc.want)
			}
		})
	}
}

// TestDetectIsNoneUnderTest is the property the goldens rest on.
//
// A test has no controlling terminal, so Detect returns None and every frame a
// test captures is the cell drawing. That is what makes the fallback the
// most-tested path in the library rather than the least, and it holds by
// construction rather than by anyone remembering to set something.
func TestDetectIsNoneUnderTest(t *testing.T) {
	if got := term.Detect(); got != term.None {
		t.Fatalf("Detect() = %v under test, want None — goldens would move", got)
	}
}

func TestEnvOverride(t *testing.T) {
	for _, tc := range []struct {
		set  string
		want term.Graphics
	}{
		{"kitty", term.Kitty},
		{"sixel", term.Sixel},
		{"KITTY", term.Kitty}, // case is not a decision
		{" sixel ", term.Sixel},
		{"none", term.None},
		{"nonsense", term.None},
	} {
		t.Run(tc.set, func(t *testing.T) {
			t.Setenv(term.EnvOverride, tc.set)
			if got := term.Detect(); got != tc.want {
				t.Errorf("with %s=%q, Detect() = %v, want %v",
					term.EnvOverride, tc.set, got, tc.want)
			}
		})
	}
}

// TestQueryOnAPipeIsNone: a pipe is not a terminal, and asking one a question
// that only a terminal can answer must not block.
func TestQueryOnAPipeIsNone(t *testing.T) {
	r, w, err := os.Pipe()
	if err != nil {
		t.Fatal(err)
	}
	defer r.Close()
	defer w.Close()
	if got := term.Query(w, term.DefaultTimeout); got != term.None {
		t.Errorf("Query(pipe) = %v, want None", got)
	}
}

// TestQueryBytesCarriesAPayload is the regression for the bug that reported
// "none" on Ghostty: a=q REQUIRES image data, and without it kitty and Ghostty
// answer nothing — which is indistinguishable from a terminal that cannot draw.
func TestQueryBytesCarriesAPayload(t *testing.T) {
	q := term.QueryBytes
	for _, want := range []string{"a=q", "s=1,v=1", "t=d", "f=24", ";AAAA", "\x1b[c"} {
		if !strings.Contains(q, want) {
			t.Errorf("query is missing %q: %q", want, q)
		}
	}
}

func TestGraphicsString(t *testing.T) {
	for g, want := range map[term.Graphics]string{
		term.None: "none", term.Sixel: "sixel", term.Kitty: "kitty",
	} {
		if got := g.String(); got != want {
			t.Errorf("String() = %q, want %q", got, want)
		}
	}
}

// TestDA1Complete is the read's stopping condition, and it has to be precise.
//
// The kitty answer arrives BEFORE the DA1 answer and carries arbitrary text, so
// "contains a c" stops early on a refusal like `ENOENT:image not found` — and
// stopping early means never seeing whether DA1 advertised Sixel. A terminal
// that supports both would be read as supporting neither.
func TestDA1Complete(t *testing.T) {
	complete := []string{
		"\x1b[?62;4c",
		"\x1b_Gi=31;OK\x1b\\\x1b[?64;1;2;4;6;17;18;21;22c",
		"\x1b_Gi=31;ENOENT:image not found\x1b\\\x1b[?62;4c",
	}
	partial := []string{
		"",
		"\x1b_Gi=31;ENOENT:image not found\x1b\\", // kitty answered, DA1 has not
		"\x1b_Gi=31;OK\x1b\\",                     // ditto
		"\x1b[?62;4",                              // mid-reply
	}
	for _, s := range complete {
		if !term.Parseable(s) {
			t.Errorf("%q should be a complete reply", s)
		}
	}
	for _, s := range partial {
		if term.Parseable(s) {
			t.Errorf("%q should not be complete — the read would stop early", s)
		}
	}
}

// TestBothProtocolsAnswered is iTerm2's actual reply: it refuses the kitty
// query AND advertises Sixel as attribute 4. Reading only the first half loses
// the Sixel answer entirely.
func TestBothProtocolsAnswered(t *testing.T) {
	reply := "\x1b_Gi=31;invalid payload\x1b\\\x1b[?64;1;2;4;6;17;18;21;22c"
	if got := term.Parse(reply); got != term.Sixel {
		t.Errorf("Parse = %v, want Sixel — DA1 advertised attribute 4", got)
	}
}
