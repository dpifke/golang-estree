package estree

import (
	"encoding/json"
	"fmt"
)

// UnaryOperator is the operator token of an UnaryExpression, which modifies a
// single operand.
type UnaryOperator string

var (
	Minus      UnaryOperator = "-"
	Plus       UnaryOperator = "+"
	Not        UnaryOperator = "!"
	BitwiseNot UnaryOperator = "~"
	TypeOf     UnaryOperator = "typeof"
	Void       UnaryOperator = "void"
	Delete     UnaryOperator = "delete"
)

func (uo UnaryOperator) GoString() string {
	switch uo {
	case Minus:
		return "Minus"
	case Plus:
		return "Plus"
	case Not:
		return "Not"
	case BitwiseNot:
		return "BitwiseNot"
	case TypeOf:
		return "TypeOf"
	case Void:
		return "Void"
	case Delete:
		return "Delete"
	}
	return fmt.Sprintf("%q", uo)
}

func (uo UnaryOperator) IsValid() bool {
	switch uo {
	case Minus, Plus, Not, BitwiseNot, TypeOf, Void, Delete:
		return true
	}
	return false
}

// UnaryExpression is an expression modifying a single operand.
type UnaryExpression struct {
	baseExpression
	Loc      SourceLocation
	Operator UnaryOperator
	Prefix   bool
	Argument Expression
}

func (UnaryExpression) Type() string                { return "UnaryExpression" }
func (ue UnaryExpression) Location() SourceLocation { return ue.Loc }

func (ue UnaryExpression) IsZero() bool {
	return ue.Loc.IsZero() &&
		ue.Operator == "" &&
		!ue.Prefix &&
		(ue.Argument == nil || ue.Argument.IsZero())
}

func (ue UnaryExpression) Walk(v Visitor) {
	if v = v.Visit(ue); v != nil {
		defer v.Visit(nil)
		if ue.Argument != nil {
			ue.Argument.Walk(v)
		}
	}
}

func (ue UnaryExpression) Errors() []error {
	c := nodeChecker{Node: ue}
	if ue.Prefix {
		c.require(ue.Argument, "unary argument")
	}
	if !ue.Operator.IsValid() {
		c.appendf("%w UnaryOperator %q", ErrWrongValue, ue.Operator)
	}
	if !ue.Prefix {
		c.require(ue.Argument, "unary argument")
	}
	return c.errors()
}

func (ue UnaryExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ue)
	x["operator"] = ue.Operator
	x["prefix"] = ue.Prefix
	x["argument"] = ue.Argument
	return json.Marshal(x)
}

func (ue *UnaryExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Loc      SourceLocation  `json:"loc"`
		Operator UnaryOperator   `json:"operator"`
		Prefix   bool            `json:"prefix"`
		Argument json.RawMessage `json:"argument"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ue.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, ue.Type(), x.Type)
	}
	if err == nil {
		ue.Loc, ue.Prefix = x.Loc, x.Prefix
		if x.Operator.IsValid() {
			ue.Operator = x.Operator
		} else {
			err = fmt.Errorf("%w UnaryExpression.Operator %q", ErrWrongValue, x.Operator)
		}
		var err2 error
		if ue.Argument, _, err = unmarshalExpression(x.Argument); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

// UpdateOperator is the operator token of an UpdateExpression, which modifies
// a single operand in-place.
type UpdateOperator string

var (
	Increment UpdateOperator = "++"
	Decrement UpdateOperator = "--"
)

func (uo UpdateOperator) GoString() string {
	switch uo {
	case Increment:
		return "Increment"
	case Decrement:
		return "Decrement"
	}
	return fmt.Sprintf("%q", uo)
}

func (uo UpdateOperator) IsValid() bool {
	switch uo {
	case Increment, Decrement:
		return true
	}
	return false
}

// UpdateExpression is an expression which modifies a single operand in-place.
type UpdateExpression struct {
	baseExpression
	Loc      SourceLocation
	Operator UpdateOperator
	Argument Expression
	Prefix   bool
}

func (UpdateExpression) Type() string                { return "UpdateExpression" }
func (ue UpdateExpression) Location() SourceLocation { return ue.Loc }

func (ue UpdateExpression) IsZero() bool {
	return ue.Loc.IsZero() &&
		ue.Operator == "" &&
		(ue.Argument == nil || ue.Argument.IsZero()) &&
		!ue.Prefix
}

func (ue UpdateExpression) Walk(v Visitor) {
	if v = v.Visit(ue); v != nil {
		defer v.Visit(nil)
		if ue.Argument != nil {
			ue.Argument.Walk(v)
		}
	}
}

func (ue UpdateExpression) Errors() []error {
	c := nodeChecker{Node: ue}
	if ue.Prefix {
		c.require(ue.Argument, "update argument")
	}
	if !ue.Operator.IsValid() {
		c.appendf("%w UpdateOperator %q", ErrWrongValue, ue.Operator)
	}
	if !ue.Prefix {
		c.require(ue.Argument, "update argument")
	}
	return c.errors()
}

func (ue UpdateExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ue)
	x["operator"] = ue.Operator
	x["argument"] = ue.Argument
	x["prefix"] = ue.Prefix
	return json.Marshal(x)
}

func (ue *UpdateExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Loc      SourceLocation  `json:"loc"`
		Operator UpdateOperator  `json:"operator"`
		Argument json.RawMessage `json:"argument"`
		Prefix   bool            `json:"prefix"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ue.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, ue.Type(), x.Type)
	}
	if err == nil {
		ue.Loc, ue.Prefix = x.Loc, x.Prefix
		if x.Operator.IsValid() {
			ue.Operator = x.Operator
		} else {
			err = fmt.Errorf("%w UnaryExpression.Operator %q", ErrWrongValue, x.Operator)
		}
		var err2 error
		if ue.Argument, _, err2 = unmarshalExpression(x.Argument); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}
