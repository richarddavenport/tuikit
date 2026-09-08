package harness

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/comp"
)

// snapModel is a two-row list that can be moved and clicked, which is enough to
// drive a script through every verb.
type snapModel struct {
	cursor int
	canvas *comp.Canvas
}

func (m *snapModel) SetSize(int, int) {}

func (m *snapModel) Init() tea.Cmd { return nil }

func (m *snapModel) Update(msg tea.Msg) (tea.Model, tea.Cmd) {
	switch msg := msg.(type) {
	case tea.KeyMsg:
		if msg.String() == "j" {
			m.cursor++
		}
	case tea.MouseMsg:
		if id := m.canvas.OwnerAt(msg.X, msg.Y); id.Name == "row" {
			m.cursor = id.Index
		}
	}
	return m, nil
}

func (m *snapModel) View() string {
	c := comp.NewCanvas(20, 3)
	for i := range 3 {
		text := "  row " + string(rune('a'+i))
		if i == m.cursor {
			text = "> row " + string(rune('a'+i))
		}
		c.Text(0, i, text, nil, comp.Region("row").At(i))
	}
	m.canvas = c
	return c.String()
}

func (m *snapModel) Canvas() *comp.Canvas { return m.canvas }

// The whole point: an agent captures the interface from the binary, with no
// test involved and nothing to know beyond --help.
func TestSnapshotWritesTheFramesAScriptAsksFor(t *testing.T) {
	dir := t.TempDir()
	frames, err := Snapshot(&snapModel{}, dir, "press j; shot moved; click row[2]; shot clicked")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if len(frames) != 2 {
		t.Fatalf("got %d frames, want 2", len(frames))
	}

	moved := read(t, dir, "moved.ansi")
	if !strings.Contains(moved, "> row b") {
		t.Errorf("the key did not drive the model:\n%s", moved)
	}
	clicked := read(t, dir, "clicked.ansi")
	if !strings.Contains(clicked, "> row c") {
		t.Errorf("the click did not land on the row it named:\n%s", clicked)
	}
}

// A caller who asked for a snapshot and gets nothing has been told nothing.
func TestAScriptWithNoShotStillCapturesOne(t *testing.T) {
	dir := t.TempDir()
	frames, err := Snapshot(&snapModel{}, dir, "press j j")
	if err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if len(frames) != 1 || frames[0].Name != "screen" {
		t.Fatalf("got %+v, want one frame called screen", frames)
	}
	if got := read(t, dir, "screen.ansi"); !strings.Contains(got, "> row c") {
		t.Errorf("the keys did not run before the shot:\n%s", got)
	}
}

// An empty script is the ordinary case: show me the first screen.
func TestAnEmptyScriptCapturesTheOpeningScreen(t *testing.T) {
	dir := t.TempDir()
	if _, err := Snapshot(&snapModel{}, dir, ""); err != nil {
		t.Fatalf("snapshot: %v", err)
	}
	if got := read(t, dir, "screen.ansi"); !strings.Contains(got, "> row a") {
		t.Errorf("got:\n%s", got)
	}
}

// The manifest is what docgen reads and what an agent reads to find out which
// frames exist without listing a directory and guessing at names.
func TestSnapshotWritesAManifest(t *testing.T) {
	dir := t.TempDir()
	if _, err := Snapshot(&snapModel{}, dir, "shot one; press j; shot two"); err != nil {
		t.Fatalf("snapshot: %v", err)
	}

	body, err := os.ReadFile(filepath.Join(dir, "frames.json"))
	if err != nil {
		t.Fatal(err)
	}
	var manifest struct {
		Frames []Frame `json:"frames"`
	}
	if err := json.Unmarshal(body, &manifest); err != nil {
		t.Fatalf("the manifest is not JSON: %v", err)
	}
	if len(manifest.Frames) != 2 {
		t.Errorf("the manifest lists %d frames", len(manifest.Frames))
	}
}

// A script that has drifted from the interface says which line and which
// region, rather than clicking column zero and reporting success.
func TestABadScriptSaysWhichLine(t *testing.T) {
	dir := t.TempDir()
	_, err := Snapshot(&snapModel{}, dir, "press j\nclick nowhere[0]\nshot never")
	if err == nil {
		t.Fatal("a script naming a region that was never drawn succeeded")
	}
	if !strings.Contains(err.Error(), "line 2") || !strings.Contains(err.Error(), "nowhere") {
		t.Errorf("the error does not say where: %v", err)
	}
	if _, err := os.Stat(filepath.Join(dir, "never.ansi")); err == nil {
		t.Error("it kept going after the failure")
	}
}

// Newlines or semicolons, so the same script works in a file and as a single
// shell argument.
func TestAScriptIsTheSameInAFileAndOnACommandLine(t *testing.T) {
	one, two := t.TempDir(), t.TempDir()
	if _, err := Snapshot(&snapModel{}, one, "press j\n# a comment\nshot x"); err != nil {
		t.Fatal(err)
	}
	if _, err := Snapshot(&snapModel{}, two, "press j; # a comment; shot x"); err != nil {
		t.Fatal(err)
	}
	if read(t, one, "x.ansi") != read(t, two, "x.ansi") {
		t.Error("the same script produced different frames")
	}
}

// --script takes a path or the script itself, so there is no second flag whose
// only job is to say which of the two this is.
func TestScriptFileTakesEither(t *testing.T) {
	dir := t.TempDir()
	path := filepath.Join(dir, "keys.txt")
	if err := os.WriteFile(path, []byte("press j\nshot x\n"), 0o600); err != nil {
		t.Fatal(err)
	}

	fromFile, err := ScriptFile(path)
	if err != nil || !strings.Contains(fromFile, "press j") {
		t.Errorf("a path did not read as one: %q, %v", fromFile, err)
	}
	inline, err := ScriptFile("press j; shot x")
	if err != nil || inline != "press j; shot x" {
		t.Errorf("a script read as a path: %q, %v", inline, err)
	}
}

func read(t *testing.T, dir, name string) string {
	t.Helper()
	body, err := os.ReadFile(filepath.Join(dir, name))
	if err != nil {
		t.Fatal(err)
	}
	return Strip(string(body))
}

// Issue 65: the silent failure is that a model loading in Init captures an
// empty screen, and nothing says so.
func TestLoadsBeforeCaptureNamesThePendingInit(t *testing.T) {
	if why, ok := LoadsBeforeCapture(&loadsInInit{}); ok {
		t.Error("a model whose Init returns a command was called ready")
	} else if !strings.Contains(why, "Init") {
		t.Errorf("the reason does not mention Init: %q", why)
	}
}

// A model with nothing pending is ready, and so is one with no Init at all —
// the check must not force every tool to grow a method.
func TestLoadsBeforeCaptureAcceptsAReadyModel(t *testing.T) {
	if _, ok := LoadsBeforeCapture(&snapModel{}); !ok {
		t.Error("a model whose Init returns nil was called not ready")
	}
	if _, ok := LoadsBeforeCapture(struct{}{}); !ok {
		t.Error("a model with no Init at all was called not ready")
	}
}

type loadsInInit struct{ snapModel }

func (*loadsInInit) Init() tea.Cmd {
	return func() tea.Msg { return nil }
}
