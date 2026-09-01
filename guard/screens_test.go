package guard

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// screenPkg writes a throwaway package for the guard to read.
func screenPkg(t *testing.T, body string) string {
	t.Helper()
	dir := t.TempDir()
	src := "package ui\n\ntype screen int\n\n" + body + "\n"
	if err := os.WriteFile(filepath.Join(dir, "model.go"), []byte(src), 0o600); err != nil {
		t.Fatal(err)
	}
	return dir
}

// The enumeration comes from the SOURCE, so a screen added at the end of the
// block is checked — which is exactly what a hand-written
// `for s := first; s <= last; s++` misses.
func TestScreensFiresOnAScreenNothingDraws(t *testing.T) {
	dir := screenPkg(t, `const (
	screenDashboard screen = iota
	screenLogs
	screenRun
)`)
	msgs := run(t, func(rec T) { Screens(rec, dir, "screen", func(v int) bool { return v < 2 }) })

	if len(msgs) != 1 {
		t.Fatalf("got %d complaints, want 1:\n%s", len(msgs), strings.Join(msgs, "\n"))
	}
	if !strings.Contains(msgs[0], "screenRun") {
		t.Errorf("the complaint does not name the screen: %s", msgs[0])
	}
}

func TestScreensIsQuietWhenEveryScreenIsDrawn(t *testing.T) {
	dir := screenPkg(t, `const (
	screenDashboard screen = iota
	screenLogs
	screenRun
)`)
	msgs := run(t, func(rec T) { Screens(rec, dir, "screen", func(int) bool { return true }) })

	if len(msgs) > 0 {
		t.Errorf("complained about a clean package:\n%s", strings.Join(msgs, "\n"))
	}
}

// Constants of another type in the same file are not screens.
func TestScreensIgnoresOtherTypes(t *testing.T) {
	dir := screenPkg(t, `type pane int

const (
	paneList pane = iota
	paneDetail
)

const (
	screenDashboard screen = iota
)`)
	var asked []int
	msgs := run(t, func(rec T) {
		Screens(rec, dir, "screen", func(v int) bool { asked = append(asked, v); return true })
	})

	if len(asked) != 1 || asked[0] != 0 {
		t.Errorf("asked about %v, want just the one screen", asked)
	}
	if len(msgs) > 0 {
		t.Errorf("complained: %s", strings.Join(msgs, "\n"))
	}
}

// Explicit values are read, because a tool is entitled to number its screens.
func TestScreensReadsExplicitValues(t *testing.T) {
	dir := screenPkg(t, `const (
	screenDashboard screen = 0
	screenRun       screen = 7
)`)
	var asked []int
	run(t, func(rec T) {
		Screens(rec, dir, "screen", func(v int) bool { asked = append(asked, v); return true })
	})

	if len(asked) != 2 || asked[1] != 7 {
		t.Errorf("asked about %v, want [0 7]", asked)
	}
}

// A value it cannot evaluate is REPORTED, not skipped. A screen silently not
// checked is the thing this exists to prevent, so failing to read one has to be
// louder than reading it and finding nothing wrong.
func TestScreensReportsAValueItCannotRead(t *testing.T) {
	dir := screenPkg(t, `const base = 4

const (
	screenDashboard screen = iota + base
)`)
	msgs := run(t, func(rec T) { Screens(rec, dir, "screen", func(int) bool { return true }) })

	if len(msgs) != 1 || !strings.Contains(msgs[0], "cannot be read") {
		t.Errorf("got %v, want a complaint about an unreadable value", msgs)
	}
}

// A guard that enumerates nothing passes silently, which is how everyone comes
// to believe a rule is on when it is not.
func TestScreensWithNoConstantsFails(t *testing.T) {
	dir := screenPkg(t, `const notAScreen = 3`)
	msgs := run(t, func(rec T) { Screens(rec, dir, "screen", func(int) bool { return true }) })

	if len(msgs) != 1 || !strings.Contains(msgs[0], "enumerates nothing") {
		t.Errorf("got %v, want a fatal about enumerating nothing", msgs)
	}
}
