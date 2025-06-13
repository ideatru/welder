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

// ElementEncoder is an interface that can encode an Element to a byte slice.
type ElementEncoder interface {
	Encode(data Elements) ([]byte, error)
}

// ElementDecoder is an interface that can decode a byte slice to an Element.
type ElementDecoder interface {
	Decode(data []byte) (Elements, error)
}
