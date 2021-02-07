package estree

import (
	"encoding/json"
	"errors"
	"fmt"
)

// DirectiveOrStatement is either a Directive or a Statement.
type DirectiveOrStatement interface {
	Node

	isDirectiveOrStatement()
}

func unmarshalDirectiveOrStatement(m json.RawMessage) (DirectiveOrStatement, error) {
	var s Statement
	if err := json.Unmarshal([]byte(m), &s); err == nil {
		return s, nil
	} else if !errors.Is(err, ErrWrongType) {
		return nil, err
	}

	var d Directive
	if err := json.Unmarshal([]byte(m), &d); err == nil {
		return d, nil
	} else if !errors.Is(err, ErrWrongType) {
		return nil, err
	}

	return nil, fmt.Errorf("%w: expected Directive or Statement, got %v", ErrWrongType, string(m))
}

// Program is a complete program source tree.
type Program struct {
	Body []DirectiveOrStatement
}

func (Program) Type() string { return "Program" }

func (p Program) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": p.Type(),
		"body": p.Body,
	})
}

func (p *Program) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string            `json:"type"`
		Body []json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != p.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, p.Type(), x.Type)
	}
	if err == nil {
		p.Body = make([]DirectiveOrStatement, len(x.Body))
		for i := range x.Body {
			if p.Body[i], err = unmarshalDirectiveOrStatement(x.Body[i]); err != nil {
				break
			}
		}
	}
	return err
}
