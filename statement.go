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
	var x struct {
		Type string `json:"type"`
	}
	if err = json.Unmarshal(m, &x); err == nil {
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
			s, match = EmptyStatement{}, true
		case DebuggerStatement{}.Type():
			s, match = DebuggerStatement{}, true
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
			err = fmt.Errorf("%w: expected Statement, got %v", ErrWrongType, string(m))
		}
		if err != nil {
			s = nil // don't return incomplete objects
		}
	}
	return
}

type baseStatement struct{}

func (baseStatement) isStatement()            {}
func (baseStatement) isDirectiveOrStatement() {}

func (bl baseStatement) marshalEmptyJSON(typ string) ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": typ,
	})
}

// ExpressionStatement is an expression statement, i.e., a statement
// consisting of a single expression.
type ExpressionStatement struct {
	baseStatement
	Expression
}

func (ExpressionStatement) Type() string { return "ExpressionStatement" }

func (es ExpressionStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       es.Type(),
		"expression": es.Expression,
	})
}

func (es *ExpressionStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Expression json.RawMessage `json:"expression"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != es.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, es.Type(), x.Type)
	}
	if err == nil {
		es.Expression, _, err = unmarshalExpression(x.Expression)
	}
	return err
}

// BlockStatement is a block statement, i.e., a sequence of statements
// surrounded by braces.
type BlockStatement struct {
	baseStatement
	Body []Statement
}

func (BlockStatement) Type() string { return "BlockStatement" }

func (bs BlockStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       bs.Type(),
		"expression": bs.Body,
	})
}

func (bs *BlockStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string            `json:"type"`
		Body []json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != bs.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, bs.Type(), x.Type)
	}
	if err == nil {
		bs.Body = make([]Statement, len(x.Body))
		for i := range x.Body {
			if bs.Body[i], _, err = unmarshalStatement(x.Body[i]); err != nil {
				break
			}
		}
	}
	return err
}

// FunctionBody is the body of a function, which is a block statement that may
// begin with directives.
type FunctionBody struct {
	baseStatement
	Body []DirectiveOrStatement
}

func (FunctionBody) Type() string { return BlockStatement{}.Type() }

func (fb FunctionBody) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       fb.Type(),
		"expression": fb.Body,
	})
}

func (fb *FunctionBody) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string            `json:"type"`
		Body []json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != fb.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, fb.Type(), x.Type)
	}
	if err == nil {
		fb.Body = make([]DirectiveOrStatement, len(x.Body))
		for i := range x.Body {
			if fb.Body[i], err = unmarshalDirectiveOrStatement(x.Body[i]); err != nil {
				break
			}
		}
	}
	return err
}

// EmptyStatement is an empty statement, i.e., a solitary semicolon.
type EmptyStatement struct {
	baseStatement
}

func (EmptyStatement) Type() string { return "EmptyStatement" }

func (es EmptyStatement) MarshalJSON() ([]byte, error) {
	return es.marshalEmptyJSON(es.Type())
}

// DebuggerStatement is a debugger statement.
type DebuggerStatement struct {
	baseStatement
}

func (DebuggerStatement) Type() string { return "DebuggerStatement" }

func (ds DebuggerStatement) MarshalJSON() ([]byte, error) {
	return ds.marshalEmptyJSON(ds.Type())
}

// WithStatement is a with statement.
type WithStatement struct {
	baseStatement
	Object Expression
	Body   Statement
}

func (WithStatement) Type() string { return "WithStatement" }

func (ws WithStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":   ws.Type(),
		"object": ws.Object,
		"body":   ws.Body,
	})
}

func (ws *WithStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type   string          `json:"type"`
		Object json.RawMessage `json:"object"`
		Body   json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ws.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ws.Type(), x.Type)
	}
	if err == nil {
		ws.Object, _, err = unmarshalExpression(x.Object)
	}
	if err == nil {
		ws.Body, _, err = unmarshalStatement(x.Body)
	}
	return err
}
