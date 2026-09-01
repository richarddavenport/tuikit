// Package fleet is democtl's engine: a small fictional service fleet, and the
// operations you can run against it.
//
// It has no UI imports, which is the split every tuikit tool keeps. It also has
// no backend — everything here is generated from a seed, deterministically, and
// that is the point rather than a shortcut. tuikit's capture harness renders
// democtl's screens as its own test fixture, and a fixture that changes between
// runs is not one. Given the same seed this package returns the same fleet, the
// same log lines and the same step timings, forever.
package fleet

import (
	"fmt"
	"math/rand"
	"time"
)

// Epoch is the clock every generated timestamp is relative to.
//
// Fixed rather than time.Now, for the same reason the data is seeded: a frame
// showing "updated 4m ago" has to say that tomorrow too, or every captured
// golden fails the day after it is written.
var Epoch = time.Date(2026, 8, 31, 9, 14, 3, 0, time.UTC)

// State is what a service is doing. The four are chosen to exercise the palette:
// each maps to a colour role, and together they cover every one a status can
// take.
type State int

// The four states, each mapping to a colour role, together covering every one a
// status line can take.
const (
	Running State = iota
	Pending
	Degraded
	Failed
)

func (s State) String() string {
	switch s {
	case Running:
		return "running"
	case Pending:
		return "pending"
	case Degraded:
		return "degraded"
	case Failed:
		return "failed"
	}
	return "unknown"
}

// Service is one deployed thing.
type Service struct {
	Name    string
	Stack   string
	Image   string
	Node    string
	Ready   int
	Want    int
	State   State
	Updated time.Time
	// Note is the one-line explanation a Degraded or Failed service owes. Empty
	// on a healthy one, because a healthy service has nothing to explain.
	Note string
}

// Fleet is the whole estate at one moment.
type Fleet struct {
	Services []Service
}

// New builds the same fleet every time for a given seed.
func New(seed int64) Fleet {
	r := rand.New(rand.NewSource(seed))
	out := make([]Service, 0, len(catalog))
	for _, s := range catalog {
		svc := Service{
			Name:    s.name,
			Stack:   s.stack,
			Image:   fmt.Sprintf("registry.example/%s:%s", s.name, tags[r.Intn(len(tags))]),
			Node:    nodes[r.Intn(len(nodes))],
			Want:    s.want,
			Ready:   s.want,
			State:   Running,
			Updated: Epoch.Add(-time.Duration(r.Intn(600)) * time.Second),
		}
		switch s.trouble {
		case troubleDegraded:
			svc.Ready, svc.State = s.want-1, Degraded
			svc.Note = "one task restarting: liveness probe timed out"
		case troubleFailed:
			svc.Ready, svc.State = 0, Failed
			svc.Note = "no suitable node: memory reservation exceeds every node"
		case troublePending:
			svc.Ready, svc.State = 0, Pending
			svc.Note = "waiting on image pull"
		}
		out = append(out, svc)
	}
	return Fleet{Services: out}
}

// Stacks returns the stack names in the order they should be read, which is the
// order they were declared rather than alphabetical — a fleet has a shape, and
// sorting it hides that.
func (f Fleet) Stacks() []string {
	seen := map[string]bool{}
	var out []string
	for _, s := range f.Services {
		if !seen[s.Stack] {
			seen[s.Stack] = true
			out = append(out, s.Stack)
		}
	}
	return out
}

// Healthy reports whether everything is where it should be — what a header
// wants to say in one word.
func (f Fleet) Healthy() bool {
	for _, s := range f.Services {
		if s.State != Running {
			return false
		}
	}
	return true
}

// LogLine is one line of a service's output.
type LogLine struct {
	At time.Time
	// Stderr is its own field rather than a level, because that is the
	// distinction a terminal actually has. Most programs write ordinary
	// progress to stderr, so this is not a synonym for "error".
	Stderr bool
	Text   string
}

// Logs returns n deterministic lines for a service, oldest first.
func (f Fleet) Logs(service string, n int) []LogLine {
	r := rand.New(rand.NewSource(int64(len(service)) * 7919))
	out := make([]LogLine, 0, n)
	for i := range n {
		tpl := logTemplates[r.Intn(len(logTemplates))]
		out = append(out, LogLine{
			At:     Epoch.Add(-time.Duration(n-i) * 3 * time.Second),
			Stderr: tpl.stderr,
			Text:   fmt.Sprintf(tpl.text, service, 100+r.Intn(800)),
		})
	}
	return out
}

// Step is one unit of a Plan: the shape azctl's playbooks and pgctl's applies
// both have, reduced to what a step list has to draw.
type Step struct {
	Name string
	// Took is how long it takes when it runs. Fixed per step so a captured run
	// reports the same durations every time.
	Took time.Duration
	// Skip marks a step whose check already passes — the idempotency case, and
	// the one a step list has to render differently or it lies about the work.
	Skip bool
	// Fails marks the step a Plan is built to break on, when it is built to
	// break. A step list that has only ever been drawn succeeding is a step list
	// with an untested failure state.
	Fails bool
}

