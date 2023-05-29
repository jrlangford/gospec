package gospectest

type TestCase[IN any, OUT any] struct {
	Description    string
	Input          IN
	ExpectedOutput OUT
}

type SpecTracer struct {
	trace []string
}

func NewSpecTracer() *SpecTracer {
	return &SpecTracer{
		trace: make([]string, 0),
	}
}

func (t *SpecTracer) Append(s string) {
	t.trace = append(t.trace, s)
}

func (t *SpecTracer) Get(i int) string {
	return t.trace[i]
}

func (t *SpecTracer) Len() int {
	return len(t.trace)
}

func (t *SpecTracer) Clear() {
	t.trace = t.trace[:0]
}
