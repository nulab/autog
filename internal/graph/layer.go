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

// Head returns the first node in this layer, i.e. the node with index 0
func (layer *Layer) Head() *Node {
	return layer.Nodes[0]
}

// Tail returns the last node in this layer, i.e. the node with the greatest index
func (layer *Layer) Tail() *Node {
	return layer.Nodes[len(layer.Nodes)-1]
}
