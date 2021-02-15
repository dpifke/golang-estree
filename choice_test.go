package estree

import (
	"testing"
)

func TestIfStatement(t *testing.T) {
	var is IfStatement
	if !is.IsZero() {
		t.Error("expected IsZero()")
	}

	is.Test = CallExpression{
		Callee: Identifier{Name: "foo"},
	}
	is.Consequent = ReturnStatement{}
	is.Alternate = DebuggerStatement{}
	if is.IsZero() {
		t.Error("expected !IsZero()")
	}
	if is.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", is.MinVersion())
	}

	var v mockVisitor
	is.Walk(&v)
	v.expect(t, is,
		is.Test, is.Test.(CallExpression).Callee, nil, nil,
		is.Consequent, nil,
		is.Alternate, nil, nil)

	testRoundtripJSON(t, is, new(IfStatement))

	if errs := is.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	is.Test = FunctionExpression{}
	if !hasError(ErrMissingNode, is.Errors()...) {
		t.Error("expected ErrMissingNode for zero Test")
	}
	is.Test = ThisExpression{}
	is.Consequent = nil
	if !hasError(ErrMissingNode, is.Errors()...) {
		t.Error("expected ErrMissingNode for nil Consequent")
	}
}

func TestSwitchStatement(t *testing.T) {
	var ss SwitchStatement
	if !ss.IsZero() {
		t.Error("expected IsZero()")
	}

	ss.Discriminant = Identifier{Name: "foo"}
	ss.Cases = []SwitchCase{
		SwitchCase{
			Test: StringLiteral{},
		},
		SwitchCase{
			Test: ThisExpression{},
			Consequent: []Statement{
				ContinueStatement{
					Label: Identifier{Name: "bar"},
				},
			},
		},
		SwitchCase{
			Consequent: []Statement{
				ReturnStatement{},
			},
		},
	}
	if ss.IsZero() {
		t.Error("expected !IsZero()")
	}
	if ss.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", ss.MinVersion())
	}

	var v mockVisitor
	ss.Walk(&v)
	v.expect(t, ss,
		ss.Discriminant, nil,
		ss.Cases[0], ss.Cases[0].Test, nil, nil,
		ss.Cases[1], ss.Cases[1].Test, nil,
		ss.Cases[1].Consequent[0],
		ss.Cases[1].Consequent[0].(ContinueStatement).Label, nil, nil, nil,
		ss.Cases[2], ss.Cases[2].Consequent[0], nil, nil, nil)

	testRoundtripJSON(t, ss, new(SwitchStatement))

	if errs := ss.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	ss.Discriminant = SequenceExpression{}
	if !hasError(ErrMissingNode, ss.Errors()...) {
		t.Error("expected ErrMissingNode for zero Test")
	}
	ss.Discriminant = nil
	if !hasError(ErrMissingNode, ss.Errors()...) {
		t.Error("expected ErrMissingNode for nil Test")
	}
}

func TestSwitchCase(t *testing.T) {
	var sc SwitchCase
	if sc.IsZero() {
		t.Error("expected !IsZero()")
	}

	sc.Test = StringLiteral{}
	sc.Consequent = []Statement{
		BreakStatement{},
	}
	if sc.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", sc.MinVersion())
	}

	var v mockVisitor
	sc.Walk(&v)
	v.expect(t, sc,
		sc.Test, nil,
		sc.Consequent[0], nil, nil)

	testRoundtripJSON(t, sc, new(SwitchCase))

	if errs := sc.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	sc.Consequent[0] = IfStatement{}
	if !hasError(ErrMissingNode, sc.Errors()...) {
		t.Error("expected ErrMissingNode for zero Consequent")
	}
	sc.Consequent[0] = nil
	if !hasError(ErrMissingNode, sc.Errors()...) {
		t.Error("expected ErrMissingNode for nil Consequent")
	}
}
