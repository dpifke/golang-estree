package estree

import (
	"encoding/json"
	"fmt"
)

type ThrowStatement struct {
	baseStatement
	Argument Expression
}

func (ThrowStatement) Type() string { return "ThrowStatement" }

func (ts ThrowStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":     ts.Type(),
		"argument": ts.Argument,
	})
}

func (ts *ThrowStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Argument json.RawMessage `json:"argument"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ts.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ts.Type(), x.Type)
	}
	if err == nil {
		ts.Argument, _, err = unmarshalExpression(x.Argument)
	}
	return err
}

type TryStatement struct {
	baseStatement
	Block     BlockStatement
	Handler   CatchClause
	Finalizer Statement // BlockStatement or nil
}

func (TryStatement) Type() string { return "TryStatement" }

func (ts TryStatement) MarshalJSON() ([]byte, error) {
	x := map[string]interface{}{
		"type":  ts.Type(),
		"block": ts.Block,
	}
	if ts.Handler.Param != nil {
		x["handler"] = ts.Handler
	}
	if ts.Finalizer != nil {
		x["finalizer"] = ts.Finalizer
	}
	return json.Marshal(x)
}

func (ts *TryStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type      string          `json:"type"`
		Block     BlockStatement  `json:"block"`
		Handler   CatchClause     `json:"handler"`
		Finalizer json.RawMessage `json:"finalizer"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ts.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ts.Type(), x.Type)
	}
	if err == nil && len(x.Finalizer) > 0 {
		ts.Finalizer, _, err = unmarshalStatement(x.Finalizer)
	}
	return err
}

type CatchClause struct {
	Param Pattern
	Body  BlockStatement
}

func (CatchClause) Type() string { return "CatchClause" }

func (cc CatchClause) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": cc.Type(),
		"body": cc.Body,
	})
}

func (cc *CatchClause) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string          `json:"type"`
		Param json.RawMessage `json:"param"`
		Body  BlockStatement  `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != cc.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, cc.Type(), x.Type)
	}
	if err == nil {
		cc.Param, _, err = unmarshalPattern(x.Param)
	}
	return err
}
