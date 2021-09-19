package parser

import (
	"fmt"
	"gitlab.fit.cvut.cz/fedorgle/gila/ast"
	"gitlab.fit.cvut.cz/fedorgle/gila/token"
)

func (p *Parser) Parse() *ast.Program {
	p.match(token.PROGRAM)
	programName := p.match(token.IDENT).Value
	p.match(token.SEMICOLON)

	var functions []*ast.Function
	for p.current.Kind == token.FUNCTION || p.current.Kind == token.PROCEDURE {
		functions = append(functions, p.toplevelFunctionDeclaration())
	}
	if p.current.Kind != token.CONST && p.current.Kind != token.VAR && p.current.Kind != token.BEGIN {
		panic("This program lacks the main function.")
	}
	mainSignature := &ast.Signature{
		Name:   "main",
		Return: ast.VOID,
	}
	mainFunction := &ast.Function{
		Signature: mainSignature,
		Body:      nil,
		Variables: make(map[string]struct{}),
		Constants: make(map[string]ast.Literal),
	}
	p.context = mainFunction
	mainFunction.Body = p.functionBody(mainSignature)
	p.context = nil
	functions = append(functions, mainFunction)
	return &ast.Program{Name: programName, Functions: functions}
}

func (p *Parser) toplevelFunctionDeclaration() *ast.Function {
	signature := p.functionSignature()
	function := &ast.Function{
		Signature: signature,
		Body:      nil,
		Variables: make(map[string]struct{}),
		Constants: make(map[string]ast.Literal),
	}
	if p.current.Kind == token.FORWARD {
		p.advance()
		p.match(token.SEMICOLON)
	} else {
		p.context = function
		function.Body = p.functionBody(signature)
		p.context = nil
	}
	return function
}

func (p *Parser) functionSignature() *ast.Signature {
	if p.current.Kind != token.FUNCTION && p.current.Kind != token.PROCEDURE {
		errorString := fmt.Sprintf("Expected a function or procedure keyword on line %d in a function declaration", p.lexer.Position.Line)
		panic(errorString)
	}
	hasReturnType := p.current.Kind == token.FUNCTION
	p.advance()
	functionName := p.match(token.IDENT).Value
	p.match(token.LPAREN)
	signature := &ast.Signature{
		Name:       functionName,
		Return:     ast.VOID,
		Parameters: nil,
	}
	for p.current.Kind != token.RPAREN && p.current.Kind != token.EOF {
		parName := p.match(token.IDENT).Value
		p.match(token.COLON)
		p.match(token.INTEGER)
		v := ast.Variable{
			Name: parName,
		}
		signature.Parameters = append(signature.Parameters, v)
		if p.current.Kind == token.SEMICOLON {
			p.advance()
		}
	}
	p.match(token.RPAREN)
	if hasReturnType {
		p.match(token.COLON)
		returnType := keywordToType(p.match(token.INTEGER))
		signature.Return = returnType
	}
	p.match(token.SEMICOLON)
	return signature
}

func (p *Parser) functionBody(s *ast.Signature) *ast.Block {
	var statements []ast.Statement
	// Add return
	if s.Return != ast.VOID {
		p.context.Variables[s.Name] = struct{}{}
		//statements = append(statements, &ast.VariableDeclaration{
		//	Name: s.Name,
		//})
	}
	// Statements
	body := p.block()
	if p.current.Kind == token.DOT || p.current.Kind == token.SEMICOLON {
		p.advance()
	}
	body.Statements = append(statements, body.Statements...)
	return body
}

func (p *Parser) block() *ast.Block {
	var statements []ast.Statement
	if p.current.Kind == token.CONST {
		statements = append(statements, p.constantDeclarations()...)
	}

	// Parse variable declarations
	if p.current.Kind == token.VAR {
		statements = append(statements, p.variableDeclarations()...)
	}
	p.match(token.BEGIN)
	for p.current.Kind != token.END && p.current.Kind != token.EOF {
		statements = append(statements, p.statement())
		for p.current.Kind == token.SEMICOLON {
			p.advance()
		}
	}
	p.match(token.END)

	return &ast.Block{Statements: statements}
}
func (p *Parser) statement() ast.Statement {
	switch p.current.Kind {
	case token.IDENT:
		if p.peek.Kind == token.ASSIGN {
			return p.assignment()
		} else if p.peek.Kind == token.LPAREN {
			return p.procedureCall()
		}
	case token.IF:
		return p.ifStatement()
	case token.VAR:
		return p.block()
	case token.CONST:
		return p.block()
	case token.BEGIN:
		return p.block()
	case token.WHILE:
		return p.whileLoop()
	case token.FOR:
		return p.forLoop()
	case token.BREAK:
		p.advance()
		return &ast.Break{}
	case token.EXIT:
		p.advance()
		return &ast.Exit{}
	default:
		panic("Invalid statement")
	}
	panic("Some branch is not implemented!!!")
}

