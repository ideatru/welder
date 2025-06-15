# Welder

<p align="center">
  <img src="./docs/logo.png" width="500" alt="Welder">
</p>

A Go library for seamless type conversion and ABI encoding between schema definitions and Ethereum smart contracts.

## Overview

Welder provides a bridge between generic schema types and Ethereum ABI types, enabling easy encoding and decoding of data for smart contract interactions. It supports complex nested structures including arrays and objects, making it ideal for applications that need to dynamically convert data schemas to Ethereum-compatible formats.

## Features

- **Type Conversion**: Convert between schema types (string, number, boolean, array, object) and Ethereum ABI types
- **Complete ABI Support**: Full support for all Ethereum ABI types including:
  - Integers: `uint8` to `uint256`, `int8` to `int256`
  - Fixed bytes: `bytes1` to `bytes32`
  - Dynamic types: `bytes`, `string`, `address`, `hash`
- **Nested Structures**: Handle complex nested arrays and objects
- **ABI Encoding**: Generate ABI-compatible byte encodings for smart contract calls
- **Type Safety**: Strong typing with compile-time validation

## Installation

```bash
go get github.com/ideatru/welder
```

## Requirements

- Go 1.24.0 or later
- Ethereum go-ethereum library

## Quick Start

### Basic Usage

```go
package main

import (
    "fmt"
    "github.com/ideatru/welder/ether"
    "github.com/ideatru/welder/types"
)

func main() {
    // Create a converter
    converter := ether.NewConverter()
    
    // Define schema elements
    elements := types.Elements{
        {
            Name: "user",
            Type: types.String,
        },
        {
            Name: "balance",
            Type: types.Number,
        },
        {
            Name: "active",
            Type: types.Bool,
        },
    }
    
    // Encode to ABI format
    encoded, err := converter.Encode(elements)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("Encoded: %s\n", encoded)
}
```

### Working with Arrays

```go
// Define an array of strings
elements := types.Elements{
    {
        Name: "usernames",
        Type: types.Array,
        Children: []types.Element{
            {
                Type: types.String,
            },
        },
    },
}
```

### Working with Objects

```go
// Define a complex object
elements := types.Elements{
    {
        Name: "user",
        Type: types.Object,
        Children: []types.Element{
            {
                Name: "name",
                Type: types.String,
            },
            {
                Name: "age",
                Type: types.Number,
            },
            {
                Name: "verified",
                Type: types.Bool,
            },
        },
    },
}
```

## API Reference

### Types

#### ElementType

```go
type ElementType string

const (
    Number = ElementType("number")
    String = ElementType("string")
    Bool   = ElementType("boolean")
    Array  = ElementType("array")
    Object = ElementType("object")
)
```

#### Element

```go
type Element struct {
    Name     string
    Type     ElementType
    Children []Element
}
```

### Converter

#### NewConverter()

Creates a new ABI converter instance.

```go
converter := ether.NewConverter()
```

#### Encode(elements types.Elements) ([]byte, error)

Encodes schema elements into ABI-compatible JSON format.

```go
encoded, err := converter.Encode(elements)
```

### ABI Elements

#### Encode(values ...any) ([]byte, error)

Encodes values directly using ABI arguments.

```go
abiElements := ether.AbiElements{...}
encoded, err := abiElements.Encode(value1, value2, value3)
```

## Project Structure

```
.
├── ether/           # ABI conversion and encoding logic
├── types/           # Schema type definitions
├── solidity/        # Smart contract components (Foundry)
├── internal/        # Internal utilities
├── go.mod          # Go module definition
└── LICENSE         # MIT License
```

## Smart Contract Integration

The project includes Solidity components built with Foundry for smart contract development and testing.

### Prerequisites for Solidity Development

- [Foundry](https://book.getfoundry.sh/getting-started/installation)

### Build Smart Contracts

```bash
cd solidity
forge build
```

### Run Tests

```bash
cd solidity
forge test
```

## Contributing

1. Fork the repository
2. Create your feature branch (`git checkout -b feature/amazing-feature`)
3. Commit your changes (`git commit -m 'Add some amazing feature'`)
4. Push to the branch (`git push origin feature/amazing-feature`)
5. Open a Pull Request

## Testing

Run the Go tests:

```bash
go test ./...
```

Run specific package tests:

```bash
go test ./ether
go test ./types
```

## Examples

### Complete Example with Smart Contract Interaction

```go
package main

import (
    "fmt"
    "github.com/ideatru/welder/ether"
    "github.com/ideatru/welder/types"
)

func main() {
    converter := ether.NewConverter()
    
    // Define a complex schema matching a smart contract function
    elements := types.Elements{
        {
            Name: "transaction",
            Type: types.Object,
            Children: []types.Element{
                {
                    Name: "from",
                    Type: types.String,
                },
                {
                    Name: "to", 
                    Type: types.String,
                },
                {
                    Name: "amount",
                    Type: types.Number,
                },
                {
                    Name: "metadata",
                    Type: types.Array,
                    Children: []types.Element{
                        {
                            Type: types.String,
                        },
                    },
                },
            },
        },
    }
    
    // Convert to ABI format
    encoded, err := converter.Encode(elements)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("ABI-encoded schema: %s\n", encoded)
}
```

## License

This project is licensed under the MIT License - see the [LICENSE](LICENSE) file for details.

## Acknowledgments

- Built with [go-ethereum](https://github.com/ethereum/go-ethereum)
- Smart contract development powered by [Foundry](https://github.com/foundry-rs/foundry)

---

**Welder** - Bridging the gap between schemas and smart contracts. 