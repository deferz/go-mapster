# Go-Mapster

Go-Mapster 是一个高性能的 Go 结构体映射库，使用 Go 1.18+ 泛型特性提供类型安全的对象映射。

## 特点

- **类型安全**：使用 Go 1.18+ 泛型特性
- **高性能**：优化的映射算法
- **类型注册**：通过显式注册类型映射关系，提高安全性和性能
- **简洁 API**：简单易用的 API 设计
- **全面支持**：支持结构体、切片、数组和映射等类型

## 安装

```bash
go get github.com/deferz/go-mapster
```

## 基本用法

### 注册类型映射

在使用 Go-Mapster 进行映射前，必须先注册源类型和目标类型之间的映射关系：

```go
// 注册从 User 到 UserDTO 的映射
mapster.NewMapperConfig[User, UserDTO]().Register()
```

### 映射到新对象

```go
package main

import (
    "fmt"
    "github.com/deferz/go-mapster"
)

type User struct {
    ID   int
    Name string
    Age  int
}

type UserDTO struct {
    ID   int
    Name string
    Age  int
}

func main() {
    // 注册映射关系
    mapster.NewMapperConfig[User, UserDTO]().Register()
    
    // 创建源对象
    user := User{ID: 1, Name: "张三", Age: 30}
    
    // 映射到新对象
    userDTO, err := mapster.Map[UserDTO](user)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User DTO: %+v\n", userDTO)
}
```

### 映射到现有对象

```go
func main() {
    // 注册映射关系
    mapster.NewMapperConfig[User, UserDTO]().Register()
    
    // 创建源对象
    user := User{ID: 1, Name: "张三", Age: 30}
    
    // 创建目标对象
    userDTO := UserDTO{}
    
    // 映射到现有对象
    err := mapster.MapTo(user, &userDTO)
    if err != nil {
        panic(err)
    }
    
    fmt.Printf("User DTO: %+v\n", userDTO)
}
```

### Map类型映射

对于Map类型，只需要注册值类型的映射关系，键类型会自动处理：

```go
package main

import (
    "fmt"
    "github.com/deferz/go-mapster"
)

type SourceValue struct {
    ID   int
    Name string
}

type TargetValue struct {
    ID   int
    Name string
}

func main() {
    // 只需要注册值类型的映射关系
    mapster.NewMapperConfig[SourceValue, TargetValue]().Register()
    
    // 创建源Map
    srcMap := map[string]SourceValue{
        "key1": {ID: 1, Name: "值1"},
        "key2": {ID: 2, Name: "值2"},
    }
    
    // 映射到新Map，注意我们没有注册键类型的映射关系
    dstMap, err := mapster.Map[map[string]TargetValue](srcMap)
    if err != nil {
        panic(err)
    }
    
    // 输出结果
    for k, v := range dstMap {
        fmt.Printf("键: %s, 值: %+v\n", k, v)
    }
    
    // 对于可转换的键类型也能自动处理
    intKeyMap := map[int]SourceValue{
        1: {ID: 1, Name: "值1"},
        2: {ID: 2, Name: "值2"},
    }
    
    // int 键类型可以自动转换为 int64
    int64KeyMap, err := mapster.Map[map[int64]TargetValue](intKeyMap)
    if err != nil {
        panic(err)
    }
    
    for k, v := range int64KeyMap {
        fmt.Printf("键: %d, 值: %+v\n", k, v)
    }
}
```

## 注意事项

1. **必须注册映射关系**：在使用 `Map` 或 `MapTo` 函数前，必须先通过 `NewMapperConfig[T, R]().Register()` 注册类型映射关系
2. **类型安全**：使用泛型确保类型安全，编译时检查类型匹配
3. **类型安全的指针处理**：`MapTo` 函数的目标参数类型为 `*T`，确保编译时类型安全
4. **Map类型映射**：对于Map类型，只需要注册值类型的映射关系，键类型会自动处理

## 性能

Go-Mapster 在各种映射场景中都表现出色：

- **基本结构体映射**：比大多数映射库快 2-5 倍
- **嵌套结构体映射**：高效处理复杂嵌套结构
- **集合映射**：优化的切片和数组映射

## 未来计划

- **自定义字段映射**：支持源字段到目标字段的自定义映射
- **值转换器**：支持自定义字段值转换
- **映射条件**：支持条件映射
- **更多缓存优化**：进一步提高性能

## 许可证

MIT License