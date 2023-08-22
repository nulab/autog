package autog

import (
	"testing"

	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/monitor"
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
		})
	})

	t.Run("phase3 GraphvizDot", func(t *testing.T) {
		t.Run("identical edge segfault in cross counting", func(t *testing.T) {
			g := graph.FromEdgeSlice([][]string{
				{"gql", "acc"},
				{"gql", "dia"},
				{"gql", "edt"},
				{"gql", "fld"},
				{"gql", "itg"},
				{"gql", "ntf"},
				{"gql", "org"},
				{"gql", "sub"},
				{"gql", "spt"},
				{"gql", "tmp"},
				{"acc", "lgc"},
				{"acc", "sub"},
				{"fld", "acc"},
				{"fld", "dia"},
				{"fld", "org"},
				{"fld", "sub"},
				{"dia", "acc"},
				{"dia", "fld"},
				{"dia", "lgc"},
				{"dia", "org"},
				{"dia", "sub"},
			})
			assert.NotPanics(t, func() {
				g = Layout(g)
			})
		})

		t.Run("wrong initialization of flat edges", func(t *testing.T) {
			g := graph.FromEdgeSlice([][]string{
				{"gql", "acc"},
				{"gql", "dia"},
				{"gql", "edt"},
				{"gql", "fld"},
				{"gql", "itg"},
				{"gql", "ntf"},
				{"gql", "org"},
				{"gql", "sub"},
				{"gql", "spt"},
				{"gql", "tmp"},
				{"acc", "lgc"},
				{"acc", "sub"},
				{"dia", "acc"},
				{"dia", "fld"},
				{"dia", "lgc"},
				{"dia", "org"},
				{"dia", "sub"},
				{"fld", "acc"},
				{"fld", "dia"},
				{"fld", "org"},
				{"fld", "sub"},
			})
			assert.NotPanics(t, func() {
				g = Layout(g, WithMonitor(monitor.NewStdout()))
			})
		})
	})
}
