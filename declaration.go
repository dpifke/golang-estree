package estree

import (
	"encoding/json"
	"fmt"
)

// Declaration is any declaration node.  Note that declarations are considered
// statements; this is because declarations can appear in any statement
// context.
type Declaration interface {
	Statement
	isDeclaration()
}

type baseDeclaration struct {
	baseStatement
}

func (baseDeclaration) MinVersion() Version { return ES5 }
func (baseDeclaration) isDeclaration()      {}

// FunctionDeclaration declares a function. Note that unlike in the parent
// interface Function, ID cannot be nil.
type FunctionDeclaration struct {
	baseDeclaration
	Loc    SourceLocation
	ID     Identifier
	Params []Pattern
	Body   FunctionBody
}

func (FunctionDeclaration) Type() string                { return "FunctionDelaration" }
func (fd FunctionDeclaration) Location() SourceLocation { return fd.Loc }

func (fd FunctionDeclaration) IsZero() bool {
	return fd.Loc.IsZero() &&
		fd.ID.IsZero() &&
		len(fd.Params) == 0 &&
		fd.Body.IsZero()
}

func (fd FunctionDeclaration) Walk(v Visitor) {
	if v = v.Visit(fd); v != nil {
		defer v.Visit(nil)
		fd.ID.Walk(v)
	}
}

func (fd FunctionDeclaration) Errors() []error {
	c := nodeChecker{Node: fd}
	c.require(fd.ID, "function name")
	c.requireEach(nodeSlice{
		Index: func(i int) Node { return fd.Params[i] },
		Len:   len(fd.Params),
	}, "function parameter")
	c.require(fd.Body, "function body")
	return c.errors()
}

func (fd FunctionDeclaration) MarshalJSON() ([]byte, error) {
	x := nodeToMap(fd)
	x["id"] = fd.ID
	x["params"] = fd.Params
	x["body"] = fd.Body
	return json.Marshal(x)
}

func (fd *FunctionDeclaration) UnmarshalJSON(b []byte) error {
	var x struct {
		Type   string            `json:"type"`
		Loc    SourceLocation    `json:"loc"`
		ID     Identifier        `json:"id"`
		Params []json.RawMessage `json:"params"`
		Body   FunctionBody      `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != fd.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, fd.Type(), x.Type)
	}
	if err == nil && len(x.Params) > 0 {
		fd.Params = make([]Pattern, len(x.Params))
		for i := range x.Params {
			if fd.Params[i], _, err = unmarshalPattern(x.Params[i]); err != nil {
				break
			}
		}
	}
	if err == nil {
		fd.Loc, fd.ID, fd.Body = x.Loc, x.ID, x.Body
		fd.Params = make([]Pattern, len(x.Params))
		for i := range x.Params {
			var err2 error
			fd.Params[i], _, err2 = unmarshalPattern(x.Params[i])
			if err == nil && err2 != nil {
				err = err2
			}
		}
	}
	return err
}

// VariableDeclarationOrExpression is used where a Node can be either a
// VariableDeclaration, or any implementation of Expression.
type VariableDeclarationOrExpression interface {
	Node
	isVariableDeclarationOrExpression()
}

func unmarshalVariableDeclarationOrExpression(m json.RawMessage) (VariableDeclarationOrExpression, error) {
	if e, match, err := unmarshalExpression(m); match {
		return e, err
	}
	var x struct {
		Type string `json:"type"`
	}
	var err error
	if err = json.Unmarshal(m, &x); err == nil {
		var vd VariableDeclaration
		if x.Type == vd.Type() {
			if err = json.Unmarshal([]byte(m), &vd); err == nil {
				return vd, nil
			}
		} else {
			// TODO: use x.Type if != "", maybe truncate string(m) if long
			err = fmt.Errorf("%w VariableDeclaration or Expression, got %v", ErrWrongType, string(m))
		}
	}

	return nil, err
}

// VariableDeclarationOrPattern is used where a Node can be either a
// VariableDeclaration, or any implementation of Pattern.
type VariableDeclarationOrPattern interface {
	Node
	isVariableDeclarationOrPattern()
}

func unmarshalVariableDeclarationOrPattern(m json.RawMessage) (VariableDeclarationOrPattern, error) {
	if e, match, err := unmarshalPattern(m); match {
		return e, err
	}
	var x struct {
		Type string `json:"type"`
	}
	var err error
	if err = json.Unmarshal(m, &x); err == nil {
		var vd VariableDeclaration
		if x.Type == vd.Type() {
			if err = json.Unmarshal([]byte(m), &vd); err == nil {
				return vd, nil
			}
		} else {
			err = fmt.Errorf("%w VariableDeclaration or Pattern, got %v", ErrWrongType, string(m))
		}
	}
	return nil, err
}

// VariableDeclarationKind is the kind of VariableDeclaration.
type VariableDeclarationKind string

var (
	Var VariableDeclarationKind = "var"
)

func (vdk VariableDeclarationKind) GoString() string {
	if vdk == Var {
		return "Var"
	}
	return fmt.Sprintf("%q", vdk)
}

func (vdk VariableDeclarationKind) IsValid() bool {
	return vdk == Var
}

// VariableDeclaration is a group of VariableDeclarators.
type VariableDeclaration struct {
	baseDeclaration
	Loc          SourceLocation
	Declarations []VariableDeclarator
	Kind         VariableDeclarationKind
}

func (vd VariableDeclaration) IsZero() bool {
	return vd.Loc.IsZero() && len(vd.Declarations) == 0 && vd.Kind == ""
}

func (VariableDeclaration) Type() string                       { return "VariableDelaration" }
func (vd VariableDeclaration) Location() SourceLocation        { return vd.Loc }
func (VariableDeclaration) isVariableDeclarationOrExpression() {}
func (VariableDeclaration) isVariableDeclarationOrPattern()    {}

func (vd VariableDeclaration) Walk(v Visitor) {
	if v = v.Visit(vd); v != nil {
		defer v.Visit(nil)
		for i := range vd.Declarations {
			vd.Declarations[i].Walk(v)
		}
	}
}

func (vd VariableDeclaration) Errors() []error {
	c := nodeChecker{Node: vd}
	if !vd.Kind.IsValid() {
		c.appendf("%w VariableDeclarationKind %q", ErrWrongValue, vd.Kind)
	}
	c.requireEach(nodeSlice{
		Index: func(i int) Node { return vd.Declarations[i] },
		Len:   len(vd.Declarations),
	}, "variable declaration")
	return c.errors()
}

func (vd VariableDeclaration) MarshalJSON() ([]byte, error) {
	x := nodeToMap(vd)
	x["declarations"] = vd.Declarations
	x["kind"] = vd
	return json.Marshal(x)
}

func (vd *VariableDeclaration) UnmarshalJSON(b []byte) error {
	var x struct {
		Type         string                  `json:"type"`
		Loc          SourceLocation          `json:"loc"`
		Declarations []VariableDeclarator    `json:"declarations"`
		Kind         VariableDeclarationKind `json:"kind"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != vd.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, vd.Type(), x.Type)
	}
	if err == nil {
		vd.Loc, vd.Declarations = x.Loc, x.Declarations
		if x.Kind.IsValid() {
			vd.Kind = x.Kind
		} else {
			err = fmt.Errorf("%w VariableDeclaration.Kind %q", ErrWrongValue, x.Kind)
		}
	}
	return err
}

