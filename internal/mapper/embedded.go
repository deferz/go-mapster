package mapper

import (
	"reflect"

	"github.com/deferz/go-mapster/internal/cache"
)

// processEmbeddedFields processes embedded fields in structs
// Embedded fields are anonymous fields in Go and require special handling
func processEmbeddedFields(src, dst reflect.Value) error {
	srcType := src.Type()
	dstType := dst.Type()

	// Get type information from cache
	// 使用 GetOrCreate 方法获取或创建类型信息
	// 注意：我们不需要显式地使用类型信息，因为在 MapValue 中已经处理了类型信息
	_ = cache.GetGlobalCache().GetOrCreate(srcType)
	_ = cache.GetGlobalCache().GetOrCreate(dstType)

	// Iterate through all fields in source struct to find embedded fields
	for i := 0; i < src.NumField(); i++ {
		srcField := srcType.Field(i)

		// Check if it's an embedded field (anonymous field)
		if srcField.Anonymous {
			srcFieldValue := src.Field(i)

			// If it's a pointer type embedded field, dereference it
			if srcFieldValue.Kind() == reflect.Ptr && !srcFieldValue.IsNil() {
				srcFieldValue = srcFieldValue.Elem()
			}

			// Find corresponding embedded field in target struct
			for j := 0; j < dst.NumField(); j++ {
				dstField := dstType.Field(j)

				// If found embedded field of same type
				if dstField.Anonymous && dstField.Type == srcField.Type {
					dstFieldValue := dst.Field(j)

					// Recursively map embedded field
					if err := MapValue(srcFieldValue, dstFieldValue); err != nil {
						return err
					}
					break
				}
			}
		}
	}

	return nil
}

// findFieldInEmbedded finds a field with the specified name in embedded fields
// This function supports mapping fields accessed through embedded fields
func findFieldInEmbedded(value reflect.Value, fieldName string) (reflect.Value, bool) {
	valueType := value.Type()

	// Get type information from cache
	typeCache := cache.GetGlobalCache()
	typeInfo := typeCache.GetOrCreate(valueType)

	// 使用缓存的嵌入字段映射
	if embeddedFieldInfo, exists := typeInfo.EmbeddedFieldsMap[fieldName]; exists {
		// 根据缓存的路径获取字段值
		fieldValue := value

		// 遍历嵌入路径
		for _, idx := range embeddedFieldInfo.EmbeddedPath {
			// 处理指针类型
			if fieldValue.Kind() == reflect.Ptr {
				// 如果是 nil 指针，初始化它
				if fieldValue.IsNil() {
					fieldValue.Set(reflect.New(fieldValue.Type().Elem()))
				}
				fieldValue = fieldValue.Elem()
			}

			// 获取下一级字段
			fieldValue = fieldValue.Field(idx)
		}

		return fieldValue, true
	}

	// 如果在缓存中没有找到，则使用递归方式查找
	// 这是一个后备方案，大多数情况下应该使用缓存
	for _, anonField := range typeInfo.AnonymousFields {
		// 获取匿名字段的值
		anonValue := value

		// 遍历嵌入路径
		for _, idx := range anonField.Index {
			// 处理指针类型
			if anonValue.Kind() == reflect.Ptr {
				if anonValue.IsNil() {
					// 如果是 nil 指针，跳过这个匿名字段
					anonValue = reflect.Value{}
					break
				}
				anonValue = anonValue.Elem()
			}

			if !anonValue.IsValid() {
				break
			}

			anonValue = anonValue.Field(idx)
		}

		// 如果匿名字段有效
		if anonValue.IsValid() {
			// 处理指针类型
			if anonValue.Kind() == reflect.Ptr {
				if anonValue.IsNil() {
					continue
				}
				anonValue = anonValue.Elem()
			}

			// 在匿名字段中查找目标字段
			if anonValue.Kind() == reflect.Struct {
				// 直接使用字段名称访问以提高性能
				if targetField := anonValue.FieldByName(fieldName); targetField.IsValid() {
					return targetField, true
				}
			}
		}
	}

	// 作为最后的手段，使用传统方式遍历所有字段
	for i := 0; i < value.NumField(); i++ {
		field := valueType.Field(i)

		// 如果是匿名字段
		if field.Anonymous {
			fieldValue := value.Field(i)

			// 处理指针类型的匿名字段
			if fieldValue.Kind() == reflect.Ptr {
				if fieldValue.IsNil() {
					continue
				}
				fieldValue = fieldValue.Elem()
			}

			// 在匿名字段中查找目标字段
			if fieldValue.Kind() == reflect.Struct {
				// 使用直接字段访问以提高性能
				if targetField := fieldValue.FieldByName(fieldName); targetField.IsValid() {
					return targetField, true
				}

				// 递归搜索更深层的嵌入字段
				if found, ok := findFieldInEmbedded(fieldValue, fieldName); ok {
					return found, true
				}
			}
		}
	}

	return reflect.Value{}, false
}
