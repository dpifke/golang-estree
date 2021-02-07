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

func unmarshalStatement(m json.RawMessage) (Statement, error) {
	var x struct {
		Type string `json:"type"`
	}
	err := json.Unmarshal(m, &x)
	if err == nil {
		switch x.Type {
		case ExpressionStatement{}.Type():
			var es ExpressionStatement
			if err = json.Unmarshal(m, &es); err == nil {
				return es, nil
			}
		case BlockStatement{}.Type():
			var bs BlockStatement
			if err = json.Unmarshal(m, &bs); err == nil {
				return bs, nil
			}
		case FunctionBody{}.Type():
			var fb FunctionBody
			if err = json.Unmarshal(m, &fb); err == nil {
				return fb, nil
			}
		case EmptyStatement{}.Type():
			return EmptyStatement{}, nil
		case DebuggerStatement{}.Type():
			return DebuggerStatement{}, nil
		case WithStatement{}.Type():
			var ws WithStatement
			if err = json.Unmarshal(m, &ws); err == nil {
				return ws, nil
			}
		case ReturnStatement{}.Type():
			var rs ReturnStatement
			if err = json.Unmarshal(m, &rs); err == nil {
				return rs, nil
			}
		case LabeledStatement{}.Type():
			var ls LabeledStatement
			if err = json.Unmarshal(m, &ls); err == nil {
				return ls, nil
			}
		case BreakStatement{}.Type():
			var bs BreakStatement
			if err = json.Unmarshal(m, &bs); err == nil {
				return bs, nil
			}
		case ContinueStatement{}.Type():
			var cs ContinueStatement
			if err = json.Unmarshal(m, &cs); err == nil {
				return cs, nil
			}
		case IfStatement{}.Type():
			var is IfStatement
			if err = json.Unmarshal(m, &is); err == nil {
				return is, nil
			}
		case SwitchStatement{}.Type():
			var ss SwitchStatement
			if err = json.Unmarshal(m, &ss); err == nil {
				return ss, nil
			}
		case ThrowStatement{}.Type():
			var ts ThrowStatement
			if err = json.Unmarshal(m, &ts); err == nil {
				return ts, nil
			}
		case TryStatement{}.Type():
			var ts TryStatement
			if err = json.Unmarshal(m, &ts); err == nil {
				return ts, nil
			}
		case WhileStatement{}.Type():
			var ws WhileStatement
			if err = json.Unmarshal(m, &ws); err == nil {
				return ws, nil
			}
		case DoWhileStatement{}.Type():
			var dws DoWhileStatement
			if err = json.Unmarshal(m, &dws); err == nil {
				return dws, nil
			}
		case ForStatement{}.Type():
			var fs ForStatement
			if err = json.Unmarshal(m, &fs); err == nil {
				return fs, nil
			}
		case ForInStatement{}.Type():
			var fis ForInStatement
			if err = json.Unmarshal(m, &fis); err == nil {
				return fis, nil
			}
		default:
			err = fmt.Errorf("%w: expected Statement, got %v", ErrWrongType, string(m))
		}
	}
	return nil, err
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
		es.Expression, err = unmarshalExpression(x.Expression)
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
			if bs.Body[i], err = unmarshalStatement(x.Body[i]); err != nil {
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
		ws.Object, err = unmarshalExpression(x.Object)
	}
	if err == nil {
		ws.Body, err = unmarshalStatement(x.Body)
	}
	return err
}
