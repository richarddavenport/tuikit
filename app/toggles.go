package app

// Toggles is a set of things that are on, keyed by a stable id.
//
// # Two tools, two purposes
//
// The cloud tool keeps `expanded map[string]bool` for which buckets in its
// tree are open. The deploy tool keeps this, under another name, for which
// secrets are unmasked:
//
//	// revealState tracks which rows are unmasked. Keyed by a stable row id so
//	// the state survives re-renders and cursor movement.
//	type revealState struct {
//		all  bool
//		rows map[string]bool
//	}
//
// The same shape, in two tools, for two unrelated things. That is the
// strongest form of the signal this package exists to act on: a convention
// arrived at twice independently is mechanism, not one tool's policy.
//
// # Keyed, not indexed
//
// By a stable id rather than a row number, because the rows are rebuilt — by a
// filter, a re-sort, a pivot, a refresh — and an index survives none of that.
// A set keyed by position quietly expands the wrong group after any of them,
// which reads as the interface having a mind of its own rather than as a bug.
//
// The zero value works.
type Toggles struct {
	// all is the global override. Kept separate from the per-key set rather
	// than written into it, because "everything, including what has not been
	// loaded yet" is not the same claim as "these forty things".
	all  bool
	keys map[string]bool
}

// On reports whether a key is on.
func (t *Toggles) On(key string) bool { return t.all || t.keys[key] }

// Toggle flips one key.
func (t *Toggles) Toggle(key string) { t.Set(key, !t.On(key)) }

// Set turns one key on or off.
//
// Setting a key while everything is on turns the global override OFF and
// switches every other key on in its place, so "all of them except this one"
// is reachable — and so the next read of On tells the truth. Leaving all set
// would make the write appear to do nothing at all.
func (t *Toggles) Set(key string, on bool) {
	if t.keys == nil {
		t.keys = map[string]bool{}
	}
	if t.all && !on {
		t.all = false
		for k := range t.keys {
			t.keys[k] = true
		}
		// Keys never named are lost here, which is why this only happens on
		// the way DOWN from all: a caller turning one thing off inside "all"
		// is telling us the set it cares about is the one it has told us
		// about.
	}
	t.keys[key] = on
}

// SetAll turns everything on or off, including keys nothing has named yet.
func (t *Toggles) SetAll(on bool) { t.all, t.keys = on, map[string]bool{} }

// ToggleAll flips the global override.
//
// Turning it off also clears the per-key set, so pressing the key twice is a
// reliable way back to nothing on. The deploy tool's note on why: without the
// clear, "hide everything" leaves whatever was individually revealed still
// showing, and the reader has no way to tell how many that is.
func (t *Toggles) ToggleAll() { t.SetAll(!t.all) }

// All reports that the global override is on. For a footer that has to say
// which key is the way back.
func (t *Toggles) All() bool { return t.all }

// Any reports that anything at all is on, which is the question a hint asks.
func (t *Toggles) Any() bool {
	if t.all {
		return true
	}
	for _, on := range t.keys {
		if on {
			return true
		}
	}
	return false
}
