package bridge

type Deserializer interface {
	Deserialize(data []byte) (any, error)
}

type Serializer interface {
	Serialize(data any) ([]byte, error)
}
