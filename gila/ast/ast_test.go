package ast

import (
	"fmt"
	"testing"
)

func Test_String(t *testing.T) {
	l := Literal{
		Value: 12,
	}
	fmt.Println(l)
	u := Unary{
		Operand:   l,
		Operation: MINUS,
	}
	fmt.Println(u)
	b := Binary{
		Left:      l,
		Right:     u,
		Operation: MULTIPLY,
	}
	fmt.Println(b)
}
