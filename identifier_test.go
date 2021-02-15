package estree

import (
	"testing"
)

func TestIdentifier(t *testing.T) {
	var i Identifier
	if !i.IsZero() {
		t.Error("expected IsZero() for uninitialized Identifier")
	}
	if !hasError(ErrWrongValue, i.Errors()...) {
		t.Error("expected ErrWrongValue for empty Identifier.Name")
	}

	i.Name = "foo"
	if i.IsZero() {
		t.Error("expected !IsZero() after setting Identifier.Name")
	}
	if errs := i.Errors(); len(errs) != 0 {
		t.Errorf("unexpected errors after setting Identifier.Name: %v", errs)
	}

	var v mockVisitor
	i.Walk(&v)
	v.expect(t, i, nil)

	testRoundtripJSON(t, i, new(Identifier))
}
