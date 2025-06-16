package types

type Deserializer[T any] interface {
	Deserialize(data T) (Elements, error)
}

type Serializer[T any] interface {
	Serialize(data Elements) (T, error)
}
