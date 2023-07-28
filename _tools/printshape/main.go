package main

import (
	"encoding/json"

	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/testfiles"
	"github.com/nulab/autog/positioning"
	"github.com/vibridi/cacooclip"
)

const exampleDiagram = "ci_router_ComplexRouter.json"

// todo: very ugly code used to set up a quick POC, eventually refactor this into its own project
func main() {
	elkg := testfiles.ReadTestFile("internal/testfiles/elk_relabeled", exampleDiagram)
	// elkg := testfiles.ReadTestFile("internal/testfiles/elk_constructed", "simple_long_edge.json")

	dg := graph.FromAdjacencyList(elkg.AdjacencyList())
	for _, n := range dg.Nodes {
		n.W = 100
		n.H = 100
	}

	clip := ClipboardShapes{
		Target:  "shapes",
		SheetId: "A97FC",
	}
	xoffset := 0.0
	shapes := map[string]*Shape{}
	for _, subg := range dg.ConnectedComponents() {
		autog.Layout(subg, autog.WithPositioning(positioning.NetworkSimplex))
		maxxoffset := 0.0
		for _, n := range subg.Nodes {
			shape := NewGroup(n.ID)
			shape.Bounds.Top = n.Y
			shape.Bounds.Left = n.X + xoffset
			shape.Bounds.Right = shape.Bounds.Left + 100.0
			shape.Bounds.Bottom = shape.Bounds.Top + 100.0
			shape.BuildConnectionPoints()
			if n.IsVirtual {
				shape.Shapes[0].LineInfo.Type = 3
			}

			clip.Shapes = append(clip.Shapes, shape)
			maxxoffset = max(maxxoffset, n.X)
			shapes[n.ID] = shape
		}
		xoffset += maxxoffset + 100.0 + 40.0 // 40 is the conn-comp distance

		for _, e := range subg.Edges {
			line := NewLine()
			setLineProperties(line, e, shapes)
			clip.Shapes = append(clip.Shapes, line)
		}
	}

	b, err := json.Marshal(clip)
	if err != nil {
		panic(err)
	}
	if err = cacooclip.Write(string(b)); err != nil {
		panic(err)
	}
}
