package phase2

type Alg uint8

func (alg Alg) Phase() int {
	return 2
}

const (
	// LongestPath computes a partition of the graph in layers by traversing nodes in topological order.
	// It may result in more flat edges and comparatively more virtual nodes, therefore more long edges too, but runs in O(N).
	// Suitable for graphs with few "flow" paths.
	LongestPath Alg = iota

	// NetworkSimplex computes a partition of the graph in layers by minimizing total edge length.
	// It results in few virtual nodes and usually no flat edges, but runs in Θ(VE). Worst case seems to be O(V^2*E)
	NetworkSimplex
)

func (alg Alg) String() (s string) {
	switch alg {
	case LongestPath:
		s = "longest path"
	case NetworkSimplex:
		s = "network simplex"
	default:
		s = "<invalid>"
	}
	return
}
