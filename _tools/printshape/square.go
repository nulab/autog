package main

// this should be "new square"
func NewGroup(text string) *Shape {
	return &Shape{
		Type: ShapeType_Group,
		Bounds: Bounds{
			Top:    -200,
			Bottom: -100,
			Left:   -250,
			Right:  -150,
		},
		Shapes: []*Shape{NewSquare(), NewText(text)},
	}
}

func NewSquare() *Shape {
	return &Shape{
		Type: ShapeType_Polygon,
		Uid:  cacooUid(),
		Bounds: Bounds{
			Top:    0,
			Bottom: 100,
			Left:   0,
			Right:  100,
		},
		ConnectionPoints: []*ConnectionPoint{
			{0, 0}, {0, 25}, {0, 50}, {0, 75}, {0, 100},
			{25, 0}, {25, 100},
			{50, 0}, {50, 100},
			{75, 0}, {75, 100},
			{100, 0}, {100, 25}, {100, 50}, {100, 75}, {100, 100},
		},
		Attr:         []*Attr{{"stencil-id", "250F0"}},
		CategoryName: "basic",
		Locked:       false,
		LineInfo: &LineInfo{
			Enabled:   true,
			Thickness: 2,
			Color:     "4A4A4A",
			Opacity:   1,
		},
		DrawInfo: &DrawInfo{
			Enabled:        true,
			FillRule:       "nonzero",
			GradientColors: []GradientColor{{"FFFFFF", 0}, {"000000", 100}},
		},
		Paths: []*Path{
			{
				Closed: true,
				Points: [][]float64{{100, 0}, {0, 0}, {0, 100}, {100, 100}, {100, 0}},
			},
		},
	}
}

func NewText(text string) *Shape {
	return &Shape{
		Type: ShapeType_Text,
		Uid:  cacooUid(),
		Bounds: Bounds{
			Top:         1,
			Bottom:      -2,
			Left:        1,
			Right:       -2,
			TopFixed:    0,
			LeftFixed:   2,
			RightFixed:  3,
			BottomFixed: 1,
		},
		Text:    text,
		Leading: 2,
		Halign:  1,
		Valign:  1,
		Styles: []*Style{
			{
				Index: 0,
				Font:  "Open Sans",
				Size:  15,
				Color: "000000",
			},
		},
	}
}
