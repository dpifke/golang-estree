package estree

import (
	"encoding/json"
	"fmt"
)

// ThrowStatement throws an exception, the result of Argument.
type ThrowStatement struct {
	baseStatement
	Loc      SourceLocation
	Argument Expression
}

func (ThrowStatement) Type() string                { return "ThrowStatement" }
func (ts ThrowStatement) Location() SourceLocation { return ts.Loc }

func (ts ThrowStatement) IsZero() bool {
	return ts.Loc.IsZero() && (ts.Argument == nil || ts.Argument.IsZero())
}

func (ts ThrowStatement) Walk(v Visitor) {
	if v = v.Visit(ts); v != nil {
		defer v.Visit(nil)
		if ts.Argument != nil {
			ts.Argument.Walk(v)
		}
	}
}

func (ts ThrowStatement) Errors() []error {
	return nil // TODO
}

func (ts ThrowStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ts)
	x["argument"] = ts.Argument
	return json.Marshal(x)
}

func (ts *ThrowStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Loc      SourceLocation  `json:"loc"`
		Argument json.RawMessage `json:"argument"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ts.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ts.Type(), x.Type)
	}
	if err == nil {
		ts.Loc = x.Loc
		ts.Argument, _, err = unmarshalExpression(x.Argument)
	}
	return err
}

// TryStatement executes Block, optionally with a Handler to execute if an
// exception is caught, and optionally with a Finalizer to execute afterwards.
type TryStatement struct {
	baseStatement
	Loc       SourceLocation
	Block     BlockStatement
	Handler   CatchClause
	Finalizer BlockStatement
}

func (TryStatement) Type() string                { return "TryStatement" }
func (ts TryStatement) Location() SourceLocation { return ts.Loc }

func (ts TryStatement) IsZero() bool {
	return ts.Loc.IsZero() &&
		ts.Block.IsZero() &&
		ts.Handler.IsZero() &&
		ts.Finalizer.IsZero()
}

func (ts TryStatement) Walk(v Visitor) {
	if v = v.Visit(ts); v != nil {
		defer v.Visit(nil)
		ts.Block.Walk(v)
		ts.Handler.Walk(v)
		ts.Finalizer.Walk(v)
	}
}

func (ts TryStatement) Errors() []error {
	return nil // TODO
}

func (ts TryStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ts)
	x["block"] = ts.Block
	if !ts.Handler.IsZero() {
		x["handler"] = ts.Handler
	}
	if !ts.Finalizer.IsZero() {
		x["finalizer"] = ts.Finalizer
	}
	return json.Marshal(x)
}

func (ts *TryStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type      string         `json:"type"`
		Loc       SourceLocation `json:"loc"`
		Block     BlockStatement `json:"block"`
		Handler   CatchClause    `json:"handler"`
		Finalizer BlockStatement `json:"finalizer"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ts.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ts.Type(), x.Type)
	}
	if err == nil {
		ts.Loc, ts.Block, ts.Handler, ts.Finalizer =
			x.Loc, x.Block, x.Handler, x.Finalizer
	}
	return err
}

// CatchClause is the catch clause following a try block.
type CatchClause struct {
	Loc   SourceLocation
	Param Pattern
	Body  BlockStatement
}

func (CatchClause) Type() string                { return "CatchClause" }
func (cc CatchClause) Location() SourceLocation { return cc.Loc }
func (CatchClause) MinVersion() Version         { return ES5 }

func (cc CatchClause) IsZero() bool {
	return cc.Loc.IsZero() &&
		(cc.Param == nil || cc.Param.IsZero()) &&
		cc.Body.IsZero()
}

func (cc CatchClause) Walk(v Visitor) {
	if v = v.Visit(cc); v != nil {
		defer v.Visit(nil)
		if cc.Param != nil {
			cc.Param.Walk(v)
		}
		cc.Body.Walk(v)
	}
}

func (cc CatchClause) Errors() []error {
	return nil // TODO
}

func (cc CatchClause) MarshalJSON() ([]byte, error) {
	x := nodeToMap(cc)
	x["param"] = cc.Param
	x["body"] = cc.Body
	return json.Marshal(x)
}

func (cc *CatchClause) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string          `json:"type"`
		Loc   SourceLocation  `json:"loc"`
		Param json.RawMessage `json:"param"`
		Body  BlockStatement  `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != cc.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, cc.Type(), x.Type)
	}
	if err == nil {
		cc.Loc, cc.Body = x.Loc, x.Body
		cc.Param, _, err = unmarshalPattern(x.Param)
	}
	return err
}
