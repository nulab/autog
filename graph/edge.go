package graph

type Edge struct {
	// Identifier of this edge's source node
	FromID string
	// Identifier of this edge's target node
	ToID string
	// Slice of points required to draw this edge.
	// Depending on the routing algorithm, the slice can be interpreted differently:
	//
	// in case of edges made of straight segments (options EdgeRoutingStraight, EdgeRoutingPieceWise, EdgeRoutingOrtho),
	// the slice contains the start point, any number of intermediate bend points and the end point.
	// If the option EdgeRoutingStraight was chosen, the slice will have length 2 and include only the start and end points;
	//
	// in case of curved edges, the slice contains the control points of a piece-wise cubic bezier spline, i.e. its length
	// is a multiple of 4, and it must be processed in chunks of 4 points each.
	//
	// The Edge struct has no field that indicates how to interpret this slice. It's assumed that the caller code
	// knows what options autog.Layout ran with.
	Points [][2]float64

	// Whether the arrow head is placed at the start or end point of this edge.
	ArrowHeadStart bool
}
