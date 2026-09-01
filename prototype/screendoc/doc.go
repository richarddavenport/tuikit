// Package screendoc is a throwaway prototype: it asks whether a screen can be
// DECLARED rather than drawn, by declaring one that already exists and diffing
// the result against the hand-written original.
//
// Not for merge. The question it answers is the one a visual builder rests on
// — a builder saves a document, and if no document can express a screen as
// ordinary as democtl's dashboard then there is nothing for a builder to save.
// Doing it against an existing screen rather than a new one is the whole point:
// a format invented alongside its first document always fits.
//
// # What a document may not contain
//
// No colours, only role names — and the roles are declared once at the top as
// styles, which is newStyles() written as data. No positions, only constraints.
// No glyphs outside the set. A builder over this format cannot express a raw
// hex colour or a box-drawing character the font may lack, which is a stronger
// guarantee than a guard: the guard catches it afterwards, the format cannot
// say it at all.
package screendoc

// Doc is one screen.
type Doc struct {
	// Screen names it, matching the screen constant it renders.
	Screen string `json:"screen"`
	// Styles is the document's whole vocabulary of appearance, built from
	// palette ROLES rather than colours. Everything below refers to these by
	// name.
	Styles map[string]Style `json:"styles"`
	// Glyphs is every non-ASCII character the document is allowed to draw,
	// listed so a guard can hold it closed without parsing the layout.
	Glyphs []string `json:"glyphs"`
	Root   Node     `json:"root"`
}

// Style is one entry in the document's vocabulary: a role, optionally a
// background role, optionally bold. Hex is unrepresentable on purpose.
type Style struct {
	Role string `json:"role"`
	BG   string `json:"bg,omitempty"`
	Bold bool   `json:"bold,omitempty"`
}

// Constraint is how a child claims space from its parent. One of these, not a
// position: a screen laid out by coordinates is a screen that breaks at 80
// columns.
type Constraint struct {
	Fixed int   `json:"fixed,omitempty"`
	Ratio []int `json:"ratio,omitempty"` // {1,3} is a third
	Min   int   `json:"min,omitempty"`
	Fill  bool  `json:"fill,omitempty"`
}

// Node is one element. A single struct rather than an interface per type
// because the format has to survive being written by a builder, read by a
// person, and diffed in review — and a tagged union in JSON is the shape all
// three tolerate.
type Node struct {
	Type string `json:"type"`
	// When gates the node on a value from the data. A screen has conditional
	// content — a note only some services have, a tab's worth of fields, a
	// title that changes while you type — and a format that cannot branch can
	// describe a picture but not a screen.
	When *When `json:"when,omitempty"`
	// ID becomes the owner ID of every cell this node draws. An identity, not
	// a position: after the list scrolls, a click still resolves to the row it
	// is over.
	ID         string      `json:"id,omitempty"`
	Constraint *Constraint `json:"constraint,omitempty"`

	// column, row
	Children []Node `json:"children,omitempty"`
	Gap      int    `json:"gap,omitempty"`
	// Reserve is rows the column deliberately does not draw into. democtl
	// leaves the terminal's last line blank; a format that cannot say so
	// cannot reproduce the screen.
	Reserve int `json:"reserve,omitempty"`

	// statusbar
	Left  []Span `json:"left,omitempty"`
	Right []Span `json:"right,omitempty"`

	// text, rule, keyhints
	Style string `json:"style,omitempty"`
	Text  string `json:"text,omitempty"`

	// box
	Title string `json:"title,omitempty"`
	// Titles is Title with conditions: the first whose When matches wins, and
	// an entry with no When is the default. Kept separate from Title so the
	// common case stays one line.
	Titles []Titled `json:"titles,omitempty"`
	// Focused names a boolean in the data. Focus is state, not layout, so the
	// document says which flag decides it rather than baking one in.
	Focused string `json:"focused,omitempty"`
	Body    []Node `json:"body,omitempty"`

	// list
	Rows     string `json:"rows,omitempty"`     // data key
	Template string `json:"template,omitempty"` // " {name:-14} {ready}/{want} {mark}"
	// StyleBy names a field on each row; StyleMap turns its value into a style
	// name. This pair is where "a failed service is red" lives — a UI decision,
	// so it belongs to the document rather than to the engine, which is not
	// allowed to know what red means.
	StyleBy  string            `json:"styleBy,omitempty"`
	StyleMap map[string]string `json:"styleMap,omitempty"`
	// GlyphBy/GlyphMap is the same idea for the mark a row ends with.
	GlyphBy  string            `json:"glyphBy,omitempty"`
	GlyphMap map[string]string `json:"glyphMap,omitempty"`
	Cursor   string            `json:"cursor,omitempty"` // data key: selected index
	// Selected is used when the pane has focus, Unfocused when it does not.
	Selected   string `json:"selected,omitempty"`
	Unfocused  string `json:"unfocused,omitempty"`
	Empty      string `json:"empty,omitempty"`
	EmptyStyle string `json:"emptyStyle,omitempty"`

	// tabs
	Names     []string `json:"names,omitempty"`
	Active    string   `json:"active,omitempty"` // data key: active index
	On        string   `json:"on,omitempty"`     // style for the active tab
	Off       string   `json:"off,omitempty"`
	OpenWith  string   `json:"openWith,omitempty"`
	CloseWith string   `json:"closeWith,omitempty"`

	// fields
	Fields []Field `json:"fields,omitempty"`
	// LabelStyle and LabelWidth are the shape of a name/value row.
	LabelStyle string `json:"labelStyle,omitempty"`
	LabelWidth int    `json:"labelWidth,omitempty"`
}

// When is a condition on a value the engine supplies. Is compares for
// equality, Not for inequality; an empty Not tests that a value is present at
// all, which is how "this service has a note" is written.
type When struct {
	Key string  `json:"key"`
	Is  string  `json:"is,omitempty"`
	Not *string `json:"not,omitempty"`
}

// Titled is one candidate title.
type Titled struct {
	When *When  `json:"when,omitempty"`
	Text string `json:"text"`
}

// Span is a run of text inside a status bar, with a style and possibly a
// condition.
type Span struct {
	Style string `json:"style,omitempty"`
	Text  string `json:"text,omitempty"`
	// Switch names a boolean in the data; Cases gives the span to use for
	// "true" and "false". A screen has conditional content, so a document that
	// cannot branch cannot describe one — and this is the part a drag-and-drop
	// builder will find hardest to offer.
	Switch string          `json:"switch,omitempty"`
	Cases  map[string]Span `json:"cases,omitempty"`
}

// Field is one name/value row in a detail pane.
type Field struct {
	Name  string `json:"name"`
	Value string `json:"value"`
	// StyleBy/StyleMap colours the value by some other value, the same way a
	// list row is coloured.
	StyleBy  string            `json:"styleBy,omitempty"`
	StyleMap map[string]string `json:"styleMap,omitempty"`
}

// Data is what the engine supplies. Deliberately narrow and stringly-typed:
// the engine has no terminal concepts, so it hands over values and rows, and
// every decision about what they LOOK like is made by the document.
type Data interface {
	Value(key string) string
	Bool(key string) bool
	Index(key string) int
	Rows(key string) []map[string]string
}
