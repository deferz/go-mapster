package mapster

import (
	"reflect"
)

// reflectionMap performs mapping using reflection
func reflectionMap[T any](src any) T {
	var target T

	srcValue := reflect.ValueOf(src)
	targetValue := reflect.ValueOf(&target).Elem()

	// Check for custom mapping configuration
	srcType := srcValue.Type()
	targetType := targetValue.Type()

	if config := GetMappingConfig(srcType, targetType); config != nil {
		mapWithConfig(srcValue, targetValue, config)
	} else {
		mapReflect(srcValue, targetValue)
	}

	return target
}

// mapWithConfig performs mapping using custom configuration
func mapWithConfig(srcValue, targetValue reflect.Value, config *MappingDefinition) {
	if !srcValue.IsValid() || !targetValue.IsValid() {
		return
	}

	srcType := srcValue.Type()
	targetType := targetValue.Type()

	if srcType.Kind() != reflect.Struct || targetType.Kind() != reflect.Struct {
		// For non-struct types, fall back to regular mapping
		mapReflect(srcValue, targetValue)
		return
	}

	// Create a map of source field names to field indices
	srcFields := make(map[string]reflect.Value)
	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		if field.IsExported() {
			srcFields[field.Name] = srcValue.Field(i)
		}
	}

	// Map fields according to configuration
	for i := 0; i < targetType.NumField(); i++ {
		targetField := targetType.Field(i)
		if !targetField.IsExported() {
			continue
		}

		targetFieldValue := targetValue.Field(i)
		if !targetFieldValue.CanSet() {
			continue
		}

		fieldName := targetField.Name

		// Check if there's a custom mapping for this field
		if mapping, exists := config.FieldMappings[fieldName]; exists {
			switch mapping.MappingType {
			case MappingTypeIgnore:
				// Skip this field
				continue

			case MappingTypeField:
				// Map from specified source field
				if srcFieldValue, exists := srcFields[mapping.SourceField]; exists {
					mapReflect(srcFieldValue, targetFieldValue)
				}

			case MappingTypeFunc:
				// Use custom mapping function
				if mapping.CustomFunc != nil {
					result := callCustomFunc(mapping.CustomFunc, srcValue.Interface())
					if result.IsValid() {
						// Try to convert the result to the target field type
						if result.Type().AssignableTo(targetFieldValue.Type()) {
							targetFieldValue.Set(result)
						} else if result.Type().ConvertibleTo(targetFieldValue.Type()) {
							targetFieldValue.Set(result.Convert(targetFieldValue.Type()))
						} else {
							// For interface{} results, try to extract the underlying value
							if result.Kind() == reflect.Interface && !result.IsNil() {
								actualValue := result.Elem()
								if actualValue.Type().AssignableTo(targetFieldValue.Type()) {
									targetFieldValue.Set(actualValue)
								} else if actualValue.Type().ConvertibleTo(targetFieldValue.Type()) {
									targetFieldValue.Set(actualValue.Convert(targetFieldValue.Type()))
								}
							}
						}
					}
				}

			case MappingTypeTransform:
				// First map from source field, then apply transform
				sourceFieldName := mapping.SourceField
				if sourceFieldName == "" {
					sourceFieldName = fieldName // Default to same field name
				}

				if srcFieldValue, exists := srcFields[sourceFieldName]; exists {
					// Apply transform if specified
					if mapping.Transform != nil {
						result := callTransformFunc(mapping.Transform, srcFieldValue.Interface())
						if result.IsValid() && result.Type().AssignableTo(targetFieldValue.Type()) {
							targetFieldValue.Set(result)
						}
					} else {
						mapReflect(srcFieldValue, targetFieldValue)
					}
				}
			}
		} else {
			// No custom mapping, use default field mapping
			if srcFieldValue, exists := srcFields[fieldName]; exists {
				mapReflect(srcFieldValue, targetFieldValue)
			}
		}
	}
}

