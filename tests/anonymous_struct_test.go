package tests

import (
	"testing"

	mapster "github.com/deferz/go-mapster"
)

// TestAnonymousStructMapping 测试匿名结构体映射
func TestAnonymousStructMapping(t *testing.T) {
	// 创建匿名源结构体
	src := struct {
		Name  string
		Age   int
		Email string
	}{
		Name:  "张三",
		Age:   30,
		Email: "zhangsan@example.com",
	}

	// 创建匿名目标结构体
	dst, err := mapster.Map[struct {
		Name  string
		Age   int
		Email string
	}](src)

	if err != nil {
		t.Fatalf("匿名结构体映射失败: %v", err)
	}

	// 验证字段值
	if dst.Name != src.Name {
		t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("期望 Age=%d, 得到 %d", src.Age, dst.Age)
	}
	if dst.Email != src.Email {
		t.Errorf("期望 Email=%s, 得到 %s", src.Email, dst.Email)
	}
}

// TestAnonymousNestedStructMapping 测试带嵌套的匿名结构体映射
func TestAnonymousNestedStructMapping(t *testing.T) {
	// 创建带嵌套的匿名源结构体
	src := struct {
		Name    string
		Age     int
		Address struct {
			Street  string
			City    string
			Country string
		}
	}{
		Name: "张三",
		Age:  30,
		Address: struct {
			Street  string
			City    string
			Country string
		}{
			Street:  "中关村大街",
			City:    "北京",
			Country: "中国",
		},
	}

	// 创建匿名目标结构体
	dst, err := mapster.Map[struct {
		Name    string
		Age     int
		Address struct {
			Street  string
			City    string
			Country string
		}
	}](src)

	if err != nil {
		t.Fatalf("带嵌套的匿名结构体映射失败: %v", err)
	}

	// 验证字段值
	if dst.Name != src.Name {
		t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("期望 Age=%d, 得到 %d", src.Age, dst.Age)
	}
	if dst.Address.Street != src.Address.Street {
		t.Errorf("期望 Street=%s, 得到 %s", src.Address.Street, dst.Address.Street)
	}
	if dst.Address.City != src.Address.City {
		t.Errorf("期望 City=%s, 得到 %s", src.Address.City, dst.Address.City)
	}
	if dst.Address.Country != src.Address.Country {
		t.Errorf("期望 Country=%s, 得到 %s", src.Address.Country, dst.Address.Country)
	}
}

// TestAnonymousToFlattenMapping 测试匿名嵌套结构体到扁平化结构体的映射
func TestAnonymousToFlattenMapping(t *testing.T) {
	// 创建带嵌套的匿名源结构体
	src := struct {
		Name    string
		Age     int
		Address struct {
			Street  string
			City    string
			Country string
		}
	}{
		Name: "张三",
		Age:  30,
		Address: struct {
			Street  string
			City    string
			Country string
		}{
			Street:  "中关村大街",
			City:    "北京",
			Country: "中国",
		},
	}

	// 创建扁平化的匿名目标结构体
	dst, err := mapster.Map[struct {
		Name    string
		Age     int
		Street  string
		City    string
		Country string
	}](src)

	if err != nil {
		t.Fatalf("匿名嵌套结构体到扁平化结构体的映射失败: %v", err)
	}

	// 验证字段值
	if dst.Name != src.Name {
		t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("期望 Age=%d, 得到 %d", src.Age, dst.Age)
	}
	if dst.Street != src.Address.Street {
		t.Errorf("期望 Street=%s, 得到 %s", src.Address.Street, dst.Street)
	}
	if dst.City != src.Address.City {
		t.Errorf("期望 City=%s, 得到 %s", src.Address.City, dst.City)
	}
	if dst.Country != src.Address.Country {
		t.Errorf("期望 Country=%s, 得到 %s", src.Address.Country, dst.Country)
	}
}

// TestAnonymousEmbeddedStructMapping 测试匿名嵌入结构体的映射
func TestAnonymousEmbeddedStructMapping(t *testing.T) {
	// 创建一个带有匿名嵌入结构体的源结构体
	type BaseInfo struct {
		ID   int
		Time string
	}

	src := struct {
		BaseInfo // 匿名嵌入
		Name     string
		Age      int
	}{
		BaseInfo: BaseInfo{
			ID:   1001,
			Time: "2023-01-01",
		},
		Name: "张三",
		Age:  30,
	}

	// 创建一个带有匿名嵌入结构体的目标结构体
	dst, err := mapster.Map[struct {
		BaseInfo // 匿名嵌入
		Name     string
		Age      int
	}](src)

	if err != nil {
		t.Fatalf("匿名嵌入结构体映射失败: %v", err)
	}

	// 验证字段值
	if dst.ID != src.ID {
		t.Errorf("期望 ID=%d, 得到 %d", src.ID, dst.ID)
	}
	if dst.Time != src.Time {
		t.Errorf("期望 Time=%s, 得到 %s", src.Time, dst.Time)
	}
	if dst.Name != src.Name {
		t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("期望 Age=%d, 得到 %d", src.Age, dst.Age)
	}
}
