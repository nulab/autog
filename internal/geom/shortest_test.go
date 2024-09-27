package geom

import (
	"slices"
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

func assertPath(t *testing.T, want, got []P) {
	require.Equal(t, len(want), len(got))
	for i := 0; i < len(got); i++ {
		if i > 0 {
			assert.True(t, got[i-1].Y < got[i].Y)
		}
		assert.Equal(t, want[i], got[i])
	}
}
