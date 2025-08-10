package mapster

import (
	"fmt"
	"reflect"

	"github.com/deferz/go-mapster/internal/cache"
	"github.com/deferz/go-mapster/internal/mapper"
)

// Map maps the source object to the target type and returns a new instance.
// This function uses generics to ensure type safety, checking type matches at compile time.
// Mapping between source and target types is automatically registered on first use.
func Map[T any](src any) (T, error) {
	var result T
	if src == nil {
		return result, fmt.Errorf("source cannot be nil")
	}

	// Get types
	typeCache := cache.GetGlobalCache()
	sourceType := reflect.TypeOf(src)
	targetType := reflect.TypeOf(result)

	// Auto-register the mapping if not already registered
	if !typeCache.IsRegistered(sourceType, targetType) {
		typeCache.RegisterMapping(sourceType, targetType)
	}

	resultPtr := &result
	if err := mapper.MapValue(reflect.ValueOf(src), reflect.ValueOf(resultPtr).Elem()); err != nil {
		return result, fmt.Errorf("mapping failed: %w", err)
	}
	return result, nil
}

// MapTo maps the source object to an existing target object.
// This function modifies the target object in place.
// Mapping between source and target types is automatically registered on first use.
// The destination parameter must be a pointer to the target type.
func MapTo[T any](src any, dst *T) error {
	if src == nil {
		return fmt.Errorf("source cannot be nil")
	}

	if dst == nil {
		return fmt.Errorf("destination cannot be nil pointer")
	}

	// Get types
	typeCache := cache.GetGlobalCache()
	sourceType := reflect.TypeOf(src)
	targetType := reflect.TypeOf(*dst)

	// Auto-register the mapping if not already registered
	if !typeCache.IsRegistered(sourceType, targetType) {
		typeCache.RegisterMapping(sourceType, targetType)
	}

	return mapper.MapValue(reflect.ValueOf(src), reflect.ValueOf(dst).Elem())
}
