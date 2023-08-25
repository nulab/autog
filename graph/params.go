package graph

import "github.com/nulab/autog/monitor"

// Params holds parameters and options that are used by the layout algorithms
// and don't strictly belong to the graph itself
type Params struct {

	// ---- phase2 options ---

	// Factor used in to determine the maximum number of iterations.
	NetworkSimplexThoroughness uint
	// If positive, factor by which thoroughness is multiplied to determine the maximum number of iterations.
	// Otherwise, ignored.
	NetworkSimplexMaxIterFactor int
	// If true, balances the network simplex layering by moving nodes to less crowded layers.
	NetworkSimplexBalance int // todo: make enum

	// ---- phase3 options ---

	// Maximum number of iterations of the GraphvizDot orderer.
	GraphvizDotMaxIter uint

	// ---- phase4 options ---

	// Spacing between layers (above and below).
	LayerSpacing float64
	// Spacing between nodes (left and right).
	NodeSpacing float64
	// Weight factor for edges in the network simplex positioner.
	NetworkSimplexAuxiliaryGraphWeightFactor int
	// Allows choosing one of the four B&K layouts. The accepted values are: 0: bottom-right, 1: bottom-left, 2: top-right, 3: top-left.
	// The directions refer to the direction in which the algorithm sweeps layers and nodes. Different directions
	// result in different aligmnent priorities, and therefore in different positioning.
	// In case the inequality 0 <= v < 4 doesn't hold, the default balanced layout is used instead.
	BrandesKoepfLayout int

	// ---- other options ---

	// Algorithm monitor, receives traces and logs during algorithm execution.
	Monitor monitor.Monitor
}
