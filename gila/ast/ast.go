package ast

import (
	"fmt"
)

// ast definitions for Abstract Syntax Tree nodes.

// Node is a base type for every node of the ast.
type Node interface {
	isNode()
}

// Program holds an array of top level declarations (Functions / Procedures).
type Program struct {
	Name      string
	Functions []*Function
}

// Function is a top level declaration of a function.
type (
	Function struct {
		Signature *Signature
		Body      Statement
		Variables map[string]struct{}
		Constants map[string]Literal
	}

	// Signature is the type of a function
	Signature struct {
		Name       string
		Return     Type
		Parameters []Variable
	}

	FunctionBody struct {
		Body Statement
	}
)

type (
	// Statement is a base type for all language constructs, that return no value.
	Statement interface {
		Node
		isStatement()
	}

	// Block is a sequence of statements
	Block struct {
		Statements []Statement
	}

	// Assignment represents an assignment of a new value to a variable
	Assignment struct {
		Variable Variable
		Value    Expression
	}

	// If represents a conditional branching statement
	If struct {
		Condition Expression
		Then      Statement
		Else      Statement
	}
	While struct {
		Condition Expression
		Body      Statement
	}

	For struct {
		Initial *Assignment
		Upto    bool
		Target  Expression
		Body    *Block
	}

	// Break the current loop
	Break struct{}

	// Exit is a return statement
	Exit struct{}

	// ProcedureCall is a call to a function, that doesn't return a value.
	ProcedureCall struct {
		Name string
		Args []Expression
	}

	VariableDeclaration struct {
		Name string
	}

	ParameterDeclaration struct {
		Name string
	}

	ConstantDeclaration struct {
		Name    string
		Literal Literal
	}
)

type (
	// Expression is a node, that returns a Value of some sort
	Expression interface {
		Node
		isExpression()
	}

	// Literal is a literal value inside the source code.
	Literal struct {
		Value int64
	}

	StringLiteral struct {
		Value string
	}

	// Variable represents a symbol, referencing a value in a program.
	Variable struct {
		Name string
	}

	// Binary represents an operation with 2 operands.
	Binary struct {
		Left, Right Expression
		Operation   Operation
	}

	// Unary represents an operation with 1 operand.
	Unary struct {
		Operand   Expression
		Operation Operation
	}

	// FunctionCall represents a call to a function, that returns something.
	FunctionCall struct {
		Name string
		Args []Expression
	}
)

// Expressions' methods

func (_ Literal) isNode()       {}
func (_ Literal) isExpression() {}
func (l Literal) String() string {
	return fmt.Sprintf("%v", l.Value)
}

func (_ StringLiteral) isNode()       {}
func (_ StringLiteral) isExpression() {}
func (l StringLiteral) String() string {
	return fmt.Sprintf("%v", l.Value)
}

func (_ Variable) isNode()       {}
func (_ Variable) isExpression() {}
func (v Variable) String() string {
	return v.Name
}

func (_ Binary) isNode()       {}
func (_ Binary) isExpression() {}
func (b Binary) String() string {
	return fmt.Sprintf("(%v %v %v)", b.Left, b.Operation, b.Right)
}

func (_ Unary) isNode()       {}
func (_ Unary) isExpression() {}
func (u Unary) String() string {
	return fmt.Sprintf("(%v %v)", u.Operation, u.Operand)
}

func (_ FunctionCall) isNode()       {}
func (_ FunctionCall) isExpression() {}
func (f FunctionCall) String() string {
	return fmt.Sprintf("(%s(%v))", f.Name, f.Args)
}

func (_ Block) isNode()      {}
func (_ Block) isStatement() {}

func (_ Break) isNode()      {}
func (_ Break) isStatement() {}

func (_ Exit) isNode()      {}
func (_ Exit) isStatement() {}

// Statements' methods
func (_ Assignment) isNode()      {}
func (_ Assignment) isStatement() {}
func (a Assignment) String() string {
	return fmt.Sprintf("%v := %v", a.Variable.Name, a.Value)
}

func (_ If) isNode()      {}
func (_ If) isStatement() {}

func (_ While) isNode()      {}
func (_ While) isStatement() {}

func (_ For) isNode()      {}
func (_ For) isStatement() {}

func (_ ProcedureCall) isNode()      {}
func (_ ProcedureCall) isStatement() {}
func (p ProcedureCall) String() string {
	return fmt.Sprintf("%v%v", p.Name, p.Args)
}
func (_ VariableDeclaration) isNode()      {}
func (_ VariableDeclaration) isStatement() {}

func (_ ParameterDeclaration) isNode()      {}
func (_ ParameterDeclaration) isStatement() {}

func (_ ConstantDeclaration) isNode()      {}
func (_ ConstantDeclaration) isStatement() {}
