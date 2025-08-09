package mapster

import (
	"fmt"
	"reflect"

	"github.com/deferz/go-mapster/internal/cache"
	"github.com/deferz/go-mapster/internal/config"
	"github.com/deferz/go-mapster/internal/mapper"
)

// Map maps the source object to the target type and returns a new instance.
// This function uses generics to ensure type safety, checking type matches at compile time.
// If the mapping between source and target types is not registered, it returns an error.
func Map[T any](src any) (T, error) {
	var result T
	if src == nil {
		return result, fmt.Errorf("source cannot be nil")
	}

	// Check if mapping is registered
	mappingCache := cache.GetGlobalMappingCache()
	sourceType := reflect.TypeOf(src)
	targetType := reflect.TypeOf(result)

	// Strict check - require registration
	if !mappingCache.IsRegistered(sourceType, targetType) {
		return result, fmt.Errorf("no mapping registered from %s to %s", sourceType, targetType)
	}

	resultPtr := &result
	if err := mapper.MapValue(reflect.ValueOf(src), reflect.ValueOf(resultPtr).Elem()); err != nil {
		return result, fmt.Errorf("mapping failed: %w", err)
	}
	return result, nil
}

// MapTo maps the source object to an existing target object.
// This function modifies the target object in place.
// If the mapping between source and target types is not registered, it returns an error.
// The destination parameter must be a pointer to the target type.
func MapTo[T any](src any, dst *T) error {
	if src == nil {
		return fmt.Errorf("source cannot be nil")
	}

	if dst == nil {
		return fmt.Errorf("destination cannot be nil pointer")
	}

	// Check if mapping is registered
	mappingCache := cache.GetGlobalMappingCache()
	sourceType := reflect.TypeOf(src)
	targetType := reflect.TypeOf(*dst)

	if !mappingCache.IsRegistered(sourceType, targetType) {
		return fmt.Errorf("no mapping registered from %s to %s", sourceType, targetType)
	}

	return mapper.MapValue(reflect.ValueOf(src), reflect.ValueOf(dst).Elem())
}

// NewMapperConfig creates and returns a configuration for mapping between types T and R.
// This function must be called and Register() must be invoked before mapping between these types.
//
// Example:
//
//	mapster.NewMapperConfig[UserEntity, UserDTO]().Register()
func NewMapperConfig[T, R any]() *config.MapperConfig[T, R] {
	return config.NewMapperConfig[T, R]()
}
