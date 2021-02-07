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
}

type LiteralOrIdentifier interface {
	Node

	isLiteralOrIdentifier()
}

type baseLiteral struct {
	baseExpression
}

func (baseLiteral) Type() string                    { return "Literal" }
func (baseLiteral) isLiteral()                      {}
func (baseLiteral) isVariableDeclarationOrLiteral() {}
func (baseLiteral) isLiteralOrIdentifier()          {}

func (bl baseLiteral) marshalJSON(value interface{}) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  bl.Type(),
		"value": value,
	})
}

func unmarshalLiteral(m json.RawMessage) (Literal, error) {
	var x interface{}
	if err := json.Unmarshal([]byte(m), x); err != nil {
		return nil, err
	}
	switch x := x.(type) {
	case string:
		return StringLiteral{Value: x}, nil
	case bool:
		return BoolLiteral{Value: x}, nil
	case nil:
		return NullLiteral{}, nil
	case float64:
		return NumberLiteral{Value: x}, nil
	}
	var re RegExpLiteral
	if err := json.Unmarshal([]byte(m), &re); err == nil {
		return re, nil
	} else if !errors.Is(err, ErrWrongType) {
		return nil, err
	}
	return nil, fmt.Errorf("%w: expected Literal, got %v", ErrWrongType, string(m))
}

type StringLiteral struct {
	baseLiteral
	Value string
}

func (sl StringLiteral) MarshalJSON() ([]byte, error) {
	return sl.marshalJSON(sl.Value)
}

type BoolLiteral struct {
	baseLiteral
	Value bool
}

func (bl BoolLiteral) MarshalJSON() ([]byte, error) {
	return bl.marshalJSON(bl.Value)
}

type NullLiteral struct {
	baseLiteral
}

func (nl NullLiteral) MarshalJSON() ([]byte, error) {
	return nl.marshalJSON(nil)
}

type NumberLiteral struct {
	baseLiteral
	Value float64
}

func (nl NumberLiteral) MarshalJSON() ([]byte, error) {
	return nl.marshalJSON(nl.Value)
}

type RegExpLiteral struct {
	baseLiteral
	Pattern string
	Flags   string
}

func (rel RegExpLiteral) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": "Literal",
		"regex": map[string]interface{}{
			"pattern": rel.Pattern,
			"flags":   rel.Flags,
		},
	})
}

func (rel *RegExpLiteral) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string `json:"type"`
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
		rel.Pattern = x.Regex.Pattern
		rel.Flags = x.Regex.Flags
	}
	return err
}
