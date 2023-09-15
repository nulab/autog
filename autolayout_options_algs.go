package autog

import (
	"github.com/nulab/autog/internal/phase1"
	"github.com/nulab/autog/internal/phase2"
	"github.com/nulab/autog/internal/phase3"
	"github.com/nulab/autog/internal/phase4"
	"github.com/nulab/autog/internal/phase5"
)

func WithCycleBreakingGreedy() Option     { return func(o *options) { o.p1 = phase1.Greedy } }
func WithCycleBreakingDFS() Option        { return func(o *options) { o.p1 = phase1.DepthFirst } }
func WithLayeringLongestPath() Option     { return func(o *options) { o.p2 = phase2.LongestPath } }
func WithLayeringNetworkSimplex() Option  { return func(o *options) { o.p2 = phase2.NetworkSimplex } }
func WithOrderingNoop() Option            { return func(o *options) { o.p3 = phase3.NoOrdering } }
func WithOrderingGraphvizDot() Option     { return func(o *options) { o.p3 = phase3.GraphvizDot } }
func WithPositioningNoop() Option         { return func(o *options) { o.p4 = phase4.NoPositioning } }
func WithPositioningVAlign() Option       { return func(o *options) { o.p4 = phase4.VerticalAlign } }
func WithPositioningBrandesKoepf() Option { return func(o *options) { o.p4 = phase4.BrandesKoepf } }
func WithPositioningSinkColoring() Option { return func(o *options) { o.p4 = phase4.SinkColoring } }
func WithPositioningPackRight() Option    { return func(o *options) { o.p4 = phase4.PackRight } }
func WithEdgeRoutingNoop() Option         { return func(o *options) { o.p5 = phase5.NoRouting } }
func WithEdgeRoutingStraight() Option     { return func(o *options) { o.p5 = phase5.Straight } }
func WithEdgeRoutingPieceWise() Option    { return func(o *options) { o.p5 = phase5.PieceWise } }
func WithEdgeRoutingOrtho() Option        { return func(o *options) { o.p5 = phase5.Ortho } }

// todo: proper unit tests, proper documentation
