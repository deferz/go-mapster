# Go Mapster

一个高性能的 Go 对象映射库，灵感来自 .NET 的 Mapster。这个库提供了简单灵活的方式来映射不同类型，配置最少。

**中文** | **[English](README.md)**

## 特性

- **🚀 手动零反射映射**：用户编写映射函数，性能提升 1.5 倍
- **零配置**：大多数映射场景通过自动字段匹配开箱即用
- **流畅的配置 API**：使用链式 API 轻松配置自定义映射
- **高性能**：优化的反射映射，智能缓存机制
- **类型安全**：利用 Go 泛型实现编译时类型检查
- **灵活**：支持自定义映射函数、转换和条件映射
- **📊 深度路径解析**：使用点标记法访问嵌套对象属性（如 `Company.Address.City`）
- **🔄 循环引用检测**：安全处理包含循环引用的复杂对象图
- **⏰ 智能时间转换**：自动 int64 ↔ time.Time 转换，支持可配置行为

## 性能表现

```
零反射映射:    474 ns/op   312 B/op    8 allocs/op  ⭐ 生成代码
配置映射:      490 ns/op   224 B/op    8 allocs/op  🔧 自定义配置
反射映射:      732 ns/op   320 B/op    8 allocs/op  🔄 自动映射
```

## 安装

```bash
go get github.com/deferz/go-mapster
```

## 快速开始

### 基础映射

```go
package main

import (
    "fmt"
    "github.com/deferz/go-mapster"
)

type User struct {
    ID        int64
    FirstName string
    LastName  string
    Email     string
    Age       int
}

type UserDTO struct {
    ID        int64
    FirstName string
    LastName  string
    Email     string
}

func main() {
    user := User{
        ID:        1,
        FirstName: "张",
        LastName:  "三",
        Email:     "zhangsan@example.com",
        Age:       30,
    }

    // 零配置映射
    dto := mapster.Map[UserDTO](user)
    fmt.Printf("映射结果: %+v\n", dto)
}
```

### 自定义映射配置

```go
func init() {
    // 配置自定义映射
    mapster.Config[User, UserDTO]().
        Map("FullName").FromFunc(func(u User) any {
            return u.FirstName + u.LastName
        }).
        Map("AgeGroup").FromFunc(func(u User) any {
            if u.Age < 18 {
                return "未成年"
            } else if u.Age < 65 {
                return "成年人"
            }
            return "老年人"
        }).
        Register()
}
```

### 手动零反射映射 🚀

为了获得最佳性能，你可以编写自己的映射函数并注册它们来完全避免反射：

```go
// 生成优化的映射函数
func mapUserToUserDTO(src User) UserDTO {
    return UserDTO{
        ID:        src.ID,
        FirstName: src.FirstName,
        LastName:  src.LastName,
        Email:     src.Email,
        FullName:  src.FirstName + " " + src.LastName, // 自定义逻辑
    }
}

func init() {
    // 注册你的手动映射器
mapster.RegisterGeneratedMapper(mapUserToUserDTO)
}

func main() {
    user := User{ID: 1, FirstName: "张", LastName: "三"}
    
    // 这会自动使用你的手动映射器（快 1.5 倍！）
    userDTO := mapster.Map[UserDTO](user)
    fmt.Printf("手动映射结果: %+v\n", userDTO)
}
```

**优势**：
- 🚀 **1.5 倍性能**：直接字段访问而非反射
- 🛡️ **类型安全**：编译时检查
- 🔄 **自动回退**：没有生成映射器时使用反射
- 🔧 **简单集成**：只需注册函数

## API 参考

### 核心函数

- `Map[T any](src any) T` - 将源对象映射到目标类型
- `MapTo[T any](src any, target *T)` - 映射到现有对象

### 配置 API

- `Config[S, T any]()` - 开始配置源类型和目标类型的映射
- `Map(field)` - 配置特定字段的映射
- `FromField(field)` - 从不同名称的源字段映射
- `FromFunc(func)` - 使用自定义映射函数
- `FromPath(path)` - 从嵌套字段映射（如 "Customer.Name"）
- `Transform(func)` - 对映射值应用转换
- `When(condition)` - 添加条件映射
- `Ignore(field)` - 忽略特定字段
- `Register()` - 注册配置

## 示例

### 字段映射

```go
mapster.Config[Source, Target]().
    Map("目标字段").FromField("源字段").
    Register()
```

### 深度路径解析

使用点标记法访问嵌套对象属性：

