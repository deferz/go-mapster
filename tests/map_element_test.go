package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestMapElementTypeRegistration tests that registering key and value types
// automatically enables mapping for maps with those types
func TestMapElementTypeRegistration(t *testing.T) {
	// Define custom types for this test
	type CustomKey struct {
		ID int
	}

	type CustomValue struct {
		Name string
	}

	// Register only the key and value types
	mapster.NewMapperConfig[CustomKey, CustomKey]().Register()
	mapster.NewMapperConfig[CustomValue, CustomValue]().Register()

	// Test data
	src := map[CustomKey]CustomValue{
		{ID: 1}: {Name: "Value 1"},
		{ID: 2}: {Name: "Value 2"},
		{ID: 3}: {Name: "Value 3"},
	}

	// Test map mapping
	dst, err := mapster.Map[map[CustomKey]CustomValue](src)
	if err != nil {
		t.Fatalf("Map mapping failed: %v", err)
	}

	// Verify mapping results
	if len(dst) != len(src) {
		t.Fatalf("Expected map length %d, got %d", len(src), len(dst))
	}

	for key, value := range src {
		if dstValue, ok := dst[key]; !ok {
			t.Errorf("Key {%d} not found in destination map", key.ID)
		} else if dstValue.Name != value.Name {
			t.Errorf("Expected Name=%s for key {%d}, got %s",
				value.Name, key.ID, dstValue.Name)
		}
	}
}

// TestMapKeyConversionWithElementRegistration tests mapping between maps with different key types
func TestMapKeyConversionWithElementRegistration(t *testing.T) {
	// Register only the key conversion
	mapster.NewMapperConfig[int, int64]().Register()

	// Test data
	src := map[int]string{
		1: "Value 1",
		2: "Value 2",
		3: "Value 3",
	}

	// Test map key conversion
	dst, err := mapster.Map[map[int64]string](src)
	if err != nil {
		t.Fatalf("Map key conversion failed: %v", err)
	}

	// Verify mapping results
	if len(dst) != len(src) {
		t.Fatalf("Expected map length %d, got %d", len(src), len(dst))
	}

	for key, value := range src {
		int64Key := int64(key)
		if dstValue, ok := dst[int64Key]; !ok {
			t.Errorf("Key %d not found in destination map", key)
		} else if dstValue != value {
			t.Errorf("Expected value=%s for key %d, got %s",
				value, key, dstValue)
		}
	}
}

// TestMapValueConversionWithElementRegistration tests mapping between maps with different value types
func TestMapValueConversionWithElementRegistration(t *testing.T) {
	// Define custom types for this test
	type SourceValue struct {
		ID   int
		Name string
	}

	type TargetValue struct {
		ID   int
		Name string
	}

	// Register only the value type conversion
	mapster.NewMapperConfig[SourceValue, TargetValue]().Register()

	// Test data
	src := map[string]SourceValue{
		"key1": {ID: 1, Name: "Value 1"},
		"key2": {ID: 2, Name: "Value 2"},
		"key3": {ID: 3, Name: "Value 3"},
	}

	// Test map value conversion
	dst, err := mapster.Map[map[string]TargetValue](src)
	if err != nil {
		t.Fatalf("Map value conversion failed: %v", err)
	}

	// Verify mapping results
	if len(dst) != len(src) {
		t.Fatalf("Expected map length %d, got %d", len(src), len(dst))
	}

	for key, value := range src {
		if dstValue, ok := dst[key]; !ok {
			t.Errorf("Key %s not found in destination map", key)
		} else {
			if dstValue.ID != value.ID {
				t.Errorf("Expected ID=%d for key %s, got %d",
					value.ID, key, dstValue.ID)
			}
			if dstValue.Name != value.Name {
				t.Errorf("Expected Name=%s for key %s, got %s",
					value.Name, key, dstValue.Name)
			}
		}
	}
}
