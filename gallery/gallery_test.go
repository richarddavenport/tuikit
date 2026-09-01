package gallery

import (
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/richarddavenport/tuikit/harness"
	"github.com/richarddavenport/tuikit/theme"
)

// A component that is not in the gallery is not finished.
//
// The rule enforced rather than written down. comp's own source is read for
// every exported type with a Draw* method — which is what a component IS here —
// and each one has to have an entry. Adding a component without one fails this
// test, instead of shipping something nobody has ever seen run.
//
// The reverse direction matters too: an entry for a component that no longer
// exists is a gallery lying about what the library has.
func TestEveryComponentIsInTheGallery(t *testing.T) {
	components := componentsIn(t, "../comp")
	if len(components) == 0 {
		t.Fatal("no components found in ../comp — the check is passing vacuously")
	}

	m := New(theme.Default)
	shown := map[string]bool{}
	for _, e := range m.entries {
		shown[e.Name] = true
	}

	for _, name := range components {
		if !shown[name] {
			t.Errorf("comp.%s has no gallery entry — a component that is not in the gallery is not finished", name)
		}
	}
	in := map[string]bool{}
	for _, name := range components {
		in[name] = true
	}
	for name := range shown {
		if !in[name] {
			t.Errorf("the gallery has an entry for %q, which comp does not have", name)
		}
	}
}

// componentsIn is every exported type in a package with a Draw* method.
func componentsIn(t *testing.T, dir string) []string {
	t.Helper()

	// ParseDir is deprecated for not considering build tags. comp has none,
	// and go/packages would pull a toolchain into a test that only wants the
	// exported names.
	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, nil, 0) //nolint:staticcheck // see above
	if err != nil {
		t.Fatalf("reading %s: %v", dir, err)
	}
	var out []string
	for _, pkg := range pkgs {
		for name, file := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			for _, decl := range file.Decls {
				fn, ok := decl.(*ast.FuncDecl)
				// Draw, DrawAt, DrawOn — a prefix, not an exact name.
				// comp.Menu draws at a point or on a region and has no plain
				// Draw, so an exact match let it into the library with no
				// gallery entry at all: the check that exists to stop that
				// silently stopped applying to it.
				if !ok || !strings.HasPrefix(fn.Name.Name, "Draw") || fn.Recv == nil || len(fn.Recv.List) != 1 {
					continue
				}
				if recv := receiverName(fn.Recv.List[0].Type); recv != "" && ast.IsExported(recv) {
					out = append(out, recv)
				}
			}
		}
	}
	return out
}

func receiverName(expr ast.Expr) string {
	switch t := expr.(type) {
	case *ast.StarExpr:
		return receiverName(t.X)
	case *ast.Ident:
		return t.Name
	}
	return ""
}

// Every entry has more than one state, because the happy path is the state that
// never needed a gallery. And every state says what it is FOR: "empty" explains
// itself, "clamped" does not.
func TestEveryEntryShowsMoreThanTheHappyPath(t *testing.T) {
	for _, e := range New(theme.Default).entries {
		if len(e.States) < 2 {
			t.Errorf("%s has %d states — the happy path is the one that never needed a gallery", e.Name, len(e.States))
		}
		for _, s := range e.States {
			if s.Draw == nil {
				t.Errorf("%s/%s draws nothing", e.Name, s.Name)
			}
			if s.Note == "" {
				t.Errorf("%s/%s does not say what it is for", e.Name, s.Name)
			}
		}
		if e.Summary == "" || e.From == "" {
			t.Errorf("%s does not say what it is for or where it came from", e.Name)
		}
		if len(e.Roles) == 0 {
			t.Errorf("%s names no colour roles — a tool cannot tell whether its palette can supply it", e.Name)
		}
	}
}

// The states a component is most likely to get wrong are the two a screenshot
// of the happy path never shows.
func TestTheAwkwardStatesAreCovered(t *testing.T) {
	want := map[string][]string{
		"List":     {"empty", "overflowing"},
		"LogPane":  {"empty"},
		"Pane":     {"no room"},
		"Tabs":     {"too wide"},
		"StepList": {"no room"},
		"Bar":      {"squeezed"},
		"Confirm":  {"a long body"},
	}
	for _, e := range New(theme.Default).entries {
		have := map[string]bool{}
		for _, s := range e.States {
			have[s.Name] = true
		}
		for _, name := range want[e.Name] {
			if !have[name] {
				t.Errorf("%s has no %q state", e.Name, name)
			}
		}
	}
}

// states walks the gallery as a person would: every component, every state.
func states(t *testing.T, w, h int, do func(name string, m *Model)) {
	t.Helper()
	m := New(theme.Default)
	for i, e := range m.entries {
		for j, s := range e.States {
			m.index.Select(i)
			m.state = j
			m.SetSize(w, h)
			do(e.Name+"-"+strings.ReplaceAll(s.Name, " ", "-"), m)
		}
	}
}

// A golden per component per state. This is the design system's regression net:
// a component whose empty state quietly changed is a change somebody has to
// notice, and now it is a test failure instead.
func TestGalleryGoldens(t *testing.T) {
	states(t, 132, 38, func(name string, m *Model) {
		t.Run(name, func(t *testing.T) {
			harness.Golden(t, "testdata", name, m.View())
		})
	})
}

// The gallery has to survive 80 columns like everything else it shows.
func TestGalleryAtEightyColumns(t *testing.T) {
	states(t, 80, 24, func(name string, m *Model) {
		t.Run(name, func(t *testing.T) {
			frame := m.View()
			if got := harness.Width(frame); got > 80 {
				t.Errorf("%s is %d columns wide", name, got)
			}
		})
	})
}

// Colour must not change the shape — the check democtl's tab strip failed.
func TestColourDoesNotChangeTheShape(t *testing.T) {
	for _, size := range []struct{ w, h int }{{132, 38}, {80, 24}} {
		states(t, size.w, size.h, func(name string, m *Model) {
			harness.ShapeSurvivesColour(t, name, func() string {
				fresh := New(theme.Default)
				fresh.index = m.index
				fresh.state = m.state
				fresh.SetSize(size.w, size.h)
				return fresh.View()
			})
		})
	}
}

// Capture every entry as a frame you can look at.
//
//	TUIKIT_GALLERY_FRAMES=/tmp/g go test ./gallery -run CaptureFrames
func TestCaptureFrames(t *testing.T) {
	dir := harness.Enabled("TUIKIT_GALLERY_FRAMES")
	if dir == "" {
		t.Skip("set TUIKIT_GALLERY_FRAMES to capture frames")
	}
	lipgloss.SetColorProfile(termenv.TrueColor)

	s := harness.Capture(t, dir, harness.Size(132, 38))
	states(t, 132, 38, func(name string, m *Model) { s.Shot(name, m) })
	s.Done()
}