```go
type Employee struct {
    Name    string
    Company *Company
}

type Company struct {
    Name    string
    Address Address
}

type EmployeeDTO struct {
    Name        string
    CompanyName string
    CompanyCity string
}

mapster.Config[Employee, EmployeeDTO]().
    Map("CompanyName").FromPath("Company.Name").
    Map("CompanyCity").FromPath("Company.Address.City").
    Register()

employee := Employee{
    Name: "张三",
    Company: &Company{
        Name: "科技公司",
        Address: Address{City: "北京"},
    },
}

dto := mapster.Map[EmployeeDTO](employee)
// 结果: {Name: "张三", CompanyName: "科技公司", CompanyCity: "北京"}
```

### 自定义函数

```go
mapster.Config[User, UserDTO]().
    Map("FullName").FromFunc(func(u User) any {
        return u.FirstName + u.LastName
    }).
    Register()
```

### 循环引用处理

安全处理复杂的对象图：

```go
type Node struct {
    ID       int
    Name     string
    Parent   *Node
    Children []*Node
}

type NodeDTO struct {
    ID         int
    Name       string
    ParentName string
    ChildCount int
}

// 安全映射避免循环引用
mapster.Config[Node, NodeDTO]().
    Map("ParentName").FromFunc(func(n Node) any {
        if n.Parent != nil {
            return n.Parent.Name
        }
        return ""
    }).
    Map("ChildCount").FromFunc(func(n Node) any {
        return len(n.Children)
    }).
    Register()

// 即使有循环引用也能安全工作
dto := mapster.Map[NodeDTO](nodeWithCircularRef)
```

### 转换

```go
mapster.Config[Order, OrderDTO]().
    Map("FormattedDate").FromField("CreatedAt").Transform(func(t time.Time) string {
        return t.Format("2006-01-02")
    }).
    Register()
```

### 条件映射

```go
mapster.Config[User, UserDTO]().
    Map("Email").When(func(u User) bool {
        return u.Email != ""
    }).FromField("Email").
    Register()
```

### 切片映射

```go
users := []User{user1, user2, user3}
dtos := make([]UserDTO, len(users))
for i, u := range users {
    dtos[i] = mapster.Map[UserDTO](u)
}
```

## 为什么选择 Go Mapster？

- 🚀 **零学习成本**：如果你懂 Go 结构体，你就懂 Mapster
- ⚡ **高性能**：优化的反射使用，未来支持代码生成
- 🛡️ **类型安全**：通过泛型实现编译时类型检查
- 🔧 **灵活**：处理从简单到复杂的映射场景
- 📦 **零依赖**：纯 Go 实现

## 性能

Go Mapster 针对高性能场景进行了优化：

- **快速映射**：基础结构体映射 ~1.2μs 每次操作
- **内存高效**：最少分配（每次操作 8-12 次分配）
- **智能缓存**：反射元数据被缓存以供重复使用
- **面向未来**：设计支持代码生成以消除反射开销

```go
// 基准测试结果 (Apple M1):
// BenchmarkBasicMapping-8     927649    1199 ns/op    416 B/op    12 allocs/op
// BenchmarkSliceMapping-8       9754  120473 ns/op  51115 B/op  1202 allocs/op
```

## 路线图

### 当前状态 ✅
- **🚀 手动零反射映射**：用户编写映射函数，性能提升 1.5 倍
- **📊 深度路径解析**：完整的 `FromPath("Address.Street")` 实现，支持点标记法
- **🔄 循环引用检测**：安全的复杂对象图处理，支持指针追踪
- **🎭 条件映射**：`When()` 条件和 `Ignore()` 字段排除
- **🔧 转换映射**：`Transform()` 值转换函数
- **⏰ 智能时间转换**：自动 int64 ↔ time.Time 转换，支持可配置行为
- **🔧 流畅的配置 API**：链式配置接口
- **🎯 自定义映射函数**：复杂逻辑支持，类型安全
- **📦 切片映射**：批量对象处理能力
- **🛡️ 类型安全的泛型 API**：利用 Go 1.18+ 泛型实现编译时类型检查

### 未来计划 📋
- **代码生成工具**：编译时映射代码生成
- **IDE 插件支持**：VS Code 扩展

想要贡献？查看我们的[贡献指南](#贡献)！

## 贡献

欢迎贡献！请随时提交 issue、功能请求或 pull request。

## 许可证

本项目采用 MIT 许可证 - 查看 LICENSE 文件了解详情。
