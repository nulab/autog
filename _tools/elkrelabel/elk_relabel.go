package main

import (
	"encoding/json"
	"os"

	"github.com/nulab/autog/internal/testfiles"
)

func main() {
	originals := testfiles.ReadTestDir("internal/testfiles/elk_original")
	for _, elkg := range originals {
		for _, n := range elkg.Nodes {
			for _, l := range n.Labels {
				l.Text = n.ID
			}
		}
		f, err := os.OpenFile("internal/testfiles/elk_relabeled/"+elkg.Name, os.O_CREATE|os.O_WRONLY, 0644)
		if err != nil {
			panic(err)
		}
		enc := json.NewEncoder(f)
		enc.SetIndent("", "\t")
		err = enc.Encode(elkg)
		if err != nil {
			panic(err)
		}
	}
}
