package tests

import (
	"testing"
	"time"

	mapster "github.com/deferz/go-mapster"
)

// TestEmbeddedStructs tests mapping with embedded structs
func TestEmbeddedStructs(t *testing.T) {
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

	now := time.Now()
	src := SourceUser{
		BaseInfo: BaseInfo{
			ID:        123,
			CreatedAt: now,
		},
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	dst, err := mapster.Map[TargetUser](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if dst.ID != src.ID {
		t.Errorf("Expected ID=%d, got %d", src.ID, dst.ID)
	}
	if !dst.CreatedAt.Equal(src.CreatedAt) {
		t.Errorf("Expected CreatedAt=%v, got %v", src.CreatedAt, dst.CreatedAt)
	}
	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
}

func TestEmbeddedDifferentTypes(t *testing.T) {
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
		ID        int
		CreatedAt time.Time
		Name      string
		Email     string
	}

	now := time.Now()
	src := SourceUser{
		BaseInfo: BaseInfo{
			ID:        123,
			CreatedAt: now,
		},
		Name:  "Jane Doe",
		Email: "jane@example.com",
	}

	dst, err := mapster.Map[TargetUser](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if dst.ID != src.ID {
		t.Errorf("Expected ID=%d, got %d", src.ID, dst.ID)
	}
	if !dst.CreatedAt.Equal(src.CreatedAt) {
		t.Errorf("Expected CreatedAt=%v, got %v", src.CreatedAt, dst.CreatedAt)
	}
	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
}

// TestNestedStructs tests mapping with nested structs
func TestNestedStructs(t *testing.T) {
	type Address struct {
		Street string
		City   string
		ZIP    string
	}

	type SourcePerson struct {
		Name    string
		Address Address // nested struct
	}

	type TargetPerson struct {
		Name    string
		Address Address // same nested struct
	}

	src := SourcePerson{
		Name: "Alice",
		Address: Address{
			Street: "123 Main St",
			City:   "Anytown",
			ZIP:    "12345",
		},
	}

	dst, err := mapster.Map[TargetPerson](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
	if dst.Address.Street != src.Address.Street {
		t.Errorf("Expected Address.Street=%s, got %s", src.Address.Street, dst.Address.Street)
	}
	if dst.Address.City != src.Address.City {
		t.Errorf("Expected Address.City=%s, got %s", src.Address.City, dst.Address.City)
	}
}

// TestDeepNestedStructs tests mapping with deeply nested structs
func TestDeepNestedStructs(t *testing.T) {
	type Level3 struct {
		Value string
	}

	type Level2 struct {
		Name  string
		Deep  Level3
		Count int
	}

	type Level1 struct {
		ID     int
		Nested Level2
	}

	src := Level1{
		ID: 1,
		Nested: Level2{
			Name:  "Level 2",
			Count: 42,
			Deep: Level3{
				Value: "Deep value",
			},
		},
	}

	dst, err := mapster.Map[Level1](src)
	if err != nil {
		t.Fatalf("Map deep nested structs failed: %v", err)
	}

	if dst.ID != src.ID {
		t.Errorf("Expected ID=%d, got %d", src.ID, dst.ID)
	}
	if dst.Nested.Name != src.Nested.Name {
		t.Errorf("Expected Nested.Name=%s, got %s", src.Nested.Name, dst.Nested.Name)
	}
	if dst.Nested.Count != src.Nested.Count {
		t.Errorf("Expected Nested.Count=%d, got %d", src.Nested.Count, dst.Nested.Count)
	}
	if dst.Nested.Deep.Value != src.Nested.Deep.Value {
		t.Errorf("Expected Nested.Deep.Value=%s, got %s", src.Nested.Deep.Value, dst.Nested.Deep.Value)
	}
}

// TestEmbeddedSameType tests mapping with embedded structs of same type
func TestEmbeddedSameType(t *testing.T) {
	type BaseInfo struct {
		ID   int
		Code string
	}

	type SourceWithEmbedded struct {
		BaseInfo // Same type in both structs
		Name     string
	}

	type TargetWithEmbedded struct {
		BaseInfo // Same type in both structs
		Name     string
	}

	src := SourceWithEmbedded{
		BaseInfo: BaseInfo{
			ID:   123,
			Code: "ABC",
		},
		Name: "Test Item",
	}

	dst, err := mapster.Map[TargetWithEmbedded](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// ID should be mapped from embedded struct
	if dst.ID != src.ID {
		t.Errorf("Expected ID=%d, got %d", src.ID, dst.ID)
	}

	// Code should be mapped
	if dst.Code != src.Code {
		t.Errorf("Expected Code=%s, got %s", src.Code, dst.Code)
	}

	// Name should be mapped
	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
}
