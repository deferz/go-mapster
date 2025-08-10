package tests

import (
	"testing"
	"time"

	mapster "github.com/deferz/go-mapster"
)

// 基础结构体定义
type PersonInfo struct {
	Name string
	Age  int
}

type ContactInfo struct {
	Email   string
	Phone   string
	Address string
}

type WorkInfo struct {
	Company  string
	Position string
	Salary   float64
}

// 测试源结构体（使用组合/匿名嵌入）
type SourceEmployee struct {
	PersonInfo  // 匿名嵌入
	ContactInfo // 匿名嵌入
	WorkInfo    // 匿名嵌入
	ID          int
	HireDate    time.Time
}

// 测试目标结构体（也使用组合/匿名嵌入）
type TargetEmployee struct {
	PersonInfo  // 匿名嵌入
	ContactInfo // 匿名嵌入
	WorkInfo    // 匿名嵌入
	ID          int
	HireDate    time.Time
}

// 扁平化的目标结构体
type FlattenEmployee struct {
	Name     string
	Age      int
	Email    string
	Phone    string
	Address  string
	Company  string
	Position string
	Salary   float64
	ID       int
	HireDate time.Time
}

// TestCompositeStructMapping 测试组合结构体映射
func TestCompositeStructMapping(t *testing.T) {
	hireDate := time.Date(2020, 1, 1, 0, 0, 0, 0, time.UTC)

	// 创建源对象
	src := SourceEmployee{
		PersonInfo: PersonInfo{
			Name: "张三",
			Age:  30,
		},
		ContactInfo: ContactInfo{
			Email:   "zhangsan@example.com",
			Phone:   "13800138000",
			Address: "北京市海淀区",
		},
		WorkInfo: WorkInfo{
			Company:  "ABC科技有限公司",
			Position: "高级工程师",
			Salary:   20000.0,
		},
		ID:       1001,
		HireDate: hireDate,
	}

	// 测试映射到相同结构的目标对象
	t.Run("组合结构体到组合结构体映射", func(t *testing.T) {
		dst, err := mapster.Map[TargetEmployee](src)
		if err != nil {
			t.Fatalf("组合结构体映射失败: %v", err)
		}

		// 验证 PersonInfo 字段
		if dst.Name != src.Name {
			t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
		}
		if dst.Age != src.Age {
			t.Errorf("期望 Age=%d, 得到 %d", src.Age, dst.Age)
		}

		// 验证 ContactInfo 字段
		if dst.Email != src.Email {
			t.Errorf("期望 Email=%s, 得到 %s", src.Email, dst.Email)
		}
		if dst.Phone != src.Phone {
			t.Errorf("期望 Phone=%s, 得到 %s", src.Phone, dst.Phone)
		}
		if dst.Address != src.Address {
			t.Errorf("期望 Address=%s, 得到 %s", src.Address, dst.Address)
		}

		// 验证 WorkInfo 字段
		if dst.Company != src.Company {
			t.Errorf("期望 Company=%s, 得到 %s", src.Company, dst.Company)
		}
		if dst.Position != src.Position {
			t.Errorf("期望 Position=%s, 得到 %s", src.Position, dst.Position)
		}
		if dst.Salary != src.Salary {
			t.Errorf("期望 Salary=%f, 得到 %f", src.Salary, dst.Salary)
		}

		// 验证普通字段
		if dst.ID != src.ID {
			t.Errorf("期望 ID=%d, 得到 %d", src.ID, dst.ID)
		}
		if !dst.HireDate.Equal(src.HireDate) {
			t.Errorf("期望 HireDate=%v, 得到 %v", src.HireDate, dst.HireDate)
		}
	})

	// 测试映射到扁平化结构体
	t.Run("组合结构体到扁平化结构体映射", func(t *testing.T) {
		dst, err := mapster.Map[FlattenEmployee](src)
		if err != nil {
			t.Fatalf("组合结构体到扁平化结构体映射失败: %v", err)
		}

		// 验证原来在 PersonInfo 中的字段
		if dst.Name != src.Name {
			t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
		}
		if dst.Age != src.Age {
			t.Errorf("期望 Age=%d, 得到 %d", src.Age, dst.Age)
		}

		// 验证原来在 ContactInfo 中的字段
		if dst.Email != src.Email {
			t.Errorf("期望 Email=%s, 得到 %s", src.Email, dst.Email)
		}
		if dst.Phone != src.Phone {
			t.Errorf("期望 Phone=%s, 得到 %s", src.Phone, dst.Phone)
		}
		if dst.Address != src.Address {
			t.Errorf("期望 Address=%s, 得到 %s", src.Address, dst.Address)
		}

		// 验证原来在 WorkInfo 中的字段
		if dst.Company != src.Company {
			t.Errorf("期望 Company=%s, 得到 %s", src.Company, dst.Company)
		}
		if dst.Position != src.Position {
			t.Errorf("期望 Position=%s, 得到 %s", src.Position, dst.Position)
		}
		if dst.Salary != src.Salary {
			t.Errorf("期望 Salary=%f, 得到 %f", src.Salary, dst.Salary)
		}

		// 验证普通字段
		if dst.ID != src.ID {
			t.Errorf("期望 ID=%d, 得到 %d", src.ID, dst.ID)
		}
		if !dst.HireDate.Equal(src.HireDate) {
			t.Errorf("期望 HireDate=%v, 得到 %v", src.HireDate, dst.HireDate)
		}
	})
}

