// Package estree provides an Abstract Syntax Tree for working with Javascript
// code, based on the ESTree specification (https://github.com/estree/estree).
//
// ESTree allows interoperability between software which parses, manipulates,
// and/or executes Javascript code.  Examples include transpilers, minifiers, and
// front-end frameworks.
//
// AST nodes all implement the Node interface, which should be your starting
// point if you're looking for documentation.
//
// Versions Supported
//
// Right now, just ES5.  ES6 on up is coming soon.
//
// Comments
//
// See https://github.com/estree/estree/issues/201.  In short, there isn't a
// standard way to represent comments, so I haven't (yet) done so.
package estree
