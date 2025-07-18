# Go Mapster 项目技术文档

## 项目概述

Go Mapster 是一个高性能的对象映射库，采用"约定优于配置"的设计理念，支持零配置自动映射和灵活的自定义配置。

## 核心架构设计

### 1. 系统架构图

```mermaid
graph TB
    UserCall[用户调用] --> TypeCheck{映射类型检查}
    TypeCheck -->|有预生成代码| GeneratedMapping[生成代码映射]
    TypeCheck -->|无预生成代码| ReflectionMapping[反射映射]
    
    ReflectionMapping --> ConfigCheck{是否有配置}
    ConfigCheck -->|有配置| ConfigBasedMapping[配置映射]
    ConfigCheck -->|无配置| DefaultMapping[默认映射]
    
    ConfigBasedMapping --> CustomFuncMapping[自定义函数映射]
    ConfigBasedMapping --> FieldNameMapping[字段名映射]
    ConfigBasedMapping --> TransformMapping[转换映射]
    
    DefaultMapping --> SameNameFieldMatching[同名字段匹配]
    SameNameFieldMatching --> TypeCompatibilityCheck[类型兼容检查]
    
    GeneratedMapping --> ReturnResult[返回结果]
    CustomFuncMapping --> ReturnResult
    FieldNameMapping --> ReturnResult
    TransformMapping --> ReturnResult
    TypeCompatibilityCheck --> ReturnResult
```

### 2. 包结构设计

```
mapster/
├── mapster.go          # 核心 API 入口
├── config.go           # 配置系统
├── reflection.go       # 反射映射引擎
├── mapster_test.go     # 单元测试
├── benchmark_test.go   # 性能测试
└── examples/           # 示例代码
```

## 核心函数设计

### 1. 主要 API 函数

#### `Map[T any](src any) T`
**功能**: 将源对象映射到目标类型
**设计思路**:
- 使用泛型确保类型安全
- 优先使用生成的映射函数
- 回退到反射映射

```go
func Map[T any](src any) T {
    // 1. 空值检查
    if src == nil {
        var zero T
        return zero
    }

    // 2. 获取类型信息
    srcType := reflect.TypeOf(src)
    targetType := reflect.TypeOf(*new(T))

    // 3. 尝试使用生成的映射器
    if mapper := getGeneratedMapper(srcType, targetType); mapper != nil {
        if mapperFunc, ok := mapper.(func(interface{}) interface{}); ok {
            return mapperFunc(src).(T)
        }
    }

    // 4. 回退到反射映射
    return reflectionMap[T](src)
}
```

#### `MapSlice[T any](src any) []T`
**功能**: 批量映射切片
**优化**: 预分配切片容量，避免动态扩容

#### `MapTo[T any](src any, target *T)`
**功能**: 映射到现有对象
**用途**: 避免额外内存分配

### 2. 配置系统设计

#### 配置数据结构

```go
type MappingDefinition struct {
    SourceType    reflect.Type
    TargetType    reflect.Type
    FieldMappings map[string]FieldMapping
}

type FieldMapping struct {
    TargetField string
    MappingType MappingType
    SourceField string
    CustomFunc  interface{}
    Transform   interface{}
    Condition   interface{}
}
```

#### 配置流程图

```mermaid
graph LR
    ConfigStart[Config[S,T]()] --> MapMethod[Map(字段名)]
    MapMethod --> MappingTypeChoice{映射类型}
    MappingTypeChoice -->|FromField| FieldMapping[字段映射]
    MappingTypeChoice -->|FromFunc| FunctionMapping[函数映射]
    MappingTypeChoice -->|Transform| TransformMapping[转换映射]
    MappingTypeChoice -->|When| ConditionalMapping[条件映射]
    FieldMapping --> RegisterConfig[Register()]
    FunctionMapping --> RegisterConfig
    TransformMapping --> RegisterConfig
    ConditionalMapping --> RegisterConfig
    RegisterConfig --> StoreGlobalConfig[存储到全局配置]
```

### 3. 反射映射引擎

#### 核心函数设计

##### `reflectionMap[T any](src any) T`
**执行流程**:
1. 创建目标对象
2. 检查是否有自定义配置
3. 选择映射策略（配置映射 vs 默认映射）

