package graph

import "github.com/nulab/autog/monitor"

// Params holds parameters and options that are used by the layout algorithms
// and don't strictly belong to the graph itself
type Params struct {

	// ---- phase3 options ---

	// Maximum number of iterations of the GraphvizDot orderer.
	GraphvizDotMaxIter int

	// ---- other options ---

	// Algorithm monitor, receives traces and logs during algorithm execution.
	Monitor monitor.Monitor
}
