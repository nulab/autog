package phase4

type Alg uint8

// Phase returns the ordinal number of this phase: 4.
func (alg Alg) Phase() int {
	return 4
}

// String returns a mnemonic representation of this algorithm.
// The exact string values are not documented and may change in the future.
func (alg Alg) String() (s string) {
	switch alg {
	case NoPositioning:
		s = "noop"
	case VerticalAlign:
		s = "vertical"
	case BrandesKoepf:
		s = "b&k"
	case NetworkSimplex:
		s = "ns"
	case SinkColoring:
		s = "sinkcoloring"
	case PackRight:
		s = "packright"
	default:
		s = "<invalid>"
	}
	return
}

const (
	NoPositioning Alg = iota
	VerticalAlign
	BrandesKoepf
	NetworkSimplex
	SinkColoring
	PackRight
	_endAlg
)
