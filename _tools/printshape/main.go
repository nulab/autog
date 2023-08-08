package main

import (
	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/testfiles"
	positioning "github.com/nulab/autog/phase4"
)

const exampleDiagram = "lib_decg_DECGPi.json" //  "ci_router_ComplexRouter.json" //  //

// todo: very ugly code used to set up a quick POC, eventually refactor this into its own project
func main() {
	elkg := testfiles.ReadTestFile("internal/testfiles/elk_relabeled", exampleDiagram)
	// elkg := testfiles.ReadTestFile("internal/testfiles/elk_constructed", "simple_long_edge.json")

	g := graph.FromElk(elkg)
	for _, n := range g.Nodes {
		n.W = 100
		n.H = 100
		if n.ID == "N13" || n.ID == "N4" {
			n.W = 150
		}
	}

	autog.Layout(g, autog.WithPositioning(positioning.SinkColoring))
	// cacooShapesJson(g)
	svgFile(g)

}
