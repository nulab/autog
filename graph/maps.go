package graph

type hashmap[K comparable, V any] map[K]V

type NodeIntMap = hashmap[*Node, int]

type EdgeIntMap = hashmap[*Edge, int]

type NodeFloatMap = hashmap[*Node, float64]

type NodeMap = hashmap[*Node, *Node]

type NodeSet = hashmap[*Node, bool]

type EdgeSet = hashmap[*Edge, bool]

func (m hashmap[K, V]) Clone() hashmap[K, V] {
	m2 := make(hashmap[K, V], len(m))
	for k, v := range m {
		m2[k] = v
	}
	return m2
}
