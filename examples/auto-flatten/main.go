package main

import (
	"fmt"
	"time"

	"github.com/deferz/go-mapster"
)

// 定义各种类型的别名
type MyInt int
type MyString string
type MyFloat64 float64
type MyBool bool
type MyTime time.Time

// 嵌套结构体 B
type StructB struct {
	Name        string
	Age         int
	MyAge       MyInt
	Email       string
	MyEmail     MyString
	Score       float64
	MyScore     MyFloat64
	IsActive    bool
	MyIsActive  MyBool
	CreatedAt   time.Time
	MyCreatedAt MyTime
}

// 嵌套结构体 D
type StructD struct {
	Address     string
	MyAddress   MyString
	Phone       string
	MyPhone     MyString
	Salary      int
	MySalary    MyInt
	Rating      float64
	MyRating    MyFloat64
	IsVIP       bool
	MyIsVIP     MyBool
	UpdatedAt   time.Time
	MyUpdatedAt MyTime
}

// 源结构体 A - 包含 B 和 D 结构体
type StructA struct {
	ID       int
	MyID     MyInt
	Title    string
	MyTitle  MyString
	Status   bool
	MyStatus MyBool
	StructB
	StructD
}

// 目标结构体 C - 包含 B 和 D 的所有字段
type StructC struct {
	// 来自 A 的字段
	ID       int
	MyID     int // 从 MyInt 转换到 int
	Title    string
	MyTitle  string // 从 MyString 转换到 string
	Status   bool
	MyStatus bool // 从 MyBool 转换到 bool

	// 来自 B 的字段
	Name        string
	MyName      string // 从 MyString 转换到 string
	Age         int
	MyAge       int // 从 MyInt 转换到 int
	Email       string
	MyEmail     string // 从 MyString 转换到 string
	Score       float64
	MyScore     float64 // 从 MyFloat64 转换到 float64
	IsActive    bool
	MyIsActive  bool // 从 MyBool 转换到 bool
	CreatedAt   time.Time
	MyCreatedAt time.Time // 从 MyTime 转换到 time.Time

	// 来自 D 的字段
	Address     string
	MyAddress   string // 从 MyString 转换到 string
	Phone       string
	MyPhone     string // 从 MyString 转换到 string
	Salary      int
	MySalary    int // 从 MyInt 转换到 int
	Rating      float64
	MyRating    float64 // 从 MyFloat64 转换到 float64
	IsVIP       bool
	MyIsVIP     bool // 从 MyBool 转换到 bool
	UpdatedAt   time.Time
	MyUpdatedAt time.Time // 从 MyTime 转换到 time.Time
}

