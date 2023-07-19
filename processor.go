package autog

import "github.com/vibridi/autog/graph"

type processor interface {
	Process(*graph.DGraph)
}
