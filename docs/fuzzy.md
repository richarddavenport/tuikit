# `fuzzy` — ranked matching that says where it matched

```go
for _, hit := range fuzzy.Rank(query, names) {
	row := comp.Row{Spans: comp.Highlight(names[hit.Index], hit.At, nil, &accent)}
}
```

| | |
| --- | --- |
| `Match(query, text) (Result, bool)` | does it match, how well, and **which runes** |
| `Rank(query, texts) []Ranked` | the matches, best first, each with its `Index` and `At` |

`At` is the point. A ranked list whose order you have to take on trust is a
worse list than an unranked one: the reader cannot tell whether the top row is
first because it is the best match or because it was first already. Marking the
letters that matched turns the order into something they can check at a glance —
which is why `comp.Highlight` takes exactly the indices `Match` returns.

Used by `comp.Palette`, `comp.Input` filters, and any list with a search over
it.

## What it cannot do

Not fzf. No multi-token queries, no negation, no field-scoped search, no
configurable scoring. One query, one string, a score and the positions.
