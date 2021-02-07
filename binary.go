package estree

import (
	"encoding/json"
	"fmt"
)

type BinaryOperator string

var (
	Equal              BinaryOperator = "=="
	NotEqual           BinaryOperator = "!="
	Identical          BinaryOperator = "==="
	NotIdentical       BinaryOperator = "!=="
	Lesser             BinaryOperator = "<"
	LesserEqual        BinaryOperator = "<="
	Greater            BinaryOperator = ">"
	GreaterEqual       BinaryOperator = ">="
	ShiftLeft          BinaryOperator = "<<"
	ShiftRight         BinaryOperator = ">>"
	UnsignedShiftRight BinaryOperator = ">>>"
	Add                BinaryOperator = "+"
	Subtract           BinaryOperator = "-"
	Multiply           BinaryOperator = "*"
	Divide             BinaryOperator = "/"
	Modulo             BinaryOperator = "%"
	BitwiseOr          BinaryOperator = "|"
	BitwiseXor         BinaryOperator = "^"
	BitwiseAnd         BinaryOperator = "&"
	In                 BinaryOperator = "in"
	InstanceOf         BinaryOperator = "instanceof"
)

func (bo BinaryOperator) GoString() string {
	switch bo {
	case Equal:
		return "Equal"
	case NotEqual:
		return "NotEqual"
	case Identical:
		return "Identical"
	case NotIdentical:
		return "NotIdentical"
	case Lesser:
		return "Lesser"
	case LesserEqual:
		return "LesserEqual"
	case Greater:
		return "Greater"
	case GreaterEqual:
		return "GreaterEqual"
	case ShiftLeft:
		return "ShiftLeft"
	case ShiftRight:
		return "ShiftRight"
	case UnsignedShiftRight:
		return "UnsignedShiftRight"
	case Add:
		return "Add"
	case Subtract:
		return "Subtract"
	case Multiply:
		return "Multiply"
	case Divide:
		return "Divide"
	case Modulo:
		return "Modulo"
	case BitwiseOr:
		return "BitwiseOr"
	case BitwiseXor:
		return "BitwiseXor"
	case BitwiseAnd:
		return "BitwiseAnd"
	case In:
		return "In"
	case InstanceOf:
		return "InstanceOf"
	}
	return fmt.Sprintf("%q", bo)
}

type BinaryExpression struct {
	baseExpression
	Operator    BinaryOperator
	Left, Right Expression
}

func (BinaryExpression) Type() string { return "BinaryExpression" }

func (be BinaryExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     be.Type(),
		"operator": be.Operator,
		"left":     be.Left,
		"right":    be.Right,
	})
}

