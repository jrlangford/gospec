package gospec

import (
	"errors"
	"fmt"
)

type OperatorLabel string

const (
	And OperatorLabel = "AND"
	Or                = "OR"
	Not               = "NOT"
)

type ExpressionTrace struct {
	ExpressionName  string
	ExpressionValue bool
}

func NewExpressionTrace(name string, val bool) *ExpressionTrace {
	return &ExpressionTrace{
		ExpressionName:  name,
		ExpressionValue: val,
	}
}

type OpTrace struct {
	opLabel OperatorLabel
	left    *ExpressionTrace
	right   *ExpressionTrace
}

func NewOpTraceWithLeft(opLabel OperatorLabel, left *ExpressionTrace) *OpTrace {
	return &OpTrace{
		opLabel: opLabel,
		left:    left,
		right:   nil,
	}
}

func (o *OpTrace) SetRight(right *ExpressionTrace) {
	o.right = right
}

func (o *OpTrace) GetLabel() OperatorLabel {
	return o.opLabel
}

func (o *OpTrace) Explain() (string, error) {

	if o.left == nil {
		return "", errors.New("Received nil left.")
	}

	switch o.opLabel {
	case And:
		s := fmt.Sprintf("[left %s right]", o.opLabel)

		leftVal := o.left.ExpressionValue
		s += fmt.Sprintf(" > left: %s evaluates to %t", o.left.ExpressionName, leftVal)
		if !leftVal {
			s += fmt.Sprintf(" > false AND X is false through short-circuit")
			return s, nil
		}

		if o.right == nil {
			return "", errors.New("Received nil right.")
		}

		rightVal := o.right.ExpressionValue
		s += fmt.Sprintf(" > right: %s evaluates to %t", o.right.ExpressionName, rightVal)

		result := leftVal && rightVal

		s += fmt.Sprintf(" > %t AND %t is %t", leftVal, rightVal, result)
		return s, nil
	case Or:
		s := fmt.Sprintf("[left %s right]", o.opLabel)

		leftVal := o.left.ExpressionValue
		s += fmt.Sprintf(" > left: %s evaluates to %t", o.left.ExpressionName, leftVal)
		if leftVal {
			s += fmt.Sprintf(" > true OR X is true through short-circuit")
			return s, nil
		}

		if o.right == nil {
			return "", errors.New("Received nil right.")
		}

		rightVal := o.right.ExpressionValue
		s += fmt.Sprintf(" > right: %s evaluates to %t", o.right.ExpressionName, rightVal)

		result := leftVal || rightVal

		s += fmt.Sprintf(" > %t OR %t is %t", leftVal, rightVal, result)
		return s, nil
	case Not:
		s := fmt.Sprintf("[%s left]", o.opLabel)

		leftVal := o.left.ExpressionValue
		s += fmt.Sprintf(" > left: %s evaluates to %t", o.left.ExpressionName, leftVal)

		result := !leftVal

		s += fmt.Sprintf(" > NOT %t is %t", leftVal, result)
		return s, nil
	}
	return "", errors.New("Unidentified operation.")
}
