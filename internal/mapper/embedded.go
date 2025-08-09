package mapper

import (
	"reflect"

	"github.com/deferz/go-mapster/internal/cache"
)

// processEmbeddedFields processes embedded fields in structs
// Embedded fields are anonymous fields in Go and require special handling
func processEmbeddedFields(src, dst reflect.Value) error {
	srcType := src.Type()
	dstType := dst.Type()

	// Get type information from cache
	typeCache := cache.GetGlobalCache()

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

	// Iterate through all fields in source struct to find embedded fields
	for i := 0; i < src.NumField(); i++ {
		srcField := srcType.Field(i)

		// Check if it's an embedded field (anonymous field)
		if srcField.Anonymous {
			srcFieldValue := src.Field(i)

			// If it's a pointer type embedded field, dereference it
			if srcFieldValue.Kind() == reflect.Ptr && !srcFieldValue.IsNil() {
				srcFieldValue = srcFieldValue.Elem()
			}

			// Find corresponding embedded field in target struct
			for j := 0; j < dst.NumField(); j++ {
				dstField := dstType.Field(j)

				// If found embedded field of same type
				if dstField.Anonymous && dstField.Type == srcField.Type {
					dstFieldValue := dst.Field(j)

					// Recursively map embedded field
					if err := MapValue(srcFieldValue, dstFieldValue); err != nil {
						return err
					}
					break
				}
			}
		}
	}

	return nil
}

// findFieldInEmbedded finds a field with the specified name in embedded fields
// This function supports mapping fields accessed through embedded fields
func findFieldInEmbedded(value reflect.Value, fieldName string) (reflect.Value, bool) {
	valueType := value.Type()

	// Get type information from cache
	typeCache := cache.GetGlobalCache()
	typeInfo := typeCache.Get(valueType)

	if typeInfo == nil {
		// If not in cache, build and store it
		typeInfo = cache.BuildTypeInfo(valueType)
		typeCache.Store(valueType, typeInfo)
	}

	// Iterate through all fields
	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)

		// If it's an embedded field
		if field.Anonymous {
			fieldValue := value.Field(i)

			// Handle pointer type embedded fields
			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					continue
				}
				fieldValue = fieldValue.Elem()
			}

			// Look for target field in embedded field
			if fieldValue.Kind() == reflect.Struct {
				// Use direct field access for better performance
				if targetField := fieldValue.FieldByName(fieldName); targetField.IsValid() {
					return targetField, true
				}

				// Recursively search deeper embedded fields
				if found, ok := findFieldInEmbedded(fieldValue, fieldName); ok {
					return found, true
				}
			}
		}
	}

	return reflect.Value{}, false
}
