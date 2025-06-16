package types

// Elements is a slice of Element.
type Elements []Element

// Element represents a schema element with a type, optional name, nullability flag,
// and optional child elements for array and object types.
type Element struct {
	Name     string      `json:"name"`
	Type     ElementType `json:"type"`
	Nullable bool        `json:"nullable"`
	Children Elements    `json:"children"`
}

type Deserializer[T any] interface {
	Deserialize(data T) (Elements, error)
}

type Serializer[T any] interface {
	Serialize(data Elements) (T, error)
}
