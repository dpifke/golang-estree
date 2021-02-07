package estree

import (
	"encoding/json"
	"fmt"
)

type ReturnStatement struct {
	baseStatement
	Argument Expression
}

func (ReturnStatement) Type() string { return "ReturnStatement" }

func (rs ReturnStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     rs.Type(),
		"argument": rs.Argument,
	})
}

func (rs *ReturnStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Argument json.RawMessage `json:"argument"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != rs.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, rs.Type(), x.Type)
	}
	if err == nil {
		rs.Argument, err = unmarshalExpression(x.Argument)
	}
	return err
}

type LabeledStatement struct {
	baseStatement
	Label Identifier
	Body  Statement
}

func (LabeledStatement) Type() string { return "LabeledStatement" }

func (ls LabeledStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  ls.Type(),
		"label": ls.Label,
		"body":  ls.Body,
	})
}

func (ls *LabeledStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Identifier Identifier      `json:"identifier"`
		Body       json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ls.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ls.Type(), x.Type)
	}
	if err == nil {
		ls.Body, err = unmarshalStatement(x.Body)
	}
	return err
}

type BreakStatement struct {
	baseStatement
	Label Identifier
}

func (BreakStatement) Type() string { return "BreakStatement" }

func (bs BreakStatement) MarshalJSON() ([]byte, error) {
	x := map[string]interface{}{
		"type": bs.Type(),
	}
	if bs.Label.Name != "" {
		x["label"] = bs.Label
	}
	return json.Marshal(x)
}

func (bs *BreakStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string     `json:"type"`
		Identifier Identifier `json:"identifier"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != bs.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, bs.Type(), x.Type)
	}
	return err
}

type ContinueStatement struct {
	baseStatement
	Label Identifier
}

func (ContinueStatement) Type() string { return "ContinueStatement" }

func (cs ContinueStatement) MarshalJSON() ([]byte, error) {
	x := map[string]interface{}{
		"type": cs.Type(),
	}
	if cs.Label.Name != "" {
		x["label"] = cs.Label
	}
	return json.Marshal(x)
}

func (cs *ContinueStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string     `json:"type"`
		Identifier Identifier `json:"identifier"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != cs.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, cs.Type(), x.Type)
	}
	return err
}
