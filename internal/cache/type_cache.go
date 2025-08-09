package cache

import (
	"reflect"
	"sync"
)

// TypeInfo stores cached reflection information about a type
type TypeInfo struct {
	Type         reflect.Type
	Fields       []FieldInfo
	FieldsMap    map[string]FieldInfo
	IsStruct     bool
	IsCollection bool
	IsMap        bool
}

// FieldInfo stores cached reflection information about a struct field
type FieldInfo struct {
	Name      string
	Index     int
	Type      reflect.Type
	IsStruct  bool
	IsPointer bool
	IsSlice   bool
	IsMap     bool
}

// TypeCache provides caching for type reflection information
type TypeCache struct {
	cache map[reflect.Type]*TypeInfo
	mutex sync.RWMutex
}

// NewTypeCache creates a new TypeCache instance
func NewTypeCache() *TypeCache {
	return &TypeCache{
		cache: make(map[reflect.Type]*TypeInfo),
	}
}

// Get retrieves cached type information or returns nil if not found
func (tc *TypeCache) Get(t reflect.Type) *TypeInfo {
	tc.mutex.RLock()
	defer tc.mutex.RUnlock()
	
	return tc.cache[t]
}

// Store caches type information
func (tc *TypeCache) Store(t reflect.Type, info *TypeInfo) {
	tc.mutex.Lock()
	defer tc.mutex.Unlock()
	
	tc.cache[t] = info
}

// BuildTypeInfo analyzes a type and builds TypeInfo
func BuildTypeInfo(t reflect.Type) *TypeInfo {
	info := &TypeInfo{
		Type:         t,
		IsStruct:     t.Kind() == reflect.Struct,
		IsCollection: t.Kind() == reflect.Slice || t.Kind() == reflect.Array,
		IsMap:        t.Kind() == reflect.Map,
	}
	
	// For struct types, cache field information
	if info.IsStruct {
		numFields := t.NumField()
		info.Fields = make([]FieldInfo, 0, numFields)
		info.FieldsMap = make(map[string]FieldInfo, numFields)
		
		for i := 0; i < numFields; i++ {
			field := t.Field(i)
			
			// Skip unexported fields
			if field.PkgPath != "" {
				continue
			}
			
			fieldType := field.Type
			fieldInfo := FieldInfo{
				Name:      field.Name,
				Index:     i,
				Type:      fieldType,
				IsStruct:  fieldType.Kind() == reflect.Struct,
				IsPointer: fieldType.Kind() == reflect.Ptr,
				IsSlice:   fieldType.Kind() == reflect.Slice,
				IsMap:     fieldType.Kind() == reflect.Map,
			}
			
			info.Fields = append(info.Fields, fieldInfo)
			info.FieldsMap[field.Name] = fieldInfo
		}
	}
	
	return info
}

// Global shared type cache instance
var globalTypeCache = NewTypeCache()

// GetGlobalCache returns the global type cache instance
func GetGlobalCache() *TypeCache {
	return globalTypeCache
}
