package phase1

type Alg uint8

// Phase returns the ordinal number of this phase: 1.
func (alg Alg) Phase() int {
	return 1
}

// String returns a mnemonic representation of this algorithm.
// The exact string values are not documented and may change in the future.
func (alg Alg) String() (s string) {
	switch alg {
	case Greedy:
		s = "greedy"
	case DepthFirst:
		s = "dfs"
	default:
		s = "<invalid>"
	}
	return
}

const (
	// Greedy implements a solution to the feedback arc set problem using a greedy heuristic. It is non-deterministic.
	Greedy Alg = iota

	// DepthFirst removes cycles using the classical DFS strategy. It gives acceptable results in a short time, but the feedback arc set
	// may be larger compared to the greedy heuristic. It is deterministic.
	DepthFirst

	_endAlg
)
