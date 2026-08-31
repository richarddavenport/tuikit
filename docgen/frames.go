package docgen

import (
	"encoding/json"
	"fmt"
	"html"
	"os"
	"path/filepath"
	"strings"

	"github.com/richarddavenport/tuikit/harness"
)

// Frames turns a capture directory into a page.
//
// The capture writes `.ansi` files and a manifest; this reads both and produces
// something a person can look at and send to someone. It is the other half of
// the loop `tuikit watch` runs.
//
// Three rules live here rather than in whoever publishes the page, because they
// are properties of terminal frames rather than editorial choices:
//
//   - The frames keep a DARK ground in both light and dark themes. The ANSI was
//     captured for a dark terminal; recolouring it reports colours the tool does
//     not have.
//   - The system monospace stack, with no webfont. Box-drawing and braille must
//     share one set of advance widths, and a webfont plus a fallback for the
//     glyphs it lacks guarantees a frame that does not line up.
//   - Every frame shows its provenance. A page that quietly mixes real and
//     staged frames misreports the tool.
//
// The chrome around them stays quiet. The frames are the loud thing.
type Frames struct {
	// Title names the page. Empty becomes the directory's name.
	Title string
	// Lede is a sentence under the title. Optional.
	Lede string
	// Groups orders and explains the frames. A frame not named by any group is
	// appended to a trailing group, so a new screen appears rather than
	// vanishing — the failure mode of a hand-maintained list.
	//
	// Grouping is given rather than inferred because it carries meaning a
	// filename cannot: what the reader is doing.
	Groups []Group
}

// Group is a set of frames with a heading.
type Group struct {
	Title string
	Lede  string
	// Frames names them in order, by the name given to harness.Shot.
	Frames []FrameNote
}

// FrameNote is one frame in a group, with what a reader needs beside it.
type FrameNote struct {
	Name string
	// Keys is the sequence that reaches this screen, as a person would type it:
	// "tab right right". Optional.
	Keys string
	// Note says what the frame is showing. Optional, and worth writing.
	Note string
}

// Page renders the capture in dir.
func (f Frames) Page(dir string) (string, error) {
	manifest, err := readManifest(dir)
	if err != nil {
		return "", err
	}

	byName := make(map[string]harness.Frame, len(manifest.Frames))
	for _, fr := range manifest.Frames {
		byName[fr.Name] = fr
	}

	title := f.Title
	if title == "" {
		title = filepath.Base(dir)
	}

	var b strings.Builder
	fmt.Fprintf(&b, "<!doctype html>\n<html lang=\"en\">\n<head>\n<meta charset=\"utf-8\">\n"+
		"<meta name=\"viewport\" content=\"width=device-width, initial-scale=1\">\n"+
		"<title>%s</title>\n<style>%s</style>\n</head>\n<body>\n",
		html.EscapeString(title), pageCSS)

	fmt.Fprintf(&b, "<header>\n<h1>%s</h1>\n", html.EscapeString(title))
	if f.Lede != "" {
		fmt.Fprintf(&b, "<p class=\"lede\">%s</p>\n", html.EscapeString(f.Lede))
	}
	fmt.Fprintf(&b, "<p class=\"facts\">%d frames · %d&#215;%d</p>\n</header>\n",
		len(manifest.Frames), manifest.Width, manifest.Height)

	shown := map[string]bool{}
	for _, g := range f.Groups {
		var notes []FrameNote
		for _, note := range g.Frames {
			if _, ok := byName[note.Name]; ok {
				notes = append(notes, note)
				shown[note.Name] = true
			}
		}
		if len(notes) == 0 {
			continue
		}
		if err := writeGroup(&b, dir, g.Title, g.Lede, notes, byName); err != nil {
			return "", err
		}
	}

	// Anything the groups did not name. A frame that exists but is not listed
	// must still appear: a page that silently drops a new screen is worse than
	// one with an untidy section.
	var rest []FrameNote
	for _, fr := range manifest.Frames {
		if !shown[fr.Name] {
			rest = append(rest, FrameNote{Name: fr.Name})
		}
	}
	if len(rest) > 0 {
		title, lede := "Ungrouped", "Captured, but not named by any group."
		if len(f.Groups) == 0 {
			title, lede = "Frames", ""
		}
		if err := writeGroup(&b, dir, title, lede, rest, byName); err != nil {
			return "", err
		}
	}

	b.WriteString("</body>\n</html>\n")
	return b.String(), nil
}

// Write renders the page to a file.
func (f Frames) Write(dir, out string) error {
	page, err := f.Page(dir)
	if err != nil {
		return err
	}
	return os.WriteFile(out, []byte(page), 0o600)
}

type manifest struct {
	Width  int             `json:"width"`
	Height int             `json:"height"`
	Frames []harness.Frame `json:"frames"`
}

// readManifest fails loudly on a directory that holds no capture. A page
// rendered from nothing is a page that says the tool has no screens.
func readManifest(dir string) (manifest, error) {
	var m manifest
	body, err := os.ReadFile(filepath.Join(dir, "frames.json"))
	if os.IsNotExist(err) {
		return m, fmt.Errorf("docgen: no frames.json in %s — capture there first", dir)
	}
	if err != nil {
		return m, err
	}
	if err := json.Unmarshal(body, &m); err != nil {
		return m, fmt.Errorf("docgen: %s/frames.json: %w", dir, err)
	}
	if len(m.Frames) == 0 {
		return m, fmt.Errorf("docgen: %s captured no frames", dir)
	}
	return m, nil
}

func writeGroup(b *strings.Builder, dir, title, lede string, notes []FrameNote, byName map[string]harness.Frame) error {
	fmt.Fprintf(b, "<section>\n<div class=\"group\"><h2>%s</h2>", html.EscapeString(title))
	if lede != "" {
		fmt.Fprintf(b, "<p>%s</p>", html.EscapeString(lede))
	}
	b.WriteString("</div>\n")

	for _, note := range notes {
		fr := byName[note.Name]
		raw, err := os.ReadFile(filepath.Join(dir, fr.File))
		if err != nil {
			return fmt.Errorf("docgen: %w", err)
		}

		b.WriteString("<figure>\n<figcaption>\n<div class=\"id\">")
		fmt.Fprintf(b, "<h3>%s</h3>", html.EscapeString(fr.Name))
		if note.Keys != "" {
			b.WriteString(`<span class="keys">`)
			for _, k := range strings.Fields(note.Keys) {
				fmt.Fprintf(b, "<kbd>%s</kbd>", html.EscapeString(k))
			}
			b.WriteString(`</span>`)
		}
		b.WriteString("</div>\n")
		if note.Note != "" {
			fmt.Fprintf(b, "<p>%s</p>\n", html.EscapeString(note.Note))
		}
		fmt.Fprintf(b, "<p class=\"dims\"><span class=\"chip %s\">%s</span>%d&#215;%d</p>\n</figcaption>\n",
			html.EscapeString(string(fr.Mode)), html.EscapeString(string(fr.Mode)), fr.Width, fr.Lines)
		b.WriteString(harness.HTML(string(raw)))
		b.WriteString("\n</figure>\n")
	}
	b.WriteString("</section>\n")
	return nil
}
