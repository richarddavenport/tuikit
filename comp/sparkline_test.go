package comp

import (
	"strings"
	"testing"
)

const sparkName Name = "spark"

func spark(w, h int, s Sparkline, values []float64) string {
	c := NewCanvas(w, h)
	s.Draw(c, c.Bounds(), values, Region(sparkName))
	return c.String()
}

// The rule bottom forked its framework's chart to get: the newest sample is in
// the LAST column, and a short series leaves the left blank.
func TestTheNewestSampleIsAtTheRightEdge(t *testing.T) {
	// Non-zero values, because a zero deliberately draws nothing — see
	// TestZeroDrawsNothing. Half of the peak is four eighths.
	got := strings.TrimRight(spark(10, 1, Sparkline{}, []float64{0.5, 1}), "\n")
	if !strings.HasSuffix(got, "▄█") {
		t.Errorf("the series is not against the right edge: %q", got)
	}
	if !strings.HasPrefix(got, "        ") {
		t.Errorf("the series was pushed to the left: %q", got)
	}
}

// A series longer than the pane drops its OLDEST samples, not its newest.
func TestALongSeriesDropsItsOldest(t *testing.T) {
	values := []float64{9, 9, 9, 9, 9, 9, 9, 9, 1}
	got := strings.TrimRight(spark(3, 1, Sparkline{}, values), "\n")
	if Width(got) != 3 {
		t.Errorf("drew %d columns into a pane of 3: %q", Width(got), got)
	}
	if !strings.HasSuffix(got, "▁") {
		t.Errorf("the newest sample is not last: %q", got)
	}
}

// A value present but tiny still draws, because a gap in a sparkline means no
// data rather than a small number.
func TestATinyValueStillDraws(t *testing.T) {
	got := strings.TrimRight(spark(4, 1, Sparkline{Max: 1000}, []float64{0, 0.001, 0, 1000}), "\n")
	if !strings.Contains(got, "▁") {
		t.Errorf("a tiny value drew nothing: %q", got)
	}
	if strings.Count(got, "█") != 1 {
		t.Errorf("want one full bar: %q", got)
	}
}

// Zero is not drawn, so a flat series and a gap look the same — which is
// correct, because both mean "nothing happened".
func TestZeroDrawsNothing(t *testing.T) {
	got := spark(4, 1, Sparkline{Max: 10}, []float64{0, 0, 0, 0})
	if strings.TrimSpace(got) != "" {
		t.Errorf("zeroes drew %q", got)
	}
}

// Every sample zero and no Max is the case an auto-scale cannot represent: a
// flat line at the bottom would look the same as a series at its peak.
func TestAnAllZeroSeriesWithNoMaxDrawsNothing(t *testing.T) {
	got := spark(4, 1, Sparkline{}, []float64{0, 0, 0})
	if strings.TrimSpace(got) != "" {
		t.Errorf("drew %q", got)
	}
}

// Auto-scale takes the largest value on screen, so the peak is always full.
func TestAutoScaleMakesThePeakFull(t *testing.T) {
	got := spark(4, 1, Sparkline{}, []float64{1, 2, 3, 4})
	if !strings.Contains(got, "█") {
		t.Errorf("nothing reached the top: %q", got)
	}
}

// Max fixes the scale, which is what makes two sparklines comparable. Without
// it, two series with different peaks look identical.
func TestMaxFixesTheScale(t *testing.T) {
	small := spark(4, 1, Sparkline{Max: 100}, []float64{1, 2, 3, 4})
	big := spark(4, 1, Sparkline{Max: 100}, []float64{50, 60, 70, 80})
	if small == big {
		t.Error("two series with the same Max drew identically")
	}
	auto1 := spark(4, 1, Sparkline{}, []float64{1, 2, 3, 4})
	auto2 := spark(4, 1, Sparkline{}, []float64{50, 100, 150, 200})
	if auto1 != auto2 {
		t.Error("auto-scale did not make proportional series identical, which is its whole behaviour")
	}
}

// A taller rect stacks full blocks and puts the partial on top.
func TestATallRectDrawsAColumnChart(t *testing.T) {
	got := spark(3, 3, Sparkline{Max: 3}, []float64{1, 2, 3})
	lines := strings.Split(got, "\n")
	if len(lines) != 3 {
		t.Fatalf("drew %d rows: %q", len(lines), got)
	}
	// The tallest column reaches the top row; the shortest does not.
	if strings.TrimSpace(lines[0]) != "█" {
		t.Errorf("the top row is %q, want only the tallest column", lines[0])
	}
	if strings.Count(strings.TrimSpace(lines[2]), "█") != 3 {
		t.Errorf("the bottom row is %q, want all three columns full", lines[2])
	}
}

// A bar grows from the BOTTOM. Drawn from the top it reads as a hanging chart
// and the shape is upside down.
func TestBarsGrowFromTheBottom(t *testing.T) {
	got := spark(1, 4, Sparkline{Max: 4}, []float64{1})
	lines := strings.Split(got, "\n")
	if strings.TrimSpace(lines[0]) != "" {
		t.Errorf("a quarter-height bar reached the top row: %q", got)
	}
	if strings.TrimSpace(lines[3]) == "" {
		t.Errorf("a quarter-height bar did not draw on the bottom row: %q", got)
	}
}

// A value above Max is clamped rather than overflowing into the row above.
func TestAValueAboveMaxIsClamped(t *testing.T) {
	got := spark(2, 1, Sparkline{Max: 10}, []float64{5, 1000})
	if strings.Count(got, "█") != 1 {
		t.Errorf("an over-max value drew %q", got)
	}
}

// Every cell says who drew it, so a click on a sparkline resolves.
func TestTheSparklineOwnsItsCells(t *testing.T) {
	c := NewCanvas(4, 1)
	Sparkline{Max: 1}.Draw(c, c.Bounds(), []float64{1, 1, 1, 1}, Region(sparkName))
	if got := c.OwnerAt(2, 0); got != Region(sparkName) {
		t.Errorf("a bar is owned by %v", got)
	}
}

// No room, and no values. The two states every component gets wrong first.
func TestSparklineWithNoRoomOrNoDataIsQuiet(t *testing.T) {
	c := NewCanvas(10, 2)
	Sparkline{}.Draw(c, Rect{}, []float64{1, 2}, Region(sparkName))
	Sparkline{}.Draw(c, c.Bounds(), nil, Region(sparkName))
	if strings.TrimSpace(c.String()) != "" {
		t.Errorf("drew %q", c.String())
	}
}
