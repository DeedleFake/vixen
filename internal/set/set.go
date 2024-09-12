package set

import (
	"iter"
	"maps"
)

type Set[T comparable] map[T]struct{}

func (set Set[T]) Add(v T) {
	set[v] = struct{}{}
}

func (set Set[T]) Has(v T) bool {
	_, ok := set[v]
	return ok
}

func (set Set[T]) Delete(v T) {
	delete(set, v)
}

func (set Set[T]) All() iter.Seq[T] {
	return maps.Keys(set)
}
