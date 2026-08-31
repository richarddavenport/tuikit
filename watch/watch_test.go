package watch

import (
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func write(t *testing.T, path, body string) {
	t.Helper()
	if err := os.MkdirAll(filepath.Dir(path), 0o755); err != nil {
		t.Fatal(err)
	}
	if err := os.WriteFile(path, []byte(body), 0o600); err != nil {
		t.Fatal(err)
	}
}

// The first call reports a change, so a watch draws something immediately
// rather than after the first edit.
func TestAWatchHasSomethingToShowBeforeYouTouchAnything(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.go"), "package a\n")

	p := NewPoller(dir)
	changed, err := p.Changed()
	if err != nil {
		t.Fatal(err)
	}
	if !changed {
		t.Error("the first poll reported nothing to build")
	}
}

func TestAnEditIsAChangeAndStillnessIsNot(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "a.go")
	write(t, path, "package a\n")

	p := NewPoller(dir)
	if _, err := p.Changed(); err != nil {
		t.Fatal(err)
	}
	if changed, _ := p.Changed(); changed {
		t.Error("nothing changed, but the poller said it did")
	}

	time.Sleep(10 * time.Millisecond)
	write(t, path, "package a\n\nvar x = 1\n")

	if changed, _ := p.Changed(); !changed {
		t.Error("an edit was not noticed")
	}
}

// A new file counts, and so does a deletion. A watch that only notices edits
// misses the two things that change a package's shape.
func TestAddingAndRemovingFilesAreChanges(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.go"), "package a\n")

	p := NewPoller(dir)
	p.Changed()

	write(t, filepath.Join(dir, "b.go"), "package a\n")
	if changed, _ := p.Changed(); !changed {
		t.Error("a new file was not noticed")
	}

	if err := os.Remove(filepath.Join(dir, "b.go")); err != nil {
		t.Fatal(err)
	}
	if changed, _ := p.Changed(); !changed {
		t.Error("a deleted file was not noticed")
	}
}

// The directories a build writes into are skipped, or the watch retriggers on
// its own output forever.
func TestOutputDirectoriesDoNotTriggerTheWatch(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.go"), "package a\n")

	p := NewPoller(dir)
	p.Changed()

	for _, out := range []string{"testdata/x.go", "bin/y.go", ".git/z.go"} {
		write(t, filepath.Join(dir, out), "package a\n")
	}
	if changed, _ := p.Changed(); changed {
		t.Error("writing into testdata, bin or .git retriggered the watch")
	}
}

// Only Go files. A watch that fires on the page it just wrote never settles.
func TestOnlyGoFilesCount(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.go"), "package a\n")

	p := NewPoller(dir)
	p.Changed()

	write(t, filepath.Join(dir, "page.html"), "<p>hello</p>")
	if changed, _ := p.Changed(); changed {
		t.Error("a non-Go file triggered a rebuild")
	}
}

// --- building and serving -----------------------------------------------

func TestBuildRunsTheCaptureThenWritesThePage(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "page.html")
	marker := filepath.Join(dir, "captured")

	c := Config{
		Capture: []string{"sh", "-c", "echo yes > " + marker},
		Frames:  dir,
		Out:     out,
		Page:    func(string) (string, error) { return "<p>page</p>", nil },
	}
	if err := c.Build(); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(marker); err != nil {
		t.Error("the capture did not run")
	}
	body, err := os.ReadFile(out)
	if err != nil || string(body) != "<p>page</p>" {
		t.Errorf("page = %q, %v", body, err)
	}
}

