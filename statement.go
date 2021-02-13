package estree

import (
	"encoding/json"
	"fmt"
)

// Statement is any statement.
type Statement interface {
	DirectiveOrStatement
	isStatement()
}

func unmarshalStatement(m json.RawMessage) (s Statement, match bool, err error) {
	// TODO: coerce null/empty m into EmptyStatement?
	var x struct {
		Type string `json:"type"`
	}
	if err = json.Unmarshal(m, &x); err != nil {
		match = true
	} else {
		switch x.Type {
		case ExpressionStatement{}.Type():
			var es ExpressionStatement
			s, match, err = es, true, json.Unmarshal(m, &es)
		case BlockStatement{}.Type():
			var bs BlockStatement
			s, match, err = bs, true, json.Unmarshal(m, &bs)
		case FunctionBody{}.Type():
			var fb FunctionBody
			s, match, err = fb, true, json.Unmarshal(m, &fb)
		case EmptyStatement{}.Type():
			var es EmptyStatement
			s, match, err = es, true, json.Unmarshal(m, &es)
		case DebuggerStatement{}.Type():
			var ds DebuggerStatement
			s, match, err = ds, true, json.Unmarshal(m, &ds)
		case WithStatement{}.Type():
			var ws WithStatement
			s, match, err = ws, true, json.Unmarshal(m, &ws)
		case ReturnStatement{}.Type():
			var rs ReturnStatement
			s, match, err = rs, true, json.Unmarshal(m, &rs)
		case LabeledStatement{}.Type():
			var ls LabeledStatement
			s, match, err = ls, true, json.Unmarshal(m, &ls)
		case BreakStatement{}.Type():
			var bs BreakStatement
			s, match, err = bs, true, json.Unmarshal(m, &bs)
		case ContinueStatement{}.Type():
			var cs ContinueStatement
			s, match, err = cs, true, json.Unmarshal(m, &cs)
		case IfStatement{}.Type():
			var is IfStatement
			s, match, err = is, true, json.Unmarshal(m, &is)
		case SwitchStatement{}.Type():
			var ss SwitchStatement
			s, match, err = ss, true, json.Unmarshal(m, &ss)
		case ThrowStatement{}.Type():
			var ts ThrowStatement
			s, match, err = ts, true, json.Unmarshal(m, &ts)
		case TryStatement{}.Type():
			var ts TryStatement
			s, match, err = ts, true, json.Unmarshal(m, &ts)
		case WhileStatement{}.Type():
			var ws WhileStatement
			s, match, err = ws, true, json.Unmarshal(m, &ws)
		case DoWhileStatement{}.Type():
			var dws DoWhileStatement
			s, match, err = dws, true, json.Unmarshal(m, &dws)
		case ForStatement{}.Type():
			var fs ForStatement
			s, match, err = fs, true, json.Unmarshal(m, &fs)
		case ForInStatement{}.Type():
			var fis ForInStatement
			s, match, err = fis, true, json.Unmarshal(m, &fis)
		default:
			err = fmt.Errorf("%w Statement, got %v", ErrWrongType, string(m))
		}
		if err != nil {
			s = nil // don't return incomplete objects
		}
	}
	return
}

type baseStatement struct{}

func (baseStatement) MinVersion() Version     { return ES5 }
func (baseStatement) isStatement()            {}
func (baseStatement) isDirectiveOrStatement() {}

// ExpressionStatement is a statement consisting of a single expression.
type ExpressionStatement struct {
	baseStatement
	Loc SourceLocation
	Expression
}

func (ExpressionStatement) Type() string                { return "ExpressionStatement" }
func (es ExpressionStatement) Location() SourceLocation { return es.Loc }

func (es ExpressionStatement) MinVersion() Version {
	return es.Expression.MinVersion()
}

func (es ExpressionStatement) IsZero() bool {
	return es.Loc.IsZero() && (es.Expression == nil || es.Expression.IsZero())
}

func (es ExpressionStatement) Walk(v Visitor) {
	if v = v.Visit(es); v != nil {
		defer v.Visit(nil)
		if es.Expression != nil {
			es.Expression.Walk(v)
		}
	}
}

func (es ExpressionStatement) Errors() []error {
	c := nodeChecker{Node: es}
	c.require(es.Expression, "expression")
	return c.errors()
}

func (es ExpressionStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(es)
	x["expression"] = es.Expression
	return json.Marshal(x)
}

