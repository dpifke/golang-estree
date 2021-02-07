package estree

import (
	"encoding/json"
	"fmt"
)

type Expression interface {
	Node

	isExpression()
	isVariableDeclarationOrExpression()
	isPatternOrExpression()
}

func unmarshalExpression(m json.RawMessage) (e Expression, match bool, err error) {
	var x struct {
		Type string `json:"type"`
	}
	if err = json.Unmarshal(m, &x); err == nil {
		switch x.Type {
		case ThisExpression{}.Type():
			e, match = ThisExpression{}, true
		case ArrayExpression{}.Type():
			var ae ArrayExpression
			e, match, err = ae, true, json.Unmarshal(m, &ae)
		case ObjectExpression{}.Type():
			var oe ObjectExpression
			e, match, err = oe, true, json.Unmarshal(m, &oe)
		case FunctionExpression{}.Type():
			var fe FunctionExpression
			e, match, err = fe, true, json.Unmarshal(m, &fe)
		case UnaryExpression{}.Type():
			var ue UnaryExpression
			e, match, err = ue, true, json.Unmarshal(m, &ue)
		case UpdateExpression{}.Type():
			var ue UpdateExpression
			e, match, err = ue, true, json.Unmarshal(m, &ue)
		case BinaryExpression{}.Type():
			var be BinaryExpression
			e, match, err = be, true, json.Unmarshal(m, &be)
		case AssignmentExpression{}.Type():
			var ae AssignmentExpression
			e, match, err = ae, true, json.Unmarshal(m, &ae)
		case LogicalExpression{}.Type():
			var le LogicalExpression
			e, match, err = le, true, json.Unmarshal(m, &le)
		case MemberExpression{}.Type():
			var me MemberExpression
			e, match, err = me, true, json.Unmarshal(m, &me)
		case ConditionalExpression{}.Type():
			var ce ConditionalExpression
			e, match, err = ce, true, json.Unmarshal(m, &ce)
		case CallExpression{}.Type():
			var ce CallExpression
			e, match, err = ce, true, json.Unmarshal(m, &ce)
		case NewExpression{}.Type():
			var ne NewExpression
			e, match, err = ne, true, json.Unmarshal(m, &ne)
		case SequenceExpression{}.Type():
			var se SequenceExpression
			e, match, err = se, true, json.Unmarshal(m, &se)
		default:
			err = fmt.Errorf("%w: expected Expression, got %v", ErrWrongType, string(m))
		}
		if err != nil {
			e = nil // don't return incomplete nodes
		}
	}
	return
}

type baseExpression struct{}

func (baseExpression) isExpression()                      {}
func (baseExpression) isVariableDeclarationOrExpression() {}
func (baseExpression) isPatternOrExpression()             {}

type ThisExpression struct {
	baseExpression
}

func (ThisExpression) Type() string { return "ThisExpression" }

func (te ThisExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": te.Type(),
	})
}

func (te *ThisExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string `json:"type"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != te.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, te.Type(), x.Type)
	}
	return err
}

type ArrayExpression struct {
	baseExpression
	Elements []Expression
}

func (ArrayExpression) Type() string { return "ArrayExpression" }

func (ae ArrayExpression) MarshalJSON() ([]byte, error) {
	x := map[string]interface{}{
		"type": ae.Type(),
	}
	if len(ae.Elements) > 0 {
		x["elements"] = ae.Elements
	}
	return json.Marshal(x)
}

func (ae *ArrayExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string            `json:"type"`
		Elements []json.RawMessage `json:"elements"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ae.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ae.Type(), x.Type)
	}
	if err == nil {
		ae.Elements = make([]Expression, len(x.Elements))
		for i := range x.Elements {
			if ae.Elements[i], _, err = unmarshalExpression(x.Elements[i]); err != nil {
				break
			}
		}
	}
	return err
}

type ObjectExpression struct {
	baseExpression
	Properties []Property
}

func (ObjectExpression) Type() string { return "ObjectExpression" }

func (oe ObjectExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       oe.Type(),
		"properties": oe.Properties,
	})
}

func (oe *ObjectExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string     `json:"type"`
		Properties []Property `json:"properties"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != oe.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, oe.Type(), x.Type)
	}
	return err
}

type PropertyKind string

var (
	Init PropertyKind = "init"
	Get  PropertyKind = "get"
	Set  PropertyKind = "set"
)

func (pk PropertyKind) GoString() string {
	switch pk {
	case Init:
		return "Init"
	case Get:
		return "Get"
	case Set:
		return "Set"
	}
	return fmt.Sprintf("%q", pk)
}

type Property struct {
	Key   LiteralOrIdentifier
	Value Expression
	Kind  PropertyKind
}

func (Property) Type() string { return "Property" }

func (p Property) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  p.Type(),
		"key":   p.Key,
		"value": p.Value,
		"kind":  p.Kind,
	})
}

func (p *Property) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string          `json:"type"`
		Key   json.RawMessage `json:"key"`
		Value json.RawMessage `json:"value"`
		Kind  PropertyKind    `json:"kind"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != p.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, p.Type(), x.Type)
	}
	if err == nil {
		switch x.Kind {
		case Init, Get, Set:
			p.Kind = x.Kind
		default:
			err = fmt.Errorf("%w for Property.Kind: %q", ErrWrongValue, x.Kind)
		}
	}
	if err == nil {
		p.Key, _, err = unmarshalLiteralOrIdentifier(x.Key)
	}
	if err == nil {
		p.Value, _, err = unmarshalExpression(x.Value)
	}
	return err
}

