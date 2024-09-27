package geom

import (
	"slices"

	"github.com/nulab/autog/internal/collectors"
)

// Shortest finds the shortest path between p1 and p2 through the inner polygon identified by the given list of rectangles.
func Shortest(p1, p2 P, rects []Rect) []P {
	// triangulate the polygon that would result from merging the input rectangles
	// this allows to special-case the triangulation and complete in O(n)
	// even though the polygon isn't strictly monotone.
	ts := Triangulate(rects)

	tmap := make(map[int]Tri, len(ts))
	// find start and end triangles
	var start, end Tri
	for _, t := range ts {
		if t.Contains(p1) {
			start = t
		}
		if t.Contains(p2) {
			end = t
		}
		tmap[t.ID] = t
	}
	if start.ID == end.ID {
		return []P{p2, p1}
	}

	// link triangles by their common side and build the dual graph as an adjacency matrix
	// pointers to common diagonals are stored at position (i,j) in the matrix
	adj := dualGraph(start, ts)

	// find list of diagonals that the path has to cross
	dlist := crossedDiagonals(start.ID, end.ID, adj, map[int]bool{})

	// append the last diagonal that has p2 as endpoint, it doesn't matter by which endpoint it's connected
	dlist = append(dlist, &Segment{dlist[len(dlist)-1].A, p2})

	predecessor := map[P]P{}

	deq := collectors.NewDeque[P](len(rects) * 2)
	deq.PushFront(p1)

	apex := deq.Front()

	if orientation(p1, dlist[0].A, dlist[0].B) == ccw {
		deq.PushFront(dlist[0].A)
		deq.PushBack(dlist[0].B)
	} else {
		deq.PushFront(dlist[0].B)
		deq.PushBack(dlist[0].A)
	}

	// helper functions to test whether the point v is outside the funnel
	outsideLeft := func(v P) bool {
		if deq.Len() < 2 {
			return true
		}
		d := orientation(deq.PeekFront(2), deq.PeekFront(1), v)
		return (deq.Front() < apex && d != ccw) || (deq.Front() >= apex && d != cw)
	}
	outsideRight := func(v P) bool {
		if deq.Len() < 2 {
			return true
		}
		d := orientation(deq.PeekBack(2), deq.PeekBack(1), v)
		return (deq.Back() > apex && d != cw) || (deq.Back() <= apex && d != ccw)
	}

	// funnel loop
	for i := 1; i < len(dlist); i++ {
		c := commonVertex(dlist[i-1], dlist[i])
		switch {
		case deq.PeekBack(1) == c:
			// left chain
			v := dlist[i].Other(c)

			out := outsideLeft(v)
			// if true, point already to the left of the current chain, widen funnel
			// otherwise shrink the funnel until v is inside the wedge formed by projecting the last two chain segments onto dlist[i]
			for !out {
				deq.PopFront()
				out = outsideLeft(v)
			}
			if deq.Front() > apex {
				apex = deq.Front()
			}
			predecessor[v] = deq.PeekFront(1)
			deq.PushFront(v)

		case deq.PeekFront(1) == c:
			// right chain
			v := dlist[i].Other(c)

			out := outsideRight(v)
			for !out {
				deq.PopBack()
				out = outsideRight(v)
			}
			if deq.Back() < apex {
				apex = deq.Back()
			}
			predecessor[v] = deq.PeekBack(1)
			deq.PushBack(v)

		default:
			panic("shortest path: funnels: disconnected triangulation diagonal")
		}
	}

	// compose the shortest path by walking up the predecessors chain until p1
	path := []P{}
	u, ok := p2, true
	for ok {
		path = append(path, u)
		u, ok = predecessor[u]
	}
	if path[len(path)-1] != p1 {
		path = append(path, p1)
	}
	return slices.Clip(path)
}

// builds the dual graph of the triangulation and returns it as an adjacency matrix, common diagonals are
// stored in matrix cells.
// In other words, when mat[i][j] != nil then there is an edge between Tri.ID = i and Tri.ID = j, and that edge crosses
// the diagonal stored in mat[i][j]
func dualGraph(start Tri, ts []Tri) collectors.Mat[*Segment] {
	pmap := map[Segment]int{}
	pmap[start.OrderedSide(0) /* AB */] = start.ID

	// build the adjacency matrix
	adj := collectors.NewMat[*Segment](len(ts) + 1)

	for _, t := range ts {
		for i := 0; i < 3; i++ {
			side := t.OrderedSide(i)
			id, ok := pmap[side]
			if !ok {
				pmap[side] = t.ID
			} else {
				adj[id][t.ID] = &side
				adj[t.ID][id] = &side
			}
		}
	}
	return adj
}

// finds the list of diagonals that the shortest path from Tri.ID = startid and Tri.ID = endid has to cross.
// The argument adj is the adjacency matrix of the dual graph of the triangulation (see also dualGraph).
func crossedDiagonals(startid, endid int, adj collectors.Mat[*Segment], visited map[int]bool) []*Segment {
	if startid == endid {
		return []*Segment{}
	}
	visited[startid] = true
	for tid, dg := range adj[startid] {
		if dg != nil && !visited[tid] {
			out := crossedDiagonals(tid, endid, adj, visited)
			if out != nil {
				return append([]*Segment{dg}, out...)
			}
		}
	}
	return nil
}

// assuming that d1 and d2 share a vertex, return that vertex
func commonVertex(d1, d2 *Segment) P {
	if d1.A == d2.A || d1.B == d2.A {
		return d2.A
	}
	return d2.B
}