##### `mapWithConfig(srcValue, targetValue, config)`
**配置映射流程**:

```mermaid
graph TD
    StartConfigMapping[开始配置映射] --> BuildSourceIndex[构建源字段索引]
    BuildSourceIndex --> IterateTargetFields[遍历目标字段]
    IterateTargetFields --> CheckFieldConfig{检查字段配置}
    CheckFieldConfig -->|忽略| SkipField[跳过字段]
    CheckFieldConfig -->|字段映射| CopyFromSourceField[从源字段复制]
    CheckFieldConfig -->|函数映射| CallCustomFunction[调用自定义函数]
    CheckFieldConfig -->|转换映射| ApplyTransform[应用转换函数]
    CheckFieldConfig -->|无配置| DefaultSameNameMapping[默认同名映射]
    SkipField --> HasMoreFields{还有字段?}
    CopyFromSourceField --> HasMoreFields
    CallCustomFunction --> HasMoreFields
    ApplyTransform --> HasMoreFields
    DefaultSameNameMapping --> HasMoreFields
    HasMoreFields -->|是| IterateTargetFields
    HasMoreFields -->|否| CompleteMaping[完成映射]
```

##### `mapReflect(srcValue, targetValue)`
**默认映射策略**:

```mermaid
graph TD
    StartDefaultMapping[开始默认映射] --> SourceTypeCheck{源类型检查}
    SourceTypeCheck -->|相同类型| DirectAssignment[直接赋值]
    SourceTypeCheck -->|结构体到结构体| FieldLevelMapping[字段级映射]
    SourceTypeCheck -->|指针处理| PointerDereference[指针解引用]
    SourceTypeCheck -->|切片映射| ElementLevelMapping[元素级映射]
    SourceTypeCheck -->|类型转换| CompatibilityCheck[兼容性检查]
    
    FieldLevelMapping --> BuildFieldIndex[建立字段索引]
    BuildFieldIndex --> SameNameFieldMatch[同名字段匹配]
    SameNameFieldMatch --> RecursiveMapping[递归映射]
    
    DirectAssignment --> Complete[完成]
    PointerDereference --> Complete
    ElementLevelMapping --> Complete
    CompatibilityCheck --> Complete
    RecursiveMapping --> Complete
```

## 性能优化策略

### 1. 反射优化

#### 字段缓存机制
```go
// 源字段索引缓存
srcFields := make(map[string]reflect.Value)
for i := 0; i < srcType.NumField(); i++ {
    field := srcType.Field(i)
    if field.IsExported() {
        srcFields[field.Name] = srcValue.Field(i)
    }
}
```

#### 类型兼容性检查
```go
// 智能类型转换
if result.Type().AssignableTo(targetFieldValue.Type()) {
    targetFieldValue.Set(result)
} else if result.Type().ConvertibleTo(targetFieldValue.Type()) {
    targetFieldValue.Set(result.Convert(targetFieldValue.Type()))
} else if result.Kind() == reflect.Interface && !result.IsNil() {
    // 处理 interface{} 返回值
    actualValue := result.Elem()
    if actualValue.Type().AssignableTo(targetFieldValue.Type()) {
        targetFieldValue.Set(actualValue)
    }
}
```

### 2. 内存优化

#### 预分配策略
- 切片映射时预分配目标切片
- 避免多次内存分配和复制

#### 零拷贝优化
- 相同类型直接赋值
- 避免不必要的数据复制

### 3. 代码生成优化（未来实现）

```mermaid
graph LR
    SourceCodeAnalysis[源码分析] --> ASTParser[AST 解析]
    ASTParser --> MappingCodeGeneration[映射代码生成]
    MappingCodeGeneration --> CompileTimeIntegration[编译时集成]
    CompileTimeIntegration --> ZeroReflectionOverhead[零反射开销]
```

## 错误处理机制

### 1. 类型安全检查
- 编译时泛型类型检查
- 运行时类型兼容性验证

### 2. 空值处理
```go
// 多层空值检查
if src == nil {
    return zero
}
if !srcValue.IsValid() || !targetValue.IsValid() {
    return
}
```

