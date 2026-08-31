// Package watch is the inner loop: change a line of view code, look at the
// frame.
//
// Capture is a building tool first, and a building tool has to be fast enough
// to leave running. A command you have to remember to type with two environment
// variables set is not a loop — you do it twice and then stop, and the frames
// go stale without anyone noticing.
//
// It also makes the agent loop symmetrical: the same directory of frames is
// what an agent reads to see what it just changed.
package watch

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"io"
	"io/fs"
	"os"
	"os/exec"
	"path/filepath"
	"sort"
	"strings"
	"time"
)

// Config is one watch.
type Config struct {
	// Dir is the package to watch. Its Go files are the trigger.
	Dir string
	// Capture runs the tool's capture. Typically `go test ./ui -run Capture`
	// with the frames directory in the environment.
	Capture []string
	// Frames is where the capture writes, and where the page is built from.
	Frames string
	// Out is the page file served at /.
	Out string
	// Page builds the page from the frames directory.
	Page func(framesDir string) (string, error)
	// Interval is how often the tree is checked. Zero means 300ms.
	Interval time.Duration
	// Log receives one line per rebuild. Nil discards.
	Log io.Writer
}

// Poller watches a directory tree by digesting it.
//
// Polling rather than fsnotify, and the reason is not laziness: this module has
// four direct dependencies and they are all Charm, and for a dev tool watching
// one package the difference between an inotify callback and a digest every
// 300ms is imperceptible. A dependency that buys nothing a person can feel is a
// dependency that only costs.
//
// Digesting the tree rather than comparing mtimes also makes a save that
// rewrites a file with identical content a non-event, which is what an editor
// with format-on-save does constantly.
type Poller struct {
	dir  string
	last string
}

// NewPoller starts watching dir. The first Changed call reports true, so a
// watch draws something immediately rather than after the first edit.
func NewPoller(dir string) *Poller { return &Poller{dir: dir} }

// Changed reports whether the tree's Go files differ from the last call.
func (p *Poller) Changed() (bool, error) {
	sum, err := digest(p.dir)
	if err != nil {
		return false, err
	}
	if sum == p.last {
		return false, nil
	}
	p.last = sum
	return true, nil
}

// digest is a hash of every Go file's path, size and modification time.
func digest(dir string) (string, error) {
	var names []string
	err := filepath.WalkDir(dir, func(path string, d fs.DirEntry, err error) error {
		if err != nil {
			return err
		}
		if d.IsDir() {
			// Skip the places a build writes into, or the watch retriggers on
			// its own output forever.
			switch d.Name() {
			case ".git", "node_modules", "testdata", "bin":
				return filepath.SkipDir
			}
			return nil
		}
		if strings.HasSuffix(path, ".go") {
			names = append(names, path)
		}
		return nil
	})
	if err != nil {
		return "", err
	}
	sort.Strings(names)

	h := sha256.New()
	for _, name := range names {
		info, err := os.Stat(name)
		if err != nil {
			continue // deleted between walking and stating; the next tick sees it
		}
		_, _ = fmt.Fprintf(h, "%s %d %d\n", name, info.Size(), info.ModTime().UnixNano())
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

// Build runs the capture and writes the page.
//
// A failure is returned rather than swallowed, and the caller shows it. Leaving
// the last good frames up when the code no longer compiles is a page that lies
// about what the tool does, which is worse than a page that says it is broken.
func (c Config) Build() error {
	if len(c.Capture) == 0 {
		return fmt.Errorf("watch: nothing to run — give a capture command")
	}
	cmd := exec.Command(c.Capture[0], c.Capture[1:]...)
	cmd.Env = os.Environ()

	if out, err := cmd.CombinedOutput(); err != nil {
		return fmt.Errorf("%s\n%s", strings.Join(c.Capture, " "), out)
	}
	page, err := c.Page(c.Frames)
	if err != nil {
		return err
	}
	return os.WriteFile(c.Out, []byte(page), 0o600)
}

func (c Config) interval() time.Duration {
	if c.Interval <= 0 {
		return 300 * time.Millisecond
	}
	return c.Interval
}

func (c Config) logf(format string, args ...any) {
	if c.Log == nil {
		return
	}
	_, _ = fmt.Fprintf(c.Log, format+"\n", args...)
}
