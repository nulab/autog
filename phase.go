package autog

import "github.com/vibridi/autog/graph"

type phase interface {
	Process(*graph.Graph)
	Cleanup()
}
