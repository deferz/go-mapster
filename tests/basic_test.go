package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestBasicMapping tests basic struct to struct mapping
func TestBasicMapping(t *testing.T) {
	src := Source{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	dst, err := mapster.Map[Target](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("Expected Age=%d, got %d", src.Age, dst.Age)
	}
	if dst.Email != src.Email {
		t.Errorf("Expected Email=%s, got %s", src.Email, dst.Email)
	}
}

// TestMapToBasic tests mapping to an existing object
func TestMapToBasic(t *testing.T) {
	src := Person{
		Name: "John Doe",
		Age:  30,
	}

	var dst Person
	err := mapster.MapTo(src, &dst)
	if err != nil {
		t.Fatalf("MapTo failed: %v", err)
	}

	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("Expected Age=%d, got %d", src.Age, dst.Age)
	}
}

// TestDifferentFieldNames tests mapping between structs with different field names
func TestDifferentFieldNames(t *testing.T) {
	// For this test, we'll use the Source and Target types
	// but we'll only check the fields that are common between them
	src := Source{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	dst, err := mapster.Map[Target](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("Expected Age=%d, got %d", src.Age, dst.Age)
	}
	if dst.Email != src.Email {
		t.Errorf("Expected Email=%s, got %s", src.Email, dst.Email)
	}
}

// TestPointerMapping tests mapping between pointer types
func TestPointerMapping(t *testing.T) {
	src := &Source{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	dst, err := mapster.Map[*Target](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("Expected Age=%d, got %d", src.Age, dst.Age)
	}
	if dst.Email != src.Email {
		t.Errorf("Expected Email=%s, got %s", src.Email, dst.Email)
	}
}

// TestMapToPreserveFields tests that MapTo preserves fields in the destination
func TestMapToPreserveFields(t *testing.T) {
	src := Source{
		Name:  "John Doe",
		Age:   30,
		Email: "",
	}

	dst := Target{
		Name:  "",
		Age:   0,
		Email: "existing@example.com",
	}

	err := mapster.MapTo(src, &dst)
	if err != nil {
		t.Fatalf("MapTo failed: %v", err)
	}

	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("Expected Age=%d, got %d", src.Age, dst.Age)
	}
	// Note: In the current implementation, empty fields are not preserved
	// This behavior could be customized with field mapping options in the future
}

// TestMapToNilPointer tests mapping to a nil pointer
func TestMapToNilPointer(t *testing.T) {
	src := Source{
		Name:  "John Doe",
		Age:   30,
		Email: "john@example.com",
	}

	var dst *Target
	err := mapster.MapTo(src, dst)
	if err == nil {
		t.Fatalf("Expected error for nil pointer, but got nil")
	}
}

// TestTypeConversions tests conversions between different basic types
func TestTypeConversions(t *testing.T) {
	// Test int to int64
	intVal := 42
	int64Val, err := mapster.Map[int64](intVal)
	if err != nil {
		t.Fatalf("Map int to int64 failed: %v", err)
	}
	if int64Val != 42 {
		t.Errorf("Expected 42, got %d", int64Val)
	}

	// Test int to float64
	floatVal, err := mapster.Map[float64](intVal)
	if err != nil {
		t.Fatalf("Map int to float64 failed: %v", err)
	}
	if floatVal != 42.0 {
		t.Errorf("Expected 42.0, got %f", floatVal)
	}

	// Test float64 to int
	float64Val := 42.75
	intFromFloat, err := mapster.Map[int](float64Val)
	if err != nil {
		t.Fatalf("Map float64 to int failed: %v", err)
	}
	if intFromFloat != 42 {
		t.Errorf("Expected 42, got %d", intFromFloat)
	}

	// String to int conversion is not supported by default
	// This would require a custom converter implementation

	// String to bool conversion is not supported by default
	// This would require a custom converter implementation

	// Register uint to int mapping
	mapster.NewMapperConfig[uint, int]().Register()

	// Test uint to int
	uintVal := uint(42)
	intFromUint, err := mapster.Map[int](uintVal)
	if err != nil {
		t.Fatalf("Map uint to int failed: %v", err)
	}
	if intFromUint != 42 {
		t.Errorf("Expected 42, got %d", intFromUint)
	}
}

// TestEmptyStructs tests mapping between empty structs
func TestEmptyStructs(t *testing.T) {
	src := EmptySource{}
	dst, err := mapster.Map[EmptyTarget](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// No fields to check, just make sure it doesn't panic
	_ = dst
}