type FunctionExpression struct {
	baseExpression
	ID     Identifier
	Params []Pattern
	Body   FunctionBody
}

func (FunctionExpression) Type() string { return "FunctionExpression" }

func (fe FunctionExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":   fe.Type(),
		"id":     fe.ID,
		"params": fe.Params,
		"body":   fe.Body,
	})
}

func (fe *FunctionExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type   string            `json:"type"`
		ID     Identifier        `json:"id"`
		Params []json.RawMessage `json:"params"`
		Body   FunctionBody      `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != fe.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, fe.Type(), x.Type)
	}
	if err == nil && len(x.Params) > 0 {
		fe.Params = make([]Pattern, len(x.Params))
		for i := range x.Params {
			if fe.Params[i], _, err = unmarshalPattern(x.Params[i]); err != nil {
				break
			}
		}
	}
	if err == nil {
		fe.ID, fe.Body = x.ID, x.Body
	}
	return err
}

type ConditionalExpression struct {
	baseExpression
	Test, Alternate, Consequent Expression
}

func (ConditionalExpression) Type() string { return "ConditionalExpression" }

func (ce ConditionalExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       ce.Type(),
		"test":       ce.Test,
		"alternate":  ce.Alternate,
		"consequent": ce.Consequent,
	})
}

func (ce *ConditionalExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Test       json.RawMessage `json:"test"`
		Alternate  json.RawMessage `json:"alternate"`
		Consequent json.RawMessage `json:"consequent"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ce.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ce.Type(), x.Type)
	}
	if err == nil {
		ce.Test, _, err = unmarshalExpression(x.Test)
	}
	if err == nil {
		ce.Alternate, _, err = unmarshalExpression(x.Alternate)
	}
	if err == nil {
		ce.Consequent, _, err = unmarshalExpression(x.Consequent)
	}
	return err
}

type CallExpression struct {
	baseExpression
	Callee    Expression
	Arguments []Expression
}

func (CallExpression) Type() string { return "CallExpression" }

func (ce CallExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":      ce.Type(),
		"callee":    ce.Callee,
		"arguments": ce.Arguments,
	})
}

func (ce *CallExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type      string            `json:"type"`
		Callee    json.RawMessage   `json:"callee"`
		Arguments []json.RawMessage `json:"arguments"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ce.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ce.Type(), x.Type)
	}
	if err == nil {
		ce.Callee, _, err = unmarshalExpression(x.Callee)
	}
	if err == nil {
		ce.Arguments = make([]Expression, len(x.Arguments))
		for i := range x.Arguments {
			if ce.Arguments[i], _, err = unmarshalExpression(x.Arguments[i]); err != nil {
				break
			}
		}
	}
	return err
}

type NewExpression struct {
	baseExpression
	Callee    Expression
	Arguments []Expression
}

func (NewExpression) Type() string { return "NewExpression" }

func (ne NewExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":      ne.Type(),
		"callee":    ne.Callee,
		"arguments": ne.Arguments,
	})
}

func (ne *NewExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type      string            `json:"type"`
		Callee    json.RawMessage   `json:"callee"`
		Arguments []json.RawMessage `json:"arguments"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ne.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ne.Type(), x.Type)
	}
	if err == nil {
		ne.Callee, _, err = unmarshalExpression(x.Callee)
	}
	if err == nil {
		ne.Arguments = make([]Expression, len(x.Arguments))
		for i := range x.Arguments {
			if ne.Arguments[i], _, err = unmarshalExpression(x.Arguments[i]); err != nil {
				break
			}
		}
	}
	return err
}

type SequenceExpression struct {
	baseExpression
	Expressions []Expression
}

func (SequenceExpression) Type() string { return "SequenceExpression" }

func (se SequenceExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":        se.Type(),
		"expressions": se.Expressions,
	})
}

func (se *SequenceExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type        string            `json:"type"`
		Expressions []json.RawMessage `json:"expressions"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != se.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, se.Type(), x.Type)
	}
	if err == nil {
		se.Expressions = make([]Expression, len(x.Expressions))
		for i := range x.Expressions {
			if se.Expressions[i], _, err = unmarshalExpression(x.Expressions[i]); err != nil {
				break
			}
		}
	}
	return err
}
