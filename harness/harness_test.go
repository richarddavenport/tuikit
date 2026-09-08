package harness

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	tea "github.com/charmbracelet/bubbletea"
	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/richarddavenport/tuikit/app"
	"github.com/richarddavenport/tuikit/comp"
)

// The trap that put Strip in a library with a test rather than a regex at each
// call site.
//
// An escape sequence ends at a final byte in 0x40-0x7E, and the SGR terminator
// is `m`. A stripper that ends a sequence at any of [A-Za-z] therefore eats the
// `m` and leaves the parameters behind — stripping exactly the color it was
// meant to preserve, and printing the digits.
func TestStripDoesNotEatTheSGRTerminator(t *testing.T) {
	colored := "\x1b[38;5;205mapi_gateway\x1b[0m  3/3"

	got := Strip(colored)
	if got != "api_gateway  3/3" {
		t.Errorf("Strip = %q", got)
	}
	for _, leftover := range []string{"38", "205", "[", "\x1b"} {
		if strings.Contains(got, leftover) {
			t.Errorf("%q survived stripping: %q", leftover, got)
		}
	}
}

func TestStripHandlesTheFormsAFrameContains(t *testing.T) {
	for _, tc := range []struct{ in, want string }{
		{"plain", "plain"},
		{"\x1b[0m", ""},
		{"\x1b[1;38;5;205mbold\x1b[0m", "bold"},
		{"\x1b[48;5;57m row \x1b[0m", " row "},
		{"a\x1b[38;5;34mb\x1b[0mc", "abc"},
		{"┌─┐\x1b[0m", "┌─┐"},
		{"\x1b[38;5;241m↑↓ move · q quit\x1b[0m", "↑↓ move · q quit"},
	} {
		if got := Strip(tc.in); got != tc.want {
			t.Errorf("Strip(%q) = %q, want %q", tc.in, got, tc.want)
		}
	}
}

// Stripping must not change the shape, or a golden is a golden of a different
// frame.
func TestStripPreservesWidth(t *testing.T) {
	style := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	line := style.Render("api_gateway") + "   " + style.Render("3/3")

	if want, got := lipgloss.Width(line), len([]rune(Strip(line))); want != got {
		t.Errorf("visible width %d became %d characters after stripping", want, got)
	}
}

func TestHTMLCarriesTheColorThrough(t *testing.T) {
	forceColor()
	frame := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("selected")

	got := HTML(frame)
	if !strings.Contains(got, "#ff5faf") {
		t.Errorf("palette index 205 did not become its hex:\n%s", got)
	}
	if !strings.Contains(got, "selected") {
		t.Errorf("the text is missing:\n%s", got)
	}
	if strings.Contains(got, "\x1b") {
		t.Errorf("an escape sequence reached the page:\n%s", got)
	}
}

// A color left switched on across a newline paints the page's background,
// which is how a frame ends up with a colored margin down the side.
func TestHTMLClosesSpansAtEveryLine(t *testing.T) {
	frame := "\x1b[38;5;205mone\ntwo\nthree"

	got := HTML(frame)
	if opens, closes := strings.Count(got, "<span"), strings.Count(got, "</span>"); opens != closes {
		t.Errorf("%d spans opened, %d closed:\n%s", opens, closes, got)
	}
}

func TestHTMLEscapesWhatTheTerminalDraws(t *testing.T) {
	got := HTML(`<script>&"`)
	for _, want := range []string{"&lt;script&gt;", "&amp;", "&#34;"} {
		if !strings.Contains(got, want) {
			t.Errorf("%s is not escaped: %s", want, got)
		}
	}
}

// --- capture ------------------------------------------------------------

type fakeModel struct {
	w, h int
	at   time.Time
	body string
}

func (m *fakeModel) SetSize(w, h int) { m.w, m.h = w, h }
func (m *fakeModel) Now(t time.Time)  { m.at = t }
func (m *fakeModel) View() string {
	return fmt.Sprintf("%s at %dx%d, %s", m.body, m.w, m.h, m.at.Format("15:04:05"))
}

func TestCaptureFixesTheSizeAndFreezesTheClock(t *testing.T) {
	dir := t.TempDir()
	at := time.Date(2026, 8, 31, 9, 14, 3, 0, time.UTC)

	s := Capture(t, dir, Size(100, 24), At(at))
	m := &fakeModel{body: "frame"}
	s.Shot("one", m)

	if m.w != 100 || m.h != 24 {
		t.Errorf("the model was not sized: %dx%d", m.w, m.h)
	}
	if !m.at.Equal(at) {
		t.Errorf("the clock was not frozen: %v", m.at)
	}
	body, err := os.ReadFile(filepath.Join(dir, "one.ansi"))
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(string(body), "09:14:03") {
		t.Errorf("the frozen time is not in the frame: %s", body)
	}
}

