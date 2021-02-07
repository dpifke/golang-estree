package estree

import (
	"encoding/json"
	"fmt"
)

type Declaration interface {
	Statement

	isDeclaration()
}

type baseDeclaration struct {
	baseStatement
}

func (baseDeclaration) isDeclaration() {}

type FunctionDeclaration struct {
	baseDeclaration
	ID     Identifier
	Params []Pattern
	Body   FunctionBody
}

func (FunctionDeclaration) Type() string { return "FunctionDelaration" }

func (fd FunctionDeclaration) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":   fd.Type(),
		"id":     fd.ID,
		"params": fd.Params,
		"body":   fd.Body,
	})
}

func (fd *FunctionDeclaration) UnmarshalJSON(b []byte) error {
	var x struct {
		Type   string            `json:"type"`
		ID     Identifier        `json:"id"`
		Params []json.RawMessage `json:"params"`
		Body   FunctionBody      `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != fd.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, fd.Type(), x.Type)
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
		fd.ID, fd.Body = x.ID, x.Body
	}
	return err
}

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
			err = fmt.Errorf("%w: expected VariableDeclaration or Expression, got %v", ErrWrongType, string(m))
		}
	}

	return nil, err
}

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
			err = fmt.Errorf("%w: expected VariableDeclaration or Pattern, got %v", ErrWrongType, string(m))
		}
	}
	return nil, err
}

type VariableDeclarationKind string

var (
	Var VariableDeclarationKind = "var"
)

func (vdk VariableDeclarationKind) GoString() string {
	if vdk == "Var" {
		return "Var"
	}
	return fmt.Sprintf("%q", vdk)
}

type VariableDeclaration struct {
	baseDeclaration
	Declarations []VariableDeclarator
	Kind         VariableDeclarationKind
}

func (VariableDeclaration) Type() string                       { return "VariableDelaration" }
func (VariableDeclaration) isVariableDeclarationOrExpression() {}
func (VariableDeclaration) isVariableDeclarationOrPattern()    {}

func (vd VariableDeclaration) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":         vd.Type(),
		"declarations": vd.Declarations,
		"kind":         vd,
	})
}

func (vd *VariableDeclaration) UnmarshalJSON(b []byte) error {
	var x struct {
		Type         string                  `json:"type"`
		Declarations []VariableDeclarator    `json:"declarations"`
		Kind         VariableDeclarationKind `json:"kind"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != vd.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, vd.Type(), x.Type)
	}
	if err == nil {
		if x.Kind == Var {
			vd.Kind = x.Kind
		} else {
			err = fmt.Errorf("invalid VariableDeclaration.Kind %q", x.Kind)
		}
	}
	return err
}

type VariableDeclarator struct {
	ID   Pattern
	Init Expression
}

func (VariableDeclarator) Type() string { return "VariableDelarator" }

func (vd VariableDeclarator) MarshalJSON() ([]byte, error) {
	x := map[string]interface{}{
		"type": vd.Type(),
		"id":   vd.ID,
	}
	if vd.Init != nil {
		x["init"] = vd.Init
	}
	return json.Marshal(x)
}

func (vd *VariableDeclarator) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string          `json:"type"`
		ID   json.RawMessage `json:"id"`
		Init json.RawMessage `json:"init"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != vd.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, vd.Type(), x.Type)
	}
	if err == nil {
		vd.ID, _, err = unmarshalPattern(x.ID)
	}
	if err == nil && len(x.Init) > 0 {
		vd.Init, _, err = unmarshalExpression(x.Init)
	}
	return err
}
