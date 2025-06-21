package ether

import (
	"fmt"
	"math/big"
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ideatru/welder/types"
	"github.com/stretchr/testify/assert"
)

func TestEtherParser_Serialize(t *testing.T) {

	type Testcase struct {
		Name     string
		Input    types.Elements
		Values   []any
		Expected []byte
	}

	testcases := []Testcase{
		{
			Name:     "string,bool,address,bytes,int,uint",
			Input:    types.Elements{{Type: types.String}, {Type: types.Bool}, {Type: types.Address}, {Type: types.Bytes}, {Type: types.Int}, {Type: types.Uint}},
			Values:   []any{"string", true, common.HexToAddress("0xB035aD4B31759d909178d32da02266BD199c7e15"), hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000006737472696e670000000000000000000000000000000000000000000000000000"), int64(1), uint64(2)},
			Expected: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000c00000000000000000000000000000000000000000000000000000000000000001000000000000000000000000b035ad4b31759d909178d32da02266bd199c7e150000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000006737472696e670000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000006000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000006737472696e670000000000000000000000000000000000000000000000000000"),
		},
		{
			Name:     "int128,uint256,bytes32",
			Input:    types.Elements{{Type: types.Int, Size: 128}, {Type: types.Uint, Size: 256}, {Type: types.Bytes, Size: 32}},
			Values:   []any{big.NewInt(1), big.NewInt(101), [32]byte(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001"))},
			Expected: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000000100000000000000000000000000000000000000000000000000000000000000650000000000000000000000000000000000000000000000000000000000000001"),
		},
		{
			Name:     "string[]",
			Input:    types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.String}}}},
			Values:   []any{[]string{"hello", "world"}},
			Expected: hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000020000000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000400000000000000000000000000000000000000000000000000000000000000080000000000000000000000000000000000000000000000000000000000000000568656c6c6f0000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000000005776f726c64000000000000000000000000000000000000000000000000000000"),
		},
		{
			Name:  "uint256[][]",
			Input: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Uint, Size: 256}}}}}},
			Values: []any{
				[][]*big.Int{
					{big.NewInt(1), big.NewInt(2), big.NewInt(3)},
					{big.NewInt(4), big.NewInt(5), big.NewInt(1000)},
				},
			},
			Expected: hexutil.MustDecode("0x00000000000000000000000000000000000000000000000000000000000000200000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000c0000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000300000000000000000000000000000000000000000000000000000000000000030000000000000000000000000000000000000000000000000000000000000004000000000000000000000000000000000000000000000000000000000000000500000000000000000000000000000000000000000000000000000000000003e8"),
		},
		{
			Name:  "(string name, int count, uint256 balance, bytes[32] signature, bool valid)",
			Input: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "count", Type: types.Int}, {Name: "balance", Type: types.Uint, Size: 256}, {Name: "signature", Type: types.Bytes, Size: 32}, {Name: "valid", Type: types.Bool}}}},
			Values: []any{struct {
				Name      string
				Count     int64
				Balance   *big.Int
				Signature [32]byte `abi:"signature"`
				Valid     bool
			}{
				Name:      "welder",
				Count:     10,
				Balance:   big.NewInt(9999),
				Signature: [32]byte(hexutil.MustDecode("0x0000000000000000000000000000000000000000000000000000000000000001")),
				Valid:     true,
			}},
			Expected: hexutil.MustDecode("0x000000000000000000000000000000000000000000000000000000000000002000000000000000000000000000000000000000000000000000000000000000a0000000000000000000000000000000000000000000000000000000000000000a000000000000000000000000000000000000000000000000000000000000270f00000000000000000000000000000000000000000000000000000000000000010000000000000000000000000000000000000000000000000000000000000001000000000000000000000000000000000000000000000000000000000000000677656c6465720000000000000000000000000000000000000000000000000000"),
		},
		{
			Name:  "(string name, ((string instrument, uint256[] prices) base, (string instrument, uint256[] prices) quote)[] pair)[][]",
			Input: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "pair", Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "base", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Uint, Size: 256}}}}}, {Name: "quote", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Uint, Size: 256}}}}}}}}}}}}}}}},
			Values: []any{[][]struct {
				Name string `abi:"name" json:"name"`
				Pair []struct {
					Base struct {
						Instrument string     `abi:"instrument" json:"instrument"`
						Prices     []*big.Int `abi:"prices" json:"prices"`
					} `abi:"base" json:"base"`
					Quote struct {
						Instrument string     `abi:"instrument" json:"instrument"`
						Prices     []*big.Int `abi:"prices" json:"prices"`
					} `abi:"quote" json:"quote"`
				} `abi:"pair" json:"pair"`
			}{
				{
					{
						Name: "CryptoExchange1",
						Pair: []struct {
							Base struct {
								Instrument string     `abi:"instrument" json:"instrument"`
								Prices     []*big.Int `abi:"prices" json:"prices"`
							} `abi:"base" json:"base"`
							Quote struct {
								Instrument string     `abi:"instrument" json:"instrument"`
								Prices     []*big.Int `abi:"prices" json:"prices"`
							} `abi:"quote" json:"quote"`
						}{
							{
								Base: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "BTC",
									Prices: []*big.Int{
										big.NewInt(68000),
										big.NewInt(68100),
										big.NewInt(68050),
									},
								},
								Quote: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "USD",
									Prices:     []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)},
								},
							},
							{
								Base: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "ETH",
									Prices:     []*big.Int{big.NewInt(3500), big.NewInt(3510), big.NewInt(3505)},
								},
								Quote: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "USD",
									Prices: []*big.Int{
										big.NewInt(1),
										big.NewInt(1),
										big.NewInt(1),
									},
								},
							},
						},
					},
					{
						Name: "CryptoExchange2",
						Pair: []struct {
							Base struct {
								Instrument string     `abi:"instrument" json:"instrument"`
								Prices     []*big.Int `abi:"prices" json:"prices"`
							} `abi:"base" json:"base"`
							Quote struct {
								Instrument string     `abi:"instrument" json:"instrument"`
								Prices     []*big.Int `abi:"prices" json:"prices"`
							} `abi:"quote" json:"quote"`
						}{
							{
								Base: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "SOL",
									Prices:     []*big.Int{big.NewInt(150), big.NewInt(152), big.NewInt(151)},
								},
								Quote: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "USDT",
									Prices: []*big.Int{
										big.NewInt(1),
										big.NewInt(1),
										big.NewInt(1),
									},
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
								Instrument string     `abi:"instrument" json:"instrument"`
								Prices     []*big.Int `abi:"prices" json:"prices"`
							} `abi:"base" json:"base"`
							Quote struct {
								Instrument string     `abi:"instrument" json:"instrument"`
								Prices     []*big.Int `abi:"prices" json:"prices"`
							} `abi:"quote" json:"quote"`
						}{
							{
								Base: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "EUR",
									Prices:     []*big.Int{big.NewInt(10710), big.NewInt(10715), big.NewInt(10712)},
								},
								Quote: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "USD",
									Prices:     []*big.Int{big.NewInt(1), big.NewInt(1), big.NewInt(1)},
								},
							},
							{
								Base: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "GBP",
									Prices:     []*big.Int{big.NewInt(12750), big.NewInt(12755), big.NewInt(12751)},
								},
								Quote: struct {
									Instrument string     `abi:"instrument" json:"instrument"`
									Prices     []*big.Int `abi:"prices" json:"prices"`
								}{
									Instrument: "JPY",
									Prices:     []*big.Int{big.NewInt(157), big.NewInt(157), big.NewInt(157)},
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

			actual, err := args.Encode(tc.Values...)
			fmt.Println(hexutil.Encode(actual))
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}