### 3. 函数调用安全
```go
// 自定义函数调用保护
func callCustomFunc(fn interface{}, src interface{}) reflect.Value {
    fnValue := reflect.ValueOf(fn)
    if fnValue.Kind() != reflect.Func {
        return reflect.Value{}
    }
    
    fnType := fnValue.Type()
    if fnType.NumIn() != 1 || fnType.NumOut() != 1 {
        return reflect.Value{}
    }
    // ... 安全调用
}
```

## 测试策略

### 1. 单元测试覆盖
- 基础映射功能
- 配置映射功能
- 边界条件测试
- 错误处理测试

### 2. 性能基准测试
```go
// 基准测试指标
BenchmarkBasicMapping-8         927649   1199 ns/op   416 B/op   12 allocs/op
BenchmarkSliceMapping-8           9754 120473 ns/op 51115 B/op 1202 allocs/op
BenchmarkWithoutConfiguration-8 2441320   490 ns/op   224 B/op    8 allocs/op
```

### 3. 测试数据设计
- 简单结构体映射
- 复杂嵌套对象
- 大型切片处理
- 特殊类型处理

## 扩展性设计

### 1. 插件化架构
- 自定义映射函数接口
- 转换器插件系统
- 验证器集成点

### 2. 代码生成接口
```go
// 为代码生成器预留的接口
type GeneratedMapper interface {
    Map(src interface{}) interface{}
    SourceType() reflect.Type
    TargetType() reflect.Type
}
```

### 3. 配置扩展点
- 自定义映射策略
- 条件映射逻辑
- 类型转换器

## 未来优化方向

### 1. 编译时代码生成
- AST 分析工具
- 模板代码生成
- Go generate 集成

### 2. 高级映射功能
- 深度嵌套对象映射
- 循环引用处理
- 动态字段映射

### 3. 性能进一步优化
- 字段访问缓存
- 类型信息预计算
- 并发安全优化

## 内部关键函数调用逻辑

### 1. 主要映射函数调用链

```mermaid
graph TD
    MapFunctionEntry[Map[T any]src any] --> NullCheck{src == nil?}
    NullCheck -->|是| ReturnZeroValue[返回零值]
    NullCheck -->|否| GetTypeInfo[获取类型信息]
    
    GetTypeInfo --> CheckGeneratedMapper[getGeneratedMapper]
    CheckGeneratedMapper --> HasGeneratedMapper{有生成映射器?}
    HasGeneratedMapper -->|是| CallGeneratedMapper[调用生成映射器]
    HasGeneratedMapper -->|否| UseReflectionMap[reflectionMap[T]]
    
    UseReflectionMap --> CreateTargetObject[创建目标对象]
    CreateTargetObject --> GetMappingConfigCall[GetMappingConfig]
    GetMappingConfigCall --> HasCustomConfig{有自定义配置?}
    HasCustomConfig -->|是| UseMapWithConfig[mapWithConfig]
    HasCustomConfig -->|否| UseMapReflect[mapReflect]
    
    UseMapWithConfig --> BuildSourceFieldIndex[构建源字段索引]
    UseMapWithConfig --> IterateTargetField[遍历目标字段]
    UseMapWithConfig --> ApplyFieldMapping[应用字段映射]
    
    UseMapReflect --> CheckTypeCompatibility[检查类型兼容性]
    UseMapReflect --> CallMapStructToStruct[mapStructToStruct]
    
    CallGeneratedMapper --> ReturnFinalResult[返回结果]
    ApplyFieldMapping --> ReturnFinalResult
    CheckTypeCompatibility --> ReturnFinalResult
    CallMapStructToStruct --> ReturnFinalResult
    ReturnZeroValue --> ReturnFinalResult
```

### 2. 配置系统调用逻辑

