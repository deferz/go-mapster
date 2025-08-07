# 变更日志

## [1.0.0] - 2025-01-07

### 🎉 首次正式发布

#### 新增功能
- **🚀 零反射代码生成**：支持生成优化的映射代码，性能提升 1.5x
- **⏰ 智能时间转换**：自动 int64 ↔ time.Time 转换，支持可配置行为
- **🔧 多层优化策略**：智能回退机制，自动选择最优映射策略
- **📊 深度路径解析**：支持嵌套对象属性访问（如 `Company.Address.City`）
- **🔄 循环引用检测**：安全处理包含循环引用的复杂对象图
- **🎭 灵活配置系统**：链式 API，支持自定义映射、转换和条件映射

#### 核心 API
- `Map[T any](src any) T` - 主要映射函数
- `MapTo[T any](src any, target *T)` - 映射到现有对象
- `Config[S, T any]()` - 配置系统
- `RegisterGeneratedMapper()` - 零反射优化
- `EnableTimeConversion()` - 时间转换配置

#### 性能表现
```
手动映射:          18 ns/op    0 B/op    0 allocs/op  ⭐ 最快
零反射映射:        422 ns/op  304 B/op    7 allocs/op  🚀 生成代码
配置映射:          490 ns/op  224 B/op    8 allocs/op  🔧 自定义配置
纯反射映射:        732 ns/op  320 B/op    8 allocs/op  🔄 自动映射
```

#### 技术特性
- 基于 Go 1.18+ 泛型，编译时类型安全
- 零外部依赖，纯 Go 实现
- 支持 Go 1.24.3+
- 完整的测试覆盖
- 详细的文档和示例

## [未发布]

### 移除
- 移除了 `MapSlice[T any](src any) []T` API
  - 原因：性能测试显示 MapSlice 没有性能提升，只是提供便利性
  - 替代方案：使用循环 + `Map` 函数进行批量映射
  - 示例：
    ```go
    // 之前
    dtos := mapster.MapSlice[UserDTO](users)
    
    // 现在
    dtos := make([]UserDTO, len(users))
    for i, u := range users {
        dtos[i] = mapster.Map[UserDTO](u)
    }
    ```

### 保留的核心 API
- `Map[T any](src any) T` - 主要映射函数
- `MapTo[T any](src any, target *T)` - 映射到现有对象
- 配置系统：`Config[S, T any]()` 等
- 零反射优化：`RegisterGeneratedMapper()`

### 性能影响
- 移除 MapSlice 后，批量映射的性能与之前相同
- 用户需要手动编写循环，但性能没有损失
- 代码更简洁，减少了不必要的 API 复杂度 