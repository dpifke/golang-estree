package estree

import (
	"encoding/json"
	"fmt"
)

type WhileStatement struct {
	baseStatement
	Test Expression
	Body Statement
}

func (WhileStatement) Type() string { return "WhileStatement" }

func (ws WhileStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": ws.Type(),
		"test": ws.Test,
		"body": ws.Body,
	})
}

func (ws *WhileStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string          `json:"type"`
		Test json.RawMessage `json:"test"`
		Body json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ws.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ws.Type(), x.Type)
	}
	if err == nil && len(x.Test) > 0 {
		ws.Test, _, err = unmarshalExpression(x.Test)
	}
	if err == nil && len(x.Body) > 0 {
		ws.Body, _, err = unmarshalStatement(x.Body)
	}
	return err
}

type DoWhileStatement struct {
	baseStatement
	Body Statement
	Test Expression
}

func (DoWhileStatement) Type() string { return "DoWhileStatement" }

func (dws DoWhileStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": dws.Type(),
		"body": dws.Body,
		"test": dws.Test,
	})
}

func (dws *DoWhileStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string          `json:"type"`
		Body json.RawMessage `json:"body"`
		Test json.RawMessage `json:"test"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != dws.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, dws.Type(), x.Type)
	}
	if err == nil && len(x.Test) > 0 {
		dws.Test, _, err = unmarshalExpression(x.Test)
	}
	if err == nil && len(x.Body) > 0 {
		dws.Body, _, err = unmarshalStatement(x.Body)
	}
	return err
}

type ForStatement struct {
	baseStatement
	Init   VariableDeclarationOrExpression
	Test   Expression
	Update Expression
	Body   Statement
}

func (ForStatement) Type() string { return "ForStatement" }

func (fs ForStatement) MarshalJSON() ([]byte, error) {
	x := map[string]interface{}{
		"type": fs.Type(),
	}
	if fs.Init != nil {
		x["init"] = fs.Init
	}
	if fs.Test != nil {
		x["test"] = fs.Test
	}
	if fs.Update != nil {
		x["update"] = fs.Update
	}
	x["body"] = fs.Body
	return json.Marshal(x)
}

func (fs *ForStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type   string          `json:"type"`
		Init   json.RawMessage `json:"init"`
		Test   json.RawMessage `json:"test"`
		Update json.RawMessage `json:"update"`
		Body   json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != fs.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, fs.Type(), x.Type)
	}
	if err == nil && len(x.Init) > 0 {
		fs.Init, err = unmarshalVariableDeclarationOrExpression(x.Init)
	}
	if err == nil && len(x.Test) > 0 {
		fs.Test, _, err = unmarshalExpression(x.Test)
	}
	if err == nil && len(x.Update) > 0 {
		fs.Update, _, err = unmarshalExpression(x.Update)
	}
	if err == nil {
		fs.Body, _, err = unmarshalStatement(x.Body)
	}
	return err
}

type ForInStatement struct {
	baseStatement
	Left  VariableDeclarationOrPattern
	Right Expression
	Body  Statement
}

func (ForInStatement) Type() string { return "ForInStatement" }

func (fis ForInStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  fis.Type(),
		"left":  fis.Left,
		"right": fis.Right,
		"body":  fis.Body,
	})
}

func (fis *ForInStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string          `json:"type"`
		Left  json.RawMessage `json:"left"`
		Right json.RawMessage `json:"right"`
		Body  json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != fis.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, fis.Type(), x.Type)
	}
	if err == nil {
		fis.Left, err = unmarshalVariableDeclarationOrPattern(x.Left)
	}
	if err == nil {
		fis.Right, _, err = unmarshalExpression(x.Right)
	}
	if err == nil {
		fis.Body, _, err = unmarshalStatement(x.Body)
	}
	return err
}
