package app

import "testing"

// The zero value works. A tool that declares a Toggles and uses it without
// building one first is the common case, and a nil map panics on write.
func TestTheZeroValueWorks(t *testing.T) {
	var tg Toggles
	if tg.On("rg-forge") {
		t.Error("a fresh set has something on")
	}
	tg.Toggle("rg-forge")
	if !tg.On("rg-forge") {
		t.Error("toggling a key in a zero-value set did not turn it on")
	}
}

// Keyed by a stable id, not by position. The rows are rebuilt by a filter, a
// re-sort, a pivot or a refresh, and an index survives none of them.
func TestAKeyIsRememberedAcrossARebuild(t *testing.T) {
	var tg Toggles
	tg.Toggle("rg-data")

	// Whatever the rows are now, the answer is about the thing, not the row.
	if !tg.On("rg-data") || tg.On("rg-forge") {
		t.Error("the set answered about a position rather than a thing")
	}
}

// the deploy tool's rule, and its reason: without the clear, "hide everything"
// leaves whatever was individually revealed still showing.
func TestToggleAllTwiceIsAReliableWayBackToNothing(t *testing.T) {
	var tg Toggles
	tg.Toggle("one")
	tg.Toggle("two")

	tg.ToggleAll()
	if !tg.On("three") {
		t.Error("all did not cover a key nothing had named")
	}
	tg.ToggleAll()
	for _, k := range []string{"one", "two", "three"} {
		if tg.On(k) {
			t.Errorf("%s survived the way back to nothing", k)
		}
	}
}

// The global override covers keys nothing has named yet, which is the whole
// difference between it and turning on the forty things you happen to know
// about.
func TestAllCoversWhatHasNotBeenSeenYet(t *testing.T) {
	var tg Toggles
	tg.SetAll(true)
	if !tg.On("a-bucket-that-arrives-later") {
		t.Error("all did not cover a key added after it")
	}
}

// Turning one thing off inside "all" has to actually do something. Leaving the
// override set would make the write appear to be ignored, which is worse than
// refusing it.
func TestTurningOneOffInsideAllLeavesTheRestOn(t *testing.T) {
	var tg Toggles
	tg.Set("one", true)
	tg.Set("two", true)
	tg.SetAll(true)
	tg.Set("one", true) // named, so it survives the way down
	tg.Set("two", true)

	tg.Set("two", false)
	if tg.On("two") {
		t.Error("turning a key off inside all did nothing")
	}
	if !tg.On("one") {
		t.Error("turning one key off turned the others off too")
	}
	if tg.All() {
		t.Error("the global override is still set after an exception")
	}
}

// Any is what a footer asks before it offers the key that undoes this.
func TestAnyIsTrueForBothKindsOfOn(t *testing.T) {
	var tg Toggles
	if tg.Any() {
		t.Error("a fresh set reports something on")
	}
	tg.Toggle("one")
	if !tg.Any() {
		t.Error("a single key did not count")
	}
	tg.SetAll(false)
	tg.SetAll(true)
	if !tg.Any() {
		t.Error("the global override did not count")
	}
}