```mermaid
graph TD
    ConfigCall[Config[S,T]] --> CreateConfigBuilder[创建 ConfigBuilder]
    CreateConfigBuilder --> MapMethodCall[Map方法调用]
    MapMethodCall --> CreateFieldConfigBuilder[创建 FieldConfigBuilder]
    
    CreateFieldConfigBuilder --> ConfigTypeChoice{配置类型}
    ConfigTypeChoice -->|FromField| SetFieldMapping[设置字段映射]
    ConfigTypeChoice -->|FromFunc| SetFunctionMapping[设置函数映射]
    ConfigTypeChoice -->|Transform| SetTransformMapping[设置转换映射]
    ConfigTypeChoice -->|When| SetConditionalMapping[设置条件映射]
    
    SetFieldMapping --> ReturnConfigBuilder[返回 ConfigBuilder]
    SetFunctionMapping --> ReturnConfigBuilder
    SetTransformMapping --> ReturnConfigBuilder
    SetConditionalMapping --> CreateConditionConfigBuilder[创建 ConditionConfigBuilder]
    CreateConditionConfigBuilder --> SetConditionLogic[设置条件逻辑]
    SetConditionLogic --> ReturnConfigBuilder
    
    ReturnConfigBuilder --> RegisterMethod[Register方法]
    RegisterMethod --> CallGetMappingKey[getMappingKey]
    RegisterMethod --> StoreToGlobalConfigs[存储到 globalConfigs]
    
    StoreToGlobalConfigs --> ConfigComplete[配置完成]
```

### 3. 反射映射详细调用流程

```mermaid
graph TD
    MapReflectEntry[mapReflect] --> TypeCheckDecision{类型检查}
    TypeCheckDecision -->|相同类型| DirectValueSet[targetValue.Set]
    TypeCheckDecision -->|结构体| CallMapStructToStruct[mapStructToStruct]
    TypeCheckDecision -->|指针| PointerHandlingLogic[指针处理逻辑]
    TypeCheckDecision -->|切片| CallMapSliceToSlice[mapSliceToSlice]
    TypeCheckDecision -->|可转换| TypeConversion[类型转换]
    
    CallMapStructToStruct --> BuildSourceFieldMapping[建立源字段映射]
    BuildSourceFieldMapping --> IterateTargetFields[遍历目标字段]
    IterateTargetFields --> FieldExistsCheck{字段存在?}
    FieldExistsCheck -->|是| RecursiveMapReflectCall[递归调用 mapReflect]
    FieldExistsCheck -->|否| SkipCurrentField[跳过字段]
    
    PointerHandlingLogic --> CheckPointerNotNull[检查指针非空]
    CheckPointerNotNull --> CreateNewPointer[创建新指针]
    CreateNewPointer --> RecursiveMappingContent[递归映射内容]
    
    CallMapSliceToSlice --> CreateTargetSlice[创建目标切片]
    CreateTargetSlice --> IterateSourceSliceElements[遍历源切片元素]
    IterateSourceSliceElements --> RecursiveMapEachElement[递归映射每个元素]
    
    DirectValueSet --> MappingComplete[完成]
    RecursiveMapReflectCall --> MappingComplete
    SkipCurrentField --> MappingComplete
    RecursiveMappingContent --> MappingComplete
    RecursiveMapEachElement --> MappingComplete
    TypeConversion --> MappingComplete
```

### 4. 自定义函数调用安全机制

```mermaid
graph TD
    CallCustomFuncEntry[callCustomFunc] --> ReflectValueOf[reflect.ValueOf]
    ReflectValueOf --> IsFunctionTypeCheck{是函数类型?}
    IsFunctionTypeCheck -->|否| ReturnEmptyValue[返回空值]
    IsFunctionTypeCheck -->|是| CheckFunctionSignature[检查函数签名]
    
    CheckFunctionSignature --> ParamCountCheck{参数数量 == 1?}
    ParamCountCheck -->|否| ReturnEmptyValue
    ParamCountCheck -->|是| ReturnCountCheck{返回值数量 == 1?}
    ReturnCountCheck -->|否| ReturnEmptyValue
    ReturnCountCheck -->|是| CheckParamTypeCompatibility[检查参数类型兼容]
    
    CheckParamTypeCompatibility --> TypeAssignableCheck{类型可分配?}
    TypeAssignableCheck -->|否| ReturnEmptyValue
    TypeAssignableCheck -->|是| FunctionValueCall[fnValue.Call]
    
    FunctionValueCall --> GetReturnValue[获取返回值]
    GetReturnValue --> TypeSafetyCheck[类型安全检查]
    TypeSafetyCheck --> ReturnFinalResult[返回结果]
    
    ReturnEmptyValue --> FunctionEnd[结束]
    ReturnFinalResult --> FunctionEnd
```

