package estree

import (
	"encoding/json"
	"fmt"
)

// WhileStatement is a while loop.
type WhileStatement struct {
	baseStatement
	Loc  SourceLocation
	Test Expression
	Body Statement
}

func (WhileStatement) Type() string                { return "WhileStatement" }
func (ws WhileStatement) Location() SourceLocation { return ws.Loc }

func (ws WhileStatement) IsZero() bool {
	return ws.Loc.IsZero() &&
		(ws.Test == nil || ws.Test.IsZero()) &&
		(ws.Body == nil || ws.Body.IsZero())
}

func (ws WhileStatement) Walk(v Visitor) {
	if v = v.Visit(ws); v != nil {
		defer v.Visit(nil)
		if ws.Test != nil {
			ws.Test.Walk(v)
		}
		if ws.Body != nil {
			ws.Body.Walk(v)
		}
	}
}

func (ws WhileStatement) Errors() []error {
	return nil // TODO
}

func (ws WhileStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ws)
	x["test"] = ws.Test
	x["body"] = ws.Body
	return json.Marshal(x)
}

func (ws *WhileStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string          `json:"type"`
		Loc  SourceLocation  `json:"loc"`
		Test json.RawMessage `json:"test"`
		Body json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ws.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ws.Type(), x.Type)
	}
	if err == nil {
		ws.Loc = x.Loc
		ws.Test, _, err = unmarshalExpression(x.Test)
		var err2 error
		if ws.Body, _, err2 = unmarshalStatement(x.Body); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

// DoWhileStatement is a do / while loop.
type DoWhileStatement struct {
	baseStatement
	Loc  SourceLocation
	Body Statement
	Test Expression
}

func (DoWhileStatement) Type() string                 { return "DoWhileStatement" }
func (dws DoWhileStatement) Location() SourceLocation { return dws.Loc }

func (dws DoWhileStatement) IsZero() bool {
	return dws.Loc.IsZero() &&
		(dws.Body == nil || dws.Body.IsZero()) &&
		(dws.Test == nil || dws.Test.IsZero())
}

func (dws DoWhileStatement) Walk(v Visitor) {
	if v = v.Visit(dws); v != nil {
		defer v.Visit(nil)
		if dws.Body != nil {
			dws.Body.Walk(v)
		}
		if dws.Test != nil {
			dws.Test.Walk(v)
		}
	}
}

func (dws DoWhileStatement) Errors() []error {
	return nil // TODO
}

func (dws DoWhileStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(dws)
	x["body"] = dws.Body
	x["test"] = dws.Test
	return json.Marshal(x)
}

func (dws *DoWhileStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string          `json:"type"`
		Loc  SourceLocation  `json:"loc"`
		Body json.RawMessage `json:"body"`
		Test json.RawMessage `json:"test"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != dws.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, dws.Type(), x.Type)
	}
	if err == nil {
		dws.Loc = x.Loc
		dws.Body, _, err = unmarshalStatement(x.Body)
		var err2 error
		if dws.Test, _, err = unmarshalExpression(x.Test); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

// ForStatement is a for loop.
type ForStatement struct {
	baseStatement
	Loc    SourceLocation
	Init   VariableDeclarationOrExpression
	Test   Expression
	Update Expression
	Body   Statement
}

func (ForStatement) Type() string                { return "ForStatement" }
func (fs ForStatement) Location() SourceLocation { return fs.Loc }

func (fs ForStatement) IsZero() bool {
	return fs.Loc.IsZero() &&
		(fs.Init == nil || fs.Init.IsZero()) &&
		(fs.Test == nil || fs.Test.IsZero()) &&
		(fs.Update == nil || fs.Update.IsZero()) &&
		(fs.Body == nil || fs.Body.IsZero())
}

func (fs ForStatement) Walk(v Visitor) {
	if v = v.Visit(fs); v != nil {
		defer v.Visit(nil)
		if fs.Init != nil {
			fs.Init.Walk(v)
		}
		if fs.Test != nil {
			fs.Test.Walk(v)
		}
		if fs.Update != nil {
			fs.Update.Walk(v)
		}
		if fs.Body != nil {
			fs.Body.Walk(v)
		}
	}
}

func (fs ForStatement) Errors() []error {
	return nil // TODO
}

func (fs ForStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(fs)
	if fs.Init != nil && !fs.Init.IsZero() {
		x["init"] = fs.Init
	}
	if fs.Test != nil && !fs.Test.IsZero() {
		x["test"] = fs.Test
	}
	if fs.Update != nil && !fs.Update.IsZero() {
		x["update"] = fs.Update
	}
	x["body"] = fs.Body
	return json.Marshal(x)
}

func (fs *ForStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type   string          `json:"type"`
		Loc    SourceLocation  `json:"loc"`
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
		fs.Loc = x.Loc
		if len(x.Init) > 0 {
			fs.Init, err = unmarshalVariableDeclarationOrExpression(x.Init)
		}
		var err2 error
		if len(x.Test) > 0 {
			if fs.Test, _, err = unmarshalExpression(x.Test); err == nil && err2 != nil {
				err = err2
			}
		}
		if len(x.Update) > 0 {
			if fs.Update, _, err = unmarshalExpression(x.Update); err == nil && err2 != nil {
				err = err2
			}
		}
		if fs.Body, _, err = unmarshalStatement(x.Body); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}

type ForInStatement struct {
	baseStatement
	Loc   SourceLocation
	Left  VariableDeclarationOrPattern
	Right Expression
	Body  Statement
}

func (ForInStatement) Type() string                 { return "ForInStatement" }
func (fis ForInStatement) Location() SourceLocation { return fis.Loc }

func (fis ForInStatement) IsZero() bool {
	return fis.Loc.IsZero() &&
		(fis.Left == nil || fis.Left.IsZero()) &&
		(fis.Right == nil || fis.Right.IsZero()) &&
		(fis.Body == nil || fis.Body.IsZero())
}

func (fis ForInStatement) Walk(v Visitor) {
	if v = v.Visit(fis); v != nil {
		defer v.Visit(nil)
		if fis.Left != nil {
			fis.Left.Walk(v)
		}
		if fis.Right != nil {
			fis.Right.Walk(v)
		}
		if fis.Body != nil {
			fis.Body.Walk(v)
		}
	}
}

func (fis ForInStatement) Errors() []error {
	return nil // TODO
}

func (fis ForInStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(fis)
	x["left"] = fis.Left
	x["right"] = fis.Right
	x["body"] = fis.Body
	return json.Marshal(x)
}

func (fis *ForInStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type  string          `json:"type"`
		Loc   SourceLocation  `json:"loc"`
		Left  json.RawMessage `json:"left"`
		Right json.RawMessage `json:"right"`
		Body  json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != fis.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, fis.Type(), x.Type)
	}
	if err == nil {
		fis.Loc = x.Loc
		fis.Left, err = unmarshalVariableDeclarationOrPattern(x.Left)
		var err2 error
		if fis.Right, _, err = unmarshalExpression(x.Right); err == nil && err2 != nil {
			err = err2
		}
		if fis.Body, _, err = unmarshalStatement(x.Body); err == nil && err2 != nil {
			err = err2
		}
	}
	return err
}
