package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestMapValueOnlyRegistration tests that registering only value types
// enables mapping for maps with those value types, without requiring key type registration
func TestMapValueOnlyRegistration(t *testing.T) {
	// Define custom types for this test
	type SourceValue struct {
		ID   int
		Name string
	}

	type TargetValue struct {
		ID   int
		Name string
	}

	// Register ONLY the value type conversion
	mapster.NewMapperConfig[SourceValue, TargetValue]().Register()

	// Test data with different key types (string and int)
	srcStringKey := map[string]SourceValue{
		"key1": {ID: 1, Name: "Value 1"},
		"key2": {ID: 2, Name: "Value 2"},
	}

	srcIntKey := map[int]SourceValue{
		1: {ID: 1, Name: "Value 1"},
		2: {ID: 2, Name: "Value 2"},
	}

	// Test map with same key type
	t.Run("Same key type", func(t *testing.T) {
		dst, err := mapster.Map[map[string]TargetValue](srcStringKey)
		if err != nil {
			t.Fatalf("Map value conversion failed: %v", err)
		}

		// Verify mapping results
		if len(dst) != len(srcStringKey) {
			t.Fatalf("Expected map length %d, got %d", len(srcStringKey), len(dst))
		}

		for key, value := range srcStringKey {
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
	})

	// Test map with different but convertible key type (int to int64)
	t.Run("Convertible key type", func(t *testing.T) {
		dst, err := mapster.Map[map[int64]TargetValue](srcIntKey)
		if err != nil {
			t.Fatalf("Map with convertible key type failed: %v", err)
		}

		// Verify mapping results
		if len(dst) != len(srcIntKey) {
			t.Fatalf("Expected map length %d, got %d", len(srcIntKey), len(dst))
		}

		for key, value := range srcIntKey {
			int64Key := int64(key)
			if dstValue, ok := dst[int64Key]; !ok {
				t.Errorf("Key %d not found in destination map", key)
			} else {
				if dstValue.ID != value.ID {
					t.Errorf("Expected ID=%d for key %d, got %d",
						value.ID, key, dstValue.ID)
				}
				if dstValue.Name != value.Name {
					t.Errorf("Expected Name=%s for key %d, got %s",
						value.Name, key, dstValue.Name)
				}
			}
		}
	})

	// Test map with completely different key types (should fail if we try to map string->int)
	t.Run("Incompatible key type", func(t *testing.T) {
		_, err := mapster.Map[map[int]TargetValue](srcStringKey)
		if err == nil {
			t.Fatalf("Expected error when mapping incompatible key types, but got no error")
		}
	})
}
