# Mapster for Go

A high-performance object mapping library for Go, inspired by .NET's Mapster. This library provides a simple and flexible way to map between different types with minimal configuration.

**[中文文档](README_zh.md)** | **English**

## Features

- **Zero Configuration**: Most mapping scenarios work out of the box with automatic field matching
- **Fluent Configuration API**: Easy to configure custom mappings using a chainable API
- **High Performance**: Uses reflection efficiently with future code generation support
- **Type Safe**: Leverages Go's generics for compile-time type safety
- **Flexible**: Supports custom mapping functions, transformations, and conditional mapping

## Installation

```bash
go get github.com/deferz/go-mapster
```

## Quick Start

### Basic Mapping

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
        FirstName: "John",
        LastName:  "Doe",
        Email:     "john@example.com",
        Age:       30,
    }

    // Zero configuration mapping
    dto := mapster.Map[UserDTO](user)
    fmt.Printf("Mapped: %+v\n", dto)
}
```

### Custom Mapping Configuration

```go
func init() {
    // Configure custom mappings
    mapster.Config[User, UserDTO]().
        Map("FullName").FromFunc(func(u User) interface{} {
            return u.FirstName + " " + u.LastName
        }).
        Map("AgeGroup").FromFunc(func(u User) interface{} {
            if u.Age < 18 {
                return "Minor"
            } else if u.Age < 65 {
                return "Adult"
            }
            return "Senior"
        }).
        Register()
}
```

## API Reference

### Core Functions

- `Map[T any](src any) T` - Maps source object to target type
- `MapSlice[T any](src any) []T` - Maps slice of objects
- `MapTo[T any](src any, target *T)` - Maps to existing object

### Configuration API

- `Config[S, T any]()` - Starts configuration for source and target types
- `Map(field)` - Configures mapping for a specific field
- `FromField(field)` - Maps from a different source field name
- `FromFunc(func)` - Uses custom mapping function
- `FromPath(path)` - Maps from nested field (e.g., "Customer.Name")
- `Transform(func)` - Applies transformation to mapped value
- `When(condition)` - Adds conditional mapping
- `Ignore(field)` - Ignores specific field
- `Register()` - Registers the configuration

## Examples

### Field Mapping

```go
mapster.Config[Source, Target]().
    Map("TargetField").FromField("SourceField").
    Register()
```

### Custom Functions

```go
mapster.Config[User, UserDTO]().
    Map("FullName").FromFunc(func(u User) interface{} {
        return u.FirstName + " " + u.LastName
    }).
    Register()
```

### Transformations

```go
mapster.Config[Order, OrderDTO]().
    Map("FormattedDate").FromField("CreatedAt").Transform(func(t time.Time) string {
        return t.Format("2006-01-02")
    }).
    Register()
```

### Conditional Mapping

```go
mapster.Config[User, UserDTO]().
    Map("Email").When(func(u User) bool {
        return u.Email != ""
    }).FromField("Email").
    Register()
```

### Slice Mapping

```go
users := []User{user1, user2, user3}
dtos := mapster.MapSlice[UserDTO](users)
```

## Why Choose Mapster for Go?

- 🚀 **Zero Learning Curve**: If you know Go structs, you know Mapster
- ⚡ **High Performance**: Optimized reflection with future code generation
- 🛡️ **Type Safe**: Compile-time type checking with generics
- 🔧 **Flexible**: Handle simple to complex mapping scenarios
- 📦 **Zero Dependencies**: Pure Go implementation

## Performance

Mapster for Go is optimized for high-performance scenarios:

- **Fast Mapping**: ~1.2μs per operation for basic struct mapping
- **Memory Efficient**: Minimal allocations (8-12 per operation)
- **Smart Caching**: Reflection metadata is cached for repeated use
- **Future Proof**: Designed for code generation to eliminate reflection overhead

```go
// Benchmark results (Apple M1):
// BenchmarkBasicMapping-8     927649    1199 ns/op    416 B/op    12 allocs/op
// BenchmarkSliceMapping-8       9754  120473 ns/op  51115 B/op  1202 allocs/op
```

## Roadmap

### Current Status ✅
- Basic reflection-based mapping
- Fluent configuration API  
- Custom mapping functions
- Slice mapping

### Coming Soon 🚧
- Code generation for zero-reflection mapping
- Enhanced nested object mapping
- Validation integration
- Comprehensive benchmarks
- Additional configuration options

Want to contribute? Check out our [Contributing Guidelines](#contributing)!

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
