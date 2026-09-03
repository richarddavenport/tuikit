package harness_test

import (
	"fmt"
	"testing"

	tea "github.com/charmbracelet/bubbletea"

	"github.com/richarddavenport/tuikit/app"
	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/harness"
)

// hintModel offers three keys and handles two of them.
type hintModel struct{ n, tab int }

func (hintModel) Init() tea.Cmd { return nil }

func (m hintModel) Update(msg tea.Msg) (app.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok {
		switch k.String() {
		case "j":
			m.n++
		case "tab":
			m.tab++
			// "x" is advertised in the footer below and handled nowhere.
		}
	}
	return m, nil
}

func (m hintModel) Draw(c *comp.Canvas, r comp.Rect) {
	c.Text(r.X, r.Y, "n="+itoa(m.n)+" tab="+itoa(m.tab), nil, comp.Region("body"))
}

func itoa(n int) string {
	if n == 0 {
		return "0"
	}
	var b []byte
	for ; n > 0; n /= 10 {
		b = append([]byte{byte('0' + n%10)}, b...)
	}
	return string(b)
}

func build() harness.Driver {
	r := app.New(hintModel{})
	r.Update(tea.WindowSizeMsg{Width: 40, Height: 6})
	return r
}

// The dead key is reported and the live ones are not.
func TestHintsFindsTheKeyThatDoesNothing(t *testing.T) {
	rec := &recorder{}

	harness.Hints(rec, build,
		comp.Hint{Key: "j", Label: "move"},
		comp.Hint{Key: "tab", Label: "next pane"},
		comp.Hint{Key: "x", Label: "explode"},
	)

	if len(rec.msgs) != 1 {
		t.Fatalf("reported %d hints, want 1: %v", len(rec.msgs), rec.msgs)
	}
	if got := rec.msgs[0]; !contains(got, `"x"`) || !contains(got, "explode") {
		t.Errorf("the message does not name the dead key and its label: %s", got)
	}
}

// A guard that never fires is indistinguishable from a clean package, so the
// live-only case must pass.
func TestHintsPassesWhenEveryKeyDoesSomething(t *testing.T) {
	rec := &recorder{}
	harness.Hints(rec, build, comp.Hint{Key: "j", Label: "move"}, comp.Hint{Key: "tab", Label: "next pane"})

	if len(rec.msgs) != 0 {
		t.Errorf("reported a live key: %v", rec.msgs)
	}
}

// A hint naming several keys is several keys, or a tool that documented a pair
// is told only about the first one.
func TestHintsChecksEveryKeyAHintNames(t *testing.T) {
	rec := &recorder{}
	harness.Hints(rec, build, comp.Hint{Key: "j/x", Label: "move"})

	if len(rec.msgs) != 1 {
		t.Fatalf("a pair reported %d, want 1 (x is dead, j is not): %v", len(rec.msgs), rec.msgs)
	}
}

type recorder struct{ msgs []string }

func (r *recorder) Helper() {}
func (r *recorder) Errorf(format string, args ...any) {
	r.msgs = append(r.msgs, sprintf(format, args...))
}
func (r *recorder) Fatalf(format string, args ...any) {
	r.msgs = append(r.msgs, sprintf(format, args...))
}

func sprintf(format string, args ...any) string { return fmt.Sprintf(format, args...) }

func contains(s, sub string) bool {
	for i := 0; i+len(sub) <= len(s); i++ {
		if s[i:i+len(sub)] == sub {
			return true
		}
	}
	return false
}
