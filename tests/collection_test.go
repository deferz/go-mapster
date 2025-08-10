package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestCollectionMapping tests mapping of collection types (slices, arrays, maps)
func TestCollectionMapping(t *testing.T) {
	// 测试切片映射
	t.Run("Slice mapping", func(t *testing.T) {
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
	})

	// 测试切片到现有切片的映射
	t.Run("MapTo slice", func(t *testing.T) {
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
	})

	// 测试映射结构体的映射
	t.Run("Map mapping", func(t *testing.T) {
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
	})

	// 测试映射到现有映射
	t.Run("MapTo map", func(t *testing.T) {
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
	})

	// 测试数组映射
	t.Run("Array to array", func(t *testing.T) {
		// 注意：我们不再需要直接注册数组类型
		// 只需要注册元素类型 (int) 就足够了

		src := [2]int{1, 2}
		dst, err := mapster.Map[[4]int](src)
		if err != nil {
			t.Fatalf("Map array failed: %v", err)
		}

		// 检查值是否正确复制
		for i := 0; i < len(src); i++ {
			if dst[i] != src[i] {
				t.Errorf("Expected dst[%d]=%d, got %d", i, src[i], dst[i])
			}
		}

		// 检查剩余元素是否为零
		for i := len(src); i < len(dst); i++ {
			if dst[i] != 0 {
				t.Errorf("Expected dst[%d]=0, got %d", i, dst[i])
			}
		}
	})

	// 测试注册元素类型自动启用数组和切片映射
	t.Run("Element type registration", func(t *testing.T) {
		// 定义测试用的自定义类型
		type CustomSource struct {
			Value int
		}

		type CustomTarget struct {
			Value int
		}

		// 映射现在在调用 Map 和 MapTo 时自动注册

		// 测试数组映射
		t.Run("Array mapping", func(t *testing.T) {
			srcArray := [3]CustomSource{
				{Value: 1},
				{Value: 2},
				{Value: 3},
			}

			dstArray, err := mapster.Map[[3]CustomTarget](srcArray)
			if err != nil {
				t.Fatalf("Array mapping failed: %v", err)
			}

			// 验证数组映射结果
			for i, item := range srcArray {
				if dstArray[i].Value != item.Value {
					t.Errorf("Expected Value=%d at index %d, got %d",
						item.Value, i, dstArray[i].Value)
				}
			}
		})

		// 测试切片映射
		t.Run("Slice mapping", func(t *testing.T) {
			srcSlice := []CustomSource{
				{Value: 1},
				{Value: 2},
				{Value: 3},
			}

			dstSlice, err := mapster.Map[[]CustomTarget](srcSlice)
			if err != nil {
				t.Fatalf("Slice mapping failed: %v", err)
			}

			// 验证切片映射结果
			if len(dstSlice) != len(srcSlice) {
				t.Fatalf("Expected slice length %d, got %d", len(srcSlice), len(dstSlice))
			}

			for i, item := range srcSlice {
				if dstSlice[i].Value != item.Value {
					t.Errorf("Expected Value=%d at index %d, got %d",
						item.Value, i, dstSlice[i].Value)
				}
			}
		})
	})

	// 测试不同大小的数组映射
	t.Run("Different array sizes", func(t *testing.T) {
		// 映射现在在调用 Map 和 MapTo 时自动注册

		// 测试数据
		srcArray := [5]int{1, 2, 3, 4, 5}

		// 映射到更小的数组
		t.Run("To smaller array", func(t *testing.T) {
			smallerArray, err := mapster.Map[[3]int](srcArray)
			if err != nil {
				t.Fatalf("Mapping to smaller array failed: %v", err)
			}

			// 验证前 3 个元素已复制
			for i := 0; i < 3; i++ {
				if smallerArray[i] != srcArray[i] {
					t.Errorf("Expected smallerArray[%d]=%d, got %d",
						i, srcArray[i], smallerArray[i])
				}
			}
		})

		// 映射到更大的数组
		t.Run("To larger array", func(t *testing.T) {
			largerArray, err := mapster.Map[[7]int](srcArray)
			if err != nil {
				t.Fatalf("Mapping to larger array failed: %v", err)
			}

			// 验证前 5 个元素已复制
			for i := 0; i < 5; i++ {
				if largerArray[i] != srcArray[i] {
					t.Errorf("Expected largerArray[%d]=%d, got %d",
						i, srcArray[i], largerArray[i])
				}
			}

			// 验证剩余元素为零
			for i := 5; i < 7; i++ {
				if largerArray[i] != 0 {
					t.Errorf("Expected largerArray[%d]=0, got %d", i, largerArray[i])
				}
			}
		})
	})

	// 测试注册键和值类型自动启用映射映射
	t.Run("Map element type registration", func(t *testing.T) {
		// 定义测试用的自定义类型
		type CustomKey struct {
			ID int
		}

		type CustomValue struct {
			Name string
		}

		// 映射现在在调用 Map 和 MapTo 时自动注册

		// 测试数据
		src := map[CustomKey]CustomValue{
			{ID: 1}: {Name: "Value 1"},
			{ID: 2}: {Name: "Value 2"},
			{ID: 3}: {Name: "Value 3"},
		}

		// 测试映射映射
		dst, err := mapster.Map[map[CustomKey]CustomValue](src)
		if err != nil {
			t.Fatalf("Map mapping failed: %v", err)
		}

		// 验证映射结果
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
	})

	// 测试不同键类型的映射映射
	t.Run("Map key conversion", func(t *testing.T) {
		// 映射现在在调用 Map 和 MapTo 时自动注册

		// 测试数据
		src := map[int]string{
			1: "Value 1",
			2: "Value 2",
			3: "Value 3",
		}

		// 测试映射键转换
		dst, err := mapster.Map[map[int64]string](src)
		if err != nil {
			t.Fatalf("Map key conversion failed: %v", err)
		}

		// 验证映射结果
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
	})

	// 测试不同值类型的映射映射
	t.Run("Map value conversion", func(t *testing.T) {
		// 定义测试用的自定义类型
		type SourceValue struct {
			ID   int
			Name string
		}

		type TargetValue struct {
			ID   int
			Name string
		}

		// 映射现在在调用 Map 和 MapTo 时自动注册

		// 测试数据
		src := map[string]SourceValue{
			"key1": {ID: 1, Name: "Value 1"},
			"key2": {ID: 2, Name: "Value 2"},
			"key3": {ID: 3, Name: "Value 3"},
		}

		// 测试映射值转换
		dst, err := mapster.Map[map[string]TargetValue](src)
		if err != nil {
			t.Fatalf("Map value conversion failed: %v", err)
		}

		// 验证映射结果
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
	})

	// 测试只注册值类型启用映射
	t.Run("Map value only registration", func(t *testing.T) {
		// 定义测试用的自定义类型
		type SourceValue struct {
			ID   int
			Name string
		}

		type TargetValue struct {
			ID   int
			Name string
		}

		// 映射现在在调用 Map 和 MapTo 时自动注册

		// 不同键类型的测试数据（字符串和整数）
		srcStringKey := map[string]SourceValue{
			"key1": {ID: 1, Name: "Value 1"},
			"key2": {ID: 2, Name: "Value 2"},
		}

		srcIntKey := map[int]SourceValue{
			1: {ID: 1, Name: "Value 1"},
			2: {ID: 2, Name: "Value 2"},
		}

		// 测试相同键类型的映射
		t.Run("Same key type", func(t *testing.T) {
			dst, err := mapster.Map[map[string]TargetValue](srcStringKey)
			if err != nil {
				t.Fatalf("Map value conversion failed: %v", err)
			}

			// 验证映射结果
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

		// 测试不同但可转换的键类型（int 到 int64）
		t.Run("Convertible key type", func(t *testing.T) {
			dst, err := mapster.Map[map[int64]TargetValue](srcIntKey)
			if err != nil {
				t.Fatalf("Map with convertible key type failed: %v", err)
			}

			// 验证映射结果
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

		// 测试完全不同的键类型（如果我们尝试映射 string->int 应该失败）
		t.Run("Incompatible key type", func(t *testing.T) {
			_, err := mapster.Map[map[int]TargetValue](srcStringKey)
			if err == nil {
				t.Fatalf("Expected error when mapping incompatible key types, but got no error")
			}
		})
	})
}

// TestPointerCollections tests collections with pointer elements
func TestPointerCollections(t *testing.T) {
	// 测试指针切片映射
	t.Run("Slice of pointers", func(t *testing.T) {
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
	})
}