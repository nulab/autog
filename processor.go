package autog

import "github.com/nulab/autog/graph"

type processor interface {
	Process(*graph.DGraph)
}
