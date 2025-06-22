package ether

import (
	"fmt"
	"github.com/ethereum/go-ethereum/crypto"
	"reflect"
	"strings"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ideatru/welder/internal/utils"
	"github.com/ideatru/welder/types"
)

// AbiElements is a slice of abi.Argument that implements encoder and decoder functionality
type AbiElements []abi.Argument

// Encode packs the provided values according to the ABI specification
// Returns the packed bytes or an error if packing fails
func (a AbiElements) Encode(values ...any) ([]byte, error) { return abi.Arguments(a).Pack(values...) }

// EncodeWithFunctionSignature encodes values with an optional Ethereum ABI function signature prepended to the data.
// If the function signature is empty, only the encoded values are returned.
// Returns the encoded byte array or an error if encoding fails.
func (a AbiElements) EncodeWithFunctionSignature(name string, values ...any) ([]byte, error) {
	data, err := a.Encode(values...)
	if err != nil {
		return nil, fmt.Errorf("failed to encode values: %w", err)
	}

	name = strings.TrimSpace(name)
	if name == "" {
		return data, nil
	}

	signature := crypto.Keccak256([]byte(name))[:4]
	return append(signature, data...), nil
}

// Decode unpacks the provided data according to the ABI specification
// Returns the unpacked values or an error if unpacking fails
func (a AbiElements) Decode(data []byte) ([]any, error) {
	return abi.Arguments(a).Unpack(data)
}

// EtherParser is responsible for converting between types.Elements and AbiElements
// It implements the types.Parser interface
type EtherParser[T AbiElements] struct{}

// NewEtherParser creates a new instance of EtherParser
func NewEtherParser[T AbiElements]() *EtherParser[T] { return &EtherParser[T]{} }

// Serialize converts types.Elements to AbiElements (T)
// Returns the serialized elements or an error if serialization fails
func (e *EtherParser[T]) Serialize(elements types.Elements) (T, error) {
	return e.serialize(elements)
}

// serialize is the internal implementation of Serialize
// Converts types.Elements to AbiElements (T)
func (e *EtherParser[T]) serialize(elements types.Elements) (T, error) {
	var (
		args = make(T, len(elements))
		err  error
	)

	for i, elem := range elements {
		args[i].Type, err = e.encode(elem)
		if err != nil {
			return nil, err
		}

		args[i].Name = elem.Name
	}

	return args, nil
}

// encode converts a types.Element to an abi.Type
// Dispatches to the appropriate type-specific encoder based on the element type
func (e *EtherParser[T]) encode(elem types.Element) (abi.Type, error) {
	var (
		ty  abi.Type
		err error
	)

	switch elem.Type {
	case types.String:
		ty, err = e.encodeString(elem)
	case types.Bytes:
		ty, err = e.encodeBytes(elem)
	case types.Address:
		ty, err = e.encodeAddress(elem)
	case types.Int, types.Uint:
		ty, err = e.encodeNumber(elem)
	case types.Bool:
		ty, err = e.encodeBool(elem)
	case types.Array:
		ty, err = e.encodeArray(elem)
	case types.Object:
		ty, err = e.encodeObject(elem)
	default:
		return emptyTy, fmt.Errorf("parser does not support %q", elem.Type)
	}

	if err != nil {
		return emptyTy, err
	}

	return ty, nil
}

// encodeString converts a string Element to an abi.Type
// Returns an error if the element is not of type String
func (e *EtherParser[T]) encodeString(elem types.Element) (abi.Type, error) {
	if elem.Type != types.String {
		return emptyTy, fmt.Errorf("`encodeString` does not support type %q", elem.Type)
	}

	return abi.Type{T: abi.StringTy}, nil
}

// encodeBytes converts a bytes Element to an abi.Type
// Handles both fixed and dynamic bytes based on element size
// Returns an error if the element is not of type Bytes
func (e *EtherParser[T]) encodeBytes(elem types.Element) (abi.Type, error) {
	if elem.Type != types.Bytes {
		return emptyTy, fmt.Errorf("`encodeBytes` does not support type %q", elem.Type)
	}

	if elem.Size <= 0 {
		return abi.Type{T: abi.BytesTy}, nil
	}

	return abi.Type{T: abi.FixedBytesTy, Size: elem.Size}, nil
}

