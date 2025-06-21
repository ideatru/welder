package builder

import (
	"reflect"

	"github.com/ideatru/welder/types"
)

var DefaultReflectReplacers = make(map[types.ElementType]ReflectFn)

func DefaultStructTag(tag string) reflect.StructTag {
	return reflect.StructTag(`json:"` + tag + `"`)
}
