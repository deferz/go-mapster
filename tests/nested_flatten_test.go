package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// 定义3层嵌套的源结构体
type Level3Source struct {
	Value3 string
	Num3   int
}

type Level2Source struct {
	Level3 Level3Source
	Value2 string
	Num2   int
}

type Level1Source struct {
	Level2 Level2Source
	Value1 string
	Num1   int
}

// 定义扁平化的目标结构体
type FlattenTarget struct {
	Value1 string
	Num1   int
	Value2 string
	Num2   int
	Value3 string
	Num3   int
}

// TestNestedFlatten 测试嵌套结构体是否能自动扁平化映射
func TestNestedFlatten(t *testing.T) {
	// 创建嵌套源结构体
	src := Level1Source{
		Value1: "level1 value",
		Num1:   1,
		Level2: Level2Source{
			Value2: "level2 value",
			Num2:   2,
			Level3: Level3Source{
				Value3: "level3 value",
				Num3:   3,
			},
		},
	}

	// 尝试映射到扁平化的目标结构体
	dst, err := mapster.Map[FlattenTarget](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// 验证第一层字段
	if dst.Value1 != src.Value1 {
		t.Errorf("Expected Value1=%s, got %s", src.Value1, dst.Value1)
	}
	if dst.Num1 != src.Num1 {
		t.Errorf("Expected Num1=%d, got %d", src.Num1, dst.Num1)
	}

	// 验证第二层字段
	if dst.Value2 != src.Level2.Value2 {
		t.Errorf("Expected Value2=%s, got %s", src.Level2.Value2, dst.Value2)
	}
	if dst.Num2 != src.Level2.Num2 {
		t.Errorf("Expected Num2=%d, got %d", src.Level2.Num2, dst.Num2)
	}

	// 验证第三层字段
	if dst.Value3 != src.Level2.Level3.Value3 {
		t.Errorf("Expected Value3=%s, got %s", src.Level2.Level3.Value3, dst.Value3)
	}
	if dst.Num3 != src.Level2.Level3.Num3 {
		t.Errorf("Expected Num3=%d, got %d", src.Level2.Level3.Num3, dst.Num3)
	}
}

// 测试使用MapTo进行扁平化映射
func TestNestedFlattenMapTo(t *testing.T) {
	// 创建嵌套源结构体
	src := Level1Source{
		Value1: "level1 value",
		Num1:   1,
		Level2: Level2Source{
			Value2: "level2 value",
			Num2:   2,
			Level3: Level3Source{
				Value3: "level3 value",
				Num3:   3,
			},
		},
	}

	// 创建目标结构体实例
	dst := FlattenTarget{}

	// 使用MapTo映射
	err := mapster.MapTo(src, &dst)
	if err != nil {
		t.Fatalf("MapTo failed: %v", err)
	}

	// 验证第一层字段
	if dst.Value1 != src.Value1 {
		t.Errorf("Expected Value1=%s, got %s", src.Value1, dst.Value1)
	}
	if dst.Num1 != src.Num1 {
		t.Errorf("Expected Num1=%d, got %d", src.Num1, dst.Num1)
	}

	// 验证第二层字段
	if dst.Value2 != src.Level2.Value2 {
		t.Errorf("Expected Value2=%s, got %s", src.Level2.Value2, dst.Value2)
	}
	if dst.Num2 != src.Level2.Num2 {
		t.Errorf("Expected Num2=%d, got %d", src.Level2.Num2, dst.Num2)
	}

	// 验证第三层字段
	if dst.Value3 != src.Level2.Level3.Value3 {
		t.Errorf("Expected Value3=%s, got %s", src.Level2.Level3.Value3, dst.Value3)
	}
	if dst.Num3 != src.Level2.Level3.Num3 {
		t.Errorf("Expected Num3=%d, got %d", src.Level2.Level3.Num3, dst.Num3)
	}
}

// 定义带有嵌套前缀的目标结构体
type PrefixTarget struct {
	Value1     string
	Num1       int
	Level2_Value2 string
	Level2_Num2   int
	Level2_Level3_Value3 string
	Level2_Level3_Num3   int
}

// 测试带有前缀的扁平化映射
func TestNestedFlattenWithPrefix(t *testing.T) {
	// 创建嵌套源结构体
	src := Level1Source{
		Value1: "level1 value",
		Num1:   1,
		Level2: Level2Source{
			Value2: "level2 value",
			Num2:   2,
			Level3: Level3Source{
				Value3: "level3 value",
				Num3:   3,
			},
		},
	}

	// 尝试映射到带前缀的扁平化目标结构体
	dst, err := mapster.Map[PrefixTarget](src)
	if err != nil {
		t.Fatalf("Map failed: %v", err)
	}

	// 验证第一层字段
	if dst.Value1 != src.Value1 {
		t.Errorf("Expected Value1=%s, got %s", src.Value1, dst.Value1)
	}
	if dst.Num1 != src.Num1 {
		t.Errorf("Expected Num1=%d, got %d", src.Num1, dst.Num1)
	}

	// 验证第二层字段（带前缀）
	if dst.Level2_Value2 != src.Level2.Value2 {
		t.Errorf("Expected Level2_Value2=%s, got %s", src.Level2.Value2, dst.Level2_Value2)
	}
	if dst.Level2_Num2 != src.Level2.Num2 {
		t.Errorf("Expected Level2_Num2=%d, got %d", src.Level2.Num2, dst.Level2_Num2)
	}

	// 验证第三层字段（带前缀）
	if dst.Level2_Level3_Value3 != src.Level2.Level3.Value3 {
		t.Errorf("Expected Level2_Level3_Value3=%s, got %s", src.Level2.Level3.Value3, dst.Level2_Level3_Value3)
	}
	if dst.Level2_Level3_Num3 != src.Level2.Level3.Num3 {
		t.Errorf("Expected Level2_Level3_Num3=%d, got %d", src.Level2.Level3.Num3, dst.Level2_Level3_Num3)
	}
}