// A model that does not care about size or time is still capturable — the
// harness asks, it does not demand.
type bareModel struct{}

func (bareModel) View() string { return "bare" }

func TestCaptureDoesNotRequireASizerOrAClock(t *testing.T) {
	s := Capture(t, t.TempDir())
	s.Shot("bare", bareModel{})

	if got := s.Frames(); len(got) != 1 || got[0].Name != "bare" {
		t.Errorf("got %+v", got)
	}
}

// The manifest is what docgen reads and what an agent reads to know which
// frames exist without listing a directory and guessing at names.
func TestDoneWritesAManifestCarryingProvenance(t *testing.T) {
	dir := t.TempDir()
	s := Capture(t, dir, Size(80, 20), WithMode(Live))
	s.Shot("live-one", bareModel{})
	s.ShotAs("staged", bareModel{}, Composed)
	s.Done()

	body, err := os.ReadFile(filepath.Join(dir, "frames.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Width  int     `json:"width"`
		Frames []Frame `json:"frames"`
	}
	if err := json.Unmarshal(body, &manifest); err != nil {
		t.Fatal(err)
	}
	if manifest.Width != 80 {
		t.Errorf("width = %d", manifest.Width)
	}
	modes := map[string]Mode{}
	for _, f := range manifest.Frames {
		modes[f.Name] = f.Mode
	}
	// A page that quietly mixes real and staged frames misreports the tool.
	if modes["live-one"] != Live || modes["staged"] != Composed {
		t.Errorf("provenance is wrong: %v", modes)
	}
}

// Capture into nowhere fails rather than quietly doing nothing — the same rule
// the guards follow.
func TestCaptureWithNoDirectoryFails(t *testing.T) {
	rec := &recorder{}
	func() {
		defer func() { _ = recover() }()
		Capture(rec, "")
	}()
	if len(rec.msgs) == 0 || !strings.Contains(rec.msgs[0], "no directory") {
		t.Errorf("got %v", rec.msgs)
	}
}

func TestEnabledReadsTheEnvironment(t *testing.T) {
	t.Setenv("TUIKIT_TEST_FRAMES", "/tmp/frames")
	if got := Enabled("TUIKIT_TEST_FRAMES"); got != "/tmp/frames" {
		t.Errorf("Enabled = %q", got)
	}
	if got := Enabled("TUIKIT_TEST_UNSET"); got != "" {
		t.Errorf("an unset variable returned %q", got)
	}
}

// --- goldens ------------------------------------------------------------

func TestGoldenComparesShapeAndReportsTheLine(t *testing.T) {
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "frame.golden"), []byte("one\ntwo\nthree"), 0o600); err != nil {
		t.Fatal(err)
	}

	rec := &recorder{}
	Golden(rec, dir, "frame", "one\nTWO\nthree")

	if len(rec.msgs) != 1 {
		t.Fatalf("want one failure, got %v", rec.msgs)
	}
	for _, part := range []string{"line 2", "was", "now", "cols"} {
		if !strings.Contains(rec.msgs[0], part) {
			t.Errorf("the report is missing %q:\n%s", part, rec.msgs[0])
		}
	}
}

// A golden holds the shape, so color is not a difference.
func TestGoldenIgnoresColor(t *testing.T) {
	forceColor()
	dir := t.TempDir()
	if err := os.WriteFile(filepath.Join(dir, "frame.golden"), []byte("selected"), 0o600); err != nil {
		t.Fatal(err)
	}

	colored := lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("selected")
	rec := &recorder{}
	Golden(rec, dir, "frame", colored)

	if len(rec.msgs) != 0 {
		t.Errorf("color was reported as a difference: %v", rec.msgs)
	}
}

// A missing golden says how to write one rather than writing it: a golden that
// appears on first run is a test that has never checked anything.
func TestAMissingGoldenSaysHowToWriteIt(t *testing.T) {
	rec := &recorder{}
	Golden(rec, t.TempDir(), "absent", "frame")

	if len(rec.msgs) != 1 || !strings.Contains(rec.msgs[0], "-update-goldens") {
		t.Errorf("got %v", rec.msgs)
	}
}

// --- helpers ------------------------------------------------------------

func forceColor() {
	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)
}

type recorder struct{ msgs []string }

func (r *recorder) Helper() {}
func (r *recorder) Errorf(format string, args ...any) {
	r.msgs = append(r.msgs, fmt.Sprintf(format, args...))
}
func (r *recorder) Fatalf(format string, args ...any) {
	r.Errorf(format, args...)
	panic("fatal")
}

