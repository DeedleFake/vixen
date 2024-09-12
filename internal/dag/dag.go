package dag

import (
	"errors"
	"iter"

	"deedles.dev/vixen/internal/set"
)

var (
	ErrCyclic = errors.New("graph is cyclic")
)

// DAG is a directed, acyclic graph. A zero-value DAG is ready to use.
//
// A DAG is not safe for concurrent use.
type DAG struct {
	edges map[string][]string
	roots set.Set[string]
}

func (dag *DAG) init() {
	if dag.edges == nil {
		dag.edges = make(map[string][]string)
	}
	if dag.roots == nil {
		dag.roots = make(set.Set[string])
	}
}

func (dag *DAG) recalcroots() {
	clear(dag.roots)

	nonroots := make(set.Set[string], len(dag.edges))
	for _, to := range dag.edges {
		for _, n := range to {
			nonroots.Add(n)
		}
	}

	for n := range dag.edges {
		if !nonroots.Has(n) {
			dag.roots.Add(n)
		}
	}
}

func (dag *DAG) cyclic() bool {
	if len(dag.roots) == 0 && len(dag.edges) > 0 {
		return true
	}

	visited := make(set.Set[string], len(dag.edges))

	for r := range dag.roots.All() {
		var count int
		check := func(n string) (bool, error) {
			if n == r {
				count++
				if count > 1 {
					return false, ErrCyclic
				}
			}

			defer visited.Add(n)
			return !visited.Has(n), nil
		}

		clear(visited)
		_, err := dag.traverse(r, func(string) bool { return true }, check)
		if err != nil {
			return true
		}
	}

	return false
}

func (dag *DAG) traverse(cur string, yield func(string) bool, visited func(string) (bool, error)) (bool, error) {
	ok, err := visited(cur)
	if err != nil {
		return false, err
	}
	if !ok {
		return true, nil
	}

	if !yield(cur) {
		return false, nil
	}

	for _, next := range dag.edges[cur] {
		ok, err := dag.traverse(next, yield, visited)
		if err != nil || !ok {
			return ok, err
		}
	}

	return true, nil
}

func (dag *DAG) Add(name string, edges ...string) error {
	dag.init()
	_, existed := dag.edges[name]

	dag.edges[name] = edges
	dag.recalcroots()
	if !dag.cyclic() {
		return nil
	}

	if !existed {
		delete(dag.edges, name)
		dag.recalcroots()
	}

	return ErrCyclic
}

// All returns an iterator that yields node names and values from dag
// in the order implied by the direction of the graph.
func (dag *DAG) All() iter.Seq[string] {
	return func(yield func(string) bool) {
		visited := make(set.Set[string], len(dag.edges))
		check := func(n string) (bool, error) {
			defer visited.Add(n)
			return !visited.Has(n), nil
		}

		for r := range dag.roots.All() {
			ok, _ := dag.traverse(r, yield, check)
			if !ok {
				return
			}
		}
	}
}
