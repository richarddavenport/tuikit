package comp

import (
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

func forceColour() {
	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)
}

const (
	list Name = "services.list"
	row  Name = "services.row"
	pane Name = "detail"
)

// A new canvas is blanks, and it serialises to exactly the size it was asked
// for. Everything else depends on this.
func TestACanvasIsExactlyTheSizeItWasAskedFor(t *testing.T) {
	c := NewCanvas(10, 3)
	got := c.String()

	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("expected 3 lines, got %d: %q", len(lines), got)
	}
	for i, line := range lines {
		if Width(line) != 10 {
			t.Errorf("line %d is %d columns, not 10: %q", i+1, Width(line), line)
		}
	}
}

// Drawing past the edge is not an error to catch. It is a coordinate that does
// not exist, which is what makes "nothing wider than the terminal" structural
// rather than something every tool has to remember to check.
func TestDrawingOutsideTheCanvasIsANoOp(t *testing.T) {
	c := NewCanvas(4, 2)
	c.Text(2, 0, "far too long", nil, Region(list))
	c.Text(-3, 1, "before the start", nil, Region(list))
	c.Set(0, 9, "x", nil, Region(list))

	for i, line := range strings.Split(c.String(), "\n") {
		if Width(line) != 4 {
			t.Errorf("line %d is %d columns, not 4: %q", i+1, Width(line), line)
		}
	}
}

// A clipped string is a PREFIX of what was asked for. Skipping a cluster that
// does not fit and carrying on would leave a hole in the middle, which reads as
// a data bug rather than as a narrow terminal.
func TestTextClipsToAPrefix(t *testing.T) {
	c := NewCanvas(6, 1)
	if n := c.Text(0, 0, "abcdefghij", nil, Region(list)); n != 6 {
		t.Errorf("drew %d columns into a 6-column canvas", n)
	}
	if got := c.String(); got != "abcdef" {
		t.Errorf("got %q", got)
	}
}

// --- wide runes ---------------------------------------------------------

// One rune is not one cell. A CJK character claims two, and the second is a
// continuation rather than a character of its own.
func TestAWideRuneClaimsTwoCells(t *testing.T) {
	c := NewCanvas(6, 1)
	if n := c.Set(0, 0, "世", nil, Region(pane)); n != 2 {
		t.Errorf("a wide rune claimed %d columns", n)
	}

	lead, _ := c.CellAt(0, 0)
	cont, _ := c.CellAt(1, 0)
	if lead.Text != "世" {
		t.Errorf("the lead cell holds %q", lead.Text)
	}
	if !cont.Continuation() {
		t.Errorf("the second cell is %q, not a continuation", cont.Text)
	}
}

// A click on the right half of a wide glyph has to hit the thing that drew it.
// The alternative is a character you can only select by aiming at its left half.
func TestOwnerAtIsTheSameOnBothHalvesOfAWideRune(t *testing.T) {
	c := NewCanvas(6, 1)
	c.Set(2, 0, "界", nil, Region(row).At(4))

	if left, right := c.OwnerAt(2, 0), c.OwnerAt(3, 0); left != right {
		t.Errorf("the halves have different owners: %v and %v", left, right)
	}
	if got := c.OwnerAt(3, 0).String(); got != "services.row[4]" {
		t.Errorf("the right half reports %q", got)
	}
}

// Serialising emits the lead once and skips its continuation. Emitting a space
// there would push everything after it a column right — the exact class of bug
// the canvas exists to make impossible.
func TestAFrameOfWideRunesSerialisesToItsDeclaredWidth(t *testing.T) {
	for _, tc := range []struct{ name, text string }{
		{"CJK", "世界世界"},
		{"an emoji", "ok 👍 done"},
		{"an emoji built from joiners", "hi 👨‍👩‍👧 there"},
		{"mixed", "世 ok 👍"},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCanvas(20, 1)
			c.Text(0, 0, tc.text, nil, Region(pane))
			if got := Width(c.String()); got != 20 {
				t.Errorf("a 20-column canvas serialised to %d columns: %q", got, c.String())
			}
		})
	}
}

// Overwriting half of a wide cluster has to take both cells. A leftover
// continuation is an orphan — it serialises as nothing, so the row silently
// loses a column.
func TestOverwritingHalfOfAWideRuneClearsBoth(t *testing.T) {
	for _, tc := range []struct {
		name string
		x    int
	}{
		{"writing over its lead", 0},
		{"writing over its continuation", 1},
	} {
		t.Run(tc.name, func(t *testing.T) {
			c := NewCanvas(4, 1)
			c.Set(0, 0, "世", nil, Region(pane))
			c.Set(tc.x, 0, "x", nil, Region(pane))

			if got := c.String(); Width(got) != 4 {
				t.Errorf("the row is %d columns after the overwrite: %q", Width(got), got)
			}
			if strings.Contains(c.String(), "世") {
				t.Errorf("half of the wide rune survived: %q", c.String())
			}
		})
	}
}

// A wide rune with one column left draws nothing. There is no half of a glyph
// to draw, and a row a column short is worse than a row with a blank in it.
func TestAWideRuneAtTheLastColumnClips(t *testing.T) {
	c := NewCanvas(3, 1)
	c.Text(0, 0, "ab", nil, Region(pane))
	if n := c.Set(2, 0, "世", nil, Region(pane)); n != 0 {
		t.Errorf("a wide rune drew %d columns with one to spare", n)
	}
	if got := c.String(); got != "ab " || Width(got) != 3 {
		t.Errorf("got %q, want %q", got, "ab ")
	}
}

// --- ownership ----------------------------------------------------------

