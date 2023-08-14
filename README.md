# autog

autog (auto-graph) is a Go library to automatically draw hierarchical diagrams.
Given a directed graph as input, autog computes the X and Y coordinates of each node and edge that result in a visually readable and pleasant layout.

As usual with graph layout algorithms, autog only outputs the final X and Y coordinates of graph elements; it does not render them.
The rendering step is left to the client application based on the preferred output format. 

A significant part of this project is a port of [ELK (Eclipse Layout Kernel)](https://projects.eclipse.org/projects/modeling.elk) to Go, 
with some deviations usually documented right in the sources. 

## Status 

This project is being actively developed, however it is currently on version `0`. Its public API and exported methods may change.

### Bug reporting

If you find a bug, please open a new issue. Make sure to provide at least the input graph that caused the bug so that we can reproduce it.

## Usage

autog is a library. Add it as a dependency of your project as with any Go module:

    $ go get -u github.com/nulab/autog@latest

### Requirements

Go 1.21

### Example

In its most simple form:

```go
package main

import (
	"fmt"
    "github.com/nulab/autog"
	"github.com/nulab/autog/graph"
)

func main() {
    // construct the graph object from an ordered adjacency list
    adj := [][]string{
        {"N1", "N2"}, // edge from node N1 to node N2
        {"N2", "N3"}, // edge from node N2 to node N3
    }
    g := graph.FromEdgeSlice(adj)
    
    // run the default autolayout pipeline
    autog.Layout(g)
	
    // print the computed node coordinates
    for _, n := range g.Nodes {
        fmt.Println(n.X, n.Y)
    }
    // print the computed edge control points
    for _, e := range g.Edges {
        fmt.Println(e.Points)
    }
}
```

Hierarchical graph layout algorithms usually comprise five main phases run in succession: 
1. cycle breaking: the input must be a DAG (Directed Acyclic Graph), so the first phase finds a set of edges that, if reversed, would
make the graph acyclic and reverses them until the end of the computation. 
2. layering: all nodes are arranged in horizontal layers so that edges point downwards
3. ordering: all nodes in each layer are reordered until the edge crossings between them are minimized
4. positioning: the actual X and Y coordinates of each node are computed so that nodes are properly spaced and don't overlap, 
while maintaining the relative node order found in the previous phase
5. edge routing: the actual X and Y coordinates of each edge's start, end and bend point are computed. 
Finally, edges that were reversed in phase 1 are restored to their original direction.

The pipeline runs a default implementation for each phase. It is possible to override individual defaults via functional options:

```go
    // import positioning "github.com/nulab/autog/phase4"
    autog.Layout(
        g,
        // override default phase4 implementation
        autog.WithPositioning(positioning.VerticalAlign),
    )
```
It is also possible to override other default algorithm parameters such as max iterations, multiplying factors, etc. 
Please see the relevant documentation on the `graph.Params` type for more information available configurations.

Different phase implementations may result in different final layouts, and different combinations may make more sense for
different graphs. Please see the documentation on each phase constant for tips about choosing the right one for your use case. 
If you are in doubt, stick to the defaults.

## Authors

* **[Gabriele V.](https://github.com/vibridi/)** - *Main contributor*

Currently there are no other contributors

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details
