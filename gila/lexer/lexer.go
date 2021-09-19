package lexer

import (
	"bufio"
	"fmt"
	"gitlab.fit.cvut.cz/fedorgle/gila/token"
	"io"
	"strconv"
	"strings"
	"unicode"
)

type Lexer struct {
	reader  *bufio.Reader
	current rune
	next    rune
	token.Position
}

func New(r io.Reader) *Lexer {
	l := Lexer{
		reader:  bufio.NewReader(r),
		current: rune(0),
		next:    rune(0),
	}
	// Prime the lexer, filling current and next runes
	l.advance()
	l.advance()
	l.Position = token.Position{
		Line: 0,
		Col:  0,
	}
	if l.current == rune(0) || l.next == rune(0) {
		panic("Given text is too small to contain any program!!")
	}
	return &l
}

func (l *Lexer) advance() {
	r, _, err := l.reader.ReadRune()
	if err != nil {
		if err == io.EOF {
			// Zero signifies the end of the buffer
			r = rune(0)
		}
	}
	if l.current == '\n' {
		l.Position.Line++
		l.Col = 0
	}
	l.Col++
	l.current = l.next
	l.next = r
}

func (l *Lexer) skipWhitespace() {
	for unicode.IsSpace(l.current) {
		l.advance()
	}
}

func (l *Lexer) NextToken() token.Token {
	l.skipWhitespace()
	switch l.current {
	case '(':
		l.advance()
		return token.Token{Kind: token.LPAREN, Position: l.Position}
	case ')':
		l.advance()
		return token.Token{Kind: token.RPAREN, Position: l.Position}
	case ';':
		l.advance()
		return token.Token{Kind: token.SEMICOLON, Position: l.Position}
	case rune(0):
		return token.Token{Kind: token.EOF, Position: l.Position}
	case '+':
		l.advance()
		return token.Token{Kind: token.PLUS, Position: l.Position}
	case '-':
		l.advance()
		return token.Token{Kind: token.MINUS, Position: l.Position}
	case '=':
		l.advance()
		return token.Token{Kind: token.EQUALS, Position: l.Position}
	case '*':
		l.advance()
		return token.Token{Kind: token.MULTIPLY, Position: l.Position}
	case '.':
		l.advance()
		return token.Token{Kind: token.DOT, Position: l.Position}
	case ',':
		l.advance()
		return token.Token{Kind: token.COMA, Position: l.Position}
	case ':':
		l.advance()
		if l.current == '=' {
			l.advance()
			return token.Token{Kind: token.ASSIGN, Position: l.Position}
		} else {
			return token.Token{Kind: token.COLON, Position: l.Position}
		}
	case '<':
		l.advance()
		if l.current == '=' {
			l.advance()
			return token.Token{Kind: token.LESSEQ, Position: l.Position}
		} else if l.current == '>' {
			l.advance()
			return token.Token{Kind: token.NOTEQUALS, Position: l.Position}
		} else {
			return token.Token{Kind: token.LESS, Position: l.Position}
		}
	case '>':
		l.advance()
		if l.current == '=' {
			l.advance()
			return token.Token{Kind: token.GREATEREQ, Position: l.Position}
		} else {
			return token.Token{Kind: token.GREATER, Position: l.Position}
		}
	case '\'':
		l.advance()
		var strLit strings.Builder
		for l.current != '\'' {
			strLit.WriteRune(l.current)
			l.advance()
		}
		l.advance()
		return token.Token{
			Kind:     token.STRLIT,
			Value:    strLit.String(),
			Position: l.Position,
		}
	default:
		if unicode.IsLetter(l.current) || l.current == '_' {
			return l.identifierOrKeyword()
		} else if unicode.IsDigit(l.current) || l.current == '&' || l.current == '$' {
			return l.numberLiteral()
		} else {
			panic("Encountered an invalid character.")
		}
	}
}
func (l *Lexer) identifierOrKeyword() token.Token {
	var val strings.Builder
	for unicode.IsLetter(l.current) || unicode.IsDigit(l.current) || l.current == '_' {
		val.WriteRune(l.current)
		l.advance()
	}
	lit := val.String()
	tt := token.IsKeywordOrIdent(lit)
	if tt == token.IDENT {
		// Not a keyword requires a literal desc
		pos := l.Position
		pos.Col -= len(lit)
		return token.Token{Kind: tt, Value: lit, Position: pos}
	} else {
		// Keywords' literals equal to themselves, no need to specify
		pos := l.Position
		pos.Col -= len(tt.String())
		return token.Token{Kind: tt, Position: pos}
	}
}

func (l *Lexer) numberLiteral() token.Token {
	var val strings.Builder
	var base = 10
	if l.current == '&' {
		base = 8
		l.advance()
	} else if l.current == '$' {
		base = 16
		l.advance()
	}
	for unicode.IsDigit(l.current) || unicode.IsLetter(l.current) {
		val.WriteRune(l.current)
		l.advance()
	}
	if base10, err := strconv.ParseUint(val.String(), base, 32); err == nil {
		return token.Token{Kind: token.NUMBER, Value: strconv.FormatUint(base10, 10), Position: l.Position}
	} else {
		errMsg := fmt.Sprintf("Invalid number literal %s in base %v", val.String(), base)
		panic(errMsg)
	}
}
