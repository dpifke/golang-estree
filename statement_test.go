package estree

import (
	"testing"
)

func TestExpressionStatement(t *testing.T) {
	var es ExpressionStatement
	if !es.IsZero() {
		t.Error("expected IsZero()")
	}

	es.Expression = CallExpression{
		Callee: Identifier{Name: "foo"},
	}
	if es.IsZero() {
		t.Error("expected !IsZero()")
	}
	if es.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", es.MinVersion())
	}

	var v mockVisitor
	es.Walk(&v)
	v.expect(t, es,
		es.Expression,
		es.Expression.(CallExpression).Callee, nil, nil, nil)

	testRoundtripJSON(t, es, new(ExpressionStatement))

	if errs := es.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	es.Expression = Identifier{}
	if !hasError(ErrMissingNode, es.Errors()...) {
		t.Error("expected ErrMissingNode for zero Expression")
	}
	es.Expression = nil
	if !hasError(ErrMissingNode, es.Errors()...) {
		t.Error("expected ErrMissingNode for nil Expression")
	}
}

func TestBlockStatement(t *testing.T) {
	var bs BlockStatement
	if bs.IsZero() {
		t.Error("expected !IsZero()")
	}

	bs.Body = append(bs.Body, VariableDeclaration{
		Kind: Var,
		Declarations: []VariableDeclarator{
			VariableDeclarator{
				ID: Identifier{Name: "foo"},
			},
		},
	})
	bs.Body = append(bs.Body, ReturnStatement{})
	if bs.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", bs.MinVersion())
	}

	var v mockVisitor
	bs.Walk(&v)
	v.expect(t, bs,
		bs.Body[0],
		bs.Body[0].(VariableDeclaration).Declarations[0],
		bs.Body[0].(VariableDeclaration).Declarations[0].ID, nil, nil, nil,
		bs.Body[1], nil, nil)

	testRoundtripJSON(t, bs, new(BlockStatement))

	if errs := bs.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	bs.Body[1] = ExpressionStatement{}
	if !hasError(ErrMissingNode, bs.Errors()...) {
		t.Error("expected ErrMissingNode for zero Body")
	}
	bs.Body[1] = nil
	if !hasError(ErrMissingNode, bs.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
}

func TestFunctionBody(t *testing.T) {
	var fb FunctionBody
	if fb.IsZero() {
		t.Error("expected !IsZero()")
	}

	fb.Body = append(fb.Body, Directive{
		Expression: StringLiteral{Value: "use strict"},
		Directive:  "use strict",
	})
	fb.Body = append(fb.Body, ReturnStatement{})
	if fb.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", fb.MinVersion())
	}

	var v mockVisitor
	fb.Walk(&v)
	v.expect(t, fb,
		fb.Body[0],
		fb.Body[0].(Directive).Expression, nil, nil,
		fb.Body[1], nil, nil)

	testRoundtripJSON(t, fb, new(FunctionBody))

	if errs := fb.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	fb.Body[1] = ExpressionStatement{}
	if !hasError(ErrMissingNode, fb.Errors()...) {
		t.Error("expected ErrMissingNode for zero Body")
	}
	fb.Body[1] = nil
	if !hasError(ErrMissingNode, fb.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
}

func TestEmptyStatement(t *testing.T) {
	var es EmptyStatement
	if es.IsZero() {
		t.Error("expected !IsZero()")
	}
	if es.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", es.MinVersion())
	}

	var v mockVisitor
	es.Walk(&v)
	v.expect(t, es, nil)

	testRoundtripJSON(t, es, new(EmptyStatement))

	if errs := es.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestDebuggerStatement(t *testing.T) {
	var ds DebuggerStatement
	if ds.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ds.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ds.MinVersion())
	}

	var v mockVisitor
	ds.Walk(&v)
	v.expect(t, ds, nil)

	testRoundtripJSON(t, ds, new(DebuggerStatement))

	if errs := ds.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestWithStatement(t *testing.T) {
	var ws WithStatement
	if !ws.IsZero() {
		t.Error("expected IsZero()")
	}

	ws.Object = ObjectExpression{
		Properties: []Property{
			Property{
				Key:   Identifier{Name: "foo"},
				Kind:  Init,
				Value: Identifier{Name: "bar"},
			},
		},
	}
	ws.Body = BlockStatement{
		Body: []Statement{
			ExpressionStatement{
				Expression: CallExpression{
					Callee: Identifier{Name: "setTimeout"},
					Arguments: []Expression{
						Identifier{Name: "foo"},
					},
				},
			},
		},
	}
	if ws.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ws.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ws.MinVersion())
	}

	var v mockVisitor
	ws.Walk(&v)
	v.expect(t, ws,
		ws.Object,
		ws.Object.(ObjectExpression).Properties[0],
		ws.Object.(ObjectExpression).Properties[0].Key, nil,
		ws.Object.(ObjectExpression).Properties[0].Value, nil, nil, nil,
		ws.Body,
		ws.Body.(BlockStatement).Body[0],
		ws.Body.(BlockStatement).Body[0].(ExpressionStatement).Expression,
		ws.Body.(BlockStatement).Body[0].(ExpressionStatement).Expression.(CallExpression).Callee, nil,
		ws.Body.(BlockStatement).Body[0].(ExpressionStatement).Expression.(CallExpression).Arguments[0],
		nil, nil, nil, nil, nil)

	testRoundtripJSON(t, ws, new(WithStatement))

	if errs := ws.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ws.Object = Identifier{}
	if !hasError(ErrMissingNode, ws.Errors()...) {
		t.Error("expected ErrMissingNode for zero Object")
	}
	ws.Object = nil
	if !hasError(ErrMissingNode, ws.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
	ws.Object = Identifier{Name: "foo"}
	ws.Body = ExpressionStatement{}
	if !hasError(ErrMissingNode, ws.Errors()...) {
		t.Error("expected ErrMissingNode for zero Body")
	}
	ws.Body = nil
	if !hasError(ErrMissingNode, ws.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
}
