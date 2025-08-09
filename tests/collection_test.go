package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestSliceMapping tests mapping slices of structs
func TestSliceMapping(t *testing.T) {
	src := []Item{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
		{ID: 3, Name: "Item 3"},
	}

	dst, err := mapster.Map[[]Item](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Fatalf("Expected slice length %d, got %d", len(src), len(dst))
	}

	for i, item := range src {
		if dst[i].ID != item.ID {
			t.Errorf("Expected ID=%d at index %d, got %d", item.ID, i, dst[i].ID)
		}
		if dst[i].Name != item.Name {
			t.Errorf("Expected Name=%s at index %d, got %s", item.Name, i, dst[i].Name)
		}
	}
}

// TestMapToSlice tests mapping slices to existing slices
func TestMapToSlice(t *testing.T) {
	src := []Item{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
		{ID: 3, Name: "Item 3"},
	}

	var dst []Item
	err := mapster.MapTo(src, &dst)
	if err != nil {
		t.Fatalf("MapTo failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Fatalf("Expected slice length %d, got %d", len(src), len(dst))
	}

	for i, item := range src {
		if dst[i].ID != item.ID {
			t.Errorf("Expected ID=%d at index %d, got %d", item.ID, i, dst[i].ID)
		}
		if dst[i].Name != item.Name {
			t.Errorf("Expected Name=%s at index %d, got %s", item.Name, i, dst[i].Name)
		}
	}
}

// TestMapMapping tests mapping maps of structs
func TestMapMapping(t *testing.T) {
	src := map[string]User{
		"user1": {ID: 1, Name: "User 1"},
		"user2": {ID: 2, Name: "User 2"},
		"user3": {ID: 3, Name: "User 3"},
	}

	dst, err := mapster.Map[map[string]User](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Fatalf("Expected map length %d, got %d", len(src), len(dst))
	}

	for key, user := range src {
		if _, ok := dst[key]; !ok {
			t.Errorf("Key %s not found in destination map", key)
			continue
		}
		if dst[key].ID != user.ID {
			t.Errorf("Expected ID=%d for key %s, got %d", user.ID, key, dst[key].ID)
		}
		if dst[key].Name != user.Name {
			t.Errorf("Expected Name=%s for key %s, got %s", user.Name, key, dst[key].Name)
		}
	}
}

// TestMapToMap tests mapping maps to existing maps
func TestMapToMap(t *testing.T) {
	src := map[string]User{
		"user1": {ID: 1, Name: "User 1"},
		"user2": {ID: 2, Name: "User 2"},
		"user3": {ID: 3, Name: "User 3"},
	}

	dst := make(map[string]User)
	err := mapster.MapTo(src, &dst)
	if err != nil {
		t.Fatalf("MapTo failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Fatalf("Expected map length %d, got %d", len(src), len(dst))
	}

	for key, user := range src {
		if _, ok := dst[key]; !ok {
			t.Errorf("Key %s not found in destination map", key)
			continue
		}
		if dst[key].ID != user.ID {
			t.Errorf("Expected ID=%d for key %s, got %d", user.ID, key, dst[key].ID)
		}
		if dst[key].Name != user.Name {
			t.Errorf("Expected Name=%s for key %s, got %s", user.Name, key, dst[key].Name)
		}
	}
}

// TestArrayToArray tests mapping arrays
func TestArrayToArray(t *testing.T) {
	// Note: We no longer need to register array types directly
	// Just registering the element type (int) is sufficient

	src := [2]int{1, 2}
	dst, err := mapster.Map[[4]int](src)
	if err != nil {
		t.Fatalf("Map array failed: %v", err)
	}

	// Check that the values were copied correctly
	for i := 0; i < len(src); i++ {
		if dst[i] != src[i] {
			t.Errorf("Expected dst[%d]=%d, got %d", i, src[i], dst[i])
		}
	}

	// Check that the remaining elements are zero
	for i := len(src); i < len(dst); i++ {
		if dst[i] != 0 {
			t.Errorf("Expected dst[%d]=0, got %d", i, dst[i])
		}
	}
}

// TestSliceOfPointers tests mapping slices of pointers
func TestSliceOfPointers(t *testing.T) {
	src := []*Item{
		{ID: 1, Name: "Item 1"},
		{ID: 2, Name: "Item 2"},
		{ID: 3, Name: "Item 3"},
	}

	dst, err := mapster.Map[[]*Item](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Fatalf("Expected slice length %d, got %d", len(src), len(dst))
	}

	for i, item := range src {
		if dst[i].ID != item.ID {
			t.Errorf("Expected ID=%d at index %d, got %d", item.ID, i, dst[i].ID)
		}
		if dst[i].Name != item.Name {
			t.Errorf("Expected Name=%s at index %d, got %s", item.Name, i, dst[i].Name)
		}
	}
}
