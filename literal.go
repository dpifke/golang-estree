package estree

import (
	"encoding/json"
	"errors"
	"fmt"
)

// Literal is a literal token.  Note that a literal can be an expression.
type Literal interface {
	Expression
	isLiteral()
	isLiteralOrIdentifier()
	isVariableDeclarationOrLiteral()
}

func unmarshalLiteral(m json.RawMessage) (l Literal, match bool, err error) {
	var x struct {
		Type  string         `json:"type"`
		Loc   SourceLocation `json:"loc"`
		Value interface{}    `json:"value"`
	}
	if err = json.Unmarshal(m, &x); err != nil {
		match = true
	} else {
		if x.Type == (baseLiteral{}).Type() {
			switch v := x.Value.(type) {
			case string:
				l, match = StringLiteral{Loc: x.Loc, Value: v}, true
			case bool:
				l, match = BoolLiteral{Loc: x.Loc, Value: v}, true
			case nil:
				l, match = NullLiteral{Loc: x.Loc}, true
			case float64:
				l, match = NumberLiteral{Loc: x.Loc, Value: v}, true
			default:
				var re RegExpLiteral
				if err = json.Unmarshal([]byte(m), &re); err == nil {
					l, match = re, true
				} else if !errors.Is(err, ErrWrongType) {
					match = true
				} else {
					err = nil
				}
			}
		} else {
			err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, baseLiteral{}.Type(), x.Type)
		}
	}
	return
}

type LiteralOrIdentifier interface {
	Node
	isLiteralOrIdentifier()
}

func unmarshalLiteralOrIdentifier(m json.RawMessage) (LiteralOrIdentifier, bool, error) {
	if l, match, err := unmarshalLiteral(m); match {
		return l, true, err
	}
	var i Identifier
	if err := i.UnmarshalJSON([]byte(m)); !errors.Is(err, ErrWrongType) {
		return i, true, err
	}
	return nil, false, fmt.Errorf("%w: expected Literal or Identifier, got %v", ErrWrongType, string(m))
}

type baseLiteral struct {
	baseExpression
}

func (baseLiteral) Type() string                    { return "Literal" }
func (baseLiteral) MinVersion() Version             { return ES5 }
func (baseLiteral) IsZero() bool                    { return false }
func (baseLiteral) Errors() []error                 { return nil }
func (baseLiteral) isLiteral()                      {}
func (baseLiteral) isLiteralOrIdentifier()          {}
func (baseLiteral) isVariableDeclarationOrLiteral() {}

type StringLiteral struct {
	baseLiteral
	Loc   SourceLocation
	Value string
}

func (sl StringLiteral) Location() SourceLocation { return sl.Loc }

func (sl StringLiteral) Walk(v Visitor) {
	if v = v.Visit(sl); v != nil {
		v.Visit(nil)
	}
}

func (sl StringLiteral) MarshalJSON() ([]byte, error) {
	x := nodeToMap(sl)
	x["value"] = sl.Value
	return json.Marshal(x)
}

type BoolLiteral struct {
	baseLiteral
	Loc   SourceLocation
	Value bool
}

func (bl BoolLiteral) Location() SourceLocation { return bl.Loc }

func (bl BoolLiteral) Walk(v Visitor) {
	if v = v.Visit(bl); v != nil {
		v.Visit(nil)
	}
}

func (bl BoolLiteral) MarshalJSON() ([]byte, error) {
	x := nodeToMap(bl)
	x["value"] = bl.Value
	return json.Marshal(x)
}

type NullLiteral struct {
	baseLiteral
	Loc SourceLocation
}

func (nl NullLiteral) Location() SourceLocation { return nl.Loc }

func (nl NullLiteral) Walk(v Visitor) {
	if v = v.Visit(nl); v != nil {
		v.Visit(nil)
	}
}

func (nl NullLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(nodeToMap(nl))
}

type NumberLiteral struct {
	baseLiteral
	Loc   SourceLocation
	Value float64
}

func (nl NumberLiteral) Location() SourceLocation { return nl.Loc }

func (nl NumberLiteral) Walk(v Visitor) {
	if v = v.Visit(nl); v != nil {
		v.Visit(nil)
	}
}

func (nl NumberLiteral) MarshalJSON() ([]byte, error) {
	x := nodeToMap(nl)
	x["value"] = nl.Value
	return json.Marshal(x)
}

type RegExpLiteral struct {
	baseLiteral
	Loc     SourceLocation
	Pattern string
	Flags   string
}

func (rel RegExpLiteral) Location() SourceLocation { return rel.Loc }

// TODO: I think empty regex should still return false for IsZero; otherwise
// we should override IsZero and Errors here.

func (rel RegExpLiteral) Walk(v Visitor) {
	if v = v.Visit(rel); v != nil {
		v.Visit(nil)
	}
}

func (rel RegExpLiteral) MarshalJSON() ([]byte, error) {
	x := nodeToMap(rel)
	x["regex"] = map[string]interface{}{
		"pattern": rel.Pattern,
		"flags":   rel.Flags,
	}
	return json.Marshal(x)
}

func (rel *RegExpLiteral) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string         `json:"type"`
		Loc   SourceLocation `json:"loc"`
		Regex struct {
			Pattern string `json:"pattern"`
			Flags   string `json:"flags"`
		} `json:"regex"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != rel.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, rel.Type(), x.Type)
	}
	if err == nil {
		rel.Loc, rel.Pattern, rel.Flags = x.Loc, x.Regex.Pattern, x.Regex.Flags
	}
	return err
}
