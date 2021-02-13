package estree

import (
	"errors"
	"fmt"
	"strings"
)

var (
	// ErrWrongType is wrapped when unmarshaling a Node whose Type is not
	// allowed at the given location.  The wrapping error will contain
	// information about the Type encountered versus the expected Type(s).
	//
	// This method is only used by the JSON unmarshaler; when constructing the
	// AST from within Go code, the compiler is responsible for enforcing type
	// checks.
	ErrWrongType = errors.New("expected")

	// ErrWrongValue is wrapped when a Node's field is invalid, e.g. when a
	// operator derived from a string has an unknown value.  The wrapping
	// error will contain more information.
	ErrWrongValue = errors.New("unrecognized")

	// ErrMissingNode is wrapped when a required field is nil or zero.  The
	// wrapping error will contain more information.
	ErrMissingNode = errors.New("missing")
)

// SyntaxError wraps an error related to a Node.
type SyntaxError struct {
	// Err is the wrapped error.  It is required.
	Err error

	// Node is the invalid Node.  It is optional.
	Node Node

	// Position provides a more specific error location in the source.  It is
	// optional, and defaults to the start of Node.
	Position Position
}

func (err SyntaxError) Error() string {
	var s strings.Builder
	pos := err.Position
	if err.Node != nil {
		loc := err.Node.Location()
		if loc.Source != "" {
			s.WriteString(loc.Source)
			s.WriteRune(':')
		}
		if pos.Line == 0 {
			pos = loc.Start
		}
	}
	if pos.Line > 0 {
		fmt.Fprintf(&s, "%d:", pos.Line)
	}
	if pos.Column >= 0 {
		fmt.Fprintf(&s, "%d:", pos.Column)
	}
	if s.Len() > 0 {
		s.WriteRune(' ')
	}
	s.WriteString(err.Err.Error())
	return s.String()
}

func (err SyntaxError) Unwrap() error { return err.Err }

// nodeChecker provides a concise way to perform validation of a Node.
type nodeChecker struct {
	Node      Node
	checked   []SyntaxError
	locations []SourceLocation
}

// appendf appends a new SyntaxError at the current location.
func (c *nodeChecker) appendf(format string, args ...interface{}) {
	c.checked = append(c.checked, SyntaxError{
		Node: c.Node,
		Err:  fmt.Errorf(format, args...),
	})
	c.locations = append(c.locations, SourceLocation{})
}

// require appends a SyntaxError if n is nil or zero, otherwise it marks the
// current location.
func (c *nodeChecker) require(n Node, what string) {
	if n == nil || n.IsZero() {
		c.appendf("%w %s", ErrMissingNode, what)
	} else {
		c.locations = append(c.locations, n.Location())
	}
}

// nodeSlice provides a generic interface for accessing a slice of objects
// which implement Node.
type nodeSlice struct {
	Index func(i int) Node
	Len   int
}

// requireEach appends a SyntaxError for each element of a nodeSlice that is
// nil or zero.
func (c *nodeChecker) requireEach(ns nodeSlice, what string) {
	for i := 0; i < ns.Len; i++ {
		if n := ns.Index(i); n == nil || n.IsZero() {
			c.appendf("%w %s at index %d", ErrMissingNode, what, i)
		} else {
			c.locations = append(c.locations, n.Location())
		}
	}
}

// optional marks the location of Node.
func (c *nodeChecker) optional(n Node) {
	if n != nil {
		c.checked = append(c.checked, SyntaxError{})
		c.locations = append(c.locations, n.Location())
	}
}

// findPosition reports the location nearest the i'th error.
func (c *nodeChecker) findPosition(i int) Position {
	if i == 0 {
		// Report first error at NodeStart
		return Position{}
	}
	for j := i - 1; j >= 0; j-- {
		// Try to position at the end of the previously checked Node.
		if c.locations[i].End.Line > 0 {
			return c.locations[j].End
		}
	}
	for j := i + 1; i < len(c.locations); i++ {
		// If no previously checked Nodes, try to position at the start of the
		// next checked Node.
		if c.locations[i].Start.Line > 0 {
			return c.locations[j].Start
		}
	}
	return Position{}
}

// errors converts the nodeChecker results to an error slice.
func (c *nodeChecker) errors() []error {
	var errs []error
	for i, err := range c.checked {
		if err.Err != nil {
			err.Position = c.findPosition(i)
			if errs == nil {
				errs = make([]error, 0, len(c.checked)-i)
			}
			errs = append(errs, err)
		}
	}
	return errs
}
