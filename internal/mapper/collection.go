package mapper

import (
	"fmt"
	"reflect"

	"github.com/deferz/go-mapster/internal/cache"
)

// mapCollection handles mapping of slices and arrays
// Features:
// 1. For slices: Creates a brand new slice with length equal to source slice
// 2. For arrays: Creates a brand new array with length equal to target array type
// 3. Target will be completely replaced, not preserving original data
func mapCollection(src, dst reflect.Value) error {
	// Verify source value is slice or array
	if src.Kind() != reflect.Slice && src.Kind() != reflect.Array {
		return fmt.Errorf("source value is not a slice or array, but %s", src.Kind())
	}

	// Verify target value is slice or array
	if dst.Kind() != reflect.Slice && dst.Kind() != reflect.Array {
		return fmt.Errorf("target value is not a slice or array, but %s", dst.Kind())
	}

	// Get type information from cache
	typeCache := cache.GetGlobalCache()
	srcType := src.Type()
	dstType := dst.Type()

	srcTypeInfo := typeCache.Get(srcType)
	if srcTypeInfo == nil {
		// If not in cache, build and store it
		srcTypeInfo = cache.BuildTypeInfo(srcType)
		typeCache.Store(srcType, srcTypeInfo)
	}

	dstTypeInfo := typeCache.Get(dstType)
	if dstTypeInfo == nil {
		// If not in cache, build and store it
		dstTypeInfo = cache.BuildTypeInfo(dstType)
		typeCache.Store(dstType, dstTypeInfo)
	}

	// Get source length
	srcLen := src.Len()

	var dstVal reflect.Value
	var mapLen int

	// Create new target value based on target type
	if dstTypeInfo.IsCollection && dst.Kind() == reflect.Slice {
		// Slice: Create new slice with same length as source
		dstVal = reflect.MakeSlice(dst.Type(), srcLen, srcLen)
		mapLen = srcLen
	} else {
		// Array: Create new array with length of target array type
		dstVal = reflect.New(dst.Type()).Elem()
		dstLen := dst.Type().Len()

		// Calculate actual number of elements to map (take minimum)
		mapLen = srcLen
		if dstLen < srcLen {
			mapLen = dstLen
		}
	}

	// Map each element
	for i := 0; i < mapLen; i++ {
		srcElem := src.Index(i)
		dstElem := dstVal.Index(i)

		// Recursively map element
		if err := MapValue(srcElem, dstElem); err != nil {
			return fmt.Errorf("failed to map element at index %d: %w", i, err)
		}
	}

	// Set new value to target
	dst.Set(dstVal)
	return nil
}

// mapMap handles mapping from Map to Map
// Supports conversion of both keys and values
func mapMap(src, dst reflect.Value) error {
	// Verify source value is a Map
	if src.Kind() != reflect.Map {
		return fmt.Errorf("source value is not a Map, but %s", src.Kind())
	}

	// Verify target value is a Map
	if dst.Kind() != reflect.Map {
		return fmt.Errorf("target value is not a Map, but %s", dst.Kind())
	}

	// Get type information from cache
	typeCache := cache.GetGlobalCache()
	srcType := src.Type()
	dstType := dst.Type()

	srcTypeInfo := typeCache.Get(srcType)
	if srcTypeInfo == nil {
		// If not in cache, build and store it
		srcTypeInfo = cache.BuildTypeInfo(srcType)
		typeCache.Store(srcType, srcTypeInfo)
	}

	dstTypeInfo := typeCache.Get(dstType)
	if dstTypeInfo == nil {
		// If not in cache, build and store it
		dstTypeInfo = cache.BuildTypeInfo(dstType)
		typeCache.Store(dstType, dstTypeInfo)
	}

	// Get target Map type information from cached type info
	dstKeyType := dstType.Key()
	dstElemType := dstType.Elem()

	// Create new target Map
	dstMap := reflect.MakeMap(dstType)

	// Iterate through all key-value pairs in source Map
	for _, key := range src.MapKeys() {
		// Get source value
		srcValue := src.MapIndex(key)

		// Create target key
		dstKey := reflect.New(dstKeyType).Elem()
		// Try direct key conversion first
		if key.Type() == dstKeyType {
			dstKey.Set(key)
		} else if key.Type().ConvertibleTo(dstKeyType) {
			dstKey.Set(key.Convert(dstKeyType))
		} else if err := MapValue(key, dstKey); err != nil {
			return fmt.Errorf("failed to map Map key: %w", err)
		}

		// Create target value
		dstValue := reflect.New(dstElemType).Elem()
		if err := MapValue(srcValue, dstValue); err != nil {
			return fmt.Errorf("failed to map Map value: %w", err)
		}

		// Set key-value pair
		dstMap.SetMapIndex(dstKey, dstValue)
	}

	// Set new Map to target value
	dst.Set(dstMap)
	return nil
}
