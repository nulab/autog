package graph

type edgeIter struct {
	in, out []*Edge
	i       int
}

func (n *Node) EdgeIter() *edgeIter {
	return &edgeIter{in: n.In, out: n.Out}
}

func (itr *edgeIter) HasNext() bool {
	return itr.i < len(itr.in)+len(itr.out)
}

func (itr *edgeIter) Next() *Edge {
	if itr.i >= len(itr.in)+len(itr.out) {
		panic("edge iterator: out of bounds")
	}
	i := itr.i
	itr.i++
	if i < len(itr.in) {
		return itr.in[i]
	}
	return itr.out[i-len(itr.in)]
}

func (n *Node) AllEdges() (visitor func(func(*Edge)) bool, next bool) {
	i := 0
	visitor = func(f func(*Edge)) bool {
		if i >= len(n.In)+len(n.Out) {
			return false
		}
		if i < len(n.In) {
			f(n.In[i])
		} else {
			f(n.Out[i-len(n.In)])
		}
		i++
		return true
	}
	next = true
	return
}
