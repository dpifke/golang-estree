package estree

import (
	"testing"
)

func TestThisExpression(t *testing.T) {
	var te ThisExpression
	if te.IsZero() {
		t.Error("expected !IsZero()")
	}
	if te.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", te.MinVersion())
	}

	var v mockVisitor
	te.Walk(&v)
	v.expect(t, te, nil)

	testRoundtripJSON(t, te, new(ThisExpression))

	if errs := te.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestArrayExpression(t *testing.T) {
	var ae ArrayExpression
	if ae.IsZero() {
		t.Error("expected !IsZero()")
	}

	ae.Elements = append(ae.Elements, ThisExpression{})
	ae.Elements = append(ae.Elements, ArrayHole{})
	if ae.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ae.MinVersion())
	}

	var v mockVisitor
	ae.Walk(&v)
	v.expect(t, ae,
		ae.Elements[0], nil,
		ae.Elements[1], nil, nil)

	testRoundtripJSON(t, ae, new(ArrayExpression))

	if errs := ae.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ae.Elements[0] = LogicalExpression{}
	if !hasError(ErrMissingNode, ae.Errors()...) {
		t.Error("expected ErrMissingNode for zero Element")
	}
	ae.Elements[0] = nil
	if !hasError(ErrMissingNode, ae.Errors()...) {
		t.Error("expected ErrMissingNode for nil Element")
	}
}

func TestObjectExpression(t *testing.T) {
	var oe ObjectExpression
	if oe.IsZero() {
		t.Error("expected !IsZero()")
	}

	oe.Properties = []Property{
		Property{
			Key:   Identifier{Name: "foo"},
			Value: StringLiteral{Value: "bar"},
			Kind:  Init,
		},
		Property{
			Key:   StringLiteral{Value: "self"},
			Value: ThisExpression{},
			Kind:  Init,
		},
	}
	if oe.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", oe.MinVersion())
	}

	var v mockVisitor
	oe.Walk(&v)
	v.expect(t, oe,
		oe.Properties[0],
		oe.Properties[0].Key, nil,
		oe.Properties[0].Value, nil, nil,
		oe.Properties[1],
		oe.Properties[1].Key, nil,
		oe.Properties[1].Value, nil, nil, nil)

	testRoundtripJSON(t, oe, new(ObjectExpression))

	if errs := oe.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	oe.Properties[0] = Property{}
	if !hasError(ErrMissingNode, oe.Errors()...) {
		t.Error("expected ErrMissingNode for zero Property")
	}
}

func TestProperty(t *testing.T) {
	var p Property
	if !p.IsZero() {
		t.Error("expected IsZero()")
	}

	p.Key = Identifier{Name: "foo"}
	p.Value = StringLiteral{Value: "bar"}
	p.Kind = Init
	if p.IsZero() {
		t.Error("expected !IsZero()")
	}
	if p.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", p.MinVersion())
	}

	var v mockVisitor
	p.Walk(&v)
	v.expect(t, p,
		p.Key, nil,
		p.Value, nil, nil)

	testRoundtripJSON(t, p, new(Property))

	if errs := p.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	p.Key = Identifier{}
	if !hasError(ErrMissingNode, p.Errors()...) {
		t.Error("expected ErrMissingNode for zero Key")
	}
	p.Key = nil
	if !hasError(ErrMissingNode, p.Errors()...) {
		t.Error("expected ErrMissingNode for nil Key")
	}
	p.Key = Identifier{Name: "foo"}
	p.Value = Identifier{}
	if !hasError(ErrMissingNode, p.Errors()...) {
		t.Error("expected ErrMissingNode for zero Value")
	}
	p.Value = nil
	if !hasError(ErrMissingNode, p.Errors()...) {
		t.Error("expected ErrMissingNode for nil Value")
	}
	p.Value = StringLiteral{Value: "bar"}
	p.Kind = "not valid"
	if !hasError(ErrWrongValue, p.Errors()...) {
		t.Error("expected ErrWrongValue for invalid PropertyKind")
	}
}

func TestFunctionExpression(t *testing.T) {
	var fe FunctionExpression
	if !fe.IsZero() {
		t.Error("expected IsZero()")
	}

	fe.ID = Identifier{Name: "foo"}
	fe.Params = []Pattern{
		Identifier{Name: "bar"},
	}
	fe.Body.Body = append(fe.Body.Body, ReturnStatement{})
	if fe.IsZero() {
		t.Error("expected !IsZero()")
	}
	if fe.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", fe.MinVersion())
	}

	var v mockVisitor
	fe.Walk(&v)
	v.expect(t, fe,
		fe.Params[0], nil,
		fe.Body, fe.Body.Body[0], nil, nil, nil)

	testRoundtripJSON(t, fe, new(FunctionExpression))

	if errs := fe.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	fe.ID = Identifier{}
	if errs := fe.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors with zero ID: %v", errs)
	}
	fe.Params[0] = Identifier{}
	if !hasError(ErrMissingNode, fe.Errors()...) {
		t.Error("expected ErrMissingNode for zero Param")
	}
	fe.Params[0] = nil
	if !hasError(ErrMissingNode, fe.Errors()...) {
		t.Error("expected ErrMissingNode for nil Param")
	}
}