### 5. 配置映射执行流程

```mermaid
graph TD
    MapWithConfigEntry[mapWithConfig] --> ValidateInputParams[验证输入参数]
    ValidateInputParams --> StructTypeCheck{结构体类型?}
    StructTypeCheck -->|否| FallbackToMapReflect[回退到 mapReflect]
    StructTypeCheck -->|是| BuildSourceFieldIndex[构建源字段索引]
    
    BuildSourceFieldIndex --> IterateTargetFields[遍历目标字段]
    IterateTargetFields --> GetFieldConfig[获取字段配置]
    GetFieldConfig --> MappingTypeDecision{配置类型}
    
    MappingTypeDecision -->|MappingTypeIgnore| SkipCurrentField[跳过字段]
    MappingTypeDecision -->|MappingTypeField| FieldNameMapping[字段名映射]
    MappingTypeDecision -->|MappingTypeFunc| CustomFunctionMapping[自定义函数映射]
    MappingTypeDecision -->|MappingTypeTransform| TransformMapping[转换映射]
    MappingTypeDecision -->|无配置| DefaultSameNameMapping[默认同名映射]
    
    FieldNameMapping --> CopyFromSourceField[从源字段复制]
    CustomFunctionMapping --> CallCustomFunc[callCustomFunc]
    TransformMapping --> MapThenTransform[先映射再转换]
    DefaultSameNameMapping --> SameNameFieldLookup[同名字段查找]
    
    CallCustomFunc --> HandleResultType[结果类型处理]
    HandleResultType --> InterfaceTypeCheck{interface{}类型?}
    InterfaceTypeCheck -->|是| ExtractActualValue[提取实际值]
    InterfaceTypeCheck -->|否| DirectTypeCheck[直接类型检查]
    
    ExtractActualValue --> TypeCompatibilityCheck[类型兼容检查]
    DirectTypeCheck --> TypeCompatibilityCheck
    TypeCompatibilityCheck --> SetTargetField[设置目标字段]
    
    CopyFromSourceField --> ContinueNextField[继续下一字段]
    MapThenTransform --> ContinueNextField
    SameNameFieldLookup --> ContinueNextField
    SetTargetField --> ContinueNextField
    SkipCurrentField --> ContinueNextField
    
    ContinueNextField --> HasMoreFieldsCheck{还有字段?}
    HasMoreFieldsCheck -->|是| IterateTargetFields
    HasMoreFieldsCheck -->|否| MappingCompleted[映射完成]
```

### 6. 类型转换安全检查

```mermaid
graph TD
    TypeConversionEntry[类型转换检查] --> GetResultType[result.Type]
    GetResultType --> GetTargetFieldType[targetFieldValue.Type]
    GetTargetFieldType --> AssignableToCheck{AssignableTo?}
    
    AssignableToCheck -->|是| DirectAssignment[直接赋值]
    AssignableToCheck -->|否| ConvertibleToCheck{ConvertibleTo?}
    ConvertibleToCheck -->|是| TypeConversionAssignment[类型转换后赋值]
    ConvertibleToCheck -->|否| InterfaceTypeCheck{Interface类型?}
    
    InterfaceTypeCheck -->|是| ResultElem[result.Elem]
    InterfaceTypeCheck -->|否| IgnoreField[忽略该字段]
    
    ResultElem --> GetActualValueType[actualValue.Type]
    GetActualValueType --> ActualTypeAssignableCheck{实际类型可分配?}
    ActualTypeAssignableCheck -->|是| SetActualValue[设置实际值]
    ActualTypeAssignableCheck -->|否| ActualTypeConvertibleCheck{实际类型可转换?}
    ActualTypeConvertibleCheck -->|是| ConvertAndSet[转换后设置]
    ActualTypeConvertibleCheck -->|否| IgnoreField
    
    DirectAssignment --> ConversionComplete[完成]
    TypeConversionAssignment --> ConversionComplete
    SetActualValue --> ConversionComplete
    ConvertAndSet --> ConversionComplete
    IgnoreField --> ConversionComplete
```

