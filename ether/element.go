package ether

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ideatru/welder/types"
)

// EtherParser implements types.Deserializer and types.Serializer
var _ interface {
	types.Deserializer[AbiElements]
	types.Serializer[AbiElements]
} = &EtherParser[AbiElements]{}

// AbiElements is a list of abi.Argument
type AbiElements []abi.Argument

// Encode encodes the values into a byte slice
func (e AbiElements) Encode(values ...any) ([]byte, error) {
	args := abi.Arguments(e)
	return args.Pack(values...)
}

// EtherParser is a converter that converts types.Elements to AbiElements
type EtherParser[T AbiElements] struct{}

// NewEtherParser creates a new AbiConverter
func NewEtherParser[T AbiElements]() *EtherParser[T] {
	return &EtherParser[T]{}
}

// Encode encodes the elements into a byte slice
func (e *EtherParser[T]) Serialize(elements types.Elements) (T, error) {
	return e.encode(elements)
}

// encode encodes the elements into AbiElements
func (e *EtherParser[T]) encode(elements types.Elements) (T, error) {
	args := make(T, len(elements))

	for i, el := range elements {
		args[i].Name = el.Name
		switch el.Type {
		case types.String:
			args[i].Type = crossTypes[types.String]
		case types.Number:
			args[i].Type = crossTypes[types.Number]
		case types.Bool:
			args[i].Type = crossTypes[types.Bool]
		case types.Array:
			ty, err := e.encodeArray(el)
			if err != nil {
				return nil, err
			}
			args[i].Type = ty
		case types.Object:
			ty, err := e.encodeObject(el)
			if err != nil {
				return nil, err
			}
			args[i].Type = ty
		default:
			return nil, fmt.Errorf("parser does not suuprot %q", el.Type)
		}
	}

	return args, nil
}

// encodeArray encodes the array into an abi.Type
func (e *EtherParser[T]) encodeArray(el types.Element) (abi.Type, error) {
	if len(el.Children) != 1 {
		return emptyTy, fmt.Errorf("array must have one child")
	}

	ty := crossTypes[types.Array]
	switch el.Children[0].Type {
	case types.String:
		childTy := crossTypes[el.Children[0].Type]
		ty.Elem = &childTy
	case types.Number:
		childTy := crossTypes[el.Children[0].Type]
		ty.Elem = &childTy
	case types.Bool:
		childTy := crossTypes[el.Children[0].Type]
		ty.Elem = &childTy
	case types.Array:
		childTy, err := e.encodeArray(el.Children[0])
		if err != nil {
			return emptyTy, err
		}
		ty.Elem = &childTy
	case types.Object:
		childTy, err := e.encodeObject(el.Children[0])
		if err != nil {
			return emptyTy, err
		}
		ty.Elem = &childTy
	default:
		return emptyTy, fmt.Errorf("unsupported type: %q", el.Children[0].Type)
	}

	return ty, nil
}

