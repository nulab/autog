# autog

[![Go Report Card](https://goreportcard.com/badge/github.com/nulab/autog)](https://goreportcard.com/report/github.com/nulab/autog) ![Unit Tests](https://github.com/nulab/autog/actions/workflows/go.yml/badge.svg) [![Go Reference](https://pkg.go.dev/badge/github.com/nulab/autog.svg)](https://pkg.go.dev/github.com/nulab/autog)


autog (auto-graph) is an open-source Go library designed to automatically create hierarchical diagrams. 
It generates visually readable and aesthetically pleasing graphical layouts of directed graphs.

autog takes a directed graph as input and calculates the X and Y coordinates of each node and edge, producing a layout that is visually intuitive. 
However, it's essential to note that autog focuses on the layout computation and does not handle rendering. 
The rendering step is delegated to the client application, which can choose the preferred output format.

autog has been inspired by, and in essence is, a Go port of [Graphviz Dot](https://graphviz.org/docs/layouts/dot/) and [ELK (Eclipse Layout Kernel)](https://projects.eclipse.org/projects/modeling.elk). 
Unlike many similar projects that only provide a frontend in their language of choice to Graphviz, autog implements the algorithms in pure Go.

# Table of Contents
- [Usage](#usage)
- [Overview](#overview)
- [Imlementation Details](#implementation-details)
  - [Cycle breaking](#cycle-breaking)
  - [Layering](#layering)
  - [Ordering](#ordering)
  - [Positioning](#positioning)
  - [Edge routing](#edge-routing)
  - [References](#references)
- [Status](#status)

## Usage

autog is a Go library that can be added as a dependency to your project using standard Go module commands. It requires Go 1.22:

    $ go get -u github.com/nulab/autog@latest

To help you get started, here's a basic usage example:

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
	// obtain a graph.Source (here by converting the input to EdgeSlice) 
	src := graph.EdgeSlice(adj)
    
    // run the default autolayout pipeline
    layout := autog.Layout(src)
	
    // print the computed node coordinates
    for _, n := range layout.Nodes {
        fmt.Println(n.X, n.Y)
    }
    // print the computed edge control points
    for _, e := range layout.Edges {
        fmt.Println(e.Points)
    }
}
```

## Overview

Hierarchical graph layout algorithms typically involve five primary phases, executed sequentially:

1. Cycle Breaking: This phase ensures that the input is a Directed Acyclic Graph (DAG). It identifies a set of edges that, if reversed, would make the graph acyclic and reverses them until the computation is complete. 
2. Layering: All nodes are arranged in horizontal layers to ensure that edges point downwards.
3. Ordering: Nodes in each layer are reordered to minimize edge crossings.
4. Positioning: Actual X and Y coordinates for each node are calculated to ensure proper spacing without overlaps while preserving the relative node order established in the previous phase.
5. Edge Routing: Actual X and Y coordinates for each edge's start, end, and bend points are determined. Any edges reversed in phase 1 are restored to their original direction.

The autog pipeline runs default implementations for each of these phases. 
However, it's possible to override individual defaults using functional options. For example:

```go
    // import positioning "github.com/nulab/autog/phase4"
    autog.Layout(
        g,
        // override default phase4 implementation
        autog.WithPositioning(autog.PositioningVAlign),
    )
```
You can also customize other algorithm parameters, such as max iterations and multiplying factors. 
Refer to the documentation on the `"github.com/nulab/autog/internal/graph".Params` type for details on all configurable parameters or 
inspect `autog.With*` functional options in `autolayout_options_params.go` to know which ones you can tweak.

Different phase implementations can yield different final layouts, and the choice of phases may vary depending on your specific graph. 
Consult the documentation for each phase constant for guidance on selecting the right one for your use case.

## Implementation details

Here is a brief write-up about some of the more juicy implementation details. Reference papers and other resources are indicated in square brackets.
See the [References](#references) section for citations and links.

### Cycle breaking

- `Greedy`: This is a solution to the feedback arc set ported from [ELK] source code, which in turn builds on the work of Eades et al. and Di Battista et al.
It may get stuck in local optima. I found it useful to break two-node cycles such as edge pairs A->B and B->A  before running this algorithm. 
**NOTE:** The algorithm is non-deterministic; it arranges non-sink non-source nodes in the arc diagram randomly before reversing edges. Indeed this produces variable results. 
This isn't great for testing/debugging, so as of v0 the random selection is disabled and replaced with `nodes[len(nodes)/2]`.

### Layering

- `NetworkSimplex`: This is the algorithm described in [GKNV] except for the incremental computation of cut values, which is still a naive O(N^2) loop.
The algorithm sets an initial layering by making all edges point downwards. Basically the initial layer of a node `v` is 1 plus the maximum layer of adjacent incoming nodes. 
It then constructs an undirected spanning tree of the underlying graph, assigns cut values and replaces tree edges with negative cut values with non-tree edges with minimum slack until all edges have non-negative cut values.
The current Dot implementation in [GVDOT] has a few optimizations not mentioned in the paper, e.g. the initial feasible tree is constructed by merging feasible trees from a min-heap. 

### Ordering

- `WMedian`: This is the weighted median ordering routine described in [GKNV], also called `mincross` in [GVDOT]. I find it easier to reason about, whereas [ELK] implements something closer to the original Sugiyama method.
autog's current implementation doesn't properly account for flat edges — i.e. edges between nodes in the same layer. That part is summarily described in [GKNV]. 
In practice, even though flat edges are theoretically possible, they rarely result from layering the graph with the network simplex method.
Edge crossings are counted with the method described in [BM], which has a O(E log V) run time (V is the layer with fewer nodes). 
[BM] omits the implementation of the preliminary radix sort step. autog attempts to implement a simple O(E+V) routine. 

### Positioning

- `VAlign` (Vertical Align): This is a simple alignment and equal spacing of all nodes around the center of the largest layer.
- `NetworkSimplex`: This is the network simplex positioning described in [GKNV], including the horizontal balancing step which is only briefly mentioned in the paper but can be found in [GVDOT]. 
Note that without the min-heap feasible tree and incremental cut values optimizations, this exhibits a poor run time. [GKNV] also mentions an index rotation strategy to choose the edge to replace, which autog doesn't implement.  
- `BrandesKoepf`: This is the O(N) heuristic described in [BK]. The implementation in [ELK] is slightly harder to follow due to Java-style OOP and partially undocumented modifications introduced with [RSCH], therefore autog's implementation follows the original [BK].
It's worth noting that [BK] indeed does not account for node size. Final layouts may present overlaps due to zero-sizes of vitual nodes or unequal sizes of real nodes. It should be extended to properly incorporate [RSCH]'s ideas. 
- `SinkColoring`: This is an original (to my knowledge) algorithm based on [BK] which accounts for node sizes. It also prioritizes long straight segments and produces an overall more orthogonal layout. Useful as a size-aware solution until [BK] is properly finalized.

### Edge routing

- `Ortho`: This draws edges as orthogonal polylines. The result is similar to edge routing in ELK. ELK actually draws edges that have a common source or target node using the same line segment up until the 
point where they start diverging. autog simplifies this a bit and computes each edge coordinates independently. Both in ELK and autog, the presence of overlapping or common edge segments makes it harder to 
see where an edge starts and finishes. Therefore, this routing strategy works best when the graph has few edges with common source/target nodes.
- `Splines` (Work in Progress): This implements cubic Bezier curves as described in [DGKN]. A spline routing algorithm was originally described in [GKNV] but that has been superseded by [DGKN]. 
The algorithm in [GVDOT] still computes bounding boxes with some resemblance to [GKNV], however I must admit the C sources are quite hard to read due to the amount of static variables and functions with side-effects. 
The general idea here does follow [DGKN]: it triangulates the polygon obtained by merging the bounding box rectangles together, finds a shortest path through this polygon using the edge's starting and end points, finally it fits a cubic Bezier spline to 
the set of points in the shortest path. autog does things a little differently due to the scarcity of details in [DGKN] and in the available literature. 
The biggest obstacle is that most resources about polygon triangulation assume the input is a **strictly** monotone polygon, whereas the shape obtained from merged rectangles is monotone but not strictly monotone. 
More formally, for each point `P[i]` in a y-monotone polygon, O(N) triangulation assumes `P[i].Y > P[i-1].Y`, whereas with adjacent rectangles we have `P[i].Y >= P[i-1].Y`. The literature seems to confirm that non-strict monotonicity still admits linear time triangulation,
however the details of how the algorithm must change to accomodate for equal Y coordinates are always omitted. My understanding is that strict monotonicity is used to guarantee linear time sorting of the points in the polygon's left-right or upper-lower chains, with a typical
"merge sorted arrays" strategy. However, in a non-strictly monotone polygon, the two point chains are *not* sorted. An additional sorting step seems to be required, therefore we apparently fall back to O(N log N), which is the same running time of triangulation of arbitrary simple polygons, 
which includes cutting the polygon in strictly monotone sub-polygons. Therefore, I decided to cut the knot in autog and triangulate the merged rectangles in linear time using a special-cased holistic approach, whose correctness I'm currently unable to prove. But it appears to work well in practice. 
Bug reports will probably help refine this routine unless a different strategy is employed.
Once the polygon is triangulated, autog finds the shortest path using a "funnel" algorithm based on a dequeue. [GVDOT] seems to follow Lee and Preparata, while autog follows [HS]. The implementation is basically the same.
With the set of points defining the shortest path, both [GVDOT] and autog fit a cubic Bezier spline to it using the method in [Sch]. As [DGKN] mentions, this cubic Bezier is actually piece-wise and not a single cubic spline. As a matter of fact [Sch] is also 
recursive: [Sch] does this to improve the fit to an arbitrary polygonal path, [GVDOT] instead attempts to fit the spline within the edges of the constraint polygon — [DGKN] calls them "barriers". After a first attempt at fitting the spline, the algorithm computes the maximum square distance
between the input path points and the corresponding points on the parametric curve, cuts che curve at that point and then calls the same routine recursively on the two paths thus obtained.
The result is indeed a set of cubic Bezier control points.

### References

- [GKNV] Gansner, Koutsofios, North, Vo, "A Technique for Drawing Directed Graphs", AT&T Bell Laboratories ([link](https://www.graphviz.org/documentation/TSE93.pdf))
- [GVDOT] Graphviz source code at https://gitlab.com/graphviz/graphviz
- [DGKN] Dobkin, Gansner, Koutsofios, North, "Implementing a General-Purpose Edge Router", Princeton University and AT&T Bell Laboratories ([link](https://dpd.cs.princeton.edu/Papers/DGKN97.pdf))
- [ELK] ELK source code at https://github.com/eclipse/elk
- [BK] Brandes and Köpf, "Fast and Simple Horizontal Coordinate Assignment", Department of Computer & Information Science, University of Konstanz ([link](https://link.springer.com/content/pdf/10.1007/3-540-45848-4_3.pdf?pdf=inline%20link))
- [RSCH] Rüegg, Schulze, Carstens, von Hanxleden, "Size- and Port-Aware Horizontal Node Coordinate Assignment", Dept. of Computer Science, Christian-Albrechts-Universität zu Kiel ([link](https://rtsys.informatik.uni-kiel.de/~biblio/downloads/papers/gd15.pdf))
- [BM] Barth and Mutzel, "Simple and Efficient Bilayer Cross Counting", Institut für Computergraphik und Algorithmen Technische Universität Wien ([link](https://pdfs.semanticscholar.org/272d/73edce86bcfac3c82945042cf6733ad281a0.pdf))
- [Sch] Schneider, "An algorithm for automatically fitting digitized curves" in Andrew S. Glassner, editor, Graphics Gems, pages 612–626. Academic Press, Boston, Mass., 1990. (can be found in internet)
- [HS] Hershberger and Snoeyink, "Computing minimum length paths of a given homotopy class", Computational Geometry Volume 4, Issue 2, June 1994, Pages 63-97 (can be found in internet)

## Status

This project is actively under development, but it is currently in version 0. 
Please be aware that the public API and exported methods may undergo changes.

- Self-loops don't break the program any more ([issue #23](https://github.com/nulab/autog/issues/23)) but are not supported. The final layout includes self-loop edges but those edges are not routed (`e.Points` is `nil`)

## Commit guidelines

Commits should be prefixed with a short name in square brackets (a tag) that summarizes 
which main area of the code was changed. The prefixes may change as the repository structure changes. 

The prefixes are:

- `[pN-name]`: phase `N`, where `N` is a number from 1 to 5 and an optional `name` mnemonic indicating a phase N's algorithm
- `[preproc]`: preprocessing code that runs before phase 1
- `[postproc]`: postprocessing code that runs after phase 5
- `[graph]`: the `graph` package(s), changes to the structures and types used throughout the program
- `[monitor]`: the `monitor` package
- `[geom]`: the `geom` package
- `[docs]`: documentation, within the code (e.g. comments for godocs) or README
- 

## Bug reporting

If you encounter a bug, please open a new issue and include at least the input graph that triggered the bug to help us reproduce and address it.

## Authors

* **[Gabriele V.](https://github.com/vibridi/)** - *Main contributor*
* Currently, there are no other contributors

## License

This project is licensed under the MIT License. For detailed licensing information, refer to the [LICENSE](LICENSE) file included in the repository.
