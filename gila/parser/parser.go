package parser

import (
	"fmt"
	"gitlab.fit.cvut.cz/fedorgle/gila/ast"
	"gitlab.fit.cvut.cz/fedorgle/gila/lexer"
	"gitlab.fit.cvut.cz/fedorgle/gila/token"
)

type Parser struct {
	lexer   *lexer.Lexer
	current token.Token
	peek    token.Token
	context *ast.Function
}

func New(lexer *lexer.Lexer) *Parser {
	p := Parser{
		lexer:   lexer,
		current: token.Token{Kind: token.EOF},
		peek:    token.Token{Kind: token.EOF},
	}
	p.advance()
	p.advance() // Prime the parser
	if p.peek.Kind == token.EOF || p.current.Kind == token.EOF {
		panic("Programs is WAY to short.")
	}
	return &p
}

// Advance moves parser to the next token
func (p *Parser) advance() token.Token {
	c := p.current
	p.current = p.peek
	p.peek = p.lexer.NextToken()
	return c
}

// Moves parser to the next token, asserting, that current and expected match
func (p *Parser) match(expected token.Type) token.Token {
	if p.current.Kind == expected {
		tok := p.current
		p.advance()
		return tok
	} else {
		errorMessage := fmt.Sprintf("Expected: \n%v\nBut got:\n%v", expected, p.current)
		panic(errorMessage)
	}
}