// encodeObject encodes the object into an abi.Type
func (e *EtherParser[T]) encodeObject(el types.Element) (abi.Type, error) {
	if len(el.Children) == 0 {
		return emptyTy, fmt.Errorf("object must have at least one child")
	}

	ty := crossTypes[types.Object]
	fields := make([]reflect.StructField, 0)
	for _, child := range el.Children {
		ty.TupleRawNames = append(ty.TupleRawNames, child.Name)

		switch child.Type {
		case types.String:
			childTy := crossTypes[types.String]
			ty.TupleElems = append(ty.TupleElems, &childTy)
			fields = append(fields, reflect.StructField{
				Name: abi.ToCamelCase(child.Name),
				Type: childTy.GetType(),
				Tag:  reflect.StructTag(`abi:"` + child.Name + `" json:"` + child.Name + `"`),
			})
		case types.Number:
			childTy := crossTypes[types.Number]
			ty.TupleElems = append(ty.TupleElems, &childTy)
			fields = append(fields, reflect.StructField{
				Name: abi.ToCamelCase(child.Name),
				Type: childTy.GetType(),
				Tag:  reflect.StructTag(`abi:"` + child.Name + `" json:"` + child.Name + `"`),
			})
		case types.Bool:
			childTy := crossTypes[types.Bool]
			ty.TupleElems = append(ty.TupleElems, &childTy)
			fields = append(fields, reflect.StructField{
				Name: abi.ToCamelCase(child.Name),
				Type: childTy.GetType(),
				Tag:  reflect.StructTag(`abi:"` + child.Name + `" json:"` + child.Name + `"`),
			})
		case types.Array:
			childTy, err := e.encodeArray(child)
			if err != nil {
				return emptyTy, err
			}

			ty.TupleElems = append(ty.TupleElems, &childTy)
			fields = append(fields, reflect.StructField{
				Name: abi.ToCamelCase(child.Name),
				Type: childTy.GetType(),
				Tag:  reflect.StructTag(`abi:"` + child.Name + `" json:"` + child.Name + `"`),
			})
		case types.Object:
			childTy, err := e.encodeObject(child)
			if err != nil {
				return emptyTy, err
			}
			ty.TupleElems = append(ty.TupleElems, &childTy)
			fields = append(fields, reflect.StructField{
				Name: abi.ToCamelCase(child.Name),
				Type: childTy.GetType(),
				Tag:  reflect.StructTag(`abi:"` + child.Name + `" json:"` + child.Name + `"`),
			})
		default:
			return emptyTy, fmt.Errorf("unsupported type: %q", child.Type)
		}
	}

	ty.TupleType = reflect.StructOf(fields)
	return ty, nil
}

func (e *EtherParser[T]) Deserialize(data T) (types.Elements, error) {
	return e.decode(data)
}

func (e *EtherParser[T]) decode(data T) (types.Elements, error) {
	elements := make(types.Elements, len(data))

	for i, arg := range data {
		elements[i].Name = arg.Name
		switch arg.Type.T {
		case abi.StringTy:
			elements[i].Type = types.String
		case abi.IntTy:
			elements[i].Type = types.Number
		case abi.BoolTy:
			elements[i].Type = types.Bool
		case abi.SliceTy:
			el, err := e.decodeArray(&arg.Type)
			if err != nil {
				return nil, err
			}
			elements[i].Type = el.Type
			elements[i].Children = el.Children
		default:
			return nil, fmt.Errorf("elements does not support type %q", arg.Type.T)
		}
	}

	return elements, nil
}

func (e *EtherParser[T]) decodeArray(arg *abi.Type) (types.Element, error) {
	if arg.Elem == nil {
		return types.Element{}, fmt.Errorf("abi argument must have `Elem`")
	}

	el := types.Element{Type: types.Array, Children: make(types.Elements, 1)}
	switch arg.Elem.T {
	case abi.StringTy:
		el.Children[0] = types.Element{Type: types.String}
	case abi.IntTy:
		el.Children[0] = types.Element{Type: types.String}
	case abi.BoolTy:
		el.Children[0] = types.Element{Type: types.String}
	case abi.SliceTy:
		childEl, err := e.decodeArray(arg.Elem)
		if err != nil {
			return types.Element{}, err
		}
		el.Type = types.Array
		el.Children[0] = childEl
	default:
		return types.Element{}, fmt.Errorf("elements does not support type %q", arg.T)
	}

	return el, nil
}

func (e *EtherParser[T]) decodeObject(arg *abi.Type) (types.Element, error) {
	if arg.Elem == nil {
		return types.Element{}, fmt.Errorf("abi argument must have `Elem`")
	}

	if len(arg.TupleElems) != len(arg.TupleRawNames) {
		return types.Element{}, fmt.Errorf("abi argument must have `TupleElems` and `TupleRawNames`")
	}

	el := types.Element{Type: types.Object}
	for i, tupleEl := range arg.TupleElems {
		childEl := types.Element{
			Name: arg.TupleRawNames[i],
		}

		switch tupleEl.T {
		case abi.StringTy:
			childEl.Type = types.String
		case abi.IntTy:
			childEl.Type = types.Number
		case abi.BoolTy:
			childEl.Type = types.Bool
		default:
			return types.Element{}, fmt.Errorf("elements does not support type %q", tupleEl.T)
		}

		el.Children = append(el.Children, childEl)
	}

	return el, nil
}
