package phase2

type Alg uint8

// Phase returns the ordinal number of this phase: 2.
func (alg Alg) Phase() int {
	return 2
}

// String returns a mnemonic representation of this algorithm.
// The exact string values are not documented and may change in the future.
func (alg Alg) String() (s string) {
	switch alg {
	case LongestPath:
		s = "longestpath"
	case NetworkSimplex:
		s = "ns"
	default:
		s = "<invalid>"
	}
	return
}

const (
	LongestPath Alg = iota
	NetworkSimplex
	_endAlg
)
