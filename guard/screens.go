package guard

import (
	"go/ast"
	"go/token"
	"strconv"
)

// Screens reports a screen constant that nothing draws.
//
// # The failure this catches
//
// A screen is a constant, and drawing it is a separate declaration — a map
// entry, historically a switch case. Adding the constant and forgetting the
// other half compiles, runs, and shows an empty window. `app.Screens` makes
// that a visible complaint rather than a blank terminal, which is better, but
// it is still something a person has to reach before anyone knows.
//
// # Why it parses rather than counting
//
// democtl checked this by hand:
//
//	for s := screenDashboard; s <= screenRun; s++ { ... }
//
// A hand-written range over the constants, which is the list that goes stale —
// a screen added after screenRun is not checked, and the check that exists to
// catch a forgotten screen is itself a place to forget one. The enumeration has
// to come from the source.
//
// # What it cannot do
//
// Screens are constants of one type, so they can be enumerated and coverage is
// decidable. STATES — screen by focus by tab by filter by modal — are
// combinatorial, and which are worth a frame is editorial. The promise is every
// screen, not every state.
//
// It understands `iota` blocks and explicit integer literals, which is what a
// screen enumeration is. A constant whose value it cannot evaluate is reported
// as such rather than skipped, because a screen silently not checked is the
// thing this exists to prevent.
func Screens(t T, dir, typeName string, drawn func(value int) bool) {
	t.Helper()

	found := false
	for _, f := range sources(t, dir) {
		for _, decl := range f.ast.Decls {
			gen, ok := decl.(*ast.GenDecl)
			if !ok || gen.Tok != token.CONST {
				continue
			}
			for _, c := range constsOfType(gen, typeName) {
				found = true
				switch {
				case !c.known:
					t.Errorf("%s: %s is a %s whose value cannot be read here — "+
						"guard.Screens understands iota and integer literals, and a screen "+
						"it cannot evaluate is a screen it cannot check",
						f.name, c.name, typeName)
				case !drawn(c.value):
					t.Errorf("%s: %s is a %s that nothing draws — "+
						"a screen constant without a view renders an empty window and says nothing about why",
						f.name, c.name, typeName)
				}
			}
		}
	}
	if !found {
		t.Fatalf("guard: no constants of type %s in %s — a guard that enumerates nothing passes silently", typeName, dir)
	}
}

// screenConst is one constant of the named type.
type screenConst struct {
	name  string
	value int
	known bool
}

// constsOfType reads one const block, carrying the type and the iota position
// down the specs the way Go does.
func constsOfType(gen *ast.GenDecl, typeName string) []screenConst {
	var out []screenConst
	typed, usesIota := false, false

	for i, spec := range gen.Specs {
		vs, ok := spec.(*ast.ValueSpec)
		if !ok {
			continue
		}
		// A spec with its own type starts a new run; one without inherits the
		// last, which is how `a T = iota` followed by bare names works.
		if id, ok := vs.Type.(*ast.Ident); ok {
			typed = id.Name == typeName
			usesIota = false
		} else if vs.Type != nil {
			typed = false
		}
		if len(vs.Values) == 1 {
			if id, ok := vs.Values[0].(*ast.Ident); ok && id.Name == "iota" {
				usesIota = true
			}
		}
		if !typed {
			continue
		}

		for _, name := range vs.Names {
			if name.Name == "_" {
				continue
			}
			c := screenConst{name: name.Name}
			switch {
			case usesIota:
				// The position in the block IS the value, which is what iota
				// means. Anything more arithmetic than that is not evaluated.
				c.value, c.known = i, len(vs.Values) <= 1
			case len(vs.Values) == 1:
				if lit, ok := vs.Values[0].(*ast.BasicLit); ok && lit.Kind == token.INT {
					if v, err := strconv.Atoi(lit.Value); err == nil {
						c.value, c.known = v, true
					}
				}
			}
			out = append(out, c)
		}
	}
	return out
}
