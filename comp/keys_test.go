package comp

import (
	"strings"
	"testing"
)

func sections() []KeySection {
	return []KeySection{
		{Name: "Everywhere", Keys: []Hint{
			{Key: "?", Label: "these keys"},
			{Key: "q", Label: "quit"},
		}},
		{Name: "The estate", Keys: []Hint{
			{Key: "j/k", Label: "move"},
			{Key: "enter", Label: "expand"},
			{Key: "/", Label: "filter"},
		}},
	}
}

func help(k Keys, w, h int) string {
	c := NewCanvas(w, h)
	k.Draw(c, Rect{X: 0, Y: 0, W: w, H: h}, Region("keys"))
	return c.String()
}

func TestKeysGroupsBySection(t *testing.T) {
	got := help(Keys{Sections: sections()}, 40, 12)
	for _, want := range []string{"Everywhere", "The estate", "?", "enter", "filter"} {
		if !strings.Contains(got, want) {
			t.Errorf("the help is missing %q:\n%s", want, got)
		}
	}
	// A blank line between sections, so it reads as two groups.
	lines := strings.Split(got, "\n")
	var everywhere, estate int
	for i, l := range lines {
		if strings.Contains(l, "Everywhere") {
			everywhere = i
		}
		if strings.Contains(l, "The estate") {
			estate = i
		}
	}
	if estate-everywhere != 4 { // heading, two keys, blank, heading
		t.Errorf("the sections are %d rows apart, want 4:\n%s", estate-everywhere, got)
	}
}

// The labels line up in one column, whatever the keys are.
func TestKeysAlignTheirLabels(t *testing.T) {
	got := help(Keys{Sections: sections()}, 40, 12)
	col := -1
	for _, line := range strings.Split(got, "\n") {
		for _, label := range []string{"these keys", "quit", "move", "expand", "filter"} {
			if i := strings.Index(line, label); i >= 0 {
				if col < 0 {
					col = i
				} else if i != col {
					t.Errorf("%q starts at %d, want %d:\n%s", label, i, col, got)
				}
			}
		}
	}
	if col < 0 {
		t.Fatal("no labels were drawn")
	}
}

// A help screen that quietly omits half the keys is worse than one that says it
// was cut, because a reader believes it.
func TestKeysSaysWhenItRanOutOfRoom(t *testing.T) {
	got := help(Keys{Sections: sections()}, 40, 4)
	if !strings.Contains(got, "more") {
		t.Errorf("the help was cut without saying so:\n%s", got)
	}
}

// As an overlay it is a centered box, so "what was I looking at" is still
// answerable.
func TestKeysAsAnOverlayIsABox(t *testing.T) {
	c := NewCanvas(60, 20)
	c.Fill(Rect{X: 0, Y: 0, W: 60, H: 20}, "x", nil, Region("behind"))
	Keys{Sections: sections(), Overlay: true, Title: "Keys"}.
		Draw(c, c.Bounds(), Region("keys"))

	got := c.String()
	if !strings.Contains(got, "┌") {
		t.Errorf("the overlay has no border:\n%s", got)
	}
	if !strings.Contains(got, "x") {
		t.Error("the overlay covered the whole interface")
	}
	if !strings.Contains(got, "Everywhere") {
		t.Error("the overlay drew no keys")
	}
}

func TestKeysWithNothingToSay(t *testing.T) {
	if got := help(Keys{}, 20, 5); strings.TrimSpace(got) != "" {
		t.Errorf("an empty help drew %q", got)
	}
	if got := help(Keys{Sections: sections()}, 20, 0); strings.TrimSpace(got) != "" {
		t.Errorf("no room drew %q", got)
	}
}

// A section with no name is still a group — the blank line does the work.
func TestAnUnnamedSection(t *testing.T) {
	got := help(Keys{Sections: []KeySection{
		{Keys: []Hint{{Key: "q", Label: "quit"}}},
	}}, 30, 5)
	if !strings.Contains(got, "quit") {
		t.Errorf("an unnamed section drew nothing:\n%s", got)
	}
}

// The marker CLEARS the row it replaces. One that only overwrote its first few
// columns read as "… moretate".
func TestTheMoreMarkerClearsItsRow(t *testing.T) {
	got := help(Keys{Sections: sections()}, 40, 4)
	for _, line := range strings.Split(got, "\n") {
		if !strings.Contains(line, "more") {
			continue
		}
		// Whatever the marker says, nothing of the row it covered may survive
		// beside it.
		if fields := strings.Fields(line); len(fields) > 3 {
			t.Errorf("the marker row is %q — something of the row it replaced is still there", line)
		}
	}
}

// The marker says HOW MANY, and which way. "there is more" and "there are nine
// more" are different amounts of help, and a reader deciding whether to scroll
// is answering the second question (issue 50).
func TestTheMarkerCountsWhatIsHidden(t *testing.T) {
	k := Keys{Sections: sections()}
	got := help(k, 40, 4)

	if !strings.Contains(got, "more") {
		t.Fatalf("nothing was marked as hidden:\n%s", got)
	}
	// Four rows, one of which the marker takes, so three lines are shown.
	hidden := k.Rows() - 3
	if !strings.Contains(got, itoa(hidden)+" more") {
		t.Errorf("the marker does not say %d are hidden:\n%s", hidden, got)
	}
}

// Offset is how a reader reaches them, which is the half that was missing: it
// said "… more" and offered no way to see it.
func TestOffsetScrollsTheKeys(t *testing.T) {
	all := sections()
	last := all[len(all)-1].Keys[len(all[len(all)-1].Keys)-1]

	top := help(Keys{Sections: all}, 40, 5)
	if strings.Contains(top, last.Label) {
		t.Fatalf("the fixture is too short to test scrolling:\n%s", top)
	}

	k := Keys{Sections: all}
	scrolled := help(Keys{Sections: all, Offset: k.Rows()}, 40, 5)
	if !strings.Contains(scrolled, last.Label) {
		t.Errorf("scrolling to the end does not reach the last binding:\n%s", scrolled)
	}
	if !strings.Contains(scrolled, "above") {
		t.Errorf("a scrolled list does not say there is something above it:\n%s", scrolled)
	}
}

// An offset past the end still draws, rather than showing an empty box: a
// caller clamping against Rows is the contract, and a caller that does not is
// not punished with a blank help screen.
func TestAnOffsetPastTheEndStillDrawsSomething(t *testing.T) {
	got := help(Keys{Sections: sections(), Offset: 999}, 40, 5)
	if strings.TrimSpace(got) == "" {
		t.Error("an over-scrolled help screen is blank")
	}
}
