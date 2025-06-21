package ether

import (
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ideatru/welder/internal/builder"
	"github.com/ideatru/welder/types"
)

var (
	emptyTy = abi.Type{}
)

var (
	EtherBuilderOptions = builder.Option{
		ReflectReplacers:  EtherReflectFns,
		StructTagReplacer: EtherStructTag,
	}
	EtherReflectFns = map[types.ElementType]builder.ReflectFn{
		types.Address: ReflectAddressFn,
		types.Bytes:   ReflectBytesFn,
	}
)

func ReflectAddressFn(elem types.Element) (reflect.Type, error) {
	return reflect.TypeOf(common.Address{}), nil
}

func ReflectBytesFn(elem types.Element) (reflect.Type, error) {
	if elem.Size == 0 {
		return reflect.TypeOf(hexutil.Bytes{}), nil
	}

	return reflect.ArrayOf(elem.Size, reflect.TypeOf(byte(0))), nil
}

func EtherStructTag(tag string) reflect.StructTag {
	return reflect.StructTag(`abi:"` + tag + `" json:"` + tag + `"`)
}
