package welder

import (
	"encoding/json"

	"github.com/ideatru/welder/ether"
	"github.com/ideatru/welder/internal/builder"
	"github.com/ideatru/welder/types"
)

// NewEthereumBuilder returns a builder configured with Ethereum-specific options.
func NewEthereumBuilder() *builder.Builder {
	return builder.New(ether.EtherBuilderOptions)
}

// EthereumWelder implements the types.Welder interface for Ethereum ABI.
type EthereumWelder struct {
	parser  types.Parser[ether.AbiElements]
	builder *builder.Builder
}

// NewEthereum creates a new EthereumWelder with default configuration.
func NewEthereum() *EthereumWelder {
	return &EthereumWelder{
		parser:  ether.NewEtherParser(),
		builder: NewEthereumBuilder(),
	}
}

// Deserialize converts Ethereum ABI elements into types.Elements.
func (w *EthereumWelder) Deserialize(data ether.AbiElements) (types.Elements, error) {
	return w.parser.Deserialize(data)
}

// Serialize converts types.Elements into Ethereum ABI elements.
func (w *EthereumWelder) Serialize(elements types.Elements) (ether.AbiElements, error) {
	return w.parser.Serialize(elements)
}

// Weld builds Go types from the schema and unmarshals data into them.
func (w *EthereumWelder) Weld(schema types.Elements, data []byte) ([]any, error) {
	result, err := w.builder.Builds(schema)
	if err != nil {
		return nil, err
	}

	if err := json.Unmarshal(data, &result); err != nil {
		return nil, err
	}

	return result, nil
}

// Builder returns the underlying builder instance.
func (w *EthereumWelder) Builder() *builder.Builder {
	return w.builder
}