// encodeAddress converts an address Element to an abi.Type
// Returns an error if the element is not of type Address
func (e *EtherParser[T]) encodeAddress(elem types.Element) (abi.Type, error) {
	if elem.Type != types.Address {
		return emptyTy, fmt.Errorf("`encodeAddress` does not support type %q", elem.Type)
	}

	return abi.Type{T: abi.AddressTy}, nil
}

// encodeNumber converts a numeric (Int or Uint) Element to an abi.Type
// Handles size specification, defaulting to 64 bits if not provided
// Returns an error if the element is not of a numeric type
func (e *EtherParser[T]) encodeNumber(elem types.Element) (abi.Type, error) {
	ty := abi.Type{}
	switch elem.Type {
	case types.Int:
		ty.T = abi.IntTy
	case types.Uint:
		ty.T = abi.UintTy
	default:
		return emptyTy, fmt.Errorf("`encodeNumber` does not support type %q", elem.Type)
	}

	if elem.Size <= 0 {
		ty.Size = 64
		return ty, nil
	}

	ty.Size = elem.Size
	return ty, nil
}

// encodeBool converts a boolean Element to an abi.Type
// Returns an error if the element is not of type Bool
func (e *EtherParser[T]) encodeBool(elem types.Element) (abi.Type, error) {
	if elem.Type != types.Bool {
		return emptyTy, fmt.Errorf("`encodeBool` does not support type %q", elem.Type)
	}

	return abi.Type{T: abi.BoolTy}, nil
}

// encodeArray converts an array Element to an abi.Type
// Handles both fixed arrays (with size) and dynamic slices
// Returns an error if the element does not have exactly one child
func (e *EtherParser[T]) encodeArray(elem types.Element) (abi.Type, error) {
	if len(elem.Children) != 1 {
		return emptyTy, fmt.Errorf("array must have one child")
	}

	ty := abi.Type{T: abi.SliceTy}
	if elem.Size > 0 {
		ty.T = abi.ArrayTy
		ty.Size = elem.Size
	}

	childTy, err := e.encode(elem.Children[0])
	if err != nil {
		return childTy, err
	}

	ty.Elem = &childTy
	return ty, nil
}

// encodeObject converts an object Element to an abi.Type (tuple)
// Creates a struct representation of the object with appropriate field types
// Returns an error if the element has no children
func (e *EtherParser[T]) encodeObject(elem types.Element) (abi.Type, error) {
	if len(elem.Children) == 0 {
		return emptyTy, fmt.Errorf("object must have at least one child")
	}

	ty := abi.Type{T: abi.TupleTy}
	fields := make([]reflect.StructField, 0)

	for _, childElem := range elem.Children {
		tupleElem, err := e.encode(childElem)
		if err != nil {
			return tupleElem, err
		}

		ty.TupleRawNames = append(ty.TupleRawNames, childElem.Name)
		ty.TupleElems = append(ty.TupleElems, &tupleElem)
		fields = append(fields, reflect.StructField{
			Name: utils.ToCamelCase(childElem.Name),
			Type: tupleElem.GetType(),
			Tag:  EtherStructTag(childElem.Name),
		})
	}

	ty.TupleType = reflect.StructOf(fields)
	return ty, nil
}

// Deserialize converts AbiElements (T) to types.Elements
// Returns the deserialized elements or an error if deserialization fails
func (e *EtherParser[T]) Deserialize(args T) (types.Elements, error) {
	return e.deserialize(args)
}

// deserialize is the internal implementation of Deserialize
// Converts AbiElements (T) to types.Elements
func (e *EtherParser[T]) deserialize(args T) (types.Elements, error) {
	elements := make(types.Elements, len(args))

	for i, arg := range args {
		elem, err := e.decode(arg.Type)
		if err != nil {
			return nil, err
		}

		elem.Name = arg.Name
		elements[i] = *elem
	}

	return elements, nil
}

// decode converts an abi.Type to a types.Element
// Dispatches to the appropriate type-specific decoder based on the ABI type
func (e *EtherParser[T]) decode(ty abi.Type) (*types.Element, error) {
	switch ty.T {
	case abi.StringTy:
		return e.decodeString(ty)
	case abi.BytesTy, abi.FixedBytesTy:
		return e.decodeBytes(ty)
	case abi.AddressTy:
		return e.decodeAddress(ty)
	case abi.IntTy, abi.UintTy:
		return e.decodeNumber(ty)
	case abi.BoolTy:
		return e.decodeBool(ty)
	case abi.SliceTy, abi.ArrayTy:
		return e.decodeArray(ty)
	case abi.TupleTy:
		return e.decodeObject(ty)
	}

	return nil, fmt.Errorf("elements does not support type %q", ty.T)
}

