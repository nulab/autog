package graph

import "maps"

type set[T any] map[*T]bool

type NodeSet = set[Node]

type EdgeSet = set[Edge]

func (s set[T]) AsList() []*T {
	return maps.Keys(s)
}
