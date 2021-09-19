package parser

import (
	"fmt"
	"gitlab.fit.cvut.cz/fedorgle/gila/ast"
	"gitlab.fit.cvut.cz/fedorgle/gila/lexer"
	"gitlab.fit.cvut.cz/fedorgle/gila/token"
	"strings"
	"testing"
)

func TestNew(t *testing.T) {
	l := lexer.New(strings.NewReader("-2 + 3 * (4 - 2)"))
	p := New(l)
	if p.current.Kind != token.MINUS || p.peek.Kind != token.NUMBER || p.peek.Value != "2" {
		t.Errorf("Could not create a parser. Expected tokens: %c, %v. Got: %v %v", '-', "2", p.current, p.peek)
	}
}

func TestParser_Parse(t *testing.T) {
	l := lexer.New(strings.NewReader("program test; begin writeln(2 + 3); end."))
	p := New(l)
	program := p.Parse()
	fmt.Println(program)
}
func Test_ParseSignature(t *testing.T) {
	l := lexer.New(strings.NewReader("function gcdi(a: integer; b: integer): integer;"))
	p := New(l)
	s := p.functionSignature()
	if s.Name != "gcdi" || len(s.Parameters) != 2 || s.Return != ast.INT {
		t.Error("Could not parse a basic function signature.")
	}
}

func Test_ParseProcedureCall(t *testing.T) {
	l := lexer.New(strings.NewReader("gcdi(1 + 3, 4 * (- 2));"))
	p := New(l)
	pc := p.procedureCall()
	if pc.Name != "gcid" || len(pc.Args) != 2 {
		t.Error("Can not parse a call")
	}
}

func Test_ParseVariableAssignment(t *testing.T) {
	l := lexer.New(strings.NewReader("x := -3 + 2 * a;"))
	p := New(l)
	a := p.assignment()
	fmt.Println(a)
	if a.Variable.Name != "x" {
		t.Error("Failed to parse assignment.")
	}
}

func Test_ConstantDeclarations(t *testing.T) {
	l := lexer.New(strings.NewReader("const A = 2; B = 12; begin"))
	p := New(l)
	cd := p.constantDeclarations()
	if len(cd) != 2 {
		t.Error("Failed to parse const declarations.")
	}
}

func Test_VariableDeclarations(t *testing.T) {
	l := lexer.New(strings.NewReader("var lol, kek: integer; x: integer; y: integer; begin"))
	p := New(l)
	vd := p.variableDeclarations()
	if len(vd) != 4 {
		t.Error("Failed to parse const declarations.")
	}
}

func Test_TrivialMath(t *testing.T) {
	l := lexer.New(strings.NewReader("(12 + 3) * 3 * (-2) - 5"))
	p := New(l)
	fmt.Println(p.parseExpression())
}

func Test_FunctionDeclaration(t *testing.T) {
	l := lexer.New(strings.NewReader("function testFunc(a: integer; b: integer): integer;\nvar tmp: integer;\nbegin\ntmp := a + b;\nif 0 then\nbegin\ngcdr := b;\nend;\nwriteln(b);\nend;"))
	p := New(l)
	f := p.toplevelFunctionDeclaration()
	fmt.Println(f)
}

func Test_ConstsInExpr(t *testing.T) {
	l := lexer.New(strings.NewReader("begin $1 mod 3 = 5 and 1 = 1 end."))
	p := New(l)
	p.match(token.BEGIN)
	ex := p.parseExpression()
	p.match(token.END)
	p.match(token.DOT)
	fmt.Println(ex)
	if p.current.Kind != token.EOF {
		t.Error("Did not parse the whole input.")
	}
}
