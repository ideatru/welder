package ether

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ideatru/welder/types"
)

// Ensure EtherParser implements the Parser interface for AbiElements.
var _ types.Parser[AbiElements] = &EtherParser[AbiElements]{}

// AbiElements is a list of abi.Argument, typically used for encoding/decoding Ethereum ABI data.
type AbiElements []abi.Argument

// Encode encodes the values into a byte slice
func (e AbiElements) Encode(values ...any) ([]byte, error) {
	args := abi.Arguments(e)
	return args.Pack(values...)
}

// EtherParser is a generic converter that facilitates translation between `types.Elements`
// and a target ABI type `T` (which must be `AbiElements`).
type EtherParser[T AbiElements] struct{}

// NewEtherParser creates and returns a new instance of EtherParser.
func NewEtherParser[T AbiElements]() *EtherParser[T] {
	return &EtherParser[T]{}
}

// Serialize converts a structured `types.Elements` representation into the target ABI type `T`.
func (e *EtherParser[T]) Serialize(elements types.Elements) (T, error) {
	return e.encode(elements)
}

// encode internally processes and converts `types.Elements` into the `AbiElements` format.
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

// encodeArray converts a `types.Element` representing an array into an `abi.Type` for array.
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

// encodeObject converts a `types.Element` representing an object into an `abi.Type` for a tuple.
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

// Deserialize converts data of type `T` (AbiElements) into a structured `types.Elements` representation.
func (e *EtherParser[T]) Deserialize(data T) (types.Elements, error) {
	return e.decode(data)
}

// decode internally processes and converts `AbiElements` into the `types.Elements` format.
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
		case abi.TupleTy:
			el, err := e.decodeObject(&arg.Type)
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

// decodeArray recursively decodes an `abi.Type` representing an array into a `types.Element`.
func (e *EtherParser[T]) decodeArray(arg *abi.Type) (types.Element, error) {
	if arg.Elem == nil {
		return types.Element{}, fmt.Errorf("abi argument must have `Elem`")
	}

	el := types.Element{Type: types.Array, Children: make(types.Elements, 1)}
	switch arg.Elem.T {
	case abi.StringTy:
		el.Children[0] = types.Element{Type: types.String}
	case abi.IntTy:
		el.Children[0] = types.Element{Type: types.Number} // Fixed: was types.String
	case abi.BoolTy:
		el.Children[0] = types.Element{Type: types.Bool} // Fixed: was types.String
	case abi.SliceTy:
		childEl, err := e.decodeArray(arg.Elem)
		if err != nil {
			return types.Element{}, err
		}
		el.Children[0] = childEl
	case abi.TupleTy:
		childEl, err := e.decodeObject(arg.Elem)
		if err != nil {
			return types.Element{}, err
		}
		el.Children[0] = childEl
	default:
		return types.Element{}, fmt.Errorf("elements does not support type %q", arg.Elem.T) // Fixed: was arg.T
	}

	return el, nil
}

// decodeObject recursively decodes an `abi.Type` representing a tuple (object) into a `types.Element`.
func (e *EtherParser[T]) decodeObject(arg *abi.Type) (types.Element, error) {
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
		case abi.SliceTy:
			tmpEl, err := e.decodeArray(tupleEl)
			if err != nil {
				return types.Element{}, err
			}
			childEl.Type = tmpEl.Type
			childEl.Children = tmpEl.Children // Fixed: maintain consistency with other decode methods
		case abi.TupleTy:
			tmpEl, err := e.decodeObject(tupleEl)
			if err != nil {
				return types.Element{}, err
			}
			childEl.Type = tmpEl.Type
			childEl.Children = tmpEl.Children
		default:
			return types.Element{}, fmt.Errorf("elements does not support type %q", tupleEl.T)
		}

		el.Children = append(el.Children, childEl)
	}

	return el, nil
}
