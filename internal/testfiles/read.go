package testfiles

import (
	"encoding/json"
	"os"

	"github.com/vibridi/autog/internal/elk"
)

func ReadTestDir(dir string) []*elk.Graph {
	fs, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	graphs := make([]*elk.Graph, len(fs))
	for i, f := range fs {
		graphs[i] = ReadTestFile(dir, f.Name())
	}
	return graphs
}

func ReadTestFile(dir, name string) *elk.Graph {
	b, err := os.ReadFile(dir + "/" + name)
	if err != nil {
		panic(err)
	}
	g := elk.Graph{
		Name: name,
	}
	err = json.Unmarshal(b, &g)
	if err != nil {
		panic(err)
	}

	return &g
}
