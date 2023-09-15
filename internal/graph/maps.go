package graph

type hashmap[K comparable, V any] map[K]V

type NodeIntMap = hashmap[*Node, int]

type EdgeIntMap = hashmap[*Edge, int]

type NodeFloatMap = hashmap[*Node, float64]

type NodeMap = hashmap[*Node, *Node]

type NodeSet = hashmap[*Node, bool]

type EdgeSet = hashmap[*Edge, bool]

type Layers = hashmap[int, *Layer]

func (m hashmap[K, V]) Clone() hashmap[K, V] {
	m2 := make(hashmap[K, V], len(m))
	for k, v := range m {
		m2[k] = v
	}
	return m2
}

func (m hashmap[K, V]) Keys() []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
