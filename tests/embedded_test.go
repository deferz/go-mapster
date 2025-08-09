package tests

import (
	"testing"
	"time"

	mapster "github.com/deferz/go-mapster"
)

// TestEmbeddedStructs tests mapping with embedded structs
func TestEmbeddedStructs(t *testing.T) {
	now := time.Now()
	src := SourceUser{
		BaseInfo: BaseInfo{
			ID:        123,
			CreatedAt: now,
		},
		Name:  "John Doe",
		Email: "john@example.com",
	}

	dst, err := mapster.Map[TargetUser](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// Check embedded struct fields
	if dst.ID != src.ID {
		t.Errorf("Expected ID=%d, got %d", src.ID, dst.ID)
	}
	if !dst.CreatedAt.Equal(src.CreatedAt) {
		t.Errorf("Expected CreatedAt=%v, got %v", src.CreatedAt, dst.CreatedAt)
	}

	// Check normal fields
	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
	if dst.Email != src.Email {
		t.Errorf("Expected Email=%s, got %s", src.Email, dst.Email)
	}
}

// TestEmbeddedDifferentTypes tests mapping with embedded structs of different types
func TestEmbeddedDifferentTypes(t *testing.T) {
	src := SourceUser{
		BaseInfo: BaseInfo{
			ID:        123,
			CreatedAt: time.Now(),
		},
		Name:  "John Doe",
		Email: "john@example.com",
	}

	dst, err := mapster.Map[TargetUser](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// Check embedded struct fields
	if dst.ID != src.ID {
		t.Errorf("Expected ID=%d, got %d", src.ID, dst.ID)
	}

	// Check normal fields
	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
	if dst.Email != src.Email {
		t.Errorf("Expected Email=%s, got %s", src.Email, dst.Email)
	}
}

// TestNestedStructs tests mapping with nested structs
func TestNestedStructs(t *testing.T) {
	src := SourcePerson{
		Name: "John Doe",
		Age:  30,
		Address: Address{
			Street:  "123 Main St",
			City:    "New York",
			Country: "USA",
		},
	}

	dst, err := mapster.Map[TargetPerson](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// Check normal fields
	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("Expected Age=%d, got %d", src.Age, dst.Age)
	}

	// Check nested struct fields
	if dst.Address.Street != src.Address.Street {
		t.Errorf("Expected Street=%s, got %s", src.Address.Street, dst.Address.Street)
	}
	if dst.Address.City != src.Address.City {
		t.Errorf("Expected City=%s, got %s", src.Address.City, dst.Address.City)
	}
	if dst.Address.Country != src.Address.Country {
		t.Errorf("Expected Country=%s, got %s", src.Address.Country, dst.Address.Country)
	}
}

// TestDeepNestedStructs tests mapping with deeply nested structs
func TestDeepNestedStructs(t *testing.T) {
	src := Level1{
		Level2: Level2{
			Level3: Level3{
				Value: "nested value",
			},
		},
	}

	dst, err := mapster.Map[Level1](src)
	if err != nil {
		t.Fatalf("Map deep nested structs failed: %v", err)
	}

	// Check deeply nested value
	if dst.Level2.Level3.Value != src.Level2.Level3.Value {
		t.Errorf("Expected Value=%s, got %s", src.Level2.Level3.Value, dst.Level2.Level3.Value)
	}
}

// TestEmbeddedSameType tests mapping with embedded structs of the same type
func TestEmbeddedSameType(t *testing.T) {
	src := SourceWithEmbedded{
		BaseEmbedded: BaseEmbedded{
			ID:   123,
			Name: "John Doe",
		},
		Email: "john@example.com",
	}

	dst, err := mapster.Map[TargetWithEmbedded](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// Check embedded struct fields
	if dst.ID != src.ID {
		t.Errorf("Expected ID=%d, got %d", src.ID, dst.ID)
	}
	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}

	// Check normal fields
	if dst.Email != src.Email {
		t.Errorf("Expected Email=%s, got %s", src.Email, dst.Email)
	}
}
