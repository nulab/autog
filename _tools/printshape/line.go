package main

import (
	"github.com/nulab/autog/graph"
	"golang.nulab-inc.com/cacoo/library/common/v6/ptr"
)

func NewLine() *Line {
	return &Line{
		Uid:        cacooUid(),
		Type:       ShapeType_Line,
		StartArrow: ptr.Of(0),
		EndArrow:   ptr.Of(8),
		Bounds:     Bounds{},
		Style:      0,
		LineInfo: &LineInfo{
			Enabled:   true,
			Thickness: 1,
			Color:     "000000",
			Opacity:   1,
		},
		Points: nil,
	}
}

func setLineProperties(line *Line, e *graph.Edge, shapes map[string]*Shape) {
	if e.IsReversed {
		a := line.StartArrow
		line.StartArrow = line.EndArrow
		line.EndArrow = a
		line.LineInfo.Color = "d2302f"
		line.LineInfo.Thickness = 2
	}

	if e.From.Layer == e.To.Layer {
		sameLayerLine(line, e, shapes)
		return
	}

	from, to := shapes[e.From.ID], shapes[e.To.ID]
	line.StartConnection = from.Shapes[0].Uid + ".8"
	line.EndConnection = to.Shapes[0].Uid + ".7"

	if from.Bounds.Left > to.Bounds.Left {
		x := line.Bounds.Left
		line.Bounds.Left = line.Bounds.Right
		line.Bounds.Right = x
	}

	fromw := from.Bounds.Right - from.Bounds.Left
	tow := to.Bounds.Right - to.Bounds.Left

	startx, starty := from.Bounds.Left+fromw/2, from.Bounds.Bottom
	endx, endy := to.Bounds.Left+tow/2, to.Bounds.Top

	line.Bounds.Left = min(startx, endx)
	line.Bounds.Right = max(startx, endx)
	if from.Bounds.Left > to.Bounds.Left {
		line.Bounds.Right = max(startx, endx)
		line.Bounds.Left = min(startx, endx)
	}
	line.Bounds.Top = min(starty, endy)
	line.Bounds.Bottom = max(starty, endy)

	line.Points = append([]*LinePoint{}, &LinePoint{0, 0}, &LinePoint{line.Bounds.Right - line.Bounds.Left, line.Bounds.Bottom - line.Bounds.Top})
	if from.Bounds.Left > to.Bounds.Left {
		line.Points = append([]*LinePoint{}, &LinePoint{line.Bounds.Right - line.Bounds.Left, 0}, &LinePoint{0, line.Bounds.Bottom - line.Bounds.Top})
	}
}

func sameLayerLine(line *Line, e *graph.Edge, shapes map[string]*Shape) {
	from, to := shapes[e.From.ID], shapes[e.To.ID]

	suffFrom, suffTo := ".13", ".2"
	if from.Bounds.Left > to.Bounds.Left {
		suffFrom, suffTo = suffTo, suffFrom
	}

	line.StartConnection = from.Shapes[0].Uid + suffFrom
	line.EndConnection = to.Shapes[0].Uid + suffTo

	fromh := from.Bounds.Bottom - from.Bounds.Top
	toh := to.Bounds.Bottom - to.Bounds.Top

	startx, starty := min(from.Bounds.Right, to.Bounds.Right), from.Bounds.Top+fromh/2
	endx, endy := max(from.Bounds.Left, to.Bounds.Left), to.Bounds.Top+toh/2

	line.Bounds.Left = min(startx, endx)
	line.Bounds.Right = max(startx, endx)
	line.Bounds.Top = min(starty, endy)
	line.Bounds.Bottom = max(starty, endy)

	p1 := &LinePoint{0, 0}
	p2 := &LinePoint{line.Bounds.Right - line.Bounds.Left, 0}

	if from.Bounds.Left > to.Bounds.Left {
		line.Points = append([]*LinePoint{}, p2, p1)
	} else {
		line.Points = append([]*LinePoint{}, p1, p2)
	}

}
