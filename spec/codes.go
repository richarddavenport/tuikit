package spec

// Exit codes, and the reason cobra is not here.
//
// swarmctl's contract is richer than ok-or-not, and the difference carries
// meaning a script acts on:
//
//	0  it worked, or a dry run found nothing to do
//	1  it failed
//	2  a dry run found DRIFT — not an error, a finding
//
// The two matters. `swarmctl diff && deploy` should not deploy when there is
// nothing to deploy, and `if swarmctl diff; then ...` should be able to tell
// "no drift" from "could not look". A tool that reports drift as failure makes
// every wrapper script wrong in the same direction.
//
// And a command that runs somebody else's program passes ITS status through
// unchanged, because a wrapper that flattens exit codes is a wrapper that
// breaks the thing it wraps.
const (
	// OK is success, and a dry run that found nothing to do.
	OK = 0
	// Fail is an error: we could not do the thing, or could not look.
	Fail = 1
	// Drift is a dry run that found something. A finding, not a failure.
	Drift = 2
)
