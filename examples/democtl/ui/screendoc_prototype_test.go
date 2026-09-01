package ui

import (
	"encoding/json"
	"fmt"
	"os"
	"strconv"
	"testing"

	"github.com/charmbracelet/lipgloss"
	"github.com/muesli/termenv"

	"github.com/richarddavenport/tuikit/examples/democtl/fleet"
	"github.com/richarddavenport/tuikit/prototype/htmlround"
	"github.com/richarddavenport/tuikit/prototype/screendoc"
)

// modelData is the engine side of the split, as narrow as the interface allows.
//
// Note what is NOT here: no colour, no width, no glyph, no style name. The
// model hands over "failed" and the document decides that failed is red and
// ends in ✗. That is the engine/UI split expressed as a data boundary rather
// than as a convention somebody has to remember.
type modelData struct{ m *Model }

func (d modelData) Bool(key string) bool {
	switch key {
	case "fleet.healthy":
		return d.m.fleet.Healthy()
	case "services.focused":
		return d.m.focus == paneList
	case "detail.focused":
		return d.m.focus == paneDetail
	}
	return false
}

func (d modelData) Index(key string) int {
	switch key {
	case "services.cursor":
		return d.m.cur
	case "detail.tab":
		return d.m.tab
	}
	return -1
}

func (d modelData) Value(key string) string {
	svc, ok := d.m.selected()
	switch key {
	case "clock":
		return d.m.now.Format("15:04:05")
	case "services.count":
		return strconv.Itoa(len(d.m.visible()))
	case "detail.tab":
		return strconv.Itoa(d.m.tab)
	case "ui.typing":
		return strconv.FormatBool(d.m.typing)
	case "ui.filtering":
		return strconv.FormatBool(d.m.filter != "" || d.m.typing)
	case "filter":
		return d.m.filter
	}
	if !ok {
		return ""
	}
	switch key {
	case "selected.name":
		return svc.Name
	case "selected.state":
		return svc.State.String()
	case "selected.replicas":
		return fmt.Sprintf("%d/%d", svc.Ready, svc.Want)
	case "selected.stack":
		return svc.Stack
	case "selected.node":
		return svc.Node
	case "selected.updated":
		return ago(d.m.now.Sub(svc.Updated))
	case "selected.note":
		return svc.Note
	case "selected.image":
		return svc.Image
	case "selected.want":
		return strconv.Itoa(svc.Want)
	}
	return ""
}

func (d modelData) Rows(key string) []map[string]string {
	if key != "services" {
		return nil
	}
	var out []map[string]string
	for _, s := range d.m.visible() {
		out = append(out, map[string]string{
			"name":  s.Name,
			"ready": strconv.Itoa(s.Ready),
			"want":  strconv.Itoa(s.Want),
			"state": s.State.String(),
		})
	}
	return out
}

// Can a screen be declared rather than drawn?
//
// The document in prototype/screendoc/testdata/dashboard.json describes the
// dashboard. This renders it against the same fleet, at the same size, and
// compares CELL FOR CELL with the hand-written View — so a colour that is
// merely close, or a row that is right but bold, fails.
//
// Against View rather than against the golden on purpose: the golden is
// stripped, so it would confirm the layout and say nothing about whether the
// document's roles resolve to the colours the screen actually has.
func TestPrototypeDashboardCanBeDeclared(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)

	raw, err := os.ReadFile("../../../prototype/screendoc/testdata/dashboard.json")
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}
	var doc screendoc.Doc
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parsing the document: %v", err)
	}

	for _, tc := range []struct {
		name string
		w, h int
		keys []string
	}{
		{"dashboard at 132x38", 132, 38, nil},
		{"dashboard at 80x24", 80, 24, nil},
		{"dashboard at 60x20", 60, 20, nil},
		{"a failed service selected", 132, 38, []string{"j", "j"}},
		{"the detail pane focused", 132, 38, []string{"tab"}},
		{"filtering", 132, 38, []string{"/", "w", "e"}},
		{"the config tab", 132, 38, []string{"tab", "right"}},
		{"a filter that matches nothing", 132, 38, []string{"/", "z", "z", "z"}},
	} {
		t.Run(tc.name, func(t *testing.T) {
			m := New(1)
			press(m, tc.keys...)
			m.SetSize(tc.w, tc.h)
			m.Now(fleet.Epoch)

			r, err := screendoc.New(doc, Palette, modelData{m})
			if err != nil {
				t.Fatalf("building the renderer: %v", err)
			}

			want := htmlround.FromANSI(m.View())
			got := htmlround.FromANSI(r.Render(doc, tc.w, tc.h))

			if d := htmlround.Diff(want, got); d != "" {
				t.Errorf("the declared screen is not the drawn one:\n%s", d)
			}
		})
	}
}

// The Events tab is the one the format cannot say, and it is worth naming
// rather than leaving out.
//
// Its rows are `muted(timestamp) + plain(text)` — two styles on one line — and
// a list row carries a single style. Written as a log rather than a failure
// because the gap is known: fixing it means a row of spans instead of a
// template string, which is a real change to the format and not a patch to
// this document.
func TestPrototypeEventsTabIsBeyondTheFormat(t *testing.T) {
	lipgloss.SetColorProfile(termenv.TrueColor)
	lipgloss.SetHasDarkBackground(true)

	raw, err := os.ReadFile("../../../prototype/screendoc/testdata/dashboard.json")
	if err != nil {
		t.Fatalf("reading the document: %v", err)
	}
	var doc screendoc.Doc
	if err := json.Unmarshal(raw, &doc); err != nil {
		t.Fatalf("parsing the document: %v", err)
	}

	m := New(1)
	press(m, "tab", "right", "right")
	m.SetSize(132, 38)
	m.Now(fleet.Epoch)

	r, err := screendoc.New(doc, Palette, modelData{m})
	if err != nil {
		t.Fatalf("building the renderer: %v", err)
	}
	if d := htmlround.Diff(htmlround.FromANSI(m.View()), htmlround.FromANSI(r.Render(doc, 132, 38))); d == "" {
		t.Log("the Events tab now round-trips — fold it into the table above")
	} else {
		t.Logf("known gap, a list row carries one style:\n%s", d)
	}
}
