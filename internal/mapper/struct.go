package mapper

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/deferz/go-mapster/internal/cache"
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

	// Get source and target type information from cache
	typeCache := cache.GetGlobalCache()
	srcType := src.Type()
	dstType := dst.Type()

	// 使用 GetOrCreate 方法获取或创建类型信息
	srcTypeInfo := typeCache.GetOrCreate(srcType)
	dstTypeInfo := typeCache.GetOrCreate(dstType)

	// Use cached field information for target struct
	for _, fieldInfo := range dstTypeInfo.Fields {
		// Get target field
		dstField := dst.Field(fieldInfo.Index)

		// Field is already verified as exported in BuildTypeInfo
		// Get field name
		fieldName := fieldInfo.Name

		// Find corresponding field in source struct using cached field map
		srcFieldInfo, exists := srcTypeInfo.FieldsMap[fieldName]
		if exists {
			srcField := src.Field(srcFieldInfo.Index)
			// Recursively map field value
			if err := MapValue(srcField, dstField); err != nil {
				return fmt.Errorf("failed to map field %s: %w", fieldName, err)
			}
		} else {
			// Try to find in embedded fields
			if embeddedField, found := findFieldInEmbedded(src, fieldName); found {
				if err := MapValue(embeddedField, dstField); err != nil {
					return fmt.Errorf("failed to map field %s from embedded: %w", fieldName, err)
				}
			} else {
				// Try to find in nested fields with flattening
				if nestedField, found := findNestedField(src, fieldName); found {
					if err := MapValue(nestedField, dstField); err != nil {
						return fmt.Errorf("failed to map nested field %s: %w", fieldName, err)
					}
				}
				// If field not found, skip (keep original value in target field)
			}
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

// findNestedField finds a field in nested structs using dot notation or prefix matching
// This function supports flattening of nested structures
func findNestedField(src reflect.Value, fieldName string) (reflect.Value, bool) {
	// 首先检查类型缓存中是否有嵌套字段映射
	typeCache := cache.GetGlobalCache()
	srcType := src.Type()
	srcTypeInfo := typeCache.GetOrCreate(srcType)

	// 使用缓存的嵌套字段映射
	if nestedFieldInfo, exists := srcTypeInfo.NestedFieldsMap[fieldName]; exists {
		// 根据缓存的路径获取字段值
		fieldValue := src

		// 遍历嵌套路径
		for i, idx := range nestedFieldInfo.IndexPath {
			// 处理指针类型
			if fieldValue.Kind() == reflect.Ptr {
				// 如果是 nil 指针，初始化它
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			// 获取下一级字段
			if i < len(nestedFieldInfo.IndexPath)-1 {
				fieldValue = fieldValue.Field(idx)
			} else {
				// 最后一个字段，直接返回
				return fieldValue.Field(idx), true
			}
		}

		return reflect.Value{}, false // 不应该到达这里
	}

	// 如果在缓存中没有找到，尝试手动查找
	// 检查是否是带前缀的字段名（用于扁平化映射）
	// 例如：Level2_Value2 应该映射到 src.Level2.Value2
	parts := strings.Split(fieldName, "_")
	if len(parts) > 1 {
		// 尝试找到嵌套路径
		return findNestedFieldByPath(src, parts)
	}

	// 尝试使用点号表示法查找嵌套字段
	// 例如：Level2.Level3.Value3
	dotParts := strings.Split(fieldName, ".")
	if len(dotParts) > 1 {
		return findNestedFieldByPath(src, dotParts)
	}

	// 如果没有明确的分隔符，尝试在所有嵌套结构体中查找该字段
	return findFieldInAllNestedStructs(src, fieldName)
}

// findNestedFieldByPath 根据路径查找嵌套字段
func findNestedFieldByPath(src reflect.Value, pathParts []string) (reflect.Value, bool) {
	current := src

	// 遍历路径的每一部分，除了最后一个（字段名）
	for i := 0; i < len(pathParts)-1; i++ {
		// 获取当前路径部分
		pathPart := pathParts[i]

		// 确保当前值是结构体
		if current.Kind() != reflect.Struct {
			return reflect.Value{}, false
		}

		// 尝试获取下一级字段
		field := current.FieldByName(pathPart)
		if !field.IsValid() {
			return reflect.Value{}, false
		}

		// 处理指针类型
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				return reflect.Value{}, false
			}
			field = field.Elem()
		}

		current = field
	}

	// 获取最终字段
	if current.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	// 获取最后一部分作为字段名
	finalField := current.FieldByName(pathParts[len(pathParts)-1])
	if !finalField.IsValid() {
		return reflect.Value{}, false
	}

	return finalField, true
}

// findFieldInAllNestedStructs 在所有嵌套结构体中查找指定字段
func findFieldInAllNestedStructs(src reflect.Value, fieldName string) (reflect.Value, bool) {
	// 确保源是结构体
	if src.Kind() != reflect.Struct {
		return reflect.Value{}, false
	}

	// 首先在当前结构体中查找
	field := src.FieldByName(fieldName)
	if field.IsValid() {
		return field, true
	}

	// 遍历所有字段，查找嵌套结构体
	for i := 0; i < src.NumField(); i++ {
		field := src.Field(i)
		fieldType := src.Type().Field(i)

		// 跳过未导出字段
		if fieldType.PkgPath != "" {
			continue
		}

		// 处理指针类型
		if field.Kind() == reflect.Ptr {
			if field.IsNil() {
				continue
			}
			field = field.Elem()
		}

		// 如果是结构体，递归查找
		if field.Kind() == reflect.Struct {
			// 在嵌套结构体中查找
			if nestedField, found := findFieldInAllNestedStructs(field, fieldName); found {
				return nestedField, true
			}
		}
	}

	return reflect.Value{}, false
}
