//go:build unit

package testfiles

import (
	"testing"

	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	"github.com/stretchr/testify/assert"
)

func TestNoRegression(t *testing.T) {
	t.Run("ELK", func(t *testing.T) {
		for _, testcase := range elkTestGraphs {
			t.Run(testcase.name, func(t *testing.T) {
				assert.NotPanics(t, func() { autog.Layout(graph.EdgeSlice(testcase.adj)) })
			})
		}
	})

	t.Run("Dot", func(t *testing.T) {
		assert.NotPanics(t, func() { autog.Layout(graph.EdgeSlice(DotAbstract)) })
	})
}