// 不同名称的组合结构体
type SourceWithDifferentNames struct {
	PersonInfo              // 匿名嵌入
	ContactData ContactInfo // 有名称的嵌入
	WorkInfo                // 匿名嵌入
	ID          int
}

type TargetWithDifferentNames struct {
	PersonInfo                 // 匿名嵌入
	ContactDetails ContactInfo // 有名称的嵌入，名称不同
	WorkInfo                   // 匿名嵌入
	ID             int
}

// TestCompositeWithNamedFields 测试带有命名字段的组合结构体映射
func TestCompositeWithNamedFields(t *testing.T) {
	// 创建源对象
	src := SourceWithDifferentNames{
		PersonInfo: PersonInfo{
			Name: "李四",
			Age:  35,
		},
		ContactData: ContactInfo{
			Email:   "lisi@example.com",
			Phone:   "13900139000",
			Address: "上海市浦东新区",
		},
		WorkInfo: WorkInfo{
			Company:  "XYZ有限公司",
			Position: "技术经理",
			Salary:   25000.0,
		},
		ID: 1002,
	}

	// 测试映射
	dst, err := mapster.Map[TargetWithDifferentNames](src)
	if err != nil {
		t.Fatalf("带命名字段的组合结构体映射失败: %v", err)
	}

	// 验证匿名嵌入的字段
	if dst.Name != src.Name {
		t.Errorf("期望 Name=%s, 得到 %s", src.Name, dst.Name)
	}
	if dst.Age != src.Age {
		t.Errorf("期望 Age=%d, 得到 %d", src.Age, dst.Age)
	}

	// 验证 WorkInfo 字段
	if dst.Company != src.Company {
		t.Errorf("期望 Company=%s, 得到 %s", src.Company, dst.Company)
	}
	if dst.Position != src.Position {
		t.Errorf("期望 Position=%s, 得到 %s", src.Position, dst.Position)
	}
	if dst.Salary != src.Salary {
		t.Errorf("期望 Salary=%f, 得到 %f", src.Salary, dst.Salary)
	}

	// 验证普通字段
	if dst.ID != src.ID {
		t.Errorf("期望 ID=%d, 得到 %d", src.ID, dst.ID)
	}

	// 验证命名字段 - 这里应该是空值，因为字段名不匹配
	if dst.ContactDetails.Email != "" || dst.ContactDetails.Phone != "" || dst.ContactDetails.Address != "" {
		t.Errorf("期望命名字段为空值，但得到了非空值: %+v", dst.ContactDetails)
	}
}
