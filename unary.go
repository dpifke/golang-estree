package estree

import (
	"encoding/json"
	"fmt"
)

type UnaryOperator string

var (
	Negative   UnaryOperator = "-"
	ToNumber   UnaryOperator = "+"
	Not        UnaryOperator = "!"
	BitwiseNot UnaryOperator = "~"
	TypeOf     UnaryOperator = "typeof"
	Void       UnaryOperator = "void"
	Delete     UnaryOperator = "delete"
)

func (uo UnaryOperator) GoString() string {
	switch uo {
	case Negative:
		return "Negative"
	case ToNumber:
		return "ToNumber"
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

type UnaryExpression struct {
	baseExpression
	Operator UnaryOperator
	Prefix   bool
	Argument Expression
}

func (UnaryExpression) Type() string { return "UnaryExpression" }

func (ue UnaryExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     ue.Type(),
		"operator": ue.Operator,
		"prefix":   ue.Prefix,
		"Argument": ue.Argument,
	})
}

func (ue *UnaryExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Operator UnaryOperator   `json:"type"`
		Prefix   bool            `json:"prefix"`
		Argument json.RawMessage `json:"argument"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ue.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ue.Type(), x.Type)
	}
	if err == nil {
		switch x.Operator {
		case Negative, ToNumber, Not, BitwiseNot, TypeOf, Void, Delete:
			ue.Operator = x.Operator
		default:
			err = fmt.Errorf("%w for UnaryExpression.Operator: %q", ErrWrongValue, x.Operator)
		}
	}
	if err == nil {
		ue.Argument, _, err = unmarshalExpression(x.Argument)
	}
	return err
}

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

type UpdateExpression struct {
	baseExpression
	Operator UpdateOperator
	Argument Expression
	Prefix   bool
}

func (UpdateExpression) Type() string { return "UpdateExpression" }

func (ue UpdateExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     ue.Type(),
		"operator": ue.Operator,
		"Argument": ue.Argument,
		"prefix":   ue.Prefix,
	})
}

func (ue *UpdateExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Operator UpdateOperator  `json:"type"`
		Argument json.RawMessage `json:"argument"`
		Prefix   bool            `json:"prefix"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ue.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ue.Type(), x.Type)
	}
	if err == nil {
		switch x.Operator {
		case Increment, Decrement:
			ue.Operator = x.Operator
		default:
			err = fmt.Errorf("%w for UnaryExpression.Operator: %q", ErrWrongValue, x.Operator)
		}
	}
	if err == nil {
		ue.Argument, _, err = unmarshalExpression(x.Argument)
	}
	return err
}
