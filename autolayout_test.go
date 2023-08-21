package autog

import (
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/phase4"
	"github.com/stretchr/testify/assert"
)

// These tests are meant to ensure a minimal degree of reliability of the default layout pipeline and
// a first barrier against regressions, pending proper unit tests
func TestLayoutCrashers(t *testing.T) {
	t.Run("phase4 SinkColoring", func(t *testing.T) {
		t.Run("#1 and #4", func(t *testing.T) {
			g := graph.FromEdgeSlice([][]string{
				{"N1", "N2"},
				{"N3", "N1"},
				{"N2", "N3"},
				{"Nh", "N1"},
				{"Nk", "N1"},
				{"Na", "N2"},
				{"Na", "N3"},
				{"N2", "Nd"},
			})
			assert.NotPanics(t, func() {
				g = Layout(g, WithPositioning(phase4.SinkColoring))
			})
			assertInvariants(t, g)
		})
	})
}

func assertInvariants(t *testing.T, g *graph.DGraph) {
	// todo
}
