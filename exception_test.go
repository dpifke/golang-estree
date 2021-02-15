package estree

import (
	"testing"
)

func TestThrowStatement(t *testing.T) {
	var ts ThrowStatement
	if !ts.IsZero() {
		t.Error("expected IsZero()")
	}

	ts.Argument = UpdateExpression{
		Operator: Increment,
		Argument: NumberLiteral{Value: 1234},
	}
	if ts.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ts.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ts.MinVersion())
	}

	var v mockVisitor
	ts.Walk(&v)
	v.expect(t, ts,
		ts.Argument,
		ts.Argument.(UpdateExpression).Argument, nil, nil, nil)

	testRoundtripJSON(t, ts, new(ThrowStatement))

	if errs := ts.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ts.Argument = LogicalExpression{}
	if !hasError(ErrMissingNode, ts.Errors()...) {
		t.Error("expected ErrMissingNode for zero Argument")
	}
	ts.Argument = nil
	if !hasError(ErrMissingNode, ts.Errors()...) {
		t.Error("expected ErrMissingNode for nil Argument")
	}
}

func TestTryStatement(t *testing.T) {
	var ts TryStatement
	if !ts.IsZero() {
		t.Error("expected IsZero()")
	}

	ts.Block = BlockStatement{
		Body: []Statement{
			ThrowStatement{
				Argument: Identifier{Name: "foo"},
			},
		},
	}
	ts.Handler = CatchClause{
		Param: Identifier{Name: "bar"},
		Body: BlockStatement{
			Body: []Statement{
				ReturnStatement{},
			},
		},
	}
	ts.Finalizer = BlockStatement{
		Body: []Statement{
			EmptyStatement{},
		},
	}
	if ts.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ts.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ts.MinVersion())
	}

	var v mockVisitor
	ts.Walk(&v)
	v.expect(t, ts,
		ts.Block,
		ts.Block.Body[0],
		ts.Block.Body[0].(ThrowStatement).Argument, nil, nil, nil,
		ts.Handler,
		ts.Handler.Param, nil,
		ts.Handler.Body,
		ts.Handler.Body.Body[0], nil, nil, nil,
		ts.Finalizer,
		ts.Finalizer.Body[0], nil, nil, nil)

	testRoundtripJSON(t, ts, new(TryStatement))

	if errs := ts.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ts.Handler = CatchClause{}
	ts.Finalizer = BlockStatement{}
	if !hasError(ErrMissingNode, ts.Errors()...) {
		t.Error("expected ErrMissingNode for zero Handler and Finalizer")
	}
}

func TestCatchClause(t *testing.T) {
	var cc CatchClause
	if !cc.IsZero() {
		t.Error("expected IsZero()")
	}

	cc.Param = Identifier{Name: "foo"}
	cc.Body.Body = []Statement{
		BreakStatement{},
	}
	if cc.IsZero() {
		t.Error("expected !IsZero()")
	}
	if cc.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", cc.MinVersion())
	}

	var v mockVisitor
	cc.Walk(&v)
	v.expect(t, cc,
		cc.Param, nil,
		cc.Body, cc.Body.Body[0], nil, nil, nil)

	testRoundtripJSON(t, cc, new(CatchClause))

	if errs := cc.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	cc.Param = Identifier{}
	if !hasError(ErrMissingNode, cc.Errors()...) {
		t.Error("expected ErrMissingNode for zero Param")
	}
	cc.Param = nil
	if !hasError(ErrMissingNode, cc.Errors()...) {
		t.Error("expected ErrMissingNode for nil Param")
	}
}
