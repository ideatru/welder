package main

import (
	"encoding/json"
	"fmt"
	"github.com/ethereum/go-ethereum/accounts/abi"
	"github.com/ideatru/welder"
	"github.com/ideatru/welder/ether"
	"github.com/ideatru/welder/types"
	"reflect"
)

var (
	// These define basic Ethereum ABI types that will be used in the EVM schema
	abiStringTy = abi.Type{T: abi.StringTy} // ABI string type
	abiNumberTy = abi.Type{T: abi.IntTy}    // ABI integer type
	abiBoolTy   = abi.Type{T: abi.BoolTy}   // ABI boolean type
)

var (
	// JsonSchema - A complex nested schema definition using Welder's type system
	// Comment describes the structure: "string exchange, (string name, ((string instrument, uint256[] prices) base,
	// (string instrument, uint256[] prices) quote)[] pair)[][]"
	JsonSchema = types.Elements{
		{Type: types.String}, // First element: a string (exchange name)
		{Type: types.Array, // Second element: an array
			Children: types.Elements{
				{Type: types.Array, // Array of arrays
					Children: types.Elements{
						{Type: types.Object, // Object within nested array
							Children: types.Elements{
								{Name: "name", Type: types.String}, // String field "name"
								{Name: "pair", Type: types.Array, // Array field "pair"
									Children: types.Elements{
										{Type: types.Object, // Object within "pair" array
											Children: types.Elements{
												{Name: "base", Type: types.Object, // "base" nested object
													Children: types.Elements{
														{Name: "instrument", Type: types.String}, // String field
														{Name: "prices", Type: types.Array, // Array field
															Children: types.Elements{
																{Type: types.Uint, Size: 256}, // Array of uint256
															},
														},
													},
												},
												{Name: "quote", Type: types.Object, // "quote" nested object
													Children: types.Elements{
														{Name: "instrument", Type: types.String}, // String field
														{Name: "prices", Type: types.Array, // Array field
															Children: types.Elements{
																{Type: types.Uint, Size: 256}, // Array of uint256
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
	}

	// EvmSchema - A different schema defined using go-ethereum's ABI types
	// Structure: "(string name, (number amount, bool valid) detail)"
	EvmSchema = ether.AbiElements{
		{Type: abi.Type{
			T:             abi.TupleTy,                // Main type is a tuple (struct)
			TupleRawNames: []string{"name", "detail"}, // Field names in the tuple
			TupleElems: []*abi.Type{
				&abiStringTy, // First field is a string
				{
					T:             abi.TupleTy,                           // Second field is another tuple
					TupleRawNames: []string{"amount", "valid"},           // Nested tuple field names
					TupleElems:    []*abi.Type{&abiNumberTy, &abiBoolTy}, // Types: number and boolean
					// Go struct type definition using reflection for the nested tuple
					TupleType: reflect.StructOf([]reflect.StructField{
						{Name: "Amount", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"amount" json:"amount"`)},
						{Name: "Valid", Type: abiBoolTy.GetType(), Tag: reflect.StructTag(`abi:"valid" json:"valid"`)},
					}),
				},
			},
			// Go struct type definition using reflection for the main tuple
			TupleType: reflect.StructOf([]reflect.StructField{
				{Name: "Name", Type: abiStringTy.GetType(), Tag: reflect.StructTag(`abi:"name" json:"name"`)},
				{Name: "Detail",
					Type: abi.Type{
						T:             abi.TupleTy,
						TupleRawNames: []string{"amount", "valid"},
						TupleElems:    []*abi.Type{&abiNumberTy, &abiBoolTy},
						TupleType: reflect.StructOf([]reflect.StructField{
							{Name: "Amount", Type: abiNumberTy.GetType(), Tag: reflect.StructTag(`abi:"amount" json:"amount"`)},
							{Name: "Valid", Type: abiBoolTy.GetType(), Tag: reflect.StructTag(`abi:"valid" json:"valid"`)},
						}),
					}.GetType(),
					Tag: reflect.StructTag(`abi:"detail" json:"detail"`),
				},
			}),
		}},
	}
)

func main() {
	// Step 1: Initialize Welder for Ethereum operations
	welder := welder.NewEthereum()

	// Step 2: Generate example JSON data that matches the JsonSchema structure
	// The Builder().Builds() method creates sample data for all elements in the schema
	jsonResult, err := welder.Builder().Builds(JsonSchema)
	if err != nil {
		panic(err)
	}
	// Print the generated example data
	PrettyPrint("JSON", jsonResult)

	// Step 3: Convert the EVM schema (go-ethereum ABI format) to Welder's schema format
	// Deserialize transforms the EVM ABI schema into Welder's schema types
	jsonSchema, err := welder.Deserialize(EvmSchema)
	if err != nil {
		panic(err)
	}

	// Step 4: Generate example data for the first element of the converted schema
	// Build() creates sample data for a single schema element
	evmResult, err := welder.Builder().Build(jsonSchema[0])
	if err != nil {
		panic(err)
	}
	// Print the generated example data
	PrettyPrint("EVM", []any{evmResult})
}

func PrettyPrint(title string, values []any) {
	fmt.Println("==========================")
	fmt.Printf("%s Result:\n", title)
	fmt.Println("--------------------------")
	data, err := json.MarshalIndent(values, "", "  ")
	if err != nil {
		fmt.Println("Error marshaling data:", err)
		return
	}
	fmt.Printf("%v\n\n", string(data))

	fmt.Printf("%s Built Types:\n", title)
	fmt.Println("--------------------------")
	for i, value := range values {
		ty := reflect.TypeOf(value)
		fmt.Printf("Type of Index(%d): %v\n", i, ty)
	}

	fmt.Print("\n\n")
}
