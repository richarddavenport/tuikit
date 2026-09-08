# Using tuikit

One page per package: what it is for, **what it can do**, and what it will not
do. Written for people building a tool. For why any of it is shaped this way,
`design/` has the reasoning and `design/decisions.md` the numbered record.

| | | |
| --- | --- | --- |
| [`comp`](comp.md) | the cell canvas and 25 components | drawing |
| [`app`](app.md) | key routing, mouse, screen stack, async | the shell |
| [`spec`](spec.md) | one command declaration, four surfaces | the command line |
| [`theme`](theme.md) | colour roles, a closed glyph set, box characters | the vocabulary |
| [`guard`](guard.md) | nine tests that hold the vocabulary closed | correctness |
| [`harness`](harness.md) | drive a model, capture frames, compare goldens | testing |
| [`term`](term.md) | what this terminal can do | capability |
| [`paint`](paint.md) | gradients, panels and bars, as images | pixels |
| [`fuzzy`](fuzzy.md) | ranked matching that reports where it matched | search |
| [`docgen`](docgen.md) | frames and the design system as HTML or Markdown | documentation |
| [`scaffold`](scaffold.md) | `tuikit new` | starting |

Two smaller packages are not in that table because you call them through a
command rather than an import: `news` backs `tuikit news`, and `watch` backs
`tuikit watch`. Both are described in [tooling](tooling.md).

## The shortest possible tool

```go
type model struct{ list comp.List }

func (m *model) Init() tea.Cmd { return nil }

func (m *model) Update(msg tea.Msg) (app.Model, tea.Cmd) {
	if k, ok := msg.(tea.KeyMsg); ok && k.String() == "j" {
		m.list.Move(1)
	}
	return m, nil
}

func (m *model) Draw(c *comp.Canvas, r comp.Rect) {
	m.list.Draw(c, r, []comp.Row{{Text: "one"}, {Text: "two"}})
}

func main() { tea.NewProgram(app.New(&model{})).Run() }
```

Three methods, and no `View() string`. `Draw` takes a canvas and a rect and has nowhere else to put
anything — there is no `View() string` to smuggle a hand-joined string through,
which is the seam that lets a golden record something a component never drew.

## The rule the whole thing is built around

**Nothing is invented.** Every component here was pulled out of tools that had
already written it — usually twice, differently — and each says in its own doc
comment which tools, which files, and what the two versions disagreed about.
When you wonder why something works the way it does, that is where the answer
is: `go doc github.com/richarddavenport/tuikit/comp.List`.
