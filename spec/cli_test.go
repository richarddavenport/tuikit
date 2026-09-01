package spec

import (
	"bytes"
	"strings"
	"testing"
)

func tree() Command {
	return Command{
		Name:  "democtl",
		Short: "a tuikit example",
		Commands: []Command{
			{
				Name:  "status",
				Short: "current state of the fleet",
				Args:  []Arg{{Name: "service", Help: "which one"}},
				Flags: []Flag{
					{Name: "watch", Short: "w", Kind: Bool, Help: "keep looking"},
					{Name: "since", Kind: String, Default: "1h", Help: "how far back"},
				},
				Run: func(c Call) int {
					c.Out.Write([]byte("service=" + c.Arg("service") +
						" watch=" + c.Flag("watch") + " since=" + c.Flag("since")))
					return OK
				},
			},
			{
				Name:  "diff",
				Short: "what a deploy would change",
				Run:   func(Call) int { return Drift },
			},
			{
				Name:   "internal",
				Short:  "still works, not advertised",
				Hidden: true,
				Run:    func(Call) int { return OK },
			},
			{
				Name:  "deploy",
				Short: "push a new spec",
				Args:  []Arg{{Name: "service", Required: true}, {Name: "tags", Variadic: true}},
				Run: func(c Call) int {
					c.Out.Write([]byte(c.Arg("service") + " " + strings.Join(c.Rest, ",")))
					return OK
				},
			},
		},
	}
}

func run(argv ...string) (int, string, string) {
	var out, errw bytes.Buffer
	code := Run(tree(), argv, &out, &errw)
	return code, out.String(), errw.String()
}

