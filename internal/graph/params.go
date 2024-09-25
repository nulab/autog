package graph

// OptionNsBalance controls which balancing strategy to use in the network simplex layerer.
type OptionNsBalance uint8

const (
	// OptionNsBalanceV represents vertical balancing in network simplex solver. Default value used in phase 2.
	OptionNsBalanceV OptionNsBalance = iota + 1
	// OptionNsBalanceH represents horizontal balancing in network simplex solver. Used in phase 4 NetworkSimplex positioner.
	OptionNsBalanceH
)

// Params holds parameters and options that are used by the layout algorithms
// and don't strictly belong to the graph itself
type Params struct {

	// Sets the same width and height to all non-virtual nodes
	NodeFixedSizeFunc func(n *Node)

	// Sets a width and height to individual non-virtual nodes
	NodeSizeFunc func(n *Node)

	// ---- phase2 options ---

	// Factor used in to determine the maximum number of iterations.
	NetworkSimplexThoroughness uint
	// If positive, factor by which thoroughness is multiplied to determine the maximum number of iterations.
	// Otherwise, ignored.
	NetworkSimplexMaxIterFactor int
	// Controls which balancing strategy to use in the network simplex layering by moving nodes to less crowded layers.
	NetworkSimplexBalance OptionNsBalance

	// ---- phase3 options ---

	// Maximum number of iterations of the WMedian orderer.
	WMedianMaxIter uint

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
}
