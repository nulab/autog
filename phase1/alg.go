package phase1

type Alg uint8

func (alg Alg) Phase() int {
	return 1
}

const (
	Greedy Alg = iota // todo: document that this is non-deterministic
	DepthFirst
)

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
