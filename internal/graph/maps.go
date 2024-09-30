package graph

import "maps"

type hashmap[K comparable, V any] map[K]V

type NodeIntMap = hashmap[*Node, int]

type EdgeIntMap = hashmap[*Edge, int]

type NodeFloatMap = hashmap[*Node, float64]

type NodeMap = hashmap[*Node, *Node]

type NodeSet = hashmap[*Node, bool]

type EdgeSet = hashmap[*Edge, bool]

func (m hashmap[K, V]) Clone() hashmap[K, V] {
	return maps.Clone(m)
}

func (m hashmap[K, V]) Keys() []K {
	ks := make([]K, 0, len(m))
	for k := range m {
		ks = append(ks, k)
	}
	return ks
}