// An owner is an identity, so the index is the item's own, not the row it
// happens to be on. Two rows of the same list scrolled to different offsets
// must still name the same service.
func TestAnOwnerIsAnIdentityNotAScreenRow(t *testing.T) {
	c := NewCanvas(20, 2)
	c.Text(0, 0, "api_gateway", nil, Region(row).At(41))
	c.Text(0, 1, "api_worker", nil, Region(row).At(42))

	if got := c.OwnerAt(3, 0); got.Index != 41 {
		t.Errorf("the first row reports index %d", got.Index)
	}
	if got := c.OwnerAt(3, 1); got.Index != 42 {
		t.Errorf("the second row reports index %d", got.Index)
	}
}

// A region that nothing drew is a region a script may not name.
func TestRegionReportsWhetherItWasDrawn(t *testing.T) {
	c := NewCanvas(20, 3)
	c.Text(4, 1, "api_gateway", nil, Region(row).At(0))

	got, ok := c.Region(Region(row).At(0))
	if !ok {
		t.Fatalf("a region that was drawn reports as missing")
	}
	if want := (Rect{4, 1, 11, 1}); got != want {
		t.Errorf("region is %+v, want %+v", got, want)
	}
	if _, ok := c.Region(Region(row).At(7)); ok {
		t.Errorf("a region nothing drew reports as present")
	}
}

func TestAnUnclaimedCellHasNoOwner(t *testing.T) {
	c := NewCanvas(4, 1)
	if got := c.OwnerAt(0, 0); !got.Zero() {
		t.Errorf("a blank cell is owned by %v", got)
	}
	if got := c.OwnerAt(99, 99); !got.Zero() {
		t.Errorf("a cell off the canvas is owned by %v", got)
	}
}

func TestIDNamesWhatACaptureScriptWrites(t *testing.T) {
	for _, tc := range []struct {
		id   ID
		want string
	}{
		{Region(list), "services.list"},
		{Region(row).At(2), "services.row[2]"},
		{ID{}, ""},
	} {
		if got := tc.id.String(); got != tc.want {
			t.Errorf("%+v is %q, want %q", tc.id, got, tc.want)
		}
	}
}

// --- serialising --------------------------------------------------------

// Runs sharing a style are one escape sequence, not one per character. A frame
// is captured and diffed; one sequence per cell makes both unreadable.
func TestRunsSharingAStyleAreGroupedOnSerialising(t *testing.T) {
	forceColour()
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))

	c := NewCanvas(11, 1)
	c.Text(0, 0, "api_gateway", &accent, Region(row).At(0))

	if got := strings.Count(c.String(), "\x1b[38;5;205m"); got != 1 {
		t.Errorf("eleven characters produced %d colour sequences:\n%q", got, c.String())
	}
}

// Style identity is by pointer, so two components holding the same style share
// a run and two holding their own do not.
func TestAStyleChangeBreaksTheRun(t *testing.T) {
	forceColour()
	accent := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	muted := lipgloss.NewStyle().Foreground(lipgloss.Color("241"))

	c := NewCanvas(6, 1)
	c.Text(0, 0, "ok", &accent, Region(pane))
	c.Text(3, 0, "no", &muted, Region(pane))

	got := c.String()
	if n := strings.Count(got, "\x1b["); n < 2 {
		t.Errorf("two styles produced %d sequences: %q", n, got)
	}
}

// --- rects --------------------------------------------------------------

func TestRect(t *testing.T) {
	r := Rect{2, 1, 10, 4}
	if !r.Contains(2, 1) || !r.Contains(11, 4) {
		t.Error("a corner is not inside the rect")
	}
	if r.Contains(12, 4) || r.Contains(2, 0) {
		t.Error("a point outside the rect is inside it")
	}
	if got, want := r.Inset(1), (Rect{3, 2, 8, 2}); got != want {
		t.Errorf("Inset(1) is %+v, want %+v", got, want)
	}
	if !(Rect{0, 0, 3, 0}).Empty() {
		t.Error("a rect with no rows is not empty")
	}
	// Inset must not produce a negative size — a pane squeezed by a narrow
	// terminal is a normal state, and a negative width crashes whatever loops
	// over it.
	if got := (Rect{0, 0, 1, 1}).Inset(4); got.W != 0 || got.H != 0 {
		t.Errorf("over-insetting produced %+v", got)
	}
}

// The box is drawn with the six characters the default glyph set allows, and it
// owns every cell of its frame so a click on a border hits the pane.
func TestBoxDrawsAFrameItOwns(t *testing.T) {
	c := NewCanvas(5, 3)
	c.Box(Rect{0, 0, 5, 3}, nil, Region(pane))

	want := "┌───┐\n│   │\n└───┘"
	if got := c.String(); got != want {
		t.Errorf("got:\n%s\nwant:\n%s", got, want)
	}
	if got := c.OwnerAt(0, 0); got != Region(pane) {
		t.Errorf("the corner is owned by %v", got)
	}
	if got := c.OwnerAt(2, 1); !got.Zero() {
		t.Errorf("the inside of the box is owned by %v", got)
	}
}

func TestBoxTooSmallToDrawDrawsNothing(t *testing.T) {
	c := NewCanvas(4, 2)
	c.Box(Rect{0, 0, 1, 1}, nil, Region(pane))
	if got := c.String(); strings.ContainsAny(got, "┌┐└┘─│") {
		t.Errorf("a 1x1 box drew something: %q", got)
	}
}

func TestFillCoversARect(t *testing.T) {
	c := NewCanvas(6, 3)
	c.Fill(Rect{1, 1, 3, 1}, "·", nil, Region(pane))

	if got, want := c.String(), "      \n ···  \n      "; got != want {
		t.Errorf("got %q, want %q", got, want)
	}
}
