package ether

import (
	"reflect"
	"testing"

	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ethereum/go-ethereum/common/hexutil"
	"github.com/ideatru/crosschain/types"
	"github.com/stretchr/testify/assert"
)

func TestAbiConverter_encode(t *testing.T) {
	var (
		abiStringTy = crossTypes[types.String]
		abiNumberTy = crossTypes[types.Number]
		abiBoolTy   = crossTypes[types.Bool]
	)

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
			Name: "(string name, ((string instrument, number[] prices) base, (string instrument, number[] prices) quote)[] pair)[][]",
			Input: types.Elements{
				{
					Type: types.Array,
					Children: types.Elements{
						{
							Type: types.Array,
							Children: types.Elements{
								{
									Type: types.Object,
									Children: types.Elements{
										{
											Name: "name",
											Type: types.String,
										},
										{
											Name: "pair",
											Type: types.Array,
											Children: types.Elements{
												{
													Type: types.Object,
													Children: types.Elements{
														{
															Name: "base",
															Type: types.Object,
															Children: types.Elements{
																{
																	Name: "instrument",
																	Type: types.String,
																},
																{
																	Name: "prices",
																	Type: types.Array,
																	Children: types.Elements{
																		{
																			Type: types.Number,
																		},
																	},
																},
															},
														},
														{
															Name: "quote",
															Type: types.Object,
															Children: types.Elements{
																{
																	Name: "instrument",
																	Type: types.String,
																},
																{
																	Name: "prices",
																	Type: types.Array,
																	Children: types.Elements{
																		{
																			Type: types.Number,
																		},
																	},
																},
															},
														},
													},
												},
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Expected: AbiElements{
				{
					Type: abi.Type{
						T: abi.SliceTy,
						Elem: &abi.Type{
							T: abi.SliceTy,
							Elem: &abi.Type{
								T:             abi.TupleTy,
								TupleRawNames: []string{"name", "pair"},
								TupleElems: []*abi.Type{
									&abiStringTy,
									{
										T: abi.SliceTy,
										Elem: &abi.Type{
											T:             abi.TupleTy,
											TupleRawNames: []string{"base", "quote"},
											TupleElems: []*abi.Type{
												{
													T:             abi.TupleTy,
													TupleRawNames: []string{"instrument", "prices"},
													TupleElems: []*abi.Type{
														&abiStringTy,
														{
															T:    abi.SliceTy,
															Elem: &abiNumberTy,
														},
													},
													TupleType: reflect.StructOf([]reflect.StructField{
														{
															Name: "Instrument",
															Type: abiStringTy.GetType(),
															Tag:  reflect.StructTag(`abi:"instrument" json:"instrument"`),
														},
														{
															Name: "Prices",
															Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(),
															Tag:  reflect.StructTag(`abi:"prices" json:"prices"`),
														},
													}),
												},
												{
													T:             abi.TupleTy,
													TupleRawNames: []string{"instrument", "prices"},
													TupleElems: []*abi.Type{
														&abiStringTy,
														{
															T:    abi.SliceTy,
															Elem: &abiNumberTy,
														},
													},
													TupleType: reflect.StructOf([]reflect.StructField{
														{
															Name: "Instrument",
															Type: abiStringTy.GetType(),
															Tag:  reflect.StructTag(`abi:"instrument" json:"instrument"`),
														},
														{
															Name: "Prices",
															Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(),
															Tag:  reflect.StructTag(`abi:"prices" json:"prices"`),
														},
													}),
												},
											},
											TupleType: reflect.StructOf([]reflect.StructField{
												{
													Name: "Base",
													Type: reflect.StructOf([]reflect.StructField{
														{
															Name: "Instrument",
															Type: abiStringTy.GetType(),
															Tag:  reflect.StructTag(`abi:"instrument" json:"instrument"`),
														},
														{
															Name: "Prices",
															Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(),
															Tag:  reflect.StructTag(`abi:"prices" json:"prices"`),
														},
													}),
													Tag: reflect.StructTag(`abi:"base" json:"base"`),
												},
												{
													Name: "Quote",
													Type: reflect.StructOf([]reflect.StructField{
														{
															Name: "Instrument",
															Type: abiStringTy.GetType(),
															Tag:  reflect.StructTag(`abi:"instrument" json:"instrument"`),
														},
														{
															Name: "Prices",
															Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(),
															Tag:  reflect.StructTag(`abi:"prices" json:"prices"`),
														},
													}),
													Tag: reflect.StructTag(`abi:"quote" json:"quote"`),
												},
											}),
										},
									},
								},
								TupleType: reflect.StructOf([]reflect.StructField{
									{
										Name: "Name",
										Type: abiStringTy.GetType(),
										Tag:  reflect.StructTag(`abi:"name" json:"name"`),
									},
									{
										Name: "Pair",
										Type: abi.Type{
											T: abi.SliceTy,
											Elem: &abi.Type{
												T:             abi.TupleTy,
												TupleRawNames: []string{"base", "quote"},
												TupleElems: []*abi.Type{
													{
														T:             abi.TupleTy,
														TupleRawNames: []string{"instrument", "prices"},
														TupleElems: []*abi.Type{
															&abiStringTy,
															{
																T:    abi.SliceTy,
																Elem: &abiNumberTy,
															},
														},
														TupleType: reflect.StructOf([]reflect.StructField{
															{
																Name: "Instrument",
																Type: abiStringTy.GetType(),
																Tag:  reflect.StructTag(`abi:"instrument" json:"instrument"`),
															},
															{
																Name: "Prices",
																Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(),
																Tag:  reflect.StructTag(`abi:"prices" json:"prices"`),
															},
														}),
													},
													{
														T:             abi.TupleTy,
														TupleRawNames: []string{"instrument", "prices"},
														TupleElems: []*abi.Type{
															&abiStringTy,
															{
																T:    abi.SliceTy,
																Elem: &abiNumberTy,
															},
														},
														TupleType: reflect.StructOf([]reflect.StructField{
															{
																Name: "Instrument",
																Type: abiStringTy.GetType(),
																Tag:  reflect.StructTag(`abi:"instrument" json:"instrument"`),
															},
															{
																Name: "Prices",
																Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(),
																Tag:  reflect.StructTag(`abi:"prices" json:"prices"`),
															},
														}),
													},
												},
												TupleType: reflect.StructOf([]reflect.StructField{
													{
														Name: "Base",
														Type: reflect.StructOf([]reflect.StructField{
															{
																Name: "Instrument",
																Type: abiStringTy.GetType(),
																Tag:  reflect.StructTag(`abi:"instrument" json:"instrument"`),
															},
															{
																Name: "Prices",
																Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(),
																Tag:  reflect.StructTag(`abi:"prices" json:"prices"`),
															},
														}),
														Tag: reflect.StructTag(`abi:"base" json:"base"`),
													},
													{
														Name: "Quote",
														Type: reflect.StructOf([]reflect.StructField{
															{
																Name: "Instrument",
																Type: abiStringTy.GetType(),
																Tag:  reflect.StructTag(`abi:"instrument" json:"instrument"`),
															},
															{
																Name: "Prices",
																Type: abi.Type{T: abi.SliceTy, Elem: &abiNumberTy}.GetType(),
																Tag:  reflect.StructTag(`abi:"prices" json:"prices"`),
															},
														}),
														Tag: reflect.StructTag(`abi:"quote" json:"quote"`),
													},
												}),
											},
										}.GetType(),
										Tag: reflect.StructTag(`abi:"pair" json:"pair"`),
									},
								}),
							},
						},
					},
				},
			},
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
	}

	converter := NewConverter()
	for _, tc := range testcases {
		t.Run(tc.Name, func(t *testing.T) {
			args, err := converter.encode(tc.Input)
			assert.NoError(t, err)

			actual, err := args.Encode(tc.Args...)
			assert.NoError(t, err)
			assert.Equal(t, tc.Expected, actual)
		})
	}
}
