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
	NoOrdering Alg = iota
	GraphvizDot
	_endAlg
)
