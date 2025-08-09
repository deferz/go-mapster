# go-mapster

go-mapster 是一个用于 Go 语言的高性能对象映射库，利用泛型提供类型安全的对象转换功能。

## 特性

- 🚀 **类型安全**: 使用 Go 1.18+ 泛型，在编译时进行类型检查
- 🎯 **简单易用**: 简洁的 API 设计，只需一行代码完成映射
- 🔧 **灵活配置**: 预留了配置接口设计（后续版本）
- 📦 **丰富的类型支持**: 结构体、切片、数组、Map、指针等
- 🏗️ **嵌入字段支持**: 自动处理 Go 的匿名字段
- ⚡ **高性能**: 最小化反射使用，优化性能

## 安装

```bash
go get github.com/deferz/go-mapster
```

## 快速开始

### 基本使用

```go
package main

import (
    "fmt"
    "log"
    mapster "github.com/deferz/go-mapster"
)

type User struct {
    Name  string
    Email string
    Age   int
}

type UserDTO struct {
    Name  string
    Email string
    Age   int
}

func main() {
    user := User{
        Name:  "张三",
        Email: "zhangsan@example.com",
        Age:   25,
    }

    // 映射到新对象
    userDTO, err := mapster.Map[UserDTO](user)
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%+v\n", userDTO)
}
```

### 映射到现有对象

```go
var userDTO UserDTO
err := mapster.MapTo(user, &userDTO)
if err != nil {
    log.Fatal(err)
}
```

## 高级功能

### 字段映射规则

当前版本使用自动字段名匹配：
- 优先精确匹配字段名
- 支持大小写不敏感匹配
- 如果找不到对应字段，会在嵌入字段中查找
- 未来版本将支持通过 API 配置自定义映射规则

### 支持的类型

#### 切片映射

```go
users := []User{{Name: "张三"}, {Name: "李四"}}
userDTOs, err := mapster.Map[[]UserDTO](users)
```

#### Map 映射

```go
userMap := map[string]User{
    "u1": {Name: "张三"},
    "u2": {Name: "李四"},
}
dtoMap, err := mapster.Map[map[string]UserDTO](userMap)
```

#### 嵌套结构体

```go
type Person struct {
    Name    string
    Address Address // 嵌套结构体
}

type PersonDTO struct {
    Name    string
    Address Address
}

// 自动映射嵌套结构
dto, err := mapster.Map[PersonDTO](person)
```

### 嵌入字段支持

go-mapster 自动处理 Go 的匿名字段（嵌入字段）：

```go
type BaseInfo struct {
    ID        int
    CreatedAt time.Time
}

type User struct {
    BaseInfo  // 嵌入字段
    Name      string
    Email     string
}

type UserDTO struct {
    BaseInfo  // 相同的嵌入字段
    Name      string
    Email     string
}

// 嵌入字段会自动映射
dto, err := mapster.Map[UserDTO](user)
```

## API 文档

### Map[T any](src any) (T, error)

将源对象映射到目标类型并返回新实例。

**参数:**
- `src`: 源对象，可以是任何类型

**返回:**
- `T`: 目标类型的新实例
- `error`: 如果映射失败则返回错误

### MapTo[T any](src any, dst *T) error

将源对象映射到现有的目标对象。

**参数:**
- `src`: 源对象，可以是任何类型
- `dst`: 指向目标对象的指针

**返回:**
- `error`: 如果映射失败则返回错误

## 错误处理

go-mapster 提供清晰的错误信息：

```go
// 源对象为 nil
_, err := mapster.Map[UserDTO](nil)
// 错误: 源对象不能为 nil

// 类型不兼容
_, err := mapster.Map[int]("string")
// 错误: 无法将类型 string 转换为 int
```

## 性能考虑

- 使用反射进行类型检查和字段访问
- 对于大量重复映射，考虑重用目标对象（使用 `MapTo`）
- 基本类型之间的转换使用 Go 的内置转换机制

## 限制

- 需要 Go 1.18 或更高版本（泛型支持）
- 不支持自定义转换函数（计划在未来版本中添加）
- 不支持深度复制（映射的是值，不是引用）

## 贡献

欢迎提交 Issue 和 Pull Request！

## 许可证

MIT License