func TestEtherParser_Deserialize(t *testing.T) {
	var (
		abiStringTy = abi.Type{T: abi.StringTy}
		abiNumberTy = abi.Type{T: abi.IntTy}
		abiBoolTy   = abi.Type{T: abi.BoolTy}
	)

	type Testcase struct {
		Name     string
		Input    AbiElements
		Expected types.Elements
	}

	testcases := []Testcase{
		{
			Name:     "addrses,bytes,bytes32",
			Input:    AbiElements{{Type: abi.Type{T: abi.AddressTy}}, {Type: abi.Type{T: abi.BytesTy}}, {Type: abi.Type{T: abi.FixedBytesTy, Size: 32}}},
			Expected: types.Elements{{Type: types.Address}, {Type: types.Bytes}, {Type: types.Bytes, Size: 32}},
		},
		{
			Name:     "uint256[100]",
			Input:    AbiElements{{Type: abi.Type{T: abi.ArrayTy, Size: 100, Elem: &abi.Type{T: abi.UintTy, Size: 256}}}},
			Expected: types.Elements{{Type: types.Array, Size: 100, Children: types.Elements{{Type: types.Uint, Size: 256}}}},
		},
		{
			Name:     "string[][]",
			Input:    AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.StringTy}}}}},
			Expected: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.String}}}}}},
		},
		{
			Name:     "(string name, number value, bool valid)",
			Input:    AbiElements{{Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "value", "valid"}, TupleElems: []*abi.Type{{T: abi.StringTy}, {T: abi.IntTy}, {T: abi.BoolTy}}}}},
			Expected: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "value", Type: types.Int}, {Name: "valid", Type: types.Bool}}}},
		},
		{
			Name:     "(string name, number value, bool valid)[]",
			Input:    AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "value", "valid"}, TupleElems: []*abi.Type{{T: abi.StringTy}, {T: abi.IntTy}, {T: abi.BoolTy}}}}}},
			Expected: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "value", Type: types.Int}, {Name: "valid", Type: types.Bool}}}}}},
		},
		{
			Name:     "(string name, (number amount, bool valid) detail)",
			Input:    AbiElements{{Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "detail"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.TupleTy, TupleRawNames: []string{"amount", "valid"}, TupleElems: []*abi.Type{&abiNumberTy, &abiBoolTy}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Amount", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"amount" json:"amount"`)}, {Name: "Valid", Type: abiBoolTy.GetType(), Tag: reflect.StructTag(`abi:"valid" json:"valid"`)}})}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Name", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"name" json:"name"`)}, {Name: "Detail", Type: abi.Type{T: abi.TupleTy, TupleRawNames: []string{"amount", "valid"}, TupleElems: []*abi.Type{&abiNumberTy, &abiBoolTy}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Amount", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"amount" json:"amount"`)}, {Name: "Valid", Type: abiBoolTy.GetType(), Tag: reflect.StructTag(`abi:"valid" json:"valid"`)}})}.GetType(), Tag: reflect.StructTag(`abi:"detail" json:"detail"`)}})}}},
			Expected: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "detail", Type: types.Object, Children: types.Elements{{Name: "amount", Type: types.Int}, {Name: "valid", Type: types.Bool}}}}}},
		},
		{
			Name:     "(string name, ((string instrument, number[] prices) base, (string instrument, number[] prices) quote)[] pair)[][]",
			Input:    AbiElements{{Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"name", "pair"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"base", "quote"}, TupleElems: []*abi.Type{{T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}, {T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Base", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"base" json:"base"`)}, {Name: "Quote", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"quote" json:"quote"`)}})}}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Name", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"name" json:"name"`)}, {Name: "Pair", Type: abi.Type{T: abi.SliceTy, Elem: &abi.Type{T: abi.TupleTy, TupleRawNames: []string{"base", "quote"}, TupleElems: []*abi.Type{{T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}, {T: abi.TupleTy, TupleRawNames: []string{"instrument", "prices"}, TupleElems: []*abi.Type{&abiStringTy, {T: abi.SliceTy, Elem: &abiNumberTy}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}})}}, TupleType: reflect.StructOf([]reflect.StructField{{Name: "Base", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"base" json:"base"`)}, {Name: "Quote", Type: reflect.StructOf([]reflect.StructField{{Name: "Instrument", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"instrument" json:"instrument"`)}, {Name: "Prices", Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(), Tag: reflect.StructTag(`abi:"prices" json:"prices"`)}}), Tag: reflect.StructTag(`abi:"quote" json:"quote"`)}})}}.GetType(), Tag: reflect.StructTag(`abi:"pair" json:"pair"`)}})}}}}},
			Expected: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "name", Type: types.String}, {Name: "pair", Type: types.Array, Children: types.Elements{{Type: types.Object, Children: types.Elements{{Name: "base", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Int}}}}}, {Name: "quote", Type: types.Object, Children: types.Elements{{Name: "instrument", Type: types.String}, {Name: "prices", Type: types.Array, Children: types.Elements{{Type: types.Int}}}}}}}}}}}}}}}},
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