func main() {
	fmt.Println("=== 自动扁平化功能示例 ===\n")

	// 创建测试数据
	now := time.Now()
	myTime := MyTime(now)

	structB := StructB{
		Name:        "John Doe",
		Age:         30,
		MyAge:       MyInt(30),
		Email:       "john@example.com",
		MyEmail:     MyString("john@example.com"),
		Score:       95.5,
		MyScore:     MyFloat64(95.5),
		IsActive:    true,
		MyIsActive:  MyBool(true),
		CreatedAt:   now,
		MyCreatedAt: myTime,
	}

	structD := StructD{
		Address:     "123 Main St",
		MyAddress:   MyString("123 Main St"),
		Phone:       "+1-555-0123",
		MyPhone:     MyString("+1-555-0123"),
		Salary:      75000,
		MySalary:    MyInt(75000),
		Rating:      4.8,
		MyRating:    MyFloat64(4.8),
		IsVIP:       false,
		MyIsVIP:     MyBool(false),
		UpdatedAt:   now.Add(time.Hour),
		MyUpdatedAt: MyTime(now.Add(time.Hour)),
	}

	src := StructA{
		ID:       1001,
		MyID:     MyInt(1001),
		Title:    "Senior Developer",
		MyTitle:  MyString("Senior Developer"),
		Status:   true,
		MyStatus: MyBool(true),
		StructB:  structB,
		StructD:  structD,
	}

	fmt.Printf("源结构体 A 数据:\n")
	fmt.Printf("  ID: %v\n", src.ID)
	fmt.Printf("  Title: %v\n", src.Title)
	fmt.Printf("  Status: %v\n", src.Status)
	fmt.Printf("  StructB.Name: %v\n", src.StructB.Name)
	fmt.Printf("  StructB.Age: %v\n", src.StructB.Age)
	fmt.Printf("  StructD.Address: %v\n", src.StructD.Address)
	fmt.Printf("  StructD.Salary: %v\n", src.StructD.Salary)
	fmt.Println()

	// 示例1: 禁用自动扁平化（默认行为）
	fmt.Println("1. 禁用自动扁平化（默认行为）:")
	mapster.DisableAutoFlatten()

	var target1 StructC
	mapster.MapTo(src, &target1)

	fmt.Printf("   来自 A 的字段:\n")
	fmt.Printf("     ID: %v\n", target1.ID)
	fmt.Printf("     MyID: %v\n", target1.MyID)
	fmt.Printf("     Title: %v\n", target1.Title)
	fmt.Printf("     MyTitle: %v\n", target1.MyTitle)
	fmt.Printf("     Status: %v\n", target1.Status)
	fmt.Printf("     MyStatus: %v\n", target1.MyStatus)

	fmt.Printf("   来自 B 的字段（未映射）:\n")
	fmt.Printf("     Name: %v (应为 'John Doe')\n", target1.Name)
	fmt.Printf("     Age: %v (应为 30)\n", target1.Age)
	fmt.Printf("     Email: %v (应为 'john@example.com')\n", target1.Email)

	fmt.Printf("   来自 D 的字段（未映射）:\n")
	fmt.Printf("     Address: %v (应为 '123 Main St')\n", target1.Address)
	fmt.Printf("     Salary: %v (应为 75000)\n", target1.Salary)
	fmt.Printf("     Phone: %v (应为 '+1-555-0123')\n", target1.Phone)
	fmt.Println()

	// 示例2: 启用自动扁平化
	fmt.Println("2. 启用自动扁平化:")
	mapster.EnableAutoFlatten()

	var target2 StructC
	mapster.MapTo(src, &target2)

	fmt.Printf("   来自 A 的字段:\n")
	fmt.Printf("     ID: %v\n", target2.ID)
	fmt.Printf("     MyID: %v\n", target2.MyID)
	fmt.Printf("     Title: %v\n", target2.Title)
	fmt.Printf("     MyTitle: %v\n", target2.MyTitle)
	fmt.Printf("     Status: %v\n", target2.Status)
	fmt.Printf("     MyStatus: %v\n", target2.MyStatus)

	fmt.Printf("   来自 B 的字段（自动扁平化）:\n")
	fmt.Printf("     Name: %v\n", target2.Name)
	fmt.Printf("     MyName: %v\n", target2.MyName)
	fmt.Printf("     Age: %v\n", target2.Age)
	fmt.Printf("     MyAge: %v\n", target2.MyAge)
	fmt.Printf("     Email: %v\n", target2.Email)
	fmt.Printf("     MyEmail: %v\n", target2.MyEmail)
	fmt.Printf("     Score: %v\n", target2.Score)
	fmt.Printf("     MyScore: %v\n", target2.MyScore)
	fmt.Printf("     IsActive: %v\n", target2.IsActive)
	fmt.Printf("     MyIsActive: %v\n", target2.MyIsActive)
	fmt.Printf("     CreatedAt: %v\n", target2.CreatedAt)
	fmt.Printf("     MyCreatedAt: %v\n", target2.MyCreatedAt)

	fmt.Printf("   来自 D 的字段（自动扁平化）:\n")
	fmt.Printf("     Address: %v\n", target2.Address)
	fmt.Printf("     MyAddress: %v\n", target2.MyAddress)
	fmt.Printf("     Phone: %v\n", target2.Phone)
	fmt.Printf("     MyPhone: %v\n", target2.MyPhone)
	fmt.Printf("     Salary: %v\n", target2.Salary)
	fmt.Printf("     MySalary: %v\n", target2.MySalary)
	fmt.Printf("     Rating: %v\n", target2.Rating)
	fmt.Printf("     MyRating: %v\n", target2.MyRating)
	fmt.Printf("     IsVIP: %v\n", target2.IsVIP)
	fmt.Printf("     MyIsVIP: %v\n", target2.MyIsVIP)
	fmt.Printf("     UpdatedAt: %v\n", target2.UpdatedAt)
	fmt.Printf("     MyUpdatedAt: %v\n", target2.MyUpdatedAt)
	fmt.Println()

	// 示例3: 使用全局配置
	fmt.Println("3. 使用全局配置:")

	// 获取当前全局配置
	currentConfig := mapster.GetGlobalConfig()
	fmt.Printf("   当前全局配置:\n")
	fmt.Printf("     时间转换: %v\n", currentConfig.EnableTimeConversion)
	fmt.Printf("     自动扁平化: %v\n", currentConfig.AutoFlatten.Enabled)
	fmt.Printf("     最大深度: %v\n", currentConfig.AutoFlatten.MaxDepth)
	fmt.Printf("     冲突策略: %v\n", currentConfig.AutoFlatten.ConflictStrategy)

	// 设置新的全局配置
	newConfig := mapster.GlobalConfig{
		EnableTimeConversion: true,
		AutoFlatten: mapster.AutoFlattenConfig{
			Enabled:          true,
			MaxDepth:         2,
			ConflictStrategy: mapster.KeepLast,
			UsePrefix:        false,
			PrefixSeparator:  "_",
		},
	}
	mapster.SetGlobalConfig(newConfig)

	var target3 StructC
	mapster.MapTo(src, &target3)

	fmt.Printf("   配置后的映射结果:\n")
	fmt.Printf("     Name: %v\n", target3.Name)
	fmt.Printf("     Age: %v\n", target3.Age)
	fmt.Printf("     Address: %v\n", target3.Address)
	fmt.Printf("     Salary: %v\n", target3.Salary)
	fmt.Println()

	// 示例4: 验证转换的正确性
	fmt.Println("4. 验证转换的正确性:")

	// 验证来自 A 的字段
	if target2.MyID == 1001 {
		fmt.Println("   ✓ MyID -> int 转换正确")
	} else {
		fmt.Println("   ✗ MyID -> int 转换错误")
	}

	if target2.MyTitle == "Senior Developer" {
		fmt.Println("   ✓ MyTitle -> string 转换正确")
	} else {
		fmt.Println("   ✗ MyTitle -> string 转换错误")
	}

	if target2.MyStatus == true {
		fmt.Println("   ✓ MyStatus -> bool 转换正确")
	} else {
		fmt.Println("   ✗ MyStatus -> bool 转换错误")
	}

	// 验证来自 B 的字段
	if target2.Name == "John Doe" {
		fmt.Println("   ✓ DataB.Name -> Name 自动扁平化正确")
	} else {
		fmt.Println("   ✗ DataB.Name -> Name 自动扁平化错误")
	}

	if target2.MyAge == 30 {
		fmt.Println("   ✓ DataB.MyAge -> MyAge 自动扁平化和转换正确")
	} else {
		fmt.Println("   ✗ DataB.MyAge -> MyAge 自动扁平化和转换错误")
	}

	if target2.MyEmail == "john@example.com" {
		fmt.Println("   ✓ DataB.MyEmail -> MyEmail 自动扁平化和转换正确")
	} else {
		fmt.Println("   ✗ DataB.MyEmail -> MyEmail 自动扁平化和转换错误")
	}

	// 验证来自 D 的字段
	if target2.Address == "123 Main St" {
		fmt.Println("   ✓ DataD.Address -> Address 自动扁平化正确")
	} else {
		fmt.Println("   ✗ DataD.Address -> Address 自动扁平化错误")
	}

	if target2.MySalary == 75000 {
		fmt.Println("   ✓ DataD.MySalary -> MySalary 自动扁平化和转换正确")
	} else {
		fmt.Println("   ✗ DataD.MySalary -> MySalary 自动扁平化和转换错误")
	}

	if target2.MyRating == 4.8 {
		fmt.Println("   ✓ DataD.MyRating -> MyRating 自动扁平化和转换正确")
	} else {
		fmt.Println("   ✗ DataD.MyRating -> MyRating 自动扁平化和转换错误")
	}

	fmt.Println("\n=== 自动扁平化功能测试完成 ===")
}
