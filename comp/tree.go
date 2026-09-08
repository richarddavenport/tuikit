package comp

import (
	"strings"

	"github.com/richarddavenport/tuikit/theme"
)

// Tree is collapse state over a flattened hierarchy.
//
// # What it is not
//
// It is not a node type, and it does not draw. Both of those were the reason a
// tree was refused the first time: the cloud tool's rows carried a bucket and
// a resource, the deploy tool's carried a service and an action and a change,
// and a shared node type would have made both of them box their data or keep
// it twice.
//
// That refusal was right about the payload and wrong about the conclusion.
// What those tools have in common is not a node — it is the ANSWER TO ONE
// QUESTION: given a flat list of things with depths, and a set of collapsed
// ones, which are visible? That question has no payload in it.
//
// So the tool keeps its own rows, [List] still draws them, [Row.Depth] still
// indents them, and this decides which ones exist this frame.
//
// # Why depth rather than parents
//
// A flattened hierarchy is what every one of these tools already has, because
// it is what a list can draw. Depth is derivable from a parent pointer and not
// the other way round, and asking for parents would mean a tool that has
// depths — which is all of them — building a graph to hand over and throwing
// it away again.
//
// # What actually needs one
//
// Four tools in the survey have a tree, checked in their source on 2026-09-04
// rather than assumed: lazygit (`pkg/gui/filetree`, eleven files, plus a
// second one for commit files), dive (`dive/filetree/`), termshark
// (`pkg/pdmltree` and `widgets/copymodetree`) and fx (collapse state over a
// JSON document).
//
// Three that were first claimed to have one do not. yazi and superfile contain
// no match for "tree" or "collapse" anywhere in their sources, and ranger's
// only collapse is `collapse_preview` — the preview COLUMN when it has nothing
// to show. All three are Miller columns: parent, current, preview, side by
// side, with the hierarchy walked rather than drawn.
//
// That is the sharper claim, and it says what this type is for. Every FILE
// MANAGER in the field chose columns over a tree. What needs a tree is a
// hierarchy you cannot walk into — a docker image's layers, a JSON document, a
// packet dissection, a git status — where the whole shape has to be on screen
// at once because comparing across it is the point.
type Tree struct {
	// Collapsed is the keys that are shut, by the tool's own identity for a
	// row — a path, an ID, a name.
	//
	// A KEY rather than an index, because indices move. A filter that removes
	// one row above a collapsed directory would otherwise open a different
	// directory, which is the bug this design exists to make impossible.
	//
	// Exported so a tool can persist it, restore it, or collapse everything on
	// first draw without asking this type for permission.
	Collapsed map[string]bool
}

// Node is one row's place in the hierarchy: how deep it sits, and what to call
// it when it is collapsed.
//
// Key may be empty for a row that can never be collapsed — a leaf. Nothing is
// stored about it, so a tool with a million leaves pays nothing for them.
type Node struct {
	Depth int
	Key   string
}

// Visible is the indices to draw, in order.
//
// A row is hidden when any ancestor is collapsed. Ancestry is read off the
// depths as they are walked: the nearest earlier row with a smaller depth is
// the parent, so a hidden subtree ends at the first row shallow enough not to
// belong to it.
//
// Returns indices rather than nodes so the caller's own slice stays the source
// of the data — the same arrangement as [List.DrawFunc], and for the same
// reason: a tree of two hundred thousand should not be copied to be drawn.
//
// # The cost, and where it would bite
//
// This walks every node, so it is O(total) per frame rather than O(visible).
// fx solves the same problem differently and better at scale: its nodes are a
// doubly-linked list of lines, and a collapsed one hands back a Collapsed
// pointer to the node AFTER its subtree, so traversal skips what is folded
// without a separate pass.
//
// Returning indices is what buys [List] a draw with no copying, and the linked
// design gives that up. So this is right for hierarchies of thousands and
// would need rethinking at millions — the map lookup on Key, once per node per
// frame, is where it would show first. No tool in the survey has a
// million-node tree.
func (t *Tree) Visible(nodes []Node) []int {
	out := make([]int, 0, len(nodes))

	// hidden is the depth at which the current collapsed subtree began, or -1.
	// One int rather than a stack: a subtree inside a collapsed subtree is
	// already hidden, so the outermost collapse is the only one that matters.
	hidden := -1

	for i, node := range nodes {
		if hidden >= 0 {
			if node.Depth > hidden {
				continue
			}
			hidden = -1
		}
		out = append(out, i)
		if node.Key != "" && t.Collapsed[node.Key] {
			hidden = node.Depth
		}
	}
	return out
}

