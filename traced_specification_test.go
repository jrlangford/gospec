package gospec_test

import (
	"testing"

	"github.com/jrlangford/gospec"
	g "github.com/jrlangford/gospec"
	th "github.com/jrlangford/gospec/gospectest"
)

type ID string

type User struct {
	id      ID
	name    string
	age     uint8
	heightM float64
}

// An IsLegalAdult is an age-based TSpecification example.
type IsLegalAdult struct {
	gospec.BaseTSpecification[User]
	ageOfMajority uint8
}

func NewIsLegalAdult(t gospec.Tracer) *IsLegalAdult {
	s := IsLegalAdult{
		ageOfMajority: 18,
	}
	s.Init(&s, t)
	return &s
}

func (s *IsLegalAdult) IsSatisfiedBy(p User) bool {
	return p.age >= s.ageOfMajority
}

func (s *IsLegalAdult) GetName() string {
	return "IsLegalAdult"
}

// An IsFlagged is a flagged-based TSpecification example, as in "user X has been
// flagged for risky bhavior".
type IsFlagged struct {
	gospec.BaseTSpecification[User]
	flaggedIDs []ID
}

func NewIsFlagged(t gospec.Tracer) *IsFlagged {
	s := IsFlagged{
		flaggedIDs: []ID{"1"},
	}
	s.Init(&s, t)
	return &s
}

func (s *IsFlagged) IsSatisfiedBy(p User) bool {
	for _, id := range s.flaggedIDs {
		if id == p.id {
			return true
		}
	}
	return false
}

func (s *IsFlagged) GetName() string {
	return "IsFlagged"
}

// An IsBob is a name-based TSpecification example.
type IsBob struct {
	gospec.BaseTSpecification[User]
	name string
}

func NewIsBob(t gospec.Tracer) *IsBob {
	s := IsBob{
		name: "Bob",
	}
	s.Init(&s, t)
	return &s
}

func (s *IsBob) IsSatisfiedBy(p User) bool {
	return s.name == p.name
}

func (s *IsBob) GetName() string {
	return "IsBob"
}

// An IsBanned is a composite TSpecification example.
type IsBanned struct {
	gospec.BaseTSpecification[User]
	compositeSpec gospec.TSpecification[User]
}

func NewIsBanned(t gospec.Tracer) *IsBanned {

	adultSpec := NewIsLegalAdult(t)
	flaggedSpec := NewIsFlagged(t)
	bobSpec := NewIsBob(t)

	s := IsBanned{
		compositeSpec: adultSpec.Not().Or(flaggedSpec).Or(bobSpec),
	}
	s.Init(&s, t)
	return &s
}

func (s *IsBanned) IsSatisfiedBy(p User) bool {
	return s.compositeSpec.IsSatisfiedBy(p)
}

func (s *IsBanned) GetName() string {
	return "IsBanned"
}

type testInput struct {
	user User
	spec gospec.TSpecification[User]
}

type testOutput struct {
	result   bool
	opLabels []gospec.OperatorLabel
}

var tracer = th.NewSpecTracer[*gospec.OpTrace]()

