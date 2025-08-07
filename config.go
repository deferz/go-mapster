package mapster

import (
	"fmt"
	"reflect"
)

// MappingType represents different types of field mapping
type MappingType int

const (
	MappingTypeField MappingType = iota
	MappingTypeFunc
	MappingTypeTransform
	MappingTypeIgnore
)

// FieldMapping represents a mapping configuration for a single field
type FieldMapping struct {
	TargetField string
	MappingType MappingType
	SourceField string
	CustomFunc  any // Function for custom mapping
	Transform   any // Transform function
	Condition   any // Condition function for conditional mapping
}

// MappingDefinition represents a complete mapping configuration between two types
type MappingDefinition struct {
	SourceType    reflect.Type
	TargetType    reflect.Type
	FieldMappings map[string]FieldMapping
}

// ConfigBuilder provides a fluent API for configuring mappings
type ConfigBuilder[S, T any] struct {
	definition *MappingDefinition
}

// GlobalConfig holds global configuration options for the mapster library
type GlobalConfig struct {
	// EnableTimeConversion enables automatic conversion between int64 and time.Time
	// Default: true
	EnableTimeConversion bool
}

// globalConfig holds the global configuration
var globalConfig = GlobalConfig{
	EnableTimeConversion: true, // 默认启用时间转换
}

// globalConfigs stores all registered mapping configurations
var globalConfigs = make(map[string]*MappingDefinition)

// Config starts a new configuration chain for specific types
func Config[S, T any]() *ConfigBuilder[S, T] {
	var s S
	var t T

	sourceType := reflect.TypeOf(s)
	targetType := reflect.TypeOf(t)

	definition := &MappingDefinition{
		SourceType:    sourceType,
		TargetType:    targetType,
		FieldMappings: make(map[string]FieldMapping),
	}

	return &ConfigBuilder[S, T]{
		definition: definition,
	}
}

// Map configures mapping for a specific field
func (c *ConfigBuilder[S, T]) Map(targetField string) *FieldConfigBuilder[S, T] {
	return &FieldConfigBuilder[S, T]{
		configBuilder: c,
		targetField:   targetField,
	}
}

// Ignore configures a field to be ignored during mapping
func (c *ConfigBuilder[S, T]) Ignore(targetField string) *ConfigBuilder[S, T] {
	c.definition.FieldMappings[targetField] = FieldMapping{
		TargetField: targetField,
		MappingType: MappingTypeIgnore,
	}
	return c
}

// Register registers the mapping configuration globally
func (c *ConfigBuilder[S, T]) Register() {
	key := getMappingKey(c.definition.SourceType, c.definition.TargetType)
	globalConfigs[key] = c.definition
}

// FieldConfigBuilder provides field-specific configuration
type FieldConfigBuilder[S, T any] struct {
	configBuilder *ConfigBuilder[S, T]
	targetField   string
}

// FromField maps from a source field with the same or different name
func (f *FieldConfigBuilder[S, T]) FromField(sourceField string) *ConfigBuilder[S, T] {
	f.configBuilder.definition.FieldMappings[f.targetField] = FieldMapping{
		TargetField: f.targetField,
		MappingType: MappingTypeField,
		SourceField: sourceField,
	}
	return f.configBuilder
}

// FromFunc maps using a custom function
func (f *FieldConfigBuilder[S, T]) FromFunc(mapperFunc func(S) any) *ConfigBuilder[S, T] {
	f.configBuilder.definition.FieldMappings[f.targetField] = FieldMapping{
		TargetField: f.targetField,
		MappingType: MappingTypeFunc,
		CustomFunc:  mapperFunc,
	}
	return f.configBuilder
}

// FromPath maps from a nested field using dot notation (e.g., "Customer.Name")
func (f *FieldConfigBuilder[S, T]) FromPath(path string) *ConfigBuilder[S, T] {
	// For now, treat path mapping as field mapping
	// TODO: Implement proper nested field mapping
	f.configBuilder.definition.FieldMappings[f.targetField] = FieldMapping{
		TargetField: f.targetField,
		MappingType: MappingTypeField,
		SourceField: path,
	}
	return f.configBuilder
}

