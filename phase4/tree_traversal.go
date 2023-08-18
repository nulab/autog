package phase4

import "github.com/nulab/autog/graph"

func execTreeTraversal(g *graph.DGraph, params graph.Params) {
	visited := graph.NodeSet{}
	layercx := map[int]float64{}

	for _, n := range g.Sources() {
		depthFirstPosition(n, visited, layercx)
	}
	for _, n := range g.Nodes {
		depthFirstPosition(n, visited, layercx)
	}

	for _, l := range g.Layers {
		for _, n := range l.Nodes {
			l.H = max(l.H, n.H)
		}
	}
}

func depthFirstPosition(n *graph.Node, visited graph.NodeSet, layercx map[int]float64) float64 {
	if visited[n] {
		return n.X + n.W/2
	}
	visited[n] = true

	mx := make([]float64, 0, len(n.Out))
	for _, e := range n.Out {
		if e.To.Layer <= e.From.Layer {
			continue
		}
		m := depthFirstPosition(e.ConnectedNode(n), visited, layercx)
		mx = append(mx, m)
	}
	if len(mx) == 0 {
		n.X += layercx[n.Layer]
		goto space
	}

	if len(mx)%2 != 0 {
		n.X = mx[len(mx)/2] - n.W/2
	} else {
		n.X = (mx[len(mx)/2-1]+mx[len(mx)/2])/2 - n.W/2
	}
space:
	layercx[n.Layer] += n.X + 20 + 40
	return n.X + n.W/2
}