var userTestCases = []th.TestCase[testInput, testOutput]{
	{
		"simple",
		testInput{
			User{
				age: 31,
			},
			NewIsLegalAdult(tracer),
		},
		testOutput{
			true,
			[]g.OperatorLabel{},
		},
	},
	{
		"not",
		testInput{
			User{
				age: 16,
			},
			NewIsLegalAdult(tracer).Not(),
		},
		testOutput{
			true,
			[]g.OperatorLabel{g.Not},
		},
	},
	{
		"and: left is true, right is true",
		testInput{
			User{
				id:  "1",
				age: 21,
			},
			NewIsLegalAdult(tracer).And(NewIsFlagged(tracer)),
		},
		testOutput{
			true,
			[]g.OperatorLabel{g.And},
		},
	},
	{
		"and: left is true, right is false",
		testInput{
			User{
				id:  "2",
				age: 21,
			},
			NewIsLegalAdult(tracer).And(NewIsFlagged(tracer)),
		},
		testOutput{
			false,
			[]g.OperatorLabel{g.And},
		},
	},

	{
		"and: left is false, right is true",
		testInput{
			User{
				id:  "1",
				age: 16,
			},
			NewIsLegalAdult(tracer).And(NewIsFlagged(tracer)),
		},
		testOutput{
			false,
			[]g.OperatorLabel{g.And},
		},
	},

	{
		"and: left is false, right is false",
		testInput{
			User{
				id:  "2",
				age: 12,
			},
			NewIsLegalAdult(tracer).And(NewIsFlagged(tracer)),
		},
		testOutput{
			false,
			[]g.OperatorLabel{g.And},
		},
	},

	{
		"or: left is true, right is true",
		testInput{
			User{
				id:  "1",
				age: 21,
			},
			NewIsLegalAdult(tracer).Or(NewIsFlagged(tracer)),
		},
		testOutput{
			true,
			[]g.OperatorLabel{g.Or},
		},
	},

	{
		"or: left is true, right is false",
		testInput{
			User{
				id:  "2",
				age: 21,
			},
			NewIsLegalAdult(tracer).Or(NewIsFlagged(tracer)),
		},
		testOutput{
			true,
			[]g.OperatorLabel{g.Or},
		},
	},

	{
		"or: left is false, right is true",
		testInput{
			User{
				id:  "1",
				age: 16,
			},
			NewIsLegalAdult(tracer).Or(NewIsFlagged(tracer)),
		},
		testOutput{
			true,
			[]g.OperatorLabel{g.Or},
		},
	},

	{
		"or: left is false, right is false",
		testInput{
			User{
				id:  "2",
				age: 16,
			},
			NewIsLegalAdult(tracer).Or(NewIsFlagged(tracer)),
		},
		testOutput{
			false,
			[]g.OperatorLabel{g.Or},
		},
	},

	{
		"not chain",
		testInput{
			User{
				age: 26,
			},
			NewIsLegalAdult(tracer).Not().Not(),
		},
		testOutput{
			true,
			[]g.OperatorLabel{g.Not, g.Not},
		},
	},

	{
		"and chain",
		testInput{
			User{
				id:   "1",
				name: "Bob",
				age:  26,
			},
			NewIsLegalAdult(tracer).And(NewIsFlagged(tracer)).And(NewIsBob(tracer)),
		},
		testOutput{
			true,
			[]g.OperatorLabel{g.And, g.And},
		},
	},

	{
		"or chain",
		testInput{
			User{
				id:   "2",
				name: "Alice",
				age:  12,
			},
			NewIsLegalAdult(tracer).Or(NewIsFlagged(tracer)).Or(NewIsBob(tracer)),
		},
		testOutput{
			false,
			[]g.OperatorLabel{g.Or, g.Or},
		},
	},

	{
		"composite specification",
		testInput{
			User{
				id:   "3",
				name: "Bob",
				age:  34,
			},
			NewIsBanned(tracer),
		},
		testOutput{
			true,
			[]g.OperatorLabel{g.Not, g.Or, g.Or},
		},
	},
}

func TestAnd(t *testing.T) {

	for _, tCase := range userTestCases {

		spec := tCase.Input.spec
		user := tCase.Input.user

		specSatisfied := spec.IsSatisfiedBy(user)

		if specSatisfied != tCase.ExpectedOutput.result {
			t.Errorf(
				"\nDescription: %s\nExpected: %v\nGot: %v\n",
				tCase.Description,
				tCase.ExpectedOutput.result,
				specSatisfied,
			)
		}

		expectedLabels := tCase.ExpectedOutput.opLabels

		tracerEntryLen := tracer.Len()
		expectedTraceLen := len(expectedLabels)

		if tracerEntryLen != expectedTraceLen {
			t.Errorf(
				"\nDescription: %s\nExpected: %v\nGot: %v\n",
				tCase.Description,
				expectedTraceLen,
				tracerEntryLen,
			)
		}

		for i := 0; i < expectedTraceLen; i++ {
			expectedLabel := expectedLabels[i]
			trace := tracer.Get(i)
			traceLabel := trace.GetLabel()

			if traceLabel != expectedLabel {
				t.Errorf(
					"\nDescription: %s\nExpected: %v\nGot: %v\n",
					tCase.Description,
					expectedLabel,
					traceLabel,
				)
			}
		}

		//fmt.Println(tCase.Description)
		//for i := 0; i < tracer.Len(); i++ {
		//	s, err := tracer.Get(i).Explain()
		//	if err != nil {
		//		fmt.Println(err)
		//		continue
		//	}
		//	fmt.Println(s)
		//}

		tracer.Clear()
	}

}
