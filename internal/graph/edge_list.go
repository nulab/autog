package graph

type EdgeList []*Edge

func (list *EdgeList) Add(e *Edge) {
	if list == nil {
		return
	}
	*list = append(*list, e)
}

func (list *EdgeList) Remove(e *Edge) {
	if list == nil {
		return
	}
	// compact in place instead of re-slicing inside the range loop, which
	// would skip the element after each match and revisit a stale last element
	n := 0
	for _, f := range *list {
		if f != e {
			(*list)[n] = f
			n++
		}
	}
	*list = (*list)[:n]
}
