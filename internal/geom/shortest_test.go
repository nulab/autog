package geom

import (
	"fmt"
	"slices"
	"strconv"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestShortest(t *testing.T) {
	rects := []Rect{
		{P{406.00, 1.00}, P{464.00, 35.00}},
		{P{110.00, 35.00}, P{437.00, 69.00}},
		{P{201.00, 69.00}, P{331.00, 103.00}},
		{P{188.00, 103.00}, P{295.00, 137.00}},
		{P{179.00, 137.00}, P{232.00, 171.00}},
		{P{193.00, 171.00}, P{429.00, 205.00}},
		{P{11.00, 205.00}, P{298.00, 239.00}},
		{P{170.00, 239.00}, P{276.00, 273.00}},
		{P{219.00, 273.00}, P{342.00, 307.00}},
		{P{324.00, 307.00}, P{375.00, 341.00}},
		{P{219.00, 341.00}, P{361.00, 375.00}},
		{P{226.00, 375.00}, P{247.00, 409.00}},
		{P{157.00, 409.00}, P{245.00, 443.00}},
		{P{3.00, 443.00}, P{181.00, 477.00}},
		{P{139.00, 477.00}, P{392.00, 511.00}},
	}
	re := len(rects) - 1
	start := P{rects[0].TL.X + rects[0].Width()/2, rects[0].TL.Y}
	end := P{rects[re].TL.X + rects[re].Width()/4, rects[re].TL.Y + rects[re].Height()/2}

	path := Shortest(start, end, rects)

	assert.Len(t, path, 13)

	// want points listed top-down
	want := []P{
		{435.00, 1.00},
		{406.00, 35.00},
		{331.00, 69.00},
		{232.00, 137.00},
		{232.00, 171.00},
		{276.00, 273.00},
		{324.00, 307.00},
		{324.00, 341.00},
		{247.00, 375.00},
		{226.00, 409.00},
		{181.00, 443.00},
		{181.00, 477.00},
		{202.25, 494.00},
	}
	slices.Reverse(path)
	assertPath(t, want, path)

	end = P{42.25, 464.00}
	path = Shortest(start, end, rects)
	slices.Reverse(path)
	want2 := append(want[:len(want)-4], P{226.00, 409.00}, P{157.00, 443.00}, end)
	assertPath(t, want2, path)

	end = P{372.25, 334.00}
	path = Shortest(start, end, rects)
	slices.Reverse(path)
	want3 := append(want[:6], end)
	assertPath(t, want3, path)

	t.Run("same triangle", func(t *testing.T) {
		end = P{start.X, start.Y + 20}
		path = Shortest(start, end, rects)
		slices.Reverse(path)
		assertPath(t, []P{start, end}, path)
	})

	t.Run("different triangles but straight line", func(t *testing.T) {
		end = P{405.00, 61.00}
		path = Shortest(start, end, rects)
		slices.Reverse(path)
		assertPath(t, []P{start, end}, path)
	})
}

func TestShortestEdgeCases(t *testing.T) {
	rects := []Rect{
		{P{112, 90}, P{200, 140}},
		{P{80, 140}, P{150, 300}},
		{P{140, 300}, P{270, 380}},
	}
	start := P{190, 140 - 1}
	end := P{200, 300 + 1}

	path := Shortest(start, end, rects)

	printall(rects, start, end, path)
}

func assertPath(t *testing.T, want, got []P) {
	require.Equal(t, len(want), len(got))
	for i := 0; i < len(got); i++ {
		if i > 0 {
			assert.True(t, got[i-1].Y < got[i].Y)
		}
		assert.Equal(t, want[i], got[i])
	}
}

func printpath(path []P) {
	for i := 1; i < len(path); i++ {
		u, v := path[i-1], path[i]
		fmt.Printf(`<path d="M %.2f,%.2f %.2f,%.2f" stroke="black" stroke-width="3" />`+"\n", u.X, u.Y, v.X, v.Y)
	}
}

func printall(rects []Rect, start, end P, path []P) {
	p := MergeRects(rects)

	s := polyline(p.Points, "red")
	fmt.Println(s)

	fmt.Println(start.String())
	fmt.Println(end.String())
	printpath(path)
}

func polyline(points []P, color string) string {
	b := strings.Builder{}
	b.WriteString(`<polyline points="`)
	for _, p := range points {
		b.WriteString(strconv.FormatFloat(p.X, 'f', 2, 64))
		b.WriteRune(',')
		b.WriteString(strconv.FormatFloat(p.Y, 'f', 2, 64))
		b.WriteRune(' ')
	}
	b.WriteString(`" fill="none" stroke="` + color + `" />`)
	return b.String()
}
