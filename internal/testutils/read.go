package testutils

import (
	"encoding/json"
	"os"
)

func ReadTestDir(dir string) []*Graph {
	fs, err := os.ReadDir(dir)
	if err != nil {
		panic(err)
	}
	graphs := make([]*Graph, len(fs))
	for i, f := range fs {
		graphs[i] = ReadTestFile(dir, f.Name())
	}
	return graphs
}

func ReadTestFile(dir, name string) *Graph {
	b, err := os.ReadFile(dir + "/" + name)
	if err != nil {
		panic(err)
	}
	g := Graph{
		Name: name,
	}
	err = json.Unmarshal(b, &g)
	if err != nil {
		panic(err)
	}

	return &g
}
