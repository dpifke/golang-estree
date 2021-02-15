package estree

import (
	"testing"
)

func TestUnaryExpression(t *testing.T) {
	var ue UnaryExpression
	if !ue.IsZero() {
		t.Error("expected IsZero()")
	}

	ue.Operator = TypeOf
	ue.Argument = ThisExpression{}
	ue.Prefix = true
	if ue.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ue.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ue.MinVersion())
	}

	var v mockVisitor
	ue.Walk(&v)
	v.expect(t, ue, ue.Argument, nil, nil)

	testRoundtripJSON(t, ue, new(UnaryExpression))

	if errs := ue.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ue.Operator = "not valid"
	if !hasError(ErrWrongValue, ue.Errors()...) {
		t.Error("expected ErrWrongValue for invalid UnaryOperator")
	}
	ue.Operator = Delete
	ue.Argument = Identifier{}
	if !hasError(ErrMissingNode, ue.Errors()...) {
		t.Error("expected ErrMissingNode for zero Argument")
	}
	ue.Argument = nil
	if !hasError(ErrMissingNode, ue.Errors()...) {
		t.Error("expected ErrMissingNode for nil Argument")
	}
}

func TestUpdateExpression(t *testing.T) {
	var ue UpdateExpression
	if !ue.IsZero() {
		t.Error("expected IsZero()")
	}

	ue.Operator = Decrement
	ue.Argument = Identifier{Name: "foo"}
	ue.Prefix = true
	if ue.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ue.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ue.MinVersion())
	}

	var v mockVisitor
	ue.Walk(&v)
	v.expect(t, ue, ue.Argument, nil, nil)

	testRoundtripJSON(t, ue, new(UpdateExpression))

	if errs := ue.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ue.Operator = "not valid"
	if !hasError(ErrWrongValue, ue.Errors()...) {
		t.Error("expected ErrWrongValue for invalid UpdateOperator")
	}
	ue.Operator = Increment
	ue.Argument = Identifier{}
	if !hasError(ErrMissingNode, ue.Errors()...) {
		t.Error("expected ErrMissingNode for zero Argument")
	}
	ue.Argument = nil
	if !hasError(ErrMissingNode, ue.Errors()...) {
		t.Error("expected ErrMissingNode for nil Argument")
	}
}