### 7. 系统函数依赖关系图

```mermaid
graph TB
    subgraph "用户 API 层"
        MapAPI[Map[T]]
        MapSliceAPI[MapSlice[T]]
        MapToAPI[MapTo[T]]
    end
    
    subgraph "配置 API 层"
        ConfigAPI[Config[S,T]]
        ConfigBuilderMap[ConfigBuilder.Map]
        FieldConfigBuilder[FieldConfigBuilder]
        RegisterConfig[Register]
    end
    
    subgraph "核心映射引擎"
        ReflectionMapEngine[reflectionMap]
        MapWithConfigEngine[mapWithConfig]
        MapReflectEngine[mapReflect]
        MapStructToStructEngine[mapStructToStruct]
        MapSliceToSliceEngine[mapSliceToSlice]
    end
    
    subgraph "辅助函数层"
        GetGeneratedMapperHelper[getGeneratedMapper]
        GetMappingConfigHelper[GetMappingConfig]
        CallCustomFuncHelper[callCustomFunc]
        CallTransformFuncHelper[callTransformFunc]
        GetMappingKeyHelper[getMappingKey]
    end
    
    subgraph "类型系统层"
        ReflectTypeOf[reflect.TypeOf]
        ReflectValueOf[reflect.ValueOf]
        TypeCompatibilityCheck[类型兼容检查]
        FieldIndexBuilder[字段索引构建]
    end
    
    MapAPI --> GetGeneratedMapperHelper
    MapAPI --> ReflectionMapEngine
    MapSliceAPI --> MapAPI
    MapToAPI --> MapAPI
    
    ConfigAPI --> ConfigBuilderMap
    ConfigBuilderMap --> FieldConfigBuilder
    FieldConfigBuilder --> RegisterConfig
    RegisterConfig --> GetMappingKeyHelper
    RegisterConfig --> GetMappingConfigHelper
    
    ReflectionMapEngine --> GetMappingConfigHelper
    ReflectionMapEngine --> MapWithConfigEngine
    ReflectionMapEngine --> MapReflectEngine
    MapWithConfigEngine --> CallCustomFuncHelper
    MapWithConfigEngine --> CallTransformFuncHelper
    MapWithConfigEngine --> FieldIndexBuilder
    MapReflectEngine --> MapStructToStructEngine
    MapReflectEngine --> MapSliceToSliceEngine
    MapReflectEngine --> TypeCompatibilityCheck
    
    GetGeneratedMapperHelper --> ReflectTypeOf
    GetMappingConfigHelper --> GetMappingKeyHelper
    CallCustomFuncHelper --> ReflectValueOf
    CallTransformFuncHelper --> ReflectValueOf
    MapStructToStructEngine --> FieldIndexBuilder
    MapSliceToSliceEngine --> MapReflectEngine
    
    MapWithConfigEngine --> ReflectTypeOf
    MapReflectEngine --> ReflectTypeOf
    MapStructToStructEngine --> ReflectValueOf
    MapSliceToSliceEngine --> ReflectValueOf
```

### 8. 错误处理调用链

```mermaid
graph TD
    FunctionEntry[函数入口] --> ParameterValidation[参数验证]
    ParameterValidation --> ParamsValidCheck{参数有效?}
    ParamsValidCheck -->|否| ReturnZeroOrError[返回零值/错误]
    ParamsValidCheck -->|是| TypeCheck[类型检查]
    
    TypeCheck --> TypeCompatibleCheck{类型兼容?}
    TypeCompatibleCheck -->|否| TryTypeConversion[尝试类型转换]
    TypeCompatibleCheck -->|是| ExecuteMappingLogic[执行映射逻辑]
    
    TryTypeConversion --> ConversionSuccessCheck{转换成功?}
    ConversionSuccessCheck -->|否| LogWarningOrSkip[记录警告/跳过]
    ConversionSuccessCheck -->|是| ExecuteMappingLogic
    
    ExecuteMappingLogic --> ReflectionOperation[反射操作]
    ReflectionOperation --> ReflectionSafeCheck{反射安全?}
    ReflectionSafeCheck -->|否| CatchPanic[捕获panic]
    ReflectionSafeCheck -->|是| SetFieldValue[设置字段值]
    
    CatchPanic --> ErrorRecovery[错误恢复]
    ErrorRecovery --> ContinueProcessOtherFields[继续处理其他字段]
    
    SetFieldValue --> HasMoreFieldsCheck{还有字段?}
    HasMoreFieldsCheck -->|是| ExecuteMappingLogic
    HasMoreFieldsCheck -->|否| ReturnResult[返回结果]
    
    ReturnZeroOrError --> ProcessEnd[结束]
    LogWarningOrSkip --> ContinueProcessOtherFields
    ContinueProcessOtherFields --> HasMoreFieldsCheck
    ReturnResult --> ProcessEnd
```

