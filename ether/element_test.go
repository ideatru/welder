package ether

import (
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ideatru/crosschain/types"
	"github.com/stretchr/testify/assert"
)

func TestAbiConverter(t *testing.T) {
	type Testcase struct {
		Name     string
		Input    types.Elements
		Expected AbiElements
	}

	testcases := []Testcase{
		{
			Name:     "number",
			Input:    types.Elements{{Type: types.Number}},
			Expected: AbiElements{{Type: crossTypes[types.Number]}},
		},
		{
			Name:     "string",
			Input:    types.Elements{{Type: types.String}},
			Expected: AbiElements{{Type: crossTypes[types.String]}},
		},
		{
			Name:     "bool",
			Input:    types.Elements{{Type: types.Bool}},
			Expected: AbiElements{{Type: crossTypes[types.Bool]}},
		},
		{
			Name:     "array",
			Input:    types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.String}}}},
			Expected: AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}}},
		},
	}

	converter := NewConverter()
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := converter.encode(tc.Input)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}
