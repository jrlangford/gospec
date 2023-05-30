# gospec

An implementation of the specification pattern in Go.

The module includes a regular implementation and a traceable one, which can 'explain' its execution.

## Quick Start

Embed BaseSpecification, instantiated to the type that will be tested, in your own specification.
Ensure the method that checks for specification satisfaction has the following signature:  `IsSatisfiedBy(t T) bool`.

```go
type User struct {
  id      string
  name    string
}

type IsBob struct {
  BaseSpecification[User]
}

func NewIsBob() *IsBob {
  s := IsBob{}
  s.Init(&s)
  return &s
}

func (s *IsBob) IsSatisfiedBy(u User) bool {
  return u.name == "Bob"
}
```

Once you have defined and instantiated several specifications you can compose a more complex one by using Boolean operators.

```go
adultSpec := NewIsLegalAdult()
flaggedSpec := NewIsFlagged()
bobSpec := NewIsBob()

compositeSpec := adultSpec.Not().Or(flaggedSpec).Or(bobSpec),

satisfied := compositeSpec.IsSatisfiedBy(userBob)
```

Use BaseTSpecification instead if you want to be able to trace the execution of Boolean operations in a composite specification. This type requires the addition of a `GetName() string` method to your specification and also the inclusion of a 'Tracer', an interface that can process operation trace instances (OpTrace).

```go
type IsFlagged struct {
  gospec.BaseTSpecification[User]
  flaggedIDs  []string
  specName    string
}

func NewIsFlagged(t gospec.Tracer) *IsFlagged {
  s := IsFlagged{
    flaggedIDs: []string{"1"},
    specName:   "IsFlaggedSpec",
  }
  s.Init(&s, t)
  return &s
}

func (s *IsFlagged) IsSatisfiedBy(u User) bool {
  for _, id := range s.flaggedIDs {
    if id == u.id {
      return true
    }
  }
return false
}

func (s *IsFlagged) GetName() string {
  return s.specName
}
```

Operation trace instances can produce a verbose explanation of their contents.

Run the logged_specification example for a basic demonstration of its functionality.

```bash
go test -v examples/logged_specification_test.go
```

## Development Status: Beta

Core features have been tested, trace features have not been tested.
