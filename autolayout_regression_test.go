package autog

import (
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/stretchr/testify/assert"
)

func TestLayoutNoRegression(t *testing.T) {
	t.Run("ELK lib_decg_DECGPi", func(t *testing.T) {
		g := graph.FromEdgeSlice([][]string{
			{"N2", "N8"},
			{"N2", "N13"},
			{"N2", "N15"},
			{"N2", "N4"},
			{"N3", "N1"},
			{"N4", "N3"},
			{"N5", "N16"},
			{"N5", "N18"},
			{"N6", "N8"},
			{"N6", "N18"},
			{"N7", "N6"},
			{"N8", "N5"},
			{"N8", "N9"},
			{"N9", "N6"},
			{"N9", "N7"},
			{"N10", "N14"},
			{"N10", "N19"},
			{"N11", "N10"},
			{"N12", "N10"},
			{"N12", "N11"},
			{"N13", "N14"},
			{"N14", "N17"},
			{"N14", "N12"},
			{"N15", "N13"},
			{"N16", "N4"},
			{"N17", "N16"},
			{"N17", "N19"},
			{"N18", "N5"},
			{"N19", "N17"},
		})
		assert.NotPanics(t, func() {
			g = Layout(g)
		})
	})

	t.Run("ELK pn_brockackerman_BrockAckerman", func(t *testing.T) {
		g := graph.FromEdgeSlice([][]string{
			{"N1", "N2"},
			{"N2", "N9"},
			{"N4", "N10"},
			{"N4", "N15"},
			{"N5", "N6"},
			{"N6", "N14"},
			{"N8", "N1"},
			{"N8", "N3"},
			{"N9", "N12"},
			{"N10", "N11"},
			{"N11", "N12"},
			{"N12", "N8"},
			{"N13", "N5"},
			{"N13", "N7"},
			{"N14", "N16"},
			{"N15", "N16"},
			{"N16", "N13"},
		})
		assert.NotPanics(t, func() {
			g = Layout(g)
		})
	})
}
