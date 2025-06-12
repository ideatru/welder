package bridge

// ElementType represents the type of a schema element.
type ElementType string

const (
	// Number represents a number type.
	Number = ElementType("number")
	// String represents a string type.
	String = ElementType("string")
	// Boolean represents a boolean type.
	Boolean = ElementType("boolean")
	// Array represents an array type.
	Array = ElementType("array")
	// Object represents an object type.
	Object = ElementType("object")
)