### 9. 性能优化决策流程

```mermaid
graph TD
    MappingRequest[映射请求] --> CachedMapperCheck{缓存中有映射器?}
    CachedMapperCheck -->|是| UseCachedMapper[使用缓存映射器]
    CachedMapperCheck -->|否| AnalyzeMappingComplexity[分析映射复杂度]
    
    AnalyzeMappingComplexity --> SimpleMappingCheck{简单映射?}
    SimpleMappingCheck -->|是| DirectReflectionMapping[直接反射映射]
    SimpleMappingCheck -->|否| CustomConfigCheck{有自定义配置?}
    
    CustomConfigCheck -->|否| StructuredReflectionMapping[结构化反射映射]
    CustomConfigCheck -->|是| ConfigDrivenMapping[配置驱动映射]
    
    DirectReflectionMapping --> FieldLevelCaching[字段级缓存]
    StructuredReflectionMapping --> TypeLevelCaching[类型级缓存]
    ConfigDrivenMapping --> ConfigLevelCaching[配置级缓存]
    
    FieldLevelCaching --> PerformanceCheck{性能满足?}
    TypeLevelCaching --> PerformanceCheck
    ConfigLevelCaching --> PerformanceCheck
    
    PerformanceCheck -->|是| CompleteMapping[完成映射]
    PerformanceCheck -->|否| MarkForCodeGeneration[标记为代码生成候选]
    
    MarkForCodeGeneration --> GenerateSpecializedMapper[生成专用映射函数]
    GenerateSpecializedMapper --> RegisterGeneratedMapper[注册生成映射器]
    RegisterGeneratedMapper --> UseCachedMapper
    
    UseCachedMapper --> CompleteMapping
```

### 10. 内存管理优化流程

```mermaid
graph TD
    StartMapping[开始映射] --> TargetObjectTypeCheck{目标对象类型}
    TargetObjectTypeCheck -->|值类型| StackAllocationStrategy[栈分配策略]
    TargetObjectTypeCheck -->|指针类型| HeapAllocationStrategy[堆分配策略]
    TargetObjectTypeCheck -->|切片类型| BatchAllocationStrategy[批量分配策略]
    
    StackAllocationStrategy --> DirectValueCopy[直接值复制]
    HeapAllocationStrategy --> SourceNilCheck{源为nil?}
    SourceNilCheck -->|是| SetToNil[设置nil]
    SourceNilCheck -->|否| AllocateNewObject[分配新对象]
    
    BatchAllocationStrategy --> PreAllocateSliceCapacity[预分配切片容量]
    PreAllocateSliceCapacity --> BatchElementMapping[批量元素映射]
    
    AllocateNewObject --> RecursiveMappingContent[递归映射内容]
    BatchElementMapping --> ElementLevelMapping[元素级映射]
    
    DirectValueCopy --> MappingComplete[完成]
    SetToNil --> MappingComplete
    RecursiveMappingContent --> MappingComplete
    ElementLevelMapping --> MappingComplete
```

## 最佳实践建议

### 1. 使用建议
- 优先使用零配置映射
- 复杂场景才使用自定义配置
- 注意函数返回类型匹配

### 2. 性能建议
- 避免频繁的类型转换
- 大切片考虑分批处理
- 复用映射配置

### 3. 安全建议
- 验证自定义函数的类型安全
- 处理可能的 panic 情况
- 注意并发访问安全