// callCustomFunc calls a custom mapping function
func callCustomFunc(fn interface{}, src interface{}) reflect.Value {
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return reflect.Value{}
	}

	fnType := fnValue.Type()
	if fnType.NumIn() != 1 || fnType.NumOut() != 1 {
		return reflect.Value{}
	}

	srcValue := reflect.ValueOf(src)
	if !srcValue.Type().AssignableTo(fnType.In(0)) {
		return reflect.Value{}
	}

	results := fnValue.Call([]reflect.Value{srcValue})
	if len(results) > 0 {
		return results[0]
	}

	return reflect.Value{}
}

// callTransformFunc calls a transform function
func callTransformFunc(fn interface{}, src interface{}) reflect.Value {
	fnValue := reflect.ValueOf(fn)
	if fnValue.Kind() != reflect.Func {
		return reflect.Value{}
	}

	fnType := fnValue.Type()
	if fnType.NumIn() != 1 || fnType.NumOut() != 1 {
		return reflect.Value{}
	}

	srcValue := reflect.ValueOf(src)
	if !srcValue.Type().AssignableTo(fnType.In(0)) {
		return reflect.Value{}
	}

	results := fnValue.Call([]reflect.Value{srcValue})
	if len(results) > 0 {
		return results[0]
	}

	return reflect.Value{}
}
func mapReflect(srcValue, targetValue reflect.Value) {
	if !srcValue.IsValid() || !targetValue.IsValid() {
		return
	}

	srcType := srcValue.Type()
	targetType := targetValue.Type()

	// Handle different type combinations
	switch {
	case srcType == targetType:
		// Same type, direct assignment
		targetValue.Set(srcValue)

	case srcType.Kind() == reflect.Struct && targetType.Kind() == reflect.Struct:
		// Struct to struct mapping
		mapStructToStruct(srcValue, targetValue)

	case srcType.Kind() == reflect.Ptr && targetType.Kind() == reflect.Ptr:
		// Pointer to pointer mapping
		if !srcValue.IsNil() {
			targetValue.Set(reflect.New(targetType.Elem()))
			mapReflect(srcValue.Elem(), targetValue.Elem())
		}

	case srcType.Kind() == reflect.Ptr && targetType.Kind() == reflect.Struct:
		// Pointer to struct mapping
		if !srcValue.IsNil() {
			mapReflect(srcValue.Elem(), targetValue)
		}

	case srcType.Kind() == reflect.Struct && targetType.Kind() == reflect.Ptr:
		// Struct to pointer mapping
		targetValue.Set(reflect.New(targetType.Elem()))
		mapReflect(srcValue, targetValue.Elem())

	case srcType.Kind() == reflect.Slice && targetType.Kind() == reflect.Slice:
		// Slice to slice mapping
		mapSliceToSlice(srcValue, targetValue)

	default:
		// Try direct assignment for compatible types
		if srcValue.Type().ConvertibleTo(targetType) {
			targetValue.Set(srcValue.Convert(targetType))
		}
	}
}

// mapStructToStruct maps fields from source struct to target struct
func mapStructToStruct(srcValue, targetValue reflect.Value) {
	srcType := srcValue.Type()
	targetType := targetValue.Type()

	// Create a map of source field names to field indices for faster lookup
	srcFields := make(map[string]int)
	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		if field.IsExported() {
			srcFields[field.Name] = i
		}
	}

	// Map fields by name
	for i := 0; i < targetType.NumField(); i++ {
		targetField := targetType.Field(i)
		if !targetField.IsExported() {
			continue
		}

		if srcFieldIndex, exists := srcFields[targetField.Name]; exists {
			srcFieldValue := srcValue.Field(srcFieldIndex)
			targetFieldValue := targetValue.Field(i)

			if targetFieldValue.CanSet() {
				mapReflect(srcFieldValue, targetFieldValue)
			}
		}
	}
}

// mapSliceToSlice maps elements from source slice to target slice
func mapSliceToSlice(srcValue, targetValue reflect.Value) {
	srcLen := srcValue.Len()
	targetSlice := reflect.MakeSlice(targetValue.Type(), srcLen, srcLen)

	for i := 0; i < srcLen; i++ {
		srcElem := srcValue.Index(i)
		targetElem := targetSlice.Index(i)
		mapReflect(srcElem, targetElem)
	}

	targetValue.Set(targetSlice)
}
