package cache

import (
	"reflect"
	"sync"
)

// AnonymousFieldInfo stores information about an anonymous (embedded) field
type AnonymousFieldInfo struct {
	Type      reflect.Type // 匹名字段的类型
	Index     []int        // 字段索引路径
	IsPointer bool         // 是否是指针类型
	Name      string       // 字段名称
}

// TypeInfo stores cached reflection information about a type
type TypeInfo struct {
	Type         reflect.Type
	Fields       []FieldInfo
	FieldsMap    map[string]FieldInfo
	IsStruct     bool
	IsCollection bool
	IsMap        bool
	// 存储此类型可以映射到的目标类型
	MappableTargets map[reflect.Type]bool
	// 存储可以映射到此类型的源类型
	MappableSources map[reflect.Type]bool
	// 匹名字段信息
	AnonymousFields []AnonymousFieldInfo
	// 匹名字段中的字段映射
	EmbeddedFieldsMap map[string]EmbeddedFieldInfo
	// 嵌套字段映射（用于扁平化）
	NestedFieldsMap map[string]NestedFieldInfo
}

// FieldInfo stores cached reflection information about a struct field
type FieldInfo struct {
	Name        string
	Index       int
	Type        reflect.Type
	IsStruct    bool
	IsPointer   bool
	IsSlice     bool
	IsMap       bool
	IsAnonymous bool // 是否是匹名字段
}

// EmbeddedFieldInfo stores information about a field in an embedded struct
type EmbeddedFieldInfo struct {
	Field        FieldInfo    // 字段信息
	EmbeddedPath []int        // 字段在嵌套结构中的路径
	ParentType   reflect.Type // 父结构体类型
}

