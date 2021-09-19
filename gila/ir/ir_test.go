package ir

import (
	"fmt"
	"gitlab.fit.cvut.cz/fedorgle/gila/lexer"
	"gitlab.fit.cvut.cz/fedorgle/gila/parser"
	"os"
	"strings"
	"testing"
)

//function gcdi(a: integer): integer;
//begin
//end;

func TestNewModule(t *testing.T) {
	input := `

program nestedBlocks;
var
	x: integer;
begin
	for x := 1 to 10 do
	begin
		writeln(x);
	end;
end.


`
	l := lexer.New(strings.NewReader(input))
	p := parser.New(l)
	fmt.Printf("Mila source:\n%s\n\nProduced LLVM IR:\n\n", input)
	program := p.Parse()
	m := NewModule(program)
	fmt.Println(m)
	m.DumpToFile("/Users/gleb/go/gila/llvm/lest.ir")
}

func Test_Final(t *testing.T) {
	baseDir := "/Users/gleb/go/gila/llvm/samples/"
	resDir := "/Users/gleb/go/gila/llvm/results/"
	testedFiles := []string{
		"consts.mila",
		"expressions.mila",
		"expressions2.mila",
		"factorial.mila",
		"factorialCycle.mila",
		"factorialRec.mila",
		"factorization.mila",
		"fibonacci.mila",
		"gcd.mila",
		"indirectrecursion.mila",
		"inputOutput.mila",
		"isprime.mila",
	}
	for _, filename := range testedFiles {
		input, err := os.Open(baseDir + filename)
		if err != nil {
			panic(err.Error())
		}
		l := lexer.New(input)
		p := parser.New(l)
		program := p.Parse()
		m := NewModule(program)
		m.DumpToFile(resDir + filename)
		fmt.Printf("%v is fine\n", filename)
	}
}
