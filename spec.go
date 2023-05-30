package gospec

type Specification[T any] interface {
	IsSatisfiedBy(t T) bool
	And(s Specification[T]) Specification[T]
	Or(s Specification[T]) Specification[T]
	Not() Specification[T]
}

type BaseSpecification[T any] struct {
	childSpec Specification[T]
}

func (c *BaseSpecification[T]) Init(s Specification[T]) {
	c.childSpec = s
}

func (c *BaseSpecification[T]) And(s Specification[T]) Specification[T] {
	return NewAndSpecification[T](c.childSpec, s)
}

func (c *BaseSpecification[T]) Or(s Specification[T]) Specification[T] {
	return NewOrSpecification[T](c.childSpec, s)
}

func (c *BaseSpecification[T]) Not() Specification[T] {
	return NewNotSpecification[T](c.childSpec)
}

type AndSpecification[T any] struct {
	BaseSpecification[T]
	left  Specification[T]
	right Specification[T]
}

func NewAndSpecification[T any](left, right Specification[T]) *AndSpecification[T] {
	s := &AndSpecification[T]{
		left:  left,
		right: right,
	}
	s.Init(s)
	return s
}

func (s *AndSpecification[T]) IsSatisfiedBy(t T) bool {
	return s.left.IsSatisfiedBy(t) && s.right.IsSatisfiedBy(t)
}

type OrSpecification[T any] struct {
	BaseSpecification[T]
	left  Specification[T]
	right Specification[T]
}

func NewOrSpecification[T any](left, right Specification[T]) *OrSpecification[T] {
	s := &OrSpecification[T]{
		left:  left,
		right: right,
	}
	s.Init(s)
	return s
}

func (s *OrSpecification[T]) IsSatisfiedBy(t T) bool {
	return s.left.IsSatisfiedBy(t) || s.right.IsSatisfiedBy(t)
}

type NotSpecification[T any] struct {
	BaseSpecification[T]
	single Specification[T]
}

func NewNotSpecification[T any](single Specification[T]) *NotSpecification[T] {
	s := &NotSpecification[T]{
		single: single,
	}
	s.Init(s)
	return s
}

func (s *NotSpecification[T]) IsSatisfiedBy(t T) bool {
	return !s.single.IsSatisfiedBy(t)
}
