package phase3

type Alg uint8

func (alg Alg) Phase() int {
	return 3
}

const (
	// NoOrdering does nothing. Nodes won't be reordered in their layers to minimize edge crossings.
	NoOrdering Alg = iota

	// GraphvizDot implements the mincross heuristic used in dot. It attempts to minimize bilayer edge crossings
	// by sweeping up and down the layers and ordering nodes based on their weighted medians.
	GraphvizDot
)

func (alg Alg) String() (s string) {
	switch alg {
	case NoOrdering:
		s = "noop"
	case GraphvizDot:
		s = "graphviz dot"
	default:
		s = "<invalid>"
	}
	return
}
