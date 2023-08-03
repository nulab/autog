package main

import (
	"fmt"
	"math"
	"sync"

	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
	"github.com/nulab/autog/internal/elk"
	"github.com/nulab/autog/internal/testfiles"
	"github.com/nulab/autog/monitor"
	"github.com/nulab/autog/phase4"
)

const (
	exampleDiagram = "ci_router_ComplexRouter.json"
	itr            = 500
)

func main() {
	elkg := testfiles.ReadTestFile("internal/testfiles/elk_relabeled", exampleDiagram)
	checkSuitability(elkg)

	ad := elkg.AdjacencyList()
	ch := make(chan monitor.Log)
	wg := &sync.WaitGroup{}
	wg.Add(itr)

	go func() {
		wg.Wait()
		close(ch)
	}()

	go func() {
		for i := 0; i < itr; i++ {
			autog.Layout(
				graph.FromAdjacencyList(ad),
				autog.WithPositioning(phase4.NoPositioning),
				autog.WithMonitor(monitor.New(relayTo(ch, wg))),
			)
		}
	}()

	maxx := math.MinInt
	minx := math.MaxInt
	tot := 0
	i := 0
	for log := range ch {
		if log.Name == "phase3/graphvizdot/crossings" {
			x := log.Value.AsInt()
			maxx = max(maxx, x)
			minx = min(minx, x)
			tot += x
			i++
		}
	}
	avg := float64(tot) / float64(i)
	fmt.Printf("iterations: %d, crossings: avg: %.02f, min: %d, max: %d\n", i, avg, minx, maxx)
}

func checkSuitability(elkg *elk.Graph) {
	if len(graph.FromAdjacencyList(elkg.AdjacencyList()).ConnectedComponents()) != 1 {
		panic("monte carlo simulation must run on a graph with a single connected component")
	}
}

func relayTo(out chan monitor.Log, wg *sync.WaitGroup) chan monitor.Log {
	in := make(chan monitor.Log)
	go func() {
		defer wg.Done()
		for log := range in {
			out <- log
		}
	}()
	return in
}
