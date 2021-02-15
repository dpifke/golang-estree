package estree

import (
	"testing"
)

func TestFunctionDeclaration(t *testing.T) {
	var fd FunctionDeclaration
	if !fd.IsZero() {
		t.Error("expected IsZero()")
	}

	fd.ID = Identifier{Name: "foo"}
	fd.Params = []Pattern{
		Identifier{Name: "bar"},
	}
	fd.Body = FunctionBody{
		Body: []DirectiveOrStatement{
			Directive{
				Expression: StringLiteral{Value: "use strict"},
				Directive:  "use strict",
			},
		},
	}
	if fd.IsZero() {
		t.Error("expected !IsZero()")
	}
	if fd.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", fd.MinVersion())
	}

	var v mockVisitor
	fd.Walk(&v)
	v.expect(t, fd,
		fd.ID, nil,
		fd.Params[0], nil,
		fd.Body, fd.Body.Body[0],
		fd.Body.Body[0].(Directive).Expression, nil, nil, nil, nil)

	testRoundtripJSON(t, fd, new(FunctionDeclaration))

	if errs := fd.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	fd.ID = Identifier{}
	if !hasError(ErrMissingNode, fd.Errors()...) {
		t.Error("expected ErrMissingNode for zero ID")
	}
}

func TestVariableDeclaration(t *testing.T) {
	var vd VariableDeclaration
	if !vd.IsZero() {
		t.Error("expected IsZero()")
	}

	vd.Kind = Var
	vd.Declarations = []VariableDeclarator{
		VariableDeclarator{
			ID:   Identifier{Name: "blah"},
			Init: NullLiteral{},
		},
	}
	if vd.IsZero() {
		t.Error("expected !IsZero()")
	}
	if vd.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", vd.MinVersion())
	}

	var v mockVisitor
	vd.Walk(&v)
	v.expect(t, vd,
		vd.Declarations[0],
		vd.Declarations[0].ID, nil,
		vd.Declarations[0].Init, nil, nil, nil)

	testRoundtripJSON(t, vd, new(VariableDeclaration))

	if errs := vd.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	vd.Kind = "not valid"
	if !hasError(ErrWrongValue, vd.Errors()...) {
		t.Error("expected ErrWrongValue for invalid Kind")
	}
	vd.Kind = Var
	vd.Declarations[0] = VariableDeclarator{}
	if !hasError(ErrMissingNode, vd.Errors()...) {
		t.Error("expected ErrMissingNode for zero Declaration")
	}
}

func TestVariableDeclarator(t *testing.T) {
	var vd VariableDeclarator
	if !vd.IsZero() {
		t.Error("expected IsZero()")
	}

	vd.ID = Identifier{Name: "blah"}
	vd.Init = UnaryExpression{
		Operator: BitwiseNot,
		Prefix:   true,
		Argument: NumberLiteral{Value: 1234},
	}
	if vd.IsZero() {
		t.Error("expected !IsZero()")
	}
	if vd.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", vd.MinVersion())
	}

	var v mockVisitor
	vd.Walk(&v)
	v.expect(t, vd,
		vd.ID, nil,
		vd.Init, vd.Init.(UnaryExpression).Argument, nil, nil, nil)

	testRoundtripJSON(t, vd, new(VariableDeclarator))

	if errs := vd.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	vd.Init = nil
	if errs := vd.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	vd.ID = Identifier{}
	if !hasError(ErrMissingNode, vd.Errors()...) {
		t.Error("expected ErrMissingNode for zero ID")
	}
}
