package welder

import "github.com/ideatru/welder/types"

type Welder interface {
	Weld(schema types.Elements, data []byte) ([]any, error)
}
