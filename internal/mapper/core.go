package mapper

import (
	"fmt"
	"reflect"

	"github.com/deferz/go-mapster/internal/cache"
)

// ValueConverter defines an interface for custom value conversion
type ValueConverter interface {
	Convert(src any) (any, error)
}

// FieldResolver defines an interface for custom field name resolution
type FieldResolver interface {
	ResolveField(srcType reflect.Type, dstType reflect.Type, fieldName string) (string, error)
}

// MapValue is the core mapping function responsible for mapping source value to target value
// Parameters:
//   - src: Reflection value of the source
//   - dst: Reflection value of the target (must be settable)
//
// Returns:
//   - error: Returns an error if mapping fails
func MapValue(src, dst reflect.Value) error {
	// Check for nil source
	if !src.IsValid() {
		return fmt.Errorf("source value is invalid")
	}

	// Check if target is settable
	if !dst.CanSet() {
		return fmt.Errorf("target value is not settable")
	}

	srcType := src.Type()
	dstType := dst.Type()

	// If types are identical, assign directly
	if srcType == dstType {
		dst.Set(src)
		return nil
	}

	// Get type information from cache for fast type checking
	typeCache := cache.GetGlobalCache()

	// Get or build source type info
	srcTypeInfo := typeCache.Get(srcType)
	if srcTypeInfo == nil {
		srcTypeInfo = cache.BuildTypeInfo(srcType)
		typeCache.Store(srcType, srcTypeInfo)
	}

	// Get or build destination type info
	dstTypeInfo := typeCache.Get(dstType)
	if dstTypeInfo == nil {
		dstTypeInfo = cache.BuildTypeInfo(dstType)
		typeCache.Store(dstType, dstTypeInfo)
	}

	// Choose mapping strategy based on cached type information
	if dstTypeInfo.IsStruct {
		return mapStruct(src, dst)
	} else if dstTypeInfo.IsCollection {
		return mapCollection(src, dst)
	} else if dstTypeInfo.IsMap {
		return mapMap(src, dst)
	} else if dst.Kind() == reflect.Ptr {
		return mapPointer(src, dst)
	} else {
		// Try basic type conversion
		return mapBasicType(src, dst)
	}
}

// mapBasicType handles conversion between basic types
func mapBasicType(src, dst reflect.Value) error {
	srcType := src.Type()
	dstType := dst.Type()

	// Handle nil pointers
	if srcType.Kind() == reflect.Ptr && src.IsNil() {
		return nil // Skip mapping for nil pointers
	}

	// Try direct conversion
	if srcType.ConvertibleTo(dstType) {
		dst.Set(src.Convert(dstType))
		return nil
	}

	return fmt.Errorf("cannot convert from %s to %s", srcType, dstType)
}

// mapPointer handles pointer type mapping
func mapPointer(src, dst reflect.Value) error {
	// If destination is nil pointer, create a new instance
	if dst.IsNil() {
		dst.Set(reflect.New(dst.Type().Elem()))
	}

	// If source is nil, do nothing
	if src.Kind() == reflect.Ptr && src.IsNil() {
		return nil
	}

	// If source is not a pointer, map to the pointer's element
	if src.Kind() != reflect.Ptr {
		return MapValue(src, dst.Elem())
	}

	// Both are pointers, map their elements
	return MapValue(src.Elem(), dst.Elem())
}
