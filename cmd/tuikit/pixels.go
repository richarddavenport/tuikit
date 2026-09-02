package main

import (
	"flag"
	"fmt"
	"image"
	"image/color"
	"image/png"
	"os"
	"strings"

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
	png := fs.String("png", "", "write the pictures to PNG files as well, so you can see what they should look like")
	raw := fs.Bool("raw", false, "show the bytes the terminal replied with")
	debug := fs.Bool("debug", false, "send one picture with replies turned on, and print what the terminal says")
	if err := fs.Parse(args); err != nil {
		os.Exit(2)
	}

	p := comp.Detect()

	// Printed on request, and automatically when the answer is "none" and
	// nobody asked for that. "It says none" is not debuggable: it is what you
	// get from a terminal that cannot draw, from one that can but was too slow,
	// and from a multiplexer answering on the terminal's behalf. The bytes tell
	// those apart and nothing else does.
	showRaw := *raw
	if p.Mode == term.None && os.Getenv(term.EnvOverride) == "" {
		showRaw = true
	}
	if showRaw {
		got, reply := term.ProbeTTY(term.DefaultTimeout)
		fmt.Printf("query        %s\n", escape(term.QueryBytes))
		if reply == "" {
			fmt.Println("reply        (nothing — no controlling terminal, or it never answered)")
		} else {
			fmt.Printf("reply        %s\n", escape(reply))
		}
		fmt.Printf("read as      %s\n\n", got)
	}

	if forced := os.Getenv(term.EnvOverride); forced != "" {
		fmt.Printf("NOTE  %s=%s is set, so nothing was detected — this is what you TOLD it.\n", term.EnvOverride, forced)
		fmt.Println("      Forcing a protocol your terminal does not speak prints stray characters")
		fmt.Println("      rather than a picture. Unset it to find out what you actually have.")
		fmt.Println()
	}
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
		how := "measured"
		if !p.CellMeasured {
			how = "ASSUMED — the terminal declined to say, so pictures will be the wrong size"
		}
		fmt.Printf("cell size    %dx%d pixels (%s)\n", p.CellW, p.CellH, how)
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

	if *debug {
		debugKitty(p)
		return
	}
	if *png != "" {
		writeReference(*png, p)
		return
	}
	fmt.Println("Not sure what you are looking at? Write the same pictures to files and open them:")
	fmt.Println("    tuikit pixels -png /tmp/ref")
	fmt.Println("Whatever your terminal drew above should look like those.")
}

// debugKitty sends one image with q=0 and prints the terminal's answer.
//
// The whole point is that the normal path is silent by necessity. q=2 keeps a
// reply from arriving in the middle of a frame and being read as a keystroke —
// but it also means a refusal goes to nobody, and an image that never appears
// looks exactly like an image drawn somewhere you cannot see it.
func debugKitty(p comp.Pixels) {
	if p.Mode != term.Kitty {
		fmt.Printf("This terminal is %s, not kitty — nothing to ask.\n", p.Mode)
		return
	}
	tty, err := os.OpenFile("/dev/tty", os.O_RDWR, 0)
	if err != nil {
		fmt.Fprintln(os.Stderr, "tuikit pixels: no controlling terminal:", err)
		os.Exit(1)
	}
	defer tty.Close() //nolint:errcheck // done with it

	img := paint.Bar{W: 40 * p.CellW, H: p.CellH, Value: 0.6, Ramp: p.Ramp, Track: track(p.Ramp)}.Image()
	seq := term.EncodeKittyVerbose(img, 99)

	fmt.Printf("image        %d bytes of RGBA, %dx%d pixels\n",
		len(img.Pix), img.Bounds().Dx(), img.Bounds().Dy())
	fmt.Printf("escape       %d bytes, first chunk: %s…\n\n",
		len(seq), escape(firstN(seq, 90)))

	fmt.Println("sending it here, with a marker line under it:")
	fmt.Println()
	// DA1 appended so the read ends as soon as the terminal has finished,
	// rather than sitting out the whole timeout on one that stays quiet.
	answer := term.SendAndRead(tty, seq+"\x1b[c", term.DefaultTimeout)
	fmt.Println("^ the bar should be on the blank line above this one")
	fmt.Println()

	switch {
	case answer == "":
		fmt.Println("reply        (nothing) — accepted silently, or ignored entirely.")
		fmt.Println("             If you see no bar, the image was placed somewhere you")
		fmt.Println("             cannot see it, or z=-1 put it under an opaque background.")
	case strings.Contains(answer, ";OK"):
		fmt.Printf("reply        %s\n", escape(answer))
		fmt.Println("             Accepted. If no bar is visible, z=-1 is drawing it under")
		fmt.Println("             the cell background rather than under the text.")
	default:
		fmt.Printf("reply        %s\n", escape(answer))
		fmt.Println("             That is a refusal — the code after the id says why.")
	}
}

