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

func (baseExpression) MinVersion() Version                { return ES5 }
func (baseExpression) isExpression()                      {}
func (baseExpression) isVariableDeclarationOrExpression() {}
func (baseExpression) isPatternOrExpression()             {}

type ThisExpression struct {
	baseExpression
	Loc SourceLocation
}

func (ThisExpression) Type() string                { return "ThisExpression" }
func (te ThisExpression) Location() SourceLocation { return te.Loc }

func (te ThisExpression) IsZero() bool {
	return false // no required members
}

func (te ThisExpression) Walk(v Visitor) {
	if v = v.Visit(te); v != nil {
		v.Visit(nil)
	}
}

func (te ThisExpression) Errors() []error {
	return nil // TODO
}

func (te ThisExpression) MarshalJSON() ([]byte, error) {
	return json.Marshal(nodeToMap(te))
}

func (te *ThisExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string         `json:"type"`
		Loc  SourceLocation `json:"loc"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != te.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, te.Type(), x.Type)
	}
	if err == nil {
		te.Loc = x.Loc
	}
	return err
}

// ArrayExpression is an array expression.
type ArrayExpression struct {
	baseExpression
	Loc SourceLocation

	// TODO: ExpressionOrNilArrayElement to support sparse arrays
	Elements []Expression
}

func (ArrayExpression) Type() string                { return "ArrayExpression" }
func (ae ArrayExpression) Location() SourceLocation { return ae.Loc }

func (ae ArrayExpression) IsZero() bool {
	return false // empty array is non-zero
}

func (ae ArrayExpression) Walk(v Visitor) {
	if v = v.Visit(ae); v != nil {
		defer v.Visit(nil)
		for _, e := range ae.Elements {
			if e != nil {
				e.Walk(v)
			}
		}
	}
}

func (ae ArrayExpression) Errors() []error {
	return nil // TODO
}

func (ae ArrayExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ae)
	x["elements"] = ae.Elements
	return json.Marshal(x)
}

func (ae *ArrayExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string            `json:"type"`
		Loc      SourceLocation    `json:"loc"`
		Elements []json.RawMessage `json:"elements"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ae.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ae.Type(), x.Type)
	}
	if err == nil {
		ae.Loc = x.Loc
		ae.Elements = make([]Expression, len(x.Elements))
		for i := range x.Elements {
			if len(x.Elements[i]) == 0 { // or x.Elements[i] unmarshals to nil?
				continue // TODO: insertNilArrayElement
			}
			var err2 error
			if ae.Elements[i], _, err2 = unmarshalExpression(x.Elements[i]); err == nil && err2 != nil {
				err = err2
			}
		}
	}
	return err
}

// ObjectExpression is an object expression.
type ObjectExpression struct {
	baseExpression
	Loc        SourceLocation
	Properties []Property
}

func (ObjectExpression) Type() string                { return "ObjectExpression" }
func (oe ObjectExpression) Location() SourceLocation { return oe.Loc }

func (oe ObjectExpression) IsZero() bool {
	return false // empty object is non-zero
}

func (oe ObjectExpression) Walk(v Visitor) {
	if v = v.Visit(oe); v != nil {
		defer v.Visit(nil)
		for _, p := range oe.Properties {
			p.Walk(v)
		}
	}
}

func (ae ObjectExpression) Errors() []error {
	return nil // TODO
}

func (oe ObjectExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(oe)
	x["properties"] = oe.Properties
	return json.Marshal(x)
}

