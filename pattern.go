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

func unmarshalPattern(m json.RawMessage) (Pattern, bool, error) {
	var i Identifier
	err := json.Unmarshal(m, &i)
	if err == nil {
		return i, true, nil
	}
	if errors.Is(err, ErrWrongType) {
		return nil, false, fmt.Errorf("%w Pattern, got %v", ErrWrongType, string(m))
	}
	return nil, true, err
}

type PatternOrExpression interface {
	Node
	isPatternOrExpression()
}

func unmarshalPatternOrExpression(m json.RawMessage) (PatternOrExpression, error) {
	if p, match, err := unmarshalPattern(m); match {
		return p, err
	}
	if e, match, err := unmarshalExpression(m); match {
		return e, err
	}
	return nil, fmt.Errorf("%w Pattern or Expression, got %v", ErrWrongType, string(m))
}

type basePattern struct {
	loc SourceLocation
}

func (p basePattern) Loc() *SourceLocation          { return &p.loc }
func (basePattern) MinVersion() Version             { return ES5 }
func (basePattern) isPattern()                      {}
func (basePattern) isVariableDeclarationOrPattern() {}
func (basePattern) isPatternOrExpression()          {}
