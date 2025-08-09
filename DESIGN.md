# go-mapster 设计文档

## 项目概述

go-mapster 是一个用 Go 语言编写的对象映射库，旨在提供类型安全、高性能的对象间转换功能。该库使用 Go 1.18+ 的泛型特性，在编译时确保类型安全，减少运行时错误。

## 核心功能

### 1. 主要 API

#### Map[T any](src any) (T, error)
- **功能**: 将源对象映射到目标类型并返回新实例
- **特点**: 使用泛型确保类型安全
- **用途**: 创建新的目标对象实例

```go
userDTO, err := mapster.Map[UserDTO](user)
```

#### MapTo[T any](src any, dst *T) error
- **功能**: 将源对象映射到现有的目标对象
- **特点**: 修改传入的目标对象，而不创建新实例
- **用途**: 更新现有对象的字段值

```go
var userDTO UserDTO
err := mapster.MapTo(user, &userDTO)
```

### 2. 支持的映射类型

根据代码注释和测试用例，该库支持以下映射场景：

1. **结构体到结构体映射**
   - 相同字段名自动映射
   - 支持嵌套结构体
   - 支持指针类型

2. **基本类型转换**
   - 支持可转换的基本类型之间的映射（如 int 到 int64）
   - 不支持字符串到数字的直接转换（需要自定义转换器）

3. **集合类型映射**
   - 切片到切片的映射
   - 数组到数组的映射（支持截断）
   - Map 类型的键值对映射

4. **指针类型处理**
   - 自动处理源对象的指针类型
   - 支持嵌套指针

## 架构设计

### 1. 模块结构

```
go-mapster/
├── mapster.go          # 主入口，提供公共 API
└── internal/
    ├── mapper/         # 核心映射逻辑
    │   ├── core.go     # 核心接口定义
    │   ├── struct.go   # 结构体映射
    │   ├── slice.go    # 切片映射
    │   ├── embedded.go # 嵌入字段映射
    │   ├── circular.go # 循环引用检测
    │   └── assign.go   # 字段赋值逻辑
    ├── convert/        # 类型转换模块
    └── cache/          # 缓存模块（性能优化）
```

### 2. 核心接口

#### ValueConverter 接口
```go
type ValueConverter interface {
    Convert(from reflect.Value, to reflect.Type) (reflect.Value, bool)
}
```
- **职责**: 处理不同类型之间的转换
- **设计思路**: 先支持内置类型和别名转换，后续支持自定义转换器注册

#### FieldResolver 接口
```go
type FieldResolver interface {
    Resolve(value reflect.Value, fieldName string) (reflect.Value, bool)
}
```
- **职责**: 解析字段值
- **设计思路**: 支持自定义字段解析器，处理特殊的字段映射逻辑

### 3. 映射管线

根据 core.go 中的注释，映射流程遵循以下管线：

```
遍历目标字段 → 解析源值 → 类型转换 → 赋值
```

管线编排原则：
1. 先查配置（如果有自定义映射配置）
2. 再走字段解析（通过 FieldResolver）
3. 再走类型转换（通过 ValueConverter）
4. 最后执行赋值

## 设计特点

### 1. 类型安全
- 使用 Go 泛型提供编译时类型检查
- 减少运行时类型断言和错误

### 2. 性能优化
- 预留了 cache 模块用于缓存映射元数据
- 减少反射操作的开销

### 3. 扩展性
- 通过接口设计支持自定义转换器
- 支持自定义字段解析器
- 模块化设计便于功能扩展

### 4. 错误处理
- 完善的错误信息返回
- 支持错误链（error wrapping）

## 测试覆盖

根据测试文件，当前已覆盖的测试场景包括：

1. **基本映射测试**
   - 相同类型映射
   - 结构体到结构体映射
   - 指针源对象处理

2. **错误处理测试**
   - nil 源对象
   - nil 目标对象
   - 不可转换类型

3. **集合类型测试**
   - 切片映射
   - 数组映射（包括截断）
   - Map 类型映射

4. **类型转换测试**
   - 基本类型转换（如 int 到 int64）
   - Map 键类型转换

## 待实现功能

基于当前代码结构，以下功能尚待实现：

1. **核心映射逻辑**
   - MapValue 函数的具体实现
   - 各种类型的映射处理器

2. **性能优化**
   - 缓存机制的实现
   - 映射元数据的缓存

3. **高级功能**
   - 自定义转换器注册
   - 字段名映射配置
   - 循环引用处理
   - 忽略字段配置

4. **类型转换扩展**
   - 时间类型转换
   - 字符串与基本类型的转换

## 使用示例

### 基本使用
```go
// 定义源和目标结构体
type User struct {
    Name string
    Age  int
}

type UserDTO struct {
    Name string
    Age  int
}

// 使用 Map 创建新实例
userDTO, err := mapster.Map[UserDTO](user)
if err != nil {
    log.Fatal(err)
}

// 使用 MapTo 映射到现有对象
var existingDTO UserDTO
err = mapster.MapTo(user, &existingDTO)
if err != nil {
    log.Fatal(err)
}
```

### 集合映射
```go
// 切片映射
users := []User{{Name: "Alice", Age: 25}, {Name: "Bob", Age: 30}}
userDTOs, err := mapster.Map[[]UserDTO](users)

// Map 映射
userMap := map[string]User{"u1": {Name: "Alice", Age: 25}}
dtoMap, err := mapster.Map[map[string]UserDTO](userMap)
```

## 总结

go-mapster 是一个现代化的 Go 对象映射库，充分利用了 Go 语言的泛型特性，提供类型安全的对象转换功能。项目采用模块化设计，具有良好的扩展性。虽然目前核心实现尚未完成，但从架构设计和接口定义来看，该项目具有清晰的设计思路和良好的发展潜力。
