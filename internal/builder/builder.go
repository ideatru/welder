package builder

import (
	"fmt"
	"math/big"
	"reflect"

	"github.com/ideatru/welder/internal/utils"
	"github.com/ideatru/welder/types"
)

type (
	// ReflectFn is a function type that converts a types.Element to a reflect.Type
	// Used for custom type mapping in the builder
	ReflectFn func(types.Element) (reflect.Type, error)

	// StructTagFn is a function type that generates a reflect.StructTag from a field name
	// Used to customize how struct field tags are generated
	StructTagFn func(string) reflect.StructTag

	// Option contains configuration options for the Builder
	// Allows customization of type reflection and struct tag generation
	Option struct {
		// ReflectReplacers maps element types to custom reflection functions
		ReflectReplacers map[types.ElementType]ReflectFn

		// StructTagReplacer is a function that generates struct tags for object fields
		StructTagReplacer StructTagFn
	}
)

// Builder is responsible for creating Go types based on element definitions
// Uses reflection to dynamically build types that match the specified elements
type Builder struct {
	// ReflectReplacers maps element types to custom reflection functions
	ReflectReplacers map[types.ElementType]ReflectFn

	// StructTagReplacer is a function that generates struct tags for object fields
	StructTagReplacer StructTagFn
}

// New creates a new Builder with the provided options
// If no options are provided, default options are used
// Returns a pointer to the configured Builder
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

func (b *Builder) Builds(elements types.Elements) ([]any, error) {
	values := make([]any, 0, len(elements))
	for _, elem := range elements {
		value, err := b.Build(elem)
		if err != nil {
			return nil, err
		}
		values = append(values, value)
	}
	return values, nil
}

// Build constructs a new instance of the type defined by the provided element
// Uses reflection to create the appropriate Go type and initializes it
// Returns the initialized value or an error if type building fails
// Recovers from panics that might occur during reflection operations
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

	value = reflect.New(ty).Interface()
	return
}

// buildType creates a reflect.Type based on the provided element
// Dispatches to the appropriate type builder based on the element's type
// Returns the reflect.Type or an error if type building fails
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

// buildString creates a reflect.Type for a string element
// Returns the type for Go's string
func (b *Builder) buildString(elem types.Element) (reflect.Type, error) {
	return reflect.TypeOf(""), nil
}

// buildUint creates a reflect.Type for an unsigned integer element
// Selects the appropriate uint type based on the element's size
// Returns a *big.Int type for sizes that don't match standard Go uint types
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

// buildInt creates a reflect.Type for a signed integer element
// Selects the appropriate int type based on the element's size
// Returns a *big.Int type for sizes that don't match standard Go int types
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

// buildFloat creates a reflect.Type for a floating-point element
// Selects float32 or float64 based on the element's size
// Returns an error for sizes that don't match standard Go float types
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

// buildBool creates a reflect.Type for a boolean element
// Returns the type for Go's bool
func (b *Builder) buildBool(elem types.Element) (reflect.Type, error) {
	return reflect.TypeOf(false), nil
}

// buildBytes creates a reflect.Type for a bytes element
// Returns the type for Go's []byte (byte slice)
func (b *Builder) buildBytes(elem types.Element) (reflect.Type, error) {
	return reflect.TypeOf([]byte{}), nil
}

// buildAddress creates a reflect.Type for an address element
// Returns an error as the default implementation doesn't support address types
// Custom implementations should be provided through ReflectReplacers
func (b *Builder) buildAddress(elem types.Element) (reflect.Type, error) {
	return nil, fmt.Errorf("`Builder` does not support type %q", types.Address)
}

// buildArray creates a reflect.Type for an array element
// Returns a slice type of the child element's type
// Returns an error if the array doesn't have exactly one child element
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

// buildObject creates a reflect.Type for an object element
// Constructs a struct type with fields matching the object's children
// Returns an error if the object has no children
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

// unwrapReplacer applies a custom type builder if available for the element type
// Falls back to the default builder if no custom builder is registered
// Returns the resulting reflect.Type or an error if building fails
func unwrapReplacer(elem types.Element, defaultFn ReflectFn, replacerFns map[types.ElementType]ReflectFn) (reflect.Type, error) {
	if fn, ok := replacerFns[elem.Type]; ok {
		return fn(elem)
	}

	return defaultFn(elem)
}
