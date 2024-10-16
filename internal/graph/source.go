package graph

// Source represent the source of graph data. It hides the implementation details of the internal DGraph struct
// and allows only this module to provide implementations.
type Source interface {
	Populate(*DGraph)
}
