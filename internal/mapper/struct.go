package mapper

import (
	"fmt"
	"reflect"

	"github.com/deferz/go-mapster/internal/cache"
)

// mapStruct handles mapping from struct to struct
// This function iterates through all fields of the target struct and attempts to find corresponding fields in the source struct
func mapStruct(src, dst reflect.Value) error {
	// Verify source value is a struct
	if src.Kind() != reflect.Struct {
		return fmt.Errorf("source value is not a struct, but %s", src.Kind())
	}

	// Verify target value is a struct
	if dst.Kind() != reflect.Struct {
		return fmt.Errorf("target value is not a struct, but %s", dst.Kind())
	}

	// Get source and target type information from cache
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

	// Use cached field information for target struct
	for _, fieldInfo := range dstTypeInfo.Fields {
		// Get target field
		dstField := dst.Field(fieldInfo.Index)

		// Field is already verified as exported in BuildTypeInfo
		// Get field name
		fieldName := fieldInfo.Name

		// Find corresponding field in source struct using cached field map
		srcFieldInfo, exists := srcTypeInfo.FieldsMap[fieldName]
		if exists {
			srcField := src.Field(srcFieldInfo.Index)
			// Recursively map field value
			if err := MapValue(srcField, dstField); err != nil {
				return fmt.Errorf("failed to map field %s: %w", fieldName, err)
			}
		} else {
			// Try to find in embedded fields
			if embeddedField, found := findFieldInEmbedded(src, fieldName); found {
				if err := MapValue(embeddedField, dstField); err != nil {
					return fmt.Errorf("failed to map field %s from embedded: %w", fieldName, err)
				}
			}
			// If field not found, skip (keep original value in target field)
		}
	}

	return nil
}

// findSourceField finds a field with the specified name in the source struct
// Only performs exact name matching
func findSourceField(src reflect.Value, fieldName string) reflect.Value {
	// Only do exact match
	return src.FieldByName(fieldName)
}
