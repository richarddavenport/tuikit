// Package app is the shell a tuikit tool runs inside: key routing, mouse
// dispatch, a screen router, and the async conventions as TYPES.
//
// The argument for types over prose is in this repo's own history. Every one
// of these conventions was written down in a design doc and then
// re-implemented, slightly differently, in each of four tools — and democtl,
// which documents capturesKeys as the central convention, still needed a
// reviewer to notice where it applied. A convention you can forget to apply is
// not a convention. A type you have to construct is.
package app

// Gen invalidates results from work that has been abandoned.
//
// # Where this came from
//
// The deploy tool's sessGen — "a tick scheduled by a closed session must not
// drive the replacement session (quick reconnects would otherwise multiply
// refresh chains)" — and democtl's gen, for a step finishing after you pressed
// Esc.
//
// The failure it prevents is nasty precisely because it looks like nothing: an
// abandoned run's result arrives and draws into its successor, so the step
// list shows a step that belongs to a different run. Nothing errors. The
// screen is simply wrong, in a way that is almost impossible to reproduce.
type Gen struct{ n int }

// Next abandons everything in flight and returns the new generation, which is
// what any work started now should carry.
func (g *Gen) Next() int { g.n++; return g.n }

// Current is the generation work started now belongs to.
func (g *Gen) Current() int { return g.n }

// Stale reports that a result belongs to a generation that has been abandoned,
// and should be dropped rather than drawn.
func (g *Gen) Stale(gen int) bool { return gen != g.n }

// Poll is a repeating background read: single-flight, and it gives up.
//
// # Where this came from
//
// The deploy tool's refreshing flag — "a slow fetch (>refreshEvery, e.g. over
// a laggy ssh link) must not stack concurrent Docker API calls" — and its
// repairing flag, so "a second failure reports instead of looping".
//
// Both halves matter and both are easy to leave out. Without single-flight, a
// link slower than the poll interval accumulates calls until something falls
// over, and it only happens to people on bad connections. Without a bound, a
// tool whose backend is down retries forever, quietly, and the only symptom is
// that the machine is warm.
type Poll struct {
	// Limit is the consecutive failures tolerated before it stops. Zero takes
	// the deploy tool's five, which is long enough to ride out a blip and short
	// enough that nobody is left polling into the void.
	Limit int

	inFlight bool
	failures int
}

const pollLimit = 5

// Start reports whether to begin a read. False means one is already in flight,
// or it has given up.
func (p *Poll) Start() bool {
	if p.inFlight || p.Stopped() {
		return false
	}
	p.inFlight = true
	return true
}

// Done records a read that worked, which also forgives earlier failures: an
// interface that stops polling because of five failures spread over an hour is
// reporting a problem that went away.
func (p *Poll) Done() {
	p.inFlight, p.failures = false, 0
}

// Failed records a read that did not.
func (p *Poll) Failed() {
	p.inFlight = false
	p.failures++
}

// Stopped reports that it has given up, which a tool should SAY rather than
// merely do — a poller that has quietly stopped looks exactly like one whose
// data has not changed.
func (p *Poll) Stopped() bool {
	limit := p.Limit
	if limit == 0 {
		limit = pollLimit
	}
	return p.failures >= limit
}

// Failures is how many consecutive reads have failed, for saying so.
func (p *Poll) Failures() int { return p.failures }

// Resume clears the failures and starts again, for when the reader asks.
func (p *Poll) Resume() { p.inFlight, p.failures = false, 0 }
