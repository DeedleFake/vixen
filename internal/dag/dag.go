package dag

import (
	"errors"
	"iter"
	"slices"

	"deedles.dev/vixen/internal/set"
)

var (
	ErrCyclic = errors.New("graph is cyclic")
)

// DAG is a directed, acyclic graph. A zero-value DAG is ready to use.
//
// A DAG is not safe for concurrent use.
type DAG[Name comparable] struct {
	edges map[Name]set.Set[Name]
	roots set.Set[Name]
}

func (dag *DAG[Name]) recalcroots() {
	if dag.roots != nil {
		return
	}

	nonroots := make(set.Set[Name], len(dag.edges))
	for _, to := range dag.edges {
		for n := range to {
			nonroots.Add(n)
		}
	}

	dag.roots = make(set.Set[Name], len(dag.edges)-len(nonroots))
	for n := range dag.edges {
		if !nonroots.Has(n) {
			dag.roots.Add(n)
		}
	}
}

func (dag *DAG[Name]) Add(from, to Name) {
	if dag.edges == nil {
		dag.edges = make(map[Name]set.Set[Name])
	}

	edges := dag.edges[from]
	if edges == nil {
		edges = make(set.Set[Name])
		dag.edges[from] = edges
	}

	edges.Add(to)
	dag.roots = nil
}

func (dag *DAG[Name]) nodes() iter.Seq[Name] {
	return func(yield func(Name) bool) {
		visited := make(set.Set[Name], len(dag.edges))
		for from, to := range dag.edges {
			if !visited.Has(from) {
				if !yield(from) {
					return
				}
				visited.Add(from)
			}

			for n := range to {
				if !visited.Has(n) {
					if !yield(n) {
						return
					}
					visited.Add(n)
				}
			}
		}
	}
}

func (dag *DAG[Name]) rtopological() iter.Seq2[Name, bool] {
	dag.recalcroots()

	return func(yield func(Name, bool) bool) {
		t := make(set.Set[Name])
		u := set.Collect(dag.nodes())

		var visit func(Name) (bool, bool)
		visit = func(n Name) (cont, ok bool) {
			if !u.Has(n) {
				return true, true
			}
			if t.Has(n) {
				return false, false
			}

			t.Add(n)

			for m := range dag.edges[n] {
				cont, ok := visit(m)
				if !cont || !ok {
					return cont, ok
				}
			}

			u.Delete(n)
			if !yield(n, true) {
				return false, true
			}
			return true, true
		}

		for len(u) > 0 {
			for n := range u {
				cont, ok := visit(n)
				if !cont {
					if !ok {
						yield(n, false)
					}
					return
				}
			}
		}
	}
}

// Topological returns the nodes of dag sorted with a topological
// ordering. The ordering is not guaranteed to be the same every time.
//
// If the graph has an error, such as being cyclic, nil and an error
// are returned.
func (dag *DAG[Name]) Topological() ([]Name, error) {
	var nodes []Name
	for name, ok := range dag.rtopological() {
		if !ok {
			return nil, ErrCyclic
		}
		nodes = append(nodes, name)
	}
	slices.Reverse(nodes)
	return nodes, nil
}
