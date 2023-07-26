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

func (p *brandesKoepfPositioner) markConflicts(g *graph.DGraph) {
	L := len(g.Layers)
	if L < 3 {
		return
	}
	layerSize := map[int]int{}
	for _, n := range g.Nodes {
		layerSize[n.Layer]++
	}

	U := 2
	for i := 1; i < L-1; i++ {
		layer := g.Layers[U]

		k0 := 0
		l := 0

		for l1 := 0; l1 < layerSize[i+1]; l1++ {
			vli := layer[l1]
			if (l1 == layerSize[i+1]-1) || incidentToInner(vli, i+1, i) {
				k1 := layerSize[i] - 1
				if incidentToInner(vli, i+1, i) {
					k1 = -1 // todo
				}
				for l <= l1 {
					vl := layer[l]
					if !incidentToInner(vl, i+1, i) {
						for true {
							k := -20 // todo
							if k < k0 || k > k1 {
								// Marked edge can't return null here, because the upper neighbor
								// relationship between v_l and upperNeighbor enforces the existence
								// of at least one edge between the two nodes
								p.markedEdges = append(p.markedEdges /*upperNeighbor.getSecond()*/, nil)
							}
						}
					}

					l++
				}
				k0 = k1
			}
		}
	}

	//
	//        // The following call succeeds since there are at least 3 layers in the graph
	//        Iterator<Layer> layerIterator = layeredGraph.getLayers().listIterator(2);
	//        for (int i = 1; i < numberOfLayers - 1; i++) {
	//            // The variable naming here follows the notation of the corresponding paper
	//            // Normally, underscores are not allowed in local variable names, but since there
	//            // is no way of properly writing indices beside underscores, Checkstyle will be
	//            // disabled here and in future methods containing indexed variables
	//            // CHECKSTYLEOFF Local Variable Names
	//            Layer currentLayer = layerIterator.next();
	//            Iterator<LNode> nodeIterator = currentLayer.getNodes().iterator();
	//
	//            int k_0 = 0;
	//            int l = 0;
	//
	//            for (int l_1 = 0; l_1 < layerSize[i + 1]; l_1++) {
	//                // In the paper, l and i are indices for the layer and the position in the layer
	//                LNode v_l_i = nodeIterator.next(); // currentLayer.getNodes().get(l_1);
	//
	//                if (l_1 == ((layerSize[i + 1]) - 1) || incidentToInnerSegment(v_l_i, i + 1, i)) {
	//                    int k_1 = layerSize[i] - 1;
	//                    if (incidentToInnerSegment(v_l_i, i + 1, i)) {
	//                        k_1 = ni.nodeIndex[ni.leftNeighbors.get(v_l_i.id).get(0).getFirst().id];
	//                    }
	//
	//                    while (l <= l_1) {
	//                        LNode v_l = currentLayer.getNodes().get(l);
	//
	//                        if (!incidentToInnerSegment(v_l, i + 1, i)) {
	//                            for (Pair<LNode, LEdge> upperNeighbor : ni.leftNeighbors.get(v_l.id)) {
	//                                int k = ni.nodeIndex[upperNeighbor.getFirst().id];
	//
	//                                if (k < k_0 || k > k_1) {
	//                                    // Marked edge can't return null here, because the upper neighbor
	//                                    // relationship between v_l and upperNeighbor enforces the existence
	//                                    // of at least one edge between the two nodes
	//                                    markedEdges.add(upperNeighbor.getSecond());
	//                                }
	//                            }
	//                        }
	//
	//                        l++;
	//                    }
	//
	//                    k_0 = k_1;
	//                }
	//            }
	//            // CHECKSTYLEON Local Variable Names
	//        }
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

func incidentToInner(n *graph.Node, i, j int) bool {
	return false
}
