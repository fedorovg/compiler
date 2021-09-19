package parser

import (
	"gitlab.fit.cvut.cz/fedorgle/gila/ast"
	"gitlab.fit.cvut.cz/fedorgle/gila/token"
	"strconv"
)

func keywordToType(t token.Token) ast.Type {
	switch t.Kind {
	case token.INTEGER:
		return ast.INT
		// TODO: Add more types
	default:
		panic("Trying to get type from an inappropriate token.")
	}
}

func tokenToOperation(t *token.Token) ast.Operation {
	switch t.Kind {
	case token.PLUS:
		return ast.PLUS
	case token.MINUS:
		return ast.MINUS
	case token.MULTIPLY:
		return ast.MULTIPLY
	case token.EQUALS:
		return ast.EQUALS
	case token.NOTEQUALS:
		return ast.NOTEQUALS
	case token.LESS:
		return ast.LESS
	case token.LESSEQ:
		return ast.LESSEQ
	case token.GREATER:
		return ast.GREATER
	case token.GREATEREQ:
		return ast.GREATEREQ
	case token.MOD:
		return ast.MOD
	case token.DIV:
		return ast.DIV
	case token.AND:
		return ast.AND
	case token.OR:
		return ast.OR
	default:
		panic("Trying to create an opeartion from an invalid Token")
	}
}

func numTokenValue(t token.Token) int {
	if t.Kind != token.NUMBER {
		panic("Token must be a numeric constant.")
	}
	switch t.Value[0] {
	case '&':
		if i, err := strconv.ParseInt(t.Value[1:], 8, 64); err == nil {
			return int(i)
		} else {
			panic("Could not get value of a number &.")
		}
	case '$':
		if i, err := strconv.ParseInt(t.Value[1:], 16, 64); err == nil {
			return int(i)
		} else {
			panic("Could not get value of a number $.")
		}
	default:
		if i, err := strconv.ParseInt(t.Value, 10, 64); err == nil {
			return int(i)
		} else {
			panic("Could not get value of a number.")
		}
	}
}
