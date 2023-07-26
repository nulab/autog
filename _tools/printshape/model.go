package main

import "golang.nulab-inc.com/cacoo/library/common/v6/rand"

const (
	ShapeType_Line    = 0
	ShapeType_Polygon = 1
	ShapeType_Group   = 3
	ShapeType_Text    = 5
)

func cacooUid() string {
	return rand.String(10) + "-" + rand.String(10)
}

type ClipboardShapes struct {
	Target  string `json:"target"`
	SheetId string `json:"sheetId"`
	Shapes  []any  `json:"shapes"`
}

type Shape struct {
	Type             int                `json:"type"`
	Uid              string             `json:"uid"`
	Bounds           Bounds             `json:"bounds"`
	ConnectionPoints []*ConnectionPoint `json:"connectionPoints,omitempty"`
	Attr             []*Attr            `json:"attr,omitempty"`
	CategoryName     string             `json:"categoryName,omitempty"`
	Locked           bool               `json:"locked"`
	LineInfo         *LineInfo          `json:"lineInfo,omitempty"`
	DrawInfo         *DrawInfo          `json:"drawInfo,omitempty"`
	Paths            []*Path            `json:"paths,omitempty"`
	Text             string             `json:"text,omitempty"`
	Leading          int                `json:"leading,omitempty"`
	Halign           int                `json:"halign,omitempty"`
	Valign           int                `json:"valign,omitempty"`
	Styles           []*Style           `json:"styles,omitempty"`
	Shapes           []*Shape           `json:"shapes,omitempty"`
	Links            []*Link            `json:"-"`
}

func (s *Shape) BuildConnectionPoints() {
	h := s.Bounds.Bottom - s.Bounds.Top
	w := s.Bounds.Right - s.Bounds.Left
	tl := s.Bounds.Top
	tr := s.Bounds.Top + w
	bl := s.Bounds.Bottom
	br := s.Bounds.Bottom + w
	s.ConnectionPoints = []*ConnectionPoint{
		{tl, tl}, {tl, tl + w/4}, {tl, tl + w/2}, {tl, tl + w/4*3}, {tl, tr},
		{tl + h/4, tl}, {tl + h/4, tr},
		{tl + h/2, tl}, {tl + h/2, tr},
		{tl + h/4*3, tl}, {tl + h/4*3, tr},
		{bl, tl}, {bl, bl + w/4}, {bl, bl + w/2}, {bl, bl + w/4*3}, {bl, br},
	}
}

type Bounds struct {
	Top         float64 `json:"top"`
	Bottom      float64 `json:"bottom"`
	Left        float64 `json:"left"`
	Right       float64 `json:"right"`
	TopFixed    float64 `json:"topFixed,omitempty"`
	LeftFixed   float64 `json:"leftFixed,omitempty"`
	RightFixed  float64 `json:"rightFixed,omitempty"`
	BottomFixed float64 `json:"bottomFixed,omitempty"`
}

type DrawInfo struct {
	Enabled        bool            `json:"enabled"`
	FillRule       string          `json:"fillRule"`
	GradientColors []GradientColor `json:"gradientColors"`
}

type GradientColor struct {
	Color  string `json:"color"`
	Offset int    `json:"offset"`
}

type LineInfo struct {
	Enabled   bool   `json:"enabled"`
	Thickness int    `json:"thickness"`
	Color     string `json:"color"`
	Opacity   int    `json:"opacity"`
	Type      int    `json:"type"`
}

type Path struct {
	Closed bool        `json:"closed"`
	Points [][]float64 `json:"points"`
}

type ConnectionPoint [2]float64

type Style struct {
	Index           int    `json:"index"`
	Font            string `json:"font"`
	Size            int    `json:"size"`
	Color           string `json:"color"`
	Bold            bool   `json:"bold"`
	Italic          bool   `json:"italic"`
	Underline       bool   `json:"underline"`
	StrikeThrough   bool   `json:"strikeThrough"`
	BackgroundColor string `json:"backgroundColor"`
}

type Link struct{}

type Attr struct {
	Name  string `json:"name"`
	Value string `json:"value"`
}

type Line struct {
	Uid             string       `json:"uid"`
	Type            int          `json:"type"`
	Bounds          Bounds       `json:"bounds"`
	Style           LineStyle    `json:"style"`
	LineInfo        *LineInfo    `json:"lineInfo"`
	Points          []*LinePoint `json:"points,omitempty"`
	StartConnection string       `json:"startConnection,omitempty"` // uid or uid.index
	EndConnection   string       `json:"endConnection,omitempty"`
	StartArrow      *int         `json:"startArrow"`
	EndArrow        *int         `json:"endArrow"`
}

type LineStyle int

type LinePoint [2]float64