// NestedFieldInfo stores information about a field in a nested struct
type NestedFieldInfo struct {
	Field      FieldInfo    // 字段信息
	NestedPath []string     // 字段在嵌套结构中的路径（字段名路径）
	IndexPath  []int        // 字段在嵌套结构中的索引路径
	ParentType reflect.Type // 父结构体类型
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

// GetOrCreate retrieves cached type information or creates it if not found
func (tc *TypeCache) GetOrCreate(t reflect.Type) *TypeInfo {
	tc.mutex.RLock()
	info, exists := tc.cache[t]
	tc.mutex.RUnlock()

	if exists {
		return info
	}

	// Create new type info
	info = BuildTypeInfo(t)

	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	// Double-check locking
	if existingInfo, exists := tc.cache[t]; exists {
		return existingInfo
	}

	tc.cache[t] = info
	return info
}

// RegisterMapping registers a mapping from source type to target type
func (tc *TypeCache) RegisterMapping(sourceType, targetType reflect.Type) {
	// Get or create source type info
	sourceInfo := tc.GetOrCreate(sourceType)

	// Get or create target type info
	targetInfo := tc.GetOrCreate(targetType)

	tc.mutex.Lock()
	defer tc.mutex.Unlock()

	// Update source type's mappable targets
	sourceInfo.MappableTargets[targetType] = true

	// Update target type's mappable sources
	targetInfo.MappableSources[sourceType] = true
}

// IsRegistered checks if a mapping from source type to target type is registered
func (tc *TypeCache) IsRegistered(sourceType, targetType reflect.Type) bool {
	// Skip nil check
	if sourceType == nil || targetType == nil {
		return false
	}

	tc.mutex.RLock()
	defer tc.mutex.RUnlock()

	// Check if source type is cached
	sourceInfo, sourceExists := tc.cache[sourceType]
	if !sourceExists || sourceInfo.MappableTargets == nil {
		return false
	}

	// Check direct mapping
	if sourceInfo.MappableTargets[targetType] {
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
		elemSourceInfo, elemSourceExists := tc.cache[srcElemType]
		if elemSourceExists && elemSourceInfo.MappableTargets != nil && elemSourceInfo.MappableTargets[dstElemType] {
			return true
		}
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
		valueSourceInfo, valueSourceExists := tc.cache[srcValueType]
		if valueSourceExists && valueSourceInfo.MappableTargets != nil && valueSourceInfo.MappableTargets[dstValueType] {
			return true
		}
	}

	return false
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
		Type:              t,
		IsStruct:          t.Kind() == reflect.Struct,
		IsCollection:      t.Kind() == reflect.Slice || t.Kind() == reflect.Array,
		IsMap:             t.Kind() == reflect.Map,
		MappableTargets:   make(map[reflect.Type]bool),
		MappableSources:   make(map[reflect.Type]bool),
		EmbeddedFieldsMap: make(map[string]EmbeddedFieldInfo),
		NestedFieldsMap:   make(map[string]NestedFieldInfo),
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
				Name:        field.Name,
				Index:       i,
				Type:        fieldType,
				IsStruct:    fieldType.Kind() == reflect.Struct,
				IsPointer:   fieldType.Kind() == reflect.Ptr,
				IsSlice:     fieldType.Kind() == reflect.Slice,
				IsMap:       fieldType.Kind() == reflect.Map,
				IsAnonymous: field.Anonymous,
			}

			info.Fields = append(info.Fields, fieldInfo)
			info.FieldsMap[field.Name] = fieldInfo

			// If this is an anonymous field, collect its fields
			if field.Anonymous {
				// 获取实际类型（处理指针类型）
				actualType := fieldType
				isPointer := fieldType.Kind() == reflect.Ptr
				if isPointer {
					actualType = fieldType.Elem()
				}

				// 只处理结构体类型的匿名字段
				if actualType.Kind() == reflect.Struct {
					// 添加到匿名字段列表
					anonInfo := AnonymousFieldInfo{
						Type:      actualType,
						Index:     []int{i},
						IsPointer: isPointer,
						Name:      field.Name,
					}
					info.AnonymousFields = append(info.AnonymousFields, anonInfo)

					// 收集匿名字段中的字段
					collectEmbeddedFields(info, actualType, []int{i}, isPointer)
				}
			}

			// 如果是结构体类型的字段（非匿名），收集嵌套字段
			if !field.Anonymous && (fieldType.Kind() == reflect.Struct || 
				(fieldType.Kind() == reflect.Ptr && fieldType.Elem().Kind() == reflect.Struct)) {
				actualType := fieldType
				isPointer := fieldType.Kind() == reflect.Ptr
				if isPointer {
					actualType = fieldType.Elem()
				}

				// 收集嵌套字段（用于扁平化映射）
				collectNestedFields(info, actualType, []string{field.Name}, []int{i}, isPointer)
			}
		}
	}

	return info
}

