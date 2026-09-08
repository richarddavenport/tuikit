// Package news answers "what has tuikit done since this tool last looked".
//
// # Why decisions and not versions
//
// A tool resolves tuikit through `replace github.com/richarddavenport/tuikit
// => ../tuikit`, so there is no version, no tag and no module cache — a pull in
// tuikit changes the tool's behavior with no upgrade event to hang release
// notes on. Something else has to be the marker.
//
// `design/decisions.md` already is one. It is numbered, append-only, and every
// entry was written to be read by whoever comes next; the record existed and
// simply had no reader. A tool records the number it reconciled with and this
// prints what came after.
//
// # Why not the git log
//
// It catches more and says less. A tool's author needs to know that config
// moved to ~/.config, not that a lint rule about preallocation was satisfied,
// and a list where those two lines look the same is a list nobody reads twice.
// The cost is real and stated in the output: a change that got no decision does
// not appear here, so the inventory is the other half of the answer.
package news

import (
	"bufio"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

// Decision is one numbered entry in design/decisions.md.
type Decision struct {
	Number int
	Title  string
	// Lede is the first paragraph, which is where these are written to say what
	// changed before they say why.
	Lede string
}

// heading matches "## 33. Config goes in ~/.config/<tool>/, and it is copied".
var heading = regexp.MustCompile(`^##\s+(\d+)\.\s+(.*\S)\s*$`)

// marker matches the line a tool keeps to say where it is up to. Case and
// punctuation are loose because it is written by hand in prose.
var marker = regexp.MustCompile(`(?i)reconciled\s+with\s+tuikit\s+through\s+decision\s+(\d+)`)

// Read parses the decisions file.
func Read(path string) ([]Decision, error) {
	f, err := os.Open(path) //nolint:gosec // a path the caller chose
	if err != nil {
		return nil, fmt.Errorf("news: %w", err)
	}
	defer func() { _ = f.Close() }()

	var out []Decision
	var lede []string
	// collecting is true between a heading and the end of its first paragraph,
	// so a decision's opening lines are captured and the rest of its several
	// hundred words are not.
	collecting := false

	sc := bufio.NewScanner(f)
	sc.Buffer(make([]byte, 0, 64*1024), 1024*1024)
	for sc.Scan() {
		line := sc.Text()
		if m := heading.FindStringSubmatch(line); m != nil {
			if len(out) > 0 {
				out[len(out)-1].Lede = strings.Join(lede, " ")
			}
			n, _ := strconv.Atoi(m[1])
			out = append(out, Decision{Number: n, Title: m[2]})
			lede, collecting = nil, true
			continue
		}
		if !collecting {
			continue
		}
		if strings.TrimSpace(line) == "" {
			if len(lede) > 0 {
				collecting = false
			}
			continue
		}
		lede = append(lede, strings.TrimSpace(line))
	}
	if err := sc.Err(); err != nil {
		return nil, fmt.Errorf("news: %w", err)
	}
	if len(out) > 0 {
		out[len(out)-1].Lede = strings.Join(lede, " ")
	}
	if len(out) == 0 {
		return nil, fmt.Errorf("news: %s has no numbered decisions", path)
	}
	return out, nil
}

// Since returns the decisions after n, in order.
func Since(ds []Decision, n int) []Decision {
	var out []Decision
	for _, d := range ds {
		if d.Number > n {
			out = append(out, d)
		}
	}
	return out
}

// Latest is the highest decision number, which is what a tool records once it
// has read everything above.
func Latest(ds []Decision) int {
	n := 0
	for _, d := range ds {
		if d.Number > n {
			n = d.Number
		}
	}
	return n
}

// Marker reads "reconciled with tuikit through decision N" out of a file.
//
// In the tool's AGENTS.md rather than a dotfile, because the number is only
// useful to whoever is about to change the tool and that is the file they are
// told to read first. A dotfile is a fact nobody sees until it is wrong.
func Marker(path string) (int, bool) {
	body, err := os.ReadFile(path) //nolint:gosec // a path the caller chose
	if err != nil {
		return 0, false
	}
	m := marker.FindSubmatch(body)
	if m == nil {
		return 0, false
	}
	n, err := strconv.Atoi(string(m[1]))
	return n, err == nil
}

// replaceLine matches the go.mod directive a tool uses to resolve tuikit.
var replaceLine = regexp.MustCompile(`(?m)^\s*replace\s+github\.com/richarddavenport/tuikit\s+=>\s+(\S+)`)

// Checkout finds the tuikit source a tool is actually building against, by
// reading the replace directive out of its go.mod.
//
// Asking the go.mod rather than a flag or an environment variable, because the
// answer already exists there and cannot be wrong: it is the path the compiler
// used. A second way to say where tuikit is, is a second way to be told about a
// tuikit the tool does not build against.
//
// The path is resolved relative to the go.mod, since that is what Go does.
func Checkout(gomod string) (string, bool) {
	body, err := os.ReadFile(gomod) //nolint:gosec // a path the caller chose
	if err != nil {
		return "", false
	}
	m := replaceLine.FindSubmatch(body)
	if m == nil {
		return "", false
	}
	path := string(m[1])
	if !filepath.IsAbs(path) {
		path = filepath.Join(filepath.Dir(gomod), path)
	}
	abs, err := filepath.Abs(path)
	if err != nil {
		return "", false
	}
	return abs, true
}
