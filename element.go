package bridge

// Element represents a schema element with a type, optional name, nullability flag,
// and optional child elements for array and object types.
type Element struct {
	Name     string      `json:"name"`
	Ty       ElementType `json:"type"`
	Nullable bool        `json:"nullable"`
	Children []Element   `json:"children,omitempty"`
}
