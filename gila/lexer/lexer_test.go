package lexer

import (
	"fmt"
	"gitlab.fit.cvut.cz/fedorgle/gila/token"
	"strings"
	"testing"
)

func TestLexerAdvance(t *testing.T) {
	l := New(strings.NewReader("ab"))
	if l.current != 'a' || l.next != 'b' {
		t.Error("Lexer doesn't read simple runes")
	}
	l.advance()
	if l.current != 'b' || l.next != rune(0) {
		t.Error("Lexer does not detect EOF.")
	}
	l.advance()
	if l.current != rune(0) || l.next != rune(0) {
		t.Error("Lexer does not detect EOF.")
	}
}

func TestLexer_NextToken(t *testing.T) {
	l := New(strings.NewReader("program kekes;228"))
	// Detects keywords and identifiers
	tok := l.NextToken()
	if tok.Kind != token.PROGRAM || tok.Value != "" {
		t.Error("Lexer fails to detect token PROGRAM!")
	}
	tok = l.NextToken()
	if tok.Kind != token.IDENT || tok.Value != "kekes" {
		t.Error("Lexer fails to detect an IDENTIFIER!")
	}
	tok = l.NextToken()
	if tok.Kind != token.SEMICOLON || tok.Value != "" {
		t.Error("Lexer fails to detect a SEMICOLON token!")
	}
	tok = l.NextToken()
	if tok.Kind != token.NUMBER || tok.Value != "228" {
		t.Error("Lexer fails to detect a NUMBER !")
	}
}

func Test_BasicSampleProgram(t *testing.T) {
	sampleProg := `
program factorial;

function facti(n : integer) : integer;
begin
    facti := 1;
    while n > 1 do
    begin
        facti := facti * n;
        dec(n);
    end;
end;    

function factr(n : integer) : integer;
begin
    if n = 1 then 
        factr := 1
    else
        factr := n * factr(n-1);
end;    

begin
    writeln(facti(5));
    writeln(factr(5));
end.
`
	l := New(strings.NewReader(sampleProg))
	for tok := l.NextToken(); tok.Kind != token.EOF; tok = l.NextToken() {
		fmt.Println(tok)
	}
}

func Test_MathExpr(t *testing.T) {
	l := New(strings.NewReader("-2 + 3 * (4 - 2)"))
	for tok := l.NextToken(); tok.Kind != token.EOF; tok = l.NextToken() {
		fmt.Println(tok)
	}
}
