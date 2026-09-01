// Package scaffold generates a tuikit tool.
//
// Decision 1 draws the line this package sits on: mechanism goes in the
// library, policy goes in the scaffolder. Anything that must work for every
// tool is a function you call; anything that is merely a good default is
// generated code the tool OWNS and can edit on day one. Repo infrastructure —
// a Makefile, a CI workflow, a linter config — cannot be a dependency, which is
// the whole reason this exists.
//
// Seeded from what azctl's migration actually needed rather than from a guess,
// which is why it came after that migration. The order of what it writes is
// the order of how much trouble each thing saves, and that order is in
// azctl/design/migrating-to-tuikit.md.
package scaffold

import (
	"embed"
	"fmt"
	"io/fs"
	"os"
	"path/filepath"
	"sort"
	"strings"
	"text/template"
)

// all: because embed skips files starting with a dot by default, and half of
// what a repo needs is dotfiles — .golangci.yml, .github/workflows. A
// scaffolder that silently omits CI is worse than one that does not run.
//
//go:embed all:templates
var templates embed.FS

// Tool is what a new tool is called and where it lives.
type Tool struct {
	// Name is the binary's name, and the directory's.
	Name string
	// Module is the go.mod path. Defaults to github.com/<user>/<name> when a
	// user is given, and to the bare name otherwise.
	Module string
	// Short is the one line that goes in the README, the CLI's help and the
	// manifest an agent reads. One sentence, said once.
	Short string

	// Tuikit is the path a `replace` points at.
	//
	// tuikit is private and unpublished, so a generated tool needs one to
	// build. A relative path by default, so a checkout beside tuikit works on
	// any machine rather than only the one it was written on — which is the
	// mistake azctl's go.mod made first.
	Tuikit string
	// GoVersion is the go directive and the version CI installs.
	GoVersion string
	// Linter is the pinned golangci-lint, run through `go run` so it is the
	// same version whether or not anything is installed.
	Linter string
}

// Defaults fills in what was not given.
func (t Tool) Defaults() Tool {
	if t.Module == "" {
		t.Module = t.Name
	}
	if t.Short == "" {
		t.Short = "a tuikit tool"
	}
	if t.Tuikit == "" {
		t.Tuikit = "../tuikit"
	}
	if t.GoVersion == "" {
		t.GoVersion = "1.25"
	}
	if t.Linter == "" {
		t.Linter = "v2.12.0"
	}
	return t
}

// Env is the tool's name as an environment variable prefix: DEMOCTL, AZCTL.
func (t Tool) Env() string {
	return strings.ToUpper(strings.ReplaceAll(t.Name, "-", "_"))
}

// Title is the tool's name for prose: capitalised, since a sentence starts with
// one and a binary name does not.
func (t Tool) Title() string {
	if t.Name == "" {
		return ""
	}
	return strings.ToUpper(t.Name[:1]) + t.Name[1:]
}

// Write generates the tool into dir/<name> and returns what it wrote.
//
// It refuses to write into a directory that already has anything in it. A
// scaffolder that overwrites is a scaffolder nobody runs twice, and "it emptied
// my repo" is not a thing to find out by doing it.
func (t Tool) Write(dir string) ([]string, error) {
	t = t.Defaults()
	if err := t.valid(); err != nil {
		return nil, err
	}

	root := filepath.Join(dir, t.Name)
	if entries, err := os.ReadDir(root); err == nil && len(entries) > 0 {
		return nil, fmt.Errorf("%s already exists and is not empty", root)
	}

	names, err := files()
	if err != nil {
		return nil, err
	}

	var written []string
	for _, name := range names {
		// The template's path IS the output path, with <name> substituted and
		// the .tmpl suffix dropped — so adding a file to the scaffolder is
		// adding a file, not adding a file and a line in a list.
		out := filepath.Join(root, strings.ReplaceAll(
			strings.TrimSuffix(strings.TrimPrefix(name, "templates/"), ".tmpl"),
			"NAME", t.Name))

		body, err := t.render(name)
		if err != nil {
			return written, err
		}
		if err := os.MkdirAll(filepath.Dir(out), 0o755); err != nil {
			return written, err
		}
		mode := os.FileMode(0o644)
		if strings.HasSuffix(out, ".sh") {
			mode = 0o755
		}
		if err := os.WriteFile(out, body, mode); err != nil {
			return written, err
		}
		written = append(written, out)
	}
	return written, nil
}

func (t Tool) valid() error {
	if t.Name == "" {
		return fmt.Errorf("a tool needs a name")
	}
	for _, r := range t.Name {
		if r >= 'a' && r <= 'z' || r >= '0' && r <= '9' || r == '-' {
			continue
		}
		// The name is a binary, a directory, a package path and a Go
		// identifier's neighbour. Rejecting the rest here beats a compile error
		// in generated code somebody did not write.
		return fmt.Errorf("%q: a tool's name is lower-case letters, digits and dashes", t.Name)
	}
	return nil
}

func (t Tool) render(name string) ([]byte, error) {
	src, err := templates.ReadFile(name)
	if err != nil {
		return nil, err
	}
	// Delimiters that no shell, YAML or Makefile uses, because these templates
	// generate all three and $() and {{ }} are all taken.
	tmpl, err := template.New(name).Delims("<<", ">>").Parse(string(src))
	if err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	var b strings.Builder
	if err := tmpl.Execute(&b, t); err != nil {
		return nil, fmt.Errorf("%s: %w", name, err)
	}
	return []byte(b.String()), nil
}

// files is every template, sorted — so a generated tool is written in the same
// order every time and two runs can be compared.
func files() ([]string, error) {
	var out []string
	err := fs.WalkDir(templates, "templates", func(path string, d fs.DirEntry, err error) error {
		if err != nil || d.IsDir() {
			return err
		}
		out = append(out, path)
		return nil
	})
	sort.Strings(out)
	return out, err
}
