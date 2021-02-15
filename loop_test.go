package estree

import (
	"testing"
)

func TestWhileStatement(t *testing.T) {
	var ws WhileStatement
	if !ws.IsZero() {
		t.Error("expected IsZero()")
	}

	ws.Test = BoolLiteral{Value: true}
	ws.Body = ExpressionStatement{
		Expression: CallExpression{
			Callee: Identifier{Name: "foo"},
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
		ws.Test, nil,
		ws.Body,
		ws.Body.(ExpressionStatement).Expression,
		ws.Body.(ExpressionStatement).Expression.(CallExpression).Callee, nil, nil, nil, nil)

	testRoundtripJSON(t, ws, new(WhileStatement))

	if errs := ws.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ws.Test = CallExpression{}
	if !hasError(ErrMissingNode, ws.Errors()...) {
		t.Error("expected ErrMissingNode for zero Test")
	}
	ws.Test = nil
	if !hasError(ErrMissingNode, ws.Errors()...) {
		t.Error("expected ErrMissingNode for nil Test")
	}
	ws.Test = BoolLiteral{}
	ws.Body = ExpressionStatement{}
	if !hasError(ErrMissingNode, ws.Errors()...) {
		t.Error("expected ErrMissingNode for zero Body")
	}
	ws.Body = nil
	if !hasError(ErrMissingNode, ws.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
}

func TestDoWhileStatement(t *testing.T) {
	var dws DoWhileStatement
	if !dws.IsZero() {
		t.Error("expected IsZero()")
	}

	dws.Test = UpdateExpression{
		Operator: Increment,
		Argument: Identifier{Name: "foo"},
	}
	dws.Body = EmptyStatement{}
	if dws.IsZero() {
		t.Error("expected !IsZero()")
	}
	if dws.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", dws.MinVersion())
	}

	var v mockVisitor
	dws.Walk(&v)
	v.expect(t, dws,
		dws.Body, nil,
		dws.Test,
		dws.Test.(UpdateExpression).Argument, nil, nil, nil)

	testRoundtripJSON(t, dws, new(DoWhileStatement))

	if errs := dws.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	dws.Test = UpdateExpression{}
	if !hasError(ErrMissingNode, dws.Errors()...) {
		t.Error("expected ErrMissingNode for zero Test")
	}
	dws.Test = nil
	if !hasError(ErrMissingNode, dws.Errors()...) {
		t.Error("expected ErrMissingNode for nil Test")
	}
	dws.Test = BoolLiteral{}
	dws.Body = ExpressionStatement{}
	if !hasError(ErrMissingNode, dws.Errors()...) {
		t.Error("expected ErrMissingNode for zero Body")
	}
	dws.Body = nil
	if !hasError(ErrMissingNode, dws.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
}

func TestForStatement(t *testing.T) {
	var fs ForStatement
	if !fs.IsZero() {
		t.Error("expected IsZero()")
	}

	fs.Init = AssignmentExpression{
		Left:     Identifier{Name: "i"},
		Operator: Assign,
		Right:    NumberLiteral{},
	}
	fs.Test = BinaryExpression{
		Left:     Identifier{Name: "i"},
		Operator: LessThan,
		Right:    NumberLiteral{Value: 10},
	}
	fs.Update = AssignmentExpression{
		Left:     Identifier{Name: "i"},
		Operator: AddAssign,
		Right:    NumberLiteral{Value: 1},
	}
	fs.Body = DebuggerStatement{}
	if fs.IsZero() {
		t.Error("expected !IsZero()")
	}
	if fs.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", fs.MinVersion())
	}

	var v mockVisitor
	fs.Walk(&v)
	v.expect(t, fs,
		fs.Init,
		fs.Init.(AssignmentExpression).Left, nil,
		fs.Init.(AssignmentExpression).Right, nil, nil,
		fs.Test,
		fs.Test.(BinaryExpression).Left, nil,
		fs.Test.(BinaryExpression).Right, nil, nil,
		fs.Update,
		fs.Update.(AssignmentExpression).Left, nil,
		fs.Update.(AssignmentExpression).Right, nil, nil,
		fs.Body, nil, nil)

	testRoundtripJSON(t, fs, new(ForStatement))

	if errs := fs.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	fs.Init = nil
	fs.Test = nil
	fs.Update = nil
	if errs := fs.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors with nil Init/Test/Update: %v", errs)
	}
	fs.Body = ExpressionStatement{}
	if !hasError(ErrMissingNode, fs.Errors()...) {
		t.Error("expected ErrMissingNode for zero Body")
	}
	fs.Body = nil
	if !hasError(ErrMissingNode, fs.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
}

func TestForInStatement(t *testing.T) {
	var fis ForInStatement
	if !fis.IsZero() {
		t.Error("expected IsZero()")
	}

	fis.Left = Identifier{Name: "foo"}
	ae := ArrayExpression{}
	ae.Elements = append(ae.Elements, BoolLiteral{Value: true})
	ae.Elements = append(ae.Elements, BoolLiteral{Value: false})
	fis.Right = ae
	fis.Body = EmptyStatement{}
	if fis.IsZero() {
		t.Error("expected !IsZero()")
	}
	if fis.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", fis.MinVersion())
	}

	var v mockVisitor
	fis.Walk(&v)
	v.expect(t, fis,
		fis.Left, nil,
		fis.Right,
		fis.Right.(ArrayExpression).Elements[0], nil,
		fis.Right.(ArrayExpression).Elements[1], nil, nil,
		fis.Body, nil, nil)

	testRoundtripJSON(t, fis, new(ForInStatement))

	if errs := fis.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	fis.Left = Identifier{}
	if !hasError(ErrMissingNode, fis.Errors()...) {
		t.Error("expected ErrMissingNode for zero Left")
	}
	fis.Left = nil
	if !hasError(ErrMissingNode, fis.Errors()...) {
		t.Error("expected ErrMissingNode for nil Left")
	}
	fis.Left = Identifier{Name: "foo"}
	fis.Right = CallExpression{}
	if !hasError(ErrMissingNode, fis.Errors()...) {
		t.Error("expected ErrMissingNode for zero Right")
	}
	fis.Right = nil
	if !hasError(ErrMissingNode, fis.Errors()...) {
		t.Error("expected ErrMissingNode for nil Right")
	}
	fis.Right = Identifier{Name: "bar"}
	fis.Body = ExpressionStatement{}
	if !hasError(ErrMissingNode, fis.Errors()...) {
		t.Error("expected ErrMissingNode for zero Body")
	}
	fis.Body = nil
	if !hasError(ErrMissingNode, fis.Errors()...) {
		t.Error("expected ErrMissingNode for nil Body")
	}
}
