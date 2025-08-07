// Package mapster provides high-performance object mapping for Go
package mapster

import (
	"fmt"
	"reflect"
)

// Map performs mapping from source type to target type
func Map[T any](src any) T {
	if src == nil {
		var zero T
		return zero
	}

	srcType := reflect.TypeOf(src)
	targetType := reflect.TypeOf(*new(T))

	// Try to use generated mapper first
	if mapper := getGeneratedMapper(srcType, targetType); mapper != nil {
		// Handle different mapper function signatures
		switch m := mapper.(type) {
		case func(any) any:
			return m(src).(T)
		default:
			// Try to call the generic mapper using reflection
			mapperValue := reflect.ValueOf(mapper)
			if mapperValue.Kind() == reflect.Func {
				results := mapperValue.Call([]reflect.Value{reflect.ValueOf(src)})
				if len(results) > 0 {
					return results[0].Interface().(T)
				}
			}
		}
	}

	// Fallback to reflection mapping
	return reflectionMap[T](src)
}

// MapTo maps source to an existing target object
func MapTo[T any](src any, target *T) {
	if src == nil || target == nil {
		return
	}

	mapped := Map[T](src)
	*target = mapped
}

// generatedMappers stores generated mapping functions
var generatedMappers = make(map[string]any)

// RegisterGeneratedMapper registers a manually written mapping function for zero-reflection performance
// Users can write their own mapping functions and register them for maximum performance
func RegisterGeneratedMapper[S, T any](mapper func(S) T) {
	key := fmt.Sprintf("%T->%T", *new(S), *new(T))
	generatedMappers[key] = mapper
}

// getGeneratedMapper retrieves a generated mapper if available
func getGeneratedMapper(srcType, targetType reflect.Type) any {
	key := fmt.Sprintf("%s->%s", srcType.String(), targetType.String())
	return generatedMappers[key]
}

// ClearGeneratedMappers clears all registered generated mappers
// This is useful for testing and benchmarking
func ClearGeneratedMappers() {
	generatedMappers = make(map[string]any)
}
