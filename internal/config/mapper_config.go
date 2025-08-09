package config

import (
	"reflect"

	"github.com/deferz/go-mapster/internal/cache"
)

// MapperConfig stores configuration for mapping between specific types
type MapperConfig[T, R any] struct {
	sourceType reflect.Type
	targetType reflect.Type
}

// NewMapperConfig creates a new MapperConfig instance
func NewMapperConfig[T, R any]() *MapperConfig[T, R] {
	var source T
	var target R

	return &MapperConfig[T, R]{
		sourceType: reflect.TypeOf(source),
		targetType: reflect.TypeOf(target),
	}
}

// Register registers the mapping configuration in the global cache
func (c *MapperConfig[T, R]) Register() {
	mappingCache := cache.GetGlobalMappingCache()
	mappingCache.Register(c.sourceType, c.targetType)
}
