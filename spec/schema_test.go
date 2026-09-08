package spec_test

import (
	"encoding/json"
	"testing"

	"github.com/richarddavenport/tuikit/spec"
)

func cmd() spec.Command {
	return spec.Command{
		Name:  "restart",
		Short: "restart a service, rolling",
		Args: []spec.Arg{
			{Name: "service", Help: "which service", Required: true},
			{Name: "extra", Help: "any more of them", Variadic: true},
		},
		Flags: []spec.Flag{
			{Name: "force", Kind: spec.Bool, Help: "do not ask"},
			{Name: "replicas", Kind: spec.Int, Help: "how many", Default: "2"},
			{Name: "timeout", Kind: spec.Duration, Help: "how long to wait", Default: "30s"},
			{Name: "note", Kind: spec.String, Help: "why"},
		},
	}
}

// Arguments and flags land in one property map: a caller that is not a shell
// has a bag of named values, not positions.
func TestSchemaCarriesArgsAndFlags(t *testing.T) {
	s := spec.SchemaOf(cmd())

	for _, name := range []string{"service", "extra", "force", "replicas", "timeout", "note"} {
		if _, ok := s.Properties[name]; !ok {
			t.Errorf("%q is not in the schema", name)
		}
	}
	if s.Type != "object" {
		t.Errorf("the schema is %q", s.Type)
	}
}

// Kinds map onto JSON Schema's, and a duration is a string.
func TestKindsBecomeJSONTypes(t *testing.T) {
	s := spec.SchemaOf(cmd())

	for name, want := range map[string]string{
		"force":    "boolean",
		"replicas": "integer",
		"note":     "string",
		"service":  "string",
		"extra":    "array",
		// "30s" is declared as a string and parsed as one. "integer" would
		// invite an agent to send 30 and mean half a minute.
		"timeout": "string",
	} {
		if got := s.Properties[name].Type; got != want {
			t.Errorf("%s is %q, want %q", name, got, want)
		}
	}
}

// Only a required ARGUMENT is required. A flag that must be given is not a
// flag, and a schema saying otherwise makes an agent ask for something the CLI
// would have defaulted.
func TestOnlyRequiredArgumentsAreRequired(t *testing.T) {
	s := spec.SchemaOf(cmd())

	if len(s.Required) != 1 || s.Required[0] != "service" {
		t.Errorf("required is %v, want [service]", s.Required)
	}
}

// The help text comes through, because it is the only thing telling a model
// what a parameter means.
func TestHelpBecomesTheDescription(t *testing.T) {
	s := spec.SchemaOf(cmd())
	if got := s.Properties["service"].Description; got != "which service" {
		t.Errorf("description is %q", got)
	}
	if got := s.Properties["replicas"].Default; got != "2" {
		t.Errorf("default is %q", got)
	}
}

// It serializes to the shape a tool-calling API expects, and stably: an agent's
// cache key is usually the serialized schema.
func TestTheSchemaSerializesStably(t *testing.T) {
	first, err := json.Marshal(spec.SchemaOf(cmd()))
	if err != nil {
		t.Fatal(err)
	}
	for i := 0; i < 5; i++ {
		again, err := json.Marshal(spec.SchemaOf(cmd()))
		if err != nil {
			t.Fatal(err)
		}
		if string(again) != string(first) {
			t.Fatalf("run %d differs:\n%s\n%s", i, first, again)
		}
	}
	if !json.Valid(first) {
		t.Error("not valid JSON")
	}
}

// A command with nothing to take still produces a usable object rather than
// null, which a strict tool-calling host rejects.
func TestACommandWithNoParametersIsStillAnObject(t *testing.T) {
	s := spec.SchemaOf(spec.Command{Name: "status"})

	body, err := json.Marshal(s)
	if err != nil {
		t.Fatal(err)
	}
	if string(body) != `{"type":"object","properties":{}}` {
		t.Errorf("empty schema serializes as %s", body)
	}
}
