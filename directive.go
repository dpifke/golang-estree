package estree

import (
	"encoding/json"
	"fmt"
)

// Directive is a directive from the prologue of a script or function.
type Directive struct {
	Loc        SourceLocation
	Expression Literal

	// Directive is the raw string source of the directive without quotes.
	Directive string
}

func (Directive) Type() string               { return "Directive" }
func (d Directive) Location() SourceLocation { return d.Loc }
func (Directive) MinVersion() Version        { return ES5 }
func (Directive) isDirectiveOrStatement()    {}

func (d Directive) IsZero() bool {
	return d.Loc.IsZero() &&
		(d.Expression == nil || d.Expression.IsZero()) &&
		d.Directive == ""
}

func (d Directive) Walk(v Visitor) {
	if v = v.Visit(d); v != nil {
		defer v.Visit(nil)
		d.Expression.Walk(v)
	}
}

func (d Directive) Errors() []error {
	return nil // TODO
}

func (d Directive) MarshalJSON() ([]byte, error) {
	x := nodeToMap(d)
	x["expression"] = d.Expression
	x["directive"] = d.Directive
	return json.Marshal(x)
}

func (d *Directive) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Loc        SourceLocation  `json:"loc"`
		Expression json.RawMessage `json:"expression"`
		Directive  string          `json:"directive"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != d.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, d.Type(), x.Type)
	}
	if err == nil {
		d.Loc, d.Directive = x.Loc, x.Directive
		d.Expression, _, err = unmarshalLiteral(x.Expression)
	}
	return err
}
