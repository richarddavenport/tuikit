package spec

import (
	"encoding/json"
	"sort"

	"github.com/richarddavenport/tuikit/theme"
)

// Manifest is the whole surface of a tool, in one call.
//
// The point is the one call. An agent working out what a tool can do otherwise
// greps a README that is a version behind, or parses --help, or reads the
// source — three answers that disagree in small ways, and no way to tell which
// is current. Everything here is generated from the same declarations the CLI
// runs and the TUI draws, so it cannot be stale without them being stale too.
//
// The palette and glyph set are in it because an agent asked to add a screen
// needs to know which colour roles exist and which characters it may print
// before it writes anything — and because those are the two things a guard will
// fail it for afterwards.
type Manifest struct {
	Tool     string        `json:"tool"`
	Version  string        `json:"version,omitempty"`
	Commands []CommandInfo `json:"commands"`
	Screens  []string      `json:"screens,omitempty"`
	Regions  []RegionInfo  `json:"regions,omitempty"`
	Roles    []RoleInfo    `json:"roles,omitempty"`
	Glyphs   []GlyphInfo   `json:"glyphs,omitempty"`
}

// CommandInfo is one command, flattened to the name you would type.
type CommandInfo struct {
	Name  string     `json:"name"`
	Short string     `json:"short,omitempty"`
	Long  string     `json:"long,omitempty"`
	Args  []ArgInfo  `json:"args,omitempty"`
	Flags []FlagInfo `json:"flags,omitempty"`
	// Screen, Target and Key are the other three surfaces, reported so an
	// agent can see that the CLI command, the screen, the menu entry and the
	// keystroke are one thing rather than four.
	Screen string `json:"screen,omitempty"`
	Target string `json:"target,omitempty"`
	Key    string `json:"key,omitempty"`
}

// ArgInfo is a positional argument.
type ArgInfo struct {
	Name     string `json:"name"`
	Help     string `json:"help,omitempty"`
	Required bool   `json:"required,omitempty"`
	Variadic bool   `json:"variadic,omitempty"`
	// Completes says values can be offered for it, without offering them:
	// completions are live, and a manifest is not the place to run a config
	// lookup.
	Completes bool `json:"completes,omitempty"`
}

// FlagInfo is an option.
type FlagInfo struct {
	Name      string `json:"name"`
	Short     string `json:"short,omitempty"`
	Kind      string `json:"kind"`
	Default   string `json:"default,omitempty"`
	Help      string `json:"help,omitempty"`
	Completes bool   `json:"completes,omitempty"`
}

// RegionInfo is a named region and what can be done to it — which is what makes
// `click services.row[2]` writable by something that has never seen the screen.
type RegionInfo struct {
	Name     string   `json:"name"`
	Commands []string `json:"commands,omitempty"`
}

// RoleInfo is one colour role, by name and reason. Never by hue.
type RoleInfo struct {
	Name string `json:"name"`
	Hex  string `json:"hex"`
	Why  string `json:"why,omitempty"`
}

// GlyphInfo is one allowed character and why it is allowed.
type GlyphInfo struct {
	Glyph string `json:"glyph"`
	Why   string `json:"why"`
}

// Describe builds the manifest.
func Describe(root Command, version string, p theme.Palette, glyphs theme.GlyphSet) Manifest {
	m := Manifest{Tool: root.Name, Version: version}

	screens := map[string]bool{}
	regions := map[string][]string{}
	walk(root, nil, func(cmd Command, path []string) {
		if cmd.Hidden {
			return
		}
		info := CommandInfo{
			Name:   join(path),
			Short:  cmd.Short,
			Long:   cmd.Long,
			Screen: cmd.Screen,
			Target: string(cmd.Target),
			Key:    cmd.Key,
		}
		for _, a := range cmd.Args {
			info.Args = append(info.Args, ArgInfo{
				Name: a.Name, Help: a.Help, Required: a.Required,
				Variadic: a.Variadic, Completes: a.Complete != nil,
			})
		}
		for _, f := range cmd.Flags {
			info.Flags = append(info.Flags, FlagInfo{
				Name: f.Name, Short: f.Short, Kind: f.Kind.String(),
				Default: f.Default, Help: f.Help, Completes: f.Complete != nil,
			})
		}
		if cmd.Run != nil || len(cmd.Args) > 0 || len(cmd.Flags) > 0 {
			m.Commands = append(m.Commands, info)
		}
		if cmd.Screen != "" {
			screens[cmd.Screen] = true
		}
		if cmd.Target != "" {
			regions[string(cmd.Target)] = append(regions[string(cmd.Target)], info.Name)
		}
	})

	m.Screens = sorted(screens)
	for _, name := range sortedKeys(regions) {
		m.Regions = append(m.Regions, RegionInfo{Name: name, Commands: regions[name]})
	}
	for _, role := range p.Roles() {
		m.Roles = append(m.Roles, RoleInfo{Name: role.Name, Hex: theme.Hex(role.Color), Why: role.Why})
	}
	for _, glyph := range sortedKeys(glyphs) {
		m.Glyphs = append(m.Glyphs, GlyphInfo{Glyph: string(glyph), Why: glyphs[glyph]})
	}
	return m
}

// JSON is the manifest as an agent reads it: indented, because it is meant to
// be looked at as well as parsed.
func (m Manifest) JSON() ([]byte, error) { return json.MarshalIndent(m, "", "  ") }

// walk visits every command with the path that reaches it.
func walk(cmd Command, path []string, fn func(Command, []string)) {
	path = append(path, cmd.Name)
	fn(cmd, path)
	for _, sub := range cmd.Commands {
		walk(sub, path, fn)
	}
}

func join(path []string) string {
	out := ""
	for i, p := range path {
		if i > 0 {
			out += " "
		}
		out += p
	}
	return out
}

func sorted(set map[string]bool) []string {
	out := make([]string, 0, len(set))
	for k := range set {
		out = append(out, k)
	}
	sort.Strings(out)
	return out
}

func sortedKeys[K ~string | ~rune, V any](m map[K]V) []K {
	out := make([]K, 0, len(m))
	for k := range m {
		out = append(out, k)
	}
	sort.Slice(out, func(i, j int) bool { return out[i] < out[j] })
	return out
}
