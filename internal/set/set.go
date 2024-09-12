package set

import "iter"

type Set[T comparable] map[T]struct{}

func Collect[T comparable](seq iter.Seq[T]) Set[T] {
	s := make(Set[T])
	for v := range seq {
		s.Add(v)
	}
	return s
}

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