// collectEmbeddedFields 收集嵌入结构体中的字段
func collectEmbeddedFields(info *TypeInfo, embeddedType reflect.Type, path []int, parentIsPointer bool) {
	// 确保是结构体类型
	if embeddedType.Kind() != reflect.Struct {
		return
	}

	// 遍历嵌入结构体的所有字段
	for i := 0; i < embeddedType.NumField(); i++ {
		field := embeddedType.Field(i)

		// 跳过未导出字段
		if field.PkgPath != "" {
			continue
		}

		fieldType := field.Type
		fieldInfo := FieldInfo{
			Name:        field.Name,
			Index:       i,
			Type:        fieldType,
			IsStruct:    fieldType.Kind() == reflect.Struct,
			IsPointer:   fieldType.Kind() == reflect.Ptr,
			IsSlice:     fieldType.Kind() == reflect.Slice,
			IsMap:       fieldType.Kind() == reflect.Map,
			IsAnonymous: field.Anonymous,
		}

		// 创建字段路径
		fieldPath := append([]int{}, path...)
		fieldPath = append(fieldPath, i)

		// 将字段添加到嵌入字段映射中
		embeddedFieldInfo := EmbeddedFieldInfo{
			Field:        fieldInfo,
			EmbeddedPath: fieldPath,
			ParentType:   embeddedType,
		}

		// 如果字段名称不存在于主结构体中，或者已经是来自另一个嵌入字段
		if _, exists := info.FieldsMap[field.Name]; !exists {
			if _, existsInEmbedded := info.EmbeddedFieldsMap[field.Name]; !existsInEmbedded {
				info.EmbeddedFieldsMap[field.Name] = embeddedFieldInfo
			}
		}

		// 如果是匿名字段，继续递归收集
		if field.Anonymous {
			anonType := fieldType
			isPointer := fieldType.Kind() == reflect.Ptr
			if isPointer {
				anonType = fieldType.Elem()
			}

			if anonType.Kind() == reflect.Struct {
				// 添加到匿名字段列表
				anonInfo := AnonymousFieldInfo{
					Type:      anonType,
					Index:     fieldPath,
					IsPointer: isPointer || parentIsPointer,
					Name:      field.Name,
				}
				info.AnonymousFields = append(info.AnonymousFields, anonInfo)

				// 递归收集
				collectEmbeddedFields(info, anonType, fieldPath, isPointer || parentIsPointer)
			}
		}
	}
}

// Global shared type cache instance
var globalTypeCache = NewTypeCache()

// GetGlobalCache returns the global type cache instance
func GetGlobalCache() *TypeCache {
	return globalTypeCache
}

// collectNestedFields 收集嵌套结构体中的字段（用于扁平化映射）
func collectNestedFields(info *TypeInfo, nestedType reflect.Type, path []string, indexPath []int, parentIsPointer bool) {
	// 确保是结构体类型
	if nestedType.Kind() != reflect.Struct {
		return
	}

	// 遍历嵌套结构体的所有字段
	for i := 0; i < nestedType.NumField(); i++ {
		field := nestedType.Field(i)

		// 跳过未导出字段
		if field.PkgPath != "" {
			continue
		}

		fieldType := field.Type
		fieldInfo := FieldInfo{
			Name:        field.Name,
			Index:       i,
			Type:        fieldType,
			IsStruct:    fieldType.Kind() == reflect.Struct,
			IsPointer:   fieldType.Kind() == reflect.Ptr,
			IsSlice:     fieldType.Kind() == reflect.Slice,
			IsMap:       fieldType.Kind() == reflect.Map,
			IsAnonymous: field.Anonymous,
		}

		// 创建字段路径
		fieldPath := append([]string{}, path...)
		fieldPath = append(fieldPath, field.Name)

		// 创建索引路径
		fieldIndexPath := append([]int{}, indexPath...)
		fieldIndexPath = append(fieldIndexPath, i)

		// 生成扁平化字段名
		// 例如: Level2_Level3_Value3
		flattenedName := ""
		for j, part := range fieldPath {
			if j > 0 {
				flattenedName += "_"
			}
			flattenedName += part
		}

		// 将字段添加到嵌套字段映射中
		nestedFieldInfo := NestedFieldInfo{
			Field:      fieldInfo,
			NestedPath: fieldPath,
			IndexPath:  fieldIndexPath,
			ParentType: nestedType,
		}

		// 添加到嵌套字段映射中
		info.NestedFieldsMap[flattenedName] = nestedFieldInfo

		// 如果是结构体类型，继续递归收集
		if fieldType.Kind() == reflect.Struct || 
			(fieldType.Kind() == reflect.Ptr && fieldType.Elem().Kind() == reflect.Struct) {
			actualType := fieldType
			isPointer := fieldType.Kind() == reflect.Ptr
			if isPointer {
				actualType = fieldType.Elem()
			}

			// 递归收集嵌套字段
			collectNestedFields(info, actualType, fieldPath, fieldIndexPath, isPointer || parentIsPointer)
		}
	}
}
