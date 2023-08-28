package phase4

type Alg uint8

func (alg Alg) Phase() int {
	return 4
}

func (alg Alg) String() (s string) {
	switch alg {
	case NoPositioning:
		s = "noop"
	case VerticalAlign:
		s = "vertical align"
	case BrandesKoepf:
		s = "b&k"
	case NetworkSimplex:
		s = "network simplex"
	case SinkColoring:
		s = "sink coloring"
	case PackRight:
		s = "pack right"
	default:
		s = "<invalid>"
	}
	return
}

const (
	// NoPositioning does nothing. Nodes won't be assigned any coordinates.
	NoPositioning Alg = iota

	// VerticalAlign aligns nodes in each layer vertically around the center of the diagram.
	// Works best for tree-like graphs with no back-edges.
	VerticalAlign

	// BrandesKoepf aligns nodes based on blocks and classes in O(V+E).
	// It results in a compact drawing but with less long straight edges.
	BrandesKoepf

	// NetworkSimplex sets X coordinates by constructing an auxiliary graph and solving it with the network simplex method.
	// Layers in the auxiliary graph are X coordinates in the main graph. Might be time-intensive for graphs above a few dozen nodes.
	NetworkSimplex

	// SinkColoring is a variant of BrandesKöpf that aligns nodes based on vertical blocks starting from the bottom.
	// It results in a larger drawing but with more long vertical edge paths. Runs in O(2kn) with 1 <= k <= maxshifts.
	SinkColoring

	// PackRight aligns nodes to the right.
	PackRight
)
