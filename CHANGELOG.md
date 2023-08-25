
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
