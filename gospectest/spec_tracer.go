package gospectest

type SpecTracer[T any] struct {
	traces []T
}

func NewSpecTracer[T any]() *SpecTracer[T] {
	return &SpecTracer[T]{
		traces: make([]T, 0),
	}
}

func (t *SpecTracer[T]) Append(item T) {
	t.traces = append(t.traces, item)
}

func (t *SpecTracer[T]) Get(i int) T {
	return t.traces[i]
}

func (t *SpecTracer[T]) Len() int {
	return len(t.traces)
}

func (t *SpecTracer[T]) Clear() {
	t.traces = t.traces[:0]
}