func (es *ExpressionStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Loc        SourceLocation  `json:"loc"`
		Expression json.RawMessage `json:"expression"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != es.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, es.Type(), x.Type)
	}
	if err == nil {
		es.Loc = x.Loc
		es.Expression, _, err = unmarshalExpression(x.Expression)
	}
	return err
}

// BlockStatement is a block statement, i.e., a sequence of statements
// surrounded by braces.
type BlockStatement struct {
	baseStatement
	Loc  SourceLocation
	Body []Statement
}

func (BlockStatement) Type() string                { return "BlockStatement" }
func (bs BlockStatement) Location() SourceLocation { return bs.Loc }
func (BlockStatement) IsZero() bool                { return false }

func (bs BlockStatement) Walk(v Visitor) {
	if v = v.Visit(bs); v != nil {
		defer v.Visit(nil)
		for _, b := range bs.Body {
			b.Walk(v)
		}
	}
}

func (bs BlockStatement) Errors() []error {
	c := nodeChecker{Node: bs}
	c.requireEach(nodeSlice{
		Index: func(i int) Node { return bs.Body[i] },
		Len:   len(bs.Body),
	}, "statement")
	return c.errors()
}

func (bs BlockStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(bs)
	x["body"] = bs.Body
	return json.Marshal(x)
}

func (bs *BlockStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string            `json:"type"`
		Loc  SourceLocation    `json:"loc"`
		Body []json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != bs.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, bs.Type(), x.Type)
	}
	if err == nil {
		bs.Loc = x.Loc
		bs.Body = make([]Statement, len(x.Body))
		for i := range x.Body {
			var err2 error
			if bs.Body[i], _, err2 = unmarshalStatement(x.Body[i]); err == nil && err2 != nil {
				err = err2
			}
		}
	}
	return err
}

// FunctionBody is the body of a function, which is a block statement that may
// begin with directives.
type FunctionBody struct {
	baseStatement
	Loc  SourceLocation
	Body []DirectiveOrStatement
}

func (FunctionBody) Type() string                { return BlockStatement{}.Type() }
func (fb FunctionBody) Location() SourceLocation { return fb.Loc }
func (FunctionBody) IsZero() bool                { return false }

func (fb FunctionBody) Walk(v Visitor) {
	if v = v.Visit(fb); v != nil {
		defer v.Visit(nil)
		for _, b := range fb.Body {
			b.Walk(v)
		}
	}
}

func (fb FunctionBody) Errors() []error {
	c := nodeChecker{Node: fb}
	c.requireEach(nodeSlice{
		Index: func(i int) Node { return fb.Body[i] },
		Len:   len(fb.Body),
	}, "directive or statement")
	return c.errors()
}

func (fb FunctionBody) MarshalJSON() ([]byte, error) {
	x := nodeToMap(fb)
	x["body"] = fb.Body
	return json.Marshal(x)
}

func (fb *FunctionBody) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string            `json:"type"`
		Loc  SourceLocation    `json:"loc"`
		Body []json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != fb.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, fb.Type(), x.Type)
	}
	if err == nil {
		fb.Loc = x.Loc
		fb.Body = make([]DirectiveOrStatement, len(x.Body))
		for i := range x.Body {
			var err2 error
			if fb.Body[i], err2 = unmarshalDirectiveOrStatement(x.Body[i]); err == nil && err2 != nil {
				err = err2
			}
		}
	}
	return err
}

// EmptyStatement is an empty statement, i.e., a solitary semicolon.
type EmptyStatement struct {
	baseStatement
	Loc SourceLocation
}

func (EmptyStatement) Type() string                { return "EmptyStatement" }
func (es EmptyStatement) Location() SourceLocation { return es.Loc }
func (EmptyStatement) IsZero() bool                { return false }
func (EmptyStatement) Errors() []error             { return nil }

func (es EmptyStatement) Walk(v Visitor) {
	if v = v.Visit(es); v != nil {
		v.Visit(nil)
	}
}

func (es EmptyStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(nodeToMap(es))
}

func (es *EmptyStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string         `json:"type"`
		Loc  SourceLocation `json:"loc"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != es.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, es.Type(), x.Type)
	}
	if err == nil {
		es.Loc = x.Loc
	}
	return err
}

// DebuggerStatement is a debugger statement.
type DebuggerStatement struct {
	baseStatement
	Loc SourceLocation
}

func (DebuggerStatement) Type() string                { return "DebuggerStatement" }
func (ds DebuggerStatement) Location() SourceLocation { return ds.Loc }
func (DebuggerStatement) IsZero() bool                { return false }
func (DebuggerStatement) Errors() []error             { return nil }

func (ds DebuggerStatement) Walk(v Visitor) {
	if v = v.Visit(ds); v != nil {
		v.Visit(nil)
	}
}

func (ds DebuggerStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(nodeToMap(ds))
}

func (ds *DebuggerStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string         `json:"type"`
		Loc  SourceLocation `json:"loc"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ds.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, ds.Type(), x.Type)
	}
	if err == nil {
		ds.Loc = x.Loc
	}
	return err
}

// WithStatement is a with statement.
type WithStatement struct {
	baseStatement
	Loc    SourceLocation
	Object Expression
	Body   Statement
}

func (WithStatement) Type() string                { return "WithStatement" }
func (ws WithStatement) Location() SourceLocation { return ws.Loc }

func (ws WithStatement) IsZero() bool {
	return ws.Loc.IsZero() &&
		(ws.Object == nil || ws.Object.IsZero()) &&
		(ws.Body == nil || ws.Body.IsZero())
}

func (ws WithStatement) Walk(v Visitor) {
	if v = v.Visit(ws); v != nil {
		defer v.Visit(nil)
		if ws.Object != nil {
			ws.Object.Walk(v)
		}
		if ws.Body != nil {
			ws.Body.Walk(v)
		}
	}
}

func (ws WithStatement) Errors() []error {
	c := nodeChecker{Node: ws}
	c.require(ws.Object, "with expression")
	c.require(ws.Body, "with body")
	return c.errors()
}

func (ws WithStatement) MarshalJSON() ([]byte, error) {
	x := nodeToMap(ws)
	x["object"] = ws.Object
	x["body"] = ws.Body
	return json.Marshal(x)
}

func (ws *WithStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type   string          `json:"type"`
		Loc    SourceLocation  `json:"loc"`
		Object json.RawMessage `json:"object"`
		Body   json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ws.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, ws.Type(), x.Type)
	}
	if err == nil {
		ws.Loc = x.Loc
		ws.Object, _, err = unmarshalExpression(x.Object)
	}
	if err == nil {
		ws.Body, _, err = unmarshalStatement(x.Body)
	}
	return err
}
