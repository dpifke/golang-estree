package estree

import (
	"encoding/json"
	"fmt"
)

type Identifier struct {
	basePattern
	Name string
}

func (Identifier) Type() string           { return "Identifier" }
func (Identifier) isLiteralOrIdentifier() {}

func (i Identifier) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type": i.Type(),
		"name": i.Name,
	})
}

func (i *Identifier) UnmarshalJSON(b []byte) error {
	var x struct {
		Type string `json:"type"`
		Name string `json:"name"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != i.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, i.Type(), x.Type)
	}
	return err
}
