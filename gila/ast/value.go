package ast

// Type of a literal or a variable
type Type int

const (
	INT Type = iota
	REAL
	STRING
	VOID
)

type Value interface {
	isValue()
	GetInt() int64
	GetFloat() float64
}

/*
	Wrappers for types, available inside Mila language.
	We create wrappers, so that they can implement a common interface.
	This interface will be used as a union type to store any of the following types.
*/
type MilaInt int

func (m MilaInt) GetInt() int64 {
	return int64(m)
}

func (m MilaInt) GetFloat() float64 {
	panic("Wrong type. Expected float, receiver an int.")
}

type MilaString string
type MilaReal float64

func (r MilaReal) GetInt() int64 {
	panic("Wrong type. Expected int, received float.")
}

func (r MilaReal) GetFloat() float64 {
	return float64(r)
}

func (m MilaInt) isValue()    {}
func (m MilaString) isValue() {}
func (m MilaReal) isValue()   {}
