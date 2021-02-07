package estree

// Function is a function declaration or expression.
type Function interface {
	Node
	FunctionID() Identifier // or nil
	FunctionParams() []Pattern
	FunctionBody() FunctionBody
}

func (fd FunctionDeclaration) FunctionID() Identifier     { return fd.ID }
func (fd FunctionDeclaration) FunctionParams() []Pattern  { return fd.Params }
func (fd FunctionDeclaration) FunctionBody() FunctionBody { return fd.Body }

func (fe FunctionExpression) FunctionID() Identifier     { return fe.ID }
func (fe FunctionExpression) FunctionParams() []Pattern  { return fe.Params }
func (fe FunctionExpression) FunctionBody() FunctionBody { return fe.Body }
