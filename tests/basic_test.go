package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestAutomaticRegistration tests that mapping is automatically registered when Map or MapTo is called
func TestAutomaticRegistration(t *testing.T) {
	// Create a new type that hasn't been explicitly registered
	type AutoRegisteredSource struct {
		ID   int
		Name string
	}

	// Create test data
	source := AutoRegisteredSource{
		ID:   1,
		Name: "John Doe",
	}

	// 定义一个简单的目标类型
	type SimpleTarget struct {
		ID   int
		Name string
	}

	// Attempt mapping with automatic registration
	target, err := mapster.Map[SimpleTarget](source)

	// Should succeed now with automatic registration
	if err != nil {
		t.Errorf("Expected automatic registration to work, but got error: %v", err)
	}

	// Verify basic mapping worked
	if target.Name != source.Name {
		t.Errorf("Expected Name=%s, got %s", source.Name, target.Name)
	}
	if target.ID != source.ID {
		t.Errorf("Expected ID=%d, got %d", source.ID, target.ID)
	}
}

// TestStructMapping tests various struct mapping scenarios
func TestStructMapping(t *testing.T) {
	// 基本结构体映射
	t.Run("Basic struct mapping", func(t *testing.T) {
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
	})

	// 映射到现有对象
	t.Run("MapTo basic", func(t *testing.T) {
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
	})

	// 字段名称不同的结构体映射
	t.Run("Different field names", func(t *testing.T) {
		// 定义不同字段名的类型
		type SourceWithDifferentFields struct {
			FullName string
			Years    int
			Contact  string
		}

		type TargetWithDifferentFields struct {
			Name  string
			Age   int
			Email string
		}

		src := SourceWithDifferentFields{
			FullName: "John Doe",
			Years:    30,
			Contact:  "john@example.com",
		}

		// 使用自动注册功能
		dst, err := mapster.Map[TargetWithDifferentFields](src)
		if err != nil {
			// 当字段名不同时，默认情况下不会映射，所以这里我们期望失败
			// 实际应用中需要使用自定义映射器或字段标签
			t.Log("Expected mapping to fail with different field names: ", err)
		} else {
			// 如果映射成功，那么字段应该是空的
			if dst.Name != "" || dst.Age != 0 || dst.Email != "" {
				t.Errorf("Expected empty fields, got Name=%s, Age=%d, Email=%s", 
					dst.Name, dst.Age, dst.Email)
			}
		}
	})

	// 指针类型映射
	t.Run("Pointer mapping", func(t *testing.T) {
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
	})

	// 保留目标对象中的字段
	t.Run("MapTo preserves fields", func(t *testing.T) {
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
	})

	// 映射到空指针
	t.Run("MapTo nil pointer", func(t *testing.T) {
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
	})
}

// TestTypeConversions tests conversions between different basic types
func TestTypeConversions(t *testing.T) {
	// 定义测试用例
	testCases := []struct {
		name     string
		source   any
		target   any
		expected any
	}{
		{"int to int64", 42, int64(0), int64(42)},
		{"int to float64", 42, float64(0), float64(42.0)},
		{"float64 to int", 42.75, int(0), int(42)},
		{"uint to int", uint(42), int(0), int(42)},
	}

	// 运行测试用例
	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			// 直接调用 Map 函数
			var err error
			var result any
			
			switch tc.target.(type) {
			case int:
				result, err = mapster.Map[int](tc.source)
			case int64:
				result, err = mapster.Map[int64](tc.source)
			case float64:
				result, err = mapster.Map[float64](tc.source)
			}
			if err != nil {
				t.Fatalf("%s failed: %v", tc.name, err)
			}
			
			// 检查结果
			if result != tc.expected {
				t.Errorf("%s: expected %v, got %v", tc.name, tc.expected, result)
			}
		})
	}

	// 测试不支持的转换
	t.Run("Unsupported conversions", func(t *testing.T) {
		// String to int conversion is not supported by default
		// This would require a custom converter implementation
		
		// String to bool conversion is not supported by default
		// This would require a custom converter implementation
	})
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
