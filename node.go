package estree

import (
	"bytes"
	"encoding/json"
	"fmt"
)

// Node represents an ESTree abstract syntax tree (AST) node.
type Node interface {
	json.Marshaler

	// Type is a string representing the AST variant type.
	Type() string

	// Location returns the SourceLocation of the Node.
	Location() SourceLocation

	// MinVersion reports the lowest version of the ECMAScript specification
	// required to fully express the syntax of this Node.
	//
	// Only the current Node is reported; to recurse into child Nodes, use
	// Walk.
	MinVersion() Version

	// IsZero indicates all fields of this Node are their uninitialized (zero
	// or nil) value, unless all fields are optional, in which case this
	// method always returns false.
	//
	// This method returns true if *any* required field is non-zero; to check
	// that *all* required fields are non-zero, use Errors.
	IsZero() bool

	// Walk performs a depth-first search of the AST using Visitor.  The
	// current Node is visited first, followed by each of its non-nil
	// children, in the order defined by the ESTree grammar.
	//
	// Uninitialized Nodes will be visited, but it is recommended that they be
	// treated the same as nil, but returning nil from Visitor.Visit if
	// Node.IsZero is true.
	Walk(Visitor)

	// Errors checks that required fields of this Node are non-nil and
	// non-zero, and that values such as operator tokens are valid.
	//
	// Only the current Node is checked; to recurse into child Nodes, use
	// Walk.
	Errors() []error
}

// nodeToMap returns a map containing a Node's Type and Loc.  Loc is omitted
// if it is the zero value.
//
// Maps are used instead of structs when marshaling to JSON, as encoding/json
// doesn't have a good way to omit zero values (omitempty works with nil
// interfaces, but not with zero-value structs).
func nodeToMap(n Node) map[string]interface{} {
	m := make(map[string]interface{}, 10)
	m["type"] = n.Type()
	if !n.Location().IsZero() {
		m["loc"] = n.Location()
	}
	return m
}

// SourceLocation contains the start and end positions of a Node.
type SourceLocation struct {
	// Source indicates the origin of the parsed source region, typically its
	// filename.
	Source string

	// Start is the position of the first character of the parsed source
	// region.
	Start Position

	// End is the position of the first character after the parsed source
	// region.
	End Position
}

// IsZero indicates sl contains no information about the source location.
func (sl SourceLocation) IsZero() bool {
	return sl.Source == "" && sl.Start.IsZero() && sl.End.IsZero()
}

func (sl SourceLocation) MarshalJSON() ([]byte, error) {
	if sl.IsZero() {
		return json.Marshal(nil)
	}
	x := map[string]interface{}{
		"source": sl.Source,
		"start":  sl.Start,
		"end":    sl.End,
	}
	if sl.Source == "" {
		x["source"] = nil
	}
	return json.Marshal(x)
}

// Position contains the line (1-indexed) and column (0-indexed) of a position
// in the parsed source region.
type Position struct {
	// Line indicates the line number.  The first line in a source region is
	// 1.
	Line int

	// Column indicates the column number.  The first column on a line is 0.
	Column int
}

// IsZero indicates p contains no information about the source location.
func (p Position) IsZero() bool {
	return p.Line == 0 && p.Column == 0
}

// Version represents a particular revision of the ECMAScript standard.
type Version int

var (
	ES5    Version = 5
	ES6    Version = 6
	ES2015 Version = ES6
	ES2016 Version = 7
	ES2017 Version = 8
	ES2018 Version = 9
	ES2019 Version = 10
	ES2020 Version = 11
	ES2021 Version = 12
)

func (v Version) String() string {
	switch v {
	case ES5:
		return "ES5"
	case ES6:
		return "ES6"
	case ES2016:
		return "ES2016"
	case ES2017:
		return "ES2017"
	case ES2018:
		return "ES2018"
	case ES2019:
		return "ES2019"
	case ES2020:
		return "ES2020"
	case ES2021:
		return "ES2021"
	}
	return fmt.Sprintf("%d", int(v))
}

// Vistor traverses the abstract syntax tree via Node.Walk.
type Visitor interface {
	// Visit is invoked for each node encountered by Node.Walk.  If a non-nil
	// Visitor is returned, Walk visits each of the children of the Node with
	// the new Visitor, followed by a call of Visit(nil).
	//
	// The final Visit(nil) call is via defer, so it's possible for a Visitor
	// to recover from a panic.
	Visit(Node) Visitor
}

// VisitorFunc allows an ordinary function to be used as a Visitor.
type VisitorFunc func(Node) Visitor

func (f VisitorFunc) Visit(n Node) Visitor {
	return f(n)
}

// isNullOrEmptyRawMessage is used when unmarshaling optional fields, to treat
// null and missing values as equivalent.
func isNullOrEmptyRawMessage(m json.RawMessage) bool {
	return len(m) == 0 || bytes.Equal(m, []byte("null"))
}
