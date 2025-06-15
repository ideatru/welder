package ether

import (
	"encoding/json"
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ideatru/welder/types"
)

// AbiElements is a list of abi.Argument
type AbiElements []abi.Argument

// Encode encodes the values into a byte slice
func (e AbiElements) Encode(values ...any) ([]byte, error) {
	args := abi.Arguments(e)
	return args.Pack(values...)
}

// AbiConverter is a converter that converts types.Elements to AbiElements
type AbiConverter struct{}

// NewConverter creates a new AbiConverter
func NewConverter() *AbiConverter {
	return &AbiConverter{}
}

// Encode encodes the elements into a byte slice
func (e *AbiConverter) Encode(elements types.Elements) ([]byte, error) {
	args, err := e.encode(elements)
	if err != nil {
		return nil, err
	}

	return json.Marshal(args)
}

// encode encodes the elements into AbiElements
func (e *AbiConverter) encode(elements types.Elements) (AbiElements, error) {
	args := make(AbiElements, len(elements))

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
			return nil, fmt.Errorf("abi converter does not suuprot %q", el.Type)
		}
	}

	return args, nil
}

// encodeArray encodes the array into an abi.Type
func (e *AbiConverter) encodeArray(el types.Element) (abi.Type, error) {
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
func (e *AbiConverter) encodeObject(el types.Element) (abi.Type, error) {
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
