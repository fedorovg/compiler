package parser

import (
	"gitlab.fit.cvut.cz/fedorgle/gila/ast"
	"gitlab.fit.cvut.cz/fedorgle/gila/token"
	"strconv"
)

func (p *Parser) parseExpression() ast.Expression {
	res := p.expr()
	return res
}

func (p *Parser) expr() ast.Expression {
	res := p.eqExpr()
	for p.current.Kind == token.AND || p.current.Kind == token.OR {
		op := tokenToOperation(&p.current)
		p.advance()
		res = &ast.Binary{
			Left:      res,
			Right:     p.eqExpr(),
			Operation: op,
		}
	}
	return res
}

func (p *Parser) eqExpr() ast.Expression {
	res := p.pmExpr()
	logicalOperators := map[token.Type]bool{
		token.EQUALS:    true,
		token.NOTEQUALS: true,
		token.GREATER:   true,
		token.GREATEREQ: true,
		token.LESS:      true,
		token.LESSEQ:    true,
	}
	for logicalOperators[p.current.Kind] {
		op := tokenToOperation(&p.current)
		p.advance()
		res = &ast.Binary{
			Left:      res,
			Right:     p.pmExpr(),
			Operation: op,
		}
	}
	return res
}

func (p *Parser) pmExpr() ast.Expression {
	res := p.term()
	for p.current.Kind == token.PLUS || p.current.Kind == token.MINUS {
		op := tokenToOperation(&p.current)
		p.advance()
		res = &ast.Binary{
			Left:      res,
			Right:     p.term(),
			Operation: op,
		}
	}
	return res
}

func (p *Parser) term() ast.Expression {
	res := p.factor()
	termOperators := map[token.Type]bool{
		token.MOD:      true,
		token.DIV:      true,
		token.MULTIPLY: true,
	}
	for termOperators[p.current.Kind] {
		op := tokenToOperation(&p.current)
		p.advance()
		res = &ast.Binary{Left: res, Right: p.factor(), Operation: op}
	}
	return res
}

func (p *Parser) factor() ast.Expression {
	switch p.current.Kind {
	case token.NUMBER:
		return p.number()
	case token.LPAREN:
		return p.parens()
	case token.MINUS:
		return p.unary()
	case token.PLUS:
		return p.unary()
	case token.IDENT:
		if p.peek.Kind != token.LPAREN { // Constant or variable
			var res ast.Expression
			if _, ok := p.context.Variables[p.current.Value]; ok {
				res = &ast.Variable{
					Name: p.current.Value,
				}
			}
			if literal, ok := p.context.Constants[p.current.Value]; ok {
				res = &literal
			}
			for _, par := range p.context.Signature.Parameters {
				if par.Name == p.current.Value {
					res = &par
					break
				}
			}
			p.advance()
			return res
		} else {
			//panic("Invalid factor after identifier.")
			return p.functionCall()
		}
	default:
		panic("Invalid factor.")
	}
}

func (p *Parser) functionCall() *ast.FunctionCall {
	procedureName := p.match(token.IDENT).Value
	p.match(token.LPAREN)
	var args []ast.Expression
	for p.current.Kind != token.RPAREN && p.current.Kind != token.EOF {
		args = append(args, p.expr())
		if p.current.Kind == token.COMA {
			p.advance()
		}
	}
	p.match(token.RPAREN)
	return &ast.FunctionCall{
		Name: procedureName,
		Args: args,
	}
}

func (p *Parser) parens() ast.Expression {
	p.advance()
	res := p.expr()
	if p.current.Kind != token.RPAREN {
		panic("Invalid syntax. Unbalanced parens")
	}
	p.advance()
	return res
}

func (p *Parser) number() *ast.Literal {
	i, err := strconv.Atoi(p.current.Value)
	if err != nil {
		panic("Invalid number token. Failed to convert.")
	}
	p.advance()
	return &ast.Literal{Value: int64(i)}
}

func (p *Parser) unary() *ast.Unary {
	op := tokenToOperation(&p.current)
	p.advance()
	return &ast.Unary{
		Operand:   p.factor(),
		Operation: op,
	}
}
