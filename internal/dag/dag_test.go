package dag_test

import (
	"slices"
	"testing"

	"deedles.dev/vixen/internal/dag"
)

func TestDAG(t *testing.T) {
	var g dag.DAG[string]
	g.Add("one", "two")
	g.Add("one", "three")
	g.Add("two", "three")

	nodes, err := g.Topological()
	if err != nil {
		t.Fatal(err)
	}
	if !slices.Equal(nodes, []string{"one", "two", "three"}) {
		t.Fatal(nodes)
	}

	g.Add("three", "one")
	_, err = g.Topological()
	if err != dag.ErrCyclic {
		t.Fatal(err)
	}
}
