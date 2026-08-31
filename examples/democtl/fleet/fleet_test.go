package fleet

import "testing"

// The harness renders democtl as its own fixture, and a fixture that changes
// between runs is not one. Everything here has to be a pure function of the
// seed — no clock, no map iteration order reaching the output, no rand without
// a source.
func TestTheSameSeedIsTheSameFleet(t *testing.T) {
	a, b := New(1), New(1)
	if len(a.Services) != len(b.Services) {
		t.Fatalf("different lengths: %d and %d", len(a.Services), len(b.Services))
	}
	for i := range a.Services {
		if a.Services[i] != b.Services[i] {
			t.Errorf("service %d differs:\n %+v\n %+v", i, a.Services[i], b.Services[i])
		}
	}
}

func TestADifferentSeedIsADifferentFleet(t *testing.T) {
	if New(1).Services[0] == New(2).Services[0] {
		t.Error("the seed is not reaching the generated data")
	}
}

func TestLogsAreStableAndOrdered(t *testing.T) {
	first, second := New(1).Logs("api_gateway", 12), New(9).Logs("api_gateway", 12)
	if len(first) != 12 {
		t.Fatalf("got %d lines, want 12", len(first))
	}
	for i := range first {
		// Logs depend on the service, not the fleet's seed: two runs looking at
		// the same service must show the same output.
		if first[i] != second[i] {
			t.Errorf("line %d differs between fleets", i)
		}
		if i > 0 && !first[i].At.After(first[i-1].At) {
			t.Errorf("line %d is not after line %d", i, i-1)
		}
	}
}

// Every state has to be reachable, or a screen showing one of them cannot be
// captured and the palette role it uses is never seen in use.
func TestEveryStateIsPresentInTheFleet(t *testing.T) {
	seen := map[State]bool{}
	for _, s := range New(1).Services {
		seen[s.State] = true
	}
	for _, want := range []State{Running, Pending, Degraded, Failed} {
		if !seen[want] {
			t.Errorf("no service is %s — that state can never be looked at", want)
		}
	}
}

// A service that is not running owes an explanation; a healthy one has nothing
// to explain and should not invent something.
func TestOnlyUnhealthyServicesCarryANote(t *testing.T) {
	for _, s := range New(1).Services {
		if s.State == Running && s.Note != "" {
			t.Errorf("%s is running but explains itself: %q", s.Name, s.Note)
		}
		if s.State != Running && s.Note == "" {
			t.Errorf("%s is %s with no explanation", s.Name, s.State)
		}
	}
}

// The broken plan has to actually break, and the healthy one has to not.
func TestDeployBreaksOnlyWhenAsked(t *testing.T) {
	for _, step := range Deploy("api_gateway", false).Steps {
		if step.Fails {
			t.Errorf("%q fails in the healthy plan", step.Name)
		}
	}
	var fails int
	for _, step := range Deploy("api_migrate", true).Steps {
		if step.Fails {
			fails++
		}
	}
	if fails != 1 {
		t.Errorf("the broken plan has %d failing steps, want 1", fails)
	}
}
