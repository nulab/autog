package phase3

type Alg uint8

// Phase returns the ordinal number of this phase: 3.
func (alg Alg) Phase() int {
	return 3
}

// String returns a mnemonic representation of this algorithm.
// The exact string values are not documented and may change in the future.
func (alg Alg) String() (s string) {
	switch alg {
	case NoOrdering:
		s = "noop"
	case GraphvizDot:
		s = "gvdot"
	default:
		s = "<invalid>"
	}
	return
}

const (
	// NoOrdering does nothing. Nodes won't be reordered in their layers to minimize edge crossings.
	NoOrdering Alg = iota

	// GraphvizDot implements the mincross heuristic used in dot. It attempts to minimize bilayer edge crossings
	// by sweeping up and down the layers and ordering nodes based on their weighted medians.
	GraphvizDot
)
