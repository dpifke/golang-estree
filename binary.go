package estree

import (
	"encoding/json"
	"fmt"
)

// BinaryOperator is the operator token of a BinaryExpression, which takes two
// operands and performs operations such as arithmetic or comparisons.
type BinaryOperator string

var (
	Equal              BinaryOperator = "=="
	NotEqual           BinaryOperator = "!="
	StrictEqual        BinaryOperator = "==="
	StrictNotEqual     BinaryOperator = "!=="
	LessThan           BinaryOperator = "<"
	LessThanOrEqual    BinaryOperator = "<="
	GreaterThan        BinaryOperator = ">"
	GreaterThanOrEqual BinaryOperator = ">="
	LeftShift          BinaryOperator = "<<"
	SignedRightShift   BinaryOperator = ">>"
	UnsignedRightShift BinaryOperator = ">>>"
	Add                BinaryOperator = "+"
	Subtract           BinaryOperator = "-"
	Multiply           BinaryOperator = "*"
	Divide             BinaryOperator = "/"
	Remainder          BinaryOperator = "%"
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
	case StrictEqual:
		return "StrictEqual"
	case StrictNotEqual:
		return "StrictNotEqual"
	case LessThan:
		return "LessThan"
	case LessThanOrEqual:
		return "LessThanOrEqual"
	case GreaterThan:
		return "GreaterThan"
	case GreaterThanOrEqual:
		return "GreaterThanOrEqual"
	case LeftShift:
		return "LeftShift"
	case SignedRightShift:
		return "SignedRightShift"
	case UnsignedRightShift:
		return "UnsignedRightShift"
	case Add:
		return "Add"
	case Subtract:
		return "Subtract"
	case Multiply:
		return "Multiply"
	case Divide:
		return "Divide"
	case Remainder:
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

func (bo BinaryOperator) IsValid() bool {
	switch bo {
	case Equal, NotEqual, StrictEqual, StrictNotEqual, LessThan,
		LessThanOrEqual, GreaterThan, GreaterThanOrEqual, LeftShift,
		SignedRightShift, UnsignedRightShift, Add, Subtract, Multiply, Divide,
		Remainder, BitwiseOr, BitwiseXor, BitwiseAnd, In, InstanceOf:
		return true
	}
	return false
}

func (bo BinaryOperator) MinVersion() Version { return ES5 }

// BinaryExpression is a binary (two operand) expression.
type BinaryExpression struct {
	baseExpression
	Loc         SourceLocation
	Operator    BinaryOperator
	Left, Right Expression
}

func (BinaryExpression) Type() string                { return "BinaryExpression" }
func (be BinaryExpression) Location() SourceLocation { return be.Loc }

func (be BinaryExpression) MinVersion() Version {
	return be.Operator.MinVersion()
}

func (be BinaryExpression) IsZero() bool {
	return be.Loc.IsZero() &&
		be.Operator == "" &&
		(be.Left == nil || be.Left.IsZero()) &&
		(be.Right == nil || be.Right.IsZero())
}

func (be BinaryExpression) Walk(v Visitor) {
	if v = v.Visit(be); v != nil {
		defer v.Visit(nil)
		if be.Left != nil {
			be.Left.Walk(v)
		}
		if be.Right != nil {
			be.Right.Walk(v)
		}
	}
}

func (be BinaryExpression) Errors() []error {
	c := nodeChecker{Node: be}
	c.require(be.Left, "left-hand expression")
	if !be.Operator.IsValid() {
		c.appendf("%w binary operator %q", ErrWrongValue, be.Operator)
	}
	c.require(be.Right, "right-hand expression")
	return c.errors()
}

func (be BinaryExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(be)
	x["operator"] = be.Operator
	x["left"] = be.Left
	x["right"] = be.Right
	return json.Marshal(x)
}

func (be *BinaryExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Loc      SourceLocation  `json:"loc"`
		Operator BinaryOperator  `json:"operator"`
		Left     json.RawMessage `json:"left"`
		Right    json.RawMessage `json:"right"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != be.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, be.Type(), x.Type)
	}
	if err == nil {
		be.Loc = x.Loc
		if x.Operator.IsValid() {
			be.Operator = x.Operator
		} else {
			err = fmt.Errorf("%w BinaryExpression.Operator %q", ErrWrongValue, x.Operator)
		}
		var err2 error
		if be.Left, _, err2 = unmarshalExpression(x.Left); err == nil && err2 != nil {
			err = err2
		}
		if be.Right, _, err = unmarshalExpression(x.Right); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

// AssignmentOperator is the operator token for an AssignmentExpression, which
// assigns a right-hand operand to a left-hand operand.
type AssignmentOperator string

var (
	Assign                   AssignmentOperator = "="
	AddAssign                AssignmentOperator = "+="
	SubtractAssign           AssignmentOperator = "-="
	MultiplyAssign           AssignmentOperator = "*="
	DivideAssign             AssignmentOperator = "/="
	RemainderAssign          AssignmentOperator = "%="
	LeftShiftAssign          AssignmentOperator = "<<="
	SignedRightShiftAssign   AssignmentOperator = ">>="
	UnsignedRightShiftAssign AssignmentOperator = ">>>="
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
	case RemainderAssign:
		return "RemainderAssign"
	case LeftShiftAssign:
		return "LeftShiftAssign"
	case SignedRightShiftAssign:
		return "SignedRightShiftAssign"
	case UnsignedRightShiftAssign:
		return "UnsignedRightShiftAssign"
	case BitwiseOrAssign:
		return "BitwiseOrAssign"
	case BitwiseXorAssign:
		return "BitwiseXorAssign"
	case BitwiseAndAssign:
		return "BitwiseAndAssign"
	}
	return fmt.Sprintf("%q", ao)
}

func (ao AssignmentOperator) IsValid() bool {
	switch ao {
	case Assign, AddAssign, SubtractAssign, MultiplyAssign, DivideAssign,
		RemainderAssign, LeftShiftAssign, SignedRightShiftAssign,
		UnsignedRightShiftAssign, BitwiseOrAssign, BitwiseXorAssign,
		BitwiseAndAssign:
		return true
	}
	return false
}

// AssignmentExpression is an expression modifying the Left operand according
// to the Right.
type AssignmentExpression struct {
	baseExpression
	Loc      SourceLocation
	Operator AssignmentOperator
	Left     PatternOrExpression
	Right    Expression
}

func (AssignmentExpression) Type() string                { return "AssignmentExpression" }
func (ae AssignmentExpression) Location() SourceLocation { return ae.Loc }

func (ae AssignmentExpression) IsZero() bool {
	return ae.Loc.IsZero() &&
		ae.Operator == "" &&
		(ae.Left == nil || ae.Left.IsZero()) &&
		(ae.Right == nil || ae.Right.IsZero())
}

func (ae AssignmentExpression) Walk(v Visitor) {
	if v = v.Visit(ae); v != nil {
		defer v.Visit(nil)
		if ae.Left != nil {
			ae.Left.Walk(v)
		}
		if ae.Right != nil {
			ae.Right.Walk(v)
		}
	}
}

func (ae AssignmentExpression) Errors() []error {
	c := nodeChecker{Node: ae}
	c.require(ae.Left, "left-hand expression in assignment")
	if !ae.Operator.IsValid() {
		c.appendf("%w assignment operator %q", ErrWrongValue, ae.Operator)
	}
	c.require(ae.Right, "right-hand expression in assignment")
	return c.errors()
}

func (ae AssignmentExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ae)
	x["operator"] = ae.Operator
	x["left"] = ae.Left
	x["right"] = ae.Right
	return json.Marshal(x)
}

func (ae *AssignmentExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string             `json:"type"`
		Loc      SourceLocation     `json:"loc"`
		Operator AssignmentOperator `json:"operator"`
		Left     json.RawMessage    `json:"left"`
		Right    json.RawMessage    `json:"right"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ae.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, ae.Type(), x.Type)
	}
	if err == nil {
		ae.Loc = x.Loc
		if x.Operator.IsValid() {
			ae.Operator = x.Operator
		} else {
			err = fmt.Errorf("%w AssignmentExpression.Operator %q", ErrWrongValue, x.Operator)
		}
		var err2 error
		if ae.Left, err2 = unmarshalPatternOrExpression(x.Left); err == nil && err2 != nil {
			err = err2
		}
		if ae.Right, _, err = unmarshalExpression(x.Right); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

// LogicalOperator is the operator token for a LogicalExpression, which
// expresses Boolean logic.
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

func (lo LogicalOperator) IsValid() bool {
	switch lo {
	case Or, And:
		return true
	}
	return false
}

// LogicalExpression is an expression evaluating Boolean logic between two
// operands.
type LogicalExpression struct {
	baseExpression
	Loc         SourceLocation
	Operator    LogicalOperator
	Left, Right Expression
}

func (LogicalExpression) Type() string                { return "LogicalExpression" }
func (le LogicalExpression) Location() SourceLocation { return le.Loc }

func (le LogicalExpression) IsZero() bool {
	return le.Loc.IsZero() &&
		le.Operator == "" &&
		(le.Left == nil || le.Left.IsZero()) &&
		(le.Right == nil || le.Right.IsZero())
}

func (le LogicalExpression) Walk(v Visitor) {
	if v = v.Visit(le); v != nil {
		defer v.Visit(nil)
		if le.Left != nil {
			le.Left.Walk(v)
		}
		if le.Right != nil {
			le.Right.Walk(v)
		}
	}
}

func (le LogicalExpression) Errors() []error {
	c := nodeChecker{Node: le}
	c.require(le.Left, "left-hand expression")
	if !le.Operator.IsValid() {
		c.appendf("%w logical operator %q", ErrWrongValue, le.Operator)
	}
	c.require(le.Right, "right-hand expression")
	return c.errors()
}

func (le LogicalExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(le)
	x["operator"] = le.Operator
	x["left"] = le.Left
	x["right"] = le.Right
	return json.Marshal(x)
}

func (le *LogicalExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Loc      SourceLocation  `json:"loc"`
		Operator LogicalOperator `json:"operator"`
		Left     json.RawMessage `json:"left"`
		Right    json.RawMessage `json:"right"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != le.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, le.Type(), x.Type)
	}
	if err == nil {
		le.Loc = x.Loc
		if x.Operator.IsValid() {
			le.Operator = x.Operator
		} else {
			err = fmt.Errorf("%w LogicalExpression.Operator %q", ErrWrongValue, x.Operator)
		}
		var err2 error
		if le.Left, _, err2 = unmarshalExpression(x.Left); err == nil && err2 != nil {
			err = err2
		}
		if le.Right, _, err = unmarshalExpression(x.Right); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

// MemberExpression is a member expression, returning a value contained within
// Object, identified by Property.
type MemberExpression struct {
	baseExpression
	Loc      SourceLocation
	Object   Expression
	Property Expression

	// Computed indicates the node corresponds to a computed (a[b]) member
	// expression, and Property is an Expression.  If Computed is false, the
	// node corresponds to a static (a.b) member expression and Property is an
	// Identifier.
	Computed bool
}

func (MemberExpression) Type() string                { return "MemberExpression" }
func (me MemberExpression) Location() SourceLocation { return me.Loc }

func (me MemberExpression) IsZero() bool {
	return me.Loc.IsZero() &&
		(me.Object == nil || me.Object.IsZero()) &&
		(me.Property == nil || me.Property.IsZero())
}

func (me MemberExpression) Walk(v Visitor) {
	if v = v.Visit(me); v != nil {
		defer v.Visit(nil)
		if me.Object != nil {
			me.Object.Walk(v)
		}
		if me.Property != nil {
			me.Property.Walk(v)
		}
	}
}

func (me MemberExpression) Errors() []error {
	c := nodeChecker{Node: me}
	c.require(me.Object, "object in member expression")
	c.require(me.Property, "property or index in member expression")
	return c.errors()
}

func (me MemberExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(me)
	x["object"] = me.Object
	x["property"] = me.Property
	x["computed"] = me.Computed
	return json.Marshal(x)
}

func (me *MemberExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Loc      SourceLocation  `json:"loc"`
		Object   json.RawMessage `json:"object"`
		Property json.RawMessage `json:"property"`
		Computed bool            `json:"computed"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != me.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, me.Type(), x.Type)
	}
	if err == nil {
		me.Loc, me.Computed = x.Loc, x.Computed
		me.Object, _, err = unmarshalExpression(x.Object)
		var err2 error
		if me.Property, _, err2 = unmarshalExpression(x.Property); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}
