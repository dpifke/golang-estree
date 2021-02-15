package estree

import (
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"strings"
	"testing"
)

func printNodes(nodes []Node) string {
	var s strings.Builder
	s.WriteRune('[')
	for i, n := range nodes {
		if i > 0 {
			s.WriteString(", ")
		}
		fmt.Fprintf(&s, "%T", n)
	}
	s.WriteRune(']')
	return s.String()
}

// mockVisitor implements Visitor by storing the list of nodes visited.
type mockVisitor []Node

func (v *mockVisitor) Visit(n Node) Visitor {
	*v = append(*v, n)
	return v
}

func (v mockVisitor) String() string {
	return "visited: " + printNodes(v)
}

// expect logs a test error if the visited nodes don't match.
func (v mockVisitor) expect(t *testing.T, nodes ...Node) {
	if !reflect.DeepEqual(nodes, []Node(v)) {
		t.Helper()
		t.Error("expected:", printNodes(nodes))
		t.Error(v)
	}
}

// testRoundtripJSON is a test helper for JSON serialization.
func testRoundtripJSON(t *testing.T, in Node, out json.Unmarshaler) {
	t.Helper()
	b, err := json.Marshal(in)
	if err != nil {
		t.Error(err)
	} else if err := out.UnmarshalJSON(b); err != nil {
		t.Error(err)
	} else if !reflect.DeepEqual(in, reflect.ValueOf(out).Elem().Interface()) {
		t.Errorf("JSON roundtrip failed marshaling/unmarshaling %T", in)
		t.Log("marshal:", string(b))
		t.Logf("unmarshal: %+v", out)
		if b2, err := json.Marshal(out); err == nil {
			t.Log("re-marshal:", string(b2))
		}
	}

	err = out.UnmarshalJSON([]byte(`{"type":"DoesNotExist","foo":"bar"}`))
	if !errors.Is(err, ErrWrongType) {
		t.Errorf("expected ErrWrongType unmarshaling %T, got %v", in, err)
	}
}
