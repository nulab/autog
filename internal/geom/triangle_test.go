package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestTri(t *testing.T) {
	t.Run("contains", func(t *testing.T) {
		cases := []struct {
			t Tri
			p P
		}{
			{Tri{0, P{0, 0}, P{20, 0}, P{20, 50}}, P{10, 5}},
			{Tri{1, P{19.2, 44.7}, P{142.6, 16.5}, P{228, 212}}, P{28.6, 46.6}},
			{Tri{2, P{20, 20}, P{200, 20}, P{300, 40}}, P{197.7, 27}},
			{Tri{3, P{50, 10}, P{80, 10}, P{100, 50}}, P{65, 10}},            // collinear parallel x axis
			{Tri{4, P{255, 210}, P{0, 60}, P{255, 60}}, P{185.615, 169.185}}, // collinear neg slope
		}
		for _, c := range cases {
			assert.Truef(t, c.t.Contains(c.p), "triangle %d does not contain point", c.t.ID)
		}
	})
}