// A terminal accumulates SGR state, so bold set by one sequence survives a
// later one that only names a color. Treating each sequence as a complete
// style renders the second run pink with the bold silently gone — which is a
// real frame democtl draws, not a synthetic case.
func TestHTMLAccumulatesStyleTheWayATerminalDoes(t *testing.T) {
	forceColor()
	got := HTML("\x1b[1m\x1b[38;5;205mbold pink\x1b[0m")

	// The text must sit inside ONE span carrying both. Asserting that the page
	// merely contains each property somewhere passes against the bug, which
	// emits a bold span and then a separate color span with the bold gone.
	if want := `<span style="color:#ff5faf;font-weight:600">bold pink`; !strings.Contains(got, want) {
		t.Errorf("the bold and the color are not on one span:\n%s", got)
	}
}

// A style that straddles a line break continues in a terminal, so it has to
// continue on the page. The span still closes at the newline — it is reopened
// on the next line rather than left hanging.
func TestHTMLCarriesStyleAcrossALineBreak(t *testing.T) {
	forceColor()
	got := HTML("\x1b[48;5;57mrow one\nrow two\x1b[0m")

	lines := strings.Split(got, "\n")
	if len(lines) != 2 {
		t.Fatalf("expected two lines:\n%s", got)
	}
	if !strings.Contains(lines[1], "background:#5f00ff") {
		t.Errorf("the background did not cross the line break:\n%s", lines[1])
	}
	if opens, closes := strings.Count(got, "<span"), strings.Count(got, "</span>"); opens != closes {
		t.Errorf("%d spans opened, %d closed:\n%s", opens, closes, got)
	}
}

// 24-bit color is not something theme produces — the palette is nine 256
// indices — but it is what lipgloss emits the moment a tool author writes a hex
// value, and reading its channels one at a time as separate SGR codes finds 95
// in the bright-color range and renders a hand-picked pink as bright magenta.
func TestHTMLDoesNotMistakeAColorChannelForACode(t *testing.T) {
	forceColor()
	got := HTML(lipgloss.NewStyle().Foreground(lipgloss.Color("#ff5faf")).Render("hand-picked"))

	if !strings.Contains(got, "#ff5faf") {
		t.Errorf("a 24-bit color did not survive:\n%s", got)
	}
	if strings.Contains(got, "#ff00ff") {
		t.Errorf("the blue channel was read as SGR 95, bright magenta:\n%s", got)
	}
}

// A background must not survive a reset, or every frame after the first styled
// run is painted.
func TestHTMLResetClearsEverything(t *testing.T) {
	forceColor()
	got := HTML("\x1b[1;48;5;57mrow\x1b[0mplain")

	if strings.Contains(got, ">plain") && strings.Contains(got, `<span style="background:#5f00ff;font-weight:600">plain`) {
		t.Errorf("the reset did not clear the state:\n%s", got)
	}
}

// A script can drive hover and double-click by name, like everything else.
type probe struct {
	canvas                    *comp.Canvas
	mouse                     app.Mouse
	hovered, doubled, pressed int
}

func (p *probe) View() string {
	c := comp.NewCanvas(20, 3)
	c.Fill(comp.Rect{X: 2, Y: 1, W: 8, H: 1}, " ", nil, comp.Region("box"))
	p.canvas = c
	return c.String()
}
func (p *probe) Canvas() *comp.Canvas { return p.canvas }
func (p *probe) Init() tea.Cmd        { return nil }
func (p *probe) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	if mouse, ok := msg.(tea.MouseMsg); ok {
		p.mouse.Route(mouse, p.canvas, app.Handler{
			Hover:       func(comp.ID) tea.Cmd { p.hovered++; return nil },
			Press:       func(comp.ID, tea.MouseMsg) tea.Cmd { p.pressed++; return nil },
			DoubleClick: func(comp.ID, tea.MouseMsg) tea.Cmd { p.doubled++; return nil },
		})
	}
	return p, nil
}

func TestHoverAndDoubleClickAreScriptable(t *testing.T) {
	p := &probe{}
	p.View() // a region has to be drawn before it can be addressed
	Script(t, p, `
	    hover box
	    doubleclick box
	`)
	if p.hovered != 1 {
		t.Errorf("hover fired %d times, want 1", p.hovered)
	}
	if p.doubled != 1 {
		t.Errorf("double-click fired %d times, want 1", p.doubled)
	}
	// The first press of a double is a real press.
	if p.pressed != 2 {
		t.Errorf("press fired %d times, want 2", p.pressed)
	}
}