func TestASubcommandGetsItsArgumentsAndFlags(t *testing.T) {
	code, out, _ := run("status", "api_gateway", "-w", "--since", "24h")

	if code != OK {
		t.Errorf("exit %d", code)
	}
	if want := "service=api_gateway watch=true since=24h"; out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

// A flag not given takes its declared default, so a command reads one value
// rather than checking whether it was set.
func TestAFlagNotGivenTakesItsDefault(t *testing.T) {
	_, out, _ := run("status", "api_gateway")
	if !strings.Contains(out, "since=1h") || !strings.Contains(out, "watch=false") {
		t.Errorf("got %q", out)
	}
}

// The exit-code contract, and the reason cobra is not here. Drift is a finding,
// not a failure: `diff && deploy` must not deploy when there is nothing to
// deploy, and `if diff; then` must tell "no drift" from "could not look".
func TestTheExitCodeContract(t *testing.T) {
	if code, _, _ := run("diff"); code != Drift {
		t.Errorf("a dry run that found drift exited %d, want %d", code, Drift)
	}
	if code, _, _ := run("status", "x"); code != OK {
		t.Errorf("a command that worked exited %d", code)
	}
	if code, _, _ := run("deploy"); code != Fail {
		t.Errorf("a missing required argument exited %d, want %d", code, Fail)
	}
}

// A wrapper that flattens somebody else's exit code breaks the thing it wraps,
// so a command returns whatever it likes and nothing here rewrites it.
func TestACommandsOwnCodeIsPassedThrough(t *testing.T) {
	root := Command{Name: "t", Commands: []Command{
		{Name: "script", Run: func(Call) int { return 137 }},
	}}
	var out, errw bytes.Buffer
	if code := Run(root, []string{"script"}, &out, &errw); code != 137 {
		t.Errorf("a script's own status became %d", code)
	}
}

// A missing required argument says WHICH one. "expected 2 args, got 1" makes
// the reader count.
func TestAMissingArgumentSaysWhichOne(t *testing.T) {
	code, _, errOut := run("deploy")
	if code != Fail {
		t.Errorf("exit %d", code)
	}
	if !strings.Contains(errOut, "missing argument service") {
		t.Errorf("the error does not name the argument: %q", errOut)
	}
	if !strings.Contains(errOut, "usage:") {
		t.Errorf("the error does not show usage: %q", errOut)
	}
}

func TestAnUnexpectedArgumentIsRejected(t *testing.T) {
	code, _, errOut := run("status", "one", "two")
	if code != Fail || !strings.Contains(errOut, `unexpected argument "two"`) {
		t.Errorf("exit %d, %q", code, errOut)
	}
}

// A variadic argument takes the rest.
func TestAVariadicArgumentTakesTheRest(t *testing.T) {
	_, out, _ := run("deploy", "api_gateway", "v2.4.1", "v2.4.0")
	if want := "api_gateway v2.4.1,v2.4.0"; out != want {
		t.Errorf("got %q, want %q", out, want)
	}
}

// Help is answered in one place, so no command can forget to support it and
// none can answer it differently.
func TestHelpIsAnsweredForEveryCommand(t *testing.T) {
	for _, argv := range [][]string{{"-h"}, {"--help"}, {"status", "-h"}, {"status", "--help"}} {
		code, out, _ := run(argv...)
		if code != OK {
			t.Errorf("%v exited %d", argv, code)
		}
		if !strings.Contains(out, "usage:") {
			t.Errorf("%v printed no usage: %q", argv, out)
		}
	}
}

// Usage is generated, so it cannot describe a flag the command does not have —
// which is what a hand-written usage string does about six months in.
func TestUsageComesFromTheDeclaration(t *testing.T) {
	_, out, _ := run("status", "-h")

	for _, want := range []string{
		"usage: democtl status [service] [flags]",
		"-w, --watch",
		"--since string",
		"(default 1h)",
		"--json",
	} {
		if !strings.Contains(out, want) {
			t.Errorf("%q is missing from usage:\n%s", want, out)
		}
	}
}

// Every command takes --json, so an agent never has to find out which ones do.
func TestEveryCommandTakesJSON(t *testing.T) {
	root := Command{Name: "t", Commands: []Command{
		{Name: "one", Run: func(c Call) int {
			if !c.JSON {
				t.Error("--json did not reach the command")
			}
			return OK
		}},
	}}
	var out, errw bytes.Buffer
	Run(root, []string{"one", "--json"}, &out, &errw)
}

// Hidden commands still work and are not advertised.
func TestAHiddenCommandWorksButIsNotListed(t *testing.T) {
	if code, _, _ := run("internal"); code != OK {
		t.Errorf("a hidden command exited %d", code)
	}
	_, out, _ := run("-h")
	if strings.Contains(out, "internal") {
		t.Errorf("a hidden command is listed:\n%s", out)
	}
	if !strings.Contains(out, "status") {
		t.Errorf("the visible commands are missing:\n%s", out)
	}
}

// A group named with no subcommand is not an error to guess about — it shows
// what it has.
func TestAGroupWithNoSubcommandShowsItsCommands(t *testing.T) {
	code, _, errOut := run()
	if code != Fail {
		t.Errorf("exit %d", code)
	}
	if !strings.Contains(errOut, "diff") {
		t.Errorf("the group did not list its commands:\n%s", errOut)
	}
}

// Flags work wherever they appear.
//
// stdlib flag stops at the first non-flag argument, so `deploy api_gateway
// --dry-run` would parse no flags at all and take --dry-run as a positional.
// cmd/tuikit already works around that by hand, with the note that "the
// directory goes after the flags" is a rule nobody remembers and nothing
// enforces. A generated CLI must not make every tool learn it again.
func TestFlagsWorkAfterPositionalArguments(t *testing.T) {
	for _, argv := range [][]string{
		{"status", "-w", "--since", "24h", "api_gateway"},
		{"status", "api_gateway", "-w", "--since", "24h"},
		{"status", "-w", "api_gateway", "--since=24h"},
	} {
		code, out, errOut := run(argv...)
		if code != OK {
			t.Errorf("%v exited %d: %s", argv, code, errOut)
			continue
		}
		if want := "service=api_gateway watch=true since=24h"; out != want {
			t.Errorf("%v gave %q, want %q", argv, out, want)
		}
	}
}

// A double dash ends the flags, because a service really can be called -w.
func TestADoubleDashEndsTheFlags(t *testing.T) {
	code, out, _ := run("status", "--", "-w")
	if code != OK || !strings.Contains(out, "service=-w") {
		t.Errorf("exit %d, got %q", code, out)
	}
	if strings.Contains(out, "watch=true") {
		t.Errorf("a positional after -- was read as a flag: %q", out)
	}
}

// An unknown flag is named, rather than being silently taken as a positional —
// which is the failure mode of a parser that stops at the first non-flag.
func TestAnUnknownFlagIsNamed(t *testing.T) {
	code, _, errOut := run("status", "--wach")
	if code != Fail {
		t.Errorf("exit %d", code)
	}
	if !strings.Contains(errOut, "unknown flag -wach") {
		t.Errorf("the error does not name the flag: %q", errOut)
	}
}

func TestAValueFlagWithNoValueSaysSo(t *testing.T) {
	code, _, errOut := run("status", "--since")
	if code != Fail || !strings.Contains(errOut, "needs a value") {
		t.Errorf("exit %d, %q", code, errOut)
	}
}