// VariableDeclarator defines a new variable identified by ID, optionally
// initialized to the result of Init.
type VariableDeclarator struct {
	Loc  SourceLocation
	ID   Pattern
	Init Expression // or nil
}

func (VariableDeclarator) Type() string                { return "VariableDelarator" }
func (vd VariableDeclarator) Location() SourceLocation { return vd.Loc }
func (VariableDeclarator) MinVersion() Version         { return ES5 }

func (vd VariableDeclarator) IsZero() bool {
	return vd.Loc.IsZero() &&
		(vd.ID == nil || vd.ID.IsZero()) &&
		(vd.Init == nil || vd.Init.IsZero())
}

func (vd VariableDeclarator) Walk(v Visitor) {
	if v = v.Visit(vd); v != nil {
		defer v.Visit(nil)
		vd.ID.Walk(v)
		if vd.Init != nil {
			vd.Init.Walk(v)
		}
	}
}

func (vd VariableDeclarator) Errors() []error {
	c := nodeChecker{Node: vd}
	c.require(vd.ID, "variable name")
	c.optional(vd.Init)
	return c.errors()
}

func (vd VariableDeclarator) MarshalJSON() ([]byte, error) {
	x := nodeToMap(vd)
	x["id"] = vd.ID
	if vd.Init != nil {
		x["init"] = vd.Init
	}
	return json.Marshal(x)
}

func (vd *VariableDeclarator) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string          `json:"type"`
		Loc  SourceLocation  `json:"loc"`
		ID   json.RawMessage `json:"id"`
		Init json.RawMessage `json:"init"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != vd.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, vd.Type(), x.Type)
	}
	if err == nil {
		vd.Loc = x.Loc
		vd.ID, _, err = unmarshalPattern(x.ID)
		if len(x.Init) > 0 {
			var err2 error
			if vd.Init, _, err2 = unmarshalExpression(x.Init); err == nil && err2 != nil {
				err = err2
			}
		}
	}
	return err
}
