// Package fuzzy matches a typed query against a list of things, and says WHERE
// it matched.
//
// The where is the part that is easy to leave out and hard to add later. A
// matcher that answers yes or no can filter a list; only one that returns
// positions can underline the letters you typed, and without that a ranked
// list is a list whose order you have to take on trust. The the deploy tool
// palette design asks for the matched substring to be underlined, and it is
// right to: it is the difference between "these five things matched" and "here
// is why".
//
// It is a package rather than part of comp because it draws nothing. comp
// holds components; this is the arithmetic underneath one.
package fuzzy

import (
	"math"
	"sort"
	"strings"
	"unicode"
)

// Result is one match: how good, and which runes of the text were used.
type Result struct {
	// Score is higher for better. Comparable only within one query — it is a
	// ranking, not a measurement, and nothing should show it to a reader.
	Score int
	// At holds RUNE indices into the text, ascending. Runes because they will
	// be used to style characters on screen, and a byte index into a
	// multi-byte character is not a position.
	At []int
}

// Bonuses and penalties. Tuned so that the shape of a match matters more than
// its length: `dsk` finding "disk usage" should beat `dsk` finding a longer
// string that happens to contain the same three letters further apart.
const (
	bonusConsecutive = 8  // the letters were typed as they appear
	bonusWordStart   = 12 // after a space, dash, slash, dot, or a camelCase hump
	bonusFirst       = 16 // the very first character of the text
	bonusExactCase   = 2  // matched without having to lower-case it
	penaltyGap       = 1  // per rune skipped between matches
)

// Match scores query against text, case-insensitively.
//
// An empty query matches everything with no highlights, which is what makes an
// unfiltered list and a filtered one the same code path — the palette shows
// all 41 commands before you type anything, and that is not a special case.
func Match(query, text string) (Result, bool) {
	if query == "" {
		return Result{}, true
	}
	return best([]rune(query), []rune(text))
}

// best is the highest-scoring alignment of q into t, by dynamic programming.
//
// # Why not greedy
//
// The obvious implementation walks forward taking the first letter that fits.
// It finds a match when one exists, and it highlights the wrong letters: `ge`
// against "generate env" underlines the g and the e of "generate" rather than
// the pair a reader would have picked.
//
// The obvious fix — a backward pass pulling each match as late as it can go —
// is worse, and was written here first. Moving each index independently pulls
// letters OUT of the consecutive runs that make a match good: `env` against
// "env of this service" lost its v to the one in "service", turning a perfect
// three-letter run into a match with a thirteen-rune gap.
//
// The shape of a match is not a property of any one letter, so no per-letter
// rule can find it. The table below scores whole alignments, which is the
// smallest thing that can weigh "these three letters are adjacent" against
// "this letter is at the start of a word".
func best(q, t []rune) (Result, bool) {
	n, m := len(q), len(t)
	if n > m {
		return Result{}, false
	}
	const none = math.MinInt / 2

	score := make([][]int, n)
	from := make([][]int, n)
	for i := range score {
		score[i], from[i] = make([]int, m), make([]int, m)
		for j := range score[i] {
			score[i][j], from[i][j] = none, -1
		}
	}
	for j := 0; j < m; j++ {
		if eq(q[0], t[j]) {
			score[0][j] = place(q[0], t, j)
		}
	}
	for i := 1; i < n; i++ {
		for j := i; j < m; j++ {
			if !eq(q[i], t[j]) {
				continue
			}
			pick, at := none, -1
			for k := i - 1; k < j; k++ {
				if score[i-1][k] == none {
					continue
				}
				v := score[i-1][k] - (j-k-1)*penaltyGap
				if k == j-1 {
					v += bonusConsecutive
				}
				if v > pick {
					pick, at = v, k
				}
			}
			if at >= 0 {
				score[i][j], from[i][j] = pick+place(q[i], t, j), at
			}
		}
	}

	end, total := -1, none
	for j := n - 1; j < m; j++ {
		if score[n-1][j] > total {
			end, total = j, score[n-1][j]
		}
	}
	if end < 0 {
		return Result{}, false
	}

	at := make([]int, n)
	for i := n - 1; i >= 0; i-- {
		at[i] = end
		end = from[i][end]
	}
	// Shorter texts win ties: between two commands containing what you typed,
	// the one with less around it is the one you meant more often than not.
	return Result{Score: total - len(t)/8, At: at}, true
}

// place is what one matched rune is worth where it landed.
func place(q rune, t []rune, j int) int {
	bonus := 0
	switch {
	case j == 0:
		bonus = bonusFirst
	case wordStart(t, j):
		bonus = bonusWordStart
	}
	if q == t[j] {
		bonus += bonusExactCase
	}
	return bonus
}

func wordStart(t []rune, i int) bool {
	prev := t[i-1]
	if strings.ContainsRune(" \t-_/.:", prev) {
		return true
	}
	// A camelCase hump: lower then upper.
	return unicode.IsLower(prev) && unicode.IsUpper(t[i])
}

func eq(a, b rune) bool { return a == b || unicode.ToLower(a) == unicode.ToLower(b) }

// Ranked is one candidate's result, with the index it came from.
type Ranked struct {
	Index int
	Result
}

// Rank matches a query against several texts and returns the survivors, best
// first.
//
// The order is stable in the input order for equal scores, so a list with
// nothing typed into it is in the order the caller built it — the palette's
// tiers stay in nearest-first order until a query actually has an opinion.
func Rank(query string, texts []string) []Ranked {
	out := make([]Ranked, 0, len(texts))
	for i, text := range texts {
		if r, ok := Match(query, text); ok {
			out = append(out, Ranked{Index: i, Result: r})
		}
	}
	sort.SliceStable(out, func(i, j int) bool { return out[i].Score > out[j].Score })
	return out
}
