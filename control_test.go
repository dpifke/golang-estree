package estree

import (
	"testing"
)

func TestReturnStatement(t *testing.T) {
	var rs ReturnStatement
	if rs.IsZero() {
		t.Error("expected !IsZero()")
	}

	rs.Argument = ThisExpression{}
	if rs.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", rs.MinVersion())
	}

	var v mockVisitor
	rs.Walk(&v)
	v.expect(t, rs, rs.Argument, nil, nil)

	testRoundtripJSON(t, rs, new(ReturnStatement))

	if errs := rs.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestLabeledStatement(t *testing.T) {
	var ls LabeledStatement
	if !ls.IsZero() {
		t.Error("expected IsZero()")
	}

	ls.Label = Identifier{Name: "foo"}
	ls.Body = EmptyStatement{}
	if ls.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ls.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ls.MinVersion())
	}

	var v mockVisitor
	ls.Walk(&v)
	v.expect(t, ls,
		ls.Label, nil,
		ls.Body, nil, nil)

	testRoundtripJSON(t, ls, new(LabeledStatement))

	if errs := ls.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ls.Label = Identifier{}
	if !hasError(ErrMissingNode, ls.Errors()...) {
		t.Error("expected ErrMissingNode for zero Label")
	}
	ls.Label.Name = "foo"
	ls.Body = nil
	if !hasError(ErrMissingNode, ls.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
}

func TestBreakStatement(t *testing.T) {
	var bs BreakStatement
	if bs.IsZero() {
		t.Error("expected !IsZero()")
	}

	bs.Label = Identifier{Name: "foo"}
	if bs.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", bs.MinVersion())
	}

	var v mockVisitor
	bs.Walk(&v)
	v.expect(t, bs, bs.Label, nil, nil)

	testRoundtripJSON(t, bs, new(BreakStatement))

	if errs := bs.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestContinueStatement(t *testing.T) {
	var cs ContinueStatement
	if cs.IsZero() {
		t.Error("expected !IsZero()")
	}

	cs.Label = Identifier{Name: "foo"}
	if cs.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", cs.MinVersion())
	}

	var v mockVisitor
	cs.Walk(&v)
	v.expect(t, cs, cs.Label, nil, nil)

	testRoundtripJSON(t, cs, new(ContinueStatement))

	if errs := cs.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}
