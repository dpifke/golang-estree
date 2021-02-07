package estree

import (
	"encoding/json"
	"fmt"
)

type IfStatement struct {
	baseStatement
	Test       Expression
	Consequent Statement
	Alternate  Statement
}

func (IfStatement) Type() string { return "IfStatement" }

func (is IfStatement) MarshalJSON() ([]byte, error) {
	x := map[string]interface{}{
		"type":       is.Type(),
		"test":       is.Test,
		"consequent": is.Consequent,
	}
	if is.Alternate != nil {
		x["alternate"] = is.Alternate
	}
	return json.Marshal(x)
}

func (is *IfStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Test       json.RawMessage `json:"test"`
		Consequent json.RawMessage `json:"consequent"`
		Alternate  json.RawMessage `json:"alternate"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != is.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, is.Type(), x.Type)
	}
	if err == nil {
		is.Test, err = unmarshalExpression(x.Test)
	}
	if err == nil {
		is.Consequent, err = unmarshalStatement(x.Consequent)
	}
	if err == nil && len(x.Alternate) > 0 {
		is.Alternate, err = unmarshalStatement(x.Alternate)
	}
	return err
}

type SwitchStatement struct {
	baseStatement
	Discriminant Expression
	Cases        []SwitchCase
}

func (SwitchStatement) Type() string { return "SwitchStatement" }

func (ss SwitchStatement) MarshalJSON() ([]byte, error) {
	return json.Marshal(map[string]interface{}{
		"type":  ss.Type(),
		"test":  ss.Discriminant,
		"cases": ss.Cases,
	})
}

func (ss *SwitchStatement) UnmarshalJSON(b []byte) error {
	var x struct {
		Type         string          `json:"type"`
		Discriminant json.RawMessage `json:"test"`
		Cases        []SwitchCase    `json:"cases"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != ss.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, ss.Type(), x.Type)
	}
	if err == nil {
		ss.Discriminant, err = unmarshalExpression(x.Discriminant)
	}
	return err
}

type SwitchCase struct {
	Test       Expression
	Consequent Statement
}

func (SwitchCase) Type() string { return "SwitchCase" }

func (sc SwitchCase) MarshalJSON() ([]byte, error) {
	x := map[string]interface{}{
		"type": sc.Type(),
	}
	if sc.Test != nil {
		x["test"] = sc.Test
	}
	if sc.Consequent != nil {
		x["consequent"] = sc.Consequent
	}
	return json.Marshal(x)
}

func (sc *SwitchCase) UnmarshalJSON(b []byte) error {
	var x struct {
		Type       string          `json:"type"`
		Test       json.RawMessage `json:"test"`
		Consequent json.RawMessage `json:"consequent"`
	}
	err := json.Unmarshal(b, &x)
	if err == nil && x.Type != sc.Type() {
		err = fmt.Errorf("%w: expected %q, got %q", ErrWrongType, sc.Type(), x.Type)
	}
	if err == nil && len(x.Test) > 0 {
		sc.Test, err = unmarshalExpression(x.Test)
	}
	if err == nil && len(x.Consequent) > 0 {
		sc.Consequent, err = unmarshalStatement(x.Consequent)
	}
	return err
}
