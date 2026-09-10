package spec

// Exit codes, and the reason cobra is not here.
//
// The deploy tool's contract is richer than ok-or-not, and the difference
// carries meaning a script acts on:
//
//	0  it worked, or a dry run found nothing to do
//	1  it failed
//	2  a dry run found DRIFT — not an error, a finding
//
// The two matters. `the deploy tool diff && deploy` should not deploy when
// there is nothing to deploy, and `if the deploy tool diff; then ...` should
// be able to tell "no drift" from "could not look". A tool that reports drift
// as failure makes every wrapper script wrong in the same direction.
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

// codes is the contract above, in the form the manifest carries.
//
// Here rather than in describe.go so the numbers and their meanings are one
// thing: a code added to the consts and forgotten here would be a code no
// caller could learn about, and TestEveryExitCodeIsInTheManifest fails if the
// two lists disagree.
var codes = []CodeInfo{
	{Code: OK, Name: "ok", Meaning: "it worked, or a dry run found nothing to do"},
	{Code: Fail, Name: "fail", Meaning: "it failed: we could not do the thing, or could not look"},
	{Code: Drift, Name: "drift", Meaning: "a dry run found something — a finding, not a failure"},
}

// Codes returns the exit-code contract.
//
// A function returning a fresh slice rather than an exported var, because an
// exported slice is writable by anyone who can read it: one caller sorting it
// in place, or rewriting a Meaning, changes what every other caller in the
// process is told. theme.DefaultGlyphs learned the same lesson — see the test
// that fails if With mutates it.
func Codes() []CodeInfo { return append([]CodeInfo(nil), codes...) }
