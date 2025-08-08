# Mapster 全局配置

Mapster 提供了全局配置功能，可以统一管理映射行为，包括时间转换和自动扁平化功能。

## 全局配置结构

```go
type GlobalConfig struct {
    // 启用时间转换（int64 <-> time.Time）
    EnableTimeConversion bool
    
    // 自动扁平化配置
    AutoFlatten AutoFlattenConfig
}

type AutoFlattenConfig struct {
    // 是否启用自动扁平化
    Enabled bool
    
    // 扁平化深度限制
    MaxDepth int
    
    // 字段名冲突处理策略
    ConflictStrategy ConflictStrategy
    
    // 是否使用前缀
    UsePrefix bool
    
    // 前缀分隔符
    PrefixSeparator string
}
```

## 默认配置

```go
var globalConfig = GlobalConfig{
    EnableTimeConversion: true,  // 默认启用时间转换
    AutoFlatten: AutoFlattenConfig{
        Enabled:         false,  // 默认禁用自动扁平化
        MaxDepth:        3,
        ConflictStrategy: KeepFirst,
        UsePrefix:       false,
        PrefixSeparator: "_",
    },
}
```

## 全局配置函数

### 获取和设置全局配置

```go
// 获取当前全局配置
config := mapster.GetGlobalConfig()

// 设置全局配置
newConfig := mapster.GlobalConfig{
    EnableTimeConversion: true,
    AutoFlatten: mapster.AutoFlattenConfig{
        Enabled:         true,
        MaxDepth:        2,
        ConflictStrategy: mapster.KeepLast,
        UsePrefix:       false,
        PrefixSeparator: "_",
    },
}
mapster.SetGlobalConfig(newConfig)
```

### 时间转换配置

```go
// 启用时间转换
mapster.EnableTimeConversion(true)

// 禁用时间转换
mapster.EnableTimeConversion(false)

// 检查时间转换是否启用
enabled := mapster.IsTimeConversionEnabled()
```

### 自动扁平化配置

#### 基本控制

```go
// 启用自动扁平化
mapster.EnableAutoFlatten()

// 禁用自动扁平化
mapster.DisableAutoFlatten()

// 检查自动扁平化是否启用
enabled := mapster.IsAutoFlattenEnabled()
```

#### 深度控制

```go
// 设置最大深度
mapster.SetAutoFlattenMaxDepth(3)

// 获取当前最大深度
depth := mapster.GetAutoFlattenMaxDepth()
```

#### 冲突策略

```go
// 设置冲突处理策略
mapster.SetAutoFlattenConflictStrategy(mapster.KeepFirst)  // 保留第一个
mapster.SetAutoFlattenConflictStrategy(mapster.KeepLast)   // 保留最后一个
mapster.SetAutoFlattenConflictStrategy(mapster.UsePrefix)  // 使用前缀
mapster.SetAutoFlattenConflictStrategy(mapster.Skip)       // 跳过冲突

// 获取当前冲突策略
strategy := mapster.GetAutoFlattenConflictStrategy()
```

#### 前缀配置

```go
// 设置是否使用前缀
mapster.SetAutoFlattenUsePrefix(true)

// 检查是否使用前缀
usePrefix := mapster.IsAutoFlattenUsePrefix()

// 设置前缀分隔符
mapster.SetAutoFlattenPrefixSeparator("_")

// 获取当前前缀分隔符
separator := mapster.GetAutoFlattenPrefixSeparator()
```

#### 完整配置

```go
// 设置完整的自动扁平化配置
config := mapster.AutoFlattenConfig{
    Enabled:         true,
    MaxDepth:        2,
    ConflictStrategy: mapster.KeepLast,
    UsePrefix:       true,
    PrefixSeparator: "_",
}
mapster.SetAutoFlattenConfig(config)

// 获取当前自动扁平化配置
config := mapster.GetAutoFlattenConfig()
```

## 使用示例

### 示例 1: 基本全局配置

```go
package main

import "github.com/deferz/go-mapster"

func main() {
    // 获取当前配置
    config := mapster.GetGlobalConfig()
    fmt.Printf("时间转换: %v\n", config.EnableTimeConversion)
    fmt.Printf("自动扁平化: %v\n", config.AutoFlatten.Enabled)
    
    // 启用自动扁平化
    mapster.EnableAutoFlatten()
    
    // 设置最大深度为 2
    mapster.SetAutoFlattenMaxDepth(2)
    
    // 现在所有的映射都会使用这些配置
    var target TargetStruct
    mapster.MapTo(source, &target)
}
```

### 示例 2: 完整配置设置

```go
package main

import "github.com/deferz/go-mapster"

func main() {
    // 设置完整的全局配置
    newConfig := mapster.GlobalConfig{
        EnableTimeConversion: true,
        AutoFlatten: mapster.AutoFlattenConfig{
            Enabled:         true,
            MaxDepth:        3,
            ConflictStrategy: mapster.KeepLast,
            UsePrefix:       false,
            PrefixSeparator: "_",
        },
    }
    mapster.SetGlobalConfig(newConfig)
    
    // 执行映射（使用新的配置）
    var target TargetStruct
    mapster.MapTo(source, &target)
}
```

### 示例 3: 应用程序初始化

```go
package main

import "github.com/deferz/go-mapster"

func init() {
    // 在应用程序启动时配置 mapster
    mapster.EnableTimeConversion(true)
    mapster.EnableAutoFlatten()
    mapster.SetAutoFlattenMaxDepth(3)
    mapster.SetAutoFlattenConflictStrategy(mapster.KeepFirst)
}

func main() {
    // 应用程序代码...
}
```

## 配置优先级

1. **全局配置**：影响所有映射操作
2. **特定类型配置**：使用 `mapster.Config[S, T]()` 创建的配置会覆盖全局配置
3. **运行时配置**：可以在运行时动态修改全局配置

## 注意事项

1. **线程安全**：全局配置是线程安全的，但建议在应用程序启动时设置
2. **性能影响**：自动扁平化会增加一些性能开销，建议根据实际需求启用
3. **向后兼容**：默认配置保持向后兼容，不会影响现有代码
4. **配置持久性**：全局配置在程序运行期间保持有效，重启后需要重新设置

## 最佳实践

1. **应用程序启动时配置**：在 `main()` 函数或 `init()` 函数中设置全局配置
2. **根据需求启用功能**：只在需要时启用自动扁平化
3. **合理设置深度**：避免设置过大的深度值，防止性能问题
4. **测试配置效果**：在生产环境使用前，充分测试配置的效果