func (be *BinaryExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Operator BinaryOperator  `json:"operator"`
		Left     json.RawMessage `json:"left"`
		Right    json.RawMessage `json:"right"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != be.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, be.Type(), x.Type)
	}
	if err == nil {
		switch x.Operator {
		case Equal, NotEqual, Identical, NotIdentical, Lesser, LesserEqual,
			Greater, GreaterEqual, ShiftLeft, ShiftRight, UnsignedShiftRight,
			Add, Subtract, Multiply, Divide, Modulo, BitwiseOr, BitwiseXor,
			BitwiseAnd, In, InstanceOf:
			be.Operator = x.Operator
		default:
			err = fmt.Errorf("%w for BinaryExpression.Operator: %q", ErrWrongValue, x.Operator)
		}
	}
	if err == nil {
		be.Left, _, err = unmarshalExpression(x.Left)
	}
	if err == nil {
		be.Right, _, err = unmarshalExpression(x.Right)
	}
	return err
}

type AssignmentOperator string

var (
	Assign                   AssignmentOperator = "="
	AddAssign                AssignmentOperator = "+="
	SubtractAssign           AssignmentOperator = "-="
	MultiplyAssign           AssignmentOperator = "*="
	DivideAssign             AssignmentOperator = "/="
	ModuloAssign             AssignmentOperator = "%="
	ShiftLeftAssign          AssignmentOperator = "<<="
	ShiftRightAssign         AssignmentOperator = ">>="
	UnsignedShiftRightAssign AssignmentOperator = ">>>="
	BitwiseOrAssign          AssignmentOperator = "|="
	BitwiseXorAssign         AssignmentOperator = "^="
	BitwiseAndAssign         AssignmentOperator = "&="
)

func (ao AssignmentOperator) GoString() string {
	switch ao {
	case Assign:
		return "Assign"
	case AddAssign:
		return "AddAssign"
	case SubtractAssign:
		return "SubtractAssign"
	case MultiplyAssign:
		return "MultiplyAssign"
	case DivideAssign:
		return "DivideAssign"
	case ModuloAssign:
		return "ModuloAssign"
	case ShiftLeftAssign:
		return "ShiftLeftAssign"
	case ShiftRightAssign:
		return "ShiftRightAssign"
	case UnsignedShiftRightAssign:
		return "UnsignedShiftRightAssign"
	case BitwiseOrAssign:
		return "BitwiseOrAssign"
	case BitwiseXorAssign:
		return "BitwiseXorAssign"
	case BitwiseAndAssign:
		return "BitwiseAndAssign"
	}
	return fmt.Sprintf("%q", ao)
}

type AssignmentExpression struct {
	baseExpression
	Operator AssignmentOperator
	Left     PatternOrExpression
	Right    Expression
}

func (AssignmentExpression) Type() string { return "AssignmentExpression" }

func (ae AssignmentExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     ae.Type(),
		"operator": ae.Operator,
		"left":     ae.Left,
		"right":    ae.Right,
	})
}

func (ae *AssignmentExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string             `json:"type"`
		Operator AssignmentOperator `json:"operator"`
		Left     json.RawMessage    `json:"left"`
		Right    json.RawMessage    `json:"right"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ae.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ae.Type(), x.Type)
	}
	if err == nil {
		switch x.Operator {
		case Assign, AddAssign, SubtractAssign, MultiplyAssign, DivideAssign,
			ModuloAssign, ShiftLeftAssign, ShiftRightAssign,
			UnsignedShiftRightAssign, BitwiseOrAssign, BitwiseXorAssign,
			BitwiseAndAssign:
			ae.Operator = x.Operator
		default:
			err = fmt.Errorf("%w for AssignmentExpression.Operator: %q", ErrWrongValue, x.Operator)
		}
	}
	if err == nil {
		ae.Left, err = unmarshalPatternOrExpression(x.Left)
	}
	if err == nil {
		ae.Right, _, err = unmarshalExpression(x.Right)
	}
	return err
}

type LogicalOperator string

var (
	Or  LogicalOperator = "||"
	And LogicalOperator = "&&"
)

func (lo LogicalOperator) GoString() string {
	switch lo {
	case Or:
		return "Or"
	case And:
		return "And"
	}
	return fmt.Sprintf("%q", lo)
}

type LogicalExpression struct {
	baseExpression
	Operator    LogicalOperator
	Left, Right Expression
}

func (LogicalExpression) Type() string { return "LogicalExpression" }

func (le LogicalExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     le.Type(),
		"operator": le.Operator,
		"left":     le.Left,
		"right":    le.Right,
	})
}

func (le *LogicalExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Operator LogicalOperator `json:"operator"`
		Left     json.RawMessage `json:"left"`
		Right    json.RawMessage `json:"right"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != le.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, le.Type(), x.Type)
	}
	if err == nil {
		switch x.Operator {
		case Or, And:
			le.Operator = x.Operator
		default:
			err = fmt.Errorf("%w for LogicalExpression.Operator: %q", ErrWrongValue, x.Operator)
		}
	}
	if err == nil {
		le.Left, _, err = unmarshalExpression(x.Left)
	}
	if err == nil {
		le.Right, _, err = unmarshalExpression(x.Right)
	}
	return err
}

type MemberExpression struct {
	baseExpression
	Object   Expression
	Property Expression
	Computed bool
}

func (MemberExpression) Type() string { return "MemberExpression" }

func (me MemberExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     me.Type(),
		"object":   me.Object,
		"property": me.Property,
		"computed": me.Computed,
	})
}

func (me *MemberExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Object   json.RawMessage `json:"object"`
		Property json.RawMessage `json:"property"`
		Computed bool            `json:"computed"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != me.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, me.Type(), x.Type)
	}
	if err == nil {
		me.Object, _, err = unmarshalExpression(x.Object)
	}
	if err == nil {
		me.Property, _, err = unmarshalExpression(x.Property)
	}
	return err
}
