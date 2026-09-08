package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
)

const detail Name = "detail"

func TestAPaneIsAFrameAndTheRoomInside(t *testing.T) {
	c := NewCanvas(12, 4)
	inner := Pane{}.Draw(c, c.Bounds(), Region(detail))

	if want := (Rect{1, 1, 10, 2}); inner != want {
		t.Errorf("inside is %+v, want %+v", inner, want)
	}
	if got, want := c.String(), "┌──────────┐\n│          │\n│          │\n└──────────┘"; got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
}

func TestATitleInTheEdgeCostsNoRow(t *testing.T) {
	c := NewCanvas(16, 4)
	inner := Pane{Title: "Detail"}.Draw(c, c.Bounds(), Region(detail))

	if inner.Y != 1 || inner.H != 2 {
		t.Errorf("a title in the edge took a row: %+v", inner)
	}
	if got := strings.Split(c.String(), "\n")[0]; got != "┌─ Detail ─────┐" {
		t.Errorf("top edge is %q", got)
	}
}

func TestATitleOnARowIsPaintedFullWidth(t *testing.T) {
	forceColor()
	title := lipgloss.NewStyle().Background(lipgloss.Color("57"))

	c := NewCanvas(16, 5)
	inner := Pane{Title: "Detail", TitleAt: TitleOnRow, TitleStyle: &title}.
		Draw(c, c.Bounds(), Region(detail))

	if inner.Y != 2 {
		t.Errorf("the title did not take a row: %+v", inner)
	}
	// Painted across, so the highlight reads as the pane's rather than the
	// text's.
	if row := strings.Split(c.String(), "\n")[1]; lipgloss.Width(row) != 16 {
		t.Errorf("the title row is %d columns: %q", lipgloss.Width(row), row)
	}
}

// A title too long is truncated rather than dropped. Losing the name of the
// pane you are looking at is worst exactly when the window is too small to work
// out what it is from anything else.
func TestALongTitleIsTruncatedNotDropped(t *testing.T) {
	for _, at := range []TitlePlacement{TitleInEdge, TitleOnRow} {
		c := NewCanvas(14, 4)
		Pane{Title: "a service with a very long name", TitleAt: at}.
			Draw(c, c.Bounds(), Region(detail))

		got := c.String()
		if !strings.Contains(got, "a serv") {
			t.Errorf("placement %d dropped the title:\n%s", at, got)
		}
		if !strings.Contains(got, "…") {
			t.Errorf("placement %d did not say the title was cut:\n%s", at, got)
		}
		for _, line := range strings.Split(got, "\n") {
			if Width(line) > 14 {
				t.Errorf("placement %d drew %d columns: %q", at, Width(line), line)
			}
		}
	}
}

// The inside is blanked, which is what makes a pane drawn over something else
// cover it — a modal is a draw rather than a composite.
func TestAPaneCoversWhatIsBehindIt(t *testing.T) {
	c := NewCanvas(12, 4)
	c.Fill(c.Bounds(), "x", nil, Region(detail))

	Pane{}.Draw(c, Rect{X: 2, Y: 1, W: 8, H: 3}, Region("modal"))

	if row := strings.Split(c.String(), "\n")[2]; !strings.Contains(row, "│      │") {
		t.Errorf("the pane did not cover what was behind it: %q", row)
	}
}

func TestAPaneWithNoRoomDrawsNothing(t *testing.T) {
	c := NewCanvas(10, 4)
	if got := (Pane{Title: "x"}).Draw(c, Rect{X: 0, Y: 0, W: 1, H: 1}, Region(detail)); !got.Empty() {
		t.Errorf("a 1x1 pane returned %+v to draw into", got)
	}
	if strings.ContainsAny(c.String(), "┌┐└┘") {
		t.Errorf("a 1x1 pane drew a frame: %q", c.String())
	}
}

func TestFocusChangesTheBorderAndTheTitle(t *testing.T) {
	forceColor()
	border := lipgloss.NewStyle().Foreground(lipgloss.Color("240"))
	focus := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	p := Pane{Title: "Detail", Border: &border, Focus: &focus}
	plain := NewCanvas(16, 4)
	p.Draw(plain, plain.Bounds(), Region(detail))

	p.Focused = true
	lit := NewCanvas(16, 4)
	p.Draw(lit, lit.Bounds(), Region(detail))

	if !strings.Contains(plain.String(), "38;5;240") {
		t.Error("an unfocused pane does not use the border style")
	}
	if !strings.Contains(lit.String(), "38;5;205") {
		t.Error("a focused pane does not use the focus style")
	}
}
