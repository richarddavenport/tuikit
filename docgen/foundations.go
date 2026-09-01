package docgen

import (
	"fmt"
	"html"
	"strings"
)

// The foundations: the three things a terminal interface is made of that a
// design tool does not know about on its own.

func (r renderer) colorsPage() string {
	var b strings.Builder
	b.WriteString(`<div style="display:grid;gap:10px;max-width:760px">`)
	for _, role := range r.pal.Roles() {
		fmt.Fprintf(&b, `<div style="display:grid;grid-template-columns:64px 132px 1fr;gap:16px;align-items:center">
  <div style="height:40px;border-radius:4px;background:%s;border:1px solid #2a2a30"></div>
  <div><code style="color:%s">%s</code><br><code style="color:#6e6e6e">%s · %s</code></div>
  <div style="color:#9a9a9a;font-size:13px">%s</div>
</div>`, hexOf(role.Color), hexOf(role.Color), html.EscapeString(role.Name),
			html.EscapeString(valueOf(role.Color)), hexOf(role.Color), html.EscapeString(role.Why))
	}
	b.WriteString(`</div>`)

	// Every role drawn in a line of interface, not only as a swatch. A colour
	// is a decision about what something MEANS, and a swatch cannot show that.
	b.WriteString(`<h2>In place</h2>`)
	b.WriteString(r.rolesInUse())

	b.WriteString(`<h2>Why indices, not hex</h2>
<p>The values are ANSI 256 palette <em>indices</em>, because that is what a terminal
understands and what every terminal has agreed on. A truecolour hex would look
right on the machine it was picked on and wrong over ssh from another. The hex
shown here is what that index resolves to in a default palette — it is for
drawing the interface <em>outside</em> a terminal, and is never what the code
sends.</p>
<p>Names are roles, never hues. <code>Accent</code> survives someone deciding the
interface should be blue; <code>pink</code> does not.</p>`)

	return r.page("Foundations", "Colour roles",
		"Each role is an ANSI 256 index. The interface has no other colours — guard.Tokens holds that closed.",
		b.String())
}

// rolesInUse draws one line per role, generated from the palette rather than
// written out, so a tool that adds an Extra role sees it here without editing
// this file — and so the "every role is drawn" guard cannot be satisfied by a
// swatch alone.
func (r renderer) rolesInUse() string {
	lines := make([]string, 0, 7+len(r.pal.Extra))
	lines = append(lines,
		r.span("accent", "> api_api")+plain("      2/2  ")+r.span("success", "ok"),
		plain("  api_worker   1/1  ")+r.span("pending", "→ replicas 3"),
		plain("  api_migrate  0/1  ")+r.span("danger", "✗ rejected"),
		plain("  api_logs     1/1  ")+r.span("stderr", "warning: retrying"),
		`<span class="selected">  api_web      3/3  the selected row</span>`,
		r.span("border", "└──────────────────────────────────────┘"),
		r.span("muted", "  j/k move · X remove · q quit"),
	)
	for _, role := range r.pal.Extra {
		lines = append(lines, plain("  ")+r.span(cssClass(role.Name), role.Name+" — "+role.Why))
	}
	return term(lines...)
}

func (r renderer) glyphsPage() string {
	var b strings.Builder
	b.WriteString(`<div style="display:grid;grid-template-columns:repeat(auto-fill,minmax(240px,1fr));gap:8px">`)
	for _, g := range r.glyphRows() {
		fmt.Fprintf(&b, `<div style="display:flex;gap:14px;align-items:baseline;padding:8px 10px;background:#0c0c0e;border:1px solid #26262c;border-radius:4px">
  <span style="font-family:var(--mono);font-size:19px;color:#e8e8e8;width:1.2em;text-align:center">%s</span>
  <span style="font-size:12px;color:#8a8a8a"><code style="color:#6e6e6e">U+%04X</code><br>%s</span>
</div>`, html.EscapeString(string(g.R)), g.R, html.EscapeString(g.Why))
	}
	b.WriteString(`</div>`)

	b.WriteString(`<h2>The set is closed</h2>
<p>A terminal font without a glyph draws a replacement box, which reads as a bug
rather than as decoration — that is exactly what happened to a block-character
edit cursor. Adding one means adding it to the tool's glyph set and deciding,
deliberately, that it is common enough.</p>
<p><strong>Not available:</strong> block elements (U+2580–U+259F), geometric shapes
beyond the plain bullet, emoji, and Nerd Font private-use icons. If a design needs
a shape that is not here, it needs a different design — not a different font.</p>
<p>The spinner is the one exception, and it is written down rather than left as a
hole: its Braille cells (U+2800–U+28FF) are drawn by <code>bubbles</code> rather
than by the tool, so no string literal contains them and the guard never sees
them. Braille is in every terminal font, which is why a spinner reaches for it
instead of the block elements that would be a box on someone's screen.</p>`)

	return r.page("Foundations", "Glyphs",
		"Every non-ASCII character the interface may print. Anything outside this set is a box on someone's terminal.",
		b.String())
}

func (r renderer) gridPage() string {
	ruler := r.span("muted", "0        1         2         3         4         5") + "\n" +
		r.span("muted", "1234567890123456789012345678901234567890123456789012")
	body := `<h2>One character, one cell</h2>` +
		term(ruler,
			plain("┌─[3] Stacks (3)───────────────────────────────────┐"),
			plain("│")+r.span("accent", "> demo 3/3")+plain("                                       │"),
			plain("└──────────────────────────────────────────────────┘")) +
		`<p>Width is measured in <em>columns</em>, never pixels — and not in bytes either.
A right arrow is three bytes and one column, which is why the code measures with
<code>lipgloss.Width</code> and a test compares against that rather than against
<code>len</code>.</p>

<h2>Nothing may be wider than the terminal</h2>
<p>The commonest bug in a TUI, and the one assertions miss: a box with no bound
drawn from content that turned out long — a description listing every affected
row, or a file path from a real machine rather than a fixture's short one. Both
happened. Both were invisible until someone looked at a captured frame.</p>

<h2>Two ways to shorten a line, and they are not interchangeable</h2>
<p>Truncating counts runes and is right for plain text. A line that already
carries colour has to be <em>clipped</em> instead: truncation counts the escape
sequences as width, so a coloured row gets cut short of the pane and, worse, cut
mid-escape.</p>

<h2>The interface paints no background</h2>
<p>Except one: the selected row of a table, where a foreground colour alone cannot
be seen against its neighbours. Everything else sits on whatever the reader's
terminal is, which is why the same design has to hold up both ways.</p>` +
		`<div style="display:grid;grid-template-columns:1fr 1fr;gap:12px">` +
		`<div class="term">` + r.span("accent", "> api_api") + plain("   2/2  ") + r.span("success", "ok") + "\n" +
		plain("  api_web") + plain("    3/3") + `</div>` +
		`<div class="term light">` + r.span("accent", "> api_api") + plain("   2/2  ") + r.span("success", "ok") + "\n" +
		plain("  api_web") + plain("    3/3") + `</div>` +
		`</div>`

	return r.page("Foundations", "The grid",
		"A character grid, not a canvas: columns rather than pixels, and a font with no ligatures.",
		body)
}
