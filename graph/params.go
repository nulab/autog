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
	NetworkSimplexBalance bool

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

	// ---- other options ---

	// Algorithm monitor, receives traces and logs during algorithm execution.
	Monitor monitor.Monitor
}
