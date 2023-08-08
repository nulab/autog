package main

import (
	"encoding/json"
	"fmt"
	"strconv"

	"github.com/nulab/autog/graph"
)

func cacooShapesJson(g *graph.DGraph) {
	clip := ClipboardShapes{
		Target:  "shapes",
		SheetId: "A97FC",
	}
	shapes := map[string]*Shape{}

	for _, n := range g.Nodes {
		shape := NewGroup(n.ID + "-L" + strconv.Itoa(n.Layer) + "-P" + strconv.Itoa(n.LayerPos))
		shape.Bounds.Top = n.Y
		shape.Bounds.Left = n.X
		shape.Bounds.Right = shape.Bounds.Left + n.W
		shape.Bounds.Bottom = shape.Bounds.Top + n.W
		shape.BuildConnectionPoints()
		if n.IsVirtual {
			shape.Shapes[0].LineInfo.Type = 3
		}

		clip.Shapes = append(clip.Shapes, shape)
		shapes[n.ID] = shape
	}

	for _, e := range g.Edges {
		line := NewLine()
		setLineProperties(line, e, shapes)
		clip.Shapes = append(clip.Shapes, line)
	}

	b, err := json.Marshal(clip)
	if err != nil {
		panic(err)
	}
	fmt.Println(string(b))
}
