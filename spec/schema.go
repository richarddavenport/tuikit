package spec

import "sort"

// Schema is a JSON Schema for one command's arguments and flags.
//
// # Why this is here and an MCP server is not
//
// The board needed to be usable by an agent that cannot run a shell — a hosted
// runtime, an Agents SDK app — and wrote an MCP server for it. About two
// hundred lines, and by its author's own account "not one of them is about the
// board": every tool it offers is generated from the [Command] tree and every
// call goes back through Run (issue 54).
//
// That is a true observation about the tree, and it is also an argument for a
// smaller thing than an MCP server. The part that is tuikit's business is the
// TRANSFORMATION — a command already declares its name, its help, its
// arguments and their kinds, and turning that into a parameter schema is a
// fact about [Command] rather than about any protocol. The JSON-RPC plumbing,
// the tool naming, the allow-list of what to expose, the transport: all of
// that is the server's, and it differs per host.
//
// So this is the half that would otherwise be written once per transport.
// An MCP server built on it is small; so is a function-calling adapter, or
// anything else that needs to tell a model what a command takes.
//
// tuikit does not gain a network protocol for one consumer.
type Schema struct {
	Type       string              `json:"type"`
	Properties map[string]Property `json:"properties"`
	Required   []string            `json:"required,omitempty"`
}

// Property is one parameter.
type Property struct {
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	// Default is the flag's default as declared, omitted when there is none.
	Default string `json:"default,omitempty"`
}

// SchemaOf builds the schema for a command.
//
// Arguments and flags share one property map, because a caller that is not a
// shell has no positional/named distinction to honor — it has a bag of named
// values. Required covers the arguments a command cannot run without; a flag
// is never required, which is what makes it a flag.
//
// A variadic argument becomes an array. Everything else follows [Kind], and an
// unknown kind is a string rather than an error: a schema that refuses to
// build leaves an agent with no description at all, which is worse than a
// loose one.
func SchemaOf(cmd Command) Schema {
	s := Schema{Type: "object", Properties: map[string]Property{}}

	for _, arg := range cmd.Args {
		p := Property{Type: "string", Description: arg.Help}
		if arg.Variadic {
			p.Type = "array"
		}
		s.Properties[arg.Name] = p
		if arg.Required {
			s.Required = append(s.Required, arg.Name)
		}
	}

	for _, flag := range cmd.Flags {
		s.Properties[flag.Name] = Property{
			Type:        jsonType(flag.Kind),
			Description: flag.Help,
			Default:     flag.Default,
		}
	}

	// Sorted, because a schema that reorders itself between runs is one nobody
	// can diff — and an agent's cache key is usually the serialized schema.
	sort.Strings(s.Required)
	return s
}

// jsonType maps a flag's kind onto JSON Schema's.
//
// Duration is a string rather than a number: it is declared as "30s" and
// parsed as one, and a schema saying "integer" would invite an agent to send
// 30 and mean half a minute.
func jsonType(k Kind) string {
	switch k {
	case Bool:
		return "boolean"
	case Int:
		return "integer"
	case String, Duration:
		return "string"
	}
	return "string"
}
