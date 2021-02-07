package estree

import (
	"encoding/json"
	"errors"
	"fmt"
)

type Pattern interface {
	Node

	isPattern()
	isVariableDeclarationOrPattern()
	isPatternOrExpression()
}

func unmarshalPattern(m json.RawMessage) (Pattern, error) {
	var i Identifier
	err := json.Unmarshal(m, &i)
	if err == nil {
		return i, nil
	}
	if errors.Is(err, ErrWrongType) {
		err = fmt.Errorf("%w: expected Pattern, got %v", ErrWrongType, string(m))
	}
	return nil, err
}

type PatternOrExpression interface {
	Node

	isPatternOrExpression()
}

func unmarshalPatternOrExpression(m json.RawMessage) (PatternOrExpression, error) {
	if p, err := unmarshalPattern(m); err == nil {
		return p, nil
	} else if !errors.Is(err, ErrWrongType) {
		return nil, err
	}
	if e, err := unmarshalExpression(m); err == nil {
		return e, nil
	} else if !errors.Is(err, ErrWrongType) {
		return nil, err
	}
	return nil, fmt.Errorf("%w: expected Pattern or Expression, got %v", ErrWrongType, string(m))
}

type basePattern struct{}

func (basePattern) isPattern()                      {}
func (basePattern) isVariableDeclarationOrPattern() {}
func (basePattern) isPatternOrExpression()          {}
