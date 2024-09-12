package dag_test

import (
	"slices"
	"testing"

	"deedles.dev/vixen/internal/dag"
)

func TestDAG(t *testing.T) {
	var g dag.DAG
	err := g.Add("one", "two", "three")
	if err != nil {
		t.Fatal(err)
	}
	err = g.Add("two", "three")
	if err != nil {
		t.Fatal(err)
	}
	err = g.Add("three")
	if err != nil {
		t.Fatal(err)
	}
	if s := slices.Collect(g.All()); !slices.Equal(s, []string{"one", "two", "three"}) {
		t.Fatal(s)
	}
	err = g.Add("three", "one")
	if err != dag.ErrCyclic {
		t.Fatal(err)
	}
}
