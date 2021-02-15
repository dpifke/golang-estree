package estree

import (
	"encoding/json"
	"fmt"
)

// IfStatement is a conditional branch.  If Test is true, Consequent will be
// executed, otherwise Alternate.  Alternate may be nil.
type IfStatement struct {
	baseStatement
	Loc        SourceLocation
	Test       Expression
	Consequent Statement
	Alternate  Statement // or nil
}

func (IfStatement) Type() string                { return "IfStatement" }
func (is IfStatement) Location() SourceLocation { return is.Loc }

func (is IfStatement) IsZero() bool {
	return is.Loc.IsZero() &&
		(is.Test == nil || is.Test.IsZero()) &&
		(is.Consequent == nil || is.Consequent.IsZero()) &&
		(is.Alternate == nil || is.Alternate.IsZero())
}

func (is IfStatement) Walk(v Visitor) {
	if v = v.Visit(is); v != nil {
		defer v.Visit(nil)
		if is.Test != nil {
			is.Test.Walk(v)
		}
		if is.Consequent != nil {
			is.Consequent.Walk(v)
		}
		if is.Alternate != nil {
			is.Alternate.Walk(v)
		}
	}
}

func (is IfStatement) Errors() []error {
	c := nodeChecker{Node: is}
	c.require(is.Test, "if expression")
	c.require(is.Consequent, "if statement block")
	c.optional(is.Alternate)
	return c.errors()
}

func (is IfStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(is)
	x["test"] = is.Test
	x["consequent"] = is.Consequent
	if is.Alternate != nil && !is.Alternate.IsZero() {
		x["alternate"] = is.Alternate
	}
	return json.Marshal(x)
}

func (is *IfStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Loc        SourceLocation  `json:"loc"`
		Test       json.RawMessage `json:"test"`
		Consequent json.RawMessage `json:"consequent"`
		Alternate  json.RawMessage `json:"alternate"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != is.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, is.Type(), x.Type)
	}
	if err == nil {
		is.Loc = x.Loc
		is.Test, _, err = unmarshalExpression(x.Test)
		var err2 error
		if is.Consequent, _, err2 = unmarshalStatement(x.Consequent); err == nil && err2 != nil {
			err = err2
		}
		if is.Alternate, _, err2 = unmarshalStatement(x.Alternate); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

// SwitchStatement is a conditional branch, consisting of zero more
// SwitchCases clauses.
type SwitchStatement struct {
	baseStatement
	Loc          SourceLocation
	Discriminant Expression
	Cases        []SwitchCase
}

func (SwitchStatement) Type() string                { return "SwitchStatement" }
func (ss SwitchStatement) Location() SourceLocation { return ss.Loc }

func (ss SwitchStatement) IsZero() bool {
	return ss.Loc.IsZero() &&
		(ss.Discriminant == nil || ss.Discriminant.IsZero()) &&
		len(ss.Cases) == 0
}

func (ss SwitchStatement) Walk(v Visitor) {
	if v = v.Visit(ss); v != nil {
		defer v.Visit(nil)
		if ss.Discriminant != nil {
			ss.Discriminant.Walk(v)
		}
		for _, sc := range ss.Cases {
			sc.Walk(v)
		}
	}
}

func (ss SwitchStatement) Errors() []error {
	c := nodeChecker{Node: ss}
	c.require(ss.Discriminant, "switch expression")
	return c.errors()
}

func (ss SwitchStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ss)
	x["test"] = ss.Discriminant
	if len(ss.Cases) > 0 {
		x["cases"] = ss.Cases
	}
	return json.Marshal(x)
}

func (ss *SwitchStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type         string          `json:"type"`
		Loc          SourceLocation  `json:"loc"`
		Discriminant json.RawMessage `json:"test"`
		Cases        []SwitchCase    `json:"cases"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ss.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, ss.Type(), x.Type)
	}
	if err == nil {
		ss.Loc, ss.Cases = x.Loc, x.Cases
		ss.Discriminant, _, err = unmarshalExpression(x.Discriminant)
	}
	return err
}

// SwitchCase is a branch of a SwitchStatement.  Consequent is executed if
// Test evaluates true.
//
// If Test is nil, this SwitchCase is the default clause.
type SwitchCase struct {
	Loc        SourceLocation
	Test       Expression // or nil
	Consequent []Statement
}

func (SwitchCase) Type() string                { return "SwitchCase" }
func (sc SwitchCase) Location() SourceLocation { return sc.Loc }
func (SwitchCase) MinVersion() Version         { return ES5 }
func (SwitchCase) IsZero() bool                { return false }

func (sc SwitchCase) Walk(v Visitor) {
	if v = v.Visit(sc); v != nil {
		defer v.Visit(nil)
		if sc.Test != nil && !sc.Test.IsZero() {
			sc.Test.Walk(v)
		}
		for _, c := range sc.Consequent {
			if c != nil {
				c.Walk(v)
			}
		}
	}
}

func (sc SwitchCase) Errors() []error {
	c := nodeChecker{Node: sc}
	c.optional(sc.Test)
	c.requireEach(nodeSlice{
		Index: func(i int) Node { return sc.Consequent[i] },
		Len:   len(sc.Consequent),
	}, "switch case statement")
	return c.errors()
}

func (sc SwitchCase) MarshalJSON() ([]byte, error) {
	x := nodeToMap(sc)
	if sc.Test != nil {
		x["test"] = sc.Test
	}
	if sc.Consequent != nil {
		x["consequent"] = sc.Consequent
	}
	return json.Marshal(x)
}

func (sc *SwitchCase) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string            `json:"type"`
		Loc        SourceLocation    `json:"loc"`
		Test       json.RawMessage   `json:"test"`
		Consequent []json.RawMessage `json:"consequent"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != sc.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, sc.Type(), x.Type)
	}
	if err == nil {
		sc.Loc = x.Loc
		sc.Test, _, err = unmarshalExpression(x.Test)
		if len(x.Consequent) == 0 {
			sc.Consequent = nil
		} else {
			sc.Consequent = make([]Statement, len(x.Consequent))
			for i := range x.Consequent {
				var err2 error
				sc.Consequent[i], _, err2 = unmarshalStatement(x.Consequent[i])
				if err == nil && err2 != nil {
					err = err2
				}
			}
		}
	}
	return err
}
