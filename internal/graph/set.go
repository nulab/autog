package graph

import "golang.org/x/exp/maps"

type Set[T any] map[*T]bool

type NodeSet = Set[Node]

type EdgeSet = Set[Edge]

func (s Set[T]) AsList() []*T {
	return maps.Keys(s)
}
