package estree

import (
	"testing"
)

func TestBinaryExpression(t *testing.T) {
	var be BinaryExpression
	if !be.IsZero() {
		t.Error("expected IsZero()")
	}

	be.Operator = Equal
	be.Left = NumberLiteral{Value: 2.0}
	be.Right = NumberLiteral{Value: 2.0}
	if be.IsZero() {
		t.Error("expected !IsZero()")
	}
	if be.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", be.MinVersion())
	}

	var v mockVisitor
	be.Walk(&v)
	v.expect(t, be,
		be.Left, nil,
		be.Right, nil, nil)

	testRoundtripJSON(t, be, new(BinaryExpression))

	if errs := be.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	be.Operator = "not valid"
	if !hasError(ErrWrongValue, be.Errors()...) {
		t.Error("expected ErrWrongValue for invalid BinaryOperator")
	}
	be.Operator = NotEqual
	be.Left = nil
	if !hasError(ErrMissingNode, be.Errors()...) {
		t.Error("expected ErrMissingNode for nil Left")
	}
	be.Left = StringLiteral{Value: "foo"}
	be.Right = FunctionExpression{}
	if !hasError(ErrMissingNode, be.Errors()...) {
		t.Error("expected ErrMissingNode for zero Right")
	}
}

func TestAssignmentExpression(t *testing.T) {
	var ae AssignmentExpression
	if !ae.IsZero() {
		t.Error("expected IsZero()")
	}

	ae.Operator = AddAssign
	ae.Left = Identifier{Name: "foo"}
	ae.Right = NumberLiteral{Value: 2.0}
	if ae.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ae.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ae.MinVersion())
	}

	var v mockVisitor
	ae.Walk(&v)
	v.expect(t, ae,
		ae.Left, nil,
		ae.Right, nil, nil)

	testRoundtripJSON(t, ae, new(AssignmentExpression))

	if errs := ae.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ae.Operator = "not valid"
	if !hasError(ErrWrongValue, ae.Errors()...) {
		t.Error("expected ErrWrongValue for invalid AssignmentOperator")
	}
	ae.Operator = SubtractAssign
	ae.Left = MemberExpression{}
	if !hasError(ErrMissingNode, ae.Errors()...) {
		t.Error("expected ErrMissingNode for zero Left")
	}
	ae.Left = MemberExpression{
		Object:   ObjectExpression{},
		Property: Identifier{Name: "foo"},
	}
	ae.Right = nil
	if !hasError(ErrMissingNode, ae.Errors()...) {
		t.Error("expected ErrMissingNode for nil Right")
	}
}

func TestLogicalExpression(t *testing.T) {
	var le LogicalExpression
	if !le.IsZero() {
		t.Error("expected IsZero()")
	}

	le.Operator = And
	le.Left = Identifier{Name: "foo"}
	le.Right = BoolLiteral{Value: false}
	if le.IsZero() {
		t.Error("expected !IsZero()")
	}
	if le.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", le.MinVersion())
	}

	var v mockVisitor
	le.Walk(&v)
	v.expect(t, le,
		le.Left, nil,
		le.Right, nil, nil)

	testRoundtripJSON(t, le, new(LogicalExpression))

	if errs := le.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	le.Operator = "not valid"
	if !hasError(ErrWrongValue, le.Errors()...) {
		t.Error("expected ErrWrongValue for invalid LogicalOperator")
	}
	le.Operator = Or
	le.Left = ConditionalExpression{}
	if !hasError(ErrMissingNode, le.Errors()...) {
		t.Error("expected ErrMissingNode for zero Left")
	}
	le.Left = ConditionalExpression{
		Test: CallExpression{
			Callee: Identifier{Name: "foo"},
		},
		Consequent: StringLiteral{Value: "foo"},
		Alternate:  StringLiteral{Value: "bar"},
	}
	le.Right = nil
	if !hasError(ErrMissingNode, le.Errors()...) {
		t.Error("expected ErrMissingNode for nil Right")
	}
}

func TestMemberExpression(t *testing.T) {
	var me MemberExpression
	if !me.IsZero() {
		t.Error("expected IsZero()")
	}

	me.Object = ThisExpression{}
	me.Property = CallExpression{
		Callee: Identifier{Name: "foo"},
	}
	me.Computed = true
	if me.IsZero() {
		t.Error("expected !IsZero()")
	}
	if me.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", me.MinVersion())
	}

	var v mockVisitor
	me.Walk(&v)
	v.expect(t, me,
		me.Object, nil,
		me.Property, me.Property.(CallExpression).Callee, nil, nil, nil)

	testRoundtripJSON(t, me, new(MemberExpression))

	if errs := me.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	me.Object = SequenceExpression{}
	if !hasError(ErrMissingNode, me.Errors()...) {
		t.Error("expected ErrMissingNode for zero Left")
	}
	me.Property = ArrayExpression{}
	me.Property = nil
	if !hasError(ErrMissingNode, me.Errors()...) {
		t.Error("expected ErrMissingNode for nil Right")
	}
}
