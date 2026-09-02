package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"os"

	"github.com/charmbracelet/lipgloss"

	"github.com/richarddavenport/tuikit/comp"
	"github.com/richarddavenport/tuikit/paint"
	"github.com/richarddavenport/tuikit/term"
	"github.com/richarddavenport/tuikit/theme"
)

// pixels is the self-test: what did this terminal say, and does the picture
// actually arrive?
//
// It exists because every other check in this repo runs where there is no
// terminal. The goldens prove the fallback, the unit tests prove the encoders
// byte by byte, and neither can prove that foot draws the thing — only a person
// looking at foot can do that. So this prints the answers it got, then draws
// the same bar twice, and asks you whether the second one looks better.
func pixels(args []string) {
	fs := flag.NewFlagSet("pixels", flag.ContinueOnError)
	fs.SetOutput(os.Stderr)
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	p := comp.Detect()
	fmt.Printf("TERM         %s\n", os.Getenv("TERM"))
	if prog := os.Getenv("TERM_PROGRAM"); prog != "" {
		fmt.Printf("TERM_PROGRAM %s\n", prog)
	}
	fmt.Printf("graphics     %s\n", p.Mode)
	if p.Mode == term.None {
		// Detect stops after the first answer, so reporting the rest as zeroes
		// would be reporting measurements that were never taken.
		fmt.Println("cell size    not asked — nothing would use it")
		fmt.Println("background   not asked")
		fmt.Println("ramp         not asked")
	} else {
		fmt.Printf("cell size    %dx%d pixels\n", p.CellW, p.CellH)
		fmt.Printf("background   %s\n", hex(p.Background))
		fmt.Printf("ramp         %s → %s   (ANSI 5 and 13, read from your theme)\n",
			hex(p.Ramp.From), hex(p.Ramp.To))
	}
	fmt.Println()

	switch p.Mode {
	case term.None:
		fmt.Println("This terminal draws characters only, so everything below is characters.")
		fmt.Println("That is not a degraded mode — it is what every golden in this repo records.")
	case term.Sixel:
		fmt.Println("Sixel: the bar below is an opaque image pasted over blanked cells.")
		fmt.Printf("Its soft edges are baked against %s, because Sixel has no alpha.\n", hex(p.Background))
	case term.Kitty:
		fmt.Println("Kitty protocol: the bar below is an RGBA image at z=-1, behind the text.")
		fmt.Println("The characters are still there — select them and you will get characters.")
	}
	fmt.Println()

	sty := barStyles()
	draw := func(label string, value float64, wantPixels bool) {
		c := comp.NewCanvas(72, 1).WithGraphics(p)
		comp.Meter{
			Value: value, Label: label, Track: comp.Name("bar"), Pixels: wantPixels,
			Filled: &sty[0], Empty: &sty[1], LabelStyle: &sty[1],
		}.Draw(c, comp.Rect{X: 0, Y: 0, W: 72, H: 1}, comp.Region("meter"))
		fmt.Println(c.String())
	}

	fmt.Println("  cells only — what every terminal shows:")
	draw("51%", 0.51, false)
	draw("53%", 0.53, false)
	fmt.Println()
	fmt.Println("  with a picture offered — identical above unless your terminal can:")
	draw("51%", 0.51, true)
	draw("53%", 0.53, true)
	fmt.Println()

	if p.Mode != term.None {
		fmt.Println("  a gradient panel, to show the ramp at full height:")
		c := comp.NewCanvas(72, 5).WithGraphics(p)
		c.Fill(comp.Rect{X: 0, Y: 0, W: 72, H: 5}, " ", nil, comp.Region("panel"))
		ramp := p.Ramp
		c.Picture(comp.Region("panel"), func(w, h int) *image.RGBA {
			return paint.Panel{W: w, H: h, Ramp: ramp, Radius: h / 3,
				Bars: []float64{.2, .5, .35, .8, .6, .95, .4, .7}}.Image()
		})
		fmt.Println(c.String())
		fmt.Println()
	}

	fmt.Printf("If the two rows above look the same and %s says none, that is correct.\n", term.EnvOverride)
	fmt.Printf("Force a mode with %s=none|sixel|kitty to compare.\n", term.EnvOverride)
}

func hex(c color.RGBA) string { return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B) }

// barStyles is the two styles the bar needs, from the default palette — which
// is ANSI 0-15, so they are your terminal's own colours too.
func barStyles() [2]lipgloss.Style {
	return [2]lipgloss.Style{
		lipgloss.NewStyle().Foreground(theme.Default.Accent),
		lipgloss.NewStyle().Foreground(theme.Default.Muted),
	}
}
