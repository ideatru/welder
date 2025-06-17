package ether

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ideatru/welder/types"
	"github.com/stretchr/testify/assert"
)

var (
	abiStringTy = crossTypes[types.String]
	abiNumberTy = crossTypes[types.Number]
	abiBoolTy   = crossTypes[types.Bool]
)

func TestAbiParser_Serialize(t *testing.T) {
	type Testcase struct {
		Name     string
		Input    types.Elements
		Expected AbiElements
	}

	testcases := []Testcase{
		{
			Name:     "number",
			Input:    types.Elements{{Type: types.Number}},
			Expected: AbiElements{{Type: abiNumberTy}},
		},
		{
			Name:     "string",
			Input:    types.Elements{{Type: types.String}},
			Expected: AbiElements{{Type: abiStringTy}},
		},
		{
			Name:     "bool",
			Input:    types.Elements{{Type: types.Bool}},
			Expected: AbiElements{{Type: abiBoolTy}},
		},
		{
			Name: "number,string,bool",
			Input: types.Elements{
				{Type: types.Number},
				{Type: types.String},
				{Type: types.Bool},
			},
			Expected: AbiElements{
				{Type: abiNumberTy},
				{Type: abiStringTy},
				{Type: abiBoolTy},
			},
		},
		{
			Name:     "string[]",
			Input:    types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.String}}}},
			Expected: AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}}},
		},
		{
			Name:     "string[][]",
			Input:    types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.String}}}}}},
			Expected: AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}}}},
		},
		{
			Name:     "(string name, number amount, bool valid)",
			Input:    types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "amount", Type: types.Number}, {Name: "valid", Type: types.Bool}}}},
			Expected: AbiElements{{Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "amount", "valid"}, TupleElems: []*abi.Type{&abiStringTy, &abiNumberTy, &abiBoolTy}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Name", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"name" json:"name"`)}, {Name: "Amount", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"amount" json:"amount"`)}, {Name: "Valid", Type: abiBoolTy.GetType(), Tag: reflect.StructTag(`abi:"valid" json:"valid"`)}})}}},
		},
		{
			Name:     "(string[] names)",
			Input:    types.Elements{{Type: types.Object, Children: types.Elements{{Name: "names", Type: types.Array, Children: types.Elements{{Type: types.String}}}}}},
			Expected: AbiElements{{Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"names"}, TupleElems: []*abi.Type{{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Names", Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}.GetType(), Tag: reflect.StructTag(`abi:"names" json:"names"`)}})}}},
		},
		{
			Name:     "(string[][] names)",
			Input:    types.Elements{{Type: types.Object, Children: types.Elements{{Name: "names", Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.String}}}}}}}},
			Expected: AbiElements{{Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"names"}, TupleElems: []*abi.Type{{T: abi.SliceTy, Elem: &abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Names", Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}}.GetType(), Tag: reflect.StructTag(`abi:"names" json:"names"`)}})}}},
		},
		{
			Name:     "(string name, number age)[]",
			Input:    types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "age", Type: types.Number}}}}}},
			Expected: AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "age"}, TupleElems: []*abi.Type{&abiStringTy, &abiNumberTy}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Name", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"name" json:"name"`)}, {Name: "Age", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"age" json:"age"`)}})}}}},
		},
		{
			Name:     "(string name, (number amount, bool valid) detail)",
			Input:    types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "detail", Type: types.Object, Children: types.Elements{{Name: "amount", Type: types.Number}, {Name: "valid", Type: types.Bool}}}}}},
			Expected: AbiElements{{Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "detail"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.TupleTy, TupleRawNames: []string{"amount", "valid"}, TupleElems: []*abi.Type{&abiNumberTy, &abiBoolTy}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Amount", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"amount" json:"amount"`)}, {Name: "Valid", Type: abiBoolTy.GetType(), Tag: reflect.StructTag(`abi:"valid" json:"valid"`)}})}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Name", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"name" json:"name"`)}, {Name: "Detail", Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"amount", "valid"}, TupleElems: []*abi.Type{&abiNumberTy, &abiBoolTy}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Amount", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"amount" json:"amount"`)}, {Name: "Valid", Type: abiBoolTy.GetType(), Tag: reflect.StructTag(`abi:"valid" json:"valid"`)}})}.GetType(), Tag: reflect.StructTag(`abi:"detail" json:"detail"`)}})}}},
		},
		{
			Name:     "(string name, ((string instrument, number[] prices) base, (string instrument, number[] prices) quote)[] pair)[][]",
			Input:    types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "pair", Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "base", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Number}}}}}, {Name: "quote", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Number}}}}}}}}}}}}}}}},
			Expected: AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "pair"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"base", "quote"}, TupleElems: []*abi.Type{{T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}, {T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Base", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"base" json:"base"`)}, {Name: "Quote", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"quote" json:"quote"`)}})}}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Name", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"name" json:"name"`)}, {Name: "Pair", Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"base", "quote"}, TupleElems: []*abi.Type{{T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}, {T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Base", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"base" json:"base"`)}, {Name: "Quote", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"quote" json:"quote"`)}})}}.GetType(), Tag: reflect.StructTag(`abi:"pair" json:"pair"`)}})}}}}},
		},
	}

	parser := NewEtherParser()
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			result, err := parser.Serialize(tc.Input)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, result)
		})
	}
}

func TestAbiElements_Encode(t *testing.T) {
	type Testcase struct {
		Name     string
		Input    types.Elements
		Args     []any
		Expected []byte
	}

	testcases := []Testcase{
		{
			Name:     "number",
			Input:    types.Elements{{Type: types.Number}},
			Args:     []any{int64(1)},
			Expected: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			Name:     "string",
			Input:    types.Elements{{Type: types.String}},
			Args:     []any{"lorem ipsum"},
			Expected: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b6c6f72656d20697073756d000000000000000000000000000000000000000000"),
		},
		{
			Name:     "bool",
			Input:    types.Elements{{Type: types.Bool}},
			Args:     []any{true},
			Expected: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			Name:     "string[]",
			Input:    types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.String}}}},
			Args:     []any{[]string{"lorem ipsum"}},
			Expected: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000b6c6f72656d20697073756d000000000000000000000000000000000000000000"),
		},
		{
			Name:     "int64[][]",
			Input:    types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Number}}}}}},
			Args:     []any{[][]int64{{1, 2}, {3, 4}}},
			Expected: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000004"),
		},
		{
			Name:  "(string instrument, number price, bool valid)",
			Input: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "price", Type: types.Number}, {Name: "valid", Type: types.Bool}}}},
			Args: []any{struct {
				Instrument string `abi:"instrument" json:"instrument"`
				Price      int64  `abi:"price" json:"price"`
				Valid      bool   `abi:"valid" json:"valid"`
			}{Instrument: "BTC", Price: 1000000000000000000, Valid: true}},
			Expected: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000034254430000000000000000000000000000000000000000000000000000000000"),
		},
		{
			Name:  "(string instrument, number price, bool valid)[]",
			Input: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "price", Type: types.Number}, {Name: "valid", Type: types.Bool}}}}}},
			Args: []any{[]struct {
				Instrument string `abi:"instrument" json:"instrument"`
				Price      int64  `abi:"price" json:"price"`
				Valid      bool   `abi:"valid" json:"valid"`
			}{{Instrument: "BTC", Price: 1000000000000000000, Valid: true}, {Instrument: "ETH", Price: 1000000000000000000, Valid: true}}},
			Expected: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000e000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000de0b6b3a764000000000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000003425443000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000600000000000000000000000000000000000000000000000000de0b6b3a7640000000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000034554480000000000000000000000000000000000000000000000000000000000"),
		},
		{
			Name:  "(string name, ((string instrument, number[] prices) base, (string instrument, number[] prices) quote)[] pair)[][]",
			Input: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "pair", Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "base", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Number}}}}}, {Name: "quote", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Number}}}}}}}}}}}}}}}},
			Args: []any{[][]struct {
				Name string `abi:"name" json:"name"`
				Pair []struct {
					Base struct {
						Instrument string  `abi:"instrument" json:"instrument"`
						Prices     []int64 `abi:"prices" json:"prices"`
					} `abi:"base" json:"base"`
					Quote struct {
						Instrument string  `abi:"instrument" json:"instrument"`
						Prices     []int64 `abi:"prices" json:"prices"`
					} `abi:"quote" json:"quote"`
				} `abi:"pair" json:"pair"`
			}{
				{
					{
						Name: "CryptoExchange1",
						Pair: []struct {
							Base struct {
								Instrument string  `abi:"instrument" json:"instrument"`
								Prices     []int64 `abi:"prices" json:"prices"`
							} `abi:"base" json:"base"`
							Quote struct {
								Instrument string  `abi:"instrument" json:"instrument"`
								Prices     []int64 `abi:"prices" json:"prices"`
							} `abi:"quote" json:"quote"`
						}{
							{
								Base: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "BTC",
									Prices:     []int64{68000, 68100, 68050},
								},
								Quote: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "USD",
									Prices:     []int64{1, 1, 1},
								},
							},
							{
								Base: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "ETH",
									Prices:     []int64{3500, 3510, 3505},
								},
								Quote: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "USD",
									Prices:     []int64{1, 1, 1},
								},
							},
						},
					},
					{
						Name: "CryptoExchange2",
						Pair: []struct {
							Base struct {
								Instrument string  `abi:"instrument" json:"instrument"`
								Prices     []int64 `abi:"prices" json:"prices"`
							} `abi:"base" json:"base"`
							Quote struct {
								Instrument string  `abi:"instrument" json:"instrument"`
								Prices     []int64 `abi:"prices" json:"prices"`
							} `abi:"quote" json:"quote"`
						}{
							{
								Base: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "SOL",
									Prices:     []int64{150, 152, 151},
								},
								Quote: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "USDT",
									Prices:     []int64{1, 1, 1},
								},
							},
						},
					},
				},
				{
					{
						Name: "ForexBrokerA",
						Pair: []struct {
							Base struct {
								Instrument string  `abi:"instrument" json:"instrument"`
								Prices     []int64 `abi:"prices" json:"prices"`
							} `abi:"base" json:"base"`
							Quote struct {
								Instrument string  `abi:"instrument" json:"instrument"`
								Prices     []int64 `abi:"prices" json:"prices"`
							} `abi:"quote" json:"quote"`
						}{
							{
								Base: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "EUR",
									Prices:     []int64{10710, 10715, 10712},
								},
								Quote: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "USD",
									Prices:     []int64{1, 1, 1},
								},
							},
							{
								Base: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "GBP",
									Prices:     []int64{12750, 12755, 12751},
								},
								Quote: struct {
									Instrument string  `abi:"instrument" json:"instrument"`
									Prices     []int64 `abi:"prices" json:"prices"`
								}{
									Instrument: "JPY",
									Prices:     []int64{157, 157, 157},
								},
							},
						},
					},
				},
			}},
			Expected: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000009000000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000005a000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000f43727970746f45786368616e6765310000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000034254430000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000109a00000000000000000000000000000000000000000000000000000000000010a0400000000000000000000000000000000000000000000000000000000000109d20000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000035553440000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000003455448000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000dac0000000000000000000000000000000000000000000000000000000000000db60000000000000000000000000000000000000000000000000000000000000db10000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000035553440000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000f43727970746f45786368616e67653200000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000140000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000003534f4c00000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000009600000000000000000000000000000000000000000000000000000000000000980000000000000000000000000000000000000000000000000000000000000097000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000800000000000000000000000000000000000000000000000000000000000000004555344540000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000c466f72657842726f6b6572410000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000280000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000034555520000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000029d600000000000000000000000000000000000000000000000000000000000029db00000000000000000000000000000000000000000000000000000000000029d800000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000355534400000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000001400000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000034742500000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000031ce00000000000000000000000000000000000000000000000000000000000031d300000000000000000000000000000000000000000000000000000000000031cf0000000000000000000000000000000000000000000000000000000000000040000000000000000000000000000000000000000000000000000000000000008000000000000000000000000000000000000000000000000000000000000000034a505900000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000003000000000000000000000000000000000000000000000000000000000000009d000000000000000000000000000000000000000000000000000000000000009d000000000000000000000000000000000000000000000000000000000000009d"),
		},
	}

	parser := NewEtherParser()
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			args, err := parser.Serialize(tc.Input)
			assert.NoError(t, err)

			actual, err := args.Encode(tc.Args...)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func TestAbiParser_Deserialize(t *testing.T) {
	type Testcase struct {
		Name     string
		Input    AbiElements
		Expected types.Elements
	}

	testcases := []Testcase{
		{
			Name:     "number",
			Input:    AbiElements{{Type: crossTypes[types.Number]}},
			Expected: types.Elements{{Type: types.Number}},
		},
		{
			Name:     "string",
			Input:    AbiElements{{Type: crossTypes[types.String]}},
			Expected: types.Elements{{Type: types.String}},
		},
		{
			Name:     "bool",
			Input:    AbiElements{{Type: crossTypes[types.Bool]}},
			Expected: types.Elements{{Type: types.Bool}},
		},
		{
			Name:     "number,string,bool",
			Input:    AbiElements{{Type: crossTypes[types.Number]}, {Type: crossTypes[types.String]}, {Type: crossTypes[types.Bool]}},
			Expected: types.Elements{{Type: types.Number}, {Type: types.String}, {Type: types.Bool}},
		},
		{
			Name:     "string[]",
			Input:    AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}}},
			Expected: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.String}}}},
		},
		{
			Name:     "string[][]",
			Input:    AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}}}},
			Expected: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.String}}}}}},
		},
		{
			Name:     "(string name, number value, bool valid)",
			Input:    AbiElements{{Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "value", "valid"}, TupleElems: []*abi.Type{{T: abi.StringTy}, {T: abi.IntTy}, {T: abi.BoolTy}}}}},
			Expected: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "value", Type: types.Number}, {Name: "valid", Type: types.Bool}}}},
		},
		{
			Name:     "(string name, number value, bool valid)[]",
			Input:    AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "value", "valid"}, TupleElems: []*abi.Type{{T: abi.StringTy}, {T: abi.IntTy}, {T: abi.BoolTy}}}}}},
			Expected: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "value", Type: types.Number}, {Name: "valid", Type: types.Bool}}}}}},
		},
		{
			Name:     "(string name, (number amount, bool valid) detail)",
			Input:    AbiElements{{Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "detail"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.TupleTy, TupleRawNames: []string{"amount", "valid"}, TupleElems: []*abi.Type{&abiNumberTy, &abiBoolTy}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Amount", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"amount" json:"amount"`)}, {Name: "Valid", Type: abiBoolTy.GetType(), Tag: reflect.StructTag(`abi:"valid" json:"valid"`)}})}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Name", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"name" json:"name"`)}, {Name: "Detail", Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"amount", "valid"}, TupleElems: []*abi.Type{&abiNumberTy, &abiBoolTy}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Amount", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"amount" json:"amount"`)}, {Name: "Valid", Type: abiBoolTy.GetType(), Tag: reflect.StructTag(`abi:"valid" json:"valid"`)}})}.GetType(), Tag: reflect.StructTag(`abi:"detail" json:"detail"`)}})}}},
			Expected: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "detail", Type: types.Object, Children: types.Elements{{Name: "amount", Type: types.Number}, {Name: "valid", Type: types.Bool}}}}}},
		},
		{
			Name:     "(string name, ((string instrument, number[] prices) base, (string instrument, number[] prices) quote)[] pair)[][]",
			Input:    AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "pair"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"base", "quote"}, TupleElems: []*abi.Type{{T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}, {T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Base", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"base" json:"base"`)}, {Name: "Quote", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"quote" json:"quote"`)}})}}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Name", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"name" json:"name"`)}, {Name: "Pair", Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"base", "quote"}, TupleElems: []*abi.Type{{T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}, {T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Base", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"base" json:"base"`)}, {Name: "Quote", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"quote" json:"quote"`)}})}}.GetType(), Tag: reflect.StructTag(`abi:"pair" json:"pair"`)}})}}}}},
			Expected: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "pair", Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "base", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Number}}}}}, {Name: "quote", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Number}}}}}}}}}}}}}}}},
		},
	}

	parser := NewEtherParser()
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			actual, err := parser.Deserialize(tc.Input)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}
