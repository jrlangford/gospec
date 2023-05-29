package gospec

type Tracer interface {
	Append(*OpTrace)
}

type TSpecification[T any] interface {
	IsSatisfiedBy(t T) bool
	And(s TSpecification[T]) TSpecification[T]
	Or(s TSpecification[T]) TSpecification[T]
	Not() TSpecification[T]
	GetName() string
}

type BaseTSpecification[T any] struct {
	TSpecification[T]
	tracer Tracer
}

func (c *BaseTSpecification[T]) Init(s TSpecification[T], t Tracer) {
	c.TSpecification = s
	c.tracer = t
}

func (c *BaseTSpecification[T]) And(s TSpecification[T]) TSpecification[T] {
	return NewAndTSpecification[T](c.TSpecification, s, c.tracer)
}

func (c *BaseTSpecification[T]) Or(s TSpecification[T]) TSpecification[T] {
	return NewOrTSpecification[T](c.TSpecification, s, c.tracer)
}

func (c *BaseTSpecification[T]) Not() TSpecification[T] {
	return NewNotTSpecification[T](c.TSpecification, c.tracer)
}

type AndTSpecification[T any] struct {
	BaseTSpecification[T]
	left  TSpecification[T]
	right TSpecification[T]
}

func NewAndTSpecification[T any](left, right TSpecification[T], t Tracer) *AndTSpecification[T] {
	s := &AndTSpecification[T]{
		left:  left,
		right: right,
	}
	s.Init(s, t)
	return s
}

func (s *AndTSpecification[T]) IsSatisfiedBy(t T) bool {

	l := s.left.IsSatisfiedBy(t)

	trace := NewOpTraceWithLeft(
		And,
		NewExpressionTrace(s.left.GetName(), l),
	)

	if !l {
		s.tracer.Append(trace)
		return l
	}

	r := s.right.IsSatisfiedBy(t)
	leftAndRight := l && r

	trace.SetRight(NewExpressionTrace(s.right.GetName(), r))
	s.tracer.Append(trace)

	return leftAndRight
}

func (s *AndTSpecification[T]) GetName() string {
	return "And Expression"
}

type OrTSpecification[T any] struct {
	BaseTSpecification[T]
	left  TSpecification[T]
	right TSpecification[T]
}

func NewOrTSpecification[T any](left, right TSpecification[T], t Tracer) *OrTSpecification[T] {
	s := &OrTSpecification[T]{
		left:  left,
		right: right,
	}
	s.Init(s, t)
	return s
}

func (s *OrTSpecification[T]) GetName() string {
	return "Or Expression"
}

func (s *OrTSpecification[T]) IsSatisfiedBy(t T) bool {

	l := s.left.IsSatisfiedBy(t)

	trace := NewOpTraceWithLeft(
		Or,
		NewExpressionTrace(s.left.GetName(), l),
	)

	if l {
		s.tracer.Append(trace)
		return l
	}

	r := s.right.IsSatisfiedBy(t)
	leftOrRight := l || r

	trace.SetRight(NewExpressionTrace(s.right.GetName(), r))
	s.tracer.Append(trace)

	return leftOrRight
}

type NotTSpecification[T any] struct {
	BaseTSpecification[T]
	single TSpecification[T]
}

func NewNotTSpecification[T any](single TSpecification[T], t Tracer) *NotTSpecification[T] {
	s := &NotTSpecification[T]{
		single: single,
	}
	s.Init(s, t)
	return s
}

func (s *NotTSpecification[T]) IsSatisfiedBy(t T) bool {

	l := s.single.IsSatisfiedBy(t)

	trace := NewOpTraceWithLeft(
		Not,
		NewExpressionTrace(s.single.GetName(), l),
	)

	not := !l

	s.tracer.Append(trace)

	return not
}

func (s *NotTSpecification[T]) GetName() string {
	return "Not Expression"
}
