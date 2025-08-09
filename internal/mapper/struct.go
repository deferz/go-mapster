package mapper

import (
	"fmt"
	"reflect"
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

	// Get target struct type information
	dstType := dst.Type()

	// Iterate through all fields of the target struct
	for i := 0; i < dst.NumField(); i++ {
		// Get target field
		dstField := dst.Field(i)
		dstFieldType := dstType.Field(i)

		// Skip unexported fields (private fields)
		if !dstField.CanSet() {
			continue
		}

		// Get field name
		fieldName := dstFieldType.Name

		// Find corresponding field in source struct
		srcField := findSourceField(src, fieldName)
		if !srcField.IsValid() {
			// Try to find in embedded fields
			if embeddedField, found := findFieldInEmbedded(src, fieldName); found {
				srcField = embeddedField
			} else {
				// If field not found, skip (keep original value in target field)
				continue
			}
		}

		// Recursively map field value
		if err := MapValue(srcField, dstField); err != nil {
			return fmt.Errorf("failed to map field %s: %w", fieldName, err)
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
