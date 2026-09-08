package guard

import (
	"strings"
	"testing"
)

// A clean engine is stdlib plus its own domain SDK. That is not an ideal — it
// is what democtl's fleet and swarmctl's internal/engine actually import.
func TestEngineAcceptsStdlibAndADomainSDK(t *testing.T) {
	dir := pkg(t, `package fleet

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/jackc/pgx/v5"
)

var _ = fmt.Sprint
var _ = rand.Int
var _ = time.Now
var _ pgx.Tx
`)
	silent(t, run(t, func(rec T) { Engine(rec, dir) }))
}

// The bar is "no terminal concepts", not "no Bubble Tea": every one of these is
// a way of learning about columns or escape sequences.
func TestEngineRejectsEveryWayOfLearningAboutTerminals(t *testing.T) {
	for _, tc := range []struct {
		path string
		says string
	}{
		{"github.com/charmbracelet/bubbletea", "does not draw"},
		{"github.com/charmbracelet/lipgloss", "does not draw"},
		{"github.com/gdamore/tcell/v2", "does not draw"},
		{"github.com/muesli/termenv", "no terminal"},
		{"github.com/muesli/reflow/wordwrap", "no columns"},
		{"github.com/mattn/go-runewidth", "no columns"},
		{"github.com/mattn/go-isatty", "must not care"},
		{"golang.org/x/term", "no terminal"},
		{"github.com/richarddavenport/tuikit/theme", "gets nothing from the framework"},
	} {
		t.Run(tc.path, func(t *testing.T) {
			dir := pkg(t, "package engine\n\nimport _ \""+tc.path+"\"\n")
			msgs := run(t, func(rec T) { Engine(rec, dir) })
			want(t, msgs, tc.path)
			want(t, msgs, tc.says)
		})
	}
}

// The reason is reported, not just the fact. A guard that says "forbidden" tells
// you what happened; one that says why tells you whether the rule is wrong,
// which is occasionally the right conclusion.
func TestEngineSaysWhyRatherThanJustNo(t *testing.T) {
	dir := pkg(t, "package engine\n\nimport _ \"github.com/charmbracelet/lipgloss\"\n")
	msgs := run(t, func(rec T) { Engine(rec, dir) })
	if len(msgs) != 1 {
		t.Fatalf("want one failure, got %v", msgs)
	}
	for _, part := range []string{"lipgloss", "the engine does not draw", "the terminal is the UI's business"} {
		if !strings.Contains(msgs[0], part) {
			t.Errorf("the message is missing %q:\n%s", part, msgs[0])
		}
	}
}

// Prefixes match on a path boundary, so a package that merely starts with the
// same letters is not condemned.
func TestEngineMatchesOnAPathBoundary(t *testing.T) {
	dir := pkg(t, "package engine\n\nimport _ \"github.com/muesli/ansible-inventory\"\n")
	silent(t, run(t, func(rec T) { Engine(rec, dir) }))
}

// A test file may import anything — it is not the engine.
func TestEngineIgnoresTestFiles(t *testing.T) {
	dir := pkg(t, "package engine\n\nimport _ \"fmt\"\n")
	write(t, dir, "engine_test.go", "package engine\n\nimport _ \"github.com/charmbracelet/lipgloss\"\n")
	silent(t, run(t, func(rec T) { Engine(rec, dir) }))
}

// Same failure mode as the other guards: pointed at nothing, fail rather than
// pass. A guard that scans nothing silently is how everyone comes to believe a
// rule is on when it is not.
func TestEngineWithNothingToScanFails(t *testing.T) {
	want(t, run(t, func(rec T) { Engine(rec, t.TempDir()) }), "no Go files")
}

// A tool may hold its engine to a different list.
func TestEngineTakesADifferentList(t *testing.T) {
	dir := pkg(t, "package engine\n\nimport _ \"github.com/charmbracelet/lipgloss\"\n")
	own := []Denied{{"example.com/forbidden", "the tool says so"}}

	silent(t, run(t, func(rec T) { Engine(rec, dir, own...) }))
}

// Every entry has to explain itself, or the message above cannot.
func TestEveryDeniedPrefixSaysWhy(t *testing.T) {
	for _, d := range TerminalPackages {
		if d.Prefix == "" || d.Why == "" {
			t.Errorf("%+v is incomplete", d)
		}
	}
}

// Issue 62: decision 1 says the fixture is a value the engine returns, and
// nothing enforced it, so in docket it quietly stopped being one.
func TestDerivedRejectsAHandBuiltAnswer(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "board_test.go", `package tui

import "example.com/tool/internal/engine"

func fixture() engine.Truth {
	return engine.Truth{Stale: true, PR: "merged 11d ago by you"}
}
`)
	got := run(t, func(rec T) { Derived(rec, dir, "engine", "Truth", "Live") })
	want(t, got, "builds engine.Truth by hand")
}

// The other direction: a fixture that asks the engine is what the rule wants,
// and must not fire.
func TestDerivedAcceptsAnAnswerFromTheEngine(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "board_test.go", `package tui

import "example.com/tool/internal/engine"

func fixture() engine.Truth {
	return engine.Reconcile(cannedInput())
}
`)
	silent(t, run(t, func(rec T) { Derived(rec, dir, "engine", "Truth", "Live") }))
}

// A type the test SHOULD construct is not listed and must not fire. This is not
// "tests may not build structs".
func TestDerivedIgnoresTypesItWasNotGiven(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "board_test.go", `package tui

import "example.com/tool/internal/engine"

var req = engine.Request{Branch: "main"}
`)
	silent(t, run(t, func(rec T) { Derived(rec, dir, "engine", "Truth") }))
}

// An empty literal is a zero value, not an invented world. Refusing it would
// push tests into a worse shape to satisfy a guard.
func TestDerivedAllowsTheZeroValue(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "board_test.go", `package tui

import "example.com/tool/internal/engine"

var nothingYet = engine.Truth{}
`)
	silent(t, run(t, func(rec T) { Derived(rec, dir, "engine", "Truth") }))
}

// A guard given nothing to check passes vacuously, which is the failure every
// guard here is written to avoid.
func TestDerivedWithNoTypesFails(t *testing.T) {
	dir := t.TempDir()
	write(t, dir, "board_test.go", "package tui\n")
	got := run(t, func(rec T) { Derived(rec, dir, "engine") })
	want(t, got, "checks nothing")
}
