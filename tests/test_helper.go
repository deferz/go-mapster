package tests

import (
	"time"
)

// Source and Target for basic tests
type Source struct {
	Name  string
	Age   int
	Email string
}

type Target struct {
	Name  string
	Age   int
	Email string
}

// Person for pointer tests
type Person struct {
	Name string
	Age  int
}

// EmptySource and EmptyTarget for empty struct tests
type EmptySource struct{}
type EmptyTarget struct{}

// Item for collection tests
type Item struct {
	ID   int
	Name string
}

// User for map tests
type User struct {
	ID   int
	Name string
}

// Types for embedded struct tests
type BaseInfo struct {
	ID        int
	CreatedAt time.Time
}

type SourceUser struct {
	BaseInfo // embedded struct
	Name     string
	Email    string
}

type TargetUser struct {
	BaseInfo // same embedded struct
	Name     string
	Email    string
}

type Address struct {
	Street  string
	City    string
	Country string
}

type SourcePerson struct {
	Name    string
	Age     int
	Address Address
}

type TargetPerson struct {
	Name    string
	Age     int
	Address Address
}

type Level3 struct {
	Value string
}

type Level2 struct {
	Level3 Level3
}

type Level1 struct {
	Level2 Level2
}

type BaseEmbedded struct {
	ID   int
	Name string
}

type SourceWithEmbedded struct {
	BaseEmbedded
	Email string
}

type TargetWithEmbedded struct {
	BaseEmbedded
	Email string
}

// 注意：由于映射现在在调用 Map 和 MapTo 时自动注册，我们不再需要显式注册类型映射关系
