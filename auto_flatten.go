package mapster

import (
	"reflect"
)

// AutoFlattenConfig 自动扁平化配置
type AutoFlattenConfig struct {
	// 是否启用自动扁平化
	Enabled bool
	// 扁平化深度限制，防止无限递归
	MaxDepth int
	// 字段名冲突时的处理策略
	ConflictStrategy ConflictStrategy
	// 是否包含嵌套结构体的字段名作为前缀
	UsePrefix bool
	// 前缀分隔符
	PrefixSeparator string
}

// ConflictStrategy 字段名冲突处理策略
type ConflictStrategy int

const (
	// KeepFirst 保留第一个遇到的字段（默认）
	KeepFirst ConflictStrategy = iota
	// KeepLast 保留最后一个遇到的字段
	KeepLast
	// UsePrefix 使用前缀区分
	UsePrefix
	// Skip 跳过冲突字段
	Skip
)

// DefaultAutoFlattenConfig 默认自动扁平化配置
func DefaultAutoFlattenConfig() *AutoFlattenConfig {
	return &AutoFlattenConfig{
		Enabled:          true,
		MaxDepth:         3,
		ConflictStrategy: KeepFirst,
		UsePrefix:        false,
		PrefixSeparator:  "_",
	}
}

// 获取全局自动扁平化配置（从全局配置中获取）
func getGlobalAutoFlattenConfig() *AutoFlattenConfig {
	return &globalConfig.AutoFlatten
}

// flattenStructFields 扁平化结构体字段
func flattenStructFields(srcType reflect.Type, config *AutoFlattenConfig) map[string]reflect.StructField {
	if config == nil {
		config = getGlobalAutoFlattenConfig()
	}

	flattened := make(map[string]reflect.StructField)
	flattenStructFieldsRecursive(srcType, "", flattened, config, 0)
	return flattened
}

// flattenStructFieldsRecursive 递归扁平化结构体字段
func flattenStructFieldsRecursive(
	srcType reflect.Type,
	prefix string,
	flattened map[string]reflect.StructField,
	config *AutoFlattenConfig,
	depth int,
) {
	if depth >= config.MaxDepth {
		return
	}

	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldName := field.Name
		if prefix != "" && config.UsePrefix {
			fieldName = prefix + config.PrefixSeparator + fieldName
		}

		// 检查字段类型
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			fieldType = fieldType.Elem()
		}

		// 如果是结构体且不是基本类型，递归扁平化
		if fieldType.Kind() == reflect.Struct && !isBasicType(fieldType) {
			// 递归扁平化嵌套结构体
			flattenStructFieldsRecursive(fieldType, fieldName, flattened, config, depth+1)
		} else {
			// 处理字段名冲突
			if _, exists := flattened[fieldName]; exists {
				switch config.ConflictStrategy {
				case KeepFirst:
					// 保留第一个，跳过当前字段
					continue
				case KeepLast:
					// 保留最后一个，覆盖现有字段
					flattened[fieldName] = field
				case UsePrefix:
					// 使用前缀区分
					prefixedName := field.Name + config.PrefixSeparator + fieldName
					flattened[prefixedName] = field
				case Skip:
					// 跳过冲突字段
					continue
				}
			} else {
				flattened[fieldName] = field
			}
		}
	}
}

// isBasicType 检查是否为基本类型
func isBasicType(t reflect.Type) bool {
	switch t.Kind() {
	case reflect.Bool, reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64, reflect.Complex64, reflect.Complex128,
		reflect.String, reflect.Array, reflect.Slice, reflect.Map, reflect.Interface:
		return true
	case reflect.Struct:
		// 检查是否为 time.Time 等特殊结构体
		if t.PkgPath() == "time" && t.Name() == "Time" {
			return true
		}
		return false
	default:
		return false
	}
}

// mapStructToStructWithAutoFlatten 支持自动扁平化的结构体映射
func mapStructToStructWithAutoFlatten(srcValue, targetValue reflect.Value, config *AutoFlattenConfig) {
	if config == nil {
		config = getGlobalAutoFlattenConfig()
	}

	if !config.Enabled {
		// 如果未启用自动扁平化，使用默认映射
		mapStructToStruct(srcValue, targetValue)
		return
	}

	targetType := targetValue.Type()

	// 创建源字段值映射
	srcFieldValues := make(map[string]reflect.Value)
	buildSrcFieldValues(srcValue, "", srcFieldValues, config, 0)

	// 映射目标字段
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

		// 尝试从扁平化字段中查找匹配
		if srcFieldValue, exists := srcFieldValues[fieldName]; exists {
			mapReflect(srcFieldValue, targetFieldValue)
		} else {
			// 尝试直接字段名匹配（向后兼容）
			if srcFieldValue := srcValue.FieldByName(fieldName); srcFieldValue.IsValid() {
				mapReflect(srcFieldValue, targetFieldValue)
			}
		}
	}
}

// buildSrcFieldValues 构建源字段值映射
func buildSrcFieldValues(
	srcValue reflect.Value,
	prefix string,
	srcFieldValues map[string]reflect.Value,
	config *AutoFlattenConfig,
	depth int,
) {
	if depth >= config.MaxDepth {
		return
	}

	srcType := srcValue.Type()

	for i := 0; i < srcType.NumField(); i++ {
		field := srcType.Field(i)
		if !field.IsExported() {
			continue
		}

		fieldValue := srcValue.Field(i)
		fieldName := field.Name

		if prefix != "" && config.UsePrefix {
			fieldName = prefix + config.PrefixSeparator + fieldName
		}

		// 检查字段类型
		fieldType := field.Type
		if fieldType.Kind() == reflect.Ptr {
			if fieldValue.IsNil() {
				continue
			}
			fieldType = fieldType.Elem()
			fieldValue = fieldValue.Elem()
		}

		// 如果是结构体且不是基本类型，递归处理
		if fieldType.Kind() == reflect.Struct && !isBasicType(fieldType) {
			buildSrcFieldValues(fieldValue, fieldName, srcFieldValues, config, depth+1)
		} else {
			// 处理字段名冲突
			if _, exists := srcFieldValues[fieldName]; exists {
				switch config.ConflictStrategy {
				case KeepFirst:
					// 保留第一个，跳过当前字段
					continue
				case KeepLast:
					// 保留最后一个，覆盖现有字段
					srcFieldValues[fieldName] = fieldValue
				case UsePrefix:
					// 使用前缀区分
					prefixedName := field.Name + config.PrefixSeparator + fieldName
					srcFieldValues[prefixedName] = fieldValue
				case Skip:
					// 跳过冲突字段
					continue
				}
			} else {
				srcFieldValues[fieldName] = fieldValue
			}
		}
	}
}
