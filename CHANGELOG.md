# 变更日志

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