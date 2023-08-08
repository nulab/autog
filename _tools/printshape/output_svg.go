package main

import (
	"os"
	"strconv"
	"time"

	svg "github.com/ajstarks/svgo"
	"github.com/nulab/autog/graph"
)

func svgFile(g *graph.DGraph) {
	f, err := os.OpenFile("_tools/printshape/output.svg", os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	canvas := svg.New(f)

	x, y := 0.0, 0.0
	var rx, ry *graph.Node
	for _, n := range g.Nodes {
		if n.X >= x {
			x = n.X
			rx = n
		}
		if n.Y >= y {
			y = n.Y
			ry = n
		}
	}
	const spacing = 40.0

	canvas.Start(int(rx.X+rx.W+spacing*2), int(ry.Y+ry.H+spacing*4))
	canvas.Text(10, 10+spacing, time.Now().String(), "font-size:20px;fill:black")
	// canvas.Def()
	// canvas.Marker("arrowhead", 3, 3, 6, 6, `orient="auto-start-reverse"`)
	// canvas.Path("M0,0 V6 L3,3 Z", "fill:black")
	// canvas.MarkerEnd()
	// canvas.DefEnd()

	for _, n := range g.Nodes {
		canvas.Rect(int(n.X)+spacing, int(n.Y)+spacing*2, int(n.W), int(n.H), "fill:none;stroke:black")

		text := n.ID + "-" + strconv.Itoa(n.LayerPos)
		canvas.Text(int(n.X)+int(n.W)/2+spacing, int(n.Y)+int(n.H)/2+spacing*2, text, "text-anchor:middle;font-size:30px;fill:black")
	}
	for _, e := range g.Edges {
		from, to := e.From, e.To

		// marker := "marker-end"
		if from.Layer > to.Layer {
			from, to = to, from
			//	marker = "marker-start"
		}

		startx := int(from.X + from.W/2 + spacing)
		starty := int(from.Y + from.H + spacing*2)
		endx := int(to.X + to.W/2 + spacing)
		endy := int(to.Y + spacing*2)

		canvas.Line(startx, starty, endx, endy, "stroke-width:2;fill:none;stroke:black") // ;"+marker+":url(#arrowhead)")
	}

	canvas.End()
}
