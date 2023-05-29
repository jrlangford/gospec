package gospec

import (
	"testing"

	th "github.com/jrlangford/gospec/gospectest"
)

type ID string

type User struct {
	id      ID
	name    string
	age     uint8
	heightM float64
}

type IsLegalAdult struct {
	BaseSpecification[User]
	ageOfMajority uint8
}

func NewIsLegalAdult() *IsLegalAdult {
	s := IsLegalAdult{
		ageOfMajority: 18,
	}
	s.Init(&s)
	return &s
}

func (s *IsLegalAdult) IsSatisfiedBy(p User) bool {
	return p.age >= s.ageOfMajority
}

type IsFlagged struct {
	BaseSpecification[User]
	flaggedIDs []ID
}

func NewIsFlagged() *IsFlagged {
	s := IsFlagged{
		flaggedIDs: []ID{"1"},
	}
	s.Init(&s)
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

type IsBob struct {
	BaseSpecification[User]
	name string
}

func NewIsBob() *IsBob {
	s := IsBob{
		name: "Bob",
	}
	s.Init(&s)
	return &s
}

func (s *IsBob) IsSatisfiedBy(p User) bool {
	return s.name == p.name
}

type IsBanned struct {
	BaseSpecification[User]
	compositeSpec Specification[User]
}

func NewIsBanned() *IsBanned {

	adultSpec := NewIsLegalAdult()
	flaggedSpec := NewIsFlagged()
	bobSpec := NewIsBob()

	s := IsBanned{
		compositeSpec: adultSpec.Not().Or(flaggedSpec).Or(bobSpec),
	}
	s.Init(&s)
	return &s
}

func (s *IsBanned) IsSatisfiedBy(p User) bool {
	return s.compositeSpec.IsSatisfiedBy(p)
}

type testInput struct {
	user User
	spec Specification[User]
}

var userTestCases = []th.TestCase[testInput, bool]{
	{
		"simple",
		testInput{
			User{
				age: 31,
			},
			NewIsLegalAdult(),
		},
		true,
	},
	{
		"not",
		testInput{
			User{
				age: 16,
			},
			NewIsLegalAdult().Not(),
		},
		true,
	},
	{
		"and: left is true, right is true",
		testInput{
			User{
				id:  "1",
				age: 21,
			},
			NewIsLegalAdult().And(NewIsFlagged()),
		},
		true,
	},
	{
		"and: left is true, right is false",
		testInput{
			User{
				id:  "2",
				age: 21,
			},
			NewIsLegalAdult().And(NewIsFlagged()),
		},
		false,
	},
	{
		"and: left is false, right is true",
		testInput{
			User{
				id:  "1",
				age: 16,
			},
			NewIsLegalAdult().And(NewIsFlagged()),
		},
		false,
	},
	{
		"and: left is false, right is false",
		testInput{
			User{
				id:  "2",
				age: 12,
			},
			NewIsLegalAdult().And(NewIsFlagged()),
		},
		false,
	},
	{
		"or: left is true, right is true",
		testInput{
			User{
				id:  "1",
				age: 21,
			},
			NewIsLegalAdult().Or(NewIsFlagged()),
		},
		true,
	},
	{
		"or: left is true, right is false",
		testInput{
			User{
				id:  "2",
				age: 21,
			},
			NewIsLegalAdult().Or(NewIsFlagged()),
		},
		true,
	},
	{
		"or: left is false, right is true",
		testInput{
			User{
				id:  "1",
				age: 16,
			},
			NewIsLegalAdult().Or(NewIsFlagged()),
		},
		true,
	},
	{
		"or: left is false, right is false",
		testInput{
			User{
				id:  "2",
				age: 16,
			},
			NewIsLegalAdult().Or(NewIsFlagged()),
		},
		false,
	},
	{
		"not chain",
		testInput{
			User{
				age: 26,
			},
			NewIsLegalAdult().Not().Not(),
		},
		true,
	},
	{
		"and chain",
		testInput{
			User{
				id:   "1",
				name: "Bob",
				age:  26,
			},
			NewIsLegalAdult().And(NewIsFlagged()).And(NewIsBob()),
		},
		true,
	},
	{
		"or chain",
		testInput{
			User{
				id:   "2",
				name: "Alice",
				age:  12,
			},
			NewIsLegalAdult().Or(NewIsFlagged()).Or(NewIsBob()),
		},
		false,
	},
	{
		"composite specification",
		testInput{
			User{
				id:   "3",
				name: "Bob",
				age:  34,
			},
			NewIsBanned(),
		},
		true,
	},
}

func TestAll(t *testing.T) {
	for _, tCase := range userTestCases {

		spec := tCase.Input.spec
		user := tCase.Input.user

		specSatisfied := spec.IsSatisfiedBy(user)

		if specSatisfied != tCase.ExpectedOutput {
			t.Errorf(
				"\nDescription: %s\nExpected: %v\nGot: %v\n",
				tCase.Description,
				tCase.ExpectedOutput,
				specSatisfied,
			)
		}
	}
}