func firstN(s string, n int) string {
	if len(s) <= n {
		return s
	}
	return s[:n]
}

// writeReference writes the pictures to PNG.
//
// This is the answer to "how do I know what it should look like". Nothing in a
// terminal can tell you whether the thing you are seeing is the thing that was
// sent: a protocol the emulator does not speak prints as stray characters, one
// it half-speaks prints as a smear, and both look like a bug in the program. A
// file you can open settles it — if the PNG is right and the terminal is not,
// the encoder is fine and the terminal is the problem.
func writeReference(prefix string, p comp.Pixels) {
	ramp := p.Ramp
	if ramp == (paint.Ramp{}) {
		ramp = comp.NewCanvas(1, 1).Ramp()
	}
	cellW, cellH := p.CellW, p.CellH
	if cellW == 0 || cellH == 0 {
		cellW, cellH = term.DefaultCellW, term.DefaultCellH
	}

	for _, ref := range []struct {
		name string
		img  *image.RGBA
	}{
		{"bar-51", paint.Bar{W: 66 * cellW, H: cellH, Value: 0.51, Ramp: ramp, Track: track(ramp)}.Image()},
		{"bar-53", paint.Bar{W: 66 * cellW, H: cellH, Value: 0.53, Ramp: ramp, Track: track(ramp)}.Image()},
		{"panel", paint.Panel{W: 72 * cellW, H: 5 * cellH, Ramp: ramp, Radius: cellH,
			Bars: []float64{.2, .5, .35, .8, .6, .95, .4, .7}}.Image()},
	} {
		// Flattened against the terminal background, because that is what a
		// Sixel terminal is actually shown — a PNG with an alpha channel would
		// open on a white page and look like a different picture.
		name := fmt.Sprintf("%s-%s.png", prefix, ref.name)
		f, err := os.Create(name)
		if err != nil {
			fmt.Fprintln(os.Stderr, "tuikit pixels:", err)
			os.Exit(1)
		}
		if err := png.Encode(f, paint.Flatten(ref.img, p.Background)); err != nil {
			fmt.Fprintln(os.Stderr, "tuikit pixels:", err)
			os.Exit(1)
		}
		_ = f.Close()
		fmt.Println("wrote", name)
	}
	fmt.Println()
	fmt.Println("Open those. The two bars differ by two percent — a difference the")
	fmt.Println("character bar cannot show at all, because it has one step per column.")
}

// track is the unfilled part of the bar: the ramp's own start, mostly clear.
func track(r paint.Ramp) color.RGBA {
	c := r.From
	c.A = 60
	return c
}

func hex(c color.RGBA) string { return fmt.Sprintf("#%02x%02x%02x", c.R, c.G, c.B) }

// escape makes control characters visible, so a reply can be read and pasted
// into a bug report rather than executed by the terminal printing it.
func escape(s string) string {
	var b strings.Builder
	for _, r := range s {
		switch {
		case r == 0x1b:
			b.WriteString("ESC")
		case r < 0x20 || r == 0x7f:
			fmt.Fprintf(&b, "\\x%02x", r)
		default:
			b.WriteRune(r)
		}
	}
	return b.String()
}

// barStyles is the two styles the bar needs, from the default palette — which
// is ANSI 0-15, so they are your terminal's own colours too.
func barStyles() [2]lipgloss.Style {
	return [2]lipgloss.Style{
		lipgloss.NewStyle().Foreground(theme.Default.Accent),
		lipgloss.NewStyle().Foreground(theme.Default.Muted),
	}
}
