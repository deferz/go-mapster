package cache

import (
	"fmt"
	"reflect"
	"sync"
)

// MappingKey represents a source type to target type mapping key
type MappingKey struct {
	SourceType reflect.Type
	TargetType reflect.Type
}

// MappingInfo stores information about a registered type mapping
type MappingInfo struct {
	SourceType   reflect.Type
	TargetType   reflect.Type
	IsRegistered bool
}

// MappingCache provides caching for type mapping information
type MappingCache struct {
	mappings map[MappingKey]*MappingInfo
	mutex    sync.RWMutex
}

// NewMappingCache creates a new MappingCache instance
func NewMappingCache() *MappingCache {
	return &MappingCache{
		mappings: make(map[MappingKey]*MappingInfo),
	}
}

// Get retrieves cached mapping information
// Returns the mapping info if found, or an error if not registered
func (mc *MappingCache) Get(sourceType, targetType reflect.Type) (*MappingInfo, error) {
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	key := MappingKey{SourceType: sourceType, TargetType: targetType}
	info, exists := mc.mappings[key]

	if !exists || !info.IsRegistered {
		return nil, fmt.Errorf("no mapping registered from %s to %s", sourceType, targetType)
	}

	return info, nil
}

// Register registers a new type mapping
func (mc *MappingCache) Register(sourceType, targetType reflect.Type) {
	mc.mutex.Lock()
	defer mc.mutex.Unlock()

	key := MappingKey{SourceType: sourceType, TargetType: targetType}
	mc.mappings[key] = &MappingInfo{
		SourceType:   sourceType,
		TargetType:   targetType,
		IsRegistered: true,
	}
}

// IsRegistered checks if a mapping is registered
func (mc *MappingCache) IsRegistered(sourceType, targetType reflect.Type) bool {
	// Skip nil check
	if sourceType == nil || targetType == nil {
		return false
	}

	// Check direct registration first
	mc.mutex.RLock()
	defer mc.mutex.RUnlock()

	key := MappingKey{SourceType: sourceType, TargetType: targetType}
	info, exists := mc.mappings[key]
	if exists && info.IsRegistered {
		return true
	}

	// For slice and array types, check if their element types are registered
	if (sourceType.Kind() == reflect.Slice || sourceType.Kind() == reflect.Array) &&
		(targetType.Kind() == reflect.Slice || targetType.Kind() == reflect.Array) {

		srcElemType := sourceType.Elem()
		dstElemType := targetType.Elem()

		// Handle pointer element types
		if srcElemType.Kind() == reflect.Ptr && dstElemType.Kind() == reflect.Ptr {
			srcElemType = srcElemType.Elem()
			dstElemType = dstElemType.Elem()
		}

		// Check if element types are registered
		elemKey := MappingKey{SourceType: srcElemType, TargetType: dstElemType}
		elemInfo, elemExists := mc.mappings[elemKey]
		return elemExists && elemInfo.IsRegistered
	}

	// For map types, check if their value types are registered
	if sourceType.Kind() == reflect.Map && targetType.Kind() == reflect.Map {
		srcValueType := sourceType.Elem()
		dstValueType := targetType.Elem()

		// Handle pointer element types
		if srcValueType.Kind() == reflect.Ptr && dstValueType.Kind() == reflect.Ptr {
			srcValueType = srcValueType.Elem()
			dstValueType = dstValueType.Elem()
		}

		// Check if value types are registered
		valueKey := MappingKey{SourceType: srcValueType, TargetType: dstValueType}
		valueInfo, valueExists := mc.mappings[valueKey]
		return valueExists && valueInfo.IsRegistered
	}

	return false
}

// Global shared mapping cache instance
var globalMappingCache = NewMappingCache()

// GetGlobalMappingCache returns the global mapping cache instance
func GetGlobalMappingCache() *MappingCache {
	return globalMappingCache
}
