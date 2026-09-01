package spec

import (
	"bytes"
	"slices"
	"strings"
	"testing"
)

func completing() Command {
	envs := func(string) []string { return []string{"staging", "production", "prod-eu"} }
	return Command{Name: "democtl", Commands: []Command{
		{
			Name: "status", Short: "state of the fleet",
			Args:  []Arg{{Name: "service", Complete: func(string) []string { return []string{"api_gateway", "api_worker"} }}},
			Flags: []Flag{{Name: "watch", Kind: Bool}, {Name: "env", Kind: String, Complete: envs}},
			Run:   func(Call) int { return OK },
		},
		{Name: "deploy", Run: func(Call) int { return OK }},
		{Name: "internal", Hidden: true, Run: func(Call) int { return OK }},
	}}
}

func TestCompletingASubcommand(t *testing.T) {
	got := Complete(completing(), []string{""})
	if !slices.Contains(got, "status") || !slices.Contains(got, "deploy") {
		t.Errorf("got %v", got)
	}
	if slices.Contains(got, "internal") {
		t.Errorf("a hidden command is offered: %v", got)
	}

	if got := Complete(completing(), []string{"de"}); !slices.Contains(got, "deploy") || slices.Contains(got, "status") {
		t.Errorf("got %v", got)
	}
}

// The same completer the CLI uses, so bash and the interface cannot offer
// different answers.
func TestCompletingAPositionalArgument(t *testing.T) {
	got := Complete(completing(), []string{"status", "api_"})
	if !slices.Contains(got, "api_gateway") || !slices.Contains(got, "api_worker") {
		t.Errorf("got %v", got)
	}
}

// After a flag that takes a value, the value is what comes next — not another
// flag and not a subcommand.
func TestCompletingAFlagsValue(t *testing.T) {
	got := Complete(completing(), []string{"status", "--env", "prod"})
	if !slices.Contains(got, "production") || !slices.Contains(got, "prod-eu") {
		t.Errorf("got %v", got)
	}
	if slices.Contains(got, "api_gateway") {
		t.Errorf("a flag's value was completed with an argument's values: %v", got)
	}
}

// A boolean flag takes no value, so what follows is the next argument.
func TestABooleanFlagDoesNotSwallowTheNextWord(t *testing.T) {
	got := Complete(completing(), []string{"status", "--watch", "api_"})
	if !slices.Contains(got, "api_gateway") {
		t.Errorf("got %v", got)
	}
}

func TestCompletingAFlagName(t *testing.T) {
	got := Complete(completing(), []string{"status", "--w"})
	if !slices.Contains(got, "--watch") {
		t.Errorf("got %v", got)
	}
	if got := Complete(completing(), []string{"status", "--"}); !slices.Contains(got, "--json") {
		t.Errorf("--json is not offered: %v", got)
	}
}

// The script is deliberately ignorant: it calls the tool back rather than
// listing anything. A generated script that enumerated the flags would be a
// fourth copy of the declaration, and the one nobody regenerates.
func TestTheCompletionScriptKnowsNothingAboutTheCommands(t *testing.T) {
	for _, shell := range []string{"bash", "zsh", "fish"} {
		var b bytes.Buffer
		if err := CompletionScript(shell, "democtl", &b); err != nil {
			t.Fatalf("%s: %v", shell, err)
		}
		script := b.String()
		if !strings.Contains(script, "__complete") {
			t.Errorf("%s does not call back into the tool:\n%s", shell, script)
		}
		for _, leaked := range []string{"status", "deploy", "--watch"} {
			if strings.Contains(script, leaked) {
				t.Errorf("%s hard-codes %q, so it goes stale:\n%s", shell, leaked, script)
			}
		}
	}
	var b bytes.Buffer
	if err := CompletionScript("csh", "democtl", &b); err == nil {
		t.Error("an unsupported shell was accepted")
	}
}
