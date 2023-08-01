package positioning

import (
	"unsafe"

	"github.com/nulab/autog/graph"
)

// todo: maybe this will become Rüegg-Schulze or BrandesKoepfExtended instead, which accounts for node sizes and ports.
// Rueegg-Schulze developed the algo for arbitrary port positioning.
// If ports aren't relevant to a particular implementation, node size still is, so the port can be set by default
// at the middle point of the node side.

type brandesKoepfPositioner struct {
	markedEdges []*graph.Edge
}

type dir uint8

const (
	top dir = iota
	bottom
	left
	right
)

type layout struct {
	v, h dir
}

func execBrandesKoepf(g *graph.DGraph) {

	layouts := [4]layout{
		{top, left},
		{top, right},
		{bottom, left},
		{bottom, right},
	}
	for _, _ = range layouts {
		verticalAlignment()
		horizontalCompaction()
	}
}

func verticalAlignment() {

}

func horizontalCompaction() {

}

// marks edges that cross inner edges, i.e. type 1 and type 2 conflicts as defined in B&K
func markConflicts(g *graph.DGraph) {
	if len(g.Layers) < 4 {
		return
	}
	marked := []*graph.Edge{}
	// sweep layers from top to bottom except the first and the last
	for i := 1; i < len(g.Layers)-1; i++ {
		k0 := 0
		for l1, v := range g.Layers[i+1].Nodes {
			ksrc := incidentToInner(v)
			if g.Layers[i+1].Tail() == v || ksrc >= 0 {
				// set k1 to the index of the last node or the index of the inner edge's source
				// if v is the tail node, ksrc is either identical to v.LayerPos or negative
				// if v belongs to an inner edge, ksrc is non-negative and identical to v.LayerPost
				// therefore max() returns the correct value for k1
				k1 := max(v.LayerPos, ksrc)
				// range over same layer nodes until v included
				for l2, w := range g.Layers[i+1].Nodes {
					if l2 > l1 {
						break
					}
					for _, e := range w.In {
						if e.SelfLoops() || e.IsFlat() {
							continue
						}
						// greater than k1 captures that e crosses an inner edge from left to right:
						// 	- k1 is set to the position of v's source;
						// 	- the strict inequality prevents erroneously marking e as conflicting with itself;
						// 	- there is a conflict even if v is the target node of both e and the inner edge
						// lesser than k0 captures that e crosses an inner edge from right to left:
						// 	- k0 is updated only after finding an inner edge with target v
						//	- at the next iteration moving to v+1, there is a crossing if the v+1 source precedes v's source
						if e.From.LayerPos < k0 || e.From.LayerPos > k1 {
							marked = append(marked, e)
						}
					}
				}
				k0 = k1
			}
		}
	}
}

type block = []*graph.Edge

var inn = map[*graph.Node]float64{}
var blox = [][]*graph.Edge{}
var blockSize = map[*block]float64{}

func innerShift(nodes []*graph.Node) {
	for _, n := range nodes {
		inn[n] = 0
		for _, b := range blox {
			left, right := 0.0, 0.0
			for _, e := range b {
				p, q := e.From, e.To
				s := inn[π(p)] + xp(p) - xp(q)
				inn[π(q)] = s
				left = min(left, s)
				right = max(right, s+width(π(q)))
			}
			for _, e := range blox {
				n := (*graph.Node)(unsafe.Pointer(e[0]))
				inn[n] -= left
			}
			blockSize[&b] = right - left
		}
	}
}

type port = *graph.Node // todo

// maps port to node
func π(port) *graph.Node {
	return nil
}

func xp(port) float64 {
	return 0.0
}

func width(port) float64 {
	return 0
}

// returns a non-negative number if n is the target node of an inner edge, i.e. an edge connecting two virtual nodes
// on adjacent layers, where the number is the position of the edge source in the upper layer;
// it returns -1 if the node isn't involved in an inner edge.
func incidentToInner(n *graph.Node) int {
	if !n.IsVirtual {
		return -1
	}
	for _, e := range n.In {
		if e.From.IsVirtual && e.From.Layer == n.Layer-1 {
			return e.From.LayerPos
		}
	}
	return -1
}
