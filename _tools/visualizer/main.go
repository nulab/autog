// Command visualizer serves a small web page to edit a graph and preview the layout computed by autog.
//
// Usage:
//
//	go run ./_tools/visualizer [-addr localhost:8080] [-timeout 5s]
//
// Each layout runs in a child process (this same binary, started with -worker),
// so a panic or an infinite loop inside autog can't bring the server down.
package main

import (
	"bytes"
	"context"
	_ "embed"
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"strings"
	"time"

	"github.com/nulab/autog"
	"github.com/nulab/autog/graph"
)

//go:embed index.html
var indexHTML []byte

// layoutRequest is the payload sent by the page.
type layoutRequest struct {
	Source   string   `json:"source"`
	Settings settings `json:"settings"`
}

type settings struct {
	CycleBreaking string `json:"cycleBreaking"`
	Layering      string `json:"layering"`
	Ordering      string `json:"ordering"`
	Positioning   string `json:"positioning"`
	EdgeRouting   string `json:"edgeRouting"`

	NodeWidth              float64 `json:"nodeWidth"`
	NodeHeight             float64 `json:"nodeHeight"`
	LayerSpacing           float64 `json:"layerSpacing"`
	NodeSpacing            float64 `json:"nodeSpacing"`
	VirtualNodeSize        float64 `json:"virtualNodeSize"`
	NetworkSimplexThor     uint    `json:"networkSimplexThoroughness"`
	BrandesKoepfLayout     int     `json:"brandesKoepfLayout"`
	NonDeterministicGreedy bool    `json:"nonDeterministicGreedy"`
	ShowVirtualNodes       bool    `json:"showVirtualNodes"`
}

type layoutResponse struct {
	Nodes []node `json:"nodes"`
	Edges []edge `json:"edges"`
	Error string `json:"error,omitempty"`
	Took  string `json:"took,omitempty"`
}

type node struct {
	ID      string  `json:"id"`
	X       float64 `json:"x"`
	Y       float64 `json:"y"`
	W       float64 `json:"w"`
	H       float64 `json:"h"`
	Virtual bool    `json:"virtual"`
}

type edge struct {
	From           string       `json:"from"`
	To             string       `json:"to"`
	Points         [][2]float64 `json:"points"`
	ArrowHeadStart bool         `json:"arrowHeadStart"`
}

var (
	cycleBreakings = map[string]autog.Option{
		"greedy": autog.WithCycleBreaking(autog.CycleBreakingGreedy),
		"dfs":    autog.WithCycleBreaking(autog.CycleBreakingDepthFirst),
	}
	layerings = map[string]autog.Option{
		"networksimplex": autog.WithLayering(autog.LayeringNetworkSimplex),
		"longestpath":    autog.WithLayering(autog.LayeringLongestPath),
	}
	orderings = map[string]autog.Option{
		"wmedian": autog.WithOrdering(autog.OrderingWMedian),
		"noop":    autog.WithOrdering(autog.OrderingNoop),
	}
	positionings = map[string]autog.Option{
		"sinkcoloring":   autog.WithPositioning(autog.PositioningSinkColoring),
		"brandeskoepf":   autog.WithPositioning(autog.PositioningBrandesKoepf),
		"networksimplex": autog.WithPositioning(autog.PositioningNetworkSimplex),
		"valign":         autog.WithPositioning(autog.PositioningVAlign),
		"packright":      autog.WithPositioning(autog.PositioningPackRight),
		"noop":           autog.WithPositioning(autog.PositioningNoop),
	}
	edgeRoutings = map[string]autog.Option{
		"polyline": autog.WithEdgeRouting(autog.EdgeRoutingPolyline),
		"straight": autog.WithEdgeRouting(autog.EdgeRoutingStraight),
		"ortho":    autog.WithEdgeRouting(autog.EdgeRoutingOrtho),
		"splines":  autog.WithEdgeRouting(autog.EdgeRoutingSplines),
		"noop":     autog.WithEdgeRouting(autog.EdgeRoutingNoop),
	}
)

func main() {
	addr := flag.String("addr", "localhost:8080", "address to listen on")
	timeout := flag.Duration("timeout", 5*time.Second, "maximum time allowed for one layout")
	worker := flag.Bool("worker", false, "internal: compute one layout from stdin and write it to stdout")
	flag.Parse()

	if *worker {
		runWorker()
		return
	}

	self, err := os.Executable()
	if err != nil {
		log.Fatal(err)
	}

	http.HandleFunc("GET /{$}", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		w.Write(indexHTML)
	})
	http.HandleFunc("POST /layout", func(w http.ResponseWriter, r *http.Request) {
		body, err := io.ReadAll(io.LimitReader(r.Body, 1<<20))
		if err != nil {
			writeJSON(w, layoutResponse{Error: err.Error()})
			return
		}
		writeJSON(w, runInChild(r.Context(), self, *timeout, body))
	})

	log.Printf("autog visualizer listening on http://%s", *addr)
	log.Fatal(http.ListenAndServe(*addr, nil))
}

