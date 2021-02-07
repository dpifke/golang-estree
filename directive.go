package estree

import (
	"encoding/json"
	"fmt"
)

// Directive is a directive from the prologue of a script or function.
type Directive struct {
	Expression Literal

	// Directive is the raw string source of the directive without quotes.
	Directive string
}

func (Directive) Type() string            { return "Directive" }
func (Directive) isDirectiveOrStatement() {}

func (d Directive) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":       d.Type(),
		"expression": d.Expression,
		"directive":  d.Directive,
	})
}

func (d *Directive) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Expression json.RawMessage `json:"expression"`
		Directive  string          `json:"directive"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != d.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, d.Type(), x.Type)
	}
	if err == nil {
		d.Expression, err = unmarshalLiteral(x.Expression)
	}
	if err == nil {
		d.Directive = x.Directive
	}
	return err
}
