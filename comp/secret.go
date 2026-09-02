package comp

import "strings"

// maskWidth is how many bullets a masked value is, whatever it holds.
//
// FIXED, and that is the entire point. Masking to the value's own length tells
// a reader the password is six characters, which is most of what a guesser
// wants to know and all of what a shoulder-surfer needs. Eight is enough to
// read as "there is something here" and short enough not to look like a value.
const maskWidth = 8

// Mask is what a secret looks like when it is not being shown.
//
// It takes no argument, which is deliberate: a function given the secret is a
// function that could return some of it, and every "just the last four" is a
// decision about a particular kind of secret that a framework has no business
// making. A tool that wants `ghp_…a1b2` builds that string itself and passes it
// as an ordinary value — then it is visibly the tool's choice.
//
// # Where the rest of this lives
//
// Three separate things, and only this one is drawing:
//
//   - WHICH secrets are revealed is [app.Toggles], keyed by a stable id so the
//     state survives a filter, a re-sort or a refresh. That already exists,
//     and swarmctl's revealState is what it was lifted from.
//   - SCRUBBING a secret out of text that flows past — a command line, a log
//     line — belongs to the tool's engine, because only the tool knows what
//     its secrets are and where they go. azctl's runner does it, and it is
//     what caught a swarm join token printing to stdout.
//   - What a hidden value LOOKS like is this.
func Mask() string { return strings.Repeat("•", maskWidth) }