func (oe *ObjectExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string         `json:"type"`
		Loc        SourceLocation `json:"loc"`
		Properties []Property     `json:"properties"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != oe.Type() {
		oe.Loc, oe.Properties = x.Loc, x.Properties
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, oe.Type(), x.Type)
	}
	return err
}

// PropertyKind is a value for Property.Kind, indicating whether the literal
// property in an ObjectExpression is an initializer, getter, or setter.
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

func (pk PropertyKind) IsValid() bool {
	switch pk {
	case Init, Get, Set:
		return true
	}
	return false
}

// Property is a literal property in an ObjectExpression.
type Property struct {
	Loc SourceLocation
	Key LiteralOrIdentifier

	// Value can be either a string or a number.
	Value Expression

	Kind PropertyKind
}

func (Property) Type() string               { return "Property" }
func (p Property) Location() SourceLocation { return p.Loc }
func (Property) MinVersion() Version        { return ES5 }

func (p Property) IsZero() bool {
	return p.Loc.IsZero() &&
		(p.Key == nil || p.Key.IsZero()) &&
		(p.Value == nil || p.Value.IsZero()) &&
		p.Kind == ""
}

func (p Property) Walk(v Visitor) {
	if v = v.Visit(p); v != nil {
		defer v.Visit(nil)
		if p.Key != nil {
			p.Key.Walk(v)
		}
		if p.Value != nil {
			p.Value.Walk(v)
		}
	}
}

func (p Property) Errors() []error {
	return nil // TODO
}

func (p Property) MarshalJSON() ([]byte, error) {
	x := nodeToMap(p)
	x["key"] = p.Key
	x["value"] = p.Value
	x["kind"] = p.Kind
	return json.Marshal(x)
}

func (p *Property) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string          `json:"type"`
		Loc   SourceLocation  `json:"loc"`
		Key   json.RawMessage `json:"key"`
		Value json.RawMessage `json:"value"`
		Kind  PropertyKind    `json:"kind"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != p.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, p.Type(), x.Type)
	}
	if err == nil {
		p.Loc = x.Loc
		if x.Kind.IsValid() {
			p.Kind = x.Kind
		} else {
			err = fmt.Errorf("%w for Property.Kind: %q", ErrWrongValue, x.Kind)
		}
		var err2 error
		if p.Key, _, err2 = unmarshalLiteralOrIdentifier(x.Key); err == nil && err2 != nil {
			err = err2
		}
		if p.Value, _, err = unmarshalExpression(x.Value); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

// FunctionExpression is a function expression (closure).
type FunctionExpression struct {
	baseExpression
	Loc    SourceLocation
	ID     Identifier
	Params []Pattern
	Body   FunctionBody
}

func (FunctionExpression) Type() string                { return "FunctionExpression" }
func (fe FunctionExpression) Location() SourceLocation { return fe.Loc }

func (fe FunctionExpression) IsZero() bool {
	return false // empty function is non-zero
}

func (fe FunctionExpression) Walk(v Visitor) {
	if v = v.Visit(fe); v != nil {
		v.Visit(nil)
	}
}

func (fe FunctionExpression) Errors() []error {
	return nil // TODO
}

func (fe FunctionExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(fe)
	x["id"] = fe.ID
	x["params"] = fe.Params
	x["body"] = fe.Body
	return json.Marshal(x)
}

func (fe *FunctionExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type   string            `json:"type"`
		Loc    SourceLocation    `json:"loc"`
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
		fe.Loc, fe.ID, fe.Body = x.Loc, x.ID, x.Body
	}
	return err
}

// ConditionalExpression is a ternary (x ? y : z) expression.
type ConditionalExpression struct {
	baseExpression
	Loc                         SourceLocation
	Test, Consequent, Alternate Expression
}

func (ConditionalExpression) Type() string                { return "ConditionalExpression" }
func (ce ConditionalExpression) Location() SourceLocation { return ce.Loc }

func (ce ConditionalExpression) IsZero() bool {
	return ce.Loc.IsZero() &&
		(ce.Test == nil || ce.Test.IsZero()) &&
		(ce.Alternate == nil || ce.Alternate.IsZero()) &&
		(ce.Consequent == nil || ce.Consequent.IsZero())
}

func (ce ConditionalExpression) Walk(v Visitor) {
	if v = v.Visit(ce); v != nil {
		defer v.Visit(nil)
		if ce.Test != nil {
			ce.Test.Walk(v)
		}
		if ce.Consequent != nil {
			ce.Consequent.Walk(v)
		}
		if ce.Alternate != nil {
			ce.Alternate.Walk(v)
		}
	}
}

func (ce ConditionalExpression) Errors() []error {
	return nil // TODO
}

func (ce ConditionalExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ce)
	x["test"] = ce.Test
	x["consequent"] = ce.Consequent
	x["alternate"] = ce.Alternate
	return json.Marshal(x)
}

