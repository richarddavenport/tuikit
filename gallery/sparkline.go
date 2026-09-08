package gallery

import (
	"math"

	"github.com/richarddavenport/tuikit/comp"
)

// sparklineEntry is a number over time, in the states that decide whether the
// component is right: short, long, flat, and two series meant to be compared.
func (m *Model) sparklineEntry(s *styles) Entry {
	// A minute of CPU, sampled every second, with a spike near the end so the
	// right-hand edge is where the interesting part is.
	cpu := make([]float64, 60)
	for i := range cpu {
		cpu[i] = 18 + 12*math.Sin(float64(i)/7) + float64(i%5)
	}
	cpu[52], cpu[53], cpu[54] = 88, 96, 91

	steady := make([]float64, 60)
	for i := range steady {
		steady[i] = 40 + 3*math.Sin(float64(i)/4)
	}

	line := func(sp comp.Sparkline, values []float64, rows int) func(*comp.Canvas, comp.Rect, bool) {
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			band := comp.Layout{Constraints: []comp.Constraint{
				comp.Length(1), comp.Length(rows), comp.Fill(1),
			}}.Rows(r)
			comp.Bar{
				Left:  []comp.Segment{{Text: "cpu", Style: &s.title}},
				Right: []comp.Segment{{Text: "last 60s", Style: &s.muted}},
			}.Draw(c, band[0], comp.Region("demo.sparkhead"))
			sp.Draw(c, band[1], values, comp.Region("demo.spark"))
		}
	}

	// Two series at very different scales, drawn both ways. This is the state
	// worth looking at: auto-scaled they are indistinguishable, and that is
	// either what you want or a bug, depending on what the reader is asking.
	pair := func(max float64) func(*comp.Canvas, comp.Rect, bool) {
		small := make([]float64, 40)
		big := make([]float64, 40)
		for i := range small {
			small[i] = 4 + 2*math.Sin(float64(i)/5)
			big[i] = 70 + 20*math.Sin(float64(i)/5)
		}
		return func(c *comp.Canvas, r comp.Rect, _ bool) {
			rows := comp.Layout{Constraints: []comp.Constraint{
				comp.Length(1), comp.Length(3), comp.Length(1), comp.Length(3), comp.Fill(1),
			}}.Rows(r)
			c.Text(rows[0].X, rows[0].Y, "  4 rps", &s.muted, comp.Region("demo.sparklabel"))
			comp.Sparkline{Max: max, Style: &s.success}.Draw(c, rows[1], small, comp.Region("demo.spark"))
			c.Text(rows[2].X, rows[2].Y, "  80 rps", &s.muted, comp.Region("demo.sparklabel2"))
			comp.Sparkline{Max: max, Style: &s.pending}.Draw(c, rows[3], big, comp.Region("demo.spark2"))
		}
	}

	return Entry{
		Name: "Sparkline",
		Summary: "A number over time. The newest sample is always the last column, " +
			"so the bar under your cursor does not move when one arrives.",
		From: "k9s wrote internal/tchart; bottom vendored 55 kB of ratatui's own Chart " +
			"rather than use it, and the specialization it forked for is this one",
		Roles:  []string{"Accent", "Success", "Pending", "Muted"},
		Glyphs: []string{"▁", "▂", "▃", "▄", "▅", "▆", "▇", "█"},
		States: []State{
			{Name: "a minute of cpu", Note: "one row is the classic sparkline; the spike is near the right because that is where now is",
				Draw: line(comp.Sparkline{Style: &s.focused}, cpu, 1)},
			{Name: "taller", Note: "a rect of more than one row stacks full blocks and puts the partial on top",
				Draw: line(comp.Sparkline{Style: &s.focused}, cpu, 5)},
			{Name: "shorter than the pane", Note: "eight samples in a wide pane leave the LEFT blank — the right edge is now, and it does not move",
				Draw: line(comp.Sparkline{Style: &s.focused}, cpu[:8], 3)},
			{Name: "auto-scaled", Note: "no Max: the largest value on screen is the top, which shows the shape and hides the size",
				Draw: pair(0)},
			{Name: "one scale", Note: "the same two series with Max set — now they are comparable, which auto-scaling makes impossible",
				Draw: pair(100)},
			{Name: "flat", Note: "a steady number is a flat band rather than a jagged one, because auto-scale has a floor of zero and not of the minimum",
				Draw: line(comp.Sparkline{Style: &s.success}, steady, 3)},
			{Name: "all zero", Note: "nothing drawn, deliberately: an auto-scaled flat line at the bottom would look the same as a series at its peak",
				Draw: line(comp.Sparkline{Style: &s.muted}, make([]float64, 40), 3)},
			{Name: "no room", Note: "the state every component gets wrong first",
				Draw: func(c *comp.Canvas, r comp.Rect, _ bool) {
					comp.Sparkline{}.Draw(c, comp.Rect{X: r.X, Y: r.Y, W: r.W, H: 0}, cpu, comp.Region("demo.spark"))
				}},
		},
	}
}