// decodeString converts an abi.Type of string to a types.Element
// Returns an error if the type is not StringTy
func (e *EtherParser[T]) decodeString(ty abi.Type) (*types.Element, error) {
	if ty.T != abi.StringTy {
		return nil, fmt.Errorf("`decodeString` does not support type %q", ty.T)
	}

	return &types.Element{Type: types.String}, nil
}

// decodeBytes converts an abi.Type of bytes to a types.Element
// Handles both fixed and dynamic bytes
// Returns an error if the type is not BytesTy or FixedBytesTy
func (e *EtherParser[T]) decodeBytes(ty abi.Type) (*types.Element, error) {
	if ty.T != abi.FixedBytesTy && ty.T != abi.BytesTy {
		return nil, fmt.Errorf("`decodeBytes` does not support type %q", ty.T)
	}

	return &types.Element{Type: types.Bytes, Size: ty.Size}, nil
}

// decodeAddress converts an abi.Type of address to a types.Element
// Returns an error if the type is not AddressTy
func (e *EtherParser[T]) decodeAddress(ty abi.Type) (*types.Element, error) {
	if ty.T != abi.AddressTy {
		return nil, fmt.Errorf("`decodeAddress` does not support type %q", ty.T)
	}

	return &types.Element{Type: types.Address}, nil
}

// decodeNumber converts an abi.Type of int/uint to a types.Element
// Handles both signed and unsigned integers with their size
// Returns an error if the type is not IntTy or UintTy
func (e *EtherParser[T]) decodeNumber(ty abi.Type) (*types.Element, error) {
	switch ty.T {
	case abi.IntTy:
		return &types.Element{Type: types.Int, Size: ty.Size}, nil
	case abi.UintTy:
		return &types.Element{Type: types.Uint, Size: ty.Size}, nil
	}

	return nil, fmt.Errorf("`decodeNumber` does not support type %q", ty.T)
}

// decodeBool converts an abi.Type of bool to a types.Element
// Returns an error if the type is not BoolTy
func (e *EtherParser[T]) decodeBool(ty abi.Type) (*types.Element, error) {
	if ty.T != abi.BoolTy {
		return nil, fmt.Errorf("`decodeBool` does not support type %q", ty.T)
	}

	return &types.Element{Type: types.Bool}, nil
}

// decodeArray converts an abi.Type of array/slice to a types.Element
// Handles both fixed arrays and dynamic slices
// Returns an error if the type is not ArrayTy or SliceTy, or if Elem is nil
func (e *EtherParser[T]) decodeArray(ty abi.Type) (*types.Element, error) {
	if ty.T != abi.SliceTy && ty.T != abi.ArrayTy {
		return nil, fmt.Errorf("`decodeArray` does not support type %q", ty.T)
	}

	if ty.Elem == nil {
		return nil, fmt.Errorf("`decodeArray` does not support `*abi.Type.Elem == nil`")
	}

	childEl, err := e.decode(*ty.Elem)
	if err != nil {
		return nil, err
	}

	elem := types.Element{Type: types.Array, Children: make(types.Elements, 1)}
	if ty.T == abi.ArrayTy {
		elem.Size = ty.Size
	}

	elem.Children[0] = *childEl
	return &elem, nil
}

// decodeObject converts an abi.Type of tuple to a types.Element
// Creates an object with all children corresponding to tuple elements
// Returns an error if the tuple has inconsistent structure
func (e *EtherParser[T]) decodeObject(ty abi.Type) (*types.Element, error) {
	if len(ty.TupleElems) != len(ty.TupleRawNames) {
		return nil, fmt.Errorf("`decodeObject` does not support invalid abi type")
	}

	elem := types.Element{Type: types.Object}
	for i, childTy := range ty.TupleElems {
		childElem, err := e.decode(*childTy)
		if err != nil {
			return nil, err
		}

		childElem.Name = ty.TupleRawNames[i]
		elem.Children = append(elem.Children, *childElem)
	}

	return &elem, nil
}
