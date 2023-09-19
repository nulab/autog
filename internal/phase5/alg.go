package phase5

type Alg uint8

func (alg Alg) Phase() int {
	return 5
}

func (alg Alg) String() (s string) {
	switch alg {
	case NoRouting:
		s = "noop"
	case Straight:
		s = "straight"
	case PieceWise:
		s = "piecewise"
	case Ortho:
		s = "ortho"
	default:
		s = "<invalid>"
	}
	return s
}

const (
	NoRouting Alg = iota
	Straight
	PieceWise
	Ortho
	_endAlg
)
