package ether

import (
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ideatru/welder/types"
)

var (
	// emptyTy represents an empty abi.Type.
	emptyTy = abi.Type{}

	// crossTypes maps internal `types.Type` to Ethereum `abi.Type` for encoding purposes.
	crossTypes = map[types.ElementType]abi.Type{
		types.String: {T: abi.StringTy},
		types.Number: {T: abi.IntTy, Size: 64},
		types.Bool:   {T: abi.BoolTy},
		types.Array:  {T: abi.SliceTy},
		types.Object: {T: abi.TupleTy},
	}
)
