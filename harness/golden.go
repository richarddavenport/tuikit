package harness

import (
	"flag"
	"os"
	"path/filepath"
	"strings"
)

// update regenerates goldens instead of comparing against them.
//
// A flag rather than an environment variable, and explicit rather than
// automatic: a golden that rewrites itself when it fails is a test that always
// passes. Run `go test ./... -update-goldens` when a layout change is intended.
var update = flag.Bool("update-goldens", false, "rewrite golden frames instead of comparing")

// Golden compares a frame against testdata/<name>.golden, colour stripped.
//
// Stripped because a golden holds the SHAPE. A diff in a pull request wants to
// show that a box moved, and one full of escape sequences shows nothing a
// reader can act on. The .ansi files from a capture keep the colour.
//
// This is what turns "did I break the layout" into an ordinary test failure.
// Across the four tools this framework came from there are 23,000 lines of
// tests and not one asserts on a rendered frame; every layout bug they have had
// was found by somebody looking.
func Golden(t T, dir, name, frame string) {
	t.Helper()

	got := Strip(frame)
	path := filepath.Join(dir, name+".golden")

	if *update {
		if err := os.MkdirAll(dir, 0o755); err != nil {
			t.Fatalf("harness: %v", err)
		}
		if err := os.WriteFile(path, []byte(got), 0o600); err != nil {
			t.Fatalf("harness: %v", err)
		}
		return
	}

	wantBytes, err := os.ReadFile(path)
	if os.IsNotExist(err) {
		t.Errorf("no golden for %q — run with -update-goldens to write %s", name, path)
		return
	}
	if err != nil {
		t.Fatalf("harness: %v", err)
	}
	if want := string(wantBytes); got != want {
		t.Errorf("%s differs from its golden:\n%s", name, diff(want, got))
	}
}

// diff reports the first differing line with a little context, and the counts.
//
// Not a full diff: a frame is a picture, and the useful question is "which row
// changed and how wide is it now", which the first difference answers. A
// forty-line unified diff of a forty-line frame is the frame twice.
func diff(want, got string) string {
	wantLines, gotLines := strings.Split(want, "\n"), strings.Split(got, "\n")

	var b strings.Builder
	if len(wantLines) != len(gotLines) {
		b.WriteString(sprintf("  height: was %d lines, now %d\n", len(wantLines), len(gotLines)))
	}
	for i := 0; i < len(wantLines) && i < len(gotLines); i++ {
		if wantLines[i] == gotLines[i] {
			continue
		}
		b.WriteString(sprintf("  line %d\n", i+1))
		b.WriteString(sprintf("    was  (%3d cols) %q\n", widthOf(wantLines[i]), wantLines[i]))
		b.WriteString(sprintf("    now  (%3d cols) %q\n", widthOf(gotLines[i]), gotLines[i]))
		return b.String()
	}
	if b.Len() == 0 {
		b.WriteString("  trailing lines differ\n")
	}
	return b.String()
}
