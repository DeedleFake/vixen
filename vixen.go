package vixen

import (
	"reflect"
	"sync"

	"deedles.dev/vixen/internal/dag"
)

var (
	m         sync.Mutex
	deps      dag.DAG[reflect.Type]
	providers map[reflect.Type]provider
)

type provider struct {
	f reflect.Value
}

func Provide(p any) {
	panic("Not implemented.")
}

func Invoke(f any) {
	panic("Not implemented.")
}

// Require is a conenience function that simply invokes with a
// requirement for the given type.
func Require[T any]() {
	Invoke(func(T) {})
}
