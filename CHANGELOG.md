
v0.5.0 / 2024-06-20
==================

  * Upgrade to Go 1.22, remove exp/constraints package (#17)
  * Update README based on latest API changes
  * [p5,geom] Implement helper package for cubic bezier spline routing (#15)

v0.4.0 / 2023-09-19
==================

  * [p3] Remove traces of old phase 3 wmedian naming
  * [p5-ortho] Compute edge control points using routable edges (fixes #13)
  * [p5-poly] Compute edge control points using routable edges
  * [p5] Fix route merging, append non-terminal nodes correctly
  * [p3-wmedian] Rename graphviz dot ordering to wmedian to reflect actual theory
  * [p5-poly] Rename piecewise edge routing strategy to polyline
  * [options] Extract alg constants into public API with documentation and refactor alg functional options accordingly
  * Fix edge merging, ensure mapping between edges and route info
  * Output only graph items with minimum set of fields necessary to draw the layout
  * [graph] Add temp option to control global node size
  * [graph] Add flag to remember where to draw the arrowhead
  * Move public graph's node and edge definitions in their own files
  * [graph] Delete connected components routine, currently unused
  * [p5] Merge long edges as a preprocessing step before routing
  * Fix or skip existing unit tests so that the current regression set passes
  * Change exported type to Layout and expose only fields relevant to the user
  * Remove NS balancing from functional options
  * Refactor functional options to avoid using internal consts; this also prevents accidentally passing invalid numerical values
  * Define EdgeSlice as implementation of graph.Source
  * Move BreakLongEdges into phase3 package as top-level func
  * Make all phase and graph packages internal
  * [p1,p5] Restore reverted edges as a post-processing step of phase 5
  * Update README with implementation details
  * [p2-nsimplex] Rewrite head component boolean logic to follow Graphviz's paper
  * Change network simplex balancing strategy to typed enum
  * [p2-nsimplex] Fix and reenable horizontal balancing used in p4 (fixes #10)
  * Add basic Makefile with test target
  * Move test adjacency lists into non-test file for reuse in other packages with test-only build tag
  * Move layout tests into internal package

v0.3.0 / 2023-08-29
==================

  * [p5] Add alg unit test
  * [p4] Add alg unit test
  * [p3] Add alg unit test
  * [p2] Add alg unit test
  * [p1] Add alg unit test
  * [p1] Document exported methods and consts
  * [p2] Document exported methods
  * [p3] Document exported methods
  * [p4-noop] Don't assign Y coordinates to match documentation
  * [p4] Document exported methods
  * Move all algs Process methods into separate files
  * [monitor] Add monitor helpers for unit tests
  * Implement processor.P interface with all algs to set monitor prefixes
  * Move monitor option out of graph params
  * [monitor] Move monitor to internal package and simplify API
  * Move processor interface into internal package
  * [p3-graphviz] Fix sorting of -1 values in wmedian
  * [p3-graphviz] Return early if crossings is zero after init
  * [p2-nsimplex] Improve clarity of feasible span in vbalance
  * [p2-nsimplex] Fix vertical balancing incorrectly reusing updated var in loop stop condition

v0.2.2 / 2023-08-25
==================

  * [p4] Remove enum for median positioner

v0.2.0 / 2023-08-25
==================

  * [p2-nsimplex] Fix checking whether node is in edge's head component
  * [p4-nsimplex] Set w/h of nodes in the auxiliary graph
  * [p2,p4-nsimplex] Change NetworkSimplexBalance parameter to integer to distinguish between vertical and horizontal balancing
  * [p1] Fail fast if after cycle breaking the graph is not acyclic
  * [p2-nsimplex] Simplify logic in lim/low assignment
  * [p5] Implement orthogonal edge router
  * [p1] Do not run cycle breaker if removal of 2-node cycles already makes the graph acyclic
  * [p4-coloring] Return false when testing if an edge crosses an priority edge on different layers
  * [p3-graphviz] Fix initialization of same-layer transitive closures
  * [p2-nsimplex] Improve error message when non-incident tree edges aren't found
  * [p4-b&k] Add optional parameter to choose a particular B&K layout
  * [p3] Ignore nil entries from radix sort when counting edge crossings
  * [p1] Keep eagerly reversed edges in the edge list
  * Add non-regression and bugfix layout tests
  * [p2] Always fill layers at the end of phase 2

v0.1.3 / 2023-08-21
==================

  * Add draft unit tests for crashers
  * Restore SinkColoring as default positioning option
  * [p4-coloring] Consider edge priorities when painting blocks (fixes #4)
  * [p4-b&k] Refactor marking edge conflicts into a top-level function to improve reuse
  * Remove node margin as a parameter from all algorithms, use only node spacing (fixes #5)

v0.1.2 / 2023-08-16
==================

  * [p5] Restore hidden edges before routing (fixes #3)