// runInChild computes the layout in a separate process and kills it after the timeout.
func runInChild(ctx context.Context, self string, timeout time.Duration, req []byte) layoutResponse {
	ctx, cancel := context.WithTimeout(ctx, timeout)
	defer cancel()

	start := time.Now()
	cmd := exec.CommandContext(ctx, self, "-worker")
	cmd.Stdin = bytes.NewReader(req)
	var stdout, stderr bytes.Buffer
	cmd.Stdout = &stdout
	cmd.Stderr = &stderr

	err := cmd.Run()
	took := time.Since(start).Round(time.Millisecond).String()
	if errors.Is(ctx.Err(), context.DeadlineExceeded) {
		return layoutResponse{Error: fmt.Sprintf("layout did not finish within %s and was killed", timeout), Took: took}
	}

	var resp layoutResponse
	if jerr := json.Unmarshal(stdout.Bytes(), &resp); jerr != nil {
		msg := strings.TrimSpace(stderr.String())
		if msg == "" && err != nil {
			msg = err.Error()
		}
		return layoutResponse{Error: "layout process failed: " + msg, Took: took}
	}
	resp.Took = took
	return resp
}

func runWorker() {
	var req layoutRequest
	if err := json.NewDecoder(os.Stdin).Decode(&req); err != nil {
		json.NewEncoder(os.Stdout).Encode(layoutResponse{Error: "bad request: " + err.Error()})
		return
	}
	resp := func() (resp layoutResponse) {
		defer func() {
			if r := recover(); r != nil {
				resp = layoutResponse{Error: fmt.Sprintf("autog panicked: %v", r)}
			}
		}()
		return computeLayout(req)
	}()
	json.NewEncoder(os.Stdout).Encode(resp)
}

func computeLayout(req layoutRequest) layoutResponse {
	edges, nodeIDs, err := parseSource(req.Source)
	if err != nil {
		return layoutResponse{Error: err.Error()}
	}
	if len(edges) == 0 {
		return layoutResponse{}
	}

	s := req.Settings
	opts := []autog.Option{
		autog.WithNodeFixedSize(s.NodeWidth, s.NodeHeight),
		autog.WithLayerSpacing(s.LayerSpacing),
		autog.WithNodeSpacing(s.NodeSpacing),
		autog.WithVirtualNodeFixedSize(s.VirtualNodeSize),
		autog.WithNetworkSimplexThoroughness(s.NetworkSimplexThor),
		autog.WithBrandesKoepfLayout(s.BrandesKoepfLayout),
		autog.WithOutputVirtualNodes(s.ShowVirtualNodes),
	}
	if s.NonDeterministicGreedy {
		opts = append(opts, autog.WithNonDeterministicGreedyCycleBreaker())
	}
	for _, choice := range []struct {
		name  string
		value string
		algs  map[string]autog.Option
	}{
		{"cycle breaking", s.CycleBreaking, cycleBreakings},
		{"layering", s.Layering, layerings},
		{"ordering", s.Ordering, orderings},
		{"positioning", s.Positioning, positionings},
		{"edge routing", s.EdgeRouting, edgeRoutings},
	} {
		opt, ok := choice.algs[choice.value]
		if !ok {
			return layoutResponse{Error: fmt.Sprintf("unknown %s algorithm %q", choice.name, choice.value)}
		}
		opts = append(opts, opt)
	}

	out := autog.Layout(graph.EdgeSlice(edges), opts...)

	resp := layoutResponse{Nodes: []node{}, Edges: []edge{}}
	for _, n := range out.Nodes {
		resp.Nodes = append(resp.Nodes, node{
			ID: n.ID, X: n.X, Y: n.Y, W: n.W, H: n.H,
			// autog doesn't flag virtual nodes in its output: any node that isn't in the input is virtual
			Virtual: !nodeIDs[n.ID],
		})
	}
	for _, e := range out.Edges {
		resp.Edges = append(resp.Edges, edge{From: e.FromID, To: e.ToID, Points: e.Points, ArrowHeadStart: e.ArrowHeadStart})
	}
	return resp
}

// parseSource reads one edge per line, written as "a -> b" or "a b".
// Blank lines and lines starting with # or // are ignored.
func parseSource(src string) ([][]string, map[string]bool, error) {
	var edges [][]string
	ids := map[string]bool{}
	for i, line := range strings.Split(src, "\n") {
		line = strings.TrimSpace(line)
		if line == "" || strings.HasPrefix(line, "#") || strings.HasPrefix(line, "//") {
			continue
		}
		fields := strings.Fields(strings.ReplaceAll(line, "->", " "))
		if len(fields) != 2 {
			return nil, nil, fmt.Errorf("line %d: expected \"from -> to\", got %q", i+1, line)
		}
		edges = append(edges, fields)
		ids[fields[0]] = true
		ids[fields[1]] = true
	}
	return edges, ids, nil
}

func writeJSON(w http.ResponseWriter, v any) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(v)
}
