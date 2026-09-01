package app

import "testing"

// A result from work that has been abandoned must not draw into its successor.
// The failure is nasty precisely because nothing errors: the screen is simply
// wrong, in a way that is almost impossible to reproduce.
func TestAnAbandonedGenerationsResultIsStale(t *testing.T) {
	var g Gen
	first := g.Current()

	if g.Stale(first) {
		t.Error("the current generation is stale")
	}

	second := g.Next()
	if !g.Stale(first) {
		t.Error("a result from the abandoned generation is not stale")
	}
	if g.Stale(second) {
		t.Error("a result from the new generation is stale")
	}
}

// Single-flight: a link slower than the poll interval must not stack calls.
func TestPollDoesNotStackReads(t *testing.T) {
	var p Poll
	if !p.Start() {
		t.Fatal("the first read did not start")
	}
	if p.Start() {
		t.Error("a second read started while the first was in flight")
	}
	p.Done()
	if !p.Start() {
		t.Error("no read started after the first finished")
	}
}

// A tool whose backend is down must not retry forever. The only symptom of that
// is a warm machine.
func TestPollGivesUp(t *testing.T) {
	var p Poll
	for range pollLimit {
		if !p.Start() {
			t.Fatalf("stopped after %d failures, before the limit", p.Failures())
		}
		p.Failed()
	}
	if !p.Stopped() {
		t.Errorf("still polling after %d failures", p.Failures())
	}
	if p.Start() {
		t.Error("a read started after it had given up")
	}
}

// Five failures spread over an hour is a problem that went away, not a reason
// to stop.
func TestASuccessForgivesEarlierFailures(t *testing.T) {
	var p Poll
	for range pollLimit - 1 {
		p.Start()
		p.Failed()
	}
	p.Start()
	p.Done()

	for range pollLimit - 1 {
		p.Start()
		p.Failed()
	}
	if p.Stopped() {
		t.Error("failures either side of a success added up")
	}
}

func TestResumeStartsAgain(t *testing.T) {
	p := Poll{Limit: 2}
	p.Start()
	p.Failed()
	p.Start()
	p.Failed()
	if !p.Stopped() {
		t.Fatal("did not stop at its limit")
	}

	p.Resume()
	if p.Stopped() || !p.Start() {
		t.Error("resuming did not start it again")
	}
}

// The limit is the caller's, because how long to keep trying depends on what is
// being read.
func TestTheLimitIsTheCallers(t *testing.T) {
	p := Poll{Limit: 1}
	p.Start()
	p.Failed()
	if !p.Stopped() {
		t.Error("a limit of one did not stop after one failure")
	}
}
