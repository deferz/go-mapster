package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestSliceMapping tests mapping slices of structs
func TestSliceMapping(t *testing.T) {
	type Item struct {
		ID   int
		Name string
	}

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
		if dst[i].ID != item.ID || dst[i].Name != item.Name {
			t.Errorf("Item %d mismatch: expected %+v, got %+v", i, item, dst[i])
		}
	}
}

// TestMapToSlice tests MapTo with slices
func TestMapToSlice(t *testing.T) {
	type Item struct {
		ID   int
		Name string
	}

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
		if dst[i].ID != item.ID || dst[i].Name != item.Name {
			t.Errorf("Item %d mismatch: expected %+v, got %+v", i, item, dst[i])
		}
	}
}

// TestMapMapping tests mapping maps
func TestMapMapping(t *testing.T) {
	type User struct {
		Name  string
		Email string
	}

	src := map[string]User{
		"user1": {Name: "User 1", Email: "user1@example.com"},
		"user2": {Name: "User 2", Email: "user2@example.com"},
	}

	dst, err := mapster.Map[map[string]User](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Fatalf("Expected map size %d, got %d", len(src), len(dst))
	}

	for key, srcUser := range src {
		dstUser, ok := dst[key]
		if !ok {
			t.Errorf("Key %s missing from destination map", key)
			continue
		}

		if dstUser.Name != srcUser.Name || dstUser.Email != srcUser.Email {
			t.Errorf("User %s mismatch: expected %+v, got %+v", key, srcUser, dstUser)
		}
	}
}

// TestMapToMap tests MapTo with maps
func TestMapToMap(t *testing.T) {
	type User struct {
		Name  string
		Email string
	}

	src := map[string]User{
		"user1": {Name: "User 1", Email: "user1@example.com"},
		"user2": {Name: "User 2", Email: "user2@example.com"},
	}

	var dst map[string]User
	err := mapster.MapTo(src, &dst)
	if err != nil {
		t.Fatalf("MapTo failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Fatalf("Expected map size %d, got %d", len(src), len(dst))
	}

	for key, srcUser := range src {
		dstUser, ok := dst[key]
		if !ok {
			t.Errorf("Key %s missing from destination map", key)
			continue
		}

		if dstUser.Name != srcUser.Name || dstUser.Email != srcUser.Email {
			t.Errorf("User %s mismatch: expected %+v, got %+v", key, srcUser, dstUser)
		}
	}
}

// TestArrayToArray tests array to array mapping
func TestArrayToArray(t *testing.T) {
	// Source array larger than destination
	srcLarger := [5]int{1, 2, 3, 4, 5}
	dstSmaller, err := mapster.Map[[3]int](srcLarger)
	if err != nil {
		t.Fatalf("Map array failed: %v", err)
	}

	// Should truncate to first 3 elements
	expected := [3]int{1, 2, 3}
	for i, v := range expected {
		if dstSmaller[i] != v {
			t.Errorf("Expected dstSmaller[%d]=%d, got %d", i, v, dstSmaller[i])
		}
	}

	// Source array smaller than destination
	srcSmaller := [2]int{1, 2}
	dstLarger, err := mapster.Map[[4]int](srcSmaller)
	if err != nil {
		t.Fatalf("Map array failed: %v", err)
	}

	// Should copy first 2 elements, rest should be zero values
	expected2 := [4]int{1, 2, 0, 0}
	for i, v := range expected2 {
		if dstLarger[i] != v {
			t.Errorf("Expected dstLarger[%d]=%d, got %d", i, v, dstLarger[i])
		}
	}
}

// TestMapKeyConversion tests map key type conversion
func TestMapKeyConversion(t *testing.T) {
	// Map with int keys to map with int64 keys
	srcMap := map[int]string{
		1: "one",
		2: "two",
		3: "three",
	}

	dstMap, err := mapster.Map[map[int64]string](srcMap)
	if err != nil {
		t.Fatalf("Map key conversion failed: %v", err)
	}

	if len(dstMap) != len(srcMap) {
		t.Fatalf("Expected map size %d, got %d", len(srcMap), len(dstMap))
	}

	// Check that keys were properly converted
	for k, v := range srcMap {
		if dstMap[int64(k)] != v {
			t.Errorf("Expected dstMap[%d]=%s, got %s", k, v, dstMap[int64(k)])
		}
	}
}

// TestSliceOfPointers tests mapping slices of pointers
func TestSliceOfPointers(t *testing.T) {
	type Item struct {
		ID   int
		Name string
	}

	// Create source slice of pointers
	item1 := Item{ID: 1, Name: "Item 1"}
	item2 := Item{ID: 2, Name: "Item 2"}
	src := []*Item{&item1, &item2}

	// Map to destination slice of pointers
	dst, err := mapster.Map[[]*Item](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	if len(dst) != len(src) {
		t.Fatalf("Expected slice length %d, got %d", len(src), len(dst))
	}

	// Verify content matches (but don't check pointers)
	// Note: Current implementation may not create new pointers
	for i, srcItem := range src {
		dstItem := dst[i]

		// Content should match
		if dstItem.ID != srcItem.ID || dstItem.Name != srcItem.Name {
			t.Errorf("Item %d mismatch: expected %+v, got %+v", i, *srcItem, *dstItem)
		}
	}
}
