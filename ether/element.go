package ether

import (
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ideatru/crosschain/types"
)

type AbiElements []abi.Argument

type AbiConverter struct{}

func NewConverter() *AbiConverter {
	return &AbiConverter{}
}

func (e *AbiConverter) Encode(elements types.Elements) ([]byte, error) {
	args, err := e.encode(elements)
	if err != nil {
		return nil, err
	}

	return json.Marshal(args)
}

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

func (e *AbiConverter) encodeArray(el types.Element) (abi.Type, error) {
	if len(el.Children) != 1 {
		return emptyTy, fmt.Errorf("array must have one child")
	}

	ty := crossTypes[types.Array]
	switch el.Children[0].Type {
	case types.Array:
		childTy, err := e.encodeArray(el.Children[0])
		if err != nil {
			return emptyTy, err
		}

		ty.Elem = &childTy
	default:
		childTy := crossTypes[el.Children[0].Type]
		ty.Elem = &childTy
	}

	return ty, nil
}

func (e *AbiConverter) encodeObject(el types.Element) (abi.Type, error) {
	if len(el.Children) == 0 {
		return emptyTy, fmt.Errorf("object must have at least one child")
	}

	ty := crossTypes[types.Object]
	ty.TupleElems = make([]*abi.Type, 0)
	for _, child := range el.Children {
		_ = child
	}

	return ty, nil
}