func (ce *ConditionalExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Loc        SourceLocation  `json:"loc"`
		Test       json.RawMessage `json:"test"`
		Consequent json.RawMessage `json:"consequent"`
		Alternate  json.RawMessage `json:"alternate"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ce.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ce.Type(), x.Type)
	}
	if err == nil {
		ce.Loc = x.Loc
		ce.Test, _, err = unmarshalExpression(x.Test)
		var err2 error
		if ce.Alternate, _, err2 = unmarshalExpression(x.Alternate); err == nil && err2 != nil {
			err = err2
		}
		if ce.Consequent, _, err = unmarshalExpression(x.Consequent); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

// CallExpression is an expression which returns the result of a function or
// method call.
type CallExpression struct {
	baseExpression
	Loc       SourceLocation
	Callee    Expression
	Arguments []Expression
}

func (CallExpression) Type() string                { return "CallExpression" }
func (ce CallExpression) Location() SourceLocation { return ce.Loc }

func (ce CallExpression) IsZero() bool {
	return ce.Loc.IsZero() &&
		(ce.Callee == nil || ce.Callee.IsZero()) &&
		len(ce.Arguments) == 0
}

func (ce CallExpression) Walk(v Visitor) {
	if v = v.Visit(ce); v != nil {
		defer v.Visit(nil)
		if ce.Callee != nil {
			ce.Callee.Walk(v)
		}
		for _, a := range ce.Arguments {
			a.Walk(v)
		}
	}
}

func (ce CallExpression) Errors() []error {
	return nil // TODO
}

func (ce CallExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ce)
	x["callee"] = ce.Callee
	x["arguments"] = ce.Arguments
	return json.Marshal(ce)
}

func (ce *CallExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type      string            `json:"type"`
		Loc       SourceLocation    `json:"loc"`
		Callee    json.RawMessage   `json:"callee"`
		Arguments []json.RawMessage `json:"arguments"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ce.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ce.Type(), x.Type)
	}
	if err == nil {
		ce.Loc = x.Loc
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

// NewExpression is an expression which calls a constructor.
type NewExpression struct {
	baseExpression
	Loc       SourceLocation
	Callee    Expression
	Arguments []Expression
}

func (NewExpression) Type() string                { return "NewExpression" }
func (ne NewExpression) Location() SourceLocation { return ne.Loc }

func (ne NewExpression) IsZero() bool {
	return ne.Loc.IsZero() &&
		(ne.Callee == nil || ne.Callee.IsZero()) &&
		len(ne.Arguments) == 0
}

func (ne NewExpression) Walk(v Visitor) {
	if v = v.Visit(ne); v != nil {
		defer v.Visit(nil)
		if ne.Callee != nil {
			ne.Callee.Walk(v)
		}
		for _, a := range ne.Arguments {
			a.Walk(v)
		}
	}
}

func (ne NewExpression) Errors() []error {
	return nil // TODO
}

func (ne NewExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ne)
	x["callee"] = ne.Callee
	x["arguments"] = ne.Arguments
	return json.Marshal(x)
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
		ne.Arguments = make([]Expression, len(x.Arguments))
		for i := range x.Arguments {
			var err2 error
			if ne.Arguments[i], _, err2 = unmarshalExpression(x.Arguments[i]); err == nil && err2 != nil {
				err = err2
			}
		}
	}
	return err
}

// SequenceExpression is a comma-separated sequence of expressions.
type SequenceExpression struct {
	baseExpression
	Loc         SourceLocation
	Expressions []Expression
}

func (SequenceExpression) Type() string                { return "SequenceExpression" }
func (se SequenceExpression) Location() SourceLocation { return se.Loc }

func (se SequenceExpression) IsZero() bool {
	return se.Loc.IsZero() && len(se.Expressions) == 0
}

func (se SequenceExpression) Walk(v Visitor) {
	if v = v.Visit(se); v != nil {
		defer v.Visit(nil)
		for _, e := range se.Expressions {
			e.Walk(v)
		}
	}
}

func (se SequenceExpression) Errors() []error {
	return nil // TODO
}

func (se SequenceExpression) MarshalJSON() ([]byte, error) {
	x := nodeToMap(se)
	x["expressions"] = se.Expressions
	return json.Marshal(x)
}

func (se *SequenceExpression) UnmarshalJSON(b []byte) error {
	var x struct {
		Type        string            `json:"type"`
		Loc         SourceLocation    `json:"loc"`
		Expressions []json.RawMessage `json:"expressions"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != se.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, se.Type(), x.Type)
	}
	if err == nil {
		se.Loc = x.Loc
		se.Expressions = make([]Expression, len(x.Expressions))
		for i := range x.Expressions {
			var err2 error
			if se.Expressions[i], _, err2 = unmarshalExpression(x.Expressions[i]); err == nil && err2 != nil {
				err = err2
			}
		}
	}
	return err
}
