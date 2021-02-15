package estree

import (
	"encoding/json"
	"errors"
	"reflect"
	"testing"
)

func testRoundtripLiteralJSON(t *testing.T, in Literal) {
	t.Helper()
	b, err := json.Marshal(in)
	if err != nil {
		t.Error(err)
	} else if out, match, err := unmarshalLiteral(b); !match || err != nil {
		t.Errorf("match is %t, err is %v", match, err)
	} else if !reflect.DeepEqual(in, out) {
		t.Errorf("JSON roundtrip failed marshaling/unmarshaling %T", in)
		t.Log("marshal:", string(b))
		t.Logf("unmarshal: %+v", out)
		if b2, err := json.Marshal(out); err == nil {
			t.Log("re-marshal:", string(b2))
		}
	}
}

func TestUnmarshalInvalidLiteral(t *testing.T) {
	b := []byte(`{"type":"MagicLiteral","value":"hocus pocus"}`)
	_, match, err := unmarshalLiteral(b)
	if match {
		t.Errorf("expected !match, err is %v", err)
	}
	if !errors.Is(err, ErrWrongType) {
		t.Errorf("expected ErrWrongType, got %v", err)
	}

	b = []byte(`{"type":"Literal","value":[1,2,3]}`)
	_, match, err = unmarshalLiteral(b)
	if !match {
		t.Errorf("expected match, err is %v", err)
	}
	if !errors.Is(err, ErrWrongType) {
		t.Errorf("expected ErrWrongType, got %v", err)
	}
}

func TestStringLiteral(t *testing.T) {
	var sl StringLiteral
	if sl.IsZero() {
		t.Error("expected !IsZero()")
	}

	sl.Value = "foo"
	if sl.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", sl.MinVersion())
	}

	var v mockVisitor
	sl.Walk(&v)
	v.expect(t, sl, nil)

	testRoundtripLiteralJSON(t, sl)

	if errs := sl.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestBoolLiteral(t *testing.T) {
	var bl BoolLiteral
	if bl.IsZero() {
		t.Error("expected !IsZero()")
	}

	bl.Value = true
	if bl.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", bl.MinVersion())
	}

	var v mockVisitor
	bl.Walk(&v)
	v.expect(t, bl, nil)

	testRoundtripLiteralJSON(t, bl)

	if errs := bl.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestNullLiteral(t *testing.T) {
	var nl NullLiteral
	if nl.IsZero() {
		t.Error("expected !IsZero()")
	}

	if nl.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", nl.MinVersion())
	}

	var v mockVisitor
	nl.Walk(&v)
	v.expect(t, nl, nil)

	testRoundtripLiteralJSON(t, nl)

	if errs := nl.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestNumberLiteral(t *testing.T) {
	var nl NumberLiteral
	if nl.IsZero() {
		t.Error("expected !IsZero()")
	}

	nl.Value = 1234.4321
	if nl.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", nl.MinVersion())
	}

	var v mockVisitor
	nl.Walk(&v)
	v.expect(t, nl, nil)

	testRoundtripLiteralJSON(t, nl)

	if errs := nl.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}

func TestRegExpLiteral(t *testing.T) {
	var rl RegExpLiteral
	if rl.IsZero() {
		t.Error("expected !IsZero()")
	}

	rl.Pattern = "^.*$"
	rl.Flags = "i"
	if rl.MinVersion() != ES5 {
		t.Errorf("expected ES5, got %s", rl.MinVersion())
	}

	var v mockVisitor
	rl.Walk(&v)
	v.expect(t, rl, nil)

	testRoundtripLiteralJSON(t, rl)

	if errs := rl.Errors(); len(errs) != 0 {
		t.Fatalf("unexpected errors: %v", errs)
	}
}
