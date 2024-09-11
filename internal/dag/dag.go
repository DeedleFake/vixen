package dag

import (
	"iter"
	"slices"
)

// DAG is a directed, acyclic graph. It uses the provided comparison
// function to order things in the graph. The function should not be
// changed once at least one item has been inserted into the graph.
// The graph is not safe for concurrent use.
type DAG[V any] struct {
	Compare func(V, V) int

	roots []*node[V]
}

type node[V any] struct {
	val  V
	next []*node[V]
}

func traverse[V any](cur *node[V], yield func(*node[V]) bool) bool {
	if !yield(cur) {
		return false
	}

	for _, n := range cur.next {
		if !traverse(n, yield) {
			return false
		}
	}

	return true
}

func (dag *DAG[V]) nodes() iter.Seq[*node[V]] {
	return func(yield func(*node[V]) bool) {
		for _, r := range dag.roots {
			if !traverse(r, yield) {
				return
			}
		}
	}
}

func (dag *DAG[V]) topographical() []*node[V] {
	return slices.SortedStableFunc(dag.nodes(), func(n1, n2 *node[V]) int {
		return dag.Compare(n1.val, n2.val)
	})
}

// All returns an iterator over the nodes of dag in a topographical
// ordering. The returned iterator operates on a snapshot of the DAG's
// current state, meaning that it is reusable and that it will not
// reflect any changes that happen to the DAG after it is created.
//
// Calling this function concurrently is not safe, but it is safe to
// use the returned iterator concurrently with other DAG operations.
func (dag *DAG[V]) All() iter.Seq[V] {
	nodes := dag.topographical()
	return func(yield func(V) bool) {
		for _, n := range nodes {
			if !yield(n.val) {
				return
			}
		}
	}
}

func (dag *DAG[V]) Add(val V) []*node[V] {
	panic("Not implemented.")
}

// Complete returns true if every element in the DAG relates to every
// other element in the DAG at least indirectly.
func (dag *DAG[V]) Complete() bool {
	return len(dag.roots) <= 1
}
