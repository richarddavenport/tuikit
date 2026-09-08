package gallery

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/richarddavenport/tuikit/comp"
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
	// The reverse direction asks a WEAKER question on purpose: does comp have a
	// type by this name at all. An entry for a type that no longer exists is a
	// gallery lying about the library; an entry for one that exists and does
	// not draw is comp.Layout, which computes rects and is worth looking at
	// precisely because you cannot see it in any other way.
	types := typesIn(t, "../comp")
	for name := range shown {
		if !types[name] {
			t.Errorf("the gallery has an entry for %q, which comp does not have", name)
		}
	}
}

// typesIn is every exported type in a package, drawing or not.
func typesIn(t *testing.T, dir string) map[string]bool {
	t.Helper()

	pkgs, err := parser.ParseDir(token.NewFileSet(), dir, nil, 0) //nolint:staticcheck // see componentsIn
	if err != nil {
		t.Fatalf("gallery: %v", err)
	}
	out := map[string]bool{}
	for _, pkg := range pkgs {
		for name, file := range pkg.Files {
			if strings.HasSuffix(name, "_test.go") {
				continue
			}
			for _, decl := range file.Decls {
				gen, ok := decl.(*ast.GenDecl)
				if !ok || gen.Tok != token.TYPE {
					continue
				}
				for _, spec := range gen.Specs {
					if ts, ok := spec.(*ast.TypeSpec); ok && ast.IsExported(ts.Name.Name) {
						out[ts.Name.Name] = true
					}
				}
			}
		}
	}
	return out
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
			t.Errorf("%s names no color roles — a tool cannot tell whether its palette can supply it", e.Name)
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
			harness.Golden(t, "testdata", name, view(m))
		})
	})
}

// A state's name becomes a golden's FILENAME, and a Go module zip refuses some
// characters outright.
//
// `gallery/testdata/Spinner-the-caller's-frames.golden` made `go get
// github.com/richarddavenport/tuikit@v0.1.0` fail for everyone with
//
//	create zip: malformed file path: invalid char '\''
//
// The tests passed, `make check` passed, and CI passed — the file is perfectly
// legal on disk. It only broke for people who were not us, which is the class
// of bug a repository finds on the day it goes public and not before.
//
// So the names are checked here, where they are written, rather than the files
// being checked where they land.
func TestAStateNameCanBeAFilename(t *testing.T) {
	// What golang.org/x/mod/zip rejects. Spaces are absent because states()
	// turns them into hyphens, so they never reach a filename; everything
	// listed here does.
	const illegal = `'"\!*[]:<>|?` + "`"
	m := New(theme.Default)
	for _, e := range m.Entries() {
		for _, s := range e.States {
			for _, r := range e.Name + s.Name {
				if strings.ContainsRune(illegal, r) {
					t.Errorf("%s/%s contains %q, which cannot be in a filename inside a Go module",
						e.Name, s.Name, r)
				}
			}
		}
	}
}

// The gallery has to survive 80 columns like everything else it shows.
func TestGalleryAtEightyColumns(t *testing.T) {
	states(t, 80, 24, func(name string, m *Model) {
		t.Run(name, func(t *testing.T) {
			frame := view(m)
			if got := harness.Width(frame); got > 80 {
				t.Errorf("%s is %d columns wide", name, got)
			}
		})
	})
}

// Color must not change the shape — the check democtl's tab strip failed.
func TestColorDoesNotChangeTheShape(t *testing.T) {
	for _, size := range []struct{ w, h int }{{132, 38}, {80, 24}} {
		states(t, size.w, size.h, func(name string, m *Model) {
			harness.ShapeSurvivesColor(t, name, func() string {
				fresh := New(theme.Default)
				fresh.index = m.index
				fresh.state = m.state
				fresh.SetSize(size.w, size.h)
				return view(fresh)
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
	states(t, 132, 38, func(name string, m *Model) { s.Shot(name, run(m)) })
	s.Done()
}

// No component draws outside the rect it was given.
//
// The canvas guarantees nothing lands outside the CANVAS — Set clips, so a
// coordinate past the edge is one that does not exist. It guarantees nothing
// about a component staying inside the RECT it was handed, and a component that
// overruns paints over its neighbor rather than failing: the frame is still
// well-formed, the goldens still pass, and the pane beside it is simply wrong.
//
// This is issue 6's guard.Width, and it lives here rather than in `guard`
// because the gallery is already the complete list of components — held closed
// by TestEveryComponentIsInTheGallery — and already draws each of them into a
// rect in every state it has. A guard in `guard` would need its own list of
// components and its own way to build them, and a second list is a list that
// drifts.
//
// It found comp.Toast on its first run: Min is 24 columns and clamp() raised
// the width UP to it inside a pane 10 wide, so the toast drew 114 cells over
// whatever was beside it.
func TestNoComponentDrawsOutsideItsRect(t *testing.T) {
	// A rect with room around it on all four sides, so an overrun in any
	// direction has somewhere to land where it can be seen.
	box := comp.Rect{X: 6, Y: 3, W: 24, H: 8}

	m := New(theme.Default)
	for _, e := range m.entries {
		for _, st := range e.States {
			if st.Overlay {
				continue // positions itself against the canvas, by contract
			}
			t.Run(e.Name+"-"+strings.ReplaceAll(st.Name, " ", "-"), func(t *testing.T) {
				c := comp.NewCanvas(48, 16)
				st.Draw(c, box, true)

				var out []string
				for y := 0; y < 16; y++ {
					for x := 0; x < 48; x++ {
						if x >= box.X && x <= box.Right() && y >= box.Y && y <= box.Bottom() {
							continue
						}
						if owner := c.OwnerAt(x, y); owner.Name != "" {
							out = append(out, fmt.Sprintf("(%d,%d) owned by %v", x, y, owner))
						}
					}
				}
				if len(out) > 0 {
					shown := out
					if len(shown) > 5 {
						shown = shown[:5]
					}
					t.Errorf("drew %d cells outside %v — it is painting over whatever is beside it.\n  %s",
						len(out), box, strings.Join(shown, "\n  "))
				}
			})
		}
	}
}
