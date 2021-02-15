package estree

import (
	"testing"
)

func TestDirective(t *testing.T) {
	var d Directive
	if !d.IsZero() {
		t.Error("expected IsZero()")
	}

	d.Expression = StringLiteral{Value: "use strict"}
	d.Directive = "use strict"
	if d.IsZero() {
		t.Error("expected !IsZero()")
	}
	if d.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", d.MinVersion())
	}

	var v mockVisitor
	d.Walk(&v)
	v.expect(t, d, d.Expression, nil, nil)

	testRoundtripJSON(t, d, new(Directive))

	if errs := d.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
	d.Expression = nil
	if !hasError(ErrMissingNode, d.Errors()...) {
		t.Error("expected ErrMissingNode for nil Expression")
	}
}
