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
	Greedy Alg = iota
	DepthFirst
	_endAlg
)
