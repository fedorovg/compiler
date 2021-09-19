package ast

type Operation int

func (o Operation) String() string {
	switch o {
	case PLUS:
		return "+"
	case MINUS:
		return "-"
	case MULTIPLY:
		return "*"
	case DIV:
		return "DIV"
	case MOD:
		return "MOD"
	case EQUALS:
		return "="
	case NOTEQUALS:
		return "<>"
	case LESS:
		return "<"
	case LESSEQ:
		return "<="
	case GREATER:
		return ">"
	case GREATEREQ:
		return ">="
	case AND:
		return "and"
	case OR:
		return "or"

	default:
		panic("Invalid Operation value.")
	}
}

const (
	PLUS Operation = iota
	MINUS
	MULTIPLY
	DIV
	MOD
	EQUALS
	NOTEQUALS
	LESS
	LESSEQ
	GREATER
	GREATEREQ
	AND
	OR
)
