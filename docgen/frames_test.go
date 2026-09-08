package docgen

import (
	"os"
	"path/filepath"
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/richarddavenport/tuikit/harness"
)

// capture writes a small real capture — real ANSI from lipgloss, not a fixture
// string, because the escape sequences are the thing being handled.
func capture(t *testing.T) string {
	t.Helper()
	lipgloss.SetColorProfile(termenv.TrueColor)

	dir := t.TempDir()
	s := harness.Capture(t, dir, harness.Size(40, 3), harness.WithMode(harness.Live))
	s.Shot("first", model{"┌──────┐\n│ " + lipgloss.NewStyle().Foreground(lipgloss.Color("205")).Render("api") + "  │\n└──────┘"})
	s.ShotAs("staged", model{lipgloss.NewStyle().Foreground(lipgloss.Color("34")).Render("ok")}, harness.Composed)
	s.Shot("third", model{"plain"})
	s.Done()
	return dir
}

type model struct{ body string }

func (m model) View() string { return m.body }

func TestPageRendersEveryFrameWithItsProvenance(t *testing.T) {
	dir := capture(t)

	page, err := Frames{
		Title: "democtl",
		Lede:  "What the tool looks like.",
		Groups: []Group{{
			Title:  "Browsing",
			Lede:   "Where you land.",
			Frames: []FrameNote{{Name: "first", Keys: "tab right", Note: "the dashboard"}},
		}},
	}.Page(dir)
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{
		"<title>democtl</title>", "What the tool looks like.",
		"Browsing", "Where you land.", "the dashboard",
		"<kbd>tab</kbd>", "<kbd>right</kbd>",
		`class="chip live"`, `class="chip composed"`,
		"#ff5faf", // the accent survived as its hex
	} {
		if !strings.Contains(page, want) {
			t.Errorf("the page is missing %q", want)
		}
	}
	if n := strings.Count(page, `class="tuikit-frame"`); n != 3 {
		t.Errorf("%d frames rendered, want 3", n)
	}
}

// A frame no group names must still appear. A page that silently drops a new
// screen is worse than one with an untidy section — the omission is invisible,
// which is the failure mode of every hand-maintained list.
func TestAFrameNoGroupNamesStillAppears(t *testing.T) {
	dir := capture(t)

	page, err := Frames{Groups: []Group{{
		Title:  "Browsing",
		Frames: []FrameNote{{Name: "first"}},
	}}}.Page(dir)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.Contains(page, "Ungrouped") {
		t.Error("no section for the frames the groups did not name")
	}
	for _, name := range []string{"first", "staged", "third"} {
		if !strings.Contains(page, ">"+name+"<") {
			t.Errorf("%s is not on the page", name)
		}
	}
}

// The zero value works: point it at a capture and get a page.
func TestTheZeroValueRendersEveryFrame(t *testing.T) {
	page, err := Frames{}.Page(capture(t))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(page, "Ungrouped") {
		t.Error("with no groups at all the section should read Frames, not Ungrouped")
	}
	if n := strings.Count(page, `class="tuikit-frame"`); n != 3 {
		t.Errorf("%d frames, want 3", n)
	}
}

// A directory with no capture in it is an error, not an empty page. A page
// rendered from nothing says the tool has no screens.
func TestAnEmptyDirectoryIsAnError(t *testing.T) {
	if _, err := (Frames{}).Page(t.TempDir()); err == nil {
		t.Fatal("no error for a directory with no capture")
	} else if !strings.Contains(err.Error(), "capture there first") {
		t.Errorf("the error does not say what to do: %v", err)
	}
}

// Over a real multi-frame capture rather than a fixture string: no escape
// sequence may reach the page, and every span must close.
func TestNoEscapeSequenceReachesThePageAndSpansBalance(t *testing.T) {
	page, err := Frames{}.Page(capture(t))
	if err != nil {
		t.Fatal(err)
	}
	if strings.Contains(page, "\x1b") {
		t.Error("an escape sequence reached the page")
	}
	if o, c := strings.Count(page, "<span"), strings.Count(page, "</span>"); o != c {
		t.Errorf("%d spans opened, %d closed", o, c)
	}
}

// The three rules from design/capture.md, enforced rather than remembered.
func TestTheFrameRulesAreInTheGeneratedCSS(t *testing.T) {
	page, err := Frames{}.Page(capture(t))
	if err != nil {
		t.Fatal(err)
	}
	frame := framePlateCSS(t, page)

	if !strings.Contains(frame, "background: #0c0c0e") {
		t.Error("the frame does not keep a dark ground")
	}
	if strings.Contains(frame, "var(--ground)") || strings.Contains(frame, "var(--surface)") {
		t.Error("the frame takes its ground from a theme token, so it would flip in light mode")
	}
	if !strings.Contains(frame, "ui-monospace") || strings.Contains(page, "fonts.googleapis.com") {
		t.Error("the frames must use the system mono stack, with no webfont")
	}
	if !strings.Contains(frame, `font-variant-ligatures: none`) {
		t.Error("ligatures would break one character to one cell")
	}
}

// The classic unreadable-artifact bug: a colour whose only definition sits
// inside a media query never applies when the viewer's theme is unset, which is
// the default.
func TestEveryColourIsDefinedOnBareRoot(t *testing.T) {
	page, err := Frames{}.Page(capture(t))
	if err != nil {
		t.Fatal(err)
	}
	bare := between(t, page, ":root {", "}")

	for _, token := range regexp.MustCompile(`--[a-z]+`).FindAllString(page, -1) {
		if strings.HasPrefix(token, "--sans") || strings.HasPrefix(token, "--mono") {
			continue
		}
		if !strings.Contains(bare, token+":") {
			t.Errorf("%s is never defined on bare :root", token)
		}
	}
}

func TestWriteProducesAFile(t *testing.T) {
	out := filepath.Join(t.TempDir(), "page.html")
	if err := (Frames{Title: "democtl"}).Write(capture(t), out); err != nil {
		t.Fatal(err)
	}
	body, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	if !strings.HasPrefix(string(body), "<!doctype html>") {
		t.Error("the file is not a document")
	}
}

func framePlateCSS(t *testing.T, page string) string {
	t.Helper()
	return between(t, page, ".tuikit-frame {", "}")
}

func between(t *testing.T, s, open, close string) string {
	t.Helper()
	i := strings.Index(s, open)
	if i < 0 {
		t.Fatalf("%q not found", open)
	}
	rest := s[i+len(open):]
	j := strings.Index(rest, close)
	if j < 0 {
		t.Fatalf("%q not closed", open)
	}
	return rest[:j]
}

// The page says its colours are not the reader's.
//
// Issue 52: it renders the sixteen as xterm's defaults, which is the only
// answer a static page can give — and until it said so, it documented a
// themeable tool in a palette nobody sees. A reviewer read "accent is bright
// magenta" off a page like this; their own index 13 was dark mustard.
func TestThePageSaysItsColoursAreNotYours(t *testing.T) {
	page, err := Frames{Title: "mytool"}.Page(capture(t))
	if err != nil {
		t.Fatal(err)
	}

	for _, want := range []string{"xterm", "terminal"} {
		if !strings.Contains(page, want) {
			t.Errorf("the page does not mention %q, so a reader reviewing colour is reviewing a tool nobody sees", want)
		}
	}
	if !strings.Contains(page, "caveat") {
		t.Error("the caveat is not marked up, so it cannot be styled apart from the lede")
	}
}
