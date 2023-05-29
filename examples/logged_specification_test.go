package gospec_test

import (
	"fmt"
	"log"
	"strings"

	"github.com/jrlangford/gospec"
)

type User struct {
	Name string
}

type UsernameContainsLetter struct {
	gospec.BaseTSpecification[User]
	letter   rune
	specName string
}

func NewUsernameContainsLetter(r rune, t gospec.Tracer) *UsernameContainsLetter {
	s := UsernameContainsLetter{
		letter:   r,
		specName: "UsernameContainsLetter[" + string(r) + "]",
	}
	s.Init(&s, t)
	return &s
}

func (s *UsernameContainsLetter) IsSatisfiedBy(p User) bool {
	return strings.ContainsRune(p.Name, s.letter)
}

func (s *UsernameContainsLetter) GetName() string {
	return s.specName
}

type TraceLogger struct{}

func NewTraceLogger() *TraceLogger {
	return &TraceLogger{}
}

func (l *TraceLogger) Append(operation *gospec.OpTrace) {
	ex, err := operation.Explain()
	if err != nil {
		log.Print(err)
	}
	log.Print(ex)
}

func Example() {
	user := User{
		Name: "Isabella",
	}

	traceLogger := NewTraceLogger()

	letterASpec := NewUsernameContainsLetter('a', traceLogger)
	letterCSpec := NewUsernameContainsLetter('c', traceLogger)
	letterESpec := NewUsernameContainsLetter('e', traceLogger)

	satisfied := letterASpec.And(letterCSpec.Not()).And(letterESpec).IsSatisfiedBy(user)
	fmt.Print(satisfied)
	// Output:
	// true
}
