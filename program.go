package estree

import (
	"encoding/json"
	"errors"
	"fmt"
)

// DirectiveOrStatement is either a Directive or a Statement.
type DirectiveOrStatement interface {
	Node
	isDirectiveOrStatement()
}

func unmarshalDirectiveOrStatement(m json.RawMessage) (DirectiveOrStatement, error) {
	if s, match, err := unmarshalStatement(m); match {
		return s, err
	}
	var d Directive
	if err := json.Unmarshal([]byte(m), &d); err == nil {
		return d, nil
	} else if !errors.Is(err, ErrWrongType) {
		return nil, err // don't return incomplete object
	}
	return nil, fmt.Errorf("%w Directive or Statement, got %v", ErrWrongType, string(m))
}

// Program is a complete program source tree.
type Program struct {
	Loc  SourceLocation
	Body []DirectiveOrStatement
}

func (Program) Type() string               { return "Program" }
func (p Program) Location() SourceLocation { return p.Loc }
func (Program) MinVersion() Version        { return ES5 }

func (p Program) IsZero() bool {
	return p.Loc.IsZero() && len(p.Body) == 0
}

func (p Program) Walk(v Visitor) {
	if v = v.Visit(p); v != nil {
		defer v.Visit(nil)
		for _, b := range p.Body {
			b.Walk(v)
		}
	}
}

func (p Program) Errors() []error {
	c := nodeChecker{Node: p}
	c.requireEach(nodeSlice{
		Index: func(i int) Node { return p.Body[i] },
		Len:   len(p.Body),
	}, "directive or statement")
	return c.errors()
}

func (p Program) MarshalJSON() ([]byte, error) {
	x := nodeToMap(p)
	if len(p.Body) > 0 {
		x["body"] = p.Body
	}
	return json.Marshal(x)
}

func (p *Program) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string            `json:"type"`
		Loc  SourceLocation    `json:"loc"`
		Body []json.RawMessage `json:"body"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != p.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, p.Type(), x.Type)
	}
	if err == nil {
		p.Loc = x.Loc
		if len(x.Body) == 0 {
			p.Body = nil
		} else {
			p.Body = make([]DirectiveOrStatement, len(x.Body))
			for i := range x.Body {
				var err2 error
				p.Body[i], err2 = unmarshalDirectiveOrStatement(x.Body[i])
				if err == nil && err2 != nil {
					err = err2
				}
			}
		}
	}
	return err
}
