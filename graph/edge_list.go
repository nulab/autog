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
	for i, f := range *list {
		if f == e {
			*list = append((*list)[:i], (*list)[i+1:]...)
		}
	}
}
