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

func unmarshalExpression(m json.RawMessage) (Expression, error) {
	var x struct {
		Type string `json:"type"`
	}
	err := json.Unmarshal(m, &x)
	if err == nil {
		switch x.Type {
		case ThisExpression{}.Type():
			return ThisExpression{}, nil
		case ArrayExpression{}.Type():
			var ae ArrayExpression
			if err = json.Unmarshal(m, &ae); err == nil {
				return ae, nil
			}
		case ObjectExpression{}.Type():
			var oe ObjectExpression
			if err = json.Unmarshal(m, &oe); err == nil {
				return oe, nil
			}
		case FunctionExpression{}.Type():
			var fe FunctionExpression
			if err = json.Unmarshal(m, &fe); err == nil {
				return fe, nil
			}
		case UnaryExpression{}.Type():
			var ue UnaryExpression
			if err = json.Unmarshal(m, &ue); err == nil {
				return ue, nil
			}
		case UpdateExpression{}.Type():
			var ue UpdateExpression
			if err = json.Unmarshal(m, &ue); err == nil {
				return ue, nil
			}
		case BinaryExpression{}.Type():
			var be BinaryExpression
			if err = json.Unmarshal(m, &be); err == nil {
				return be, nil
			}
		case AssignmentExpression{}.Type():
			var ae AssignmentExpression
			if err = json.Unmarshal(m, &ae); err == nil {
				return ae, nil
			}
		case LogicalExpression{}.Type():
			var le LogicalExpression
			if err = json.Unmarshal(m, &le); err == nil {
				return le, nil
			}
		case MemberExpression{}.Type():
			var me MemberExpression
			if err = json.Unmarshal(m, &me); err == nil {
				return me, nil
			}
		case ConditionalExpression{}.Type():
			var ce ConditionalExpression
			if err = json.Unmarshal(m, &ce); err == nil {
				return ce, nil
			}
		case CallExpression{}.Type():
			var ce CallExpression
			if err = json.Unmarshal(m, &ce); err == nil {
				return ce, nil
			}
		case NewExpression{}.Type():
			var ne NewExpression
			if err = json.Unmarshal(m, &ne); err == nil {
				return ne, nil
			}
		case SequenceExpression{}.Type():
			var se SequenceExpression
			if err = json.Unmarshal(m, &se); err == nil {
				return se, nil
			}
		default:
			err = fmt.Errorf("%w: expected Expression, got %v", ErrWrongType, string(m))
		}
	}
	return nil, err
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
			if ae.Elements[i], err = unmarshalExpression(x.Elements[i]); err != nil {
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
		"value": p.Value,
		"kind":  p.Kind,
	})
}

func (p *Property) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string          `json:"type"`
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
		p.Value, err = unmarshalExpression(x.Value)
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
			if fe.Params[i], err = unmarshalPattern(x.Params[i]); err != nil {
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
		ce.Test, err = unmarshalExpression(x.Test)
	}
	if err == nil {
		ce.Alternate, err = unmarshalExpression(x.Alternate)
	}
	if err == nil {
		ce.Consequent, err = unmarshalExpression(x.Consequent)
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
		ce.Callee, err = unmarshalExpression(x.Callee)
	}
	if err == nil {
		ce.Arguments = make([]Expression, len(x.Arguments))
		for i := range x.Arguments {
			if ce.Arguments[i], err = unmarshalExpression(x.Arguments[i]); err != nil {
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
		ne.Callee, err = unmarshalExpression(x.Callee)
	}
	if err == nil {
		ne.Arguments = make([]Expression, len(x.Arguments))
		for i := range x.Arguments {
			if ne.Arguments[i], err = unmarshalExpression(x.Arguments[i]); err != nil {
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
			if se.Expressions[i], err = unmarshalExpression(x.Expressions[i]); err != nil {
				break
			}
		}
	}
	return err
}
