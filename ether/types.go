package ether

import (
	"fmt"
	"reflect"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ideatru/welder/internal/builder"
	"github.com/ideatru/welder/types"
)

var (
	// emptyTy represents an empty Ethereum ABI type
	// Used as a default return value for error cases in encoding/decoding operations
	emptyTy = abi.Type{}
)

var (
	// EtherBuilderOptions defines the configuration options for the Ethereum builder
	// Contains custom reflect functions and struct tag formatting for Ethereum types
	EtherBuilderOptions = builder.Option{
		ReflectReplacers:  EtherReflectFns,
		StructTagReplacer: EtherStructTag,
	}

	// EtherReflectFns maps element types to custom reflection functions
	// Provides special handling for Ethereum-specific types like Address and Bytes
	EtherReflectFns = map[types.ElementType]builder.ReflectFn{
		types.Address: ReflectAddressFn,
		types.Bytes:   ReflectBytesFn,
		types.Float:   ReflectFloatFn,
	}
)

// ReflectFloatFn returns an error for all float types, since the EVM doesn't support floats.
func ReflectFloatFn(elem types.Element) (reflect.Type, error) {
	return nil, fmt.Errorf("EVM compatibility does not support float types: %s", elem.Type)
}

// ReflectAddressFn provides the reflection type for Ethereum addresses
// Returns the reflect.Type for common.Address regardless of input element properties
func ReflectAddressFn(elem types.Element) (reflect.Type, error) {
	return reflect.TypeOf(common.Address{}), nil
}

// ReflectBytesFn provides the reflection type for byte arrays and slices
// Returns:
//   - hexutil.Bytes type for dynamic byte arrays (when Size is 0)
//   - Fixed-size array of bytes for fixed-size byte arrays (when Size > 0)
func ReflectBytesFn(elem types.Element) (reflect.Type, error) {
	if elem.Size == 0 {
		return reflect.TypeOf(hexutil.Bytes{}), nil
	}

	return reflect.ArrayOf(elem.Size, reflect.TypeOf(byte(0))), nil
}

// EtherStructTag creates a struct tag for Ethereum ABI and JSON serialization
// Takes a field name and generates a struct tag with both ABI and JSON tags using that name
// Format: `abi:"fieldName" json:"fieldName"`
func EtherStructTag(tag string) reflect.StructTag {
	return reflect.StructTag(`abi:"` + tag + `" json:"` + tag + `"`)
}
