package types

// ElementType represents the type of a schema element.
type ElementType string

const (
	// Int represents a number type.
	Int = ElementType("int")
	// Uint represents a unsigned-number type.
	Uint = ElementType("uint")
	// Float represents a float type.
	Float = ElementType("float")
	// String represents a string type.
	String = ElementType("string")
	// Bytes represents a bytes type.
	Bytes = ElementType("bytes")
	// Address represents an address type.
	Address = ElementType("address")
	// Bool represents a boolean type.
	Bool = ElementType("boolean")
	// Array represents an array type.
	Array = ElementType("array")
	// Object represents an object type.
	Object = ElementType("object")
)
