package spec

import (
	"encoding/json"
	"strings"
	"testing"

	"github.com/richarddavenport/tuikit/theme"
)

func manifest(t *testing.T) Manifest {
	t.Helper()
	root := tree()
	root.Commands = append(root.Commands, Command{
		Name: "logs", Short: "read a service's output",
		Screen: "logs", Target: "services.row", Key: "L",
		Run: func(Call) int { return OK },
	})
	return Describe(root, "v1.2.3", theme.Default, theme.DefaultGlyphs)
}

// The whole surface in one call. An agent otherwise greps a README a version
// behind, parses --help, or reads the source — three answers that disagree in
// small ways with no way to tell which is current.
func TestTheManifestNamesEveryCommand(t *testing.T) {
	m := manifest(t)

	names := map[string]CommandInfo{}
	for _, c := range m.Commands {
		names[c.Name] = c
	}
	for _, want := range []string{"democtl status", "democtl diff", "democtl deploy", "democtl logs"} {
		if _, ok := names[want]; !ok {
			t.Errorf("%q is missing from the manifest", want)
		}
	}
	status := names["democtl status"]
	if len(status.Flags) != 2 || status.Flags[0].Name != "watch" || status.Flags[0].Short != "w" {
		t.Errorf("status's flags are %+v", status.Flags)
	}
	if len(status.Args) != 1 || status.Args[0].Name != "service" {
		t.Errorf("status's arguments are %+v", status.Args)
	}
}

// A hidden command is not advertised here either, or "hidden" would mean
// "hidden from people".
func TestAHiddenCommandIsNotInTheManifest(t *testing.T) {
	for _, c := range manifest(t).Commands {
		if strings.Contains(c.Name, "internal") {
			t.Errorf("a hidden command is in the manifest: %q", c.Name)
		}
	}
}

// The four surfaces are reported together, so an agent can see that the CLI
// command, the screen, the menu entry and the keystroke are one thing.
func TestACommandReportsAllFourSurfaces(t *testing.T) {
	m := manifest(t)
	var logs CommandInfo
	for _, c := range m.Commands {
		if c.Name == "democtl logs" {
			logs = c
		}
	}
	if logs.Screen != "logs" || logs.Target != "services.row" || logs.Key != "L" {
		t.Errorf("logs reports %+v", logs)
	}
	if len(m.Screens) == 0 || m.Screens[0] != "logs" {
		t.Errorf("screens are %v", m.Screens)
	}
}

// Regions say what can be done to them, which is what makes `click
// services.row[2]` writable by something that has never seen the screen.
func TestRegionsSayWhatCanBeDoneToThem(t *testing.T) {
	for _, r := range manifest(t).Regions {
		if r.Name != "services.row" {
			continue
		}
		if len(r.Commands) != 1 || r.Commands[0] != "democtl logs" {
			t.Errorf("services.row offers %v", r.Commands)
		}
		return
	}
	t.Error("no regions in the manifest")
}

// The roles and glyphs are in it because an agent asked to add a screen needs
// to know which exist before it writes anything — and because those are the two
// things a guard will fail it for afterwards.
func TestTheVocabularyIsInTheManifest(t *testing.T) {
	m := manifest(t)

	if len(m.Roles) < 9 {
		t.Errorf("only %d roles", len(m.Roles))
	}
	for _, r := range m.Roles {
		if r.Name == "" || r.Hex == "" || r.Why == "" {
			t.Errorf("role %+v is missing something", r)
		}
		if r.Name == "Accent" && !strings.HasPrefix(r.Hex, "#") {
			t.Errorf("Accent's color is %q", r.Hex)
		}
	}
	if len(m.Glyphs) == 0 {
		t.Error("no glyphs")
	}
	for _, g := range m.Glyphs {
		if g.Why == "" {
			t.Errorf("glyph %q has no reason", g.Glyph)
		}
	}
}

// Completions are LIVE, so the manifest says a value can be offered without
// running a config lookup to produce one.
func TestTheManifestSaysWhatCompletesWithoutCompletingIt(t *testing.T) {
	root := Command{Name: "t", Commands: []Command{{
		Name: "one",
		Args: []Arg{{Name: "env", Complete: func(string) []string { t.Error("the manifest ran a completer"); return nil }}},
		Run:  func(Call) int { return OK },
	}}}
	m := Describe(root, "", theme.Default, theme.DefaultGlyphs)

	if !m.Commands[0].Args[0].Completes {
		t.Error("an argument with a completer does not say so")
	}
}

// It has to survive a round trip, because something else is going to parse it.
func TestTheManifestIsValidJSON(t *testing.T) {
	raw, err := manifest(t).JSON()
	if err != nil {
		t.Fatal(err)
	}
	var back Manifest
	if err := json.Unmarshal(raw, &back); err != nil {
		t.Fatalf("the manifest does not parse: %v", err)
	}
	if back.Tool != "democtl" || back.Version != "v1.2.3" {
		t.Errorf("round trip lost the tool: %+v", back)
	}
}