// Plan is a named sequence of steps.
type Plan struct {
	Name  string
	Steps []Step
}

// Deploy is the plan democtl runs against a service. broken asks for the
// variant that fails partway, which is the interesting one to look at.
func Deploy(service string, broken bool) Plan {
	steps := []Step{
		{Name: "resolve image digest", Took: 400 * time.Millisecond},
		{Name: "check node capacity", Took: 260 * time.Millisecond, Skip: true},
		{Name: "push service spec", Took: 1200 * time.Millisecond},
		{Name: "wait for tasks to converge", Took: 4300 * time.Millisecond},
		{Name: "verify health endpoint", Took: 900 * time.Millisecond},
	}
	if broken {
		steps[3].Fails = true
	}
	return Plan{Name: "deploy " + service, Steps: steps}
}

// --- the fixture itself -------------------------------------------------

type trouble int

const (
	troubleNone trouble = iota
	troubleDegraded
	troubleFailed
	troublePending
)

var catalog = []struct {
	name    string
	stack   string
	want    int
	trouble trouble
}{
	// The first seven are the fleet democtl was written against, kept first and
	// in order so the screens that name a row by position still mean what they
	// meant.
	{"api_gateway", "api", 3, troubleNone},
	{"api_worker", "api", 2, troubleDegraded},
	{"api_migrate", "api", 1, troubleFailed},
	{"web_frontend", "web", 4, troubleNone},
	{"web_assets", "web", 1, troubleNone},
	{"data_indexer", "data", 2, troublePending},
	{"data_archiver", "data", 1, troubleNone},

	// The rest exist so the list OVERFLOWS. A demo whose list fits the pane it
	// is demonstrated in demonstrates nothing: a wheel that correctly does
	// nothing looks exactly like a wheel that is broken, and the viewport, the
	// count and the marker saying where an off-screen selection went are all
	// invisible. Trouble is sprinkled deep on purpose, so scrolling finds
	// something rather than more of the same.
	{"api_scheduler", "api", 2, troubleNone},
	{"api_webhooks", "api", 1, troubleNone},
	{"api_ratelimit", "api", 2, troubleNone},
	{"web_admin", "web", 2, troubleNone},
	{"web_docs", "web", 1, troubleNone},
	{"web_preview", "web", 1, troublePending},
	{"data_etl", "data", 3, troubleNone},
	{"data_warehouse", "data", 2, troubleNone},
	{"data_replica", "data", 2, troubleDegraded},
	{"data_backup", "data", 1, troubleNone},
	{"edge_router", "edge", 3, troubleNone},
	{"edge_cache", "edge", 4, troubleNone},
	{"edge_tls", "edge", 2, troubleNone},
	{"edge_waf", "edge", 2, troubleDegraded},
	{"auth_session", "auth", 3, troubleNone},
	{"auth_tokens", "auth", 2, troubleNone},
	{"auth_directory", "auth", 1, troubleNone},
	{"auth_mfa", "auth", 1, troublePending},
	{"media_upload", "media", 2, troubleNone},
	// Fifteen characters, one over what the list column shows. Deliberate: it
	// is the fixture's proof that a name too long is truncated with an
	// ellipsis rather than pushing the columns beside it out of line.
	{"media_transcode", "media", 4, troubleFailed},
	{"media_thumbs", "media", 2, troubleNone},
	{"media_cdn", "media", 3, troubleNone},
	{"ops_metrics", "ops", 2, troubleNone},
	{"ops_logs", "ops", 3, troubleNone},
	{"ops_alerts", "ops", 1, troubleNone},
	{"ops_tracing", "ops", 2, troubleNone},
	{"ops_registry", "ops", 1, troubleDegraded},
	{"ops_backup", "ops", 1, troubleNone},
	{"search_ingest", "search", 2, troubleNone},
	{"search_query", "search", 3, troubleNone},
	{"search_suggest", "search", 1, troubleNone},
	{"search_reindex", "search", 1, troublePending},
}

var (
	nodes = []string{"node-a", "node-b", "node-c"}
	tags  = []string{"v2.4.1", "v2.4.0", "v2.3.9", "sha-4f1a2c"}
)

var logTemplates = []struct {
	text   string
	stderr bool
}{
	{"%s: listening on :%d", false},
	{"%s: handled request in %dms", false},
	{"%s: cache warm, %d entries", false},
	{"%s: retrying upstream after %dms", true},
	{"%s: connection reset by peer (attempt %d)", true},
	{"%s: flushed %d records", false},
}
