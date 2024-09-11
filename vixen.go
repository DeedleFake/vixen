package vixen

func Provide(provider any) {
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
