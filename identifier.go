package estree

import (
	"encoding/json"
	"fmt"
)

type Identifier struct {
	basePattern
	Loc  SourceLocation
	Name string
}

func (Identifier) Type() string               { return "Identifier" }
func (i Identifier) Location() SourceLocation { return i.Loc }
func (Identifier) isLiteralOrIdentifier()     {}

func (i Identifier) IsZero() bool {
	return i.Loc.IsZero() && i.Name == ""
}

func (i Identifier) Walk(v Visitor) {
	if v = v.Visit(i); v != nil {
		v.Visit(nil)
	}
}

func (i Identifier) Errors() []error {
	c := nodeChecker{Node: i}
	if i.Name == "" {
		c.appendf("%w empty identifier not allowed", ErrWrongValue)
	}
	// TODO: other validity checks?
	return c.errors()
}

func (i Identifier) MarshalJSON() ([]byte, error) {
	x := nodeToMap(i)
	x["name"] = i.Name
	return json.Marshal(x)
}

func (i *Identifier) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string         `json:"type"`
		Loc  SourceLocation `json:"loc"`
		Name string         `json:"name"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != i.Type() {
		err = fmt.Errorf("%w %s, got %q", ErrWrongType, i.Type(), x.Type)
	}
	if err == nil {
		i.Loc, i.Name = x.Loc, x.Name
	}
	return err
}
