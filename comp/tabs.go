package comp

import "github.com/charmbracelet/lipgloss"

// Tabs is a strip of names, one of them current.
//
// # Where this came from
//
// swarmctl (pane.go:365), pgctl (pane.go:60) and democtl each have one. They
// agree on the shape — style the active one, style the rest quietly — and
// differ on two things.
//
// The chevrons. swarmctl wraps the WHOLE strip in ‹ ›; democtl wraps the ACTIVE
// TAB. democtl's is the accident: the active tab is already styled, so the
// chevrons around it repeat what the colour says, and repeat it in the one
// place a reader is already looking. swarmctl's say something no styling can —
// that the strip CYCLES, and that ‹ and › are the keys — which is a fact about
// the keymap rather than about the state. Its comment says so outright. So the
// chevrons go around the strip.
//
// The count. swarmctl alone can write Pending(3), for a tab with things behind
// it. That is a real requirement rather than an accident, and cheap: a tab with
// nothing pending simply has no count.
//
// pgctl's third difference — active-when-focused styled differently from
// active-when-not — is kept, from swarmctl, which does the same. A tab strip
// that looks identical whether or not you can press the keys is a strip that
// invites you to press them.
type Tabs struct {
	Tabs   []Tab
	Active int
	// Focused changes how the current tab is drawn, so a strip you cannot
	// operate does not look like one you can.
	Focused bool

	// Style is a tab that is not current; Selected is the current one, and
	// FocusSelected replaces it when focused. Chrome draws the chevrons and
	// the separators.
	Style, Selected, FocusSelected, Chrome *lipgloss.Style
}

// Tab is one entry. A Count above zero is drawn beside the name, for a tab with
// things waiting behind it.
type Tab struct {
	Name  string
	Count int
}

// Draw renders the strip along the top of r and returns the columns it used.
//
// Each tab owns its own region, indexed, so a click lands on the tab rather
// than on the strip.
func (t Tabs) Draw(c *Canvas, r Rect, name Name) int {
	c = c.Clip(r)
	x := r.X
	x += c.Text(x, r.Y, "‹", t.Chrome, Region(name))

	for i, tab := range t.Tabs {
		if i > 0 {
			x += c.Text(x, r.Y, "·", t.Chrome, Region(name))
		}
		label := " " + tab.Name
		if tab.Count > 0 {
			label += "(" + itoa(tab.Count) + ")"
		}
		label += " "

		style := t.Style
		if i == t.Active {
			style = t.Selected
			if t.Focused && t.FocusSelected != nil {
				style = t.FocusSelected
			}
		}
		x += c.Text(x, r.Y, label, style, Region(name).At(i))
	}
	x += c.Text(x, r.Y, "›", t.Chrome, Region(name))
	return x - r.X
}
