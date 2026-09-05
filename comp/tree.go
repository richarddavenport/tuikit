package comp

// Tree is collapse state over a flattened hierarchy.
//
// # What it is not
//
// It is not a node type, and it does not draw. Both of those were the reason a
// tree was refused the first time: azctl's rows carried a bucket and a resource,
// swarmctl's carried a service and an action and a change, and a shared node
// type would have made both of them box their data or keep it twice.
//
// That refusal was right about the payload and wrong about the conclusion. What
// those tools have in common is not a node — it is the ANSWER TO ONE QUESTION:
// given a flat list of things with depths, and a set of collapsed ones, which
// are visible? That question has no payload in it.
//
// So the tool keeps its own rows, [List] still draws them, [Row.Depth] still
// indents them, and this decides which ones exist this frame.
//
// # Why depth rather than parents
//
// A flattened hierarchy is what every one of these tools already has, because
// it is what a list can draw. Depth is derivable from a parent pointer and not
// the other way round, and asking for parents would mean a tool that has
// depths — which is all of them — building a graph to hand over and throwing it
// away again.
//
// Read across seven tools that each wrote this: lazygit (`pkg/gui/filetree`,
// eleven files), yazi, superfile, dive, fx, termshark, ranger. Six of the
// fourteen surveyed needed a tree; the seventh is lazygit's second one, for
// commit files.
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
// one. [Row.Lead] is where it goes, and [Row.LeadStyle] keeps its colour under
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
// shutting one hides nothing and leaves a marker beside a row with no children.
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
