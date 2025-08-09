package tests

import (
	"testing"

	"github.com/deferz/go-mapster"
)

// Source and target types for testing
type SimpleSource struct {
	ID   int
	Name string
	Age  int
}

type SimpleTarget struct {
	ID   int
	Name string
	Age  int
}

// Test basic CreateMap functionality
func TestCreateMap(t *testing.T) {
	// Register mapping configuration
	mapster.NewMapperConfig[SimpleSource, SimpleTarget]().Register()

	// Create test data
	source := SimpleSource{
		ID:   1,
		Name: "John Doe",
		Age:  30,
	}

	// Perform mapping
	target, err := mapster.Map[SimpleTarget](source)
	if err != nil {
		t.Fatalf("Mapping failed: %v", err)
	}

	// Verify mapping results
	if target.ID != source.ID {
		t.Errorf("Expected ID=%d, got %d", source.ID, target.ID)
	}

	if target.Name != source.Name {
		t.Errorf("Expected Name=%s, got %s", source.Name, target.Name)
	}

	if target.Age != source.Age {
		t.Errorf("Expected Age=%d, got %d", source.Age, target.Age)
	}
}

// Test MapTo functionality
func TestMapTo(t *testing.T) {
	// Register mapping configuration
	mapster.NewMapperConfig[SimpleSource, SimpleTarget]().Register()

	// Create test data
	source := SimpleSource{
		ID:   1,
		Name: "John Doe",
		Age:  30,
	}

	// Create target object
	var target SimpleTarget

	// Perform mapping
	err := mapster.MapTo(source, &target)
	if err != nil {
		t.Fatalf("Mapping failed: %v", err)
	}

	// Verify mapping results
	if target.ID != source.ID {
		t.Errorf("Expected ID=%d, got %d", source.ID, target.ID)
	}

	if target.Name != source.Name {
		t.Errorf("Expected Name=%s, got %s", source.Name, target.Name)
	}

	if target.Age != source.Age {
		t.Errorf("Expected Age=%d, got %d", source.Age, target.Age)
	}
}

// Test mapping without registration
func TestMappingWithoutRegistration(t *testing.T) {
	// Create a new type that hasn't been registered
	type UnregisteredSource struct {
		ID   int
		Name string
	}

	// Create test data
	source := UnregisteredSource{
		ID:   1,
		Name: "John Doe",
	}

	// Attempt mapping without registration
	_, err := mapster.Map[SimpleTarget](source)

	// Should fail with error about no registered mapping
	if err == nil {
		t.Errorf("Expected error for unregistered mapping, but got nil")
	}
}