func (p *Parser) constantDeclarations() []ast.Statement {
	p.match(token.CONST)
	var declarations []ast.Statement
	for p.current.Kind != token.VAR && p.current.Kind != token.BEGIN && p.current.Kind != token.EOF {
		name := p.match(token.IDENT).Value
		p.match(token.EQUALS)
		value := p.number()
		p.context.Constants[name] = *value
		declarations = append(
			declarations,
			&ast.ConstantDeclaration{
				Name:    name,
				Literal: *value,
			},
		)
		p.match(token.SEMICOLON)
	}
	return declarations
}

func (p *Parser) variableDeclarations() []ast.Statement {
	p.match(token.VAR)
	var declarations []ast.Statement
	for p.current.Kind != token.VAR && p.current.Kind != token.BEGIN && p.current.Kind != token.EOF {
		var names []string
		for p.current.Kind != token.COLON && p.current.Kind != token.EOF {
			names = append(names, p.match(token.IDENT).Value)
			if p.current.Kind == token.COMA {
				p.advance()
			}
		}
		p.match(token.COLON)
		p.match(token.INTEGER)
		p.match(token.SEMICOLON)
		for _, name := range names {
			declarations = append(
				declarations,
				&ast.VariableDeclaration{
					Name: name,
				},
			)
			p.context.Variables[name] = struct{}{}
		}
	}
	return declarations
}
func (p *Parser) assignment() *ast.Assignment {
	variableName := p.match(token.IDENT).Value
	p.match(token.ASSIGN)
	value := p.expr()
	return &ast.Assignment{
		Variable: ast.Variable{Name: variableName},
		Value:    value,
	}
}

func (p *Parser) ifStatement() *ast.If {
	p.match(token.IF)
	condition := p.expr()
	p.match(token.THEN)
	thenBranch := p.statement()
	if p.current.Kind == token.DOT || p.current.Kind == token.SEMICOLON {
		p.advance()
	}
	i := &ast.If{
		Condition: condition,
		Then:      thenBranch,
		Else:      nil,
	}
	if p.current.Kind == token.ELSE {
		p.advance()
		i.Else = p.statement()
		if p.current.Kind == token.DOT || p.current.Kind == token.SEMICOLON {
			p.advance()
		}
	}
	return i
}

func (p *Parser) whileLoop() *ast.While {
	p.match(token.WHILE)
	condition := p.expr()
	p.match(token.DO)
	stmnt := p.statement()
	if p.current.Kind == token.DOT || p.current.Kind == token.SEMICOLON {
		p.advance()
	}
	return &ast.While{
		Condition: condition,
		Body:      stmnt,
	}
}

func (p *Parser) forLoop() *ast.For {
	p.match(token.FOR)
	variableName := p.match(token.IDENT).Value
	p.match(token.ASSIGN)
	value := p.expr()

	assignment := &ast.Assignment{
		Variable: ast.Variable{Name: variableName},
		Value:    value,
	}

	loopDirection := p.advance()
	if loopDirection.Kind != token.TO && loopDirection.Kind != token.DOWNTO {
		panic("Error in for loop header. Expected 'to' or 'downto' keywords.")
	}

	target := p.expr()
	p.match(token.DO)
	body := p.block()
	for p.current.Kind == token.SEMICOLON {
		p.advance()
	}

	return &ast.For{
		Initial: assignment,
		Upto:    loopDirection.Kind == token.TO,
		Target:  target,
		Body:    body,
	}
}

func (p *Parser) procedureCall() *ast.ProcedureCall {
	procedureName := p.match(token.IDENT).Value
	p.match(token.LPAREN)
	if p.current.Kind == token.STRLIT {
		t := p.advance()
		p.match(token.RPAREN)
		return &ast.ProcedureCall{
			Name: procedureName,
			Args: []ast.Expression{ast.StringLiteral{Value: t.Value}},
		}
	}
	var args []ast.Expression
	for p.current.Kind != token.RPAREN && p.current.Kind != token.EOF {
		args = append(args, p.expr())
		if p.current.Kind == token.COMA {
			p.advance()
		}
	}
	p.match(token.RPAREN)
	return &ast.ProcedureCall{
		Name: procedureName,
		Args: args,
	}
}
