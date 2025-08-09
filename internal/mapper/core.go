package mapper

import (
	"fmt"
	"reflect"
)

// ValueConverter handles type conversion (first built-in aliases/time, then support for registration)
type ValueConverter interface {
	Convert(from reflect.Value, to reflect.Type) (reflect.Value, bool)
}

// FieldResolver handles field resolution (supports custom resolvers)
type FieldResolver interface {
	Resolve(value reflect.Value, fieldName string) (reflect.Value, bool)
}

// MapValue is the core mapping function responsible for mapping source value to target value
// Parameters:
//   - src: Reflection value of the source
//   - dst: Reflection value of the target (must be settable)
//
// Returns:
//   - error: Returns an error if mapping fails
func MapValue(src, dst reflect.Value) error {
	// Handle nil source value
	if !src.IsValid() {
		return fmt.Errorf("source value is invalid")
	}

	// Ensure target value can be set
	if !dst.CanSet() {
		return fmt.Errorf("target value cannot be set")
	}

	// Handle pointer type source value
	if src.Kind() == reflect.Ptr {
		if src.IsNil() {
			// If source pointer is nil, set target to zero value
			dst.Set(reflect.Zero(dst.Type()))
			return nil
		}
		// Dereference pointer
		src = src.Elem()
	}

	// Get source and target types
	srcType := src.Type()
	dstType := dst.Type()

	// If types are identical, assign directly
	if srcType == dstType {
		dst.Set(src)
		return nil
	}

	// Choose mapping strategy based on target type
	switch dst.Kind() {
	case reflect.Struct:
		return mapStruct(src, dst)
	case reflect.Slice, reflect.Array:
		// Use mapCollection for both slices and arrays
		return mapCollection(src, dst)
	case reflect.Map:
		return mapMap(src, dst)
	case reflect.Ptr:
		return mapPointer(src, dst)
	default:
		// Try basic type conversion
		return mapBasicType(src, dst)
	}
}

// mapBasicType handles conversion between basic types
func mapBasicType(src, dst reflect.Value) error {
	srcType := src.Type()
	dstType := dst.Type()

	// Check if direct conversion is possible
	if src.Type().ConvertibleTo(dst.Type()) {
		dst.Set(src.Convert(dst.Type()))
		return nil
	}

	return fmt.Errorf("cannot convert type %s to %s", srcType, dstType)
}

// mapPointer handles mapping to pointer types
func mapPointer(src, dst reflect.Value) error {
	// Create new pointer of target type if nil
	if dst.IsNil() {
		dst.Set(reflect.New(dst.Type().Elem()))
	}

	// Map to the value pointed by the pointer
	return MapValue(src, dst.Elem())
}
