package ether

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ideatru/welder/types"
)

var _ types.Parser[AbiElements] = NewEtherParser()

type AbiElements []abi.Argument

func (a AbiElements) Encode(values ...any) ([]byte, error) { return abi.Arguments(a).Pack(values...) }

type EtherParser[T AbiElements] struct{}

func NewEtherParser[T AbiElements]() *EtherParser[T] { return &EtherParser[T]{} }

func (e *EtherParser[T]) Serialize(elements types.Elements) (T, error) {
	return e.serialize(elements)
}

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

func (e *EtherParser[T]) encodeString(elem types.Element) (abi.Type, error) {
	if elem.Type != types.String {
		return emptyTy, fmt.Errorf("`encodeString` does not support type %q", elem.Type)
	}

	return abi.Type{T: abi.StringTy}, nil
}

func (e *EtherParser[T]) encodeBytes(elem types.Element) (abi.Type, error) {
	if elem.Type != types.Bytes {
		return emptyTy, fmt.Errorf("`encodeBytes` does not support type %q", elem.Type)
	}

	if elem.Size <= 0 {
		return abi.Type{T: abi.BytesTy}, nil
	}

	return abi.Type{T: abi.FixedBytesTy, Size: elem.Size}, nil
}

func (e *EtherParser[T]) encodeAddress(elem types.Element) (abi.Type, error) {
	if elem.Type != types.Address {
		return emptyTy, fmt.Errorf("`encodeAddress` does not support type %q", elem.Type)
	}

	return abi.Type{T: abi.AddressTy}, nil
}

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

func (e *EtherParser[T]) encodeBool(elem types.Element) (abi.Type, error) {
	if elem.Type != types.Bool {
		return emptyTy, fmt.Errorf("`encodeBool` does not support type %q", elem.Type)
	}

	return abi.Type{T: abi.BoolTy}, nil
}

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
			Name: abi.ToCamelCase(childElem.Name),
			Type: tupleElem.GetType(),
			Tag:  reflect.StructTag(`abi:"` + childElem.Name + `" json:"` + childElem.Name + `"`),
		})
	}

	ty.TupleType = reflect.StructOf(fields)
	return ty, nil
}

func (e *EtherParser[T]) Deserialize(args T) (types.Elements, error) {
	return e.deserialize(args)
}

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

func (e *EtherParser[T]) decodeString(ty abi.Type) (*types.Element, error) {
	if ty.T != abi.StringTy {
		return nil, fmt.Errorf("`decodeString` does not support type %q", ty.T)
	}

	return &types.Element{Type: types.String}, nil
}

func (e *EtherParser[T]) decodeBytes(ty abi.Type) (*types.Element, error) {
	if ty.T != abi.FixedBytesTy && ty.T != abi.BytesTy {
		return nil, fmt.Errorf("`decodeBytes` does not support type %q", ty.T)
	}

	return &types.Element{Type: types.Bytes, Size: ty.Size}, nil
}

func (e *EtherParser[T]) decodeAddress(ty abi.Type) (*types.Element, error) {
	if ty.T != abi.AddressTy {
		return nil, fmt.Errorf("`decodeAddress` does not support type %q", ty.T)
	}

	return &types.Element{Type: types.Address}, nil
}

func (e *EtherParser[T]) decodeNumber(ty abi.Type) (*types.Element, error) {
	switch ty.T {
	case abi.IntTy:
		return &types.Element{Type: types.Int, Size: ty.Size}, nil
	case abi.UintTy:
		return &types.Element{Type: types.Uint, Size: ty.Size}, nil
	}

	return nil, fmt.Errorf("`decodeNumber` does not support type %q", ty.T)
}

func (e *EtherParser[T]) decodeBool(ty abi.Type) (*types.Element, error) {
	if ty.T != abi.BoolTy {
		return nil, fmt.Errorf("`decodeBool` does not support type %q", ty.T)
	}

	return &types.Element{Type: types.Bool}, nil
}

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
