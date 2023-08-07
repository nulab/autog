package graph

import "fmt"

type Layer struct {
	Nodes []*Node
	Index int
	Size
}

func (layer *Layer) String() string {
	return fmt.Sprint(layer.Nodes)
}

func (layer *Layer) Len() int {
	return len(layer.Nodes)
}

// Tail returns the last node in this layer, i.e. the node with greatest index
func (layer *Layer) Tail() *Node {
	return layer.Nodes[len(layer.Nodes)-1]
}
