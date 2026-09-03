package docgen

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
)

// The whole point: a document with an image per frame, and every image on disk
// where the document says it is. A link to a file that is not there renders as
// a broken image, which is worse than no document.
func TestMarkdownWritesAnImagePerFrameAndLinksItCorrectly(t *testing.T) {
	dir := capture(t)
	out := filepath.Join(t.TempDir(), "screens.md")

	if err := (Frames{Title: "democtl"}).Markdown(dir, out); err != nil {
		t.Fatal(err)
	}

	doc, err := os.ReadFile(out)
	if err != nil {
		t.Fatal(err)
	}
	for _, name := range []string{"first", "staged", "third"} {
		link := "(screens/" + name + ".svg)"
		if !strings.Contains(string(doc), link) {
			t.Errorf("%s is not linked in the document", name)
		}
		// The link is relative to the document, so resolve it the way a reader
		// would rather than trusting the path that was written.
		if _, err := os.Stat(filepath.Join(filepath.Dir(out), "screens", name+".svg")); err != nil {
			t.Errorf("%s is linked but not written: %v", name, err)
		}
	}
}

// The images go beside the document, in a directory named after it, so moving
// or deleting the document is one operation rather than two.
func TestTheImagesAreNamedAfterTheDocument(t *testing.T) {
	dir := capture(t)
	base := t.TempDir()
	out := filepath.Join(base, "how-it-looks.md")

	if err := (Frames{}).Markdown(dir, out); err != nil {
		t.Fatal(err)
	}
	if _, err := os.Stat(filepath.Join(base, "how-it-looks")); err != nil {
		t.Errorf("the images are not beside the document: %v", err)
	}
}

// A note is what a reader who cannot see the image gets. "third" tells them
// nothing they could not have guessed from the heading.
func TestTheNoteBecomesTheAltText(t *testing.T) {
	dir := capture(t)
	out := filepath.Join(t.TempDir(), "screens.md")

	err := Frames{Groups: []Group{{
		Title:  "Looking around",
		Frames: []FrameNote{{Name: "first", Note: "the list, with a service selected", Keys: "j j"}},
	}}}.Markdown(dir, out)
	if err != nil {
		t.Fatal(err)
	}

	doc, _ := os.ReadFile(out)
	if !strings.Contains(string(doc), "![the list, with a service selected](") {
		t.Errorf("the note is not the alt text:\n%s", doc)
	}
	if !strings.Contains(string(doc), "## Looking around") {
		t.Error("the group's title is not a heading")
	}
	if !strings.Contains(string(doc), "`j`") {
		t.Error("the keys that reach the screen are not in the document")
	}
}

// The same rule the page keeps: a frame that exists but no group names must
// still appear. A document that silently drops a new screen is worse than one
// with an untidy section.
func TestAFrameNoGroupNamesStillAppearsInMarkdown(t *testing.T) {
	dir := capture(t)
	out := filepath.Join(t.TempDir(), "screens.md")

	err := Frames{Groups: []Group{{
		Title:  "One of them",
		Frames: []FrameNote{{Name: "first"}},
	}}}.Markdown(dir, out)
	if err != nil {
		t.Fatal(err)
	}

	doc, _ := os.ReadFile(out)
	for _, name := range []string{"first", "staged", "third"} {
		if !strings.Contains(string(doc), name+".svg") {
			t.Errorf("%s was dropped from the document", name)
		}
	}
	if !strings.Contains(string(doc), "Ungrouped") {
		t.Error("the frames no group named are not marked as such")
	}
}

// A staged frame says so, because a document that quietly mixes real and
// composed screens misreports the tool. Same rule as the page.
func TestStagedFramesSaySo(t *testing.T) {
	dir := capture(t)
	out := filepath.Join(t.TempDir(), "screens.md")

	if err := (Frames{}).Markdown(dir, out); err != nil {
		t.Fatal(err)
	}
	doc, _ := os.ReadFile(out)
	if !strings.Contains(string(doc), "composed") {
		t.Errorf("the composed frame does not declare itself:\n%s", doc)
	}
}

// No escape sequence reaches the document or the images. A frame's raw ANSI in
// a Markdown file renders as mojibake.
func TestNoEscapeSequenceReachesTheDocument(t *testing.T) {
	dir := capture(t)
	base := t.TempDir()
	out := filepath.Join(base, "screens.md")

	if err := (Frames{}).Markdown(dir, out); err != nil {
		t.Fatal(err)
	}

	err := filepath.Walk(base, func(path string, info os.FileInfo, err error) error {
		if err != nil || info.IsDir() {
			return err
		}
		body, err := os.ReadFile(path)
		if err != nil {
			return err
		}
		if strings.ContainsRune(string(body), 0x1b) {
			t.Errorf("%s contains a raw escape sequence", filepath.Base(path))
		}
		return nil
	})
	if err != nil {
		t.Fatal(err)
	}
}

// An empty directory is an error, not an empty document — the same call the
// page makes. A document rendered from nothing says the tool has no screens.
func TestMarkdownFromAnEmptyDirectoryIsAnError(t *testing.T) {
	out := filepath.Join(t.TempDir(), "screens.md")
	if err := (Frames{}).Markdown(t.TempDir(), out); err == nil {
		t.Error("a capture directory with no frames produced a document")
	}
}
