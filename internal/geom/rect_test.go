package geom

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestRect(t *testing.T) {
	r := Rect{TL: P{10, 20}, BR: P{50, 30}}
	assert.Equal(t, 40.0, r.Width())
	assert.Equal(t, 10.0, r.Height())
	assert.Equal(t, `<rect class="rect" x="10.00" y="20.00" width="40.00" height="10.00" style="fill: lightgrey; stroke: black;" />`, r.SVG())

}
