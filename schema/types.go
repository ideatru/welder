package schema

type ElementType string

const (
	Number  = ElementType("number")
	String  = ElementType("string")
	Boolean = ElementType("bool")
	Array   = ElementType("array")
	Object  = ElementType("object")
)
