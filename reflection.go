package mapster

import (
	"fmt"
	"reflect"
	"strings"
	"time"
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

	// Dereference pointer if necessary
	if srcValue.Kind() == reflect.Ptr {
		if srcValue.IsNil() {
			return
		}
		srcValue = srcValue.Elem()
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
				// Map from specified source field or path
				sourceField := mapping.SourceField

				// Check if it's a path (contains dots)
				if strings.Contains(sourceField, ".") {
					// Use path resolution
					if pathValue, err := getValueByPath(srcValue, sourceField); err == nil && pathValue.IsValid() {
						mapReflect(pathValue, targetFieldValue)
					}
				} else {
					// Simple field mapping
					if srcFieldValue, exists := srcFields[sourceField]; exists {
						mapReflect(srcFieldValue, targetFieldValue)
					}
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
							// For any results, try to extract the underlying value
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
func callCustomFunc(fn any, src any) reflect.Value {
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
func callTransformFunc(fn any, src any) reflect.Value {
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

	// Try type alias conversion first (before struct-to-struct mapping)
	if converted := tryTypeAliasConversion(srcValue, targetType); converted.IsValid() {
		targetValue.Set(converted)
		return
	}

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
		// Try custom type conversions
		if converted := tryCustomTypeConversion(srcValue, targetType); converted.IsValid() {
			targetValue.Set(converted)
			return
		}

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

// getValueByPath extracts a value from a nested structure using dot notation
// e.g., "Address.Street", "Company.Address.City"
func getValueByPath(src reflect.Value, path string) (reflect.Value, error) {
	if path == "" {
		return src, nil
	}

	parts := strings.Split(path, ".")
	current := src

	for i, part := range parts {
		if !current.IsValid() {
			return reflect.Value{}, fmt.Errorf("invalid value at path segment '%s' (index %d)", part, i)
		}

		// Dereference pointers
		if current.Kind() == reflect.Ptr {
			if current.IsNil() {
				return reflect.Value{}, fmt.Errorf("nil pointer at path segment '%s' (index %d)", part, i)
			}
			current = current.Elem()
		}

		// Handle different types
		switch current.Kind() {
		case reflect.Struct:
			// Find field by name
			fieldValue := current.FieldByName(part)
			if !fieldValue.IsValid() {
				return reflect.Value{}, fmt.Errorf("field '%s' not found in struct at path segment %d", part, i)
			}
			current = fieldValue

		case reflect.Map:
			// Handle map access
			key := reflect.ValueOf(part)
			if !key.Type().AssignableTo(current.Type().Key()) {
				return reflect.Value{}, fmt.Errorf("key type mismatch at path segment '%s' (index %d)", part, i)
			}
			mapValue := current.MapIndex(key)
			if !mapValue.IsValid() {
				return reflect.Value{}, fmt.Errorf("key '%s' not found in map at path segment %d", part, i)
			}
			current = mapValue

		case reflect.Interface:
			// Unwrap interface
			if current.IsNil() {
				return reflect.Value{}, fmt.Errorf("nil interface at path segment '%s' (index %d)", part, i)
			}
			current = current.Elem()
			// Recursively process remaining path
			remainingPath := strings.Join(parts[i:], ".")
			return getValueByPath(current, remainingPath)

		default:
			return reflect.Value{}, fmt.Errorf("cannot navigate into type %v at path segment '%s' (index %d)", current.Kind(), part, i)
		}
	}

	return current, nil
}

// CircularReferenceDetector tracks visited objects to prevent infinite recursion
type CircularReferenceDetector struct {
	visited      map[uintptr]bool
	maxDepth     int
	currentDepth int
}

// NewCircularReferenceDetector creates a new detector with specified max depth
func NewCircularReferenceDetector(maxDepth int) *CircularReferenceDetector {
	return &CircularReferenceDetector{
		visited:      make(map[uintptr]bool),
		maxDepth:     maxDepth,
		currentDepth: 0,
	}
}

// Enter checks if we can safely enter an object (returns false if circular reference detected)
func (d *CircularReferenceDetector) Enter(value reflect.Value) bool {
	d.currentDepth++

	// Check max depth
	if d.currentDepth > d.maxDepth {
		return false
	}

	// Only track pointer types and interfaces that could cause cycles
	if value.Kind() != reflect.Ptr && value.Kind() != reflect.Interface {
		return true
	}

	if !value.IsValid() || value.IsNil() {
		return true
	}

	// Get the pointer address
	ptr := value.Pointer()
	if d.visited[ptr] {
		return false // Circular reference detected
	}

	d.visited[ptr] = true
	return true
}

// Exit decrements the depth counter and removes tracking for the object
func (d *CircularReferenceDetector) Exit(value reflect.Value) {
	d.currentDepth--

	if value.Kind() == reflect.Ptr || value.Kind() == reflect.Interface {
		if value.IsValid() && !value.IsNil() {
			ptr := value.Pointer()
			delete(d.visited, ptr)
		}
	}
}

// mapReflectWithCircularDetection performs mapping with circular reference detection
func mapReflectWithCircularDetection(srcValue, targetValue reflect.Value, detector *CircularReferenceDetector) {
	// Check for circular references
	if !detector.Enter(srcValue) {
		// Circular reference detected, skip mapping
		return
	}
	defer detector.Exit(srcValue)

	// Perform regular mapping
	mapReflect(srcValue, targetValue)
}

// tryTypeAliasConversion attempts to convert between type aliases and their underlying types
func tryTypeAliasConversion(srcValue reflect.Value, targetType reflect.Type) reflect.Value {
	if !srcValue.IsValid() {
		return reflect.Value{}
	}

	srcType := srcValue.Type()

	// Handle type alias conversions
	// For type aliases, we need to check if the underlying types are the same
	// and if they are convertible to each other

	// Try direct conversion first (this handles both directions)
	if srcValue.Type().ConvertibleTo(targetType) {
		return srcValue.Convert(targetType)
	}

	// Handle pointer type aliases
	if srcType.Kind() == reflect.Ptr && targetType.Kind() == reflect.Ptr {
		if srcType.Elem().ConvertibleTo(targetType.Elem()) {
			if !srcValue.IsNil() {
				convertedElem := srcValue.Elem().Convert(targetType.Elem())
				result := reflect.New(targetType.Elem())
				result.Elem().Set(convertedElem)
				return result
			}
		}
	}

	return reflect.Value{}
}

// tryCustomTypeConversion attempts to convert between types that are not directly convertible
// but have logical conversion rules (e.g., int64 <-> time.Time)
func tryCustomTypeConversion(srcValue reflect.Value, targetType reflect.Type) reflect.Value {
	if !srcValue.IsValid() {
		return reflect.Value{}
	}

	// Check if time conversion is enabled
	if !globalConfig.EnableTimeConversion {
		return reflect.Value{}
	}

	srcType := srcValue.Type()

	// int64 -> time.Time conversion
	if srcType.Kind() == reflect.Int64 && targetType == reflect.TypeOf(time.Time{}) {
		if srcValue.CanInt() {
			timestamp := srcValue.Int()
			// Handle different timestamp formats
			var convertedTime time.Time
			if timestamp > 1e15 { // Likely nanoseconds (16+ digits)
				convertedTime = time.Unix(0, timestamp)
			} else if timestamp > 1e12 { // Likely milliseconds (13+ digits)
				convertedTime = time.Unix(timestamp/1000, (timestamp%1000)*int64(time.Millisecond))
			} else { // Likely seconds (10-12 digits)
				convertedTime = time.Unix(timestamp, 0)
			}
			return reflect.ValueOf(convertedTime)
		}
	}

	// time.Time -> int64 conversion
	if srcType == reflect.TypeOf(time.Time{}) && targetType.Kind() == reflect.Int64 {
		if timeValue, ok := srcValue.Interface().(time.Time); ok {
			// Convert to Unix timestamp (seconds since epoch)
			timestamp := timeValue.Unix()
			return reflect.ValueOf(timestamp)
		}
	}

	// int -> time.Time conversion
	if srcType.Kind() == reflect.Int && targetType == reflect.TypeOf(time.Time{}) {
		if srcValue.CanInt() {
			timestamp := srcValue.Int()
			convertedTime := time.Unix(timestamp, 0)
			return reflect.ValueOf(convertedTime)
		}
	}

	// time.Time -> int conversion
	if srcType == reflect.TypeOf(time.Time{}) && targetType.Kind() == reflect.Int {
		if timeValue, ok := srcValue.Interface().(time.Time); ok {
			timestamp := int(timeValue.Unix())
			return reflect.ValueOf(timestamp)
		}
	}

	return reflect.Value{}
}
