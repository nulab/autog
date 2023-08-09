package main

import (
	"os"
	"strconv"
	"time"

	svg "github.com/ajstarks/svgo"
	"github.com/nulab/autog/graph"
)

const skipVirtual = true

func svgFile(g *graph.DGraph) {
	f, err := os.OpenFile("_tools/printshape/output.svg", os.O_RDWR|os.O_TRUNC, 0755)
	if err != nil {
		panic(err)
	}
	canvas := svg.New(f)

	var left, right, top, bottom float64
	for _, n := range g.Nodes {
		left = min(left, n.X)
		right = max(right, n.X+n.W)
		top = min(top, n.Y)
		bottom = max(bottom, n.Y+n.H)
	}
	const spacing = 40.0

	canvas.Start(int(right-left+spacing*2), int(bottom-top+spacing*4))
	canvas.Text(10, 10+spacing, time.Now().String(), "font-size:20px;fill:black")
	canvas.Def()
	canvas.Marker("arrowhead", 7, 2, 8, 8, `orient="auto-start-reverse"`)
	canvas.Path("M0,0 V4 L8,2 Z", "fill:black")
	canvas.MarkerEnd()
	canvas.DefEnd()

	for _, n := range g.Nodes {
		if skipVirtual && n.IsVirtual {
			continue
		}
		canvas.Rect(int(n.X)+spacing, int(n.Y)+spacing*2, int(n.W), int(n.H), "fill:none;stroke:black")

		text := n.ID + "-" + strconv.Itoa(n.LayerPos)
		canvas.Text(int(n.X)+int(n.W)/2+spacing, int(n.Y)+int(n.H)/2+spacing*2, text, "text-anchor:middle;font-size:30px;fill:black")
	}

	for _, e := range g.Edges {
		if len(e.Points) == 0 {
			continue
		}
		var xs, ys []int
		for _, p := range e.Points {
			xs = append(xs, int(p[0])+spacing)
			ys = append(ys, int(p[1])+spacing*2)
		}

		marker := "marker-end"
		if e.From.Layer > e.To.Layer {
			marker = "marker-start"
		}
		canvas.Polyline(xs, ys, "stroke-width:2;fill:none;stroke:black;"+marker+":url(#arrowhead)")
	}

	canvas.End()
}
