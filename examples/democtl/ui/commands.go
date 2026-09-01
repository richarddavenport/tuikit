package ui

import (
	"fmt"
	"path/filepath"
	"sort"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/harness"

	"github.com/richarddavenport/tuikit/app"
	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
	"github.com/richarddavenport/tuikit/spec"
)

// Commands is democtl, declared once.
//
// The CLI below, the context menu on a service row, the keystroke that does the
// same thing, and the entry in `democtl describe --json` all come from here.
// Not four descriptions kept in step — one, read four ways.
//
// The Key on a command with a Target is not optional and guard.Reachable says
// so: an agent cannot click, and a multiplexer may eat the right-click before
// democtl ever sees it.
func Commands(seed int64) spec.Command {
	f := fleet.New(seed)

	return spec.Command{
		Name:  "democtl",
		Short: "a fictional service fleet, and tuikit's example tool",
		Long: "Everything here is generated from a seed, so it needs no backend and\n" +
			"shows the same thing every run.",
		Commands: []spec.Command{
			{
				Name:  "tui",
				Short: "Open the interface",
				// A Screen, so spec gives it --snapshot and --script and an
				// agent can capture democtl without writing a test.
				Screen: "dashboard",
				Run:    func(c spec.Call) int { return open(seed, c) },
			},
			{
				Name:  "status",
				Short: "Current state",
				Args:  []spec.Arg{{Name: "service", Help: "one service, or all of them", Complete: names(f)}},
				Flags: []spec.Flag{{Name: "failing", Kind: spec.Bool, Help: "only what needs attention"}},
				Run:   func(c spec.Call) int { return status(f, c) },
			},
			{
				Name:  "logs",
				Short: "View logs",
				Args:  []spec.Arg{{Name: "service", Required: true, Complete: names(f)}},
				// The other three surfaces: a screen to open, a region whose
				// menu it belongs in, and the key that does the same thing.
				Screen: "logs",
				Target: regServicesRow,
				Key:    "L",
				Run:    func(c spec.Call) int { return logsOf(f, c) },
			},
			{
				Name:   "deploy",
				Short:  "Deploy",
				Args:   []spec.Arg{{Name: "service", Required: true, Complete: names(f)}},
				Flags:  []spec.Flag{{Name: "dry-run", Kind: spec.Bool, Help: "say what would change"}},
				Screen: "run",
				Target: regServicesRow,
				Key:    "D",
				Run:    func(c spec.Call) int { return deploy(f, c) },
			},
		},
	}
}

// names completes a service name from the fleet — the same completer the shell
// and the interface both use, so they cannot offer different answers.
func names(f fleet.Fleet) spec.Completer {
	return func(string) []string {
		out := make([]string, 0, len(f.Services))
		for _, s := range f.Services {
			out = append(out, s.Name)
		}
		sort.Strings(out)
		return out
	}
}

func status(f fleet.Fleet, c spec.Call) int {
	for _, s := range f.Services {
		if want := c.Arg("service"); want != "" && s.Name != want {
			continue
		}
		if c.Bool("failing") && s.State == fleet.Running {
			continue
		}
		_, _ = fmt.Fprintf(c.Out, "%-16s %d/%d %s\n", s.Name, s.Ready, s.Want, s.State)
	}
	return spec.OK
}

func logsOf(f fleet.Fleet, c spec.Call) int {
	name := c.Arg("service")
	if !known(f, name) {
		_, _ = fmt.Fprintf(c.Err, "democtl: no service %q\n", name)
		return spec.Fail
	}
	for _, l := range f.Logs(name, 40) {
		_, _ = fmt.Fprintf(c.Out, "%s %s\n", l.At.Format("15:04:05"), l.Text)
	}
	return spec.OK
}

// deploy is where the exit-code contract earns its keep: a dry run that found
// something to do reports Drift, which is a finding rather than a failure. A
// tool that reported it as failure would make every `deploy --dry-run && ...`
// wrong in the same direction.
func deploy(f fleet.Fleet, c spec.Call) int {
	name := c.Arg("service")
	if !known(f, name) {
		_, _ = fmt.Fprintf(c.Err, "democtl: no service %q\n", name)
		return spec.Fail
	}
	plan := fleet.Deploy(name, false)
	for _, step := range plan.Steps {
		if c.Bool("dry-run") && step.Skip {
			continue
		}
		_, _ = fmt.Fprintf(c.Out, "%s\n", step.Name)
	}
	if c.Bool("dry-run") {
		return spec.Drift
	}
	return spec.OK
}

func known(f fleet.Fleet, name string) bool {
	for _, s := range f.Services {
		if s.Name == name {
			return true
		}
	}
	return false
}

// act runs the TUI's version of a command, by the key that reaches it.
//
// By KEY, so the menu entry and the keystroke cannot run different code: the
// menu is built from the same declarations and carries the same key, so
// choosing "View logs" and pressing L are the same call.
func (m *Model) act(key string) tea.Cmd {
	switch key {
	case "L":
		if svc, ok := m.selected(); ok {
			m.log.Follow = true
			m.stack.Push(app.Screen(screenLogs), "logs "+svc.Name)
		}
	case "D":
		m.confirmDeploy()
	}
	return nil
}

// open runs the interface, or captures it.
//
// The whole of what a tool writes to get --snapshot: build the model, and if a
// directory was asked for, drive it rather than opening it. Decision 10's
// second capture mechanism, in nine lines.
func open(seed int64, c spec.Call) int {
	m := New(seed)
	dir := c.Flag("snapshot")
	if dir == "" {
		p := tea.NewProgram(m, tea.WithAltScreen(), tea.WithMouseCellMotion())
		if _, err := p.Run(); err != nil {
			_, _ = fmt.Fprintln(c.Err, "democtl:", err)
			return spec.Fail
		}
		return spec.OK
	}

	script, err := harness.ScriptFile(c.Flag("script"))
	if err != nil {
		_, _ = fmt.Fprintln(c.Err, "democtl:", err)
		return spec.Fail
	}
	frames, err := harness.Snapshot(m, dir, script)
	if err != nil {
		_, _ = fmt.Fprintln(c.Err, "democtl:", err)
		return spec.Fail
	}
	for _, f := range frames {
		_, _ = fmt.Fprintln(c.Out, filepath.Join(dir, f.File))
	}
	return spec.OK
}
