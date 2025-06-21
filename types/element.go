package types

// Elements is a slice of Element.
type Elements []Element

// Element represents a schema element with a type, optional name, nullability flag,
// and optional child elements for array and object types.
type Element struct {
	Name     string      `json:"name"`
	Type     ElementType `json:"type"`
	Size     int         `json:"size"`
	Children Elements    `json:"children"`
}

type Parser[T any] interface {
	Deserializer[T]
	Serializer[T]
}

type Deserializer[T any] interface {
	Deserialize(data T) (Elements, error)
}

type Serializer[T any] interface {
	Serialize(data Elements) (T, error)
}
