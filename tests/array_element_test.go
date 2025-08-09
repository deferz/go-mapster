package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestArrayElementTypeRegistration tests that registering element types
// automatically enables mapping for arrays and slices of those types
func TestArrayElementTypeRegistration(t *testing.T) {
	// Define custom types for this test
	type CustomSource struct {
		Value int
	}

	type CustomTarget struct {
		Value int
	}

	// Register only the element types
	mapster.NewMapperConfig[CustomSource, CustomTarget]().Register()

	// Test data
	srcArray := [3]CustomSource{
		{Value: 1},
		{Value: 2},
		{Value: 3},
	}

	// Test array mapping
	dstArray, err := mapster.Map[[3]CustomTarget](srcArray)
	if err != nil {
		t.Fatalf("Array mapping failed: %v", err)
	}

	// Verify array mapping results
	for i, item := range srcArray {
		if dstArray[i].Value != item.Value {
			t.Errorf("Expected Value=%d at index %d, got %d",
				item.Value, i, dstArray[i].Value)
		}
	}

	// Test slice mapping
	srcSlice := []CustomSource{
		{Value: 1},
		{Value: 2},
		{Value: 3},
	}

	dstSlice, err := mapster.Map[[]CustomTarget](srcSlice)
	if err != nil {
		t.Fatalf("Slice mapping failed: %v", err)
	}

	// Verify slice mapping results
	if len(dstSlice) != len(srcSlice) {
		t.Fatalf("Expected slice length %d, got %d", len(srcSlice), len(dstSlice))
	}

	for i, item := range srcSlice {
		if dstSlice[i].Value != item.Value {
			t.Errorf("Expected Value=%d at index %d, got %d",
				item.Value, i, dstSlice[i].Value)
		}
	}
}

// TestDifferentArraySizes tests mapping between arrays of different sizes
func TestDifferentArraySizes(t *testing.T) {
	// Register only int element type
	mapster.NewMapperConfig[int, int]().Register()

	// Test data
	srcArray := [5]int{1, 2, 3, 4, 5}

	// Map to smaller array
	smallerArray, err := mapster.Map[[3]int](srcArray)
	if err != nil {
		t.Fatalf("Mapping to smaller array failed: %v", err)
	}

	// Verify first 3 elements were copied
	for i := 0; i < 3; i++ {
		if smallerArray[i] != srcArray[i] {
			t.Errorf("Expected smallerArray[%d]=%d, got %d",
				i, srcArray[i], smallerArray[i])
		}
	}

	// Map to larger array
	largerArray, err := mapster.Map[[7]int](srcArray)
	if err != nil {
		t.Fatalf("Mapping to larger array failed: %v", err)
	}

	// Verify first 5 elements were copied
	for i := 0; i < 5; i++ {
		if largerArray[i] != srcArray[i] {
			t.Errorf("Expected largerArray[%d]=%d, got %d",
				i, srcArray[i], largerArray[i])
		}
	}

	// Verify remaining elements are zero
	for i := 5; i < 7; i++ {
		if largerArray[i] != 0 {
			t.Errorf("Expected largerArray[%d]=0, got %d", i, largerArray[i])
		}
	}
}