// A failing capture must not leave the page it wrote last time.
func TestAFailedCaptureDoesNotWriteAPage(t *testing.T) {
	dir := t.TempDir()
	out := filepath.Join(dir, "page.html")
	write(t, out, "<p>stale</p>")

	c := Config{
		Capture: []string{"sh", "-c", "echo 'ui/view.go:42: undefined: thing' >&2; exit 1"},
		Frames:  dir,
		Out:     out,
		Page:    func(string) (string, error) { return "<p>fresh</p>", nil },
	}
	err := c.Build()
	if err == nil {
		t.Fatal("a failing capture reported success")
	}
	if !strings.Contains(err.Error(), "undefined: thing") {
		t.Errorf("the compiler's message was swallowed: %v", err)
	}
	if body, _ := os.ReadFile(out); string(body) != "<p>stale</p>" {
		t.Error("the page was overwritten by a failed build")
	}
}

func TestBuildWithNothingToRunSaysSo(t *testing.T) {
	err := Config{Page: func(string) (string, error) { return "", nil }}.Build()
	if err == nil || !strings.Contains(err.Error(), "give a capture command") {
		t.Errorf("got %v", err)
	}
}

// The whole loop: serve, then change a file and watch the page follow.
func TestTheLoopRebuildsAndTheVersionMoves(t *testing.T) {
	dir := t.TempDir()
	src := filepath.Join(dir, "a.go")
	write(t, src, "package a\n")

	var n int
	c := Config{
		Dir:      dir,
		Capture:  []string{"true"},
		Frames:   dir,
		Out:      filepath.Join(dir, "page.html"),
		Interval: 20 * time.Millisecond,
		Page: func(string) (string, error) {
			n++
			return fmt.Sprintf("<p>build %d</p>", n), nil
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr := make(chan string, 1)
	go func() { _ = c.Run(ctx, "127.0.0.1:0", func(a string) { addr <- a }) }()

	base := <-addr
	first := waitFor(t, base, "build 1")
	if !strings.Contains(first, "location.reload") {
		t.Error("the page carries no reload script")
	}

	time.Sleep(30 * time.Millisecond)
	write(t, src, "package a\n\nvar x = 1\n")
	waitFor(t, base, "build 2")

	if v := get(t, base+"/version"); v == "0" {
		t.Errorf("the version never moved: %s", v)
	}
}

// A build that fails puts the error on the page, rather than leaving the last
// good frames up to describe a tool that no longer exists.
func TestABrokenBuildPutsTheErrorOnThePage(t *testing.T) {
	dir := t.TempDir()
	write(t, filepath.Join(dir, "a.go"), "package a\n")

	c := Config{
		Dir:      dir,
		Capture:  []string{"sh", "-c", "echo 'ui/view.go:42: undefined: thing' >&2; exit 1"},
		Frames:   dir,
		Out:      filepath.Join(dir, "page.html"),
		Interval: 20 * time.Millisecond,
		Page:     func(string) (string, error) { return "<p>never</p>", nil },
		Log:      io.Discard,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	addr := make(chan string, 1)
	go func() { _ = c.Run(ctx, "127.0.0.1:0", func(a string) { addr <- a }) }()

	page := waitFor(t, <-addr, "undefined: thing")
	for _, want := range []string{"The capture failed", "not shown", "location.reload"} {
		if !strings.Contains(page, want) {
			t.Errorf("the error page is missing %q", want)
		}
	}
	if strings.Contains(page, "<p>never</p>") {
		t.Error("a page was rendered from a failed build")
	}
}

func waitFor(t *testing.T, base, substr string) string {
	t.Helper()
	deadline := time.Now().Add(4 * time.Second)
	var last string
	for time.Now().Before(deadline) {
		last = get(t, base)
		if strings.Contains(last, substr) {
			return last
		}
		time.Sleep(20 * time.Millisecond)
	}
	t.Fatalf("never saw %q; last page was:\n%s", substr, last)
	return ""
}

func get(t *testing.T, url string) string {
	t.Helper()
	// A connection refused while the server is still coming up is ordinary, so
	// an empty body means "not yet" and waitFor tries again.
	resp, err := http.Get(url) //nolint:noctx // a test against its own loopback server
	if err != nil {
		return ""
	}
	defer resp.Body.Close()
	body, _ := io.ReadAll(resp.Body)
	return string(body)
}
