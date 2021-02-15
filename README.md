# pifke.org/estree

Go (Golang) implementation of the [ESTree Specification](https://github.com/estree/estree),
an Abstract Syntax Tree for working with Javascript code.

ESTree allows interoperability between software which parses, manipulates,
and/or executes Javascript code.  Examples include transpilers, minifiers, and
front-end frameworks.

This library provides type-checked objects representing AST nodes, which can
be serialized to and from JSON in the format used by compilers such as Acorn,
Babel, Typescript, etc.

## Example

The following Javascript:

```
function hello(name) {
  var greeting = "Hello";
  console.log(greeting + " " + name);
}
```

...can be represented as:

```
Program{
  Body: []DirectiveOrStatement{
    FunctionDeclaration{
      ID:     Identifier{Name: "hello"},
      Params: []Pattern{
        Identifier{Name: "name"},
      },
      Body:   FunctionBody{
        Body: []DirectiveOrStatement{
          VariableDeclaration{
            Declarations: []VariableDeclarator{
              VariableDeclarator{
                ID:   Identifier{Name: "greeting"},
                Init: StringLiteral{Value: "Hello"},
              },
            },
            Kind: Var,
          },
          ExpressionStatement{
            Expression: CallExpression{
              Callee: MemberExpression{
                Object:   Identifier{Name: "console"},
                Property: Identifier{Name: "log"},
              },
              Arguments: []Expression{
                BinaryExpression{
                  Left: BinaryExpression{
                    Left:     Identifier{Name: "greeting"},
                    Operator: Add,
                    Right:    StringLiteral{Value: " "},
                  },
                  Operator: Add,
                  Right:    Identifier{Name: "name"},
                },
              },
            },
          },
        },
      },
    },
  },
}
```

## Roadmap

Only the ES5 nodes have been implemented so far.  The remainder are coming
very soon.

I'll declare this library 1.0 once I have some real-world experience using it;
I built this with a particular project in mind, so this will hopefully also be
soon.

## Downloading

If you use this library in your own code, please use the canonical URL in your
Go code, instead of Github:

```
go get pifke.org/estree
```

As opposed to the pifke.org URL, I make no guarantee this Github repository
will exist or be up-to-date in the future.

## License

MIT (Expat), see LICENSE.txt.  If for some reason you'd like to use this code
and that doesn't work for you, please contact me.

The ESTree Specification itself is Copyright Mozilla Contributors and ESTree
Contributors, [Creative Commons Sharealike](https://creativecommons.org/licenses/by-sa/2.5/).

## Author

Dave Pifke.  My email address is my first name "at" my last name "dot org."

I'm [@dave@pifke.social](https://pifke.social/dave) in the Fediverse,
`@dave:pifke.chat` in Matrix, and `dpifke` on Freenode.  My PGP key is
available from [my web site](https://pifke.org/dpifke.asc).
