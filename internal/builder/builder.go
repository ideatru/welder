package builder

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ideatru/welder/internal/utils"
	"github.com/ideatru/welder/types"
)

type (
	//
	ReflectFn func(types.Element) (reflect.Type, error)
	//
	StructTagFn func(string) reflect.StructTag
	//
	Option struct {
		ReflectReplacers  map[types.ElementType]ReflectFn
		StructTagReplacer StructTagFn
	}
)

type Builder struct {
	ReflectReplacers  map[types.ElementType]ReflectFn
	StructTagReplacer StructTagFn
}

func New(opts ...Option) *Builder {
	var opt Option
	if len(opts) > 0 {
		opt = opts[0]
	}

	builder := new(Builder)
	if opt.ReflectReplacers == nil {
		builder.ReflectReplacers = DefaultReflectReplacers
	} else {
		builder.ReflectReplacers = opt.ReflectReplacers
	}

	if opt.StructTagReplacer == nil {
		opt.StructTagReplacer = DefaultStructTag
	} else {
		builder.StructTagReplacer = opt.StructTagReplacer
	}

	return builder
}

func (b *Builder) Build(elem types.Element) (value any, err error) {
	defer func() {
		if r := recover(); r != nil {
			err = fmt.Errorf("%v", r)
			value = nil
			return
		}
	}()

	ty, err := b.buildType(elem)
	if err != nil {
		return nil, err
	}

	value = reflect.ValueOf(ty).Interface()
	return
}

func (b *Builder) buildType(elem types.Element) (reflect.Type, error) {
	switch elem.Type {
	case types.String:
		return unwrapReplacer(elem, b.buildString, b.ReflectReplacers)
	case types.Int:
		return unwrapReplacer(elem, b.buildInt, b.ReflectReplacers)
	case types.Uint:
		return unwrapReplacer(elem, b.buildUint, b.ReflectReplacers)
	case types.Float:
		return unwrapReplacer(elem, b.buildFloat, b.ReflectReplacers)
	case types.Bool:
		return unwrapReplacer(elem, b.buildBool, b.ReflectReplacers)
	case types.Bytes:
		return unwrapReplacer(elem, b.buildBytes, b.ReflectReplacers)
	case types.Address:
		return unwrapReplacer(elem, b.buildAddress, b.ReflectReplacers)
	case types.Array:
		return unwrapReplacer(elem, b.buildArray, b.ReflectReplacers)
	case types.Object:
		return unwrapReplacer(elem, b.buildObject, b.ReflectReplacers)
	}

	return nil, fmt.Errorf("`Builder does not support type %q", elem.Type)
}

func (b *Builder) buildString(elem types.Element) (reflect.Type, error) {
	return reflect.TypeOf(""), nil
}

func (b *Builder) buildUint(elem types.Element) (reflect.Type, error) {
	switch elem.Size {
	case 8:
		return reflect.TypeOf(uint8(0)), nil
	case 16:
		return reflect.TypeOf(uint16(0)), nil
	case 32:
		return reflect.TypeOf(uint32(0)), nil
	case 0, 64:
		return reflect.TypeOf(uint64(0)), nil
	default:
		return reflect.TypeOf(big.NewInt(0)), nil
	}

}

func (b *Builder) buildInt(elem types.Element) (reflect.Type, error) {
	switch elem.Size {
	case 8:
		return reflect.TypeOf(int8(0)), nil
	case 16:
		return reflect.TypeOf(int16(0)), nil
	case 32:
		return reflect.TypeOf(int32(0)), nil
	case 0, 64:
		return reflect.TypeOf(int64(0)), nil
	default:
		return reflect.TypeOf(big.NewInt(0)), nil
	}
}

func (b *Builder) buildFloat(elem types.Element) (reflect.Type, error) {
	switch elem.Size {
	case 32:
		return reflect.TypeOf(float32(0)), nil
	case 0, 64:
		return reflect.TypeOf(float64(0)), nil
	default:
		return nil, fmt.Errorf("`Builder` does not support type %q with size %v", elem.Type, elem.Size)
	}
}

func (b *Builder) buildBool(elem types.Element) (reflect.Type, error) {
	return reflect.TypeOf(false), nil
}

func (b *Builder) buildBytes(elem types.Element) (reflect.Type, error) {
	return reflect.TypeOf([]byte{}), nil
}

func (b *Builder) buildAddress(elem types.Element) (reflect.Type, error) {
	return nil, fmt.Errorf("`Builder` does not support type %q", types.Address)
}

func (b *Builder) buildArray(elem types.Element) (reflect.Type, error) {
	if len(elem.Children) != 1 {
		return nil, fmt.Errorf("array must have one child")
	}

	ty, err := b.buildType(elem.Children[0])
	if err != nil {
		return nil, err
	}

	return reflect.SliceOf(ty), nil
}

func (b *Builder) buildObject(elem types.Element) (reflect.Type, error) {
	if len(elem.Children) == 0 {
		return nil, fmt.Errorf("object must have at least one child")
	}

	fields := make([]reflect.StructField, 0)
	for _, childElem := range elem.Children {
		childTy, err := b.buildType(childElem)
		if err != nil {
			return nil, err
		}

		field := reflect.StructField{
			Name: utils.ToCamelCase(childElem.Name),
			Type: childTy,
			Tag:  b.StructTagReplacer(childElem.Name),
		}

		fields = append(fields, field)
	}

	return reflect.StructOf(fields), nil
}

func unwrapReplacer(elem types.Element, defaultFn ReflectFn, replacerFns map[types.ElementType]ReflectFn) (reflect.Type, error) {
	if fn, ok := replacerFns[elem.Type]; ok {
		return fn(elem)
	}

	return defaultFn(elem)
}
