package harness

import (
	"strings"

	"github.com/richarddavenport/tuikit/comp"
)

// Hints fails when a frame advertises a key that does nothing.
//
// A footer is a promise. `q abort` is read by someone deciding whether to press
// it, and a footer naming a key nothing handles is a lie the reader has no way
// to check — the interface still renders perfectly, so no golden and no
// assertion notices.
//
// It presses each advertised key against a copy of the state the caller has
// already driven to, and reports the ones that change nothing at all.
//
// # What this does NOT catch
//
// A key that does something OTHER than what its label says. That was the bug
// which prompted this (azctl's runner promised `q abort (the running step
// finishes)` while `q` cancelled the run and quit the program): `q` did plenty,
// it just did not do what the footer said. No mechanical check reaches that,
// because it is a claim about English. Decision 32 is the discipline that
// covers it — a comment or a label describing behaviour is an untested
// assertion, and the answer is a test named after the sentence.
//
// So this catches the lesser sibling: the DEAD advertised key. Worth having,
// and worth knowing it is not the whole of the problem.
//
// # False positives
//
// A key that legitimately does nothing in this state — `G follow` when already
// following — will be reported. That is why the caller passes the hints rather
// than the frame being scraped: test the state where the key is live, or leave
// the inert one out and say why in the test.
func Hints(t T, build func() Driver, hints ...comp.Hint) {
	t.Helper()

	for _, hint := range hints {
		for _, key := range keysIn(hint.Key) {
			// A fresh model per key: pressing them in sequence means the
			// second key is judged against whatever the first one did.
			m := build()
			redraw(m)
			before := m.View()

			Press(m, key)
			if m.View() == before {
				t.Errorf("the frame offers %q (%s) and pressing it changes nothing — "+
					"a footer naming a key that does nothing is a lie a reader believes",
					key, hint.Label)
			}
		}
	}
}

// keysIn splits a hint that names several keys for one action: "j/k", "↑↓".
//
// The arrow pair is one glyph run rather than a separated list, so it is split
// per rune. Splitting every hint per rune instead would turn "tab" into three
// keys, which is why this looks at what the string is made of.
func keysIn(hint string) []string {
	if fields := strings.FieldsFunc(hint, func(r rune) bool {
		return r == '/' || r == ' ' || r == ','
	}); len(fields) > 1 {
		return fields
	}
	if arrows := arrowsIn(hint); len(arrows) > 0 {
		return arrows
	}
	return []string{hint}
}

func arrowsIn(hint string) []string {
	var out []string
	for _, r := range hint {
		switch r {
		case '↑', '↓', '←', '→':
			out = append(out, string(r))
		default:
			return nil
		}
	}
	if len(out) < 2 {
		return nil
	}
	return out
}
