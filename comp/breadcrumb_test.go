package comp

import (
	"strings"
	"testing"
)

func crumbs(b Breadcrumb, w int) (*Canvas, string) {
	c := NewCanvas(w, 1)
	b.Name = "trail"
	b.Draw(c, Rect{X: 0, Y: 0, W: w, H: 1})
	return c, c.String()
}

func TestBreadcrumbDrawsTheTrail(t *testing.T) {
	_, got := crumbs(Breadcrumb{Crumbs: []string{"democtl", "api_gateway", "Logs"}}, 60)
	if want := "democtl · api_gateway · Logs"; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}

func TestBreadcrumbOfOne(t *testing.T) {
	if _, got := crumbs(Breadcrumb{Crumbs: []string{"democtl"}}, 40); got != "democtl" {
		t.Errorf("got %q", got)
	}
}

// Too narrow, it elides from the LEFT: where you are and how far in are the two
// things a breadcrumb is for, and the middle is what you can lose.
func TestBreadcrumbElidesFromTheLeft(t *testing.T) {
	b := Breadcrumb{Crumbs: []string{"democtl", "services", "api_gateway", "Logs"}}
	_, got := crumbs(b, 26)

	if !strings.HasSuffix(got, "Logs") {
		t.Errorf("the current screen was dropped: %q", got)
	}
	if !strings.HasPrefix(got, "…") {
		t.Errorf("nothing says the trail was cut: %q", got)
	}
	if Width(got) > 26 {
		t.Errorf("%q is %d columns in a 26-column rect", got, Width(got))
	}
}

// The last crumb survives alone. A breadcrumb that dropped the screen you are
// on to keep the ones you are not is answering the wrong question.
func TestTheCurrentScreenAlwaysSurvives(t *testing.T) {
	b := Breadcrumb{Crumbs: []string{"democtl", "services", "api_gateway", "Logs"}}
	_, got := crumbs(b, 8)
	if !strings.Contains(got, "Logs") {
		t.Errorf("got %q, want it to keep Logs", got)
	}
}

// A click resolves to a DEPTH by name, rather than by counting columns.
func TestEachCrumbIsItsOwnRegion(t *testing.T) {
	c, _ := crumbs(Breadcrumb{
		Crumbs: []string{"democtl", "api_gateway", "Logs"}, Item: "crumb",
	}, 60)

	for i, want := range []string{"democtl", "api_gateway", "Logs"} {
		r, ok := c.Region(Region("crumb").At(i))
		if !ok {
			t.Errorf("crumb %d (%s) is not a region", i, want)
			continue
		}
		if got := Width(want); r.W != got {
			t.Errorf("crumb %d is %d columns, want %d", i, r.W, got)
		}
	}
	// The separators belong to the strip, not to either crumb beside them.
	if owner := c.OwnerAt(7, 0); owner.Name != "trail" {
		t.Errorf("the separator is owned by %v", owner)
	}
}

// Without Item the whole trail is still one region, so a click has somewhere
// to land.
func TestAnUnnamedTrailIsStillOneRegion(t *testing.T) {
	c, _ := crumbs(Breadcrumb{Crumbs: []string{"a", "b"}}, 20)
	if _, ok := c.Region(Region("trail")); !ok {
		t.Error("the strip is not a region")
	}
}

func TestBreadcrumbInNoRoom(t *testing.T) {
	for _, w := range []int{0, 1} {
		if _, got := crumbs(Breadcrumb{Crumbs: []string{"democtl", "Logs"}}, w); strings.Contains(got, "democtl") {
			t.Errorf("width %d drew the root: %q", w, got)
		}
	}
	c := NewCanvas(20, 1)
	if y := (Breadcrumb{Name: "trail"}).Draw(c, Rect{X: 0, Y: 0, W: 20, H: 1}); y != 0 {
		t.Errorf("an empty trail advanced to row %d", y)
	}
}
