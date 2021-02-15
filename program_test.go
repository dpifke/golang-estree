package estree

import (
	"testing"
)

func TestProgram(t *testing.T) {
	var p Program
	if !p.IsZero() {
		t.Error("expected IsZero()")
	}

	p.Body = append(p.Body, Directive{
		Expression: StringLiteral{Value: "use strict"},
		Directive:  "use strict",
	})
	p.Body = append(p.Body, ThrowStatement{
		Argument: NewExpression{
			Callee: Identifier{Name: "Exception"},
		},
	})
	if p.IsZero() {
		t.Error("expected !IsZero()")
	}
	if p.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", p.MinVersion())
	}

	var v mockVisitor
	p.Walk(&v)
	v.expect(t, p,
		p.Body[0],
		p.Body[0].(Directive).Expression, nil, nil,
		p.Body[1],
		p.Body[1].(ThrowStatement).Argument,
		p.Body[1].(ThrowStatement).Argument.(NewExpression).Callee, nil, nil, nil, nil)

	testRoundtripJSON(t, p, new(Program))

	if errs := p.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	p.Body[0] = Directive{}
	if !hasError(ErrMissingNode, p.Errors()...) {
		t.Error("expected ErrMissingNode for zero Body")
	}
	p.Body[0] = nil
	if !hasError(ErrMissingNode, p.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
}