// Toggle opens a shut key and shuts an open one.
func (t *Tree) Toggle(key string) {
	if key == "" {
		return
	}
	if t.Collapsed == nil {
		t.Collapsed = map[string]bool{}
	}
	if t.Collapsed[key] {
		delete(t.Collapsed, key)
		return
	}
	t.Collapsed[key] = true
}

// IsCollapsed reports whether a key is shut, for a tool drawing the marker.
//
// The marker is the TOOL's, not this type's: ▸ and ▾ are one choice, + and -
// are another, and a file manager that wants neither should not have to accept
// one. [Row.Lead] is where it goes, and [Row.LeadStyle] keeps its color under
// the selection.
func (t *Tree) IsCollapsed(key string) bool { return t.Collapsed[key] }

// HasChildren reports whether the node at i has any, which is what decides
// between drawing a marker and drawing the blank that keeps the names in line.
//
// Read off the depths rather than stored, so it cannot disagree with them.
func HasChildren(nodes []Node, i int) bool {
	return i >= 0 && i+1 < len(nodes) && nodes[i+1].Depth > nodes[i].Depth
}

// Expand opens every key, and Collapse shuts every one that has children.
//
// The pair a tool binds to a key for "open all"/"close all", written here
// because doing it correctly means knowing that a leaf must not be collapsed —
// shutting one hides nothing and leaves a marker beside a row with no
// children.
func (t *Tree) Expand() { t.Collapsed = nil }

// Collapse shuts every node that has children.
func (t *Tree) Collapse(nodes []Node) {
	t.Collapsed = map[string]bool{}
	for i, node := range nodes {
		if node.Key != "" && HasChildren(nodes, i) {
			t.Collapsed[node.Key] = true
		}
	}
}

// Branches is the connector prefix for each node, for a tree that draws its
// shape rather than indenting.
//
// # Why a function and not [Row.Depth]
//
// Depth cannot answer it. Which connector a row gets depends on facts about
// the rows BETWEEN it and its parent:
//
//	init
//	├─ systemd-journald
//	├─ sshd
//	│  └─ sshd: richard
//	└─ firefox
//
// `sshd: richard` needs a │ in its first column because `sshd` still has a
// sibling coming. A row at the same depth under `firefox` would get a space.
// Two rows, one depth, different prefixes.
//
// htop keeps a bitmask per row for this — `Row.indent` in `Row.h`, built by
// `Table_buildTreeBranch`. dive keeps three string constants in
// `dive/filetree/file_tree.go`. Two tools, which is the bar for extracting.
//
// # Where the result goes
//
// [Row.Indent], for a list. Or anywhere: htop draws its tree INSIDE the
// command column, because a prefix in front of the PID makes every numeric
// column ragged, and a plain string can go there. Returning strings rather
// than drawing is what makes both possible.
//
// Takes the same []Node as [Tree.Visible] and returns one prefix per node, so
// a caller that filtered with Visible indexes this with the same indices.
func Branches(nodes []Node, ch theme.Chrome) []string {
	out := make([]string, len(nodes))
	// open[d] is whether the ancestor at depth d still has siblings below, and
	// therefore whether a vertical line passes through this row's column d.
	var open []bool

	for i, node := range nodes {
		d := max(0, node.Depth)
		if d > len(open) {
			d = len(open)
		}
		open = open[:d]

		var b strings.Builder
		// open[0] is the root level, which has no column: a depth-1 row starts
		// at the left edge with its own connector. The columns drawn are for
		// ancestors at depths 1..d-1.
		ancestors := open
		if len(ancestors) > 0 {
			ancestors = ancestors[1:]
		}
		for _, passing := range ancestors {
			if passing {
				b.WriteString(ch.Tree.Vertical)
			} else {
				b.WriteString(ch.Tree.Gap)
			}
		}
		if d > 0 {
			if lastChild(nodes, i, d) {
				b.WriteString(ch.Tree.Last)
			} else {
				b.WriteString(ch.Tree.Branch)
			}
		}
		out[i] = b.String()
		open = append(open, !lastChild(nodes, i, d))
	}
	return out
}

// lastChild reports whether i is the final node at its depth under its parent.
//
// Read off the depths rather than stored, the same way [HasChildren] is: a
// stored flag can disagree with the shape it describes, and a derived one
// cannot. The scan stops at the first node shallower than i, which is where
// the parent's run of children ends.
func lastChild(nodes []Node, i, depth int) bool {
	for j := i + 1; j < len(nodes); j++ {
		switch d := max(0, nodes[j].Depth); {
		case d < depth:
			return true
		case d == depth:
			return false
		}
	}
	return true
}
