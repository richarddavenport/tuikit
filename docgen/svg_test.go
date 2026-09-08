package docgen

import (
	"encoding/xml"
	"strings"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"
)

// frame is real ANSI from lipgloss rather than a hand-written escape string,
// because the sequences are the thing being read.
func frame(t *testing.T) string {
	t.Helper()
	lipgloss.SetColorProfile(termenv.TrueColor)
	pink := lipgloss.NewStyle().Foreground(lipgloss.Color("205"))
	sel := lipgloss.NewStyle().Foreground(lipgloss.Color("0")).Background(lipgloss.Color("7")).Bold(true)
	return "┌──────┐\n│ " + pink.Render("api") + "  │\n│ " + sel.Render("web") + "  │\n└──────┘"
}

// An SVG that is not well-formed XML renders as nothing at all, with no error
// anyone sees — the failure this test exists to make loud.
func TestTheSVGIsWellFormedXML(t *testing.T) {
	if err := xml.Unmarshal([]byte(SVG(frame(t))), new(any)); err != nil {
		t.Fatalf("the SVG does not parse: %v", err)
	}
}

// The color has to survive, or the image reports a tool that draws in gray.
func TestTheSVGCarriesTheColorThrough(t *testing.T) {
	out := SVG(frame(t))

	pink := hexOf(lipgloss.Color("205"))
	if !strings.Contains(out, `fill="`+pink+`"`) {
		t.Errorf("the foreground %s is not in the image", pink)
	}
	// A background is a rect, because text has no background in SVG. Without
	// this a selected row loses the one thing that marks it as selected.
	if !strings.Contains(out, "<rect") || strings.Count(out, "<rect") < 2 {
		t.Errorf("no background rect was drawn: %s", out)
	}
	if !strings.Contains(out, `font-weight="600"`) {
		t.Error("bold was dropped")
	}
}

// The grid rule: every run states the width it must occupy, so the frame lines
// up whatever font the reader has. A renderer that trusts the font's advance
// width is one box-drawing character away from a frame that does not meet.
func TestEverySpanStatesItsWidth(t *testing.T) {
	out := SVG(frame(t))

	texts := strings.Count(out, "<text")
	lengths := strings.Count(out, "textLength=")
	if texts == 0 {
		t.Fatal("no text was drawn")
	}
	if texts != lengths {
		t.Errorf("%d text runs but %d state a width — the rest are at the font's mercy", texts, lengths)
	}
	// spacingAndGlyphs would meet the width by distorting the glyphs, which
	// bends box-drawing characters into not joining.
	if strings.Contains(out, "spacingAndGlyphs") {
		t.Error("the glyphs are being stretched to fit; adjust the spacing instead")
	}
}

// No stylesheet. An SVG in a Markdown document is rendered by a sanitiser, and
// a stripped <style> block leaves a frame that is all one color with nothing
// to say why.
func TestTheSVGCarriesNoStylesheet(t *testing.T) {
	if out := SVG(frame(t)); strings.Contains(out, "<style") {
		t.Error("the image depends on a stylesheet a sanitiser may strip")
	}
}

// A frame draws whatever the tool draws, and < and & are characters a tool
// draws. Unescaped, the first one ends the image.
func TestTheSVGEscapesWhatTheTerminalDrew(t *testing.T) {
	out := SVG(`a < b & c "d" 'e'`)

	if err := xml.Unmarshal([]byte(out), new(any)); err != nil {
		t.Fatalf("a frame containing < and & produced broken XML: %v", err)
	}
	if strings.Contains(out, "a < b") {
		t.Error("the frame's own angle bracket reached the markup unescaped")
	}
}

// The image states its size in cells-turned-pixels, so a 132-column frame and
// an 80-column one are not scaled to the same width by the document.
func TestTheImageIsSizedFromTheGrid(t *testing.T) {
	narrow := SVG("ab\ncd")
	wide := SVG(strings.Repeat("x", 80))

	if w := attr(narrow, "width"); w == "" {
		t.Fatal("the image has no width")
	}
	if attr(narrow, "width") == attr(wide, "width") {
		t.Errorf("a 2-column frame and an 80-column one are the same width (%s)", attr(narrow, "width"))
	}
	if attr(narrow, "viewBox") == "" {
		t.Error("no viewBox, so the image cannot scale in a document")
	}
}

// An image with no text alternative is one a reader using a screen reader is
// simply not shown.
func TestTheImageHasATitle(t *testing.T) {
	if !strings.Contains(SVG("hello"), "<title>") {
		t.Error("the image has no title")
	}
}

// A trailing run of spaces is not drawn. It cannot be seen, and a frame is
// mostly blank — writing them is a file several times the size for nothing.
func TestTrailingBlanksAreNotDrawn(t *testing.T) {
	padded := SVG("hi" + strings.Repeat(" ", 200))

	if strings.Contains(padded, strings.Repeat(" ", 20)) {
		t.Error("a run of trailing spaces was written into the image")
	}
	if !strings.Contains(padded, ">hi<") {
		t.Error("the text itself was lost")
	}
}

func attr(svg, name string) string {
	_, rest, ok := strings.Cut(svg, name+`="`)
	if !ok {
		return ""
	}
	val, _, _ := strings.Cut(rest, `"`)
	return val
}
