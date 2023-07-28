package graph

func (n *Node) VisitEdges(visit func(*Edge)) {
	fn, next := n.allEdges()
	for next {
		next = fn(visit)
	}
}

func (n *Node) allEdges() (visitor func(func(*Edge)) bool, next bool) {
	i := 0
	visitor = func(yield func(*Edge)) bool {
		if i >= len(n.In)+len(n.Out) {
			return false
		}
		if i < len(n.In) {
			yield(n.In[i])
		} else {
			yield(n.Out[i-len(n.In)])
		}
		i++
		return true
	}
	next = true
	return
}
