package main

import (
	"fmt"
	"gitlab.fit.cvut.cz/fedorgle/gila/ir"
	"gitlab.fit.cvut.cz/fedorgle/gila/lexer"
	"gitlab.fit.cvut.cz/fedorgle/gila/parser"
	"os"
	"io/ioutil"
	"strings"
)

func main()  {
	// A hack to make lexer wotk with the test script
	bytes, _ := ioutil.ReadAll(os.Stdin)
	l := lexer.New(strings.NewReader(string(bytes)))
	p := parser.New(l)
	program := p.Parse()
	m := ir.NewModule(program)
	fmt.Println(m)
}
