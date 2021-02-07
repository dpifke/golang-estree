package estree

import (
	"encoding/json"
	"errors"
)

type Node interface {
	json.Marshaler

	Type() string
}

var (
	ErrWrongType  = errors.New("wrong type")
	ErrWrongValue = errors.New("unrecognized value")
)
