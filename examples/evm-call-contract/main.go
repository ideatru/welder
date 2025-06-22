package main

import (
	"context"
	"encoding/json"
	"fmt"

	"github.com/ethereum/go-ethereum"
	"github.com/ethereum/go-ethereum/common"
	"github.com/ethereum/go-ethereum/ethclient"
	"github.com/ideatru/welder"       // Welder main package for contract interaction
	"github.com/ideatru/welder/types" // Types used to define schemas
)

var (
	// RPC endpoint for connecting to the Optimism Sepolia testnet
	RPCEndpoint = "https://sepolia.optimism.io"

	// Smart contract address we're going to call
	ContractAddress = common.HexToAddress("0xEC1Dd7E79f631cf61e2937Bff1CcAf56B0EDB7c3")

	// Function signature in the format "functionName(paramType1,paramType2,...)"
	// This describes the function we want to call and its parameter types
	FunctionSignature = "getData(string,(address,string,(uint256,string)))"

	// Schema definition describes the structure of our data
	// This maps directly to the function parameters in our FunctionSignature
	Schema = types.Elements{
		{Type: types.String}, // First parameter: a simple string
		{
			Type: types.Object, // Second parameter: a complex object/struct
			Children: types.Elements{
				{Type: types.Address, Name: "owner"}, // Address field
				{Type: types.String, Name: "name"},   // String field
				{Type: types.Object, Name: "balance", Children: types.Elements{ // Nested object
					{Type: types.Uint, Size: 256, Name: "amount"}, // uint256 field
					{Type: types.String, Name: "currency"},        // String field
				}},
			},
		},
	}

	// JSON data that we want to send to the contract
	// Note how it matches the schema structure defined above
	Payload = []byte(`["Hello, World!!!",{"owner": "0xf39Fd6e51aad88F6F4ce6aB8827279cffFb92266", "name": "Ether", "balance": {"amount": 1000000000000000000, "currency": "ETH"}}]`)
)

func main() {
	// Step 1: Create a new Welder instance for Ethereum interaction
	welder := welder.NewEthereum()

	// Step 2: Serialize the schema to create an Arguments object
	// This generates the ABI types information from our schema
	args, err := welder.Serialize(Schema)
	if err != nil {
		panic(err)
	}

	// Step 3: "Weld" the JSON payload to match our schema
	// This transforms our JSON data into the format expected by the Ethereum ABI
	params, err := welder.Weld(Schema, Payload)
	if err != nil {
		panic(err)
	}

	// Step 4: Encode the data with the function signature
	// This produces the actual call data bytes to send to the contract
	data, err := args.EncodeWithFunctionSignature(FunctionSignature, params)
	if err != nil {
		panic(err)
	}

	// Step 5: Create a call message for the Ethereum contract
	callMsg := ethereum.CallMsg{
		To:   &ContractAddress, // Target contract address
		Data: data,             // Encoded function call data
	}

	// Step 6: Connect to the Ethereum network via RPC
	client, err := ethclient.Dial(RPCEndpoint)
	if err != nil {
		panic(err)
	}

	// Step 7: Execute the contract call (read-only call, not a transaction)
	calledData, err := client.CallContract(context.Background(), callMsg, nil)
	if err != nil {
		panic(err)
	}

	// Step 8: Print the raw response data
	fmt.Println("Called Data:")
	fmt.Println(string(calledData))

	// Step 9: Decode the response data using our Arguments object
	// This transforms the Ethereum ABI encoded response back to a readable format
	result, err := args.Decode(calledData)
	if err != nil {
		panic(err)
	}

	// Step 10: Convert the result to JSON and print it
	buffer, _ := json.Marshal(result)
	fmt.Println("Decoded Result:")
	fmt.Println(string(buffer))
}