// Transform applies a transformation function to the mapped value
func (f *FieldConfigBuilder[S, T]) Transform(transformFunc any) *ConfigBuilder[S, T] {
	// Get existing mapping or create a new one
	mapping, exists := f.configBuilder.definition.FieldMappings[f.targetField]
	if !exists {
		mapping = FieldMapping{
			TargetField: f.targetField,
			MappingType: MappingTypeField,
			SourceField: f.targetField, // Default to same field name
		}
	}

	mapping.Transform = transformFunc
	f.configBuilder.definition.FieldMappings[f.targetField] = mapping
	return f.configBuilder
}

// When adds a condition for conditional mapping
func (f *FieldConfigBuilder[S, T]) When(conditionFunc func(S) bool) *ConditionConfigBuilder[S, T] {
	return &ConditionConfigBuilder[S, T]{
		fieldConfigBuilder: f,
		conditionFunc:      conditionFunc,
	}
}

// ConditionConfigBuilder handles conditional mapping configuration
type ConditionConfigBuilder[S, T any] struct {
	fieldConfigBuilder *FieldConfigBuilder[S, T]
	conditionFunc      func(S) bool
}

// FromField sets the source field for conditional mapping
func (c *ConditionConfigBuilder[S, T]) FromField(sourceField string) *ConfigBuilder[S, T] {
	c.fieldConfigBuilder.configBuilder.definition.FieldMappings[c.fieldConfigBuilder.targetField] = FieldMapping{
		TargetField: c.fieldConfigBuilder.targetField,
		MappingType: MappingTypeField,
		SourceField: sourceField,
		Condition:   c.conditionFunc,
	}
	return c.fieldConfigBuilder.configBuilder
}

// FromFunc sets a custom function for conditional mapping
func (c *ConditionConfigBuilder[S, T]) FromFunc(mapperFunc func(S) any) *ConfigBuilder[S, T] {
	c.fieldConfigBuilder.configBuilder.definition.FieldMappings[c.fieldConfigBuilder.targetField] = FieldMapping{
		TargetField: c.fieldConfigBuilder.targetField,
		MappingType: MappingTypeFunc,
		CustomFunc:  mapperFunc,
		Condition:   c.conditionFunc,
	}
	return c.fieldConfigBuilder.configBuilder
}

// getMappingKey generates a unique key for a mapping configuration
func getMappingKey(sourceType, targetType reflect.Type) string {
	return fmt.Sprintf("%s->%s", sourceType.String(), targetType.String())
}

// GetMappingConfig retrieves a mapping configuration for given types
func GetMappingConfig(sourceType, targetType reflect.Type) *MappingDefinition {
	// Try exact match first
	key := getMappingKey(sourceType, targetType)
	if config := globalConfigs[key]; config != nil {
		return config
	}

	// Try to dereference pointer types and match again
	actualSourceType := sourceType
	if sourceType.Kind() == reflect.Ptr {
		actualSourceType = sourceType.Elem()
	}

	actualTargetType := targetType
	if targetType.Kind() == reflect.Ptr {
		actualTargetType = targetType.Elem()
	}

	// Try with dereferenced types
	if actualSourceType != sourceType || actualTargetType != targetType {
		key = getMappingKey(actualSourceType, actualTargetType)
		return globalConfigs[key]
	}

	return nil
}

// SetGlobalConfig sets the global configuration options
func SetGlobalConfig(config GlobalConfig) {
	globalConfig = config
}

// GetGlobalConfig returns the current global configuration
func GetGlobalConfig() GlobalConfig {
	return globalConfig
}

// EnableTimeConversion enables or disables automatic time conversion
func EnableTimeConversion(enable bool) {
	globalConfig.EnableTimeConversion = enable
}

// IsTimeConversionEnabled returns whether time conversion is currently enabled
func IsTimeConversionEnabled() bool {
	return globalConfig.EnableTimeConversion
}
