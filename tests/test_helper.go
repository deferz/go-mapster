package tests

import (
	"testing"
	"time"

	mapster "github.com/deferz/go-mapster"
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

// RegisterTestTypes registers all the test types used in the test suite
func RegisterTestTypes(t *testing.T) {
	// From basic_test.go
	mapster.NewMapperConfig[Source, Target]().Register()
	mapster.NewMapperConfig[Person, Person]().Register()
	mapster.NewMapperConfig[*Source, *Target]().Register()
	mapster.NewMapperConfig[int, int64]().Register()
	mapster.NewMapperConfig[int, float64]().Register()
	mapster.NewMapperConfig[float64, int]().Register()
	// String conversions removed as they require custom converters
	mapster.NewMapperConfig[EmptySource, EmptyTarget]().Register()
	mapster.NewMapperConfig[uint, int]().Register()

	// From collection_test.go
	mapster.NewMapperConfig[Item, Item]().Register()
	// No need to register slice and array types separately now
	// No need to register map types separately now, just register key and value types
	mapster.NewMapperConfig[string, string]().Register() // For map keys
	mapster.NewMapperConfig[User, User]().Register()     // For map values

	// From embedded_test.go
	mapster.NewMapperConfig[BaseInfo, BaseInfo]().Register()
	mapster.NewMapperConfig[SourceUser, TargetUser]().Register()
	mapster.NewMapperConfig[Address, Address]().Register()
	mapster.NewMapperConfig[SourcePerson, TargetPerson]().Register()
	mapster.NewMapperConfig[Level3, Level3]().Register()
	mapster.NewMapperConfig[Level2, Level2]().Register()
	mapster.NewMapperConfig[Level1, Level1]().Register()
	mapster.NewMapperConfig[BaseEmbedded, BaseEmbedded]().Register()
	mapster.NewMapperConfig[SourceWithEmbedded, TargetWithEmbedded]().Register()
}

// init is called before any tests are run
func init() {
	// Register all test types
	RegisterTestTypes(nil)
}
