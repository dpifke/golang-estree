package estree

import (
	"encoding/json"
	"fmt"
)

// ReturnStatement is a return from a function.
type ReturnStatement struct {
	baseStatement
	Loc      SourceLocation
	Argument Expression
}

func (ReturnStatement) Type() string                { return "ReturnStatement" }
func (rs ReturnStatement) Location() SourceLocation { return rs.Loc }

func (rs ReturnStatement) IsZero() bool {
	return rs.Loc.IsZero() &&
		(rs.Argument == nil || rs.Argument.IsZero())
}

func (rs ReturnStatement) Walk(v Visitor) {
	if v = v.Visit(rs); v != nil {
		defer v.Visit(nil)
		if rs.Argument != nil {
			rs.Argument.Walk(v)
		}
	}
}

func (rs ReturnStatement) Errors() []error {
	return nil // TODO
}

func (rs ReturnStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(rs)
	x["argument"] = rs.Argument
	return json.Marshal(x)
}

func (rs *ReturnStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type     string          `json:"type"`
		Loc      SourceLocation  `json:"loc"`
		Argument json.RawMessage `json:"argument"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != rs.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, rs.Type(), x.Type)
	}
	if err == nil {
		rs.Loc = x.Loc
		rs.Argument, _, err = unmarshalExpression(x.Argument)
	}
	return err
}

// LabeledStatement is a Statement prefixed by an identifier, e.g. a break or
// continue label.
type LabeledStatement struct {
	baseStatement
	Loc   SourceLocation
	Label Identifier
	Body  Statement
}

func (LabeledStatement) Type() string                { return "LabeledStatement" }
func (ls LabeledStatement) Location() SourceLocation { return ls.Loc }

func (ls LabeledStatement) IsZero() bool {
	return ls.Loc.IsZero() &&
		ls.Label.IsZero() &&
		(ls.Body == nil || ls.Body.IsZero())
}

func (ls LabeledStatement) Walk(v Visitor) {
	if v = v.Visit(ls); v != nil {
		defer v.Visit(nil)
		ls.Label.Walk(v)
		if ls.Body != nil {
			ls.Body.Walk(v)
		}
	}
}

func (ls LabeledStatement) Errors() []error {
	return nil // TODO
}

func (ls LabeledStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ls)
	x["label"] = ls.Label
	x["body"] = ls.Body
	return json.Marshal(x)
}

func (ls *LabeledStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string          `json:"type"`
		Loc   SourceLocation  `json:"loc"`
		Label Identifier      `json:"label"`
		Body  json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ls.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ls.Type(), x.Type)
	}
	if err == nil {
		ls.Loc, ls.Label = x.Loc, x.Label
		ls.Body, _, err = unmarshalStatement(x.Body)
	}
	return err
}

// BreakStatement exits a loop or SwitchStatement.
type BreakStatement struct {
	baseStatement
	Loc   SourceLocation
	Label Identifier
}

func (BreakStatement) Type() string                { return "BreakStatement" }
func (bs BreakStatement) Location() SourceLocation { return bs.Loc }

func (bs BreakStatement) IsZero() bool {
	return bs.Loc.IsZero() && bs.Label.IsZero()
}

func (bs BreakStatement) Walk(v Visitor) {
	if v = v.Visit(bs); v != nil {
		defer v.Visit(nil)
		bs.Label.Walk(v)
	}
}

func (bs BreakStatement) Errors() []error {
	return nil // TODO
}

func (bs BreakStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(bs)
	if bs.Label.Name != "" {
		x["label"] = bs.Label
	}
	return json.Marshal(x)
}

func (bs *BreakStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string         `json:"type"`
		Loc   SourceLocation `json:"loc"`
		Label Identifier     `json:"label"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != bs.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, bs.Type(), x.Type)
	}
	if err == nil {
		bs.Loc, bs.Label = x.Loc, x.Label
	}
	return err
}

// ContinueStatement skips the remainder of the current loop.
type ContinueStatement struct {
	baseStatement
	Loc   SourceLocation
	Label Identifier
}

func (ContinueStatement) Type() string                { return "ContinueStatement" }
func (cs ContinueStatement) Location() SourceLocation { return cs.Loc }

func (cs ContinueStatement) IsZero() bool {
	return cs.Loc.IsZero() && cs.Label.IsZero()
}

func (cs ContinueStatement) Walk(v Visitor) {
	if v = v.Visit(cs); v != nil {
		defer v.Visit(nil)
		cs.Label.Walk(v)
	}
}

func (cs ContinueStatement) Errors() []error {
	return nil // TODO
}

func (cs ContinueStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(cs)
	if cs.Label.Name != "" {
		x["label"] = cs.Label
	}
	return json.Marshal(x)
}

func (cs *ContinueStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string         `json:"type"`
		Loc   SourceLocation `json:"loc"`
		Label Identifier     `json:"label"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != cs.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, cs.Type(), x.Type)
	}
	if err == nil {
		cs.Loc, cs.Label = x.Loc, x.Label
	}
	return err
}
