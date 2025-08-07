# Mapster for Go

A high-performance object mapping library for Go, inspired by .NET's Mapster. This library provides a simple and flexible way to map between different types with minimal configuration.

**[中文文档](README_zh.md)** | **English**

## ✨ Key Features

- **🚀 Zero-Reflection Code Generation**: Generate optimized mapping code at build time for 1.5x performance boost
- **🎯 Type-Safe Generic API**: Leverages Go 1.18+ generics for compile-time type safety
- **🔧 Fluent Configuration**: Intuitive, chainable API for complex mapping scenarios
- **⚡ High Performance**: Both zero-reflection generated code and optimized reflection-based mapping
- **🎭 Flexible Field Mapping**: Support for custom fields, transformations, and conditional mapping
- **📊 Deep Path Resolution**: Access nested object properties using dot notation (e.g., `Company.Address.City`)
- **🔄 Circular Reference Detection**: Safe handling of complex object graphs with circular references
- **📦 Batch Processing**: Efficient slice and array mapping capabilities
- **⏰ Smart Time Conversion**: Automatic int64 ↔ time.Time conversion with configurable behavior

## Performance

```
Zero-Reflection:  474 ns/op   312 B/op    8 allocs/op  ⭐ Generated mappers
Configuration:    490 ns/op   224 B/op    8 allocs/op  🔧 Custom config  
Reflection:       732 ns/op   320 B/op    8 allocs/op  🔄 Auto-mapping
```

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
        Map("FullName").FromFunc(func(u User) any {
            return u.FirstName + " " + u.LastName
        }).
        Map("AgeGroup").FromFunc(func(u User) any {
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

### Zero-Reflection Code Generation 🚀

For maximum performance, you can register generated mappers that avoid reflection entirely:

```go
// Generate optimized mapper function
func mapUserToUserDTO(src User) UserDTO {
    return UserDTO{
        ID:        src.ID,
        FirstName: src.FirstName,
        LastName:  src.LastName,
        Email:     src.Email,
        FullName:  src.FirstName + " " + src.LastName, // Custom logic
    }
}

func init() {
    // Register the generated mapper
    mapster.RegisterGeneratedMapper(mapUserToUserDTO)
}

func main() {
    user := User{ID: 1, FirstName: "John", LastName: "Doe"}
    
    // This will automatically use the generated mapper (1.5x faster!)
    userDTO := mapster.Map[UserDTO](user)
    fmt.Printf("Generated mapping: %+v\n", userDTO)
}
```

**Benefits**:
- 🚀 **1.5x Performance**: Direct field access instead of reflection
- 🛡️ **Type Safety**: Compile-time checking
- 🔄 **Auto Fallback**: Uses reflection if no generated mapper exists
- 🔧 **Easy Integration**: Just register the function

## API Reference

### Core Functions

- `Map[T any](src any) T` - Maps source object to target type
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

### Time Conversion API

- `EnableTimeConversion(bool)` - Enable/disable automatic time conversion
- `IsTimeConversionEnabled()` - Check if time conversion is enabled
- `SetGlobalConfig(GlobalConfig)` - Set global configuration options
- `GetGlobalConfig()` - Get current global configuration

## Examples

### Field Mapping

```go
mapster.Config[Source, Target]().
    Map("TargetField").FromField("SourceField").
    Register()
```

### Deep Path Resolution

Access nested object properties using dot notation:

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
    Name: "John Doe",
    Company: &Company{
        Name: "Tech Corp",
        Address: Address{City: "San Francisco"},
    },
}

dto := mapster.Map[EmployeeDTO](employee)
// Result: {Name: "John Doe", CompanyName: "Tech Corp", CompanyCity: "San Francisco"}
```

### Custom Functions

```go
mapster.Config[User, UserDTO]().
    Map("FullName").FromFunc(func(u User) any {
        return u.FirstName + " " + u.LastName
    }).
    Register()
```

### Circular Reference Handling

Safely handle complex object graphs:

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

// Safe mapping avoids circular references
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

// Works safely even with circular references
dto := mapster.Map[NodeDTO](nodeWithCircularRef)
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
dtos := make([]UserDTO, len(users))
for i, u := range users {
    dtos[i] = mapster.Map[UserDTO](u)
}
```

### Time Conversion

Automatic conversion between `int64` timestamps and `time.Time`:

```go
// Database model with int64 timestamps
type UserModel struct {
    ID        int64
    Name      string
    CreatedAt int64  // Unix timestamp
    UpdatedAt int64  // Unix timestamp
}

// API response with time.Time
type UserResponse struct {
    ID        int64
    Name      string
    CreatedAt time.Time
    UpdatedAt time.Time
}

func main() {
    userModel := UserModel{
        ID:        1,
        Name:      "John Doe",
        CreatedAt: time.Now().Unix(),
        UpdatedAt: time.Now().Unix(),
    }

    // Automatic conversion (enabled by default)
    userResponse := mapster.Map[UserResponse](userModel)
    fmt.Printf("CreatedAt: %s\n", userResponse.CreatedAt.Format("2006-01-02 15:04:05"))

    // Disable time conversion for performance
    mapster.EnableTimeConversion(false)
    userResponse2 := mapster.Map[UserResponse](userModel)
    fmt.Printf("CreatedAt (disabled): %s\n", userResponse2.CreatedAt.Format("2006-01-02 15:04:05"))
}
```

#### Time Conversion Configuration

```go
// Enable/disable globally
mapster.EnableTimeConversion(true)  // Default: true
mapster.EnableTimeConversion(false) // Disable for performance

// Check current status
enabled := mapster.IsTimeConversionEnabled()

// Use global configuration
mapster.SetGlobalConfig(mapster.GlobalConfig{
    EnableTimeConversion: false,
})
```

#### Custom Configuration Priority

Custom field mappings take priority over global time conversion:

```go
// Custom configuration for specific field
mapster.Config[UserModel, UserResponse]().
    Map("CreatedAt").FromFunc(func(src UserModel) any {
        return time.Date(2020, 1, 1, 12, 0, 0, 0, time.UTC)
    }).
    Register()

// CreatedAt uses custom config, UpdatedAt uses global time conversion
userResponse := mapster.Map[UserResponse](userModel)
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
- **🚀 Zero-Reflection Code Generation**: 1.5x performance boost with generated mappers
- **Basic reflection-based mapping**: Automatic field matching
- **Fluent configuration API**: Chainable configuration interface
- **Custom mapping functions**: Complex logic support
- **Slice mapping**: Batch object processing
- **Basic nested object mapping**: Automatic struct-in-struct mapping
- **⏰ Smart Time Conversion**: Automatic int64 ↔ time.Time conversion with configurable behavior

### Enhanced Features in Development 🚧
- **Deep path mapping**: Complete `FromPath("Address.Street")` implementation
- **Flattening mappings**: Smart nested-to-flat structure mapping
- **Circular reference handling**: Safe mapping without infinite recursion
- **Dynamic field mapping**: Runtime field discovery and mapping
- **Validation integration**: Data validation during mapping process
- **Advanced configuration options**: Conditional mapping, field ignoring, etc.

### Future Plans 📋
- **Compile-time code generation tools**: Automated mapper generation
- **IDE plugin support**: VS Code extensions
- **Performance analysis tools**: Mapping performance monitoring
- **Community contribution templates**: Standardized contribution workflow

Want to contribute? Check out our [Contributing Guidelines](#contributing)!

## Contributing

Contributions are welcome! Please feel free to submit issues, feature requests, or pull requests.

## License

This project is licensed under the MIT License - see the LICENSE file for details.