func TestConditionalExpression(t *testing.T) {
	var ce ConditionalExpression
	if !ce.IsZero() {
		t.Error("expected IsZero()")
	}

	ce.Test = BinaryExpression{
		Left:     Identifier{Name: "foo"},
		Operator: GreaterThan,
		Right:    NumberLiteral{Value: 1234},
	}
	ce.Consequent = StringLiteral{Value: "yes"}
	ce.Alternate = StringLiteral{Value: "no"}
	if ce.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ce.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ce.MinVersion())
	}

	var v mockVisitor
	ce.Walk(&v)
	v.expect(t, ce,
		ce.Test,
		ce.Test.(BinaryExpression).Left, nil,
		ce.Test.(BinaryExpression).Right, nil, nil,
		ce.Consequent, nil,
		ce.Alternate, nil, nil)

	testRoundtripJSON(t, ce, new(ConditionalExpression))

	if errs := ce.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ce.Test = Identifier{}
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for zero Test")
	}
	ce.Test = nil
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for nil Test")
	}
	ce.Test = BoolLiteral{Value: true}
	ce.Consequent = Identifier{}
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for zero Consequent")
	}
	ce.Consequent = nil
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for nil Consequent")
	}
	ce.Consequent = Identifier{Name: "foo"}
	ce.Alternate = Identifier{}
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for zero Alternate")
	}
	ce.Alternate = nil
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for nil Alternate")
	}
}

func TestCallExpression(t *testing.T) {
	var ce CallExpression
	if !ce.IsZero() {
		t.Error("expected IsZero()")
	}

	ce.Callee = MemberExpression{
		Object:   ThisExpression{},
		Property: Identifier{Name: "foo"},
	}
	ce.Arguments = append(ce.Arguments, NumberLiteral{Value: 1})
	ce.Arguments = append(ce.Arguments, NumberLiteral{Value: 2})
	ce.Arguments = append(ce.Arguments, NumberLiteral{Value: 3})
	if ce.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ce.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ce.MinVersion())
	}

	var v mockVisitor
	ce.Walk(&v)
	v.expect(t, ce,
		ce.Callee,
		ce.Callee.(MemberExpression).Object, nil,
		ce.Callee.(MemberExpression).Property, nil, nil,
		ce.Arguments[0], nil,
		ce.Arguments[1], nil,
		ce.Arguments[2], nil, nil)

	testRoundtripJSON(t, ce, new(CallExpression))

	if errs := ce.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ce.Callee = Identifier{}
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for zero Callee")
	}
	ce.Callee = nil
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for nil Callee")
	}
	ce.Callee = Identifier{Name: "foo"}
	ce.Arguments[1] = Identifier{}
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for zero Argument")
	}
	ce.Arguments[1] = nil
	if !hasError(ErrMissingNode, ce.Errors()...) {
		t.Error("expected ErrMissingNode for nil Argument")
	}
}

func TestNewExpression(t *testing.T) {
	var ne NewExpression
	if !ne.IsZero() {
		t.Error("expected IsZero()")
	}

	ne.Callee = Identifier{Name: "foo"}
	ne.Arguments = append(ne.Arguments, ThisExpression{})
	ne.Arguments = append(ne.Arguments, Identifier{Name: "bar"})
	if ne.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ne.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ne.MinVersion())
	}

	var v mockVisitor
	ne.Walk(&v)
	v.expect(t, ne,
		ne.Callee, nil,
		ne.Arguments[0], nil,
		ne.Arguments[1], nil, nil)

	testRoundtripJSON(t, ne, new(NewExpression))

	if errs := ne.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ne.Callee = Identifier{}
	if !hasError(ErrMissingNode, ne.Errors()...) {
		t.Error("expected ErrMissingNode for zero Callee")
	}
	ne.Callee = nil
	if !hasError(ErrMissingNode, ne.Errors()...) {
		t.Error("expected ErrMissingNode for nil Callee")
	}
	ne.Callee = Identifier{Name: "foo"}
	ne.Arguments[1] = Identifier{}
	if !hasError(ErrMissingNode, ne.Errors()...) {
		t.Error("expected ErrMissingNode for zero Argument")
	}
	ne.Arguments[1] = nil
	if !hasError(ErrMissingNode, ne.Errors()...) {
		t.Error("expected ErrMissingNode for nil Argument")
	}
}

func TestSequenceExpression(t *testing.T) {
	var se SequenceExpression
	if !se.IsZero() {
		t.Error("expected IsZero()")
	}

	se.Expressions = append(se.Expressions, Identifier{Name: "foo"})
	se.Expressions = append(se.Expressions, Identifier{Name: "bar"})
	if se.IsZero() {
		t.Error("expected !IsZero()")
	}
	if se.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", se.MinVersion())
	}

	var v mockVisitor
	se.Walk(&v)
	v.expect(t, se,
		se.Expressions[0], nil,
		se.Expressions[1], nil, nil)

	testRoundtripJSON(t, se, new(SequenceExpression))

	if errs := se.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	se.Expressions[1] = Identifier{}
	if !hasError(ErrMissingNode, se.Errors()...) {
		t.Error("expected ErrMissingNode for zero Expression")
	}
	se.Expressions[1] = nil
	if !hasError(ErrMissingNode, se.Errors()...) {
		t.Error("expected ErrMissingNode for nil Expression")
	}
}
