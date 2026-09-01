package docgen

import (
	"regexp"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/theme"
)

// The point of generating the design system rather than drawing it: it cannot
// depict a tool that does not exist.
//
// A mockup made in a design tool is free to use any character in Unicode. A
// terminal is not. So every character the bundle shows inside a terminal block
// has to be one the interface is actually allowed to print — otherwise the
// design system is quietly proposing a box on someone's screen.
func TestTheBundleOnlyDrawsCharactersATerminalCanPrint(t *testing.T) {
	d := DesignSystem{Tool: "democtl"}
	offenders := map[rune]string{}
	for _, page := range d.Pages() {
		for _, block := range terminalBlocks(page.HTML) {
			for _, r := range block {
				if !theme.DefaultGlyphs.Printable(r) {
					offenders[r] = page.Path
				}
			}
		}
	}
	for r, where := range offenders {
		t.Errorf("%s draws %q (U+%04X), which is not in the glyph set — "+
			"a font without it shows a box, so the design system is proposing something unbuildable",
			where, r, r)
	}
}

// Every card has to reach the Design System pane, which means the @dsCard
// marker has to be on the first line. A card that renders perfectly and never
// appears is the failure mode worth a test.
func TestEveryCardIsMarkedForTheDesignSystemPane(t *testing.T) {
	marker := regexp.MustCompile(`^<!-- @dsCard group="([^"]+)" -->$`)
	groups := map[string]bool{}
	for _, page := range (DesignSystem{}).Pages() {
		first, _, _ := strings.Cut(page.HTML, "\n")
		m := marker.FindStringSubmatch(first)
		if m == nil {
			t.Errorf("%s starts with %q, not a @dsCard marker", page.Path, first)
			continue
		}
		groups[m[1]] = true
	}
	if !groups["Foundations"] {
		t.Error("no card is grouped under Foundations")
	}
}

// A role that exists in the palette but is never drawn in a line of interface
// is a swatch nobody can see in use — and span panics on a role that does not
// exist, so this is the other half of that guard.
func TestEveryRoleIsDrawnAndNotOnlySwatched(t *testing.T) {
	var all strings.Builder
	for _, page := range (DesignSystem{}).Pages() {
		all.WriteString(page.HTML)
	}
	for _, role := range theme.Default.Roles() {
		if !strings.Contains(all.String(), `class="`+classFor(role.Name)+`"`) {
			t.Errorf("%s is in the palette but never drawn — show it in use or drop the role", role.Name)
		}
	}
}

// A tool's own palette has to reach the page, or the bundle documents tuikit's
// defaults while the tool ships something else — the drift this package exists
// to prevent.
func TestAToolsOwnPaletteIsWhatGetsDrawn(t *testing.T) {
	pal := theme.Default
	pal.Accent = lipgloss.Color("33") // #0087ff
	pal.Extra = []theme.Role{{Name: "Info", Color: lipgloss.Color("39"), Why: "a note the tool wants to make"}}

	html := DesignSystem{Tool: "democtl", Palette: pal}.Pages()[0].HTML

	if !strings.Contains(html, "#0087ff") {
		t.Error("the overridden accent is not in the colours page")
	}
	if strings.Contains(html, "#ff5faf") {
		t.Error("tuikit's default accent is still being drawn")
	}
	if !strings.Contains(html, `class="info"`) {
		t.Error("an Extra role got a swatch but is never drawn in use")
	}
}

// A tool's added glyph has to reach the page too, and the printable guard above
// has to accept it — otherwise adding a glyph deliberately fails the bundle.
func TestAnAddedGlyphIsDocumentedAndAllowed(t *testing.T) {
	glyphs := theme.DefaultGlyphs.With('×', "multiplication sign — replica counts")
	pages := DesignSystem{Glyphs: glyphs}.Pages()

	if !strings.Contains(pages[1].HTML, "U+00D7") {
		t.Error("the added glyph is missing from the glyphs page")
	}
	for _, page := range pages {
		for _, block := range terminalBlocks(page.HTML) {
			for _, r := range block {
				if !glyphs.Printable(r) {
					t.Errorf("%s draws %q, outside the tool's own set", page.Path, r)
				}
			}
		}
	}
}

// The zero value has to render, because that is what a tool gets before it has
// opinions.
func TestTheZeroValueRendersTheDefaults(t *testing.T) {
	pages := (DesignSystem{}).Pages()
	if len(pages) != 3 {
		t.Fatalf("got %d pages, want 3", len(pages))
	}
	for _, p := range pages {
		if len(p.HTML) < 500 {
			t.Errorf("%s rendered %d bytes — that is not a page", p.Path, len(p.HTML))
		}
	}
}

func TestSpanPanicsOnARoleThatDoesNotExist(t *testing.T) {
	defer func() {
		if recover() == nil {
			t.Error("a component drew with a role the palette does not have, and nothing stopped it")
		}
	}()
	(DesignSystem{}).renderer().span("chartreuse", "text")
}

// terminalBlocks returns the text inside each .term block, tags stripped. Only
// those blocks are a terminal; the prose around them is a web page and may say
// whatever it likes.
func terminalBlocks(page string) []string {
	var out []string
	re := regexp.MustCompile(`(?s)<div class="term[^"]*">(.*?)</div>`)
	for _, m := range re.FindAllStringSubmatch(page, -1) {
		text := regexp.MustCompile(`<[^>]*>`).ReplaceAllString(m[1], "")
		text = strings.NewReplacer("&lt;", "<", "&gt;", ">", "&amp;", "&", "&#34;", `"`, "&#39;", "'").Replace(text)
		out = append(out, text)
	}
	return out
}
