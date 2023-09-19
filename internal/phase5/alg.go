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
	case Polyline:
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
	Polyline
	Ortho
	_endAlg
)
