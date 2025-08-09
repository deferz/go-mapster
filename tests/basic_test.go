package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestBasicMapping tests basic struct to struct mapping
func TestBasicMapping(t *testing.T) {
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

// TestMapToBasic tests basic MapTo functionality
func TestMapToBasic(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

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
	type Source struct {
		FirstName string
		LastName  string
		Age       int
	}

	type Target struct {
		// Different field names
		Name string // We expect this to remain empty since there's no matching field
		Age  int    // This should map correctly
	}

	src := Source{
		FirstName: "John",
		LastName:  "Doe",
		Age:       30,
	}

	dst, err := mapster.Map[Target](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// Name should be empty because no matching field exists
	if dst.Name != "" {
		t.Errorf("Expected Name to be empty, got %s", dst.Name)
	}

	// Age should map correctly
	if dst.Age != src.Age {
		t.Errorf("Expected Age=%d, got %d", src.Age, dst.Age)
	}
}

// TestPointerMapping tests mapping with pointer fields
func TestPointerMapping(t *testing.T) {
	type Source struct {
		Name *string
		Age  *int
	}

	type Target struct {
		Name *string
		Age  *int
	}

	name := "John Doe"
	age := 30

	src := Source{
		Name: &name,
		Age:  &age,
	}

	dst, err := mapster.Map[Target](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if dst.Name == nil {
		t.Fatal("Expected Name not to be nil")
	}
	if dst.Age == nil {
		t.Fatal("Expected Age not to be nil")
	}

	if *dst.Name != name {
		t.Errorf("Expected Name=%s, got %s", name, *dst.Name)
	}
	if *dst.Age != age {
		t.Errorf("Expected Age=%d, got %d", age, *dst.Age)
	}

	// Ensure pointers are different (deep copy)
	if dst.Name == src.Name {
		t.Error("Expected Name pointers to be different")
	}
	if dst.Age == src.Age {
		t.Error("Expected Age pointers to be different")
	}
}

// TestMapToPreserveFields tests that MapTo preserves fields in the destination
// that don't exist in the source
func TestMapToPreserveFields(t *testing.T) {
	type Source struct {
		Name string
	}

	type Target struct {
		Name    string
		Age     int
		Country string
	}

	src := Source{
		Name: "John Doe",
	}

	dst := Target{
		Age:     30,
		Country: "USA",
	}

	err := mapster.MapTo(src, &dst)
	if err != nil {
		t.Fatalf("MapTo failed: %v", err)
	}

	// Name should be updated
	if dst.Name != src.Name {
		t.Errorf("Expected Name=%s, got %s", src.Name, dst.Name)
	}

	// Age and Country should be preserved
	if dst.Age != 30 {
		t.Errorf("Expected Age=30, got %d", dst.Age)
	}
	if dst.Country != "USA" {
		t.Errorf("Expected Country=USA, got %s", dst.Country)
	}
}

// TestMapToNilPointer tests MapTo with nil pointer fields
func TestMapToNilPointer(t *testing.T) {
	type Source struct {
		Name *string
		Age  *int
	}

	type Target struct {
		Name *string
		Age  *int
	}

	name := "John Doe"
	src := Source{
		Name: &name,
		Age:  nil, // Nil pointer
	}

	var dst Target
	err := mapster.MapTo(src, &dst)
	if err != nil {
		t.Fatalf("MapTo failed: %v", err)
	}

	// Name should be mapped
	if dst.Name == nil {
		t.Fatal("Expected Name not to be nil")
	}
	if *dst.Name != name {
		t.Errorf("Expected Name=%s, got %s", name, *dst.Name)
	}

	// Age should be nil
	if dst.Age != nil {
		t.Errorf("Expected Age to be nil, got %d", *dst.Age)
	}
}

// TestTypeConversions tests basic type conversions
func TestTypeConversions(t *testing.T) {
	// Int to int64
	intVal := 42
	int64Val, err := mapster.Map[int64](intVal)
	if err != nil {
		t.Fatalf("Map int to int64 failed: %v", err)
	}
	if int64Val != 42 {
		t.Errorf("Expected int64Val=42, got %d", int64Val)
	}

	// Int to float64
	floatVal, err := mapster.Map[float64](intVal)
	if err != nil {
		t.Fatalf("Map int to float64 failed: %v", err)
	}
	if floatVal != 42.0 {
		t.Errorf("Expected floatVal=42.0, got %f", floatVal)
	}

	// Float to int
	float64Val := 42.5
	intFromFloat, err := mapster.Map[int](float64Val)
	if err != nil {
		t.Fatalf("Map float64 to int failed: %v", err)
	}
	if intFromFloat != 42 {
		t.Errorf("Expected intFromFloat=42, got %d", intFromFloat)
	}

	// Uint to int
	uintVal := uint(100)
	intFromUint, err := mapster.Map[int](uintVal)
	if err != nil {
		t.Fatalf("Map uint to int failed: %v", err)
	}
	if intFromUint != 100 {
		t.Errorf("Expected intFromUint=100, got %d", intFromUint)
	}
}

// TestIncompatibleTypes tests mapping between incompatible types
func TestIncompatibleTypes(t *testing.T) {
	// String to int should fail
	_, err := mapster.Map[int]("not a number")
	if err == nil {
		t.Error("Expected error for string to int conversion, got nil")
	}

	// Struct to int should fail
	type Person struct {
		Name string
	}
	_, err = mapster.Map[int](Person{Name: "John"})
	if err == nil {
		t.Error("Expected error for struct to int conversion, got nil")
	}
}

// TestMapToErrors tests error cases for MapTo
func TestMapToErrors(t *testing.T) {
	type Person struct {
		Name string
		Age  int
	}

	// Test nil source
	var dst Person
	err := mapster.MapTo(nil, &dst)
	if err == nil {
		t.Error("Expected error for nil source, got nil")
	}

	// Test nil destination
	src := Person{Name: "John", Age: 30}
	var nilDst *Person = nil
	err = mapster.MapTo(src, nilDst)
	if err == nil {
		t.Error("Expected error for nil destination, got nil")
	}
}

// TestEmptyStructs tests mapping with empty structs
func TestEmptyStructs(t *testing.T) {
	type EmptySource struct{}
	type EmptyTarget struct{}

	src := EmptySource{}
	dst, err := mapster.Map[EmptyTarget](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// Just checking that it doesn't panic
	_ = dst
}
