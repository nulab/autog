package autog

import (
	"github.com/nulab/autog/internal/phase1"
	"github.com/nulab/autog/internal/phase2"
	"github.com/nulab/autog/internal/phase3"
	"github.com/nulab/autog/internal/phase4"
	"github.com/nulab/autog/internal/phase5"
)

const (
	// CycleBreakingGreedy breaks cycles using the greedy heuristic by Eades et al.
	// It reverses fewer edges on average but is non-deterministic.
	CycleBreakingGreedy = phase1.Greedy

	// CycleBreakingDepthFirst breaks cycles with a classic DFS approach.
	// It's deterministic but may end up farther from the global optimum, i.e. reverse more edges than necessary.
	CycleBreakingDepthFirst = phase1.DepthFirst
)

const (
	// LayeringLongestPath computes a partition of the graph in layers by traversing nodes in topological order.
	// It may result in more flat edges and comparatively more virtual nodes, therefore more long edges, but runs in O(N).
	// Suitable for graphs with few "flow paths".
	LayeringLongestPath = phase2.LongestPath

	// LayeringNetworkSimplex computes a partition of the graph in layers by minimizing total edge length.
	// It results in few virtual nodes and usually no flat edges, but runs in Θ(VE). Worst case seems to be O(V^2*E).
	LayeringNetworkSimplex = phase2.NetworkSimplex
)

const (
	// OrderingNoop does not reorder nodes and edge crossings are not minimized.
	// It can be used as a no-op phase 3 for testing purposes.
	OrderingNoop = phase3.NoOrdering

	// OrderingWMedian implements a crossing minimization strategy based on Graphviz Dot's weighted median heuristic.
	// It sweeps layers two by two and reorders nodes based on the weighted median of their upper or lower neighbors.
	// It gives acceptable results in close to linear time.
	OrderingWMedian = phase3.GraphvizDot // todo: eventually resolve naming mismatch, even though the end user doesn't see it
)

const (
	// PositioningNoop does nothing. Nodes won't be assigned any coordinates.
	// It can be used as a no-op phase 4 for testing purposes.
	PositioningNoop = phase4.NoPositioning

	// PositioningVAlign aligns nodes in each layer vertically around the center of the diagram.
	// Works best for tree-like graphs with no back-edges.
	PositioningVAlign = phase4.VerticalAlign

	// PositioningNetworkSimplex sets X coordinates by constructing an auxiliary graph and layering it with the network simplex method.
	// Layers in the auxiliary graph are the X coordinates in the main graph.
	// Might be time-intensive for graphs above a few dozen nodes, as the current implementation is missing some of the optimizations
	// mentioned in Graphviz seminal paper.
	PositioningNetworkSimplex = phase4.NetworkSimplex

	// PositioningBrandesKoepf aligns nodes based on their partition in blocks and classes.
	// It runs in O(V+E) time. It results in a compact drawing but with less long straight edges.
	// The current implementation is not aware of node sizes.
	PositioningBrandesKoepf = phase4.BrandesKoepf

	// PositioningSinkColoring implements a variant of the Brandes-Köpf heuristic, aligning nodes bottom-up.
	// It results in a wider layout but with more long vertical edge paths. Runs in O(2kn) with 1 <= k <= maxshifts.
	PositioningSinkColoring = phase4.SinkColoring

	// PositioningPackRight aligns nodes to the right. Runs in linear time.
	PositioningPackRight = phase4.PackRight
)

const (
	// EdgeRoutingNoop outputs a layout without edges. It can be used as a no-op phase 5 for testing purposes.
	EdgeRoutingNoop = phase5.NoRouting

	// EdgeRoutingStraight computes the start and end point of each edge, thus drawing edges as straight lines.
	// Unsuitable for graphs with many long edges or flat edges between non-consecutive nodes,
	// as straight edges may end up crossing the nodes' interior.
	EdgeRoutingStraight = phase5.Straight

	// EdgeRoutingPieceWise outputs edges as polygonal paths.
	// Non-terminal points are located where the virtual nodes introduced in the layering phase would be.
	// Edges can be drawn with curved elbows if non-terminal points are interpreted as
	// the second control point of a quadratic bezier curve (P1),
	// however the rendering code must determine the location of P0 and P2 to ensure curve continuity.
	EdgeRoutingPieceWise = phase5.PieceWise // todo: rename to polyline

	// EdgeRoutingOrtho outputs edges as polygonal paths with orthogonal segments only, i.e. all edges bend at 90 degrees.
	// Dense graphs look tidier, but it's harder to understand where edges start and finish.
	// Suitable when there's few sets of edges with the same target node.
	EdgeRoutingOrtho = phase5.Ortho
)

func WithCycleBreaking(alg phase1.Alg) Option {
	return func(o *options) {
		o.p1 = alg
	}
}

func WithLayering(alg phase2.Alg) Option {
	return func(o *options) {
		o.p2 = alg
	}
}

func WithOrdering(alg phase3.Alg) Option {
	return func(o *options) {
		o.p3 = alg
	}
}

func WithPositioning(alg phase4.Alg) Option {
	return func(o *options) {
		o.p4 = alg
	}
}

func WithEdgeRouting(alg phase5.Alg) Option {
	return func(o *options) {
		o.p5 = alg
	}
}
