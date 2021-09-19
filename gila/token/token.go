package token

import (
	"fmt"
)

const (
	EOF Type = iota
	IDENT
	NUMBER
	PLUS
	MINUS
	PROGRAM
	BEGIN
	END
	LPAREN
	RPAREN
	SEMICOLON
	CONST
	VAR
	DOT
	EQUALS
	MULTIPLY
	COMA
	FUNCTION
	PROCEDURE
	FORWARD
	COLON
	ASSIGN
	INTEGER
	IF
	THEN
	ELSE
	MOD
	DIV
	NOTEQUALS
	LESS
	GREATER
	LESSEQ
	GREATEREQ
	WHILE
	DO
	BREAK
	EXIT
	FOR
	TO
	DOWNTO
	OR
	AND
	STRLIT
)

var tokens = []string{
	EOF:       "eof",
	IDENT:     "identifier",
	NUMBER:    "number",
	PLUS:      "+",
	MINUS:     "-",
	PROGRAM:   "program",
	BEGIN:     "begin",
	END:       "end",
	LPAREN:    "(",
	RPAREN:    ")",
	COLON:     ":",
	SEMICOLON: ";",
	CONST:     "const",
	VAR:       "var",
	DOT:       ".",
	EQUALS:    "=",
	LESS:      "<",
	GREATER:   ">",
	NOTEQUALS: "<>",
	LESSEQ:    "<=",
	GREATEREQ: ">=",
	MULTIPLY:  "*",
	COMA:      ",",
	ASSIGN:    ":=",
	FUNCTION:  "function",
	PROCEDURE: "procedure",
	FORWARD:   "forward",
	INTEGER:   "integer",
	IF:        "if",
	THEN:      "then",
	ELSE:      "else",
	DIV:       "div",
	MOD:       "mod",
	WHILE:     "while",
	DO:        "do",
	BREAK:     "break",
	EXIT:      "exit",
	FOR:       "for",
	TO:        "to",
	DOWNTO:    "downto",
	AND:       "and",
	OR:        "or",
}

var keywords = map[string]Type{
	tokens[PROGRAM]:   PROGRAM,
	tokens[BEGIN]:     BEGIN,
	tokens[END]:       END,
	tokens[CONST]:     CONST,
	tokens[VAR]:       VAR,
	tokens[FUNCTION]:  FUNCTION,
	tokens[PROCEDURE]: PROCEDURE,
	tokens[FORWARD]:   FORWARD,
	tokens[INTEGER]:   INTEGER,
	tokens[WHILE]:     WHILE,
	tokens[IF]:        IF,
	tokens[THEN]:      THEN,
	tokens[ELSE]:      ELSE,
	tokens[MOD]:       MOD,
	tokens[DIV]:       DIV,
	tokens[DO]:        DO,
	tokens[BREAK]:     BREAK,
	tokens[EXIT]:      EXIT,
	tokens[FOR]:       FOR,
	tokens[TO]:        TO,
	tokens[DOWNTO]:    DOWNTO,
	tokens[AND]:       AND,
	tokens[OR]:        OR,
}

type Type int

type Position struct {
	Line int
	Col  int
}

type Token struct {
	Kind     Type
	Value    string
	Position Position
}

func (t Token) String() string {
	return fmt.Sprintf("Token: [ %s | %s ] at line %d char %d", t.Value, tokens[t.Kind], t.Position.Line, t.Position.Col)
}

func (t Type) String() string {
	return tokens[t]
}

func IsKeywordOrIdent(value string) Type {
	if t, ok := keywords[value]; ok {
		return t
	}
	return IDENT
}
