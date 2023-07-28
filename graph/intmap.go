package graph

type intmap[T any] map[*T]int

type NodeMap = intmap[Node]

type EdgeMap = intmap[Edge]

func (m intmap[T]) Clone() intmap[T] {
	m2 := intmap[T]{}
	for k, v := range m {
		m2[k] = v
	}
	return m2
}
